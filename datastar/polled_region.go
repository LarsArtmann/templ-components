package datastar

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// defaultPolledInterval is the default poll interval for PolledRegion
// (mirrors htmx.PolledRegion's defaultPollInterval).
const defaultPolledInterval = "10s"

// PolledRegionProps configures an auto-refreshing Datastar region. The region
// fetches its own URL on an interval (data-on-interval) and the endpoint
// responds with a fragment whose root id matches the region id, so the
// runtime's default outer-mode patch replaces the region — the same
// self-replacing model as htmx.PolledRegion.
type PolledRegionProps struct {
	utils.BaseProps

	// URL is the endpoint fetched on every tick. It should return the
	// refreshed region fragment (root element id = this region's id).
	// An empty URL renders NO interval attribute — the region stays a plain
	// static container (the runtime throws FetchNoUrlProvided for @get('')).
	URL string

	// Every is the poll interval. Accepted shapes: "500ms", "10s", "2m",
	// "1h" (m/h are converted to seconds — the runtime's duration parser
	// would silently read them as MILLISECONDS, a 1ms self-poll storm from a
	// natural-looking typo), or a bare number (milliseconds, matching the
	// runtime). An empty Every defaults to "10s". Anything else passes
	// through verbatim (exotic timing stays the runtime's business).
	Every string

	// Live sets the aria-live politeness for screen-reader announcements.
	// Defaults to LivePolite.
	Live LivePoliteness

	// ShowTimestamp renders an "Updated HH:MM:SS" footer so operators can
	// see at a glance that polling is active (the timestamp ticks forward on
	// each successful poll because the whole region re-renders).
	ShowTimestamp bool

	// Now overrides the wall clock used for the timestamp footer. Nil (the
	// zero default) means time.Now. Inject a fixed clock in visual goldens
	// and tests so ShowTimestamp coverage is deterministic.
	Now func() time.Time

	// TimeFormat is the Go time format string for the timestamp footer.
	// Default: "15:04:05" (time-only). Use time.RFC3339 for full date+time.
	TimeFormat string
}

// DefaultPolledRegionProps returns sensible defaults for a polled region.
func DefaultPolledRegionProps() PolledRegionProps {
	return PolledRegionProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Every:         defaultPolledInterval,
		Live:          LivePolite,
		ShowTimestamp: true,
	}
}

// polledIntervalValue normalizes the Every prop into the runtime duration
// token emitted inside the data-on-interval__duration modifier.
//
// Why m/h are converted: the pinned runtime parses durations as
// `<n>ms` / `<n>s` / bare parseFloat-ms — a trailing "m" or "h" truncates at
// parseFloat ("5m" → 5ms), so a natural-looking "5m" typo would poll every
// 5 milliseconds. Hours become seconds (same rule as the htmx twin, whose
// parseInterval has the identical hazard) and minutes become seconds too,
// which the runtime's parser cannot express natively.
func polledIntervalValue(every string) string {
	if strings.TrimSpace(every) == "" {
		return defaultPolledInterval
	}

	for _, conv := range []struct {
		re      *regexp.Regexp
		seconds float64
	}{
		{polledHoursRe, 3600},
		{polledMinutesRe, 60},
	} {
		match := conv.re.FindStringSubmatch(every)
		if match == nil {
			continue
		}

		n, err := strconv.ParseFloat(match[1], 64)
		if err != nil || n <= 0 {
			return every
		}

		return strconv.FormatFloat(n*conv.seconds, 'f', -1, 64) + "s"
	}

	return every
}

var (
	polledHoursRe   = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)h$`)
	polledMinutesRe = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)m$`)
)

// polledRegionAttrs builds the interval attribute for a polled region:
//
//	data-on-interval__duration.10s="@get('/url')"
//
// The duration lives in the ATTRIBUTE NAME (the runtime's modifier grammar:
// the name splits on __ into modifier groups), so the attribute is emitted
// through a templ.Attributes map with a dynamic key. An empty URL returns
// nil — the region renders inert (no @get(”) expressions, which the runtime
// rejects).
func polledRegionAttrs(url, every string) templ.Attributes {
	if strings.TrimSpace(url) == "" {
		return nil
	}

	return templ.Attributes{
		"data-on-interval__duration." + polledIntervalValue(every): actionExpr("get", url),
	}
}
