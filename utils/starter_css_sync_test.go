package utils

import (
	"os"
	"testing"
)

// TestStarterCSSMatchesTemplates pins the `tc init` starter files to the
// canonical consumer templates. The starter copies are what `tc init`
// scaffolds; the templates/ files are what the README tells consumers to
// copy. They drifted apart once (the starter predated the `.go`-scanning
// @source lesson and several custom.css sections — TODO_LIST #290,
// 2026-09-23), which silently shipped new consumers a worse entry point.
// The starter files must be byte-identical to the canonical ones: fix
// drift by copying templates/ over starter/, never by editing starter/.
func TestStarterCSSMatchesTemplates(t *testing.T) {
	t.Parallel()

	pairs := [][2]string{
		{"../cmd/tc/_sources/starter/app.css", "../templates/app.css"},
		{"../cmd/tc/_sources/starter/custom.css", "../templates/custom.css"},
	}

	for _, pair := range pairs {
		starter, err := os.ReadFile(pair[0])
		if err != nil {
			t.Fatalf("read %s: %v", pair[0], err)
		}

		canonical, err := os.ReadFile(pair[1])
		if err != nil {
			t.Fatalf("read %s: %v", pair[1], err)
		}

		if string(starter) != string(canonical) {
			t.Errorf(
				"%s drifted from %s — copy the canonical file over the starter (never hand-edit the starter)",
				pair[0],
				pair[1],
			)
		}
	}
}
