package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var testg *gin.Engine

func TestPing(t *testing.T) {
	if testg == nil {
		testg = CreateRoute()
	}

	gin.SetMode(gin.TestMode)
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
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest(method, p, nil)
				testg.ServeHTTP(w, c.Request)
				assert.Equal(t, http.StatusOK, w.Code)
			}
		}

	}
}
