package tuistyles

import (
	"strings"

	"github.com/orchard9/tui-styles/internal/measure"
)

// Width sets the width of the styled text box in cells.
//
// Negative values are clamped to 0. Returns a new Style with width set,
// leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Width(80)
//	fmt.Println(s.Render("Text in 80-cell box"))
func (s Style) Width(w int) Style {
	if w < 0 {
		w = 0
	}
	s2 := s
	s2.width = &w
	return s2
}

// Height sets the height of the styled text box in lines.
//
// Negative values are clamped to 0. Returns a new Style with height set,
// leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Height(10)
//	fmt.Println(s.Render("Text in 10-line box"))
func (s Style) Height(h int) Style {
	if h < 0 {
		h = 0
	}
	s2 := s
	s2.height = &h
	return s2
}

// MaxWidth sets the maximum width in cells before text wrapping.
//
// Negative values are clamped to 0. Returns a new Style with maxWidth set,
// leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().MaxWidth(100)
//	fmt.Println(s.Render("Long text that will wrap at 100 cells"))
func (s Style) MaxWidth(w int) Style {
	if w < 0 {
		w = 0
	}
	s2 := s
	s2.maxWidth = &w
	return s2
}

// MaxHeight sets the maximum height in lines before text truncation.
//
// Negative values are clamped to 0. Returns a new Style with maxHeight set,
// leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().MaxHeight(20)
//	fmt.Println(s.Render("Long text that will be truncated after 20 lines"))
func (s Style) MaxHeight(h int) Style {
	if h < 0 {
		h = 0
	}
	s2 := s
	s2.maxHeight = &h
	return s2
}

// Align sets horizontal text alignment.
//
// Accepts Left, Center, or Right positions. Returns a new Style with align set,
// leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Width(80).Align(Center)
//	fmt.Println(s.Render("Centered text"))
func (s Style) Align(p Position) Style {
	s2 := s
	s2.align = &p
	return s2
}

// AlignVertical sets vertical text alignment.
//
// Accepts Top, Center, or Bottom positions. Returns a new Style with alignVertical set,
// leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Height(10).AlignVertical(Center)
//	fmt.Println(s.Render("Vertically centered text"))
func (s Style) AlignVertical(p Position) Style {
	s2 := s
	s2.alignVertical = &p
	return s2
}

