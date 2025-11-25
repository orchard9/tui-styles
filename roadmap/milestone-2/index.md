# Style API & Builder Pattern

**Status**: Complete
**Owner**: TBD
**Duration**: 1-2 weeks
**Dependencies**: [milestone-1]

## Goals

1. Implement a complete, immutable Style struct with all text styling, layout, and border fields
2. Create a fluent builder API with 30+ methods supporting method chaining
3. Establish the copy-on-write pattern as the foundation for all future styling operations

## Success Criteria

- [ ] Style struct defined with all fields (text attributes, colors, sizing, alignment, spacing, borders, content)
- [ ] All builder methods implemented (Bold, Italic, Foreground, Padding, Border, etc.)
- [ ] Immutability verified - original Style unchanged after method calls
- [ ] Method chaining works fluently (e.g., `s.Bold(true).Foreground(Red).Padding(2)`)
- [ ] CSS-style shorthand methods handle 1/2/4 argument variations correctly
- [ ] Zero linter warnings (golangci-lint passes)
- [ ] 100% unit test coverage for builder methods
- [ ] Benchmarks show negligible allocation overhead from copying

## Scope

**In Scope**:
- Complete Style struct definition with 30+ fields
- Text attribute methods (Bold, Italic, Underline, Strikethrough, Faint, Blink, Reverse)
- Color methods (Foreground, Background, SetString)
- Layout methods (Width, Height, MaxWidth, MaxHeight, Align, AlignVertical)
- Spacing methods with CSS shorthand (Padding, Margin, individual edges)
- Border methods (Border, BorderForeground, BorderBackground, individual edges)
- NewStyle() constructor
- Input validation for all methods
- Unit tests for each method
- Immutability verification tests

**Out of Scope**:
- Rendering logic (Milestone 3)
- String width calculations (Milestone 3)
- ANSI escape sequence generation (Milestone 3)
- Advanced border characters/corners (Milestone 4)
- Style inheritance/composition (Milestone 5)
- Performance optimizations beyond basic copying (later milestones)

## Technical Approach

**Immutability Pattern**: All builder methods return a new Style with modified fields, never mutating the receiver. This follows functional programming principles and prevents spooky action at a distance.

**Copy-on-Write**:
```go
func (s Style) Bold(v bool) Style {
    s2 := s  // Shallow copy (all fields copied)
    s2.bold = &v  // Modify copy
    return s2
}
```

**Pointer Fields for Optionality**: Use pointers for boolean attributes and dimensions to distinguish "not set" (nil) from "explicitly false/zero". This allows proper zero-value defaults.

**CSS-Style Shorthand**: Padding/Margin methods accept variadic args:
- 1 arg: all sides
- 2 args: top/bottom, left/right
- 4 args: top, right, bottom, left

**Validation**: Methods validate inputs (no negative dimensions, nil Color checks) and return sensible defaults on invalid input rather than panicking.

**Testing Strategy**:
- Unit test each builder method in isolation
- Test method chaining combinations
- Verify immutability (assert original Style unchanged)
- Test CSS shorthand variations (1/2/4 args)
- Benchmark allocation overhead
- Table-driven tests for comprehensive input coverage

## Risks

1. **Copy Performance**: Shallow copying on every method call could impact performance if Style has many fields
   - **Mitigation**: Benchmark early. If problematic, consider arena allocation or copy-on-write optimization in later milestone. Most real-world usage won't chain 100+ methods.

2. **API Explosion**: 30+ builder methods is a large surface area
   - **Mitigation**: Group by category in source files (text.go, color.go, layout.go, spacing.go, border.go). Comprehensive tests catch inconsistencies.

3. **Nil Pointer Confusion**: Developers may not understand when fields are nil vs zero
   - **Mitigation**: Document clearly in godoc. Provide Has() methods in later milestone if needed.

4. **CSS Shorthand Edge Cases**: Padding(1, 2, 3) with 3 args is invalid CSS - how to handle?
   - **Mitigation**: Panic on invalid arg count (fail fast) or document as unsupported. Follow principle of least surprise.

## Timeline

**Phase 1 - Quick Wins**: 2 days
- Define Style struct skeleton
- Implement text attribute methods (7 methods, straightforward pattern)
- Basic tests proving immutability works

**Phase 2 - Core Features**: 5 days
- Color methods (Foreground, Background, SetString)
- Layout methods (Width, Height, MaxWidth, MaxHeight, Align, AlignVertical)
- Spacing methods (Padding/Margin shorthand + 8 individual edge methods)
- Border methods (Border + BorderForeground/Background + 4 edge toggles)
- Comprehensive unit tests for all methods

**Phase 3 - Polish**: 2 days
- Method chaining tests
- CSS shorthand edge case tests
- Benchmarks for copy overhead
- Linter fixes and godoc polish
- Integration test demonstrating full builder API

**Total**: 9 days (1.5 weeks with buffer)

## Tasks by Phase

### Phase 1: Quick Wins
[Tasks will be added here]

### Phase 2: Core Features
[Tasks will be added here]

### Phase 3: Polish
[Tasks will be added here]
