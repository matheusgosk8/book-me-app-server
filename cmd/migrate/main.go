package main

import (
	"context"
	"log"

	"github.com/matheusgosk8/book-me-server/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Iniciando migração...")

	// 1. Abre a conexão usando a função existente
	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Erro ao conectar para migração: %v", err)
	}
	defer client.Close()

	// 2. Roda as migrações do Ent
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("❌ Falha ao aplicar schema: %v", err)
	}

	// 3. Se quiser rodar o Seed (popular dados iniciais), coloque aqui:
	if err := db.SeedDatabase(client); err != nil {
		log.Printf("⚠️ Erro ao semear o banco: %v", err)
	}

	log.Println("✅ Banco de dados atualizado e semeado com sucesso!")
}
