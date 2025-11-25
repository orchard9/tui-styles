## Purpose

Implement golden file snapshot testing to validate visual terminal output, ensuring rendering changes are intentional and preventing visual regressions.

## Acceptance Criteria

- [ ] Golden test framework set up with testify/golden or custom implementation
- [ ] Snapshot tests for all border styles (single, double, rounded, thick)
- [ ] Snapshot tests for text alignment variations
- [ ] Snapshot tests for complex compositions (nested boxes, styled text)
- [ ] Golden files committed to testdata/ directory
- [ ] Update mechanism for intentional changes (`UPDATE_GOLDEN=1 go test`)
- [ ] All golden tests pass on CI

## Technical Approach

**Golden Test Setup**:
1. Create testdata/ directory structure for golden files
2. Implement golden test helper using testify/golden or custom logic
3. Capture actual output and compare with golden files
4. Provide update mechanism for approved changes

**Directory Structure**:
```
testdata/
  golden/
    border_single.txt
    border_double.txt
    border_rounded.txt
    alignment_center.txt
    composition_nested.txt
```

**Golden Test Pattern**:
```go
func TestBorder_Single_Golden(t *testing.T) {
    output := border.Box("Hello", border.Single, 20, 5)
    golden.Assert(t, "border_single.txt", output)
}
```

**Test Coverage**:
- All border styles with various dimensions
- Text alignment at different widths
- Style combinations (bold + color + underline)
- Nested box compositions
- Dashboard-like complex layouts

**Files to Create/Modify**:
- testdata/golden/*.txt (golden files)
- internal/testing/golden/golden.go (helper)
- border/border_golden_test.go (new)
- render/render_golden_test.go (new)

**Dependencies**:
- github.com/stretchr/testify/assert (for golden helper)

## Testing Strategy

**Golden File Generation**:
- Run tests with UPDATE_GOLDEN=1 to create baseline
- Manually review golden files for correctness
- Commit golden files to version control
- Future runs compare against committed files

**Validation**:
- Golden tests fail when output doesn't match
- Diff output shows exact character differences
- Easy to review visual changes in PRs
- Works on CI (deterministic output)

**CI Integration**:
- Ensure testdata/ directory is committed
- Golden tests run on every PR
- Failures show diffs for review

## Notes

**Golden Test Implementation**:
```go
package golden

func Assert(t *testing.T, goldenFile string, actual string) {
    path := filepath.Join("testdata", "golden", goldenFile)
    if os.Getenv("UPDATE_GOLDEN") == "1" {
        os.WriteFile(path, []byte(actual), 0644)
        return
    }
    expected, _ := os.ReadFile(path)
    assert.Equal(t, string(expected), actual)
}
```

**Why Golden Tests**:
- Catch unintentional visual regressions
- Easy to review rendering changes in diffs
- Complement unit tests with visual validation


