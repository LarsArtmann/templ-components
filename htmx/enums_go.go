package htmx

// SwapStyleIsValid reports whether v is one of the defined SwapStyle constants.
func SwapStyleIsValid(v SwapStyle) bool {
	return validSwapStyles[v]
}
