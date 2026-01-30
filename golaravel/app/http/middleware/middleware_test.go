package middleware

import (
        "net/http"
        "net/http/httptest"
        "testing"
        "time"

        "golaravel/app/http/request"
        "golaravel/app/http/response"
        "golaravel/app/routing"
)

func createTestHandler() routing.HandlerFunc {
        return func(req *request.Request, res *response.Response) error {
                return res.String("OK")
        }
}

func TestLoggerMiddleware(t *testing.T) {
        handler := Logger()(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        err := handler(req, res)
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200, got: %d", w.Code)
        }
}

func TestRecoveryMiddleware(t *testing.T) {
        panicHandler := func(req *request.Request, res *response.Response) error {
                panic("test panic")
        }

        handler := Recovery()(panicHandler)

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        err := handler(req, res)
        if err != nil {
                t.Fatalf("Expected recovery to handle panic, got error: %v", err)
        }

        if w.Code != http.StatusInternalServerError {
                t.Errorf("Expected status 500 after panic, got: %d", w.Code)
        }
}

func TestCORSMiddleware(t *testing.T) {
        handler := CORS("*")(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.Header.Set("Origin", "http://example.com")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Header().Get("Access-Control-Allow-Origin") != "*" {
                t.Error("Expected CORS header to be set")
        }
}

func TestCORSMiddlewareSpecificOrigin(t *testing.T) {
        handler := CORS("http://allowed.com")(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.Header.Set("Origin", "http://allowed.com")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Header().Get("Access-Control-Allow-Origin") != "http://allowed.com" {
                t.Error("Expected CORS header for allowed origin")
        }
}

func TestCORSMiddlewareOptions(t *testing.T) {
        handler := CORS("*")(createTestHandler())

        httpReq := httptest.NewRequest("OPTIONS", "/test", nil)
        httpReq.Header.Set("Origin", "http://example.com")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusNoContent {
                t.Errorf("Expected status 204 for OPTIONS, got: %d", w.Code)
        }
}

func TestJSONMiddleware(t *testing.T) {
        jsonHandler := func(req *request.Request, res *response.Response) error {
                return res.JSON(map[string]string{"status": "ok"})
        }
        handler := JSON()(jsonHandler)

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Header().Get("Content-Type") != "application/json" {
                t.Errorf("Expected Content-Type to be application/json, got: %s", w.Header().Get("Content-Type"))
        }
}

func TestAuthMiddleware(t *testing.T) {
        validateToken := func(token string) bool {
                return token == "valid-token"
        }

        handler := Auth(validateToken)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.Header.Set("Authorization", "Bearer valid-token")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200 for valid token, got: %d", w.Code)
        }
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
        validateToken := func(token string) bool {
                return token == "valid-token"
        }

        handler := Auth(validateToken)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.Header.Set("Authorization", "Bearer invalid-token")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusUnauthorized {
                t.Errorf("Expected status 401 for invalid token, got: %d", w.Code)
        }
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
        validateToken := func(token string) bool {
                return true
        }

        handler := Auth(validateToken)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusUnauthorized {
                t.Errorf("Expected status 401 for missing token, got: %d", w.Code)
        }
}

func TestBasicAuthMiddleware(t *testing.T) {
        users := map[string]string{
                "admin": "secret",
        }

        handler := BasicAuth(users)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.SetBasicAuth("admin", "secret")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200 for valid credentials, got: %d", w.Code)
        }
}

func TestBasicAuthMiddlewareInvalid(t *testing.T) {
        users := map[string]string{
                "admin": "secret",
        }

        handler := BasicAuth(users)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.SetBasicAuth("admin", "wrong")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusUnauthorized {
                t.Errorf("Expected status 401 for invalid credentials, got: %d", w.Code)
        }
}

func TestBasicAuthMiddlewareMissing(t *testing.T) {
        users := map[string]string{
                "admin": "secret",
        }

        handler := BasicAuth(users)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusUnauthorized {
                t.Errorf("Expected status 401 for missing auth, got: %d", w.Code)
        }
        if w.Header().Get("WWW-Authenticate") == "" {
                t.Error("Expected WWW-Authenticate header")
        }
}

