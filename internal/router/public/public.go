package public

import (
    "github.com/go-chi/chi"
    handlers "github.com/matheusgosk8/book-me-server/internal/handlers"
)

func PublicRouter() chi.Router {
    r := chi.NewRouter()
    r.Get("/status", handlers.Status)
    r.Post("/register", handlers.Register)
    r.Post("/login", handlers.Login) 

    return r
}