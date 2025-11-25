## Purpose

Define the complete Style struct with all 30+ fields needed for text styling, layout, and borders. This establishes the data model for the entire builder API and ensures zero values work correctly as defaults.

## Acceptance Criteria

- [ ] Style struct defined in style.go with all fields
- [ ] Text attribute fields (bold, italic, underline, strikethrough, faint, blink, reverse) as *bool pointers
- [ ] Color fields (foreground, background) as *Color pointers
- [ ] Layout fields (width, height, maxWidth, maxHeight) as *int pointers
- [ ] Alignment fields (align, alignVertical) as *Position pointers
- [ ] Spacing fields (paddingTop/Right/Bottom/Left, marginTop/Right/Bottom/Left) as *int pointers
- [ ] Border fields (borderType as *BorderType, borderTop/Right/Bottom/Left as *bool, borderForeground/Background as *Color)
- [ ] NewStyle() constructor returns zero-value Style
- [ ] All pointer fields default to nil (not set)
- [ ] Struct documented with godoc explaining immutability pattern

## Technical Approach

Create a new `style.go` file with the complete Style struct definition. Use pointer fields for all optional attributes to distinguish "not set" (nil) from "explicitly false/zero".

**Struct Design**:
```go
type Style struct {
    // Text attributes
    bold          *bool
    italic        *bool
    underline     *bool
    strikethrough *bool
    faint         *bool
    blink         *bool
    reverse       *bool

    // Colors
    foreground *Color
    background *Color

    // Layout
    width     *int
    height    *int
    maxWidth  *int
    maxHeight *int

    // Alignment
    align         *Position
    alignVertical *Position

    // Spacing
    paddingTop    *int
    paddingRight  *int
    paddingBottom *int
    paddingLeft   *int
    marginTop     *int
    marginRight   *int
    marginBottom  *int
    marginLeft    *int

    // Borders
    borderType       *BorderType
    borderTop        *bool
    borderRight      *bool
    borderBottom     *bool
    borderLeft       *bool
    borderForeground *Color
    borderBackground *Color
}
```

**Constructor**:
```go
// NewStyle returns a new Style with all fields unset (nil).
// Styles are immutable - all builder methods return a new Style.
func NewStyle() Style {
    return Style{}
}
```

**Files to Create/Modify**:
- style.go (create new file)

**Dependencies**:
- color.go (Color type from milestone-1)
- position.go (Position type from milestone-1)
- border.go (BorderType from milestone-1)

## Testing Strategy

**Unit Tests** (in style_test.go):
- TestNewStyle_ZeroValue: Verify NewStyle() returns struct with all nil fields
- TestStyle_FieldCount: Use reflection to assert struct has expected number of fields (prevents accidental field removal)
- TestStyle_PointerFields: Verify all fields are pointers (optionality requirement)

**Documentation**:
- Add package-level godoc explaining immutability pattern
- Document each field group with comments
- Include example usage in godoc

## Notes

- Lowercase field names enforce encapsulation - only builder methods can modify
- Pointer fields add 8 bytes overhead per field but enable optionality
- Struct size will be ~256 bytes (32 fields * 8 bytes) - acceptable for copy-on-write
- Consider adding field groups as comments (// Text attributes, // Colors, etc.) for readability


