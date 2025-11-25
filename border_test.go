package tuistyles

import (
	"testing"
)

func TestNormalBorder(t *testing.T) {
	border := NormalBorder()
	if border.Top != "─" || border.Bottom != "─" {
		t.Errorf("NormalBorder horizontal lines incorrect")
	}
	if border.Left != "│" || border.Right != "│" {
		t.Errorf("NormalBorder vertical lines incorrect")
	}
	if border.TopLeft != "┌" || border.TopRight != "┐" {
		t.Errorf("NormalBorder top corners incorrect")
	}
	if border.BottomLeft != "└" || border.BottomRight != "┘" {
		t.Errorf("NormalBorder bottom corners incorrect")
	}
}

func TestRoundedBorder(t *testing.T) {
	border := RoundedBorder()
	if border.TopLeft != "╭" || border.TopRight != "╮" {
		t.Errorf("RoundedBorder top corners not rounded")
	}
	if border.BottomLeft != "╰" || border.BottomRight != "╯" {
		t.Errorf("RoundedBorder bottom corners not rounded")
	}
}

func TestBlockBorder(t *testing.T) {
	border := BlockBorder()
	// All characters should be solid blocks
	expected := "█"
	fields := []string{
		border.Top, border.Bottom, border.Left, border.Right,
		border.TopLeft, border.TopRight, border.BottomLeft, border.BottomRight,
	}
	for i, field := range fields {
		if field != expected {
			t.Errorf("BlockBorder field %d = %q, want %q", i, field, expected)
		}
	}
}

func TestOuterHalfBlockBorder(t *testing.T) {
	border := OuterHalfBlockBorder()
	if border.Top != "▀" {
		t.Errorf("OuterHalfBlockBorder.Top = %q, want %q", border.Top, "▀")
	}
	if border.Bottom != "▄" {
		t.Errorf("OuterHalfBlockBorder.Bottom = %q, want %q", border.Bottom, "▄")
	}
	if border.Left != "▌" {
		t.Errorf("OuterHalfBlockBorder.Left = %q, want %q", border.Left, "▌")
	}
	if border.Right != "▐" {
		t.Errorf("OuterHalfBlockBorder.Right = %q, want %q", border.Right, "▐")
	}
}

func TestInnerHalfBlockBorder(t *testing.T) {
	border := InnerHalfBlockBorder()
	if border.Top != "▄" {
		t.Errorf("InnerHalfBlockBorder.Top = %q, want %q", border.Top, "▄")
	}
	if border.Bottom != "▀" {
		t.Errorf("InnerHalfBlockBorder.Bottom = %q, want %q", border.Bottom, "▀")
	}
}

func TestThickBorder(t *testing.T) {
	border := ThickBorder()
	if border.Top != "━" || border.Bottom != "━" {
		t.Errorf("ThickBorder horizontal lines not thick")
	}
	if border.Left != "┃" || border.Right != "┃" {
		t.Errorf("ThickBorder vertical lines not thick")
	}
	if border.TopLeft != "┏" || border.BottomRight != "┛" {
		t.Errorf("ThickBorder corners not thick")
	}
}

func TestDoubleBorder(t *testing.T) {
	border := DoubleBorder()
	if border.Top != "═" || border.Bottom != "═" {
		t.Errorf("DoubleBorder horizontal lines not double")
	}
	if border.Left != "║" || border.Right != "║" {
		t.Errorf("DoubleBorder vertical lines not double")
	}
	if border.TopLeft != "╔" || border.BottomRight != "╝" {
		t.Errorf("DoubleBorder corners not double")
	}
}

func TestHiddenBorder(t *testing.T) {
	border := HiddenBorder()
	// All characters should be spaces
	expected := " "
	fields := []string{
		border.Top, border.Bottom, border.Left, border.Right,
		border.TopLeft, border.TopRight, border.BottomLeft, border.BottomRight,
	}
	for i, field := range fields {
		if field != expected {
			t.Errorf("HiddenBorder field %d = %q, want %q (space)", i, field, expected)
		}
	}
}

func TestAllBordersHaveAllFields(t *testing.T) {
	borders := []struct {
		name string
		b    Border
	}{
		{"NormalBorder", NormalBorder()},
		{"RoundedBorder", RoundedBorder()},
		{"BlockBorder", BlockBorder()},
		{"OuterHalfBlockBorder", OuterHalfBlockBorder()},
		{"InnerHalfBlockBorder", InnerHalfBlockBorder()},
		{"ThickBorder", ThickBorder()},
		{"DoubleBorder", DoubleBorder()},
		{"HiddenBorder", HiddenBorder()},
	}

	for _, tt := range borders {
		t.Run(tt.name, func(t *testing.T) {
			// Verify no empty fields
			if tt.b.Top == "" || tt.b.Bottom == "" ||
				tt.b.Left == "" || tt.b.Right == "" ||
				tt.b.TopLeft == "" || tt.b.TopRight == "" ||
				tt.b.BottomLeft == "" || tt.b.BottomRight == "" {
				t.Errorf("%s has empty fields", tt.name)
			}
		})
	}
}
