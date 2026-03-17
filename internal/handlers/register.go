package handlers

import (
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Register(res http.ResponseWriter, req *http.Request) {

	log.Info("Starting to register new user")

	//campos email, e name
	newUser, err := utils.BodyParser[models.User](req)
	newUser.Id = utils.GenerateID()

	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response = models.RegisterResponse{
		User:    newUser,
		Code:    200,
		Message: "User register successfully",
	}
	utils.ServerResponse(res, response)
}
