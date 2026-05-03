// Package feedback provides tests for feedback components like Alert, Toast, Spinner, and Loading.
package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAlertRender(t *testing.T) {
	t.Parallel()
	t.Run("error alert", func(t *testing.T) {
		t.Parallel()
		props := AlertProps{
			Title:       "Error",
			Message:     "Something failed",
			Type:        AlertError,
		}
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, "Error")
		utils.AssertContains(t, output, "Something failed")
		utils.AssertContains(t, output, "bg-red-50")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("dismissible alert", func(t *testing.T) {
		t.Parallel()
		props := AlertProps{
			Title:       "Warning",
			Type:        AlertWarning,
			Dismissible: true,
		}
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, `aria-label="Dismiss"`)
		utils.AssertContains(t, output, `data-dismiss="alert"`)
		utils.AssertNotContains(t, output, `onclick=`)
	})
}

func TestToastContainerRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer("test-nonce"))
	utils.AssertContains(t, output, "tc-toast-container")
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}

func TestToastRender(t *testing.T) {
	t.Parallel()
	props := ToastProps{
		Message: "Saved!",
		Type:    ToastSuccess,
	}
	output := utils.Render(t, Toast(props))
	utils.AssertContains(t, output, "Saved!")
	utils.AssertContains(t, output, `role="alert"`)
	utils.AssertContains(t, output, `data-dismiss="toast"`)
	utils.AssertNotContains(t, output, `onclick=`)
}

func TestInlineErrorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineError("Required"))
	utils.AssertContains(t, output, "Required")
	utils.AssertContains(t, output, "text-red-600")
}

func TestInlineSuccessRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineSuccess("Done"))
	utils.AssertContains(t, output, "Done")
	utils.AssertContains(t, output, "text-green-600")
}

func TestSpinnerRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Spinner(SpinnerMedium, "text-blue-600"))
	utils.AssertContains(t, output, "<svg")
	utils.AssertContains(t, output, "animate-spin")
	utils.AssertContains(t, output, "h-6 w-6")
}

func TestLoadingOverlayRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay("Loading...", true, 45))
	utils.AssertContains(t, output, "Loading...")
	utils.AssertContains(t, output, "fixed inset-0")
	utils.AssertContains(t, output, "width: 45%")
}

func TestInlineLoadingRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineLoading("Saving..."))
	utils.AssertContains(t, output, "Saving...")
	utils.AssertContains(t, output, "animate-spin")
}

func TestSkeletonRender(t *testing.T) {
	t.Parallel()
	t.Run("text variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonText))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "w-3/4")
		utils.AssertContains(t, output, `aria-busy="true"`)
	})
	t.Run("card variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonCard))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "space-y-4")
	})
}

func TestSkeletonGroupRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(
		t,
		SkeletonGroup([]SkeletonVariant{SkeletonTitle, SkeletonText, SkeletonText}),
	)
	utils.AssertContains(t, output, "space-y-3")
	utils.AssertContains(t, output, `aria-busy="true"`)
}

func TestProgressBarRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current:   50,
		Total:     100,
		Label:     "Upload",
		ShowLabel: true,
	}))
	utils.AssertContains(t, output, "Upload")
	utils.AssertContains(t, output, `role="progressbar"`)
	utils.AssertContains(t, output, `aria-valuenow="50"`)
}

func TestStepIndicatorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StepIndicator(StepIndicatorProps{
		Steps:       []string{"Details", "Review", "Confirm"},
		CurrentStep: 1,
	}))
	utils.AssertContains(t, output, "Details")
	utils.AssertContains(t, output, "Review")
	utils.AssertContains(t, output, "Confirm")
	utils.AssertContains(t, output, `aria-current="step"`)
}
