package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(a.facade.RequireAuthentication)

		r.Route("/authors", func(r chi.Router) {
			r.Get("/", a.controller.GetAllAuthors)
			r.Post("/", a.controller.CreateAuthor)
			r.Get("/top", a.controller.GetTopAuthors)
			r.Get("/{authorId}", a.controller.GetAuthor)
			r.Delete("/{authorId}", a.controller.DeleteAuthor)
		})

		r.Route("/books", func(r chi.Router) {
			r.Get("/", a.controller.GetAllBooks)
			r.Post("/", a.controller.CreateBook)
			r.Get("/{bookId}", a.controller.GetBook)
			r.Delete("/{bookId}", a.controller.DeleteBook)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", a.controller.GetAllUsers)
		r.Post("/", a.controller.CreateUser)
		r.Get("/{userId}", a.controller.GetUser)
		r.Get("/{userId}/rent/{bookId}", a.controller.RentBook)
		r.Get("/{userId}/return/{bookId}", a.controller.ReturnBook)
		r.Post("/login", a.controller.LoginUser)
		r.Get("/logout", a.controller.LogoutUser)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.host, a.port)),
	))

	return r
}
