package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// ---------------------------------------------------------------------------
// Toast: AriaLabel branch
// ---------------------------------------------------------------------------

func TestToastAriaLabelBranch(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toast(ToastProps{
		Message:   "Saved!",
		Type:      FeedbackSuccess,
		BaseProps: utils.BaseProps{AriaLabel: "Success notification"},
	}))
	utils.AssertContains(t, output, `aria-label="Success notification"`)
}

// ---------------------------------------------------------------------------
// ProgressBar: Total==0, negative Current, no label, indeterminate+AriaLabel, invalid size
// ---------------------------------------------------------------------------

func TestProgressBarTotalZero(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current: 5,
		Total:   0,
	}))
	// Should not panic, should render a progress bar
	utils.AssertContains(t, output, "role=\"progressbar\"")
}

func TestProgressBarNegativeCurrent(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current: -10,
		Total:   100,
	}))
	utils.AssertContains(t, output, "role=\"progressbar\"")
}

func TestProgressBarNoLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current:   50,
		Total:     100,
		ShowLabel: false,
		Label:     "",
	}))
	// No label row div should appear (only the bar itself)
	utils.AssertNotContains(t, output, "text-sm font-medium")
}

func TestProgressBarIndeterminateWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Indeterminate: true,
		BaseProps: utils.BaseProps{
			ID:        "pb-indet",
			AriaLabel: "Loading data",
		},
	}))
	utils.AssertContainsAll(t, output, `id="pb-indet"`, `aria-label="Loading data"`)
}

func TestProgressBarInvalidSizeFallback(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current: 50,
		Total:   100,
		Size:    ProgressBarSize("bogus"),
	}))
	utils.AssertContains(t, output, "role=\"progressbar\"")
}

func TestProgressBarAllSizes(t *testing.T) {
	t.Parallel()

	for _, size := range []ProgressBarSize{ProgressBarSizeSM, ProgressBarSizeMD, ProgressBarSizeLG} {
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 50,
			Total:   100,
			Size:    size,
		}))
		utils.AssertContains(t, output, "role=\"progressbar\"")
	}
}

func TestProgressBarCustomColor(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ProgressBar(ProgressBarProps{
		Current: 50,
		Total:   100,
		Color:   "bg-green-600",
	}))
	utils.AssertContains(t, output, "bg-green-600")
}

// ---------------------------------------------------------------------------
// LoadingOverlay: no ID, AriaLabel empty, ShowProgress=false, Progress clamping
// ---------------------------------------------------------------------------

func TestLoadingOverlayMinimal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		Message: "Please wait",
	}))
	utils.AssertContains(t, output, "Please wait")
}

func TestLoadingOverlayNoAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		Message: "Working...",
	}))
	// Should fall back to Message for aria-label
	utils.AssertContains(t, output, `aria-label="Working..."`)
}

func TestLoadingOverlayShowProgressFalse(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{}))
	utils.AssertContains(t, output, "Please wait...")
}

func TestLoadingOverlayProgressClamping(t *testing.T) {
	t.Parallel()
	output1 := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		ShowProgress: true,
		Progress:     150,
	}))
	utils.AssertContains(t, output1, "100%")

	output2 := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		ShowProgress: true,
		Progress:     -20,
	}))
	utils.AssertContains(t, output2, "0%")
}

// ---------------------------------------------------------------------------
// InlineLoading: entirely untested
// ---------------------------------------------------------------------------

func TestInlineLoading(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineLoading("Saving changes..."))
	utils.AssertContains(t, output, "Saving changes...")
}

// ---------------------------------------------------------------------------
// Spinner: invalid size fallback
// ---------------------------------------------------------------------------

func TestSpinnerInvalidSizeFallback(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Spinner(SpinnerProps{
		Size: SpinnerSize("bogus"),
	}))
	// Should still render without panic
	utils.AssertContains(t, output, "animate-spin")
}

func TestSpinnerCustomColor(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Spinner(SpinnerProps{
		Color: "text-red-600",
	}))
	utils.AssertContains(t, output, "text-red-600")
}

// ---------------------------------------------------------------------------
// SkeletonCardGrid: both branches (n<=0 fallback, n>0)
// ---------------------------------------------------------------------------

func TestSkeletonCardGridNormal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SkeletonCardGrid(3))
	utils.AssertContains(t, output, "role=\"status\"")
}

