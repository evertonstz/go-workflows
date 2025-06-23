.PHONY: build run format lint

VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)

build:
	go build -ldflags "-X main.Version=$(VERSION)" .

run:
	go run -ldflags "-X main.Version=$(VERSION)" . $(ARGS)

format:
	goimports -w .

lint:
	golangci-lint run ./...