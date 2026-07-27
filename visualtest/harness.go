package visualtest

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
)

// AssertScreenshot renders component in an isolated page, captures an element
// screenshot of #tc-root, and compares the pixels against the golden PNG in
// testdata. With -update, the golden is rewritten instead.
//
//	name uniquely identifies the golden file (e.g. "button/primary_dark"). It
//	must not contain path separators other than a single subdirectory level.
func AssertScreenshot(t *testing.T, name string, component templ.Component, opts ...Options) {
	t.Helper()
	o := resolveOptions(opts)

	// Build the page once (also validates CSS + component render early).
	page, err := renderHTML(component, o)
	if err != nil {
		t.Fatalf("visualtest[%s]: build page: %v", name, err)
	}

	ctx, cancel := newBrowser(t)
	defer cancel()

	actual, err := capture(ctx, page, o)
	if err != nil {
		t.Fatalf("visualtest[%s]: capture: %v", name, err)
	}

	if *update {
		writeGolden(t, name, actual)
		return
	}

	golden, exists := readGolden(t, name)
	if !exists {
		writeGolden(t, name, actual)
		t.Errorf("visualtest[%s]: no golden yet — wrote %s (re-run without -update to verify)", name, goldenPath(name))
		return
	}

	actualImg, err := png.Decode(bytes.NewReader(actual))
	if err != nil {
		t.Fatalf("visualtest[%s]: decode actual: %v", name, err)
	}

	maxMismatchPct := o.MaxMismatch * 100
	result, diff := comparePixels(golden, actualImg, o.PerPixelTolerance, maxMismatchPct)
	if !result.Match {
		writeFailureArtifacts(t, name, actual, diff)
		t.Errorf("visualtest[%s]: visual mismatch — %dx%d, %.4f%% pixels differ (max %.4f%%).\n"+
			"Inspect testdata/.fail/%s.{actual,diff}.png, then run `go test -tags=visual -update` if the change is intended.",
			name, result.Width, result.Height, result.MismatchPct, maxMismatchPct, name)
		return
	}
	t.Logf("visualtest[%s]: OK (%.4f%% mismatched, threshold %.4f%%)", name, result.MismatchPct, maxMismatchPct)
}

// resolveOptions merges the variadic option list into one Options with defaults.
func resolveOptions(opts []Options) Options {
	merged := Options{}
	for _, o := range opts {
		merged.Dark = merged.Dark || o.Dark
		merged.RTL = merged.RTL || o.RTL
		if o.Viewport.Width != 0 {
			merged.Viewport.Width = o.Viewport.Width
		}
		if o.Viewport.Height != 0 {
			merged.Viewport.Height = o.Viewport.Height
		}
		if o.MaxMismatch != 0 {
			merged.MaxMismatch = o.MaxMismatch
		}
		if o.PerPixelTolerance != 0 {
			merged.PerPixelTolerance = o.PerPixelTolerance
		}
	}
	return defaultOptions(merged)
}

// capture serves the page on an ephemeral in-process HTTP server (a data:
// URL cannot hold the ~80KB compiled CSS), waits for layout to settle, and
// screenshots the #tc-root element.
func capture(ctx context.Context, page string, opts Options) ([]byte, error) {
	const rootSel = "#tc-root"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, page)
	}))
	defer srv.Close()

	timeoutCtx, cancel := context.WithTimeout(ctx, captureTimeout)
	defer cancel()

	var screenshot []byte
	tasks := []chromedp.Action{
		chromedp.EmulateViewport(int64(opts.Viewport.Width), int64(opts.Viewport.Height)),
		chromedp.Navigate(srv.URL),
		// Wait for the root to exist and the document to finish loading so web
		// fonts / CSS settle before the screenshot.
		chromedp.WaitVisible(rootSel, chromedp.ByQuery),
		chromedp.Sleep(settleDelay),
		chromedp.Screenshot(rootSel, &screenshot, chromedp.ByQuery, chromedp.NodeVisible),
	}
	if err := chromedp.Run(timeoutCtx, tasks...); err != nil {
		return nil, fmt.Errorf("chromedp actions: %w", err)
	}
	return screenshot, nil
}

const (
	captureTimeout = 20 * time.Second
	settleDelay    = 200 * time.Millisecond
)

// decodePNG is a small helper for tests in this package that need to inspect a
// captured image directly.
func decodePNG(b []byte) (image.Image, error) { //nolint:unused // helper for future test introspection
	return png.Decode(bytes.NewReader(b))
}
