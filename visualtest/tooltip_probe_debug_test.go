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
		`(function(){var t=document.querySelector("#tt-edge-tooltip");var w=t.parentElement;var cs=getComputedStyle(t);var ws=getComputedStyle(w);return JSON.stringify({ua:navigator.userAgent.match(/Chrome\/[\d.]+/)[0],sAnchor:CSS.supports("anchor-name: --x"),sArea:CSS.supports("position-area: top"),sTry:CSS.supports("position-try-fallbacks: flip-block"),wAnchorName:ws.anchorName,wHasAttr:w.hasAttribute("data-tc-tooltip"),anchor:cs.anchorName,posAnchor:cs.positionAnchor,area:cs.positionArea||cs.insetArea,top:cs.top,bottom:cs.bottom,insetAuto:cs.inset})})()`,
		&diag,
	)); err != nil {
		t.Fatalf("diag: %v", err)
	}
	t.Logf("diag: %s", diag)

	var probe string
	if err := chromedp.Do(ctx, evalExprInto(
		`(function(){
		var host=document.createElement("div");document.body.appendChild(host);
		host.innerHTML="<div id=\"pa\" style=\"position:relative;width:100px;height:20px;anchor-name:--pa\">A</div>";
		var mk=function(id,inner){var d=document.createElement("div");d.id=id;d.style.position="absolute";d.style.width="40px";d.style.height="10px";d.style.positionAnchor="--pa";d.style.positionArea="top";d.style.inset="auto";d.textContent=inner;return d};
		var child=mk("pc","child");host.querySelector("#pa").appendChild(child);
		var sib=mk("ps","sib");host.appendChild(sib);
		var rc=document.getElementById("pc").getBoundingClientRect();
		var rs=document.getElementById("ps").getBoundingClientRect();
		return JSON.stringify({childTop:rc.top,sibTop:rs.top,anchorTop:document.getElementById("pa").getBoundingClientRect().top});
		})()`,
		&probe,
	)); err != nil {
		t.Fatalf("probe: %v", err)
	}
	t.Logf("probe: %s", probe)
}
