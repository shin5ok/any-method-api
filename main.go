package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var paths = []string{"/", "/{p1}", "/{p1}/{p2}"}
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

func CreateRoute() *chi.Mux {
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Healthcheck", "always ok")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
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
	if port == "" {
		port = "8080"
	}

	r := CreateRoute()
	log.Info().Msgf("Server starting on port %s", port)
	http.ListenAndServe(":"+port, r)
}

func commonHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	headers := r.Header
	path := r.URL.Path
	code := http.StatusOK
	var sleep time.Duration

	resultData := map[string]interface{}{
		"method":          method,
		"request_headers": headers,
		"path":            path,
		"sleep":           sleep,
	}

	if Dummy != "" {
		dummy()
	}

	// if ForceSleep != "" {
	// 	sleep = randSleeping()
	// 	resultData["sleep"] = sleep
	// }

	if Rand500int >= 1 {
		if isRand500(Rand500int) {
			code = http.StatusServiceUnavailable
			resultData = map[string]interface{}{}
		}
	}

	log.Info().Str("Path", path).Str("Method", method).Send()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resultData)
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
