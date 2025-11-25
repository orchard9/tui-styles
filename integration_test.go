package tuistyles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIntegration_ButtonStyle demonstrates a realistic button styling scenario.
func TestIntegration_ButtonStyle(t *testing.T) {
	white, _ := NewColor("white")
	blue, _ := NewColor("blue")

	button := NewStyle().
		Bold(true).
		Foreground(white).
		Background(blue).
		Padding(1, 3). // Vertical: 1, Horizontal: 3
		Border(RoundedBorder()).
		Align(Center)

	// Verify all attributes set correctly
	require.NotNil(t, button.bold)
	require.True(t, *button.bold)

	require.NotNil(t, button.foreground)
	require.Equal(t, white, *button.foreground)

	require.NotNil(t, button.background)
	require.Equal(t, blue, *button.background)

	require.Equal(t, 1, *button.paddingTop)
	require.Equal(t, 3, *button.paddingLeft)
	require.Equal(t, 1, *button.paddingBottom)
	require.Equal(t, 3, *button.paddingRight)

	require.NotNil(t, button.borderType)
	require.NotNil(t, button.align)
	require.Equal(t, Center, *button.align)
}

// TestIntegration_CardStyle demonstrates a card container styling scenario.
func TestIntegration_CardStyle(t *testing.T) {
	gray, _ := NewColor("gray")

	card := NewStyle().
		Width(80).
		Padding(2).
		Margin(1).
		Border(NormalBorder()).
		BorderForeground(gray)

	require.Equal(t, 80, *card.width)
	require.Equal(t, 2, *card.paddingTop)
	require.Equal(t, 2, *card.paddingRight)
	require.Equal(t, 2, *card.paddingBottom)
	require.Equal(t, 2, *card.paddingLeft)

	require.Equal(t, 1, *card.marginTop)
	require.Equal(t, 1, *card.marginRight)
	require.Equal(t, 1, *card.marginBottom)
	require.Equal(t, 1, *card.marginLeft)

	require.NotNil(t, card.borderType)
	require.Equal(t, gray, *card.borderForeground)
}

