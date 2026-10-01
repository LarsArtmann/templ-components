package feedback_test

import (
	"github.com/larsartmann/templ-components/feedback"
)

func ExampleInlineError() {
	_ = feedback.InlineError("Something went wrong")
	// Output:
}

func ExampleInlineSuccess() {
	_ = feedback.InlineSuccess("Changes saved")
	// Output:
}

func ExampleInlineLoading() {
	_ = feedback.InlineLoading("Saving...")
	// Output:
}

func ExampleLoadingOverlay() {
	_ = feedback.LoadingOverlay(feedback.DefaultLoadingOverlayProps())
	// Output:
}

func ExampleSkeleton() {
	_ = feedback.Skeleton(feedback.SkeletonText)
	// Output:
}

func ExampleSkeletonGroup() {
	_ = feedback.SkeletonGroup([]feedback.SkeletonVariant{feedback.SkeletonTitle, feedback.SkeletonText})
	// Output:
}

func ExampleProgressBar() {
	_ = feedback.ProgressBar(feedback.DefaultProgressBarProps())
	// Output:
}

func ExampleStepIndicator() {
	_ = feedback.StepIndicator(feedback.DefaultStepIndicatorProps())
	// Output:
}

func ExampleToastContainer() {
	_ = feedback.ToastContainer("")
	// Output:
}
