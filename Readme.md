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

Contribuições
- Abra issues ou envie pull requests para este repositório.