// TestIntegration_HeaderStyle demonstrates a header styling scenario.
func TestIntegration_HeaderStyle(t *testing.T) {
	white, _ := NewColor("white")
	darkblue, _ := NewColor("blue")

	header := NewStyle().
		Bold(true).
		Foreground(white).
		Background(darkblue).
		Width(100).
		Padding(1, 2).
		BorderBottom(true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	require.True(t, *header.bold)
	require.Equal(t, white, *header.foreground)
	require.Equal(t, darkblue, *header.background)
	require.Equal(t, 100, *header.width)
	require.Equal(t, 1, *header.paddingTop)
	require.Equal(t, 2, *header.paddingLeft)

	require.True(t, *header.borderBottom, "Only bottom border should be enabled")
	require.False(t, *header.borderTop)
	require.False(t, *header.borderLeft)
	require.False(t, *header.borderRight)
}

// TestIntegration_AlertStyle demonstrates an alert box styling scenario.
func TestIntegration_AlertStyle(t *testing.T) {
	red, _ := NewColor("red")
	white, _ := NewColor("white")

	alert := NewStyle().
		Bold(true).
		Foreground(white).
		Background(red).
		Padding(1).
		Border(ThickBorder()).
		BorderForeground(red).
		Width(60).
		Align(Left)

	require.True(t, *alert.bold)
	require.Equal(t, white, *alert.foreground)
	require.Equal(t, red, *alert.background)
	require.Equal(t, 1, *alert.paddingTop)
	require.NotNil(t, alert.borderType)
	require.Equal(t, red, *alert.borderForeground)
	require.Equal(t, 60, *alert.width)
	require.Equal(t, Left, *alert.align)
}

// TestIntegration_ComplexChain demonstrates a very long method chain.
func TestIntegration_ComplexChain(t *testing.T) {
	black, _ := NewColor("black")
	white, _ := NewColor("white")
	gray, _ := NewColor("gray")

	complexStyle := NewStyle().
		Bold(true).
		Italic(false).
		Underline(false).
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

	// Verify text attributes
	require.True(t, *complexStyle.bold)
	require.False(t, *complexStyle.italic)
	require.False(t, *complexStyle.underline)

	// Verify colors
	require.Equal(t, black, *complexStyle.foreground)
	require.Equal(t, white, *complexStyle.background)

	// Verify layout
	require.Equal(t, 80, *complexStyle.width)
	require.Equal(t, 24, *complexStyle.height)
	require.Equal(t, 100, *complexStyle.maxWidth)
	require.Equal(t, 50, *complexStyle.maxHeight)
	require.Equal(t, Center, *complexStyle.align)
	require.Equal(t, Center, *complexStyle.alignVertical)

	// Verify spacing
	require.Equal(t, 2, *complexStyle.paddingTop)
	require.Equal(t, 4, *complexStyle.paddingLeft)
	require.Equal(t, 1, *complexStyle.marginTop)
	require.Equal(t, 2, *complexStyle.marginLeft)

	// Verify borders
	require.NotNil(t, complexStyle.borderType)
	require.Equal(t, gray, *complexStyle.borderForeground)
	require.True(t, *complexStyle.borderTop)
	require.True(t, *complexStyle.borderRight)
	require.True(t, *complexStyle.borderBottom)
	require.True(t, *complexStyle.borderLeft)
}

// TestIntegration_OrderIndependence verifies chaining order doesn't affect final result.
func TestIntegration_OrderIndependence(t *testing.T) {
	red, _ := NewColor("red")

	// Chain in one order
	s1 := NewStyle().
		Bold(true).
		Foreground(red).
		Width(80)

	// Chain in different order
	s2 := NewStyle().
		Width(80).
		Bold(true).
		Foreground(red)

	// Both should have same final state
	require.Equal(t, *s1.bold, *s2.bold)
	require.Equal(t, *s1.foreground, *s2.foreground)
	require.Equal(t, *s1.width, *s2.width)
}

// TestIntegration_StyleReuse demonstrates creating base styles and extending them.
func TestIntegration_StyleReuse(t *testing.T) {
	white, _ := NewColor("white")
	blue, _ := NewColor("blue")
	red, _ := NewColor("red")

	// Base button style
	baseButton := NewStyle().
		Bold(true).
		Foreground(white).
		Padding(1, 3).
		Border(RoundedBorder())

	// Primary button (blue)
	primaryButton := baseButton.Background(blue)

	// Danger button (red)
	dangerButton := baseButton.Background(red)

	// Base should have no background
	require.Nil(t, baseButton.background)

	// Primary should be blue
	require.NotNil(t, primaryButton.background)
	require.Equal(t, blue, *primaryButton.background)
	require.True(t, *primaryButton.bold, "Should inherit bold")

	// Danger should be red
	require.NotNil(t, dangerButton.background)
	require.Equal(t, red, *dangerButton.background)
	require.True(t, *dangerButton.bold, "Should inherit bold")
}

// TestIntegration_AllMethodCategories verifies all method categories work together.
func TestIntegration_AllMethodCategories(t *testing.T) {
	white, _ := NewColor("white")
	blue, _ := NewColor("blue")
	gray, _ := NewColor("gray")

	s := NewStyle().
		// Text attributes
		Bold(true).
		Italic(true).
		Underline(false).
		// Colors
		Foreground(white).
		Background(blue).
		// Layout
		Width(80).
		Height(24).
		Align(Center).
		AlignVertical(Center).
		// Spacing
		Padding(2).
		Margin(1).
		// Borders
		Border(RoundedBorder()).
		BorderForeground(gray)

	// Verify one field from each category
	require.True(t, *s.bold, "Text attributes should work")
	require.Equal(t, white, *s.foreground, "Colors should work")
	require.Equal(t, 80, *s.width, "Layout should work")
	require.Equal(t, 2, *s.paddingTop, "Spacing should work")
	require.NotNil(t, s.borderType, "Borders should work")
}
