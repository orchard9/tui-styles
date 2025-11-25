// Package ansi provides ANSI escape code generation for terminal styling.
// It handles color codes (16-color, 256-color, true color), text attributes
// (bold, italic, underline), and sequence composition.
package ansi

import (
	"fmt"
	"strconv"
	"strings"
)

// ANSI color name to code mapping
var ansiColorNames = map[string]int{
	// Standard colors (0-7)
	"black":   0,
	"red":     1,
	"green":   2,
	"yellow":  3,
	"blue":    4,
	"magenta": 5,
	"cyan":    6,
	"white":   7,
	// Bright colors (8-15)
	"bright-black":   8,
	"bright-red":     9,
	"bright-green":   10,
	"bright-yellow":  11,
	"bright-blue":    12,
	"bright-magenta": 13,
	"bright-cyan":    14,
	"bright-white":   15,
	// Aliases for bright colors
	"gray":        8,
	"grey":        8,
	"bright-gray": 15,
	"bright-grey": 15,
}

// IsValidANSIName returns true if the name is a valid ANSI color name
func IsValidANSIName(name string) bool {
	_, exists := ansiColorNames[strings.ToLower(name)]
	return exists
}

// ColorToANSI converts a color string to ANSI escape sequence (foreground)
func ColorToANSI(color string, background bool) string {
	// Handle hex colors (#RRGGBB)
	if strings.HasPrefix(color, "#") {
		r, g, b, err := hexToRGB(color)
		if err != nil {
			return ""
		}
		if background {
			return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
		}
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}

	// Handle ANSI color names
	if code, exists := ansiColorNames[strings.ToLower(color)]; exists {
		if background {
			if code < 8 {
				return fmt.Sprintf("\x1b[%dm", 40+code)
			}
			return fmt.Sprintf("\x1b[%dm", 100+code-8)
		}
		if code < 8 {
			return fmt.Sprintf("\x1b[%dm", 30+code)
		}
		return fmt.Sprintf("\x1b[%dm", 90+code-8)
	}

	// Handle ANSI 256-color codes
	if code, err := strconv.Atoi(color); err == nil {
		if code >= 0 && code <= 255 {
			if background {
				return fmt.Sprintf("\x1b[48;5;%dm", code)
			}
			return fmt.Sprintf("\x1b[38;5;%dm", code)
		}
	}

	return ""
}

// hexToRGB converts hex color to RGB values
func hexToRGB(hex string) (r, g, b int, err error) {
	hex = strings.TrimPrefix(hex, "#")

	// Expand 3-digit hex to 6-digit (#RGB -> #RRGGBB)
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color length: %s", hex)
	}

	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %s", hex)
	}

	// These conversions are safe because values are masked with 0xFF (max 255)
	r = int((values >> 16) & 0xFF) //nolint:gosec // G115: Safe conversion, masked to 0-255
	g = int((values >> 8) & 0xFF)  //nolint:gosec // G115: Safe conversion, masked to 0-255
	b = int(values & 0xFF)         //nolint:gosec // G115: Safe conversion, masked to 0-255

	return r, g, b, nil
}

// Reset returns the ANSI reset sequence
func Reset() string {
	return "\x1b[0m"
}

// Text attribute ANSI codes

// Bold returns the ANSI code for bold/bright text
func Bold() string {
	return "\x1b[1m"
}

// Faint returns the ANSI code for faint/dim text
func Faint() string {
	return "\x1b[2m"
}

// Italic returns the ANSI code for italic text
func Italic() string {
	return "\x1b[3m"
}

// Underline returns the ANSI code for underlined text
func Underline() string {
	return "\x1b[4m"
}

// Blink returns the ANSI code for blinking text
func Blink() string {
	return "\x1b[5m"
}

// Reverse returns the ANSI code for reverse video (swap fg/bg)
func Reverse() string {
	return "\x1b[7m"
}

// Strikethrough returns the ANSI code for strikethrough text
func Strikethrough() string {
	return "\x1b[9m"
}

// NoBold returns the ANSI code to disable bold
func NoBold() string {
	return "\x1b[22m"
}

// NoItalic returns the ANSI code to disable italic
func NoItalic() string {
	return "\x1b[23m"
}

// NoUnderline returns the ANSI code to disable underline
func NoUnderline() string {
	return "\x1b[24m"
}

// NoBlink returns the ANSI code to disable blink
func NoBlink() string {
	return "\x1b[25m"
}

// NoReverse returns the ANSI code to disable reverse video
func NoReverse() string {
	return "\x1b[27m"
}

// NoStrikethrough returns the ANSI code to disable strikethrough
func NoStrikethrough() string {
	return "\x1b[29m"
}

// ForegroundColor returns the ANSI escape sequence for the given foreground color
func ForegroundColor(color string) string {
	return ColorToANSI(color, false)
}

// BackgroundColor returns the ANSI escape sequence for the given background color
func BackgroundColor(color string) string {
	return ColorToANSI(color, true)
}
