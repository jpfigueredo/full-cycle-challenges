FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app/service-b
# copy module files for service-b
COPY service-b/go.mod service-b/go.sum ./
# copy shared module to the parent relative path expected by replace (../internal/shared)
COPY internal/shared ../internal/shared
RUN go mod download
COPY service-b/. .
RUN go build -o main server/main.go

FROM alpine:latest
WORKDIR /app/service-b
COPY --from=builder /app/service-b/main .
CMD ["./main"]
