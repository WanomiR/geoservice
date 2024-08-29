package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello world"))
	})

	return r
}
