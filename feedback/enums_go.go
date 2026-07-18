package feedback

// SpinnerSizeIsValid reports whether v is one of the defined SpinnerSize constants.
func SpinnerSizeIsValid(v SpinnerSize) bool {
	_, ok := spinnerSizeLookup[v]

	return ok
}

// ProgressBarSizeIsValid reports whether v is one of the defined ProgressBarSize constants.
func ProgressBarSizeIsValid(v ProgressBarSize) bool {
	_, ok := progressHeightLookup[v]

	return ok
}

// SkeletonVariantIsValid reports whether v is one of the defined SkeletonVariant constants.
func SkeletonVariantIsValid(v SkeletonVariant) bool {
	switch v {
	case SkeletonText, SkeletonTextShort, SkeletonTitle, SkeletonAvatar, SkeletonImage, SkeletonCard, SkeletonTableRow:
		return true
	default:
		return false
	}
}

// StepIndicatorOrientationIsValid reports whether v is one of the defined StepIndicatorOrientation constants.
func StepIndicatorOrientationIsValid(v StepIndicatorOrientation) bool {
	switch v {
	case StepHorizontal, StepVertical:
		return true
	default:
		return false
	}
}
