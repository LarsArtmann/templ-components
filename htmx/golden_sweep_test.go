package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils"
)

// TestGoldenSweepHTMXComponents provides golden HTML snapshots for every htmx
// component that previously relied on brittle substring assertions. PolledRegion
// already has its own sweep (golden_sweep_polled_test.go).
func TestGoldenSweepHTMXComponents(t *testing.T) {
	t.Parallel()

	spinner := feedback.Spinner(feedback.DefaultSpinnerProps())

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "loading_indicator", HTML: utils.Render(t, LoadingIndicator(spinner))},
		{Name: "inline_loading_overlay", HTML: utils.Render(t, InlineLoadingOverlay("saving-overlay", spinner))},
		{Name: "loading_button", HTML: utils.Render(t, LoadingButton("Save", "Saving…", spinner))},
		{Name: "csrf_token", HTML: utils.Render(t, CSRFToken("abc123secret"))},
		{Name: "confirm_delete", HTML: utils.Render(t, ConfirmDelete(ConfirmDeleteProps{
			Delete:  "/api/items/42",
			Target:  "#item-42",
			Confirm: "Delete this item?",
		}))},
		{Name: "confirm_delete_no_confirm", HTML: utils.Render(t, ConfirmDelete(ConfirmDeleteProps{
			Delete: "/api/items/42",
			Target: "#item-42",
		}))},
		{Name: "swap_oob_default", HTML: utils.Render(t, SwapOOB(SwapOOBProps{
			Selector: "#toast-container",
		}))},
		{Name: "swap_oob_beforeend", HTML: utils.Render(t, SwapOOB(SwapOOBProps{
			Selector:  "#toast-container",
			SwapStyle: SwapBeforeEnd,
		}))},
		{Name: "swap_oob_no_selector", HTML: utils.Render(t, SwapOOB(SwapOOBProps{}))},
		{Name: "global_error_handling", HTML: utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))},
		{Name: "global_error_handling_custom", HTML: utils.Render(t, GlobalErrorHandling(ErrorHandlingConfig{
			MaxErrorHistory: 10,
			MaxRetries:      5,
			RetryDelayMS:    2000,
		}))},
		{Name: "view_transitions_global", HTML: utils.Render(t, ViewTransitions(ViewTransitionsProps{Global: true}))},
		{Name: "view_transitions_off", HTML: utils.Render(t, ViewTransitions(DefaultViewTransitionsProps()))},
	})
}
