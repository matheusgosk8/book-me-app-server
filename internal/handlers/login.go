package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// 1. Busca o usuário incluindo o campo user_type[cite: 5]
	u, err := db.Client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(r.Context())

	if err != nil {
		log.Errorf("Tentativa de login falhou (Email não encontrado): %v", err)
		http.Error(w, "E-mail ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 2. Valida a senha usando CompareHashAndPassword[cite: 5]
	err = bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(req.Senha))
	if err != nil {
		log.Warnf("Senha incorreta para o email: %s", req.Email)
		http.Error(w, "E-mail ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 3. Gera Access Token e Refresh Token com UserType para o Middleware[cite: 9]
	accessToken, refreshToken, err := utils.GenerateTokens(u.ID.String(), u.UserType)
	if err != nil {
		log.Errorf("Erro ao gerar tokens: %v", err)
		http.Error(w, "Erro interno ao gerar acesso", http.StatusInternalServerError)
		return
	}

	// 4. Salva o Refresh Token no Banco de Dados[cite: 5]
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