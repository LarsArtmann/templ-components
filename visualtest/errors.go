package visualtest

import "errors"

// Sentinel errors for the visualtest package.
var (
	// errNoBrowser is returned when no Chromium binary is available. Tests
	// that encounter it skip rather than fail.
	errNoBrowser = errors.New("no Chromium binary found")
)
