package tuistyles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBorder_NoEdges tests Border with no edge arguments (all enabled).
func TestBorder_NoEdges(t *testing.T) {
	rounded := RoundedBorder()
	s := NewStyle().Border(rounded)

	// Border type should be set
	require.NotNil(t, s.borderType)
	require.Equal(t, rounded, *s.borderType)

	// All edges should be enabled
	require.NotNil(t, s.borderTop)
	require.True(t, *s.borderTop)

	require.NotNil(t, s.borderRight)
	require.True(t, *s.borderRight)

	require.NotNil(t, s.borderBottom)
	require.True(t, *s.borderBottom)

	require.NotNil(t, s.borderLeft)
	require.True(t, *s.borderLeft)
}

// TestBorder_OneEdge tests Border with one edge argument (all edges same value).
func TestBorder_OneEdge(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{"all_enabled", true},
		{"all_disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thick := ThickBorder()
			s := NewStyle().Border(thick, tt.value)

			require.NotNil(t, s.borderType)
			require.Equal(t, thick, *s.borderType)

			require.Equal(t, tt.value, *s.borderTop)
			require.Equal(t, tt.value, *s.borderRight)
			require.Equal(t, tt.value, *s.borderBottom)
			require.Equal(t, tt.value, *s.borderLeft)
		})
	}
}

// TestBorder_FourEdges tests Border with four edge arguments.
func TestBorder_FourEdges(t *testing.T) {
	normal := NormalBorder()
	s := NewStyle().Border(normal, true, false, true, false)

	require.NotNil(t, s.borderType)
	require.Equal(t, normal, *s.borderType)

	require.True(t, *s.borderTop, "Top should be enabled")
	require.False(t, *s.borderRight, "Right should be disabled")
	require.True(t, *s.borderBottom, "Bottom should be enabled")
	require.False(t, *s.borderLeft, "Left should be disabled")
}

// TestBorder_InvalidArgCount verifies panic on invalid edge argument counts.
func TestBorder_InvalidArgCount(t *testing.T) {
	tests := []struct {
		name  string
		edges []bool
	}{
		{"two_args", []bool{true, false}},
		{"three_args", []bool{true, false, true}},
		{"five_args", []bool{true, false, true, false, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Panics(t, func() {
				NewStyle().Border(RoundedBorder(), tt.edges...)
			}, "Expected panic for %d edge args", len(tt.edges))
		})
	}
}

// TestBorder_Immutability verifies Border doesn't mutate the original.
func TestBorder_Immutability(t *testing.T) {
	rounded := RoundedBorder()
	thick := ThickBorder()

	s1 := NewStyle().Border(rounded)
	s2 := s1.Border(thick)

	// s1 should still be rounded
	require.NotNil(t, s1.borderType)
	require.Equal(t, rounded, *s1.borderType)

	// s2 should be thick
	require.NotNil(t, s2.borderType)
	require.Equal(t, thick, *s2.borderType)
}

// TestBorderForeground tests border foreground color.
func TestBorderForeground(t *testing.T) {
	red, _ := NewColor("red")
	s := NewStyle().BorderForeground(red)

	require.NotNil(t, s.borderForeground)
	require.Equal(t, red, *s.borderForeground)
}

// TestBorderBackground tests border background color.
func TestBorderBackground(t *testing.T) {
	blue, _ := NewColor("blue")
	s := NewStyle().BorderBackground(blue)

	require.NotNil(t, s.borderBackground)
	require.Equal(t, blue, *s.borderBackground)
}

// TestBorderColors_Immutability verifies border color methods don't mutate.
func TestBorderColors_Immutability(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	s1 := NewStyle().BorderForeground(red)
	s2 := s1.BorderForeground(blue)

	// s1 should still be red
	require.Equal(t, red, *s1.borderForeground)

	// s2 should be blue
	require.Equal(t, blue, *s2.borderForeground)
}

// TestIndividualBorderEdges tests individual edge methods.
func TestIndividualBorderEdges(t *testing.T) {
	tests := []struct {
		name   string
		setter func(Style, bool) Style
		getter func(Style) *bool
		value  bool
	}{
		{"BorderTop_true", func(s Style, v bool) Style { return s.BorderTop(v) }, func(s Style) *bool { return s.borderTop }, true},
		{"BorderTop_false", func(s Style, v bool) Style { return s.BorderTop(v) }, func(s Style) *bool { return s.borderTop }, false},
		{"BorderRight_true", func(s Style, v bool) Style { return s.BorderRight(v) }, func(s Style) *bool { return s.borderRight }, true},
		{"BorderRight_false", func(s Style, v bool) Style { return s.BorderRight(v) }, func(s Style) *bool { return s.borderRight }, false},
		{"BorderBottom_true", func(s Style, v bool) Style { return s.BorderBottom(v) }, func(s Style) *bool { return s.borderBottom }, true},
		{"BorderBottom_false", func(s Style, v bool) Style { return s.BorderBottom(v) }, func(s Style) *bool { return s.borderBottom }, false},
		{"BorderLeft_true", func(s Style, v bool) Style { return s.BorderLeft(v) }, func(s Style) *bool { return s.borderLeft }, true},
		{"BorderLeft_false", func(s Style, v bool) Style { return s.BorderLeft(v) }, func(s Style) *bool { return s.borderLeft }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := tt.setter(s, tt.value)

			// Original unchanged
			require.Nil(t, tt.getter(s))

			// New Style has correct value
			require.NotNil(t, tt.getter(s2))
			require.Equal(t, tt.value, *tt.getter(s2))
		})
	}
}

// TestBorder_OverrideEdges verifies individual edge methods can override Border().
func TestBorder_OverrideEdges(t *testing.T) {
	rounded := RoundedBorder()

	s := NewStyle().
		Border(rounded). // All edges enabled
		BorderTop(false) // Disable top

	require.NotNil(t, s.borderType)
	require.Equal(t, rounded, *s.borderType)

	require.False(t, *s.borderTop, "Top should be disabled")
	require.True(t, *s.borderRight, "Right should still be enabled")
	require.True(t, *s.borderBottom, "Bottom should still be enabled")
	require.True(t, *s.borderLeft, "Left should still be enabled")
}

// TestBorder_CompleteChain tests full border styling chain.
func TestBorder_CompleteChain(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")
	rounded := RoundedBorder()

	s := NewStyle().
		Border(rounded).
		BorderForeground(red).
		BorderBackground(blue).
		BorderTop(false)

	require.Equal(t, rounded, *s.borderType)
	require.Equal(t, red, *s.borderForeground)
	require.Equal(t, blue, *s.borderBackground)
	require.False(t, *s.borderTop)
	require.True(t, *s.borderRight)
	require.True(t, *s.borderBottom)
	require.True(t, *s.borderLeft)
}

// TestBorder_WithOtherStyles verifies border methods chain with other Style methods.
func TestBorder_WithOtherStyles(t *testing.T) {
	red, _ := NewColor("red")
	thick := ThickBorder()

	s := NewStyle().
		Bold(true).
		Width(80).
		Padding(2).
		Border(thick).
		Foreground(red)

	require.True(t, *s.bold)
	require.Equal(t, 80, *s.width)
	require.Equal(t, 2, *s.paddingTop)
	require.Equal(t, thick, *s.borderType)
	require.Equal(t, red, *s.foreground)
}
