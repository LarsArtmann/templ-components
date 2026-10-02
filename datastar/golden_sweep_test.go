package datastar

import (
	"testing"
	"time"

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
			Version: DatastarVersion1_0_3,
		}))},
		{Name: "sdk_script_cdn_nonce", HTML: utils.Render(t, SDKScript(SDKScriptProps{
			BaseProps: utils.BaseProps{Nonce: "test-nonce-123"},
			Version:   DatastarVersion1_0_3,
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
		{Name: "live_region_retry_always", HTML: utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/metrics",
			AutoStart: true,
			Retry:     RetryAlways,
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

func TestGoldenSweepPolledRegion(t *testing.T) {
	t.Parallel()

	// Timestamp is pinned via Now so the footer is deterministic; the
	// convention matches the htmx twin's golden sweep.
	fixed := time.Date(2026, 10, 2, 14, 30, 5, 0, time.UTC)

	propsWithClock := func(p PolledRegionProps) PolledRegionProps {
		p.Now = func() time.Time { return fixed }

		return p
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		// Constructor defaults + URL: interval attribute, polite, timestamp footer.
		{Name: "polled_region_default", HTML: utils.Render(t, PolledRegion(propsWithClock(PolledRegionProps{
			URL:           "/partials/stats",
			ShowTimestamp: true,
		})))},
		// Custom interval + aria-label + inert children slot.
		{Name: "polled_region_custom", HTML: utils.Render(t, PolledRegion(propsWithClock(PolledRegionProps{
			BaseProps:     utils.BaseProps{ID: "activity-feed", AriaLabel: "Activity"},
			URL:           "/api/activity",
			Every:         "2m",
			Live:          LiveAssertive,
			ShowTimestamp: true,
		})))},
		// Empty URL renders inert (no interval attribute at all).
		{Name: "polled_region_inert", HTML: utils.Render(t, PolledRegion(propsWithClock(PolledRegionProps{})))},
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
