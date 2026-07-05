package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenSpinner(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Spinner(SpinnerProps{Size: SpinnerMD, Color: "text-blue-600"}))
	golden.Assert(t, "spinner_md", output)
}

func TestGoldenAlertError(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Alert(AlertProps{
		Title:   "Error",
		Message: "Something failed",
		Type:    AlertError,
	}))
	golden.Assert(t, "alert_error", output)
}

func TestGoldenAlertDismissible(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Alert(AlertProps{
		Title:       "Warning",
		Type:        AlertWarning,
		Dismissible: true,
	}))
	golden.Assert(t, "alert_dismissible", output)
}

func TestGoldenToast(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toast(ToastProps{
		Message: "Saved!",
		Type:    ToastSuccess,
	}))
	golden.Assert(t, "toast_success", output)
}

func TestGoldenProgressBar(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current:   50,
		Total:     100,
		Label:     "Upload",
		ShowLabel: true,
	}))
	golden.Assert(t, "progressbar", output)
}

func TestGoldenStepIndicator(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StepIndicator(StepIndicatorProps{
		Steps:       []string{"Details", "Review", "Confirm"},
		CurrentStep: 1,
	}))
	golden.Assert(t, "step_indicator", output)
}

func TestGoldenLoadingOverlay(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		Message:      "Loading...",
		ShowProgress: true,
		Progress:     45,
	}))
	golden.Assert(t, "loading_overlay", output)
}

func TestGoldenSkeleton(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Skeleton(SkeletonText))
	golden.Assert(t, "skeleton_text", output)
}

func TestGoldenSkeletonCardGrid(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SkeletonCardGrid(3))
	golden.Assert(t, "skeleton_card_grid", output)
}
