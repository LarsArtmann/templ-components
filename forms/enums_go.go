package forms

// ToggleSizeIsValid reports whether v is one of the defined ToggleSize constants.
func ToggleSizeIsValid(v ToggleSize) bool {
	_, ok := toggleSizeMap[v]
	return ok
}

// InputTypeIsValid reports whether v is one of the defined InputType constants.
func InputTypeIsValid(v InputType) bool {
	return validInputTypes[v]
}

// FormMethodIsValid reports whether v is one of the defined FormMethod constants.
func FormMethodIsValid(v FormMethod) bool {
	return validFormMethods[v]
}
