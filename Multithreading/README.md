# Multithreading — Race entre APIs CEP

Projeto CLI Go que busca um CEP simultaneamente em duas APIs públicas brasileiras, aceitando o resultado da que responder primeiro com sucesso.

## Requirements

- Go 1.20+

## Files

- `main.go` — main program entrypoint
- `fetcher.go` — API fetchers implementations
- `race.go` — race logic service
- `model.go` — shared structs
- `go.mod` — module definition

## Principais características

- Concorre simultaneamente entre:
  - <https://brasilapi.com.br/api/cep/v1/><cep>
  - <http://viacep.com.br/ws/><cep>/json/
- Timeout global de 1 segundo.
- Aceita primeiro resultado bem-sucedido, cancela requisições concorrentes restantes.
- Exibe resultado e qual API forneceu os dados.

## Estrutura do projeto

```bash
.
├── main.go        # Entrada CLI, orquestra chamada race
├── fetcher.go     # Implementação fetchers das APIs
├── race.go        # Lógica de concorrência e cancelamento
├── model.go       # Structs compartilhados
├── go.mod
└── README.md
```

## Fluxo do sistema

```
CLI -> race service -> [API1, API2 concorrentes]
                      |
                      -> primeiro sucesso -> cancela concorrentes -> retorna resposta
```

## Uso típico

```bash
go build
go run *.go 01153000
```

### Exemplo de saída

```bash
Source: viacep
CEP: 01153-000
Street: Rua Vitorino Carmilo
Neighborhood: Barra Funda
City: São Paulo
State: SP
```

## Design notes

- Uso extensivo de goroutines e canais para concorrência.
- context.WithTimeout controla timeout global.
- Cancelamento de requisições lentas com context.WithCancel.
- Interface para fetchers permite testes e extensibilidade.
- Separação clara entre lógica de fetch, corrida e CLI.

## Próximos passos recomendados para ambos projetos

- Adicionar testes unitários e de integração (mocks para APIs e DB).
- Melhorar validação e tratamento de erros.
- Implementar logging estruturado.
- Expandir para casos de uso reais (ex.: cache, retries, backoff).
