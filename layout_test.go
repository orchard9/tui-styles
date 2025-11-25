package tuistyles

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDimensionMethods tests all dimension methods with table-driven approach.
func TestDimensionMethods(t *testing.T) {
	tests := []struct {
		name     string
		setter   func(Style, int) Style
		getter   func(Style) *int
		input    int
		expected int
	}{
		// Width tests
		{"Width_Positive", func(s Style, v int) Style { return s.Width(v) }, func(s Style) *int { return s.width }, 80, 80},
		{"Width_Zero", func(s Style, v int) Style { return s.Width(v) }, func(s Style) *int { return s.width }, 0, 0},
		{"Width_Negative", func(s Style, v int) Style { return s.Width(v) }, func(s Style) *int { return s.width }, -10, 0},

		// Height tests
		{"Height_Positive", func(s Style, v int) Style { return s.Height(v) }, func(s Style) *int { return s.height }, 24, 24},
		{"Height_Zero", func(s Style, v int) Style { return s.Height(v) }, func(s Style) *int { return s.height }, 0, 0},
		{"Height_Negative", func(s Style, v int) Style { return s.Height(v) }, func(s Style) *int { return s.height }, -5, 0},

		// MaxWidth tests
		{"MaxWidth_Positive", func(s Style, v int) Style { return s.MaxWidth(v) }, func(s Style) *int { return s.maxWidth }, 100, 100},
		{"MaxWidth_Zero", func(s Style, v int) Style { return s.MaxWidth(v) }, func(s Style) *int { return s.maxWidth }, 0, 0},
		{"MaxWidth_Negative", func(s Style, v int) Style { return s.MaxWidth(v) }, func(s Style) *int { return s.maxWidth }, -20, 0},

		// MaxHeight tests
		{"MaxHeight_Positive", func(s Style, v int) Style { return s.MaxHeight(v) }, func(s Style) *int { return s.maxHeight }, 50, 50},
		{"MaxHeight_Zero", func(s Style, v int) Style { return s.MaxHeight(v) }, func(s Style) *int { return s.maxHeight }, 0, 0},
		{"MaxHeight_Negative", func(s Style, v int) Style { return s.MaxHeight(v) }, func(s Style) *int { return s.maxHeight }, -15, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := tt.setter(s, tt.input)

			// Verify original unchanged (immutability)
			require.Nil(t, tt.getter(s), "Original Style was mutated")

			// Verify new Style has correct value
			require.NotNil(t, tt.getter(s2), "Value not set")
			require.Equal(t, tt.expected, *tt.getter(s2), "Value incorrect")
		})
	}
}

// TestWidth_Immutability verifies Width doesn't mutate the original Style.
func TestWidth_Immutability(t *testing.T) {
	s1 := NewStyle().Width(80)
	s2 := s1.Width(100)

	// s1 should still be 80
	require.NotNil(t, s1.width)
	require.Equal(t, 80, *s1.width, "s1 width was mutated")

	// s2 should be 100
	require.NotNil(t, s2.width)
	require.Equal(t, 100, *s2.width, "s2 width incorrect")
}

// TestHeight_Immutability verifies Height doesn't mutate the original Style.
func TestHeight_Immutability(t *testing.T) {
	s1 := NewStyle().Height(10)
	s2 := s1.Height(20)

	// s1 should still be 10
	require.NotNil(t, s1.height)
	require.Equal(t, 10, *s1.height, "s1 height was mutated")

	// s2 should be 20
	require.NotNil(t, s2.height)
	require.Equal(t, 20, *s2.height, "s2 height incorrect")
}

// TestAlign tests horizontal alignment.
func TestAlign(t *testing.T) {
	tests := []struct {
		name     string
		position Position
	}{
		{"Left", Left},
		{"Center", Center},
		{"Right", Right},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := s.Align(tt.position)

			// Verify original unchanged
			require.Nil(t, s.align, "Original Style was mutated")

			// Verify new Style has alignment set
			require.NotNil(t, s2.align, "Alignment not set")
			require.Equal(t, tt.position, *s2.align, "Alignment incorrect")
		})
	}
}

// TestAlignVertical tests vertical alignment.
func TestAlignVertical(t *testing.T) {
	tests := []struct {
		name     string
		position Position
	}{
		{"Top", Top},
		{"Center", Center},
		{"Bottom", Bottom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := s.AlignVertical(tt.position)

			// Verify original unchanged
			require.Nil(t, s.alignVertical, "Original Style was mutated")

			// Verify new Style has alignment set
			require.NotNil(t, s2.alignVertical, "Vertical alignment not set")
			require.Equal(t, tt.position, *s2.alignVertical, "Vertical alignment incorrect")
		})
	}
}

