package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestGoldenSweepSDKScript(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "sdk_script_default", HTML: utils.Render(t, SDKScript(DefaultSDKScriptProps()))},
		{Name: "sdk_script_self_hosted", HTML: utils.Render(t, SDKScript(SDKScriptProps{
			Src: "/static/datastar.js",
		}))},
		{Name: "sdk_script_custom_cdn", HTML: utils.Render(t, SDKScript(SDKScriptProps{
			CDN:     "https://unpkg.com",
			Version: DatastarVersion1_0_2,
		}))},
	})
}

func TestGoldenSweepLiveRegion(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "live_region_default", HTML: utils.Render(t, LiveRegion(DefaultLiveRegionProps()))},
		{Name: "live_region_with_url", HTML: utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/metrics",
			AutoStart: true,
		}))},
		{Name: "live_region_manual", HTML: utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/data",
			AutoStart: false,
		}))},
		{Name: "live_region_assertive", HTML: utils.Render(t, LiveRegion(LiveRegionProps{
			URL:  "/stream/alerts",
			Live: LiveAssertive,
		}))},
	})
}

func TestGoldenSweepIndicator(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "indicator_default", HTML: utils.Render(t, Indicator(IndicatorProps{
			Signal: "saving",
		}))},
	})
}

func TestGoldenSweepSSEErrorHandling(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "sse_error_handling_default", HTML: utils.Render(t, SSEErrorHandling(SSEErrorHandlingConfig{
			Nonce:      "test-nonce",
			DurationMS: 6000,
		}))},
	})
}
