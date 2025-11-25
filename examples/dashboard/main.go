// Package main demonstrates advanced layout composition with a dashboard UI.
package main

import (
	"fmt"

	tuistyles "github.com/orchard9/tui-styles"
)

func main() {
	// Define colors
	headerBG, _ := tuistyles.NewColor("#5555FF")
	cyan, _ := tuistyles.NewColor("cyan")
	green, _ := tuistyles.NewColor("green")
	red, _ := tuistyles.NewColor("red")
	yellow, _ := tuistyles.NewColor("yellow")
	gray, _ := tuistyles.NewColor("gray")
	white, _ := tuistyles.NewColor("#FFFFFF")

	// Header
	headerStyle := tuistyles.NewStyle().
		Bold(true).
		Foreground(white).
		Background(headerBG).
		Padding(1, 2).
		Width(82).
		Align(tuistyles.Center)

	header := headerStyle.Render("🎨 TUI Styles Dashboard - System Monitor")

	// Left panel: Metrics
	metricsTitle := tuistyles.NewStyle().
		Bold(true).
		Foreground(cyan).
		Render("📊 Metrics")

	metricsContent := `
Users:     1,234
Active:      567
Pending:     123
Inactive:    544

CPU:        45%
Memory:     67%
Disk:       82%
Network:    23%`

	leftPanel := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(cyan).
		Padding(1).
		Width(38).
		Height(14).
		Render(metricsTitle + metricsContent)

	// Right panel: Status
	statusTitle := tuistyles.NewStyle().
		Bold(true).
		Foreground(green).
		Render("✓ System Status")

	statusOK := tuistyles.NewStyle().Foreground(green).Render("✓")
	statusFail := tuistyles.NewStyle().Foreground(red).Render("✗")
	statusWarn := tuistyles.NewStyle().Foreground(yellow).Render("⚠")

	statusContent := fmt.Sprintf(`

Database:   %s OK
API Server: %s OK
Cache:      %s Down
Queue:      %s OK
Search:     %s Degraded
Storage:    %s OK

Load Avg:   0.45, 0.67, 0.82
Uptime:     47 days
`,
		statusOK, statusOK, statusFail,
		statusOK, statusWarn, statusOK)

	rightPanel := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(green).
		Padding(1).
		Width(38).
		Height(14).
		Render(statusTitle + statusContent)

	// Recent logs
	logsTitle := tuistyles.NewStyle().
		Bold(true).
		Foreground(yellow).
		Render("📜 Recent Logs")

	logsContent := `
2025-11-23 23:45:12 INFO  Server started
2025-11-23 23:45:15 INFO  Connected to DB
2025-11-23 23:45:20 WARN  High memory usage
2025-11-23 23:45:25 ERROR Cache timeout
2025-11-23 23:45:30 INFO  Cache reconnected`

	logsPanel := tuistyles.NewStyle().
		Border(tuistyles.RoundedBorder()).
		BorderForeground(yellow).
		Padding(1).
		Width(80).
		Render(logsTitle + "\n" + logsContent)

	// Footer
	footerStyle := tuistyles.NewStyle().
		Foreground(gray).
		Padding(1, 0).
		Width(82).
		Align(tuistyles.Center)

	footer := footerStyle.Render("Press Q to quit • Press R to refresh • Press H for help")

	// Compose dashboard
	panels := tuistyles.JoinHorizontal(tuistyles.Top, leftPanel, "  ", rightPanel)
	dashboard := tuistyles.JoinVertical(
		tuistyles.Center,
		header,
		"",
		panels,
		"",
		logsPanel,
		"",
		footer,
	)

	// Print dashboard
	fmt.Println(dashboard)
	fmt.Println()

	// Additional demo: nested boxes
	fmt.Println("=== Nested Box Composition ===")
	fmt.Println()

	innerBox := tuistyles.NewStyle().
		Foreground(white).
		Background(headerBG).
		Padding(1, 2).
		Render("Inner Content")

	middleBox := tuistyles.NewStyle().
		Border(tuistyles.DoubleBorder()).
		BorderForeground(cyan).
		Padding(2).
		Render(innerBox)

	outerBox := tuistyles.NewStyle().
		Border(tuistyles.ThickBorder()).
		BorderForeground(yellow).
		Padding(2).
		Render(middleBox)

	fmt.Println(outerBox)

	fmt.Println()
	fmt.Println("=== Demo Complete ===")
}
