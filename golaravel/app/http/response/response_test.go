package response

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	if res == nil {
		t.Fatal("Expected response to be created")
	}
	if res.statusCode != http.StatusOK {
		t.Errorf("Expected default status to be 200, got: %d", res.statusCode)
	}
}

func TestResponseRaw(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	if res.Raw() != w {
		t.Error("Expected Raw() to return original ResponseWriter")
	}
}

func TestResponseStatus(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	res.Status(201)
	if res.statusCode != 201 {
		t.Errorf("Expected status to be 201, got: %d", res.statusCode)
	}
}

func TestResponseStatusChaining(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	result := res.Status(404)
	if result != res {
		t.Error("Expected Status to return same response for chaining")
	}
}

func TestResponseHeader(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	res.Header("X-Custom", "value").Send([]byte("test"))

	if w.Header().Get("X-Custom") != "value" {
		t.Error("Expected custom header to be set")
	}
}

func TestResponseSend(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.Send([]byte("Hello World"))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Body.String() != "Hello World" {
		t.Errorf("Expected body to be 'Hello World', got: %s", w.Body.String())
	}
}

func TestResponseSendOnlyOnce(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	res.Send([]byte("First"))
	err := res.Send([]byte("Second"))

	if err == nil {
		t.Error("Expected error when sending twice")
	}
}

func TestResponseString(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.String("Hello")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Body.String() != "Hello" {
		t.Errorf("Expected body to be 'Hello', got: %s", w.Body.String())
	}
	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("Expected Content-Type to be text/plain")
	}
}

func TestResponseHTML(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.HTML("<h1>Hello</h1>")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Body.String() != "<h1>Hello</h1>" {
		t.Errorf("Expected body to be '<h1>Hello</h1>', got: %s", w.Body.String())
	}
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("Expected Content-Type to be text/html")
	}
	if w.Header().Get("Cache-Control") != "no-cache" {
		t.Error("Expected Cache-Control to be no-cache")
	}
}

func TestResponseJSON(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	data := map[string]string{"message": "hello"}
	err := res.JSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := `{"message":"hello"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body to be %s, got: %s", expected, w.Body.String())
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type to be application/json")
	}
}

func TestResponseJSONPretty(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	data := map[string]string{"a": "b"}
	err := res.JSONPretty(data)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "{\n  \"a\": \"b\"\n}"
	if w.Body.String() != expected {
		t.Errorf("Expected pretty JSON, got: %s", w.Body.String())
	}
}

func TestResponseNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.NoContent()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got: %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Error("Expected empty body")
	}
}

func TestResponseNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.NotFound("Resource not found")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got: %d", w.Code)
	}
}

func TestResponseNotFoundDefaultMessage(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	res.NotFound("")

	if w.Body.String() != `{"error":"Not Found"}` {
		t.Errorf("Expected default message, got: %s", w.Body.String())
	}
}

func TestResponseBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.BadRequest("Invalid input")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got: %d", w.Code)
	}
}

func TestResponseUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.Unauthorized("Not authenticated")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got: %d", w.Code)
	}
}

func TestResponseForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.Forbidden("Access denied")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got: %d", w.Code)
	}
}

func TestResponseServerError(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	err := res.ServerError("Something went wrong")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got: %d", w.Code)
	}
}

func TestResponseValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	errors := map[string][]string{
		"email": {"Email is required", "Email must be valid"},
		"name":  {"Name is required"},
	}

	err := res.ValidationError(errors)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got: %d", w.Code)
	}
}

func TestResponseDownload(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	content := []byte("file content")
	err := res.Download(content, "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if w.Header().Get("Content-Disposition") != `attachment; filename="test.txt"` {
		t.Error("Expected Content-Disposition header")
	}
	if w.Header().Get("Content-Type") != "application/octet-stream" {
		t.Error("Expected Content-Type to be application/octet-stream")
	}
	if w.Body.String() != "file content" {
		t.Errorf("Expected body to be 'file content', got: %s", w.Body.String())
	}
}

func TestResponseIsSent(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	if res.IsSent() {
		t.Error("Expected IsSent to be false initially")
	}

	res.Send([]byte("test"))

	if !res.IsSent() {
		t.Error("Expected IsSent to be true after Send")
	}
}

func TestResponseCookie(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	cookie := &http.Cookie{Name: "session", Value: "abc123"}
	res.Cookie(cookie).Send([]byte("test"))

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got: %d", len(cookies))
	}
	if cookies[0].Name != "session" || cookies[0].Value != "abc123" {
		t.Error("Expected cookie to be set correctly")
	}
}

func TestResponseStatusWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	res := New(w)

	res.Status(201).JSON(map[string]string{"status": "created"})

	if w.Code != 201 {
		t.Errorf("Expected status 201, got: %d", w.Code)
	}
}
