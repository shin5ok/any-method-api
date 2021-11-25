package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

var paths = []string{"/", "/:p1", "/:p1/:p2"}
var port = os.Getenv("PORT")
var forceSleep = os.Getenv("FORCE_SLEEP")

func CreateRoute() *gin.Engine {
	g := gin.Default()

	for _, path := range paths {
		g.GET(path, commonHandler)
		g.POST(path, commonHandler)
		g.PUT(path, commonHandler)
		g.HEAD(path, commonHandler)
		g.PATCH(path, commonHandler)
		g.OPTIONS(path, commonHandler)
		g.DELETE(path, commonHandler)
	}

	return g

}

func main() {
	if port == "" {
		port = "8080"
	}

	g := CreateRoute()

	p := ginprom.New(
		ginprom.Engine(g),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	g.Use(p.Instrument())
	g.Run(":" + port)
}

func commonHandler(c *gin.Context) {
	method := c.Request.Method
	headers := c.Request.Header
	isSleep := c.Request.FormValue("s")
	var r time.Duration
	if isSleep != "" || forceSleep != "" {
		r = randSleeping()
	}
	path := c.Request.URL.Path
	// log.Printf("path: %s", path)
	c.JSON(http.StatusOK, gin.H{"method": method, "request_headers": headers, "path": path, "sleep": r})
}

func randSleeping() time.Duration {
	r := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(r)
	return r
}
