## Purpose

Implement color builder methods (Foreground, Background, SetString) following the established copy-on-write pattern. SetString adds convenience for parsing color strings like "#FF0000" or "red".

## Acceptance Criteria

- [ ] Foreground(Color) method implemented and returns new Style
- [ ] Background(Color) method implemented and returns new Style
- [ ] SetString(string) method implemented - parses color string and sets foreground
- [ ] All methods follow copy-on-write pattern
- [ ] SetString validates color strings and returns error for invalid input
- [ ] Unit tests for each method verify value set correctly
- [ ] Unit tests verify immutability preserved
- [ ] Tests cover hex colors, named colors, and invalid inputs

## Technical Approach

Implement three methods in a new `color.go` file. Foreground and Background follow the standard pattern. SetString adds parsing logic using the Color type's parsing capabilities.

**Foreground and Background**:
```go
// Foreground sets the foreground (text) color.
func (s Style) Foreground(c Color) Style {
    s2 := s
    s2.foreground = &c
    return s2
}

// Background sets the background color.
func (s Style) Background(c Color) Style {
    s2 := s
    s2.background = &c
    return s2
}
```

**SetString with Validation**:
```go
// SetString parses a color string and sets it as the foreground color.
// Accepts hex colors ("#FF0000"), named colors ("red"), or ANSI codes ("1").
// Returns error if the color string is invalid.
func (s Style) SetString(colorStr string) (Style, error) {
    c, err := ParseColor(colorStr)
    if err != nil {
        return s, fmt.Errorf("invalid color string %q: %w", colorStr, err)
    }
    return s.Foreground(c), nil
}
```

**Files to Create/Modify**:
- color.go (create new file with 3 methods)
- color_test.go (create new file with tests)

**Dependencies**:
- color.go from milestone-1 (Color type and ParseColor function)
- style.go (Style struct)

## Testing Strategy

**Unit Tests** (in color_test.go):
- TestForeground_SetColor: Verify foreground color set correctly
- TestForeground_Immutability: Verify original Style unchanged
- TestBackground_SetColor: Verify background color set correctly
- TestBackground_Immutability: Verify original Style unchanged
- TestSetString_ValidHex: Parse "#FF0000" successfully
- TestSetString_ValidNamed: Parse "red" successfully
- TestSetString_ValidANSI: Parse "1" (red ANSI code) successfully
- TestSetString_Invalid: Return error for "not-a-color"
- TestSetString_Immutability: Verify original unchanged on error

**Table-Driven Tests**:
```go
func TestSetString(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"hex", "#FF0000", false},
        {"named", "red", false},
        {"ansi", "1", false},
        {"invalid", "not-a-color", true},
        {"empty", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewStyle()
            s2, err := s.SetString(tt.input)
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
                require.NotNil(t, s2.foreground)
            }
        })
    }
}
```

## Notes

- SetString is convenience method - users can always use Foreground(ParseColor(...))
- Error handling on SetString is important - don't panic on invalid input
- Consider adding GetForeground()/GetBackground() helper methods in later milestone
- Color validation depends on ParseColor implementation from milestone-1
- If ParseColor doesn't exist yet, create stub that panics with TODO comment


