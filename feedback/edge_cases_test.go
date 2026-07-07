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
		utils.AssertContains(t, output, `role="status"`)
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

	t.Run("auto-dismiss fires even without ID", func(t *testing.T) {
		t.Parallel()
		// DefaultToastProps sets Duration: 5000 but no ID. Before the fix,
		// the absence of an ID silently disabled auto-dismiss.
		output := utils.Render(t, Toast(DefaultToastProps()))
		utils.AssertContains(t, output, `id="tc-toast-`)
		utils.AssertContains(t, output, "setTimeout")
	})

	t.Run("duration zero omits setTimeout", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			BaseProps: utils.BaseProps{ID: "manual-toast"},
			Message:   "Manual",
			Type:      ToastInfo,
		}))
		utils.AssertNotContains(t, output, "setTimeout")
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
		utils.AssertContains(t, output, `aria-valuenow="0"`)
	})

	t.Run("overflow current clamps aria-valuenow to total", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 250,
			Total:   100,
		}))
		utils.AssertContains(t, output, `aria-valuenow="100"`)
		utils.AssertContains(t, output, `aria-valuemax="100"`)
		utils.AssertContains(t, output, "100%")
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

func TestLoadingOverlayEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props LoadingOverlayProps
		want  []string
	}{
		{"no message", LoadingOverlayProps{}, []string{"Please wait..."}},
		{"custom id/class", LoadingOverlayProps{BaseProps: utils.BaseProps{ID: "lo", Class: "mt-2"}, Message: "Loading"}, []string{`id="lo"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, LoadingOverlay(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestSkeletonEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		variant SkeletonVariant
		want    []string
	}{
		{"card variant", SkeletonCard, []string{"space-y-4"}},
		{"avatar variant", SkeletonAvatar, []string{"rounded-full"}},
		{"text-short variant", SkeletonTextShort, []string{"w-1/2"}},
		{"title variant", SkeletonTitle, []string{"h-6"}},
		{"image variant", SkeletonImage, []string{"h-48"}},
		{"table-row variant", SkeletonTableRow, []string{"grid-cols-4"}},
		{"unknown variant", SkeletonVariant("unknown"), []string{"animate-pulse"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Skeleton(tt.variant))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestProgressBarIndeterminate(t *testing.T) {
	t.Parallel()

	t.Run("indeterminate renders aria-busy", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Indeterminate: true,
		}))
		utils.AssertContains(t, output, `aria-busy="true"`)
		utils.AssertContains(t, output, "animate-indeterminate-progress")
		utils.AssertNotContains(t, output, "aria-valuenow")
	})

	t.Run("indeterminate with label hides percentage", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Label:         "Loading...",
			ShowLabel:     true,
			Indeterminate: true,
		}))
		utils.AssertContains(t, output, "Loading...")
		utils.AssertNotContains(t, output, "%")
	})

	t.Run("determinate shows percentage", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current:   50,
			Total:     100,
			ShowLabel: true,
		}))
		utils.AssertContains(t, output, "50%")
		utils.AssertContains(t, output, `aria-valuenow="50"`)
	})
}

func TestStepIndicatorVertical(t *testing.T) {
	t.Parallel()

	t.Run("vertical layout uses flex-col", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps:       []string{"Details", "Review", "Confirm"},
			CurrentStep: 1,
			Orientation: StepVertical,
		}))
		utils.AssertContains(t, output, "flex-col")
		utils.AssertContains(t, output, "Details")
		utils.AssertContains(t, output, `aria-current="step"`)
		utils.AssertNotContains(t, output, "flex-1 flex items-center")
	})

	t.Run("horizontal is default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps: []string{"A", "B"},
		}))
		utils.AssertContains(t, output, "flex-1")
		utils.AssertNotContains(t, output, "space-y-4")
	})
}

func TestDefaultLoadingOverlayProps(t *testing.T) {
	t.Parallel()
	props := DefaultLoadingOverlayProps()
	if props.ShowProgress {
		t.Error("DefaultLoadingOverlayProps ShowProgress should be false")
	}
	if props.Message != "" {
		t.Error("DefaultLoadingOverlayProps Message should be empty")
	}
	output := utils.Render(t, LoadingOverlay(props))
	utils.AssertContains(t, output, "fixed")
	utils.AssertContains(t, output, `role="dialog"`)
}
