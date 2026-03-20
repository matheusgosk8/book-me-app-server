package router

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors" // <--- Adicione este import
	"github.com/matheusgosk8/book-me-server/internal/middleware"
	public "github.com/matheusgosk8/book-me-server/internal/router/public"
)

func Router(r *chi.Mux) {

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimiddle.StripSlashes)
	r.Use(middleware.LogRoute)

	r.Mount("/public", public.PublicRouter())
}