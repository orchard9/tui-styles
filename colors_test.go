package tuistyles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestForeground_SetColor verifies that Foreground sets the color correctly.
func TestForeground_SetColor(t *testing.T) {
	red, err := NewColor("red")
	require.NoError(t, err)

	s := NewStyle()
	s2 := s.Foreground(red)

	// Original should be unchanged
	require.Nil(t, s.foreground, "Original Style was mutated")

	// New Style should have foreground set
	require.NotNil(t, s2.foreground, "Foreground not set")
	require.Equal(t, red, *s2.foreground, "Foreground color incorrect")
}

// TestForeground_Immutability verifies that Foreground doesn't mutate the original.
func TestForeground_Immutability(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	s1 := NewStyle().Foreground(red)
	s2 := s1.Foreground(blue)

	// s1 should still be red
	require.NotNil(t, s1.foreground)
	require.Equal(t, red, *s1.foreground, "s1 foreground was mutated")

	// s2 should be blue
	require.NotNil(t, s2.foreground)
	require.Equal(t, blue, *s2.foreground, "s2 foreground incorrect")
}

// TestBackground_SetColor verifies that Background sets the color correctly.
func TestBackground_SetColor(t *testing.T) {
	blue, err := NewColor("blue")
	require.NoError(t, err)

	s := NewStyle()
	s2 := s.Background(blue)

	// Original should be unchanged
	require.Nil(t, s.background, "Original Style was mutated")

	// New Style should have background set
	require.NotNil(t, s2.background, "Background not set")
	require.Equal(t, blue, *s2.background, "Background color incorrect")
}

// TestBackground_Immutability verifies that Background doesn't mutate the original.
func TestBackground_Immutability(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	s1 := NewStyle().Background(red)
	s2 := s1.Background(blue)

	// s1 should still be red
	require.NotNil(t, s1.background)
	require.Equal(t, red, *s1.background, "s1 background was mutated")

	// s2 should be blue
	require.NotNil(t, s2.background)
	require.Equal(t, blue, *s2.background, "s2 background incorrect")
}

// TestForegroundBackground_Independent verifies foreground and background are independent.
func TestForegroundBackground_Independent(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	s := NewStyle().Foreground(red).Background(blue)

	require.NotNil(t, s.foreground)
	require.Equal(t, red, *s.foreground, "Foreground incorrect")

	require.NotNil(t, s.background)
	require.Equal(t, blue, *s.background, "Background incorrect")
}

// TestSetString tests SetString with various color formats.
func TestSetString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid inputs
		{"hex_6digit", "#FF0000", false},
		{"hex_3digit", "#F00", false},
		{"named_red", "red", false},
		{"named_blue", "blue", false},
		{"named_green", "green", false},
		{"named_uppercase", "RED", false},
		{"ansi_code_1", "1", false},
		{"ansi_code_196", "196", false},

		// Invalid inputs
		{"invalid_name", "not-a-color", true},
		{"invalid_hex", "#GG0000", true},
		{"empty", "", true},
		{"invalid_code", "256", true},
		{"invalid_code_negative", "-1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStyle()
			s2, err := s.SetString(tt.input)

			if tt.wantErr {
				require.Error(t, err, "Expected error for input %q", tt.input)
				// Original should be unchanged even on error
				require.Nil(t, s.foreground, "Original Style was mutated on error")
			} else {
				require.NoError(t, err, "Unexpected error for input %q", tt.input)
				require.NotNil(t, s2.foreground, "Foreground not set for valid input %q", tt.input)
				// Original should be unchanged
				require.Nil(t, s.foreground, "Original Style was mutated")
			}
		})
	}
}

// TestSetString_Immutability verifies SetString doesn't mutate the original on error.
func TestSetString_Immutability(t *testing.T) {
	s := NewStyle()
	_, err := s.SetString("invalid-color")

	require.Error(t, err, "Expected error for invalid color")
	require.Nil(t, s.foreground, "Original Style was mutated on error")
}

// TestSetString_ChainAfterError verifies we can chain after SetString error.
func TestSetString_ChainAfterError(t *testing.T) {
	s := NewStyle()
	s2, err := s.SetString("invalid")
	require.Error(t, err)

	// Should still be able to use the returned Style (which is the original)
	red, _ := NewColor("red")
	s3 := s2.Foreground(red)

	require.NotNil(t, s3.foreground)
	require.Equal(t, red, *s3.foreground)
}

// TestColorMethods_Chaining verifies color methods work in chains.
func TestColorMethods_Chaining(t *testing.T) {
	red, _ := NewColor("red")
	blue, _ := NewColor("blue")

	s := NewStyle().
		Bold(true).
		Foreground(red).
		Italic(true).
		Background(blue).
		Underline(true)

	// Verify all attributes set
	require.NotNil(t, s.bold)
	require.True(t, *s.bold)

	require.NotNil(t, s.italic)
	require.True(t, *s.italic)

	require.NotNil(t, s.underline)
	require.True(t, *s.underline)

	require.NotNil(t, s.foreground)
	require.Equal(t, red, *s.foreground)

	require.NotNil(t, s.background)
	require.Equal(t, blue, *s.background)
}

// TestSetString_ThenChain verifies SetString works in method chains.
func TestSetString_ThenChain(t *testing.T) {
	s, err := NewStyle().Bold(true).SetString("red")
	require.NoError(t, err)

	s2 := s.Italic(true)

	// Verify all attributes set
	require.NotNil(t, s2.bold)
	require.True(t, *s2.bold)

	require.NotNil(t, s2.italic)
	require.True(t, *s2.italic)

	require.NotNil(t, s2.foreground)
}

// TestColorMethods_AllColorTypes tests all color type variations.
func TestColorMethods_AllColorTypes(t *testing.T) {
	tests := []struct {
		name      string
		colorStr  string
		wantError bool
	}{
		{"hex_red", "#FF0000", false},
		{"hex_3digit", "#F00", false},
		{"named_red", "red", false},
		{"ansi_code", "196", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := NewColor(tt.colorStr)
			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Test Foreground
			s := NewStyle().Foreground(color)
			require.NotNil(t, s.foreground)
			require.Equal(t, color, *s.foreground)

			// Test Background
			s2 := NewStyle().Background(color)
			require.NotNil(t, s2.background)
			require.Equal(t, color, *s2.background)

			// Test SetString
			s3, err := NewStyle().SetString(tt.colorStr)
			require.NoError(t, err)
			require.NotNil(t, s3.foreground)
		})
	}
}
