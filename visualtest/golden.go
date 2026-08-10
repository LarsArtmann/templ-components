package visualtest

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// update is the shared -update flag. When set, golden PNGs are rewritten with
// the current render instead of being compared. Mirrors utils/golden's DX.
//
//nolint:gochecknoglobals // CLI flag for golden updates
var update = flag.Bool("update", false, "update visual golden PNGs instead of comparing")

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
	defer f.Close() //nolint:errcheck,gosec // read-only test fixture

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

// writeFailureArtifacts saves the actual PNG and a red-pixel diff image next to
// the golden so a human can inspect what changed. Lives under testdata/.fail/.
func writeFailureArtifacts(t *testing.T, name string, actual []byte, diff *image.RGBA) {
	t.Helper()

	if len(actual) > 0 {
		p := filepath.Join(goldenDir, ".fail", name+".actual.png")

		_ = os.MkdirAll(filepath.Dir(p), 0o750)
		if err := os.WriteFile(p, actual, 0o600); err != nil {
			t.Logf("visualtest: write actual %s: %v", p, err)
		}

		t.Logf("visualtest: wrote actual render to %s", p)
	}

	if diff != nil {
		p := filepath.Join(goldenDir, ".fail", name+".diff.png")
		_ = os.MkdirAll(filepath.Dir(p), 0o750)

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

// cleanFailureArtifacts removes stale actual/diff images from a previous
// failing run once the test passes again. Keeps testdata/.fail/ honest — only
// genuinely-failing tests leave artifacts behind.
func cleanFailureArtifacts(name string) {
	for _, suffix := range []string{".actual.png", ".diff.png"} {
		_ = os.Remove(filepath.Join(goldenDir, ".fail", name+suffix))
	}
}
