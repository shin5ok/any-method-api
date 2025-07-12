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
