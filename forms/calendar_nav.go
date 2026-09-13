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
// htmx additionally gets hx-swap="outerHTML": the endpoint re-renders the
// WHOLE calendar, and outer-HTML self-replacement is the library's proven
// re-binding pattern (kanban, LoadMore) — an innerHTML swap of a region that
// CONTAINS the triggers leaves the swapped-in anchors unprocessed by htmx
// (verified 2026-09-13 in a real browser). Point Target at the calendar's
// own id.
func calendarMonthNavAction(base *wire.Action, year int, month time.Month) templ.Attributes {
	if base == nil || base.URL == "" {
		return nil
	}

	nav := *base
	nav.URL = strings.ReplaceAll(nav.URL, "{year}", strconv.Itoa(year))
	nav.URL = strings.ReplaceAll(nav.URL, "{month}", strconv.Itoa(int(month)))

	attrs := nav.Attributes()
	if attrs != nil && nav.Transport != wire.TransportDatastar {
		attrs["hx-swap"] = "outerHTML"
	}

	return attrs
}

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
