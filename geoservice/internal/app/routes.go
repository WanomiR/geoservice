package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http/pprof"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	// group for applying middleware
	r.Group(func(r chi.Router) {

		// apply middleware
		r.Use(a.ZapLogger)
		r.Use(a.RequestsCounter)
		r.Use(a.RequestsLatency)

		// main api endpoints
		r.Route("/address", func(r chi.Router) {
			r.Post("/search", a.controllers.Geo.AddressSearch)
			r.Post("/geocode", a.controllers.Geo.AddressGeocode)
		})

		// authorization endpoints
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", a.controllers.Auth.Register)
			r.Post("/login", a.controllers.Auth.Login)
			r.Get("/logout", a.controllers.Auth.Logout)
		})

	})

	// profiling endpoints
	r.Route("/debug/pprof/", func(r chi.Router) {
		// available only to authorized users
		r.Use(a.services.Auth.RequireAuthorization)

		r.Get("/", pprof.Index)
		r.Get("/{cmd}", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
	})

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.config.host, a.config.port)),
	))

	return r
}
