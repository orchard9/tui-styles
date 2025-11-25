// Package main demonstrates complex rendering and layout compositions in tui-styles.
package main

import (
	"fmt"

	"github.com/orchard9/tui-styles"
)

func main() {
	fmt.Println("=== TUI Styles Rendering Showcase ===")
	fmt.Println()

	// Colors for demos
	red, _ := tuistyles.NewColor("red")
	blue, _ := tuistyles.NewColor("blue")
	green, _ := tuistyles.NewColor("green")
	yellow, _ := tuistyles.NewColor("yellow")
	cyan, _ := tuistyles.NewColor("cyan")

	// 1. Vertical Alignment Demo
	fmt.Println("1. Vertical Alignment (Height: 5)")
	fmt.Println("────────────────────────────────────")

	topAlign := tuistyles.NewStyle().
		Width(15).
		Height(5).
		AlignVertical(tuistyles.Top).
		Border(tuistyles.NormalBorder()).
		BorderForeground(red).
		Render("Top")

	centerAlign := tuistyles.NewStyle().
		Width(15).
		Height(5).
		AlignVertical(tuistyles.Center).
		Border(tuistyles.NormalBorder()).
		BorderForeground(blue).
		Render("Center")

	bottomAlign := tuistyles.NewStyle().
		Width(15).
		Height(5).
		AlignVertical(tuistyles.Bottom).
		Border(tuistyles.NormalBorder()).
		BorderForeground(green).
		Render("Bottom")

	alignRow := tuistyles.JoinHorizontal(tuistyles.Top, topAlign, "  ", centerAlign, "  ", bottomAlign)
	fmt.Println(alignRow)
	fmt.Println()

	// 2. JoinHorizontal Demo
	fmt.Println("2. JoinHorizontal (Top, Center, Bottom)")
	fmt.Println("────────────────────────────────────")

	box1 := tuistyles.NewStyle().
		Background(red).
		Padding(1).
		Render("Box 1\nShort")

	box2 := tuistyles.NewStyle().
		Background(blue).
		Padding(1).
		Render("Box 2\nMedium\nHeight")

	box3 := tuistyles.NewStyle().
		Background(green).
		Padding(1).
		Render("Box 3\nShort")

	topRow := tuistyles.JoinHorizontal(tuistyles.Top, box1, " ", box2, " ", box3)
	fmt.Println("Top Alignment:")
	fmt.Println(topRow)
	fmt.Println()

	centerRow := tuistyles.JoinHorizontal(tuistyles.Center, box1, " ", box2, " ", box3)
	fmt.Println("Center Alignment:")
	fmt.Println(centerRow)
	fmt.Println()

	// 3. JoinVertical Demo
	fmt.Println("3. JoinVertical (Left, Center, Right)")
	fmt.Println("────────────────────────────────────")

	row1 := tuistyles.NewStyle().
		Foreground(red).
		Render("Short")

	row2 := tuistyles.NewStyle().
		Foreground(blue).
		Render("Medium Length Row")

	row3 := tuistyles.NewStyle().
		Foreground(green).
		Render("Very Long Content Here")

	leftStack := tuistyles.JoinVertical(tuistyles.Left, row1, row2, row3)
	fmt.Println("Left Alignment:")
	fmt.Println(leftStack)
	fmt.Println()

	centerStack := tuistyles.JoinVertical(tuistyles.Center, row1, row2, row3)
	fmt.Println("Center Alignment:")
	fmt.Println(centerStack)
	fmt.Println()

	rightStack := tuistyles.JoinVertical(tuistyles.Right, row1, row2, row3)
	fmt.Println("Right Alignment:")
	fmt.Println(rightStack)
	fmt.Println()

	// 4. Place Demo
	fmt.Println("4. Place (9 Positions)")
	fmt.Println("────────────────────────────────────")

	boxStyle := tuistyles.NewStyle().
		Width(60).
		Height(7).
		Border(tuistyles.RoundedBorder()).
		BorderForeground(cyan)

	positions := []struct {
		h, v tuistyles.Position
		name string
	}{
		{tuistyles.Left, tuistyles.Top, "TL"},
		{tuistyles.Center, tuistyles.Top, "TC"},
		{tuistyles.Right, tuistyles.Top, "TR"},
		{tuistyles.Left, tuistyles.Center, "ML"},
		{tuistyles.Center, tuistyles.Center, "MC"},
		{tuistyles.Right, tuistyles.Center, "MR"},
		{tuistyles.Left, tuistyles.Bottom, "BL"},
		{tuistyles.Center, tuistyles.Bottom, "BC"},
		{tuistyles.Right, tuistyles.Bottom, "BR"},
	}

	for _, pos := range positions {
		content := tuistyles.NewStyle().
			Foreground(yellow).
			Bold(true).
			Render(pos.name)

		boxContent := tuistyles.Place(58, 5, pos.h, pos.v, content)
		box := boxStyle.Render(boxContent)

		fmt.Printf("Position: %s\n", pos.name)
		fmt.Println(box)
		fmt.Println()
	}

	// 5. Dashboard Layout Demo
	fmt.Println("5. Dashboard Layout (Complex Composition)")
	fmt.Println("────────────────────────────────────")

	headerStyle := tuistyles.NewStyle().
		Foreground(yellow).
		Bold(true).
		Width(70).
		Padding(1).
		Border(tuistyles.ThickBorder()).
		BorderForeground(yellow).
		Align(tuistyles.Center)

	cardStyle := tuistyles.NewStyle().
		Width(20).
		Height(8).
		Padding(1).
		Border(tuistyles.RoundedBorder()).
		BorderForeground(blue).
		AlignVertical(tuistyles.Center)

	footerStyle := tuistyles.NewStyle().
		Foreground(green).
		Width(70).
		Padding(1).
		Border(tuistyles.NormalBorder()).
		BorderForeground(green).
		Align(tuistyles.Right)

	header := headerStyle.Render("Dashboard Application")

	card1 := cardStyle.Render("CPU Usage\n───\n45%")
	card2 := cardStyle.Render("Memory\n───\n2.1 GB")
	card3 := cardStyle.Render("Disk I/O\n───\n125 MB/s")

	cardsRow := tuistyles.JoinHorizontal(tuistyles.Top, card1, " ", card2, " ", card3)

	footer := footerStyle.Render("Status: Ready | v1.0.0")

	dashboard := tuistyles.JoinVertical(tuistyles.Left, header, "", cardsRow, "", footer)

	fmt.Println(dashboard)
	fmt.Println()

	// 6. Mixed Alignment Demo
	fmt.Println("6. Mixed Alignment (Horizontal + Vertical)")
	fmt.Println("────────────────────────────────────")

	mixedStyle := tuistyles.NewStyle().
		Width(30).
		Height(10).
		Background(blue).
		Foreground(yellow).
		Padding(1).
		Border(tuistyles.DoubleBorder()).
		BorderForeground(cyan).
		Align(tuistyles.Center).
		AlignVertical(tuistyles.Center)

	mixed := mixedStyle.Render("Centered\nBoth Ways")
	fmt.Println(mixed)
	fmt.Println()

	fmt.Println("=== End of Showcase ===")
}
