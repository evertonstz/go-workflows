.PHONY: build run format format-check lint lint-fix test test-verbose test-integration test-integration-update test-cover test-cover-html test-cover-summary test-race test-clean help

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

test-integration: ## Run Bubble Tea integration tests with teatest
	go test -v -run "TestApp" ./

test-integration-update: ## Update golden files for integration tests
	go test -v -run "TestApp_FullOutput" ./ -update

test-cover: ## Run tests with coverage
	go test -cover ./...

test-cover-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-cover-summary: ## Show coverage summary
	go test -coverprofile=coverage.out ./...
	@echo "Coverage Summary:"
	@go tool cover -func=coverage.out

test-race: ## Run tests with race detection
	go test -race ./...

test-clean: ## Clean test cache and coverage files
	go clean -testcache
	rm -f coverage.out coverage.html

# Code quality
format: ## Format Go code
	go tool goimports -local github.com/evertonstz/go-workflows -w .

format-check: ## Check if code is formatted correctly
	@echo "Checking code formatting..."
	@if [ "$$(go tool goimports -local github.com/evertonstz/go-workflows -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted correctly. Files that need formatting:"; \
		go tool goimports -local github.com/evertonstz/go-workflows -l .; \
		echo "Please run 'make format' to fix formatting."; \
		exit 1; \
	fi
	@echo "Code formatting is correct ✓"

lint: ## Run golangci-lint
	golangci-lint run --config=.golangci.yml ./...

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --config=.golangci.yml --fix ./...

vuln-check: ## Check for vulnerabilities
	go tool govulncheck ./...

# Development workflows
dev: format test ## Format and test (quick dev workflow)
ci: format-check lint vuln-check test-race test-cover ## Full CI pipeline