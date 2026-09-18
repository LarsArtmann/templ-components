package visualtest

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestDemoErrorPageGoBack drives the real browser through the error pages'
// "Go back" affordance: the shared goBackScript singleton must call
// history.back() so the user lands on the page they came from — the exact
// behavior string assertions can only guess at. NotFound404 (route
// /errors/404-page, ShowGoBack default true) exercises the same singleton
// the ErrorPage WayOut button uses.
func TestDemoErrorPageGoBack(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	base := server.BaseURL()

	// Land on page A, then push an error route on top of the history stack.
	if err := chromedp.Run(ctx,
		chromedp.Navigate(base+"/"),
		chromedp.WaitReady("body"),
		chromedp.Navigate(base+"/errors/404-page"),
		chromedp.WaitVisible("[data-tc-go-back]", chromedp.ByQuery),
	); err != nil {
		t.Fatalf("visualtest[errorpage]: navigate to /errors/404-page: %v", err)
	}

	// Click "Go back" and wait until the browser is back on the demo index.
	// history.back() is async AND the navigation invalidates the JS execution
	// context the click's poll would run in — a single chromedp.Poll here
	// races the navigation and dies with "Cannot find context". Retry the
	// evaluation in a deadline loop instead (the demoClickUntil pattern).
	if err := chromedp.Run(ctx, chromedp.Click("[data-tc-go-back]", chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[errorpage]: click go-back: %v", err)
	}

	deadline := time.Now().Add(demoFlowTimeout)
	for {
		var back bool
		evalErr := chromedp.Run(ctx, chromedp.Evaluate(`window.location.pathname === '/'`, &back))
		if evalErr == nil && back {
			break
		}

		if !time.Now().Before(deadline) {
			var url string
			_ = chromedp.Run(ctx, chromedp.Location(&url))

			t.Fatalf("visualtest[errorpage]: go-back never returned to / (stuck at %s)", url)
		}

		time.Sleep(300 * time.Millisecond)
	}
}
