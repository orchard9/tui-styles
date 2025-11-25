#!/bin/bash
# Install development tools for Creator API
# This script installs all necessary Go development tools

set -e

echo "Installing development tools for Creator API..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.23+ first."
    exit 1
fi

# Install Go tools
echo "Installing Go development tools..."

echo "- Installing air (hot reload)..."
go install github.com/air-verse/air@latest

echo "- Installing golangci-lint (linter)..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "- Installing goimports (import formatter)..."
go install golang.org/x/tools/cmd/goimports@latest

echo "- Installing deadcode (unused code detector)..."
go install golang.org/x/tools/cmd/deadcode@latest

echo "- Installing gocyclo (complexity analyzer)..."
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

echo "- Installing gosec (security scanner)..."
go install github.com/securego/gosec/v2/cmd/gosec@latest

echo ""
echo "All development tools installed successfully!"
echo ""
echo "Next steps:"
echo "  1. Run 'make quality' to verify everything works"
echo "  2. Run 'make pre-submit' before submitting to P4"
echo "  3. Start developing with 'make dev' for auto-reload"
