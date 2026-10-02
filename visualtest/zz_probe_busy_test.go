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
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`window.__ev=[];['datastar-fetch','error'].forEach(n=>document.addEventListener(n,e=>window.__ev.push(n+(e.detail?':'+JSON.stringify(e.detail).slice(0,120):'')),true));`, nil),
	); err != nil {
		t.Fatal(err)
	}

	clickErr := chromedp.Run(ctx, chromedp.Click(`button[data-on\:click*="/api/wire/busy"]`, chromedp.ByQuery))

	time.Sleep(3 * time.Second)

	var out string
	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`JSON.stringify({ev: window.__ev, region: (document.getElementById('wire-busy-datastar-out')||{}).innerText})`,
		&out)); err != nil {
		t.Fatal(err)
	}
	t.Logf("CLICK ERR: %v", clickErr)
	t.Log("PROBE:", out)
}
