package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

// TestToneAdapterEquivalence is the flagship Phase-0 proof, carried through
// the ADR-0049 collapse: for every FeedbackType, the canonical tone table
// must reproduce the classes feedback SHIPPED BYTE-FOR-BYTE — border,
// background, text, and icon. The shipped literals below were the
// feedbackStyleMap anchor dialect before styles.go began deriving from
// tone.For; if this test fails, pixels changed. Semantic severities stay
// palette classes (ADR-0044 decision 3) — the accent family swap only ever
// touched brand-axis sites, never these.
func TestToneAdapterEquivalence(t *testing.T) {
	t.Parallel()

	wantByType := map[FeedbackType][4]string{
		FeedbackSuccess: {
			"border-green-200 dark:border-green-800",
			"bg-green-50 dark:bg-green-900/20",
			"text-green-800 dark:text-green-200",
			"text-green-600 dark:text-green-400",
		},
		FeedbackError: {
			"border-red-200 dark:border-red-800",
			"bg-red-50 dark:bg-red-900/20",
			"text-red-800 dark:text-red-200",
			"text-red-600 dark:text-red-400",
		},
		FeedbackWarning: {
			"border-amber-200 dark:border-amber-800",
			"bg-amber-50 dark:bg-amber-900/20",
			"text-amber-800 dark:text-amber-200",
			"text-amber-600 dark:text-amber-400",
		},
		FeedbackInfo: {
			"border-blue-200 dark:border-blue-800",
			"bg-blue-50 dark:bg-blue-900/20",
			"text-blue-800 dark:text-blue-200",
			"text-blue-600 dark:text-blue-400",
		},
	}

	for ft, want := range wantByType {
		surfaces := tone.For(ft.Tone())

		checks := []struct {
			name    string
			derived string
			shipped string
		}{
			{"Border", surfaces.Border.Classes(), want[0]},
			{"BG", surfaces.BGTint.Classes(), want[1]},
			{"Text", surfaces.Text.Classes(), want[2]},
			{"Icon", surfaces.Icon.Classes(), want[3]},
		}

		for _, c := range checks {
			if c.derived != c.shipped {
				t.Errorf(
					"FeedbackType %s: tone table %s = %q, feedback ships %q",
					ft, c.name, c.derived, c.shipped,
				)
			}
		}
	}

	// The neutral fallback dialect must equal the Neutral tone's surfaces —
	// same byte-for-byte proof for the default.
	neutral := tone.For(tone.ToneNeutral)

	for _, c := range []struct {
		name    string
		derived string
		shipped string
	}{
		{"Border", neutral.Border.Classes(), "border-gray-200 dark:border-gray-700"},
		{"BG", neutral.BGTint.Classes(), "bg-gray-50 dark:bg-gray-800/50"},
		{"Text", neutral.Text.Classes(), "text-gray-800 dark:text-gray-200"},
		{"Icon", neutral.Icon.Classes(), "text-gray-600 dark:text-gray-400"},
	} {
		if c.derived != c.shipped {
			t.Errorf(
				"neutral fallback: tone table %s = %q, feedback ships %q",
				c.name, c.derived, c.shipped,
			)
		}
	}
}
