## Purpose

Implement the ANSI escape code generation engine that converts Style attributes and colors into terminal-ready escape sequences. This is the foundational layer that enables all visual styling in the library.

## Acceptance Criteria

- [ ] Text attribute codes generated correctly (bold, italic, underline, dim, strikethrough, blink, reverse, hidden)
- [ ] Foreground color codes generated for all color types (hex, ANSI name, ANSI code)
- [ ] Background color codes generated for all color types
- [ ] Reset sequences generated (full reset and attribute-specific)
- [ ] Efficient string building with minimal allocations
- [ ] Unit tests for all ANSI code generation functions
- [ ] Zero linter warnings

## Technical Approach

Create `internal/ansi/` package with escape code generation functions:

**Core Functions**:
- `Attribute(attr Attribute) string` - Convert text attribute to ANSI code
- `ForegroundColor(c Color) string` - Convert color to foreground ANSI sequence
- `BackgroundColor(c Color) string` - Convert color to background ANSI sequence
- `Reset() string` - Full style reset
- `ResetAttribute(attr Attribute) string` - Reset specific attribute

**Color Conversion**:
- Hex colors → RGB true color sequences: `\x1b[38;2;R;G;Bm`
- ANSI names → 16-color codes: `\x1b[30-37m` (foreground), `\x1b[40-47m` (background)
- ANSI codes → direct passthrough with wrapping

**Text Attributes**:
- Map Attribute enum to SGR codes:
  - Bold: `1`, Dim: `2`, Italic: `3`, Underline: `4`
  - Blink: `5`, Reverse: `7`, Hidden: `8`, Strikethrough: `9`

**Implementation Strategy**:
1. Create package structure `internal/ansi/ansi.go`
2. Implement color conversion functions (hex parsing, ANSI mapping)
3. Implement attribute conversion (simple enum → code mapping)
4. Add string building helpers (efficient concatenation)
5. Write comprehensive unit tests

**Files to Create/Modify**:
- internal/ansi/ansi.go (core generation logic)
- internal/ansi/ansi_test.go (unit tests)
- internal/ansi/doc.go (package documentation)

**Dependencies**:
- Standard library only (fmt, strings, strconv)
- No external dependencies for this package

## Testing Strategy

**Unit Tests**:
- Text attributes: Verify each attribute generates correct SGR code
- Hex colors: Test RGB extraction and true color sequence generation
- ANSI colors: Test all 16 named colors (foreground and background)
- ANSI codes: Test direct code passthrough
- Reset sequences: Verify full and attribute-specific resets
- Edge cases: Empty colors, invalid hex, boundary values

**Test Examples**:
```go
func TestAttribute(t *testing.T) {
    assert.Equal(t, "\x1b[1m", Attribute(Bold))
    assert.Equal(t, "\x1b[4m", Attribute(Underline))
}

func TestForegroundColor(t *testing.T) {
    // Hex color
    c := Color{Type: ColorTypeHex, Hex: "#FF5733"}
    assert.Equal(t, "\x1b[38;2;255;87;51m", ForegroundColor(c))

    // ANSI name
    c = Color{Type: ColorTypeANSI, ANSIName: "red"}
    assert.Equal(t, "\x1b[31m", ForegroundColor(c))
}
```

**Manual Testing**:
- Print colored text to terminal to verify visual output
- Test on multiple terminal emulators (iTerm2, Terminal.app, Linux terminals)

## Notes

- This is a pure function package with no side effects - easy to test
- Keep allocation-heavy operations (string building) efficient with strings.Builder
- Standard ANSI codes work on all modern terminals (Windows 10+, macOS, Linux)
- True color support (24-bit) is widely available but not universal - document graceful degradation
- Reference: ANSI/ECMA-48 specification for SGR (Select Graphic Rendition) codes

