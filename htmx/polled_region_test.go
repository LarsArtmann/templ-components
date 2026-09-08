package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestPolledRegionRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PolledRegion(PolledRegionProps{
		URL:   "/partials/stats",
		Every: "10s",
		Eager: true,
		Swap:  SwapInnerHTML,
		Live:  PolledLiveAssertive,
		BaseProps: utils.BaseProps{
			ID: "stats-region",
		},
	}))
	utils.AssertContains(t, output, `id="stats-region"`)
	utils.AssertContains(t, output, `hx-get="/partials/stats"`)
	utils.AssertContains(t, output, `hx-trigger="load, every 10s"`)
	utils.AssertContains(t, output, `hx-swap="innerHTML"`)
	utils.AssertContains(t, output, `aria-live="assertive"`)
}

func TestPolledRegionNotEager(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PolledRegion(PolledRegionProps{
		URL:   "/api/health",
		Every: "5s",
	}))
	utils.AssertNotContains(t, output, "load,")
	utils.AssertContains(t, output, `hx-trigger="every 5s"`)
	utils.AssertContains(t, output, `hx-swap="outerHTML"`)
	utils.AssertContains(t, output, `aria-live="polite"`)
}

func TestPolledRegionTimestamp(t *testing.T) {
	t.Parallel()
	t.Run("shows timestamp by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PolledRegion(DefaultPolledRegionProps()))
		utils.AssertContains(t, output, "Updated")
		utils.AssertContains(t, output, "<time")
	})
	t.Run("hides timestamp when disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PolledRegion(PolledRegionProps{
			URL:           "/stats",
			Every:         "10s",
			ShowTimestamp: false,
		}))
		utils.AssertNotContains(t, output, "Updated")
	})
}

func TestPolledRegionDefaults(t *testing.T) {
	t.Parallel()

	props := DefaultPolledRegionProps()
	if props.Every != "10s" {
		t.Errorf("expected Every=10s, got %s", props.Every)
	}

	if props.Swap != SwapOuterHTML {
		t.Errorf("expected Swap=outerHTML, got %s", props.Swap)
	}

	if props.Live != PolledLivePolite {
		t.Errorf("expected Live=polite, got %s", props.Live)
	}

	if !props.ShowTimestamp {
		t.Error("expected ShowTimestamp=true")
	}
}

func TestPolledLiveIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value PolledLive
		want  bool
	}{
		{PolledLivePolite, true},
		{PolledLiveAssertive, true},
		{PolledLiveOff, true},
		{"bogus", false},
		{"", false},
	}
	for _, tt := range tests {
		got := PolledLiveIsValid(tt.value)
		if got != tt.want {
			t.Errorf("PolledLiveIsValid(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestPolledRegionInvalidValuesFallBack(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PolledRegion(PolledRegionProps{
		URL:   "/stats",
		Every: "1s",
		Swap:  "bogus",
		Live:  "bogus",
	}))
	utils.AssertContains(t, output, `hx-swap="outerHTML"`)
	utils.AssertContains(t, output, `aria-live="polite"`)
}

func TestPolledRegionCustomTimeFormat(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PolledRegion(PolledRegionProps{
		URL:           "/stats",
		Every:         "10s",
		ShowTimestamp: true,
		TimeFormat:    "2006-01-02 15:04:05",
	}))
	utils.AssertContains(t, output, "Updated")
	utils.AssertContains(t, output, "20") // year component appears in the output
}

func TestPolledRegionCustomTrigger(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PolledRegion(PolledRegionProps{
		URL:     "/stats",
		Trigger: "stats-refresh from:body",
	}))
	utils.AssertContains(t, output, `hx-trigger="stats-refresh from:body"`)
	// Every and Eager should be ignored when Trigger is set
	utils.AssertNotContains(t, output, "every")
	utils.AssertNotContains(t, output, "load,")
}

// TestPolledRegionBusyCue pins the aria-busy loading cue (parity with
// datastar.LiveRegion): eager regions render busy + the marker + the
// clearing script; non-eager and custom-trigger regions render neither.
func TestPolledRegionBusyCue(t *testing.T) {
	t.Parallel()

	t.Run("eager region renders busy cue and clearing script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PolledRegion(PolledRegionProps{
			URL:   "/stats",
			Every: "5s",
			Eager: true,
			BaseProps: utils.BaseProps{
				Nonce: "test-nonce",
			},
		}))
		utils.AssertContains(t, output, `aria-busy="true"`)
		utils.AssertContains(t, output, "data-tc-polled-busy")
		utils.AssertContains(t, output, `nonce="test-nonce"`)
		utils.AssertContains(t, output, "htmx:afterRequest")
		utils.AssertContains(t, output, "window.tcPolledBusyAttached")
	})

	t.Run("non-eager region renders no busy cue", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PolledRegion(PolledRegionProps{
			URL:   "/stats",
			Every: "5s",
		}))
		utils.AssertNotContains(t, output, "aria-busy")
		utils.AssertNotContains(t, output, "data-tc-polled-busy")
		utils.AssertNotContains(t, output, "tcPolledBusyAttached")
	})

	t.Run("custom trigger region renders no busy cue", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PolledRegion(PolledRegionProps{
			URL:     "/stats",
			Trigger: "stats-refresh from:body",
			Eager:   true, // ignored when Trigger is set — timing is consumer-owned
		}))
		utils.AssertNotContains(t, output, "aria-busy")
		utils.AssertNotContains(t, output, "tcPolledBusyAttached")
	})

	t.Run("empty nonce omits the nonce attribute entirely", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PolledRegion(PolledRegionProps{
			URL:   "/stats",
			Every: "5s",
			Eager: true,
		}))
		utils.AssertContains(t, output, "tcPolledBusyAttached")
		utils.AssertNotContains(t, output, `nonce=""`)
		utils.AssertNotContains(t, output, "nonce=")
	})
}
