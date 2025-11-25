## Purpose

Achieve comprehensive test coverage (>80%) across render, border, and layout packages to ensure correctness and enable confident refactoring.

## Acceptance Criteria

- [ ] render/render_test.go with >80% coverage
- [ ] border/border_test.go with >80% coverage
- [ ] layout/layout_test.go with >80% coverage
- [ ] Test coverage for all public APIs in each package
- [ ] Edge case tests (empty strings, zero dimensions, nil values)
- [ ] Complex composition tests (nested rendering, mixed styles)
- [ ] All tests pass with `go test -v ./...`
- [ ] Coverage report shows >80% overall

## Acceptance Criteria (Continued)

**Render Package Tests**:
- [ ] Text rendering with various styles
- [ ] Alignment (left, center, right, justify)
- [ ] Padding and margin application
- [ ] Width constraints and text wrapping
- [ ] Multi-line text handling

**Border Package Tests**:
- [ ] All border styles (single, double, rounded, thick, hidden)
- [ ] Box rendering with titles
- [ ] Border colors and styling
- [ ] Padding inside borders
- [ ] Minimum and maximum dimensions

**Layout Package Tests**:
- [ ] Centering (horizontal, vertical, both)
- [ ] Size constraints (min, max, exact)
- [ ] Overflow handling (clip, wrap, scroll markers)
- [ ] Margin and spacing utilities

## Technical Approach

**Test File Structure**:
```
render/
  render_test.go    - Rendering logic tests
border/
  border_test.go    - Border and box tests
layout/
  layout_test.go    - Layout and sizing tests
```

**Testing Patterns**:
1. Table-driven tests for various input combinations
2. Testify assertions for readable test code
3. Helper functions for common test setups
4. Edge case matrix (empty, single char, very long, Unicode)

**Render Package Test Cases**:
- Render text with style application
- Alignment with various widths
- Padding calculation and application
- Text wrapping at word boundaries
- Multi-line text with consistent styling
- Unicode character handling (emoji, CJK)

**Border Package Test Cases**:
- Each border style renders correctly
- Box dimensions match expected output
- Title placement (top-left, top-center, top-right)
- Border colors combine with content styles
- Nested boxes render without conflicts
- Min/max width enforcement

**Layout Package Test Cases**:
- Center text horizontally at various widths
- Center text vertically at various heights
- Combined centering
- Size constraint enforcement
- Overflow clipping vs wrapping
- Margin calculations

**Files to Create/Modify**:
- render/render_test.go (new)
- border/border_test.go (new)
- layout/layout_test.go (new)

**Dependencies**:
- github.com/stretchr/testify/assert
- github.com/stretchr/testify/require

## Testing Strategy

**Unit Test Coverage**:
- Every exported function has at least one test
- Internal helper functions tested via public API
- Edge cases explicitly tested
- Error conditions validated

**Table-Driven Tests**:
```go
func TestRender_Alignment(t *testing.T) {
    tests := []struct {
        name     string
        text     string
        align    Alignment
        width    int
        expected string
    }{
        {"left short", "hello", AlignLeft, 10, "hello     "},
        {"center short", "hi", AlignCenter, 10, "    hi    "},
        {"right short", "test", AlignRight, 10, "      test"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Render(tt.text, tt.align, tt.width)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Edge Case Matrix**:
- Empty string inputs
- Zero/negative dimensions
- Nil style pointers
- Very long text (>1000 chars)
- Unicode characters (emoji, wide chars)
- ANSI codes in input text

**Coverage Validation**:
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Integration Tests**:
- Render styled text in box with alignment
- Nested boxes with different border styles
- Complex dashboard-like composition
- Performance with large text blocks

## Notes

**Coverage Goals by Package**:
- style: >80% (Phase 1 baseline: 70%)
- render: >80%
- border: >80%
- layout: >80%
- color: >90% (simple package)

**Testing Philosophy**:
- Test behavior, not implementation
- Focus on public API contracts
- Edge cases prevent production bugs
- Table-driven tests improve maintainability

**Common Test Helpers**:
```go
// Helper to create test style
func testStyle() Style {
    return New().Foreground(Red)
}

// Helper to assert box dimensions
func assertBoxSize(t *testing.T, box string, expectedWidth, expectedHeight int) {
    lines := strings.Split(box, "\n")
    assert.Equal(t, expectedHeight, len(lines))
    for _, line := range lines {
        assert.Equal(t, expectedWidth, runewidth.StringWidth(line))
    }
}
```

**Deferred to Future**:
- Property-based testing with quick/gopter (v1.1+)
- Fuzz testing for parser robustness (v1.1+)
- Benchmark-driven optimization (current task 007)

**Reference**:
- Go testing: https://go.dev/doc/tutorial/add-a-test
- Table-driven tests: https://go.dev/wiki/TableDrivenTests
- testify: https://github.com/stretchr/testify

