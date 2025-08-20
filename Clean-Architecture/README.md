# Clean-Architecture — Orders Service

Resumo rápido

Projeto exemplo (Clean Architecture) para gerenciar orders com três interfaces:
- REST (Gin)
- GraphQL (gqlgen)
- gRPC

O repositório já contém implementações de repositório, serviços e handlers para Orders e Patients.

Requisitos
- Docker & Docker Compose
- Go 1.20+ (para desenvolvimento local)
- protoc (se desejar regenerar os arquivos gRPC)

Como rodar (modo recomendado: Docker)

1) Subir containers (API + Postgres)

```bash
docker compose up -d --build
```

A API será exposta nas seguintes portas do host:
- REST: http://localhost:8080 (endpoints: `/health`, `/orders`, `/orders/:id`, ...)
- GraphQL: http://localhost:8081 (GraphQL playground disponível em `/` e endpoint `/query`)
- gRPC: :50051

Testes rápidos
- Usando o arquivo `api.http` (VSCode REST Client ou HTTP client): execute as requests em `api.http` para criar/listar orders.
- Health: `GET http://localhost:8080/health`
- List orders REST: `GET http://localhost:8080/orders`
- List orders GraphQL: POST para `http://localhost:8081/query` com a query:

  { listOrders { id item amount patientId medication dosage status createdAt updatedAt } }

gRPC
- O proto está em `internal/api/grpc/order.proto`.
- Arquivos gerados já estão em `internal/api/grpc/orderpb`.
- Para regenerar os stubs (local):

```bash
# requer protoc e plugins instalados
make proto
```

ou

```bash
protoc --go_out=. --go-grpc_out=. -I=. internal/api/grpc/order.proto
```

GQLGen
- Schemas em `internal/api/graphql/schema.graphqls`.
- Para (re)gerar o código GraphQL localmente:

```bash
go run github.com/99designs/gqlgen generate
```

Migrations / Banco de dados
- O projeto utiliza `gorm.AutoMigrate` no startup para criar as tabelas `orders` e `patients` automaticamente.
- Se for necessária migração SQL explícita, crie arquivos em um diretório `migrations/` e aplique conforme sua ferramenta preferida (opcional).

Variáveis de ambiente relevantes (usadas no Docker Compose)
- DB_HOST (ex: `db` quando usar docker compose)
- DB_PORT (ex: 5432)
- DB_USER
- DB_PASSWORD
- DB_NAME
- SERVER_PORT (porta REST; por padrão 8080)

Arquivos úteis
- `api.http` — requests de exemplo para health, orders e patients (create/list/get/update/delete)
- `docker-compose.yml` — levanta `db` (Postgres) e `api`
- `internal/api/graphql` — servidor GraphQL e schema
- `internal/api/grpc` — proto e servidor gRPC

Comandos úteis de desenvolvimento local

```bash
# instalar dependências
make tidy

# gerar proto
make proto

# gerar gqlgen
go run github.com/99designs/gqlgen generate

# build local
go build -o bin/orders ./cmd/api

# executar local (sem Docker)
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=orders go run ./cmd/api
```

Contato/Entrega
- Documento `api.http` com exemplos para criar e listar orders.
- Porta REST: 8080
- Porta GraphQL: 8081
- Porta gRPC: 50051