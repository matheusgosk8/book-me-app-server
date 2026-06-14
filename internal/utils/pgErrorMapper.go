package utils

import (
	"strings"
)

// PGErrorMapping define a estrutura para mapear padrões de erro para mensagens amigáveis
type PGErrorMapping struct {
	Pattern string
	Message string
}

// Mapas de erros comuns do PostgreSQL baseados em códigos de erro padrão (SQLSTATE)
var (
	// Erros de violação de restrição única (Unique Violation - 23505)
	UniqueViolations = []PGErrorMapping{
		{Pattern: "users_email_key", Message: "Este e-mail já está em uso."},
		{Pattern: "users_cpf_key", Message: "Este CPF já está cadastrado."},
		{Pattern: "users_cnpj_key", Message: "Este CNPJ já está cadastrado."},
	}

	// Erros de violação de chave estrangeira (Foreign Key Violation - 23503)
	ForeignKeyViolations = []PGErrorMapping{
		// Ent generated constraint symbols
		{Pattern: "services_categories_services", Message: "A categoria selecionada não existe."},
		{Pattern: "services_category_id_fkey", Message: "A categoria selecionada não existe."},
		// Service -> Address FK variations
		{Pattern: "services_addresses_services", Message: "O endereço selecionado não existe."},
		{Pattern: "services_address_services", Message: "O endereço selecionado não existe."},
		{Pattern: "services_address_id_fkey", Message: "O endereço selecionado não existe."},
		// Address -> User FK variations
		{Pattern: "addresses_users_addresses", Message: "O usuário associado ao endereço não foi encontrado."},
		{Pattern: "addresses_users_id_fkey", Message: "O usuário associado ao endereço não foi encontrado."},
	}
)

// LookupPGMessage tenta encontrar uma mensagem amigável baseada no conteúdo da string de erro
func LookupPGMessage(errStr string, mappings []PGErrorMapping) (string, bool) {
	for _, m := range mappings {
		if strings.Contains(errStr, m.Pattern) {
			return m.Message, true
		}
	}
	return "", false
}
