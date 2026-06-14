package auth

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/utils"
	service "github.com/matheusgosk8/book-me-server/internal/service/auth"
)

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ServerError(w, http.StatusBadRequest, err)
		return
	}

	loginService := service.NewLoginService()
	output, err := loginService.Login(r.Context(), service.LoginInput{
		Email: req.Email,
		Senha: req.Senha,
	})

	if err != nil {
		utils.ServerError(w, http.StatusUnauthorized, err)
		return
	}

	response := map[string]interface{}{
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"user":          output.User,
	}

	utils.ServerSuccess(w, http.StatusOK, "Login realizado com sucesso", response)
}
