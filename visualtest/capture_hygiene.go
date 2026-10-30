package visualtest

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
)

// Route captures must start from NEUTRAL pointer/scroll/keyboard state. A
// fresh tab does not guarantee it: the CDP virtual pointer keeps its position
// across navigations, so a capture following a hover test can bake `:hover`
// styles (open tooltips, shifted cards) into a "rest state" golden. The
// 2026-10-30 icons_light stray capture (sidebar brand text hovering at the
// Display row) is the artifact class this kills.
//
// resetCaptureNeutralState runs EARLY in the chain (pointer home + blur +
// scroll top, before the settle waits, so reveal/scroll effects settle from a
// known origin); assertCaptureNeutralState runs LAST and FAILS the capture if
// the state drifted back — an assert (not a silent reset) so the artifact
// class surfaces instead of being papered over per run.

// resetCaptureNeutralState parks the virtual pointer at the origin, blurs any
// focused element, and scrolls the page to the top.
func resetCaptureNeutralState() chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		if _, err := cdp.Call(ctx, tab, input.DispatchMouseEvent, input.DispatchMouseEventParams{ //nolint:exhaustruct_v5 // omitted CDP fields take protocol defaults
			Type: input.DispatchMouseEventTypeMouseMoved,
			X:    0,
			Y:    0,
		}); err != nil {
			return fmt.Errorf("park pointer: %w", err)
		}

		return nil
	})
}

// neutralStateJS blurs focus and returns the scroll position to the origin.
const neutralStateJS = `(() => {
	if (document.activeElement && document.activeElement !== document.body) {
		document.activeElement.blur();
	}
	window.scrollTo(0, 0);
})()`

// neutralStatePredicate asserts the page is capture-ready: no scroll, no focus.
const neutralStatePredicate = `window.scrollY === 0 &&
	(!document.activeElement || document.activeElement === document.body || document.activeElement === document.documentElement)`

// neutralCaptureActions resets focus/scroll and then ASSERTS the page is
// capture-ready (no scroll, no focused element). Append these at the END of a
// capture chain, after the settle waits.
func neutralCaptureActions() []chromedp.Action[chromedp.Void] {
	return []chromedp.Action[chromedp.Void]{
		evalVoid(neutralStateJS),
		pollTrue(neutralStatePredicate),
	}
}
