package utils

import (
	"context"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/matheusgosk8/book-me-server/ent"
)

func RunMigration() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:docker@localhost:5433/book-me-be?sslmode=disable"
	}

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Falha ao abrir conexão Ent: %v", err)
	}

	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Falha ao rodar migration: %v", err)
	}

	log.Println("Migração Ent executada com sucesso!")
}
