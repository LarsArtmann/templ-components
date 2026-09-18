package errorpage

import "testing"

// TestErrorDetailVariantIsValid guards the ErrorDetailVariant closed set.
// The drift guard requires every IsValid to ship with a test in the same
// commit. Unknown values are intentionally "valid = false" while RENDERING
// falls back to Tinted (map+fallback convention, zero-value compatibility).
func TestErrorDetailVariantIsValid(t *testing.T) {
	t.Parallel()

	valid := []ErrorDetailVariant{ErrorDetailTinted, ErrorDetailNeutral}
	for _, v := range valid {
		if !ErrorDetailVariantIsValid(v) {
			t.Errorf("ErrorDetailVariantIsValid(%q) = false, want true", v)
		}
	}

	invalid := []ErrorDetailVariant{"", "Tinted", "NEUTRAL", "solid", "accent"}
	for _, v := range invalid {
		if ErrorDetailVariantIsValid(v) {
			t.Errorf("ErrorDetailVariantIsValid(%q) = true, want false", v)
		}
	}
}
