package visualtest

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// goldenMaxMismatchPercent is the shared pixel-mismatch tolerance (0.1%).
const goldenMaxMismatchPercent = 0.1

// update is the shared -update flag. When set, golden PNGs are rewritten with
// the current render instead of being compared. Mirrors utils/golden's DX.
//
//nolint:gochecknoglobals // CLI flag for golden updates
var update = flag.Bool("update", false, "update visual golden PNGs instead of comparing")

// updateGoldens reports whether goldens should be rewritten. The nil check
// keeps embedded/ad-hoc test mains from panicking if the flag was never
// registered.
func updateGoldens() bool { return update != nil && *update }

// goldenDir is where golden screenshots live, relative to the test package.
const goldenDir = "testdata"

// readGolden decodes the golden PNG for name. Returns nil, false if no golden
// exists yet (first run).
func readGolden(t *testing.T, name string) (image.Image, bool) {
	t.Helper()

	p := goldenPath(name)

	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}

		t.Fatalf("visualtest: open golden %s: %v", p, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("visualtest: decode golden %s: %v", p, err)
	}

	return img, true
}

// writeGolden encodes actual as the new golden PNG for name.
func writeGolden(t *testing.T, name string, actual []byte) {
	t.Helper()

	p := goldenPath(name)
	if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
		t.Fatalf("visualtest: mkdir golden dir: %v", err)
	}

	if err := os.WriteFile(p, actual, 0o600); err != nil {
		t.Fatalf("visualtest: write golden %s: %v", p, err)
	}

	t.Logf("visualtest: updated golden %s", p)
}

// failureRunDirsKept bounds how many past failure-run directories survive in
// testdata/.fail/ — enough history to compare two failed runs, few enough that
// the dir never grows unbounded.
const failureRunDirsKept = 10

// failureRunDir is THIS test process's artifact directory:
// testdata/.fail/<UTC timestamp>-<pid>/. A second run therefore can never wipe
// the evidence of a first one mid-investigation (the 2026-10-30 session lost
// exactly that evidence when a re-run started while diffs were being read).
// Directories past failureRunDirsKept are pruned on creation.
var failureRunDir = sync.OnceValue(func() string {
	name := time.Now().UTC().Format("20060102-150405") + "-" + strconv.Itoa(os.Getpid())
	dir := filepath.Join(goldenDir, ".fail", name)
	pruneOldFailureRunDirs(filepath.Join(goldenDir, ".fail"), name)

	return dir
})

// pruneOldFailureRunDirs keeps only the newest failureRunDirsKept run dirs
// (name included). Best-effort: prune errors are non-fatal by design.
func pruneOldFailureRunDirs(failDir, keep string) {
	entries, err := os.ReadDir(failDir)
	if err != nil {
		return
	}

	var runs []string
	for _, e := range entries {
		if e.IsDir() && isFailureRunDirName(e.Name()) {
			runs = append(runs, e.Name())
		}
	}

	sort.Strings(runs) // timestamp prefix == lexicographic age

	excess := len(runs) - failureRunDirsKept
	for i := range excess {
		_ = os.RemoveAll(filepath.Join(failDir, runs[i]))
	}
}

// isFailureRunDirName matches the <YYYYMMDD-HHMMSS>-<pid> run-dir shape.
func isFailureRunDirName(name string) bool {
	dash := strings.LastIndexByte(name, '-')
	if dash <= 0 {
		return false
	}

	pid := name[dash+1:]
	if pid == "" {
		return false
	}

	for _, r := range pid {
		if r < '0' || r > '9' {
			return false
		}
	}

	stamp := name[:dash]
	_, err := time.Parse("20060102-150405", stamp)

	return err == nil
}

// writeFailureArtifacts saves the actual PNG and a red-pixel diff image into
// this run's failure directory so a human can inspect what changed. Lives
// under testdata/.fail/<run>/ (per-run — see failureRunDir).
func writeFailureArtifacts(t *testing.T, name string, actual []byte, diff *image.RGBA) {
	t.Helper()

	if len(actual) > 0 {
		p := filepath.Join(failureRunDir(), name+".actual.png")

		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			t.Logf("visualtest: create fail dir %s: %v", filepath.Dir(p), err)

			return
		}

		if err := os.WriteFile(p, actual, 0o600); err != nil {
			t.Logf("visualtest: write actual %s: %v", p, err)
		}

		t.Logf("visualtest: wrote actual render to %s", p)
	}

	if diff != nil {
		p := filepath.Join(failureRunDir(), name+".diff.png")
		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			t.Logf("visualtest: create fail dir %s: %v", filepath.Dir(p), err)

			return
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, diff); err != nil {
			t.Logf("visualtest: encode diff %s: %v", p, err)

			return
		}

		if err := os.WriteFile(p, buf.Bytes(), 0o600); err != nil {
			t.Logf("visualtest: write diff %s: %v", p, err)
		}

		t.Logf("visualtest: wrote pixel diff to %s", p)
	}
}

func goldenPath(name string) string {
	return filepath.Join(goldenDir, name+".png")
}

// cleanFailureArtifacts removes this run's actual/diff images once the test
// passes again. Keeps testdata/.fail/ honest — only genuinely-failing tests
// leave artifacts behind, and past runs' evidence is never touched.
func cleanFailureArtifacts(name string) {
	for _, suffix := range []string{".actual.png", ".diff.png"} {
		_ = os.Remove(filepath.Join(failureRunDir(), name+suffix))
	}
}

// failureArtifactPath names this run's diff artifact for a failing golden, for
// use in human-facing error messages (the .actual.png sits beside it).
func failureArtifactPath(name string) string {
	return filepath.Join(failureRunDir(), name+".diff.png")
}

// assertGoldenMatch decodes actual and pixel-compares it against the golden
// (0.1% tolerance), writing failure artifacts on mismatch — the shared tail
// of every golden-capturing test.
func assertGoldenMatch(t *testing.T, name string, golden image.Image, actual []byte) {
	t.Helper()

	decoded, err := png.Decode(bytes.NewReader(actual))
	if err != nil {
		t.Fatalf("visualtest[%s]: decode actual: %v", name, err)
	}

	result, diff := comparePixels(golden, decoded, goldenMaxMismatchPercent, goldenMaxMismatchPercent*percentMultiplier)
	if !result.Match {
		writeFailureArtifacts(t, name, actual, diff)
		t.Errorf("visualtest[%s]: visual mismatch — %s (max %.4f%%).\n"+
			"Inspect testdata/.fail/%s.{actual,diff}.png, then run `go test -update` if the change is intended.",
			name, result, goldenMaxMismatchPercent, name)
	}
}
