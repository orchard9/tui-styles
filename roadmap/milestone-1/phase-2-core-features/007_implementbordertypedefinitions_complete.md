## Purpose

Implement `BorderType` definitions with character mappings for all standard border styles. This provides the foundation for rendering borders around styled content.

## Acceptance Criteria

- [ ] `Border` struct defined with character fields (Top, Bottom, Left, Right, corners)
- [ ] Standard border styles implemented as functions (NormalBorder, RoundedBorder, etc.)
- [ ] All 8 border types from spec defined
- [ ] Unit tests validating character definitions
- [ ] Documentation for each border style
- [ ] All code passes `golangci-lint` with zero warnings

## Technical Approach

**Border Type Definition** (`border.go`):
```go
package tuistyles

// Border defines characters for drawing borders
type Border struct {
    Top         string
    Bottom      string
    Left        string
    Right       string
    TopLeft     string
    TopRight    string
    BottomLeft  string
    BottomRight string

    // Optional: Middle pieces for complex borders
    MiddleLeft  string
    MiddleRight string
    Middle      string
}

// NormalBorder returns standard box-drawing border
func NormalBorder() Border {
    return Border{
        Top:         "─",
        Bottom:      "─",
        Left:        "│",
        Right:       "│",
        TopLeft:     "┌",
        TopRight:    "┐",
        BottomLeft:  "└",
        BottomRight: "┘",
    }
}

// RoundedBorder returns border with rounded corners
func RoundedBorder() Border {
    return Border{
        Top:         "─",
        Bottom:      "─",
        Left:        "│",
        Right:       "│",
        TopLeft:     "╭",
        TopRight:    "╮",
        BottomLeft:  "╰",
        BottomRight: "╯",
    }
}

// BlockBorder returns solid block border
func BlockBorder() Border {
    return Border{
        Top:         "█",
        Bottom:      "█",
        Left:        "█",
        Right:       "█",
        TopLeft:     "█",
        TopRight:    "█",
        BottomLeft:  "█",
        BottomRight: "█",
    }
}

// OuterHalfBlockBorder returns outer half-block border
func OuterHalfBlockBorder() Border {
    return Border{
        Top:         "▀",
        Bottom:      "▄",
        Left:        "▌",
        Right:       "▐",
        TopLeft:     "▛",
        TopRight:    "▜",
        BottomLeft:  "▙",
        BottomRight: "▟",
    }
}

// InnerHalfBlockBorder returns inner half-block border
func InnerHalfBlockBorder() Border {
    return Border{
        Top:         "▄",
        Bottom:      "▀",
        Left:        "▐",
        Right:       "▌",
        TopLeft:     "▗",
        TopRight:    "▖",
        BottomLeft:  "▝",
        BottomRight: "▘",
    }
}

// ThickBorder returns thick box-drawing border
func ThickBorder() Border {
    return Border{
        Top:         "━",
        Bottom:      "━",
        Left:        "┃",
        Right:       "┃",
        TopLeft:     "┏",
        TopRight:    "┓",
        BottomLeft:  "┗",
        BottomRight: "┛",
    }
}

// DoubleBorder returns double-line border
func DoubleBorder() Border {
    return Border{
        Top:         "═",
        Bottom:      "═",
        Left:        "║",
        Right:       "║",
        TopLeft:     "╔",
        TopRight:    "╗",
        BottomLeft:  "╚",
        BottomRight: "╝",
    }
}

// HiddenBorder returns invisible border (spaces)
func HiddenBorder() Border {
    return Border{
        Top:         " ",
        Bottom:      " ",
        Left:        " ",
        Right:       " ",
        TopLeft:     " ",
        TopRight:    " ",
        BottomLeft:  " ",
        BottomRight: " ",
    }
}
```

**Files to Create/Modify**:
- `border.go` - Border struct and standard border definitions
- `border_test.go` - Unit tests for border types

