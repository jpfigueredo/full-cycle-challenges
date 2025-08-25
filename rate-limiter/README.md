# Rate Limiter Challenge

## Descrição
Este projeto implementa um rate limiter em Go para controlar requisições por IP ou token de acesso (via header API_KEY). Prioriza configs do token sobre IP, bloqueia por tempo configurado se limite excedido, e usa Redis para armazenamento (com strategy para troca fácil). Alinha com Clean Architecture (entities/usecases/repositories/adapters) e SOLID (interfaces para DIP/OCP).

### Como Funciona
- **Lógica Core**: No middleware, extrai IP e token do request. Use case verifica bloqueio, incrementa contagem na janela (ex.: 1s), e bloqueia se > max. Token sobrepõe IP (ex.: max IP=5, token=10 usa 10).
- **Storage**: Redis para contagens (INCR/EXPIRE) e bloqueios (SET/EX). Prefixos: "rate:key" para contagem, "block:key" para bloqueio.
- **Configs**: Via .env ou env vars no Docker. Ex.: MAX_REQUESTS_PER_SECOND=5 (IP), MAX_TOKEN_REQUESTS_PER_SECOND=10 (token), BLOCK_DURATION_SECONDS=300 (bloqueio 5min), WINDOW_SECONDS=1 (janela).
- **Resposta em Excesso**: HTTP 429 com mensagem "you have reached the maximum number of requests or actions allowed within a certain time frame".
- **Troca de Storage**: Implemente RateLimiterRepository (interface em repository), injete no NewRateLimiterUseCase (ex.: in-memory map com mutex para testes).

### Configuração
- Copie .env.example para .env e ajuste valores.
- Rode local: `go run cmd/server/main.go` (Redis em localhost:6379).
- Docker: `docker-compose up --build` (app na 8080, Redis interno).

### Testes
- Unitários/Integração: `go test ./... -cover` (cobertura >80%, usa miniredis para mock).
- Load Test: Use ab: `ab -n 20 -c 5 -H "API_KEY: mytoken" http://localhost:8080/ping` (ajuste max no env para ver 429).

### Exemplos
- IP Limite: Com max=5, 6ª req em 1s retorna 429, bloqueia por 5min.
- Token: Com token max=10, ignora IP limite, usa 10.
- Monitor: Use GetLimitState no use case para estado (count/blockedUntil).