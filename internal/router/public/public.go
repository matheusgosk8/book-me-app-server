package public

import (
	"github.com/go-chi/chi"
	auth "github.com/matheusgosk8/book-me-server/internal/handlers/auth"
	categories "github.com/matheusgosk8/book-me-server/internal/handlers/categories"
	system "github.com/matheusgosk8/book-me-server/internal/handlers/system"
)

func PublicRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/status", system.Status)
	r.Post("/register", auth.Register)
	r.Post("/login", auth.LoginHandler)
	r.Get("/categories", categories.ListCategoriesHandler)
	return r
}
