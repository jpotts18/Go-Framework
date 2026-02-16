package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	writer     http.ResponseWriter
	request    *http.Request
	statusCode int
	headers    map[string]string
	cookies    []*http.Cookie
	sent       bool
}

func New(w http.ResponseWriter, r ...*http.Request) *Response {
	res := &Response{
		writer:     w,
		statusCode: http.StatusOK,
		headers:    make(map[string]string),
		cookies:    make([]*http.Cookie, 0),
	}
	if len(r) > 0 {
		res.request = r[0]
	}
	return res
}

func (r *Response) Raw() http.ResponseWriter {
	return r.writer
}

func (r *Response) Status(code int) *Response {
	r.statusCode = code
	return r
}

func (r *Response) Header(key, value string) *Response {
	r.headers[key] = value
	return r
}

func (r *Response) Cookie(cookie *http.Cookie) *Response {
	r.cookies = append(r.cookies, cookie)
	return r
}

func (r *Response) applyHeaders() {
	for key, value := range r.headers {
		r.writer.Header().Set(key, value)
	}
	for _, cookie := range r.cookies {
		http.SetCookie(r.writer, cookie)
	}
}

func (r *Response) Send(content []byte) error {
	if r.sent {
		return fmt.Errorf("response already sent")
	}

	r.applyHeaders()
	r.writer.WriteHeader(r.statusCode)
	_, err := r.writer.Write(content)
	r.sent = true
	return err
}

func (r *Response) String(content string) error {
	r.Header("Content-Type", "text/plain; charset=utf-8")
	return r.Send([]byte(content))
}

func (r *Response) HTML(content string) error {
	r.Header("Content-Type", "text/html; charset=utf-8")
	r.Header("Cache-Control", "no-cache")
	return r.Send([]byte(content))
}

func (r *Response) JSON(data interface{}) error {
	r.Header("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.Send(jsonBytes)
}

func (r *Response) JSONPretty(data interface{}) error {
	r.Header("Content-Type", "application/json")

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return r.Send(jsonBytes)
}

func (r *Response) Redirect(url string, code int) error {
	if r.sent {
		return fmt.Errorf("response already sent")
	}

	r.applyHeaders()
	if r.request != nil {
		http.Redirect(r.writer, r.request, url, code)
	} else {
		r.writer.Header().Set("Location", url)
		r.writer.WriteHeader(code)
	}
	r.sent = true
	return nil
}

func (r *Response) RedirectPermanent(url string) error {
	return r.Redirect(url, http.StatusMovedPermanently)
}

func (r *Response) RedirectTemporary(url string) error {
	return r.Redirect(url, http.StatusFound)
}

func (r *Response) NoContent() error {
	r.statusCode = http.StatusNoContent
	r.applyHeaders()
	r.writer.WriteHeader(r.statusCode)
	r.sent = true
	return nil
}

func (r *Response) NotFound(message string) error {
	if message == "" {
		message = "Not Found"
	}
	return r.Status(http.StatusNotFound).JSON(map[string]string{
		"error": message,
	})
}

func (r *Response) BadRequest(message string) error {
	if message == "" {
		message = "Bad Request"
	}
	return r.Status(http.StatusBadRequest).JSON(map[string]string{
		"error": message,
	})
}

func (r *Response) Unauthorized(message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return r.Status(http.StatusUnauthorized).JSON(map[string]string{
		"error": message,
	})
}

func (r *Response) Forbidden(message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return r.Status(http.StatusForbidden).JSON(map[string]string{
		"error": message,
	})
}

func (r *Response) ServerError(message string) error {
	if message == "" {
		message = "Internal Server Error"
	}
	return r.Status(http.StatusInternalServerError).JSON(map[string]string{
		"error": message,
	})
}

func (r *Response) ValidationError(errors map[string][]string) error {
	return r.Status(http.StatusUnprocessableEntity).JSON(map[string]interface{}{
		"message": "Validation failed",
		"errors":  errors,
	})
}

func (r *Response) Download(content []byte, filename string) error {
	r.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	r.Header("Content-Type", "application/octet-stream")
	return r.Send(content)
}

func (r *Response) IsSent() bool {
	return r.sent
}
