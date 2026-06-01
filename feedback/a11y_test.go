package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDarkModeClasses(t *testing.T) {
	t.Parallel()

	t.Run("alert has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultAlertProps()
		props.Title = "Test"
		props.Message = "Message"
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("toast has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultToastProps()
		props.Message = "Test"
		output := utils.Render(t, Toast(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("loadingoverlay has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultLoadingOverlayProps()
		props.Message = "Loading"
		output := utils.Render(t, LoadingOverlay(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("progressbar has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultProgressBarProps()
		props.Current = 50
		props.Total = 100
		output := utils.Render(t, ProgressBar(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("stepindicator has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultStepIndicatorProps()
		props.Steps = []string{"A", "B", "C"}
		props.CurrentStep = 1
		output := utils.Render(t, StepIndicator(props))
		utils.AssertContains(t, output, "dark:")
	})
}
