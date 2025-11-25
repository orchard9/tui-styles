package tuistyles

import (
	"strings"

	"github.com/orchard9/tui-styles/internal/ansi"
	"github.com/orchard9/tui-styles/internal/measure"
)

// Render applies the style to the given string and returns the ANSI-styled result.
//
// This method applies text attributes (bold, italic, etc.) and colors (foreground,
// background) to the input string. For multi-line strings, each line is styled
// independently with proper ANSI reset/reapplication at line boundaries.
//
// If MaxWidth is set, lines are truncated with ellipsis (...) if they exceed the width.
// Padding is rendered as colored spaces (using background color if set).
//
// Example:
//
//	red, _ := NewColor("red")
//	s := NewStyle().Bold(true).Foreground(red).Padding(2)
//	styled := s.Render("Hello, World!")
//	fmt.Println(styled) // Prints bold red text with 2 cells padding
func (s Style) Render(str string) string {
	// Allow rendering if we have padding or border, even with empty content
	if str == "" && !s.hasPadding() && !s.hasBorder() {
		return ""
	}

	// Apply basic rendering first
	var content string
	if str != "" {
		// Check if we need per-line rendering (multi-line or width constraints)
		if strings.Contains(str, "\n") || s.maxWidth != nil {
			content = s.renderMultiLine(str)
		} else {
			// Single line, no constraints - simple render
			content = s.renderSingleLine(str)
		}
	}

	// Apply horizontal alignment if width is set (before padding)
	if s.width != nil && s.align != nil {
		content = s.applyHorizontalAlignment(content)
	}

	// Apply vertical alignment if height is set (before padding)
	if s.height != nil {
		content = s.applyVerticalAlignment(content)
	}

	// Apply padding if set (after alignment)
	if s.hasPadding() {
		content = s.applyPadding(content)
	}

	// Apply border if set (wraps everything)
	if s.hasBorder() {
		content = s.applyBorder(content)
	}

	return content
}

// renderSingleLine applies styling to a single line of text
func (s Style) renderSingleLine(str string) string {
	var b strings.Builder
	b.Grow(len(str) + 50) // Pre-allocate for string + ANSI codes

	// Apply ANSI codes
	b.WriteString(s.stylePrefix())

	// Apply width constraint if set
	if s.maxWidth != nil && *s.maxWidth > 0 {
		width := measure.Width(str)
		if width > *s.maxWidth {
			str = measure.Truncate(str, *s.maxWidth, "...")
		}
	}

	b.WriteString(str)

	// Reset if any style was applied
	if s.hasAnyStyle() {
		b.WriteString(ansi.Reset())
	}

	return b.String()
}

// renderMultiLine applies styling to multi-line text, styling each line independently
func (s Style) renderMultiLine(str string) string {
	lines := strings.Split(str, "\n")
	styledLines := make([]string, len(lines))

	for i, line := range lines {
		// Apply width constraint per line if set
		if s.maxWidth != nil && *s.maxWidth > 0 {
			width := measure.Width(line)
			if width > *s.maxWidth {
				line = measure.Truncate(line, *s.maxWidth, "...")
			}
		}

		// Style each line independently
		if line == "" {
			// Empty line - just add newline separator later
			styledLines[i] = ""
		} else {
			var b strings.Builder
			b.Grow(len(line) + 50)

			// Apply ANSI codes to this line
			b.WriteString(s.stylePrefix())
			b.WriteString(line)

			// Reset after each line to prevent style bleed
			if s.hasAnyStyle() {
				b.WriteString(ansi.Reset())
			}

			styledLines[i] = b.String()
		}
	}

	return strings.Join(styledLines, "\n")
}

