package tone_test

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

// TestValidityRejectsUnknownValues pins the closed-set contract: values that
// are not part of the enum must be reported invalid (consumers gate lookups
// and v2 migrations on these).
func TestValidityRejectsUnknownValues(t *testing.T) {
	t.Parallel()

	if tone.HueIsValid(tone.Hue("not-a-hue")) {
		t.Error(`HueIsValid("not-a-hue") = true, want false`)
	}

	if tone.ToneIsValid(tone.Tone("not-a-tone")) {
		t.Error(`ToneIsValid("not-a-tone") = true, want false`)
	}
}
