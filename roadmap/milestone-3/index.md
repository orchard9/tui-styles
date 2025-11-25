# Rendering Engine & Layout Utilities

**Status**: Complete
**Owner**: TBD
**Duration**: 2-3 weeks
**Dependencies**: [milestone-1, milestone-2]

## Goals

1. Transform styled strings into terminal-ready ANSI escape sequences
2. Implement accurate string width measurement for terminal layout
3. Build composable layout utilities for complex terminal UIs
4. Deliver production-ready rendering that handles all edge cases

## Success Criteria

- [ ] ANSI code generation correct for all text attributes and color formats
- [ ] Width measurement accurate (strips ANSI, handles Unicode/CJK/emojis)
- [ ] Style.Render() produces pixel-perfect terminal output
- [ ] Multi-line content styled correctly with padding, borders, alignment
- [ ] All 3 layout functions (JoinHorizontal, JoinVertical, Place) working
- [ ] Visual output matches specification examples
- [ ] Zero linter warnings, comprehensive test coverage
- [ ] Performance benchmarks show <1ms for typical rendering operations

## Scope

**In Scope**:
- ANSI escape code generation (internal/ansi/)
- String width measurement with ANSI stripping (internal/measure/)
- Core rendering engine (Style.Render(), Style.String())
- Multi-line handling with per-line styling
- Padding rendering (colored spaces)
- Border rendering (box drawing characters)
- Alignment (horizontal and vertical)
- Margin rendering
- Height handling (vertical padding/alignment)
- Background colors for content and padding
- Layout composition utilities (JoinHorizontal, JoinVertical, Place)
- Edge case handling (empty strings, zero dimensions, Unicode)

**Out of Scope**:
- Interactive features (cursor positioning, user input) - Milestone 4
- Advanced layout containers (flex, grid) - Future milestone
- Color transformations (darken, lighten, blend) - Future milestone
- Theme system and predefined styles - Milestone 5

## Technical Approach

**ANSI Generation Strategy**:
- Internal package `internal/ansi/` for escape code generation
- Map Color types (hex, ANSI names, ANSI codes) to sequences
- Efficient string building with minimal allocations
- Reset sequence optimization (only emit when needed)

**String Measurement Strategy**:
- Use mattn/go-runewidth for Unicode width calculation
- Regex-based ANSI escape code stripping before measurement
- Cache stripped strings to avoid redundant processing
- Handle edge cases (zero-width joiners, variation selectors)

**Rendering Architecture**:
- Builder pattern for constructing styled output
- Apply styles in layers: content → padding → borders → margin
- Per-line processing for multi-line content
- Vertical alignment via calculated padding distribution
- Horizontal alignment via space/truncation

**Layout Composition**:
- Treat styled strings as rectangular blocks
- JoinHorizontal: side-by-side with height normalization
- JoinVertical: stacking with width normalization
- Place: absolute positioning within bounding box

**Performance Considerations**:
- Minimize string allocations (use strings.Builder)
- Avoid redundant ANSI code emission
- Pre-compute dimensions before rendering
- Benchmark typical use cases (table rendering, boxes)

## Risks

1. **Unicode Width Calculation Accuracy**: Terminal width for complex Unicode (emojis, CJK) varies by terminal
   - Mitigation: Use battle-tested go-runewidth library, document known edge cases, provide override mechanism

2. **ANSI Code Compatibility**: Older terminals may not support all attributes (dim, strikethrough)
   - Mitigation: Generate standard ANSI codes, document compatibility, add capability detection in future milestone

3. **Performance for Large Content**: Rendering very long strings or many lines could be slow
   - Mitigation: Benchmark early, optimize hot paths, consider streaming rendering in future

4. **Border Character Rendering**: Box drawing characters may not render correctly in all terminals/fonts
   - Mitigation: Use standard Unicode box drawing, provide ASCII fallback option, document requirements

5. **Multi-line Alignment Complexity**: Vertical alignment with borders and padding has many edge cases
   - Mitigation: Comprehensive test suite with visual examples, incremental implementation

## Timeline

