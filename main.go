package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var paths = []string{"/", "/:p1", "/:p1/:p2"}
var port = os.Getenv("PORT")
var ForceSleep = os.Getenv("SLEEP")
var Dummy = os.Getenv("DUMMY")
var Rand500int = 0
var Rand500div = os.Getenv("RAND500DIV")

func init() {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Info().Msg("init")

}

func CreateRoute() *gin.Engine {
	log.Info().Str(
		"random", randDiv,
	).Msg("begin to create the route")

	if randDiv != "" {
		if n, err := strconv.Atoi(randDiv); err == nil {
			randValue = n
		} else {
			log.Error().Msg("error:" + err.Error())
		}
	}
	g := gin.Default()

	g.GET("/test", func(c *gin.Context) {
		c.Header("X-Healthcheck", "always ok")
		c.JSON(http.StatusOK, gin.H{})
	})

	for _, path := range paths {
		r.Get(path, commonHandler)
		r.Post(path, commonHandler)
		r.Put(path, commonHandler)
		r.Head(path, commonHandler)
		r.Patch(path, commonHandler)
		r.Options(path, commonHandler)
		r.Delete(path, commonHandler)
	}

	// Add prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func main() {
	if servicePort == "" {
		servicePort = "8080"
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
	path := c.Request.URL.Path
	code := http.StatusOK
	var r time.Duration
	var resultData = gin.H{"method": method, "request_headers": headers, "path": path, "sleep": r}
	if Dummy != "" {
		dummy()
	}
	// if ForceSleep != "" {
	// 	r = randSleeping()
	// 	resultData["sleep"] = r
	// }
	if Rand500int >= 1 {
		if isRand500(Rand500int) {
			code = http.StatusServiceUnavailable
			resultData = gin.H{}
		}
	}
	log.Info().Str("Path", path).Str("Method", method).Send()
	c.JSON(code, resultData)
}

func randSleeping() time.Duration {
	r := time.Duration(genRand()) * time.Millisecond
	time.Sleep(r)
	return r
}

func isDefeating(n int) bool {
	if n == 0 {
		return false
	}
	return genRand()%n == 0
}

func genRand() int {
	// generating random int from 500 to 1000
	r := rand.Intn(1000-500) + 500
	return r
}
