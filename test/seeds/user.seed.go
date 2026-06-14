package seeds

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/user"
)

// GetUserByEmail retorna o usuário com o e-mail informado ou um erro caso não exista
func GetUserByEmail(client *ent.Client, email string) (*ent.User, error) {
	ctx := context.Background()
	return client.User.Query().Where(user.EmailEqualFold(email)).Only(ctx)
}

// ExistsUserByEmail retorna true se o usuário com o e-mail informado existir
func ExistsUserByEmail(client *ent.Client, email string) (bool, error) {
	ctx := context.Background()
	count, err := client.User.Query().Where(user.EmailEqualFold(email)).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
