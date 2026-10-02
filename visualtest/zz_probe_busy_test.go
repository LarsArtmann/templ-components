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
			t.Logf("CONSOLE %v", e)
		case *runtime.EventExceptionThrown:
			t.Logf("EXCEPTION: %v", e.ExceptionDetails)
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.BaseURL()+"/wire?transport=datastar"),
		chromedp.Sleep(4*time.Second),
		chromedp.Evaluate(`window.__ev=[];document.addEventListener('datastar-fetch',()=>window.__ev.push('fetch'),true);
			const b=[...document.querySelectorAll('button')].find(el => (el.getAttribute('data-on:click')||'').includes('/api/wire/busy'));
			b.__manualFired=false; b.addEventListener('click',()=>{b.__manualFired=true;});`, nil),
	); err != nil {
		t.Fatal(err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(`button[data-on\:click*="/api/wire/busy"]`, chromedp.ByQuery)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1500 * time.Millisecond)

	var out string
	_ = chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify({manual: (function(){const b=[...document.querySelectorAll('button')].find(el => (el.getAttribute('data-on:click')||'').includes('/api/wire/busy')); return b.__manualFired;})(), ev: window.__ev.length})`, &out))
	t.Log("AFTER TRUSTED CLICK:", out)
}
