# Project Foundation & Core Types

**Status**: Complete
**Owner**: Solo Developer
**Duration**: 1-2 weeks
**Dependencies**: []

## Goals

1. Establish project structure with proper Go module and directory organization
2. Define foundational types (Color, Position, BorderType) with validation
3. Configure development tooling for quality and consistency

## Success Criteria

- [ ] Go module initialized and builds successfully (go build ./...)
- [ ] All core types defined with proper validation and error handling
- [ ] Project structure matches specification (internal/ansi, internal/measure, examples/)
- [ ] Development tooling configured (golangci-lint, gofmt, go test)
- [ ] Zero linter warnings from golangci-lint
- [ ] Basic README with project overview and usage examples
- [ ] All code follows CODING_GUIDELINES.md standards

## Scope

**In Scope**:
- Go module initialization (go.mod, go.sum)
- Directory structure (internal/ansi, internal/measure, examples/)
- Core type definitions (Color, AdaptiveColor, Position, BorderType)
- Color validation and ANSI mapping functions
- Basic project documentation (README.md, CONTRIBUTING.md)
- Development tooling configuration (.golangci.yml, Makefile)
- Git repository setup (.gitignore)

**Out of Scope**:
- Style struct implementation (Milestone 2)
- Border rendering logic (Milestone 2)
- Layout utilities (JoinHorizontal, JoinVertical, Place) (Milestone 3)
- Advanced ANSI rendering (Milestone 2)
- Examples and benchmarks (will add incrementally)

## Technical Approach

**Go Module Structure**:
- Use Go 1.21+ for modern generics and error handling
- Follow standard Go project layout (internal/ for private packages)
- Use semantic versioning (start at v0.1.0 for development)

**Type Safety**:
- Implement Color as string type with validation
- Use iota-based enums for Position (Left, Center, Right, Top, Bottom)
- BorderType as struct with character definitions for each border piece

**Validation Strategy**:
- Color validation: Support hex (#RRGGBB), ANSI names (red, blue), ANSI codes (0-255)
- Use error returns for invalid inputs (no panics in library code)
- Helper functions for color conversion and normalization

**Testing Approach**:
- Unit tests for each core type
- Table-driven tests for color validation
- Test edge cases (empty strings, invalid hex, out-of-range ANSI codes)

## Risks

1. **Color Parsing Complexity**: Hex/ANSI name/code parsing may have edge cases
   - *Mitigation*: Comprehensive test suite with table-driven tests, reference lipgloss implementation

2. **Terminal Capability Detection**: AdaptiveColor requires background detection
   - *Mitigation*: Start with simple heuristic (check TERM env var), defer full detection to Phase 2

3. **ANSI Standard Variations**: Different terminals support different ANSI codes
   - *Mitigation*: Document supported codes, test on common terminals (iTerm2, Terminal.app, Alacritty)

## Timeline

**Phase 1 - Quick Wins**: 2-3 days
- Project structure setup
- Basic tooling configuration
- Initial documentation

**Phase 2 - Core Features**: 3-4 days
- Type definitions and validation
- Color parsing and conversion
- Unit tests

**Phase 3 - Polish**: 1-2 days
- Documentation completion
- Quality checks and refinement
- Example code

**Total**: 6-9 days (1-2 weeks allowing for iteration)

## Tasks by Phase

### Phase 1: Quick Wins
1. Initialize Go module and project structure
2. Configure development tooling
3. Create basic documentation structure

### Phase 2: Core Features
4. Implement Color type with validation
5. Implement AdaptiveColor type
6. Implement Position enum
7. Implement BorderType definitions
8. Add comprehensive unit tests

### Phase 3: Polish
9. Complete README with examples
10. Add contributing guidelines
11. Perform final quality checks
12. Create example programs
