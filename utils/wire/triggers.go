package wire

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Trigger extends an Action's trigger axis beyond DOM events: interval
// polling and viewport reveal. Both dialects express them — htmx as
// hx-trigger tokens ("every 5s", "revealed", "intersect once threshold:0.5"),
// Datastar as the data-on-interval / data-on-intersect plugin attributes —
// and both spellings are bundle-decoded (see docs/datastar-runtime-facts.md).
// The fields are independent: an Action may poll, reveal, or both (htmx
// merges them into one comma-separated hx-trigger; Datastar renders one
// attribute per trigger source).

// NormalizedInterval converts a duration string into the spelling both
// runtime trigger engines parse correctly: "<n>ms", "<n>s", or a bare
// number (milliseconds).
//
// Why m/h are converted: neither parser handles them — htmx's parseInterval
// and Datastar's duration parser both fall through to parseFloat and read
// the leading digits as MILLISECONDS ("5m" → 5ms, "1h" → 1ms), so a
// natural-looking "every 5m" would poll five times a SECOND... five times a
// millisecond. Minutes and hours become seconds. Invalid or exotic values
// pass through verbatim — the consumer's explicit garbage stays theirs.
func NormalizedInterval(every string) string {
	if strings.TrimSpace(every) == "" {
		return ""
	}

	for _, conv := range []struct {
		re      *regexp.Regexp
		seconds float64
	}{
		{intervalHoursRe, 3600},
		{intervalMinutesRe, 60},
	} {
		match := conv.re.FindStringSubmatch(every)
		if match == nil {
			continue
		}

		n, err := strconv.ParseFloat(match[1], 64)
		if err != nil || n <= 0 {
			return every
		}

		return strconv.FormatFloat(n*conv.seconds, 'f', -1, 64) + "s"
	}

	return every
}

var (
	intervalHoursRe   = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)h$`)
	intervalMinutesRe = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)m$`)
)

// Reveal configures the viewport-reveal trigger: the exchange fires when the
// element scrolls into view — the lazy-load / infinite-scroll pattern. The
// zero value is htmx's "revealed" shorthand (fire ONCE on first entry — the
// common lazy-load case, per the zero-value-equals-documented-default rule);
// set EveryEntry to re-fire on every viewport entry.
type Reveal struct {
	// ThresholdPercent is how much of the element must be visible before the
	// trigger fires: 0 (the default) = any part enters the viewport, 50 =
	// half visible, 100 = fully visible. Values are clamped to [0, 100].
	// Datastar renders named modifiers at the special values (50 → half,
	// 100 → full) and threshold.<n> otherwise; htmx renders
	// threshold:<fraction> only when nonzero.
	ThresholdPercent int

	// EveryEntry re-fires the trigger on EVERY viewport entry instead of the
	// zero-value fire-once behavior (htmx "revealed" = intersect once; the
	// Datastar __once modifier is omitted only when this is set).
	EveryEntry bool

	// Exit fires when the element LEAVES the viewport instead of entering.
	// Datastar renders __exit (bundle-verified). htmx has no exit semantics:
	// it degrades to the entry trigger (documented, like PatchModeReplace's
	// degradation) — a leave-only action cannot be expressed under htmx.
	Exit bool
}

// htmxTrigger renders the merged hx-trigger value for the action's trigger
// sources — the explicit Event (with debounce/changed modifiers) plus the
// interval and reveal triggers, comma-separated per htmx's multi-trigger
// grammar. Returns "" when no trigger applies (htmx then uses its element
// defaults).
func (a Action) htmxTrigger() string {
	var tokens []string

	if a.Event != EventUnspecified && EventIsValid(a.Event) {
		trigger := string(a.Event)
		if a.DebounceMS > 0 {
			if a.Event == EventInput || a.Event == EventChange || a.Event == EventKeyUp {
				trigger += " changed"
			}

			trigger += " delay:" + strconv.Itoa(a.DebounceMS) + "ms"
		}

		if a.ThrottleMS > 0 {
			trigger += " throttle:" + strconv.Itoa(a.ThrottleMS) + "ms"
		}

		tokens = append(tokens, trigger)
	}

	if reveal := a.htmxRevealToken(); reveal != "" {
		tokens = append(tokens, reveal)
	}

	if interval := NormalizedInterval(a.Interval); interval != "" {
		tokens = append(tokens, "every "+interval)
	}

	return strings.Join(tokens, ", ")
}

// htmxRevealToken renders the reveal trigger under htmx: "revealed" for the
// zero value (intersect once), explicit intersect modifiers otherwise.
// Exit degrades to the entry trigger — htmx cannot fire on leave.
func (a Action) htmxRevealToken() string {
	if a.Reveal == nil {
		return ""
	}

	reveal := *a.Reveal

	// "revealed" IS htmx's intersect-once shorthand; Exit degrades to entry.
	if reveal.ThresholdPercent <= 0 && !reveal.EveryEntry {
		return "revealed"
	}

	token := "intersect"

	if !reveal.EveryEntry {
		token += " once"
	}

	if percent := clampPercent(reveal.ThresholdPercent); percent > 0 {
		token += " threshold:" + strconv.FormatFloat(float64(percent)/100, 'f', -1, 64)
	}

	return token
}

// datastarTriggerAttrs renders the Datastar trigger attributes — one entry
// per configured trigger source, since each plugin owns its attribute:
//
//	data-on-interval__duration.5s="@get('/url')"
//	data-on-intersect__once="@get('/url')"
//
// The interval duration lives in the ATTRIBUTE NAME (the runtime's modifier
// grammar splits the name on __ into modifier groups). Returns nil when no
// trigger applies.
func (a Action) datastarTriggerAttrs() templ.Attributes {
	attrs := templ.Attributes{}

	if interval := NormalizedInterval(a.Interval); interval != "" {
		attrs["data-on-interval__duration."+interval] = datastarActionExpr(
			a.method(),
			a.URL,
			a.ContentType,
			a.Selector,
			a.Swap,
		)
	}

	if a.Reveal != nil {
		key := "data-on-intersect"

		reveal := *a.Reveal

		// Zero value fires once (the "revealed" twin); __once is omitted only
		// when the consumer opts into every-entry firing.
		if !reveal.EveryEntry {
			key += "__once"
		}

		switch percent := clampPercent(reveal.ThresholdPercent); {
		case percent == 100:
			key += "__full"
		case percent == 50:
			key += "__half"
		case percent > 0:
			key += "__threshold." + strconv.Itoa(percent)
		}

		if reveal.Exit {
			key += "__exit"
		}

		attrs[key] = datastarActionExpr(a.method(), a.URL, a.ContentType, a.Selector, a.Swap)
	}

	if len(attrs) == 0 {
		return nil
	}

	return attrs
}

// clampPercent bounds a reveal threshold to [0, 100].
func clampPercent(percent int) int {
	switch {
	case percent < 0:
		return 0
	case percent > 100:
		return 100
	default:
		return percent
	}
}
