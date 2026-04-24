package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type LoginResponse struct {
	Message      string      `json:"message"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// 1. Decodifica o JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// 2. Busca o usuário
	u, err := db.Client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(r.Context())

	if err != nil {
		log.Errorf("Tentativa de login falhou (Email não encontrado): %v", err)
		http.Error(w, "E-mail ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 3. Valida a senha
	if u.Senha != req.Senha {
		log.Warnf("Senha incorreta para o email: %s", req.Email)
		http.Error(w, "E-mail ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 4. Gera Access Token e Refresh Token
	accessToken, refreshToken, err := utils.GenerateTokens(u.ID.String())
	if err != nil {
		log.Errorf("Erro ao gerar tokens: %v", err)
		http.Error(w, "Erro interno ao gerar acesso", http.StatusInternalServerError)
		return
	}

	// 5. Salva o Refresh Token no Banco de Dados
	_, err = db.Client.User.
		UpdateOne(u).
		SetRefreshToken(refreshToken).
		Save(r.Context())

	if err != nil {
		log.Errorf("Erro ao salvar refresh token no banco: %v", err)
		http.Error(w, "Erro ao finalizar login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	u.Senha = "" // Segurança: oculta a senha

	response := LoginResponse{
		Message:      "Login realizado com sucesso",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         u,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}