package visualtest

import (
	"fmt"
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

	// Click "Go back" and poll until the browser is back on the demo index —
	// history.back() is async, so the URL check must poll.
	if err := chromedp.Run(ctx, chromedp.Click("[data-tc-go-back]", chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[errorpage]: click go-back: %v", err)
	}

	condition := "window.location.pathname === '/'"
	if err := chromedp.Run(ctx, pollTrue(condition, chromedp.WithPollingTimeout(3*time.Second))); err != nil {
		var url string
		_ = chromedp.Run(ctx, chromedp.Location(&url))

		t.Fatalf("visualtest[errorpage]: go-back never returned to / (stuck at %s): %v", url, err)
	}
}
