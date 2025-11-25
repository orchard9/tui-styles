## Purpose

Implement accurate string width measurement that strips ANSI escape codes and correctly handles Unicode characters (CJK, emojis) for terminal layout calculations. This is critical for alignment, padding, and layout composition.

## Acceptance Criteria

- [ ] ANSI escape codes stripped correctly using regex
- [ ] Width calculation accurate for ASCII strings
- [ ] Width calculation accurate for Unicode (CJK characters count as 2 cells)
- [ ] Width calculation accurate for emojis and other wide characters
- [ ] go-runewidth library integrated and tested
- [ ] Unit tests covering all character types
- [ ] Performance acceptable (<10μs for typical strings)
- [ ] Zero linter warnings

## Technical Approach

Create `internal/measure/` package for string width measurement:

**Core Functions**:
- `Width(s string) int` - Calculate display width after stripping ANSI codes
- `StripANSI(s string) string` - Remove all ANSI escape sequences
- `Truncate(s string, maxWidth int) string` - Truncate to width while preserving ANSI codes
- `Pad(s string, width int, padChar rune) string` - Pad to width with ANSI-awareness

**ANSI Stripping Strategy**:
- Regex pattern: `\x1b\[[0-9;]*[a-zA-Z]` (matches SGR sequences)
- Compile regex once at package init for performance
- Handle edge cases (incomplete sequences, nested codes)

**Width Calculation Strategy**:
- Strip ANSI codes first
- Use `github.com/mattn/go-runewidth` for accurate Unicode width
- go-runewidth handles:
  - East Asian Wide and Fullwidth characters (2 cells)
  - Emojis and grapheme clusters (2 cells)
  - Zero-width joiners and combining characters (0 cells)
  - Ambiguous width characters (configurable)

**Truncate Strategy**:
- Calculate width while tracking ANSI code positions
- Insert ANSI codes at appropriate positions in truncated string
- Ensure reset codes are preserved if truncated mid-style

**Implementation Steps**:
1. Create `internal/measure/measure.go`
2. Implement ANSI regex and stripping function
3. Integrate go-runewidth for Unicode width
4. Implement Width() as strip + runewidth calculation
5. Implement Truncate() with ANSI preservation
6. Write comprehensive unit tests

**Files to Create/Modify**:
- internal/measure/measure.go (core measurement logic)
- internal/measure/measure_test.go (unit tests)
- internal/measure/doc.go (package documentation)
- go.mod (add go-runewidth dependency)
- go.sum (dependency checksums)

**Dependencies**:
- github.com/mattn/go-runewidth (v0.0.15 or later)
- Standard library: regexp, strings, unicode/utf8

## Testing Strategy

**Unit Tests**:
- ANSI stripping: Plain text, colored text, multiple ANSI codes
- Width measurement ASCII: Simple strings, empty strings
- Width measurement Unicode: CJK characters (Chinese, Japanese, Korean)
- Width measurement emojis: Single emoji, emoji sequences, skin tone modifiers
- Truncate: Truncate plain text, truncate with ANSI codes, truncate at exact width
- Edge cases: Empty strings, strings with only ANSI codes, invalid UTF-8

**Test Examples**:
```go
func TestWidth(t *testing.T) {
    // ASCII
    assert.Equal(t, 5, Width("hello"))

    // ANSI codes should be ignored
    assert.Equal(t, 5, Width("\x1b[31mhello\x1b[0m"))

    // CJK characters (2 cells each)
    assert.Equal(t, 6, Width("你好好"))  // 3 characters = 6 cells

    // Emoji (2 cells)
    assert.Equal(t, 2, Width("👍"))
}

func TestStripANSI(t *testing.T) {
    input := "\x1b[1;31mBold Red\x1b[0m"
    assert.Equal(t, "Bold Red", StripANSI(input))
}
```

**Performance Tests**:
- Benchmark Width() on typical strings (50-100 chars)
- Target: <10μs per call for 100-char strings
- Verify regex compilation happens once (not per call)

**Manual Testing**:
- Test with actual terminal output
- Verify alignment with CJK text
- Verify emoji rendering matches width calculation

## Notes

- go-runewidth is battle-tested (used by many TUI libraries)
- ANSI regex must be compiled at package init for performance
- Some terminals may render emojis differently (font/rendering engine variations)
- Truncate is tricky - needs careful handling of ANSI code positions
- Consider caching stripped strings if performance becomes an issue
- Reference: Unicode East Asian Width standard (UAX #11)
- Ambiguous width characters: Default to narrow (1 cell) - matches most terminals

