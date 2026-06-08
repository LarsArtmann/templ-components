package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/utils"
)

func TestGlobalErrorHandlingA11y(t *testing.T) {
	t.Parallel()

	t.Run("script uses nonce", func(t *testing.T) {
		t.Parallel()
		cfg := ErrorHandlingConfig{Nonce: "secure-nonce"}
		output := utils.Render(t, GlobalErrorHandling(cfg))
		utils.AssertContains(t, output, `nonce="secure-nonce"`)
	})

	t.Run("registers error event listeners", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, output, "htmx:sendError")
		utils.AssertContains(t, output, "htmx:responseError")
		utils.AssertContains(t, output, "htmx:afterRequest")
	})

	t.Run("has retry logic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, output, "MAX_RETRIES")
	})

	t.Run("exposes error history", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, output, "htmxErrorHistory")
	})

	t.Run("uses toast container for aria-live announcements", func(t *testing.T) {
		t.Parallel()
		toastOutput := utils.Render(t, feedback.ToastContainer("n"))
		utils.AssertContains(t, toastOutput, `aria-live="polite"`)
		utils.AssertContains(t, toastOutput, "tc-toast-container")

		errorOutput := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, errorOutput, "tcShowToast")
	})

	t.Run("configurable values override defaults", func(t *testing.T) {
		t.Parallel()
		cfg := ErrorHandlingConfig{
			Nonce:           "n",
			MaxErrorHistory: 20,
			MaxRetries:      5,
			RetryDelayMS:    2000,
		}
		output := utils.Render(t, GlobalErrorHandling(cfg))
		utils.AssertContains(t, output, "MAX_ERROR_HISTORY = 20")
		utils.AssertContains(t, output, "MAX_RETRIES = 5")
		utils.AssertContains(t, output, "RETRY_DELAY_MS = 2000")
	})

	t.Run("default config values", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, output, "MAX_ERROR_HISTORY = 10")
		utils.AssertContains(t, output, "MAX_RETRIES = 2")
		utils.AssertContains(t, output, "RETRY_DELAY_MS = 1000")
	})
}

func TestHTMXDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("loading indicator has dark classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			LoadingIndicator(
				feedback.Spinner(
					feedback.SpinnerProps{Size: feedback.SpinnerLG, Color: "text-blue-600 dark:text-blue-400"},
				),
			),
		)
		utils.AssertContains(t, output, "dark:")
	})
}
