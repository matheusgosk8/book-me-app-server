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
	Message string      `json:"message"`
	Token   string      `json:"token"`
	User    interface{} `json:"user"` 
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// 1. Decodifica o JSON do Thunder Client/App
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	u, err := db.Client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Only(r.Context())

	if err != nil {
		log.Errorf("Tentativa de login falhou (Email não encontrado): %v", err)
		http.Error(w, "E-mail ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 3. Valida a senha (Em breve colocaremos Bcrypt aqui)
	if u.Senha != req.Senha {
		log.Warnf("Senha incorreta para o email: %s", req.Email)
		http.Error(w, "E-mail ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// 4. GERA O JWT
	token, err := utils.GenerateJWT(u.ID.String()) // Convertemos o UUID para string
	if err != nil {
		log.Errorf("Erro ao gerar JWT: %v", err)
		http.Error(w, "Erro interno ao gerar acesso", http.StatusInternalServerError)
		return
	}

	// 5. Resposta de Sucesso
	w.Header().Set("Content-Type", "application/json")
	u.Senha = "" // Segurança: oculta a senha
	
	response := LoginResponse{
		Message: "Login realizado com sucesso",
		Token:   token,
		User:    u,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}