// TestAlign_Immutability verifies Align doesn't mutate the original Style.
func TestAlign_Immutability(t *testing.T) {
	s1 := NewStyle().Align(Left)
	s2 := s1.Align(Right)

	// s1 should still be Left
	require.NotNil(t, s1.align)
	require.Equal(t, Left, *s1.align, "s1 align was mutated")

	// s2 should be Right
	require.NotNil(t, s2.align)
	require.Equal(t, Right, *s2.align, "s2 align incorrect")
}

// TestAlignVertical_Immutability verifies AlignVertical doesn't mutate the original.
func TestAlignVertical_Immutability(t *testing.T) {
	s1 := NewStyle().AlignVertical(Top)
	s2 := s1.AlignVertical(Bottom)

	// s1 should still be Top
	require.NotNil(t, s1.alignVertical)
	require.Equal(t, Top, *s1.alignVertical, "s1 alignVertical was mutated")

	// s2 should be Bottom
	require.NotNil(t, s2.alignVertical)
	require.Equal(t, Bottom, *s2.alignVertical, "s2 alignVertical incorrect")
}

// TestLayoutMethods_Chaining verifies layout methods work in chains.
func TestLayoutMethods_Chaining(t *testing.T) {
	s := NewStyle().
		Width(80).
		Height(24).
		MaxWidth(100).
		MaxHeight(50).
		Align(Center).
		AlignVertical(Center)

	// Verify all layout fields set
	require.NotNil(t, s.width)
	require.Equal(t, 80, *s.width)

	require.NotNil(t, s.height)
	require.Equal(t, 24, *s.height)

	require.NotNil(t, s.maxWidth)
	require.Equal(t, 100, *s.maxWidth)

	require.NotNil(t, s.maxHeight)
	require.Equal(t, 50, *s.maxHeight)

	require.NotNil(t, s.align)
	require.Equal(t, Center, *s.align)

	require.NotNil(t, s.alignVertical)
	require.Equal(t, Center, *s.alignVertical)
}

// TestLayoutWithTextAttributes verifies layout methods chain with text attributes.
func TestLayoutWithTextAttributes(t *testing.T) {
	s := NewStyle().
		Bold(true).
		Width(80).
		Italic(true).
		Align(Center)

	require.NotNil(t, s.bold)
	require.True(t, *s.bold)

	require.NotNil(t, s.italic)
	require.True(t, *s.italic)

	require.NotNil(t, s.width)
	require.Equal(t, 80, *s.width)

	require.NotNil(t, s.align)
	require.Equal(t, Center, *s.align)
}

// TestLayoutWithColors verifies layout methods chain with color methods.
func TestLayoutWithColors(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	s := NewStyle().
		Foreground(red).
		Width(80).
		Background(blue).
		Height(24)

	require.NotNil(t, s.foreground)
	require.Equal(t, red, *s.foreground)

	require.NotNil(t, s.background)
	require.Equal(t, blue, *s.background)

	require.NotNil(t, s.width)
	require.Equal(t, 80, *s.width)

	require.NotNil(t, s.height)
	require.Equal(t, 24, *s.height)
}

// TestDimensionValidation_EdgeCases tests edge cases for dimension validation.
func TestDimensionValidation_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"LargePositive", 10000, 10000},
		{"SmallNegative", -1, 0},
		{"LargeNegative", -9999, 0},
		{"Zero", 0, 0},
		{"One", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test all dimension methods
			s1 := NewStyle().Width(tt.value)
			require.Equal(t, tt.expected, *s1.width, "Width validation failed")

			s2 := NewStyle().Height(tt.value)
			require.Equal(t, tt.expected, *s2.height, "Height validation failed")

			s3 := NewStyle().MaxWidth(tt.value)
			require.Equal(t, tt.expected, *s3.maxWidth, "MaxWidth validation failed")

			s4 := NewStyle().MaxHeight(tt.value)
			require.Equal(t, tt.expected, *s4.maxHeight, "MaxHeight validation failed")
		})
	}
}

// TestLayoutMethods_IndependentBranches verifies independent Style branches.
func TestLayoutMethods_IndependentBranches(t *testing.T) {
	base := NewStyle()
	s1 := base.Width(80)
	s2 := base.Height(24)
	s3 := base.Align(Center)

	// Base should be unchanged
	require.Nil(t, base.width)
	require.Nil(t, base.height)
	require.Nil(t, base.align)

	// s1 should only have width
	require.NotNil(t, s1.width)
	require.Equal(t, 80, *s1.width)
	require.Nil(t, s1.height)
	require.Nil(t, s1.align)

	// s2 should only have height
	require.Nil(t, s2.width)
	require.NotNil(t, s2.height)
	require.Equal(t, 24, *s2.height)
	require.Nil(t, s2.align)

	// s3 should only have align
	require.Nil(t, s3.width)
	require.Nil(t, s3.height)
	require.NotNil(t, s3.align)
	require.Equal(t, Center, *s3.align)
}

