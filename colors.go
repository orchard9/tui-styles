package tuistyles

import "fmt"

// Foreground sets the foreground (text) color.
//
// Returns a new Style with foreground set to c, leaving the original unchanged.
//
// Example:
//
//	red, _ := NewColor("red")
//	s := NewStyle().Foreground(red)
//	fmt.Println(s.Render("Red text"))
func (s Style) Foreground(c Color) Style {
	s2 := s
	s2.foreground = &c
	return s2
}

// Background sets the background color.
//
// Returns a new Style with background set to c, leaving the original unchanged.
//
// Example:
//
//	blue, _ := NewColor("blue")
//	s := NewStyle().Background(blue)
//	fmt.Println(s.Render("Text with blue background"))
func (s Style) Background(c Color) Style {
	s2 := s
	s2.background = &c
	return s2
}

// SetString parses a color string and sets it as the foreground color.
//
// Accepts hex colors ("#FF0000"), named colors ("red"), or ANSI codes ("1").
// Returns error if the color string is invalid.
//
// Returns a new Style with foreground set, leaving the original unchanged.
//
// Example:
//
//	s, err := NewStyle().SetString("red")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(s.Render("Red text"))
func (s Style) SetString(colorStr string) (Style, error) {
	c, err := NewColor(colorStr)
	if err != nil {
		return s, fmt.Errorf("invalid color string %q: %w", colorStr, err)
	}
	return s.Foreground(c), nil
}
