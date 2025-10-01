# CEP-Clima Challenge

Este projeto é um sistema em Go que recebe um CEP brasileiro válido (8 dígitos), consulta a localização via ViaCEP, obtém a temperatura atual via WeatherAPI e retorna as temperaturas em Celsius, Fahrenheit e Kelvin. O sistema é containerizado com Docker e deployado no Google Cloud Run.

## Requisitos
- Go 1.24.6 ou superior (para desenvolvimento local).
- Docker e Docker Compose (para testes containerizados).
- Chave API válida da WeatherAPI (gratuita em weatherapi.com) – adicione em `.env`.
- Conta GCP free tier para validar o deploy (se necessário).

## Configuração Local
1. Clone o repositório e navegue para a pasta.
2. Crie `.env` baseado em `.env.example` e preencha com o seu `WEATHER_API_KEY`.
3. Rode `go mod tidy` para dependências.

## Como Rodar Localmente
- Sem Docker: `go run cmd/server/main.go`. O server roda em `localhost:8080`.
- Com Docker: `docker-compose up --build`. O app fica disponível em `localhost:8080`.

## Testes Automatizados
Rode `go test ./internal/...` para validar units (validação CEP, conversões, mocks de APIs). Todos devem passar para confirmação de lógica.

## Como Testar e Validar (Para Avaliadores)
Para aprovar, teste os cenários abaixo. Use curl ou browser. O endpoint é `/clima/:cep`.

### Cenário 1: Sucesso (HTTP 200)
- Request: `curl localhost:8080/clima/01001000` (CEP válido de São Paulo).
- Esperado: JSON como `{"temp_C":25.3,"temp_F":77.5,"temp_K":298.4}` (valores reais variam; verifique se são números arredondados, com conversões corretas: F = C*1.8+32, K = C+273).
- Validação: Confirme que a localização é consultada corretamente e temperaturas convertidas.

### Cenário 2: CEP Inválido (Formato Incorreto, HTTP 422)
- Request: `curl localhost:8080/clima/1234567` (menos de 8 dígitos) ou `curl localhost:8080/clima/123456789` (mais de 8) ou `curl localhost:8080/clima/abc-defgh` (não numérico).
- Esperado: `{"message":"invalid zipcode"}`.
- Validação: A limpeza remove não-dígitos e checa exatamente 8 dígitos.

### Cenário 3: CEP Não Encontrado (HTTP 404)
- Request: `curl localhost:8080/clima/00000000` (CEP inexistente).
- Esperado: `{"message":"can not find zipcode"}`.
- Validação: ViaCEP retorna vazio ou erro, mapeado para 404.

### Cenário 4: Erro Interno (HTTP 500)
- Simule falha: Use chave API inválida ou CEP que leva a erro na WeatherAPI (ex.: localização inválida).
- Esperado: `{"message":"internal error"}`.
- Validação: Erros de rede/API são capturados sem crash.

### Validação com Docker
- Rode `docker-compose up` e teste os curls acima – deve replicar o local.
- Inspecione logs: `docker-compose logs` para debugs (ex.: CEP limpo, localização, clima).

## Deploy no Google Cloud Run
O app está deployado em https://cep-clima-797444264306.us-central1.run.app/clima/:cep (ex: 01001000). Teste os mesmos cenários substituindo localhost pelo URL.

- Para recriar: Siga os passos no histórico (build/push/deploy com gcloud).
- Validação: Confirme que o free tier não gera custos (monitore no console GCP).

Se todos cenários passarem e testes rodarem sem falhas, o projeto atende aos requisitos!