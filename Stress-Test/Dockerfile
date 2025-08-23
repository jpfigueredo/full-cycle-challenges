FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o stress-test ./cmd

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/stress-test .
ENTRYPOINT ["./stress-test"]
