package ansi

import (
	"testing"
)

func TestIsValidANSIName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid red", "red", true},
		{"valid RED", "RED", true},
		{"valid Red", "Red", true},
		{"valid blue", "blue", true},
		{"valid bright-red", "bright-red", true},
		{"valid gray", "gray", true},
		{"valid grey", "grey", true},
		{"invalid notacolor", "notacolor", false},
		{"invalid redd", "redd", false},
		{"invalid empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidANSIName(tt.input)
			if got != tt.want {
				t.Errorf("IsValidANSIName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestColorToANSI(t *testing.T) {
	tests := []struct {
		name       string
		color      string
		background bool
		want       string
	}{
		// Foreground ANSI names
		{"fg red", "red", false, "\x1b[31m"},
		{"fg green", "green", false, "\x1b[32m"},
		{"fg blue", "blue", false, "\x1b[34m"},
		{"fg black", "black", false, "\x1b[30m"},
		{"fg white", "white", false, "\x1b[37m"},
		{"fg bright-red", "bright-red", false, "\x1b[91m"},
		{"fg bright-blue", "bright-blue", false, "\x1b[94m"},
		{"fg gray", "gray", false, "\x1b[90m"},

		// Background ANSI names
		{"bg red", "red", true, "\x1b[41m"},
		{"bg green", "green", true, "\x1b[42m"},
		{"bg blue", "blue", true, "\x1b[44m"},
		{"bg bright-red", "bright-red", true, "\x1b[101m"},

		// Foreground hex colors
		{"fg hex red", "#FF0000", false, "\x1b[38;2;255;0;0m"},
		{"fg hex green", "#00FF00", false, "\x1b[38;2;0;255;0m"},
		{"fg hex blue", "#0000FF", false, "\x1b[38;2;0;0;255m"},
		{"fg hex black", "#000000", false, "\x1b[38;2;0;0;0m"},
		{"fg hex white", "#FFFFFF", false, "\x1b[38;2;255;255;255m"},
		{"fg hex 3-digit", "#F00", false, "\x1b[38;2;255;0;0m"},

		// Background hex colors
		{"bg hex red", "#FF0000", true, "\x1b[48;2;255;0;0m"},
		{"bg hex blue", "#0000FF", true, "\x1b[48;2;0;0;255m"},

		// Foreground 256-color codes
		{"fg code 0", "0", false, "\x1b[38;5;0m"},
		{"fg code 255", "255", false, "\x1b[38;5;255m"},
		{"fg code 196", "196", false, "\x1b[38;5;196m"},

		// Background 256-color codes
		{"bg code 196", "196", true, "\x1b[48;5;196m"},
		{"bg code 42", "42", true, "\x1b[48;5;42m"},

		// Invalid cases
		{"invalid color", "invalid", false, ""},
		{"invalid code", "256", false, ""},
		{"invalid hex", "#GGGGGG", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColorToANSI(tt.color, tt.background)
			if got != tt.want {
				t.Errorf("ColorToANSI(%q, %v) = %q, want %q",
					tt.color, tt.background, got, tt.want)
			}
		})
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantR   int
		wantG   int
		wantB   int
		wantErr bool
	}{
		{"red", "#FF0000", 255, 0, 0, false},
		{"green", "#00FF00", 0, 255, 0, false},
		{"blue", "#0000FF", 0, 0, 255, false},
		{"black", "#000000", 0, 0, 0, false},
		{"white", "#FFFFFF", 255, 255, 255, false},
		{"custom", "#7D56F4", 125, 86, 244, false},
		{"3-digit red", "#F00", 255, 0, 0, false},
		{"3-digit green", "#0F0", 0, 255, 0, false},
		{"3-digit custom", "#FA0", 255, 170, 0, false},
		{"no hash", "FF0000", 255, 0, 0, false},
		{"invalid hex", "#GGGGGG", 0, 0, 0, true},
		{"too short", "#FF", 0, 0, 0, true},
		{"too long", "#FFFFFFF", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, err := hexToRGB(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("hexToRGB(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if r != tt.wantR || g != tt.wantG || b != tt.wantB {
					t.Errorf("hexToRGB(%q) = (%d, %d, %d), want (%d, %d, %d)",
						tt.hex, r, g, b, tt.wantR, tt.wantG, tt.wantB)
				}
			}
		})
	}
}

func TestReset(t *testing.T) {
	want := "\x1b[0m"
	got := Reset()
	if got != want {
		t.Errorf("Reset() = %q, want %q", got, want)
	}
}

func TestTextAttributes(t *testing.T) {
	tests := []struct {
		name string
		fn   func() string
		want string
	}{
		{"Bold", Bold, "\x1b[1m"},
		{"Faint", Faint, "\x1b[2m"},
		{"Italic", Italic, "\x1b[3m"},
		{"Underline", Underline, "\x1b[4m"},
		{"Blink", Blink, "\x1b[5m"},
		{"Reverse", Reverse, "\x1b[7m"},
		{"Strikethrough", Strikethrough, "\x1b[9m"},
		{"NoBold", NoBold, "\x1b[22m"},
		{"NoItalic", NoItalic, "\x1b[23m"},
		{"NoUnderline", NoUnderline, "\x1b[24m"},
		{"NoBlink", NoBlink, "\x1b[25m"},
		{"NoReverse", NoReverse, "\x1b[27m"},
		{"NoStrikethrough", NoStrikethrough, "\x1b[29m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn()
			if got != tt.want {
				t.Errorf("%s() = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestForegroundColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
	}{
		{"red name", "red", "\x1b[31m"},
		{"hex color", "#FF0000", "\x1b[38;2;255;0;0m"},
		{"256 color code", "196", "\x1b[38;5;196m"},
		{"bright-blue", "bright-blue", "\x1b[94m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ForegroundColor(tt.color)
			if got != tt.want {
				t.Errorf("ForegroundColor(%q) = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}

func TestBackgroundColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
	}{
		{"red name", "red", "\x1b[41m"},
		{"hex color", "#FF0000", "\x1b[48;2;255;0;0m"},
		{"256 color code", "196", "\x1b[48;5;196m"},
		{"bright-blue", "bright-blue", "\x1b[104m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BackgroundColor(tt.color)
			if got != tt.want {
				t.Errorf("BackgroundColor(%q) = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}
