package visualtest

import (
	"testing"

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

	var diag string
	if err := chromedp.Do(ctx, evalExprInto(
		`(function(){var attached=typeof window.tcTooltipAriaAttached;var b=document.querySelector("#tt-edge-tooltip");var r=b.getBoundingClientRect();b.style.marginLeft="";var r2=b.getBoundingClientRect();return JSON.stringify({attached:attached,w0:r.width,left0:r.left,right0:r.right,left1:r2.left,right1:r2.right,focusTest:(function(){document.getElementById("tt-edge-btn").focus();return "fired"})()})})()`,
		&diag,
	)); err != nil {
		t.Fatalf("diag: %v", err)
	}
	t.Logf("diag: %s", diag)

	if err := chromedp.Do(ctx, pollTrue(
		`(function(){var b=document.querySelector("#tt-edge-tooltip");var r=b.getBoundingClientRect();return r.right<=window.innerWidth-8;})()`,
		pollFast()...,
	)); err != nil {
		var late string
		_ = chromedp.Do(ctx, evalExprInto(
			`(function(){var b=document.querySelector("#tt-edge-tooltip");return JSON.stringify({right:b.getBoundingClientRect().right,ml:b.style.marginLeft,active:document.activeElement&&document.activeElement.id})})()`,
			&late,
		))
		t.Logf("late: %s", late)
		t.Fatalf("clamp never applied after focus")
	}
	t.Logf("clamp applied after focus")
}
