package display

import (
	"testing"
	"time"

	"github.com/larsartmann/templ-components/utils"
)

// TestFormatRelativeTimeBoundaries pins the formatRelativeTime bucket
// boundaries with an injected clock (M20/F089): every fencepost — 59s/60s,
// 59m/1h, 23h/24h, 6d/7d, 29d/30d — plus future timestamps (symmetric
// absolute distance) and the 30-day fallback to an absolute date.
func TestFormatRelativeTimeBoundaries(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 13, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		delta time.Duration
		want  string
	}{
		{"zero distance", 0, "just now"},
		{"one second", 1 * time.Second, "just now"},
		{"59 seconds", 59 * time.Second, "just now"},
		{"60 seconds", 60 * time.Second, "1 minute ago"},
		{"61 seconds", 61 * time.Second, "1 minute ago"},
		{"90 seconds", 90 * time.Second, "1 minute ago"},
		{"119 seconds", 119 * time.Second, "1 minute ago"},
		{"2 minutes", 120 * time.Second, "2 minutes ago"},
		{"59 minutes", 59 * time.Minute, "59 minutes ago"},
		{"1 hour", 1 * time.Hour, "1 hour ago"},
		{"90 minutes", 90 * time.Minute, "1 hour ago"},
		{"2 hours", 2 * time.Hour, "2 hours ago"},
		{"23 hours", 23 * time.Hour, "23 hours ago"},
		{"24 hours", 24 * time.Hour, "1 day ago"},
		{"25 hours", 25 * time.Hour, "1 day ago"},
		{"47 hours", 47 * time.Hour, "1 day ago"},
		{"2 days", 48 * time.Hour, "2 days ago"},
		{"6 days", 6 * 24 * time.Hour, "6 days ago"},
		{"7 days", 7 * 24 * time.Hour, "1 week ago"},
		{"8 days", 8 * 24 * time.Hour, "1 week ago"},
		{"13 days", 13 * 24 * time.Hour, "1 week ago"},
		{"14 days", 14 * 24 * time.Hour, "2 weeks ago"},
		{"29 days", 29 * 24 * time.Hour, "4 weeks ago"},
		{"30 days", 30 * 24 * time.Hour, "Aug 14, 2026"}, // absolute fallback
		{"60 days", 60 * 24 * time.Hour, "Jul 15, 2026"},
	}

	// Only the absolute-date fallback differs between directions: a future
	// date prints itself, not the mirrored past date.
	futureOverrides := map[string]string{
		"30 days": "Oct 13, 2026",
		"60 days": "Nov 12, 2026",
	}

	for _, tt := range tests {
		t.Run(tt.name+" ago", func(t *testing.T) {
			t.Parallel()

			got := formatRelativeTime(now.Add(-tt.delta), now)
			if got != tt.want {
				t.Errorf("formatRelativeTime(now-%s) = %q, want %q", tt.delta, got, tt.want)
			}
		})

		// Future timestamps mirror: distance is absolute, wording stays
		// "... ago" by design (the title attribute carries the absolute time).
		t.Run(tt.name+" future", func(t *testing.T) {
			t.Parallel()

			want := tt.want
			if override, ok := futureOverrides[tt.name]; ok {
				want = override
			}

			got := formatRelativeTime(now.Add(tt.delta), now)
			if got != want {
				t.Errorf("formatRelativeTime(now+%s) = %q, want %q (symmetric distance)", tt.delta, got, want)
			}
		})
	}
}

// TestRelativeTimeNowInjection pins the render-time clock injection: a
// pinned Now produces a deterministic relative string AND datetime attribute
// regardless of wall-clock time (the M20 determinism-injection preference).
func TestRelativeTimeNowInjection(t *testing.T) {
	t.Parallel()

	pinned := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)

	props := DefaultRelativeTimeProps()
	props.Time = pinned.Add(-2 * time.Hour)
	props.Now = pinned
	props.AutoRefresh = false // keep the script out of the assertion surface

	html := utils.Render(t, RelativeTime(props))

	utils.AssertContains(t, html, `datetime="2026-02-28T22:00:00Z"`)
	utils.AssertContains(t, html, ">2 hours ago<")
}
