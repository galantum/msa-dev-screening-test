package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyRequestSuccess(t *testing.T) {
	// Mock User Service
	mockUserService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "username": "manager", "email": "manager@example.com", "age": "40"}`))
	}))
	defer mockUserService.Close()

	// Test Gateway
	req := httptest.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()
	proxyRequest(w, req, mockUserService.URL)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != `{"id": "1", "username": "manager", "email": "manager@example.com", "age": "40"}` {
		t.Fatalf("Unexpected response: %s", body)
	}
}

func TestProxyRequestFailure(t *testing.T) {
	// Test Gateway with non-existent service
	req := httptest.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()
	proxyRequest(w, req, "http://localhost:9999/user/profile")

	resp := w.Result()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("Expected status 503, got %d", resp.StatusCode)
	}
}
