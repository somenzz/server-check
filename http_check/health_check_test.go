package http_check

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHealthWithOrigin(t *testing.T) {
	// Start a local HTTP server to receive the request and verify the Origin header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "https://example.com" {
			t.Errorf("Expected Origin header 'https://example.com', got '%s'", origin)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	ctx := context.Background()
	// Test health check passing the origin
	ok := CheckHealth(ctx, ts.URL, "GET", 200, "ok", "https://example.com")
	if !ok {
		t.Errorf("Expected CheckHealth to return true, got false")
	}

	// Test health check passing the empty origin
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			t.Errorf("Expected no Origin header, got '%s'", origin)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts2.Close()

	ok2 := CheckHealth(ctx, ts2.URL, "GET", 200, "ok", "")
	if !ok2 {
		t.Errorf("Expected CheckHealth to return true with empty origin, got false")
	}
}
