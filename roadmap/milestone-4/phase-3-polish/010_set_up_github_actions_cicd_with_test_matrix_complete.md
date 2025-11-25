## Purpose

Establish automated CI/CD pipeline with GitHub Actions to run tests and linting across multiple Go versions and operating systems, ensuring code quality and cross-platform compatibility.

## Acceptance Criteria

- [ ] .github/workflows/test.yml created with test matrix
- [ ] .github/workflows/lint.yml created for linting
- [ ] Test matrix covers Go 1.21, 1.22, 1.23 on Linux, macOS, Windows
- [ ] All tests pass on all platforms
- [ ] Lint workflow runs golangci-lint with zero warnings
- [ ] Status badges added to README
- [ ] CI runs on pull requests and main branch pushes

## Technical Approach

**Test Workflow** (.github/workflows/test.yml):
```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        go-version: ['1.21', '1.22', '1.23']
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: go test -v -race -cover ./...
```

**Lint Workflow** (.github/workflows/lint.yml):
```yaml
name: Lint

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest
```

**Coverage Workflow** (optional, .github/workflows/coverage.yml):
- Run tests with coverage
- Upload to codecov.io
- Add coverage badge to README

**Files to Create/Modify**:
- .github/workflows/test.yml (new)
- .github/workflows/lint.yml (new)
- .github/workflows/coverage.yml (new, optional)
- README.md (add status badges)

**Dependencies**:
- GitHub Actions (built-in)
- golangci-lint action

## Testing Strategy

**Workflow Validation**:
- Push to branch and verify workflows run
- Intentionally break a test and verify CI catches it
- Intentionally introduce lint warning and verify lint workflow fails
- Test on all matrix combinations

**Badge Integration**:
```markdown
[![Test](https://github.com/user/repo/workflows/Test/badge.svg)](...)
[![Lint](https://github.com/user/repo/workflows/Lint/badge.svg)](...)
[![Coverage](https://codecov.io/gh/user/repo/badge.svg)](...)
```

## Notes

**Best Practices**:
- Run tests with -race flag to detect race conditions
- Use -cover to track coverage
- Cache Go modules for faster CI runs
- Set timeout-minutes to prevent hung tests
- Run on pull requests and main branch

**GitHub Actions Tips**:
```yaml
- uses: actions/cache@v4
  with:
    path: |
      ~/go/pkg/mod
      ~/.cache/go-build
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

**Status Badges**:
- Test badge: Shows pass/fail status
- Lint badge: Shows lint status
- Coverage badge: Shows coverage %
- Go version badge: Shows supported versions

**CI Optimization**:
- Cache dependencies
- Run tests in parallel
- Fast fail strategy for quick feedback
- Separate lint from test (faster feedback)


