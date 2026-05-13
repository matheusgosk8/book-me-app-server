package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/matheusgosk8/book-me-server/internal/utils"
)

// AuthMiddleware intercepta a requisição para validar o JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Pega o header "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token não fornecido", http.StatusUnauthorized)
			return
		}

		// 2. Verifica o formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		// 3. Valida o token usando a função utilitária
		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		// 4. Segurança extra: Garante que o token é de acesso, não de refresh
		if claims.Type != "access" {
			http.Error(w, "Tipo de token inválido para esta rota", http.StatusUnauthorized)
			return
		}

		// 5. Injeta o ID e o TIPO do usuário no Contexto da requisição
		// Importante: UserType deve estar presente nos claims para evitar erros nos Handlers
		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "user_type", claims.UserType)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
