// DataState component: the honesty ladder for data-driven regions —
// disabled → unavailable → empty → content — extracted from dnsblockd's
// per-card ladders (9 accepted-duplication markers across 4 view files,
// 2026-10-10 analysis, TODO_LIST #393).
package display

import (
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// DataStateState is one rung of the data-state ladder.
type DataStateState string

// Data state constants.
const (
	// DataStateContent renders the children — data loaded, everything fine.
	// The zero value.
	DataStateContent DataStateState = "content"
	// DataStateEmpty renders the empty placeholder — loaded, nothing there.
	DataStateEmpty DataStateState = "empty"
	// DataStateUnavailable renders the failure placeholder — the query
	// failed or is retrying.
	DataStateUnavailable DataStateState = "unavailable"
	// DataStateDisabled renders the feature-off placeholder — the data
	// source is switched off by configuration.
	DataStateDisabled DataStateState = "disabled"
)

// DataStateStateIsValid reports whether s is one of the defined states.
// Unknown values fall back to content (graceful — never a panic).
func DataStateStateIsValid(s DataStateState) bool {
	return s == DataStateContent || s == DataStateEmpty || s == DataStateUnavailable || s == DataStateDisabled
}

// DataStateProps configures the data-state ladder for one region.
type DataStateProps struct {
	utils.BaseProps

	// State selects the rung. DataStateContent (the zero value) renders
	// Children; every other state renders the placeholder.
	State DataStateState
	// Title and Description carry the chosen state's copy — per-state truth,
	// so they sit on the props, not in the component.
	Title       string
	Description string
	// TitleTag overrides the placeholder's heading element (defaults to h3,
	// like EmptyState).
	TitleTag HeadingTagType
	// Icon overrides the state's default icon (empty: inbox, unavailable:
	// exclamation-triangle, disabled: no-symbol).
	Icon icons.Name
}

// DefaultDataStateProps returns sensible defaults.
func DefaultDataStateProps() DataStateProps {
	return DataStateProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		State: DataStateContent,
	}
}
