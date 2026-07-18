package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/utils"
)

// Coverage: exercise low-coverage branches in htmx package.
// GlobalErrorHandling (73.7%), LoadingButton (75.6%).

func TestGlobalErrorHandlingCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("renders with custom config", func(t *testing.T) {
		t.Parallel()

		cfg := DefaultErrorHandlingConfig()
		cfg.MaxErrorHistory = 10
		cfg.MaxRetries = 5
		cfg.RetryDelayMS = 2000
		output := utils.Render(t, GlobalErrorHandling(cfg))
		utils.AssertContains(t, output, "script")
	})

	t.Run("renders error announcer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, GlobalErrorHandling(DefaultErrorHandlingConfig()))
		utils.AssertContains(t, output, "tc-error-announcer")
		utils.AssertContains(t, output, `aria-live="polite"`)
	})
}

func TestLoadingButtonCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("renders with default spinner", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingButton("Save", "Saving...", feedback.Spinner(feedback.DefaultSpinnerProps())))
		utils.AssertContains(t, output, "Save")
		utils.AssertContains(t, output, "Saving...")
	})
}
