package testbootstrap

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
	"github.com/matheusgosk8/book-me-server/ent/address"
	"github.com/matheusgosk8/book-me-server/ent/user"
	"github.com/matheusgosk8/book-me-server/internal/db"
	log "github.com/sirupsen/logrus"
)

// SetupTestDB inicializa a conexão com o banco de dados para os testes
func SetupTestDB() *ent.Client {
	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados de teste: %v", err)
	}
	return client
}

// CleanTestUser remove o usuário de teste e seus endereços do banco de dados
func CleanTestUser(client *ent.Client, email string) {
	ctx := context.Background()

	// 1. Encontrar o usuário pelo e-mail
	u, err := client.User.Query().Where(user.EmailEqualFold(email)).Only(ctx)
	if err == nil {
		// 2. Remover endereços vinculados
		_, err = client.Address.Delete().Where(address.HasUserWith(user.ID(u.ID))).Exec(ctx)
		if err != nil {
			log.Warnf("Erro ao limpar endereços do usuário %s: %v", email, err)
		}

		// 3. Remover o usuário
		err = client.User.DeleteOne(u).Exec(ctx)
		if err != nil {
			log.Warnf("Erro ao limpar usuário %s: %v", email, err)
		} else {
			log.Infof("Usuário de teste %s e seus endereços foram removidos.", email)
		}
	}
}
