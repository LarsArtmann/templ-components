package visualtest

import (
	"fmt"
	"testing"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
)

func TestTooltipProbeDebug(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser")
	}

	srv := tooltipAnchorE2EServer(t)
	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx, chromedp.Navigate(srv.URL+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("nav: %v", err)
	}

	var box overlayBox
	if err := chromedp.Do(ctx, evalExprInto(
		`(() => { const r = document.querySelector('#tt-edge-btn').getBoundingClientRect(); return {left:r.left, top:r.top, right:r.right, bottom:r.bottom}; })()`,
		&box,
	)); err != nil {
		t.Fatalf("box: %v", err)
	}
	t.Logf("button box: %+v", box)

	if err := chromedp.Do(ctx, mouseEvent(input.DispatchMouseEventTypeMouseMoved, box.insideX(), box.insideY())); err != nil {
		t.Fatalf("move: %v", err)
	}

	var diag string
	if err := chromedp.Do(ctx, evalExprInto(
		`(function(){var b=document.querySelector('#tt-edge-btn');var w=b.parentElement;var t=w.querySelector('[role=tooltip]');return JSON.stringify({hover:b.matches(':hover'),wClass:w.className,tDisplay:t?getComputedStyle(t).display:'none',tExists:!!t})})()`,
		&diag,
	)); err != nil {
		t.Fatalf("diag: %v", err)
	}
	t.Logf("diag: %s", diag)

	_ = fmt.Sprint()
}
