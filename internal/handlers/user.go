package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	"github.com/matheusgosk8/book-me-server/ent/user"
)

func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Recupera o ID do usuário que o Middleware salvou no contexto
	ctxUserID := r.Context().Value("user_id").(string)

	// 2. Busca no banco
	u, err := db.Client.User.
		Query().
		Where(user.IDEQ(utils.ParseUUID(ctxUserID))).
		Only(r.Context())

	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	// 3. Retorna os dados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    u.ID,
		"nome":  u.Nome,
		"email": u.Email,
		"type":  u.UserType,
	})
}