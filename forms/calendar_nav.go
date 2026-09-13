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
func calendarMonthNavAction(base *wire.Action, year int, month time.Month) templ.Attributes {
	if base == nil || base.URL == "" {
		return nil
	}

	nav := *base
	nav.URL = strings.ReplaceAll(nav.URL, "{year}", strconv.Itoa(year))
	nav.URL = strings.ReplaceAll(nav.URL, "{month}", strconv.Itoa(int(month)))

	return nav.Attributes()
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
