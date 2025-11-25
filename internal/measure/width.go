// Package measure provides Unicode-aware text width measurement.
// It handles ANSI escape sequences, CJK characters, emoji, and zero-width
// characters to accurately measure terminal cell widths.
package measure

import (
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

// ansiRegex matches ANSI escape sequences to strip them before measuring
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Width returns the visible width of a string in terminal cells.
// It strips ANSI escape codes and accounts for Unicode character widths:
// - ASCII characters: 1 cell
// - CJK characters: 2 cells
// - Emoji: typically 2 cells (depends on terminal)
// - Zero-width characters: 0 cells
func Width(s string) int {
	// Strip ANSI escape codes first
	stripped := StripANSI(s)

	// Use runewidth to calculate actual width
	return runewidth.StringWidth(stripped)
}

// StripANSI removes all ANSI escape sequences from a string.
// This is useful for measuring the actual visible width of styled text.
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// WidthPerLine returns the width of each line in a multi-line string.
// Lines are split by newline characters (\n).
func WidthPerLine(s string) []int {
	lines := strings.Split(s, "\n")
	widths := make([]int, len(lines))
	for i, line := range lines {
		widths[i] = Width(line)
	}
	return widths
}

// MaxWidth returns the maximum width among all lines in a multi-line string.
// Returns 0 for empty strings.
func MaxWidth(s string) int {
	widths := WidthPerLine(s)
	maxWidth := 0
	for _, w := range widths {
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

// LineCount returns the number of lines in a string.
// An empty string has 1 line (not 0).
func LineCount(s string) int {
	if s == "" {
		return 1
	}
	return strings.Count(s, "\n") + 1
}

// Truncate truncates a string to fit within the specified width.
// If truncated, it appends the tail string (e.g., "...").
// The tail itself must fit within the width.
func Truncate(s string, width int, tail string) string {
	if width <= 0 {
		return ""
	}

	stripped := StripANSI(s)
	currentWidth := runewidth.StringWidth(stripped)

	if currentWidth <= width {
		return s // No truncation needed
	}

	tailWidth := runewidth.StringWidth(tail)
	if tailWidth >= width {
		// Tail too long, truncate tail itself
		return runewidth.Truncate(tail, width, "")
	}

	targetWidth := width - tailWidth
	truncated := runewidth.Truncate(stripped, targetWidth, "")

	return truncated + tail
}
