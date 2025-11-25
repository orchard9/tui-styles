## Purpose

Thoroughly test CSS shorthand methods (Padding, Margin, Border) with edge cases, invalid inputs, and boundary conditions. Ensures robust error handling and clear failure modes.

## Acceptance Criteria

- [ ] All CSS shorthand variations (1/2/4 args) tested exhaustively
- [ ] Invalid arg counts (0, 3, 5+) trigger clear panic messages
- [ ] Negative values handled correctly (clamped to 0)
- [ ] Maximum values tested (INT_MAX)
- [ ] Empty variadic calls handled correctly
- [ ] Panic messages clearly explain what went wrong
- [ ] Tests serve as negative documentation (what not to do)

## Technical Approach

Create a dedicated `shorthand_test.go` file with comprehensive edge case coverage for Padding, Margin, and Border shorthand methods.

**Edge Cases to Test**:

1. **Invalid Arg Counts - Padding**:
```go
func TestPadding_InvalidArgCount(t *testing.T) {
    tests := []struct {
        name     string
        args     []int
        wantPanic bool
        panicMsg  string
    }{
        {"zero args", []int{}, true, "accepts 1, 2, or 4 arguments, got 0"},
        {"three args", []int{1, 2, 3}, true, "accepts 1, 2, or 4 arguments, got 3"},
        {"five args", []int{1, 2, 3, 4, 5}, true, "accepts 1, 2, or 4 arguments, got 5"},
        {"ten args", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, true, "accepts 1, 2, or 4 arguments, got 10"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            defer func() {
                r := recover()
                if tt.wantPanic {
                    require.NotNil(t, r, "expected panic")
                    require.Contains(t, r.(string), tt.panicMsg)
                }
            }()
            _ = NewStyle().Padding(tt.args...)
            if tt.wantPanic {
                t.Fatal("expected panic, got none")
            }
        })
    }
}
```

2. **Invalid Arg Counts - Margin**:
```go
func TestMargin_InvalidArgCount(t *testing.T) {
    // Same pattern as TestPadding_InvalidArgCount
    tests := []struct {
        name     string
        args     []int
        wantPanic bool
    }{
        {"zero args", []int{}, true},
        {"three args", []int{1, 2, 3}, true},
        {"five args", []int{1, 2, 3, 4, 5}, true},
    }
    // ... (test implementation)
}
```

3. **Invalid Arg Counts - Border**:
```go
func TestBorder_InvalidEdgeCount(t *testing.T) {
    tests := []struct {
        name     string
        edges    []bool
        wantPanic bool
    }{
        {"two edges", []bool{true, false}, true},
        {"three edges", []bool{true, false, true}, true},
        {"five edges", []bool{true, false, true, false, true}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.wantPanic {
                require.Panics(t, func() {
                    _ = NewStyle().Border(Rounded, tt.edges...)
                })
            } else {
                require.NotPanics(t, func() {
                    _ = NewStyle().Border(Rounded, tt.edges...)
                })
            }
        })
    }
}
```

4. **Negative Values**:
```go
func TestShorthand_NegativeValues(t *testing.T) {
    tests := []struct {
        name     string
        fn       func() Style
        expected int
    }{
        {"padding one negative", func() Style { return NewStyle().Padding(-5) }, 0},
        {"padding two negative", func() Style { return NewStyle().Padding(-1, -2) }, 0},
        {"padding four negative", func() Style { return NewStyle().Padding(-1, -2, -3, -4) }, 0},
        {"margin one negative", func() Style { return NewStyle().Margin(-5) }, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := tt.fn()
            // All padding/margin fields should be clamped to 0
            if s.paddingTop != nil {
                require.Equal(t, tt.expected, *s.paddingTop)
            }
            if s.marginTop != nil {
                require.Equal(t, tt.expected, *s.marginTop)
            }
        })
    }
}
```

5. **Maximum Values**:
```go
func TestShorthand_MaximumValues(t *testing.T) {
    const maxInt = int(^uint(0) >> 1)  // Platform-specific max int

    s := NewStyle().
        Padding(maxInt).
        Margin(maxInt)

    require.Equal(t, maxInt, *s.paddingTop)
    require.Equal(t, maxInt, *s.marginTop)
}
```

6. **Zero Values**:
```go
func TestShorthand_ZeroValues(t *testing.T) {
    s := NewStyle().
        Padding(0).
        Margin(0)

    // Zero is valid and should be set
    require.NotNil(t, s.paddingTop)
    require.Equal(t, 0, *s.paddingTop)
    require.NotNil(t, s.marginTop)
    require.Equal(t, 0, *s.marginTop)
}
```

7. **Mixed Positive and Negative**:
```go
func TestShorthand_MixedValues(t *testing.T) {
    s := NewStyle().Padding(5, -2)

    // Positive preserved, negative clamped
    require.Equal(t, 5, *s.paddingTop)
    require.Equal(t, 0, *s.paddingLeft)  // -2 clamped to 0
}
```

8. **Border Edge Combinations**:
```go
func TestBorder_EdgeCombinations(t *testing.T) {
    tests := []struct {
        name        string
        edges       []bool
        wantTop     bool
        wantRight   bool
        wantBottom  bool
        wantLeft    bool
    }{
        {"all true", []bool{true, true, true, true}, true, true, true, true},
        {"all false", []bool{false, false, false, false}, false, false, false, false},
        {"top only", []bool{true, false, false, false}, true, false, false, false},
        {"bottom only", []bool{false, false, true, false}, false, false, true, false},
        {"vertical only", []bool{true, false, true, false}, true, false, true, false},
        {"horizontal only", []bool{false, true, false, true}, false, true, false, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := NewStyle().Border(Rounded, tt.edges...)
            require.Equal(t, tt.wantTop, *s.borderTop)
            require.Equal(t, tt.wantRight, *s.borderRight)
            require.Equal(t, tt.wantBottom, *s.borderBottom)
            require.Equal(t, tt.wantLeft, *s.borderLeft)
        })
    }
}
```

**Files to Create/Modify**:
- shorthand_test.go (create new file)

**Dependencies**:
- spacing.go (Padding, Margin methods from task 006)
- border.go (Border method from task 007)

## Testing Strategy

**Edge Case Categories**:
- Invalid arg counts (0, 3, 5+ args)
- Negative values (single, multiple, mixed)
- Zero values (explicit zero setting)
- Maximum values (INT_MAX)
- Border edge combinations (all permutations)

**Panic Testing**:
- Use require.Panics() for expected panics
- Verify panic messages contain helpful context
- Test panic recovery doesn't leave corrupted state

**Boundary Testing**:
- Test minimum (negative → 0), zero, and maximum values
- Verify clamping behavior is consistent
- Ensure no integer overflow issues

## Notes

- Comprehensive edge case testing prevents production surprises
- Clear panic messages help developers debug quickly
- Negative value clamping is safer than panicking
- These tests document the API contract for edge cases
- Consider adding fuzzing tests (go test -fuzz) for additional coverage

