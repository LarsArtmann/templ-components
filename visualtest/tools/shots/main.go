// Command shots captures full-page screenshots of a running templ-components
// demo server for manual visual inspection. It is the sanctioned replacement
// for ad-hoc audit tooling: light + dark captures per route, fresh browser
// per page (a long-lived shared browser degrades across very tall captures
// and can hang for minutes), one JPEG per route at a size image viewers accept.
//
// Usage (or `nix run .#shots` which pins Chromium and the toolchain):
//
//	GOWORK=off CHROMEDP_CHROME_PATH=/path/to/chromium go run ./tools/shots \
//	  -base http://localhost:8901 -out /tmp/tc-shots [-mode both] [-width 1280] [-page index]
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/visualtest/internal/browser"
)

type page struct {
	name string
	path string
}

//nolint:gochecknoglobals // declarative page list; a package-level table is the point
var pages = []page{
	{"index", "/"},
	// Every page of the multi-page demo gets a capture; the wire page also
	// renders per ?transport= (audit f14), so both single-transport variants
	// are captured for eyeballing.
	{"layout", "/layout"},
	{"display", "/display"},
	{"feedback", "/feedback"},
	{"forms", "/forms"},
	{"navigation", "/navigation"},
	{"icons", "/icons"},
	{"htmx", "/htmx"},
	{"datastar", "/datastar"},
	{"wire", "/wire"},
	{"wire-htmx", "/wire?transport=htmx"},
	{"wire-datastar", "/wire?transport=datastar"},
	{"kanban", "/kanban"},
	{"echarts", "/echarts"},
	{"error-pages", "/error-pages"},
	{"recipes", "/recipes"},
	{"users", "/users"},
	{"recipes-dashboard", "/recipes/dashboard"},
	{"recipes-settings", "/recipes/settings"},
	{"recipes-login", "/recipes/login"},
	{"recipes-auth", "/recipes/auth"},
}

const (
	settle = 600 * time.Millisecond

	// shotsViewportHeight pairs with the -width flag: full-page captures
	// grow past it, the height only seeds the initial viewport.
	shotsViewportHeight = 900

	// defaultShotsWidth is the desktop capture width matching the visual
	// suite's viewportDesktopWidth.
	defaultShotsWidth = 1280

	// shotsTimeout bounds one page's browser session (navigate + settle +
	// captures + dark-class toggle) before the context is cancelled.
	shotsTimeout = 120 * time.Second

	// screenshotQuality is the image quality passed to FullScreenshot.
	screenshotQuality = 92

	// Capture modes; the dark mode toggles the page's .dark class in place.
	modeLight = "light"
	modeDark  = "dark"
)

func main() {
	base := flag.String("base", "http://localhost:8901", "demo server base URL")
	out := flag.String("out", "/tmp/tc-shots", "output directory")
	mode := flag.String("mode", "both", "light | dark | both")
	width := flag.Int("width", defaultShotsWidth, "viewport width")
	only := flag.String("page", "", "capture a single page by name (e.g. index)")
	selftest := flag.Bool("selftest", false, "verify allocator + base-URL reachability, then exit")

	flag.Parse()

	if *selftest {
		if err := browser.RequireChromium(); err != nil {
			log.Fatalf("selftest FAIL: %v", err)
		}

		resp, err := browser.HTTPGetSelftest(*base + "/health")

		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("selftest FAIL: demo %s: %v", *base, err)
		}

		_ = resp.Body.Close()

		browser.OK(os.Stdout, "shots", "allocator + demo health")

		return
	}

	if err := os.MkdirAll(*out, 0o750); err != nil {
		log.Fatal(err)
	}

	modes := []string{modeLight, modeDark}

	switch *mode {
	case modeLight:
		modes = []string{modeLight}
	case modeDark:
		modes = []string{modeDark}
	}

	for _, p := range pages {
		if *only != "" && p.name != *only {
			continue
		}

		// Fresh browser per page: a long-lived shared browser degrades across
		// very tall captures (multi-minute hangs); an isolated instance
		// captures each page in ~2s.
		if err := capturePage(browser.ExecPath(), *base, *out, p, modes, *width); err != nil {
			log.Fatalf("capture %s: %v", p.name, err)
		}

		fmt.Fprintf(os.Stdout, "captured %s\n", p.name)
	}

	fmt.Fprintln(os.Stdout, "done")
}

