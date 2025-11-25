# TUI Styles

A Go library for styling terminal output with a fluent/builder pattern API.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/orchard9/tui-styles)](https://goreportcard.com/report/github.com/orchard9/tui-styles)

## Features

- 🎨 **Text Styling**: Bold, italic, underline, strikethrough, and more
- 🌈 **Colors**: Hex, ANSI names, 256-color support, adaptive colors
- 📦 **Borders**: 8 predefined border styles with Unicode box drawing
- 📏 **Layout**: Padding, margins, alignment, width/height constraints
- 🧩 **Composition**: Join and place utilities for complex layouts
- 🔒 **Immutable**: Thread-safe builder pattern
- ⚡ **Fast**: Sub-millisecond rendering

## Installation

```bash
go get github.com/orchard9/tui-styles
```

## Quick Start

```go
package main

import (
    "fmt"
    tuistyles "github.com/orchard9/tui-styles"
)

func main() {
    red, _ := tuistyles.NewColor("#FF0000")
    blue, _ := tuistyles.NewColor("blue")

    style := tuistyles.NewStyle().
        Bold(true).
        Foreground(red).
        Background(blue).
        Padding(2).
        Border(tuistyles.RoundedBorder()).
        Width(50).
        Align(tuistyles.Center)

    fmt.Println(style.Render("Hello, World!"))
}
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/orchard9/tui-styles)
- [Examples](examples/)
- [CONTRIBUTING](CONTRIBUTING.md)
- [CHANGELOG](CHANGELOG.md)

## Examples

### Text Attributes

```go
bold := tuistyles.NewStyle().Bold(true)
italic := tuistyles.NewStyle().Italic(true)
underline := tuistyles.NewStyle().Underline(true)

fmt.Println(bold.Render("Bold"))
fmt.Println(italic.Render("Italic"))
fmt.Println(underline.Render("Underline"))
```

### Colors

```go
// Hex colors
red, _ := tuistyles.NewColor("#FF0000")
shortRed, _ := tuistyles.NewColor("#F00")  // Expands to #FF0000

// ANSI color names
blue, _ := tuistyles.NewColor("blue")
green, _ := tuistyles.NewColor("GREEN")   // Case-insensitive

// 256-color codes
orange, _ := tuistyles.NewColor("214")

style := tuistyles.NewStyle().
    Foreground(red).
    Background(blue)

fmt.Println(style.Render("Colored Text"))
```

### Borders

```go
// 8 predefined border styles
borders := []tuistyles.Border{
    tuistyles.NormalBorder(),      // ─│┌┐└┘
    tuistyles.RoundedBorder(),     // ─│╭╮╰╯
    tuistyles.ThickBorder(),       // ━┃┏┓┗┛
    tuistyles.DoubleBorder(),      // ═║╔╗╚╝
    tuistyles.BlockBorder(),       // █
    tuistyles.OuterHalfBlockBorder(),
    tuistyles.InnerHalfBlockBorder(),
    tuistyles.HiddenBorder(),      // Invisible
}

for _, border := range borders {
    style := tuistyles.NewStyle().
        Border(border).
        Padding(1)
    fmt.Println(style.Render("Border"))
}
```

### Alignment

```go
// Horizontal alignment
leftStyle := tuistyles.NewStyle().
    Width(40).
    Align(tuistyles.Left).
    Border(tuistyles.NormalBorder())

centerStyle := tuistyles.NewStyle().
    Width(40).
    Align(tuistyles.Center).
    Border(tuistyles.NormalBorder())

rightStyle := tuistyles.NewStyle().
    Width(40).
    Align(tuistyles.Right).
    Border(tuistyles.NormalBorder())

fmt.Println(leftStyle.Render("Left"))
fmt.Println(centerStyle.Render("Center"))
fmt.Println(rightStyle.Render("Right"))
```

### Layout Composition

```go
// Side-by-side panels
leftPanel := tuistyles.NewStyle().
    Border(tuistyles.RoundedBorder()).
    Width(30).
    Height(10).
    Render("Left Panel")

rightPanel := tuistyles.NewStyle().
    Border(tuistyles.RoundedBorder()).
    Width(30).
    Height(10).
    Render("Right Panel")

row := tuistyles.JoinHorizontal(tuistyles.Top, leftPanel, rightPanel)

// Vertical stacking
header := tuistyles.NewStyle().
    Width(60).
    Align(tuistyles.Center).
    Render("Header")

footer := tuistyles.NewStyle().
    Width(60).
    Align(tuistyles.Center).
    Render("Footer")

page := tuistyles.JoinVertical(tuistyles.Center, header, row, footer)

fmt.Println(page)
```

### Complete Example

See [examples/rendering/main.go](examples/rendering/main.go) for a complete showcase, or run:

```bash
go run examples/basic/main.go
go run examples/borders/main.go
go run examples/alignment/main.go
go run examples/rendering/main.go
```

## Performance

TUI Styles is optimized for speed:

| Operation | Time | Allocations |
|-----------|------|-------------|
| Simple text | ~50ns | 0 |
| Borders | <1ms | 6 |
| Complex layouts | <5ms | ~70 |

Run benchmarks: `go test -bench=. -benchmem`

## API Overview

### Text Attributes
- `Bold(bool)`, `Italic(bool)`, `Underline(bool)`
- `Strikethrough(bool)`, `Faint(bool)`, `Blink(bool)`, `Reverse(bool)`

### Colors
- `Foreground(Color)`, `Background(Color)`
- `BorderForeground(Color)`, `BorderBackground(Color)`

### Layout
- `Width(int)`, `Height(int)`, `MaxWidth(int)`, `MaxHeight(int)`
- `Align(Position)`, `AlignVertical(Position)`

### Spacing
- `Padding(values ...int)` - CSS-style shorthand
- `Margin(values ...int)` - CSS-style shorthand
- `PaddingTop`, `PaddingRight`, `PaddingBottom`, `PaddingLeft`
- `MarginTop`, `MarginRight`, `MarginBottom`, `MarginLeft`

### Borders
- `Border(Border, ...bool)` - Set border with optional sides
- `BorderTop`, `BorderRight`, `BorderBottom`, `BorderLeft`

### Rendering
- `Render(string) string` - Apply style to text

### Layout Utilities
- `JoinHorizontal(Position, ...string) string`
- `JoinVertical(Position, ...string) string`
- `Place(width, height int, hPos, vPos Position, content string) string`

## Development

### Prerequisites

- Go 1.21 or higher
- golangci-lint for linting

### Build and Test

```bash
# Run tests
go test ./... -v -race -cover

# Run linters
golangci-lint run

# Run benchmarks
go test -bench=. -benchmem

# Run examples
go run examples/basic/main.go
```

### Code Quality

- Test coverage: **96%**
- Zero linter warnings
- Comprehensive benchmarks
- Golden snapshot tests

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) file

## Credits

Inspired by [lipgloss](https://github.com/charmbracelet/lipgloss) by Charm.

## Related Projects

- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling library (inspiration)
- [glamour](https://github.com/charmbracelet/glamour) - Markdown rendering for the terminal
- [bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
