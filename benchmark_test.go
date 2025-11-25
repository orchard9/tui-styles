package tuistyles

import "testing"

// BenchmarkStyleCopy benchmarks the cost of copying a Style struct.
func BenchmarkStyleCopy(b *testing.B) {
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s
	}
}

// BenchmarkBold benchmarks the Bold method (single field modification).
func BenchmarkBold(b *testing.B) {
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Bold(true)
	}
}

// BenchmarkMethodChain_Short benchmarks a short method chain (3 methods).
func BenchmarkMethodChain_Short(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewStyle().Bold(true).Italic(true).Underline(true)
	}
}

// BenchmarkMethodChain_Medium benchmarks a medium method chain (10 methods).
func BenchmarkMethodChain_Medium(b *testing.B) {
	red, _ := NewColor("red")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewStyle().
			Bold(true).
			Italic(true).
			Foreground(red).
			Width(80).
			Height(24).
			Padding(2).
			Margin(1).
			Border(RoundedBorder()).
			Align(Center).
			AlignVertical(Center)
	}
}

// BenchmarkMethodChain_Long benchmarks a long method chain (20+ methods).
func BenchmarkMethodChain_Long(b *testing.B) {
	black, _ := NewColor("black")
	white, _ := NewColor("white")
	gray, _ := NewColor("gray")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewStyle().
			Bold(true).
			Italic(false).
			Underline(true).
			Strikethrough(false).
			Faint(false).
			Foreground(black).
			Background(white).
			Width(80).
			Height(24).
			MaxWidth(100).
			MaxHeight(50).
			Align(Center).
			AlignVertical(Center).
			Padding(2, 4).
			Margin(1, 2).
			Border(RoundedBorder()).
			BorderForeground(gray).
			BorderTop(true).
			BorderRight(true).
			BorderBottom(true).
			BorderLeft(true)
	}
}

// BenchmarkPaddingShorthand benchmarks CSS shorthand padding.
func BenchmarkPaddingShorthand(b *testing.B) {
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Padding(2)
	}
}

// BenchmarkPaddingIndividual benchmarks individual padding methods.
func BenchmarkPaddingIndividual(b *testing.B) {
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.PaddingTop(2).PaddingRight(2).PaddingBottom(2).PaddingLeft(2)
	}
}

// BenchmarkBorderVariadic benchmarks Border with variadic arguments.
func BenchmarkBorderVariadic(b *testing.B) {
	rounded := RoundedBorder()
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Border(rounded, true, false, true, false)
	}
}

// BenchmarkColorMethods benchmarks color setting methods.
func BenchmarkColorMethods(b *testing.B) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Foreground(red).Background(blue)
	}
}

// BenchmarkLayoutMethods benchmarks layout methods.
func BenchmarkLayoutMethods(b *testing.B) {
	s := NewStyle()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Width(80).Height(24).Align(Center).AlignVertical(Center)
	}
}

// BenchmarkStyleReuse benchmarks creating derived styles from a base style.
func BenchmarkStyleReuse(b *testing.B) {
	white, _ := NewColor("white")
	blue, _ := NewColor("blue")
	red, _ := NewColor("red")

	baseButton := NewStyle().
		Bold(true).
		Foreground(white).
		Padding(1, 3).
		Border(RoundedBorder())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = baseButton.Background(blue)
		_ = baseButton.Background(red)
	}
}

// BenchmarkCompleteStyle benchmarks creating a fully styled element.
func BenchmarkCompleteStyle(b *testing.B) {
	white, _ := NewColor("white")
	blue, _ := NewColor("blue")
	gray, _ := NewColor("gray")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewStyle().
			Bold(true).
			Foreground(white).
			Background(blue).
			Width(80).
			Padding(2).
			Margin(1).
			Border(RoundedBorder()).
			BorderForeground(gray).
			Align(Center)
	}
}
