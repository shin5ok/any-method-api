package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	paths       = []string{"/", "/:p1", "/:p1/:p2"}
	randValue   int
	promPort    = "10080"
	servicePort = os.Getenv("PORT")
	randDiv     = os.Getenv("RAND_DIV")
	defectMode  = os.Getenv("MODE")
	metricName  = "common_handler_latency_hist"
)

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
	g.Use(ginzerolog.Logger("gin"))

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

	forExporter := gin.Default()

	forService := CreateRoute()
	// use metric middleware without expose metric path
	gaugeMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Histogram,
		Name:        "common_handler_latency_hist",
		Description: "an example of gauge type metric",
		Buckets:     []float64{0.1, 10, 50, 100, 500, 1000, 2000, 3000},
		Labels:      []string{"actual_ms"},
	}

	m := ginmetrics.GetMonitor()
	m.AddMetric(gaugeMetric)
	m.SetMetricPath("/metrics")
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.UseWithoutExposingEndpoint(forService)
	m.Expose(forExporter)

	go func() {
		_ = forExporter.Run(":" + promPort)
	}()

	_ = forService.Run(":" + servicePort)

}

func commonHandler(c *gin.Context) {
	start := time.Now()
	defer func(start time.Time) {
		go func() {
			end := time.Now()
			sub := end.Sub(start)
			uriStr := strings.Replace(c.Request.RequestURI, "/", "_", -1)
			label := fmt.Sprintf("%s_%s_actual_ms", strings.ToLower(c.Request.Method), strings.ToLower(uriStr))
			err := ginmetrics.GetMonitor().GetMetric(metricName).Observe([]string{label}, float64(sub.Milliseconds()))

			if err != nil {
				log.Error().Err(err).Send()
			}
		}()
	}(start)
	method := c.Request.Method
	headers := c.Request.Header
	path := c.Request.URL.Path
	code := http.StatusOK
	var r time.Duration
	var resultData = gin.H{"method": method, "request_headers": headers, "path": path, "sleep": r}
	if isDefeating(randValue) {
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
	log.Info().Str("URI Path", path).Str("Method", method).Send()
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
	r := rand.Intn(1000-500) + 500
	return r
}
