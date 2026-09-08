package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// prerenderRoutes maps each prerendered file to the live route serving the
// same page. The prerender pipeline and the live handlers must render the
// SAME markup (modulo the documented normalizations below) — a drift means
// one of them went stale after a demo change (#167).
var prerenderRoutes = map[string]string{
	"index.html":             "/",
	"forms/index.html":       "/forms",
	"recipes/dashboard.html": "/recipes/dashboard",
	"recipes/settings.html":  "/recipes/settings",
	"recipes/login.html":     "/recipes/login",
	"recipes/auth.html":      "/recipes/auth",
	"users/index.html":       "/users",
}

var (
	// liveCSSLinkRe strips the live-only stylesheet link (prerender sets
	// CSSPath: "" — it inlines fonts via HeadContent instead).
	liveCSSLinkRe = regexp.MustCompile(`<link rel="stylesheet" href="/css/app.css">`)
	// datetimeAttrRe normalizes time-dependent datetime attributes
	// (PolledRegion stamps time.Now() at render).
	datetimeAttrRe = regexp.MustCompile(`datetime="[^"]*"`)
	// updatedAtTextRe normalizes the PolledRegion "Updated HH:MM:SS" footer.
	updatedAtTextRe = regexp.MustCompile(`Updated [0-9:]{8}`)
	// ensureIDRe normalizes EnsureID's random suffixes (crypto/rand per
	// render — same normalization class as utils/golden). Prefixes may be
	// multi-word (tc-mobile-menu-<hex>).
	ensureIDRe = regexp.MustCompile(`tc-[a-z0-9]+(?:-[a-z0-9]+)*-[0-9a-f]{16}`)
)

// normalizePrerenderDiff removes the expected prerender/live differences so
// any REMAINING delta is real drift.
func normalizePrerenderDiff(page string) string {
	page = liveCSSLinkRe.ReplaceAllString(page, "")
	page = datetimeAttrRe.ReplaceAllString(page, `datetime="T"`)
	page = updatedAtTextRe.ReplaceAllString(page, "Updated T")
	page = ensureIDRe.ReplaceAllString(page, "tc-X-NORMALIZED")

	return page
}

// firstDiffLines returns up to 3 line pairs around the first difference.
func firstDiffLines(a, b string) string {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")

	for i := 0; i < len(aLines) && i < len(bLines); i++ {
		if aLines[i] != bLines[i] {
			start := i - 1
			if start < 0 {
				start = 0
			}

			endA := min(i+2, len(aLines))
			endB := min(i+2, len(bLines))

			return fmt.Sprintf(
				"first diff at line %d:\n--- prerender\n%s\n--- live\n%s",
				i+1,
				strings.Join(aLines[start:endA], "\n"),
				strings.Join(bLines[start:endB], "\n"),
			)
		}
	}

	if len(aLines) != len(bLines) {
		return fmt.Sprintf("line counts differ: prerender %d, live %d", len(aLines), len(bLines))
	}

	return "contents differ (no line diff located)"
}

func TestPrerenderMatchesLiveServer(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := prerender(dir); err != nil {
		t.Fatalf("prerender: %v", err)
	}

	server := httptest.NewServer(newMux())
	t.Cleanup(server.Close)

	client := server.Client()

	for file, route := range prerenderRoutes {
		t.Run(route, func(t *testing.T) {
			t.Parallel()

			preBytes, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(file)))
			if err != nil {
				t.Fatalf("read prerendered page: %v", err)
			}

			resp, err := client.Get(server.URL + route)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("live route %s: status %d", route, resp.StatusCode)
			}

			liveBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			pre := normalizePrerenderDiff(string(preBytes))
			live := normalizePrerenderDiff(string(liveBytes))

			if os.Getenv("TC_PRERENDER_DEBUG") != "" {
				_ = os.WriteFile("/tmp/tc-pre.html", []byte(pre), 0o644)
				_ = os.WriteFile("/tmp/tc-live.html", []byte(live), 0o644)
			}

			if pre != live {
				t.Errorf("prerender %s drifted from live %s:\n%s", file, route, firstDiffLines(pre, live))
			}
		})
	}
}
