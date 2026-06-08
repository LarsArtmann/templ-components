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
			Title:   "Error",
			Message: "Something failed",
			Type:    AlertError,
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
	output := utils.Render(t, Spinner(SpinnerProps{Size: SpinnerMD, Color: "text-blue-600"}))
	utils.AssertContains(t, output, "<svg")
	utils.AssertContains(t, output, "animate-spin")
	utils.AssertContains(t, output, "h-6 w-6")
}

func TestLoadingOverlayRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		Message:      "Loading...",
		ShowProgress: true,
		Progress:     45,
	}))
	utils.AssertContains(t, output, "Loading...")
	utils.AssertContains(t, output, "fixed")
	utils.AssertContains(t, output, "inset-0")
	utils.AssertContains(t, output, "width: 45%")
}

func TestLoadingOverlayCoverage(t *testing.T) {
	t.Parallel()

	t.Run("custom ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
			BaseProps: utils.BaseProps{ID: "overlay", Class: "bg-opacity-75"},
			Message:   "Wait",
		}))
		utils.AssertContains(t, output, `id="overlay"`)
		utils.AssertContains(t, output, "bg-opacity-75")
	})

	t.Run("progress bar shows percentage text", func(t *testing.T) {
		t.Parallel()
		assertLoadingOverlayProgress(t, "Processing", 75, "75%")
		utils.AssertContains(t, renderLoadingOverlayWithProgress(t, "Processing", 75), "width: 75%")
	})

	t.Run("aria-label uses message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
			Message: "Uploading files",
		}))
		utils.AssertContains(t, output, `aria-label="Uploading files"`)
	})
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

func TestDefaultProgressBarProps(t *testing.T) {
	t.Parallel()
	props := DefaultProgressBarProps()
	if props.Size != ProgressBarSizeMD {
		t.Errorf("Size = %q, want %q", props.Size, ProgressBarSizeMD)
	}
	if props.Color == "" {
		t.Error("Color should not be empty")
	}
}

func TestSkeletonVariants(t *testing.T) {
	t.Parallel()
	t.Run("avatar skeleton", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonAvatar))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "rounded-full")
	})
	t.Run("text-short skeleton", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonTextShort))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "w-1/2")
	})
	t.Run("title skeleton", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonTitle))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "h-6")
	})
	t.Run("image skeleton", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonImage))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "h-48")
	})
	t.Run("table-row skeleton", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonTableRow))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "grid-cols-4")
	})
	t.Run("unknown variant uses default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonVariant("unknown")))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "w-full")
	})
}

func TestToastAllTypes(t *testing.T) {
	t.Parallel()
	for _, tt := range []ToastType{ToastSuccess, ToastError, ToastWarning, ToastInfo} {
		t.Run("toast type "+string(tt), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Toast(ToastProps{
				Type:    tt,
				Message: "Test message",
			}))
			utils.AssertContains(t, output, "Test message")
			utils.AssertContains(t, output, `role="alert"`)
		})
	}
}

func TestAlertDismissScript(t *testing.T) {
	t.Parallel()
	t.Run("alert has dismiss script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			BaseProps:   utils.BaseProps{Nonce: "test-nonce"},
			Title:       "Dismiss me",
			Type:        AlertInfo,
			Dismissible: true,
		}))
		utils.AssertContains(t, output, `data-dismiss="alert"`)
		utils.AssertContains(t, output, `nonce="test-nonce"`)
	})
}

func TestDefaultPropsConstructors(t *testing.T) {
	t.Parallel()
	t.Run("DefaultAlertProps", func(t *testing.T) {
		t.Parallel()
		props := DefaultAlertProps()
		if props.Type != AlertInfo {
			t.Errorf("DefaultAlertProps().Type = %q, want %q", props.Type, AlertInfo)
		}
	})
	t.Run("DefaultToastProps", func(t *testing.T) {
		t.Parallel()
		props := DefaultToastProps()
		if props.Type != ToastInfo {
			t.Errorf("DefaultToastProps().Type = %q, want %q", props.Type, ToastInfo)
		}
	})
	t.Run("DefaultStepIndicatorProps", func(t *testing.T) {
		t.Parallel()
		props := DefaultStepIndicatorProps()
		if props.Steps != nil {
			t.Error("DefaultStepIndicatorProps().Steps should be nil")
		}
	})
	t.Run("DefaultLoadingOverlayProps", func(t *testing.T) {
		t.Parallel()
		props := DefaultLoadingOverlayProps()
		if props.Message != "" {
			t.Errorf("DefaultLoadingOverlayProps().Message = %q, want empty", props.Message)
		}
		if props.ShowProgress {
			t.Error("DefaultLoadingOverlayProps().ShowProgress = true, want false")
		}
	})
	t.Run("DefaultProgressBarProps", func(t *testing.T) {
		t.Parallel()
		props := DefaultProgressBarProps()
		if props.Current != 0 || props.Total != 0 {
			t.Errorf(
				"DefaultProgressBarProps() = {Current: %d, Total: %d}, want zeros",
				props.Current,
				props.Total,
			)
		}
	})
}

func TestLoadingOverlayNoProgress(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{Message: "Loading..."}))
	utils.AssertContains(t, output, "Loading...")
	utils.AssertContains(t, output, "Please wait...")
	utils.AssertNotContains(t, output, "width:")
}

func TestInlineLoadingEmptyMessage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineLoading(""))
	utils.AssertContains(t, output, "animate-spin")
}

func TestInlineErrorEmptyMessage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineError(""))
	utils.AssertContains(t, output, "text-red-600")
}

func TestInlineSuccessEmptyMessage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineSuccess(""))
	utils.AssertContains(t, output, "text-green-600")
}

func TestDefaultSpinnerProps(t *testing.T) {
	t.Parallel()
	props := DefaultSpinnerProps()
	if props.Size != SpinnerMD {
		t.Errorf("Size = %q, want %q", props.Size, SpinnerMD)
	}
}

func TestSpinnerWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Spinner(SpinnerProps{
		Size:  SpinnerSM,
		Color: "text-green-600",
		BaseProps: utils.BaseProps{
			ID:        "spin-1",
			Class:     "my-class",
			AriaLabel: "Loading data",
		},
	}))
	utils.AssertContains(t, output, `id="spin-1"`)
	utils.AssertContains(t, output, "my-class")
	utils.AssertContains(t, output, `aria-label="Loading data"`)
	utils.AssertContains(t, output, "text-green-600")
	utils.AssertContains(t, output, "w-4")
	utils.AssertContains(t, output, "h-4")
}
