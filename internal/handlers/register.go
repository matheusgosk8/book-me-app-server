package handlers

import (
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
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

	var response = models.RegisterResponse{
		User:    newUser,
		Token:   token,
		Code:    200,
		Message: "User register successfully",
	}

	utils.ServerResponse(res, response)
}
