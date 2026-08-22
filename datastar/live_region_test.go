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

	utils.AssertContainsAll(t, output, "data-init", "@get(", "/stream/metrics")
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

func TestLiveRegionEmptyURLDegradesGracefully(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		AutoStart: true,
	}))

	// @get('') makes the runtime throw FetchNoUrlProvided on every page
	// load — an empty URL must render a plain container instead.
	utils.AssertNotContains(t, output, "data-init")
	utils.AssertContains(t, output, "<div")
	utils.AssertContains(t, output, `aria-live="polite"`)
}

func TestLiveRegionRetryAlways(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:       "/stream/metrics",
		AutoStart: true,
		Retry:     RetryAlways,
	}))

	utils.AssertContains(t, output, `data-init="@get(&#39;/stream/metrics&#39;, {retry: &#39;always&#39;})"`)
}

func TestLiveRegionRetryModes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		retry RetryMode
		want  string
	}{
		{"empty defaults to auto", "", `data-init="@get(&#39;/stream/data&#39;)"`},
		{"auto is bare", RetryAuto, `data-init="@get(&#39;/stream/data&#39;)"`},
		{"error mode", RetryError, `data-init="@get(&#39;/stream/data&#39;, {retry: &#39;error&#39;})"`},
		{"never mode", RetryNever, `data-init="@get(&#39;/stream/data&#39;, {retry: &#39;never&#39;})"`},
		{"bogus degrades to auto", "bogus", `data-init="@get(&#39;/stream/data&#39;)"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, LiveRegion(LiveRegionProps{
				URL:       "/stream/data",
				AutoStart: true,
				Retry:     tt.retry,
			}))

			utils.AssertContains(t, output, tt.want)
		})
	}
}

func TestLiveRegionCancellation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		retry        RetryMode
		cancellation RequestCancellation
		want         string
	}{
		{
			"none is the bare runtime default",
			RetryAuto, CancellationNone,
			`data-init="@get(&#39;/stream/events&#39;)"`,
		},
		{
			"cleanup alone",
			RetryAuto, CancellationCleanup,
			`data-init="@get(&#39;/stream/events&#39;, {requestCancellation: &#39;cleanup&#39;})"`,
		},
		{
			"always plus cleanup — the swapped-region combination",
			RetryAlways,
			CancellationCleanup,
			`data-init="@get(&#39;/stream/events&#39;, {retry: &#39;always&#39;, requestCancellation: &#39;cleanup&#39;})"`,
		},
		{
			"error retry plus cleanup keeps both options",
			RetryError,
			CancellationCleanup,
			`data-init="@get(&#39;/stream/events&#39;, {retry: &#39;error&#39;, requestCancellation: &#39;cleanup&#39;})"`,
		},
		{
			"bogus cancellation degrades to none",
			RetryAuto, "bogus",
			`data-init="@get(&#39;/stream/events&#39;)"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, LiveRegion(LiveRegionProps{
				URL:          "/stream/events",
				AutoStart:    true,
				Retry:        tt.retry,
				Cancellation: tt.cancellation,
			}))

			utils.AssertContains(t, output, tt.want)
		})
	}
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
