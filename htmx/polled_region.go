package htmx

import "github.com/larsartmann/templ-components/utils"

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
	Eager bool

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
}

// DefaultPolledRegionProps returns sensible defaults for a polled region.
func DefaultPolledRegionProps() PolledRegionProps {
	return PolledRegionProps{ //nolint:exhaustruct // intentionally minimal defaults
		Every:         "10s",
		Swap:          SwapOuterHTML,
		Live:          PolledLivePolite,
		ShowTimestamp: true,
	}
}

// PolledLiveIsValid reports whether v is one of the defined PolledLive constants.
func PolledLiveIsValid(v PolledLive) bool {
	return v == PolledLivePolite || v == PolledLiveAssertive || v == PolledLiveOff
}