// capturePage screenshots one page per requested mode on a single tab,
// toggling the dark class in place between captures (no re-navigation).
// A main-frame response with status >= 400 is an error — without this check
// the tool happily captures the server's error page as a "golden" route
// capture (it captured 404 pages without complaint until 2026-09-08).
func capturePage(execPath, base, out string, p page, modes []string, width int) error {
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(execPath),
			chromedp.Flag("window-size", fmt.Sprintf("%d,%d", width, shotsViewportHeight)),
			chromedp.Flag("force-device-scale-factor", "1"),
		)...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	ctx, cancelTimeout := context.WithTimeout(browserCtx, shotsTimeout)
	defer cancelTimeout()

	var docStatus int64

	start := time.Now()

	// RunResponse returns the main-frame document response of the navigation in
	// the steps — the v0.20 replacement for the old ListenTarget + network.Enable
	// docStatus handshake. resp is non-nil whenever the navigation itself
	// succeeded, even if a later step failed, so error messages keep the status.
	resp, runErr := chromedp.RunResponse(ctx, captureTasks(base, p.path, width)...)
	if resp != nil {
		docStatus = resp.Status
	}

	if runErr != nil {
		return fmt.Errorf(
			"chromium (execPath=%s, docStatus=%d) after %s: %w",
			execPath,
			docStatus,
			time.Since(start).Round(time.Second),
			runErr,
		)
	}

	if err := rejectErrorPage(p, execPath, docStatus); err != nil {
		return fmt.Errorf("capture %s (execPath=%s, docStatus=%d): %w", p.path, execPath, docStatus, err)
	}

	light, dark, err := captureModeShots(ctx, modes)
	if err != nil {
		return fmt.Errorf("capture (execPath=%s, docStatus=%d): %w", execPath, docStatus, err)
	}

	for _, mode := range modes {
		buf := light
		if mode == modeDark {
			buf = dark
		}

		if err := os.WriteFile(filepath.Join(out, p.name+"_"+mode+".png"), buf, 0o600); err != nil {
			return fmt.Errorf("write capture (execPath=%s, docStatus=%d): %w", execPath, docStatus, err)
		}
	}

	return nil
}

// captureModeShots screenshots the settled page in light mode and — when
// dark is among the modes — toggles the class in place (no re-navigation),
// waits the double settle, and screenshots again.
func captureModeShots(ctx context.Context, modes []string) ([]byte, []byte, error) {
	light, err := chromedp.Run(ctx, chromedp.FullScreenshot(screenshotQuality))
	if err != nil {
		return nil, nil, fmt.Errorf("capture light: %w", err)
	}

	if !slices.Contains(modes, modeDark) {
		return light, nil, nil
	}

	if err := chromedp.Do(ctx,
		chromedp.Evaluate[chromedp.Void](`document.documentElement.classList.add('dark');`),
		chromedp.Sleep(2*settle),
	); err != nil {
		return nil, nil, fmt.Errorf("toggle dark: %w", err)
	}

	dark, err := chromedp.Run(ctx, chromedp.FullScreenshot(screenshotQuality))
	if err != nil {
		return nil, nil, fmt.Errorf("capture dark: %w", err)
	}

	return light, dark, nil
}

// captureTasks builds the chromedp action sequence: viewport, navigation (its
// main-frame response is returned by RunResponse), then wait for layout to
// settle. Dark-mode toggling and the screenshots happen in capturePage so the
// PNG bytes come from the typed FullScreenshot action.
func captureTasks(base, path string, width int) []chromedp.Action[chromedp.Void] {
	return []chromedp.Action[chromedp.Void]{
		chromedp.EmulateViewport(int64(width), shotsViewportHeight),
		chromedp.Navigate(base + path),
		chromedp.WaitReady(chromedp.CSS("body")),
		chromedp.Sleep(settle),
	}
}

// rejectErrorPage refuses to capture when the main-frame response was an HTTP
// error: without this check the tool happily captures the server's error page
// as a "golden" route capture (it captured 404 pages without complaint until
// 2026-09-08).
func rejectErrorPage(p page, execPath string, docStatus int64) error {
	if docStatus < 400 {
		return nil
	}

	return fmt.Errorf( //nolint:err113 // terminal status message; there is nothing to wrap
		"route %s (execPath=%s) returned HTTP %d — refusing to capture an error page (route list stale?)",
		p.path,
		execPath,
		docStatus,
	)
}
