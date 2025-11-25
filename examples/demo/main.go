package main

import (
	"fmt"
	"strings"

	tuistyles "github.com/orchard9/tui-styles"
)

func main() {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        TUI STYLES DEMO - v1.0.0                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝\n")

	// Section 1: Text Attributes
	printSection("1. TEXT ATTRIBUTES", demoTextAttributes())

	// Section 2: Colors
	printSection("2. COLORS", demoColors())

	// Section 3: Borders
	printSection("3. BORDER STYLES", demoBorders())

	// Section 4: Padding & Margins
	printSection("4. PADDING & MARGINS", demoPaddingMargin())

	// Section 5: Alignment
	printSection("5. ALIGNMENT", demoAlignment())

	// Section 6: Layout Composition
	printSection("6. LAYOUT COMPOSITION", demoLayout())

	// Section 7: Real-World Example - Dashboard
	printSection("7. DASHBOARD EXAMPLE", demoDashboard())

	fmt.Println("\n╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          TUI Styles - Complete Terminal Styling Library for Go             ║")
	fmt.Println("║                  github.com/orchard9/tui-styles                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝\n")
}

func printSection(title, content string) {
	header := tuistyles.NewStyle().
		Bold(true).
		Foreground(tuistyles.Color("#FFFF00")).
		Render(title)

	fmt.Println("\n" + header)
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println(content)
}

func demoTextAttributes() string {
	var parts []string

	// Bold
	bold := tuistyles.NewStyle().
		Bold(true).
		Render("Bold Text")
	parts = append(parts, bold)

	// Italic
	italic := tuistyles.NewStyle().
		Italic(true).
		Render("Italic Text")
	parts = append(parts, italic)

	// Underline
	underline := tuistyles.NewStyle().
		Underline(true).
		Render("Underlined Text")
	parts = append(parts, underline)

	// Strikethrough
	strikethrough := tuistyles.NewStyle().
		Strikethrough(true).
		Render("Strikethrough Text")
	parts = append(parts, strikethrough)

	// Faint
	faint := tuistyles.NewStyle().
		Faint(true).
		Render("Faint Text")
	parts = append(parts, faint)

	// Blink
	blink := tuistyles.NewStyle().
		Blink(true).
		Render("Blinking Text")
	parts = append(parts, blink)

	// Reverse
	reverse := tuistyles.NewStyle().
		Reverse(true).
		Render("Reversed Text")
	parts = append(parts, reverse)

	// Combined
	combined := tuistyles.NewStyle().
		Bold(true).
		Italic(true).
		Underline(true).
		Foreground(tuistyles.Color("#FF00FF")).
		Render("Combined Attributes")
	parts = append(parts, combined)

	return strings.Join(parts, "  •  ")
}

func demoColors() string {
	var rows []string

	// Hex colors
	hex1 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#FF0000")).
		Background(tuistyles.Color("#000000")).
		Padding(0, 2).
		Render("Hex Red")

	hex2 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#00FF00")).
		Background(tuistyles.Color("#000000")).
		Padding(0, 2).
		Render("Hex Green")

	hex3 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#0000FF")).
		Background(tuistyles.Color("#000000")).
		Padding(0, 2).
		Render("Hex Blue")

	rows = append(rows, "Hex Colors: "+tuistyles.JoinHorizontal(tuistyles.Top, hex1, " ", hex2, " ", hex3))

	// ANSI color names
	ansi1 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("cyan")).
		Background(tuistyles.Color("black")).
		Padding(0, 2).
		Render("ANSI Cyan")

	ansi2 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("magenta")).
		Background(tuistyles.Color("black")).
		Padding(0, 2).
		Render("ANSI Magenta")

	ansi3 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("yellow")).
		Background(tuistyles.Color("black")).
		Padding(0, 2).
		Render("ANSI Yellow")

	rows = append(rows, "ANSI Names: "+tuistyles.JoinHorizontal(tuistyles.Top, ansi1, " ", ansi2, " ", ansi3))

	// 256-color codes
	color1 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("214")). // Orange
		Padding(0, 1).
		Render("214")

	color2 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("51")). // Cyan
		Padding(0, 1).
		Render("51")

	color3 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("201")). // Pink
		Padding(0, 1).
		Render("201")

	rows = append(rows, "256-Color Codes: "+tuistyles.JoinHorizontal(tuistyles.Top, color1, " ", color2, " ", color3))

	// Background colors
	bg := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#FFFFFF")).
		Background(tuistyles.Color("#FF5555")).
		Padding(0, 2).
		Render("Colored Background")

	rows = append(rows, "\nBackground: "+bg)

	return strings.Join(rows, "\n")
}

