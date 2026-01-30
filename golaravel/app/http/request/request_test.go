package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequest(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := New(httpReq)

	if req == nil {
		t.Fatal("Expected request to be created")
	}
	if req.Method() != "GET" {
		t.Errorf("Expected method to be GET, got: %s", req.Method())
	}
	if req.Path() != "/test" {
		t.Errorf("Expected path to be /test, got: %s", req.Path())
	}
}

func TestRequestRaw(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := New(httpReq)

	if req.Raw() != httpReq {
		t.Error("Expected Raw() to return original request")
	}
}

func TestRequestParams(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := New(httpReq)

	req.SetParam("id", "123")
	req.SetParam("name", "john")

	if req.Param("id") != "123" {
		t.Errorf("Expected param id to be '123', got: %s", req.Param("id"))
	}
	if req.Param("name") != "john" {
		t.Errorf("Expected param name to be 'john', got: %s", req.Param("name"))
	}
	if req.Param("nonexistent") != "" {
		t.Error("Expected nonexistent param to be empty")
	}

	params := req.Params()
	if len(params) != 2 {
		t.Errorf("Expected 2 params, got: %d", len(params))
	}
}

func TestRequestQuery(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?name=john&age=30", nil)
	req := New(httpReq)

	if req.Query("name") != "john" {
		t.Errorf("Expected query name to be 'john', got: %s", req.Query("name"))
	}
	if req.Query("age") != "30" {
		t.Errorf("Expected query age to be '30', got: %s", req.Query("age"))
	}
	if req.Query("nonexistent") != "" {
		t.Error("Expected nonexistent query to be empty")
	}
}

func TestRequestQueryDefault(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?name=john", nil)
	req := New(httpReq)

	if req.QueryDefault("name", "default") != "john" {
		t.Error("Expected existing query to return actual value")
	}
	if req.QueryDefault("nonexistent", "default") != "default" {
		t.Error("Expected nonexistent query to return default value")
	}
}

func TestRequestQueryAll(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?a=1&b=2", nil)
	req := New(httpReq)

	queryAll := req.QueryAll()
	if queryAll.Get("a") != "1" {
		t.Error("Expected QueryAll to contain 'a'")
	}
	if queryAll.Get("b") != "2" {
		t.Error("Expected QueryAll to contain 'b'")
	}
}

func TestRequestHeaders(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("X-Custom-Header", "custom-value")
	httpReq.Header.Set("Content-Type", "application/json")
	req := New(httpReq)

	if req.Header("X-Custom-Header") != "custom-value" {
		t.Errorf("Expected custom header value, got: %s", req.Header("X-Custom-Header"))
	}
	if req.Header("Content-Type") != "application/json" {
		t.Errorf("Expected content-type, got: %s", req.Header("Content-Type"))
	}
}

func TestRequestBody(t *testing.T) {
	body := []byte(`{"name":"john"}`)
	httpReq := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req := New(httpReq)

	bodyContent, err := req.Body()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if string(bodyContent) != `{"name":"john"}` {
		t.Errorf("Expected body content, got: %s", string(bodyContent))
	}

	bodyContent2, _ := req.Body()
	if string(bodyContent2) != `{"name":"john"}` {
		t.Error("Expected body to be cached and return same content on second call")
	}
}

func TestRequestJSON(t *testing.T) {
	body := []byte(`{"name":"john","age":30}`)
	httpReq := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req := New(httpReq)

	var data map[string]interface{}
	err := req.JSON(&data)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if data["name"] != "john" {
		t.Errorf("Expected name to be 'john', got: %v", data["name"])
	}
	if data["age"].(float64) != 30 {
		t.Errorf("Expected age to be 30, got: %v", data["age"])
	}
}

func TestRequestInput(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?query=value", nil)
	req := New(httpReq)
	req.SetParam("param", "param-value")

	if req.Input("param") != "param-value" {
		t.Error("Expected Input to return param value first")
	}
	if req.Input("query") != "value" {
		t.Error("Expected Input to return query value")
	}
	if req.Input("nonexistent") != "" {
		t.Error("Expected Input to return empty for nonexistent")
	}
}

func TestRequestHas(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?name=john", nil)
	req := New(httpReq)

	if !req.Has("name") {
		t.Error("Expected Has to return true for existing query")
	}
	if req.Has("nonexistent") {
		t.Error("Expected Has to return false for nonexistent query")
	}
}

func TestRequestFilled(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?name=john&empty=", nil)
	req := New(httpReq)

	if !req.Filled("name") {
		t.Error("Expected Filled to return true for non-empty value")
	}
	if req.Filled("empty") {
		t.Error("Expected Filled to return false for empty value")
	}
	if req.Filled("nonexistent") {
		t.Error("Expected Filled to return false for nonexistent")
	}
}

