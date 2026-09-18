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

// TestErrorMaxWidthIsValid guards the ErrorMaxWidth closed set. Unknown or
// empty values render XL (the pre-field look) via map+fallback.
func TestErrorMaxWidthIsValid(t *testing.T) {
	t.Parallel()

	valid := []ErrorMaxWidth{ErrorMaxWidthLG, ErrorMaxWidthXL, ErrorMaxWidth2XL, ErrorMaxWidth4XL}
	for _, w := range valid {
		if !ErrorMaxWidthIsValid(w) {
			t.Errorf("ErrorMaxWidthIsValid(%q) = false, want true", w)
		}
	}

	invalid := []ErrorMaxWidth{"", "XL", "3xl", "full"}
	for _, w := range invalid {
		if ErrorMaxWidthIsValid(w) {
			t.Errorf("ErrorMaxWidthIsValid(%q) = true, want false", w)
		}

		if got := errorMaxWidthClass(w); got != "max-w-xl" {
			t.Errorf("errorMaxWidthClass(%q) = %q, want XL fallback", w, got)
		}
	}
}

// TestResolvedWayOut pins the WayOutAction dual-read precedence: the typed
// bundle wins ENTIRELY when its Text is set — no field mixing between the
// legacy strings and the bundle.
func TestResolvedWayOut(t *testing.T) {
	t.Parallel()

	legacyOnly := ErrorPageProps{WayOut: "Retry", WayOutHref: "/"}
	if text, href := legacyOnly.resolvedWayOut(); text != "Retry" || href != "/" {
		t.Errorf("legacy fields: resolved = (%q, %q), want (Retry, /)", text, href)
	}

	actionWins := ErrorPageProps{
		WayOut:       "Retry",
		WayOutHref:   "/",
		WayOutAction: WayOutAction{Text: "View status", Href: "/status"},
	}
	if text, href := actionWins.resolvedWayOut(); text != "View status" || href != "/status" {
		t.Errorf("action precedence: resolved = (%q, %q), want (View status, /status) — no mixing", text, href)
	}

	actionNoHref := ErrorPageProps{WayOutAction: WayOutAction{Text: "Go back"}}
	if text, href := actionNoHref.resolvedWayOut(); text != "Go back" || href != "" {
		t.Errorf("action without href: resolved = (%q, %q), want (Go back, empty)", text, href)
	}
}
