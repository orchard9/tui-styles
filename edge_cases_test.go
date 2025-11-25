package tuistyles

import (
	"strings"
	"testing"

	"github.com/orchard9/tui-styles/internal/measure"
)

// TestEdgeCases_EmptyStrings tests handling of empty strings across all functions
func TestEdgeCases_EmptyStrings(t *testing.T) {
	t.Run("Render empty string", func(t *testing.T) {
		s := NewStyle()
		output := s.Render("")
		if output != "" {
			t.Errorf("Expected empty string, got %q", output)
		}
	})

	t.Run("Render empty string with styles", func(t *testing.T) {
		s := NewStyle().Bold(true)
		output := s.Render("")
		if output != "" {
			t.Errorf("Expected empty string even with styles, got %q", output)
		}
	})

	t.Run("Render empty string with padding", func(t *testing.T) {
		s := NewStyle().Padding(2)
		output := s.Render("")
		// With padding but no content, should still render padding lines
		if output == "" {
			t.Error("Expected padding to render even with empty content")
		}
	})

	t.Run("JoinHorizontal with empty strings", func(t *testing.T) {
		output := JoinHorizontal(Top, "", "", "")
		if output != "" {
			t.Errorf("Expected empty result, got %q", output)
		}
	})

	t.Run("JoinVertical with empty strings", func(t *testing.T) {
		output := JoinVertical(Left, "", "", "")
		// Should join with newlines
		lines := strings.Split(output, "\n")
		if len(lines) != 3 {
			t.Errorf("Expected 3 lines, got %d", len(lines))
		}
	})

	t.Run("Place empty content", func(t *testing.T) {
		output := Place(10, 5, Center, Center, "")
		lines := strings.Split(output, "\n")
		if len(lines) != 5 {
			t.Errorf("Expected 5 lines, got %d", len(lines))
		}
		for _, line := range lines {
			if len(line) != 10 {
				t.Errorf("Expected line width 10, got %d", len(line))
			}
		}
	})
}

// TestEdgeCases_ZeroDimensions tests handling of zero dimensions
func TestEdgeCases_ZeroDimensions(t *testing.T) {
	t.Run("Width(0)", func(t *testing.T) {
		s := NewStyle().Width(0)
		if s.width == nil || *s.width != 0 {
			t.Error("Width(0) should set width to 0")
		}
	})

	t.Run("Height(0)", func(t *testing.T) {
		s := NewStyle().Height(0)
		if s.height == nil || *s.height != 0 {
			t.Error("Height(0) should set height to 0")
		}
	})

	t.Run("Render with width 0", func(_ *testing.T) {
		s := NewStyle().Width(0).Align(Center)
		output := s.Render("test")
		// Width 0 should not cause panic
		_ = output
	})

	t.Run("Render with height 0", func(t *testing.T) {
		s := NewStyle().Height(0)
		output := s.Render("test\nlines")
		// Height 0 should truncate to empty
		if output != "" {
			t.Errorf("Expected empty string with height 0, got %q", output)
		}
	})

	t.Run("Place with zero dimensions", func(t *testing.T) {
		output1 := Place(0, 10, Left, Top, "content")
		if output1 != "" {
			t.Error("Place with width 0 should return empty string")
		}

		output2 := Place(10, 0, Left, Top, "content")
		if output2 != "" {
			t.Error("Place with height 0 should return empty string")
		}
	})
}

// TestEdgeCases_NegativeValues tests handling of negative values (should clamp to 0)
func TestEdgeCases_NegativeValues(t *testing.T) {
	tests := []struct {
		name     string
		setter   func(Style, int) Style
		getter   func(Style) *int
		input    int
		expected int
	}{
		{"Width negative", func(s Style, v int) Style { return s.Width(v) }, func(s Style) *int { return s.width }, -10, 0},
		{"Height negative", func(s Style, v int) Style { return s.Height(v) }, func(s Style) *int { return s.height }, -5, 0},
		{"MaxWidth negative", func(s Style, v int) Style { return s.MaxWidth(v) }, func(s Style) *int { return s.maxWidth }, -20, 0},
		{"MaxHeight negative", func(s Style, v int) Style { return s.MaxHeight(v) }, func(s Style) *int { return s.maxHeight }, -15, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := tt.setter(s, tt.input)
			if s2 == (Style{}) {
				t.Fatal("Setter returned empty Style")
			}
			val := tt.getter(s2)
			if val == nil {
				t.Fatal("Value not set")
			}
			if *val != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, *val)
			}
		})
	}
}

// TestEdgeCases_UnicodeAndANSI tests handling of Unicode and ANSI codes
func TestEdgeCases_UnicodeAndANSI(t *testing.T) {
	t.Run("CJK characters with width", func(t *testing.T) {
		s := NewStyle().Width(10).Align(Center)
		output := s.Render("你好") // 2 chars, 4 cells wide
		stripped := measure.StripANSI(output)
		width := measure.Width(stripped)
		if width != 10 {
			t.Errorf("Expected width 10, got %d", width)
		}
	})

	t.Run("Emoji with alignment", func(t *testing.T) {
		s := NewStyle().Width(10).Align(Left)
		output := s.Render("🎉")
		stripped := measure.StripANSI(output)
		width := measure.Width(stripped)
		if width != 10 {
			t.Errorf("Expected width 10, got %d", width)
		}
	})

	t.Run("Preserve ANSI codes in layout functions", func(t *testing.T) {
		red, _ := NewColor("red")
		styled := NewStyle().Foreground(red).Render("test")

		output := JoinHorizontal(Top, styled, "plain")
		if !strings.Contains(output, "\x1b[") {
			t.Error("ANSI codes should be preserved in JoinHorizontal")
		}

		output2 := JoinVertical(Left, styled, "plain")
		if !strings.Contains(output2, "\x1b[") {
			t.Error("ANSI codes should be preserved in JoinVertical")
		}
	})
}

