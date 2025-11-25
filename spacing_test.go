package tuistyles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPadding_CSSShorthand tests CSS shorthand variations.
func TestPadding_CSSShorthand(t *testing.T) {
	tests := []struct {
		name   string
		args   []int
		top    int
		right  int
		bottom int
		left   int
	}{
		{"one_arg_all_sides", []int{5}, 5, 5, 5, 5},
		{"two_args_vertical_horizontal", []int{2, 4}, 2, 4, 2, 4},
		{"four_args_clockwise", []int{1, 2, 3, 4}, 1, 2, 3, 4},
		{"one_arg_zero", []int{0}, 0, 0, 0, 0},
		{"two_args_with_zero", []int{0, 5}, 0, 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle().Padding(tt.args...)

			require.NotNil(t, s.paddingTop)
			require.Equal(t, tt.top, *s.paddingTop)

			require.NotNil(t, s.paddingRight)
			require.Equal(t, tt.right, *s.paddingRight)

			require.NotNil(t, s.paddingBottom)
			require.Equal(t, tt.bottom, *s.paddingBottom)

			require.NotNil(t, s.paddingLeft)
			require.Equal(t, tt.left, *s.paddingLeft)
		})
	}
}

// TestMargin_CSSShorthand tests CSS shorthand variations.
func TestMargin_CSSShorthand(t *testing.T) {
	tests := []struct {
		name   string
		args   []int
		top    int
		right  int
		bottom int
		left   int
	}{
		{"one_arg_all_sides", []int{3}, 3, 3, 3, 3},
		{"two_args_vertical_horizontal", []int{1, 2}, 1, 2, 1, 2},
		{"four_args_clockwise", []int{5, 6, 7, 8}, 5, 6, 7, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle().Margin(tt.args...)

			require.NotNil(t, s.marginTop)
			require.Equal(t, tt.top, *s.marginTop)

			require.NotNil(t, s.marginRight)
			require.Equal(t, tt.right, *s.marginRight)

			require.NotNil(t, s.marginBottom)
			require.Equal(t, tt.bottom, *s.marginBottom)

			require.NotNil(t, s.marginLeft)
			require.Equal(t, tt.left, *s.marginLeft)
		})
	}
}

// TestPadding_InvalidArgCount verifies panic on invalid argument counts.
func TestPadding_InvalidArgCount(t *testing.T) {
	tests := []struct {
		name string
		args []int
	}{
		{"zero_args", []int{}},
		{"three_args", []int{1, 2, 3}},
		{"five_args", []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Panics(t, func() {
				NewStyle().Padding(tt.args...)
			}, "Expected panic for %d args", len(tt.args))
		})
	}
}

// TestMargin_InvalidArgCount verifies panic on invalid argument counts.
func TestMargin_InvalidArgCount(t *testing.T) {
	tests := []struct {
		name string
		args []int
	}{
		{"zero_args", []int{}},
		{"three_args", []int{1, 2, 3}},
		{"five_args", []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Panics(t, func() {
				NewStyle().Margin(tt.args...)
			}, "Expected panic for %d args", len(tt.args))
		})
	}
}

// TestPadding_NegativeValues verifies negative values are clamped to 0.
func TestPadding_NegativeValues(t *testing.T) {
	tests := []struct {
		name string
		args []int
	}{
		{"one_negative", []int{-5}},
		{"two_negative", []int{-1, -2}},
		{"four_negative", []int{-1, -2, -3, -4}},
		{"mixed", []int{-1, 5, -3, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle().Padding(tt.args...)

			// All negative values should be clamped to 0
			if tt.args[0] < 0 {
				require.Equal(t, 0, *s.paddingTop, "Negative top not clamped")
			}

			// For shorthand methods, we test the clamping happens
			require.NotNil(t, s.paddingTop)
			require.NotNil(t, s.paddingRight)
			require.NotNil(t, s.paddingBottom)
			require.NotNil(t, s.paddingLeft)

			// None should be negative
			require.GreaterOrEqual(t, *s.paddingTop, 0)
			require.GreaterOrEqual(t, *s.paddingRight, 0)
			require.GreaterOrEqual(t, *s.paddingBottom, 0)
			require.GreaterOrEqual(t, *s.paddingLeft, 0)
		})
	}
}

// TestIndividualPaddingMethods tests each individual padding edge method.
func TestIndividualPaddingMethods(t *testing.T) {
	tests := []struct {
		name   string
		setter func(Style, int) Style
		getter func(Style) *int
		value  int
		expect int
	}{
		{"PaddingTop_positive", func(s Style, v int) Style { return s.PaddingTop(v) }, func(s Style) *int { return s.paddingTop }, 5, 5},
		{"PaddingTop_negative", func(s Style, v int) Style { return s.PaddingTop(v) }, func(s Style) *int { return s.paddingTop }, -3, 0},
		{"PaddingRight_positive", func(s Style, v int) Style { return s.PaddingRight(v) }, func(s Style) *int { return s.paddingRight }, 10, 10},
		{"PaddingRight_negative", func(s Style, v int) Style { return s.PaddingRight(v) }, func(s Style) *int { return s.paddingRight }, -1, 0},
		{"PaddingBottom_positive", func(s Style, v int) Style { return s.PaddingBottom(v) }, func(s Style) *int { return s.paddingBottom }, 3, 3},
		{"PaddingBottom_negative", func(s Style, v int) Style { return s.PaddingBottom(v) }, func(s Style) *int { return s.paddingBottom }, -5, 0},
		{"PaddingLeft_positive", func(s Style, v int) Style { return s.PaddingLeft(v) }, func(s Style) *int { return s.paddingLeft }, 7, 7},
		{"PaddingLeft_negative", func(s Style, v int) Style { return s.PaddingLeft(v) }, func(s Style) *int { return s.paddingLeft }, -2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := tt.setter(s, tt.value)

			// Original unchanged
			require.Nil(t, tt.getter(s), "Original was mutated")

			// New Style has correct value
			require.NotNil(t, tt.getter(s2))
			require.Equal(t, tt.expect, *tt.getter(s2))
		})
	}
}

