package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAlertEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty title renders without h3", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Message: "Just a message",
			Type:    AlertInfo,
		}))
		utils.AssertContains(t, output, "Just a message")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("empty message renders empty paragraph", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Title: "Title Only",
			Type:  AlertInfo,
		}))
		utils.AssertContains(t, output, "Title Only")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("unknown type falls back to info styling", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Message: "Fallback",
			Type:    FeedbackType("unknown"),
		}))
		utils.AssertContains(t, output, "Fallback")
		utils.AssertContains(t, output, "border-gray-200")
	})

	t.Run("alert with ID propagates to root element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			BaseProps: utils.BaseProps{ID: "my-alert"},
			Message:   "Test",
			Type:      AlertInfo,
		}))
		utils.AssertContains(t, output, `id="my-alert"`)
	})
}

func TestToastEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty message renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			Title: "No Body",
			Type:  ToastInfo,
		}))
		utils.AssertContains(t, output, "No Body")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("unknown type falls back to info styling", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			Message: "Fallback",
			Type:    FeedbackType("unknown"),
		}))
		utils.AssertContains(t, output, "Fallback")
		utils.AssertContains(t, output, "border-gray-200")
	})

	t.Run("toast with ID propagates to root element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			BaseProps: utils.BaseProps{ID: "my-toast"},
			Message:   "Test",
			Type:      ToastInfo,
		}))
		utils.AssertContains(t, output, `id="my-toast"`)
	})
}

func TestProgressBarEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("zero total shows zero percent", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 50,
			Total:   0,
		}))
		utils.AssertContains(t, output, "0%")
	})

	t.Run("negative current shows zero percent", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: -10,
			Total:   100,
		}))
		utils.AssertContains(t, output, "0%")
	})
}

func TestStepIndicatorEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty steps renders empty nav", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps: []string{},
		}))
		utils.AssertContains(t, output, `aria-label="Progress"`)
	})

	t.Run("current step out of bounds renders numbered circles", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps:       []string{"Step 1", "Step 2"},
			CurrentStep: 99,
		}))
		utils.AssertContains(t, output, "Step 1")
		utils.AssertContains(t, output, "Step 2")
	})
}
