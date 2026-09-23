// Command siteshots captures full-page screenshots of the website's static
// dist/ output for manual visual inspection: light + dark, desktop + mobile,
// one fresh browser per page (a long-lived shared browser degrades across
// very tall captures). It also runs a search smoke check against the header
// docs-search combobox and reports hit counts.
//
// Usage:
//
//	GOWORK=off CHROMEDP_CHROME_PATH=/path/to/chromium go run ./tools/siteshots \
//	  -dist ../website/dist -out /tmp/site-shots
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/visualtest/internal/browser"
	"github.com/larsartmann/templ-components/visualtest/internal/distserver"
)

// routes mirrors the site's key page shapes: landing, sales, docs pages with
// and without code blocks, and the 404.
//
//nolint:gochecknoglobals // declarative route list; a package-level table is the point
var routes = []string{
	"/",
	"/sales",
	"/getting-started/installation",
	"/guides/theming",
	"/guides/dark-mode",
	"/404.html",
}

const (
	settle = 700 * time.Millisecond

	navigationSettle = 300 * time.Millisecond

	viewportHeight = 900
	desktopWidth   = 1440
	mobileWidth    = 390
	mobileHeight   = 844

	screenshotQuality = 90

	routeTimeout = 60 * time.Second
)

// errSearchNoHits is the search-smoke failure sentinel (static by design —
// err113 forbids dynamic error construction).
var errSearchNoHits = errors.New("search smoke FAILED: 0 hits for query")

// scrollRevealJS scrolls through the full document in steps (then returns to
// the top) so IntersectionObserver-driven scroll reveals fire for every
// section before the full-page capture.
const scrollRevealJS = `(() => new Promise((resolve) => {
	let y = 0;
	const step = innerHeight * 0.8;
	const max = document.body.scrollHeight;
	const tick = () => {
		y += step;
		if (y >= max) { scrollTo(0, 0); setTimeout(resolve, 400); return; }
		scrollTo(0, y);
		setTimeout(tick, 120);
	};
	tick();
}))()`

func main() {
	dist := flag.String("dist", "../website/dist", "website dist directory to serve")
	out := flag.String("out", "/tmp/site-shots", "screenshot output directory")
	selftest := flag.Bool("selftest", false, "verify allocator + dist listener come-up, then exit")

	flag.Parse()

	if *selftest {
		if browser.ExecPath() == "" {
			log.Fatal("selftest FAIL: no Chromium (set CHROMEDP_CHROME_PATH)")
		}

		addr, server, err := serveDist(*dist)
		if err != nil {
			log.Fatalf("selftest FAIL: dist listener: %v", err)
		}

		resp, err := httpGetSelftest("http://" + addr + "/sales.html")
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("selftest FAIL: dist fetch: %v", err)
		}

		_ = resp.Body.Close()
		_ = server.Close()

		fmt.Println("siteshots selftest OK (allocator + dist listener)")

		return
	}

	if err := run(*dist, *out); err != nil {
		log.Fatal(err)
	}
}

func run(dist, out string) error {
	if err := os.MkdirAll(out, 0o755); err != nil { //nolint:gosec // CLI-controlled output dir
		return fmt.Errorf("create output dir: %w", err)
	}

	addr, server, err := serveDist(dist)
	if err != nil {
		return err
	}

	defer func() { _ = server.Close() }()

	base := "http://" + addr
	fmt.Fprintf(os.Stdout, "siteshots: serving %s at %s\n", dist, base)

	for _, route := range routes {
		if err := captureRoute(base, route, out); err != nil {
			return fmt.Errorf("capture %s: %w", route, err)
		}
	}

	return searchSmoke(base, out)
}

// serveDist starts a loopback HTTP server with Firebase cleanUrls semantics:
// /foo resolves to foo.html when the exact file is absent.
func serveDist(dist string) (string, *http.Server, error) {
	var netListenConfig net.ListenConfig

	listener, err := netListenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("listen: %w", err)
	}

	server := &http.Server{Handler: distserver.Handler(dist)} //nolint:gosec // loopback-only dev server

	go func() { _ = server.Serve(listener) }()

	return listener.Addr().String(), server, nil
}

