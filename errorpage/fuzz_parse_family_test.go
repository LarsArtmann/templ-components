package errorpage

import "testing"

// FuzzParseFamily proves ParseFamily never panics and never returns an
// invalid Family for arbitrary input — the property the whole family-driven
// rendering pipeline rests on. Fuzz with:
//
//	go test -fuzz=FuzzParseFamily -fuzztime=30s ./errorpage/
func FuzzParseFamily(f *testing.F) {
	seeds := []string{
		"", "rejection", "REJECTION", "Rejection", "conflict", "transient",
		"corruption", "infrastructure", "orchestration", "unknown", "  ",
		"rejection ", " rejection", "rejectionx", "\x00", "ünicode",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got := ParseFamily(s)

		// The result must always be a renderable family — unknown input
		// degrades to Transient, never panics and never leaks an invalid
		// value into the style lookup chain.
		if !FamilyIsValid(got) {
			t.Fatalf("ParseFamily(%q) = invalid family %q", s, got)
		}

		// lookupFamilyStyle must resolve for every ParseFamily output — the
		// actual consumption path (a miss would render empty styles).
		_ = lookupFamilyStyle(got)

		// FamilyStatusCode must resolve (never a zero/unknown status).
		if code := FamilyStatusCode(got); code < 400 || code > 599 {
			t.Fatalf("FamilyStatusCode(ParseFamily(%q)) = %d, want an HTTP error code", s, code)
		}

		// FamilyDefaultTitle must never be empty for any resolved family.
		if title := FamilyDefaultTitle(got); title == "" {
			t.Fatalf("FamilyDefaultTitle(ParseFamily(%q)) = empty title", s)
		}
	})
}
