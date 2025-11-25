# Task 009: JoinHorizontal layout utility

**Status**: pending
**Phase**: phase-2-core-features
**Milestone**: milestone-3
**Dependencies**: []
**Assigned Agent**: none
**Confidence**: TBD (assign during pre-planning)

## Purpose

Implement JoinHorizontal layout utility that places multiple styled strings side-by-side with configurable vertical alignment. Essential for creating multi-column layouts, tables, and complex terminal UIs.

## Acceptance Criteria

- [ ] JoinHorizontal(pos Position, strs ...string) string function implemented
- [ ] Position type supports: Top, Center, Bottom vertical alignment
- [ ] Height normalization (all blocks padded to tallest height)
- [ ] ANSI codes preserved correctly in joined output
- [ ] Handles empty strings gracefully
- [ ] Handles single string (returns as-is)
- [ ] Works with strings of different heights
- [ ] Unit tests covering all alignment positions
- [ ] Zero linter warnings

## Technical Approach

Create `layout.go` with horizontal joining functionality:

**Function Signature**:
```go
// Position represents vertical or horizontal alignment
type Position int

const (
    Top Position = iota
    Center
    Bottom
    Left   // For JoinVertical
    Right  // For JoinVertical
)

// JoinHorizontal joins strings horizontally with vertical alignment
func JoinHorizontal(pos Position, strs ...string) string
```

**Algorithm**:
1. Split each string into lines (handle multi-line strings)
2. Calculate height of each block (number of lines)
3. Find maximum height across all blocks
4. Normalize heights based on position:
   - Top: Add empty lines at bottom
   - Center: Distribute empty lines top/bottom
   - Bottom: Add empty lines at top
5. Measure width of each line (use internal/measure.Width)
6. Join lines horizontally (concatenate corresponding lines from each block)
7. Return joined result

**Height Normalization**:
```go
// For Top alignment:
// Block 1: line1, line2, "", ""
// Block 2: line1, line2, line3, line4

// For Center alignment:
// Block 1: "", line1, line2, ""
// Block 2: line1, line2, line3, line4

// For Bottom alignment:
// Block 1: "", "", line1, line2
// Block 2: line1, line2, line3, line4
```

**Width Measurement**:
- Use `internal/measure.Width()` to get display width
- Required for calculating proper spacing
- Each block maintains its own width

**Edge Cases**:
- Empty strings: Treat as zero-width blocks
- Single string: Return as-is (no joining needed)
- All empty strings: Return empty string
- Mixed ANSI-styled and plain strings: Preserve all ANSI codes

**Implementation Steps**:
1. Create `layout.go` with Position type
2. Implement line splitting and height calculation
3. Implement height normalization for each Position
4. Implement horizontal line concatenation
5. Handle edge cases
6. Write comprehensive unit tests

**Files to Create/Modify**:
- layout.go (JoinHorizontal implementation)
- layout_test.go (unit tests)

**Dependencies**:
- internal/measure (Width calculation, StripANSI - task 002)
- strings package (Split, Join, Repeat)

## Testing Strategy

**Unit Tests**:
- Alignment Top: Short blocks aligned to top
- Alignment Center: Short blocks centered vertically
- Alignment Bottom: Short blocks aligned to bottom
- Same height blocks: No padding needed
- Different heights: Verify padding calculation
- Empty strings: Handle gracefully
- Single string: Return as-is
- ANSI codes: Verify preservation in joined output
- Three or more blocks: Test with multiple blocks

**Test Examples**:
```go
func TestJoinHorizontalTop(t *testing.T) {
    block1 := "line1\nline2"
    block2 := "a\nb\nc\nd"

    result := JoinHorizontal(Top, block1, block2)
    lines := strings.Split(result, "\n")

    assert.Equal(t, 4, len(lines))
    assert.Contains(t, lines[0], "line1")
    assert.Contains(t, lines[0], "a")
    assert.Contains(t, lines[1], "line2")
    assert.Contains(t, lines[1], "b")
    // lines[2] and lines[3] should have padding for block1
}

func TestJoinHorizontalCenter(t *testing.T) {
    short := "line"
    tall := "a\nb\nc\nd"

    result := JoinHorizontal(Center, short, tall)
    lines := strings.Split(result, "\n")

    // Short block should be centered (padded top and bottom)
    assert.Equal(t, 4, len(lines))
    // lines[0]: padding + "a"
    // lines[1-2]: "line" + "b", "c"
    // lines[3]: padding + "d"
}

func TestJoinHorizontalEmpty(t *testing.T) {
    result := JoinHorizontal(Top, "", "content", "")
    assert.Equal(t, "content", strings.TrimSpace(result))
}

func TestJoinHorizontalANSI(t *testing.T) {
    styled := "\x1b[1mBold\x1b[0m"
    plain := "Plain"

    result := JoinHorizontal(Top, styled, plain)
    assert.Contains(t, result, "\x1b[1m")
    assert.Contains(t, result, "Bold")
    assert.Contains(t, result, "Plain")
}
```

**Manual Testing**:
- Print side-by-side boxes to terminal
- Verify visual alignment (top, center, bottom)
- Test with styled content (borders, colors, padding)

**Visual Test Program**:
```go
func main() {
    box1 := NewStyle().
        Border(true).BorderStyle(NormalBorder).
        Padding(1).
        Render("Short")

    box2 := NewStyle().
        Border(true).BorderStyle(RoundedBorder).
        Padding(1).
        Render("Tall\nBox\nWith\nMultiple\nLines")

    fmt.Println("Top Alignment:")
    fmt.Println(JoinHorizontal(Top, box1, box2))

    fmt.Println("\nCenter Alignment:")
    fmt.Println(JoinHorizontal(Center, box1, box2))

    fmt.Println("\nBottom Alignment:")
    fmt.Println(JoinHorizontal(Bottom, box1, box2))
}
```

## Notes

- Critical for building complex layouts (dashboards, tables, side-by-side panels)
- Height normalization is key - all blocks must have same number of lines
- Empty line padding must preserve ANSI codes for background colors
- Width measurement must account for ANSI codes (use internal/measure)
- Efficient implementation: Minimize string allocations, use strings.Builder
- Similar to Lipgloss JoinHorizontal but adapted to our Style system
- Future enhancement: Add spacing parameter (gap between blocks)
- Future enhancement: Support custom padding character
- This is a pure layout utility - no styling applied by JoinHorizontal itself
