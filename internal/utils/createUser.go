package utils

import (
	"context"
	"fmt"

	"github.com/matheusgosk8/book-me-server/ent"
	api "github.com/matheusgosk8/book-me-server/internal/models"
	log "github.com/sirupsen/logrus"
)

func CreateUser(ctx context.Context, client *ent.Client, params api.User) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetCep(params.Cep).
		SetCidade(params.Cidade).
		SetCnpj(params.Cnpj).
		SetConfirmaSenha(params.ConfirmaSenha).
		SetCpf(params.Cpf).
		SetEmail(params.Email).
		SetEstado(params.Estado).
		SetLogradouro(params.Logradouro).
		SetNome(params.Nome).
		SetRua(params.Rua).
		SetSenha(params.Senha).
		SetTelefone(params.Telefone).
		SetUserType(params.UserType).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}