// stylePrefix returns the ANSI prefix codes for this style
func (s Style) stylePrefix() string {
	var b strings.Builder

	// Apply text attributes
	if s.bold != nil && *s.bold {
		b.WriteString(ansi.Bold())
	}
	if s.faint != nil && *s.faint {
		b.WriteString(ansi.Faint())
	}
	if s.italic != nil && *s.italic {
		b.WriteString(ansi.Italic())
	}
	if s.underline != nil && *s.underline {
		b.WriteString(ansi.Underline())
	}
	if s.blink != nil && *s.blink {
		b.WriteString(ansi.Blink())
	}
	if s.reverse != nil && *s.reverse {
		b.WriteString(ansi.Reverse())
	}
	if s.strikethrough != nil && *s.strikethrough {
		b.WriteString(ansi.Strikethrough())
	}

	// Apply foreground color
	if s.foreground != nil {
		b.WriteString(s.foreground.ToANSI())
	}

	// Apply background color
	if s.background != nil {
		b.WriteString(s.background.ToANSIBackground())
	}

	return b.String()
}

// String returns a string representation of the Style.
// For now, this returns an empty string. In future iterations,
// this may be used with a Value() builder method.
func (s Style) String() string {
	return ""
}

// hasAnyStyle returns true if any style attribute is set
func (s Style) hasAnyStyle() bool {
	return s.bold != nil || s.faint != nil || s.italic != nil ||
		s.underline != nil || s.blink != nil || s.reverse != nil ||
		s.strikethrough != nil || s.foreground != nil || s.background != nil
}

// hasPadding returns true if any padding is set
func (s Style) hasPadding() bool {
	return (s.paddingTop != nil && *s.paddingTop > 0) ||
		(s.paddingBottom != nil && *s.paddingBottom > 0) ||
		(s.paddingLeft != nil && *s.paddingLeft > 0) ||
		(s.paddingRight != nil && *s.paddingRight > 0)
}

// applyPadding adds padding around content as colored spaces
func (s Style) applyPadding(content string) string {
	// Get padding values (default to 0 if not set)
	paddingTop := 0
	if s.paddingTop != nil {
		paddingTop = *s.paddingTop
	}
	paddingBottom := 0
	if s.paddingBottom != nil {
		paddingBottom = *s.paddingBottom
	}
	paddingLeft := 0
	if s.paddingLeft != nil {
		paddingLeft = *s.paddingLeft
	}
	paddingRight := 0
	if s.paddingRight != nil {
		paddingRight = *s.paddingRight
	}

	if paddingTop == 0 && paddingBottom == 0 && paddingLeft == 0 && paddingRight == 0 {
		return content
	}

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Calculate content width (max line width after stripping ANSI)
	contentWidth := 0
	for _, line := range lines {
		w := measure.Width(line)
		if w > contentWidth {
			contentWidth = w
		}
	}

	// Create padding space (with background color if set)
	paddingSpace := s.makePaddingSpace(1)

	// Build result
	var b strings.Builder

	// Top padding lines
	if paddingTop > 0 {
		topLine := s.makePaddingSpace(contentWidth + paddingLeft + paddingRight)
		for i := 0; i < paddingTop; i++ {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(topLine)
		}
		if len(lines) > 0 || paddingBottom > 0 {
			b.WriteString("\n")
		}
	}

	// Content lines with left/right padding
	for i, line := range lines {
		if i > 0 {
			b.WriteString("\n")
		}

		// Left padding
		if paddingLeft > 0 {
			b.WriteString(strings.Repeat(paddingSpace, paddingLeft))
		}

		// Content
		b.WriteString(line)

		// Right padding (account for content width variation)
		if paddingRight > 0 {
			lineWidth := measure.Width(line)
			// Pad to content width, then add paddingRight
			if lineWidth < contentWidth {
				b.WriteString(strings.Repeat(" ", contentWidth-lineWidth))
			}
			b.WriteString(strings.Repeat(paddingSpace, paddingRight))
		}
	}

	// Bottom padding lines
	if paddingBottom > 0 {
		bottomLine := s.makePaddingSpace(contentWidth + paddingLeft + paddingRight)
		for i := 0; i < paddingBottom; i++ {
			b.WriteString("\n")
			b.WriteString(bottomLine)
		}
	}

	return b.String()
}

