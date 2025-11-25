package tuistyles

// Position represents horizontal or vertical alignment
type Position int

const (
	// Left represents left horizontal alignment
	Left Position = iota
	// Center represents center alignment (horizontal or vertical)
	Center
	// Right represents right horizontal alignment
	Right

	// Top represents top vertical alignment
	Top
	// Bottom represents bottom vertical alignment
	Bottom
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
	return p >= Left && p <= Bottom
}

// IsHorizontal returns true if position is Left, Center, or Right
func (p Position) IsHorizontal() bool {
	return p == Left || p == Center || p == Right
}

// IsVertical returns true if position is Top, Center, or Bottom
func (p Position) IsVertical() bool {
	return p == Top || p == Center || p == Bottom
}
