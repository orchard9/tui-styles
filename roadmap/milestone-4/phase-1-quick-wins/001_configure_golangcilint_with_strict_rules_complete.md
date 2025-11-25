## Purpose

Establish strict linting standards to catch code quality issues, enforce Go idioms, and ensure consistent code style across the codebase before v1.0 release.

## Acceptance Criteria

- [ ] .golangci.yml configuration file created with strict linter rules
- [ ] Enabled linters include: govet, errcheck, staticcheck, unused, gosimple, ineffassign, gofmt, goimports
- [ ] Additional strict linters: revive (style), gocritic (diagnostics), gosec (security)
- [ ] All existing code passes linting with zero warnings
- [ ] Makefile includes `make lint` target for easy invocation
- [ ] CI-ready configuration (exit on first error for fast feedback)

## Technical Approach

**Linter Configuration**:
1. Create `.golangci.yml` with enabled linters and custom rules
2. Configure revive for Go style best practices
3. Enable gocritic for performance and style diagnostics
4. Add gosec for security vulnerability scanning
5. Set timeout to 5m for large codebases

**Strictness Levels**:
- Enable all default linters (govet, errcheck, staticcheck, unused, gosimple, ineffassign)
- Add style linters: gofmt, goimports, revive
- Add diagnostic linters: gocritic, unconvert, prealloc
- Add security linters: gosec
- Configure max-issues-per-linter: 0 (show all issues)

**Integration**:
- Add `make lint` to Makefile invoking `golangci-lint run`
- Document linter usage in CONTRIBUTING.md
- Prepare for GitHub Actions integration in later task

**Files to Create/Modify**:
- .golangci.yml (new)
- Makefile (add lint target)
- All Go source files (fix existing warnings)

**Dependencies**:
- golangci-lint binary (installation documented in README)

## Testing Strategy

**Validation**:
- Run `golangci-lint run` and verify zero warnings on current codebase
- Test on all packages: style, render, border, layout
- Verify lint failures are caught (intentionally introduce error and confirm detection)
- Confirm `make lint` works from project root

**CI Preparation**:
- Ensure configuration works in clean environment
- Test with golangci-lint v1.55+ (latest stable)

## Notes

**Recommended Linters**:
- **govet**: Go's official vet tool
- **errcheck**: Unchecked errors
- **staticcheck**: Advanced static analysis
- **unused**: Unused code detection
- **gosimple**: Simplification suggestions
- **ineffassign**: Ineffectual assignments
- **gofmt**: Code formatting
- **goimports**: Import organization
- **revive**: Fast, configurable, extensible Go linter
- **gocritic**: Opinionated linter with performance and style checks
- **gosec**: Security-focused linter

**Configuration Example**:
```yaml
linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - gofmt
    - goimports
    - revive
    - gocritic
    - gosec

run:
  timeout: 5m
  tests: true

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
```

**References**:
- golangci-lint documentation: https://golangci-lint.run/
- Recommended linters: https://golangci-lint.run/usage/linters/


