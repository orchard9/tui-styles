package tuistyles

// Bold sets the bold/bright text attribute.
//
// Returns a new Style with bold set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Bold(true)
//	fmt.Println(s.Render("Bold text"))
func (s Style) Bold(v bool) Style {
	s2 := s      // Shallow copy (all pointer fields copied)
	s2.bold = &v // Allocate new bool pointer with value v
	return s2
}

// Italic sets the italic/slanted text attribute.
//
// Returns a new Style with italic set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Italic(true)
//	fmt.Println(s.Render("Italic text"))
func (s Style) Italic(v bool) Style {
	s2 := s
	s2.italic = &v
	return s2
}

// Underline sets the underlined text attribute.
//
// Returns a new Style with underline set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Underline(true)
//	fmt.Println(s.Render("Underlined text"))
func (s Style) Underline(v bool) Style {
	s2 := s
	s2.underline = &v
	return s2
}

// Strikethrough sets the strikethrough/crossed-out text attribute.
//
// Returns a new Style with strikethrough set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Strikethrough(true)
//	fmt.Println(s.Render("Strikethrough text"))
func (s Style) Strikethrough(v bool) Style {
	s2 := s
	s2.strikethrough = &v
	return s2
}

// Faint sets the faint/dim text attribute.
//
// Returns a new Style with faint set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Faint(true)
//	fmt.Println(s.Render("Faint text"))
func (s Style) Faint(v bool) Style {
	s2 := s
	s2.faint = &v
	return s2
}

// Blink sets the blinking text attribute.
//
// Note: Blinking text is rarely supported by modern terminals.
//
// Returns a new Style with blink set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Blink(true)
//	fmt.Println(s.Render("Blinking text"))
func (s Style) Blink(v bool) Style {
	s2 := s
	s2.blink = &v
	return s2
}

// Reverse sets the reverse video attribute (swap foreground/background colors).
//
// Returns a new Style with reverse set to v, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Reverse(true)
//	fmt.Println(s.Render("Reversed text"))
func (s Style) Reverse(v bool) Style {
	s2 := s
	s2.reverse = &v
	return s2
}
