package tone

// Surfaces is the per-tone rendering kit: which classes express a semantic
// tone on each surface kind. Values are anchored to the feedback dialect
// (Alert/Toast — the oldest and most-used semantic renderer), so
// For(ToneSuccess).Text is exactly the Text class feedbackStyleMap carries
// for success today.
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

// toneTable maps every Tone to its surfaces. The hue backing each tone:
// neutral=gray, success=green, warning=yellow (majority dialect — see the
// package doc), danger=red, info=blue, accent=blue (the ADR-0044 accent
// slot until v2 swaps it to the accent token family).
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
		Dot:      hueTable[HueYellow].BGSolid,
		Text:     Pair{"text-yellow-800", "dark:text-yellow-200"},
		TextSoft: Pair{"text-yellow-600", darkTextYellow},
		Border:   Pair{"border-yellow-200", "dark:border-yellow-800"},
		BGTint:   Pair{"bg-yellow-50", "dark:bg-yellow-900/20"},
		BGSolid:  hueTable[HueYellow].BGSolidStrong,
		Icon:     Pair{"text-yellow-600", darkTextYellow},
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
		TextSoft: Pair{"text-blue-600", darkTextBlue},
		Border:   Pair{"border-blue-200", "dark:border-blue-800"},
		BGTint:   Pair{"bg-blue-50", "dark:bg-blue-900/20"},
		BGSolid:  hueTable[HueBlue].BGSolidStrong,
		Icon:     Pair{"text-blue-600", darkTextBlue},
	},
	ToneAccent: {
		Dot:      hueTable[HueBlue].BGSolid,
		Text:     Pair{"text-blue-800", "dark:text-blue-200"},
		TextSoft: Pair{"text-blue-600", darkTextBlue},
		Border:   Pair{"border-blue-200", "dark:border-blue-800"},
		BGTint:   Pair{"bg-blue-50", "dark:bg-blue-900/20"},
		BGSolid:  hueTable[HueBlue].BGSolidStrong,
		Icon:     Pair{"text-blue-600", darkTextBlue},
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
