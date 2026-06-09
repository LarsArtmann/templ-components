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

	t.Run("registers error event listeners and defaults", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, output, "htmx:sendError")
		utils.AssertContains(t, output, "htmx:responseError")
		utils.AssertContains(t, output, "htmx:afterRequest")
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
