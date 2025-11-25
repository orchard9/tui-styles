package tuistyles

import (
	"testing"
)

func TestNewColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Color
		wantErr bool
	}{
		// Valid hex colors
		{"valid hex 6-digit uppercase", "#FF0000", "#FF0000", false},
		{"valid hex 6-digit lowercase", "#ff0000", "#FF0000", false},
		{"valid hex 6-digit mixed", "#Ff0000", "#FF0000", false},
		{"valid hex 3-digit", "#F00", "#FF0000", false},
		{"valid hex 3-digit lowercase", "#f00", "#FF0000", false},
		{"valid hex black", "#000000", "#000000", false},
		{"valid hex white", "#FFFFFF", "#FFFFFF", false},

		// Valid ANSI names
		{"valid ANSI name red", "red", "red", false},
		{"valid ANSI name RED", "RED", "red", false},
		{"valid ANSI name Red", "Red", "red", false},
		{"valid ANSI name blue", "blue", "blue", false},
		{"valid ANSI name green", "green", "green", false},
		{"valid ANSI name yellow", "yellow", "yellow", false},
		{"valid ANSI name magenta", "magenta", "magenta", false},
		{"valid ANSI name cyan", "cyan", "cyan", false},
		{"valid ANSI name white", "white", "white", false},
		{"valid ANSI name black", "black", "black", false},
		{"valid ANSI bright-red", "bright-red", "bright-red", false},
		{"valid ANSI gray", "gray", "gray", false},
		{"valid ANSI grey", "grey", "grey", false},

		// Valid ANSI codes
		{"valid ANSI code 0", "0", "0", false},
		{"valid ANSI code 255", "255", "255", false},
		{"valid ANSI code 196", "196", "196", false},
		{"valid ANSI code 42", "42", "42", false},

		// Invalid cases
		{"invalid hex missing #", "FF0000", "", true},
		{"invalid hex characters", "#GGGGGG", "", true},
		{"invalid hex too short", "#FF", "", true},
		{"invalid hex too long", "#FFFFFFF", "", true},
		{"invalid ANSI name", "notacolor", "", true},
		{"invalid ANSI name typo", "redd", "", true},
		{"invalid ANSI code negative", "-1", "", true},
		{"invalid ANSI code too high", "256", "", true},
		{"invalid ANSI code too high 2", "300", "", true},
		{"empty string", "", "", true},
		{"whitespace only", "   ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NewColor(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestColorToANSI(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
	}{
		// ANSI color names
		{"ANSI red", "red", "\x1b[31m"},
		{"ANSI green", "green", "\x1b[32m"},
		{"ANSI blue", "blue", "\x1b[34m"},
		{"ANSI yellow", "yellow", "\x1b[33m"},
		{"ANSI magenta", "magenta", "\x1b[35m"},
		{"ANSI cyan", "cyan", "\x1b[36m"},
		{"ANSI white", "white", "\x1b[37m"},
		{"ANSI black", "black", "\x1b[30m"},
		{"ANSI bright-red", "bright-red", "\x1b[91m"},
		{"ANSI bright-blue", "bright-blue", "\x1b[94m"},
		{"ANSI gray", "gray", "\x1b[90m"},

		// Hex colors (RGB true color)
		{"hex red", "#FF0000", "\x1b[38;2;255;0;0m"},
		{"hex green", "#00FF00", "\x1b[38;2;0;255;0m"},
		{"hex blue", "#0000FF", "\x1b[38;2;0;0;255m"},
		{"hex black", "#000000", "\x1b[38;2;0;0;0m"},
		{"hex white", "#FFFFFF", "\x1b[38;2;255;255;255m"},
		{"hex 3-digit red", "#F00", "\x1b[38;2;255;0;0m"},
		{"hex custom", "#7D56F4", "\x1b[38;2;125;86;244m"},

		// ANSI 256-color codes
		{"code 0", "0", "\x1b[38;5;0m"},
		{"code 255", "255", "\x1b[38;5;255m"},
		{"code 196", "196", "\x1b[38;5;196m"},
		{"code 42", "42", "\x1b[38;5;42m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewColor(tt.color)
			if err != nil {
				t.Fatalf("NewColor(%q) failed: %v", tt.color, err)
			}
			got := c.ToANSI()
			if got != tt.want {
				t.Errorf("ToANSI() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColorToANSIBackground(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
	}{
		// ANSI color names
		{"ANSI red bg", "red", "\x1b[41m"},
		{"ANSI blue bg", "blue", "\x1b[44m"},
		{"ANSI bright-red bg", "bright-red", "\x1b[101m"},

		// Hex colors
		{"hex red bg", "#FF0000", "\x1b[48;2;255;0;0m"},
		{"hex blue bg", "#0000FF", "\x1b[48;2;0;0;255m"},

		// ANSI 256-color codes
		{"code 196 bg", "196", "\x1b[48;5;196m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewColor(tt.color)
			if err != nil {
				t.Fatalf("NewColor(%q) failed: %v", tt.color, err)
			}
			got := c.ToANSIBackground()
			if got != tt.want {
				t.Errorf("ToANSIBackground() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeHex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"3-digit lowercase", "#f00", "#FF0000"},
		{"3-digit uppercase", "#F00", "#FF0000"},
		{"3-digit mixed", "#Fa0", "#FFAA00"},
		{"6-digit lowercase", "#ff0000", "#FF0000"},
		{"6-digit uppercase", "#FF0000", "#FF0000"},
		{"6-digit mixed", "#Ff00Aa", "#FF00AA"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeHex(tt.input)
			if got != tt.want {
				t.Errorf("normalizeHex(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewAdaptiveColor(t *testing.T) {
	tests := []struct {
		name    string
		light   string
		dark    string
		wantErr bool
	}{
		{"valid hex colors", "#000000", "#FFFFFF", false},
		{"valid ANSI names", "black", "white", false},
		{"valid ANSI codes", "0", "255", false},
		{"mixed valid", "#FF0000", "blue", false},
		{"invalid light hex", "#GGGGGG", "#FFFFFF", true},
		{"invalid dark hex", "#000000", "#GGGGGG", true},
		{"invalid light name", "notacolor", "#FFFFFF", true},
		{"invalid dark name", "#000000", "notacolor", true},
		{"both invalid", "bad", "worse", true},
		{"empty light", "", "#FFFFFF", true},
		{"empty dark", "#000000", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAdaptiveColor(tt.light, tt.dark)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAdaptiveColor(%q, %q) error = %v, wantErr %v",
					tt.light, tt.dark, err, tt.wantErr)
			}
		})
	}
}

func TestAdaptiveColorToColor(t *testing.T) {
	ac, err := NewAdaptiveColor("#000000", "#FFFFFF")
	if err != nil {
		t.Fatalf("NewAdaptiveColor failed: %v", err)
	}

	t.Run("dark terminal", func(t *testing.T) {
		t.Setenv("TERM_BACKGROUND", "dark")

		got := ac.ToColor()
		want, _ := NewColor("#FFFFFF")
		if got != want {
			t.Errorf("ToColor() in dark terminal = %q, want %q", got, want)
		}
	})

	t.Run("light terminal", func(t *testing.T) {
		t.Setenv("TERM_BACKGROUND", "light")

		got := ac.ToColor()
		want, _ := NewColor("#000000")
		if got != want {
			t.Errorf("ToColor() in light terminal = %q, want %q", got, want)
		}
	})

	t.Run("no env vars (default dark)", func(t *testing.T) {
		// Clear env vars that might affect detection
		t.Setenv("TERM_BACKGROUND", "")
		t.Setenv("COLORFGBG", "")

		got := ac.ToColor()
		want, _ := NewColor("#FFFFFF") // Should use dark color
		if got != want {
			t.Errorf("ToColor() with no env = %q, want %q (dark default)", got, want)
		}
	})
}
