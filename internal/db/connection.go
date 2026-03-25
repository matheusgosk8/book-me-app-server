package db

import (
	"context"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/matheusgosk8/book-me-server/ent"
)

var Client *ent.Client

func ConnectDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:docker@localhost:5433/book-me-be?sslmode=disable"
	}

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Errorf("Falha ao abrir conexão Ent: %v", err)
		return err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("Falha ao criar schema: %v", err)
		return err
	}

	log.Info("Conexão e schema Ent criados com sucesso!")
	Client = client
	return nil
}
