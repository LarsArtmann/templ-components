package tone

// Utilities is the per-hue class grammar: seven utility slots, each a
// complete light/dark Pair. The shade numbers are the library's uniform
// grammar, anchored to what the components actually render today:
//
//	Text           600 / dark:400   (chart series text, icons on tint)
//	TextStrong     700 / dark:400   (trend arrows, log tags)
//	TextOnTint     800 / dark:200   (body text on 50-tinted backgrounds)
//	BGSolid        500 / dark:400   (status dots, kanban lane dots)
//	BGSolidStrong  600 / dark:400   (segment bars, solid fills)
//	BGTint         50  / dark:900/30 (stat tiles, icon chips)
//	Border         200 / dark:800   (alert and card borders)
//
// Gray carries two deliberate exceptions, kept to match shipped pixels:
// BGSolid is 400/dark:500 (the dot convention renders grays one step
// darker) and Border uses gray-700 in dark mode (the feedback neutral
// fallback). Cells no component uses yet are inert until adopted — the
// tables describe the grammar, the components choose the slots.
type Utilities struct {
	Text           Pair
	TextStrong     Pair
	TextOnTint     Pair
	BGSolid        Pair
	BGSolidStrong  Pair
	BGTint         Pair
	Border         Pair
}

// hueTable is the single source for per-hue literals. Keys are every Hue in
// Hues(); completeness is guarded by TestHueTableExhaustive.
//
//nolint:gochecknoglobals // package-level lookup table by design
var hueTable = map[Hue]Utilities{
	HueBlue: {
		Text:           Pair{"text-blue-600", "dark:text-blue-400"},
		TextStrong:     Pair{"text-blue-700", "dark:text-blue-400"},
		TextOnTint:     Pair{"text-blue-800", "dark:text-blue-200"},
		BGSolid:        Pair{"bg-blue-500", "dark:bg-blue-400"},
		BGSolidStrong:  Pair{"bg-blue-600", "dark:bg-blue-400"},
		BGTint:         Pair{"bg-blue-50", "dark:bg-blue-900/30"},
		Border:         Pair{"border-blue-200", "dark:border-blue-800"},
	},
	HueEmerald: {
		Text:           Pair{"text-emerald-600", "dark:text-emerald-400"},
		TextStrong:     Pair{"text-emerald-700", "dark:text-emerald-400"},
		TextOnTint:     Pair{"text-emerald-800", "dark:text-emerald-200"},
		BGSolid:        Pair{"bg-emerald-500", "dark:bg-emerald-400"},
		BGSolidStrong:  Pair{"bg-emerald-600", "dark:bg-emerald-400"},
		BGTint:         Pair{"bg-emerald-50", "dark:bg-emerald-900/30"},
		Border:         Pair{"border-emerald-200", "dark:border-emerald-800"},
	},
	HueAmber: {
		Text:           Pair{"text-amber-600", "dark:text-amber-400"},
		TextStrong:     Pair{"text-amber-700", "dark:text-amber-400"},
		TextOnTint:     Pair{"text-amber-800", "dark:text-amber-200"},
		BGSolid:        Pair{"bg-amber-500", "dark:bg-amber-400"},
		BGSolidStrong:  Pair{"bg-amber-600", "dark:bg-amber-400"},
		BGTint:         Pair{"bg-amber-50", "dark:bg-amber-900/30"},
		Border:         Pair{"border-amber-200", "dark:border-amber-800"},
	},
	HueRose: {
		Text:           Pair{"text-rose-600", "dark:text-rose-400"},
		TextStrong:     Pair{"text-rose-700", "dark:text-rose-400"},
		TextOnTint:     Pair{"text-rose-800", "dark:text-rose-200"},
		BGSolid:        Pair{"bg-rose-500", "dark:bg-rose-400"},
		BGSolidStrong:  Pair{"bg-rose-600", "dark:bg-rose-400"},
		BGTint:         Pair{"bg-rose-50", "dark:bg-rose-900/30"},
		Border:         Pair{"border-rose-200", "dark:border-rose-800"},
	},
	HueViolet: {
		Text:           Pair{"text-violet-600", "dark:text-violet-400"},
		TextStrong:     Pair{"text-violet-700", "dark:text-violet-400"},
		TextOnTint:     Pair{"text-violet-800", "dark:text-violet-200"},
		BGSolid:        Pair{"bg-violet-500", "dark:bg-violet-400"},
		BGSolidStrong:  Pair{"bg-violet-600", "dark:bg-violet-400"},
		BGTint:         Pair{"bg-violet-50", "dark:bg-violet-900/30"},
		Border:         Pair{"border-violet-200", "dark:border-violet-800"},
	},
	HueCyan: {
		Text:           Pair{"text-cyan-600", "dark:text-cyan-400"},
		TextStrong:     Pair{"text-cyan-700", "dark:text-cyan-400"},
		TextOnTint:     Pair{"text-cyan-800", "dark:text-cyan-200"},
		BGSolid:        Pair{"bg-cyan-500", "dark:bg-cyan-400"},
		BGSolidStrong:  Pair{"bg-cyan-600", "dark:bg-cyan-400"},
		BGTint:         Pair{"bg-cyan-50", "dark:bg-cyan-900/30"},
		Border:         Pair{"border-cyan-200", "dark:border-cyan-800"},
	},
	HueOrange: {
		Text:           Pair{"text-orange-600", "dark:text-orange-400"},
		TextStrong:     Pair{"text-orange-700", "dark:text-orange-400"},
		TextOnTint:     Pair{"text-orange-800", "dark:text-orange-200"},
		BGSolid:        Pair{"bg-orange-500", "dark:bg-orange-400"},
		BGSolidStrong:  Pair{"bg-orange-600", "dark:bg-orange-400"},
		BGTint:         Pair{"bg-orange-50", "dark:bg-orange-900/30"},
		Border:         Pair{"border-orange-200", "dark:border-orange-800"},
	},
	HuePink: {
		Text:           Pair{"text-pink-600", "dark:text-pink-400"},
		TextStrong:     Pair{"text-pink-700", "dark:text-pink-400"},
		TextOnTint:     Pair{"text-pink-800", "dark:text-pink-200"},
		BGSolid:        Pair{"bg-pink-500", "dark:bg-pink-400"},
		BGSolidStrong:  Pair{"bg-pink-600", "dark:bg-pink-400"},
		BGTint:         Pair{"bg-pink-50", "dark:bg-pink-900/30"},
		Border:         Pair{"border-pink-200", "dark:border-pink-800"},
	},
	HueGray: {
		Text:           Pair{"text-gray-600", "dark:text-gray-400"},
		TextStrong:     Pair{"text-gray-700", "dark:text-gray-400"},
		TextOnTint:     Pair{"text-gray-800", "dark:text-gray-200"},
		BGSolid:        Pair{"bg-gray-400", "dark:bg-gray-500"}, // dot convention: grays render one step darker
		BGSolidStrong:  Pair{"bg-gray-600", "dark:bg-gray-400"},
		BGTint:         Pair{"bg-gray-50", "dark:bg-gray-900/30"},
		Border:         Pair{"border-gray-200", "dark:border-gray-700"}, // feedback-neutral fallback dialect
	},
	HueGreen: {
		Text:           Pair{"text-green-600", "dark:text-green-400"},
		TextStrong:     Pair{"text-green-700", "dark:text-green-400"},
		TextOnTint:     Pair{"text-green-800", "dark:text-green-200"},
		BGSolid:        Pair{"bg-green-500", "dark:bg-green-400"},
		BGSolidStrong:  Pair{"bg-green-600", "dark:bg-green-400"},
		BGTint:         Pair{"bg-green-50", "dark:bg-green-900/30"},
		Border:         Pair{"border-green-200", "dark:border-green-800"},
	},
	HueYellow: {
		Text:           Pair{"text-yellow-600", "dark:text-yellow-400"},
		TextStrong:     Pair{"text-yellow-700", "dark:text-yellow-400"},
		TextOnTint:     Pair{"text-yellow-800", "dark:text-yellow-200"},
		BGSolid:        Pair{"bg-yellow-500", "dark:bg-yellow-400"},
		BGSolidStrong:  Pair{"bg-yellow-600", "dark:bg-yellow-400"},
		BGTint:         Pair{"bg-yellow-50", "dark:bg-yellow-900/30"},
		Border:         Pair{"border-yellow-200", "dark:border-yellow-800"},
	},
	HueRed: {
		Text:           Pair{"text-red-600", "dark:text-red-400"},
		TextStrong:     Pair{"text-red-700", "dark:text-red-400"},
		TextOnTint:     Pair{"text-red-800", "dark:text-red-200"},
		BGSolid:        Pair{"bg-red-500", "dark:bg-red-400"},
		BGSolidStrong:  Pair{"bg-red-600", "dark:bg-red-400"},
		BGTint:         Pair{"bg-red-50", "dark:bg-red-900/30"},
		Border:         Pair{"border-red-200", "dark:border-red-800"},
	},
	HuePurple: {
		Text:           Pair{"text-purple-600", "dark:text-purple-400"},
		TextStrong:     Pair{"text-purple-700", "dark:text-purple-400"},
		TextOnTint:     Pair{"text-purple-800", "dark:text-purple-200"},
		BGSolid:        Pair{"bg-purple-500", "dark:bg-purple-400"},
		BGSolidStrong:  Pair{"bg-purple-600", "dark:bg-purple-400"},
		BGTint:         Pair{"bg-purple-50", "dark:bg-purple-900/30"},
		Border:         Pair{"border-purple-200", "dark:border-purple-800"},
	},
}

// UtilitiesFor returns the class grammar for h. Unknown hues fall back to
// Gray (the neutral default) — same contract as utils.Lookup.
func UtilitiesFor(h Hue) Utilities {
	if u, ok := hueTable[h]; ok {
		return u
	}

	return hueTable[HueGray]
}
