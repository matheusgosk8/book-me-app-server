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

	//Responsabilidade de dto retirada do body parser
	payload, err := utils.BodyParser[registerPayload](req)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.User == nil {
		http.Error(res, "user is required", http.StatusBadRequest)
		return
	}
	if payload.Address == nil {
		http.Error(res, "address is required", http.StatusBadRequest)
		return
	}

	// Validando campos obrigatórios
	newUser := payload.User
	newUserAddress := payload.Address

	// Mapear para DTO do validador e validar
	msgs := vld.ValidateUser(vld.UserDTO{
		Email: newUser.Email,
		Senha: newUser.Senha,
		Nome:  newUser.Nome,
	})
	if msgs != nil {
		res.WriteHeader(http.StatusBadRequest)
		utils.ServerResponse(res, msgs)
		return
	}

	newUser.Id = utils.GenerateID()
	newUserAddress.Id = utils.GenerateID()

	token, err := utils.GenerateJWT(newUser.Id)
	if err != nil {
		log.Errorf("ERRO AO GERAR JWT: %v", err) // <--- ESTA LINHA VAI SALVAR A GENTE
		utils.InternalErrorHandler(res, err)
		return
	}

	// Hash da senha
	hashed, err := bcrypt.GenerateFromPassword([]byte(newUser.Senha), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalErrorHandler(res, err)
		return
	}
	newUser.Senha = string(hashed)

	// Salvar usuário e endereço em transação
	//Em transaction, se uma das operações falhar, nenhuma é persistida, garantindo a integridade dos dados.
	_, _, err = repositories.CreateUserWithAddress(req.Context(), db.Client, *newUser, *newUserAddress)
	if err != nil {
		utils.InternalErrorHandler(res, err)
		return
	}

	// Return only minimal user info (id, nome, email)
	var response = models.RegisterResponse{
		User: &models.RegisterUserResponse{
			Id:    newUser.Id,
			Nome:  newUser.Nome,
			Email: newUser.Email,
		},
		Token:   token,
		Code:    200,
		Message: "User register successfully",
	}

	utils.ServerResponse(res, response)
}