**Phase 1 - Quick Wins**: 3-4 days
- ANSI generation (simpler than it looks)
- String measurement foundation
- Basic rendering (no layout)

**Phase 2 - Core Features**: 7-9 days
- Multi-line handling
- Padding and borders
- Alignment logic
- Layout utilities

**Phase 3 - Polish**: 3-5 days
- Edge case handling
- Performance optimization
- Visual testing
- Documentation

**Total**: 13-18 days (2-3 weeks)

## Tasks by Phase

### Phase 1: Quick Wins (3-4 days)

Foundation layer - ANSI generation, measurement, basic rendering.

1. **Task 001**: ANSI escape code generation
   - Internal package for converting colors/attributes to ANSI sequences
   - Pure function package, zero external dependencies
   - Foundation for all visual styling

2. **Task 002**: String width measurement with ANSI stripping
   - Accurate width calculation (strips ANSI, handles Unicode/CJK/emojis)
   - Integrate go-runewidth library
   - Critical for layout calculations

3. **Task 003**: Basic Style.Render() implementation
   - Apply ANSI codes to strings
   - Multi-line handling (per-line styling)
   - No complex layout yet (just colors + attributes)

### Phase 2: Core Features (7-9 days)

Complete rendering engine with layout utilities.

4. **Task 004**: Multi-line rendering with line-by-line styling
   - Enhanced multi-line handling with width constraints
   - Truncation with ellipsis
   - Production-ready multi-line support

5. **Task 005**: Padding rendering (colored spaces)
   - Implement padding (top, bottom, left, right)
   - Colored padding spaces (background color)
   - Width/height calculations with padding

6. **Task 006**: Border rendering (box drawing characters)
   - Unicode box drawing characters (normal, rounded, thick, double)
   - Partial borders (individual sides)
   - Border colors separate from content

7. **Task 007**: Horizontal alignment implementation
   - Left, center, right alignment within width
   - Works with padding and borders
   - Space padding or truncation

8. **Task 008**: Vertical alignment and height handling
   - Top, center, bottom alignment within height
   - Vertical padding distribution
   - Height normalization

9. **Task 009**: JoinHorizontal layout utility
   - Side-by-side block composition
   - Vertical alignment (top, center, bottom)
   - Height normalization

10. **Task 010**: JoinVertical layout utility
    - Stacked block composition
    - Horizontal alignment (left, center, right)
    - Width normalization

### Phase 3: Polish (3-5 days)

Edge cases, optimization, testing, documentation.

11. **Task 011**: Place layout utility (absolute positioning)
    - Position content within bounding box
    - Horizontal and vertical positioning
    - Useful for overlays and absolute layouts

12. **Task 012**: Edge case handling and validation
    - Empty strings, zero dimensions
    - Invalid color values
    - Boundary conditions
    - Error handling and graceful degradation

13. **Task 013**: Performance optimization and benchmarks
    - Minimize allocations (strings.Builder)
    - ANSI code optimization
    - Benchmarks for typical use cases
    - Target: <1ms for typical rendering

14. **Task 014**: Visual testing and documentation
    - Visual test suite (render to terminal)
    - Example programs demonstrating features
    - API documentation (godoc)
    - Usage guide and best practices

## Dependencies

**Task Dependencies**:
- Task 003 depends on: 001, 002
- Task 004 depends on: 002, 003
- Task 005 depends on: 002, 003
- Task 006 depends on: 002, 004, 005
- Task 007 depends on: 002, 004
- Task 008 depends on: 002, 004, 007
- Task 009 depends on: 002, 004
- Task 010 depends on: 002, 004
- Task 011 depends on: 002, 009, 010
- Task 012 depends on: all core features (tasks 003-011)
- Task 013 depends on: all core features (tasks 003-011)
- Task 014 depends on: all tasks (001-013)

**Critical Path**: 001 → 002 → 003 → 004 → 005/006/007 → 008/009/010 → 011 → 012/013 → 014

**Parallel Work Opportunities**:
- Tasks 005, 006, 007 can be parallelized after task 004
- Tasks 009, 010 can be parallelized after task 004
- Tasks 012, 013 can be parallelized after task 011
