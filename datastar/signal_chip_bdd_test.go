package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- SignalChip Behavior ---

func TestSignalChipUserSeesActiveFilter(t *testing.T) {
	t.Parallel()

	t.Run("the pressed chip is marked in the server-rendered HTML", func(t *testing.T) {
		t.Parallel()

		// Before JavaScript, and without it entirely, class and aria-pressed
		// already tell which filter is active.
		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: "Open", Selected: true,
		}))

		utils.AssertContainsAll(t, output, `aria-pressed="true"`, "bg-accent-600")
	})

	t.Run("an unpressed chip is marked unpressed", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: "Open",
		}))

		utils.AssertContains(t, output, `aria-pressed="false"`)

		utils.AssertContains(t, output, "bg-white")
	})
}

func TestSignalChipUserClicksToFilter(t *testing.T) {
	t.Parallel()

	t.Run("clicking a chip assigns its value to the signal", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "done", Label: "Done",
		}))

		utils.AssertContains(t, output, `data-on:click="$stateFilter=&#39;done&#39;"`)
	})

	t.Run("every chip in a group follows the same signal", func(t *testing.T) {
		t.Parallel()

		// One signal, many chips: each click rewrites it and each chip's look
		// is keyed on the same comparison, so they stay in sync with no
		// server round-trip and no component-side JavaScript.
		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			Signal: "stateFilter",
			Items:  []SignalChipItem{{Label: "Open", Value: "open"}, {Label: "Done", Value: "done"}},
		}))

		utils.AssertContainsAll(t, output,
			`data-on:click="$stateFilter=&#39;open&#39;"`,
			`data-on:click="$stateFilter=&#39;done&#39;"`,
			`$stateFilter===&#39;open&#39;`,
			`$stateFilter===&#39;done&#39;`,
		)
	})

	t.Run("the pressed look is swapped by the runtime, not by a re-render", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{Signal: "stateFilter", Value: "open", Label: "Open"}))

		// Both looks travel inside data-class so the runtime can swap them.
		classExpr := signalChipAttrExpr(output, "data-class")

		utils.AssertContainsAll(t, classExpr, signalChipSelectedClass, signalChipUnselectedClass)

		utils.AssertNotContains(t, output, "<script")
	})
}

func TestSignalChipUserKeepsActiveFilterOnLoad(t *testing.T) {
	t.Parallel()

	t.Run("the pressed chip starts the runtime in the rendered state", func(t *testing.T) {
		t.Parallel()

		// An undeclared signal reads as empty, so without the initial value
		// the runtime would strip the pressed look from the selected chip on
		// load. The selected chip declares it.
		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			Signal:   "stateFilter",
			Selected: "open",
			Items:    []SignalChipItem{{Label: "Open", Value: "open"}},
		}))

		utils.AssertContains(t, output, `data-signals="{&#39;stateFilter&#39;:&#39;open&#39;}"`)
	})
}

func TestSignalChipUserCanClearTheFilter(t *testing.T) {
	t.Parallel()

	t.Run("an empty-value chip resets the signal to empty", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{Signal: "stateFilter", Label: "All"}))

		utils.AssertContains(t, output, `data-on:click="$stateFilter=&#39;&#39;"`)

		utils.AssertContains(
			t,
			output,
			`data-attr:aria-pressed="$stateFilter===&#39;&#39; ? &#39;true&#39; : &#39;false&#39;"`,
		)
	})
}

func TestSignalChipUserWithBrokenConfigGetsAHarmlessChip(t *testing.T) {
	t.Parallel()

	t.Run("an unbindable signal name renders a static chip", func(t *testing.T) {
		t.Parallel()

		// Nothing is emitted for the runtime to misparse, and the chip still
		// reads as a (pressable-looking) button with its server state.
		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "state filter", Value: "open", Label: "Open", Selected: true,
		}))

		utils.AssertNotContains(t, output, "data-on")

		utils.AssertContains(t, output, `aria-pressed="true"`)
	})
}
