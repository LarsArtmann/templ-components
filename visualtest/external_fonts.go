package visualtest

import (
	"context"
	"fmt"

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
var externalWebFontPatterns = []*network.BlockPattern{
	{URLPattern: "https://fonts.googleapis.com/*", Block: true},
	{URLPattern: "https://fonts.gstatic.com/*", Block: true},
}

// blockExternalWebFonts fails every request to the Google Fonts CDN for the
// life of the tab, so pages render with the pinned system faces only.
func blockExternalWebFonts() chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		if _, err := cdp.Call(ctx, tab, network.Enable, network.EnableParams{}); err != nil {
			return fmt.Errorf("enable network domain: %w", err)
		}

		if _, err := cdp.Call(
			ctx,
			tab,
			network.SetBlockedURLs,
			network.SetBlockedURLsParams{URLPatterns: externalWebFontPatterns},
		); err != nil {
			return fmt.Errorf("block external web fonts: %w", err)
		}

		return nil
	})
}
