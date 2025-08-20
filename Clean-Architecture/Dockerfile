FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go get github.com/99designs/gqlgen/codegen/config
RUN go get github.com/99designs/gqlgen/internal/imports
RUN go get github.com/99designs/gqlgen/api
RUN go get github.com/99designs/gqlgen
RUN go run github.com/99designs/gqlgen generate
RUN go build -o server ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY ./orders.db ./orders.db
EXPOSE 8080
CMD ["./server"]
