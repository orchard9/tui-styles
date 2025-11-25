## Purpose

Add comprehensive unit tests for all core types implemented in Phase 2. Ensure 100% coverage of validation logic, edge cases, and error handling. This establishes quality standards and prevents regressions.

## Acceptance Criteria

- [ ] All Color validation tests passing (hex, ANSI name, ANSI code)
- [ ] All AdaptiveColor terminal detection tests passing
- [ ] All Position enum tests passing
- [ ] All Border type tests passing
- [ ] Test coverage > 90% for all core type files
- [ ] All tests pass with `go test -race` (race detector)
- [ ] Benchmark tests for Color conversion performance
- [ ] All tests documented with clear test names

## Technical Approach

**Test Coverage Strategy**:

1. **Color Type Tests** (`color_test.go`):
   - Table-driven tests for NewColor validation
   - Test ToANSI conversion for all color formats
   - Test hex normalization (#RGB → #RRGGBB)
   - Test case-insensitive ANSI names
   - Test boundary conditions (ANSI code 0, 255, 256)
   - Benchmark hex to ANSI conversion

2. **AdaptiveColor Tests** (`color_test.go`):
   - Test NewAdaptiveColor validation
   - Test ToColor with mocked terminal detection
   - Test fallback to dark variant
   - Test both TERM_BACKGROUND and COLORFGBG env vars

3. **Position Tests** (`position_test.go`):
   - Test Position.String() for all values
   - Test Position.IsValid() for valid and invalid values
   - Test IsHorizontal() and IsVertical() helpers
   - Test invalid enum values

4. **Border Tests** (`border_test.go`):
   - Test all 8 border types are defined
   - Verify Unicode characters are correct
   - Ensure no empty fields (except HiddenBorder spaces)
   - Test border struct field access

**Test Organization**:
```go
// color_test.go
package tuistyles

import "testing"

// Validation tests
func TestNewColor(t *testing.T) { ... }
func TestNewColorInvalid(t *testing.T) { ... }
func TestColorNormalization(t *testing.T) { ... }

// Conversion tests
func TestColorToANSI(t *testing.T) { ... }
func TestColorToANSIEdgeCases(t *testing.T) { ... }

// AdaptiveColor tests
func TestNewAdaptiveColor(t *testing.T) { ... }
func TestAdaptiveColorToColor(t *testing.T) { ... }

// Benchmarks
func BenchmarkNewColor(b *testing.B) { ... }
func BenchmarkColorToANSI(b *testing.B) { ... }
```

**Coverage Target**: Aim for >90% coverage on all core type files
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Files to Create/Modify**:
- `color_test.go` - Expand with comprehensive tests
- `position_test.go` - Expand with comprehensive tests
- `border_test.go` - Expand with comprehensive tests
- `internal/ansi/codes_test.go` - Test ANSI conversion logic
- `internal/ansi/terminal_test.go` - Test terminal detection

**Dependencies**:
- Tasks 004, 005, 006, 007 must be complete
- Standard library: `testing`

## Testing Strategy

**Test Execution**:
```bash
# Run all tests
make test

# Run with race detector
go test -race ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestNewColor ./...

# Run benchmarks
go test -bench=. ./...
```

**Edge Cases to Cover**:

**Color Tests**:
- Empty string → error
- Whitespace only → error
- Invalid hex (#GGGGGG) → error
- Short hex (#F00) → normalize to #FF0000
- ANSI code 256 → error (out of range)
- ANSI code -1 → error (out of range)
- Unknown color name → error

**AdaptiveColor Tests**:
- Invalid light color → error
- Invalid dark color → error
- Missing env vars → default to dark
- Multiple env vars set → TERM_BACKGROUND takes precedence

**Position Tests**:
- Invalid enum value (999) → IsValid() = false
- Center is both horizontal and vertical

**Border Tests**:
- All 8 border types defined
- No empty fields (except intentional spaces)

**Performance Benchmarks**:
```go
func BenchmarkNewColor(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = NewColor("#FF0000")
    }
}

func BenchmarkColorToANSI(b *testing.B) {
    c, _ := NewColor("#FF0000")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = c.ToANSI()
    }
}
```

## Notes

**Table-Driven Tests**: Use table-driven test pattern for all validation tests. This makes it easy to add new test cases and improves readability.

**Test Naming**: Follow Go conventions:
- `TestNewColor` - Happy path
- `TestNewColorInvalid` - Error cases
- `TestColorToANSI` - Conversion logic
- `BenchmarkNewColor` - Performance tests

**Race Detector**: Always run tests with `-race` flag to catch concurrency issues (even though Phase 1 code is mostly synchronous).

**Test Fixtures**: Consider creating test helper functions for common setup:
```go
func mustColor(t *testing.T, s string) Color {
    t.Helper()
    c, err := NewColor(s)
    if err != nil {
        t.Fatalf("failed to create color %q: %v", s, err)
    }
    return c
}
```

**Coverage Tools**:
- Run `go tool cover -html=coverage.out` to visualize coverage
- Aim for >90% coverage, but don't sacrifice test quality for coverage percentage
- Uncovered lines should be documented (e.g., unreachable error paths)

**CI Integration**: Ensure tests run in CI pipeline:
```yaml
# .github/workflows/test.yml
- name: Run tests
  run: |
    go test -race -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out
```

**Reference**: See existing test files in tasks 004-007 for test patterns.



