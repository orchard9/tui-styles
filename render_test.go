package tuistyles

import (
	"strings"
	"testing"

	"github.com/orchard9/tui-styles/internal/ansi"
	"github.com/orchard9/tui-styles/internal/measure"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name  string
		style func() Style
		input string
		want  string
	}{
		{
			name:  "empty string",
			style: NewStyle,
			input: "",
			want:  "",
		},
		{
			name:  "no style applied",
			style: NewStyle,
			input: "hello",
			want:  "hello",
		},
		{
			name:  "bold text",
			style: func() Style { return NewStyle().Bold(true) },
			input: "bold",
			want:  ansi.Bold() + "bold" + ansi.Reset(),
		},
		{
			name:  "italic text",
			style: func() Style { return NewStyle().Italic(true) },
			input: "italic",
			want:  ansi.Italic() + "italic" + ansi.Reset(),
		},
		{
			name:  "underline text",
			style: func() Style { return NewStyle().Underline(true) },
			input: "underline",
			want:  ansi.Underline() + "underline" + ansi.Reset(),
		},
		{
			name:  "strikethrough text",
			style: func() Style { return NewStyle().Strikethrough(true) },
			input: "strikethrough",
			want:  ansi.Strikethrough() + "strikethrough" + ansi.Reset(),
		},
		{
			name:  "faint text",
			style: func() Style { return NewStyle().Faint(true) },
			input: "faint",
			want:  ansi.Faint() + "faint" + ansi.Reset(),
		},
		{
			name:  "blink text",
			style: func() Style { return NewStyle().Blink(true) },
			input: "blink",
			want:  ansi.Blink() + "blink" + ansi.Reset(),
		},
		{
			name:  "reverse text",
			style: func() Style { return NewStyle().Reverse(true) },
			input: "reverse",
			want:  ansi.Reverse() + "reverse" + ansi.Reset(),
		},
		{
			name: "bold and italic combined",
			style: func() Style {
				return NewStyle().Bold(true).Italic(true)
			},
			input: "bold italic",
			want:  ansi.Bold() + ansi.Italic() + "bold italic" + ansi.Reset(),
		},
		{
			name: "all text attributes",
			style: func() Style {
				return NewStyle().
					Bold(true).
					Italic(true).
					Underline(true).
					Strikethrough(true)
			},
			input: "all attrs",
			want: ansi.Bold() + ansi.Italic() + ansi.Underline() +
				ansi.Strikethrough() + "all attrs" + ansi.Reset(),
		},
		{
			name: "foreground color red",
			style: func() Style {
				red, _ := NewColor("red")
				return NewStyle().Foreground(red)
			},
			input: "red text",
			want:  "\x1b[31mred text\x1b[0m",
		},
		{
			name: "foreground hex color",
			style: func() Style {
				purple, _ := NewColor("#7D56F4")
				return NewStyle().Foreground(purple)
			},
			input: "purple",
			want:  "\x1b[38;2;125;86;244mpurple\x1b[0m",
		},
		{
			name: "background color blue",
			style: func() Style {
				blue, _ := NewColor("blue")
				return NewStyle().Background(blue)
			},
			input: "blue bg",
			want:  "\x1b[44mblue bg\x1b[0m",
		},
		{
			name: "foreground and background",
			style: func() Style {
				red, _ := NewColor("red")
				blue, _ := NewColor("blue")
				return NewStyle().Foreground(red).Background(blue)
			},
			input: "red on blue",
			want:  "\x1b[31m\x1b[44mred on blue\x1b[0m",
		},
		{
			name: "bold red text",
			style: func() Style {
				red, _ := NewColor("red")
				return NewStyle().Bold(true).Foreground(red)
			},
			input: "bold red",
			want:  ansi.Bold() + "\x1b[31mbold red\x1b[0m",
		},
		{
			name: "complex styling",
			style: func() Style {
				fg, _ := NewColor("#FF0000")
				bg, _ := NewColor("#0000FF")
				return NewStyle().
					Bold(true).
					Italic(true).
					Foreground(fg).
					Background(bg)
			},
			input: "complex",
			want: ansi.Bold() + ansi.Italic() +
				"\x1b[38;2;255;0;0m\x1b[48;2;0;0;255mcomplex\x1b[0m",
		},
		{
			name: "multi-line text with per-line styling",
			style: func() Style {
				return NewStyle().Bold(true)
			},
			input: "line1\nline2\nline3",
			want: ansi.Bold() + "line1" + ansi.Reset() + "\n" +
				ansi.Bold() + "line2" + ansi.Reset() + "\n" +
				ansi.Bold() + "line3" + ansi.Reset(),
		},
		{
			name: "Unicode text",
			style: func() Style {
				return NewStyle().Bold(true)
			},
			input: "你好世界",
			want:  ansi.Bold() + "你好世界" + ansi.Reset(),
		},
		{
			name: "emoji text",
			style: func() Style {
				red, _ := NewColor("red")
				return NewStyle().Foreground(red)
			},
			input: "Hello 👋",
			want:  "\x1b[31mHello 👋\x1b[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			got := s.Render(tt.input)
			if got != tt.want {
				t.Errorf("Render(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRenderNoResetWhenNoStyle(t *testing.T) {
	// When no style is applied, no Reset should be added
	s := NewStyle()
	got := s.Render("plain")
	if strings.Contains(got, ansi.Reset()) {
		t.Errorf("Render with no style should not include Reset, got %q", got)
	}
}

func TestHasAnyStyle(t *testing.T) {
	tests := []struct {
		name  string
		style func() Style
		want  bool
	}{
		{"no style", NewStyle, false},
		{"bold set", func() Style { return NewStyle().Bold(true) }, true},
		{"italic set", func() Style { return NewStyle().Italic(true) }, true},
		{"foreground set", func() Style {
			c, _ := NewColor("red")
			return NewStyle().Foreground(c)
		}, true},
		{"background set", func() Style {
			c, _ := NewColor("blue")
			return NewStyle().Background(c)
		}, true},
		{"multiple attrs", func() Style {
			return NewStyle().Bold(true).Italic(true)
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			got := s.hasAnyStyle()
			if got != tt.want {
				t.Errorf("hasAnyStyle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	// String() should return empty for now
	s := NewStyle().Bold(true)
	got := s.String()
	if got != "" {
		t.Errorf("String() = %q, want empty string", got)
	}
}

func TestRenderMaxWidth(t *testing.T) {
	tests := []struct {
		name     string
		style    func() Style
		input    string
		wantText string // ANSI-stripped result
	}{
		{
			name:     "no truncation needed",
			style:    func() Style { return NewStyle().MaxWidth(20) },
			input:    "short",
			wantText: "short",
		},
		{
			name:     "truncate with ellipsis",
			style:    func() Style { return NewStyle().MaxWidth(10) },
			input:    "this is a very long line",
			wantText: "this is...",
		},
		{
			name: "multi-line with MaxWidth",
			style: func() Style {
				return NewStyle().Bold(true).MaxWidth(8)
			},
			input:    "line one\nshort\nvery long line here",
			wantText: "line one\nshort\nvery ...",
		},
		{
			name:     "MaxWidth with CJK",
			style:    func() Style { return NewStyle().MaxWidth(6) },
			input:    "你好世界",
			wantText: "你...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			got := s.Render(tt.input)

			// Strip ANSI to check text content
			stripped := stripANSIForTest(got)

			if stripped != tt.wantText {
				t.Errorf("Render(%q) text = %q, want %q", tt.input, stripped, tt.wantText)
			}
		})
	}
}

func TestRenderEmptyLines(t *testing.T) {
	s := NewStyle().Bold(true)
	got := s.Render("line1\n\nline3")

	// Should have three parts separated by newlines, with middle one empty
	parts := strings.Split(got, "\n")
	if len(parts) != 3 {
		t.Errorf("Expected 3 parts, got %d", len(parts))
	}

	// Middle line should be empty (no ANSI codes for empty content)
	if parts[1] != "" {
		t.Errorf("Middle line should be empty, got %q", parts[1])
	}
}

// Helper to strip ANSI codes for testing
func stripANSIForTest(s string) string {
	// Simple ANSI stripper for tests
	result := s
	result = strings.ReplaceAll(result, ansi.Bold(), "")
	result = strings.ReplaceAll(result, ansi.Italic(), "")
	result = strings.ReplaceAll(result, ansi.Underline(), "")
	result = strings.ReplaceAll(result, ansi.Faint(), "")
	result = strings.ReplaceAll(result, ansi.Blink(), "")
	result = strings.ReplaceAll(result, ansi.Reverse(), "")
	result = strings.ReplaceAll(result, ansi.Strikethrough(), "")
	result = strings.ReplaceAll(result, ansi.Reset(), "")

	// Remove color codes (simplified - just remove common patterns)
	for i := 0; i < 256; i++ {
		result = strings.ReplaceAll(result, ansi.ForegroundColor(string(rune('0'+i))), "")
		result = strings.ReplaceAll(result, ansi.BackgroundColor(string(rune('0'+i))), "")
	}

	// Remove RGB sequences more thoroughly
	for strings.Contains(result, "\x1b[") {
		start := strings.Index(result, "\x1b[")
		end := strings.Index(result[start:], "m")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}

	return result
}

func TestRenderPadding(t *testing.T) {
	tests := []struct {
		name    string
		style   func() Style
		input   string
		checkFn func(t *testing.T, output string)
	}{
		{
			name: "left padding only",
			style: func() Style {
				return NewStyle().PaddingLeft(2)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				if !strings.HasPrefix(stripped, "  hello") {
					t.Errorf("Expected 2 spaces before 'hello', got %q", stripped)
				}
			},
		},
		{
			name: "right padding only",
			style: func() Style {
				return NewStyle().PaddingRight(2)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				if !strings.HasSuffix(stripped, "hello  ") {
					t.Errorf("Expected 2 spaces after 'hello', got %q", stripped)
				}
			},
		},
		{
			name: "all sides padding",
			style: func() Style {
				return NewStyle().Padding(1)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(measure.StripANSI(output), "\n")
				if len(lines) != 3 {
					t.Errorf("Expected 3 lines (top padding, content, bottom padding), got %d", len(lines))
				}
				// Check middle line has left/right padding
				if !strings.Contains(lines[1], " hello ") {
					t.Errorf("Middle line should have left/right padding around 'hello', got %q", lines[1])
				}
			},
		},
		{
			name: "padding with background color",
			style: func() Style {
				bg, _ := NewColor("blue")
				return NewStyle().Background(bg).Padding(1)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// Should contain ANSI background color codes
				if !strings.Contains(output, "\x1b[44m") { // blue background
					t.Errorf("Expected blue background ANSI code in output")
				}
			},
		},
		{
			name: "multi-line with padding",
			style: func() Style {
				return NewStyle().Padding(1, 2)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(measure.StripANSI(output), "\n")
				// Should have: top padding, line1 with l/r padding, line2 with l/r padding, bottom padding
				if len(lines) < 4 {
					t.Errorf("Expected at least 4 lines, got %d", len(lines))
				}
			},
		},
		{
			name: "asymmetric padding",
			style: func() Style {
				return NewStyle().
					PaddingTop(1).
					PaddingRight(3).
					PaddingBottom(2).
					PaddingLeft(1)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(measure.StripANSI(output), "\n")
				// 1 top + 1 content + 2 bottom = 4 lines
				if len(lines) != 4 {
					t.Errorf("Expected 4 lines (1 top, 1 content, 2 bottom), got %d", len(lines))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			output := s.Render(tt.input)
			tt.checkFn(t, output)
		})
	}
}

func TestRenderBorder(t *testing.T) {
	tests := []struct {
		name    string
		style   func() Style
		input   string
		checkFn func(t *testing.T, output string)
	}{
		{
			name: "normal border all sides",
			style: func() Style {
				return NewStyle().Border(NormalBorder())
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				if len(lines) != 3 {
					t.Errorf("Expected 3 lines (top border, content, bottom border), got %d", len(lines))
				}
				stripped := stripANSIForTest(output)
				if !strings.Contains(stripped, "┌") {
					t.Errorf("Expected top-left corner '┌' in output")
				}
				if !strings.Contains(stripped, "┐") {
					t.Errorf("Expected top-right corner '┐' in output")
				}
				if !strings.Contains(stripped, "└") {
					t.Errorf("Expected bottom-left corner '└' in output")
				}
				if !strings.Contains(stripped, "┘") {
					t.Errorf("Expected bottom-right corner '┘' in output")
				}
				if !strings.Contains(stripped, "│") {
					t.Errorf("Expected vertical border '│' in output")
				}
			},
		},
		{
			name: "rounded border",
			style: func() Style {
				return NewStyle().Border(RoundedBorder())
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				stripped := stripANSIForTest(output)
				if !strings.Contains(stripped, "╭") {
					t.Errorf("Expected rounded top-left corner '╭'")
				}
				if !strings.Contains(stripped, "╮") {
					t.Errorf("Expected rounded top-right corner '╮'")
				}
				if !strings.Contains(stripped, "╰") {
					t.Errorf("Expected rounded bottom-left corner '╰'")
				}
				if !strings.Contains(stripped, "╯") {
					t.Errorf("Expected rounded bottom-right corner '╯'")
				}
			},
		},
		{
			name: "thick border",
			style: func() Style {
				return NewStyle().Border(ThickBorder())
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				stripped := stripANSIForTest(output)
				if !strings.Contains(stripped, "┏") {
					t.Errorf("Expected thick top-left corner '┏'")
				}
				if !strings.Contains(stripped, "┃") {
					t.Errorf("Expected thick vertical '┃'")
				}
				if !strings.Contains(stripped, "━") {
					t.Errorf("Expected thick horizontal '━'")
				}
			},
		},
		{
			name: "double border",
			style: func() Style {
				return NewStyle().Border(DoubleBorder())
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				stripped := stripANSIForTest(output)
				if !strings.Contains(stripped, "╔") {
					t.Errorf("Expected double top-left corner '╔'")
				}
				if !strings.Contains(stripped, "║") {
					t.Errorf("Expected double vertical '║'")
				}
				if !strings.Contains(stripped, "═") {
					t.Errorf("Expected double horizontal '═'")
				}
			},
		},
		{
			name: "multi-line with border",
			style: func() Style {
				return NewStyle().Border(NormalBorder())
			},
			input: "line1\nline2\nline3",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				// Top border + 3 content lines + bottom border = 5 lines
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines, got %d", len(lines))
				}
			},
		},
		{
			name: "border with foreground color",
			style: func() Style {
				red, _ := NewColor("red")
				return NewStyle().Border(NormalBorder()).BorderForeground(red)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// Should contain red foreground ANSI code
				if !strings.Contains(output, "\x1b[31m") {
					t.Errorf("Expected red foreground ANSI code in border")
				}
			},
		},
		{
			name: "border with background color",
			style: func() Style {
				blue, _ := NewColor("blue")
				return NewStyle().Border(NormalBorder()).BorderBackground(blue)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// Should contain blue background ANSI code
				if !strings.Contains(output, "\x1b[44m") {
					t.Errorf("Expected blue background ANSI code in border")
				}
			},
		},
		{
			name: "partial border - top and bottom only",
			style: func() Style {
				return NewStyle().Border(NormalBorder(), true, false, true, false)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				stripped := stripANSIForTest(output)
				lines := strings.Split(stripped, "\n")
				// Should have horizontal lines but no vertical bars around content
				if len(lines) != 3 {
					t.Errorf("Expected 3 lines, got %d", len(lines))
				}
				// Check that content line doesn't start/end with vertical bar
				if strings.HasPrefix(lines[1], "│") || strings.HasSuffix(lines[1], "│") {
					t.Errorf("Content line should not have vertical borders, got %q", lines[1])
				}
			},
		},
		{
			name: "partial border - left and right only",
			style: func() Style {
				return NewStyle().Border(NormalBorder(), false, true, false, true)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				// Should only have 1 line (content with left/right borders, no top/bottom)
				if len(lines) != 1 {
					t.Errorf("Expected 1 line (no top/bottom borders), got %d", len(lines))
				}
			},
		},
		{
			name: "border with padding",
			style: func() Style {
				return NewStyle().Padding(1).Border(NormalBorder())
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// Border should wrap padding
				// Padding adds space, border wraps the padded content
				lines := strings.Split(output, "\n")
				// Top border + top padding + content + bottom padding + bottom border = 5 lines
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines (border wraps padding), got %d", len(lines))
				}
			},
		},
		{
			name: "border without any sides enabled",
			style: func() Style {
				return NewStyle().Border(NormalBorder(), false)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// No border should be rendered
				stripped := stripANSIForTest(output)
				if stripped != "test" {
					t.Errorf("Expected just 'test' with no borders, got %q", stripped)
				}
			},
		},
		{
			name: "empty content with border",
			style: func() Style {
				return NewStyle().Border(NormalBorder())
			},
			input: "",
			checkFn: func(t *testing.T, output string) {
				// Should render border around empty content
				if output == "" {
					t.Errorf("Expected border to be rendered even with empty content")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			output := s.Render(tt.input)
			tt.checkFn(t, output)
		})
	}
}

func TestBorderTypes(t *testing.T) {
	// Test all border types render correctly
	borders := []struct {
		name   string
		border Border
	}{
		{"NormalBorder", NormalBorder()},
		{"RoundedBorder", RoundedBorder()},
		{"ThickBorder", ThickBorder()},
		{"DoubleBorder", DoubleBorder()},
		{"BlockBorder", BlockBorder()},
		{"OuterHalfBlockBorder", OuterHalfBlockBorder()},
		{"InnerHalfBlockBorder", InnerHalfBlockBorder()},
		{"HiddenBorder", HiddenBorder()},
	}

	for _, b := range borders {
		t.Run(b.name, func(t *testing.T) {
			s := NewStyle().Border(b.border)
			output := s.Render("test")

			// Just verify it renders without panic and produces output
			if output == "" {
				t.Errorf("Border %s produced empty output", b.name)
			}

			// Verify structure (3 lines: top, content, bottom)
			lines := strings.Split(output, "\n")
			if len(lines) != 3 {
				t.Errorf("Border %s produced %d lines, expected 3", b.name, len(lines))
			}
		})
	}
}

func TestBorderWithUnicode(t *testing.T) {
	s := NewStyle().Border(RoundedBorder())

	tests := []struct {
		name  string
		input string
	}{
		{"CJK characters", "你好世界"},
		{"emoji", "Hello 👋 World"},
		{"mixed", "Test 测试 🎉"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := s.Render(tt.input)

			// Verify border is properly sized for Unicode content
			lines := strings.Split(output, "\n")
			if len(lines) != 3 {
				t.Errorf("Expected 3 lines, got %d", len(lines))
			}

			// Measure widths (ANSI-aware)
			stripped := stripANSIForTest(output)
			strippedLines := strings.Split(stripped, "\n")

			topWidth := measure.Width(strippedLines[0])
			contentWidth := measure.Width(strippedLines[1])
			bottomWidth := measure.Width(strippedLines[2])

			// All lines should have same width
			if topWidth != contentWidth || contentWidth != bottomWidth {
				t.Errorf("Border lines have inconsistent widths: top=%d, content=%d, bottom=%d",
					topWidth, contentWidth, bottomWidth)
			}
		})
	}
}

func TestBorderOverrideIndividualEdges(t *testing.T) {
	s := NewStyle().
		Border(NormalBorder()). // All sides enabled
		BorderTop(false)        // Disable top

	output := s.Render("test")
	lines := strings.Split(output, "\n")

	// Should have 2 lines: content with borders + bottom border
	// (no top border)
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines (no top border), got %d: %v", len(lines), lines)
	}
}

func TestHorizontalAlignment(t *testing.T) {
	tests := []struct {
		name    string
		style   func() Style
		input   string
		checkFn func(t *testing.T, output string)
	}{
		{
			name: "left alignment",
			style: func() Style {
				return NewStyle().Width(20).Align(Left)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				// Should be "hello               " (5 chars + 15 spaces = 20 total)
				if !strings.HasPrefix(stripped, "hello") {
					t.Errorf("Left aligned text should start with 'hello', got %q", stripped)
				}
				width := measure.Width(stripped)
				if width != 20 {
					t.Errorf("Expected width 20, got %d", width)
				}
			},
		},
		{
			name: "center alignment",
			style: func() Style {
				return NewStyle().Width(20).Align(Center)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				// "hello" is 5 chars, padding is 15, split as 7 left + 8 right
				// Result: "       hello        " (7 spaces + "hello" + 8 spaces = 20)
				if !strings.Contains(stripped, "hello") {
					t.Errorf("Output should contain 'hello', got %q", stripped)
				}
				width := measure.Width(stripped)
				if width != 20 {
					t.Errorf("Expected width 20, got %d", width)
				}
				// Check that "hello" is roughly centered
				idx := strings.Index(stripped, "hello")
				if idx < 5 || idx > 10 {
					t.Errorf("Text should be centered, but starts at index %d in %q", idx, stripped)
				}
			},
		},
		{
			name: "right alignment",
			style: func() Style {
				return NewStyle().Width(20).Align(Right)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				// Should be "               hello" (15 spaces + 5 chars = 20 total)
				if !strings.HasSuffix(stripped, "hello") {
					t.Errorf("Right aligned text should end with 'hello', got %q", stripped)
				}
				width := measure.Width(stripped)
				if width != 20 {
					t.Errorf("Expected width 20, got %d", width)
				}
			},
		},
		{
			name: "left alignment with multi-line",
			style: func() Style {
				return NewStyle().Width(15).Align(Left)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 2 {
					t.Errorf("Expected 2 lines, got %d", len(lines))
				}
				// Both lines should be left-aligned and 15 chars wide
				for i, line := range lines {
					if measure.Width(line) != 15 {
						t.Errorf("Line %d width expected 15, got %d: %q", i, measure.Width(line), line)
					}
				}
			},
		},
		{
			name: "center alignment with multi-line",
			style: func() Style {
				return NewStyle().Width(20).Align(Center)
			},
			input: "short\nlonger line",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 2 {
					t.Errorf("Expected 2 lines, got %d", len(lines))
				}
				// Both lines should be 20 chars wide
				for i, line := range lines {
					if measure.Width(line) != 20 {
						t.Errorf("Line %d width expected 20, got %d: %q", i, measure.Width(line), line)
					}
				}
			},
		},
		{
			name: "alignment with background color",
			style: func() Style {
				bg, _ := NewColor("blue")
				return NewStyle().Width(20).Align(Center).Background(bg)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// Should contain blue background ANSI code
				if !strings.Contains(output, "\x1b[44m") {
					t.Errorf("Expected blue background ANSI code")
				}
				// Alignment padding should also have blue background
				stripped := measure.StripANSI(output)
				width := measure.Width(stripped)
				if width != 20 {
					t.Errorf("Expected width 20, got %d", width)
				}
			},
		},
		{
			name: "alignment with width exactly matching content",
			style: func() Style {
				return NewStyle().Width(5).Align(Center)
			},
			input: "hello",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				// No padding needed, should be exactly "hello"
				if stripped != "hello" {
					t.Errorf("Expected 'hello', got %q", stripped)
				}
			},
		},
		{
			name: "alignment with content exceeding width",
			style: func() Style {
				return NewStyle().Width(3).Align(Center)
			},
			input: "hello world",
			checkFn: func(t *testing.T, output string) {
				// Content exceeds width, should keep as is (no truncation by alignment)
				stripped := measure.StripANSI(output)
				if stripped != "hello world" {
					t.Errorf("Content should be kept as is when exceeding width, got %q", stripped)
				}
			},
		},
		{
			name: "alignment with Unicode content",
			style: func() Style {
				return NewStyle().Width(20).Align(Center)
			},
			input: "你好",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				width := measure.Width(stripped)
				if width != 20 {
					t.Errorf("Expected width 20 for Unicode content, got %d", width)
				}
				if !strings.Contains(stripped, "你好") {
					t.Errorf("Output should contain '你好', got %q", stripped)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			output := s.Render(tt.input)
			tt.checkFn(t, output)
		})
	}
}

func TestAlignmentWithPadding(t *testing.T) {
	// Alignment should work correctly with padding
	s := NewStyle().
		Width(20).
		Align(Center).
		Padding(1)

	output := s.Render("test")
	stripped := measure.StripANSI(output)

	// With padding(1), we have:
	// - Top padding line
	// - Content line (with left/right padding + alignment)
	// - Bottom padding line
	lines := strings.Split(stripped, "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines with padding, got %d", len(lines))
	}

	// Each line should be same width (20 + 2 for left/right padding = 22)
	for i, line := range lines {
		width := measure.Width(line)
		if width != 22 {
			t.Errorf("Line %d: expected width 22, got %d: %q", i, width, line)
		}
	}
}

func TestAlignmentWithBorder(t *testing.T) {
	// Alignment should work correctly with borders
	s := NewStyle().
		Width(15).
		Align(Center).
		Border(NormalBorder())

	output := s.Render("hello")
	stripped := measure.StripANSI(output)

	lines := strings.Split(stripped, "\n")
	// Top border + content + bottom border = 3 lines
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines with border, got %d", len(lines))
	}

	// Content line should have borders and be centered within width
	contentLine := lines[1]
	// Should be: "│     hello     │" (borders + centered "hello" in 15 chars)
	if !strings.Contains(contentLine, "hello") {
		t.Errorf("Content line should contain 'hello', got %q", contentLine)
	}
	if !strings.Contains(contentLine, "│") {
		t.Errorf("Content line should have vertical borders, got %q", contentLine)
	}
}

