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

func renderLoadingIndicator(t *testing.T) string {
	t.Helper()
	return utils.Render(
		t,
		LoadingIndicator(
			feedback.Spinner(feedback.SpinnerLG, "text-blue-600 dark:text-blue-400"),
		),
	)
}

// --- LoadingIndicator Behavior ---

func TestLoadingIndicatorUserSeesLoadingFeedback(t *testing.T) {
	t.Parallel()

	output := renderLoadingIndicator(t)
	for _, tt := range []struct {
		name string
		want string
	}{
		{"has indicator id", `id="tc-loading-indicator"`},
		{"has htmx-indicator class", "htmx-indicator"},
		{"has role=status", `role="status"`},
		{"has aria-live", `aria-live="polite"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertContains(t, output, tt.want)
		})
	}
}

// --- InlineLoadingOverlay Behavior ---

func TestInlineLoadingOverlayUserSeesLocalLoadingState(t *testing.T) {
	t.Parallel()

	t.Run("user sees loading overlay on specific form area", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			InlineLoadingOverlay(
				"my-form-loading",
				feedback.Spinner(feedback.SpinnerMD, "text-blue-600 dark:text-blue-400"),
			),
		)
		utils.AssertContains(t, output, `id="my-form-loading"`)
		utils.AssertContains(t, output, "htmx-indicator")
	})
}

// --- LoadingButton Behavior ---

func TestLoadingButtonUserSeesButtonStateChange(t *testing.T) {
	t.Parallel()

	t.Run("user sees button with default and loading text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			LoadingButton(
				"Save",
				"Saving...",
				feedback.Spinner(feedback.SpinnerSM, "htmx-indicator"),
			),
		)
		utils.AssertContains(t, output, "Save")
		utils.AssertContains(t, output, "Saving...")
	})
}

func renderConfirmDelete(t *testing.T, endpoint, target, confirmMsg string) string {
	t.Helper()
	return utils.Render(t, ConfirmDelete(endpoint, target, confirmMsg))
}

func assertConfirmDeleteContains(
	t *testing.T,
	endpoint, target, confirmMsg string,
	extraContains ...string,
) {
	t.Helper()
	output := renderConfirmDelete(t, endpoint, target, confirmMsg)
	utils.AssertContains(t, output, `hx-delete="`+endpoint+`"`)
	utils.AssertContains(t, output, `hx-target="`+target+`"`)
	utils.AssertContains(t, output, confirmMsg)
	utils.AssertContains(t, output, "Delete")
	for _, s := range extraContains {
		utils.AssertContains(t, output, s)
	}
}

// --- ConfirmDelete Behavior ---

func TestConfirmDeleteUserGetsConfirmationDialog(t *testing.T) {
	t.Parallel()

	t.Run("user sees delete button with confirmation", func(t *testing.T) {
		t.Parallel()
		assertConfirmDeleteContains(
			t,
			"/api/items/1",
			"#item-1",
			"Are you sure you want to delete?",
		)
	})
}

func renderCSRFToken(t *testing.T, token string) string {
	t.Helper()
	return utils.Render(t, CSRFToken(token))
}

// --- CSRFToken Behavior ---

func TestCSRFTokenProtectsFormSubmissions(t *testing.T) {
	t.Parallel()

	t.Run("form includes hidden CSRF token input", func(t *testing.T) {
		t.Parallel()
		output := renderCSRFToken(t, "secret-token-123")
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
		cfg := ErrorHandlingConfig{Nonce: "test-nonce-abc"}
		output := utils.Render(t, GlobalErrorHandling(cfg))
		utils.AssertContains(t, output, `<script nonce="test-nonce-abc"`)
	})

	t.Run("error handler registers HTMX event listeners", func(t *testing.T) {
		t.Parallel()
		cfg := ErrorHandlingConfig{Nonce: "nonce"}
		output := utils.Render(t, GlobalErrorHandling(cfg))
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
			{
				"LoadingIndicator",
				func() templ.Component { return LoadingIndicator(feedback.Spinner(feedback.SpinnerLG, "text-blue-600")) },
			},
			{
				"InlineLoadingOverlay",
				func() templ.Component {
					return InlineLoadingOverlay(
						"test",
						feedback.Spinner(feedback.SpinnerMD, "text-blue-600"),
					)
				},
			},
			{"LoadingButton", func() templ.Component {
				return LoadingButton(
					"Go",
					"Going...",
					feedback.Spinner(feedback.SpinnerSM, "htmx-indicator"),
				)
			}},
			{
				"ConfirmDelete",
				func() templ.Component { return ConfirmDelete("/del", "#t", "Sure?") },
			},
			{"CSRFToken", func() templ.Component { return CSRFToken("tok") }},
			{
				"GlobalErrorHandling",
				func() templ.Component { return GlobalErrorHandling(ErrorHandlingConfig{Nonce: "n"}) },
			},
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

func TestCSRFTokenEmptyString(t *testing.T) {
	t.Parallel()
	t.Run("empty token still renders hidden input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CSRFToken(""))
		utils.AssertContains(t, output, `type="hidden"`)
		utils.AssertContains(t, output, `name="csrf_token"`)
		utils.AssertContains(t, output, `value=""`)
	})
}

func TestSwapOOBValidatesSwapStyle(t *testing.T) {
	t.Parallel()
	t.Run("valid swapStyle renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SwapOOB("#target", "beforeend"))
		utils.AssertContains(t, output, `hx-swap-oob="beforeend:#target"`)
	})
}
