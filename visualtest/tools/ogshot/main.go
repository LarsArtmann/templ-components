// Command ogshot renders the /sales social-preview card (1200×630 PNG) into
// website/public/og/. The card is a purpose-built HTML composition using the
// site's compiled CSS, captured headless at exactly 1200×630 — regenerate it
// whenever the brand styling changes:
//
//	GOWORK=off CHROMEDP_CHROME_PATH=/path/to/chromium go run ./tools/ogshot \
//	  -dist ../website/dist -out ../website/public/og/sales.png
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

const (
	ogWidth  = 1200
	ogHeight = 630

	settle        = 700 * time.Millisecond
	routeTimeout  = 60 * time.Second
	shotQuality   = 95
	ogPageName    = "tc-og-card.html"
	screenshotFmt = 0o644
)

// ogCardHTML is the social-preview composition. It deliberately avoids
// derived counts (a static PNG cannot track them) — tagline and tech row
// only. Served from the dist root so /assets/app.css resolves.
const ogCardHTML = `<!doctype html><html lang="en"><head><meta charset="utf-8">
<link rel="stylesheet" href="/assets/app.css">
<style>html,body{margin:0;padding:0}</style></head>
<body class="font-sans antialiased bg-gray-950 text-white">
<div class="relative flex flex-col justify-between overflow-hidden p-16" style="width:1200px;height:630px">
  <div class="pointer-events-none absolute rounded-full bg-blue-600/25 blur-[120px] h-[520px] w-[520px] -top-32 -end-24"></div>
  <div class="pointer-events-none absolute rounded-full bg-indigo-600/15 blur-[100px] h-[360px] w-[360px] -bottom-24 -start-16"></div>
  <div class="relative flex items-center gap-3">
    <svg width="44" height="44" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="templ-components logo"><rect width="28" height="28" rx="6" fill="#3b82f6"></rect><path d="M10 9L5 14l5 5M18 9l5 5-5 5M15.5 7l-3 14" stroke="white" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"></path></svg>
    <span class="text-2xl font-semibold tracking-tight text-white">templ-components</span>
  </div>
  <div class="relative">
    <h1 class="font-bold leading-[1.08] text-[64px] tracking-tighter text-white m-0">Build complete UIs from Go —<br>without <span class="text-blue-400">the frontend stack</span>.</h1>
    <p class="text-[26px] leading-snug text-gray-400 mt-6 mb-0">Server-rendered, CSP-safe, dark-mode tested —<br>the website selling it is rendered by it.</p>
  </div>
  <div class="relative flex items-center gap-4 text-lg text-gray-300">
    <span class="rounded-full border border-gray-700 bg-gray-900 px-5 py-2">templ</span>
    <span class="rounded-full border border-gray-700 bg-gray-900 px-5 py-2">HTMX</span>
    <span class="rounded-full border border-gray-700 bg-gray-900 px-5 py-2">Tailwind v4</span>
    <span class="rounded-full border border-gray-700 bg-gray-900 px-5 py-2">MIT licensed</span>
  </div>
</div>
</body></html>`

func main() {
	dist := flag.String("dist", "../website/dist", "website dist directory to serve")
	out := flag.String("out", "../website/public/og/sales.png", "output PNG path")

	flag.Parse()

	if err := run(*dist, *out); err != nil {
		log.Fatal(err)
	}
}

func run(dist, out string) error {
	addr, cleanup, err := serveDist(dist)
	if err != nil {
		return err
	}

	defer cleanup()

	base := "http://" + addr

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

	var png []byte

	if err := chromedp.Run(tabCtx,
		chromedp.EmulateViewport(ogWidth, ogHeight),
		chromedp.Navigate(base+"/"+ogPageName),
		chromedp.WaitReady("body"),
		chromedp.Sleep(settle),
		chromedp.FullScreenshot(&png, shotQuality),
	); err != nil {
		return fmt.Errorf("capture OG card: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := os.WriteFile(out, png, screenshotFmt); err != nil { //nolint:gosec // CLI-controlled output path
		return fmt.Errorf("write %s: %w", out, err)
	}

	fmt.Fprintln(os.Stdout, "ogshot: wrote", out)

	return nil
}

// serveDist serves the dist directory with the OG card page injected, and
// returns a cleanup that shuts the server down and removes the injected file.
func serveDist(dist string) (string, func(), error) {
	cardPath := filepath.Join(dist, ogPageName)

	if err := os.WriteFile(cardPath, []byte(ogCardHTML), screenshotFmt); err != nil { //nolint:gosec // CLI-controlled dist root
		return "", nil, fmt.Errorf("write OG card page: %w", err)
	}

	var netListenConfig net.ListenConfig

	listener, err := netListenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		_ = os.Remove(cardPath)

		return "", nil, fmt.Errorf("listen: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clean, _, _ := strings.Cut(r.URL.Path, "?")
		if !strings.HasSuffix(clean, ".html") && !strings.Contains(clean, ".") {
			cleanPath := filepath.Join(dist, clean+".html")
			if _, err := os.Stat(cleanPath); err == nil { //nolint:gosec // CLI-controlled dist root
				http.ServeFile(w, r, cleanPath) //nolint:gosec // CLI-controlled dist root

				return
			}
		}

		http.FileServer(http.Dir(dist)).ServeHTTP(w, r)
	})

	server := &http.Server{Handler: mux} //nolint:gosec,exhaustruct_v5 // loopback-only dev server

	go func() { _ = server.Serve(listener) }()

	cleanup := func() {
		_ = server.Close()
		_ = os.Remove(cardPath)
	}

	return listener.Addr().String(), cleanup, nil
}

// chromePath resolves the browser binary: CHROMEDP_CHROME_PATH (set by
// `nix run .#visual`-style wrappers) or "chromium" from PATH.
func chromePath() string {
	if p := os.Getenv("CHROMEDP_CHROME_PATH"); p != "" {
		return p
	}

	return "chromium"
}
