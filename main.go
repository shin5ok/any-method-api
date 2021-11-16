package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	g.GET("/", commonHandler)
	g.POST("/", commonHandler)
	g.PUT("/", commonHandler)
	g.HEAD("/", commonHandler)
	g.PATCH("/", commonHandler)
	g.OPTIONS("/", commonHandler)
	g.DELETE("/", commonHandler)
	g.GET("/:p", commonHandler)
	g.POST("/:p", commonHandler)
	g.PUT("/:p", commonHandler)
	g.HEAD("/:p", commonHandler)
	g.PATCH("/:p", commonHandler)
	g.OPTIONS("/:p", commonHandler)
	g.DELETE("/:p", commonHandler)

	g.Run(":8080")
}

func commonHandler(c *gin.Context) {
	method := c.Request.Method
	headers := c.Request.Header
	path := c.Request.URL.Path
	log.Printf("path: %s", path)
	c.JSON(http.StatusOK, gin.H{"method": method, "headers": headers, "path": path})
}
