package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var paths = []string{"/", "/:p"}
var port = os.Getenv("PORT")

func main() {
	g := gin.Default()

	for path := range paths {
		p := string(path)
		g.GET(p, commonHandler)
		g.POST(p, commonHandler)
		g.PUT(p, commonHandler)
		g.HEAD(p, commonHandler)
		g.PATCH(p, commonHandler)
		g.OPTIONS(p, commonHandler)
		g.DELETE(p, commonHandler)
	}

	g.Run(":" + port)
}

func commonHandler(c *gin.Context) {
	method := c.Request.Method
	headers := c.Request.Header
	path := c.Request.URL.Path
	log.Printf("path: %s", path)
	c.JSON(http.StatusOK, gin.H{"method": method, "headers": headers, "path": path})
}
