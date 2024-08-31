package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "proxy/docs"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(a.ZapLogger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/address", func(r chi.Router) {
			r.Post("/search", nil)
			r.Post("/geocode", nil)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", "localhost", a.config.Port)),
	))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello world"))
	})

	return r
}
