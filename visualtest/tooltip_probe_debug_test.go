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
		`(function(){var b=document.querySelector("#tt-center-tooltip");var a=document.querySelector("[data-tc-tooltip-anchor]");var btn=document.querySelector("#tt-center-btn");var cs=getComputedStyle(b);var r=b.getBoundingClientRect();var ra=a.getBoundingClientRect();var rb=btn.getBoundingClientRect();return JSON.stringify({bLeft:r.left,bRight:r.right,bW:r.width,aLeft:ra.left,aW:ra.width,btnLeft:rb.left,btnW:rb.width,left:cs.left,transform:cs.transform,area:cs.positionArea,margin:cs.margin,justify:cs.justifySelf})})()`,
		&diag,
	)); err != nil {
		t.Fatalf("diag: %v", err)
	}
	t.Logf("diag: %s", diag)
}
