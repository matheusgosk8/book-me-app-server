package router

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/matheusgosk8/book-me-server/internal/config"
	"github.com/matheusgosk8/book-me-server/internal/handlers/address"
	"github.com/matheusgosk8/book-me-server/internal/handlers/auth"
	"github.com/matheusgosk8/book-me-server/internal/handlers/categories"
	"github.com/matheusgosk8/book-me-server/internal/handlers/services"
	"github.com/matheusgosk8/book-me-server/internal/handlers/user"
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
	r.Post("/refresh", auth.RefreshHandler)

	// ROTAS PROTEGIDAS
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/me", user.GetMeHandler)
		r.Get("/addresses/me", address.ListMyAddresses)
		//SEÇÃO PROVIDER
		r.Post("/provider/services", services.CreateServiceHandler)
		r.Get("/provider/services", services.ListServices)
		r.Get("/provider/services/me", services.ListMyServices)
		r.Patch("/provider/services/{id}", services.UpdateServiceHandler)
		r.Delete("/provider/services/{id}", services.DeleteServiceHandler)

		//SEÇÃO CUSTOMER
		r.Get("/customer/services", services.ListServices)

		//Categories
		r.Get("/categories", categories.ListCategoriesHandler)

	})
}
