package layout

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// Lookup maps and helper functions live in split.templ (not this file) so
// Tailwind's content scanner can detect the class strings — it only scans
// *.templ files, not *.go files. See Grid, Card, Nav for the same pattern.

// SplitRatio is a typed enum for the relative widths of the Main and Aside
// columns in a Split. Unknown values fall back to SplitRatio1To3.
type SplitRatio string

const (
	// SplitRatio1To3 — Main takes 2/3, Aside takes 1/3 (default, classic
	// article+sidebar pattern).
	SplitRatio1To3 SplitRatio = "1-3"
	// SplitRatio1To4 — Main takes 3/4, Aside takes 1/4 (wide main, narrow aside).
	SplitRatio1To4 SplitRatio = "1-4"
	// SplitRatio1To2 — Main and Aside are equal width.
	SplitRatio1To2 SplitRatio = "1-2"
	// SplitRatioDefault is the canonical default (1To3).
	SplitRatioDefault SplitRatio = SplitRatio1To3
)

// SplitRatioIsValid reports whether r is a recognized SplitRatio.
func SplitRatioIsValid(r SplitRatio) bool {
	return r == SplitRatio1To2 || r == SplitRatio1To3 || r == SplitRatio1To4
}

// AsidePosition controls whether the Aside renders before (start) or after
// (end) the Main content. Both are logical positions — they auto-mirror in
// RTL contexts because the grid uses source order, not physical left/right.
type AsidePosition string

const (
	// AsidePositionEnd renders Aside AFTER Main (default — sidebar on the
	// right in LTR, left in RTL).
	AsidePositionEnd AsidePosition = "end"
	// AsidePositionStart renders Aside BEFORE Main (sidebar on the left in
	// LTR, right in RTL).
	AsidePositionStart AsidePosition = "start"
	// AsidePositionDefault is the canonical default (End).
	AsidePositionDefault AsidePosition = AsidePositionEnd
)

// AsidePositionIsValid reports whether p is a recognized AsidePosition.
func AsidePositionIsValid(p AsidePosition) bool {
	return p == AsidePositionStart || p == AsidePositionEnd
}

// SplitProps configures a 2-column content + aside layout — the second most
// common 2D pattern after AppShell. Use it for article+sidebar widgets,
// detail+metadata, or any main-content-with-secondary-content layout.
//
// Both columns use minmax(0,1fr) per grid track to prevent grid blowout
// from wide content (see ADR-0016). The layout collapses to a single
// stacked column below the md breakpoint (mobile-friendly by default).
type SplitProps struct {
	utils.BaseProps

	// Main is the primary content column. Required.
	Main templ.Component
	// Aside is the secondary content column (sidebar, metadata, related
	// links). Required.
	Aside templ.Component
	// AsidePosition controls whether Aside renders before or after Main.
	// Default AsidePositionEnd (aside after main in source order + visual).
	AsidePosition AsidePosition
	// Ratio controls the relative column widths. Default SplitRatio1To3
	// (main 2/3, aside 1/3).
	Ratio SplitRatio
	// Gap controls spacing between the two columns. Default GridGapLG (1.5rem).
	Gap string
	// ContainerAware, when true, wraps the Split in a @container div and swaps
	// md: breakpoint classes for @md: so the 2-column collapse responds to the
	// parent container's width, not the viewport. Useful when Split is placed
	// inside a sidebar, card body, drawer, or other constrained layout. Default
	// false (viewport breakpoints). See ADR-0018.
	ContainerAware bool
}

// DefaultSplitProps returns sensible defaults: aside at end, 1To3 ratio, gap-6.
//
//nolint:exhaustruct // constructor intentionally sets only non-zero defaults
func DefaultSplitProps() SplitProps {
	return SplitProps{
		AsidePosition: AsidePositionDefault,
		Ratio:         SplitRatioDefault,
		Gap:           "gap-6",
	}
}
