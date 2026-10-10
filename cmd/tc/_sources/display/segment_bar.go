// SegmentBar component: a stacked proportion bar with a legend — language
// shares, evidence breakdowns, storage composition. Extracted from
// mr-sync's language bars and nsfw-classifier's evidence proportions
// (2026-10-10 analysis, TODO_LIST #395).
package display

import (
	"strconv"
	"strings"

	"github.com/larsartmann/templ-components/utils"
)

// SegmentBarSegment is one proportion of the bar.
type SegmentBarSegment struct {
	// Label names the segment in the legend and the accessible summary.
	Label string
	// Value is the segment's magnitude (relative to the sum of all values).
	// Negative values clamp to zero.
	Value float64
	// Color optionally overrides the deterministic palette — a full class
	// pair such as "bg-emerald-600 dark:bg-emerald-400".
	Color string
}

// SegmentBarProps configures a stacked proportion bar.
type SegmentBarProps struct {
	utils.BaseProps

	// Segments stack in the given order. A bar with no positive-value
	// segments renders nothing (an empty proportion bar is no bar).
	Segments []SegmentBarSegment
	// NoLegend hides the dot+label+percent legend. The zero value shows it.
	NoLegend bool
}

// DefaultSegmentBarProps returns sensible defaults.
func DefaultSegmentBarProps() SegmentBarProps {
	return SegmentBarProps{} //nolint:exhaustruct // no non-zero defaults
}

//nolint:gochecknoglobals // Package-level lookup table: the shared chart palette as background class pairs
var segmentBarPalette = []string{
	"bg-blue-600 dark:bg-blue-400",
	"bg-emerald-600 dark:bg-emerald-400",
	"bg-amber-600 dark:bg-amber-400",
	"bg-rose-600 dark:bg-rose-400",
	"bg-violet-600 dark:bg-violet-400",
	"bg-cyan-600 dark:bg-cyan-400",
	"bg-orange-600 dark:bg-orange-400",
	"bg-pink-600 dark:bg-pink-400",
}

// segmentBarColor resolves a segment's background class: explicit override
// first, then the deterministic palette keyed by label (stable across
// renders and filter orders — the same label always gets the same color).
func segmentBarColor(seg SegmentBarSegment) string {
	if seg.Color != "" {
		return seg.Color
	}

	hash := 0

	for _, r := range seg.Label {
		hash = (hash*31 + int(r)) % len(segmentBarPalette)
	}

	return segmentBarPalette[hash]
}

// segmentBarTotal sums the clamped segment values.
func segmentBarTotal(segments []SegmentBarSegment) float64 {
	total := 0.0

	for _, seg := range segments {
		if seg.Value > 0 {
			total += seg.Value
		}
	}

	return total
}

// segmentBarPercent formats one segment's share of total with one decimal.
func segmentBarPercent(value, total float64) string {
	if total <= 0 {
		return "0.0%"
	}

	return strconv.FormatFloat(max(value, 0)/total*100, 'f', 1, 64) + "%"
}

// segmentBarAriaLabel composes the accessible summary ("Go 62.3%, TS 37.7%").
func segmentBarAriaLabel(segments []SegmentBarSegment, total float64) string {
	var parts []string

	for _, seg := range segments {
		if seg.Value > 0 && seg.Label != "" {
			parts = append(parts, seg.Label+" "+segmentBarPercent(seg.Value, total))
		}
	}

	return strings.Join(parts, ", ")
}