func demoBorders() string {
	var borders []string

	borderTypes := []struct {
		name   string
		border tuistyles.Border
		color  tuistyles.Color
	}{
		{"Normal", tuistyles.NormalBorder(), tuistyles.Color("cyan")},
		{"Rounded", tuistyles.RoundedBorder(), tuistyles.Color("green")},
		{"Thick", tuistyles.ThickBorder(), tuistyles.Color("red")},
		{"Double", tuistyles.DoubleBorder(), tuistyles.Color("magenta")},
	}

	for _, bt := range borderTypes {
		box := tuistyles.NewStyle().
			Border(bt.border).
			BorderForeground(bt.color).
			Padding(1, 2).
			Width(12).
			Align(tuistyles.Center).
			Render(bt.name)
		borders = append(borders, box)
	}

	row1 := tuistyles.JoinHorizontal(tuistyles.Top, borders[0], " ", borders[1], " ", borders[2], " ", borders[3])

	// Second row with more exotic borders
	block := tuistyles.NewStyle().
		Border(tuistyles.BlockBorder()).
		BorderForeground(tuistyles.Color("yellow")).
		Padding(1, 2).
		Width(12).
		Align(tuistyles.Center).
		Render("Block")

	outer := tuistyles.NewStyle().
		Border(tuistyles.OuterHalfBlockBorder()).
		BorderForeground(tuistyles.Color("blue")).
		Padding(1, 2).
		Width(12).
		Align(tuistyles.Center).
		Render("Outer")

	inner := tuistyles.NewStyle().
		Border(tuistyles.InnerHalfBlockBorder()).
		BorderForeground(tuistyles.Color("green")).
		Padding(1, 2).
		Width(12).
		Align(tuistyles.Center).
		Render("Inner")

	hidden := tuistyles.NewStyle().
		Border(tuistyles.HiddenBorder()).
		Padding(1, 2).
		Width(12).
		Align(tuistyles.Center).
		Render("Hidden")

	row2 := tuistyles.JoinHorizontal(tuistyles.Top, block, " ", outer, " ", inner, " ", hidden)

	return row1 + "\n\n" + row2
}

func demoPaddingMargin() string {
	// No padding
	noPad := tuistyles.NewStyle().
		Border(tuistyles.NormalBorder()).
		BorderForeground(tuistyles.Color("cyan")).
		Render("No Padding")

	// With padding
	withPad := tuistyles.NewStyle().
		Border(tuistyles.NormalBorder()).
		BorderForeground(tuistyles.Color("green")).
		Padding(2).
		Render("Padding: 2")

	// Asymmetric padding
	asymPad := tuistyles.NewStyle().
		Border(tuistyles.NormalBorder()).
		BorderForeground(tuistyles.Color("magenta")).
		PaddingTop(1).
		PaddingBottom(1).
		PaddingLeft(3).
		PaddingRight(3).
		Render("Asymmetric")

	// With background color
	coloredPad := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("yellow")).
		Background(tuistyles.Color("#334455")).
		Foreground(tuistyles.Color("#FFFFFF")).
		Padding(1, 2).
		Render("Colored Padding")

	return tuistyles.JoinHorizontal(tuistyles.Top, noPad, "  ", withPad, "  ", asymPad, "  ", coloredPad)
}

func demoAlignment() string {
	var rows []string

	// Horizontal alignment
	width := 20

	left := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("cyan")).
		Width(width).
		Align(tuistyles.Left).
		Padding(1).
		Render("Left")

	center := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("green")).
		Width(width).
		Align(tuistyles.Center).
		Padding(1).
		Render("Center")

	right := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("magenta")).
		Width(width).
		Align(tuistyles.Right).
		Padding(1).
		Render("Right")

	rows = append(rows, "Horizontal: "+tuistyles.JoinHorizontal(tuistyles.Top, left, " ", center, " ", right))

	// Vertical alignment
	height := 6

	top := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("red")).
		Width(width).
		Height(height).
		AlignVertical(tuistyles.Top).
		Align(tuistyles.Center).
		Padding(1).
		Render("Top")

	vcenter := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("yellow")).
		Width(width).
		Height(height).
		AlignVertical(tuistyles.Center).
		Align(tuistyles.Center).
		Padding(1).
		Render("Center")

	bottom := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("blue")).
		Width(width).
		Height(height).
		AlignVertical(tuistyles.Bottom).
		Align(tuistyles.Center).
		Padding(1).
		Render("Bottom")

	rows = append(rows, "\nVertical:   "+tuistyles.JoinHorizontal(tuistyles.Top, top, " ", vcenter, " ", bottom))

	return strings.Join(rows, "\n")
}

