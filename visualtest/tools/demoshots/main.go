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

const settle = 600 * time.Millisecond

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
	browserCtx, cancelBrowserTimeout := context.WithTimeout(browserCtx, 30*time.Minute)
	defer cancelBrowserTimeout()

	if err := chromedp.Run(browserCtx); err != nil {
		log.Fatalf("start browser: %v", err)
	}

	for _, p := range pages {
		if err := capturePage(browserCtx, *base, *out, p); err != nil {
			log.Fatalf("capture %s: %v", p.name, err)
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

// capturePage screenshots one page in light mode, then toggles the dark
// class on the same tab and screenshots again (no second navigation).
func capturePage(parent context.Context, base, out string, p page) error {
	ctx, cancel := chromedp.NewContext(parent)
	defer cancel()
	ctx, cancelTimeout := context.WithTimeout(ctx, 180*time.Second)
	defer cancelTimeout()

	var light, dark []byte

	start := time.Now()

	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 900),
		chromedp.Navigate(base+p.path),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(settle),
		chromedp.FullScreenshot(&light, 92),
		chromedp.Evaluate(`document.documentElement.classList.add('dark');`, nil),
		chromedp.Sleep(2*settle),
		chromedp.FullScreenshot(&dark, 92),
	)
	if err != nil {
		return fmt.Errorf("after %s: %w", time.Since(start).Round(time.Second), err)
	}

	if err := os.WriteFile(filepath.Join(out, p.name+"_light.png"), light, 0o644); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(out, p.name+"_dark.png"), dark, 0o644)
}
