package main

import (
	"fmt"

	"github.com/orchard9/tui-styles"
)

func main() {
	fmt.Println("TUI Styles - Basic Examples")
	fmt.Println("============================\n")

	// Example 1: Basic colors
	fmt.Println("1. Basic Colors:")
	redColor, _ := tuistyles.NewColor("red")
	fmt.Printf("   Red (ANSI name):    %sHello%s\n", redColor.ToANSI(), "\x1b[0m")

	hexColor, _ := tuistyles.NewColor("#FF5733")
	fmt.Printf("   Hex color:          %sWorld%s\n", hexColor.ToANSI(), "\x1b[0m")

	codeColor, _ := tuistyles.NewColor("196")
	fmt.Printf("   ANSI 256 code (196): %sTest%s\n", codeColor.ToANSI(), "\x1b[0m")

	fmt.Println()

	// Example 2: Adaptive colors
	fmt.Println("2. Adaptive Colors:")
	adaptiveColor, _ := tuistyles.NewAdaptiveColor("#000000", "#FFFFFF")
	selectedColor := adaptiveColor.ToColor()
	fmt.Printf("   Adaptive (based on terminal): %sAdaptive Text%s\n",
		selectedColor.ToANSI(), "\x1b[0m")
	fmt.Println()

	// Example 3: Position enum
	fmt.Println("3. Position Enum:")
	positions := []tuistyles.Position{
		tuistyles.Left,
		tuistyles.Center,
		tuistyles.Right,
		tuistyles.Top,
		tuistyles.Bottom,
	}
	for _, pos := range positions {
		fmt.Printf("   Position %s: IsHorizontal=%v, IsVertical=%v\n",
			pos.String(), pos.IsHorizontal(), pos.IsVertical())
	}
	fmt.Println()

	// Example 4: Border types
	fmt.Println("4. Border Types:")
	borders := []struct {
		name   string
		border tuistyles.Border
	}{
		{"Normal", tuistyles.NormalBorder()},
		{"Rounded", tuistyles.RoundedBorder()},
		{"Thick", tuistyles.ThickBorder()},
		{"Double", tuistyles.DoubleBorder()},
		{"Block", tuistyles.BlockBorder()},
	}

	for _, b := range borders {
		fmt.Printf("   %s Border:\n", b.name)
		fmt.Printf("     %s%s%s%s%s\n",
			b.border.TopLeft, b.border.Top, b.border.Top, b.border.Top, b.border.TopRight)
		fmt.Printf("     %s   %s\n", b.border.Left, b.border.Right)
		fmt.Printf("     %s%s%s%s%s\n",
			b.border.BottomLeft, b.border.Bottom, b.border.Bottom, b.border.Bottom, b.border.BottomRight)
		fmt.Println()
	}

	// Example 5: Background colors
	fmt.Println("5. Background Colors:")
	bgColor, _ := tuistyles.NewColor("blue")
	fgColor, _ := tuistyles.NewColor("white")
	fmt.Printf("   %s%s Background Text %s\n",
		bgColor.ToANSIBackground(), fgColor.ToANSI(), "\x1b[0m")

	fmt.Println("\nFor full API documentation, see: https://github.com/orchard9/tui-styles")
}
