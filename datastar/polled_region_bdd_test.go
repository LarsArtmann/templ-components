package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- PolledRegion Behavior ---

func TestPolledRegionUserGetsAutoRefreshingRegion(t *testing.T) {
	t.Parallel()

	t.Run("user sees the interval poll wired to their endpoint", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{
			URL:   "/partials/stats",
			Every: "10s",
		}))

		// The interval lives in the attribute NAME (runtime modifier
		// grammar); the fetch expression is the attribute value.
		utils.AssertContains(t, output, `data-on-interval__duration.10s="@get('/partials/stats')"`)
	})

	t.Run("user sees fresh initial content, not a loading state", func(t *testing.T) {
		t.Parallel()

		// Initial content is server-rendered into the children slot: the
		// region is never empty while waiting for the first tick, so no
		// eager-fetch mechanism exists (a leading interval + the self-patch
		// model would refetch eagerly forever).
		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x"}))
		utils.AssertNotContains(t, output, "leading")
		utils.AssertNotContains(t, output, "aria-busy")
	})

	t.Run("operator can verify polling is active via the timestamp", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", ShowTimestamp: true}))
		utils.AssertContains(t, output, "Updated")
	})

	t.Run("endpoint responses replace the region via the runtime's outer patch", func(t *testing.T) {
		t.Parallel()

		// Contract note pinned as documentation: the region carries its id
		// (consumer-set or inherited markup) and the endpoint returns a
		// fragment whose root id matches, so the runtime's default outer
		// mode self-patches — the component itself renders no patching
		// attributes beyond the interval trigger.
		output := utils.Render(t, PolledRegion(PolledRegionProps{
			BaseProps: utils.BaseProps{ID: "stats"},
			URL:       "/partials/stats",
		}))
		utils.AssertContains(t, output, `id="stats"`)
		utils.AssertNotContains(t, output, "hx-")
	})
}

func TestPolledRegionUserTypoCannotCausePollStorm(t *testing.T) {
	t.Parallel()

	// "5m" and "1h" LOOK like five minutes / one hour. The runtime would
	// parse both as MILLISECONDS (truncating at parseFloat) — a 1ms
	// self-poll storm. The component converts m/h to seconds first.
	t.Run("minutes and hours become seconds", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", Every: "5m"}))
		utils.AssertContains(t, output, `data-on-interval__duration.300s=`)

		output = utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", Every: "1h"}))
		utils.AssertContains(t, output, `data-on-interval__duration.3600s=`)
	})
}
