package visualtest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestMain ensures the shared Chromium process is cleaned up after all tests
// complete. The browser is lazily initialized by ensureAllocator on the first
// AssertScreenshot call; this teardown runs regardless of whether any visual
// test actually ran.
func TestMain(m *testing.M) {
	pruneStaleFailureArtifacts()

	code := m.Run()

	ShutdownBrowser()

	os.Exit(code)
}

// pruneStaleFailureArtifacts empties testdata/.fail/ at the START of a run
// (TODO #202). Failure artifacts are per-run debugging evidence: the per-test
// cleanFailureArtifacts path only removes files for tests that re-ran and
// passed, so renamed goldens, aborted runs, and cross-session residue
// (including empty subdirectories) accumulated indefinitely and misled
// debugging sessions into inspecting artifacts from a previous run. Pruning
// up front guarantees everything under .fail/ after a run came from THIS run.
// Only this package writes goldens, so there is no concurrent-writer race.
func pruneStaleFailureArtifacts() {
	failDir := filepath.Join(goldenDir, ".fail")

	entries, err := os.ReadDir(failDir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		fmt.Fprintf(os.Stderr, "visualtest: read stale .fail dir %s: %v\n", failDir, err)

		return
	}

	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(failDir, entry.Name())); err != nil {
			fmt.Fprintf(os.Stderr, "visualtest: prune stale .fail entry %s: %v\n", entry.Name(), err)
		}
	}
}
