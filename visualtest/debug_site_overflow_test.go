package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestDebugSiteOverflow is a throwaway diagnostic: lists every element whose
// layout rect extends past the 320px viewport, excluding ink-only overflow
// (elements inside overflow-clipped ancestors).
func TestDebugSiteOverflow(t *testing.T) {
	base := requireSiteDist(t)

	ctx, cancel := newTab(t)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var raw string

	if err := chromedp.Run(timeoutCtx,
		chromedp.EmulateViewport(320, 900),
		chromedp.Navigate(base+"/index.html"),
		chromedp.WaitReady("body"),
		chromedp.Sleep(settleDelay),
		chromedp.Evaluate(`(() => {
			const cw = document.documentElement.clientWidth;
			const over = document.documentElement.scrollWidth - cw;
			const out = [];
			for (const el of document.querySelectorAll('body *')) {
				const r = el.getBoundingClientRect();
				if (r.right <= cw + 1) continue;
				// skip if an ancestor clips overflow
				let clipped = false;
				for (let a = el.parentElement; a; a = a.parentElement) {
					const s = getComputedStyle(a);
					if (s.overflowX === 'hidden' || s.overflowX === 'clip' || s.overflowX === 'auto' || s.overflowX === 'scroll') { clipped = true; break; }
				}
				if (clipped) continue;
				out.push({tag: el.tagName, id: el.id, cls: (el.getAttribute('class')||'').slice(0,60), right: Math.round(r.right), w: Math.round(r.width)});
			}
			return JSON.stringify({clientWidth: cw, scrollWidth: document.documentElement.scrollWidth, over, out: out.slice(0, 12)});
		})()`, &raw),
	); err != nil {
		t.Fatalf("debug probe: %v", err)
	}

	t.Logf("overflow report: %s", raw)
}
