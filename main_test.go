package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var testg *gin.Engine

func TestPing(t *testing.T) {
	if testg == nil {
		testg = CreateRoute()
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	testCombi := []map[string][]string{
		{"GET": []string{"/", "/p"}},
		{"POST": []string{"/", "/p"}},
		{"PUT": []string{"/", "/p"}},
		{"DELETE": []string{"/", "/p"}},
		{"OPTIONS": []string{"/", "/p"}},
		{"HEAD": []string{"/", "/p"}},
	}
	for _, tests := range testCombi {

		for k, paths := range tests {
			fmt.Println(k, paths)
			for _, p := range paths {
				c.Request, _ = http.NewRequest(http.MethodGet, p, nil)
				testg.ServeHTTP(w, c.Request)
			}
		}

		// assert.Equal(t, http.StatusOK, w.Code)

	}
}
