## Purpose

Implement the `Color` type with validation and ANSI code mapping. This is the foundational type for all color operations in the library, supporting hex colors, ANSI color names, and ANSI 256-color codes.

## Acceptance Criteria

- [ ] `Color` type defined as string with validation
- [ ] `NewColor(string) (Color, error)` constructor validates input
- [ ] Support hex colors (`#RRGGBB`, `#RGB`)
- [ ] Support ANSI color names (red, blue, green, yellow, etc.)
- [ ] Support ANSI 256-color codes (0-255)
- [ ] `ToANSI()` method converts Color to ANSI escape sequence
- [ ] Comprehensive unit tests with table-driven approach
- [ ] All code passes `golangci-lint` with zero warnings

## Technical Approach

**Color Type Definition** (`color.go`):
```go
package tuistyles

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

// Color represents a terminal color (hex, ANSI name, or ANSI code)
type Color string

// NewColor creates a Color with validation
func NewColor(s string) (Color, error) {
    // Validate hex (#RRGGBB or #RGB)
    if strings.HasPrefix(s, "#") {
        if !isValidHex(s) {
            return "", fmt.Errorf("invalid hex color: %s", s)
        }
        return Color(s), nil
    }

    // Validate ANSI color name
    if isValidANSIName(s) {
        return Color(s), nil
    }

    // Validate ANSI 256-color code
    if code, err := strconv.Atoi(s); err == nil {
        if code < 0 || code > 255 {
            return "", fmt.Errorf("ANSI color code out of range: %d", code)
        }
        return Color(s), nil
    }

    return "", fmt.Errorf("invalid color: %s", s)
}

// ToANSI converts Color to ANSI foreground escape sequence
func (c Color) ToANSI() string {
    // Implementation using internal/ansi package
    return ansi.ColorToANSI(string(c))
}
```

**ANSI Mapping** (`internal/ansi/codes.go`):
```go
package ansi

var ansiColorNames = map[string]int{
    "black":   0,
    "red":     1,
    "green":   2,
    "yellow":  3,
    "blue":    4,
    "magenta": 5,
    "cyan":    6,
    "white":   7,
    // Bright colors
    "bright-black":   8,
    "bright-red":     9,
    "bright-green":   10,
    "bright-yellow":  11,
    "bright-blue":    12,
    "bright-magenta": 13,
    "bright-cyan":    14,
    "bright-white":   15,
}

func ColorToANSI(color string) string {
    // Convert hex to RGB, then to ANSI 256 approximation
    // Convert ANSI name to code
    // Return raw code if numeric
}
```

**Files to Create/Modify**:
- `color.go` - Color type definition, constructor, methods
- `color_test.go` - Unit tests for Color
- `internal/ansi/codes.go` - ANSI color mappings and conversion
- `internal/ansi/codes_test.go` - ANSI mapping tests

**Dependencies**:
- Standard library: `fmt`, `regexp`, `strconv`, `strings`
- No external dependencies

## Testing Strategy

**Unit Tests** (`color_test.go`):
```go
func TestNewColor(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid hex 6-digit", "#FF0000", false},
        {"valid hex 3-digit", "#F00", false},
        {"valid ANSI name", "red", false},
        {"valid ANSI code", "196", false},
        {"invalid hex", "#GGGGGG", true},
        {"invalid code", "256", true},
        {"invalid name", "notacolor", true},
        {"empty string", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewColor(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewColor(%q) error = %v, wantErr %v",
                    tt.input, err, tt.wantErr)
            }
        })
    }
}

func TestColorToANSI(t *testing.T) {
    tests := []struct {
        color string
        want  string  // ANSI escape sequence
    }{
        {"red", "\x1b[31m"},
        {"#FF0000", "\x1b[38;2;255;0;0m"},  // RGB
        {"196", "\x1b[38;5;196m"},          // 256-color
    }

    for _, tt := range tests {
        t.Run(tt.color, func(t *testing.T) {
            c, err := NewColor(tt.color)
            if err != nil {
                t.Fatalf("NewColor failed: %v", err)
            }
            if got := c.ToANSI(); got != tt.want {
                t.Errorf("ToANSI() = %q, want %q", got, tt.want)
            }
        })
    }
}
```

**Edge Cases to Test**:
- Lowercase/uppercase hex (#ff0000 vs #FF0000)
- 3-digit hex expansion (#F00 → #FF0000)
- Case-insensitive ANSI names (Red, RED, red)
- Boundary ANSI codes (0, 255, 256)
- Empty strings and whitespace

## Notes

**Hex Color Normalization**:
- Accept both `#RGB` and `#RRGGBB` formats
- Normalize `#RGB` to `#RRGGBB` internally (#F00 → #FF0000)
- Store normalized form for consistency

**ANSI Color Approximation**:
- True hex colors (#RRGGBB) will be approximated to ANSI 256-color palette
- Use Euclidean distance in RGB space to find closest ANSI color
- Reference: [ANSI 256-color lookup table](https://www.ditig.com/256-colors-cheat-sheet)

**Performance**: Color validation happens at construction time, not at render time. Cache ANSI codes if needed.

**Reference**: See `spec.md` Section 1.1 for Color API specification.

**lipgloss Reference**: Review [lipgloss color.go](https://github.com/charmbracelet/lipgloss/blob/master/color.go) for proven patterns.


