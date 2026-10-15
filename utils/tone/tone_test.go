package tone_test

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

// TestToneTableExhaustive: every Tone in Tones() must carry a complete
// Surfaces entry (all seven pairs, every dark half prefixed dark:), and
// unknown tones must fall back to Neutral.
func TestToneTableExhaustive(t *testing.T) {
	t.Parallel()

	for _, tn := range tone.Tones() {
		if !tone.ToneIsValid(tn) {
			t.Errorf("%s: ToneIsValid false for a member of Tones()", tn)
		}

		s := tone.For(tn)

		pairs := map[string]tone.Pair{
			"Dot": s.Dot, "Text": s.Text, "TextSoft": s.TextSoft,
			"Border": s.Border, "BGTint": s.BGTint, "BGSolid": s.BGSolid, "Icon": s.Icon,
		}

		for name, p := range pairs {
			if !p.Complete() {
				t.Errorf("tone %s: surface %s is incomplete (Light=%q Dark=%q)", tn, name, p.Light, p.Dark)

				continue
			}

			if !strings.HasPrefix(p.Dark, "dark:") {
				t.Errorf(
					"tone %s: surface %s dark half %q lacks the dark: prefix (scanner-safety contract)",
					tn,
					name,
					p.Dark,
				)
			}

			if got := p.Classes(); got != p.Light+" "+p.Dark {
				t.Errorf("tone %s: surface %s Classes() = %q, want %q", tn, name, got, p.Light+" "+p.Dark)
			}
		}
	}

	if got := tone.For(tone.Tone("bogus")); got != tone.For(tone.ToneNeutral) {
		t.Error("For(unknown) does not fall back to Neutral")
	}
}

// TestHueTableExhaustive: every Hue in Hues() must carry a complete
// Utilities entry, and unknown hues must fall back to Gray.
func TestHueTableExhaustive(t *testing.T) {
	t.Parallel()

	for _, h := range tone.Hues() {
		if !tone.HueIsValid(h) {
			t.Errorf("%s: HueIsValid false for a member of Hues()", h)
		}

		u := tone.UtilitiesFor(h)

		pairs := map[string]tone.Pair{
			"Text": u.Text, "TextMuted": u.TextMuted, "TextStrong": u.TextStrong, "TextOnTint": u.TextOnTint,
			"BGSolid": u.BGSolid, "BGSolidStrong": u.BGSolidStrong,
			"BGTint": u.BGTint, "Border": u.Border,
		}

		for name, p := range pairs {
			if !p.Complete() {
				t.Errorf("hue %s: utility %s is incomplete (Light=%q Dark=%q)", h, name, p.Light, p.Dark)

				continue
			}

			if !strings.HasPrefix(p.Dark, "dark:") {
				t.Errorf("hue %s: utility %s dark half %q lacks the dark: prefix", h, name, p.Dark)
			}
		}
	}

	if got := tone.UtilitiesFor(tone.Hue("bogus")); got != tone.UtilitiesFor(tone.HueGray) {
		t.Error("UtilitiesFor(unknown) does not fall back to Gray")
	}
}

// TestHueUtilityAnchors pins the cells that components render today,
// byte-for-byte, so the grammar cannot drift from shipped pixels. Anchors:
// chart series text + segment bars (chart_geometry/segment_bar), status and
// kanban dots (status_dot/kanban), feedback alert surfaces (styles.go),
// stat tiles (card.templ), trend arrows (card.templ), log tags
// (scrollback.templ).
func TestHueUtilityAnchors(t *testing.T) {
	t.Parallel()

	anchors := []struct {
		name string
		got  string
		want string
	}{
		{"chart blue text", tone.UtilitiesFor(tone.HueBlue).Text.Classes(), "text-blue-600 dark:text-blue-400"},
		{"chart rose text", tone.UtilitiesFor(tone.HueRose).Text.Classes(), "text-rose-600 dark:text-rose-400"},
		{
			"segment violet fill",
			tone.UtilitiesFor(tone.HueViolet).BGSolidStrong.Classes(),
			"bg-violet-600 dark:bg-violet-400",
		},
		{"status dot success", tone.UtilitiesFor(tone.HueGreen).BGSolid.Classes(), "bg-green-500 dark:bg-green-400"},
		{
			"status dot warning (amber dialect)",
			tone.UtilitiesFor(tone.HueAmber).BGSolid.Classes(),
			"bg-amber-500 dark:bg-amber-400",
		},
		{
			"status dot neutral (gray exception)",
			tone.UtilitiesFor(tone.HueGray).BGSolid.Classes(),
			"bg-gray-400 dark:bg-gray-500",
		},
		{"kanban purple dot", tone.UtilitiesFor(tone.HuePurple).BGSolid.Classes(), "bg-purple-500 dark:bg-purple-400"},
		{
			"trend up (green strong)",
			tone.UtilitiesFor(tone.HueGreen).TextStrong.Classes(),
			"text-green-700 dark:text-green-400",
		},
		{
			"trend warn (amber strong)",
			tone.UtilitiesFor(tone.HueAmber).TextStrong.Classes(),
			"text-amber-700 dark:text-amber-400",
		},
		{
			"stat tile purple tint",
			tone.UtilitiesFor(tone.HuePurple).BGTint.Classes(),
			"bg-purple-50 dark:bg-purple-900/30",
		},
		{
			"feedback success border",
			tone.UtilitiesFor(tone.HueGreen).Border.Classes(),
			"border-green-200 dark:border-green-800",
		},
		{"scrollback info tag", tone.UtilitiesFor(tone.HueBlue).Text.Classes(), "text-blue-600 dark:text-blue-400"},
		{
			"scrollback neutral tag (muted convention)",
			tone.UtilitiesFor(tone.HueGray).TextMuted.Classes(),
			"text-gray-500 dark:text-gray-400",
		},
	}

	for _, a := range anchors {
		if a.got != a.want {
			t.Errorf("%s = %q, want %q", a.name, a.got, a.want)
		}
	}
}

// TestChartHuesOrder pins the canonical series cycle shared by line, pie,
// and segment charts (chart_geometry.go order).
func TestChartHuesOrder(t *testing.T) {
	t.Parallel()

	want := []tone.Hue{
		tone.HueBlue, tone.HueEmerald, tone.HueAmber, tone.HueRose,
		tone.HueViolet, tone.HueCyan, tone.HueOrange, tone.HuePink,
	}

	got := tone.ChartHues()
	if len(got) != len(want) {
		t.Fatalf("ChartHues() has %d entries, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ChartHues()[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

// TestPairIncompleteAndJoin exercises the Pair edges: empty halves join
// without dangling spaces, Complete reports honestly.
func TestPairIncompleteAndJoin(t *testing.T) {
	t.Parallel()

	lightOnly := tone.Pair{Light: "text-gray-500"}
	if lightOnly.Complete() {
		t.Error("light-only Pair reports Complete")
	}

	if got := lightOnly.Classes(); got != "text-gray-500" {
		t.Errorf("light-only Classes() = %q, want %q", got, "text-gray-500")
	}

	empty := tone.Pair{}
	if got := empty.Classes(); got != "" {
		t.Errorf("empty Pair Classes() = %q, want empty", got)
	}
}
