## Purpose

Implement the seven text attribute builder methods (Bold, Italic, Underline, Strikethrough, Faint, Blink, Reverse) using the copy-on-write pattern. This establishes the core immutability pattern that all other builder methods will follow.

## Acceptance Criteria

- [ ] Bold(bool) method implemented and returns new Style
- [ ] Italic(bool) method implemented and returns new Style
- [ ] Underline(bool) method implemented and returns new Style
- [ ] Strikethrough(bool) method implemented and returns new Style
- [ ] Faint(bool) method implemented and returns new Style
- [ ] Blink(bool) method implemented and returns new Style
- [ ] Reverse(bool) method implemented and returns new Style
- [ ] All methods follow copy-on-write pattern (shallow copy, modify field, return copy)
- [ ] Unit test for each method verifies value is set correctly
- [ ] Unit test for each method verifies original Style is unchanged (immutability)

## Technical Approach

Implement all seven methods in a new `text.go` file. Each method follows the same pattern:
1. Shallow copy the receiver Style (Go copies all fields automatically)
2. Modify the specific field on the copy
3. Return the modified copy

**Example Implementation**:
```go
// Bold sets the bold text attribute.
// Returns a new Style with bold set to v, leaving the original unchanged.
func (s Style) Bold(v bool) Style {
    s2 := s  // Shallow copy (all pointer fields copied)
    s2.bold = &v  // Allocate new bool pointer with value v
    return s2
}
```

**Why Shallow Copy Works**: Style contains only pointer fields. Shallow copy duplicates the pointers (not the values they point to), which is safe because we replace the pointer entirely (`s2.bold = &v`) rather than modifying the pointed-to value.

**Files to Create/Modify**:
- text.go (create new file with 7 methods)
- text_test.go (create new file with 14+ tests)

**Dependencies**:
- style.go (Style struct from task 001)

## Testing Strategy

**Unit Tests** (in text_test.go):
- TestBold_SetTrue: Call Bold(true), assert field is true
- TestBold_SetFalse: Call Bold(false), assert field is false
- TestBold_Immutability: Call Bold(true) on original, assert original unchanged
- (Repeat pattern for Italic, Underline, Strikethrough, Faint, Blink, Reverse)

**Table-Driven Test**:
```go
func TestTextAttributes(t *testing.T) {
    tests := []struct {
        name   string
        setter func(Style, bool) Style
        getter func(Style) *bool
    }{
        {"Bold", func(s Style, v bool) Style { return s.Bold(v) }, func(s Style) *bool { return s.bold }},
        {"Italic", func(s Style, v bool) Style { return s.Italic(v) }, func(s Style) *bool { return s.italic }},
        // ... etc
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewStyle()
            s2 := tt.setter(s, true)
            if tt.getter(s) != nil { t.Error("original modified") }
            if *tt.getter(s2) != true { t.Error("value not set") }
        })
    }
}
```

## Notes

- These methods establish the pattern for all 30+ builder methods
- Shallow copy is safe here because we replace pointers, not modify pointed-to values
- Each method allocates a new bool pointer - minor overhead but necessary for optionality
- Consider using `require` package for clearer test assertions (e.g., `require.True(t, *s2.bold)`)
- This task proves the immutability pattern works before scaling to 30+ methods


