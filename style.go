// Package tuistyles provides a fluent API for styling terminal output.
//
// TUI Styles allows you to create styled terminal output using an immutable
// builder pattern similar to lipgloss. It supports text attributes (bold, italic),
// colors (hex, ANSI names, 256-color codes), borders, padding, margins, alignment,
// and layout composition utilities.
//
// # Quick Start
//
//	style := NewStyle().
//	    Bold(true).
//	    Foreground(Color("#FF0000")).
//	    Background(Color("blue")).
//	    Padding(2).
//	    Border(RoundedBorder()).
//	    Width(50).
//	    Align(Center)
//
//	fmt.Println(style.Render("Hello, World!"))
//
// # Features
//
//   - Text Attributes: Bold, Italic, Underline, Strikethrough, Faint, Blink, Reverse
//   - Colors: Hex (#RRGGBB), ANSI names (red, blue), 256-color codes (0-255)
//   - Adaptive Colors: Automatically select colors based on terminal background
//   - Borders: 8 predefined border styles with Unicode box drawing
//   - Spacing: Padding and margin with CSS-style shorthand
//   - Alignment: Horizontal (Left, Center, Right) and Vertical (Top, Center, Bottom)
//   - Layout Utilities: JoinHorizontal, JoinVertical, Place for composition
//   - Immutability: All style methods return new instances (thread-safe)
//   - Performance: Sub-millisecond rendering for typical operations
//
// # Color Support
//
// Colors can be specified in multiple formats:
//
//	// Hex colors (3 or 6 digits)
//	red := Color("#FF0000")
//	shortRed := Color("#F00")  // Expands to #FF0000
//
//	// ANSI color names (case-insensitive)
//	blue := Color("blue")
//	green := Color("GREEN")
//
//	// ANSI 256-color codes
//	orange := Color("214")
//
// # Border Styles
//
// Eight predefined border types are available:
//
//   - NormalBorder() - Standard box drawing (─│┌┐└┘)
//   - RoundedBorder() - Rounded corners (─│╭╮╰╯)
//   - ThickBorder() - Heavy lines (━┃┏┓┗┛)
//   - DoubleBorder() - Double lines (═║╔╗╚╝)
//   - BlockBorder() - Solid blocks (█)
//   - OuterHalfBlockBorder() - Outer half-blocks
//   - InnerHalfBlockBorder() - Inner half-blocks
//   - HiddenBorder() - Invisible borders (spaces)
//
// # Layout Composition
//
// Combine styled strings using layout utilities:
//
//	// Side-by-side composition with vertical alignment
//	row := JoinHorizontal(Top, leftPanel, rightPanel)
//
//	// Vertical stacking with horizontal alignment
//	page := JoinVertical(Center, header, row, footer)
//
//	// Absolute positioning in a box
//	centered := Place(80, 24, Center, Center, content)
//
// # Thread Safety
//
// All Style methods are immutable and return new instances, making them
// safe to use concurrently from multiple goroutines.
//
// # Performance
//
// TUI Styles is optimized for performance:
//
//   - Simple text rendering: ~50ns per operation
//   - Complex styles with borders: <1ms typical
//   - Layout composition: <5ms for multi-panel layouts
package tuistyles

// Style represents an immutable text styling configuration.
//
// All fields are pointers to enable optionality - nil indicates "not set"
// while a non-nil pointer indicates an explicit value (even if zero/false).
//
// Styles follow a copy-on-write pattern: all builder methods return a new
// Style with the modified field, never mutating the receiver. This prevents
// spooky action at a distance and enables safe concurrent usage.
type Style struct {
	// Text attributes control font styling
	bold          *bool // Bold/bright text
	italic        *bool // Italic/slanted text
	underline     *bool // Underlined text
	strikethrough *bool // Strikethrough/crossed-out text
	faint         *bool // Faint/dim text
	blink         *bool // Blinking text (rarely supported)
	reverse       *bool // Reverse video (swap foreground/background)

	// Colors define foreground and background colors
	foreground *Color // Text color
	background *Color // Background color

	// Layout defines dimensions and constraints
	width     *int // Fixed width in cells
	height    *int // Fixed height in lines
	maxWidth  *int // Maximum width in cells
	maxHeight *int // Maximum height in lines

	// Alignment controls text positioning
	align         *Position // Horizontal alignment (Left, Center, Right)
	alignVertical *Position // Vertical alignment (Top, Center, Bottom)

	// Spacing controls padding and margins
	paddingTop    *int // Padding above content (cells)
	paddingRight  *int // Padding right of content (cells)
	paddingBottom *int // Padding below content (cells)
	paddingLeft   *int // Padding left of content (cells)
	marginTop     *int // Margin above element (lines)
	marginRight   *int // Margin right of element (cells)
	marginBottom  *int // Margin below element (lines)
	marginLeft    *int // Margin left of element (cells)

	// Borders control border rendering
	borderType       *Border // Border style (Rounded, Thick, etc.)
	borderTop        *bool   // Render top border edge
	borderRight      *bool   // Render right border edge
	borderBottom     *bool   // Render bottom border edge
	borderLeft       *bool   // Render left border edge
	borderForeground *Color  // Border line color
	borderBackground *Color  // Border background color
}

// NewStyle returns a new Style with all fields unset (nil).
//
// Use builder methods to configure the style:
//
//	s := NewStyle().Bold(true).Foreground(color.Red)
//
// Styles are immutable - all builder methods return a new Style, leaving
// the original unchanged.
func NewStyle() Style {
	return Style{}
}
