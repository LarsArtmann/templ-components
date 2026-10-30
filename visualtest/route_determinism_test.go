package visualtest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRouteCaptureIsDeterministic is the BEHAVIOURAL complement to
// TestCapturePathsBlockExternalFonts: the source guard proves the capture
// chain still calls the determinism helpers, this test proves the chain's
// output is actually byte-stable. It captures one rich demo route twice in one
// process (fresh tab each, shared server) and requires identical bytes —
// the property the 2026-10-30 webfont race silently broke (two runs, two
// different glyph sets, ~1% pixel mismatch on a changing subset of routes).
//
// On mismatch both PNGs land in testdata/.fail/ with their SHA-256 in the
// failure message, so the evidence survives for inspection. Run via
// `nix run .#visual -- -run TestRouteCaptureIsDeterministic`.
func TestRouteCaptureIsDeterministic(t *testing.T) {
	base := demoRouteBase()

	if !strings.HasPrefix(base, "http") {
		t.Fatalf("demo route server did not start: %s", base)
	}

	url := base + "/display"

	first := captureRouteScreenshot(t, url, false, false, ViewportDesktop, true)
	second := captureRouteScreenshot(t, url, false, false, ViewportDesktop, true)

	firstSum := sha256.Sum256(first)
	secondSum := sha256.Sum256(second)

	firstHex := hex.EncodeToString(firstSum[:])
	secondHex := hex.EncodeToString(secondSum[:])

	t.Logf("captured /display twice: sha256 %s… vs %s… (%d vs %d bytes)",
		firstHex[:16], secondHex[:16], len(first), len(second))

	if bytes.Equal(first, second) {
		return
	}

	failDir := filepath.Join("testdata", ".fail")

	if err := os.MkdirAll(failDir, 0o755); err != nil {
		t.Fatalf("mkdir failure dir: %v", err)
	}

	for name, shot := range map[string][]byte{
		"determinism_display_first.png":  first,
		"determinism_display_second.png": second,
	} {
		path := filepath.Join(failDir, name)
		if err := os.WriteFile(path, shot, 0o644); err != nil { //nolint:gosec // fixture output, fixed name
			t.Logf("write %s: %v", path, err)
		}
	}

	t.Fatalf("route capture is NOT deterministic: sha256 first=%s second=%s — both PNGs written to testdata/.fail/",
		firstHex, secondHex)
}
