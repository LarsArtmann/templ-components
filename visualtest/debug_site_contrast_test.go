package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestDebugSiteContrast is a throwaway diagnostic for the dark-token probe.
func TestDebugSiteContrast(t *testing.T) {
	base := requireSiteDist(t)

	ctx, cancel := newTab(t)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var raw string

	if err := chromedp.Run(timeoutCtx,
		chromedp.EmulateViewport(1200, 900),
		chromedp.Navigate(base+"/sales"),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(`localStorage.setItem('theme','dark'); document.documentElement.classList.add('dark'); document.documentElement.style.colorScheme='dark'; true`, nil),
		chromedp.Sleep(settleDelay),
		chromedp.Evaluate(`(() => {
			const html = document.documentElement;
			const out = [];
			for (const el of document.querySelectorAll('h3.break-after-avoid, p.text-text-muted')) {
				const cs = getComputedStyle(el);
				let bgEl = el, chain = [];
				while (bgEl && bgEl !== document.body) {
					const bg = getComputedStyle(bgEl).backgroundColor;
					chain.push(bgEl.tagName + ':' + bg);
					bgEl = bgEl.parentElement;
				}
				out.push({cls: el.getAttribute('class').slice(0, 60), color: cs.color, opacity: cs.opacity, bgs: chain.slice(0, 4)});
			}
			return JSON.stringify({htmlDark: html.classList.contains('dark'), n: out.length, out: out.slice(0, 3)});
		})()`, &raw),
	); err != nil {
		t.Fatalf("debug probe: %v", err)
	}

	t.Logf("contrast report: %s", raw)
}