// TestJoinHorizontal tests the JoinHorizontal layout function
func TestJoinHorizontal(t *testing.T) {
	tests := []struct {
		name    string
		pos     Position
		strs    []string
		checkFn func(*testing.T, string)
	}{
		{
			name: "empty input",
			pos:  Top,
			strs: []string{},
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "", output, "Empty input should return empty string")
			},
		},
		{
			name: "single string",
			pos:  Top,
			strs: []string{"hello"},
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "hello", output, "Single string should be returned as-is")
			},
		},
		{
			name: "two strings same height",
			pos:  Top,
			strs: []string{"left", "right"},
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "leftright", output, "Should join side-by-side")
			},
		},
		{
			name: "two multi-line strings same height",
			pos:  Top,
			strs: []string{"L1\nL2", "R1\nR2"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 2, len(lines), "Should have 2 lines")
				require.Equal(t, "L1R1", lines[0])
				require.Equal(t, "L2R2", lines[1])
			},
		},
		{
			name: "different heights - top alignment",
			pos:  Top,
			strs: []string{"L1\nL2\nL3", "R1"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines (max height)")
				require.Equal(t, "L1R1", lines[0])
				require.Contains(t, lines[1], "L2")
				require.Contains(t, lines[2], "L3")
				// R1 should be on first line, spaces on remaining lines
			},
		},
		{
			name: "different heights - center alignment",
			pos:  Center,
			strs: []string{"L1\nL2\nL3", "R1"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				// With 1 line content and 3 total, R should be at center (line 1)
				require.Contains(t, lines[1], "R1")
			},
		},
		{
			name: "different heights - bottom alignment",
			pos:  Bottom,
			strs: []string{"L1\nL2\nL3", "R1"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				// R should be at bottom (line 2)
				require.Contains(t, lines[2], "R1")
			},
		},
		{
			name: "three strings different heights",
			pos:  Top,
			strs: []string{"A", "B1\nB2", "C1\nC2\nC3"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines (max height)")
			},
		},
		{
			name: "preserve width per column",
			pos:  Top,
			strs: []string{"short\nS", "longer line\nL"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				// First column should be padded to "short" width (5)
				// Second column should be padded to "longer line" width (11)
				// Total width should be consistent
				width1 := len(lines[0])
				width2 := len(lines[1])
				require.Equal(t, width1, width2, "All lines should have same width")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := JoinHorizontal(tt.pos, tt.strs...)
			tt.checkFn(t, output)
		})
	}
}

// TestJoinVertical tests the JoinVertical layout function
func TestJoinVertical(t *testing.T) {
	tests := []struct {
		name    string
		pos     Position
		strs    []string
		checkFn func(*testing.T, string)
	}{
		{
			name: "empty input",
			pos:  Left,
			strs: []string{},
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "", output, "Empty input should return empty string")
			},
		},
		{
			name: "single string",
			pos:  Left,
			strs: []string{"hello"},
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "hello", output, "Single string should be returned as-is")
			},
		},
		{
			name: "two strings same width",
			pos:  Left,
			strs: []string{"top", "bot"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 2, len(lines), "Should have 2 lines")
				require.Equal(t, "top", lines[0])
				require.Equal(t, "bot", lines[1])
			},
		},
		{
			name: "two multi-line strings same width",
			pos:  Left,
			strs: []string{"T1\nT2", "B1\nB2"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 4, len(lines), "Should have 4 lines total")
				require.Equal(t, "T1", lines[0])
				require.Equal(t, "T2", lines[1])
				require.Equal(t, "B1", lines[2])
				require.Equal(t, "B2", lines[3])
			},
		},
		{
			name: "different widths - left alignment",
			pos:  Left,
			strs: []string{"short", "longer line"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 2, len(lines), "Should have 2 lines")
				// Both lines should be padded to same width (11 = "longer line")
				require.Equal(t, 11, len(lines[0]), "First line should be padded to max width")
				require.Equal(t, 11, len(lines[1]), "Second line should match max width")
				require.True(t, strings.HasPrefix(lines[0], "short"), "Should be left-aligned")
			},
		},
		{
			name: "different widths - center alignment",
			pos:  Center,
			strs: []string{"short", "longer line"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 2, len(lines), "Should have 2 lines")
				// Both lines should be same width
				require.Equal(t, len(lines[0]), len(lines[1]), "Lines should have same width")
				// "short" should be centered
				require.Contains(t, lines[0], "short")
				// Should have spaces on both sides
				require.True(t, strings.Contains(lines[0], " short"), "Should have leading space for centering")
			},
		},
		{
			name: "different widths - right alignment",
			pos:  Right,
			strs: []string{"short", "longer line"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 2, len(lines), "Should have 2 lines")
				// Both lines should be same width
				require.Equal(t, len(lines[0]), len(lines[1]), "Lines should have same width")
				// "short" should be right-aligned
				require.True(t, strings.HasSuffix(lines[0], "short"), "Should be right-aligned")
			},
		},
		{
			name: "three strings different widths",
			pos:  Left,
			strs: []string{"A", "BB", "CCC"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				// All lines should be padded to width 3
				for i, line := range lines {
					require.Equal(t, 3, len(line), "Line %d should be padded to max width", i)
				}
			},
		},
		{
			name: "multi-line strings with different widths",
			pos:  Center,
			strs: []string{"T1\nT2 wider", "B1\nB2"},
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 4, len(lines), "Should have 4 lines total")
				// All lines should be same width (8 = "T2 wider")
				for i, line := range lines {
					require.Equal(t, 8, len(line), "Line %d should be padded to max width", i)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := JoinVertical(tt.pos, tt.strs...)
			tt.checkFn(t, output)
		})
	}
}

