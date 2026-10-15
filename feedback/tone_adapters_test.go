package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

// TestToneAdapterEquivalence is the flagship Phase-0 proof: for every
// FeedbackType, the canonical tone table must reproduce the feedback style
// map BYTE-FOR-BYTE — border, background, text, and icon. The feedback map
// is the tone table's anchor dialect; if this test fails, either the table
// or the map drifted, and either way pixels changed.
func TestToneAdapterEquivalence(t *testing.T) {
	t.Parallel()

	for ft, style := range feedbackStyleMap {
		surfaces := tone.For(ft.Tone())

		checks := []struct {
			name    string
			derived string
			shipped string
		}{
			{"Border", surfaces.Border.Classes(), style.Border},
			{"BG", surfaces.BGTint.Classes(), style.BG},
			{"Text", surfaces.Text.Classes(), style.Text},
			{"Icon", surfaces.Icon.Classes(), style.Icon},
		}

		for _, c := range checks {
			if c.derived != c.shipped {
				t.Errorf(
					"FeedbackType %s: tone table %s = %q, feedbackStyleMap ships %q",
					ft, c.name, c.derived, c.shipped,
				)
			}
		}
	}

	// The neutral fallback pair (feedbackStyleDefault) must equal the
	// Neutral tone's surfaces — same proof for the default dialect.
	if feedbackStyleDefault.Border != tone.For(tone.ToneNeutral).Border.Classes() ||
		feedbackStyleDefault.BG != tone.For(tone.ToneNeutral).BGTint.Classes() ||
		feedbackStyleDefault.Text != tone.For(tone.ToneNeutral).Text.Classes() ||
		feedbackStyleDefault.Icon != tone.For(tone.ToneNeutral).Icon.Classes() {
		t.Error(
			"feedbackStyleDefault does not equal tone.For(ToneNeutral) surfaces — " +
				"the neutral fallback dialect drifted from the tone table",
		)
	}
}
