## Purpose

Implement border rendering using Unicode box drawing characters to create visual boundaries around styled content. Supports different border styles (normal, rounded, thick, double) and handles corner/edge character selection correctly.

## Acceptance Criteria

- [ ] Border.Top, Border.Bottom, Border.Left, Border.Right implemented
- [ ] BorderStyle types supported (normal, rounded, thick, double)
- [ ] Corner characters rendered correctly (top-left, top-right, bottom-left, bottom-right)
- [ ] Edge characters rendered correctly (horizontal, vertical)
- [ ] Border colors applied (foreground color for border characters)
- [ ] Width calculation accounts for border thickness (left + right = 2 cells)
- [ ] Height calculation accounts for border thickness (top + bottom = 2 lines)
- [ ] Multi-line content properly bordered
- [ ] Unit tests for all border configurations
- [ ] Zero linter warnings

## Technical Approach

Extend rendering logic in `render.go` to handle borders:

**Border Character Sets**:
```go
type BorderStyle struct {
    TopLeft, TopRight, BottomLeft, BottomRight rune
    Horizontal, Vertical rune
}

var (
    NormalBorder = BorderStyle{
        TopLeft: '┌', TopRight: '┐', BottomLeft: '└', BottomRight: '┘',
        Horizontal: '─', Vertical: '│',
    }
    RoundedBorder = BorderStyle{
        TopLeft: '╭', TopRight: '╮', BottomLeft: '╰', BottomRight: '╯',
        Horizontal: '─', Vertical: '│',
    }
    ThickBorder = BorderStyle{
        TopLeft: '┏', TopRight: '┓', BottomLeft: '┗', BottomRight: '┛',
        Horizontal: '━', Vertical: '┃',
    }
    DoubleBorder = BorderStyle{
        TopLeft: '╔', TopRight: '╗', BottomLeft: '╚', BottomRight: '╝',
        Horizontal: '═', Vertical: '║',
    }
)
```

**Rendering Strategy**:
1. Calculate content dimensions (width, height in cells/lines)
2. Render top border: `TopLeft + Horizontal * (width) + TopRight`
3. Render each content line: `Vertical + content + Vertical`
4. Render bottom border: `BottomLeft + Horizontal * (width) + BottomRight`

**Border Configuration**:
- Style methods: `BorderTop(bool)`, `BorderBottom(bool)`, `BorderLeft(bool)`, `BorderRight(bool)`
- Border style method: `BorderStyle(BorderStyle)`
- Border color method: `BorderForeground(Color)` (separate from content color)

**Width/Height Calculation**:
- Content width + left border (1 if enabled) + right border (1 if enabled)
- Content height + top border (1 if enabled) + bottom border (1 if enabled)
- Use `internal/measure.Width()` for content width measurement

**Partial Borders**:
- Support individual border sides (e.g., top + bottom only, left + right only)
- Adjust corner characters based on which sides are enabled
- If only top: no corners, just horizontal line
- If only left + right: vertical bars, no top/bottom lines

**Border Color**:
- Apply BorderForeground color to border characters using ANSI codes
- Content and border can have different colors

**Implementation Steps**:
1. Define BorderStyle types and character sets
2. Add border configuration to Style struct
3. Implement border rendering logic in `render.go`
4. Calculate dimensions with border thickness
5. Handle partial borders (selective sides)
6. Apply border colors
7. Write unit tests for all configurations

**Files to Create/Modify**:
- border.go (BorderStyle definitions and border-related methods)
- render.go (integrate border rendering into main rendering flow)
- border_test.go (unit tests)
- style.go (add border fields to Style struct)

**Dependencies**:
- internal/measure (width calculation - task 002)
- Task 004 (multi-line rendering)
- Task 005 (padding rendering - borders wrap padding)

## Testing Strategy

**Unit Tests**:
- Border styles: Test each style (normal, rounded, thick, double)
- Partial borders: Top only, bottom only, left only, right only, combinations
- Border colors: Verify ANSI codes applied to border characters
- Width calculation: Content width + border thickness
- Height calculation: Content height + border thickness
- Multi-line content: Verify borders wrap all lines correctly
- Empty content: Handle gracefully (just borders)

**Test Examples**:
```go
func TestBorderNormal(t *testing.T) {
    s := NewStyle().
        BorderTop(true).BorderBottom(true).
        BorderLeft(true).BorderRight(true).
        BorderStyle(NormalBorder)
    result := s.Render("hello")

    lines := strings.Split(result, "\n")
    assert.Equal(t, 3, len(lines)) // top border + content + bottom border
    assert.Contains(t, lines[0], "┌") // top-left corner
    assert.Contains(t, lines[0], "┐") // top-right corner
    assert.Contains(t, lines[1], "│") // left border
    assert.Contains(t, lines[2], "└") // bottom-left corner
}

func TestBorderPartial(t *testing.T) {
    s := NewStyle().BorderTop(true).BorderBottom(true)
    result := s.Render("content")

    lines := strings.Split(result, "\n")
    assert.Equal(t, 3, len(lines))
    // No vertical bars, just horizontal lines
    assert.NotContains(t, result, "│")
}

func TestBorderColor(t *testing.T) {
    s := NewStyle().
        BorderTop(true).BorderLeft(true).
        BorderForeground("red")
    result := s.Render("content")

    // Border characters should have red ANSI code
    assert.Contains(t, result, "\x1b[31m") // red foreground
}
```

**Manual Testing**:
- Print bordered boxes to terminal
- Verify visual appearance of all border styles
- Test with different content sizes (short, long, multi-line)
- Verify on multiple terminals (iTerm2, Terminal.app)

**Visual Test Program**:
```go
func main() {
    styles := []struct{
        name string
        style BorderStyle
    }{
        {"Normal", NormalBorder},
        {"Rounded", RoundedBorder},
        {"Thick", ThickBorder},
        {"Double", DoubleBorder},
    }

    for _, test := range styles {
        s := NewStyle().
            BorderTop(true).BorderBottom(true).
            BorderLeft(true).BorderRight(true).
            BorderStyle(test.style)
        fmt.Println(test.name + ":")
        fmt.Println(s.Render("Hello World"))
        fmt.Println()
    }
}
```

## Notes

- Unicode box drawing characters (U+2500 - U+257F) widely supported in modern terminals
- Some terminals/fonts may not render box drawing characters correctly - document requirements
- Consider ASCII fallback option for compatibility (using +, -, |, etc.)
- Border rendering adds complexity to width/height calculations - ensure accuracy
- Borders wrap padding (if both enabled: content → padding → border)
- Border characters count as 1 cell width each (East Asian width: Neutral)
- Performance consideration: Border rendering adds minimal overhead (just character assembly)
- Future enhancement: Custom border characters via BorderStyle.Custom()
- Reference: Lipgloss library for border implementation patterns


