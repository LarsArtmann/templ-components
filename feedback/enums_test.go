package feedback

import "testing"

func TestIsValidEnums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{"SpinnerSize SM", func() bool { return SpinnerSizeIsValid(SpinnerSM) }, true},
		{"SpinnerSize invalid", func() bool { return SpinnerSizeIsValid(SpinnerSize("bogus")) }, false},
		{"ProgressBarSize MD", func() bool { return ProgressBarSizeIsValid(ProgressBarSizeMD) }, true},
		{"ProgressBarSize invalid", func() bool { return ProgressBarSizeIsValid(ProgressBarSize("bogus")) }, false},
		{"SkeletonVariant Text", func() bool { return SkeletonVariantIsValid(SkeletonText) }, true},
		{"SkeletonVariant invalid", func() bool { return SkeletonVariantIsValid(SkeletonVariant("bogus")) }, false},
		{"FeedbackType Success", func() bool { return FeedbackTypeIsValid(FeedbackSuccess) }, true},
		{"FeedbackType invalid", func() bool { return FeedbackTypeIsValid(FeedbackType("bogus")) }, false},
		{"StepOrientation Horizontal", func() bool {
			return StepIndicatorOrientationIsValid(StepHorizontal)
		}, true},
		{"StepOrientation invalid", func() bool {
			return StepIndicatorOrientationIsValid(StepIndicatorOrientation("bogus"))
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.fn(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
