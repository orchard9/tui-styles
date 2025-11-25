## Purpose

Create CONTRIBUTING.md to guide future contributors on development workflow, coding standards, and contribution process. This ensures consistency and quality for external contributions.

## Acceptance Criteria

- [ ] CONTRIBUTING.md created with comprehensive guidelines
- [ ] Development setup instructions
- [ ] Code style and quality standards
- [ ] Testing requirements
- [ ] Pull request process
- [ ] Issue reporting guidelines
- [ ] Code of conduct reference

## Technical Approach

**CONTRIBUTING.md Structure**:

```markdown
# Contributing to TUI Styles

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- golangci-lint (for linting)
- git

### Development Setup

1. **Fork and Clone**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/tui-styles.git
   cd tui-styles
   ```

2. **Install Dependencies**:
   ```bash
   go mod download
   ```

3. **Install Development Tools**:
   ```bash
   # macOS
   brew install golangci-lint

   # Linux/Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

4. **Verify Setup**:
   ```bash
   make build
   make test
   make lint
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/my-feature
# OR
git checkout -b fix/bug-description
```

### 2. Make Changes

- Write code following our [coding standards](#coding-standards)
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detector
go test -race ./...

# Lint code
make lint

# Format code
make fmt
```

### 4. Commit Your Changes

Use conventional commit messages:

```bash
git commit -m "feat: add Color RGB conversion"
git commit -m "fix: correct Border character mapping"
git commit -m "docs: update README examples"
```

**Commit Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding/updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build/tooling changes

### 5. Submit Pull Request

1. Push your branch:
   ```bash
   git push origin feature/my-feature
   ```

2. Open a pull request on GitHub

3. Fill out the PR template with:
   - Description of changes
   - Related issue number (if applicable)
   - Test evidence (screenshots, test output)

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` and `goimports` for formatting
- Keep functions under 50 lines where possible
- Avoid deeply nested logic (max 3 levels)

### Documentation

- Add godoc comments for all exported types/functions
- Use complete sentences in comments
- Include examples in godoc where helpful

```go
// Color represents a terminal color (hex, ANSI name, or ANSI code).
//
// Color supports three formats:
//   - Hex colors: "#FF0000" or "#F00"
//   - ANSI names: "red", "blue", "green"
//   - ANSI codes: "0" through "255"
//
// Example:
//   red, err := NewColor("#FF0000")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(red.ToANSI() + "Red text" + "\x1b[0m")
type Color string
```

### Testing

- Write table-driven tests for validation logic
- Test both happy path and error cases
- Aim for >90% test coverage
- Use meaningful test names (`TestNewColorWithHex`, not `TestColor1`)
- Add benchmarks for performance-critical code

```go
func TestNewColor(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid hex", "#FF0000", false},
        {"invalid hex", "#GGGGGG", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewColor(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewColor(%q) error = %v, wantErr %v",
                    tt.input, err, tt.wantErr)
            }
        })
    }
}
```

### Error Handling

- Use descriptive error messages
- Wrap errors with context using `fmt.Errorf` with `%w`
- Don't panic in library code (return errors)

```go
if !isValidHex(s) {
    return "", fmt.Errorf("invalid hex color: %s", s)
}
```

## Quality Checklist

Before submitting a PR, ensure:

- [ ] All tests pass (`make test`)
- [ ] No linter warnings (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] New code has tests (>90% coverage)
- [ ] Public APIs have godoc comments
- [ ] CHANGELOG.md updated (if user-facing change)
- [ ] Examples updated (if API changed)

## Reporting Issues

### Bug Reports

Include:
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Minimal code example

### Feature Requests

Include:
- Use case description
- Proposed API (if applicable)
- Alternative solutions considered

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/0/code_of_conduct/).

## Questions?

- Open an issue for questions
- Tag with `question` label

Thank you for contributing!
```

**Files to Create/Modify**:
- `CONTRIBUTING.md` - Complete contributing guidelines

**Dependencies**:
- None (documentation only)

## Testing Strategy

**Review Checklist**:
- [ ] All instructions are accurate and testable
- [ ] Links to external resources are valid
- [ ] Commit message examples follow conventional commits
- [ ] Code examples compile and follow project standards
- [ ] Quality checklist covers all requirements
- [ ] Issue reporting template is helpful

**Validation**:
- Test development setup instructions on clean environment
- Verify all `make` commands work as documented
- Ensure commit message examples are correct

## Notes

**Conventional Commits**: Enforce conventional commit format for automated CHANGELOG generation (can add commitlint in future).

**Code of Conduct**: Link to Contributor Covenant rather than duplicating. Keep it simple for Phase 1.

**PR Template**: Consider adding `.github/pull_request_template.md` in future to standardize PR submissions.

**Issue Templates**: Consider adding `.github/ISSUE_TEMPLATE/` in future for bug reports and feature requests.

**Keep It Practical**: Focus on workflow and quality standards. Don't over-document edge cases.

**Reference Projects**: See `reference-code/go-api/CONTRIBUTING.md` for well-structured contribution guidelines.




