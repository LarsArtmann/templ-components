package datastar

import "github.com/larsartmann/templ-components/utils"

// LivePoliteness is the aria-live politeness setting for a live region.
// Controls how screen readers announce streamed content updates.
type LivePoliteness string

const (
	// LivePolite announces updates after the current speech finishes (default).
	LivePolite LivePoliteness = "polite"
	// LiveAssertive announces updates immediately, interrupting current speech.
	// Use sparingly — only for critical alerts.
	LiveAssertive LivePoliteness = "assertive"
	// LiveOff disables screen-reader announcements for this region.
	LiveOff LivePoliteness = "off"
)

// validLivePoliteness are the allowed aria-live values for LiveRegion.
//
//nolint:gochecknoglobals // Package-level validation set
var validLivePoliteness = map[LivePoliteness]bool{
	LivePolite: true, LiveAssertive: true, LiveOff: true,
}

// livePolitenessValue returns p if valid, otherwise falls back to polite.
func livePolitenessValue(p LivePoliteness) LivePoliteness {
	if validLivePoliteness[p] {
		return p
	}

	return LivePolite
}

// LiveRegionProps configures an SSE-powered auto-updating region.
// This is the Datastar equivalent of htmx.PolledRegion — but instead of
// polling on an interval, the server pushes patches via a long-lived SSE
// connection.
//
// The URL endpoint should return a text/event-stream response. The server
// patches the region's children (or signals) using go-datastar
// (github.com/larsartmann/go-datastar). Patch a child by selector — see the
// example on LiveRegion in live_region.templ for the full handler.
//
//	@datastar.LiveRegion(datastar.LiveRegionProps{
//	    URL:       "/stream/metrics",
//	    AutoStart: true,
//	}) {
//	    @display.StatCard(display.StatCardProps{
//	        BaseProps: utils.BaseProps{ID: "metrics"},
//	        Label:     "Users",
//	        Value:     "—",
//	    })
//	}
type LiveRegionProps struct {
	utils.BaseProps

	// URL is the SSE endpoint that streams element/signal patches.
	// The server opens a long-lived connection and pushes updates.
	URL string

	// AutoStart opens the SSE connection on page load via data-init.
	// When false, the consumer triggers the stream manually (e.g. via a
	// button with data-on:click="@get('/stream')").
	AutoStart bool

	// Live sets the aria-live politeness for screen-reader announcements.
	// Defaults to LivePolite.
	Live LivePoliteness
}

// DefaultLiveRegionProps returns sensible defaults for a live region.
func DefaultLiveRegionProps() LiveRegionProps {
	return LiveRegionProps{ //nolint:exhaustruct // intentionally minimal defaults
		AutoStart: true,
		Live:      LivePolite,
	}
}
