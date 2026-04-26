package main

import (
    "context"
    "log"

    "github.com/matheusgosk8/book-me-server/internal/db"
    "entgo.io/ent/dialect/sql/schema" // Adicione este import
    _ "github.com/lib/pq"
)

func main() {
    log.Println("Iniciando migração...")

    client, err := db.ConnectDB()
    if err != nil {
        log.Fatalf("❌ Erro ao conectar para migração: %v", err)
    }
    defer client.Close()

    ctx := context.Background()
    
    // Configurações para forçar o banco a aceitar o novo Schema
    if err := client.Schema.Create(
        ctx,
        schema.WithForeignKeys(true),
        schema.WithDropColumn(true), // Permite que o Ent remova colunas velhas (como service_type)
        schema.WithDropIndex(true),
    ); err != nil {
        log.Fatalf("❌ Falha ao aplicar schema: %v", err)
    }

    if err := db.SeedDatabase(client); err != nil {
        log.Printf("⚠️ Erro ao semear o banco: %v", err)
    }

    log.Println("✅ Banco de dados atualizado e semeado com sucesso!")
}