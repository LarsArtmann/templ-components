package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- LoadingButton Behavior ---

func TestLoadingButtonUserSeesBusySwap(t *testing.T) {
	t.Parallel()

	t.Run("user wires one signal name across button and label", func(t *testing.T) {
		t.Parallel()

		// The contract: the button carries data-indicator:<signal> (runtime
		// sets it during the action) and the label reacts via data-show on
		// the same signal. The component pins both halves of the expression.
		output := utils.Render(t, LoadingButton(LoadingButtonProps{
			Signal:      "saving",
			DefaultText: "Save",
			LoadingText: "Saving…",
		}))
		utils.AssertContainsAll(t, output, `data-show="$saving"`, `data-show="!$saving"`)
	})

	t.Run("a mistyped signal cannot pin the busy state", func(t *testing.T) {
		t.Parallel()

		// Empty/whitespace signal degrades to false/true — the runtime's
		// bare "$" object is always truthy and would show the spinner
		// forever (the Indicator degradation rule).
		output := utils.Render(t, LoadingButton(LoadingButtonProps{DefaultText: "Save"}))
		utils.AssertContains(t, output, `data-show="false"`)
	})
}

func TestLoadingButtonUserGetsRuntimeDrivenSwap(t *testing.T) {
	t.Parallel()

	// Zero component-side JavaScript: the swap is data-show + the runtime's
	// indicator signal, so HTMX re-renders and CSP nonce requirements do
	// not apply to this component at all.
	output := utils.Render(t, LoadingButton(LoadingButtonProps{Signal: "s", DefaultText: "Go"}))
	utils.AssertNotContains(t, output, "<script")
}
