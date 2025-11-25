package tuistyles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestImmutability_SingleMethod verifies that calling a single builder method
// doesn't mutate the original Style.
func TestImmutability_SingleMethod(t *testing.T) {
	original := NewStyle()
	modified := original.Bold(true)

	// Original should be unchanged
	require.Nil(t, original.bold, "Original Style was mutated")

	// Modified should have new value
	require.NotNil(t, modified.bold, "Modified Style missing value")
	require.True(t, *modified.bold, "Modified Style has incorrect value")
}

// TestImmutability_Chaining verifies that method chaining creates independent
// Style instances at each step.
func TestImmutability_Chaining(t *testing.T) {
	s1 := NewStyle()
	s2 := s1.Bold(true)
	s3 := s2.Italic(true)

	// s1 should be unchanged
	require.Nil(t, s1.bold, "s1.bold should be nil")
	require.Nil(t, s1.italic, "s1.italic should be nil")

	// s2 should only have bold
	require.NotNil(t, s2.bold, "s2.bold should be set")
	require.True(t, *s2.bold, "s2.bold should be true")
	require.Nil(t, s2.italic, "s2.italic should be nil")

	// s3 should have both
	require.NotNil(t, s3.bold, "s3.bold should be set")
	require.True(t, *s3.bold, "s3.bold should be true")
	require.NotNil(t, s3.italic, "s3.italic should be set")
	require.True(t, *s3.italic, "s3.italic should be true")
}

// TestImmutability_IndependentBranches verifies that creating multiple Styles
// from the same base Style doesn't cause them to interfere with each other.
func TestImmutability_IndependentBranches(t *testing.T) {
	base := NewStyle()
	branch1 := base.Bold(true)
	branch2 := base.Italic(true)
	branch3 := base.Underline(true)

	// Base should be unchanged
	require.Nil(t, base.bold, "base.bold should be nil")
	require.Nil(t, base.italic, "base.italic should be nil")
	require.Nil(t, base.underline, "base.underline should be nil")

	// Branch 1 should only have bold
	require.NotNil(t, branch1.bold, "branch1.bold should be set")
	require.True(t, *branch1.bold, "branch1.bold should be true")
	require.Nil(t, branch1.italic, "branch1.italic should be nil")
	require.Nil(t, branch1.underline, "branch1.underline should be nil")

	// Branch 2 should only have italic
	require.Nil(t, branch2.bold, "branch2.bold should be nil")
	require.NotNil(t, branch2.italic, "branch2.italic should be set")
	require.True(t, *branch2.italic, "branch2.italic should be true")
	require.Nil(t, branch2.underline, "branch2.underline should be nil")

	// Branch 3 should only have underline
	require.Nil(t, branch3.bold, "branch3.bold should be nil")
	require.Nil(t, branch3.italic, "branch3.italic should be nil")
	require.NotNil(t, branch3.underline, "branch3.underline should be set")
	require.True(t, *branch3.underline, "branch3.underline should be true")
}

// TestImmutability_AllTextAttributes verifies that all 7 text attribute methods
// maintain immutability.
func TestImmutability_AllTextAttributes(t *testing.T) {
	tests := []struct {
		name   string
		setter func(Style, bool) Style
		getter func(Style) *bool
	}{
		{"Bold", func(s Style, v bool) Style { return s.Bold(v) }, func(s Style) *bool { return s.bold }},
		{"Italic", func(s Style, v bool) Style { return s.Italic(v) }, func(s Style) *bool { return s.italic }},
		{"Underline", func(s Style, v bool) Style { return s.Underline(v) }, func(s Style) *bool { return s.underline }},
		{"Strikethrough", func(s Style, v bool) Style { return s.Strikethrough(v) }, func(s Style) *bool { return s.strikethrough }},
		{"Faint", func(s Style, v bool) Style { return s.Faint(v) }, func(s Style) *bool { return s.faint }},
		{"Blink", func(s Style, v bool) Style { return s.Blink(v) }, func(s Style) *bool { return s.blink }},
		{"Reverse", func(s Style, v bool) Style { return s.Reverse(v) }, func(s Style) *bool { return s.reverse }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := NewStyle()
			modified := tt.setter(original, true)

			// Original should be unchanged
			require.Nil(t, tt.getter(original), "Original Style was mutated by %s", tt.name)

			// Modified should have value set
			require.NotNil(t, tt.getter(modified), "Modified Style missing value for %s", tt.name)
			require.True(t, *tt.getter(modified), "Modified Style has incorrect value for %s", tt.name)
		})
	}
}

