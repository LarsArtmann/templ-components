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
// The clone sets PreventDefault: NavLink-style anchors whose href is the
// no-JS fallback must not ALSO navigate when the Datastar runtime handles
// the click (the runtime auto-preventDefaults only form+submit — decoded
// from the pinned v1.0.3 bundle); htmx intercepts wired clicks itself.
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
	nav.URL = calendarMonthNavURL(base, year, month)
	nav.PreventDefault = true

	attrs := nav.Attributes()
	if attrs != nil && nav.Transport != wire.TransportDatastar {
		attrs["hx-swap"] = calendarMonthNavSwap
	}

	return attrs
}

// calendarMonthNavURL substitutes {year}/{month} placeholders in the action
// URL with the navigation target month. Empty when the action is nil or
// carries no URL (an unwired action must stay inert).
func calendarMonthNavURL(base *wire.Action, year int, month time.Month) string {
	if base == nil || base.URL == "" {
		return ""
	}

	url := strings.ReplaceAll(base.URL, "{year}", strconv.Itoa(year))

	return strings.ReplaceAll(url, "{month}", strconv.Itoa(int(month)))
}

// calendarMonthNavHref resolves the arrow's href: the explicit
// HrefPrev/HrefNext when set, otherwise the substituted MonthNav URL —
// progressive enhancement, so without JS the arrow navigates to the same
// endpoint the runtime patches. Never returns a roleless aria-labelled
// anchor: an <a> without href carries no implicit role, and aria-label on
// it is an axe aria-prohibited-attr serious violation (found by the demo
// axe sweep 2026-09-14). Callers render the arrow only when this returns
// a non-empty href.
func calendarMonthNavHref(explicit string, base *wire.Action, year int, month time.Month) string {
	if explicit != "" {
		return explicit
	}

	return calendarMonthNavURL(base, year, month)
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
