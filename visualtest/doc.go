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
	"testing"

	"github.com/chromedp/chromedp"
)

// newBrowser launches a dedicated Chromium process for one test and returns a
// tab context plus its cancel func (which also stops that browser). If no
// browser binary is available the test is skipped. One process per test keeps
// failures isolated and avoids cross-test tab interference; startup is ~1s.
func newBrowser(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	chromePath := os.Getenv("CHROMEDP_CHROME_PATH")
	if chromePath == "" {
		t.Skipf("visual tests skipped: %v (set CHROMEDP_CHROME_PATH to a Chromium binary)", errNoBrowser)
	}
	if _, err := os.Stat(chromePath); err != nil {
		t.Skipf("visual tests skipped: CHROMEDP_CHROME_PATH %q not accessible: %v", chromePath, err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.NoSandbox, // required inside Nix builds / containers
		chromedp.DisableGPU,
		chromedp.Flag("font-render-hinting", "none"), // deterministic text rasterization
		chromedp.Flag("disable-background-timer-throttling", true),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Combine cancels so one defer tears down both the tab and the process.
	combinedCancel := func() {
		cancel()
		allocCancel()
	}
	return ctx, combinedCancel
}
