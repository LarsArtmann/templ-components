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
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(
			`window.__dsEvents=[];['datastar-fetch','datastar-patch-elements'].forEach(n=>document.addEventListener(n,()=>window.__dsEvents.push(n),true));`,
			nil,
		),
	); err != nil {
		t.Fatal(err)
	}

	var clicked string
	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`(() => { const b=[...document.querySelectorAll('button')].find(el => (el.getAttribute('data-on:click')||'').includes('/api/wire/busy')); if(b) b.click(); return 'clicked=' + !!b; })()`,
		&clicked,
	)); err != nil {
		t.Fatal(err)
	}

	if err := chromedp.Run(ctx, chromedp.Sleep(3*time.Second)); err != nil {
		t.Fatal(err)
	}

	var out string
	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`JSON.stringify({region: (document.getElementById('wire-busy-datastar-out')||{}).innerText, events: window.__dsEvents})`,
		&out,
	)); err != nil {
		t.Fatal(err)
	}
	t.Log("PROBE:", clicked, out)
}
