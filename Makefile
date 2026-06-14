.PHONY: test-unit test-e2e

# Roda apenas testes unitários (ignora os que têm tag e2e)
test-unit:
	go test -v ./internal/...

# Roda todos os testes que possuem a tag //go:build e2e
test-e2e:
	go test -v -tags=e2e ./internal/handlers/...

# Roda absolutamente tudo
test-all:
	go test -v -tags=e2e ./...
