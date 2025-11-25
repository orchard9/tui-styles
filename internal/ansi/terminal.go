package ansi

import (
	"os"
	"strconv"
	"strings"
)

// IsLightTerminal returns true if terminal has a light background
// Uses heuristics: TERM_BACKGROUND env var, COLORFGBG env var
// Defaults to false (dark terminal) if detection is uncertain
func IsLightTerminal() bool {
	// Check TERM_BACKGROUND env var (user can explicitly set)
	if bg := os.Getenv("TERM_BACKGROUND"); bg != "" {
		return strings.ToLower(bg) == "light"
	}

	// Check COLORFGBG env var (format: "foreground;background")
	// Background values: 0-6 = dark, 7-15 = light (simplified heuristic)
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) == 2 {
			// Simple heuristic: background value > 6 indicates light terminal
			if bg, err := strconv.Atoi(parts[1]); err == nil && bg > 6 {
				return true
			}
		}
	}

	// Default to dark terminal (conservative choice, most dev terminals are dark)
	return false
}
