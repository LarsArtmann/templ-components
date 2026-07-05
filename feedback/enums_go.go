package feedback

// SpinnerSizeIsValid reports whether v is one of the defined SpinnerSize constants.
func SpinnerSizeIsValid(v SpinnerSize) bool {
	_, ok := spinnerSizeLookup[string(v)]
	return ok
}

// ProgressBarSizeIsValid reports whether v is one of the defined ProgressBarSize constants.
func ProgressBarSizeIsValid(v ProgressBarSize) bool {
	_, ok := progressHeightLookup[string(v)]
	return ok
}
