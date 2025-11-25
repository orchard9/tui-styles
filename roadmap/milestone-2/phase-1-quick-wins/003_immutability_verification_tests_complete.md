## Purpose

Create comprehensive tests that verify the immutability pattern works correctly across multiple method calls and chaining scenarios. This ensures the copy-on-write pattern is solid before building the remaining 20+ methods.

## Acceptance Criteria

- [ ] Test verifies single method call doesn't mutate original Style
- [ ] Test verifies method chaining doesn't mutate intermediate Styles
- [ ] Test verifies multiple independent branches from same Style don't affect each other
- [ ] Test verifies all pointer fields remain independent after copying
- [ ] Test demonstrates fluent API works (s.Bold(true).Italic(true))
- [ ] All tests pass with go test -race (no data races)

## Technical Approach

Create a dedicated test file `immutability_test.go` with comprehensive verification tests that prove the copy-on-write pattern works in all scenarios.

**Test Scenarios**:

1. **Single Method Immutability**:
```go
func TestImmutability_SingleMethod(t *testing.T) {
    original := NewStyle()
    modified := original.Bold(true)

    // Original should be unchanged
    require.Nil(t, original.bold)
    // Modified should have new value
    require.NotNil(t, modified.bold)
    require.True(t, *modified.bold)
}
```

2. **Method Chaining**:
```go
func TestImmutability_Chaining(t *testing.T) {
    s1 := NewStyle()
    s2 := s1.Bold(true)
    s3 := s2.Italic(true)

    // s1 should be unchanged
    require.Nil(t, s1.bold)
    require.Nil(t, s1.italic)

    // s2 should only have bold
    require.NotNil(t, s2.bold)
    require.Nil(t, s2.italic)

    // s3 should have both
    require.NotNil(t, s3.bold)
    require.NotNil(t, s3.italic)
}
```

3. **Independent Branches**:
```go
func TestImmutability_IndependentBranches(t *testing.T) {
    base := NewStyle()
    branch1 := base.Bold(true)
    branch2 := base.Italic(true)

    // Branches should be independent
    require.NotNil(t, branch1.bold)
    require.Nil(t, branch1.italic)

    require.Nil(t, branch2.bold)
    require.NotNil(t, branch2.italic)
}
```

4. **Fluent API**:
```go
func TestFluentAPI(t *testing.T) {
    // Verify method chaining compiles and works
    s := NewStyle().
        Bold(true).
        Italic(true).
        Underline(true)

    require.NotNil(t, s.bold)
    require.NotNil(t, s.italic)
    require.NotNil(t, s.underline)
}
```

**Files to Create/Modify**:
- immutability_test.go (create new file)

**Dependencies**:
- text.go (text attribute methods from task 002)
- github.com/stretchr/testify/require (for cleaner assertions)

## Testing Strategy

**Unit Tests**:
- TestImmutability_SingleMethod: Verify one method call doesn't mutate
- TestImmutability_Chaining: Verify chain creates independent Styles
- TestImmutability_IndependentBranches: Verify branches don't affect each other
- TestImmutability_AllTextAttributes: Verify all 7 text methods are immutable
- TestFluentAPI: Verify method chaining syntax works
- TestRaceConditions: Run with -race flag to detect data races

**Property-Based Tests** (optional, advanced):
```go
func TestProperty_ImmutabilityHolds(t *testing.T) {
    // Generate random sequences of method calls
    // Verify original always unchanged
    // Requires gopter or similar property testing library
}
```

## Notes

- Use `require` package for cleaner test assertions and early failure
- Run all tests with `go test -race` to detect concurrency issues
- These tests serve as documentation of immutability guarantees
- If tests fail, the copy-on-write pattern needs fixing before proceeding
- Consider adding benchmark to measure copy overhead (e.g., BenchmarkStyleCopy)


