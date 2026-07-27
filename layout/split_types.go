package layout

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

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
	_, ok := splitRatioLookup[r]

	return ok
}

//nolint:gochecknoglobals // Package-level lookup table for split ratios
var splitRatioLookup = map[SplitRatio]string{
	SplitRatio1To2: "md:grid-cols-2",
	SplitRatio1To3: "md:grid-cols-3",
	SplitRatio1To4: "md:grid-cols-4",
}

// splitRatioLookupContainer mirrors splitRatioLookup with @container-keyed
// breakpoints (@md:) so container-aware Splits respond to their parent's
// width instead of the viewport. See ADR-0018.
var splitRatioLookupContainer = map[SplitRatio]string{
	SplitRatio1To2: "@md:grid-cols-2",
	SplitRatio1To3: "@md:grid-cols-3",
	SplitRatio1To4: "@md:grid-cols-4",
}

// splitRatioCols returns the grid-template-columns class for a ratio (e.g.
// md:grid-cols-2). Falls back to SplitRatioDefault. When containerAware is
// true, returns @container-keyed variants so the layout responds to the
// parent container's width, not the viewport.
func splitRatioCols(r SplitRatio, containerAware bool) string {
	if containerAware {
		return utils.Lookup(splitRatioLookupContainer, r, splitRatioLookupContainer[SplitRatioDefault])
	}
	return utils.Lookup(splitRatioLookup, r, splitRatioLookup[SplitRatioDefault])
}

// splitRatioMainSpanLookup maps each SplitRatio to the Main column's span
// class. Complete Tailwind literals (never dynamically concatenated —
// Tailwind's content scanner needs the full token).
//
//nolint:gochecknoglobals // Package-level lookup table
var splitRatioMainSpanLookup = map[SplitRatio]string{
	SplitRatio1To2: "md:col-span-1",
	SplitRatio1To3: "md:col-span-2",
	SplitRatio1To4: "md:col-span-3",
}

// splitRatioMainSpanLookupContainer mirrors splitRatioMainSpanLookup with
// @container-keyed breakpoints (@md:). See ADR-0018.
//
//nolint:gochecknoglobals // Package-level lookup table
var splitRatioMainSpanLookupContainer = map[SplitRatio]string{
	SplitRatio1To2: "@md:col-span-1",
	SplitRatio1To3: "@md:col-span-2",
	SplitRatio1To4: "@md:col-span-3",
}

// splitRatioMainSpan returns the col-span-N class for the Main column.
// At 1To3 (3 cols), Main spans 2. At 1To4 (4 cols), Main spans 3. At 1To2,
// Main spans 1 (equal). Both Main and Aside use minmax(0,1fr) per track to
// guard against grid blowout from wide content (ADR-0016).
func splitRatioMainSpan(r SplitRatio, containerAware bool) string {
	if containerAware {
		return utils.Lookup(splitRatioMainSpanLookupContainer, r, splitRatioMainSpanLookupContainer[SplitRatio1To3])
	}
	return utils.Lookup(splitRatioMainSpanLookup, r, splitRatioMainSpanLookup[SplitRatio1To3])
}

// splitAsideSpan returns the col-span-1 class for the Aside column, swapping
// viewport-keyed (md:) for container-keyed (@md:) when container-aware.
func splitAsideSpan(containerAware bool) string {
	if containerAware {
		return "@md:col-span-1"
	}
	return "md:col-span-1"
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
