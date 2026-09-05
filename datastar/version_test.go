package datastar

import (
	"testing"

	"github.com/larsartmann/go-datastar/static"
)

// TestDatastarVersionMatchesStatic verifies that DatastarVersion1_0_3 is derived
// from static.Version at compile time. This drift-guard ensures the CDN URL and
// the embedded JS bundle can never diverge — if static.Version is bumped,
// DatastarVersion1_0_3 automatically follows.
func TestDatastarVersionMatchesStatic(t *testing.T) {
	t.Parallel()

	if string(DatastarVersion1_0_3) != static.Version {
		t.Errorf(
			"DatastarVersion1_0_3 (%q) does not match static.Version (%q) — "+
				"the version constant should be derived via DatastarVersion(static.Version)",
			DatastarVersion1_0_3,
			static.Version,
		)
	}
}
