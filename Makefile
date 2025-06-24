.PHONY: build run format lint lint-fix

VERSION := $(shell git describe --tags --always)

build:
	go build -ldflags "-X main.Version=$(VERSION)" .

run:
	go run -ldflags "-X main.Version=$(VERSION)" . $(ARGS)

format:
	goimports -local github.com/evertonstz/go-workflows -w .

lint:
	golangci-lint run --config=.golangci.yml ./...

lint-fix:
	golangci-lint run --config=.golangci.yml --fix ./...