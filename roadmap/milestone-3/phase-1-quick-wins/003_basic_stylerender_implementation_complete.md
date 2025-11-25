## Purpose

Implement the basic rendering engine that applies ANSI escape codes to strings based on Style configuration. This is the foundational rendering without complex layout features (no borders, padding, alignment yet - just colors and text attributes).

## Acceptance Criteria

- [ ] Style.Render(string) string method implemented
- [ ] Style.String() string method implemented (renders pre-set content)
- [ ] Text attributes applied correctly (bold, italic, underline, etc.)
- [ ] Foreground colors applied correctly (all color types)
- [ ] Background colors applied correctly (all color types)
- [ ] ANSI codes properly reset at end of styled strings
- [ ] Multi-line strings handled (split into lines, style per line)
- [ ] Unit tests for all rendering scenarios
- [ ] Zero linter warnings

## Technical Approach

Implement rendering methods in `render.go` at package root:

**Core Methods**:
```go
// Render applies the style to the given string
func (s Style) Render(str string) string

// String renders the pre-set content stored in the style
func (s Style) String() string
```

**Rendering Strategy**:
1. Build ANSI prefix (opening codes for attributes and colors)
2. Apply prefix to content
3. Build ANSI suffix (reset codes)
4. Return: prefix + content + suffix

**ANSI Code Assembly**:
- Collect all active attributes into slice
- Convert each attribute to ANSI code using `internal/ansi.Attribute()`
- Add foreground color code if set using `internal/ansi.ForegroundColor()`
- Add background color code if set using `internal/ansi.BackgroundColor()`
- Concatenate all codes: `\x1b[1m\x1b[31m` (bold + red)
- Alternatively: Combine into single sequence: `\x1b[1;31m` (more efficient)

**Multi-line Handling**:
- Split input string by newline (`\n`)
- Apply style to each line independently
- Rejoin with newlines
- This ensures ANSI codes don't span lines (some terminals require this)

**Reset Strategy**:
- Always append `\x1b[0m` (full reset) at end of styled string
- This ensures style doesn't bleed into subsequent output

**Implementation Steps**:
1. Create `render.go` at package root
2. Implement `buildANSIPrefix()` helper - collects all ANSI codes
3. Implement `Render(string)` - applies prefix/suffix to input
4. Implement `String()` - delegates to `Render(s.content)`
5. Handle multi-line strings (split, style per line, rejoin)
6. Write unit tests covering all style combinations

**Files to Create/Modify**:
- render.go (rendering methods)
- render_test.go (unit tests)
- style.go (if any helper methods needed)

**Dependencies**:
- internal/ansi (ANSI code generation - task 001)
- internal/measure (for future width calculations - task 002)
- Standard library: strings (for multi-line splitting)

## Testing Strategy

**Unit Tests**:
- Single attribute: Bold, italic, underline, etc.
- Multiple attributes: Bold + italic + underline
- Foreground color: Hex, ANSI name, ANSI code
- Background color: All color types
- Combined: Foreground + background + attributes
- Multi-line: String with newlines styled per line
- Empty strings: Should handle gracefully
- Style.String(): Render pre-set content

**Test Examples**:
```go
func TestRenderTextAttribute(t *testing.T) {
    s := NewStyle().Bold(true)
    result := s.Render("hello")
    expected := "\x1b[1mhello\x1b[0m"
    assert.Equal(t, expected, result)
}

func TestRenderColor(t *testing.T) {
    s := NewStyle().Foreground("#FF5733")
    result := s.Render("colored")
    assert.Contains(t, result, "\x1b[38;2;255;87;51m")
    assert.Contains(t, result, "colored")
    assert.Contains(t, result, "\x1b[0m")
}

func TestRenderMultiline(t *testing.T) {
    s := NewStyle().Bold(true)
    result := s.Render("line1\nline2")
    // Each line should be styled independently
    lines := strings.Split(result, "\n")
    assert.Equal(t, 2, len(lines))
    assert.Contains(t, lines[0], "\x1b[1m")
    assert.Contains(t, lines[1], "\x1b[1m")
}

func TestString(t *testing.T) {
    s := NewStyle().Bold(true).SetString("content")
    result := s.String()
    expected := "\x1b[1mcontent\x1b[0m"
    assert.Equal(t, expected, result)
}
```

**Manual Testing**:
- Print styled strings to terminal
- Verify visual appearance matches expected styling
- Test on multiple terminals (iTerm2, Terminal.app)
- Verify no style bleeding between consecutive prints

**Visual Test Program**:
```go
func main() {
    s := NewStyle().Bold(true).Foreground("red")
    fmt.Println(s.Render("Bold Red Text"))

    s2 := NewStyle().Background("yellow").Foreground("black")
    fmt.Println(s2.Render("Yellow Background"))
}
```

## Notes

- This task focuses on basic rendering only - no layout features yet
- Multi-line handling is simple for now (per-line styling)
- More complex multi-line rendering (with padding, borders) comes in Phase 2
- ANSI code optimization: Combine codes into single sequence for efficiency
- Some terminals may not support all attributes (dim, strikethrough) - document limitations
- Efficient string building with strings.Builder to minimize allocations
- Dependencies: Must complete tasks 001 (ANSI generation) and 002 (measurement) first

