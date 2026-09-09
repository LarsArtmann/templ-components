package visualtest

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// The route-level goldens (#163) pin the ASSEMBLED demo pages — the tier that
// would have caught the recipes-dashboard collapse (component goldens were
// all green while the composed page was broken). The demo binary is built and
// served once for the whole suite; each route is captured full-page (the
// index page, ~29k px tall, is captured above-the-fold instead).
//
// Run via `nix run .#visual`; update goldens with -update.

//nolint:gochecknoglobals // process-wide demo server shared by all route tests
var demoRouteBase = sync.OnceValue(func() string {
	binDir, err := os.MkdirTemp("", "tc-route-golden-")
	if err != nil {
		return fmt.Sprintf("TMPDIR-FAILED: %v", err)
	}

	bin := filepath.Join(binDir, "tc-demo-route")

	build := exec.CommandContext(context.Background(), "go", "build", "-o", bin, "./examples/demo")
	build.Dir = ".."
	build.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2", "GOWORK=off")

	if out, err := build.CombinedOutput(); err != nil {
		return fmt.Sprintf("BUILD-FAILED: %v: %s", err, out)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Sprintf("LISTEN-FAILED: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	server := exec.Command(bin)
	server.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	if err := server.Start(); err != nil {
		return fmt.Sprintf("START-FAILED: %v", err)
	}

	base := fmt.Sprintf("http://127.0.0.1:%d", port)

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(base + "/health") //nolint:noctx // one-shot readiness probe with explicit deadline
		if err == nil {
			_ = resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return base
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return "READY-TIMEOUT"
})

// assertRouteScreenshot navigates to url, optionally flips dark mode, waits
// for finite animations to settle, captures (full-page or viewport), and
// compares against the routes/<name> golden — mirroring AssertScreenshot's
// compare/update flow for server-rendered pages.
func assertRouteScreenshot(t *testing.T, name, url string, dark, fullPage bool) {
	t.Helper()

	ctx, cancel := newTab(t)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 40*time.Second)
	defer cancelTimeout()

	var shot []byte

	tasks := []chromedp.Action{
		chromedp.EmulateViewport(int64(ViewportDesktop.Width), int64(ViewportDesktop.Height)),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	}

	if dark {
		tasks = append(tasks,
			chromedp.Evaluate(`document.documentElement.classList.add('dark');`, nil),
			chromedp.Sleep(500*time.Millisecond),
		)
	}

	tasks = append(tasks, waitAnimationsSettled(), chromedp.Sleep(settleDelay))

	if fullPage {
		// quality 100 = PNG (any other value yields JPEG — the golden
		// pipeline decodes PNG; the shots tool deliberately uses JPEG).
		tasks = append(tasks, chromedp.FullScreenshot(&shot, 100))
	} else {
		tasks = append(tasks, chromedp.CaptureScreenshot(&shot))
	}

	if err := chromedp.Run(timeoutCtx, tasks...); err != nil {
		t.Fatalf("route golden[%s]: capture: %v", name, err)
	}

	if *update {
		writeGolden(t, "routes/"+name, shot)

		return
	}

	golden, exists := readGolden(t, "routes/"+name)
	if !exists {
		writeGolden(t, "routes/"+name, shot)
		t.Errorf("route golden[%s]: no golden yet — wrote %s (re-run without -update to verify)", name, goldenPath("routes/"+name))

		return
	}

	actualImg, err := png.Decode(bytes.NewReader(shot))
	if err != nil {
		t.Fatalf("route golden[%s]: decode actual: %v", name, err)
	}

	result, diff := comparePixels(golden, actualImg, defaultOptions(Options{}).Threshold, defaultOptions(Options{}).MaxMismatch*percentMultiplier)
	if !result.Match {
		writeFailureArtifacts(t, "routes/"+name, shot, diff)
		t.Errorf("route golden[%s]: visual mismatch — %s (max %.4f%%).\n"+
			"Inspect testdata/.fail/routes.%s.{actual,diff}.png, then run `nix run .#visual -- -update -run TestDemoRouteGoldens` if intended.",
			name, result, defaultOptions(Options{}).MaxMismatch*percentMultiplier, name)

		return
	}

	cleanFailureArtifacts("routes/" + name)
	t.Logf("route golden[%s]: OK (%.4f%% mismatched)", name, result.MismatchPct)
}

// TestDemoRouteGoldens captures the demo routes end-to-end. Sub-routes are
// full-page; the index page is above-the-fold (its ~29k px full-page capture
// is what `nix run .#shots` is for). Runs sequentially: one shared demo
// server, one capture at a time.
func TestDemoRouteGoldens(t *testing.T) {
	base := demoRouteBase()

	if !strings.HasPrefix(base, "http") {
		t.Fatalf("demo route server did not start: %s", base)
	}

	for _, route := range []struct {
		name     string
		path     string
		dark     bool
		fullPage bool
	}{
		{"dashboard_light", "/recipes/dashboard", false, true},
		{"dashboard_dark", "/recipes/dashboard", true, true},
		{"settings_light", "/recipes/settings", false, true},
		{"login_light", "/recipes/login", false, true},
		{"auth_light", "/recipes/auth", false, true},
		{"forms_light", "/forms", false, true},
		{"users_light", "/users", false, true},
		{"index_fold_light", "/", false, false},
	} {
		assertRouteScreenshot(t, route.name, base+route.path, route.dark, route.fullPage)
	}
}
