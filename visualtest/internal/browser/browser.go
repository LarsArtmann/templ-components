// Package browser resolves the Chromium binary for the capture tools. The
// lookup lives in exactly one place so every tool honors
// CHROMEDP_CHROME_PATH identically.
package browser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ExecPath resolves the browser binary: CHROMEDP_CHROME_PATH (set by
// `nix run .#visual`-style wrappers) or "chromium" from PATH.
func ExecPath() string {
	if p := os.Getenv("CHROMEDP_CHROME_PATH"); p != "" {
		return p
	}

	return "chromium"
}

// selftestTimeout bounds every selftest request.
const selftestTimeout = 10 * time.Second

// ErrNoChromium is the selftest failure for a missing browser binary.
var ErrNoChromium = errors.New("no Chromium (set CHROMEDP_CHROME_PATH)")

// RequireChromium reports ErrNoChromium when no browser binary resolves.
func RequireChromium() error {
	if ExecPath() == "" {
		return ErrNoChromium
	}

	return nil
}

// HTTPGetSelftest performs a selftest's loopback GET with context and a
// bounded client — one home for the shape all capture tools share.
func HTTPGetSelftest(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), selftestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: selftestTimeout}

	return client.Do(req)
}

// OK prints the tool's selftest success line.
func OK(w io.Writer, tool, details string) {
	fmt.Fprintf(w, "%s selftest OK (%s)\n", tool, details)
}
