#!/bin/bash
#
# Installation script for flux-cli
#
# Usage:
#   ./install.sh          - Install to /usr/local/bin (requires sudo)
#   ./install.sh --local  - Install to ~/.local/bin (no sudo required)

set -e

BINARY_NAME="flux-cli"
DEFAULT_INSTALL_DIR="/usr/local/bin"
LOCAL_INSTALL_DIR="$HOME/.local/bin"

# Parse arguments
INSTALL_DIR="$DEFAULT_INSTALL_DIR"
NEEDS_SUDO=true

if [[ "$1" == "--local" ]]; then
    INSTALL_DIR="$LOCAL_INSTALL_DIR"
    NEEDS_SUDO=false
    echo "Installing to local directory: $INSTALL_DIR"
else
    echo "Installing to system directory: $INSTALL_DIR (requires sudo)"
fi

# Create install directory if it doesn't exist
if [ "$NEEDS_SUDO" = true ]; then
    sudo mkdir -p "$INSTALL_DIR"
else
    mkdir -p "$INSTALL_DIR"
fi

# Build the binary
echo "Building $BINARY_NAME..."
go build -ldflags "-w -s" -o "bin/$BINARY_NAME" .

# Install the binary
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
if [ "$NEEDS_SUDO" = true ]; then
    sudo cp "bin/$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    cp "bin/$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "✓ Installation complete!"
echo ""
echo "Run '$BINARY_NAME --help' to get started"

# Check if local install dir is in PATH
if [[ "$INSTALL_DIR" == "$LOCAL_INSTALL_DIR" ]]; then
    if [[ ":$PATH:" != *":$LOCAL_INSTALL_DIR:"* ]]; then
        echo ""
        echo "⚠ Warning: $LOCAL_INSTALL_DIR is not in your PATH"
        echo "Add this line to your shell config (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
fi
