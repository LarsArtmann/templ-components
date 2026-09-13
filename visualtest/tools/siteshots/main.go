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
)

// routes mirrors the site's key page shapes: landing, docs pages with and
// without code blocks, and the 404.
//
//nolint:gochecknoglobals // declarative route list; a package-level table is the point
var routes = []string{
	"/",
	"/getting-started/installation",
	"/guides/theming",
	"/guides/dark-mode",
	"/404.html",
}

const (
	settle = 700 * time.Millisecond

	viewportHeight = 900
	desktopWidth   = 1440
	mobileWidth    = 390
	mobileHeight   = 844

	routeTimeout = 60 * time.Second
)

// chromePath resolves the browser binary: CHROMEDP_CHROME_PATH (set by
// `nix run .#visual`-style wrappers) or "chromium" from PATH.
func chromePath() string {
	if p := os.Getenv("CHROMEDP_CHROME_PATH"); p != "" {
		return p
	}

	return "chromium"
}

func main() {
	dist := flag.String("dist", "../website/dist", "website dist directory to serve")
	out := flag.String("out", "/tmp/site-shots", "screenshot output directory")
	flag.Parse()

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

	defer server.Close()

	base := "http://" + addr
	fmt.Printf("siteshots: serving %s at %s\n", dist, base)

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
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("listen: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clean := strings.SplitN(r.URL.Path, "?", 2)[0]
		if !strings.HasSuffix(clean, ".html") && !strings.Contains(clean, ".") {
			if _, err := os.Stat(filepath.Join(dist, clean+".html")); err == nil {
				http.ServeFile(w, r, filepath.Join(dist, clean+".html")) //nolint:gosec // CLI-controlled dist root

				return
			}
		}

		http.FileServer(http.Dir(dist)).ServeHTTP(w, r) //nolint:gosec // CLI-controlled dist root
	})

	server := &http.Server{Handler: mux} //nolint:gosec,exhaustruct_v5 // loopback-only dev server

	go func() { _ = server.Serve(listener) }() //nolint:gosec // intended long-running loop

	return listener.Addr().String(), server, nil
}

// captureRoute takes light+dark × desktop+mobile full-page screenshots of one
// route, each in a fresh browser session.
func captureRoute(base, route, out string) error {
	for _, theme := range []string{"light", "dark"} {
		for _, viewport := range []struct{ name string; width, height int }{
			{"desktop", desktopWidth, viewportHeight},
			{"mobile", mobileWidth, mobileHeight},
		} {
			name := routeName(route) + "-" + theme + "-" + viewport.name + ".png"

			if err := screenshot(base+route, theme, viewport.width, viewport.height, filepath.Join(out, name)); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}

			fmt.Println("wrote", name)
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
			chromedp.ExecPath(chromePath()),
			chromedp.Flag("headless", true),
		)...,
	)
	defer allocCancel()

	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()

	tabCtx, tabCancel = context.WithTimeout(tabCtx, routeTimeout)
	defer tabCancel()

	actions := []chromedp.Action{
		chromedp.EmulateViewport(int64(width), int64(height)),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(settle),
	}

	if theme == "dark" {
		actions = append(actions, chromedp.Evaluate(
			"document.documentElement.classList.add('dark')", nil,
		), chromedp.Sleep(settle))
	}

	var png []byte

	actions = append(actions, chromedp.FullScreenshot(&png, 90))

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
			chromedp.ExecPath(chromePath()),
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
		chromedp.FullScreenshot(&png, 90),
	)
	if err != nil {
		return fmt.Errorf("search smoke: %w", err)
	}

	if hits == 0 {
		return fmt.Errorf("search smoke FAILED: 0 hits for query")
	}

	//nolint:gosec // CLI-controlled screenshot output
	if err := os.WriteFile(out+"/search-smoke.png", png, 0o644); err != nil {
		return fmt.Errorf("write search smoke: %w", err)
	}

	fmt.Printf("search smoke PASS: %d hits (search-smoke.png)\n", hits)

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
