package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// Temporary debug test — delete after use.
func TestDebugWireFormValidation(t *testing.T) {
	srv := wireFormE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTimeout()

	var (
		ok        bool
		preserved string
		dump      string
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

	step("navigate", chromedp.Navigate(srv.URL+"/"))
	step("gate", chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &ok))
	step("fill name", chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="name"]`), "Ada Lovelace", chromedp.NodeVisible))
	step("fill email", chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="email"]`), "ada@example", chromedp.NodeVisible))
	step("submit #1", chromedp.Click(formSel(wireFormHTMXRegion, `button[type="submit"]`), chromedp.NodeVisible))
	step("poll summary", chromedp.Poll(regionHasText(wireFormHTMXRegion, "1 error found"), &ok))
	step("evaluate value", chromedp.Evaluate(formValueExpr(wireFormHTMXRegion, `input[name="email"]`), &preserved))
	t.Logf("preserved=%q", preserved)
	step("dump form attrs", chromedp.Evaluate(`Array.from(document.querySelectorAll('#wire-form-htmx-region form')).map(f => f.getAttribute('hx-post') + '|' + f.getAttribute('hx-trigger') + '|' + (f.getAttribute('class')||'').slice(0,30)).join(' ;; ')`, &dump))
	step("htmx.process region", chromedp.Evaluate(`htmx.process(document.querySelector('#wire-form-htmx-region')); 'processed'`, &dump))
	t.Logf("form attrs: %s", dump)
	step("fix email", chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="email"]`), ".com", chromedp.NodeVisible))
	step("submit #2", chromedp.Click(formSel(wireFormHTMXRegion, `button[type="submit"]`), chromedp.NodeVisible))
	step("poll success", chromedp.Poll(regionHasText(wireFormHTMXRegion, "Subscribed Ada Lovelace (ada@example.com)"), &ok))
}
