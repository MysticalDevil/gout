package gout

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddlewareGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := newContext(rec, req)
	c.handlers = []HandlerFunc{
		CORS(),
		func(c *Context) {
			c.Status(http.StatusOK)
		},
	}

	c.Next()

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected Access-Control-Allow-Origin '*', got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("expected Access-Control-Allow-Credentials to be empty by default, got %q", got)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestCORSMiddlewareOptionsRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()
	c := newContext(rec, req)
	c.handlers = []HandlerFunc{
		CORS(),
		func(c *Context) {
			c.Status(http.StatusOK)
		},
	}

	c.Next()

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 for OPTIONS preflight, got %d", rec.Code)
	}
}
