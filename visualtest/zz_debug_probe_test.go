package visualtest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestDatastarPolledDebugProbe is a throwaway diagnostic: loads the polled
// e2e page, captures console.error / error events / unhandled rejections
// into window.__errs, and dumps the region's outerHTML. Delete after
// root-causing.
func TestDatastarPolledDebugProbe(t *testing.T) {
	t.Parallel()

	srv, ticks := datastarPolledRegionE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var (
		ready     bool
		outerHTML string
		scriptOK  bool
		errs      string
		tickText  string
	)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Evaluate(`(function(){
			window.__errs = [];
			const origErr = console.error;
			console.error = function(){ window.__errs.push(Array.from(arguments).join(' ')); origErr.apply(console, arguments); };
			window.addEventListener('error', function(e){ window.__errs.push('ERR: ' + e.message); });
			window.addEventListener('unhandledrejection', function(e){ window.__errs.push('REJ: ' + ((e.reason && e.reason.message) || String(e.reason))); });
			return true;
		})()`, nil),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`window.__dsReady===true`, &ready),
		chromedp.Evaluate(`!!document.querySelector('script[src="/datastar.js"]')`, &scriptOK),
		chromedp.Evaluate(`document.getElementById('tick-value') ? document.getElementById('tick-value').textContent : 'MISSING'`, &tickText),
		chromedp.Evaluate(`document.getElementById('polled-region') ? document.getElementById('polled-region').outerHTML : 'MISSING'`, &outerHTML),
		chromedp.Evaluate(`(window.__errs || []).join(' || ') || 'no errors captured'`, &errs),
	); err != nil {
		t.Fatalf("probe: %v", err)
	}

	t.Logf("ready=%v scriptOK=%v ticks=%d tickText=%q", ready, scriptOK, ticks.Load(), tickText)
	t.Logf("errors=%s", strings.Clone(errs))
	t.Logf("regionHTML=%s", strings.Clone(outerHTML))
}
