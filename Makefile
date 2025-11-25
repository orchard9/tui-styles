.PHONY: build test lint fmt clean coverage help

# Default target
help:
	@echo "TUI Styles - Development Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Build all packages"
	@echo "  test     - Run all tests with race detection"
	@echo "  lint     - Run golangci-lint"
	@echo "  fmt      - Format code with gofmt and goimports"
	@echo "  coverage - Generate test coverage report"
	@echo "  clean    - Clean build artifacts and coverage files"

build:
	@echo "Building all packages..."
	go build ./...

test:
	@echo "Running tests with race detection..."
	go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "goimports not installed, skipping (install with: go install golang.org/x/tools/cmd/goimports@latest)"

coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f coverage.out coverage.html