// TestIndividualMarginMethods tests each individual margin edge method.
func TestIndividualMarginMethods(t *testing.T) {
	tests := []struct {
		name   string
		setter func(Style, int) Style
		getter func(Style) *int
		value  int
		expect int
	}{
		{"MarginTop_positive", func(s Style, v int) Style { return s.MarginTop(v) }, func(s Style) *int { return s.marginTop }, 2, 2},
		{"MarginTop_negative", func(s Style, v int) Style { return s.MarginTop(v) }, func(s Style) *int { return s.marginTop }, -1, 0},
		{"MarginRight_positive", func(s Style, v int) Style { return s.MarginRight(v) }, func(s Style) *int { return s.marginRight }, 4, 4},
		{"MarginRight_negative", func(s Style, v int) Style { return s.MarginRight(v) }, func(s Style) *int { return s.marginRight }, -2, 0},
		{"MarginBottom_positive", func(s Style, v int) Style { return s.MarginBottom(v) }, func(s Style) *int { return s.marginBottom }, 1, 1},
		{"MarginBottom_negative", func(s Style, v int) Style { return s.MarginBottom(v) }, func(s Style) *int { return s.marginBottom }, -3, 0},
		{"MarginLeft_positive", func(s Style, v int) Style { return s.MarginLeft(v) }, func(s Style) *int { return s.marginLeft }, 6, 6},
		{"MarginLeft_negative", func(s Style, v int) Style { return s.MarginLeft(v) }, func(s Style) *int { return s.marginLeft }, -4, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2 := tt.setter(s, tt.value)

			// Original unchanged
			require.Nil(t, tt.getter(s), "Original was mutated")

			// New Style has correct value
			require.NotNil(t, tt.getter(s2))
			require.Equal(t, tt.expect, *tt.getter(s2))
		})
	}
}

// TestSpacing_Chaining verifies spacing methods chain correctly.
func TestSpacing_Chaining(t *testing.T) {
	s := NewStyle().
		Padding(2).
		MarginTop(1).
		PaddingLeft(5)

	// Verify padding shorthand set all sides
	require.Equal(t, 2, *s.paddingTop)
	require.Equal(t, 2, *s.paddingRight)
	require.Equal(t, 2, *s.paddingBottom)

	// Verify individual override worked
	require.Equal(t, 5, *s.paddingLeft, "PaddingLeft should override shorthand")

	// Verify margin
	require.Equal(t, 1, *s.marginTop)
	require.Nil(t, s.marginRight, "MarginRight should not be set")
	require.Nil(t, s.marginBottom, "MarginBottom should not be set")
	require.Nil(t, s.marginLeft, "MarginLeft should not be set")
}

// TestSpacing_Immutability verifies spacing maintains immutability.
func TestSpacing_Immutability(t *testing.T) {
	s1 := NewStyle().Padding(2)
	s2 := s1.Padding(4)

	// s1 should still be 2 on all sides
	require.Equal(t, 2, *s1.paddingTop)
	require.Equal(t, 2, *s1.paddingRight)
	require.Equal(t, 2, *s1.paddingBottom)
	require.Equal(t, 2, *s1.paddingLeft)

	// s2 should be 4 on all sides
	require.Equal(t, 4, *s2.paddingTop)
	require.Equal(t, 4, *s2.paddingRight)
	require.Equal(t, 4, *s2.paddingBottom)
	require.Equal(t, 4, *s2.paddingLeft)
}

// TestSpacing_WithOtherMethods verifies spacing chains with other Style methods.
func TestSpacing_WithOtherMethods(t *testing.T) {
	red, _ := NewColor("red")

	s := NewStyle().
		Bold(true).
		Padding(2).
		Foreground(red).
		Margin(1, 2).
		Width(80)

	// Verify all attributes set
	require.NotNil(t, s.bold)
	require.True(t, *s.bold)

	require.Equal(t, 2, *s.paddingTop)
	require.Equal(t, 2, *s.paddingRight)
	require.Equal(t, 2, *s.paddingBottom)
	require.Equal(t, 2, *s.paddingLeft)

	require.NotNil(t, s.foreground)
	require.Equal(t, red, *s.foreground)

	require.Equal(t, 1, *s.marginTop)
	require.Equal(t, 2, *s.marginRight)
	require.Equal(t, 1, *s.marginBottom)
	require.Equal(t, 2, *s.marginLeft)

	require.Equal(t, 80, *s.width)
}
