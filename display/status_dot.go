// StatusDot and LivePill components: the colored status indicator family —
// a dot (optionally pulsing) and a dot+text pill — extracted from the
// hand-rolled staleness/live pills in dnsblockd's dashboard and three other
// consumers (2026-10-10 analysis, TODO_LIST #392).
package display

import (
	"github.com/larsartmann/templ-components/utils"
)

// StatusTone is the semantic color of a status indicator.
type StatusTone string

// Status tone constants.
const (
	// StatusToneNeutral is gray — connecting, paused, inert. The zero value.
	StatusToneNeutral StatusTone = "neutral"
	// StatusToneSuccess is green — live, online, healthy.
	StatusToneSuccess StatusTone = "success"
	// StatusToneWarning is amber — stale, reconnecting, waiting.
	StatusToneWarning StatusTone = "warning"
	// StatusToneDanger is red — offline, failed.
	StatusToneDanger StatusTone = "danger"
	// StatusToneInfo is blue — informational state.
	StatusToneInfo StatusTone = "info"
)

// StatusToneIsValid reports whether t is one of the defined tones. Unknown
// tones render neutral (graceful degradation).
func StatusToneIsValid(t StatusTone) bool {
	return t == StatusToneNeutral || t == StatusToneSuccess || t == StatusToneWarning ||
		t == StatusToneDanger || t == StatusToneInfo
}

//nolint:gochecknoglobals // Package-level lookup table for tone background classes
var statusToneBgMap = map[StatusTone]string{
	StatusToneNeutral: "bg-gray-400 dark:bg-gray-500",
	StatusToneSuccess: "bg-green-500 dark:bg-green-400",
	StatusToneWarning: "bg-amber-500 dark:bg-amber-400",
	StatusToneDanger:  "bg-red-500 dark:bg-red-400",
	StatusToneInfo:    "bg-blue-500 dark:bg-blue-400",
}

// statusToneBg resolves the tone's background class (neutral fallback).
func statusToneBg(tone StatusTone) string {
	return utils.Lookup(statusToneBgMap, tone, statusToneBgMap[StatusToneNeutral])
}

// StatusDotProps configures a standalone status dot.
type StatusDotProps struct {
	utils.BaseProps

	// Tone selects the semantic color (neutral gray by default).
	Tone StatusTone
	// Pulse renders the radiating ring — live indicators. The zero value
	// is a static dot; the animation carries a motion-reduce fallback.
	Pulse bool
	// Label, when set, renders screen-reader text naming the state — a
	// colored dot that is the ONLY status indicator must be announced.
	// Without a Label the dot renders aria-hidden (decorative next to
	// visible text you provide).
	Label string
}

// DefaultStatusDotProps returns sensible defaults.
func DefaultStatusDotProps() StatusDotProps {
	return StatusDotProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Tone: StatusToneNeutral,
	}
}

// LivePillProps configures the dot+text status pill.
type LivePillProps struct {
	utils.BaseProps

	// Tone selects the dot's semantic color.
	Tone StatusTone
	// Text is the visible label ("Live", "Reconnecting…").
	Text string
	// Pulse renders the radiating ring on the pill's dot.
	Pulse bool
}

// DefaultLivePillProps returns sensible defaults.
func DefaultLivePillProps() LivePillProps {
	return LivePillProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Tone: StatusToneSuccess,
		Text: "Live",
	}
}
