package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestPolledRegionRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PolledRegion(PolledRegionProps{
		URL:    "/partials/stats",
		Every:  "10s",
		Eager:  true,
		Swap:   SwapInnerHTML,
		Live:   PolledLiveAssertive,
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
			URL:          "/stats",
			Every:        "10s",
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
		URL:  "/stats",
		Every: "1s",
		Swap: "bogus",
		Live: "bogus",
	}))
	utils.AssertContains(t, output, `hx-swap="outerHTML"`)
	utils.AssertContains(t, output, `aria-live="polite"`)
}
