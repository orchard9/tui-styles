# Tui Styles CLI Makefile

BINARY_NAME=tui-styles-cli
BUILD_DIR=bin
LOCAL_INSTALL_DIR=$(HOME)/.local/bin
SYSTEM_INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-w -s"

.PHONY: all build clean test install install-system uninstall uninstall-system fmt lint help

## all: Build the CLI tool
all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_NAME)"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	@echo "✓ Cleaned"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests passed"

## coverage: Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out
	@echo "✓ Coverage report generated"

## install: Install the binary to local directory (no sudo required)
install: build
	@echo "Installing $(BINARY_NAME) to $(LOCAL_INSTALL_DIR)..."
	@mkdir -p $(LOCAL_INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(LOCAL_INSTALL_DIR)/
	@chmod +x $(LOCAL_INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed $(LOCAL_INSTALL_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "Run '$(BINARY_NAME) --help' to get started"
	@if echo $$PATH | grep -q "$(LOCAL_INSTALL_DIR)"; then \
		echo "✓ $(LOCAL_INSTALL_DIR) is in your PATH"; \
	else \
		echo ""; \
		echo "⚠  $(LOCAL_INSTALL_DIR) is not in your PATH"; \
		echo "   Add this line to your shell config (~/.bashrc, ~/.zshrc, etc.):"; \
		echo "   export PATH=\"\$$HOME/.local/bin:\$$PATH\""; \
		echo ""; \
		echo "   Or use the full path: $(LOCAL_INSTALL_DIR)/$(BINARY_NAME)"; \
	fi

## install-system: Install the binary to system directory (requires sudo)
install-system: build
	@echo "Installing $(BINARY_NAME) to $(SYSTEM_INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(SYSTEM_INSTALL_DIR)/
	@sudo chmod +x $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME)"
	@echo "  Run '$(BINARY_NAME) --help' to get started"

## uninstall: Remove the binary from local directory
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(LOCAL_INSTALL_DIR)..."
	@rm -f $(LOCAL_INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled"

## uninstall-system: Remove the binary from system directory (requires sudo)
uninstall-system:
	@echo "Uninstalling $(BINARY_NAME) from $(SYSTEM_INSTALL_DIR)..."
	@sudo rm -f $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Code formatted"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "✓ Linting passed"; \
	else \
		echo "⚠ golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
	fi

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies downloaded"

## verify: Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify
	@echo "✓ Dependencies verified"

## help: Show this help message
help:
	@echo "Tui Styles CLI - Build Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'
