# TUI Styles - Architecture

## Overview

TUI Styles is a terminal styling library built on the principle of **immutability** and the **builder pattern**. This document describes the high-level design, core components, and architectural decisions.

## Design Principles

### 1. Immutability

All `Style` methods return new instances rather than mutating the receiver. This ensures:
- Thread safety
- Predictable behavior
- Safe style reuse and composition

```go
baseStyle := styles.NewStyle().Bold(true)
redStyle := baseStyle.Foreground(styles.Color("red"))   // baseStyle unchanged
blueStyle := baseStyle.Foreground(styles.Color("blue")) // independent copy
```

### 2. Builder Pattern

Fluent method chaining for ergonomic style construction:

```go
style := styles.NewStyle().
    Bold(true).
    Foreground(styles.Color("#FAFAFA")).
    Padding(2, 4).
    Border(styles.RoundedBorder(), true, true, true, true)
```

### 3. Type Safety

Strong typing for colors, positions, and borders prevents invalid configurations at compile time:

```go
type Position int // Left, Center, Right, Top, Bottom
type BorderType struct { ... } // Predefined border styles
```

## Core Components

### Public API

#### Style (style.go)
The fundamental building block holding all styling state:
- Text attributes (bold, italic, underline, etc.)
- Colors (foreground, background)
- Spacing (padding, margin)
- Borders (type, colors, visibility per side)
- Sizing and alignment

Methods: `NewStyle()`, `Render(string)`, `String()`, attribute setters.

#### Color (color.go)
Represents terminal colors with validation:
- Hex codes: `#RRGGBB`
- ANSI names: `red`, `blue`, `brightYellow`, etc.
- ANSI codes: `0-255`
- Adaptive colors: `AdaptiveColor{Light: "#000", Dark: "#FFF"}`

#### Position (position.go)
Enum for alignment:
- Horizontal: `Left`, `Center`, `Right`
- Vertical: `Top`, `Center`, `Bottom`

#### BorderType (border.go)
Defines border characters for all positions:
- Predefined styles: `NormalBorder()`, `RoundedBorder()`, `ThickBorder()`, etc.
- Custom borders supported

#### Layout Utilities
Package-level composition functions:
- `JoinHorizontal(Position, ...string)` - Horizontal concatenation with alignment
- `JoinVertical(Position, ...string)` - Vertical stacking with alignment
- `Place(width, height, hPos, vPos, content)` - Content placement in box

### Internal Packages

#### internal/ansi
ANSI escape code generation:
- Color codes (foreground, background)
- Attribute codes (bold, italic, etc.)
- Reset sequences
- Terminal capability detection

**Key functions**:
- `ColorCode(color Color, background bool) string`
- `AttributeCode(attr Attribute) string`
- `Reset() string`

#### internal/measure
String width calculation ignoring ANSI codes:
- Uses `mattn/go-runewidth` or equivalent logic
- Handles multi-byte characters (CJK, emojis)
- Splits strings into lines for rendering

**Key functions**:
- `Width(string) int`
- `Height(string) int`
- `Truncate(string, int) string`

## Data Flow

### Rendering Pipeline

1. **Style Construction**: User builds `Style` via method chaining
2. **Content Preparation**: `Render(content)` receives string to style
3. **Width Calculation**: `internal/measure` determines actual string width
4. **Layout Application**: Apply padding, alignment, borders
5. **ANSI Generation**: `internal/ansi` generates escape codes
6. **Output Assembly**: Combine ANSI codes + content + reset codes

### Example Flow

```
User: style.Render("Hello")
  ↓
Calculate content width (5 chars)
  ↓
Apply padding (add spaces based on Style.padding)
  ↓
Apply border (draw border chars if Style.border enabled)
  ↓
Apply colors/attributes (ANSI codes)
  ↓
Return: "\x1b[1m\x1b[38;2;250;250;250m  Hello  \x1b[0m"
```

## Module Structure

```
github.com/orchard9/tui-styles/
├── style.go              # Public Style type and methods
├── color.go              # Color type and validation
├── position.go           # Position enum
├── border.go             # BorderType definitions
├── layout.go             # JoinHorizontal, JoinVertical, Place
├── internal/
│   ├── ansi/
│   │   └── codes.go      # ANSI escape code generation
│   └── measure/
│       └── measure.go    # String width calculation
└── examples/
    └── basic/
        └── main.go       # Example programs
```

## Implementation Requirements

### Immutability

Every mutating method creates a new `Style`:

```go
func (s Style) Bold(enabled bool) Style {
    newStyle := s // copy all fields
    newStyle.bold = enabled
    return newStyle
}
```

### ANSI Awareness

Rendering must account for ANSI codes when measuring width:

```go
// Wrong: len("\x1b[1mHello\x1b[0m") = 14
// Correct: measure.Width("\x1b[1mHello\x1b[0m") = 5
```

### Multi-line Handling

Styles apply per-line for backgrounds and borders:

```go
style.Background(Color("red")).Render("Line 1\nLine 2")
// Each line gets red background, not just first line
```

## Dependencies

- **Go standard library only** for Phase 1
- Future consideration: `mattn/go-runewidth` for CJK character support

## Testing Strategy

- **Unit tests**: All core types and methods
- **Table-driven tests**: Color parsing, border rendering
- **Property-based tests**: Immutability guarantees
- **Snapshot tests**: Rendering output (visual regression)

## Performance Considerations

- **Avoid allocations**: Reuse buffers in rendering
- **Lazy evaluation**: Only generate ANSI codes when needed
- **String builder**: Use `strings.Builder` for concatenation

## Future Extensions

- Terminal capability detection (true color support)
- Custom border characters
- Gradient colors
- Hyperlinks (OSC 8 sequences)

## References

- [Specification](spec.md) - Complete API specification
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Inspiration and reference implementation
