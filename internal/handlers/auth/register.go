package auth

import (
	"errors"
	"net/http"

	"github.com/matheusgosk8/book-me-server/internal/models"
	service "github.com/matheusgosk8/book-me-server/internal/service/auth"
	"github.com/matheusgosk8/book-me-server/internal/utils"
	vld "github.com/matheusgosk8/book-me-server/internal/validator"
	log "github.com/sirupsen/logrus"
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

	if payload.User == nil || payload.Address == nil {
		http.Error(res, "user and address are required", http.StatusBadRequest)
		return
	}

	// 1. Validação (Ainda no handler ou pode ir pro service)
	msgs := vld.ValidateUser(vld.UserDTO{
		Email: payload.User.Email,
		Senha: payload.User.Senha,
		Nome:  payload.User.Nome,
	})
	if msgs != nil {
		res.WriteHeader(http.StatusBadRequest)
		utils.ServerResponse(res, msgs)
		return
	}

	// 2. Chamada do Service (Lógica de Negócio)
	userService := service.NewUserService()
	output, err := userService.Register(req.Context(), service.RegisterInput{
		User:    *payload.User,
		Address: *payload.Address,
	})

	if err != nil {
		// 3. Tratamento de erro especializado (Postgres)
		msg, status := utils.HandlePGError(err)
		log.WithError(err).Errorf("[Register] Erro no serviço: %s", msg)

		// Retorno padronizado: { "statusCode": 400/409, "message": "..." }
		utils.ServerError(res, status, errors.New(msg))
		return
	}

	// 4. Resposta
	utils.ServerSuccess(res, http.StatusCreated, "User registered successfully", models.RegisterResponse{
		User:         output.User,
		Token:        output.AccessToken,
		RefreshToken: output.RefreshToken,
		Code:         http.StatusCreated,
		Message:      "User registered successfully",
	})
}
