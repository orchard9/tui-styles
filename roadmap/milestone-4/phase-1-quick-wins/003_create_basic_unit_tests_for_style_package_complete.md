## Purpose

Establish foundational unit test coverage for the style package to validate API correctness, catch regressions, and provide confidence for refactoring.

## Acceptance Criteria

- [ ] style_test.go created with testify/assert
- [ ] Test coverage for Style creation (New, NewWithDefaults)
- [ ] Test coverage for style methods (Bold, Italic, Underline, Foreground, Background)
- [ ] Test coverage for method chaining and style composition
- [ ] Test coverage for ANSI code generation (ToANSI)
- [ ] Test coverage for edge cases (nil styles, empty strings, zero values)
- [ ] All tests pass with `go test -v ./style`
- [ ] Test coverage >70% for style package

## Technical Approach

**Test Structure**:
1. Create `style/style_test.go` using testify/assert
2. Group tests by functionality (creation, modifiers, rendering, edge cases)
3. Use table-driven tests for ANSI code validation
4. Test method chaining returns correct Style values

**Test Categories**:
- **Creation Tests**: New(), copy semantics, default values
- **Modifier Tests**: Bold(), Italic(), Underline(), Strikethrough()
- **Color Tests**: Foreground(), Background() with various color types
- **ANSI Tests**: ToANSI() generates correct escape codes
- **Chaining Tests**: Multiple methods chain correctly
- **Edge Cases**: nil styles, empty strings, zero dimensions

**Key Test Scenarios**:
```go
func TestStyle_Bold(t *testing.T) {
    s := style.New().Bold()
    assert.True(t, s.bold)
    assert.Contains(t, s.ToANSI(), "1") // Bold code
}

func TestStyle_Chaining(t *testing.T) {
    s := style.New().Bold().Foreground(color.Red).Underline()
    ansi := s.ToANSI()
    assert.Contains(t, ansi, "1")  // Bold
    assert.Contains(t, ansi, "31") // Red foreground
    assert.Contains(t, ansi, "4")  // Underline
}
```

**Files to Create/Modify**:
- style/style_test.go (new)

**Dependencies**:
- github.com/stretchr/testify/assert

## Testing Strategy

**Unit Tests**:
- Test each exported function and method
- Test internal ANSI code generation logic
- Validate Style immutability (methods return new Style)
- Test copy semantics for style composition

**Edge Case Tests**:
- nil Style handling
- Empty string rendering
- Default color values
- Multiple calls to same modifier (e.g., Bold().Bold())

**Coverage Validation**:
- Run `go test -cover ./style` and verify >70%
- Identify untested branches with `go test -coverprofile=coverage.out`
- Generate HTML coverage report: `go tool cover -html=coverage.out`

## Notes

**Testify Usage**:
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestStyle_New(t *testing.T) {
    s := style.New()
    assert.NotNil(t, s)
    assert.False(t, s.bold)
    assert.False(t, s.italic)
}
```

**Table-Driven Tests for ANSI Codes**:
```go
func TestStyle_ANSI(t *testing.T) {
    tests := []struct {
        name     string
        style    Style
        expected string
    }{
        {"bold", style.New().Bold(), "\x1b[1m"},
        {"red", style.New().Foreground(color.Red), "\x1b[31m"},
        {"bold red", style.New().Bold().Foreground(color.Red), "\x1b[1;31m"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, tt.style.ToANSI())
        })
    }
}
```

**Coverage Goals**:
- Focus on public API first
- Defer complex rendering tests to Phase 2
- Aim for >70% in Phase 1, >80% by end of milestone

**Reference**:
- testify documentation: https://github.com/stretchr/testify
- Go testing best practices: https://go.dev/doc/tutorial/add-a-test