// captureRoute takes light+dark × desktop+mobile full-page screenshots of one
// route, each in a fresh browser session.
func captureRoute(base, route, out string) error {
	for _, theme := range []string{"light", "dark"} {
		for _, viewport := range []struct {
			name          string
			width, height int
		}{
			{"desktop", desktopWidth, viewportHeight},
			{"mobile", mobileWidth, mobileHeight},
		} {
			name := routeName(route) + "-" + theme + "-" + viewport.name + ".png"

			if err := screenshot(
				base+route,
				theme,
				viewport.width,
				viewport.height,
				filepath.Join(out, name),
			); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}

			fmt.Fprintln(os.Stdout, "wrote", name)
		}
	}

	return nil
}

func screenshot(url, theme string, width, height int, target string) error {
	ctx, cancel := context.WithTimeout(context.Background(), routeTimeout)
	defer cancel()

	allocCtx, allocCancel := chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(browser.ExecPath()),
			chromedp.Flag("headless", true),
		)...,
	)
	defer allocCancel()

	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()

	tabCtx, tabCancel = context.WithTimeout(tabCtx, routeTimeout)
	defer tabCancel()

	// Pin the theme: headless Chromium reports prefers-color-scheme: dark by
	// default, so the ThemeScript FOUC guard would put every capture in dark
	// mode. Store the theme, reload (ThemeScript re-applies), then enforce the
	// class as a belt-and-suspenders guarantee before the capture.
	setTheme := fmt.Sprintf(
		`(() => { localStorage.setItem('theme', %q);
		document.documentElement.classList.toggle('dark', %q === 'dark'); })()`,
		theme, theme,
	)

	var png []byte

	actions := []chromedp.Action{
		chromedp.EmulateViewport(int64(width), int64(height)),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(navigationSettle),
		chromedp.Evaluate(setTheme, nil),
		chromedp.Reload(),
		chromedp.WaitReady("body"),
		chromedp.Sleep(settle),
		chromedp.Evaluate(setTheme, nil),
		chromedp.Evaluate(scrollRevealJS, nil),
		chromedp.Sleep(settle),
		chromedp.FullScreenshot(&png, screenshotQuality),
	}

	if err := chromedp.Run(tabCtx, actions...); err != nil {
		return err
	}

	//nolint:gosec // CLI-controlled screenshot output
	return os.WriteFile(target, png, 0o644)
}

// searchSmoke types a query into the header docs search and asserts results
// render, capturing the open panel as evidence.
func searchSmoke(base, out string) error {
	ctx, cancel := context.WithTimeout(context.Background(), routeTimeout)
	defer cancel()

	allocCtx, allocCancel := chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(browser.ExecPath()),
			chromedp.Flag("headless", true),
		)...,
	)
	defer allocCancel()

	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()

	tabCtx, tabCancel = context.WithTimeout(tabCtx, routeTimeout)
	defer tabCancel()

	var hits int

	queryJS := `(() => { const el = document.querySelector('#doc-search-input');
		el.focus(); el.value = 'component';
		el.dispatchEvent(new Event('input', {bubbles: true})); })()`

	var png []byte

	err := chromedp.Run(tabCtx,
		chromedp.Navigate(base+"/getting-started/installation"),
		chromedp.WaitReady("body"),
		chromedp.Sleep(settle),
		chromedp.Evaluate(queryJS, nil),
		chromedp.Sleep(2*settle),
		chromedp.Evaluate(`document.querySelectorAll('#doc-search-results .doc-search-hit').length`, &hits),
		chromedp.FullScreenshot(&png, screenshotQuality),
	)
	if err != nil {
		return fmt.Errorf("search smoke: %w", err)
	}

	if hits == 0 {
		return fmt.Errorf("query 'component': %w", errSearchNoHits)
	}

	//nolint:gosec // CLI-controlled screenshot output
	if err := os.WriteFile(out+"/search-smoke.png", png, 0o644); err != nil {
		return fmt.Errorf("write search smoke: %w", err)
	}

	fmt.Fprintf(os.Stdout, "search smoke PASS: %d hits (search-smoke.png)\n", hits)

	return nil
}

func routeName(route string) string {
	name := strings.Trim(route, "/")

	name = strings.ReplaceAll(name, "/", "-")
	if name == "" {
		name = "index"
	}

	return name
}
// httpGetSelftest performs the selftest's loopback GET with context and a
// bounded client — shared shape for all capture tools.
func httpGetSelftest(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	return client.Do(req)
}
