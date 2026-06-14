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
	testseeds "github.com/matheusgosk8/book-me-server/test/seeds"
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
	u, err := testseeds.GetUserByEmail(client, email)
	if err != nil {
		t.Fatalf("Esperava usuário criado no banco, erro: %v", err)
	}
	if u.Email != email {
		t.Fatalf("E-mail do usuário criado difere: esperado %s, recebeu %s", email, u.Email)
	}

	// 6. Tentar registrar novamente com o mesmo e-mail e esperar 409
	body2, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Erro ao serializar payload para segunda requisição: %v", err)
	}
	req2 := httptest.NewRequest(http.MethodPost, "/public/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	auth.Register(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Fatalf("Esperava status %d para e-mail duplicado, recebeu %d. Body: %s", http.StatusConflict, rr2.Code, rr2.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Erro ao decodificar resposta de erro: %v", err)
	}
	msg, _ := resp["message"].(string)
	code, _ := resp["statusCode"].(float64)
	if msg != "Este e-mail já está em uso." || int(code) != http.StatusConflict {
		t.Fatalf("Resposta de erro inesperada para e-mail duplicado: %v", resp)
	}
}
