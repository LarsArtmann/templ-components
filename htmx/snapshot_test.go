// Package htmx provides tests for HTMX helper components.
package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/utils"
)

func TestLoadingIndicatorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(
		t,
		LoadingIndicator(feedback.Spinner(feedback.SpinnerLG, "text-blue-600 dark:text-blue-400")),
	)
	utils.AssertContains(t, output, "tc-loading-indicator")
	utils.AssertContains(t, output, "animate-spin")
	utils.AssertContains(t, output, "htmx-indicator")
}

func TestInlineLoadingOverlayRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(
		t,
		InlineLoadingOverlay(
			"form-loader",
			feedback.Spinner(feedback.SpinnerMD, "text-blue-600 dark:text-blue-400"),
		),
	)
	utils.AssertContains(t, output, `id="form-loader"`)
	utils.AssertContains(t, output, "htmx-indicator")
	utils.AssertContains(t, output, "animate-spin")
}

func TestLoadingButtonRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(
		t,
		LoadingButton("Save", "Saving...", feedback.Spinner(feedback.SpinnerSM, "htmx-indicator")),
	)
	utils.AssertContains(t, output, "Save")
	utils.AssertContains(t, output, "Saving...")
	utils.AssertContains(t, output, "animate-spin")
	utils.AssertContains(t, output, "tc-btn-loading")
}

func TestConfirmDeleteRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ConfirmDelete("/api/items/1", "#list", "Delete this item?"))
	utils.AssertContains(t, output, `hx-delete="/api/items/1"`)
	utils.AssertContains(t, output, `hx-target="#list"`)
	utils.AssertContains(t, output, "Delete this item?")
	utils.AssertContains(t, output, "Delete")
}

func TestSwapOOBRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SwapOOB("#notification", "innerHTML"))
	utils.AssertContains(t, output, `hx-swap-oob="innerHTML:#notification"`)
}

func TestCSRFTokenRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CSRFToken("abc123"))
	utils.AssertContains(t, output, `name="csrf_token"`)
	utils.AssertContains(t, output, `value="abc123"`)
	utils.AssertContains(t, output, "hidden")
}
