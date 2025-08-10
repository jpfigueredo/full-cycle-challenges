# cep-race

Small CLI in Go that races two Brazilian CEP APIs (BrasilAPI and ViaCEP) and accepts the fastest successful response.

## Requirements

- Go 1.20+

## Files
- `main.go` — main program entrypoint
- `fetcher.go` — API fetchers implementations
- `race.go` — race logic service
- `model.go` — shared structs
- `go.mod` — module definition

## Behavior / Requirements implemented
- Send two requests **concurrently** to:
  - `https://brasilapi.com.br/api/cep/v1/<cep>`
  - `http://viacep.com.br/ws/<cep>/json/`
- Accept the **first successful** response and discard the slower request.
- Global timeout: **1 second**. If neither API returns successfully within 1 second, print a timeout error.
- Print the address fields and which API supplied them.

## Usage
```bash
go build
go run *.go 01153000
```

Example output:
```bash
Source: viacep
CEP: 01153-000
Street: Rua Vitorino Carmilo
Neighborhood: Barra Funda
City: São Paulo
State: SP
```

Design notes (short)
- Two or more goroutines run the HTTP requests concurrently and send results to a channel.
- A parent context.WithTimeout(..., 1*time.Second) enforces the overall timeout.
- Per-request context.WithCancel(parentCtx) allows cancellation of the slower request when a winner arrives.
- The program picks the first successful response. If both APIs error, the first error received is reported. If none respond within 1s, a timeout is printed.
- Interfaces allow easy extension to more fetchers or mocks for testing.
- Separation of concerns makes code easier to maintain and test.

## Next steps / improvements
- Normalize CEP input (strip non-digit characters).
- Add retries with exponential backoff for transient network errors (careful with the 1s deadline).
- Add unit tests using httptest.Server to simulate providers with configurable latencies and responses.
- Wrap the CLI with proper argument parsing and validation.