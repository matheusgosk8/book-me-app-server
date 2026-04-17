package repositories

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
	api "github.com/matheusgosk8/book-me-server/internal/models"
	log "github.com/sirupsen/logrus"
)

// CreateUserWithAddress creates a user and an address in a single transaction.
func CreateUserWithAddress(ctx context.Context, client *ent.Client, uParams api.User, aParams api.Address) (*ent.User, *ent.Address, error) {
	var usr *ent.User
	var addr *ent.Address
	err := RunInTx(ctx, client, func(ctx context.Context, tx *ent.Tx) error {
		var err error
		usr, err = tx.User.
			Create().
			SetCep(uParams.Cep).
			SetCidade(uParams.Cidade).
			SetCnpj(uParams.Cnpj).
			SetConfirmaSenha(uParams.ConfirmaSenha).
			SetCpf(uParams.Cpf).
			SetEmail(uParams.Email).
			SetEstado(uParams.Estado).
			SetLogradouro(uParams.Logradouro).
			SetNome(uParams.Nome).
			SetRua(uParams.Rua).
			SetSenha(uParams.Senha).
			SetTelefone(uParams.Telefone).
			SetUserType(uParams.UserType).
			Save(ctx)
		if err != nil {
			return err
		}

		addr, err = tx.Address.
			Create().
			SetStreet(aParams.Street).
			SetCity(aParams.City).
			SetState(aParams.State).
			SetPostalCode(aParams.PostalCode).
			SetCountry(aParams.Country).
			Save(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	log.Println("user and address created in transaction: ", usr, addr)
	return usr, addr, nil
}
