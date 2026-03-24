package router

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/matheusgosk8/book-me-server/internal/middleware"
	"github.com/matheusgosk8/book-me-server/internal/config"
	public "github.com/matheusgosk8/book-me-server/internal/router/public"

)

func Router(r *chi.Mux) {

	r.Use(cors.Handler(config.GetCORSConfig()))

	r.Use(cors.Handler(config.GetCORSConfig()))
	r.Use(chimiddle.StripSlashes)
	r.Use(middleware.LogRoute)

	r.Mount("/public", public.PublicRouter())
}