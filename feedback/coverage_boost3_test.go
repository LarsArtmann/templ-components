package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// Coverage: exercise low-coverage branches in feedback package.
// StepIndicator (71.8%), ProgressBar (72.8%), LoadingOverlay (72.6%), Spinner (73.8%).

func TestStepIndicatorCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("vertical orientation", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Orientation: StepVertical,
			Steps:       []string{"Step 1", "Step 2", "Step 3"},
			CurrentStep: 1,
		}))
		utils.AssertContains(t, output, "Step 1")
		utils.AssertContains(t, output, "Step 2")
		utils.AssertContains(t, output, "Step 3")
	})

	t.Run("horizontal orientation", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			Orientation: StepHorizontal,
			Steps:       []string{"Create", "Review", "Deploy"},
			CurrentStep: 1,
		}))
		utils.AssertContains(t, output, "Create")
		utils.AssertContains(t, output, "Deploy")
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			BaseProps:   utils.BaseProps{AriaLabel: "Checkout progress"},
			Steps:       []string{"Cart"},
			CurrentStep: 0,
		}))
		utils.AssertContains(t, output, `aria-label="Checkout progress"`)
	})

	t.Run("propagates custom class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StepIndicator(StepIndicatorProps{
			BaseProps:   utils.BaseProps{Class: "my-steps"},
			Steps:       []string{"X"},
			CurrentStep: 0,
		}))
		utils.AssertContains(t, output, "my-steps")
	})
}

func TestProgressBarCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("indeterminate bar", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Indeterminate: true,
		}))
		utils.AssertContains(t, output, `aria-busy="true"`)
	})

	t.Run("with label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Label:   "Uploading...",
			Current: 50,
			Total:   100,
		}))
		utils.AssertContains(t, output, "Uploading...")
	})

	t.Run("aria-valuenow clamped to total", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 150,
			Total:   100,
		}))
		utils.AssertContains(t, output, `aria-valuenow="100"`)
	})

	t.Run("aria-valuenow clamped to zero", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: -10,
			Total:   100,
		}))
		utils.AssertContains(t, output, `aria-valuenow="0"`)
	})

	t.Run("custom total", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 3,
			Total:   7,
		}))
		utils.AssertContains(t, output, `aria-valuemax="7"`)
	})
}

func TestLoadingOverlayCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("renders overlay with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
			BaseProps: utils.BaseProps{AriaLabel: "Loading page"},
		}))
		utils.AssertContains(t, output, "Loading page")
	})

	t.Run("propagates class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{
			BaseProps: utils.BaseProps{Class: "my-overlay"},
		}))
		utils.AssertContains(t, output, "my-overlay")
	})
}

func TestSpinnerCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("large spinner", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Spinner(SpinnerProps{Size: SpinnerLG}))
		utils.AssertContains(t, output, "animate-spin")
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Spinner(SpinnerProps{
			BaseProps: utils.BaseProps{AriaLabel: "Loading data"},
		}))
		utils.AssertContains(t, output, `aria-label="Loading data"`)
	})
}

func TestSkeletonCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("text variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonText))
		utils.AssertContains(t, output, "animate-pulse")
	})

	t.Run("card variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonCard))
		utils.AssertContains(t, output, "animate-pulse")
	})

	t.Run("circle variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonAvatar))
		utils.AssertContains(t, output, "animate-pulse")
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("skeleton group", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SkeletonGroup([]SkeletonVariant{SkeletonText, SkeletonCard, SkeletonAvatar}))
		utils.AssertContains(t, output, `role="status"`)
	})
}
