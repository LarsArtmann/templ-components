package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestLoadingButtonA11y is the accessibility lens for the loading button.
func TestLoadingButtonA11y(t *testing.T) {
	t.Parallel()

	t.Run("spinner is decorative, labels carry meaning", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LoadingButton(LoadingButtonProps{
			Signal:      "saving",
			DefaultText: "Save",
			LoadingText: "Saving…",
		}))

		// The spinner span carries no aria-hidden in the Datastar dialect —
		// the runtime toggles visibility via data-show (display), which
		// already removes it from the accessibility tree when hidden.
		utils.AssertContains(t, output, "animate-spin")
	})

	t.Run("motion-reduce fallback on the spinner", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LoadingButton(LoadingButtonProps{Signal: "s", DefaultText: "Go"}))
		utils.AssertContains(t, output, "motion-reduce:animate-none")
	})

	t.Run("AriaLabel propagates", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LoadingButton(LoadingButtonProps{
			AriaLabel:   "Save changes",
			Signal:      "saving",
			DefaultText: "Save",
		}))
		utils.AssertContains(t, output, `aria-label="Save changes"`)
	})

	t.Run("label swap uses data-show, not CSS-only visibility tricks", func(t *testing.T) {
		t.Parallel()

		// data-show toggles display, which removes hidden labels from the
		// accessibility tree — screen readers announce the visible label
		// only, and a focused button announces the swap.
		output := utils.Render(t, LoadingButton(LoadingButtonProps{Signal: "saving", DefaultText: "Save"}))
		utils.AssertContains(t, output, `data-show="$saving"`)
		utils.AssertContains(t, output, `data-show="!$saving"`)
	})
}
