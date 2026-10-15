package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestPositionIsValid(t *testing.T) {
	t.Parallel()

	valid := []Position{
		PositionTop, PositionBottom, PositionStart, PositionEnd,
		PositionLeft, PositionRight,
	}
	for _, pos := range valid {
		if !PositionIsValid(pos) {
			t.Errorf("PositionIsValid(%q) = false, want true", pos)
		}
	}
	if PositionIsValid(Position("diagonal")) {
		t.Error(`PositionIsValid("diagonal") = true, want false`)
	}
	if PositionIsValid(Position("")) {
		t.Error(`PositionIsValid("") = true, want false`)
	}
}

// TestPositionAliasesAreOneType pins the unification contract: the
// per-component position types are aliases of Position, and every legacy
// constant is the same value as its canonical spelling, so pre-unification
// consumer source keeps compiling and behaving identically.
func TestPositionAliasesAreOneType(t *testing.T) {
	t.Parallel()

	var typeIdentity = map[string]Position{
		"TooltipPositionTop":    TooltipPositionTop,
		"TooltipPositionBottom": TooltipPositionBottom,
		"TooltipPositionStart":  TooltipPositionStart,
		"TooltipPositionEnd":    TooltipPositionEnd,
		"TooltipPositionLeft":   TooltipPositionLeft,
		"TooltipPositionRight":  TooltipPositionRight,
		"PopoverPositionTop":    PopoverPositionTop,
		"PopoverPositionBottom": PopoverPositionBottom,
		"PopoverPositionStart":  PopoverPositionStart,
		"PopoverPositionEnd":    PopoverPositionEnd,
		"PopoverPositionLeft":   PopoverPositionLeft,
		"PopoverPositionRight":  PopoverPositionRight,
		"HoverCardPositionTop":    HoverCardPositionTop,
		"HoverCardPositionBottom": HoverCardPositionBottom,
		"HoverCardPositionStart":  HoverCardPositionStart,
		"HoverCardPositionEnd":    HoverCardPositionEnd,
	}

	canonical := map[string]Position{
		"TooltipPositionTop":    PositionTop,
		"TooltipPositionBottom": PositionBottom,
		"TooltipPositionStart":  PositionStart,
		"TooltipPositionEnd":    PositionEnd,
		"TooltipPositionLeft":   PositionLeft,
		"TooltipPositionRight":  PositionRight,
		"PopoverPositionTop":    PositionTop,
		"PopoverPositionBottom": PositionBottom,
		"PopoverPositionStart":  PositionStart,
		"PopoverPositionEnd":    PositionEnd,
		"PopoverPositionLeft":   PositionLeft,
		"PopoverPositionRight":  PositionRight,
		"HoverCardPositionTop":    PositionTop,
		"HoverCardPositionBottom": PositionBottom,
		"HoverCardPositionStart":  PositionStart,
		"HoverCardPositionEnd":    PositionEnd,
	}

	for name, legacy := range typeIdentity {
		if legacy != canonical[name] {
			t.Errorf("%s = %q, want canonical %q", name, legacy, canonical[name])
		}
	}

	// The wrapper IsValid functions are pure delegates of PositionIsValid.
	for _, pos := range []Position{"top", "bottom", "start", "end", "left", "right", "bogus"} {
		if TooltipPositionIsValid(pos) != PositionIsValid(pos) ||
			PopoverPositionIsValid(pos) != PositionIsValid(pos) ||
			HoverCardPositionIsValid(pos) != PositionIsValid(pos) {
			t.Errorf("wrapper IsValid disagreement at %q", pos)
		}
	}
}

