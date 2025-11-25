package tuistyles

import "fmt"

// Padding sets padding using CSS shorthand notation.
//
// Accepts 1, 2, or 4 arguments:
//   - 1 arg: all sides (top, right, bottom, left)
//   - 2 args: vertical (top/bottom), horizontal (left/right)
//   - 4 args: top, right, bottom, left (clockwise from top)
//
// Negative values are clamped to 0. Panics if given 3 or 5+ arguments.
//
// Returns a new Style with padding set, leaving the original unchanged.
//
// Example:
//
//	s1 := NewStyle().Padding(2)           // All sides = 2
//	s2 := NewStyle().Padding(1, 3)        // Top/bottom=1, left/right=3
//	s3 := NewStyle().Padding(1, 2, 3, 4)  // Top=1, right=2, bottom=3, left=4
func (s Style) Padding(values ...int) Style {
	switch len(values) {
	case 1:
		// All sides
		return s.PaddingTop(values[0]).
			PaddingRight(values[0]).
			PaddingBottom(values[0]).
			PaddingLeft(values[0])
	case 2:
		// Vertical, horizontal
		return s.PaddingTop(values[0]).
			PaddingRight(values[1]).
			PaddingBottom(values[0]).
			PaddingLeft(values[1])
	case 4:
		// Top, right, bottom, left
		return s.PaddingTop(values[0]).
			PaddingRight(values[1]).
			PaddingBottom(values[2]).
			PaddingLeft(values[3])
	default:
		panic(fmt.Sprintf("Padding() accepts 1, 2, or 4 arguments, got %d", len(values)))
	}
}

// Margin sets margin using CSS shorthand notation.
//
// Accepts 1, 2, or 4 arguments:
//   - 1 arg: all sides (top, right, bottom, left)
//   - 2 args: vertical (top/bottom), horizontal (left/right)
//   - 4 args: top, right, bottom, left (clockwise from top)
//
// Negative values are clamped to 0. Panics if given 3 or 5+ arguments.
//
// Returns a new Style with margin set, leaving the original unchanged.
//
// Example:
//
//	s1 := NewStyle().Margin(1)           // All sides = 1
//	s2 := NewStyle().Margin(2, 4)        // Top/bottom=2, left/right=4
//	s3 := NewStyle().Margin(1, 2, 3, 4)  // Top=1, right=2, bottom=3, left=4
func (s Style) Margin(values ...int) Style {
	switch len(values) {
	case 1:
		// All sides
		return s.MarginTop(values[0]).
			MarginRight(values[0]).
			MarginBottom(values[0]).
			MarginLeft(values[0])
	case 2:
		// Vertical, horizontal
		return s.MarginTop(values[0]).
			MarginRight(values[1]).
			MarginBottom(values[0]).
			MarginLeft(values[1])
	case 4:
		// Top, right, bottom, left
		return s.MarginTop(values[0]).
			MarginRight(values[1]).
			MarginBottom(values[2]).
			MarginLeft(values[3])
	default:
		panic(fmt.Sprintf("Margin() accepts 1, 2, or 4 arguments, got %d", len(values)))
	}
}

// PaddingTop sets top padding in cells.
//
// Negative values are clamped to 0. Returns a new Style with paddingTop set,
// leaving the original unchanged.
func (s Style) PaddingTop(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.paddingTop = &v
	return s2
}

// PaddingRight sets right padding in cells.
//
// Negative values are clamped to 0. Returns a new Style with paddingRight set,
// leaving the original unchanged.
func (s Style) PaddingRight(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.paddingRight = &v
	return s2
}

// PaddingBottom sets bottom padding in cells.
//
// Negative values are clamped to 0. Returns a new Style with paddingBottom set,
// leaving the original unchanged.
func (s Style) PaddingBottom(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.paddingBottom = &v
	return s2
}

// PaddingLeft sets left padding in cells.
//
// Negative values are clamped to 0. Returns a new Style with paddingLeft set,
// leaving the original unchanged.
func (s Style) PaddingLeft(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.paddingLeft = &v
	return s2
}

// MarginTop sets top margin in lines.
//
// Negative values are clamped to 0. Returns a new Style with marginTop set,
// leaving the original unchanged.
func (s Style) MarginTop(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.marginTop = &v
	return s2
}

// MarginRight sets right margin in cells.
//
// Negative values are clamped to 0. Returns a new Style with marginRight set,
// leaving the original unchanged.
func (s Style) MarginRight(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.marginRight = &v
	return s2
}

// MarginBottom sets bottom margin in lines.
//
// Negative values are clamped to 0. Returns a new Style with marginBottom set,
// leaving the original unchanged.
func (s Style) MarginBottom(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.marginBottom = &v
	return s2
}

// MarginLeft sets left margin in cells.
//
// Negative values are clamped to 0. Returns a new Style with marginLeft set,
// leaving the original unchanged.
func (s Style) MarginLeft(v int) Style {
	if v < 0 {
		v = 0
	}
	s2 := s
	s2.marginLeft = &v
	return s2
}
