## Purpose

Implement layout builder methods (Width, Height, MaxWidth, MaxHeight, Align, AlignVertical) for controlling text box dimensions and alignment. Includes validation to prevent negative dimensions.

## Acceptance Criteria

- [ ] Width(int) method implemented with validation (>= 0)
- [ ] Height(int) method implemented with validation (>= 0)
- [ ] MaxWidth(int) method implemented with validation (>= 0)
- [ ] MaxHeight(int) method implemented with validation (>= 0)
- [ ] Align(Position) method implemented for horizontal alignment
- [ ] AlignVertical(Position) method implemented for vertical alignment
- [ ] All methods follow copy-on-write pattern
- [ ] Negative dimensions are rejected (panic or clamp to 0)
- [ ] Unit tests verify values set correctly
- [ ] Unit tests verify validation works
- [ ] Unit tests verify immutability

## Technical Approach

Implement six methods in a new `layout.go` file. Dimension methods validate input to prevent negative values. Alignment methods accept Position enum (Left, Center, Right, Top, Middle, Bottom).

**Dimension Methods with Validation**:
```go
// Width sets the width of the styled text box.
// Negative values are clamped to 0.
func (s Style) Width(w int) Style {
    if w < 0 {
        w = 0
    }
    s2 := s
    s2.width = &w
    return s2
}

// Height sets the height of the styled text box.
// Negative values are clamped to 0.
func (s Style) Height(h int) Style {
    if h < 0 {
        h = 0
    }
    s2 := s
    s2.height = &h
    return s2
}

// MaxWidth sets the maximum width before wrapping.
// Negative values are clamped to 0.
func (s Style) MaxWidth(w int) Style {
    if w < 0 {
        w = 0
    }
    s2 := s
    s2.maxWidth = &w
    return s2
}

// MaxHeight sets the maximum height before truncating.
// Negative values are clamped to 0.
func (s Style) MaxHeight(h int) Style {
    if h < 0 {
        h = 0
    }
    s2 := s
    s2.maxHeight = &h
    return s2
}
```

**Alignment Methods**:
```go
// Align sets horizontal text alignment (Left, Center, Right).
func (s Style) Align(p Position) Style {
    s2 := s
    s2.align = &p
    return s2
}

// AlignVertical sets vertical text alignment (Top, Middle, Bottom).
func (s Style) AlignVertical(p Position) Style {
    s2 := s
    s2.alignVertical = &p
    return s2
}
```

**Files to Create/Modify**:
- layout.go (create new file with 6 methods)
- layout_test.go (create new file with tests)

**Dependencies**:
- position.go from milestone-1 (Position enum)
- style.go (Style struct)

## Testing Strategy

**Unit Tests** (in layout_test.go):
- TestWidth_Positive: Set width to 80, verify set
- TestWidth_Zero: Set width to 0, verify set
- TestWidth_Negative: Set width to -10, verify clamped to 0
- TestWidth_Immutability: Verify original unchanged
- (Repeat pattern for Height, MaxWidth, MaxHeight)
- TestAlign_Left: Set horizontal alignment to Left
- TestAlign_Center: Set horizontal alignment to Center
- TestAlign_Right: Set horizontal alignment to Right
- TestAlignVertical_Top: Set vertical alignment to Top
- TestAlignVertical_Middle: Set vertical alignment to Middle
- TestAlignVertical_Bottom: Set vertical alignment to Bottom

**Table-Driven Test for Validation**:
```go
func TestDimensionValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"positive", 80, 80},
        {"zero", 0, 0},
        {"negative", -10, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewStyle().Width(tt.input)
            require.Equal(t, tt.expected, *s.width)
        })
    }
}
```

## Notes

- Clamping negative values to 0 is safer than panicking
- Align/AlignVertical accept Position enum - validate in rendering layer later
- Consider adding GetWidth()/GetHeight() helpers in later milestone
- Width vs MaxWidth: Width is exact size, MaxWidth is upper bound before wrapping
- Validation prevents nonsensical layouts (negative dimensions)


