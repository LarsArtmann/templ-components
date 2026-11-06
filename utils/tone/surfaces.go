package tone

// Surfaces is the per-tone rendering kit: which classes express a semantic
// tone on each surface kind. Values are anchored to the feedback dialect
// (Alert/Toast — the oldest and most-used semantic renderer), so
// For(ToneSuccess).Text is exactly the Text class feedback renders for
// success today (byte-pinned by feedback.TestToneAdapterEquivalence).
type Surfaces struct {
	// Dot is a solid indicator dot (status dots, kanban lane dots).
	Dot Pair
	// Text is body text on a tinted background (alert message text).
	Text Pair
	// TextSoft is de-emphasized text or an icon on a tinted background.
	TextSoft Pair
	// Border outlines a tinted container (alerts, tinted cards).
	Border Pair
	// BGTint is a soft tinted background (alert fills).
	BGTint Pair
	// BGSolid is a solid saturated fill (progress bars, pills).
	BGSolid Pair
	// Icon is an icon on a tinted background (alert icons).
	Icon Pair
}

// toneTable maps every Tone to its surfaces. The pigment backing each tone:
// neutral=gray, success=green, warning=amber (the #428 decision — the
// one semantic warning hue; categorical yellow survives on the Hue axis
// for Kanban/StatCard), danger=red, info=blue (semantic info keeps the
// palette per ADR-0044 decision 3), accent=the accent token family
// (ADR-0044/0049 emission — the v2 swap this table was built for).
//
// Neutral deviates from the uniform tint grammar where feedback's fallback
// ships it: BGTint uses gray-800/50 in dark mode (not 900/30). Completeness
// is guarded by TestToneTableExhaustive; byte-equivalence with the feedback
// map is guarded by feedback.TestToneAdapterEquivalence.
//
//nolint:gochecknoglobals // package-level lookup table by design
var toneTable = map[Tone]Surfaces{
	ToneNeutral: {
		Dot:      hueTable[HueGray].BGSolid,
		Text:     Pair{"text-gray-800", "dark:text-gray-200"},
		TextSoft: Pair{"text-gray-600", darkTextGray},
		Border:   Pair{"border-gray-200", "dark:border-gray-700"},
		BGTint:   Pair{"bg-gray-50", "dark:bg-gray-800/50"},
		BGSolid:  hueTable[HueGray].BGSolidStrong,
		Icon:     Pair{"text-gray-600", darkTextGray},
	},
	ToneSuccess: {
		Dot:      hueTable[HueGreen].BGSolid,
		Text:     Pair{"text-green-800", "dark:text-green-200"},
		TextSoft: Pair{"text-green-600", darkTextGreen},
		Border:   Pair{"border-green-200", "dark:border-green-800"},
		BGTint:   Pair{"bg-green-50", "dark:bg-green-900/20"},
		BGSolid:  hueTable[HueGreen].BGSolidStrong,
		Icon:     Pair{"text-green-600", darkTextGreen},
	},
	ToneWarning: {
		Dot:      hueTable[HueAmber].BGSolid,
		Text:     Pair{"text-amber-800", "dark:text-amber-200"},
		TextSoft: Pair{"text-amber-600", darkTextAmber},
		Border:   Pair{"border-amber-200", "dark:border-amber-800"},
		BGTint:   Pair{"bg-amber-50", "dark:bg-amber-900/20"},
		BGSolid:  hueTable[HueAmber].BGSolidStrong,
		Icon:     Pair{"text-amber-600", darkTextAmber},
	},
	ToneDanger: {
		Dot:      hueTable[HueRed].BGSolid,
		Text:     Pair{"text-red-800", "dark:text-red-200"},
		TextSoft: Pair{"text-red-600", darkTextRed},
		Border:   Pair{"border-red-200", "dark:border-red-800"},
		BGTint:   Pair{"bg-red-50", "dark:bg-red-900/20"},
		BGSolid:  hueTable[HueRed].BGSolidStrong,
		Icon:     Pair{"text-red-600", darkTextRed},
	},
	ToneInfo: {
		Dot:      hueTable[HueBlue].BGSolid,
		Text:     Pair{"text-blue-800", "dark:text-blue-200"},
		TextSoft: Pair{textBlue600, darkTextBlue},
		Border:   Pair{"border-blue-200", "dark:border-blue-800"},
		BGTint:   Pair{"bg-blue-50", "dark:bg-blue-900/20"},
		BGSolid:  hueTable[HueBlue].BGSolidStrong,
		Icon:     Pair{textBlue600, darkTextBlue},
	},
	ToneAccent: {
		Dot:      Pair{"bg-accent-500", "dark:bg-accent-400"},
		Text:     Pair{textAccent600, darkTextAccent},
		TextSoft: Pair{textAccent600, darkTextAccent},
		Border:   Pair{"border-accent-200", "dark:border-accent-800"},
		BGTint:   Pair{"bg-accent-50", "dark:bg-accent-900/20"},
		BGSolid:  Pair{"bg-accent-600", "dark:bg-accent-400"},
		Icon:     Pair{textAccent600, darkTextAccent},
	},
}

// For returns the rendering surfaces for t. Unknown tones fall back to
// Neutral — same contract as utils.Lookup.
func For(t Tone) Surfaces {
	if s, ok := toneTable[t]; ok {
		return s
	}

	return toneTable[ToneNeutral]
}
