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

// TestGoldenAlertSuccess and TestGoldenAlertInfo complete golden coverage of
// all four FeedbackType variants (error + dismissible-warning already had
// goldens). Each variant has a distinct icon + color set; a snapshot per type
// catches a recolor or wrong-icon regression that a single-type snapshot would
// miss.
func TestGoldenAlertSuccess(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Alert(AlertProps{
		Title:   "Saved",
		Message: "Your changes are stored.",
		Type:    AlertSuccess,
	}))
	golden.Assert(t, "alert_success", output)
}

func TestGoldenAlertInfo(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Alert(AlertProps{
		Title:   "Heads up",
		Message: "A new version is available.",
		Type:    AlertInfo,
	}))
	golden.Assert(t, "alert_info", output)
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
		BaseProps: utils.BaseProps{ID: "toast-success"},
		Message:   "Saved!",
		Type:      ToastSuccess,
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
	output := utils.Render(t, SkeletonCardGrid(SkeletonCardGridProps{Count: 3}))
	golden.Assert(t, "skeleton_card_grid", output)
}
