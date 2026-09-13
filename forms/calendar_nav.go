package forms

import (
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils/wire"
)

// calendarMonthNavAction clones the Calendar MonthNav action for one
// direction, substituting {year}/{month} placeholders with the navigation
// target month. It never mutates the consumer's action (copy semantics —
// the kanban lesson: components must not rewrite consumer-supplied specs).
//
// htmx additionally gets hx-swap="outerHTML settle:0s": the endpoint
// re-renders the WHOLE calendar, and outer-HTML self-replacement is the
// library's proven re-binding pattern (kanban, LoadMore). settle:0s makes
// htmx process the swapped-in anchors in the same JS task as the swap —
// with the default 20ms settle delay the new arrows stay inert for that
// window, and month arrows are exactly the kind of control users click in
// quick bursts (e2e-verified 2026-09-13: a click inside the window silently
// no-ops). Point Target at the calendar's own id.
func calendarMonthNavAction(base *wire.Action, year int, month time.Month) templ.Attributes {
	if base == nil || base.URL == "" {
		return nil
	}

	nav := *base
	nav.URL = strings.ReplaceAll(nav.URL, "{year}", strconv.Itoa(year))
	nav.URL = strings.ReplaceAll(nav.URL, "{month}", strconv.Itoa(int(month)))

	attrs := nav.Attributes()
	if attrs != nil && nav.Transport != wire.TransportDatastar {
		attrs["hx-swap"] = calendarMonthNavSwap
	}

	return attrs
}

// calendarMonthNavSwap self-replaces the calendar and settles synchronously.
const calendarMonthNavSwap = "outerHTML settle:0s"

// calendarPrevMonth returns the year/month before the given one
// (December wraps to the previous year).
func calendarPrevMonth(year int, month time.Month) (int, time.Month) {
	if month == time.January {
		return year - 1, time.December
	}

	return year, month - 1
}

// calendarNextMonth returns the year/month after the given one
// (January wraps to the next year).
func calendarNextMonth(year int, month time.Month) (int, time.Month) {
	if month == time.December {
		return year + 1, time.January
	}

	return year, month + 1
}
