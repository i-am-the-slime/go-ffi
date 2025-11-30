package purescript_httpurple

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestOkResponse(t *testing.T) {
	exports := Foreign("HTTPurple")
	ok := exports["ok"].(func(Any) Any)
	
	resp := ok("Hello, World!").(Dict)
	
	if resp["status"] != 200 {
		t.Errorf("Expected status 200, got %v", resp["status"])
	}
	
	if resp["body"] != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got %v", resp["body"])
	}
}

func TestNotFound(t *testing.T) {
	exports := Foreign("HTTPurple")
	notFound := exports["notFound"].(Dict)
	
	if notFound["status"] != 404 {
		t.Errorf("Expected status 404, got %v", notFound["status"])
	}
}

func TestJsonResponse(t *testing.T) {
	exports := Foreign("HTTPurple")
	json := exports["json"].(func(Any) Any)
	
	resp := json(`{"message":"test"}`).(Dict)
	
	if resp["status"] != 200 {
		t.Errorf("Expected status 200, got %v", resp["status"])
	}
	
	headers := resp["headers"].(Dict)
	contentType := headers["Content-Type"].(string)
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}
}

func TestHtmlResponse(t *testing.T) {
	exports := Foreign("HTTPurple")
	html := exports["html"].(func(Any) Any)
	
	resp := html("<h1>Hello</h1>").(Dict)
	
	headers := resp["headers"].(Dict)
	contentType := headers["Content-Type"].(string)
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}
}

func TestRedirect(t *testing.T) {
	exports := Foreign("HTTPurple")
	redirect := exports["redirect"].(func(Any) Any)
	
	resp := redirect("/new-location").(Dict)
	
	if resp["status"] != 302 {
		t.Errorf("Expected status 302, got %v", resp["status"])
	}
	
	headers := resp["headers"].(Dict)
	location := headers["Location"].(string)
	if location != "/new-location" {
		t.Errorf("Expected location '/new-location', got %s", location)
	}
}

func TestWithStatus(t *testing.T) {
	exports := Foreign("HTTPurple")
	ok := exports["ok"].(func(Any) Any)
	withStatus := exports["withStatus"].(func(Any, Any) Any)
	
	resp := ok("test")
	modified := withStatus(201, resp).(Dict)
	
	if modified["status"] != 201 {
		t.Errorf("Expected status 201, got %v", modified["status"])
	}
}

func TestWithHeader(t *testing.T) {
	exports := Foreign("HTTPurple")
	ok := exports["ok"].(func(Any) Any)
	withHeader := exports["withHeader"].(func(Any, Any, Any) Any)
	
	resp := ok("test")
	modified := withHeader("X-Custom", "value", resp).(Dict)
	
	headers := modified["headers"].(Dict)
	if headers["X-Custom"] != "value" {
		t.Errorf("Expected X-Custom header, got %v", headers)
	}
}

func TestWithBody(t *testing.T) {
	exports := Foreign("HTTPurple")
	ok := exports["ok"].(func(Any) Any)
	withBody := exports["withBody"].(func(Any, Any) Any)
	
	resp := ok("original")
	modified := withBody("modified", resp).(Dict)
	
	if modified["body"] != "modified" {
		t.Errorf("Expected body 'modified', got %v", modified["body"])
	}
}

func TestRequestMethod(t *testing.T) {
	exports := Foreign("HTTPurple")
	method := exports["method"].(func(Any) Any)
	
	req := httptest.NewRequest("POST", "/test", nil)
	wrapped := wrapRequest(req)
	
	result := method(wrapped).(string)
	if result != "POST" {
		t.Errorf("Expected method POST, got %s", result)
	}
}

func TestRequestPath(t *testing.T) {
	exports := Foreign("HTTPurple")
	path := exports["path"].(func(Any) Any)
	
	req := httptest.NewRequest("GET", "/api/users/123", nil)
	wrapped := wrapRequest(req)
	
	result := path(wrapped).(string)
	if result != "/api/users/123" {
		t.Errorf("Expected path '/api/users/123', got %s", result)
	}
}

func TestRequestQuery(t *testing.T) {
	exports := Foreign("HTTPurple")
	query := exports["query"].(func(Any) Any)
	
	req := httptest.NewRequest("GET", "/test?name=value&foo=bar", nil)
	wrapped := wrapRequest(req)
	
	result := query(wrapped).(Dict)
	if result["name"] != "value" || result["foo"] != "bar" {
		t.Errorf("Expected query params, got %v", result)
	}
}

func TestRequestHeaders(t *testing.T) {
	exports := Foreign("HTTPurple")
	headers := exports["headers"].(func(Any) Any)
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Custom", "test-value")
	wrapped := wrapRequest(req)
	
	result := headers(wrapped).(Dict)
	if result["x-custom"] != "test-value" {
		t.Errorf("Expected header x-custom, got %v", result)
	}
}

func TestRequestHeader(t *testing.T) {
	exports := Foreign("HTTPurple")
	header := exports["header"].(func(Any, Any) Any)
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")
	wrapped := wrapRequest(req)
	
	// Test existing header
	result := header("authorization", wrapped).(Dict)
	if val, ok := result["value0"].(string); !ok || val != "Bearer token123" {
		t.Errorf("Expected Just 'Bearer token123', got %v", result)
	}
	
	// Test missing header
	result2 := header("nonexistent", wrapped).(Dict)
	if len(result2) != 0 {
		t.Error("Expected Nothing for missing header")
	}
}

func TestRequestBody(t *testing.T) {
	exports := Foreign("HTTPurple")
	body := exports["body"].(func(Any) Any)
	
	bodyContent := "request body content"
	req := httptest.NewRequest("POST", "/test", strings.NewReader(bodyContent))
	wrapped := wrapRequest(req)
	
	bodyEffect := body(wrapped).(func() Any)
	result := bodyEffect().(string)
	
	if result != bodyContent {
		t.Errorf("Expected body '%s', got %s", bodyContent, result)
	}
}

func TestApplyResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	
	response := Dict{
		"status": 201,
		"headers": Dict{
			"Content-Type":  "application/json",
			"X-Custom-Header": "custom-value",
		},
		"body": `{"success":true}`,
	}
	
	applyResponse(recorder, response)
	
	if recorder.Code != 201 {
		t.Errorf("Expected status 201, got %d", recorder.Code)
	}
	
	if ct := recorder.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}
	
	if custom := recorder.Header().Get("X-Custom-Header"); custom != "custom-value" {
		t.Errorf("Expected X-Custom-Header custom-value, got %s", custom)
	}
	
	body, _ := io.ReadAll(recorder.Body)
	if string(body) != `{"success":true}` {
		t.Errorf("Expected body '{\"success\":true}', got %s", string(body))
	}
}

func TestStatusCodes(t *testing.T) {
	exports := Foreign("HTTPurple")
	
	tests := []struct {
		name     string
		funcName string
		expected int
	}{
		{"created", "created", 201},
		{"accepted", "accepted", 202},
		{"noContent", "noContent", 204},
		{"badRequest", "badRequest", 400},
		{"unauthorized", "unauthorized", 401},
		{"forbidden", "forbidden", 403},
		{"notFound", "notFound", 404},
		{"internalServerError", "internalServerError", 500},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp Dict
			
			if fn, ok := exports[tt.funcName].(func(Any) Any); ok {
				resp = fn("test").(Dict)
			} else {
				resp = exports[tt.funcName].(Dict)
			}
			
			if resp["status"] != tt.expected {
				t.Errorf("%s: expected status %d, got %v", tt.name, tt.expected, resp["status"])
			}
		})
	}
}

