// SegmentedControl component: a single-select segmented button group in two
// zero-JS modes — link mode (server-round-trip filtering, the FilterChips
// sibling with segmented visuals) and radio mode (a native radiogroup that
// submits with its enclosing form, e.g. FilterBar). Extracted from
// dnsblockd's duration/preset button groups and nsfw-classifier's threshold
// presets (2026-10-10 analysis, TODO_LIST #394).
package forms

import (
	"github.com/larsartmann/templ-components/utils"
)

// SegmentedControlItem is one option.
type SegmentedControlItem struct {
	// Label is the visible text.
	Label string
	// Value is the option's submitted value (radio mode).
	Value string
	// Href, when set on ANY item, switches the whole control to link mode —
	// every item should carry one (a mixed control is a consumer bug;
	// items without an Href render inert links in link mode).
	Href string
}

// SegmentedControlProps configures a segmented single-select group.
type SegmentedControlProps struct {
	utils.BaseProps

	// Name is the radio group's field name (radio mode). Ignored in link
	// mode.
	Name string
	// Items render in FORWARD DOM order (first item first) so radiogroup
	// arrow keys move down the value order (the Rating rule).
	Items []SegmentedControlItem
	// ActiveValue marks the selected item: the checked radio (radio mode)
	// or aria-current link (link mode).
	ActiveValue string
}

// DefaultSegmentedControlProps returns sensible defaults.
func DefaultSegmentedControlProps() SegmentedControlProps {
	return SegmentedControlProps{ //nolint:exhaustruct_v5 // no non-zero defaults
		Name: "segment",
	}
}

// segmentedControlIsLinkMode reports whether any item carries an Href.
func segmentedControlIsLinkMode(items []SegmentedControlItem) bool {
	for _, item := range items {
		if item.Href != "" {
			return true
		}
	}

	return false
}
