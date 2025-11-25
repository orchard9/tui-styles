# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-23

### Added
- **Complete Style API** with 32 builder methods (Milestone 2)
  - Text attributes: `Bold`, `Italic`, `Underline`, `Strikethrough`, `Faint`, `Blink`, `Reverse`
  - Color methods: `Foreground`, `Background`, `BorderForeground`, `BorderBackground`
  - Layout methods: `Width`, `Height`, `MaxWidth`, `MaxHeight`, `Align`, `AlignVertical`
  - Spacing methods: `Padding`, `Margin` (with CSS-style shorthand), individual edge methods
  - Border methods: `Border`, `BorderTop`, `BorderRight`, `BorderBottom`, `BorderLeft`

- **Color support** (Milestone 1)
  - Hex colors: `#RRGGBB` and `#RGB` formats
  - ANSI color names: red, blue, green, etc. (case-insensitive)
  - ANSI 256-color codes: 0-255
  - Adaptive colors with terminal background detection

- **Borders** (Milestone 1)
  - 8 predefined border types: Normal, Rounded, Thick, Double, Block, OuterHalfBlock, InnerHalfBlock, Hidden
  - Unicode box drawing characters
  - Colored borders (foreground and background)
  - Partial borders (individual sides)

- **Rendering engine** (Milestone 3)
  - ANSI escape code generation
  - Unicode-aware width measurement (CJK, emoji support)
  - Multi-line rendering with `MaxWidth` constraints
  - Padding rendering with colored spaces
  - Border rendering with 8 border types
  - Horizontal and vertical alignment

- **Layout utilities** (Milestone 3)
  - `JoinHorizontal(pos Position, strs ...string)` - Side-by-side composition
  - `JoinVertical(pos Position, strs ...string)` - Vertical stacking
  - `Place(width, height, hPos, vPos Position, content string)` - Absolute positioning

- **Testing & Quality** (Milestone 4)
  - Comprehensive test suite (96% coverage)
  - Golden snapshot tests for visual validation
  - Performance benchmarks (all <1ms typical operations)
  - Zero linter warnings with golangci-lint strict rules
  - GitHub Actions CI/CD with test matrix (Go 1.21-1.23, Linux/macOS/Windows)

- **Documentation** (Milestone 4)
  - Complete package-level godoc with examples
  - Professional README with badges, installation, features, examples
  - 5 runnable examples: basic, borders, alignment, rendering, dashboard
  - CONTRIBUTING.md with development guidelines
  - API reference at pkg.go.dev

### Changed
- Initial stable release
- API frozen for v1.0 (breaking changes will require v2.0)

### Performance
- Simple text rendering: ~50ns per operation
- Border rendering: <1ms typical
- Layout composition: <5ms for multi-panel layouts
- Zero allocations for simple style operations

### Technical Details
- Go 1.21+ required
- Zero external dependencies (uses only Go standard library + mattn/go-runewidth for Unicode)
- Immutable builder pattern (all methods return new Style instances)
- Thread-safe (all operations are immutable)

## [Unreleased]

### Planned for v1.1
- More border styles (ASCII-only variants for compatibility)
- Gradient color support
- Shadow effects
- Animation utilities

### Planned for v2.0
- Breaking API changes (if needed based on community feedback)
- Advanced layout algorithms (grid, flexbox)
- Terminal capability detection
- Integration with other TUI frameworks
