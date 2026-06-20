package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestProgressBarCoverage(t *testing.T) {
	t.Parallel()

	t.Run("current exceeds total caps at 100%", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 150, Total: 100, ShowLabel: true,
		}))
		utils.AssertContains(t, output, "100%")
	})

	t.Run("with label shows percentage", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 45, Total: 100, Label: "Uploading", ShowLabel: true,
		}))
		utils.AssertContains(t, output, "Uploading")
		utils.AssertContains(t, output, "45%")
	})

	t.Run("label without showLabel still renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 30, Total: 100, Label: "Processing",
		}))
		utils.AssertContains(t, output, "Processing")
		utils.AssertContains(t, output, "30%")
	})

	t.Run("indeterminate with label hides percentage", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Indeterminate: true, Label: "Loading", ShowLabel: true,
		}))
		utils.AssertContains(t, output, "Loading")
		utils.AssertNotContains(t, output, "%")
	})

	t.Run("size SM", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Size: ProgressBarSizeSM,
		}))
		utils.AssertContains(t, output, "h-1.5")
	})

	t.Run("size LG", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Size: ProgressBarSizeLG,
		}))
		utils.AssertContains(t, output, "h-4")
	})

	t.Run("with ID and AriaLabel", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			BaseProps: utils.BaseProps{ID: "pb1", AriaLabel: "Upload progress"},
			Current:   50, Total: 100,
		}))
		utils.AssertContains(t, output, `id="pb1"`)
		utils.AssertContains(t, output, `aria-label="Upload progress"`)
	})

	t.Run("with custom color", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 50, Total: 100, Color: "bg-green-600 dark:bg-green-500",
		}))
		utils.AssertContains(t, output, "bg-green-600")
	})

	t.Run("default props renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(DefaultProgressBarProps()))
		utils.AssertContains(t, output, "bg-blue-600")
	})
}

func TestToastCoverage(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name      string
		fbType    FeedbackType
		wantClass string
	}{
		{"success", FeedbackSuccess, "green"},
		{"error", FeedbackError, "red"},
		{"warning", FeedbackWarning, "yellow"},
		{"info", FeedbackInfo, "blue"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Toast(ToastProps{
				Message: tt.name + " message",
				Type:    tt.fbType,
			}))
			utils.AssertContains(t, output, tt.wantClass)
			utils.AssertContains(t, output, tt.name+" message")
		})
	}

	t.Run("with title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			Title: "Heads up", Message: "Check this", Type: FeedbackInfo,
		}))
		utils.AssertContains(t, output, "Heads up")
		utils.AssertContains(t, output, "Check this")
	})

	t.Run("with nonce renders nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			Message: "test", Type: FeedbackInfo,
			BaseProps: utils.BaseProps{Nonce: "nonce123"},
		}))
		utils.AssertContains(t, output, `nonce="nonce123"`)
	})
}

func TestAlertCoverage(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name      string
		alertType FeedbackType
		wantClass string
	}{
		{"success", FeedbackSuccess, "green"},
		{"error", FeedbackError, "red"},
		{"warning", FeedbackWarning, "yellow"},
		{"info", FeedbackInfo, "blue"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Alert(AlertProps{
				Message: tt.name + " alert",
				Type:    tt.alertType,
			}))
			utils.AssertContains(t, output, tt.wantClass)
			utils.AssertContains(t, output, tt.name+" alert")
		})
	}

	t.Run("with title and dismissible", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Title: "Warning Title", Message: "Be careful", Type: FeedbackWarning, Dismissible: true,
		}))
		utils.AssertContains(t, output, "Warning Title")
		utils.AssertContains(t, output, "data-dismiss")
	})

	t.Run("with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Message: "test", Type: FeedbackInfo, Dismissible: true,
			BaseProps: utils.BaseProps{Nonce: "n1"},
		}))
		utils.AssertContains(t, output, `nonce="n1"`)
	})
}

func TestStepIndicatorCoverage(t *testing.T) {
	t.Parallel()

	t.Run("vertical orientation", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps:       []string{"Step 1", "Step 2", "Step 3"},
			CurrentStep: 1,
			Orientation: StepVertical,
		}))
		utils.AssertContains(t, output, "Step 1")
		utils.AssertContains(t, output, "flex-col")
	})

	t.Run("last step current", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps:       []string{"A", "B", "C"},
			CurrentStep: 2,
		}))
		utils.AssertContains(t, output, "C")
	})

	t.Run("with ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Steps:     []string{"One"},
			BaseProps: utils.BaseProps{ID: "steps", Class: "mt-4"},
		}))
		utils.AssertContains(t, output, `id="steps"`)
		utils.AssertContains(t, output, "mt-4")
	})

	t.Run("default props renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(DefaultStepIndicatorProps()))
		utils.AssertContains(t, output, `aria-label="Progress"`)
	})
}

func TestSpinnerCoverage(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		size SpinnerSize
		want string
	}{
		{"SM", SpinnerSM, "h-4"},
		{"MD", SpinnerMD, "h-6"},
		{"LG", SpinnerLG, "h-8"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Spinner(SpinnerProps{Size: tt.size}))
			utils.AssertContains(t, output, tt.want)
		})
	}

	t.Run("custom color and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Spinner(SpinnerProps{
			Color:     "text-red-600",
			BaseProps: utils.BaseProps{Class: "my-2"},
		}))
		utils.AssertContains(t, output, "text-red-600")
		utils.AssertContains(t, output, "my-2")
	})
}

func TestSkeletonGroupCoverage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SkeletonGroup([]SkeletonVariant{SkeletonCard, SkeletonTextShort}))
	utils.AssertContains(t, output, "animate-pulse")
}

func TestToastContainerCoverage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer("test-nonce"))
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}
