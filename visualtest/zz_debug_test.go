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
	var trace string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`window.__events=[]; ['htmx:afterSwap','htmx:afterSettle','htmx:load','htmx:afterProcessNode','htmx:configRequest','htmx:beforeSwap','htmx:error'].forEach(function(n){ document.addEventListener(n, function(e){ window.__events.push(n + (e.target && e.target.id ? ':'+e.target.id : '')); }); }); 'armed'`, &trace)); err != nil {
		t.Fatalf("arm trace: %v", err)
	}
	step("arm errors", chromedp.Evaluate(`window.__errs=[]; window.onerror=function(m,s,l,c){ window.__errs.push('onerror:'+m+' @'+l+':'+c); }; window.addEventListener('unhandledrejection', function(e){ window.__errs.push('rejection:'+String(e.reason)); }); 'armed-errs'`, &trace))
	step("gate", chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &ok))
	step("fill name", chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="name"]`), "Ada Lovelace", chromedp.NodeVisible))
	setEmail("ada@example")
	step("submit #1", chromedp.Click(formSel(wireFormHTMXRegion, `button[type="submit"]`), chromedp.NodeVisible))
	step("poll summary", chromedp.Poll(regionHasText(wireFormHTMXRegion, "1 error found"), &ok))
	var (
		attrs string
		ident string
	)
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`window.__liveForm = document.querySelector('#wire-form-htmx-region form'); 'tagged'`, &ident),
	); err != nil {
		t.Fatalf("tag: %v", err)
	}
	if err := chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify({events: window.__events, formHTML: document.querySelector('#wire-form-htmx-region form').outerHTML.slice(0,300)})`, &attrs)); err != nil {
		t.Fatalf("attrs: %v", err)
	}
	t.Logf("AFTER SWAP: %s", attrs)
	setEmail("ada@example.com")
	var (
		probe     string
		timerFired bool
	)
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`window.__timerFired=false; setTimeout(function(){ window.__timerFired=true; }, 50); JSON.stringify({settleDelay: htmx.config.defaultSettleDelay, globalViewTransitions: htmx.config.globalViewTransitions, defaultSwapStyle: htmx.config.defaultSwapStyle, htmxVersion: htmx.version})`, &probe),
		chromedp.Sleep(300*time.Millisecond),
		chromedp.Evaluate(`window.__timerFired`, &timerFired),
	); err != nil {
		t.Fatalf("probe: %v", err)
	}
	t.Logf("TIMER PROBE: fired=%v config=%s", timerFired, probe)
	var preTrace string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify({events: window.__events, sameNodeAsSwapped: window.__liveForm === document.querySelector('#wire-form-htmx-region form')})`, &preTrace)); err != nil {
		t.Fatalf("pretrace: %v", err)
	}
	t.Logf("PRE-SUBMIT2: %s", preTrace)
	var errs string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify(window.__errs)`, &errs)); err != nil {
		t.Fatalf("errs: %v", err)
	}
	t.Logf("PAGE ERRORS: %s", errs)
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
	var (
		dispatched string
		finalState string
	)
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`window.__gotSubmit=false; window.__gotSubmitDefaultPrevented=null; var f=document.querySelector('#wire-form-htmx-region form'); f.addEventListener('submit', function(e){ window.__gotSubmit=true; window.__gotSubmitDefaultPrevented=e.defaultPrevented; }); 'spy'`, &dispatched),
	); err != nil {
		t.Fatalf("spy: %v", err)
	}
	_ = finalState
	step("poll success", chromedp.Poll(regionHasText(wireFormHTMXRegion, "Subscribed Ada Lovelace (ada@example.com)"), &ok))
}
