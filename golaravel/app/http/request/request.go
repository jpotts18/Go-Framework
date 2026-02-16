package request

import (
        "context"
        "encoding/json"
        "io"
        "mime/multipart"
        "net/http"
        "net/url"
        "strings"
)

type Request struct {
        raw        *http.Request
        params     map[string]string
        query      url.Values
        body       []byte
        bodyParsed bool
        ctx        context.Context
}

func New(r *http.Request) *Request {
        return &Request{
                raw:    r,
                params: make(map[string]string),
                query:  r.URL.Query(),
                ctx:    r.Context(),
        }
}

func (r *Request) Context() context.Context {
        if r.ctx == nil {
                return r.raw.Context()
        }
        return r.ctx
}

func (r *Request) WithContext(ctx context.Context) *Request {
        r.ctx = ctx
        return r
}

func (r *Request) Raw() *http.Request {
        return r.raw
}

func (r *Request) Method() string {
        return r.raw.Method
}

func (r *Request) Path() string {
        return r.raw.URL.Path
}

func (r *Request) URL() *url.URL {
        return r.raw.URL
}

func (r *Request) SetParam(key, value string) {
        r.params[key] = value
}

func (r *Request) Param(key string) string {
        return r.params[key]
}

func (r *Request) Params() map[string]string {
        return r.params
}

func (r *Request) Query(key string) string {
        return r.query.Get(key)
}

func (r *Request) QueryDefault(key, defaultVal string) string {
        val := r.query.Get(key)
        if val == "" {
                return defaultVal
        }
        return val
}

func (r *Request) QueryAll() url.Values {
        return r.query
}

func (r *Request) Header(key string) string {
        return r.raw.Header.Get(key)
}

func (r *Request) Headers() http.Header {
        return r.raw.Header
}

func (r *Request) Cookie(name string) (*http.Cookie, error) {
        return r.raw.Cookie(name)
}

func (r *Request) Cookies() []*http.Cookie {
        return r.raw.Cookies()
}

func (r *Request) Body() ([]byte, error) {
        if r.bodyParsed {
                return r.body, nil
        }

        body, err := io.ReadAll(r.raw.Body)
        if err != nil {
                return nil, err
        }

        r.body = body
        r.bodyParsed = true
        return r.body, nil
}

func (r *Request) JSON(v interface{}) error {
        body, err := r.Body()
        if err != nil {
                return err
        }
        return json.Unmarshal(body, v)
}

func (r *Request) FormValue(key string) string {
        return r.raw.FormValue(key)
}

func (r *Request) PostFormValue(key string) string {
        return r.raw.PostFormValue(key)
}

func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
        return r.raw.FormFile(key)
}

func (r *Request) Input(key string) string {
        if val := r.Param(key); val != "" {
                return val
        }
        if val := r.Query(key); val != "" {
                return val
        }
        if val := r.FormValue(key); val != "" {
                return val
        }
        return ""
}

func (r *Request) All() map[string]interface{} {
        result := make(map[string]interface{})

        for k, v := range r.params {
                result[k] = v
        }

        for k, v := range r.query {
                if len(v) == 1 {
                        result[k] = v[0]
                } else {
                        result[k] = v
                }
        }

        contentType := r.Header("Content-Type")
        if strings.Contains(contentType, "application/json") {
                var jsonData map[string]interface{}
                if err := r.JSON(&jsonData); err == nil {
                        for k, v := range jsonData {
                                result[k] = v
                        }
                }
        } else if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "multipart/form-data") {
                if err := r.raw.ParseForm(); err == nil {
                        for k, v := range r.raw.PostForm {
                                if len(v) == 1 {
                                        result[k] = v[0]
                                } else {
                                        result[k] = v
                                }
                        }
                }
        }

        return result
}

func (r *Request) Has(key string) bool {
        return r.Input(key) != ""
}

func (r *Request) Filled(key string) bool {
        val := r.Input(key)
        return val != "" && strings.TrimSpace(val) != ""
}

func (r *Request) IsMethod(method string) bool {
        return strings.EqualFold(r.Method(), method)
}

func (r *Request) IsAjax() bool {
        return r.Header("X-Requested-With") == "XMLHttpRequest"
}

func (r *Request) WantsJSON() bool {
        accept := r.Header("Accept")
        return strings.Contains(accept, "application/json")
}

func (r *Request) IP() string {
        xff := r.Header("X-Forwarded-For")
        if xff != "" {
                parts := strings.Split(xff, ",")
                return strings.TrimSpace(parts[0])
        }

        xri := r.Header("X-Real-IP")
        if xri != "" {
                return xri
        }

        addr := r.raw.RemoteAddr
        if idx := strings.LastIndex(addr, ":"); idx != -1 {
                return addr[:idx]
        }
        return addr
}

func (r *Request) UserAgent() string {
        return r.Header("User-Agent")
}

func (r *Request) Accepts(contentTypes ...string) bool {
        accept := r.Header("Accept")
        for _, ct := range contentTypes {
                if strings.Contains(accept, ct) {
                        return true
                }
        }
        return false
}

func (r *Request) BearerToken() string {
        auth := r.Header("Authorization")
        if strings.HasPrefix(auth, "Bearer ") {
                return strings.TrimPrefix(auth, "Bearer ")
        }
        return ""
}
