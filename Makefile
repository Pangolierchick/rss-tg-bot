SERVICE_NAME := rss-tg-bot
BIN_DIR := ./bin

all: build

create_bin_dir:
	mkdir -p $(BIN_DIR)

build: create_bin_dir
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)/main.go

run:
	go run ./cmd/$(SERVICE_NAME)/main.go
