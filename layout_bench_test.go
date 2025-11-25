package tuistyles

import (
	"strings"
	"testing"
)

// BenchmarkJoinHorizontal benchmarks horizontal joining
func BenchmarkJoinHorizontal(b *testing.B) {
	tests := []struct {
		name string
		strs []string
	}{
		{
			name: "two short strings",
			strs: []string{"left", "right"},
		},
		{
			name: "two multi-line strings",
			strs: []string{"L1\nL2\nL3", "R1\nR2\nR3"},
		},
		{
			name: "five columns",
			strs: []string{"A", "B", "C", "D", "E"},
		},
		{
			name: "different heights",
			strs: []string{"Short", "Multi\nLine\nContent\nHere", "Medium\nLength"},
		},
		{
			name: "styled content",
			strs: func() []string {
				red, _ := NewColor("red")
				blue, _ := NewColor("blue")
				return []string{
					NewStyle().Foreground(red).Render("Left Column"),
					NewStyle().Foreground(blue).Render("Right Column"),
				}
			}(),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = JoinHorizontal(Top, tt.strs...)
			}
		})
	}
}

// BenchmarkJoinVertical benchmarks vertical stacking
func BenchmarkJoinVertical(b *testing.B) {
	tests := []struct {
		name string
		strs []string
	}{
		{
			name: "two short strings",
			strs: []string{"top", "bottom"},
		},
		{
			name: "five rows",
			strs: []string{"Row1", "Row2", "Row3", "Row4", "Row5"},
		},
		{
			name: "different widths",
			strs: []string{"Short", "Much longer content here", "Med"},
		},
		{
			name: "multi-line strings",
			strs: []string{"T1\nT2", "M1\nM2\nM3", "B1"},
		},
		{
			name: "styled content",
			strs: func() []string {
				red, _ := NewColor("red")
				blue, _ := NewColor("blue")
				green, _ := NewColor("green")
				return []string{
					NewStyle().Foreground(red).Render("Top Section"),
					NewStyle().Foreground(blue).Render("Middle Section"),
					NewStyle().Foreground(green).Render("Bottom Section"),
				}
			}(),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = JoinVertical(Left, tt.strs...)
			}
		})
	}
}

// BenchmarkPlace benchmarks content placement
func BenchmarkPlace(b *testing.B) {
	tests := []struct {
		name    string
		width   int
		height  int
		content string
	}{
		{
			name:    "small box",
			width:   10,
			height:  5,
			content: "X",
		},
		{
			name:    "medium box",
			width:   40,
			height:  10,
			content: "Centered Content",
		},
		{
			name:    "large box",
			width:   80,
			height:  24,
			content: "Title\nSubtitle\nContent",
		},
		{
			name:    "wide box",
			width:   120,
			height:  10,
			content: strings.Repeat("Wide content ", 5),
		},
		{
			name:    "tall box",
			width:   40,
			height:  50,
			content: strings.Repeat("Line\n", 10),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = Place(tt.width, tt.height, Center, Center, tt.content)
			}
		})
	}
}

// BenchmarkVerticalAlignment benchmarks vertical alignment performance
func BenchmarkVerticalAlignment(b *testing.B) {
	tests := []struct {
		name   string
		height int
		vAlign Position
		input  string
	}{
		{
			name:   "top align small",
			height: 5,
			vAlign: Top,
			input:  "test",
		},
		{
			name:   "center align medium",
			height: 10,
			vAlign: Center,
			input:  "Line1\nLine2\nLine3",
		},
		{
			name:   "bottom align large",
			height: 20,
			vAlign: Bottom,
			input:  strings.Repeat("Content\n", 5),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			s := NewStyle().Height(tt.height).Width(20).AlignVertical(tt.vAlign)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = s.Render(tt.input)
			}
		})
	}
}

// BenchmarkComplexLayout benchmarks realistic complex layouts
func BenchmarkComplexLayout(b *testing.B) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")
	green, _ := NewColor("green")

	headerStyle := NewStyle().
		Foreground(red).
		Bold(true).
		Width(40).
		Padding(1).
		Border(ThickBorder()).
		Align(Center)

	contentStyle := NewStyle().
		Foreground(blue).
		Width(40).
		Padding(2).
		Border(RoundedBorder())

	footerStyle := NewStyle().
		Foreground(green).
		Width(40).
		Padding(1).
		Border(NormalBorder()).
		Align(Right)

	b.Run("three section layout", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			header := headerStyle.Render("Application Title")
			content := contentStyle.Render("Main Content\nSecond Line\nThird Line")
			footer := footerStyle.Render("Status: Ready")

			_ = JoinVertical(Left, header, content, footer)
		}
	})

	b.Run("dashboard layout", func(b *testing.B) {
		cardStyle := NewStyle().
			Width(20).
			Height(8).
			Padding(1).
			Border(RoundedBorder()).
			AlignVertical(Center)

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			card1 := cardStyle.Render("Card 1\nMetric: 100")
			card2 := cardStyle.Render("Card 2\nMetric: 200")
			card3 := cardStyle.Render("Card 3\nMetric: 300")

			row := JoinHorizontal(Top, card1, card2, card3)
			_ = Place(80, 24, Center, Center, row)
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory efficiency
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("render many times", func(b *testing.B) {
		s := NewStyle().
			Width(40).
			Padding(1).
			Border(NormalBorder())

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = s.Render("test content")
		}
	})

	b.Run("join many strings", func(b *testing.B) {
		strs := make([]string, 10)
		for i := range strs {
			strs[i] = "Column"
		}

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = JoinHorizontal(Top, strs...)
		}
	})

	b.Run("place in large box", func(b *testing.B) {
		content := "Placed Content"

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Place(100, 50, Center, Center, content)
		}
	})
}
