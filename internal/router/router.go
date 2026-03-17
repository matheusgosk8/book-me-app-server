package router

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/matheusgosk8/book-me-server/internal/middleware"
	public "github.com/matheusgosk8/book-me-server/internal/router/public"
)

func Router(r *chi.Mux) {

	r.Use(chimiddle.StripSlashes)
	r.Use(middleware.LogRoute)

	r.Mount("/public", public.PublicRouter())

	//Routes
	// 	r.Route("/public", func(router chi.Router) {
	// 		router.Get("/status", routes.Status)
	// 		router.Post("/register", routes.Register)
	// 	})
}
