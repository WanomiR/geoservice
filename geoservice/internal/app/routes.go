package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/http/pprof"
)

var requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "total_requests_count",
	Help: "Number of requests received.",
}, []string{"requests"})

func init() {
	prometheus.MustRegister(requestsTotal)
}

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	// group for gathering metrics, doesn't include the `metrics` endpoint
	r.Group(func(r chi.Router) {

		// count number of requests for all endpoints
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				requestsTotal.WithLabelValues("requests").Inc()
			})
		})

		r.Route("/address", func(r chi.Router) {
			r.Post("/search", a.controllers.Geo.AddressSearch)
			r.Post("/geocode", a.controllers.Geo.AddressGeocode)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Get("/login", a.controllers.Auth.Login)
			r.Get("/logout", a.controllers.Auth.Logout)
		})

		r.Route("/debug/pprof/", func(r chi.Router) {
			r.Use(a.services.Auth.RequireAuthorization)

			r.Get("/", pprof.Index)
			r.Get("/{cmd}", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
		})
	})

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.config.host, a.config.port)),
	))

	return r
}
