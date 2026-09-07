// Command demoshots captures full-page screenshots of the running demo
// server (light + dark mode) for manual visual inspection.
//
// Usage:
//
//	GOWORK=off CHROMEDP_CHROME_PATH=/path/to/chromium go run ./tools/demoshots \
//	  -base http://localhost:8901 -out /tmp/tc-shots
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type page struct {
	name string
	path string
}

var pages = []page{
	{"index", "/"},
	{"forms", "/forms"},
	{"users", "/users"},
	{"recipes-dashboard", "/recipes/dashboard"},
	{"recipes-settings", "/recipes/settings"},
	{"recipes-login", "/recipes/login"},
	{"recipes-auth", "/recipes/auth"},
}

const settle = 700 * time.Millisecond

func main() {
	base := flag.String("base", "http://localhost:8901", "demo server base URL")
	out := flag.String("out", "/tmp/tc-shots", "output directory for PNGs")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(chromePath()),
			chromedp.Flag("window-size", "1280,900"),
			chromedp.Flag("force-device-scale-factor", "1"),
		)...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()
	browserCtx, cancelBrowserTimeout := context.WithTimeout(browserCtx, 10*time.Minute)
	defer cancelBrowserTimeout()

	if err := chromedp.Run(browserCtx); err != nil {
		log.Fatalf("start browser: %v", err)
	}

	for _, p := range pages {
		for _, mode := range []string{"light", "dark"} {
			if err := capturePage(browserCtx, *base, *out, p, mode); err != nil {
				log.Fatalf("capture %s (%s): %v", p.name, mode, err)
			}
			fmt.Printf("  captured %s_%s\n", p.name, mode)
		}
	}

	fmt.Println("done")
}

func chromePath() string {
	if p := os.Getenv("CHROMEDP_CHROME_PATH"); p != "" {
		return p
	}

	return "chromium"
}

func capturePage(parent context.Context, base, out string, p page, mode string) error {
	ctx, cancel := chromedp.NewContext(parent)
	defer cancel()
	ctx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)
	defer cancelTimeout()

	var png []byte

	tasks := []chromedp.Action{
		chromedp.EmulateViewport(1280, 900),
		chromedp.Navigate(base + p.path),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(settle),
	}

	if mode == "dark" {
		tasks = append(tasks,
			chromedp.Evaluate(`document.documentElement.classList.add('dark'); localStorage.setItem('theme','dark');`, nil),
			chromedp.Sleep(settle),
		)
	}

	tasks = append(tasks, chromedp.Sleep(settle), chromedp.FullScreenshot(&png, 92))

	if err := chromedp.Run(ctx, tasks...); err != nil {
		return err
	}

	name := p.name + "_" + mode + ".png"

	return os.WriteFile(filepath.Join(out, name), png, 0o644)
}

// nolint:unused // kept for reference: section-level capture helper
func captureSection(parent context.Context, base, out, name, sel string) error {
	ctx, cancel := chromedp.NewContext(parent)
	defer cancel()
	ctx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)
	defer cancelTimeout()

	var png []byte

	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 900),
		chromedp.Navigate(base),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Sleep(settle),
		chromedp.Screenshot(sel, &png, chromedp.ByID, chromedp.NodeVisible),
	)
	if err != nil {
		return err
	}

	if strings.TrimSpace(sel) == "" {
		return fmt.Errorf("empty selector")
	}

	return os.WriteFile(filepath.Join(out, "section_"+name+".png"), png, 0o644)
}
