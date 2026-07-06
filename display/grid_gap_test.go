package display

import (
	"fmt"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestGridGapIsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		gap  GridGap
		want bool
	}{
		{"valid sm", GridGapSM, true},
		{"valid md", GridGapMD, true},
		{"valid lg", GridGapLG, true},
		{"valid xl", GridGapXL, true},
		{"empty string", "", false},
		{"invalid huge", "huge", false},
		{"invalid tiny", "tiny", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GridGapIsValid(tt.gap); got != tt.want {
				t.Errorf("GridGapIsValid(%q) = %v, want %v", tt.gap, got, tt.want)
			}
		})
	}
}

func TestGridGapClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		gap  GridGap
		want string
	}{
		{"sm maps to gap-2", GridGapSM, "gap-2"},
		{"md maps to gap-4", GridGapMD, "gap-4"},
		{"lg maps to gap-6", GridGapLG, "gap-6"},
		{"xl maps to gap-8", GridGapXL, "gap-8"},
		{"unknown falls back to gap-4", GridGap("bogus"), "gap-4"},
		{"empty falls back to gap-4", "", "gap-4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("gridGapClass(%q)", tt.gap),
				gridGapClass(tt.gap),
				tt.want,
			)
		})
	}
}

func TestGridWithGap(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{Cols: GridCols3, Gap: GridGapLG}))
	utils.AssertContains(t, output, "gap-6")
}

func TestGridDefaultGap(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{Cols: GridCols3}))
	utils.AssertContains(t, output, "gap-4")
}
