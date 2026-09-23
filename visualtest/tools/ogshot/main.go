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
	"time"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/visualtest/internal/browser"
	"github.com/larsartmann/templ-components/visualtest/internal/distserver"
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

// ogCardHTML is the social-preview composition. It is styled with a
// self-contained <style> block — deliberately NOT Tailwind classes, because
// the compiled app.css only contains utilities Tailwind saw in .templ files,
// so any class unique to this card would silently not exist. Served from the
// dist root so /assets/app.css (which loads the Space Grotesk/JetBrains Mono
// webfonts) resolves. Derived counts are avoided on purpose: a static PNG
// cannot track them.
const ogCardHTML = `<!doctype html><html lang="en"><head><meta charset="utf-8">
<link rel="stylesheet" href="/assets/app.css">
<style>
  html,body{margin:0;padding:0}
  .card{position:relative;display:flex;flex-direction:column;justify-content:space-between;
    width:1200px;height:630px;overflow:hidden;box-sizing:border-box;padding:64px;
    background:#030712;color:#fff;font-family:'Space Grotesk',-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif}
  .glow-a{position:absolute;top:-128px;right:-96px;width:520px;height:520px;border-radius:9999px;
    background:rgba(37,99,235,.25);filter:blur(120px);pointer-events:none}
  .glow-b{position:absolute;bottom:-96px;left:-64px;width:360px;height:360px;border-radius:9999px;
    background:rgba(79,70,229,.15);filter:blur(100px);pointer-events:none}
  .brand{position:relative;display:flex;align-items:center;gap:14px}
  .brand-name{font-size:26px;font-weight:600;letter-spacing:-.02em}
  h1{position:relative;margin:0;font-size:64px;line-height:1.08;font-weight:700;letter-spacing:-.03em}
  h1 .accent{color:#60a5fa}
  .tagline{position:relative;margin:24px 0 0;font-size:26px;line-height:1.35;color:#9ca3af}
  .tech{position:relative;display:flex;gap:16px}
  .tech span{font-size:19px;color:#d1d5db;border:1px solid #374151;border-radius:9999px;
    background:#111827;padding:10px 22px}
</style></head>
<body>
<div class="card">
  <div class="glow-a"></div><div class="glow-b"></div>
  <div class="brand">
    <svg width="44" height="44" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="templ-components logo"><rect width="28" height="28" rx="6" fill="#2563eb"></rect><path d="M10 9L5 14l5 5M18 9l5 5-5 5M15.5 7l-3 14" stroke="white" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"></path></svg>
    <span class="brand-name">templ-components</span>
  </div>
  <div>
    <h1>Build complete UIs from Go &mdash;<br>without <span class="accent">the frontend stack</span>.</h1>
    <p class="tagline">Server-rendered, CSP-safe, dark-mode tested &mdash;<br>the website selling it is rendered by it.</p>
  </div>
  <div class="tech"><span>templ</span><span>HTMX</span><span>Tailwind v4</span><span>MIT licensed</span></div>
</div>
</body></html>`

func main() {
	dist := flag.String("dist", "../website/dist", "website dist directory to serve")
	out := flag.String("out", "../website/public/og/sales.png", "output PNG path")
	selftest := flag.Bool("selftest", false, "verify allocator + dist listener come-up, then exit")

	flag.Parse()

	if *selftest {
		if err := browser.RequireChromium(); err != nil {
			log.Fatalf("selftest FAIL: %v", err)
		}

		if err := distserver.Selftest(*dist, "sales.html"); err != nil {
			log.Fatalf("selftest FAIL: %v", err)
		}

		browser.OK(os.Stdout, "ogshot", "allocator + dist listener")

		return
	}

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
			chromedp.ExecPath(browser.ExecPath()),
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

	if err := os.MkdirAll(filepath.Dir(out), 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := os.WriteFile(out, png, screenshotFmt); err != nil {
		return fmt.Errorf("write %s: %w", out, err)
	}

	fmt.Fprintln(os.Stdout, "ogshot: wrote", out)

	return nil
}

// serveDist serves the dist directory with the OG card page injected, and
// returns a cleanup that shuts the server down and removes the injected file.
func serveDist(dist string) (string, func(), error) {
	cardPath := filepath.Join(dist, ogPageName)

	if err := os.WriteFile(cardPath, []byte(ogCardHTML), screenshotFmt); err != nil {
		return "", nil, fmt.Errorf("write OG card page: %w", err)
	}

	var netListenConfig net.ListenConfig

	listener, err := netListenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		_ = os.Remove(cardPath)

		return "", nil, fmt.Errorf("listen: %w", err)
	}

	server := &http.Server{Handler: distserver.Handler(dist)} //nolint:gosec // loopback-only dev server

	go func() { _ = server.Serve(listener) }()

	cleanup := func() {
		_ = server.Close()
		_ = os.Remove(cardPath)
	}

	return listener.Addr().String(), cleanup, nil
}
