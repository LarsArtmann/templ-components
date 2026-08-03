package echarts

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestEChartRendersElement(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		Element: `<div id="chart1" style="width:600px;height:400px;"></div>`,
		Script:  `var chart=echarts.init(document.getElementById('chart1'));`,
		Nonce:   "abc123",
	}))
	utils.AssertContains(t, html, `id="chart1"`)
	utils.AssertContains(t, html, "tc-echarts")
}

func TestEChartRendersScriptWithNonce(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		Element: `<div id="c"></div>`,
		Script:  `echarts.init(document.getElementById('c'));`,
		Nonce:   "test-nonce",
	}))
	utils.AssertContains(t, html, `nonce="test-nonce"`)
	utils.AssertContains(t, html, "echarts.init")
}

func TestEChartEmptyElementRendersNothing(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		Element: "",
		Script:  "console.log('test');",
		Nonce:   "nonce",
	}))
	if html != "" {
		t.Errorf("expected empty output for empty Element, got %q", html)
	}
}

func TestEChartDarkModeBridgeInjected(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		Element:        `<div id="c"></div>`,
		Script:         `echarts.init(document.getElementById('c'));`,
		Nonce:          "bridge-nonce",
		DarkModeBridge: true,
	}))
	utils.AssertContains(t, html, "tcEChartsDarkBridge")
	utils.AssertContains(t, html, "MutationObserver")
}

func TestEChartDarkModeBridgeDisabled(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		Element:        `<div id="c"></div>`,
		Script:         `console.log('x');`,
		Nonce:          "nonce",
		DarkModeBridge: false,
	}))
	utils.AssertNotContains(t, html, "tcEChartsDarkBridge")
}

func TestEChartAllScriptsHaveNonce(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		Element:        `<div id="c"></div>`,
		Script:         `echarts.init(document.getElementById('c'));`,
		Nonce:          "csp-nonce",
		DarkModeBridge: true,
	}))

	scriptCount := strings.Count(html, "<script")
	nonceCount := strings.Count(html, `nonce="csp-nonce"`)
	if scriptCount != nonceCount {
		t.Errorf("expected %d script tags to all have nonce, got %d nonce attrs", scriptCount, nonceCount)
	}
}

func TestEChartsScriptURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version EChartsVersion
		cdn     string
		want    string
	}{
		{"default", "", "", "https://cdn.jsdelivr.net/npm/echarts@5.5.0/dist/echarts.min.js"},
		{"custom version", "5.4.3", "", "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"},
		{"custom cdn", "5.5.0", "https://unpkg.com", "https://unpkg.com/echarts@5.5.0/dist/echarts.min.js"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := echartsScriptURL(tt.version, tt.cdn)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSDKScriptRendersCorrectURL(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, SDKScript(SDKScriptProps{
		Version: EChartsVersion5_5,
		Nonce:   "sdk-nonce",
	}))
	utils.AssertContains(t, html, "echarts@5.5.0/dist/echarts.min.js")
	utils.AssertContains(t, html, `nonce="sdk-nonce"`)
}

func TestSDKScriptCustomCDN(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, SDKScript(SDKScriptProps{
		Version: EChartsVersion5_5,
		CDN:     "https://my-cdn.example.com",
		Nonce:   "nonce",
	}))
	utils.AssertContains(t, html, "https://my-cdn.example.com/echarts@5.5.0/dist/echarts.min.js")
}

func TestSDKScriptSelfHosted(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, SDKScript(SDKScriptProps{
		Src:   "/static/echarts.min.js",
		Nonce: "nonce",
	}))
	utils.AssertContains(t, html, `src="/static/echarts.min.js"`)
	utils.AssertNotContains(t, html, "cdn.jsdelivr.net")
}

func TestEChartBasePropsPropagation(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, EChart(EChartsProps{
		BaseProps: utils.BaseProps{
			Class:     "max-w-2xl",
			ID:        "my-echart",
			AriaLabel: "Sales chart",
		},
		Element: `<div id="c"></div>`,
		Nonce:   "nonce",
	}))
	utils.AssertContains(t, html, "max-w-2xl")
	utils.AssertContains(t, html, `id="my-echart"`)
	utils.AssertContains(t, html, `aria-label="Sales chart"`)
}