// Benchmark rendering performance
func BenchmarkRender(b *testing.B) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	benchmarks := []struct {
		name  string
		style Style
		input string
	}{
		{
			name:  "plain text",
			style: NewStyle(),
			input: "hello world",
		},
		{
			name:  "bold text",
			style: NewStyle().Bold(true),
			input: "hello world",
		},
		{
			name:  "colored text",
			style: NewStyle().Foreground(red),
			input: "hello world",
		},
		{
			name:  "complex style",
			style: NewStyle().Bold(true).Italic(true).Foreground(red).Background(blue),
			input: "hello world with complex styling",
		},
		{
			name:  "multi-line",
			style: NewStyle().Bold(true),
			input: "line1\nline2\nline3\nline4\nline5",
		},
		{
			name:  "border simple",
			style: NewStyle().Border(NormalBorder()),
			input: "hello",
		},
		{
			name:  "border with color",
			style: NewStyle().Border(RoundedBorder()).BorderForeground(red),
			input: "styled border",
		},
		{
			name:  "border with padding",
			style: NewStyle().Padding(2).Border(ThickBorder()),
			input: "padded and bordered",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				bm.style.Render(bm.input)
			}
		})
	}
}

// TestVerticalAlignment tests vertical alignment and height handling
func TestVerticalAlignment(t *testing.T) {
	tests := []struct {
		name    string
		style   func() Style
		input   string
		checkFn func(*testing.T, string)
	}{
		{
			name: "top alignment (default)",
			style: func() Style {
				return NewStyle().Height(5).Width(10)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines, got %d", len(lines))
				}
				// Content should be at top
				if !strings.Contains(lines[0], "line1") {
					t.Errorf("First line should contain 'line1', got %q", lines[0])
				}
				if !strings.Contains(lines[1], "line2") {
					t.Errorf("Second line should contain 'line2', got %q", lines[1])
				}
				// Remaining lines should be empty (spaces)
				for i := 2; i < 5; i++ {
					if measure.Width(lines[i]) != 10 {
						t.Errorf("Empty line %d should be width 10, got %d", i, measure.Width(lines[i]))
					}
				}
			},
		},
		{
			name: "top alignment explicit",
			style: func() Style {
				return NewStyle().Height(5).Width(10).AlignVertical(Top)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines, got %d", len(lines))
				}
				// Content should be at top
				if !strings.Contains(lines[0], "line1") {
					t.Errorf("First line should contain 'line1', got %q", lines[0])
				}
			},
		},
		{
			name: "center alignment",
			style: func() Style {
				return NewStyle().Height(5).Width(10).AlignVertical(Center)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines, got %d", len(lines))
				}
				// With 2 lines of content and 5 total, padding = 3
				// Center: top=1, bottom=2
				// So content should be at lines[1] and lines[2]
				if !strings.Contains(lines[1], "line1") {
					t.Errorf("Line 1 should contain 'line1', got %q", lines[1])
				}
				if !strings.Contains(lines[2], "line2") {
					t.Errorf("Line 2 should contain 'line2', got %q", lines[2])
				}
			},
		},
		{
			name: "bottom alignment",
			style: func() Style {
				return NewStyle().Height(5).Width(10).AlignVertical(Bottom)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines, got %d", len(lines))
				}
				// Content should be at bottom (lines[3] and lines[4])
				if !strings.Contains(lines[3], "line1") {
					t.Errorf("Line 3 should contain 'line1', got %q", lines[3])
				}
				if !strings.Contains(lines[4], "line2") {
					t.Errorf("Line 4 should contain 'line2', got %q", lines[4])
				}
			},
		},
		{
			name: "truncate when content exceeds height",
			style: func() Style {
				return NewStyle().Height(2).Width(10)
			},
			input: "line1\nline2\nline3\nline4",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 2 {
					t.Errorf("Expected 2 lines (truncated), got %d", len(lines))
				}
				if !strings.Contains(lines[0], "line1") {
					t.Errorf("First line should contain 'line1', got %q", lines[0])
				}
				if !strings.Contains(lines[1], "line2") {
					t.Errorf("Second line should contain 'line2', got %q", lines[1])
				}
			},
		},
		{
			name: "vertical alignment with background color",
			style: func() Style {
				bg, _ := NewColor("red")
				return NewStyle().Height(3).Width(10).AlignVertical(Center).Background(bg)
			},
			input: "test",
			checkFn: func(t *testing.T, output string) {
				// Should contain red background ANSI code
				if !strings.Contains(output, "\x1b[41m") {
					t.Errorf("Expected red background ANSI code")
				}
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 3 {
					t.Errorf("Expected 3 lines, got %d", len(lines))
				}
				// All lines should be width 10
				for i, line := range lines {
					if measure.Width(line) != 10 {
						t.Errorf("Line %d width expected 10, got %d", i, measure.Width(line))
					}
				}
			},
		},
		{
			name: "height exactly matching content",
			style: func() Style {
				return NewStyle().Height(2).Width(10)
			},
			input: "line1\nline2",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				// Exact match, no padding needed
				if len(lines) != 2 {
					t.Errorf("Expected 2 lines, got %d", len(lines))
				}
			},
		},
		{
			name: "vertical alignment without width (auto-detect content width)",
			style: func() Style {
				return NewStyle().Height(4).AlignVertical(Center)
			},
			input: "short\nmedium line",
			checkFn: func(t *testing.T, output string) {
				stripped := measure.StripANSI(output)
				lines := strings.Split(stripped, "\n")
				if len(lines) != 4 {
					t.Errorf("Expected 4 lines, got %d", len(lines))
				}
				// Should detect width from longest line ("medium line" = 11 chars)
				maxWidth := 0
				for _, line := range lines {
					w := measure.Width(line)
					if w > maxWidth {
						maxWidth = w
					}
				}
				// Empty lines should match the content width
				if maxWidth != 11 {
					t.Errorf("Expected max width 11, got %d", maxWidth)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.style()
			output := s.Render(tt.input)
			tt.checkFn(t, output)
		})
	}
}
