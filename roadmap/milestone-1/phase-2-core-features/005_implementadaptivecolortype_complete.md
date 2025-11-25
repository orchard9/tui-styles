## Purpose

Implement the `AdaptiveColor` type that automatically selects between light and dark color variants based on the terminal's background. This enables themes that adapt to user terminal preferences.

## Acceptance Criteria

- [ ] `AdaptiveColor` struct defined with `Light` and `Dark` fields
- [ ] Terminal background detection implemented (dark/light heuristic)
- [ ] `ToColor()` method returns appropriate Color based on terminal
- [ ] Unit tests with mocked terminal detection
- [ ] Fallback to Dark variant if detection fails
- [ ] All code passes `golangci-lint` with zero warnings

## Technical Approach

**AdaptiveColor Type** (`color.go`):
```go
// AdaptiveColor automatically selects color based on terminal background
type AdaptiveColor struct {
    Light Color  // Color for light terminal backgrounds
    Dark  Color  // Color for dark terminal backgrounds
}

// NewAdaptiveColor creates an AdaptiveColor with validation
func NewAdaptiveColor(light, dark string) (AdaptiveColor, error) {
    lightColor, err := NewColor(light)
    if err != nil {
        return AdaptiveColor{}, fmt.Errorf("invalid light color: %w", err)
    }

    darkColor, err := NewColor(dark)
    if err != nil {
        return AdaptiveColor{}, fmt.Errorf("invalid dark color: %w", err)
    }

    return AdaptiveColor{Light: lightColor, Dark: darkColor}, nil
}

// ToColor returns the appropriate color for the current terminal
func (ac AdaptiveColor) ToColor() Color {
    if isLightTerminal() {
        return ac.Light
    }
    return ac.Dark
}
```

**Terminal Detection** (`internal/ansi/terminal.go`):
```go
package ansi

import (
    "os"
    "strings"
)

// IsLightTerminal returns true if terminal has light background
// Uses heuristics: TERM_BACKGROUND env var, terminal type detection
func IsLightTerminal() bool {
    // Check TERM_BACKGROUND env var (if set by user)
    if bg := os.Getenv("TERM_BACKGROUND"); bg != "" {
        return strings.ToLower(bg) == "light"
    }

    // Check COLORFGBG env var (format: "foreground;background")
    // Background values: 0-6 = dark, 7-15 = light
    if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
        parts := strings.Split(colorfgbg, ";")
        if len(parts) == 2 {
            // Simple heuristic: bg > 6 is light
            if bg, err := strconv.Atoi(parts[1]); err == nil && bg > 6 {
                return true
            }
        }
    }

    // Default to dark terminal (conservative choice)
    return false
}
```

**Files to Create/Modify**:
- `color.go` - Add AdaptiveColor type
- `color_test.go` - Add AdaptiveColor tests
- `internal/ansi/terminal.go` - Terminal background detection
- `internal/ansi/terminal_test.go` - Terminal detection tests

**Dependencies**:
- Task 004 (Color type must exist first)
- Standard library: `os`, `strings`, `strconv`

## Testing Strategy

**Unit Tests** (`color_test.go`):
```go
func TestNewAdaptiveColor(t *testing.T) {
    tests := []struct {
        name      string
        light     string
        dark      string
        wantErr   bool
    }{
        {"valid colors", "#000000", "#FFFFFF", false},
        {"invalid light", "invalid", "#FFFFFF", true},
        {"invalid dark", "#000000", "invalid", true},
        {"both invalid", "bad", "worse", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewAdaptiveColor(tt.light, tt.dark)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewAdaptiveColor() error = %v, wantErr %v",
                    err, tt.wantErr)
            }
        })
    }
}

func TestAdaptiveColorToColor(t *testing.T) {
    ac, _ := NewAdaptiveColor("#000000", "#FFFFFF")

    // Test with mocked terminal detection
    t.Run("dark terminal", func(t *testing.T) {
        os.Setenv("TERM_BACKGROUND", "dark")
        defer os.Unsetenv("TERM_BACKGROUND")

        got := ac.ToColor()
        want, _ := NewColor("#FFFFFF")
        if got != want {
            t.Errorf("ToColor() = %v, want %v (dark)", got, want)
        }
    })

    t.Run("light terminal", func(t *testing.T) {
        os.Setenv("TERM_BACKGROUND", "light")
        defer os.Unsetenv("TERM_BACKGROUND")

        got := ac.ToColor()
        want, _ := NewColor("#000000")
        if got != want {
            t.Errorf("ToColor() = %v, want %v (light)", got, want)
        }
    })
}
```

**Terminal Detection Tests** (`internal/ansi/terminal_test.go`):
```go
func TestIsLightTerminal(t *testing.T) {
    tests := []struct {
        name       string
        envVars    map[string]string
        want       bool
    }{
        {"explicit light", map[string]string{"TERM_BACKGROUND": "light"}, true},
        {"explicit dark", map[string]string{"TERM_BACKGROUND": "dark"}, false},
        {"COLORFGBG light", map[string]string{"COLORFGBG": "0;15"}, true},
        {"COLORFGBG dark", map[string]string{"COLORFGBG": "15;0"}, false},
        {"no env vars", map[string]string{}, false},  // Default to dark
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set env vars
            for k, v := range tt.envVars {
                os.Setenv(k, v)
                defer os.Unsetenv(k)
            }

            if got := IsLightTerminal(); got != tt.want {
                t.Errorf("IsLightTerminal() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Notes

**Terminal Detection Limitations**:
- No foolproof way to detect terminal background without user cooperation
- Heuristics are best-effort (TERM_BACKGROUND, COLORFGBG)
- Some terminals don't expose background info
- Users can override with TERM_BACKGROUND env var

**Future Enhancement**:
- Query terminal background via ANSI escape sequences (complex, not all terminals support)
- Provide manual override in Style API (`.WithColorProfile()`)
- Defer advanced detection to Phase 2 if too complex

**Design Decision**: Default to dark terminal if uncertain (most developer terminals are dark)

**Reference**: See `spec.md` Section 1.1 for AdaptiveColor API specification.

**lipgloss Reference**: Review how [lipgloss handles adaptive colors](https://github.com/charmbracelet/lipgloss/blob/master/color.go#L65).



