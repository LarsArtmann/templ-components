package visualtest

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestProbeBusyButtonSelector(t *testing.T) {
	server := StartDemoServer(t)
	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.BaseURL()+"/wire?transport=datastar"),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		t.Fatal(err)
	}

	var out string
	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`JSON.stringify({
			scripts: [...document.querySelectorAll('script[src]')].map(s => s.getAttribute('src')),
			inlineScripts: [...document.querySelectorAll('script:not([src])')].length,
		})`, &out)); err != nil {
		t.Fatal(err)
	}
	t.Log("PROBE:", out)
}
