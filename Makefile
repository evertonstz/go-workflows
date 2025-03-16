.PHONY: build run format lint

build:
	go build -v -i main.go

run:
	go run .

format:
	goimports -w .

lint:
	golangci-lint run ./...