// TestEdgeCases_ContentExceedsDimensions tests truncation and clipping
func TestEdgeCases_ContentExceedsDimensions(t *testing.T) {
	t.Run("Content exceeds width with border", func(t *testing.T) {
		s := NewStyle().Width(5).Border(NormalBorder())
		output := s.Render("This is a very long line that should be truncated")
		// Should not panic and should produce valid output
		if output == "" {
			t.Error("Expected non-empty output")
		}
	})

	t.Run("Content exceeds height", func(t *testing.T) {
		s := NewStyle().Height(2)
		output := s.Render("L1\nL2\nL3\nL4\nL5")
		lines := strings.Split(output, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines (truncated), got %d", len(lines))
		}
	})

	t.Run("Place with content larger than box", func(t *testing.T) {
		content := "Very long content that exceeds box width\nAnd multiple lines too"
		output := Place(10, 2, Left, Top, content)
		lines := strings.Split(output, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines, got %d", len(lines))
		}
		// Lines should be clipped to width
		for i, line := range lines {
			if len(line) > 10 {
				t.Errorf("Line %d exceeds width 10: got %d", i, len(line))
			}
		}
	})
}

// TestEdgeCases_MultilineEdgeCases tests edge cases with multi-line content
func TestEdgeCases_MultilineEdgeCases(t *testing.T) {
	t.Run("Single newline", func(t *testing.T) {
		s := NewStyle()
		output := s.Render("\n")
		lines := strings.Split(output, "\n")
		if len(lines) != 2 {
			t.Errorf("Single newline should produce 2 lines, got %d", len(lines))
		}
	})

	t.Run("Multiple consecutive newlines", func(t *testing.T) {
		s := NewStyle()
		output := s.Render("a\n\n\nb")
		lines := strings.Split(output, "\n")
		if len(lines) != 4 {
			t.Errorf("Expected 4 lines, got %d", len(lines))
		}
	})

	t.Run("Trailing newlines", func(t *testing.T) {
		s := NewStyle()
		output := s.Render("test\n\n")
		lines := strings.Split(output, "\n")
		// Should preserve trailing newlines
		if len(lines) < 2 {
			t.Errorf("Expected at least 2 lines with trailing newlines, got %d", len(lines))
		}
	})

	t.Run("Empty lines with padding", func(t *testing.T) {
		s := NewStyle().Padding(1)
		output := s.Render("line1\n\nline3")
		if output == "" {
			t.Error("Expected non-empty output")
		}
	})
}

// TestEdgeCases_BorderAndPaddingCombinations tests complex combinations
func TestEdgeCases_BorderAndPaddingCombinations(t *testing.T) {
	t.Run("Border with padding and alignment", func(t *testing.T) {
		s := NewStyle().
			Width(20).
			Height(5).
			Padding(1).
			Border(RoundedBorder()).
			Align(Center).
			AlignVertical(Center)

		output := s.Render("test")
		if output == "" {
			t.Error("Expected non-empty output")
		}

		// Should not panic and should produce coherent output
		lines := strings.Split(output, "\n")
		if len(lines) == 0 {
			t.Error("Expected multiple lines")
		}
	})

	t.Run("Border with zero padding", func(t *testing.T) {
		s := NewStyle().Padding(0).Border(ThickBorder())
		output := s.Render("test")
		if !strings.Contains(output, "test") {
			t.Error("Content should be present")
		}
	})

	t.Run("Partial borders", func(t *testing.T) {
		s := NewStyle().
			Border(NormalBorder()).
			BorderTop(true).
			BorderBottom(false).
			BorderLeft(false).
			BorderRight(false)

		output := s.Render("test")
		if output == "" {
			t.Error("Expected non-empty output")
		}
	})
}

// TestEdgeCases_ColorEdgeCases tests color-related edge cases
func TestEdgeCases_ColorEdgeCases(t *testing.T) {
	t.Run("Nil colors", func(t *testing.T) {
		s := NewStyle()
		// Foreground and background should be nil initially
		if s.foreground != nil {
			t.Error("Initial foreground should be nil")
		}
		if s.background != nil {
			t.Error("Initial background should be nil")
		}
	})

	t.Run("Background color with alignment", func(t *testing.T) {
		bg, _ := NewColor("blue")
		s := NewStyle().Width(20).Background(bg).Align(Center)
		output := s.Render("test")
		// Background should extend to full width
		if !strings.Contains(output, "\x1b[44m") {
			t.Error("Expected blue background ANSI code")
		}
	})

	t.Run("Background color with vertical alignment", func(t *testing.T) {
		bg, _ := NewColor("red")
		s := NewStyle().Height(5).Width(10).Background(bg).AlignVertical(Center)
		output := s.Render("test")
		// Background should be present in empty lines too
		if !strings.Contains(output, "\x1b[41m") {
			t.Error("Expected red background ANSI code")
		}
	})
}