// TestTooltipLogicalPositionsMatchLegacyInLTR pins the design promise: the
// logical start/end entries are pixel-identical to the legacy physical
// left/right entries in LTR, expressed with logical utilities so they mirror
// correctly in RTL.
func TestTooltipLogicalPositionsMatchLegacyInLTR(t *testing.T) {
	t.Parallel()

	start := tooltipLookupPosition(PositionStart)
	left := tooltipLookupPosition(PositionLeft)
	if start.Tooltip != "inset-e-full top-1/2 -translate-y-1/2 me-2" {
		t.Errorf("start tooltip classes = %q, want logical inset-e-full form", start.Tooltip)
	}
	if left.Tooltip != "right-full top-1/2 -translate-y-1/2 me-2" {
		t.Errorf("legacy left tooltip classes = %q, want unchanged physical form", left.Tooltip)
	}
	if start.Arrow != "inset-s-full top-1/2 -translate-y-1/2 border-s-gray-900 dark:border-s-gray-700" {
		t.Errorf("start arrow classes = %q, want logical inset-s-full form", start.Arrow)
	}
	if left.Arrow != "left-full top-1/2 -translate-y-1/2 border-s-gray-900 dark:border-s-gray-700" {
		t.Errorf("legacy left arrow classes = %q, want unchanged physical form", left.Arrow)
	}

	end := tooltipLookupPosition(PositionEnd)
	right := tooltipLookupPosition(PositionRight)
	if end.Tooltip != "inset-s-full top-1/2 -translate-y-1/2 ms-2" {
		t.Errorf("end tooltip classes = %q, want logical inset-s-full form", end.Tooltip)
	}
	if right.Tooltip != "left-full top-1/2 -translate-y-1/2 ms-2" {
		t.Errorf("legacy right tooltip classes = %q, want unchanged physical form", right.Tooltip)
	}
	if end.Arrow != "inset-e-full top-1/2 -translate-y-1/2 border-e-gray-900 dark:border-e-gray-700" {
		t.Errorf("end arrow classes = %q, want logical inset-e-full form", end.Arrow)
	}
	if right.Arrow != "right-full top-1/2 -translate-y-1/2 border-e-gray-900 dark:border-e-gray-700" {
		t.Errorf("legacy right arrow classes = %q, want unchanged physical form", right.Arrow)
	}
}

func TestDropdownAlign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pos  DropdownPosition
		want string
	}{
		{DropdownPositionStart, "start"},
		{DropdownPositionEnd, "end"},
		{DropdownPositionLeft, "start"},
		{DropdownPositionRight, "end"},
		{DropdownPosition(""), "start"},
	}
	for _, tt := range tests {
		if have := dropdownAlign(tt.pos); have != tt.want {
			t.Errorf("dropdownAlign(%q) = %q, want %q", tt.pos, have, tt.want)
		}
	}
}

func TestDropdownPositionIsValid(t *testing.T) {
	t.Parallel()

	for _, pos := range []DropdownPosition{
		DropdownPositionStart, DropdownPositionEnd,
		DropdownPositionLeft, DropdownPositionRight,
	} {
		if !DropdownPositionIsValid(pos) {
			t.Errorf("DropdownPositionIsValid(%q) = false, want true", pos)
		}
	}
	if DropdownPositionIsValid(DropdownPosition("center")) {
		t.Error(`DropdownPositionIsValid("center") = true, want false`)
	}
}

func TestDropdownAlignWire(t *testing.T) {
	t.Parallel()

	t.Run("end-aligned dropdown renders data-tc-align end", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			Label:    "Actions",
			Position: DropdownPositionEnd,
			Items:    []DropdownItem{{Text: "Edit", Href: "/edit"}},
		}))
		utils.AssertContains(t, output, `data-tc-align="end"`)
	})

	t.Run("start-aligned dropdown renders data-tc-align start", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			Label:    "Actions",
			Position: DropdownPositionStart,
			Items:    []DropdownItem{{Text: "Edit", Href: "/edit"}},
		}))
		utils.AssertContains(t, output, `data-tc-align="start"`)
	})

	t.Run("legacy left maps to start", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DefaultDropdownProps()))
		utils.AssertContains(t, output, `data-tc-align="start"`)
	})
}

// TestPopoverPositionJSResolvesLogicalSides pins the shared positioner script
// contract: start/end placements are resolved to physical left/right via the
// anchor's computed direction BEFORE the placement math, so the existing
// clamping and flip logic covers them unchanged.
func TestPopoverPositionJSResolvesLogicalSides(t *testing.T) {
	t.Parallel()

	js := popoverPositionJS()
	utils.AssertContains(t, js, `if(pos==="start"||pos==="end")`)
	utils.AssertContains(t, js, `var rtl=getComputedStyle(t).direction==="rtl"`)
	utils.AssertContains(t, js, `pos=(pos==="start")?(rtl?"right":"left"):(rtl?"left":"right")`)
	// The physical branches the resolution feeds into are still present.
	utils.AssertContains(t, js, `if(pos==="top")`)
	utils.AssertContains(t, js, `else if(pos==="left")`)
	utils.AssertContains(t, js, `else if(pos==="right")`)
}
