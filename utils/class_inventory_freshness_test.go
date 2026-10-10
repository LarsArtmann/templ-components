package utils

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

// TestClassInventoryFreshness guards templates/templ-components-classes.txt
// (ADR-0045): the tracked scan corpus module consumers @source must never lag
// the sources it represents. Unlike TestCSSFreshness (informational locally —
// its mtime heuristic is fragile), this guard REGENERATES the inventory with
// scripts/gen-class-inventory.sh (atomic in-place write) and byte-compares
// before/after, so it is exact and always enforced. If the test fails, the
// tree now carries the fresh inventory — commit it alongside the source
// change that caused the drift.
func TestClassInventoryFreshness(t *testing.T) {
	t.Parallel()

	const (
		script = "../scripts/gen-class-inventory.sh"
		out    = "../templates/templ-components-classes.txt"
	)

	if _, err := exec.LookPath("bash"); err != nil {
		t.Skipf("bash not available: %v", err)

		return
	}

	committed, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf(
			"read committed inventory: %v\nIf templates/templ-components-classes.txt was deliberately removed, delete this test and the ADR-0045 release.sh step too.",
			err,
		)
	}

	if err := exec.CommandContext(t.Context(), "bash", script).Run(); err != nil {
		t.Fatalf("run gen-class-inventory.sh: %v", err)
	}

	fresh, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("re-read regenerated inventory: %v", err)
	}

	if !bytes.Equal(bytes.TrimSpace(committed), bytes.TrimSpace(fresh)) {
		t.Errorf(
			"templates/templ-components-classes.txt is STALE — the class inventory did not match " +
				"the library sources. The tree now carries the regenerated file; commit it with the " +
				"source change (release.sh also regenerates it at cut time).",
		)
	}
}
