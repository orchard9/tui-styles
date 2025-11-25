// Package main demonstrates border styles and colors in tui-styles.
package main

import (
	"fmt"

	tuistyles "github.com/orchard9/tui-styles"
)

func main() {
	fmt.Println("=== Border Styles Demo ===")
	fmt.Println()

	// Demo all 8 border types
	borders := []struct {
		name   string
		border tuistyles.Border
	}{
		{"Normal Border", tuistyles.NormalBorder()},
		{"Rounded Border", tuistyles.RoundedBorder()},
		{"Thick Border", tuistyles.ThickBorder()},
		{"Double Border", tuistyles.DoubleBorder()},
		{"Block Border", tuistyles.BlockBorder()},
		{"Outer Half Block", tuistyles.OuterHalfBlockBorder()},
		{"Inner Half Block", tuistyles.InnerHalfBlockBorder()},
		{"Hidden Border", tuistyles.HiddenBorder()},
	}

	for _, b := range borders {
		s := tuistyles.NewStyle().
			Border(b.border).
			Padding(1)

		fmt.Printf("%s:\n", b.name)
		fmt.Println(s.Render("Hello, World!"))
		fmt.Println()
	}

	// Demo colored borders
	fmt.Println()
	fmt.Println("=== Colored Borders ===")
	fmt.Println()

	red, _ := tuistyles.NewColor("red")
	blue, _ := tuistyles.NewColor("blue")
	green, _ := tuistyles.NewColor("green")
	yellow, _ := tuistyles.NewColor("yellow")

	coloredStyles := []struct {
		name  string
		style tuistyles.Style
	}{
		{
			"Red Border",
			tuistyles.NewStyle().
				Border(tuistyles.RoundedBorder()).
				BorderForeground(red).
				Padding(1),
		},
		{
			"Blue Background Border",
			tuistyles.NewStyle().
				Border(tuistyles.ThickBorder()).
				BorderBackground(blue).
				Padding(1),
		},
		{
			"Green Border + Yellow BG Content",
			tuistyles.NewStyle().
				Border(tuistyles.DoubleBorder()).
				BorderForeground(green).
				Background(yellow).
				Padding(1),
		},
	}

	for _, cs := range coloredStyles {
		fmt.Printf("%s:\n", cs.name)
		fmt.Println(cs.style.Render("Styled Border"))
		fmt.Println()
	}

	// Demo partial borders
	fmt.Println()
	fmt.Println("=== Partial Borders ===")
	fmt.Println()

	partialStyles := []struct {
		name  string
		style tuistyles.Style
	}{
		{
			"Top & Bottom Only",
			tuistyles.NewStyle().
				Border(tuistyles.NormalBorder(), true, false, true, false).
				Padding(1),
		},
		{
			"Left & Right Only",
			tuistyles.NewStyle().
				Border(tuistyles.NormalBorder(), false, true, false, true).
				Padding(1),
		},
		{
			"Top Only",
			tuistyles.NewStyle().
				Border(tuistyles.ThickBorder()).
				BorderTop(true).
				BorderRight(false).
				BorderBottom(false).
				BorderLeft(false).
				Padding(1),
		},
	}

	for _, ps := range partialStyles {
		fmt.Printf("%s:\n", ps.name)
		fmt.Println(ps.style.Render("Partial Border"))
		fmt.Println()
	}

	// Demo multi-line content
	fmt.Println()
	fmt.Println("=== Multi-line Content ===")
	fmt.Println()

	multiLine := `Line 1
Line 2
Line 3`

	multiLineStyle := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(blue).
		Padding(1, 2)

	fmt.Println(multiLineStyle.Render(multiLine))

	// Demo complex composition
	fmt.Println()
	fmt.Println("=== Complex Composition ===")
	fmt.Println()

	purple, _ := tuistyles.NewColor("#7D56F4")
	pink, _ := tuistyles.NewColor("#F72798")

	card := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(purple).
		Padding(2, 4).
		Bold(true).
		Foreground(pink)

	fmt.Println(card.Render("Beautiful Card"))

	// Demo with Unicode
	fmt.Println()
	fmt.Println("=== Unicode Content ===")
	fmt.Println()

	unicodeStyle := tuistyles.NewStyle().
		Border(tuistyles.DoubleBorder()).
		BorderForeground(green).
		Padding(1)

	fmt.Println(unicodeStyle.Render("你好世界 🎉"))

	fmt.Println("\n=== Demo Complete ===")
}
