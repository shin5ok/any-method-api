package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
		"random", Rand500div,
	).Msg("begin to create the route")
	if Rand500div != "" {
		if n, err := strconv.Atoi(Rand500div); err == nil {
			Rand500int = n
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
	c.JSON(code, resultData)
}

func randSleeping() time.Duration {
	r := time.Duration(genRand()) * time.Millisecond
	time.Sleep(r)
	return r
}

func isRand500(n int) bool {
	return genRand()%n == 0
}

func dummy() {
	genRand()
}

func genRand() int {
	r := rand.Intn(100)
	// log.Printf("generated rand value: %d\n", r)
	return r
}
