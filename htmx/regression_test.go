package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestInlineLoadingOverlayAccessibility verifies the a11y properties added in
// the Round-2 fix: role="status", aria-live="polite", and sr-only loading text.
func TestInlineLoadingOverlayAccessibility(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineLoadingOverlay("test-loading", nil))
	utils.AssertContains(t, output, `role="status"`)
	utils.AssertContains(t, output, `aria-live="polite"`)
	utils.AssertContains(t, output, "Loading…")
	utils.AssertContains(t, output, `id="test-loading"`)
}

// TestLoadingIndicatorAccessibility verifies the global loading indicator has
// the same a11y properties.
func TestLoadingIndicatorAccessibility(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingIndicator(nil))
	utils.AssertContains(t, output, `role="status"`)
	utils.AssertContains(t, output, `aria-live="polite"`)
	utils.AssertContains(t, output, "Loading…")
	utils.AssertContains(t, output, `id="tc-loading-indicator"`)
}

// TestLoadingButtonHidesDefaultText verifies the [.htmx-request_&]:hidden
// arbitrary variant is present so default text hides during HTMX requests.
func TestLoadingButtonHidesDefaultText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingButton("Save", "Saving...", nil))
	utils.AssertContains(t, output, "[.htmx-request_&]:hidden")
	utils.AssertContains(t, output, "Save")
	utils.AssertContains(t, output, "Saving...")
	utils.AssertContains(t, output, "htmx-indicator")
}

// TestLoadingButtonNilSpinner does not panic with nil spinner.
func TestLoadingButtonNilSpinner(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingButton("Save", "Saving...", nil))
	if output == "" {
		t.Error("expected non-empty output")
	}
}
