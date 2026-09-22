// Package browser resolves the Chromium binary for the capture tools. The
// lookup lives in exactly one place so every tool honors
// CHROMEDP_CHROME_PATH identically.
package browser

import "os"

// ExecPath resolves the browser binary: CHROMEDP_CHROME_PATH (set by
// `nix run .#visual`-style wrappers) or "chromium" from PATH.
func ExecPath() string {
	if p := os.Getenv("CHROMEDP_CHROME_PATH"); p != "" {
		return p
	}

	return "chromium"
}
