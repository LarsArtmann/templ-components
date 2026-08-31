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

	options := resolveOptions(opts)

	// Build the page once (also validates CSS + component render early).
	page, err := renderHTML(component, options)
	if err != nil {
		t.Fatalf("visualtest[%s]: build page: %v", name, err)
	}

	ctx, cancel := newTab(t)
	defer cancel()

	actual, err := capture(ctx, page, options)
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

	maxMismatchPct := options.MaxMismatch * percentMultiplier

	result, diff := comparePixels(golden, actualImg, options.Threshold, maxMismatchPct)
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
// For *bool fields (Dark, RTL), a non-nil value from a later option overrides an
// earlier one; nil means "unset" so composing {Dark:Bool(true)} with
// {RTL:Bool(true)} yields dark+RTL. FullViewport ORs together.
func resolveOptions(opts []Options) Options {
	merged := Options{} //nolint:exhaustruct_v5 // zero value is the intended "unset" sentinel; merged below

	for _, opt := range opts {
		if opt.Dark != nil {
			merged.Dark = opt.Dark
		}

		if opt.RTL != nil {
			merged.RTL = opt.RTL
		}

		if opt.Viewport.Width != 0 {
			merged.Viewport.Width = opt.Viewport.Width
		}

		if opt.Viewport.Height != 0 {
			merged.Viewport.Height = opt.Viewport.Height
		}

		if opt.MaxMismatch != 0 {
			merged.MaxMismatch = opt.MaxMismatch
		}

		if opt.Threshold != 0 {
			merged.Threshold = opt.Threshold
		}

		if opt.State != StateRest {
			merged.State = opt.State
		}

		merged.FullViewport = merged.FullViewport || opt.FullViewport

		if opt.WaitSelector != "" {
			merged.WaitSelector = opt.WaitSelector
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
	// :hover / :focus-visible styles are baked into the screenshot. StateRest
	// is intentionally a no-op (the default — no interaction needed).
	switch opts.State {
	case StateRest: // no-op: the default, nothing to dispatch
	case StateHover:
		tasks = append(tasks, hoverAction(rootSel))
	case StateFocus:
		tasks = append(tasks, focusAction(rootSel))
	case StateClick:
		tasks = append(tasks, clickAction(rootSel))
	case StateContext:
		tasks = append(tasks, contextAction(rootSel))
	}

	// After the interaction, optionally wait for a selector to become visible
	// — e.g. an overlay menu opened by StateClick — then settle and capture.
	if opts.WaitSelector != "" {
		tasks = append(tasks, chromedp.WaitVisible(opts.WaitSelector, chromedp.ByQuery))
		// Wait for CSS @starting-style transitions to finish so the overlay is
		// captured in its settled state, not mid-slide. Without this, parallel
		// test load can delay the transition start past the fixed settleDelay,
		// capturing the drawer off-screen (a ~90% false mismatch).
		tasks = append(tasks, waitAnimationSettled(opts.WaitSelector))
	}

	var capture chromedp.Action = chromedp.Screenshot(rootSel, &screenshot, chromedp.ByQuery, chromedp.NodeVisible)
	if opts.FullViewport {
		// Full-viewport capture: top-layer overlays (Popover API menus,
		// <dialog>) paint outside #tc-root's box, so an element screenshot
		// would crop them.
		capture = chromedp.CaptureScreenshot(&screenshot)
	}

	tasks = append(
		tasks,
		chromedp.Sleep(settleDelay),
		capture,
	)
	if err := chromedp.Run(timeoutCtx, tasks...); err != nil {
		return nil, fmt.Errorf("chromedp actions: %w", err)
	}

	return screenshot, nil
}

// waitAnimationSettled blocks until all CSS animations/transitions on the first
// element matching sel have finished, with a best-effort timeout. It uses a
// two-phase approach: (1) wait up to animRegisterDelay for @starting-style
// transitions to register (getAnimations() is transiently empty before the
// browser creates the transition), then (2) poll until all animations report
// "finished". If no animations appear within the registration window, the
// function returns immediately (the element genuinely has no transition).
func waitAnimationSettled(sel string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if err := chromedp.Sleep(animInitialDelay).Do(ctx); err != nil {
			return fmt.Errorf("wait animation settled: initial sleep: %w", err)
		}

		expr := fmt.Sprintf(
			`(()=>{const el=document.querySelector(%q);if(!el)return"none";`+
				`const a=el.getAnimations();if(a.length===0)return"empty";`+
				`return a.every(x=>x.playState==='finished')?"done":"running";})()`, sel,
		)

		registerDeadline := time.Now().Add(animRegisterDelay)
		settleDeadline := time.Now().Add(animSettleDelay)
		animSeen := false

		for time.Now().Before(settleDeadline) {
			var state string

			if err := chromedp.Evaluate(expr, &state).Do(ctx); err != nil {
				return fmt.Errorf("wait animation settled: evaluate: %w", err)
			}

			switch state {
			case "none", "done":
				return nil
			case "running":
				animSeen = true
			case "empty":
				if animSeen || !time.Now().Before(registerDeadline) {
					return nil
				}
			}

			if err := chromedp.Sleep(animPollDelay).Do(ctx); err != nil {
				return fmt.Errorf("wait animation settled: poll sleep: %w", err)
			}
		}

		return nil
	})
}

