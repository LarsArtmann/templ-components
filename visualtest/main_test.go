package visualtest

import (
	"os"
	"testing"
)

// TestMain ensures the shared Chromium process is cleaned up after all tests
// complete. The browser is lazily initialized by ensureAllocator on the first
// AssertScreenshot call; this teardown runs regardless of whether any visual
// test actually ran.
func TestMain(m *testing.M) {
	code := m.Run()

	ShutdownBrowser()

	os.Exit(code)
}