func TestRateLimiterMiddleware(t *testing.T) {
        handler := RateLimiter(2, time.Minute)(createTestHandler())

        for i := 0; i < 2; i++ {
                httpReq := httptest.NewRequest("GET", "/test", nil)
                httpReq.RemoteAddr = "192.168.1.1:12345"
                w := httptest.NewRecorder()

                req := request.New(httpReq)
                res := response.New(w)

                handler(req, res)

                if w.Code != http.StatusOK {
                        t.Errorf("Request %d: Expected status 200, got: %d", i+1, w.Code)
                }
        }

        httpReq := httptest.NewRequest("GET", "/test", nil)
        httpReq.RemoteAddr = "192.168.1.1:12345"
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != 429 {
                t.Errorf("Expected status 429 after rate limit exceeded, got: %d", w.Code)
        }
}

func TestSecureHeadersMiddleware(t *testing.T) {
        handler := SecureHeaders()(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        expectedHeaders := []string{
                "X-Content-Type-Options",
                "X-Frame-Options",
                "X-XSS-Protection",
                "Referrer-Policy",
        }

        for _, header := range expectedHeaders {
                if w.Header().Get(header) == "" {
                        t.Errorf("Expected %s header to be set", header)
                }
        }
}

func TestTimeoutMiddleware(t *testing.T) {
        slowHandler := func(req *request.Request, res *response.Response) error {
                time.Sleep(200 * time.Millisecond)
                return res.String("OK")
        }

        handler := Timeout(100 * time.Millisecond)(slowHandler)

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != 504 {
                t.Errorf("Expected status 504 for timeout, got: %d", w.Code)
        }
}

func TestTimeoutMiddlewareFastHandler(t *testing.T) {
        handler := Timeout(time.Second)(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200 for fast handler, got: %d", w.Code)
        }
}

func TestContentTypeMiddleware(t *testing.T) {
        handler := ContentType("application/json")(createTestHandler())

        httpReq := httptest.NewRequest("POST", "/test", nil)
        httpReq.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200 for correct content type, got: %d", w.Code)
        }
}

func TestContentTypeMiddlewareWrong(t *testing.T) {
        handler := ContentType("application/json")(createTestHandler())

        httpReq := httptest.NewRequest("POST", "/test", nil)
        httpReq.Header.Set("Content-Type", "text/plain")
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != 415 {
                t.Errorf("Expected status 415 for wrong content type, got: %d", w.Code)
        }
}

func TestContentTypeMiddlewareGetRequest(t *testing.T) {
        handler := ContentType("application/json")(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        handler(req, res)

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200 for GET request, got: %d", w.Code)
        }
}

func TestTrimStringsMiddleware(t *testing.T) {
        handler := TrimStrings()(createTestHandler())

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        err := handler(req, res)
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }

        if w.Code != http.StatusOK {
                t.Errorf("Expected status 200, got: %d", w.Code)
        }
}

func TestMiddlewareChaining(t *testing.T) {
        callOrder := []string{}

        middleware1 := func(next routing.HandlerFunc) routing.HandlerFunc {
                return func(req *request.Request, res *response.Response) error {
                        callOrder = append(callOrder, "m1-before")
                        err := next(req, res)
                        callOrder = append(callOrder, "m1-after")
                        return err
                }
        }

        middleware2 := func(next routing.HandlerFunc) routing.HandlerFunc {
                return func(req *request.Request, res *response.Response) error {
                        callOrder = append(callOrder, "m2-before")
                        err := next(req, res)
                        callOrder = append(callOrder, "m2-after")
                        return err
                }
        }

        handler := func(req *request.Request, res *response.Response) error {
                callOrder = append(callOrder, "handler")
                return res.String("OK")
        }

        chained := middleware1(middleware2(handler))

        httpReq := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()

        req := request.New(httpReq)
        res := response.New(w)

        chained(req, res)

        expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
        if len(callOrder) != len(expected) {
                t.Errorf("Expected %d calls, got: %d", len(expected), len(callOrder))
        }

        for i, v := range expected {
                if callOrder[i] != v {
                        t.Errorf("Expected call %d to be '%s', got: '%s'", i, v, callOrder[i])
                }
        }
}
