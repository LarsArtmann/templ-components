package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Temporary debug test — delete after use.
func TestDebugWireFormValidation(t *testing.T) {
	srv := wireFormE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTimeout()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			if e.Request.URL != srv.URL+"/" {
				t.Logf("REQ %s %s", e.Request.Method, e.Request.URL)
			}
		case *network.EventResponseReceived:
			if e.Response.URL != srv.URL+"/" {
				t.Logf("RESP %d %s", e.Response.Status, e.Response.URL)
			}
		}
	})

	var (
		ok   bool
		dump string
	)

	step := func(label string, actions ...chromedp.Action) {
		t.Helper()

		if err := chromedp.Run(ctx, actions...); err != nil {
			t.Fatalf("%s: %v", label, err)
		}

		if err := chromedp.Run(ctx,
			chromedp.Evaluate(`location.pathname + ' | ' + (document.querySelector('#wire-form-htmx-region')?.innerText?.slice(0,120) ?? 'NO REGION')`, &dump),
		); err != nil {
			t.Fatalf("%s (dump): %v", label, err)
		}

		t.Logf("after %s: %s", label, dump)
	}

	setEmail := func(value string) {
		t.Helper()

		step("set email="+value, chromedp.Evaluate(
			`var i=document.querySelector('#wire-form-htmx-region input[name="email"]'); i.value=`+jsString(value)+
				`; i.dispatchEvent(new Event('input',{bubbles:true})); i.value`, &dump))
	}

	step("navigate", chromedp.Navigate(srv.URL+"/"))
	step("gate", chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &ok))
	step("fill name", chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="name"]`), "Ada Lovelace", chromedp.NodeVisible))
	setEmail("ada@example")
	step("submit #1", chromedp.Click(formSel(wireFormHTMXRegion, `button[type="submit"]`), chromedp.NodeVisible))
	step("poll summary", chromedp.Poll(regionHasText(wireFormHTMXRegion, "1 error found"), &ok))
	var attrs string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('#wire-form-htmx-region form').outerHTML.slice(0, 500)`, &attrs)); err != nil {
		t.Fatalf("attrs: %v", err)
	}
	t.Logf("SWAPPED FORM: %s", attrs)
	setEmail("ada@example.com")
	var btnType string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`window.__gotSubmit=false; window.__gotSubmitDefaultPrevented=null; var f=document.querySelector('#wire-form-htmx-region form'); f.addEventListener('submit', function(e){ window.__gotSubmit=true; window.__gotSubmitDefaultPrevented=e.defaultPrevented; }); var b=f.querySelector('button'); b ? b.type : 'NO BUTTON'`, &btnType)); err != nil {
		t.Fatalf("arm spy: %v", err)
	}
	t.Logf("button type: %s", btnType)
	step("click submit #2", chromedp.Click(formSel(wireFormHTMXRegion, `button[type="submit"]`), chromedp.NodeVisible))
	var spyRes string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`'gotSubmit=' + window.__gotSubmit + ' prevented=' + window.__gotSubmitDefaultPrevented`, &spyRes)); err != nil {
		t.Fatalf("read spy: %v", err)
	}
	t.Logf("SPY RESULT: %s", spyRes)
	var dispatched string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('#wire-form-htmx-region form').dispatchEvent(new Event('submit', {bubbles:true, cancelable:true})); 'dispatched'`, &dispatched)); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	t.Logf("dispatched: %s", dispatched)
	step("poll success", chromedp.Poll(regionHasText(wireFormHTMXRegion, "Subscribed Ada Lovelace (ada@example.com)"), &ok))
}
