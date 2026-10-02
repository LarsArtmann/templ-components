package visualtest

import (
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func TestProbeBusyButtonSelector(t *testing.T) {
	server := StartDemoServer(t)
	ctx, cancel := newFlowTab(t)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			for _, a := range e.Args {
				if a.Description != "" {
					t.Logf("CONSOLE[%s]: %s", e.Type, a.Description)
				}
			}
		case *runtime.EventExceptionThrown:
			if e.ExceptionDetails != nil && e.ExceptionDetails.Exception != nil {
				t.Logf("EXCEPTION: %s | %v", e.ExceptionDetails.Exception.Description, e.ExceptionDetails.StackTrace)
			}
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.BaseURL()+"/wire?transport=datastar"),
		chromedp.Sleep(4*time.Second),
		chromedp.Evaluate(`window.__fetches=[]; const of=window.fetch; window.fetch=function(...a){window.__fetches.push(String(a[0])); return of.apply(this,a);};`, nil),
	); err != nil {
		t.Fatal(err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(`button[data-on\:click*="/api/wire/busy"]`, chromedp.ByQuery)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2500 * time.Millisecond)

	var out string
	_ = chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify({fetches: window.__fetches, region: (document.getElementById('wire-busy-datastar-out')||{}).innerText.slice(0,25)})`, &out))
	t.Log("AFTER TRUSTED CLICK:", out)
}
