## Purpose

Implement spacing builder methods with CSS-style shorthand support. Padding(...int) and Margin(...int) accept 1/2/4 arguments (all, vertical+horizontal, top+right+bottom+left). Individual edge methods provide fine-grained control.

## Acceptance Criteria

- [ ] Padding(...int) method with CSS shorthand (1/2/4 args)
- [ ] Margin(...int) method with CSS shorthand (1/2/4 args)
- [ ] PaddingTop/Right/Bottom/Left(int) individual edge methods
- [ ] MarginTop/Right/Bottom/Left(int) individual edge methods
- [ ] Invalid arg counts (3, 5+) handled gracefully (panic with clear message)
- [ ] Negative values clamped to 0
- [ ] All methods follow copy-on-write pattern
- [ ] Unit tests for each shorthand variation (1/2/4 args)
- [ ] Unit tests for individual edge methods
- [ ] Unit tests verify immutability

## Technical Approach

Implement 10 methods in a new `spacing.go` file. Shorthand methods use variadic args and switch on len(args). Individual edge methods are straightforward setters.

**CSS Shorthand Implementation**:
```go
// Padding sets padding using CSS shorthand:
// 1 arg: all sides
// 2 args: top/bottom, left/right
// 4 args: top, right, bottom, left
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

// Margin sets margin using CSS shorthand (same pattern as Padding).
func (s Style) Margin(values ...int) Style {
    // Same implementation as Padding
}
```

**Individual Edge Methods**:
```go
// PaddingTop sets top padding. Negative values are clamped to 0.
func (s Style) PaddingTop(v int) Style {
    if v < 0 {
        v = 0
    }
    s2 := s
    s2.paddingTop = &v
    return s2
}

// PaddingRight sets right padding. Negative values are clamped to 0.
func (s Style) PaddingRight(v int) Style {
    if v < 0 {
        v = 0
    }
    s2 := s
    s2.paddingRight = &v
    return s2
}

// (Repeat for PaddingBottom, PaddingLeft, MarginTop, MarginRight, MarginBottom, MarginLeft)
```

**Files to Create/Modify**:
- spacing.go (create new file with 10 methods)
- spacing_test.go (create new file with tests)

**Dependencies**:
- style.go (Style struct)

## Testing Strategy

**Unit Tests** (in spacing_test.go):

**Shorthand Tests**:
- TestPadding_OneArg: Padding(2) sets all sides to 2
- TestPadding_TwoArgs: Padding(1, 2) sets top/bottom=1, left/right=2
- TestPadding_FourArgs: Padding(1, 2, 3, 4) sets each side correctly
- TestPadding_InvalidArgCount: Padding(1, 2, 3) panics with clear message
- TestPadding_Negative: Padding(-1) clamps all sides to 0
- (Repeat pattern for Margin)

**Individual Edge Tests**:
- TestPaddingTop_Set: PaddingTop(5) sets top padding to 5
- TestPaddingTop_Negative: PaddingTop(-1) clamps to 0
- TestPaddingTop_Immutability: Original unchanged
- (Repeat for all 8 individual edge methods)

**Chaining Test**:
```go
func TestSpacing_Chaining(t *testing.T) {
    s := NewStyle().
        Padding(2).
        MarginTop(1)

    require.Equal(t, 2, *s.paddingTop)
    require.Equal(t, 2, *s.paddingRight)
    require.Equal(t, 1, *s.marginTop)
    require.Nil(t, s.marginRight) // Not set
}
```

**Table-Driven Test for CSS Shorthand**:
```go
func TestPadding_CSSShorthand(t *testing.T) {
    tests := []struct {
        name   string
        args   []int
        top    int
        right  int
        bottom int
        left   int
    }{
        {"one arg", []int{5}, 5, 5, 5, 5},
        {"two args", []int{2, 4}, 2, 4, 2, 4},
        {"four args", []int{1, 2, 3, 4}, 1, 2, 3, 4},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewStyle().Padding(tt.args...)
            require.Equal(t, tt.top, *s.paddingTop)
            require.Equal(t, tt.right, *s.paddingRight)
            require.Equal(t, tt.bottom, *s.paddingBottom)
            require.Equal(t, tt.left, *s.paddingLeft)
        })
    }
}
```

## Notes

- CSS shorthand is familiar to web developers - improves ergonomics
- Panic on invalid arg count (3, 5+) is acceptable - fail fast
- Shorthand methods chain individual methods - ensures consistency
- Negative value clamping prevents nonsensical spacing
- Individual edge methods provide fine-grained control when needed
- Consider documenting CSS shorthand pattern in package-level godoc


