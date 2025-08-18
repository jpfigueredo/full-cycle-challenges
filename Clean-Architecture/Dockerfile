FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY ./orders.db ./orders.db
EXPOSE 8080
CMD ["./server"]
