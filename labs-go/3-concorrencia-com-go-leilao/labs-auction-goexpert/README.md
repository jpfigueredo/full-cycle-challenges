# labs-auction-goexpert

Projeto de leilões desenvolvido em Go — inclui funcionalidade de fechamento automático de leilões após um intervalo configurável via variável de ambiente `AUCTION_INTERVAL`.

## Como rodar (com Docker Compose)

1. Subir containers (MongoDB + app):

   docker compose up -d

   Observação: o `docker-compose.yml` já carrega variáveis de `cmd/auction/.env` e também replica as variáveis principais para os containers via `environment` (MONGODB_URL, MONGODB_DB, AUCTION_INTERVAL, MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD).

2. Verificar containers:

   docker ps -a

3. Executar os testes de integração (opcional):

- Opção A — rodando os testes a partir do host (recomendado quando o Mongo expõe a porta 27017 para o host):

  Antes de executar os testes, exporte as variáveis apontando para `localhost` (caso queira usar um banco de teste separado):

  ```bash
  export MONGODB_URL='mongodb://admin:admin@localhost:27017/auctions_test?authSource=admin'
  export MONGODB_DB='auctions_test'
  export AUCTION_INTERVAL='2s'
  go test ./internal/infra/database/auction -run TestAutoCloseAuctions -v
  ```

- Opção B — rodando os testes dentro do container (usa os hostnames e variáveis do compose):

  ```bash
  docker compose exec app sh -c 'export MONGODB_DB=auctions_test; export AUCTION_INTERVAL=2s; go test ./internal/infra/database/auction -run TestAutoCloseAuctions -v'
  ```

## Notas sobre problemas comuns

- Se o teste falhar com erro de conexão ao Mongo (`lookup mongodb: no such host`), você está executando o teste no host, mas o `MONGODB_URL` no `.env` utiliza o hostname `mongodb` que só resolve dentro da rede do Docker. Use a Opção A acima (alterando o host para `localhost`) ou rode os testes dentro do container (Opção B).

- Se houver erro de autenticação, confirme as credenciais em `cmd/auction/.env` (`MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD`) e que o `MONGODB_URL` inclui essas credenciais (ex.: `mongodb://admin:admin@localhost:27017/auctions_test?authSource=admin`).

- `AUCTION_INTERVAL` é lido pela aplicação na inicialização. Para testes rápidos use `2s`.

## Exemplos rápidos

- Executar servidor localmente (sem Docker): exporte `MONGODB_URL` apontando para um Mongo em execução e rode `go run cmd/auction/main.go`.

## Perguntas

Se quiser, posso:
- adicionar um target Makefile para os comandos mais comuns;
- criar um `docker-compose.override.yml` para desenvolvimento com valores diferentes (ex.: DB de teste);
- adicionar scripts de CI para executar os testes de integração em runners com Mongo disponível.