// TestFluentAPI verifies that the fluent method chaining syntax works correctly.
func TestFluentAPI(t *testing.T) {
	// Verify method chaining compiles and works
	s := NewStyle().
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true).
		Faint(true).
		Blink(true).
		Reverse(true)

	// All attributes should be set
	require.NotNil(t, s.bold, "bold should be set")
	require.True(t, *s.bold, "bold should be true")

	require.NotNil(t, s.italic, "italic should be set")
	require.True(t, *s.italic, "italic should be true")

	require.NotNil(t, s.underline, "underline should be set")
	require.True(t, *s.underline, "underline should be true")

	require.NotNil(t, s.strikethrough, "strikethrough should be set")
	require.True(t, *s.strikethrough, "strikethrough should be true")

	require.NotNil(t, s.faint, "faint should be set")
	require.True(t, *s.faint, "faint should be true")

	require.NotNil(t, s.blink, "blink should be set")
	require.True(t, *s.blink, "blink should be true")

	require.NotNil(t, s.reverse, "reverse should be set")
	require.True(t, *s.reverse, "reverse should be true")
}

// TestFluentAPI_Mixed verifies that fluent API works with mixed true/false values.
func TestFluentAPI_Mixed(t *testing.T) {
	s := NewStyle().
		Bold(true).
		Italic(false).
		Underline(true)

	require.NotNil(t, s.bold, "bold should be set")
	require.True(t, *s.bold, "bold should be true")

	require.NotNil(t, s.italic, "italic should be set")
	require.False(t, *s.italic, "italic should be false")

	require.NotNil(t, s.underline, "underline should be set")
	require.True(t, *s.underline, "underline should be true")
}

// TestImmutability_OverwriteValues verifies that calling the same method twice
// creates independent copies.
func TestImmutability_OverwriteValues(t *testing.T) {
	s1 := NewStyle().Bold(true)
	s2 := s1.Bold(false)

	// s1 should still be true
	require.NotNil(t, s1.bold, "s1.bold should be set")
	require.True(t, *s1.bold, "s1.bold should be true")

	// s2 should be false
	require.NotNil(t, s2.bold, "s2.bold should be set")
	require.False(t, *s2.bold, "s2.bold should be false")
}

// TestImmutability_ComplexTree verifies immutability with a complex branching tree.
func TestImmutability_ComplexTree(t *testing.T) {
	// Create a complex tree of Styles
	//        base
	//       /    \
	//      s1    s2
	//     / \     \
	//    s3  s4   s5

	base := NewStyle()
	s1 := base.Bold(true)
	s2 := base.Italic(true)
	s3 := s1.Underline(true)
	s4 := s1.Strikethrough(true)
	s5 := s2.Faint(true)

	// Verify base is unchanged
	require.Nil(t, base.bold)
	require.Nil(t, base.italic)
	require.Nil(t, base.underline)
	require.Nil(t, base.strikethrough)
	require.Nil(t, base.faint)

	// Verify s1 (bold only)
	require.NotNil(t, s1.bold)
	require.True(t, *s1.bold)
	require.Nil(t, s1.italic)
	require.Nil(t, s1.underline)
	require.Nil(t, s1.strikethrough)

	// Verify s2 (italic only)
	require.Nil(t, s2.bold)
	require.NotNil(t, s2.italic)
	require.True(t, *s2.italic)
	require.Nil(t, s2.underline)
	require.Nil(t, s2.faint)

	// Verify s3 (bold + underline)
	require.NotNil(t, s3.bold)
	require.True(t, *s3.bold)
	require.NotNil(t, s3.underline)
	require.True(t, *s3.underline)
	require.Nil(t, s3.strikethrough)

	// Verify s4 (bold + strikethrough)
	require.NotNil(t, s4.bold)
	require.True(t, *s4.bold)
	require.NotNil(t, s4.strikethrough)
	require.True(t, *s4.strikethrough)
	require.Nil(t, s4.underline)

	// Verify s5 (italic + faint)
	require.NotNil(t, s5.italic)
	require.True(t, *s5.italic)
	require.NotNil(t, s5.faint)
	require.True(t, *s5.faint)
	require.Nil(t, s5.bold)
}

// TestRaceConditions verifies there are no data races when copying Styles.
// Run with: go test -race
func TestRaceConditions(t *testing.T) {
	base := NewStyle().Bold(true).Italic(true)

	// Create multiple goroutines that create copies concurrently
	done := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func() {
			// Each goroutine creates its own copy
			s1 := base.Underline(true)
			s2 := base.Strikethrough(true)

			// Verify copies are correct
			require.NotNil(t, s1.bold)
			require.NotNil(t, s1.italic)
			require.NotNil(t, s1.underline)
			require.Nil(t, s1.strikethrough)

			require.NotNil(t, s2.bold)
			require.NotNil(t, s2.italic)
			require.Nil(t, s2.underline)
			require.NotNil(t, s2.strikethrough)

			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify base is still unchanged
	require.NotNil(t, base.bold)
	require.NotNil(t, base.italic)
	require.Nil(t, base.underline)
	require.Nil(t, base.strikethrough)
}
