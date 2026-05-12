package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Token inválido", http.StatusBadRequest)
		return
	}

	// 1. Valida o Token e extrai o ID do usuário
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil || claims.Type != "refresh" {
		http.Error(w, "Refresh token inválido ou expirado", http.StatusUnauthorized)
		return
	}

	// 2. Busca o usuário no banco para conferir o token e obter o UserType[cite: 5]
	u, err := db.Client.User.
		Query().
		Where(user.IDEQ(utils.ParseUUID(claims.UserId))).
		Only(r.Context())

	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	// 3. SEGURANÇA: Verifica se o token enviado é o mesmo que está no banco[cite: 5]
	if u.RefreshToken != req.RefreshToken {
		log.Warnf("Tentativa de refresh com token antigo/revogado para user: %s", u.Email)
		http.Error(w, "Token revogado", http.StatusUnauthorized)
		return
	}

	// Alteração: Adicionado u.UserType como segundo argumento para satisfazer a nova assinatura[cite: 9]
	newAccessToken, _, err := utils.GenerateTokens(u.ID.String(), u.UserType)
	if err != nil {
		log.WithError(err).Error("[Refresh] Falha ao regenerar tokens")
		http.Error(w, "Erro ao renovar acesso", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}