# Contributing to TUI Styles

Thank you for considering contributing to TUI Styles! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- **Go 1.21 or higher**: Check with `go version`
- **golangci-lint**: Install via `brew install golangci-lint` (macOS) or follow [installation guide](https://golangci-lint.run/usage/install/)
- **Git**: For version control

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/orchard9/tui-styles.git
   cd tui-styles
   ```

2. **Verify setup**:
   ```bash
   make test
   make lint
   ```

3. **Install pre-commit hooks** (optional but recommended):
   ```bash
   echo '#!/bin/bash\nmake fmt lint test' > .git/hooks/pre-commit
   chmod +x .git/hooks/pre-commit
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code style guidelines below

3. **Run tests frequently**:
   ```bash
   make test
   ```

4. **Format code**:
   ```bash
   make fmt
   ```

5. **Lint code**:
   ```bash
   make lint
   ```

6. **Commit changes**:
   ```bash
   git add .
   git commit -m "feat: add new border style"
   ```

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring (no behavior change)
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

Examples:
```
feat: add gradient color support
fix: correct ANSI color mapping for bright colors
docs: update README with usage examples
test: add edge case tests for Position enum
```

## Code Style Guidelines

### Go Conventions

- **Follow Go idioms**: Read [Effective Go](https://go.dev/doc/effective_go)
- **Use gofmt**: Code must be formatted with `gofmt` (enforced by CI)
- **Keep functions small**: Aim for < 50 lines per function
- **Avoid deep nesting**: Max 3 levels of indentation
- **Document exports**: All exported types, functions, and constants must have godoc comments

### Package Organization

- **Public API**: Export only what's necessary
- **Internal packages**: Use `internal/` for implementation details
- **Test files**: Place tests in `*_test.go` files alongside code

### Error Handling

- **Return errors**: Don't panic in library code (except for truly exceptional cases)
- **Wrap errors**: Use `fmt.Errorf("context: %w", err)` for context
- **Validate inputs**: Check inputs in constructors and return errors for invalid values

Example:
```go
func NewColor(s string) (Color, error) {
    if s == "" {
        return "", fmt.Errorf("color cannot be empty")
    }
    // ... validation logic
    return Color(s), nil
}
```

### Testing

- **Table-driven tests**: Use for multiple test cases
- **Test coverage**: Aim for > 90% coverage for new code
- **Test edge cases**: Empty strings, nil values, boundary conditions
- **Naming**: Use descriptive test names (`TestNewColor_EmptyString_ReturnsError`)

Example:
```go
func TestNewColor(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid hex", "#FF0000", false},
        {"empty string", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewColor(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Testing Requirements

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make coverage

# Run specific package tests
go test -v ./internal/ansi

# Run specific test
go test -v -run TestNewColor
```

### Test Requirements

- **All new code must have tests**
- **All tests must pass** before submitting PR
- **Zero race conditions**: Tests run with `-race` flag
- **Test coverage > 90%** for new code
- **No flaky tests**: Tests must be deterministic

## Code Review Process

### Before Submitting

- [ ] All tests passing (`make test`)
- [ ] Zero lint warnings (`make lint`)
- [ ] Code formatted (`make fmt`)
- [ ] Documentation updated (if API changed)
- [ ] CHANGELOG.md updated (for user-facing changes)

### Pull Request Guidelines

1. **Create PR** against `main` branch
2. **Describe changes** clearly in PR description
3. **Link issues** if applicable (Fixes #123)
4. **Keep PRs focused**: One feature/fix per PR
5. **Respond to feedback** promptly

### Review Checklist

Reviewers will check:
- Code quality and style
- Test coverage and quality
- Documentation completeness
- Performance implications
- API design (for public API changes)

## Release Process

(For maintainers)

1. **Update CHANGELOG.md** with release notes
2. **Update version** in documentation
3. **Tag release**: `git tag v0.1.0`
4. **Push tag**: `git push origin v0.1.0`
5. **Create GitHub release** with notes

## Getting Help

- **Questions**: Open a [Discussion](https://github.com/orchard9/tui-styles/discussions)
- **Bug reports**: Open an [Issue](https://github.com/orchard9/tui-styles/issues)
- **Feature requests**: Open an [Issue](https://github.com/orchard9/tui-styles/issues) with `enhancement` label

## Code of Conduct

- **Be respectful**: Treat all contributors with respect
- **Be constructive**: Provide helpful feedback
- **Be patient**: Everyone is learning
- **Be inclusive**: Welcome contributors of all backgrounds

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make TUI Styles better for everyone. We appreciate your time and effort!
