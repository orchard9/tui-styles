package tuistyles

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/orchard9/tui-styles/internal/ansi"
)

// Color represents a terminal color (hex, ANSI name, or ANSI code)
type Color string

var hexColorRegex = regexp.MustCompile(`^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$`)

// NewColor creates a Color with validation
func NewColor(s string) (Color, error) {
	if s == "" {
		return "", fmt.Errorf("color cannot be empty")
	}

	// Validate hex (#RRGGBB or #RGB)
	if strings.HasPrefix(s, "#") {
		if !hexColorRegex.MatchString(s) {
			return "", fmt.Errorf("invalid hex color: %s", s)
		}
		// Normalize to uppercase and expand 3-digit to 6-digit
		normalized := normalizeHex(s)
		return Color(normalized), nil
	}

	// Validate ANSI color name (case-insensitive)
	if ansi.IsValidANSIName(s) {
		return Color(strings.ToLower(s)), nil
	}

	// Validate ANSI 256-color code
	if code, err := strconv.Atoi(s); err == nil {
		if code < 0 || code > 255 {
			return "", fmt.Errorf("ANSI color code out of range (0-255): %d", code)
		}
		return Color(s), nil
	}

	return "", fmt.Errorf("invalid color: %s (must be hex, ANSI name, or ANSI code 0-255)", s)
}

// ToANSI converts Color to ANSI foreground escape sequence
func (c Color) ToANSI() string {
	return ansi.ColorToANSI(string(c), false)
}

// ToANSIBackground converts Color to ANSI background escape sequence
func (c Color) ToANSIBackground() string {
	return ansi.ColorToANSI(string(c), true)
}

// normalizeHex converts hex color to uppercase and expands 3-digit to 6-digit
func normalizeHex(hex string) string {
	hex = strings.ToUpper(hex)
	if len(hex) == 4 { // #RGB
		return fmt.Sprintf("#%c%c%c%c%c%c",
			hex[1], hex[1], hex[2], hex[2], hex[3], hex[3])
	}
	return hex
}

// AdaptiveColor automatically selects color based on terminal background
type AdaptiveColor struct {
	Light Color // Color for light terminal backgrounds
	Dark  Color // Color for dark terminal backgrounds
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

// ToColor returns the appropriate color for the current terminal background
func (ac AdaptiveColor) ToColor() Color {
	if ansi.IsLightTerminal() {
		return ac.Light
	}
	return ac.Dark
}
