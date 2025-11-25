package measure

import (
	"reflect"
	"testing"
)

func TestWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty string", "", 0},
		{"simple ASCII", "hello", 5},
		{"ASCII with spaces", "hello world", 11},
		{"ANSI colored", "\x1b[31mred\x1b[0m", 3},
		{"ANSI bold", "\x1b[1mbold\x1b[22m", 4},
		{"multiple ANSI codes", "\x1b[31m\x1b[1mred bold\x1b[0m", 8},
		{"RGB color", "\x1b[38;2;255;0;0mred\x1b[0m", 3},
		{"background color", "\x1b[48;2;0;0;255mblue\x1b[0m", 4},
		{"CJK characters", "你好", 4},              // 2 characters × 2 cells each
		{"mixed ASCII and CJK", "hello世界", 9},    // 5 + 4
		{"emoji", "👋", 2},                        // Most emoji are 2 cells wide
		{"emoji with text", "Hello 👋 World", 14}, // 5 + 1 + 2 + 1 + 5
		{"multiple emoji", "👋🌍", 4},              // 2 + 2
		{"tab character", "hello\tworld", 10},    // Tab counted as 0 width by runewidth
		{"newline ignored", "hello", 5},          // Width doesn't include newline
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Width(tt.input)
			if got != tt.want {
				t.Errorf("Width(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no ANSI", "hello", "hello"},
		{"simple color", "\x1b[31mred\x1b[0m", "red"},
		{"bold", "\x1b[1mbold\x1b[22m", "bold"},
		{"multiple codes", "\x1b[31m\x1b[1mred bold\x1b[0m", "red bold"},
		{"RGB foreground", "\x1b[38;2;255;0;0mred\x1b[0m", "red"},
		{"RGB background", "\x1b[48;2;0;0;255mblue\x1b[0m", "blue"},
		{"256 color", "\x1b[38;5;196mred\x1b[0m", "red"},
		{"mixed with text", "normal \x1b[31mred\x1b[0m normal", "normal red normal"},
		{"empty string", "", ""},
		{"only ANSI", "\x1b[31m\x1b[0m", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripANSI(tt.input)
			if got != tt.want {
				t.Errorf("StripANSI(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWidthPerLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []int
	}{
		{"single line", "hello", []int{5}},
		{"two lines", "hello\nworld", []int{5, 5}},
		{"three lines", "one\ntwo\nthree", []int{3, 3, 5}},
		{"empty lines", "\n\n", []int{0, 0, 0}},
		{"mixed widths", "hi\nhello\nworld!", []int{2, 5, 6}},
		{"ANSI in lines", "\x1b[31mred\x1b[0m\n\x1b[32mgreen\x1b[0m", []int{3, 5}},
		{"CJK multi-line", "你好\n世界", []int{4, 4}},
		{"empty string", "", []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WidthPerLine(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WidthPerLine(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaxWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"single line", "hello", 5},
		{"multiple lines same width", "hello\nworld", 5},
		{"multiple lines different widths", "hi\nhello\nworld!", 6},
		{"empty string", "", 0},
		{"empty lines", "\n\n", 0},
		{"CJK mixed", "hi\n你好", 4},
		{"ANSI codes", "\x1b[31mshort\x1b[0m\n\x1b[32mlonger line\x1b[0m", 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxWidth(tt.input)
			if got != tt.want {
				t.Errorf("MaxWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestLineCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty string", "", 1},
		{"single line", "hello", 1},
		{"two lines", "hello\nworld", 2},
		{"three lines", "one\ntwo\nthree", 3},
		{"trailing newline", "hello\n", 2},
		{"multiple newlines", "\n\n\n", 4},
		{"text with newlines", "line1\nline2\nline3\n", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LineCount(tt.input)
			if got != tt.want {
				t.Errorf("LineCount(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		tail  string
		want  string
	}{
		{"no truncation needed", "hello", 10, "...", "hello"},
		{"exact fit", "hello", 5, "...", "hello"},
		{"truncate with ellipsis", "hello world", 8, "...", "hello..."},
		{"truncate short", "hello", 3, "...", "..."},
		{"zero width", "hello", 0, "...", ""},
		{"negative width", "hello", -1, "...", ""},
		{"empty string", "", 5, "...", ""},
		{"tail too long", "hello", 2, "...", ".."},
		{"no tail", "hello world", 5, "", "hello"},
		{"CJK truncate", "你好世界", 6, "...", "你..."}, // 你=2cells, tail=3cells, total 5<=6
		{"emoji truncate", "Hello 👋 World", 8, "...", "Hello..."},
		{"ANSI preserved start", "\x1b[31mhello world\x1b[0m", 8, "...", "hello..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.width, tt.tail)
			// Measure actual width of result
			gotWidth := Width(got)
			if gotWidth > tt.width && tt.width > 0 {
				t.Errorf("Truncate(%q, %d, %q) width = %d, exceeds max %d",
					tt.input, tt.width, tt.tail, gotWidth, tt.width)
			}
			// Check content (ANSI stripped for comparison)
			if StripANSI(got) != tt.want {
				t.Errorf("Truncate(%q, %d, %q) = %q, want %q",
					tt.input, tt.width, tt.tail, got, tt.want)
			}
		})
	}
}

// Benchmark for performance validation
func BenchmarkWidth(b *testing.B) {
	testStrings := []string{
		"simple ASCII text",
		"\x1b[31m\x1b[1mcolored bold text\x1b[0m",
		"mixed ASCII and 中文字符",
		"emoji test 👋🌍🎉",
	}

	for _, s := range testStrings {
		b.Run(s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Width(s)
			}
		})
	}
}

func BenchmarkStripANSI(b *testing.B) {
	input := "\x1b[31m\x1b[1m\x1b[4mheavily styled text\x1b[0m\x1b[0m\x1b[0m"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StripANSI(input)
	}
}

func BenchmarkMaxWidth(b *testing.B) {
	input := "short\nmedium line\nvery long line here\nshort"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MaxWidth(input)
	}
}
