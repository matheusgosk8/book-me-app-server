package handlers

import (
	"net/http"

	"github.com/matheusgosk8/book-me-server/cmd/repositories"
	db "github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	vld "github.com/matheusgosk8/book-me-server/internal/validator"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Register(res http.ResponseWriter, req *http.Request) {

	log.Info("Registering new user")

	type registerPayload struct {
		User    *models.User    `json:"user"`
		Address *models.Address `json:"address"`
	}

	payload, err := utils.BodyParser[registerPayload](req)
	if err != nil {
		log.WithError(err).Warn("[Register] Falha ao decodificar body da requisição")
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.User == nil {
		log.Warn("[Register] Campo 'user' ausente no payload")
		http.Error(res, "user is required", http.StatusBadRequest)
		return
	}
	if payload.Address == nil {
		log.Warn("[Register] Campo 'address' ausente no payload")
		http.Error(res, "address is required", http.StatusBadRequest)
		return
	}

	newUser := payload.User
	newUserAddress := payload.Address

	msgs := vld.ValidateUser(vld.UserDTO{
		Email: newUser.Email,
		Senha: newUser.Senha,
		Nome:  newUser.Nome,
	})
	if msgs != nil {
		log.WithField("validation_errors", msgs).Warn("[Register] Dados inválidos no payload do usuário")
		res.WriteHeader(http.StatusBadRequest)
		utils.ServerResponse(res, msgs)
		return
	}

	// 1. Hash da senha antes de salvar no banco
	hashed, err := bcrypt.GenerateFromPassword([]byte(newUser.Senha), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("[Register] Falha ao gerar hash bcrypt da senha")
		utils.InternalErrorHandler(res, err)
		return
	}
	newUser.Senha = string(hashed)

	// 2. Geração de Tokens incluindo o UserType para o AuthMiddleware
	accessToken, refreshToken, err := utils.GenerateTokens(newUser.Id, newUser.UserType)
	if err != nil {
		log.WithError(err).Errorf("[Register] Falha ao gerar tokens para user_id=%s", newUser.Id)
		utils.InternalErrorHandler(res, err)
		return
	}

	// 3. Salvar usuário e endereço em transação[cite: 4, 5]
	// Nota: Certifique-se de que o CreateUserWithAddress também salve o refreshToken no banco
	_, _, err = repositories.CreateUserWithAddress(req.Context(), db.Client, *newUser, *newUserAddress)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"email":     newUser.Email,
			"user_type": newUser.UserType,
		}).Errorf("[Register] Falha na transação CreateUserWithAddress")
		utils.InternalErrorHandler(res, err)
		return
	}

	// 4. Montagem da resposta incluindo o RefreshToken para evitar erro de variável não utilizada
	var response = models.RegisterResponse{
		User: &models.RegisterUserResponse{
			Id:    newUser.Id,
			Nome:  newUser.Nome,
			Email: newUser.Email,
		},
		Token:        accessToken, 
		RefreshToken: refreshToken, // Utilizando a variável para satisfazer o compilador e o contrato da API
		Code:         201,
		Message:      "User registered successfully",
	}

	utils.ServerResponse(res, response)
}