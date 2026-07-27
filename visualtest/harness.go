package visualtest

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/input"
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

	ctx, cancel := newTab(t)
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

	result, diff := comparePixels(golden, actualImg, o.Threshold, maxMismatchPct)
	if !result.Match {
		writeFailureArtifacts(t, name, actual, diff)
		t.Errorf("visualtest[%s]: visual mismatch — %s (max %.4f%%).\n"+
			"Inspect testdata/.fail/%s.{actual,diff}.png, then run `go test -update` if the change is intended.",
			name, result, maxMismatchPct, name)

		return
	}

	cleanFailureArtifacts(name)
	t.Logf("visualtest[%s]: OK (%.4f%% mismatched, threshold %.4f%%)", name, result.MismatchPct, maxMismatchPct)
}

// resolveOptions merges the variadic option list into one Options with defaults.
// Later options win for non-bool fields; bools OR together so composing
// {Dark:true} with {RTL:true} yields dark+RTL.
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

		if o.Threshold != 0 {
			merged.Threshold = o.Threshold
		}

		if o.State != StateRest {
			merged.State = o.State
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
	}

	// Apply the requested interaction state before settling + capture so
	// :hover / :focus-visible styles are baked into the screenshot.
	switch opts.State {
	case StateHover:
		tasks = append(tasks, hoverAction(rootSel))
	case StateFocus:
		tasks = append(tasks, focusAction(rootSel))
	}

	tasks = append(tasks,
		chromedp.Sleep(settleDelay),
		chromedp.Screenshot(rootSel, &screenshot, chromedp.ByQuery, chromedp.NodeVisible),
	)
	if err := chromedp.Run(timeoutCtx, tasks...); err != nil {
		return nil, fmt.Errorf("chromedp actions: %w", err)
	}

	return screenshot, nil
}

const (
	captureTimeout = 20 * time.Second
	settleDelay    = 200 * time.Millisecond
)

// hoverAction returns a chromedp action that moves the mouse to the center of
// the element matching sel, triggering real :hover styles (synthetic
// mouseover events do not). The element's bounding box is read via JS so the
// coordinates are correct regardless of scroll position.
func hoverAction(sel string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// sel is interpolated as a quoted JS string literal; querySelector needs
		// no arguments object (arrow functions don't bind `arguments`).
		js := `(() => {
			const e = document.querySelector(` + fmt.Sprintf("%q", sel) + `);
			if (!e) return null;
			const r = e.getBoundingClientRect();
			return [r.x + r.width/2, r.y + r.height/2];
		})()`

		var coords []float64
		if err := chromedp.Evaluate(js, &coords).Do(ctx); err != nil {
			return fmt.Errorf("hover: get %s rect: %w", sel, err)
		}

		if len(coords) != 2 {
			return fmt.Errorf("hover: element %q not found", sel)
		}

		return input.DispatchMouseEvent(input.MouseMoved, coords[0], coords[1]).Do(ctx)
	})
}

// focusAction focuses the first focusable descendant of the element matching
// sel. The wrapper (#tc-root) is a <div> and not itself focusable, so we
// descend to the interactive element (button/a/input) it wraps. No-op if no
// focusable element exists.
func focusAction(sel string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		js := `(() => {
			const root = document.querySelector(` + fmt.Sprintf("%q", sel) + `);
			if (!root) return false;
			const f = root.querySelector('button, a[href], input, select, textarea, [tabindex]:not([tabindex="-1"])');
			if (f) { f.focus(); return true; }
			return false;
		})()`

		var focused bool
		if err := chromedp.Evaluate(js, &focused).Do(ctx); err != nil {
			return fmt.Errorf("focus: query %s: %w", sel, err)
		}

		if !focused {
			return fmt.Errorf("focus: no focusable element under %q", sel)
		}

		return nil
	})
}
