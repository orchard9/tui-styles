## Purpose

Implement the `Position` enum type for alignment and positioning. This type is used throughout the library for horizontal alignment (Left, Center, Right) and vertical alignment (Top, Center, Bottom).

## Acceptance Criteria

- [ ] `Position` type defined using iota-based enum pattern
- [ ] Constants defined: `Left`, `Center`, `Right`, `Top`, `Bottom`
- [ ] `String()` method for human-readable representation
- [ ] `IsValid()` method to check valid Position values
- [ ] Unit tests for all Position values
- [ ] All code passes `golangci-lint` with zero warnings

## Technical Approach

**Position Type Definition** (`position.go`):
```go
package tuistyles

// Position represents horizontal or vertical alignment
type Position int

const (
    // Horizontal positions
    Left Position = iota
    Center
    Right

    // Vertical positions (share Center with horizontal)
    Top    = iota + 100  // Offset to avoid overlap
    // Center is already defined (value 1)
    Bottom = iota + 100
)

// String returns human-readable position name
func (p Position) String() string {
    switch p {
    case Left:
        return "Left"
    case Center:
        return "Center"
    case Right:
        return "Right"
    case Top:
        return "Top"
    case Bottom:
        return "Bottom"
    default:
        return "Unknown"
    }
}

// IsValid checks if Position is a valid enum value
func (p Position) IsValid() bool {
    switch p {
    case Left, Center, Right, Top, Bottom:
        return true
    default:
        return false
    }
}

// IsHorizontal returns true if position is Left, Center, or Right
func (p Position) IsHorizontal() bool {
    return p == Left || p == Center || p == Right
}

// IsVertical returns true if position is Top, Center, or Bottom
func (p Position) IsVertical() bool {
    return p == Top || p == Center || p == Bottom
}
```

**Alternative: Separate H/V Types** (if Center overlap is confusing):
```go
type HPosition int
const (
    Left HPosition = iota
    HCenter
    Right
)

type VPosition int
const (
    Top VPosition = iota
    VCenter
    Bottom
)
```

**Files to Create/Modify**:
- `position.go` - Position enum definition
- `position_test.go` - Unit tests for Position

**Dependencies**:
- None (standalone type)

## Testing Strategy

**Unit Tests** (`position_test.go`):
```go
func TestPositionString(t *testing.T) {
    tests := []struct {
        pos  Position
        want string
    }{
        {Left, "Left"},
        {Center, "Center"},
        {Right, "Right"},
        {Top, "Top"},
        {Bottom, "Bottom"},
        {Position(999), "Unknown"},
    }

    for _, tt := range tests {
        t.Run(tt.want, func(t *testing.T) {
            if got := tt.pos.String(); got != tt.want {
                t.Errorf("Position.String() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestPositionIsValid(t *testing.T) {
    tests := []struct {
        pos  Position
        want bool
    }{
        {Left, true},
        {Center, true},
        {Right, true},
        {Top, true},
        {Bottom, true},
        {Position(999), false},
        {Position(-1), false},
    }

    for _, tt := range tests {
        t.Run(tt.pos.String(), func(t *testing.T) {
            if got := tt.pos.IsValid(); got != tt.want {
                t.Errorf("Position.IsValid() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestPositionTypeCheck(t *testing.T) {
    tests := []struct {
        pos        Position
        horizontal bool
        vertical   bool
    }{
        {Left, true, false},
        {Center, true, true},  // Center is both
        {Right, true, false},
        {Top, false, true},
        {Bottom, false, true},
    }

    for _, tt := range tests {
        t.Run(tt.pos.String(), func(t *testing.T) {
            if got := tt.pos.IsHorizontal(); got != tt.horizontal {
                t.Errorf("IsHorizontal() = %v, want %v", got, tt.horizontal)
            }
            if got := tt.pos.IsVertical(); got != tt.vertical {
                t.Errorf("IsVertical() = %v, want %v", got, tt.vertical)
            }
        })
    }
}
```

## Notes

**Design Decision: Shared Center vs Separate**:

**Option A: Shared Center** (simpler API):
- Single `Center` constant used for both H/V
- Less API surface area
- Slight risk of confusion (is it H or V?)

**Option B: Separate HCenter/VCenter** (more explicit):
- Two separate constants: `HCenter`, `VCenter`
- More verbose but crystal clear
- Requires separate HPosition/VPosition types

**Recommendation**: Start with **Option A (Shared Center)** for API simplicity. The context (Align vs AlignVertical) makes intent clear. Can refactor to Option B if user feedback indicates confusion.

**iota Offset Pattern**:
Using `iota + 100` for Top/Bottom ensures no value collision with Left/Center/Right. This allows Center to be shared.

**Type Safety**: Go doesn't enforce enum validation, so we provide `IsValid()` method. Consider using it in Style setters if strict validation is desired.

**Reference**: See `spec.md` Section 2.3 for Position usage in Align() and AlignVertical() methods.

**lipgloss Reference**: Review how [lipgloss handles positions](https://github.com/charmbracelet/lipgloss/blob/master/position.go).


