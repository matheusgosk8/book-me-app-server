package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/db"
	log "github.com/sirupsen/logrus"
)

type RegisterRequest struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Senha    string `json:"senha"`
	Telefone string `json:"telefone"`
	UserType string `json:"user_type"` // Ex: "client" ou "provider"
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("Erro ao decodificar JSON: %v", err)
		http.Error(w, "Formato de dados inválido", http.StatusBadRequest)
		return
	}

	newUser, err := db.Client.User.Create().
		SetNome(req.Nome).
		SetEmail(req.Email).
		SetSenha(req.Senha).
		SetTelefone(req.Telefone).
		SetUserType(req.UserType).
		Save(r.Context())

	if err != nil {
		log.Errorf("Erro ao criar usuário no banco: %v", err)
		http.Error(w, "Erro ao processar cadastro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	// Oculta a senha no retorno
	newUser.Senha = ""
	json.NewEncoder(w).Encode(newUser)
}
