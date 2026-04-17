package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserDTO struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required,min=6"`
	Nome  string `json:"nome" validate:"required"`
}

func ValidateUser(u UserDTO) map[string]string {
	v := validator.New()
	err := v.Struct(u)
	if err == nil {
		return nil
	}
	verrs := err.(validator.ValidationErrors)
	msgs := map[string]string{}
	for _, e := range verrs {
		switch e.Tag() {
		case "required":
			msgs[e.Field()] = "campo obrigatório"
		case "email":
			msgs[e.Field()] = "email inválido"
		case "min":
			msgs[e.Field()] = fmt.Sprintf("mínimo %s caracteres", e.Param())
		default:
			msgs[e.Field()] = e.Error()
		}
	}
	return msgs
}
