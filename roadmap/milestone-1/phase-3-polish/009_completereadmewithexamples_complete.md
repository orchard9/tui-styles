## Purpose

Complete the README.md with comprehensive documentation, usage examples, and API overview. This serves as the primary entry point for developers discovering the library and provides clear guidance on using core types.

## Acceptance Criteria

- [ ] README.md expanded with complete usage examples
- [ ] Installation instructions (go get command)
- [ ] Quick start guide with Color, Position, Border examples
- [ ] API overview section with links to godoc
- [ ] Development setup instructions
- [ ] Contributing guidelines link
- [ ] License information
- [ ] Badges (Go version, Go report card, license)

## Technical Approach

**README.md Structure**:

```markdown
# TUI Styles

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/orchard9/tui-styles)](https://goreportcard.com/report/github.com/orchard9/tui-styles)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go library for terminal styling with an immutable builder pattern API, inspired by [lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- 🎨 **Colors**: Hex colors, ANSI names, 256-color codes, adaptive light/dark
- 📏 **Positioning**: Left, Center, Right (horizontal), Top, Center, Bottom (vertical)
- 🖼️ **Borders**: 8 standard border styles (Normal, Rounded, Thick, Double, etc.)
- 🔧 **Type-Safe**: Strong typing with validation at construction time
- ⚡ **Performance**: Efficient ANSI code generation with minimal allocations
- 🧪 **Tested**: >90% test coverage with comprehensive unit tests

## Installation

```bash
go get github.com/orchard9/tui-styles
```

**Requirements**: Go 1.21 or higher

## Quick Start

### Colors

```go
import "github.com/orchard9/tui-styles"

// Create colors from hex
red, _ := tuistyles.NewColor("#FF0000")
fmt.Println(red.ToANSI() + "Red text" + "\x1b[0m")

// Use ANSI color names
blue, _ := tuistyles.NewColor("blue")

// Use 256-color codes
magenta, _ := tuistyles.NewColor("201")

// Adaptive colors (light/dark terminal)
textColor, _ := tuistyles.NewAdaptiveColor("#000000", "#FFFFFF")
fmt.Println(textColor.ToColor().ToANSI() + "Adapts to terminal" + "\x1b[0m")
```

### Positions

```go
// Horizontal alignment
left := tuistyles.Left
center := tuistyles.Center
right := tuistyles.Right

// Vertical alignment
top := tuistyles.Top
bottom := tuistyles.Bottom

// Check position type
fmt.Println(center.IsHorizontal())  // true
fmt.Println(center.IsVertical())    // true (Center is both)
```

### Borders

```go
// Get standard border types
normal := tuistyles.NormalBorder()
rounded := tuistyles.RoundedBorder()
thick := tuistyles.ThickBorder()
double := tuistyles.DoubleBorder()

// Access border characters
fmt.Println(normal.TopLeft)     // "┌"
fmt.Println(rounded.TopLeft)    // "╭"
```

## API Overview

### Core Types

- **`Color`**: Represents terminal colors (hex, ANSI name, ANSI code)
  - `NewColor(string) (Color, error)` - Validate and create color
  - `ToANSI() string` - Convert to ANSI escape sequence

- **`AdaptiveColor`**: Auto-select color based on terminal background
  - `NewAdaptiveColor(light, dark string) (AdaptiveColor, error)`
  - `ToColor() Color` - Resolve to appropriate color

- **`Position`**: Enum for alignment (Left, Center, Right, Top, Bottom)
  - `String() string` - Human-readable name
  - `IsValid() bool` - Check if valid enum value
  - `IsHorizontal() bool` / `IsVertical() bool` - Type checks

- **`Border`**: Border character definitions
  - `NormalBorder()`, `RoundedBorder()`, `ThickBorder()`, etc.
  - 8 standard border styles

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/orchard9/tui-styles).

## Development

### Setup

```bash
# Clone repository
git clone https://github.com/orchard9/tui-styles.git
cd tui-styles

# Install dependencies (none for core library)
go mod download

# Install development tools
brew install golangci-lint  # macOS
# OR
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build & Test

```bash
# Build
make build

# Run tests
make test

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Lint
make lint

# Format code
make fmt
```

## Roadmap

- ✅ **Phase 1**: Core types (Color, Position, Border)
- 🚧 **Phase 2**: Style struct with immutable builder pattern
- 📋 **Phase 3**: Layout utilities (JoinHorizontal, JoinVertical, Place)
- 📋 **Phase 4**: Advanced features (inline styles, custom borders)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by [lipgloss](https://github.com/charmbracelet/lipgloss) by Charm.
```

**Files to Create/Modify**:
- `README.md` - Complete rewrite with examples
- Add badges for Go version, Go Report Card, license

**Dependencies**:
- All Phase 2 tasks must be complete (examples use core types)

## Testing Strategy

**Documentation Review**:
- [ ] All code examples are valid and compile
- [ ] Links to pkg.go.dev work (after publishing)
- [ ] Badges render correctly on GitHub
- [ ] Installation instructions are accurate
- [ ] Quick start examples demonstrate key features
- [ ] API overview matches actual implementation

**Validation**:
```bash
# Test that examples compile
go build examples/colors/main.go
go build examples/positions/main.go
go build examples/borders/main.go

# Verify README renders correctly
grip README.md  # Preview in browser
```

## Notes

**Example Programs**: Create simple example programs in `examples/` directory to demonstrate usage:
- `examples/colors/main.go` - Color examples
- `examples/positions/main.go` - Position examples
- `examples/borders/main.go` - Border examples

**Badges**: Add shields.io badges for:
- Go version requirement (1.21+)
- Go Report Card (automated code quality)
- License (MIT)
- Build status (if CI configured)

**Godoc**: Ensure all public types/functions have godoc comments. The `pkg.go.dev` link will generate automatically when module is published.

**Keep It Simple**: README should be scannable. Use clear headings, code examples, and minimal prose. Developers should understand core features in <5 minutes.

**Reference Projects**: See `reference-code/go-api/README.md` for well-structured Go README patterns.




