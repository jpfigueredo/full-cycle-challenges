FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app/service-a
# copy module files for service-a
COPY service-a/go.mod service-a/go.sum ./
# copy shared module to the parent relative path expected by replace (../internal/shared)
COPY internal/shared ../internal/shared
RUN go mod download
COPY service-a/. .
RUN go build -o main server/main.go

FROM alpine:latest
WORKDIR /app/service-a
COPY --from=builder /app/service-a/main .
CMD ["./main"]