// makePaddingSpace creates a padding space string of specified width with background color
func (s Style) makePaddingSpace(width int) string {
	if width <= 0 {
		return ""
	}

	var b strings.Builder

	// Apply background color if set
	if s.background != nil {
		b.WriteString(s.background.ToANSIBackground())
	}

	// Write spaces
	b.WriteString(strings.Repeat(" ", width))

	// Reset if background was applied
	if s.background != nil {
		b.WriteString(ansi.Reset())
	}

	return b.String()
}

// hasBorder returns true if any border is configured
func (s Style) hasBorder() bool {
	if s.borderType == nil {
		return false
	}

	// Check if any border side is enabled
	return (s.borderTop != nil && *s.borderTop) ||
		(s.borderRight != nil && *s.borderRight) ||
		(s.borderBottom != nil && *s.borderBottom) ||
		(s.borderLeft != nil && *s.borderLeft)
}

// applyBorder wraps content with border characters
func (s Style) applyBorder(content string) string {
	if s.borderType == nil {
		return content
	}

	border := *s.borderType
	lines := strings.Split(content, "\n")

	// Measure content width (ANSI-aware)
	contentWidth := 0
	for _, line := range lines {
		w := measure.Width(line)
		if w > contentWidth {
			contentWidth = w
		}
	}

	// Determine which borders are enabled (default to true if not specified)
	topEnabled := s.borderTop == nil || *s.borderTop
	rightEnabled := s.borderRight == nil || *s.borderRight
	bottomEnabled := s.borderBottom == nil || *s.borderBottom
	leftEnabled := s.borderLeft == nil || *s.borderLeft

	var result strings.Builder

	// Top border
	if topEnabled {
		result.WriteString(s.renderBorderLine(border, contentWidth, true, leftEnabled, rightEnabled))
		result.WriteString("\n")
	}

	// Content lines with left/right borders
	for _, line := range lines {
		// Left border
		if leftEnabled {
			result.WriteString(s.styleBorderChar(border.Left))
		}

		// Content
		result.WriteString(line)

		// Pad line to contentWidth if needed
		lineWidth := measure.Width(line)
		if lineWidth < contentWidth {
			result.WriteString(strings.Repeat(" ", contentWidth-lineWidth))
		}

		// Right border
		if rightEnabled {
			result.WriteString(s.styleBorderChar(border.Right))
		}

		result.WriteString("\n")
	}

	// Bottom border
	if bottomEnabled {
		result.WriteString(s.renderBorderLine(border, contentWidth, false, leftEnabled, rightEnabled))
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// renderBorderLine creates a horizontal border line (top or bottom)
func (s Style) renderBorderLine(border Border, contentWidth int, isTop bool, leftEnabled bool, rightEnabled bool) string {
	var b strings.Builder

	// Select appropriate horizontal character and corners
	var horizontal, leftCorner, rightCorner string
	if isTop {
		horizontal = border.Top
		leftCorner = border.TopLeft
		rightCorner = border.TopRight
	} else {
		horizontal = border.Bottom
		leftCorner = border.BottomLeft
		rightCorner = border.BottomRight
	}

	// Left corner
	if leftEnabled {
		b.WriteString(s.styleBorderChar(leftCorner))
	}

	// Horizontal line
	b.WriteString(s.styleBorderChar(strings.Repeat(horizontal, contentWidth)))

	// Right corner
	if rightEnabled {
		b.WriteString(s.styleBorderChar(rightCorner))
	}

	return b.String()
}

// styleBorderChar applies border colors to a border character
func (s Style) styleBorderChar(char string) string {
	if s.borderForeground == nil && s.borderBackground == nil {
		return char
	}

	var b strings.Builder

	// Apply border colors
	if s.borderForeground != nil {
		b.WriteString(s.borderForeground.ToANSI())
	}
	if s.borderBackground != nil {
		b.WriteString(s.borderBackground.ToANSIBackground())
	}

	b.WriteString(char)

	// Reset if any border color was applied
	if s.borderForeground != nil || s.borderBackground != nil {
		b.WriteString(ansi.Reset())
	}

	return b.String()
}

// applyHorizontalAlignment aligns content horizontally within the specified width
func (s Style) applyHorizontalAlignment(content string) string {
	if s.width == nil || s.align == nil {
		return content
	}

	targetWidth := *s.width
	lines := strings.Split(content, "\n")
	var result strings.Builder

	for i, line := range lines {
		lineWidth := measure.Width(line)

		if lineWidth >= targetWidth {
			// Line is already at or exceeds target width, keep as is
			result.WriteString(line)
		} else {
			// Line is shorter than target width, apply alignment
			padding := targetWidth - lineWidth

			switch *s.align {
			case Left:
				// Left align: content on left, padding on right
				result.WriteString(line)
				result.WriteString(s.makeAlignmentSpace(padding))
			case Center:
				// Center align: distribute padding on both sides
				leftPad := padding / 2
				rightPad := padding - leftPad
				result.WriteString(s.makeAlignmentSpace(leftPad))
				result.WriteString(line)
				result.WriteString(s.makeAlignmentSpace(rightPad))
			case Right:
				// Right align: padding on left, content on right
				result.WriteString(s.makeAlignmentSpace(padding))
				result.WriteString(line)
			default:
				// Default to left alignment
				result.WriteString(line)
				result.WriteString(s.makeAlignmentSpace(padding))
			}
		}

		// Add newline separator between lines (but not after the last line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// makeAlignmentSpace creates alignment padding spaces (respects background color if set)
func (s Style) makeAlignmentSpace(width int) string {
	if width <= 0 {
		return ""
	}

	var b strings.Builder

	// Apply background color if set (alignment padding should match content background)
	if s.background != nil {
		b.WriteString(s.background.ToANSIBackground())
	}

	// Write spaces
	b.WriteString(strings.Repeat(" ", width))

	// Reset if background was applied
	if s.background != nil {
		b.WriteString(ansi.Reset())
	}

	return b.String()
}

// applyVerticalAlignment aligns content vertically within the specified height
func (s Style) applyVerticalAlignment(content string) string {
	if s.height == nil {
		return content
	}

	targetHeight := *s.height
	lines := strings.Split(content, "\n")
	currentHeight := len(lines)

	if currentHeight >= targetHeight {
		// Truncate if too tall
		return strings.Join(lines[:targetHeight], "\n")
	}

	// Calculate empty line width (use Style.width if set, else measure content)
	emptyLineWidth := 0
	if s.width != nil {
		emptyLineWidth = *s.width
	} else {
		for _, line := range lines {
			w := measure.Width(line)
			if w > emptyLineWidth {
				emptyLineWidth = w
			}
		}
	}

	// Ensure all existing content lines match the target width (pad with background color if needed)
	for i := range lines {
		lineWidth := measure.Width(lines[i])
		if lineWidth < emptyLineWidth {
			// Pad line to match width
			padding := emptyLineWidth - lineWidth
			lines[i] += s.makeAlignmentSpace(padding)
		}
	}

	// Create empty line with background color if set
	emptyLine := s.makeAlignmentSpace(emptyLineWidth)

	paddingLines := targetHeight - currentHeight

	// Default to Top alignment if not specified
	vAlign := Top
	if s.alignVertical != nil {
		vAlign = *s.alignVertical
	}

	switch vAlign {
	case Top:
		for i := 0; i < paddingLines; i++ {
			lines = append(lines, emptyLine)
		}
	case Center:
		topPad := paddingLines / 2
		bottomPad := paddingLines - topPad
		for i := 0; i < topPad; i++ {
			lines = append([]string{emptyLine}, lines...)
		}
		for i := 0; i < bottomPad; i++ {
			lines = append(lines, emptyLine)
		}
	case Bottom:
		for i := 0; i < paddingLines; i++ {
			lines = append([]string{emptyLine}, lines...)
		}
	}

	return strings.Join(lines, "\n")
}
