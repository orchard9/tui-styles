# /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker Development Makefile

# Variables
DOCKER_COMPOSE := docker compose -f docker-compose.yml
AIR := air
MIGRATE := orchard9-migrate

# Include .env if it exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

.DEFAULT_GOAL := help

#===============================================================================
# HELP
#===============================================================================

.PHONY: help
help: ## Show available commands
	@echo "/var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

#===============================================================================
# MAIN DEVELOPMENT WORKFLOW
#===============================================================================

.PHONY: start
start: ## Start complete development environment
	@echo "Starting complete development environment..."
	@$(MAKE) up
	@sleep 3
	@echo "Development environment started!"
	@echo "  - Use 'make dev' to start hot reload"
	@echo "  - Use 'make logs' to view logs"
	@echo "  - Use 'make stop' to stop everything"

.PHONY: stop
stop: ## Stop everything
	@echo "Stopping development environment..."
	@$(MAKE) down
	@pkill -f air 2>/dev/null || true
	@echo "Development environment stopped!"

.PHONY: restart
restart: ## Restart everything cleanly
	@echo "Restarting development environment..."
	@$(MAKE) stop
	@sleep 2
	@$(MAKE) start

.PHONY: status
status: ## Check development environment status
	@echo "Development Environment Status:"
	@echo "==============================="
	@$(DOCKER_COMPOSE) ps

#===============================================================================
# DOCKER SERVICES
#===============================================================================

.PHONY: up
up: ## Start docker services only
	@echo "Starting docker services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Docker services started!"

.PHONY: down
down: ## Stop docker services only
	@echo "Stopping docker services..."
	@$(DOCKER_COMPOSE) down
	@echo "Docker services stopped!"

.PHONY: logs
logs: ## Show docker logs
	@$(DOCKER_COMPOSE) logs -f --tail=100

#===============================================================================
# DEVELOPMENT
#===============================================================================

.PHONY: dev
dev: ## Start Air hot reload in background
	@echo "Starting Air hot reload..."
	@mkdir -p tmp
	@$(AIR) -c .air.toml

.PHONY: dev-bg
dev-bg: ## Start Air hot reload in background
	@echo "Starting Air hot reload in background..."
	@mkdir -p tmp
	@nohup $(AIR) -c .air.toml > tmp/air.log 2>&1 &
	@echo "Air started in background. Use 'make air-logs' to view logs."

.PHONY: air-logs
air-logs: ## Show Air logs
	@if [ -f tmp/air.log ]; then \
		tail -f tmp/air.log; \
	else \
		echo "No Air logs found. Start Air with 'make dev-bg' first."; \
	fi

#===============================================================================
# RESET AND CLEANUP
#===============================================================================

.PHONY: reset
reset: ## Nuclear reset - stop everything and reset docker volumes
	@echo "Nuclear reset - stopping everything and clearing data..."
	@pkill -f air 2>/dev/null || true
	@$(DOCKER_COMPOSE) down -v
	@rm -rf tmp/ bin/
	@echo "Starting fresh environment..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Reset complete!"

.PHONY: clean
clean: ## Clean build artifacts and logs
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ tmp/ coverage.out coverage.html
	@echo "Clean complete!"

#===============================================================================
# CI/BUILD/TEST
#===============================================================================

.PHONY: ci
ci: ## Run full CI pipeline (NEVER MODIFY THIS)
	@echo "Running CI pipeline..."
	@go mod tidy
	@go fmt ./...
	@go vet ./...
	@golangci-lint run --allow-parallel-runners
	@go test ./... -race -coverprofile=coverage.out
	@go build -o bin/server ./cmd/server
	@echo "CI completed"

.PHONY: test
test: ## Run tests
	@go test ./... -race

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@go test ./... -race -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: build
build: ## Build binaries
	@mkdir -p bin
	@go build -o bin/server ./cmd/server

.PHONY: fmt
fmt: ## Format code
	@go fmt ./...

.PHONY: lint
lint: ## Run linters
	@go vet ./...
	@golangci-lint run --allow-parallel-runners

.PHONY: complexity
complexity: ## Check code complexity
	@golangci-lint run --allow-parallel-runners --enable-only=gocyclo,gocognit --timeout=60s

.PHONY: file-length
file-length: ## Check file and function length
	@golangci-lint run --allow-parallel-runners --enable-only=funlen,lll --timeout=60s

.PHONY: deadcode
deadcode: ## Detect dead/unused code
	@golangci-lint run --allow-parallel-runners --enable-only=unused,ineffassign,unparam --timeout=60s

.PHONY: circular
circular: ## Check for circular dependencies
	@echo "Checking for circular dependencies..."
	@go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./... | \
		awk '{print $$1}' | sort | uniq -c | sort -rn | head -20 || true
	@echo "Checking for import cycles..."
	@go list -f '{{if .Error}}{{.ImportPath}}: {{.Error.Err}}{{end}}' ./... | grep -i cycle || echo "No import cycles detected"

.PHONY: setup
setup: ## Install dependencies and tools
	@echo "Installing dependencies and tools..."
	@go mod download
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Setup complete!"

#===============================================================================
# DOCKER
#===============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@docker build -t /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker:latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@docker run -p 8080:8080 --env-file .env /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker:latest

#===============================================================================
# UTILITIES
#===============================================================================

.PHONY: deps
deps: ## Download dependencies
	@go mod download && go mod tidy

#===============================================================================
# ALIASES
#===============================================================================

.PHONY: dev-up
dev-up: up ## Alias for up

.PHONY: dev-down
dev-down: down ## Alias for down
