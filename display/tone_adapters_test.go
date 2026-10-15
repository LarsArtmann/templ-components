package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

// TestStatusToneAdapterEquivalence: the tone adapter + DotHue must derive
// exactly what statusToneBgMap ships today — zero visual change, proven.
func TestStatusToneAdapterEquivalence(t *testing.T) {
	t.Parallel()

	for _, st := range []StatusTone{
		StatusToneNeutral, StatusToneSuccess, StatusToneWarning,
		StatusToneDanger, StatusToneInfo,
	} {
		derived := tone.UtilitiesFor(st.DotHue()).BGSolid.Classes()
		shipped := statusToneBgMap[st]

		if derived != shipped {
			t.Errorf(
				"StatusTone %s: DotHue().BGSolid = %q, statusToneBgMap ships %q — adapter drifted from pixels",
				st, derived, shipped,
			)
		}
	}

	wantTones := map[StatusTone]tone.Tone{
		StatusToneNeutral: tone.ToneNeutral,
		StatusToneSuccess: tone.ToneSuccess,
		StatusToneWarning: tone.ToneWarning, // amber hue, yellow-canonical tone: the #428 split
		StatusToneDanger:  tone.ToneDanger,
		StatusToneInfo:    tone.ToneInfo,
	}

	for st, want := range wantTones {
		if got := st.Tone(); got != want {
			t.Errorf("StatusTone %s.Tone() = %s, want %s", st, got, want)
		}
	}
}

// TestKanbanToneAdapterEquivalence: KanbanTone.Hue() must derive exactly
// what kanbanToneLookup ships today.
func TestKanbanToneAdapterEquivalence(t *testing.T) {
	t.Parallel()

	for kt, shipped := range kanbanToneLookup {
		derived := tone.UtilitiesFor(kt.Hue()).BGSolid.Classes()

		if derived != shipped {
			t.Errorf(
				"KanbanTone %s: Hue().BGSolid = %q, kanbanToneLookup ships %q",
				kt, derived, shipped,
			)
		}
	}
}

// TestStatToneAdapterEquivalence: StatTone.Hue() must derive exactly the
// tint/icon strings statTileLookup ships (tile + icon halves of the
// "tile|icon" value).
func TestStatToneAdapterEquivalence(t *testing.T) {
	t.Parallel()

	for st, shipped := range statTileLookup {
		u := tone.UtilitiesFor(st.Hue())

		if tint, _, ok := cutStatTile(shipped); ok && tint != u.BGTint.Classes() {
			t.Errorf(
				"StatTone %s: Hue().BGTint = %q, statTileLookup tile = %q",
				st, u.BGTint.Classes(), tint,
			)
		}
	}
}

// cutStatTile splits a statTileLookup "tile|icon" value.
func cutStatTile(v string) (string, string, bool) {
	for i := len(v) - 1; i >= 0; i-- {
		if v[i] == '|' {
			return v[:i], v[i+1:], true
		}
	}

	return "", "", false
}

// TestTrendDirectionAdapterEquivalence: the trend text classes must be
// derivable from the adapter tone/hue grammar (green/amber use the Strong
// slot, red uses the Text slot — pinned as shipped).
func TestTrendDirectionAdapterEquivalence(t *testing.T) {
	t.Parallel()

	derived := map[TrendDirection]string{
		TrendUp:   tone.UtilitiesFor(tone.HueGreen).TextStrong.Classes(),
		TrendDown: tone.UtilitiesFor(tone.HueRed).Text.Classes(),
		TrendWarn: tone.UtilitiesFor(tone.HueAmber).TextStrong.Classes(),
	}

	for td, want := range map[TrendDirection]string{
		TrendUp:   "text-green-700 dark:text-green-400",
		TrendDown: "text-red-600 dark:text-red-400",
		TrendWarn: "text-amber-700 dark:text-amber-400",
	} {
		if derived[td] != want {
			t.Errorf("TrendDirection %s derived %q, want shipped %q", td, derived[td], want)
		}

		if got := td.Tone(); got == "" {
			t.Errorf("TrendDirection %s.Tone() empty", td)
		}
	}
}

// TestScrollbackToneAdapterEquivalence: the scrollback tag classes must be
// derivable from the hue grammar (neutral/blue/red use Text; green/amber
// use TextStrong — pinned as shipped).
func TestScrollbackToneAdapterEquivalence(t *testing.T) {
	t.Parallel()

	derived := map[ScrollbackTone]string{
		ScrollbackToneNeutral: tone.UtilitiesFor(tone.HueGray).TextMuted.Classes(),
		ScrollbackToneInfo:    tone.UtilitiesFor(tone.HueBlue).Text.Classes(),
		ScrollbackToneSuccess: tone.UtilitiesFor(tone.HueGreen).TextStrong.Classes(),
		ScrollbackToneWarning: tone.UtilitiesFor(tone.HueAmber).TextStrong.Classes(),
		ScrollbackToneDanger:  tone.UtilitiesFor(tone.HueRed).Text.Classes(),
	}

	for st, shipped := range scrollbackToneLookup {
		if derived[st] != shipped {
			t.Errorf(
				"ScrollbackTone %s: derived %q, scrollbackToneLookup ships %q",
				st, derived[st], shipped,
			)
		}
	}
}

// TestBadgeTypeAdapterEquivalence: BadgeType's dot is a single light-mode
// class (no dark half — deliberate), and must equal the hue's BGSolid light
// half; the semantic mapping is pinned.
func TestBadgeTypeAdapterEquivalence(t *testing.T) {
	t.Parallel()

	hueFor := map[BadgeType]tone.Hue{
		BadgePrimary: tone.HueBlue,
		BadgeSuccess: tone.HueGreen,
		BadgeWarning: tone.HueYellow,
		BadgeError:   tone.HueRed,
		BadgeInfo:    tone.HueBlue,
		BadgeNeutral: tone.HueGray,
	}

	// Badge dots are single light-mode classes (no dark half — deliberate).
	// The gray neutral dot is bg-gray-500 (one step darker than the uniform
	// grammar's 400 — the same gray-dot convention as statusToneBgMap).
	wantDots := map[BadgeType]string{
		BadgePrimary: tone.UtilitiesFor(tone.HueBlue).BGSolid.Light,
		BadgeSuccess: tone.UtilitiesFor(tone.HueGreen).BGSolid.Light,
		BadgeWarning: tone.UtilitiesFor(tone.HueYellow).BGSolid.Light,
		BadgeError:   tone.UtilitiesFor(tone.HueRed).BGSolid.Light,
		BadgeInfo:    tone.UtilitiesFor(tone.HueBlue).BGSolid.Light,
		BadgeNeutral: "bg-gray-500",
	}

	for bt, style := range badgeStyleMap {
		hue, ok := hueFor[bt]
		if !ok {
			t.Errorf("BadgeType %s missing from the adapter hue map", bt)

			continue
		}

		if style.Dot != wantDots[bt] {
			t.Errorf(
				"BadgeType %s: dot %q != expected %q (hue %s)",
				bt, style.Dot, wantDots[bt], hue,
			)
		}
	}

	wantTones := map[BadgeType]tone.Tone{
		BadgePrimary: tone.ToneAccent,
		BadgeSuccess: tone.ToneSuccess,
		BadgeWarning: tone.ToneWarning,
		BadgeError:   tone.ToneDanger,
		BadgeInfo:    tone.ToneInfo,
		BadgeNeutral: tone.ToneNeutral,
	}

	for bt, want := range wantTones {
		if got := bt.Tone(); got != want {
			t.Errorf("BadgeType %s.Tone() = %s, want %s", bt, got, want)
		}
	}
}
