package public

import (
    "github.com/go-chi/chi"
    "github.com/go-chi/cors"
    handlers "github.com/matheusgosk8/book-me-server/internal/handlers"
)

func PublicRouter() chi.Router {
    r := chi.NewRouter()
    
    r.Use(cors.Handler(cors.Options{
        // Mude de localhost para "*" para o celular conseguir conectar
        AllowedOrigins:   []string{"*"}, 
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    r.Get("/status", handlers.Status)
    r.Post("/register", handlers.Register)
    r.Post("/login", handlers.Login) 

    return r
}