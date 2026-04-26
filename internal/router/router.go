package router

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/matheusgosk8/book-me-server/internal/config"
	"github.com/matheusgosk8/book-me-server/internal/handlers"
	"github.com/matheusgosk8/book-me-server/internal/middleware"
	public "github.com/matheusgosk8/book-me-server/internal/router/public"
)

func Router(r *chi.Mux) {
	// Middlewares Globais
	r.Use(cors.Handler(config.GetCORSConfig()))
	r.Use(chimiddle.StripSlashes)
	r.Use(middleware.LogRoute)

	// ROTAS PÚBLICAS
	r.Mount("/public", public.PublicRouter())
	r.Post("/refresh", handlers.RefreshHandler)

	// ROTAS PROTEGIDAS
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/me", handlers.GetMeHandler)
	})
}
