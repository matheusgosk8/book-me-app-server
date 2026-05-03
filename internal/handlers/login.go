package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// 1. Decodifica o JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ServerError(w, http.StatusBadRequest, errors.New("Dados inválidos"))
		return
	}

	// 2. Busca o usuário
	u, err := db.Client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(r.Context())

	if err != nil {
		log.Errorf("Tentativa de login falhou (Email não encontrado): %v", err)
		utils.ServerError(w, http.StatusUnauthorized, errors.New("E-mail ou senha incorretos"))
		return
	}

	// 3. Valida a senha
	if u.Senha != req.Senha {
		log.Warnf("Senha incorreta para o email: %s", req.Email)
		//Usar o compare
		err := bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(req.Senha))
		if err != nil {
			log.Warnf("Senha incorreta para o email: %s", req.Email)
			utils.ServerError(w, http.StatusUnauthorized, errors.New("E-mail ou senha incorretos"))
			return
		}
	}

	// 4. Gera Access Token e Refresh Token
	accessToken, refreshToken, err := utils.GenerateTokens(u.ID.String())
	if err != nil {
		log.Errorf("Erro ao gerar tokens: %v", err)
		utils.ServerError(w, http.StatusInternalServerError, err)
		return
	}

	// 5. Salva o Refresh Token no Banco de Dados
	_, err = db.Client.User.
		UpdateOne(u).
		SetRefreshToken(refreshToken).
		Save(r.Context())

	if err != nil {
		log.Errorf("Erro ao salvar refresh token no banco: %v", err)
		utils.ServerError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	u.Senha = "" // Segurança: oculta a senha

	userLoginResponse := models.LoginUserResponse{
		Id:    u.ID.String(),
		Nome:  u.Nome,
		Email: u.Email,
		Role:  u.UserType,
	}

	response := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userLoginResponse,
	}

	utils.ServerSuccess(w, http.StatusOK, "Login realizado com sucesso", response)
}
