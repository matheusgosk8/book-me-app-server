package public

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	handlers "github.com/matheusgosk8/book-me-server/internal/handlers"
)

func PublicRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/status", handlers.Status)
	r.Post("/register", handlers.Register)

	// Criar rota pública de login
	// mocke os dados do usuário -> id, nome, email, tipo de usuário (cliente ou prestador) e token JWT
	// r.Post("/login", handlers.Login)

	return r

}
