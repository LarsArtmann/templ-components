// Package tone is the library's canonical color language: one Tone axis
// ("why is this thing colored"), one Hue axis ("which hue family"), and the
// two tables that map them to complete Tailwind class literals.
//
// Before this package existed, the same semantic was modeled thirteen
// different ways (10 enums, 12 lookup maps, 3 palettes, ~10 fixed-color
// sites) — see docs/reviews/2026-10-15_color-domain-model.html. Phase 0 is
// ADDITIVE ONLY: existing enums gain adapters, existing lookup maps keep
// rendering, and equivalence tests prove the tables describe reality
// byte-for-byte. The maps collapse INTO these tables at v2 (ADR-0039/0044,
// TODO #431).
//
// # Scanner-safety contract
//
// Every class string in this package is a COMPLETE LITERAL, including the
// dark: prefixed half. Tailwind's scanner must see the whole token to emit
// CSS; variants are never concatenated at runtime (see AGENTS.md "Tailwind
// variant prefixes must be complete literals"). Pair exists to make the
// light/dark halves structural, not to build strings dynamically.
//
// # Warning dialect (the known split brain)
//
// Two warning hues ship today: yellow (Feedback/Alert, Badge, Kanban,
// StatCard) and amber (StatusDot, ScrollbackLog, errorpage rejection).
// The tone table canonicalizes YELLOW (the majority dialect) until the
// #428 decision lands; amber-rendering components keep their hue via the
// Hue axis and their own lookups.
package tone

// Tone is the semantic reason something is colored: what the color MEANS,
// not which pigment it uses. Six values cover every color decision the
// library makes; hue-level (categorical) decisions use Hue instead.
type Tone string

const (
	// ToneNeutral is the absence of semantic emphasis (gray).
	ToneNeutral Tone = "neutral"
	// ToneSuccess confirms a positive outcome (green).
	ToneSuccess Tone = "success"
	// ToneWarning signals caution short of failure (yellow today; see the
	// package doc's warning-dialect note and TODO #428).
	ToneWarning Tone = "warning"
	// ToneDanger signals failure or destruction (red). Canonical spelling
	// is "danger"; the FeedbackError value ("error") maps here via
	// FeedbackType.Tone() and collapses at v2.
	ToneDanger Tone = "danger"
	// ToneInfo is neutral-emphasis communication (blue).
	ToneInfo Tone = "info"
	// ToneAccent is the brand/primary slot (blue today; becomes the ADR-0044
	// accent token family at v2).
	ToneAccent Tone = "accent"
)

// ToneDefault is the canonical default (Neutral).
const ToneDefault = ToneNeutral

// ToneIsValid reports whether t is a recognized Tone.
func ToneIsValid(t Tone) bool {
	switch t {
	case ToneNeutral, ToneSuccess, ToneWarning, ToneDanger, ToneInfo, ToneAccent:
		return true
	default:
		return false
	}
}

// Tones returns every Tone, in declaration (semantic severity) order.
func Tones() []Tone {
	return []Tone{ToneNeutral, ToneSuccess, ToneWarning, ToneDanger, ToneInfo, ToneAccent}
}
