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
//	go test ./... -tags=visual            # run with an explicit browser
//	go test ./... -tags=visual -update    # regenerate golden PNGs
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
	"time"

	"github.com/chromedp/chromedp"
)

// browserAllocator is shared across all tests in a run: one Chromium process,
// many tabs (one per test). Created lazily so tests without a browser skip
// instead of failing at init.
var (
	browserOnce sync.Once
	allocCtx    context.Context
	allocCancel context.CancelFunc
	browserErr  error
)

// newBrowser prepares the shared Chromium allocator and returns a fresh tab
// context plus its cancel func. The test's cleanup is the caller's
// responsibility. If no browser is available the test is skipped.
func newBrowser(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	browserOnce.Do(func() {
		chromePath := os.Getenv("CHROMEDP_CHROME_PATH")
		if chromePath == "" {
			browserErr = errNoBrowser
			return
		}
		if _, err := os.Stat(chromePath); err != nil {
			browserErr = errNoBrowser
			return
		}

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(chromePath),
			chromedp.NoSandbox, // required inside Nix builds / containers
			chromedp.DisableGPU,
			chromedp.Flag("font-render-hinting", "none"), // deterministic text rasterization
			chromedp.Flag("disable-background-timer-throttling", true),
		)
		allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	})

	if browserErr != nil {
		t.Skipf("visual tests skipped: %v (set CHROMEDP_CHROME_PATH to a Chromium binary)", browserErr)
	}

	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	// Warm up: ensure the first tab is ready before returning so navigation
	// actions do not race the browser startup.
	timeoutCtx, timeoutCancel := context.WithTimeout(tabCtx, browserStartupTimeout)
	if err := chromedp.Run(timeoutCtx); err != nil {
		tabCancel()
		timeoutCancel()
		t.Fatalf("visualtest: start browser tab: %v", err)
	}
	timeoutCancel()

	return tabCtx, tabCancel
}

const browserStartupTimeout = 30 * time.Second
