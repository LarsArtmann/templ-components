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
		output := utils.Render(t, GlobalErrorHandling("secure-nonce"))
		utils.AssertContains(t, output, `nonce="secure-nonce"`)
	})

	t.Run("registers error event listeners", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling("n"))
		utils.AssertContains(t, output, "htmx:sendError")
		utils.AssertContains(t, output, "htmx:responseError")
		utils.AssertContains(t, output, "htmx:afterRequest")
	})

	t.Run("has retry logic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling("n"))
		utils.AssertContains(t, output, "MAX_RETRIES")
	})

	t.Run("exposes error history", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling("n"))
		utils.AssertContains(t, output, "htmxErrorHistory")
	})
}

func TestHTMXDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("loading indicator has dark classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingIndicator(feedback.Spinner(feedback.SpinnerLG, "text-blue-600 dark:text-blue-400")))
		utils.AssertContains(t, output, "dark:")
	})
}
