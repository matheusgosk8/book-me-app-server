package payloads

import "github.com/matheusgosk8/book-me-server/internal/models"

func GetRegisterPayload() map[string]interface{} {
	return map[string]interface{}{
		"user": models.User{
			Id:            "test-user-id",
			Nome:          "Test User",
			Email:         "test@e2e.com",
			Senha:         "password123",
			ConfirmaSenha: "password123",
			UserType:      "customer",
			Telefone:      "11999999999",
		},
		"address": models.Address{
			Street:     "Rua de Teste",
			City:       "São Paulo",
			State:      "SP",
			PostalCode: "01234567",
			Country:    "Brasil",
		},
	}
}
