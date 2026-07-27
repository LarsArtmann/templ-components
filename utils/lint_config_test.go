package utils

import (
	"os"
	"strings"
	"testing"
)

// TestGolangciDisabledLinters prevents the recurring regression where
// godoclint, ireturn, and testableexamples re-enter the golangci-lint enable
// list. These three linters are fundamentally incompatible with a templ
// library (see AGENTS.md "Disabled linters") and have re-appeared three times
// across releases, each time sending `golangci-lint run` from 0 to 71 findings.
// This test fails the moment one of them shows up as an enabled linter again.
func TestGolangciDisabledLinters(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../.golangci.yml")
	if err != nil {
		t.Fatalf("read .golangci.yml: %v", err)
	}

	for _, name := range []string{"godoclint", "ireturn", "testableexamples"} {
		for line := range strings.SplitSeq(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			// A YAML enabled-linter entry looks like "- ireturn". Match the
			// list-item form so a name appearing only in a comment or doc
			// string does not trip the guard.
			if strings.EqualFold(trimmed, "- "+name) {
				t.Errorf(".golangci.yml re-enables disabled linter %q — a documented regression (AGENTS.md). "+
					"Remove it from the enable list.", name)
			}
		}
	}

	// The dead `ireturn:` settings block (allow: error/empty/anon/stdlib/generic)
	// must also stay gone — it has no effect once ireturn is disabled and only
	// invites someone to re-enable ireturn assuming the config is tuned for it.
	if strings.Contains(string(data), "ireturn:") {
		t.Errorf(".golangci.yml still has an `ireturn:` settings block — delete it; ireturn is a disabled linter.")
	}
}
