//go:build e2e
// +build e2e

package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	testbootstrap "github.com/matheusgosk8/book-me-server/cmd/test"
	"github.com/matheusgosk8/book-me-server/internal/handlers/auth"
	"github.com/matheusgosk8/book-me-server/test/payloads"
)

func TestRegisterHandler_E2E(t *testing.T) {
	// 1. Setup do Banco e Payload
	client := testbootstrap.SetupTestDB()
	payload := payloads.GetRegisterPayload()
	email := "test@e2e.com"

	// Limpeza preventiva antes e depois do teste
	testbootstrap.CleanTestUser(client, email)
	defer testbootstrap.CleanTestUser(client, email)

	// 2. Preparar requisição HTTP
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Erro ao serializar payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/public/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// 3. Executar o Handler
	auth.Register(rr, req)

	// 4. Validações de Resposta
	if rr.Code != http.StatusCreated {
		t.Errorf("Status incorreto: esperado %d, recebeu %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	// 5. Validar se o usuário foi realmente criado no banco de dados
	// (Opcional, mas recomendado para um teste E2E real)
}
