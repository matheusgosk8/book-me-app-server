package handlers

import (
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := map[string]interface{}{
		"message": "Login realizado com sucesso!",
		"user": map[string]interface{}{
			"id":    1,
			"name":  "Wâiner Jardel",
			"email": "teste@teste.com",
			"role": "prestador",
		},
		"token": "token-de-teste-123",
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}