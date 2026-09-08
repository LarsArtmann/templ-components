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
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/chromedp/chromedp"
)

type page struct {
	name string
	path string
}

var pages = []page{
	{"index", "/"},
	// The wire section renders per ?transport= (audit f14): capture both
	// single-transport variants so the htmx/datastar dialects are eyeballed,
	// not just the default both-view.
	{"index-htmx", "/?transport=htmx"},
	{"index-datastar", "/?transport=datastar"},
	{"forms", "/forms"},
	{"users", "/users"},
	{"recipes-dashboard", "/recipes/dashboard"},
	{"recipes-settings", "/recipes/settings"},
	{"recipes-login", "/recipes/login"},
	{"recipes-auth", "/recipes/auth"},
}

const settle = 600 * time.Millisecond

func main() {
	base := flag.String("base", "http://localhost:8901", "demo server base URL")
	out := flag.String("out", "/tmp/tc-shots", "output directory")
	mode := flag.String("mode", "both", "light | dark | both")
	width := flag.Int("width", 1280, "viewport width")
	only := flag.String("page", "", "capture a single page by name (e.g. index)")

	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	modes := []string{"light", "dark"}

	switch *mode {
	case "light":
		modes = []string{"light"}
	case "dark":
		modes = []string{"dark"}
	}

	for _, p := range pages {
		if *only != "" && p.name != *only {
			continue
		}

		// Fresh browser per page: a long-lived shared browser degrades across
		// very tall captures (multi-minute hangs); an isolated instance
		// captures each page in ~2s.
		if err := capturePage(chromePath(), *base, *out, p, modes, *width); err != nil {
			log.Fatalf("capture %s: %v", p.name, err)
		}

		fmt.Printf("captured %s\n", p.name)
	}

	fmt.Println("done")
}

func chromePath() string {
	if p := os.Getenv("CHROMEDP_CHROME_PATH"); p != "" {
		return p
	}

	return "chromium"
}

// capturePage screenshots one page per requested mode on a single tab,
// toggling the dark class in place between captures (no re-navigation).
func capturePage(execPath, base, out string, p page, modes []string, width int) error {
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(execPath),
			chromedp.Flag("window-size", fmt.Sprintf("%d,900", width)),
			chromedp.Flag("force-device-scale-factor", "1"),
		)...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	ctx, cancelTimeout := context.WithTimeout(browserCtx, 120*time.Second)
	defer cancelTimeout()

	var light, dark []byte

	tasks := []chromedp.Action{
		chromedp.EmulateViewport(int64(width), 900),
		chromedp.Navigate(base + p.path),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(settle),
	}

	has := func(m string) bool {
		return slices.Contains(modes, m)
	}

	if has("light") {
		tasks = append(tasks, chromedp.FullScreenshot(&light, 92))
	}

	if has("dark") {
		tasks = append(tasks,
			chromedp.Evaluate(`document.documentElement.classList.add('dark');`, nil),
			chromedp.Sleep(2*settle),
			chromedp.FullScreenshot(&dark, 92),
		)
	}

	start := time.Now()

	if err := chromedp.Run(ctx, tasks...); err != nil {
		return fmt.Errorf("after %s: %w", time.Since(start).Round(time.Second), err)
	}

	for _, mode := range modes {
		buf := light
		if mode == "dark" {
			buf = dark
		}

		if err := os.WriteFile(filepath.Join(out, p.name+"_"+mode+".png"), buf, 0o644); err != nil {
			return err
		}
	}

	return nil
}