func TestRequestIsMethod(t *testing.T) {
	httpReq := httptest.NewRequest("POST", "/test", nil)
	req := New(httpReq)

	if !req.IsMethod("POST") {
		t.Error("Expected IsMethod to return true for POST")
	}
	if !req.IsMethod("post") {
		t.Error("Expected IsMethod to be case-insensitive")
	}
	if req.IsMethod("GET") {
		t.Error("Expected IsMethod to return false for GET")
	}
}

func TestRequestIsAjax(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := New(httpReq)

	if req.IsAjax() {
		t.Error("Expected IsAjax to return false without header")
	}

	httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	req2 := New(httpReq)

	if !req2.IsAjax() {
		t.Error("Expected IsAjax to return true with header")
	}
}

func TestRequestWantsJSON(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("Accept", "application/json")
	req := New(httpReq)

	if !req.WantsJSON() {
		t.Error("Expected WantsJSON to return true")
	}

	httpReq2 := httptest.NewRequest("GET", "/test", nil)
	httpReq2.Header.Set("Accept", "text/html")
	req2 := New(httpReq2)

	if req2.WantsJSON() {
		t.Error("Expected WantsJSON to return false")
	}
}

func TestRequestIP(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.RemoteAddr = "192.168.1.1:12345"
	req := New(httpReq)

	if req.IP() != "192.168.1.1" {
		t.Errorf("Expected IP to be '192.168.1.1', got: %s", req.IP())
	}
}

func TestRequestIPFromForwardedFor(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	req := New(httpReq)

	if req.IP() != "10.0.0.1" {
		t.Errorf("Expected IP to be '10.0.0.1', got: %s", req.IP())
	}
}

func TestRequestIPFromRealIP(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("X-Real-IP", "172.16.0.1")
	req := New(httpReq)

	if req.IP() != "172.16.0.1" {
		t.Errorf("Expected IP to be '172.16.0.1', got: %s", req.IP())
	}
}

func TestRequestUserAgent(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("User-Agent", "TestAgent/1.0")
	req := New(httpReq)

	if req.UserAgent() != "TestAgent/1.0" {
		t.Errorf("Expected UserAgent to be 'TestAgent/1.0', got: %s", req.UserAgent())
	}
}

func TestRequestAccepts(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("Accept", "application/json, text/html")
	req := New(httpReq)

	if !req.Accepts("application/json") {
		t.Error("Expected Accepts to return true for application/json")
	}
	if !req.Accepts("text/html") {
		t.Error("Expected Accepts to return true for text/html")
	}
	if req.Accepts("text/plain") {
		t.Error("Expected Accepts to return false for text/plain")
	}
}

func TestRequestBearerToken(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.Header.Set("Authorization", "Bearer abc123token")
	req := New(httpReq)

	if req.BearerToken() != "abc123token" {
		t.Errorf("Expected BearerToken to be 'abc123token', got: %s", req.BearerToken())
	}
}

func TestRequestBearerTokenEmpty(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := New(httpReq)

	if req.BearerToken() != "" {
		t.Error("Expected BearerToken to be empty without Authorization header")
	}

	httpReq2 := httptest.NewRequest("GET", "/test", nil)
	httpReq2.Header.Set("Authorization", "Basic abc123")
	req2 := New(httpReq2)

	if req2.BearerToken() != "" {
		t.Error("Expected BearerToken to be empty for non-Bearer auth")
	}
}

func TestRequestAll(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?query1=value1&query2=value2", nil)
	req := New(httpReq)
	req.SetParam("param1", "paramvalue1")

	all := req.All()

	if all["param1"] != "paramvalue1" {
		t.Error("Expected All to contain param1")
	}
	if all["query1"] != "value1" {
		t.Error("Expected All to contain query1")
	}
	if all["query2"] != "value2" {
		t.Error("Expected All to contain query2")
	}
}

func TestRequestURL(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test?foo=bar", nil)
	req := New(httpReq)

	url := req.URL()
	if url.Path != "/test" {
		t.Errorf("Expected URL path to be '/test', got: %s", url.Path)
	}
	if url.RawQuery != "foo=bar" {
		t.Errorf("Expected URL query to be 'foo=bar', got: %s", url.RawQuery)
	}
}

func TestRequestCookie(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	req := New(httpReq)

	cookie, err := req.Cookie("session")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cookie.Value != "abc123" {
		t.Errorf("Expected cookie value to be 'abc123', got: %s", cookie.Value)
	}
}

func TestRequestCookies(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/test", nil)
	httpReq.AddCookie(&http.Cookie{Name: "session", Value: "abc"})
	httpReq.AddCookie(&http.Cookie{Name: "token", Value: "xyz"})
	req := New(httpReq)

	cookies := req.Cookies()
	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got: %d", len(cookies))
	}
}
