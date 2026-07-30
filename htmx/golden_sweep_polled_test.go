package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for PolledRegion — the htmx package previously had no golden
// tests at all. Timestamp is disabled (ShowTimestamp: false) because the
// component calls time.Now() internally; timestamp rendering is covered by
// the assertion tests in polled_region_test.go.

func TestGoldenSweepPolledRegion(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "polled_region_default", HTML: utils.Render(t, PolledRegion(PolledRegionProps{
			URL: "/partials/stats",
		}))},
		{Name: "polled_region_eager", HTML: utils.Render(t, PolledRegion(PolledRegionProps{
			URL:   "/partials/live",
			Every: "5s",
			Eager: true,
		}))},
		{Name: "polled_region_with_id", HTML: utils.Render(t, PolledRegion(PolledRegionProps{
			BaseProps: utils.BaseProps{ID: "activity-feed"},
			URL:       "/api/activity",
			Every:     "30s",
		}))},
		{Name: "polled_region_assertive", HTML: utils.Render(t, PolledRegion(PolledRegionProps{
			URL:   "/api/alerts",
			Every: "1s",
			Live:  PolledLiveAssertive,
		}))},
	})
}