const (
	captureTimeout    = 20 * time.Second
	settleDelay       = 200 * time.Millisecond
	animInitialDelay  = 80 * time.Millisecond
	animRegisterDelay = 300 * time.Millisecond
	animSettleDelay   = 800 * time.Millisecond
	animPollDelay     = 40 * time.Millisecond
)

// hoverAction returns a chromedp action that moves the mouse to the center of
// the first interactive descendant of the element matching sel (button, link,
// etc.), falling back to the root element itself if none is found. This mirrors
// focusAction: the wrapper (#tc-root) is a plain <div>, so hovering its centre
// does not trigger :hover styles on the inner button/link that the test intends
// to capture. A real mouse-move event is dispatched (synthetic mouseover events
// do not trigger :hover) at the element's centre so coordinates are correct
// regardless of scroll position.
func hoverAction(sel string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Descend to the first interactive child like focusAction; fall back to
		// the root so components with no inner interactive element still get a
		// hover at the root centre.
		script := `(() => {
			const root = document.querySelector(` + fmt.Sprintf("%q", sel) + `);
			if (!root) return null;
			const e = root.querySelector('button, a[href], input, select, textarea, [role="button"], [tabindex]:not([tabindex="-1"])') || root;
			const r = e.getBoundingClientRect();
			return [r.x + r.width/2, r.y + r.height/2];
		})()`

		var coords []float64
		if err := chromedp.Evaluate(script, &coords).Do(ctx); err != nil {
			return fmt.Errorf("hover: get %s rect: %w", sel, err)
		}

		if len(coords) != 2 {
			return wrapSelector(errHoverElementNotFound, sel)
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
		script := `(() => {
			const root = document.querySelector(` + fmt.Sprintf("%q", sel) + `);
			if (!root) return false;
			const f = root.querySelector('button, a[href], input, select, textarea, [tabindex]:not([tabindex="-1"])');
			if (f) { f.focus(); return true; }
			return false;
		})()`

		var focused bool
		if err := chromedp.Evaluate(script, &focused).Do(ctx); err != nil {
			return fmt.Errorf("focus: query %s: %w", sel, err)
		}

		if !focused {
			return wrapSelector(errFocusNoElement, sel)
		}

		return nil
	})
}

// clickAction clicks the first interactive descendant of the element matching
// sel. It prefers a Popover API trigger ([popovertarget]) so Dropdown/Popover/
// ContextMenu open natively, then falls back to any button/link. A real mouse
// click is dispatched at the element's centre so the native :active state and
// popovertarget invoker both fire (synthetic .click() is enough for the
// invoker, but a dispatched event also exercises hover/active paint).
func clickAction(sel string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		script := `(() => {
			const root = document.querySelector(` + fmt.Sprintf("%q", sel) + `);
			if (!root) return null;
			const e = root.querySelector('[popovertarget], button, a[href], [role="button"]');
			if (!e) return null;
			const r = e.getBoundingClientRect();
			return [r.x + r.width/2, r.y + r.height/2];
		})()`

		var coords []float64
		if err := chromedp.Evaluate(script, &coords).Do(ctx); err != nil {
			return fmt.Errorf("click: query %s: %w", sel, err)
		}

		if len(coords) != 2 {
			return wrapSelector(errClickNoElement, sel)
		}

		// Press + release at centre: this triggers the popovertarget invoker
		// (opening the menu) and a real :active paint cycle.
		if err := input.DispatchMouseEvent(input.MousePressed, coords[0], coords[1]).
			WithButton(input.Left).WithClickCount(1).Do(ctx); err != nil {
			return fmt.Errorf("click: press: %w", err)
		}

		return input.DispatchMouseEvent(input.MouseReleased, coords[0], coords[1]).
			WithButton(input.Left).WithClickCount(1).Do(ctx)
	})
}

// contextAction dispatches a contextmenu (right-click) event on the first
// ContextMenu trigger (or any interactive descendant) under sel. A synthetic
// MouseEvent('contextmenu') is used because it reliably fires the singleton
// handler that calls showPopover(); a real right-button press is flaky under
// headless Chromium.
func contextAction(sel string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		script := `(() => {
			const root = document.querySelector(` + fmt.Sprintf("%q", sel) + `);
			if (!root) return false;
			const t = root.querySelector('[data-tc-ctxmenu-trigger], button, a[href], [tabindex]');
			if (!t) return false;
			t.dispatchEvent(new MouseEvent('contextmenu', {bubbles: true, cancelable: true}));
			return true;
		})()`

		var ok bool
		if err := chromedp.Evaluate(script, &ok).Do(ctx); err != nil {
			return fmt.Errorf("context: query %s: %w", sel, err)
		}

		if !ok {
			return wrapSelector(errContextNoTrigger, sel)
		}

		return nil
	})
}
