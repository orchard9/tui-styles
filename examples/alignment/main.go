// Package main demonstrates horizontal alignment features of tui-styles.
package main

import (
	"fmt"

	tuistyles "github.com/orchard9/tui-styles"
)

func main() {
	fmt.Println("=== Horizontal Alignment Demo ===")
	fmt.Println()

	// Demo left, center, right alignment
	width := 40

	alignments := []struct {
		name string
		pos  tuistyles.Position
	}{
		{"Left Alignment", tuistyles.Left},
		{"Center Alignment", tuistyles.Center},
		{"Right Alignment", tuistyles.Right},
	}

	for _, a := range alignments {
		s := tuistyles.NewStyle().
			Width(width).
			Align(a.pos).
			Border(tuistyles.NormalBorder()).
			Padding(1)

		fmt.Println(a.name + ":")
		fmt.Println(s.Render("Hello, World!"))
		fmt.Println()
	}

	// Demo alignment with colors
	fmt.Println()
	fmt.Println("=== Colored Alignment ===")
	fmt.Println()

	red, _ := tuistyles.NewColor("red")
	blue, _ := tuistyles.NewColor("blue")
	green, _ := tuistyles.NewColor("green")

	coloredStyles := []struct {
		name  string
		align tuistyles.Position
		fg    tuistyles.Color
		bg    tuistyles.Color
	}{
		{"Left + Red BG", tuistyles.Left, red, blue},
		{"Center + Green BG", tuistyles.Center, green, blue},
		{"Right + Blue BG", tuistyles.Right, blue, green},
	}

	for _, cs := range coloredStyles {
		s := tuistyles.NewStyle().
			Width(width).
			Align(cs.align).
			Foreground(cs.fg).
			Background(cs.bg).
			Border(tuistyles.RoundedBorder()).
			Padding(1)

		fmt.Println(cs.name + ":")
		fmt.Println(s.Render("Styled Text"))
		fmt.Println()
	}

	// Demo multi-line alignment
	fmt.Println()
	fmt.Println("=== Multi-line Alignment ===")
	fmt.Println()

	multiLine := `Short
Medium line
Very long line here
End`

	multiLineStyles := []struct {
		name  string
		align tuistyles.Position
	}{
		{"Left Aligned Lines", tuistyles.Left},
		{"Center Aligned Lines", tuistyles.Center},
		{"Right Aligned Lines", tuistyles.Right},
	}

	for _, mls := range multiLineStyles {
		s := tuistyles.NewStyle().
			Width(width).
			Align(mls.align).
			Border(tuistyles.DoubleBorder()).
			Padding(1)

		fmt.Println(mls.name + ":")
		fmt.Println(s.Render(multiLine))
		fmt.Println()
	}

	// Demo alignment grid (all 3 positions)
	fmt.Println()
	fmt.Println("=== Alignment Grid ===")
	fmt.Println()

	purple, _ := tuistyles.NewColor("#7D56F4")
	boxWidth := 25

	leftStyle := tuistyles.NewStyle().
		Width(boxWidth).
		Align(tuistyles.Left).
		Border(tuistyles.RoundedBorder()).
		BorderForeground(purple).
		Padding(1)

	centerStyle := tuistyles.NewStyle().
		Width(boxWidth).
		Align(tuistyles.Center).
		Border(tuistyles.RoundedBorder()).
		BorderForeground(purple).
		Padding(1)

	rightStyle := tuistyles.NewStyle().
		Width(boxWidth).
		Align(tuistyles.Right).
		Border(tuistyles.RoundedBorder()).
		BorderForeground(purple).
		Padding(1)

	fmt.Println(leftStyle.Render("Left"))
	fmt.Println(centerStyle.Render("Center"))
	fmt.Println(rightStyle.Render("Right"))

	// Demo alignment with Unicode
	fmt.Println()
	fmt.Println()
	fmt.Println("=== Unicode Content Alignment ===")
	fmt.Println()

	unicodeStyle := tuistyles.NewStyle().
		Width(30).
		Align(tuistyles.Center).
		Border(tuistyles.ThickBorder()).
		Padding(1)

	fmt.Println(unicodeStyle.Render("你好世界"))
	fmt.Println(unicodeStyle.Render("Hello 👋 World"))

	fmt.Println()
	fmt.Println("=== Demo Complete ===")
}