func demoLayout() string {
	// JoinHorizontal demo
	box1 := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("cyan")).
		Padding(1).
		Width(15).
		Height(4).
		Render("Box 1\nLine 2")

	box2 := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("green")).
		Padding(1).
		Width(15).
		Height(6).
		Render("Box 2\nTaller\nBox")

	box3 := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("magenta")).
		Padding(1).
		Width(15).
		Height(5).
		Render("Box 3\nMedium")

	horizontal := tuistyles.JoinHorizontal(tuistyles.Center, box1, " ", box2, " ", box3)

	label1 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("yellow")).
		Render("JoinHorizontal (Center alignment):")

	// JoinVertical demo
	line1 := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("red")).
		Padding(0, 1).
		Render("Short")

	line2 := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("yellow")).
		Padding(0, 1).
		Render("Medium Length")

	line3 := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("blue")).
		Padding(0, 1).
		Render("Very Long Content Here")

	vertical := tuistyles.JoinVertical(tuistyles.Center, line1, line2, line3)

	label2 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("yellow")).
		Render("JoinVertical (Center alignment):")

	// Place demo
	content := tuistyles.NewStyle().
		Bold(true).
		Foreground(tuistyles.Color("#FF00FF")).
		Render("★ Placed ★")

	placed := tuistyles.Place(30, 5, tuistyles.Center, tuistyles.Center, content)
	placedBox := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("cyan")).
		Render(placed)

	label3 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("yellow")).
		Render("Place (30x5 box, center-center):")

	return label1 + "\n" + horizontal + "\n\n" + label2 + "\n" + vertical + "\n\n" + label3 + "\n" + placedBox
}

func demoDashboard() string {
	// Header
	header := tuistyles.NewStyle().
		Bold(true).
		Foreground(tuistyles.Color("#FFFFFF")).
		Background(tuistyles.Color("#5555FF")).
		Padding(1, 2).
		Width(76).
		Align(tuistyles.Center).
		Render("📊 System Dashboard 📊")

	// Status panel (left)
	statusContent := `Status: ✓ Online
Uptime: 99.9%
CPU: 45%
Memory: 2.1 GB
Disk: 450 GB`

	statusPanel := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("green")).
		Padding(1).
		Width(36).
		Height(8).
		Render(statusContent)

	// Metrics panel (right)
	metricsContent := `Active Users: 1,234
Requests/sec: 567
Avg Response: 45ms
Error Rate: 0.02%
Cache Hit: 94%`

	metricsPanel := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("cyan")).
		Padding(1).
		Width(36).
		Height(8).
		Render(metricsContent)

	panels := tuistyles.JoinHorizontal(tuistyles.Top, statusPanel, "  ", metricsPanel)

	// Alerts section
	alert1 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#00FF00")).
		Render("✓ Database: Healthy")

	alert2 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#00FF00")).
		Render("✓ API: Responding")

	alert3 := tuistyles.NewStyle().
		Foreground(tuistyles.Color("#FFFF00")).
		Render("⚠ Cache: High Load")

	alerts := tuistyles.JoinHorizontal(tuistyles.Top, alert1, "   ", alert2, "   ", alert3)

	alertsBox := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(tuistyles.Color("yellow")).
		Padding(1).
		Width(76).
		Render(alerts)

	// Footer
	footer := tuistyles.NewStyle().
		Foreground(tuistyles.Color("gray")).
		Padding(1, 0).
		Width(76).
		Align(tuistyles.Center).
		Render("Last updated: 2025-11-23 14:30:00 UTC  •  Press R to refresh")

	// Compose everything
	dashboard := tuistyles.JoinVertical(
		tuistyles.Center,
		header,
		"",
		panels,
		"",
		alertsBox,
		"",
		footer,
	)

	return dashboard
}
