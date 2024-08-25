package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var requestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: prometheusNamespace,
	Name:      "total_requests_count",
	Help:      "Total number of requests received.",
}, []string{"type"})

var requestsLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: prometheusNamespace,
	Name:      "response_latency_seconds",
	Help:      "Histogram of response times in seconds.",
	Buckets:   []float64{0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048},
}, []string{"route", "method"})

func init() {
	prometheus.MustRegister(requestsCount)
	prometheus.MustRegister(requestsLatency)
}

func (a *App) RequestsCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		requestsCount.With(prometheus.Labels{"type": "request"}).Inc()
	})
}

func (a *App) RequestsLatency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		requestsLatency.
			With(prometheus.Labels{"route": r.Pattern, "method": r.Method}).
			Observe(time.Since(start).Seconds())
	})
}
