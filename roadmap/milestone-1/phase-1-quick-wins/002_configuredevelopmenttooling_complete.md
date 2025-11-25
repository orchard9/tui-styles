## Purpose

Configure development tooling for code quality, linting, formatting, and testing. This ensures consistent code style, catches common errors, and establishes quality standards for the project.

## Acceptance Criteria

- [ ] `.golangci.yml` configured with strict linting rules
- [ ] `Makefile` created with common development commands
- [ ] `go fmt` formatting enforced
- [ ] Run `make lint` produces zero warnings
- [ ] Run `make test` executes all tests successfully
- [ ] Run `make build` compiles the project

## Technical Approach

**Tooling Setup**:

1. **golangci-lint Configuration** (`.golangci.yml`):
   - Enable linters: `gofmt`, `goimports`, `govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`
   - Disable noisy linters for library code
   - Set timeout to 5 minutes
   - Configure exclusions for test files where appropriate

2. **Makefile Targets**:
   ```makefile
   .PHONY: build test lint fmt clean

   build:
       go build ./...

   test:
       go test -v -race -coverprofile=coverage.out ./...

   lint:
       golangci-lint run --timeout=5m

   fmt:
       go fmt ./...
       goimports -w .

   clean:
       go clean
       rm -f coverage.out
   ```

3. **Pre-commit Hook** (optional, documented in CONTRIBUTING.md):
   - Run `make fmt lint test` before committing
   - Ensure quality checks pass locally

**Files to Create/Modify**:
- `.golangci.yml` - Linter configuration (already exists, verify settings)
- `Makefile` - Development workflow commands
- `.editorconfig` - Editor consistency (optional)

**Dependencies**:
- `golangci-lint` (install via `brew install golangci-lint` or Go)
- `goimports` (install via `go install golang.org/x/tools/cmd/goimports@latest`)

## Testing Strategy

**Validation Commands**:
```bash
# Verify golangci-lint is installed
golangci-lint --version

# Verify Makefile targets work
make fmt
make lint
make test
make build

# Check that linting catches issues (create a test violation)
echo "package main; var unused int" > test.go
make lint  # Should report unused variable
rm test.go
```

**Success Criteria**:
- All Makefile targets execute without errors
- Linting produces zero warnings on clean code
- Test command runs successfully (even with no tests yet)

## Notes

**golangci-lint Configuration**: Use existing `.golangci.yml` from project root, but verify it includes:
- `gofmt`: Enforce standard formatting
- `govet`: Catch suspicious constructs
- `errcheck`: Ensure errors are checked
- `staticcheck`: Advanced static analysis
- `unused`: Catch dead code

**Reference Projects**:
- See `reference-code/go-api/` for example `.golangci.yml`
- See `reference-code/go-worker/` for Makefile patterns

**CODING_GUIDELINES.md**: Ensure configuration aligns with project standards (no warnings allowed, conventional commits).


