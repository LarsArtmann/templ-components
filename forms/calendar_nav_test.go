package forms

import (
	"testing"
	"time"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestCalendarMonthNavWire pins the dual-transport month navigation (#157):
// prev/next anchors carry the MonthNav action's dialect attributes with the
// {year}/{month} placeholders substituted per direction, the consumer's
// action is never mutated, and December/January wrap the year.
func TestCalendarMonthNavWire(t *testing.T) {
	t.Parallel()

	props := DefaultCalendarProps()
	props.Year = 2026
	props.Month = time.July
	props.BaseProps = utils.BaseProps{ID: "cal"}

	t.Run("htmx dialect renders hx-get with substituted month", func(t *testing.T) {
		t.Parallel()

		action := &wire.Action{ //nolint:exhaustruct // URL+Target are the wiring surface
			URL:    "/api/calendar?year={year}&month={month}",
			Target: "#cal",
		}
		props.MonthNav = action

		html := utils.Render(t, Calendar(props))
		utils.AssertContainsAll(t, html,
			`hx-get="/api/calendar?year=2026&amp;month=6"`,
			`hx-get="/api/calendar?year=2026&amp;month=8"`,
			`hx-target="#cal"`,
			`hx-swap="outerHTML settle:0s"`,
			`aria-label="Previous month"`,
			`aria-label="Next month"`,
		)
		utils.AssertNotContains(t, html, "{year}")
		utils.AssertNotContains(t, html, "{month}")
	})

	t.Run("datastar dialect renders @get binding", func(t *testing.T) {
		t.Parallel()

		props.MonthNav = &wire.Action{ //nolint:exhaustruct // URL is the wiring surface
			Transport: wire.TransportDatastar,
			URL:       "/api/calendar?year={year}&month={month}",
		}

		html := utils.Render(t, Calendar(props))
		utils.AssertContainsAll(t, html,
			`data-on:click="@get(&#39;/api/calendar?year=2026&amp;month=6&#39;)"`,
			`data-on:click="@get(&#39;/api/calendar?year=2026&amp;month=8&#39;)"`,
		)
	})

	t.Run("december and january wrap the year", func(t *testing.T) {
		t.Parallel()

		december := props
		december.Month = time.December
		december.MonthNav = &wire.Action{URL: "/c?y={year}&m={month}"} //nolint:exhaustruct // URL is the wiring surface

	html := utils.Render(t, Calendar(december))
	utils.AssertContains(t, html, `/c?y=2027&amp;m=1`)

	january := props
	january.Month = time.January
	january.MonthNav = &wire.Action{URL: "/c?y={year}&m={month}"} //nolint:exhaustruct // URL is the wiring surface

	html = utils.Render(t, Calendar(january))
	utils.AssertContains(t, html, `/c?y=2025&amp;m=12`)
	})

	t.Run("consumer action is never mutated", func(t *testing.T) {
		t.Parallel()

		const template = "/api/calendar?year={year}&month={month}"
		action := &wire.Action{URL: template} //nolint:exhaustruct // URL is the wiring surface
		props.MonthNav = action

		_ = utils.Render(t, Calendar(props))

		if action.URL != template {
			t.Errorf("MonthNav action URL mutated: %q, want %q — components copy consumer specs, never rewrite them", action.URL, template)
		}
	})

	t.Run("href fallback renders alongside wire attributes", func(t *testing.T) {
		t.Parallel()

		props.HrefPrev = "/calendar?y=2026&m=6"
		props.HrefNext = "/calendar?y=2026&m=8"
		props.MonthNav = &wire.Action{URL: "/api/calendar?year={year}&month={month}"} //nolint:exhaustruct // URL is the wiring surface

		html := utils.Render(t, Calendar(props))
		utils.AssertContainsAll(t, html, `href="/calendar?y=2026&amp;m=6"`, `hx-get=`)
	})

	t.Run("nil MonthNav keeps href-only behavior", func(t *testing.T) {
		t.Parallel()

		props.MonthNav = nil
		props.HrefPrev = "/calendar?y=2026&m=6"

		html := utils.Render(t, Calendar(props))
		utils.AssertContains(t, html, `href="/calendar?y=2026&amp;m=6"`)
		utils.AssertNotContains(t, html, "hx-get")
	})
}

// TestGoldenSweepCalendarMonthNav pins the wired calendar markup.
func TestGoldenSweepCalendarMonthNav(t *testing.T) {
	t.Parallel()

	props := DefaultCalendarProps()
	props.Year = 2026
	props.Month = time.July
	props.BaseProps = utils.BaseProps{ID: "cal"}
	props.MonthNav = &wire.Action{ //nolint:exhaustruct // URL+Target are the wiring surface
		URL:    "/api/calendar?year={year}&month={month}",
		Target: "#cal",
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "calendar_month_nav_htmx", HTML: utils.Render(t, Calendar(props))},
	})

	props.MonthNav = &wire.Action{ //nolint:exhaustruct // URL is the wiring surface
		Transport: wire.TransportDatastar,
		URL:       "/api/calendar?year={year}&month={month}",
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "calendar_month_nav_datastar", HTML: utils.Render(t, Calendar(props))},
	})
}
