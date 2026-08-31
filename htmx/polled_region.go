package htmx

import "github.com/larsartmann/templ-components/utils"

// defaultPollInterval is the default HTMX polling interval for PolledRegion.
const defaultPollInterval = "10s"

// PolledLive is the aria-live politeness setting for a polled region.
type PolledLive string

const (
	// PolledLivePolite announces updates after the current speech finishes (default).
	PolledLivePolite PolledLive = "polite"
	// PolledLiveAssertive announces updates immediately, interrupting current speech.
	PolledLiveAssertive PolledLive = "assertive"
	// PolledLiveOff disables screen-reader announcements for this region.
	PolledLiveOff PolledLive = "off"
)

// PolledRegionProps configures an auto-refreshing HTMX region.
type PolledRegionProps struct {
	utils.BaseProps

	// URL is the HTMX endpoint to fetch for refreshing the region content.
	URL string

	// Every is the poll interval as an HTMX timing string (e.g. "10s", "2m", "500ms").
	Every string

	// Eager fires the first fetch on initial load in addition to polling.
	// When false, the first fetch happens after the first interval elapses.
	// Ignored when Trigger is set.
	Eager bool

	// Trigger overrides the auto-generated hx-trigger value. When set,
	// Every and Eager are ignored and this string is used verbatim as
	// hx-trigger. Useful for custom triggers like SSE events
	// (e.g. "stats-refresh from:body").
	Trigger string

	// Swap controls how the fetched content replaces the region.
	// Defaults to SwapOuterHTML.
	Swap SwapStyle

	// Live sets the aria-live politeness for screen-reader announcements.
	// Defaults to PolledLivePolite.
	Live PolledLive

	// ShowTimestamp renders a "Updated HH:MM:SS" footer so operators can see
	// at a glance that polling is active (the timestamp ticks forward on each
	// successful poll because the whole region re-renders).
	ShowTimestamp bool

	// TimeFormat is the Go time format string for the timestamp footer.
	// Default: "15:04:05" (time-only). Use time.RFC3339 for full date+time.
	TimeFormat string
}

// DefaultPolledRegionProps returns sensible defaults for a polled region.
func DefaultPolledRegionProps() PolledRegionProps {
	return PolledRegionProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Every:         defaultPollInterval,
		Swap:          SwapOuterHTML,
		Live:          PolledLivePolite,
		ShowTimestamp: true,
	}
}

// PolledLiveIsValid reports whether v is one of the defined PolledLive constants.
func PolledLiveIsValid(v PolledLive) bool {
	return v == PolledLivePolite || v == PolledLiveAssertive || v == PolledLiveOff
}
