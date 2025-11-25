## Purpose

Validate that method chaining works fluently across all builder methods with realistic usage patterns. Ensures the fluent API ergonomics match developer expectations and that complex chains don't break immutability.

## Acceptance Criteria

- [ ] Integration tests cover realistic styling scenarios (10+ chained methods)
- [ ] Tests verify all method categories work together (text + color + layout + spacing + border)
- [ ] Tests verify chaining order doesn't matter for final result
- [ ] Tests demonstrate common UI patterns (buttons, cards, headers, alerts)
- [ ] All tests pass with go test -race
- [ ] Example code serves as documentation for users

## Technical Approach

Create a dedicated `integration_test.go` file with realistic styling scenarios that demonstrate the full builder API. Tests should read like documentation examples.

**Realistic Scenarios**:

1. **Button Style**:
```go
func TestIntegration_ButtonStyle(t *testing.T) {
    button := NewStyle().
        Bold(true).
        Foreground(White).
        Background(Blue).
        Padding(1, 3).  // Vertical: 1, Horizontal: 3
        Border(Rounded).
        Align(Center)

    // Verify all fields set correctly
    require.NotNil(t, button.bold)
    require.True(t, *button.bold)
    require.NotNil(t, button.foreground)
    require.NotNil(t, button.background)
    require.Equal(t, 1, *button.paddingTop)
    require.Equal(t, 3, *button.paddingLeft)
    require.NotNil(t, button.borderType)
    require.NotNil(t, button.align)
}
```

2. **Card Style**:
```go
func TestIntegration_CardStyle(t *testing.T) {
    card := NewStyle().
        Width(80).
        Height(20).
        Padding(2).
        Margin(1).
        Border(Double).
        BorderForeground(Gray).
        AlignVertical(Top)

    require.Equal(t, 80, *card.width)
    require.Equal(t, 20, *card.height)
    require.Equal(t, 2, *card.paddingTop)
    require.Equal(t, 1, *card.marginTop)
    require.NotNil(t, card.borderType)
    require.NotNil(t, card.borderForeground)
    require.NotNil(t, card.alignVertical)
}
```

3. **Alert Style**:
```go
func TestIntegration_AlertStyle(t *testing.T) {
    alert := NewStyle().
        Bold(true).
        Foreground(Red).
        Background(LightRed).
        PaddingLeft(2).
        BorderLeft(true).
        BorderForeground(DarkRed).
        MaxWidth(120)

    require.NotNil(t, alert.bold)
    require.NotNil(t, alert.foreground)
    require.NotNil(t, alert.background)
    require.Equal(t, 2, *alert.paddingLeft)
    require.NotNil(t, alert.borderLeft)
    require.NotNil(t, alert.borderForeground)
    require.Equal(t, 120, *alert.maxWidth)
}
```

4. **Header Style**:
```go
func TestIntegration_HeaderStyle(t *testing.T) {
    header := NewStyle().
        Bold(true).
        Underline(true).
        Foreground(Cyan).
        PaddingBottom(1).
        BorderBottom(true).
        BorderForeground(DarkCyan).
        Align(Center)

    require.NotNil(t, header.bold)
    require.NotNil(t, header.underline)
    require.NotNil(t, header.foreground)
    require.Equal(t, 1, *header.paddingBottom)
    require.NotNil(t, header.borderBottom)
    require.NotNil(t, header.borderForeground)
}
```

5. **Complex Chaining**:
```go
func TestIntegration_ComplexChaining(t *testing.T) {
    // Chain 15+ methods together
    s := NewStyle().
        Bold(true).
        Italic(true).
        Underline(true).
        Foreground(Blue).
        Background(White).
        Width(100).
        Height(30).
        MaxWidth(120).
        Padding(2, 4).
        Margin(1).
        Border(Rounded).
        BorderForeground(Gray).
        Align(Center).
        AlignVertical(Middle)

    // Verify all 14 operations applied correctly
    require.NotNil(t, s.bold)
    require.NotNil(t, s.italic)
    require.NotNil(t, s.underline)
    require.NotNil(t, s.foreground)
    require.NotNil(t, s.background)
    require.Equal(t, 100, *s.width)
    require.Equal(t, 30, *s.height)
    require.Equal(t, 120, *s.maxWidth)
    require.Equal(t, 2, *s.paddingTop)
    require.Equal(t, 4, *s.paddingLeft)
    require.Equal(t, 1, *s.marginTop)
    require.NotNil(t, s.borderType)
    require.NotNil(t, s.borderForeground)
    require.NotNil(t, s.align)
    require.NotNil(t, s.alignVertical)
}
```

6. **Order Independence**:
```go
func TestIntegration_OrderIndependence(t *testing.T) {
    // Build same style in different order
    s1 := NewStyle().
        Bold(true).
        Foreground(Red).
        Padding(2)

    s2 := NewStyle().
        Padding(2).
        Foreground(Red).
        Bold(true)

    // Results should be equivalent
    require.Equal(t, *s1.bold, *s2.bold)
    require.Equal(t, s1.foreground, s2.foreground)
    require.Equal(t, *s1.paddingTop, *s2.paddingTop)
}
```

7. **Branching Styles**:
```go
func TestIntegration_BranchingStyles(t *testing.T) {
    // Create base style and branch into variations
    baseButton := NewStyle().
        Padding(1, 3).
        Border(Rounded).
        Align(Center)

    primaryButton := baseButton.
        Foreground(White).
        Background(Blue)

    dangerButton := baseButton.
        Foreground(White).
        Background(Red)

    // Base should remain unchanged
    require.Nil(t, baseButton.foreground)
    require.Nil(t, baseButton.background)

    // Variants should have colors
    require.NotNil(t, primaryButton.foreground)
    require.NotNil(t, primaryButton.background)
    require.NotNil(t, dangerButton.foreground)
    require.NotNil(t, dangerButton.background)
}
```

**Files to Create/Modify**:
- integration_test.go (create new file with realistic scenarios)

**Dependencies**:
- All builder methods from tasks 002-007

## Testing Strategy

**Integration Scenarios**:
- Button style (bold, colors, padding, border, alignment)
- Card style (dimensions, spacing, border)
- Alert style (colors, left border accent, max width)
- Header style (text attributes, bottom border)
- Complex chaining (15+ methods)
- Order independence (different orders, same result)
- Branching styles (base style + variations)

**Test Goals**:
- Demonstrate real-world usage patterns
- Serve as documentation examples
- Verify immutability across complex chains
- Ensure fluent API ergonomics

**Code Quality**:
- Tests should read like documentation
- Use descriptive variable names (button, card, alert, header)
- Add comments explaining the styling choices
- Keep tests focused and readable

## Notes

- These tests serve dual purpose: validation and documentation
- Consider extracting test examples to package-level godoc
- Users should be able to copy-paste these patterns
- Integration tests catch regressions in method interactions
- If tests fail, check immutability (pointer aliasing issues)


