package datastar

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestSignalChipA11y is the accessibility lens for the signal chip.
func TestSignalChipA11y(t *testing.T) {
	t.Parallel()

	render := func() string {
		return utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			AriaLabel: "Filter by state",
			Signal:    "stateFilter",
			Selected:  "open",
			Items: []SignalChipItem{
				{Label: "All"},
				{Label: "Open", Value: "open"},
			},
		}))
	}

	t.Run("pressed state is exposed on every chip, in both states", func(t *testing.T) {
		t.Parallel()

		output := render()

		utils.AssertContains(t, output, `aria-pressed="true"`)

		utils.AssertContains(t, output, `aria-pressed="false"`)
	})

	t.Run("the reactive pressed state stays a valid token", func(t *testing.T) {
		t.Parallel()

		// A boolean data-attr expression makes the runtime REMOVE the
		// attribute when false and set an empty value when true, which leaves
		// assistive tech without a state. The string form always yields one
		// of the four aria-pressed tokens.
		output := render()

		utils.AssertContains(t, output, `data-attr:aria-pressed="$stateFilter===&#39;open&#39; ? &#39;true&#39; : &#39;false&#39;"`)
	})

	t.Run("the chip row is a named group", func(t *testing.T) {
		t.Parallel()

		output := render()

		utils.AssertContains(t, output, `role="group"`)

		utils.AssertContains(t, output, `aria-label="Filter by state"`)
	})

	t.Run("an unnamed row still gets an accessible name", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{Signal: "stateFilter"}))

		if !strings.Contains(output, "aria-label=") {
			t.Error("unnamed chip row renders no aria-label")
		}
	})

	t.Run("chips are native buttons that never submit a form", func(t *testing.T) {
		t.Parallel()

		// A native button is keyboard and switch operable for free. Without
		// type="button" a chip inside a form would submit it.
		output := render()

		if got := strings.Count(output, `<button type="button"`); got != 2 {
			t.Errorf(`<button type="button" rendered %d times, want 2`, got)
		}
	})

	t.Run("AriaLabel on a chip overrides its accessible name", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			BaseProps: utils.BaseProps{AriaLabel: "Show open findings"},
			Signal:    "stateFilter",
			Value:     "open",
			Label:     "Open",
		}))

		utils.AssertContains(t, output, `aria-label="Show open findings"`)
	})

	t.Run("the count is text, announced with the label", func(t *testing.T) {
		t.Parallel()

		count := 12
		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: "Open", Count: &count,
		}))

		utils.AssertContainsAll(t, output, `>Open <span`, `>12</span></button>`)

		utils.AssertContainsAll(t, output, "ms-1.5", "tabular-nums")

		// A purely decorative badge would be hidden; the count is information.
		utils.AssertNotContains(t, output, `aria-hidden="true"`)
	})

	t.Run("focus ring and reduced motion", func(t *testing.T) {
		t.Parallel()

		output := render()

		utils.AssertContainsAll(t, output, "focus-visible:outline-2", "motion-reduce:transition-none", "motion-reduce:duration-0")
	})

	t.Run("logical properties only (RTL)", func(t *testing.T) {
		t.Parallel()

		count := 3

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: "Open", Count: &count,
		}))

		utils.AssertContains(t, output, "ms-1.5")

		utils.AssertNotContains(t, output, "ml-")
	})

	t.Run("both looks carry dark variants", func(t *testing.T) {
		t.Parallel()

		output := render()

		utils.AssertContainsAll(t, output, "dark:bg-accent-600", "dark:bg-gray-800", "dark:text-gray-300")
	})

	t.Run("the pressed pairing keeps white text at WCAG AA contrast", func(t *testing.T) {
		t.Parallel()

		// visualtest.TestDatastarSignalChipContrast measures the real ratios
		// in a browser, at rest and on hover, in both themes. This pins the
		// class pairing that measurement rests on: white sits on the accent
		// 600 shade at rest and the 700 shade on hover, in BOTH themes.
		utils.AssertContainsAll(t, signalChipSelectedClass,
			"text-white", "dark:text-white",
			"bg-accent-600", "hover:bg-accent-700", "dark:bg-accent-600", "dark:hover:bg-accent-700")

		// A white label on the 500 or 400 shade measures 3.8:1 or worse.
		utils.AssertNotContains(t, signalChipSelectedClass, "bg-accent-500")
		utils.AssertNotContains(t, signalChipSelectedClass, "bg-accent-400")
	})

	t.Run("the count is set apart by weight, never by opacity", func(t *testing.T) {
		t.Parallel()

		// A count at 70% opacity measured 3.3:1 on the pressed chip.
		utils.AssertContains(t, signalChipCountClass, "font-normal")

		utils.AssertNotContains(t, signalChipCountClass, "opacity")
	})
}
