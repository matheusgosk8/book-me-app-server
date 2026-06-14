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

go generate ./ent  
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

### Mapeamento, FKs e erro Postgres 23503

Se ao criar/atualizar um registro você vir um erro como:

```json
{
	"message": "ent: constraint failed: pq: insert or update on table \"services\" violates foreign key constraint \"services_addresses_services\" (23503)",
	"statusCode": 400
}
```

Isso é uma violação de chave estrangeira (Postgres error 23503). Causas e checagens rápidas:

- **Entidade referenciada não existe**: o `address_id` que você está enviando não está presente na tabela `addresses`. Verifique com:

```sql
SELECT id FROM addresses WHERE id = '<ADDRESS_ID>';
```

- **Formato inválido (UUID)**: confirme que o valor é um UUID válido.
- **Ordem de operações**: se você está criando `address` e `service` na mesma operação, garanta que o `address` foi persistido antes de referenciá-lo (ou crie via edge dentro da mesma transação usando o objeto criado).
- **Uso correto do Ent**: prefira usar as helpers do Ent para ligar entidades por edge ao invés de manipular a FK manualmente quando possível.

Exemplos de criação usando Ent (patterns recomendados):

1) Criar `Service` apontando para um `Address` já existente (por id):

```go
// supondo addressID é uuid.UUID
svc, err := client.Service.
		Create().
		SetTitle("Limpeza").
		SetPriceBase(100.0).
		SetAddressID(addressID).
		Save(ctx)

// ou usando o objeto Address já carregado:
svc, err := client.Service.
		Create().
		SetTitle("Limpeza").
		SetPriceBase(100.0).
		SetAddress(address).
		Save(ctx)
```

2) Criar `Address` e `Service` em sequência (garantir persistência do endereço primeiro):

```go
addr, err := client.Address.Create().SetStreet("Rua").Save(ctx)
svc, err := client.Service.Create().SetTitle("X").SetAddress(addr).Save(ctx)
```

3) Se quiser criar ambos numa transação, use a mesma transação/ent.Client para garantir que o referenciado exista no commit.

Logs e debugging

- Faça log do payload recebido (especialmente o `address_id`) antes de persistir — já existe um log no handler de criação de serviço (`log.Printf("Payload do serviço: %+v", svcReq)`).
- Se o erro persistir, rode um `SELECT` no banco para validar a existência do id e o nome da constraint para mapear qual FK está falhando.

Se quiser, posso:

- rodar `go generate ./ent` e a migração aqui (se fornecer `DATABASE_URL`/acesso),
- ou inspecionar um exemplo de payload que causa o erro e ajudar a consertar o mapeamento no handler/mapper.

**Ac**
-[] Automatizar fluxo de migrações com o Air

Contribuições
- Abra issues ou envie pull requests para este repositório.

-Mudanças aleatórias
- Ci teste 1