**Dependencies**:
- None (standalone type, no rendering logic yet)

## Testing Strategy

**Unit Tests** (`border_test.go`):
```go
func TestBorderCharacters(t *testing.T) {
    tests := []struct {
        name   string
        border Border
        want   map[string]string
    }{
        {
            name:   "NormalBorder",
            border: NormalBorder(),
            want: map[string]string{
                "Top":         "─",
                "TopLeft":     "┌",
                "TopRight":    "┐",
                "Bottom":      "─",
                "BottomLeft":  "└",
                "BottomRight": "┘",
                "Left":        "│",
                "Right":       "│",
            },
        },
        {
            name:   "RoundedBorder",
            border: RoundedBorder(),
            want: map[string]string{
                "TopLeft":     "╭",
                "TopRight":    "╮",
                "BottomLeft":  "╰",
                "BottomRight": "╯",
            },
        },
        {
            name:   "HiddenBorder",
            border: HiddenBorder(),
            want: map[string]string{
                "Top":    " ",
                "Left":   " ",
                "Right":  " ",
                "Bottom": " ",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            for field, want := range tt.want {
                var got string
                switch field {
                case "Top":
                    got = tt.border.Top
                case "Bottom":
                    got = tt.border.Bottom
                case "Left":
                    got = tt.border.Left
                case "Right":
                    got = tt.border.Right
                case "TopLeft":
                    got = tt.border.TopLeft
                case "TopRight":
                    got = tt.border.TopRight
                case "BottomLeft":
                    got = tt.border.BottomLeft
                case "BottomRight":
                    got = tt.border.BottomRight
                }
                if got != want {
                    t.Errorf("%s.%s = %q, want %q", tt.name, field, got, want)
                }
            }
        })
    }
}

// Test all border types are defined
func TestAllBorderTypes(t *testing.T) {
    borders := []struct {
        name string
        fn   func() Border
    }{
        {"NormalBorder", NormalBorder},
        {"RoundedBorder", RoundedBorder},
        {"BlockBorder", BlockBorder},
        {"OuterHalfBlockBorder", OuterHalfBlockBorder},
        {"InnerHalfBlockBorder", InnerHalfBlockBorder},
        {"ThickBorder", ThickBorder},
        {"DoubleBorder", DoubleBorder},
        {"HiddenBorder", HiddenBorder},
    }

    for _, b := range borders {
        t.Run(b.name, func(t *testing.T) {
            border := b.fn()
            // Verify no fields are empty (except HiddenBorder uses spaces)
            if border.Top == "" || border.Bottom == "" ||
               border.Left == "" || border.Right == "" {
                t.Errorf("%s has empty fields", b.name)
            }
        })
    }
}
```

## Notes

**Unicode Box-Drawing Characters**:
- Normal: Uses standard box-drawing characters (U+2500 series)
- Thick: Uses heavy box-drawing characters (U+2501 series)
- Double: Uses double-line characters (U+2550 series)
- Rounded: Uses rounded corner characters (U+256D-U+2570)
- Blocks: Uses block element characters (U+2580 series)

**Terminal Compatibility**:
- Most modern terminals support Unicode box-drawing
- Older terminals may need fallback to ASCII (-, |, +)
- Consider environment variable (TERM_CHARSET) for fallback
- Defer fallback logic to Phase 2 (rendering implementation)

**Border Rendering**: This task only defines border character mappings. Actual rendering logic (drawing borders around content, respecting padding/margin) will be in Milestone 2 (Style struct implementation).

**Middle Pieces**: Some border types need middle connectors for complex layouts (e.g., table borders). The `Middle`, `MiddleLeft`, `MiddleRight` fields support this but are optional for Phase 1.

**Reference**: See `spec.md` Section 2.5 for Border API specification.

**lipgloss Reference**: Review [lipgloss border.go](https://github.com/charmbracelet/lipgloss/blob/master/border.go) for proven character mappings.


