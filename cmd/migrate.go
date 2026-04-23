// Arquivo mantido para compatibilidade; não é um entrypoint.
package main

import "github.com/matheusgosk8/book-me-server/internal/utils"

// Helper para executar migração a partir deste pacote sem ser entrypoint.
func RunMigrationCmd() {
	utils.RunMigration()
}