func TestSkeletonCardGridZeroOrNegative(t *testing.T) {
	t.Parallel()

	for _, count := range []int{0, -1, -5} {
		output := utils.Render(t, SkeletonCardGrid(count))
		// Should fall back to 1 card, not panic
		utils.AssertContains(t, output, "role=\"status\"")
	}
}

// ---------------------------------------------------------------------------
// Alert: InlineError, InlineSuccess
// ---------------------------------------------------------------------------

func TestInlineErrorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineError("Email is required"))
	utils.AssertContainsAll(t, output, "Email is required", "alert")
}

func TestInlineSuccessCoverage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineSuccess("Saved successfully"))
	utils.AssertContainsAll(t, output, "Saved successfully", "success")
}

// ---------------------------------------------------------------------------
// StepIndicator: CurrentStep == 0 (all future steps)
// ---------------------------------------------------------------------------

func TestStepIndicatorCurrentStepZero(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StepIndicator(StepIndicatorProps{
		Steps:       []string{"Upload", "Review", "Confirm"},
		CurrentStep: 0,
	}))
	utils.AssertContains(t, output, "Upload")
	utils.AssertContains(t, output, "Review")
}

func TestStepIndicatorVerticalWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StepIndicator(StepIndicatorProps{
		Steps:       []string{"Step 1", "Step 2"},
		CurrentStep: 1,
		Orientation: StepVertical,
		BaseProps: utils.BaseProps{
			ID:        "step-v",
			AriaLabel: "Onboarding progress",
		},
	}))
	utils.AssertContainsAll(t, output, `id="step-v"`, "Step 1", "Step 2")
}

func TestStepIndicatorLastStep(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StepIndicator(StepIndicatorProps{
		Steps:       []string{"A", "B", "C"},
		CurrentStep: 2, // last step complete
	}))
	utils.AssertContains(t, output, "C")
}

// ---------------------------------------------------------------------------
// Alert: all types + dismissible without nonce
// ---------------------------------------------------------------------------

func TestAlertAllTypes(t *testing.T) {
	t.Parallel()

	for _, atype := range []AlertType{FeedbackSuccess, FeedbackError, FeedbackWarning, FeedbackInfo} {
		output := utils.Render(t, Alert(AlertProps{
			Message: "Test message",
			Type:    atype,
		}))
		utils.AssertContains(t, output, "Test message")
	}
}

func TestAlertDismissibleWithNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Alert(AlertProps{
		Message:     "Dismiss me",
		Type:        FeedbackInfo,
		Dismissible: true,
		BaseProps:   utils.BaseProps{Nonce: "abc123"},
	}))
	utils.AssertContainsAll(t, output, "Dismiss me", `nonce="abc123"`)
}

// ---------------------------------------------------------------------------
// Toast: all types
// ---------------------------------------------------------------------------

func TestToastEachType(t *testing.T) {
	t.Parallel()

	for _, ttype := range []ToastType{FeedbackSuccess, FeedbackError, FeedbackWarning, FeedbackInfo} {
		output := utils.Render(t, Toast(ToastProps{
			Message: "Notification",
			Type:    ttype,
		}))
		utils.AssertContains(t, output, "Notification")
	}
}

func TestToastWithTitleAndDuration(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toast(ToastProps{
		Title:    "Heads up",
		Message:  "This is a warning",
		Type:     FeedbackWarning,
		Duration: 5000,
		BaseProps: utils.BaseProps{
			ID:    "toast-1",
			Nonce: "nonce-abc",
		},
	}))
	utils.AssertContainsAll(t, output, "Heads up", "This is a warning", `id="toast-1"`)
}

// ---------------------------------------------------------------------------
// Skeleton: unknown variant default case
// ---------------------------------------------------------------------------

func TestSkeletonUnknownVariantDefault(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Skeleton(SkeletonVariant("totally-unknown")))
	utils.AssertContains(t, output, "animate-pulse")
}

func TestSkeletonAllVariants(t *testing.T) {
	t.Parallel()

	for _, v := range []SkeletonVariant{
		SkeletonText, SkeletonTextShort, SkeletonTitle,
		SkeletonAvatar, SkeletonImage, SkeletonCard, SkeletonTableRow,
	} {
		output := utils.Render(t, Skeleton(v))
		utils.AssertContains(t, output, "animate-pulse")
	}
}
