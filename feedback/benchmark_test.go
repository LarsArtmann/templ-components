package feedback

import (
	"bytes"
	"context"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func BenchmarkFeedbackRenders(b *testing.B) {
	b.Run("Alert render", func(b *testing.B) {
		props := AlertProps{Title: "Error", Message: "Something failed", Type: AlertError}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Alert(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Toast render", func(b *testing.B) {
		props := ToastProps{Message: "Saved!", Type: ToastSuccess}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Toast(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Spinner render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Spinner(SpinnerProps{Size: SpinnerMD, Color: "text-blue-600"}).Render(context.Background(), &buf)
		}
	})

	b.Run("ProgressBar render", func(b *testing.B) {
		props := ProgressBarProps{Current: 50, Total: 100, Label: "Upload"}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = ProgressBar(props).Render(context.Background(), &buf)
		}
	})

	b.Run("StepIndicator render", func(b *testing.B) {
		props := StepIndicatorProps{Steps: []string{"Details", "Review", "Confirm"}, CurrentStep: 1}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = StepIndicator(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Skeleton render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Skeleton(SkeletonText).Render(context.Background(), &buf)
		}
	})

	b.Run("Class merge", func(b *testing.B) {
		for b.Loop() {
			utils.Class("px-4 py-2 bg-red-500", "px-6 bg-blue-500 text-white")
		}
	})
}
