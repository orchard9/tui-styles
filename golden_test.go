package tuistyles

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldenSnapshots(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		input string
	}{
		{
			name:  "bold_text",
			style: NewStyle().Bold(true),
			input: "Bold Text",
		},
		{
			name:  "italic_text",
			style: NewStyle().Italic(true),
			input: "Italic Text",
		},
		{
			name: "colored_text",
			style: func() Style {
				red, _ := NewColor("#FF0000")
				blue, _ := NewColor("#0000FF")
				return NewStyle().Foreground(red).Background(blue)
			}(),
			input: "Colored Text",
		},
		{
			name: "padded_box",
			style: func() Style {
				return NewStyle().Padding(1).Border(NormalBorder())
			}(),
			input: "Padded",
		},
		{
			name: "colored_border",
			style: func() Style {
				red, _ := NewColor("red")
				return NewStyle().
					Padding(1).
					Border(RoundedBorder()).
					BorderForeground(red)
			}(),
			input: "Red Border",
		},
		{
			name: "aligned_center",
			style: func() Style {
				return NewStyle().
					Width(20).
					Align(Center).
					Border(RoundedBorder())
			}(),
			input: "Centered",
		},
		{
			name: "aligned_right",
			style: func() Style {
				return NewStyle().
					Width(20).
					Align(Right).
					Border(RoundedBorder())
			}(),
			input: "Right",
		},
		{
			name: "thick_border",
			style: func() Style {
				return NewStyle().
					Padding(1, 2).
					Border(ThickBorder())
			}(),
			input: "Thick",
		},
		{
			name: "double_border",
			style: func() Style {
				return NewStyle().
					Padding(1, 2).
					Border(DoubleBorder())
			}(),
			input: "Double",
		},
		{
			name: "block_border",
			style: func() Style {
				return NewStyle().
					Padding(1).
					Border(BlockBorder())
			}(),
			input: "Block",
		},
		{
			name: "multiline_box",
			style: func() Style {
				return NewStyle().
					Padding(1).
					Border(RoundedBorder())
			}(),
			input: "Line 1\nLine 2\nLine 3",
		},
		{
			name: "complex_style",
			style: func() Style {
				purple, _ := NewColor("#7D56F4")
				pink, _ := NewColor("#F72798")
				return NewStyle().
					Bold(true).
					Foreground(pink).
					Border(RoundedBorder()).
					BorderForeground(purple).
					Padding(2, 4).
					Width(30).
					Align(Center)
			}(),
			input: "Complex",
		},
		{
			name: "unicode_content",
			style: func() Style {
				return NewStyle().
					Border(DoubleBorder()).
					Padding(1)
			}(),
			input: "你好世界 🎉",
		},
		{
			name: "partial_border_top_bottom",
			style: func() Style {
				return NewStyle().
					Border(NormalBorder(), true, false, true, false).
					Padding(1)
			}(),
			input: "Top/Bottom",
		},
		{
			name: "partial_border_left_right",
			style: func() Style {
				return NewStyle().
					Border(NormalBorder(), false, true, false, true).
					Padding(1)
			}(),
			input: "Left/Right",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.style.Render(tt.input)
			goldenFile := filepath.Join("testdata", "golden", tt.name+".golden")

			if *update {
				// Test data directory permissions
				err := os.MkdirAll(filepath.Dir(goldenFile), 0750)
				require.NoError(t, err)
				// Test golden files are not sensitive
				err = os.WriteFile(goldenFile, []byte(output), 0600)
				require.NoError(t, err)
				t.Logf("Updated golden file: %s", goldenFile)
				return
			}

			// Test golden file path is controlled (not user input)
			expected, err := os.ReadFile(goldenFile) //nolint:gosec
			require.NoError(t, err, "failed to read golden file: %s", goldenFile)

			if output != string(expected) {
				t.Errorf("output mismatch:\n=== GOT ===\n%s\n=== WANT ===\n%s", output, expected)
			}
		})
	}
}
