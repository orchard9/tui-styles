## Purpose

Initialize the Go module and create the foundational directory structure for the TUI Styles library. This establishes the project foundation and ensures proper organization for internal packages, examples, and test files.

## Acceptance Criteria

- [ ] Go module initialized with `go mod init github.com/orchard9/tui-styles`
- [ ] Directory structure created: `internal/ansi/`, `internal/measure/`, `examples/`
- [ ] `.gitignore` file created with Go-specific exclusions
- [ ] Project builds successfully with `go build ./...`
- [ ] Go version set to 1.21+ in go.mod

## Technical Approach

**Directory Structure**:
```
tui-styles/
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── .gitignore                # Git exclusions
├── style.go                  # Public API (future task)
├── color.go                  # Color type (Phase 2)
├── position.go               # Position enum (Phase 2)
├── border.go                 # Border types (Phase 2)
├── internal/
│   ├── ansi/                 # ANSI escape code generation
│   │   └── codes.go          # ANSI color/attribute mappings
│   └── measure/              # String width calculation
│       └── measure.go        # Width measurement (ignoring ANSI)
└── examples/
    └── basic/                # Example programs (Phase 3)
```

**Files to Create/Modify**:
- `go.mod` - Module definition with Go 1.21+ requirement
- `go.sum` - Will be created by Go toolchain
- `.gitignore` - Exclude binaries, test coverage, IDE files
- `internal/ansi/.gitkeep` - Placeholder for directory
- `internal/measure/.gitkeep` - Placeholder for directory
- `examples/.gitkeep` - Placeholder for directory

**Commands**:
```bash
go mod init github.com/orchard9/tui-styles
mkdir -p internal/ansi internal/measure examples
touch internal/ansi/.gitkeep internal/measure/.gitkeep examples/.gitkeep
```

**Dependencies**:
- None (pure Go standard library for Phase 1)

## Testing Strategy

**Validation**:
- Run `go mod tidy` to verify module is valid
- Run `go build ./...` to ensure project compiles
- Verify directory structure matches specification
- Check `.gitignore` excludes expected files (test with `git status`)

**No Unit Tests Required**:
This is a structural task with no code logic to test.

## Notes

**Go Version Rationale**: Go 1.21+ provides:
- Improved error handling with `errors.Join`
- Better generics support (if needed for future work)
- Performance improvements in the compiler

**Internal Package Strategy**: Using `internal/` prevents external import of implementation details, following Go best practices for library design.

**Reference**: See CODING_GUIDELINES.md for Go project structure standards.


