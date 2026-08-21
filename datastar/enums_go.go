package datastar

// DatastarVersionIsValid reports whether v is one of the defined DatastarVersion constants.
func DatastarVersionIsValid(v DatastarVersion) bool {
	switch v {
	case DatastarVersion1_0_2:
		return true
	default:
		return false
	}
}

// LivePolitenessIsValid reports whether v is one of the defined LivePoliteness constants.
func LivePolitenessIsValid(v LivePoliteness) bool {
	return validLivePoliteness[v]
}

// RetryModeIsValid reports whether v is one of the defined RetryMode constants.
func RetryModeIsValid(v RetryMode) bool {
	return validRetryModes[v]
}
