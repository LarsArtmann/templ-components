// Package visualtest provides pixel-level visual regression testing for
// templ-components using a headless Chromium driven by chromedp.
//
// Unlike the HTML-string golden tests in internal/golden (which normalize CSS
// class order and compare text), these tests render each component in a real
// browser, capture a screenshot, and diff pixels against a committed golden
// PNG. This catches regressions the string tests cannot: layout shifts,
// dark-mode color breakage, RTL mirroring failures, and responsive collapse.
//
// # Running
//
// Visual tests require a Chromium binary. Provide it via CHROMEDP_CHROME_PATH
// (set automatically by `nix run .#visual`), or they skip gracefully:
//
//	nix run .#visual                      # run all visual tests
//	go test ./...                         # run with an explicit browser (from visualtest/)
//	go test ./... -update                 # regenerate golden PNGs
//
// The -update flag rewrites every golden with the current render, just like
// internal/golden. Review the diff before committing.
//
// # Why a separate module?
//
// chromedp and its transitive dependencies are heavy and test-only. Putting
// them in the library's go.mod would pollute every consumer's dependency
// graph. This package lives in its own module (./visualtest/go.mod) with a
// local replace directive, so `go get` of the main library never pulls them.
package visualtest

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/chromedp/chromedp"
)

var (
	allocatorOnce  sync.Once
	sharedAllocCtx context.Context
	allocCancel    context.CancelFunc
	browserReady   bool
)

// ensureAllocator launches a single shared Chromium process on first use.
// Subsequent calls reuse the same allocator — each test gets a new tab
// (chromedp.NewContext) which is lightweight (~10ms) compared to a full
// browser launch (~1s). The browser process lives for the entire test run
// and is cleaned up by TestMain.
func ensureAllocator(t *testing.T) {
	t.Helper()

	allocatorOnce.Do(func() {
		chromePath := os.Getenv("CHROMEDP_CHROME_PATH")
		if chromePath == "" {
			return
		}

		if _, err := os.Stat(chromePath); err != nil {
			return
		}

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(chromePath),
			chromedp.NoSandbox,
			chromedp.DisableGPU,
			chromedp.Flag("font-render-hinting", "none"),
			chromedp.Flag("disable-background-timer-throttling", true),
		)

		sharedAllocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
		browserReady = true
	})

	if !browserReady {
		chromePath := os.Getenv("CHROMEDP_CHROME_PATH")

		if chromePath == "" {
			t.Skipf("visual tests skipped: %v (set CHROMEDP_CHROME_PATH)", errNoBrowser)

			return
		}

		t.Skipf("visual tests skipped: CHROMEDP_CHROME_PATH %q not accessible", chromePath)
	}
}

// newTab returns a fresh tab context derived from the shared Chromium
// allocator. Each call creates a new browser tab (~10ms). The cancel func
// closes only the tab, not the browser process.
func newTab(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	ensureAllocator(t)

	ctx, cancel := chromedp.NewContext(sharedAllocCtx)

	return ctx, cancel
}

// ShutdownBrowser closes the shared Chromium process. Called by TestMain
// after all tests complete.
func ShutdownBrowser() {
	if allocCancel != nil {
		allocCancel()
	}
}
