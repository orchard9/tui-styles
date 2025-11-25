## Purpose

Enhance multi-line rendering to properly handle complex scenarios with width constraints, truncation, and wrapping. This builds on the basic multi-line support from task 003 to add production-ready multi-line handling.

## Acceptance Criteria

- [ ] Multi-line strings split and styled per line correctly
- [ ] Width constraint enforcement (MaxWidth property)
- [ ] Line truncation with ellipsis (...) when exceeding MaxWidth
- [ ] Word wrapping support (optional, break long lines)
- [ ] Each line independently styled with ANSI codes
- [ ] Line height calculation accurate
- [ ] Empty lines handled (preserves spacing)
- [ ] Unit tests for all multi-line scenarios
- [ ] Zero linter warnings

## Technical Approach

Extend rendering in `render.go` with enhanced multi-line logic:

**Multi-line Processing**:
1. Split input by newlines
2. Apply width constraint to each line (truncate or wrap)
3. Style each line independently
4. Reassemble with newlines

**Width Constraint**:
- Add `MaxWidth(int)` method to Style
- Truncate lines exceeding MaxWidth using internal/measure.Truncate()
- Add ellipsis (...) to indicate truncation

**Line Styling**:
- Apply ANSI prefix to each line
- Append ANSI reset to each line
- Ensures terminal compatibility (some terminals reset per line)

**Files to Create/Modify**:
- render.go (enhance multi-line rendering)
- render_test.go (add multi-line tests)
- style.go (add MaxWidth field)

**Dependencies**:
- Task 002 (measurement for width calculation)
- Task 003 (basic rendering foundation)

## Testing Strategy

**Unit Tests**:
- Multi-line without width constraint
- Multi-line with MaxWidth (truncation)
- Multi-line with empty lines
- Very long lines (wrapping)
- ANSI preservation across lines

**Test Example**:
```go
func TestMultiLineMaxWidth(t *testing.T) {
    s := NewStyle().Bold(true).MaxWidth(10)
    result := s.Render("short\nverylonglinehere")
    lines := strings.Split(result, "\n")

    assert.Equal(t, 2, len(lines))
    assert.True(t, measure.Width(lines[1]) <= 10)
    assert.Contains(t, lines[1], "...")
}
```

## Notes

- Builds on basic multi-line from task 003
- Critical for text-heavy terminal UIs
- Truncation strategy: simple cut with ellipsis (word wrapping is advanced, optional)
- Dependencies: Must complete task 002 (measure) and 003 (basic render) first

