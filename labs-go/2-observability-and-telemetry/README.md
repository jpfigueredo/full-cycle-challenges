# CEP-Clima Challenge

Este repositório contém uma implementação em Go de dois serviços:

- `service-a` — responsável por receber o input (CEP) via HTTP POST e repassar para o Serviço B.
- `service-b` — responsável por orquestrar a busca da localização (via ViaCEP) e da temperatura atual (via WeatherAPI), converter as temperaturas e retornar o resultado junto com a cidade.

O projeto também tem integração com OpenTelemetry (OTEL) e Zipkin para tracing distribuído.

## Requisitos implementados

- Serviço A
  - Recebe POST `{ "cep": "29902555" }` no endpoint `/clima`.
  - Valida que o CEP é uma string com 8 dígitos; se inválido retorna HTTP 422 com mensagem `invalid zipcode`.
  - Encaminha requisição para Serviço B e propaga trace (OTEL).

- Serviço B
  - Recebe GET `/clima/:cep` com CEP válido de 8 dígitos.
  - Consulta ViaCEP (`https://viacep.com.br/`) para obter a localidade.
  - Consulta WeatherAPI (`https://www.weatherapi.com/`) para obter a temperatura em Celsius e converte para Fahrenheit e Kelvin.
  - Retorna HTTP 200 com JSON: `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }`.
  - Em caso de CEP inválido retorna HTTP 422 `invalid zipcode`.
  - Em caso de CEP não encontrado retorna HTTP 404 `can not find zipcode`.

- Observabilidade
  - Tracing distribuído entre `service-a` e `service-b` via OpenTelemetry.
  - Exportador Zipkin incluído; o serviço Zipkin está configurado no `docker-compose.yml` e fica disponível em `http://localhost:9411`.

## Requisitos técnicos

- Go (recomendado 1.24.x)
- Docker e Docker Compose
- Chave da WeatherAPI (adicione em `.env`)

## Estrutura do repositório (resumida)

- `internal/shared` — tipos compartilhados (`entity`), repositórios de integração com APIs externas e lógica de usecase comum.
- `service-a` — pequeno front que valida input e chama service-b.
- `service-b` — orquestra as integrações externas e retorna o payload final.
- `docker-compose.yml` — orquestra `zipkin`, `service-a` e `service-b`.

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com a sua chave WeatherAPI:

```
WEATHER_API_KEY=SEU_API_KEY_AQUI
```

(O `docker-compose.yml` usa esse `.env` para ambos os serviços.)

## Como rodar localmente (sem Docker)

1. Garanta que você está na raiz do projeto.
2. Sincronize o workspace Go: `go work sync`.
3. Em cada módulo atualize dependências e gere `go.sum`:
   - `cd service-a && go mod tidy`
   - `cd service-b && go mod tidy`
   - `cd internal/shared && go mod tidy` (opcional)
4. Rodar os serviços localmente:
   - `cd service-b && go run server/main.go` (porta padrão `8081`)
   - Em outra janela: `cd service-a && go run server/main.go` (porta padrão `8080`)
5. Teste via curl/postman:
   - `curl -X POST -H "Content-Type: application/json" -d '{"cep":"01001000"}' http://localhost:8080/clima`

## Como rodar com Docker Compose

1. Preencha `.env` com `WEATHER_API_KEY`.
2. Build e subir os serviços:

```
docker compose up --build
```

3. Acesse Zipkin em `http://localhost:9411` para visualizar traces.

Observações e troubleshooting para builds Docker:
- O Docker builder precisa ter acesso ao proxy do Go (proxy.golang.org). Se ocorrerem erros de TLS ou timeouts durante `go mod download` dentro do builder, tente:
  - Reexecutar o build (`docker compose build --no-cache`) — downloads transitórios podem falhar intermitentemente.
  - Forçar uso da rede do host para o passo de build (exemplo manual): `docker build --network=host -f service-b/Dockerfile .` (ou ajustar o compose para usar build options).
- Garanta que os arquivos `service-a/go.sum` e `service-b/go.sum` estejam atualizados (execute `go mod tidy` localmente antes do build). Arquivos `go.sum` faltantes causam erros do tipo "missing go.sum entry" durante o build.

## Testes automatizados

Rode os testes unitários:

```
go test ./internal/... ./service-b/... ./service-a/...
```

Ou apenas para o shared:

```
go test ./internal/shared/...
```

## Endpoints e exemplos

- Serviço A
  - POST /clima
    - Body: `{ "cep": "29902555" }`
    - 200: encaminha e retorna o JSON do service-b
    - 422: `{ "message": "invalid zipcode" }`

- Serviço B
  - GET /clima/:cep
    - 200: `{ "city": "...", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }`
    - 422: `{ "message": "invalid zipcode" }`
    - 404: `{ "message": "can not find zipcode" }`

Exemplo curl:

```
curl -X POST -H "Content-Type: application/json" -d '{"cep":"01001000"}' http://localhost:8080/clima
```

Ou diretamente no Service B:

```
curl http://localhost:8081/clima/01001000
```

## Checklist de entrega (mapeado ao enunciado)

- [x] Receber CEP no Service A e validá-lo (8 dígitos string).
- [x] Encaminhar requisição ao Service B via HTTP.
- [x] Service B consulta ViaCEP e WeatherAPI e retorna temperaturas em C/F/K e a cidade.
- [x] Retornar códigos HTTP corretos para sucesso e erros (200, 422, 404).
- [x] Implementação básica de tracing com OpenTelemetry e exporter Zipkin (Zipkin rodando via docker-compose).
- [ ] (Recomendado) Confirmar que `go.sum` de todos os módulos foi gerado e commitado antes de subir via Docker.
- [ ] (Recomendado) Remover duplicação de `service-b/entity/cep_clima.go` caso exista — a fonte canônica das entidades deve ser `internal/shared/entity`.

## O que falta / próximas ações sugeridas

1. Atualizar e commitar `go.sum` nos módulos `service-a` e `service-b` (execute `go mod tidy` em cada pasta). Isso evita erros de "missing go.sum entry" durante builds Docker.
2. Remover o arquivo duplicado `service-b/entity/cep_clima.go` (se ainda existir) para evitar confusão de tipos duplicados. O código atual já referencia `internal/shared/entity`.
3. Rodar a suíte de testes (`go test`) e validar comportamento.
4. Se houver instabilidade em builds dentro do Docker (TLS/timeout), executar build com `--network=host` ou reexecutar até completar os downloads.

Se você quiser, eu posso:
- Aplicar a remoção do arquivo duplicado `service-b/entity/cep_clima.go` para manter apenas `internal/shared/entity`.
- Rodar (guiar) os comandos `go mod tidy` e gerar instruções automáticas para commitar `go.sum`.

---

Se quiser que eu faça as alterações (remover arquivo duplicado e/ou rodar edições adicionais), diga qual ação prefere e eu aplico as mudanças.