package errorpage

import (
	"fmt"
	"testing"
)

// TestFamilyFromStatus pins the canonical status → family mapping, including
// the collapse decisions the forward map forces (Infrastructure → Transient at
// 503, Orchestration → Corruption at 500) and the documented defaults.
func TestFamilyFromStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status int
		want   Family
		why    string
	}{
		{409, FamilyConflict, "the status IS the conflict definition"},
		{429, FamilyTransient, "retry-later semantics"},
		{503, FamilyTransient, "retry-later semantics; wins over Infrastructure (collapse)"},
		{502, FamilyInfrastructure, "gateway cannot reach upstream"},
		{504, FamilyInfrastructure, "gateway timeout"},
		{500, FamilyCorruption, "unknown 5xx → Corruption (FromError precedent)"},
		{501, FamilyCorruption, "server itself is broken"},
		{505, FamilyCorruption, "server itself is broken"},
		{511, FamilyCorruption, "server itself is broken"},
		{599, FamilyCorruption, "server itself is broken"},
		{400, FamilyRejection, "client-side problem"},
		{401, FamilyRejection, "client-side problem"},
		{403, FamilyRejection, "client-side problem"},
		{404, FamilyRejection, "client-side problem"},
		{410, FamilyRejection, "client-side problem"},
		{422, FamilyRejection, "client-side problem"},
		{428, FamilyRejection, "client-side problem"},
		{499, FamilyRejection, "client-side problem"},
		{200, FamilyRejection, "non-error codes: documented default"},
		{301, FamilyRejection, "non-error codes: documented default"},
		{0, FamilyRejection, "non-error codes: documented default"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.status), func(t *testing.T) {
			t.Parallel()

			got := FamilyFromStatus(tt.status)
			if got != tt.want {
				t.Errorf("FamilyFromStatus(%d) = %q, want %q (%s)", tt.status, got, tt.want, tt.why)
			}
		})
	}
}

// TestFamilyFromStatusSweep asserts the whole 100-599 range maps to a VALID
// family and to the documented class bucket — no code can fall through to an
// invalid or surprising family.
func TestFamilyFromStatusSweep(t *testing.T) {
	t.Parallel()

	for status := 100; status < 600; status++ {
		family := FamilyFromStatus(status)

		if !FamilyIsValid(family) {
			t.Errorf("FamilyFromStatus(%d) = %q: not a valid Family", status, family)
		}

		var wantBucket Family
		switch status {
		case 409:
			wantBucket = FamilyConflict
		case 429, 503:
			wantBucket = FamilyTransient
		case 502, 504:
			wantBucket = FamilyInfrastructure
		default:
			if status >= 500 {
				wantBucket = FamilyCorruption
			} else {
				wantBucket = FamilyRejection
			}
		}

		if family != wantBucket {
			t.Errorf("FamilyFromStatus(%d) = %q, want bucket %q", status, family, wantBucket)
		}
	}
}

// TestFamilyStatusCodeRoundTrip proves the STATUS-level round-trip for every
// family: FamilyStatusCode(FamilyFromStatus(FamilyStatusCode(f))) ==
// FamilyStatusCode(f). The FAMILY-level round-trip is impossible by
// construction — Corruption/Orchestration share 500 and Transient/
// Infrastructure share 503 — and those two collapses are asserted explicitly.
func TestFamilyStatusCodeRoundTrip(t *testing.T) {
	t.Parallel()

	for family := range familyStatusCodeMap {
		status := FamilyStatusCode(family)
		roundTripped := FamilyStatusCode(FamilyFromStatus(status))

		if roundTripped != status {
			t.Errorf(
				"status round-trip broke for %q: %d → %q → %d",
				family, status, FamilyFromStatus(status), roundTripped,
			)
		}
	}

	// The two forced collapses, pinned so a future mapping change is a
	// conscious decision, not an accident.
	if got := FamilyFromStatus(FamilyStatusCode(FamilyInfrastructure)); got != FamilyTransient {
		t.Errorf("Infrastructure collapses to %q at 503, want Transient", got)
	}

	if got := FamilyFromStatus(FamilyStatusCode(FamilyOrchestration)); got != FamilyCorruption {
		t.Errorf("Orchestration collapses to %q at 500, want Corruption", got)
	}
}
