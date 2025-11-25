## Purpose

Ensure 100% test coverage for all builder methods with comprehensive test suites. This task fills any gaps left from individual task tests and adds integration scenarios combining multiple methods.

## Acceptance Criteria

- [ ] Test coverage report shows 100% coverage for all builder methods
- [ ] All edge cases covered (nil values, zero values, boundary conditions)
- [ ] Table-driven tests for systematic coverage
- [ ] Negative test cases (invalid inputs, panics)
- [ ] Performance benchmarks for common operations
- [ ] All tests pass with go test -race
- [ ] Test names clearly describe what they test

## Technical Approach

Review existing test files (text_test.go, color_test.go, layout_test.go, spacing_test.go, border_test.go) and identify gaps. Add missing tests in a comprehensive suite file or augment existing files.

**Coverage Analysis**:
```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
# Identify untested lines/branches
```

**Test Organization**:
1. **Per-Method Tests**: Each builder method has dedicated tests (should be mostly complete from tasks 002-007)
2. **Cross-Method Tests**: Test interactions between methods (e.g., Padding + Border)
3. **Edge Case Tests**: Nil checks, zero values, maximum values
4. **Error Path Tests**: Invalid inputs, panics with expected messages

**Missing Test Scenarios to Add**:

1. **Nil Safety**:
```go
func TestStyle_NilFieldsDefault(t *testing.T) {
    s := NewStyle()
    // Verify all fields are nil by default
    require.Nil(t, s.bold)
    require.Nil(t, s.foreground)
    require.Nil(t, s.width)
    // ... etc for all 30+ fields
}
```

2. **Zero Values**:
```go
func TestStyle_ZeroValues(t *testing.T) {
    s := NewStyle().Width(0).Height(0).Padding(0)
    // Verify zero values are set (not nil)
    require.NotNil(t, s.width)
    require.Equal(t, 0, *s.width)
}
```

3. **Maximum Values**:
```go
func TestStyle_MaximumDimensions(t *testing.T) {
    s := NewStyle().Width(10000).Height(10000)
    require.Equal(t, 10000, *s.width)
    require.Equal(t, 10000, *s.height)
}
```

4. **Method Combinations**:
```go
func TestStyle_CommonCombinations(t *testing.T) {
    // Test realistic usage patterns
    s := NewStyle().
        Bold(true).
        Foreground(Red).
        Padding(2).
        Border(Rounded).
        Width(80)

    require.NotNil(t, s.bold)
    require.NotNil(t, s.foreground)
    require.Equal(t, 2, *s.paddingTop)
    require.NotNil(t, s.borderType)
    require.Equal(t, 80, *s.width)
}
```

5. **Panic Tests**:
```go
func TestStyle_InvalidInputsPanic(t *testing.T) {
    tests := []struct {
        name string
        fn   func()
    }{
        {"padding 3 args", func() { NewStyle().Padding(1, 2, 3) }},
        {"margin 5 args", func() { NewStyle().Margin(1, 2, 3, 4, 5) }},
        {"border 2 edges", func() { NewStyle().Border(Rounded, true, false) }},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            require.Panics(t, tt.fn)
        })
    }
}
```

**Files to Create/Modify**:
- All existing test files (text_test.go, color_test.go, layout_test.go, spacing_test.go, border_test.go)
- style_test.go (add comprehensive integration tests)

**Dependencies**:
- All builder method implementations from tasks 002-007

## Testing Strategy

**Coverage Goals**:
- 100% line coverage for all builder methods
- 100% branch coverage (all if/switch cases tested)
- All panics tested with require.Panics()

**Test Execution**:
```bash
# Run all tests with race detector
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem
```

**Benchmarks to Add**:
```go
func BenchmarkStyle_SingleMethod(b *testing.B) {
    s := NewStyle()
    for i := 0; i < b.N; i++ {
        _ = s.Bold(true)
    }
}

func BenchmarkStyle_MethodChaining(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = NewStyle().
            Bold(true).
            Foreground(Red).
            Padding(2).
            Border(Rounded).
            Width(80)
    }
}

func BenchmarkStyle_Copy(b *testing.B) {
    s := NewStyle().Bold(true).Foreground(Red).Width(80)
    for i := 0; i < b.N; i++ {
        _ = s.Italic(true)
    }
}
```

## Notes

- Use `go test -cover` early and often to track coverage progress
- Aim for 100% but document any intentionally untested code with comments
- Table-driven tests are preferred for systematic coverage
- Benchmarks should show negligible allocation overhead (<500ns per method call)
- Run tests with -race on every commit to catch concurrency bugs early
- Consider using `gotestsum` for prettier test output

