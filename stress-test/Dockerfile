FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o stress-test ./cmd/main.go

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/stress-test /stress-test
ENTRYPOINT ["/stress-test"]
