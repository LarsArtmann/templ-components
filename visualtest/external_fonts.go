package visualtest

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// The demo and the website both link Google-hosted webfonts (Inter, JetBrains
// Mono, Space Grotesk) with `display=swap`. That makes every screenshot
// dependent on a third party: the first frame renders the fallback face, the
// swap lands whenever the CDN answers, and an offline or slow runner captures a
// different page than a warm one. Route goldens flaked at ~1% pixel mismatch on
// a changing subset of pages because of exactly this — the same page compared
// itself against a golden captured with a different set of faces loaded.
//
// The flake pin (flake.nix `visual`) already installs a PURE fontconfig with
// Inter/DejaVu only, so the deterministic rendering this suite wants is the
// system-face one. Blocking the CDN requests is what makes the pin actually
// decide the pixels: no swap can happen, and no run depends on network reach
// or latency. Both origins are blocked before navigation in every
// pixel-comparison path (component goldens, demo routes, site routes).

// The two hosts a Google Fonts <link> resolves to: the stylesheet and the font
// files themselves.
const (
	googleFontsCSSOrigin  = "https://fonts.googleapis.com/*"
	googleFontsFileOrigin = "https://fonts.gstatic.com/*"

	// fontsLoadTimeout bounds the wait for every declared face to finish
	// loading; fontsLoadPollInterval keeps that wait responsive.
	fontsLoadTimeout      = 10 * time.Second
	fontsLoadPollInterval = 50 * time.Millisecond
)

// externalWebFontPatterns is the CDP block list for both Google Fonts origins.
func externalWebFontPatterns() []*network.BlockPattern {
	return []*network.BlockPattern{
		{URLPattern: googleFontsCSSOrigin, Block: true},
		{URLPattern: googleFontsFileOrigin, Block: true},
	}
}

// blockExternalWebFonts fails every request to the Google Fonts CDN for the
// life of the tab, so pages render with the pinned system faces only.
func blockExternalWebFonts() chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		//nolint:exhaustruct_v5 // omitted protocol fields take their defaults
		enable := network.EnableParams{}
		if _, err := cdp.Call(ctx, tab, network.Enable, enable); err != nil {
			return fmt.Errorf("enable network domain: %w", err)
		}

		if _, err := cdp.Call(
			ctx,
			tab,
			network.SetBlockedURLs,
			network.SetBlockedURLsParams{URLPatterns: externalWebFontPatterns()},
		); err != nil {
			return fmt.Errorf("block external web fonts: %w", err)
		}

		return nil
	})
}

// waitFontsSettled blocks until the document has finished loading every
// declared face, so a capture can never race `font-display: swap`. The website
// self-hosts its woff2 files (same-origin, but still asynchronous); the demo
// relies on system faces, which report loaded immediately. Called AFTER the
// settle chain so a load that starts late is still observed.
func waitFontsSettled() chromedp.Action[chromedp.Void] {
	return pollTrue(
		`document.readyState === 'complete' && document.fonts.status === 'loaded'`,
		chromedp.WithPollingTimeout(fontsLoadTimeout),
		chromedp.WithPollingInterval(fontsLoadPollInterval),
	)
}
