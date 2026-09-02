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

// A whitespace-only URL is morally empty: @get('   ') would fetch a garbage
// relative URL (or throw) on every page load. It must degrade exactly like
// the empty URL case.
func TestLiveRegionWhitespaceURLDegradesGracefully(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:       "   ",
		AutoStart: true,
	}))

	utils.AssertNotContains(t, output, "data-init")
	utils.AssertContains(t, output, "<div")
}

// Auto-started regions carry the aria-busy loading cue and the singleton
// script that clears it once the first patch arrives. Manual regions (no
// auto-connect) must not claim a busy state — nothing is loading yet.
func TestLiveRegionBusyState(t *testing.T) {
	t.Parallel()

	t.Run("auto-start marks the region aria-busy", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/metrics",
			AutoStart: true,
		}))

		utils.AssertContains(t, output, `aria-busy="true"`)
		utils.AssertContains(t, output, "data-tc-live-busy")
		utils.AssertContains(t, output, "tcLiveBusyAttached")
	})

	t.Run("manual start renders no busy state and no script", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/metrics",
			AutoStart: false,
		}))

		utils.AssertNotContains(t, output, "aria-busy")
		utils.AssertNotContains(t, output, "tcLiveBusyAttached")
	})

	t.Run("empty URL renders no busy state and no script", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(DefaultLiveRegionProps()))

		utils.AssertNotContains(t, output, "aria-busy")
		utils.AssertNotContains(t, output, "tcLiveBusyAttached")
	})

	t.Run("clearing script carries the CSP nonce", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			BaseProps: utils.BaseProps{Nonce: "busy-nonce"},
			URL:       "/stream/metrics",
			AutoStart: true,
		}))

		utils.AssertContains(t, output, `<script nonce="busy-nonce">`)
	})
}

// URL-with-quotes must escape the single quotes so a crafted URL cannot
// inject Datastar expression syntax into the data-init attribute. LiveRegion
// shares the escaping used by datastar.Get — this test pins the inheritance.
func TestLiveRegionURLEscaping(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LiveRegion(LiveRegionProps{
		URL:       "/stream?q='",
		AutoStart: true,
	}))

	// The injected quote becomes \' inside the expression; templ renders
	// both the expression quotes and the escaped quote as &#39;.
	utils.AssertContains(t, output, `@get(&#39;/stream?q=\&#39;&#39;)`)
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