// JoinHorizontal joins styled strings side-by-side with vertical alignment.
//
// pos determines how to align strings of different heights (Top, Center, or Bottom).
// All strings are placed next to each other horizontally, with their heights normalized
// to match the tallest string. Shorter strings are padded with spaces according to the
// specified vertical position.
//
// Example:
//
//	left := NewStyle().Background(red).Render("Left\nBox")
//	right := NewStyle().Background(blue).Render("Right\nBox\nThree")
//	combined := JoinHorizontal(Top, left, right)
//	// Result: left box aligned to top, right box below it
func JoinHorizontal(pos Position, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	// Split all strings into lines
	allLines := make([][]string, len(strs))
	widths := make([]int, len(strs))
	maxHeight := 0

	for i, str := range strs {
		allLines[i] = strings.Split(str, "\n")
		if len(allLines[i]) > maxHeight {
			maxHeight = len(allLines[i])
		}

		// Calculate max width for this column
		for _, line := range allLines[i] {
			w := measure.Width(line)
			if w > widths[i] {
				widths[i] = w
			}
		}
	}

	// Pad shorter strings vertically
	for i := range allLines {
		lines := allLines[i]
		currentHeight := len(lines)

		if currentHeight < maxHeight {
			padding := maxHeight - currentHeight
			emptyLine := strings.Repeat(" ", widths[i])

			switch pos {
			case Top:
				// Pad bottom
				for j := 0; j < padding; j++ {
					lines = append(lines, emptyLine)
				}
			case Center:
				// Pad both sides
				topPad := padding / 2
				bottomPad := padding - topPad
				for j := 0; j < topPad; j++ {
					lines = append([]string{emptyLine}, lines...)
				}
				for j := 0; j < bottomPad; j++ {
					lines = append(lines, emptyLine)
				}
			case Bottom:
				// Pad top
				for j := 0; j < padding; j++ {
					lines = append([]string{emptyLine}, lines...)
				}
			default:
				// Default to Top
				for j := 0; j < padding; j++ {
					lines = append(lines, emptyLine)
				}
			}

			allLines[i] = lines
		}

		// Pad lines to column width
		for j := range allLines[i] {
			lineWidth := measure.Width(allLines[i][j])
			if lineWidth < widths[i] {
				allLines[i][j] += strings.Repeat(" ", widths[i]-lineWidth)
			}
		}
	}

	// Join lines horizontally
	var result strings.Builder
	for row := 0; row < maxHeight; row++ {
		for col := 0; col < len(allLines); col++ {
			result.WriteString(allLines[col][row])
		}
		if row < maxHeight-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// JoinVertical stacks styled strings vertically with horizontal alignment.
//
// pos determines how to align strings of different widths (Left, Center, or Right).
// All strings are stacked on top of each other vertically, with their widths normalized
// to match the widest string. Narrower strings are padded with spaces according to the
// specified horizontal position.
//
// Example:
//
//	top := NewStyle().Background(red).Render("Top Box")
//	bottom := NewStyle().Background(blue).Render("Bottom Box - Wider")
//	combined := JoinVertical(Left, top, bottom)
//	// Result: boxes stacked vertically, left-aligned
func JoinVertical(pos Position, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	// Find max width
	maxWidth := 0
	for _, str := range strs {
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			w := measure.Width(line)
			if w > maxWidth {
				maxWidth = w
			}
		}
	}

	// Stack vertically with alignment
	var result strings.Builder
	first := true

	for _, str := range strs {
		if !first {
			result.WriteString("\n")
		}
		first = false

		lines := strings.Split(str, "\n")
		for i, line := range lines {
			lineWidth := measure.Width(line)
			padding := maxWidth - lineWidth

			if padding > 0 {
				switch pos {
				case Left:
					result.WriteString(line)
					result.WriteString(strings.Repeat(" ", padding))
				case Center:
					leftPad := padding / 2
					rightPad := padding - leftPad
					result.WriteString(strings.Repeat(" ", leftPad))
					result.WriteString(line)
					result.WriteString(strings.Repeat(" ", rightPad))
				case Right:
					result.WriteString(strings.Repeat(" ", padding))
					result.WriteString(line)
				default:
					// Default to Left
					result.WriteString(line)
					result.WriteString(strings.Repeat(" ", padding))
				}
			} else {
				result.WriteString(line)
			}

			if i < len(lines)-1 {
				result.WriteString("\n")
			}
		}
	}

	return result.String()
}

// Place positions content within a box of specified dimensions.
//
// hPos and vPos determine the placement (e.g., Top-Left, Center-Center, Bottom-Right).
// The content is positioned within a box of the given width and height, with the remaining
// space filled with spaces. If content exceeds the box dimensions, it will be clipped.
//
// Example:
//
//	content := NewStyle().Foreground(red).Render("Centered")
//	placed := Place(40, 10, Center, Center, content)
//	// Result: "Centered" appears in the middle of a 40x10 box
func Place(width, height int, hPos, vPos Position, content string) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	lines := strings.Split(content, "\n")

	// Measure content dimensions
	contentHeight := len(lines)
	contentWidth := 0
	for _, line := range lines {
		w := measure.Width(line)
		if w > contentWidth {
			contentWidth = w
		}
	}

	// Calculate start position
	var startRow, startCol int

	switch vPos {
	case Top:
		startRow = 0
	case Center:
		startRow = (height - contentHeight) / 2
		if startRow < 0 {
			startRow = 0
		}
	case Bottom:
		startRow = height - contentHeight
		if startRow < 0 {
			startRow = 0
		}
	default:
		startRow = 0
	}

	switch hPos {
	case Left:
		startCol = 0
	case Center:
		startCol = (width - contentWidth) / 2
		if startCol < 0 {
			startCol = 0
		}
	case Right:
		startCol = width - contentWidth
		if startCol < 0 {
			startCol = 0
		}
	default:
		startCol = 0
	}

	// Create box filled with spaces
	box := make([]string, height)
	for i := range box {
		box[i] = strings.Repeat(" ", width)
	}

	// Place content (clip if exceeds bounds)
	for i, line := range lines {
		row := startRow + i
		if row < 0 || row >= height {
			continue
		}

		lineWidth := measure.Width(line)
		endCol := startCol + lineWidth

		if endCol > width {
			// Truncate line
			lineWidth = width - startCol
			if lineWidth <= 0 {
				continue
			}
			// Use measure.Truncate for ANSI-aware truncation
			line = measure.Truncate(line, lineWidth, "")
		}

		// Place line in box (handle ANSI codes properly)
		if startCol >= 0 && startCol < width {
			// Build the line with the placed content
			box[row] = box[row][:startCol] + line + box[row][startCol+lineWidth:]
		}
	}

	return strings.Join(box, "\n")
}
