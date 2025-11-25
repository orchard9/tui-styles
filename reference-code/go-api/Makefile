# Creator API Makefile
# Common development tasks for the Masquerade Creator API service

.PHONY: help build run test lint fmt fmt-check deadcode complexity security quality pre-submit install-hooks clean dev docker-build install-tools

# Variables
BINARY_NAME=creator-api
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/server/main.go
COVERAGE_FILE=coverage.out

# Default target: show help
help: ## Display this help message
	@echo "Creator API - Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary created at $(BINARY_PATH)"

run: ## Run the application
	@echo "Starting Creator API server..."
	@go run $(MAIN_PATH)

dev: ## Run with auto-reload (requires air: make install-tools)
	@echo "Starting development server with auto-reload..."
	@air

test: ## Run all tests with coverage
	@echo "Running tests..."
	@go test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@echo "Coverage report: $(COVERAGE_FILE)"
	@echo "View HTML coverage: go tool cover -html=$(COVERAGE_FILE)"

test-coverage: test ## Run tests and open coverage report in browser
	@go tool cover -html=$(COVERAGE_FILE)

lint: ## Run golangci-lint (requires golangci-lint: make install-tools)
	@echo "Running linters..."
	@golangci-lint run ./...

fmt: ## Format code with gofmt and goimports
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .
	@echo "Code formatted successfully"

fmt-check: ## Check if code is formatted (without modifying)
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files are not formatted:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "All files are properly formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

deadcode: ## Detect unused code (requires deadcode: make install-tools)
	@echo "Detecting dead code..."
	@deadcode -test ./...

complexity: ## Report code complexity (requires gocyclo: make install-tools)
	@echo "Analyzing code complexity..."
	@gocyclo -over 10 .

security: ## Run security checks (requires gosec: make install-tools)
	@echo "Running security scans..."
	@gosec -quiet ./...

tidy: ## Tidy and verify go.mod
	@echo "Tidying go.mod..."
	@go mod tidy
	@go mod verify

clean: ## Remove build artifacts and cache
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f $(COVERAGE_FILE)
	@go clean -cache -testcache
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t masquerade-creator-api:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -p 34075:34075 --env-file .env masquerade-creator-api:latest

install-tools: ## Install development tools (air, golangci-lint, goimports, deadcode, gocyclo, gosec)
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/tools/cmd/deadcode@latest
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Tools installed successfully"

quality: ## Run all quality checks (lint, fmt-check, vet, test, deadcode, complexity, security)
	@echo "Running comprehensive quality checks..."
	@$(MAKE) fmt-check
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) deadcode
	@$(MAKE) complexity
	@$(MAKE) security
	@echo "All quality checks passed!"

pre-submit: ## Run pre-submit validation for P4 (Perforce) workflows
	@echo "Running pre-submit validation..."
	@bash scripts/pre-submit.sh

install-hooks: ## Install pre-commit hooks (for git workflows)
	@echo "Installing pre-commit hooks..."
	@pip install pre-commit || echo "Failed to install pre-commit. Install manually: pip install pre-commit"
	@pre-commit install
	@pre-commit install --hook-type commit-msg
	@echo "Pre-commit hooks installed. Run 'pre-commit run --all-files' to test."

all: clean fmt quality build ## Run all quality checks and build

.DEFAULT_GOAL := help
