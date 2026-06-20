package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAlertFullCoverage(t *testing.T) {
	t.Parallel()
	for _, ft := range []FeedbackType{FeedbackSuccess, FeedbackError, FeedbackWarning, FeedbackInfo} {
		t.Run(string(ft)+" dismissible with BaseProps", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Alert(AlertProps{
				BaseProps: utils.BaseProps{
					ID:        "alert-" + string(ft),
					Class:     "mt-2",
					AriaLabel: string(ft) + " alert",
				},
				Title:       string(ft) + " title",
				Message:     string(ft) + " message",
				Type:        ft,
				Dismissible: true,
			}))
			utils.AssertContains(t, output, string(ft)+" title")
			utils.AssertContains(t, output, "data-dismiss")
		})
	}
}

func TestToastFullCoverage(t *testing.T) {
	t.Parallel()
	for _, ft := range []FeedbackType{FeedbackSuccess, FeedbackError, FeedbackWarning, FeedbackInfo} {
		t.Run(string(ft)+" toast", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Toast(ToastProps{
				BaseProps: utils.BaseProps{ID: "toast-" + string(ft)},
				Title:     string(ft) + " toast",
				Message:   "message",
				Type:      ft,
				Duration:  3000,
			}))
			utils.AssertContains(t, output, string(ft)+" toast")
		})
	}
}

func TestSpinnerFullCoverage(t *testing.T) {
	t.Parallel()
	for _, size := range []SpinnerSize{SpinnerSM, SpinnerMD, SpinnerLG} {
		t.Run("size_"+string(size), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Spinner(SpinnerProps{
				BaseProps: utils.BaseProps{ID: "spin-" + string(size), Class: "mr-1", AriaLabel: "Loading"},
				Size:      size,
			}))
			utils.AssertContains(t, output, "animate-spin")
		})
	}
}

func TestProgressBarFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with label and BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			BaseProps: utils.BaseProps{ID: "pb-1", Class: "w-full", AriaLabel: "Upload progress"},
			Current:   75,
			Total:     100,
			Label:     "Uploading",
			ShowLabel: true,
		}))
		utils.AssertContains(t, output, "75%")
		utils.AssertContains(t, output, "Uploading")
	})
	t.Run("indeterminate", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Indeterminate: true,
		}))
		utils.AssertContains(t, output, "aria-busy")
	})
	for _, size := range []ProgressBarSize{ProgressBarSizeSM, ProgressBarSizeMD, ProgressBarSizeLG} {
		t.Run("size_"+string(size), func(t *testing.T) {
			t.Parallel()
			utils.Render(t, ProgressBar(ProgressBarProps{Current: 50, Total: 100, Size: size}))
		})
	}
}

func TestStepIndicatorFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("horizontal", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			BaseProps:   utils.BaseProps{ID: "steps-h", AriaLabel: "Progress"},
			Steps:       []string{"Account", "Profile", "Verify", "Done"},
			CurrentStep: 2,
		}))
		utils.AssertContains(t, output, "Profile")
	})
	t.Run("vertical", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			BaseProps:   utils.BaseProps{ID: "steps-v"},
			Steps:       []string{"Step 1", "Step 2"},
			CurrentStep: 1,
			Orientation: StepVertical,
		}))
		utils.AssertContains(t, output, "Step 1")
	})
}

func TestLoadingOverlayFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with message and progress", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
			BaseProps:    utils.BaseProps{ID: "ovl-1", AriaLabel: "Loading"},
			Message:      "Processing files",
			ShowProgress: true,
			Progress:     60,
		}))
		utils.AssertContains(t, output, "Processing files")
	})
}

func TestSkeletonFullCoverage(t *testing.T) {
	t.Parallel()
	for _, v := range []SkeletonVariant{SkeletonText, SkeletonTextShort, SkeletonTitle, SkeletonAvatar, SkeletonImage, SkeletonCard, SkeletonTableRow} {
		t.Run(string(v), func(t *testing.T) {
			t.Parallel()
			utils.Render(t, Skeleton(v))
		})
	}
	t.Run("skeleton group", func(t *testing.T) {
		t.Parallel()
		utils.Render(t, SkeletonGroup([]SkeletonVariant{SkeletonText, SkeletonAvatar, SkeletonCard}))
	})
}
