package datastar

import "github.com/larsartmann/go-datastar/static"

// DatastarVersionIsValid reports whether v is one of the defined DatastarVersion constants.
// Delegates to the upstream-owned static.VersionIsValid (single source of
// truth since go-datastar/static v0.8.0).
func DatastarVersionIsValid(v DatastarVersion) bool {
	return static.VersionIsValid(string(v))
}

// LivePolitenessIsValid reports whether v is one of the defined LivePoliteness constants.
func LivePolitenessIsValid(v LivePoliteness) bool {
	return validLivePoliteness[v]
}

// RetryModeIsValid reports whether v is one of the defined RetryMode constants.
func RetryModeIsValid(v RetryMode) bool {
	return validRetryModes[v]
}
