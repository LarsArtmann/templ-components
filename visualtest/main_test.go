package visualtest

import (
	"os"
	"testing"
)

// TestMain tears the shared Chromium process down after all tests complete.
// The browser is lazily initialized by ensureAllocator on the first
// AssertScreenshot call; this teardown runs regardless of whether any visual
// test actually ran.
//
// Failure evidence lives in per-run directories (testdata/.fail/<UTC
// timestamp>-<pid>/, see failureRunDir) and is retained for the last
// failureRunDirsKept runs, so a later run can never destroy an earlier one's
// debugging evidence (the pre-2026-10-30 start-of-run wipe did exactly that,
// twice in one session).
func TestMain(m *testing.M) {
	code := m.Run()

	ShutdownBrowser()

	os.Exit(code)
}
