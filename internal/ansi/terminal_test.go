package ansi

import (
	"testing"
)

func TestIsLightTerminal(t *testing.T) {
	tests := []struct {
		name        string
		termBg      string
		colorFgBg   string
		want        bool
		description string
	}{
		{
			name:        "explicit light",
			termBg:      "light",
			colorFgBg:   "",
			want:        true,
			description: "TERM_BACKGROUND=light should return true",
		},
		{
			name:        "explicit Light (case insensitive)",
			termBg:      "Light",
			colorFgBg:   "",
			want:        true,
			description: "TERM_BACKGROUND=Light should return true",
		},
		{
			name:        "explicit dark",
			termBg:      "dark",
			colorFgBg:   "",
			want:        false,
			description: "TERM_BACKGROUND=dark should return false",
		},
		{
			name:        "explicit Dark (case insensitive)",
			termBg:      "Dark",
			colorFgBg:   "",
			want:        false,
			description: "TERM_BACKGROUND=Dark should return false",
		},
		{
			name:        "COLORFGBG light background",
			termBg:      "",
			colorFgBg:   "0;15",
			want:        true,
			description: "COLORFGBG with background > 6 should return true",
		},
		{
			name:        "COLORFGBG light background (7)",
			termBg:      "",
			colorFgBg:   "0;7",
			want:        true,
			description: "COLORFGBG with background = 7 should return true",
		},
		{
			name:        "COLORFGBG dark background",
			termBg:      "",
			colorFgBg:   "15;0",
			want:        false,
			description: "COLORFGBG with background <= 6 should return false",
		},
		{
			name:        "COLORFGBG dark background (6)",
			termBg:      "",
			colorFgBg:   "15;6",
			want:        false,
			description: "COLORFGBG with background = 6 should return false",
		},
		{
			name:        "COLORFGBG invalid format",
			termBg:      "",
			colorFgBg:   "invalid",
			want:        false,
			description: "Invalid COLORFGBG should return false (default)",
		},
		{
			name:        "COLORFGBG single value",
			termBg:      "",
			colorFgBg:   "15",
			want:        false,
			description: "COLORFGBG with single value should return false (default)",
		},
		{
			name:        "no env vars",
			termBg:      "",
			colorFgBg:   "",
			want:        false,
			description: "No env vars should return false (default to dark)",
		},
		{
			name:        "TERM_BACKGROUND overrides COLORFGBG",
			termBg:      "dark",
			colorFgBg:   "0;15",
			want:        false,
			description: "TERM_BACKGROUND should take precedence over COLORFGBG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.termBg != "" {
				t.Setenv("TERM_BACKGROUND", tt.termBg)
			}
			if tt.colorFgBg != "" {
				t.Setenv("COLORFGBG", tt.colorFgBg)
			}

			got := IsLightTerminal()
			if got != tt.want {
				t.Errorf("IsLightTerminal() = %v, want %v\n%s",
					got, tt.want, tt.description)
			}
		})
	}
}
