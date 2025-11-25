package tuistyles

import "testing"

// TestTextAttributes tests all text attribute methods using table-driven tests.
// This verifies that each method sets the correct field and maintains immutability.
func TestTextAttributes(t *testing.T) {
	tests := []struct {
		name   string
		setter func(Style, bool) Style
		getter func(Style) *bool
	}{
		{
			"Bold",
			func(s Style, v bool) Style { return s.Bold(v) },
			func(s Style) *bool { return s.bold },
		},
		{
			"Italic",
			func(s Style, v bool) Style { return s.Italic(v) },
			func(s Style) *bool { return s.italic },
		},
		{
			"Underline",
			func(s Style, v bool) Style { return s.Underline(v) },
			func(s Style) *bool { return s.underline },
		},
		{
			"Strikethrough",
			func(s Style, v bool) Style { return s.Strikethrough(v) },
			func(s Style) *bool { return s.strikethrough },
		},
		{
			"Faint",
			func(s Style, v bool) Style { return s.Faint(v) },
			func(s Style) *bool { return s.faint },
		},
		{
			"Blink",
			func(s Style, v bool) Style { return s.Blink(v) },
			func(s Style) *bool { return s.blink },
		},
		{
			"Reverse",
			func(s Style, v bool) Style { return s.Reverse(v) },
			func(s Style) *bool { return s.reverse },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test setting to true
			t.Run("SetTrue", func(t *testing.T) {
				s := NewStyle()
				s2 := tt.setter(s, true)

				// Verify original is unchanged (immutability)
				if tt.getter(s) != nil {
					t.Errorf("Original Style was modified: expected nil, got %v", *tt.getter(s))
				}

				// Verify new Style has value set
				if tt.getter(s2) == nil {
					t.Errorf("New Style field is nil, expected true")
				} else if *tt.getter(s2) != true {
					t.Errorf("New Style field = %v, expected true", *tt.getter(s2))
				}
			})

			// Test setting to false
			t.Run("SetFalse", func(t *testing.T) {
				s := NewStyle()
				s2 := tt.setter(s, false)

				// Verify original is unchanged (immutability)
				if tt.getter(s) != nil {
					t.Errorf("Original Style was modified: expected nil, got %v", *tt.getter(s))
				}

				// Verify new Style has value set
				if tt.getter(s2) == nil {
					t.Errorf("New Style field is nil, expected false")
				} else if *tt.getter(s2) != false {
					t.Errorf("New Style field = %v, expected false", *tt.getter(s2))
				}
			})

			// Test chaining (multiple calls)
			t.Run("Chaining", func(t *testing.T) {
				s := NewStyle()
				s2 := tt.setter(s, true)
				s3 := tt.setter(s2, false)

				// Verify original is unchanged
				if tt.getter(s) != nil {
					t.Errorf("Original Style was modified: expected nil, got %v", *tt.getter(s))
				}

				// Verify intermediate Style is unchanged
				if tt.getter(s2) == nil || *tt.getter(s2) != true {
					t.Errorf("Intermediate Style was modified")
				}

				// Verify final Style has correct value
				if tt.getter(s3) == nil || *tt.getter(s3) != false {
					t.Errorf("Final Style has incorrect value")
				}
			})
		})
	}
}

// TestBold_MethodChaining tests that Bold can be chained with other methods.
func TestBold_MethodChaining(t *testing.T) {
	s := NewStyle().Bold(true).Italic(true).Underline(true)

	if s.bold == nil || *s.bold != true {
		t.Errorf("Bold not set correctly after chaining")
	}
	if s.italic == nil || *s.italic != true {
		t.Errorf("Italic not set correctly after chaining")
	}
	if s.underline == nil || *s.underline != true {
		t.Errorf("Underline not set correctly after chaining")
	}
}

// TestTextAttributes_CopyIndependence verifies that modifying a copy doesn't affect the original.
func TestTextAttributes_CopyIndependence(t *testing.T) {
	s1 := NewStyle().Bold(true)
	s2 := s1.Italic(true)
	s3 := s1.Underline(true)

	// s1 should only have bold set
	if s1.bold == nil || *s1.bold != true {
		t.Errorf("s1.bold not set correctly")
	}
	if s1.italic != nil {
		t.Errorf("s1.italic should be nil, got %v", *s1.italic)
	}
	if s1.underline != nil {
		t.Errorf("s1.underline should be nil, got %v", *s1.underline)
	}

	// s2 should have bold and italic set
	if s2.bold == nil || *s2.bold != true {
		t.Errorf("s2.bold not set correctly")
	}
	if s2.italic == nil || *s2.italic != true {
		t.Errorf("s2.italic not set correctly")
	}
	if s2.underline != nil {
		t.Errorf("s2.underline should be nil, got %v", *s2.underline)
	}

	// s3 should have bold and underline set
	if s3.bold == nil || *s3.bold != true {
		t.Errorf("s3.bold not set correctly")
	}
	if s3.italic != nil {
		t.Errorf("s3.italic should be nil, got %v", *s3.italic)
	}
	if s3.underline == nil || *s3.underline != true {
		t.Errorf("s3.underline not set correctly")
	}
}
