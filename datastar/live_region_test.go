package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestLiveRegionAutoStart(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:       "/stream/metrics",
		AutoStart: true,
	}))

	utils.AssertContains(t, output, `data-init="@get('/stream/metrics')"`)
	utils.AssertContains(t, output, `aria-live="polite"`)
}

func TestLiveRegionManualStart(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:       "/stream/metrics",
		AutoStart: false,
	}))

	utils.AssertNotContains(t, output, "data-init")
}

func TestLiveRegionWithID(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		BaseProps: utils.BaseProps{ID: "live-stats"},
		URL:       "/stream/stats",
		AutoStart: true,
	}))

	utils.AssertContains(t, output, `id="live-stats"`)
}

func TestLiveRegionAssertive(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:  "/stream/alerts",
		Live: LiveAssertive,
	}))

	utils.AssertContains(t, output, `aria-live="assertive"`)
}

func TestLiveRegionInvalidPolitenessFallback(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:  "/stream/data",
		Live: "bogus",
	}))

	utils.AssertContains(t, output, `aria-live="polite"`)
}

func TestLiveRegionDefaults(t *testing.T) {
	t.Parallel()

	props := DefaultLiveRegionProps()

	if !props.AutoStart {
		t.Error("expected AutoStart=true")
	}

	if props.Live != LivePolite {
		t.Errorf("expected Live=%s, got %s", LivePolite, props.Live)
	}
}

func TestLiveRegionRendersChildren(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:       "/stream/data",
		AutoStart: true,
	}))

	// The div should be present even with no children
	utils.AssertContains(t, output, "<div")
	utils.AssertContains(t, output, "</div>")
}
