package tuistyles

import (
	"testing"
)

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
			got := tt.pos.String()
			if got != tt.want {
				t.Errorf("Position(%d).String() = %q, want %q", tt.pos, got, tt.want)
			}
		})
	}
}

func TestPositionIsValid(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want bool
	}{
		{"Left valid", Left, true},
		{"Center valid", Center, true},
		{"Right valid", Right, true},
		{"Top valid", Top, true},
		{"Bottom valid", Bottom, true},
		{"Negative invalid", Position(-1), false},
		{"Too high invalid", Position(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.IsValid()
			if got != tt.want {
				t.Errorf("Position(%d).IsValid() = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}

func TestPositionIsHorizontal(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want bool
	}{
		{"Left is horizontal", Left, true},
		{"Center is horizontal", Center, true},
		{"Right is horizontal", Right, true},
		{"Top not horizontal", Top, false},
		{"Bottom not horizontal", Bottom, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.IsHorizontal()
			if got != tt.want {
				t.Errorf("Position(%s).IsHorizontal() = %v, want %v",
					tt.pos.String(), got, tt.want)
			}
		})
	}
}

func TestPositionIsVertical(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want bool
	}{
		{"Top is vertical", Top, true},
		{"Center is vertical", Center, true},
		{"Bottom is vertical", Bottom, true},
		{"Left not vertical", Left, false},
		{"Right not vertical", Right, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.IsVertical()
			if got != tt.want {
				t.Errorf("Position(%s).IsVertical() = %v, want %v",
					tt.pos.String(), got, tt.want)
			}
		})
	}
}

func TestPositionValues(t *testing.T) {
	// Verify constant values are unique
	positions := map[Position]string{
		Left:   "Left",
		Center: "Center",
		Right:  "Right",
		Top:    "Top",
		Bottom: "Bottom",
	}

	if len(positions) != 5 {
		t.Error("Position constants are not unique")
	}

	// Verify expected values
	if Left != 0 {
		t.Errorf("Left = %d, want 0", Left)
	}
	if Center != 1 {
		t.Errorf("Center = %d, want 1", Center)
	}
	if Right != 2 {
		t.Errorf("Right = %d, want 2", Right)
	}
	if Top != 3 {
		t.Errorf("Top = %d, want 3", Top)
	}
	if Bottom != 4 {
		t.Errorf("Bottom = %d, want 4", Bottom)
	}
}