// TestPlace tests the Place layout function
func TestPlace(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		hPos    Position
		vPos    Position
		content string
		checkFn func(*testing.T, string)
	}{
		{
			name:    "zero width",
			width:   0,
			height:  5,
			hPos:    Left,
			vPos:    Top,
			content: "test",
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "", output, "Zero width should return empty string")
			},
		},
		{
			name:    "zero height",
			width:   10,
			height:  0,
			hPos:    Left,
			vPos:    Top,
			content: "test",
			checkFn: func(t *testing.T, output string) {
				require.Equal(t, "", output, "Zero height should return empty string")
			},
		},
		{
			name:    "top-left placement",
			width:   10,
			height:  3,
			hPos:    Left,
			vPos:    Top,
			content: "TL",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				require.True(t, strings.HasPrefix(lines[0], "TL"), "Content should be at top-left")
				require.Equal(t, 10, len(lines[0]), "Lines should be width 10")
			},
		},
		{
			name:    "center-center placement",
			width:   10,
			height:  5,
			hPos:    Center,
			vPos:    Center,
			content: "C",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 5, len(lines), "Should have 5 lines")
				// Single char should be at line 2 (middle of 5)
				require.Contains(t, lines[2], "C", "Content should be in middle line")
				// Should be centered horizontally
				require.Equal(t, 10, len(lines[2]), "Line should be width 10")
			},
		},
		{
			name:    "bottom-right placement",
			width:   10,
			height:  3,
			hPos:    Right,
			vPos:    Bottom,
			content: "BR",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				// Content should be at bottom (line 2)
				require.True(t, strings.HasSuffix(lines[2], "BR"), "Content should be at bottom-right")
			},
		},
		{
			name:    "multi-line content centered",
			width:   15,
			height:  5,
			hPos:    Center,
			vPos:    Center,
			content: "L1\nL2",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 5, len(lines), "Should have 5 lines")
				// 2 lines of content in 5 total, centered vertically = lines 1-2
				require.Contains(t, lines[1], "L1", "First content line should be at line 1")
				require.Contains(t, lines[2], "L2", "Second content line should be at line 2")
			},
		},
		{
			name:    "content exceeds width - truncate",
			width:   5,
			height:  3,
			hPos:    Left,
			vPos:    Top,
			content: "This is too long",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				// First line should be truncated to width 5
				require.Equal(t, 5, len(lines[0]), "Line should be truncated to width 5")
			},
		},
		{
			name:    "content exceeds height - clip",
			width:   10,
			height:  2,
			hPos:    Left,
			vPos:    Top,
			content: "L1\nL2\nL3\nL4",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 2, len(lines), "Should have 2 lines (clipped)")
				require.Contains(t, lines[0], "L1", "First line should be L1")
				require.Contains(t, lines[1], "L2", "Second line should be L2")
			},
		},
		{
			name:    "empty content",
			width:   10,
			height:  5,
			hPos:    Center,
			vPos:    Center,
			content: "",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 5, len(lines), "Should have 5 lines")
				// All lines should be spaces
				for i, line := range lines {
					require.Equal(t, 10, len(line), "Line %d should be width 10", i)
					require.Equal(t, strings.Repeat(" ", 10), line, "Line %d should be all spaces", i)
				}
			},
		},
		{
			name:    "all nine positions",
			width:   9,
			height:  3,
			hPos:    Center,
			vPos:    Center,
			content: "X",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				require.Equal(t, 3, len(lines), "Should have 3 lines")
				// X should be in center (line 1, position 4)
				require.Contains(t, lines[1], "X", "Content should be in middle line")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Place(tt.width, tt.height, tt.hPos, tt.vPos, tt.content)
			tt.checkFn(t, output)
		})
	}
}
