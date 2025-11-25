package tuistyles

import "fmt"

// Border sets the border type and optionally which edges to draw.
//
// Accepts 0, 1, or 4 edge arguments:
//   - 0 args: all edges enabled (true)
//   - 1 arg: all edges set to same value
//   - 4 args: top, right, bottom, left
//
// Returns a new Style with border set, leaving the original unchanged.
//
// Example:
//
//	s1 := NewStyle().Border(RoundedBorder())             // All edges enabled
//	s2 := NewStyle().Border(ThickBorder(), false)        // All edges disabled
//	s3 := NewStyle().Border(NormalBorder(), true, false, true, false)  // Top and bottom only
func (s Style) Border(borderType Border, edges ...bool) Style {
	s2 := s
	s2.borderType = &borderType

	switch len(edges) {
	case 0:
		// All edges enabled by default
		t := true
		s2.borderTop = &t
		s2.borderRight = &t
		s2.borderBottom = &t
		s2.borderLeft = &t
	case 1:
		// All edges set to same value
		s2.borderTop = &edges[0]
		s2.borderRight = &edges[0]
		s2.borderBottom = &edges[0]
		s2.borderLeft = &edges[0]
	case 4:
		// Top, right, bottom, left
		s2.borderTop = &edges[0]
		s2.borderRight = &edges[1]
		s2.borderBottom = &edges[2]
		s2.borderLeft = &edges[3]
	default:
		panic(fmt.Sprintf("Border() accepts 0, 1, or 4 edge arguments, got %d", len(edges)))
	}

	return s2
}

// BorderForeground sets the border line color.
//
// Returns a new Style with borderForeground set, leaving the original unchanged.
//
// Example:
//
//	red, _ := NewColor("red")
//	s := NewStyle().Border(RoundedBorder()).BorderForeground(red)
func (s Style) BorderForeground(c Color) Style {
	s2 := s
	s2.borderForeground = &c
	return s2
}

// BorderBackground sets the border background color.
//
// Returns a new Style with borderBackground set, leaving the original unchanged.
//
// Example:
//
//	blue, _ := NewColor("blue")
//	s := NewStyle().Border(RoundedBorder()).BorderBackground(blue)
func (s Style) BorderBackground(c Color) Style {
	s2 := s
	s2.borderBackground = &c
	return s2
}

// BorderTop sets whether the top border edge is drawn.
//
// Returns a new Style with borderTop set, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Border(RoundedBorder()).BorderTop(false)  // No top border
func (s Style) BorderTop(v bool) Style {
	s2 := s
	s2.borderTop = &v
	return s2
}

// BorderRight sets whether the right border edge is drawn.
//
// Returns a new Style with borderRight set, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Border(RoundedBorder()).BorderRight(false)  // No right border
func (s Style) BorderRight(v bool) Style {
	s2 := s
	s2.borderRight = &v
	return s2
}

// BorderBottom sets whether the bottom border edge is drawn.
//
// Returns a new Style with borderBottom set, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Border(RoundedBorder()).BorderBottom(false)  // No bottom border
func (s Style) BorderBottom(v bool) Style {
	s2 := s
	s2.borderBottom = &v
	return s2
}

// BorderLeft sets whether the left border edge is drawn.
//
// Returns a new Style with borderLeft set, leaving the original unchanged.
//
// Example:
//
//	s := NewStyle().Border(RoundedBorder()).BorderLeft(false)  // No left border
func (s Style) BorderLeft(v bool) Style {
	s2 := s
	s2.borderLeft = &v
	return s2
}
