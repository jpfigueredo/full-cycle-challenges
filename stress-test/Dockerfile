FROM golang:1.24 AS builder
WORKDIR /app

# copy go.mod and download dependencies first (better caching)
COPY go.mod ./
RUN go mod tidy

# copy the source code
COPY . .

# build binary
RUN go build -o stress-test ./cmd/main.go

# final image
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/stress-test /stress-test
ENTRYPOINT ["/stress-test"]
