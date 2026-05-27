package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
)

var testRouter *chi.Mux

func TestPing(t *testing.T) {
	if testRouter == nil {
		testRouter = CreateRoute()
	}

	testCombi := []map[string][]string{
		{"GET": []string{"/", "/p", "/p/v"}},
		{"POST": []string{"/", "/p", "/p/v"}},
		{"PUT": []string{"/", "/p", "/p/v"}},
		{"DELETE": []string{"/", "/p", "/p/v"}},
		{"OPTIONS": []string{"/", "/p", "/p/v"}},
		{"HEAD": []string{"/", "/p", "/p/v"}},
	}
	for _, tests := range testCombi {
		for method, paths := range tests {
			fmt.Println(method, paths)
			for _, p := range paths {
				req := httptest.NewRequest(method, p, nil)
				w := httptest.NewRecorder()
				testRouter.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			}
		}
	}

	failTestCombi := []map[string][]string{
		{"GET": []string{"/p/v/w"}},
		{"POST": []string{"/p/v/w"}},
		{"PUT": []string{"/p/v/w"}},
		{"DELETE": []string{"/p/v/w"}},
		{"OPTIONS": []string{"/p/v/w"}},
	}

	for _, tests := range failTestCombi {
		for method, paths := range tests {
			fmt.Println(method, paths)
			for _, p := range paths {
				req := httptest.NewRequest(method, p, nil)
				w := httptest.NewRecorder()
				testRouter.ServeHTTP(w, req)
				assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
			}
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	if testRouter == nil {
		testRouter = CreateRoute()
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", res.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", res.Header.Get("X-XSS-Protection"))
	assert.Equal(t, "default-src 'none'; frame-ancestors 'none';", res.Header.Get("Content-Security-Policy"))
	assert.Equal(t, "no-referrer", res.Header.Get("Referrer-Policy"))
}

func TestDisallowedMethods(t *testing.T) {
	if testRouter == nil {
		testRouter = CreateRoute()
	}

	disallowed := []string{"TRACE", "CONNECT"}
	for _, method := range disallowed {
		req := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestHeaderMasking(t *testing.T) {
	if testRouter == nil {
		testRouter = CreateRoute()
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer my-secret-token")
	req.Header.Set("Cookie", "session=12345")
	req.Header.Set("X-Api-Key", "my-api-key")
	req.Header.Set("X-Custom-Header", "safe-value")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&body)
	assert.Equal(t, nil, err)

	respHeaders, ok := body["request_headers"].(map[string]interface{})
	assert.Equal(t, true, ok)

	// Assert sensitive headers are masked
	assert.Equal(t, []interface{}{"[REDACTED]"}, respHeaders["Authorization"])
	assert.Equal(t, []interface{}{"[REDACTED]"}, respHeaders["Cookie"])
	assert.Equal(t, []interface{}{"[REDACTED]"}, respHeaders["X-Api-Key"])

	// Assert non-sensitive header is kept intact
	assert.Equal(t, []interface{}{"safe-value"}, respHeaders["X-Custom-Header"])
}

func TestBodySizeLimit(t *testing.T) {
	if testRouter == nil {
		testRouter = CreateRoute()
	}

	// Generate a body slightly larger than 1MB
	largeBody := make([]byte, (1<<20)+100)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(largeBody))

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
