# Book-me Server

API em Go parte do monorepo Book-me.

Pré-requisitos
- Go instalado (verifique com `go version`)
- `air` (opcional, para hot-reload)

Instalação de dependências
```powershell
cd C:\Users\mtsgo\Desktop\book-me\book-me-server
go mod download
# ou atualizar/resolver módulos
go mod tidy
```

Instalar `air` (opcional, recomendado para desenvolvimento)
```powershell
go install github.com/cosmtrek/air@latest
# garantir GOPATH/bin no PATH (PowerShell):
$env:PATH += ";" + (go env GOPATH) + "\bin"
```

Executar a API

- Com hot-reload (recomendado):
```powershell
air -c .air.toml
```
Isso usa a configuração em `.air.toml` e sobe a API automaticamente na porta 8000.

- Sem hot-reload (rápido):
```powershell
go run ./cmd/main.go
```

A API por padrão escuta em http://localhost:8000

Gerar e migrar schema (Ent)

1. Crie a base do schema com o comando (substitua `NomeDoSchema`):
```powershell
go run -mod=mod entgo.io/ent/cmd/ent new NomeDoSchema
# Exemplo: go run -mod=mod entgo.io/ent/cmd/ent new Address
```

2. Abra e popule o arquivo gerado em `ent/schema/NomeDoSchema.go` com os campos desejados.

3. Rode a migração manualmente para aplicar o schema no banco:
```powershell
# (opcional) defina a variável de ambiente DATABASE_URL se necessário
$env:DATABASE_URL="postgres://postgres:docker@localhost:5433/book-me-be?sslmode=disable"
go run ./cmd/migrate/main.go
```

Observações:
- Evite manter `client.Schema.Create(...)` sendo chamado automaticamente no startup do servidor se quiser controlar migrações manualmente.
- Para desenvolvimento com hot-reload você pode criar uma configuração separada do `air` que execute o comando de migração, mas o método direto acima é suficiente na maioria dos casos.

**Ac**
-[] Automatizar fluxo de migrações com o Air

Contribuições
- Abra issues ou envie pull requests para este repositório.

-Mudanças aleatórias
