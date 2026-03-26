package handlers

import (
	"net/http"

	db "github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Register(res http.ResponseWriter, req *http.Request) {

	log.Info("Starting to register new user")

	newUser, err := utils.BodyParser[models.User](req)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	newUser.Id = utils.GenerateID()

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

	// Salvar usuário no banco
	_, err = utils.CreateUser(req.Context(), db.Client, *newUser)
	if err != nil {
		utils.InternalErrorHandler(res, err)
		return
	}

	var response = models.RegisterResponse{
		User:    newUser,
		Token:   token,
		Code:    200,
		Message: "User register successfully",
	}

	utils.ServerResponse(res, response)
}
