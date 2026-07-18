package forms

import (
	"testing"
	"time"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultCalendarProps(t *testing.T) {
	t.Parallel()

	p := DefaultCalendarProps()

	now := time.Now()
	if p.Year != now.Year() || p.Month != now.Month() {
		t.Error("expected current year/month defaults")
	}
}

func TestCalendarBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Calendar(CalendarProps{
		Year:  2026,
		Month: time.July,
	}))
	utils.AssertContains(t, output, "July")
	utils.AssertContains(t, output, "2026")
	utils.AssertContains(t, output, `role="grid"`)
}

func TestCalendarGolden(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Calendar(CalendarProps{
		Year:  2026,
		Month: time.January,
	}))
	golden.Assert(t, "calendar_basic", output)
}

func TestCalendarDayLinks(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Calendar(CalendarProps{
		Year:     2026,
		Month:    time.July,
		HrefBase: "/cal?y={year}&m={month}&d={day}",
	}))
	utils.AssertContains(t, output, `/cal?y=2026&amp;m=7&amp;d=1`)
	utils.AssertContains(t, output, `/cal?y=2026&amp;m=7&amp;d=15`)
}

func TestCalendarPrevNextLinks(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Calendar(CalendarProps{
		Year:     2026,
		Month:    time.July,
		HrefPrev: "/cal?y=2026&m=6",
		HrefNext: "/cal?y=2026&m=8",
	}))
	utils.AssertContains(t, output, `href="/cal?y=2026&amp;m=6"`)
	utils.AssertContains(t, output, `href="/cal?y=2026&amp;m=8"`)
	utils.AssertContains(t, output, `aria-label="Previous month"`)
	utils.AssertContains(t, output, `aria-label="Next month"`)
}

func TestCalendarSelectedDate(t *testing.T) {
	t.Parallel()

	selected := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	output := utils.Render(t, Calendar(CalendarProps{
		Year:         2026,
		Month:        time.July,
		SelectedDate: &selected,
	}))
	utils.AssertContains(t, output, "bg-blue-600")
	utils.AssertContains(t, output, "text-white")
}

func TestCalendarMinMaxDates(t *testing.T) {
	t.Parallel()

	minDate := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	output := utils.Render(t, Calendar(CalendarProps{
		Year:    2026,
		Month:   time.July,
		MinDate: &minDate,
		MaxDate: &maxDate,
	}))
	utils.AssertContains(t, output, "cursor-not-allowed")
}

func TestCalendarEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("no href base renders spans not links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Calendar(CalendarProps{
			Year:  2026,
			Month: time.February,
		}))
		utils.AssertNotContains(t, output, `<a href`)
	})

	t.Run("custom class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Calendar(CalendarProps{
			BaseProps: utils.BaseProps{Class: "my-cal"},
			Year:      2026,
			Month:     time.July,
		}))
		utils.AssertContains(t, output, "my-cal")
	})

	t.Run("february leap year renders 29 days", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Calendar(CalendarProps{
			Year:  2024,
			Month: time.February,
		}))
		utils.AssertContains(t, output, ">29<")
	})
}

func TestCalendarDayClassHelpers(t *testing.T) {
	t.Parallel()

	if got := calendarDayClass(false, false, false); got == "" {
		t.Error("expected non-empty class")
	}

	if got := calendarDayClass(true, false, false); got == "" {
		t.Error("expected non-empty selected class")
	}

	if got := calendarDayClass(false, true, false); got == "" {
		t.Error("expected non-empty disabled class")
	}

	if got := calendarDayClass(false, false, true); got == "" {
		t.Error("expected non-empty today class")
	}
}

func TestCalendarReplacePlaceholders(t *testing.T) {
	t.Parallel()

	got := calendarReplacePlaceholders("/cal?y={year}&m={month}&d={day}", 2026, 7, 15)
	if got != "/cal?y=2026&m=7&d=15" {
		t.Errorf("got %q", got)
	}
}

func ExampleCalendar() {
	selected := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	_ = Calendar(CalendarProps{
		Year:         2026,
		Month:        time.July,
		SelectedDate: &selected,
		HrefBase:     "/calendar?y={year}&m={month}&d={day}",
	})
	// Output:
}
