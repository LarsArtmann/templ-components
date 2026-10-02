package wire

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestNormalizedInterval(t *testing.T) {
	t.Parallel()

	tests := []struct{ in, want string }{
		{"", ""},
		{"  ", ""},
		{"10s", "10s"},
		{"500ms", "500ms"},
		{"2500", "2500"},
		{"2m", "120s"},
		{"1.5m", "90s"},
		{"1h", "3600s"},
		{"0.5h", "1800s"},
		// Invalid conversions pass through verbatim.
		{"-5m", "-5m"},
		{"0h", "0h"},
		{"soon", "soon"},
	}

	for _, tc := range tests {
		if got := NormalizedInterval(tc.in); got != tc.want {
			t.Errorf("NormalizedInterval(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestTriggerHTMXIntervalOnly(t *testing.T) {
	t.Parallel()

	attrs := Action{URL: "/x", Interval: "5s"}.Attributes()

	assertAttrContains(t, attrs, "hx-trigger", "every 5s")

	// An interval-only Action must NOT also fire on click: htmx fires the
	// element default only when hx-trigger is ABSENT — here it is present.
	if trigger, _ := attrs["hx-trigger"].(string); strings.Contains(trigger, "click") {
		t.Errorf("interval-only action must not render a click token; got %q", trigger)
	}
}

func TestTriggerHTMXReveal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reveal Reveal
		want   string
	}{
		{"zero reveal is htmx's revealed shorthand", Reveal{}, "revealed"},
		{"once + threshold", Reveal{ThresholdPercent: 50}, "intersect once threshold:0.5"},
		{"fraction trimmed", Reveal{ThresholdPercent: 75}, "intersect once threshold:0.75"},
		{"every-entry intersect", Reveal{ThresholdPercent: 100, EveryEntry: true}, "intersect threshold:1"},
		{"exit degrades to entry trigger", Reveal{Exit: true}, "revealed"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			attrs := Action{URL: "/x", Reveal: &tc.reveal}.Attributes()
			assertAttrContains(t, attrs, "hx-trigger", tc.want)
		})
	}
}

func TestTriggerHTMXMergedTokens(t *testing.T) {
	t.Parallel()

	attrs := Action{
		URL:        "/x",
		Event:      EventInput,
		DebounceMS: 300,
		Interval:   "10s",
	}.Attributes()

	assertAttrContains(t, attrs, "hx-trigger", "input changed delay:300ms, every 10s")
}

func TestTriggerDatastarIntervalOnly(t *testing.T) {
	t.Parallel()

	attrs := Action{Transport: TransportDatastar, URL: "/x", Interval: "5s"}.Attributes()

	assertAttrContains(t, attrs, "data-on-interval__duration.5s", "@get('/x')")

	// No event attribute: an interval-only Action is a pure poller.
	for key := range attrs {
		if strings.HasPrefix(key, "data-on:") {
			t.Errorf("interval-only action rendered an event attribute %q", key)
		}
	}
}

func TestTriggerDatastarReveal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reveal Reveal
		want   string
	}{
		{"zero reveal fires once (revealed twin)", Reveal{}, "data-on-intersect__once"},
		{"every entry omits once", Reveal{EveryEntry: true}, "data-on-intersect"},
		{"half", Reveal{ThresholdPercent: 50}, "data-on-intersect__once__half"},
		{"full", Reveal{ThresholdPercent: 100}, "data-on-intersect__once__full"},
		{"explicit threshold", Reveal{ThresholdPercent: 75}, "data-on-intersect__once__threshold.75"},
		{"threshold clamped high", Reveal{ThresholdPercent: 150}, "data-on-intersect__once__full"},
		{"exit", Reveal{Exit: true}, "data-on-intersect__once__exit"},
		{"combined mods", Reveal{EveryEntry: true, ThresholdPercent: 50, Exit: true}, "data-on-intersect__half__exit"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			attrs := Action{Transport: TransportDatastar, URL: "/x", Reveal: &tc.reveal}.Attributes()
			assertAttrContains(t, attrs, tc.want, "@get('/x')")
		})
	}
}

func TestTriggerDatastarEventPlusInterval(t *testing.T) {
	t.Parallel()

	attrs := Action{
		Transport: TransportDatastar,
		URL:       "/x",
		Event:     EventClick,
		Interval:  "30s",
	}.Attributes()

	assertAttrContains(t, attrs, "data-on:click", "@get('/x')")
	assertAttrContains(t, attrs, "data-on-interval__duration.30s", "@get('/x')")
}

func TestTriggerEmptyURLInert(t *testing.T) {
	t.Parallel()

	attrs := Action{Interval: "5s", Reveal: &Reveal{}}.Attributes()
	if attrs != nil {
		t.Errorf("empty URL must stay inert with triggers set; got %v", attrs)
	}

	attrs = Action{Transport: TransportDatastar, Interval: "5s"}.Attributes()
	if attrs != nil {
		t.Errorf("empty URL must stay inert under datastar with triggers set; got %v", attrs)
	}
}

func TestTriggerIntervalNormalizationDialects(t *testing.T) {
	t.Parallel()

	// The m/h poison applies to BOTH dialects: "5m" must never survive as
	// 5ms on either side.
	htmxAttrs := Action{URL: "/x", Interval: "5m"}.Attributes()
	assertAttrContains(t, htmxAttrs, "hx-trigger", "every 300s")

	dsAttrs := Action{Transport: TransportDatastar, URL: "/x", Interval: "1h"}.Attributes()
	assertAttrContains(t, dsAttrs, "data-on-interval__duration.3600s", "")
}

func TestTriggerSwapAndTargetCompose(t *testing.T) {
	t.Parallel()

	// Trigger sources compose with the rest of the contract.
	htmxAttrs := Action{URL: "/x", Interval: "10s", Target: "#region", Swap: PatchModeOuter}.Attributes()
	assertAttrContains(t, htmxAttrs, "hx-target", "#region")
	assertAttrContains(t, htmxAttrs, "hx-swap", "outerHTML")

	reveal := Reveal{}
	dsAttrs := Action{
		Transport: TransportDatastar,
		URL:       "/x",
		Reveal:    &reveal,
		Selector:  "#card",
		Swap:      PatchModeInner,
	}.Attributes()
	// Selector renders only under form encoding, and Swap renders nothing
	// under Datastar (response-header driven) — neither leaks into the
	// trigger expression.
	assertAttrContains(t, dsAttrs, "data-on-intersect__once", "@get('/x')")

	if value := fmt.Sprint(dsAttrs["data-on-intersect__once"]); strings.Contains(value, "mode:") || strings.Contains(value, "selector:") {
		t.Fatalf("trigger expression must not carry client options the runtime ignores: %q", value)
	}
}

// assertAttrContains fails unless attrs carries key with a value containing
// want (an empty want only requires the key's presence).
func assertAttrContains(t *testing.T, attrs templ.Attributes, key, want string) {
	t.Helper()

	value, ok := attrs[key]
	if !ok {
		t.Fatalf("attribute %q missing; have keys: %v", key, attrKeys(attrs))
	}

	if want != "" && !strings.Contains(fmt.Sprint(value), want) {
		t.Errorf("%s = %v, want it to contain %q", key, value, want)
	}
}

func attrKeys(attrs templ.Attributes) []string {
	keys := make([]string, 0, len(attrs))

	for key := range attrs {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
