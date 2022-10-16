package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var paths = []string{"/", "/:p1", "/:p1/:p2"}
var randValue int
var promPort = os.Getenv("PROM_PORT")
var servicePort = os.Getenv("PORT")
var randDiv = os.Getenv("RAND_DIV")
var defectMode = os.Getenv("MODE")

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
	if servicePort == "" {
		servicePort = "8080"
	}

	forService := CreateRoute()
	forExporter := gin.Default()

	m := ginmetrics.GetMonitor()
	// use metric middleware without expose metric path
	m.UseWithoutExposingEndpoint(forService)
	m.SetMetricPath("/metrics")
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Expose(forExporter)

	go func() {
		_ = forExporter.Run(":10080")
	}()

	_ = forService.Run(":" + servicePort)

}

func commonHandler(c *gin.Context) {
	method := c.Request.Method
	headers := c.Request.Header
	path := c.Request.URL.Path
	code := http.StatusOK
	var r time.Duration
	var resultData = gin.H{"method": method, "request_headers": headers, "path": path, "sleep": r}
	if isRand(randValue) {
		switch defectMode {
		case "sleep":
			r := randSleeping()
			resultData["sleep"] = r.String()
			log.Info().Str("wait duration", r.String()).Send()
		case "error":
			code = http.StatusServiceUnavailable
			resultData = gin.H{}
		default:
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

func isRand(n int) bool {
	if n == 0 {
		return false
	}
	return genRand()%n == 0
}

func genRand() int {
	r := rand.Intn(1000-500) + 500
	return r
}
