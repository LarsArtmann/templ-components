package echarts

import (
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

	// Every <script tag in the output must have nonce=
	for _, line := range splitByScript(html) {
		if line == "" {
			continue
		}

		if !containsStr(line, `nonce="csp-nonce"`) {
			t.Errorf("script tag missing nonce: %s", line)
		}
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

func splitByScript(s string) []string {
	return splitBy(s, "<script")
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && indexOf(s, sub) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}

	return -1
}

func splitBy(s, sep string) []string {
	var result []string

	for {
		idx := indexOf(s, sep)
		if idx < 0 {
			if s != "" {
				result = append(result, s)
			}

			break
		}

		if idx > 0 {
			result = append(result, s[:idx])
		}

		s = s[idx:]
	}

	// Find end of this script tag section
	if len(s) > 0 {
		endIdx := indexOf(s, "</script>")
		if endIdx >= 0 {
			result = append(result, s[:endIdx+len("</script>")])
			s = s[endIdx+len("</script>"):]
		} else {
			result = append(result, s)
			s = ""
		}
	}

	_ = s

	return result
}
