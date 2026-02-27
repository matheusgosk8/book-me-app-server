package router

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	routes "github.com/matheusgosk8/book-me-server/internal/handlers"
	"github.com/matheusgosk8/book-me-server/internal/middleware"
)

func Router(r *chi.Mux) {

	r.Use(chimiddle.StripSlashes)

	r.Use(middleware.LogRoute)

	//Routes
	r.Route("/", func(router chi.Router) {
		router.Get("/status", routes.Status)
		router.Post("/register", routes.Register)
	})
}
