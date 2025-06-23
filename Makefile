.PHONY: build run format lint

VERSION := $(shell git describe --tags --always)

build:
	go build -ldflags "-X main.Version=$(VERSION)" .

run:
	go run -ldflags "-X main.Version=$(VERSION)" . $(ARGS)

format:
	goimports -w .

lint:
	golangci-lint run ./...