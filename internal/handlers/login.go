package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Dados mockados para teste
	userID := "1"
	userName := "Jâiner Jardel"
	userRole := "provider"

	// Gerando o JWT real usando .env
	token, err := utils.GenerateJWT(userID)
	if err != nil {
		log.Errorf("Erro ao gerar JWT no Login: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Erro ao gerar token"})
		return
	}

	response := LoginResponse{
		Message: "Login successful",
		Token:   token,
		User: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Role string `json:"role"`
		}{
			ID:   1,
			Name: userName,
			Role: userRole,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}