## Purpose

Implement border builder methods for controlling border type, individual edges, and border colors. Provides flexible border configuration with sensible defaults.

## Acceptance Criteria

- [ ] Border(BorderType, ...bool) method sets border type and optionally individual edges
- [ ] BorderForeground(Color) method sets border color
- [ ] BorderBackground(Color) method sets border background color
- [ ] BorderTop(bool), BorderRight(bool), BorderBottom(bool), BorderLeft(bool) individual edge methods
- [ ] Border() with no edge args enables all sides by default
- [ ] All methods follow copy-on-write pattern
- [ ] Unit tests for each method
- [ ] Unit tests verify immutability

## Technical Approach

Implement 7 methods in a new `border.go` file. Border() accepts variadic bool args for flexible edge control (0 args = all edges, 1 arg = all edges, 4 args = individual edges).

**Border Method with Variadic Edge Control**:
```go
// Border sets the border type and optionally which edges to draw.
// 0 args: all edges enabled
// 1 arg: all edges set to same value
// 4 args: top, right, bottom, left
func (s Style) Border(borderType BorderType, edges ...bool) Style {
    s2 := s
    s2.borderType = &borderType

    switch len(edges) {
    case 0:
        // All edges enabled by default
        t := true
        s2.borderTop = &t
        s2.borderRight = &t
        s2.borderBottom = &t
        s2.borderLeft = &t
    case 1:
        // All edges set to same value
        s2.borderTop = &edges[0]
        s2.borderRight = &edges[0]
        s2.borderBottom = &edges[0]
        s2.borderLeft = &edges[0]
    case 4:
        // Individual edges: top, right, bottom, left
        s2.borderTop = &edges[0]
        s2.borderRight = &edges[1]
        s2.borderBottom = &edges[2]
        s2.borderLeft = &edges[3]
    default:
        panic(fmt.Sprintf("Border() accepts 0, 1, or 4 edge arguments, got %d", len(edges)))
    }

    return s2
}
```

**Border Color Methods**:
```go
// BorderForeground sets the border foreground color.
func (s Style) BorderForeground(c Color) Style {
    s2 := s
    s2.borderForeground = &c
    return s2
}

// BorderBackground sets the border background color.
func (s Style) BorderBackground(c Color) Style {
    s2 := s
    s2.borderBackground = &c
    return s2
}
```

**Individual Edge Methods**:
```go
// BorderTop enables/disables the top border edge.
func (s Style) BorderTop(v bool) Style {
    s2 := s
    s2.borderTop = &v
    return s2
}

// BorderRight enables/disables the right border edge.
func (s Style) BorderRight(v bool) Style {
    s2 := s
    s2.borderRight = &v
    return s2
}

// BorderBottom enables/disables the bottom border edge.
func (s Style) BorderBottom(v bool) Style {
    s2 := s
    s2.borderBottom = &v
    return s2
}

// BorderLeft enables/disables the left border edge.
func (s Style) BorderLeft(v bool) Style {
    s2 := s
    s2.borderLeft = &v
    return s2
}
```

**Files to Create/Modify**:
- border.go (create new file with 7 methods)
- border_test.go (create new file with tests)

**Dependencies**:
- border.go from milestone-1 (BorderType enum)
- color.go from milestone-1 (Color type)
- style.go (Style struct)

## Testing Strategy

**Unit Tests** (in border_test.go):

**Border Method Tests**:
- TestBorder_NoEdges: Border(Rounded) enables all sides
- TestBorder_OneEdge: Border(Rounded, true) enables all sides
- TestBorder_OneEdgeFalse: Border(Rounded, false) disables all sides
- TestBorder_FourEdges: Border(Rounded, true, false, true, false) sets individual edges
- TestBorder_InvalidArgCount: Border(Rounded, true, false) panics with clear message
- TestBorder_Immutability: Original unchanged

**Border Color Tests**:
- TestBorderForeground_SetColor: Verify color set correctly
- TestBorderForeground_Immutability: Original unchanged
- TestBorderBackground_SetColor: Verify color set correctly
- TestBorderBackground_Immutability: Original unchanged

**Individual Edge Tests**:
- TestBorderTop_Enable: BorderTop(true) enables top
- TestBorderTop_Disable: BorderTop(false) disables top
- (Repeat for BorderRight, BorderBottom, BorderLeft)

**Table-Driven Test for Edge Variations**:
```go
func TestBorder_EdgeVariations(t *testing.T) {
    tests := []struct {
        name    string
        edges   []bool
        wantTop bool
        wantRight bool
        wantBottom bool
        wantLeft bool
    }{
        {"no args (default all)", []bool{}, true, true, true, true},
        {"one arg true", []bool{true}, true, true, true, true},
        {"one arg false", []bool{false}, false, false, false, false},
        {"four args mixed", []bool{true, false, true, false}, true, false, true, false},
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

## Notes

- Border() variadic design inspired by CSS border shorthand
- Default behavior (no edges) is "all enabled" - most common use case
- Individual edge methods allow fine-tuned control (e.g., bottom border only)
- BorderForeground/Background enable colored borders
- Consider adding BorderStyle() alias for Border() in later milestone if users find naming confusing


