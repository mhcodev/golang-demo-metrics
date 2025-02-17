package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var apiRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "api",
		Name:      "request_total",
		Help:      "Total number of API requests",
	},
	[]string{"method", "path"},
)

var requestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "api",
		Name:      "request_duration_seconds",
		Help:      "Duration of HTTP requests",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

func RecordMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/metrics") {
			next.ServeHTTP(w, r)
			return
		}

		apiRequests.WithLabelValues(r.Method, r.URL.Path).Inc()

		max := 10.0
		min := 0.005

		duration := min + rand.Float64()*(max-min)
		duration = math.Round(duration*100) / 100
		requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		next.ServeHTTP(w, r)
	})
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetCountryListHandler(w http.ResponseWriter, r *http.Request) {
	countries := []Country{
		{ID: 1, Name: "Chile"},
		{ID: 2, Name: "Perú"},
		{ID: 3, Name: "Brasil"},
		{ID: 4, Name: "Argentina"},
		{ID: 5, Name: "México"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func main() {
	mux := http.NewServeMux()
	prometheus.MustRegister(apiRequests)
	prometheus.MustRegister(requestDuration)

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/country", GetCountryListHandler)

	defaultMiddleware := RecordMetrics(mux)

	if err := http.ListenAndServe(":4000", defaultMiddleware); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
