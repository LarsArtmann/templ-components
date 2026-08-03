package echarts

import (
	"context"
	"fmt"
	"html"
	"io"

	"github.com/a-h/templ"
)

// darkModeBridgeJS returns the singleton JavaScript that syncs ECharts chart
// themes with the Tailwind `.dark` class on <html>. It:
//   - Scans all ECharts instances on the page
//   - Applies dark or light color overrides via setOption
//   - Uses MutationObserver to react to runtime theme toggles
//
// The bridge is idempotent — the window.tcEChartsDarkBridge guard ensures it
// only runs once even if multiple EChart components are on the page.
func darkModeBridgeJS() string {
	return `if(!window.tcEChartsDarkBridge){window.tcEChartsDarkBridge=true;` +
		`function tcECHT(){return document.documentElement.classList.contains('dark')?'dark':'light';}` +
		`function tcECHApply(){` +
		`var t=tcECHT();var dark=t==='dark';` +
		`var txtColor=dark?'#9ca3af':'#374151';` +
		`var axisColor=dark?'#4b5563':'#d1d5db';` +
		`var labelColor=dark?'#9ca3af':'#6b7280';` +
		`document.querySelectorAll('[_echarts_instance_]').forEach(function(el){` +
		`var c=echarts.getInstanceByDom(el);` +
		`if(!c)return;` +
		`c.setOption({backgroundColor:'transparent',textStyle:{color:txtColor},` +
		`xAxis:{axisLine:{lineStyle:{color:axisColor}},axisLabel:{color:labelColor},splitLine:{lineStyle:{color:dark?'#374151':'#e5e7eb'}}},` +
		`yAxis:{axisLine:{lineStyle:{color:axisColor}},axisLabel:{color:labelColor},splitLine:{lineStyle:{color:dark?'#374151':'#e5e7eb'}}}},` +
		`{merge:true});});}` +
		`tcECHApply();` +
		`new MutationObserver(function(){tcECHApply();}).observe(document.documentElement,{attributes:true,attributeFilter:['class']});` +
		`}` + "\n"
}

// darkModeBridgeComponent renders the dark mode bridge script in a CSP-safe
// <script nonce> tag. Singleton via the window.tcEChartsDarkBridge guard.
func darkModeBridgeComponent(nonce string) templ.Component {
	escapedNonce := html.EscapeString(nonce)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, "<script nonce=\"%s\">\n%s</script>\n", escapedNonce, darkModeBridgeJS()); err != nil {
			return fmt.Errorf("write echarts dark mode bridge: %w", err)
		}

		return nil
	})
}
