package repositories

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
	api "github.com/matheusgosk8/book-me-server/internal/models"
	log "github.com/sirupsen/logrus"
)

// optionalString maps empty input to nil so Ent skips the field (NULL in DB instead of "").

// CreateUserWithAddress creates a user and an address in a single transaction.
func CreateUserWithAddress(ctx context.Context, client *ent.Client, uParams api.User, aParams api.Address) (*ent.User, *ent.Address, error) {
	var usr *ent.User
	var addr *ent.Address
	err := RunInTx(ctx, client, func(ctx context.Context, tx *ent.Tx) error {
		var err error

		log.WithField("email", uParams.Email).Info("[TX] Criando usuário no banco")
		usr, err = tx.User.
			Create().
			SetNillableCep(optionalString(uParams.Cep)).
			SetNillableCidade(optionalString(uParams.Cidade)).
			SetNillableCnpj(optionalString(uParams.Cnpj)).
			SetNillableConfirmaSenha(optionalString(uParams.ConfirmaSenha)).
			SetNillableCpf(optionalString(uParams.Cpf)).
			SetEmail(uParams.Email).
			SetNillableEstado(optionalString(uParams.Estado)).
			SetNillableLogradouro(optionalString(uParams.Logradouro)).
			SetNome(uParams.Nome).
			SetNillableRua(optionalString(uParams.Rua)).
			SetSenha(uParams.Senha).
			SetNillableTelefone(optionalString(uParams.Telefone)).
			SetUserType(uParams.UserType).
			Save(ctx)
		if err != nil {
			log.WithError(err).WithField("email", uParams.Email).Error("[TX] Falha ao criar usuário — rollback será executado")
			return err
		}
		log.WithField("user_id", usr.ID).Info("[TX] Usuário criado, criando endereço")

		addr, err = tx.Address.
			Create().
			SetStreet(aParams.Street).
			SetCity(aParams.City).
			SetState(aParams.State).
			SetPostalCode(aParams.PostalCode).
			SetCountry(aParams.Country).
			SetUser(usr).
			Save(ctx)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"user_id": usr.ID,
				"city":    aParams.City,
				"street":  aParams.Street,
			}).Error("[TX] Falha ao criar endereço — rollback será executado")
			return err
		}
		return nil
	})
	if err != nil {
		log.WithError(err).Error("[CreateUserWithAddress] Transação falhou e foi revertida")
		return nil, nil, err
	}
	log.WithFields(log.Fields{
		"user_id": usr.ID,
		"addr_id": addr.ID,
	}).Info("[CreateUserWithAddress] Usuário e endereço criados com sucesso")
	return usr, addr, nil
}
