// Package htmx provides behavior-driven tests for HTMX integration components.
// These tests verify end-user experience: loading states, error handling, form submissions.
package htmx

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/utils"
)

// --- LoadingIndicator Behavior ---

func TestLoadingIndicatorUserSeesLoadingFeedback(t *testing.T) {
	t.Parallel()

	t.Run("user sees loading overlay when HTMX request starts", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingIndicator(feedback.Spinner(feedback.SpinnerLG, "text-blue-600 dark:text-blue-400")))
		utils.AssertContains(t, output, `id="tc-loading-indicator"`)
		utils.AssertContains(t, output, "htmx-indicator")
	})

	t.Run("loading indicator is accessible to screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingIndicator(feedback.Spinner(feedback.SpinnerLG, "text-blue-600 dark:text-blue-400")))
		utils.AssertContains(t, output, `role="status"`)
		utils.AssertContains(t, output, `aria-live="polite"`)
	})
}

// --- InlineLoadingOverlay Behavior ---

func TestInlineLoadingOverlayUserSeesLocalLoadingState(t *testing.T) {
	t.Parallel()

	t.Run("user sees loading overlay on specific form area", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, InlineLoadingOverlay("my-form-loading", feedback.Spinner(feedback.SpinnerMD, "text-blue-600 dark:text-blue-400")))
		utils.AssertContains(t, output, `id="my-form-loading"`)
		utils.AssertContains(t, output, "htmx-indicator")
	})
}

// --- LoadingButton Behavior ---

func TestLoadingButtonUserSeesButtonStateChange(t *testing.T) {
	t.Parallel()

	t.Run("user sees button with default and loading text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingButton("Save", "Saving...", feedback.Spinner(feedback.SpinnerSM, "htmx-indicator")))
		utils.AssertContains(t, output, "Save")
		utils.AssertContains(t, output, "Saving...")
	})
}

// --- ConfirmDelete Behavior ---

func TestConfirmDeleteUserGetsConfirmationDialog(t *testing.T) {
	t.Parallel()

	t.Run("user sees delete button with confirmation", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			ConfirmDelete("/api/items/1", "#item-1", "Are you sure you want to delete?"),
		)
		utils.AssertContains(t, output, `hx-delete="/api/items/1"`)
		utils.AssertContains(t, output, `hx-target="#item-1"`)
		utils.AssertContains(t, output, `hx-confirm="Are you sure you want to delete?"`)
		utils.AssertContains(t, output, "Delete")
	})
}

// --- CSRFToken Behavior ---

func TestCSRFTokenProtectsFormSubmissions(t *testing.T) {
	t.Parallel()

	t.Run("form includes hidden CSRF token input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CSRFToken("secret-token-123"))
		utils.AssertContains(t, output, `type="hidden"`)
		utils.AssertContains(t, output, `name="csrf_token"`)
		utils.AssertContains(t, output, `value="secret-token-123"`)
	})
}

// --- SwapOOB Behavior ---

func TestSwapOOBUpdatesMultipleElements(t *testing.T) {
	t.Parallel()

	t.Run("user gets out-of-band update for targeted element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SwapOOB("#toast-container", "beforeend"))
		utils.AssertContains(t, output, `hx-swap-oob="beforeend:#toast-container"`)
	})
}

// --- GlobalErrorHandling Behavior ---

func TestGlobalErrorHandlingUserGetsErrorFeedback(t *testing.T) {
	t.Parallel()

	t.Run("error handling script is included with CSP nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling("test-nonce-abc"))
		utils.AssertContains(t, output, `<script nonce="test-nonce-abc"`)
	})

	t.Run("error handler registers HTMX event listeners", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling("nonce"))
		utils.AssertContains(t, output, "htmx:responseError")
		utils.AssertContains(t, output, "htmx:sendError")
	})
}

// --- Edge Cases ---

func TestHTMXComponentsRenderValidHTML(t *testing.T) {
	t.Parallel()

	t.Run("all HTMX components render without errors", func(t *testing.T) {
		t.Parallel()
		components := []struct {
			name string
			comp func() templ.Component
		}{
			{"LoadingIndicator", func() templ.Component { return LoadingIndicator(feedback.Spinner(feedback.SpinnerLG, "text-blue-600")) }},
			{
				"InlineLoadingOverlay",
				func() templ.Component { return InlineLoadingOverlay("test", feedback.Spinner(feedback.SpinnerMD, "text-blue-600")) },
			},
			{"LoadingButton", func() templ.Component { return LoadingButton("Go", "Going...", feedback.Spinner(feedback.SpinnerSM, "htmx-indicator")) }},
			{
				"ConfirmDelete",
				func() templ.Component { return ConfirmDelete("/del", "#t", "Sure?") },
			},
			{"CSRFToken", func() templ.Component { return CSRFToken("tok") }},
			{"GlobalErrorHandling", func() templ.Component { return GlobalErrorHandling("n") }},
		}
		for _, tc := range components {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, tc.comp())
				if !strings.Contains(output, "<") {
					t.Errorf(
						"expected HTML output for %s, got: %s",
						tc.name,
						output[:min(len(output), 100)],
					)
				}
			})
		}
	})
}
