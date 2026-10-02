package datastar

import (
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestPolledRegionInert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props PolledRegionProps
	}{
		{"empty URL", PolledRegionProps{}},
		{"whitespace URL", PolledRegionProps{URL: "   "}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, PolledRegion(tc.props))
			utils.AssertNotContains(t, output, "data-on-interval")
			utils.AssertNotContains(t, output, "@get(")
		})
	}
}

func TestPolledRegionIntervalAttribute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		every    string
		expected string
	}{
		{"default constructor", DefaultPolledRegionProps().Every, `data-on-interval__duration.10s=`},
		{"empty Every falls back to default", "", `data-on-interval__duration.10s=`},
		{"whitespace Every falls back to default", "  ", `data-on-interval__duration.10s=`},
		{"seconds pass through", "5s", `data-on-interval__duration.5s=`},
		{"milliseconds pass through", "500ms", `data-on-interval__duration.500ms=`},
		{"bare number is runtime milliseconds", "2500", `data-on-interval__duration.2500=`},
		{"minutes convert to seconds", "2m", `data-on-interval__duration.120s=`},
		{"decimal minutes convert", "1.5m", `data-on-interval__duration.90s=`},
		{"hours convert to seconds", "1h", `data-on-interval__duration.3600s=`},
		{"decimal hours convert", "0.5h", `data-on-interval__duration.1800s=`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			props := PolledRegionProps{URL: "/partials/stats", Every: tc.every}
			output := utils.Render(t, PolledRegion(props))
			utils.AssertContains(t, output, tc.expected)
		})
	}
}

func TestPolledIntervalValuePoisonousTypos(t *testing.T) {
	t.Parallel()

	// The runtime parses ONLY ms/s/bare (parseFloat-ms). Without the
	// conversion, "5m" would poll every 5 MILLISECONDS.
	tests := []struct{ in, want string }{
		{"", "10s"},
		{"10s", "10s"},
		{"500ms", "500ms"},
		{"2500", "2500"},
		{"2m", "120s"},
		{"1h", "3600s"},
		// Invalid conversions pass through verbatim — the consumer's
		// explicit garbage stays theirs (same policy as the htmx twin).
		{"-5m", "-5m"},
		{"0h", "0h"},
		{"soon", "soon"},
		{"every 5s", "every 5s"},
	}

	for _, tc := range tests {
		if got := polledIntervalValue(tc.in); got != tc.want {
			t.Errorf("polledIntervalValue(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestPolledRegionURLPassesThroughActionExpr(t *testing.T) {
	t.Parallel()

	props := PolledRegionProps{URL: "/partials/stats?a=1"}
	output := utils.Render(t, PolledRegion(props))
	// templ escapes the expression quotes inside the attribute value.
	utils.AssertContains(t, output, `="@get(&#39;/partials/stats?a=1&#39;)"`)

	t.Run("single quotes are escaped", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x?name=o'brien"}))
		// The injected quote becomes \' inside the expression; templ renders
		// the expression quotes AND the injected quote as &#39; (the
		// backslash stays literal).
		utils.AssertContains(t, output, `@get(&#39;/x?name=o\&#39;brien&#39;)`)
	})
}

func TestPolledRegionDefaults(t *testing.T) {
	t.Parallel()

	props := DefaultPolledRegionProps()
	if props.Every != "10s" {
		t.Errorf("Every = %q, want 10s", props.Every)
	}

	if props.Live != LivePolite {
		t.Errorf("Live = %q, want polite", props.Live)
	}

	if !props.ShowTimestamp {
		t.Error("ShowTimestamp = false, want true")
	}
}

func TestPolledRegionTimestamp(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 10, 2, 14, 30, 5, 0, time.UTC)

	props := DefaultPolledRegionProps()
	props.URL = "/partials/stats"
	props.Now = func() time.Time { return fixed }

	output := utils.Render(t, PolledRegion(props))
	utils.AssertContains(t, output, `datetime="2026-10-02T14:30:05Z"`)
	utils.AssertContains(t, output, "Updated 14:30:05")

	t.Run("hidden by default in a literal", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x"}))
		if strings.Contains(output, "Updated") {
			t.Errorf("timestamp rendered without ShowTimestamp:\n%s", output)
		}
	})

	t.Run("custom format", func(t *testing.T) {
		t.Parallel()

		props := PolledRegionProps{
			URL:           "/x",
			ShowTimestamp: true,
			Now:           func() time.Time { return fixed },
			TimeFormat:    time.RFC3339,
		}
		output := utils.Render(t, PolledRegion(props))
		utils.AssertContains(t, output, "Updated 2026-10-02T14:30:05Z")
	})
}

func TestPolledRegionLiveFallback(t *testing.T) {
	t.Parallel()

	// Unknown politeness falls back to polite (shared livePolitenessValue).
	props := PolledRegionProps{URL: "/x", Live: LivePoliteness("shouty")}
	output := utils.Render(t, PolledRegion(props))
	utils.AssertContains(t, output, `aria-live="polite"`)
}

func TestPolledRegionBasePropsPropagation(t *testing.T) {
	t.Parallel()

	props := PolledRegionProps{
		BaseProps: utils.BaseProps{
			ID:        "stats-region",
			Class:     "min-h-40",
			AriaLabel: "Live statistics",
			Attrs:     templ.Attributes{"data-test": "polled"},
		},
		URL: "/partials/stats",
	}
	output := utils.Render(t, PolledRegion(props))

	for _, want := range []string{
		`id="stats-region"`,
		`aria-label="Live statistics"`,
		`data-test="polled"`,
		"min-h-40",
	} {
		utils.AssertContains(t, output, want)
	}
}
