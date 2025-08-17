APP_NAME := orders
BIN_DIR  := bin

.PHONY: all build run test clean tidy

all: build

tidy:
	go mod tidy

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/app

run: build
	./$(BIN_DIR)/$(APP_NAME)

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)
