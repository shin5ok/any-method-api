package main

import (
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
		{"GET": []string{"/", "/p"}},
		{"POST": []string{"/", "/p"}},
		{"PUT": []string{"/", "/p"}},
		{"DELETE": []string{"/", "/p"}},
		{"OPTIONS": []string{"/", "/p"}},
		{"HEAD": []string{"/", "/p"}},
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
}
