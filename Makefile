.PHONY: build run format lint lint-fix test test-verbose test-cover test-cover-html test-race test-clean help

VERSION := $(shell git describe --tags --always)

# Default target
help: ## Show available commands
	@echo "go-workflows - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -ldflags "-X main.Version=$(VERSION)" .

run: ## Run the application
	go run -ldflags "-X main.Version=$(VERSION)" . $(ARGS)

# Testing commands (The Go Way)
test: ## Run all tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -cover ./...

test-cover-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-race: ## Run tests with race detection
	go test -race ./...

test-clean: ## Clean test cache and coverage files
	go clean -testcache
	rm -f coverage.out coverage.html

# Code quality
format: ## Format Go code
	goimports -local github.com/evertonstz/go-workflows -w .

lint: ## Run golangci-lint
	golangci-lint run --config=.golangci.yml ./...

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --config=.golangci.yml --fix ./...

# Development workflows
dev: format test ## Format and test (quick dev workflow)
ci: format lint test-race test-cover ## Full CI pipeline