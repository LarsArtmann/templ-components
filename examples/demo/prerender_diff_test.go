package main

import (
	"time"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// prerenderAutoID matches the EnsureID-generated tokens (tc-<prefix>-<16
// hex>) so the diff tolerates fresh random IDs between the two renders —
// structure must match, randomness must not.
var prerenderAutoID = regexp.MustCompile(`tc-[a-z-]+-[0-9a-f]{16}`)

// prerenderCSRFToken matches the kanban move form's hidden session-CSRF
// input; the live server mints a per-visitor token while the prerender
// mints a build-scoped one, so the value is normalized away.
var prerenderCSRFToken = regexp.MustCompile(`(name="csrf_token" value=")[^"]*(")`)


// normalizePrerender strips the two BY-DESIGN differences between a
// prerendered page and the live server's response: the live stylesheet link
// (prerender pages embed fonts only, CSSPath="") and fresh random EnsureID
// tokens. Everything else must match byte-for-byte — if it does not, the
// static snapshot lies about what the live demo serves.
func normalizePrerender(html string) string {
	html = strings.ReplaceAll(html, `<link rel="stylesheet" href="/css/app.css">`, "")
	html = prerenderCSRFToken.ReplaceAllString(html, `${1}CSRF-NORMALIZED${2}`)

	return prerenderAutoID.ReplaceAllString(html, "tc-NORMALIZED")
}

func TestPrerenderMatchesLiveServer(t *testing.T) {
	dir := t.TempDir()

	if err := prerender(dir); err != nil {
		t.Fatalf("prerender: %v", err)
	}

	srv := httptest.NewServer(newMux())
	t.Cleanup(srv.Close)

	// Cookie-jar client: the live server issues the session-CSRF cookie on
	// first response and renders the move form's hidden token from it — a
	// cookieless client gets a tokenless form and the comparison would
	// always drift. A jar is what a real browser does (backlog #229).
	jar, jarErr := cookiejar.New(nil)
	if jarErr != nil {
		t.Fatalf("cookie jar: %v", jarErr)
	}
	client := &http.Client{Jar: jar}

	routes := []struct {
		file string
		path string
	}{
		{"index.html", "/"},
		{"forms/index.html", "/forms"},
		{"recipes/dashboard.html", "/recipes/dashboard"},
		{"recipes/settings.html", "/recipes/settings"},
		{"recipes/login.html", "/recipes/login"},
		{"recipes/auth.html", "/recipes/auth"},
		{"users/index.html", "/users"},
	}

	for _, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			pre, err := os.ReadFile(filepath.Join(dir, route.file))
			if err != nil {
				t.Fatalf("read prerendered %s: %v", route.file, err)
			}

			// Retry-tolerant (backlog #250): the live fetch can transiently
			// answer non-200 or truncate under machine load; the invariant
			// under test is PRERENDER DRIFT, not transport reliability.
			var (
				live []byte
				code int
			)
			for attempt := range 3 {
				req, reqErr := http.NewRequest(http.MethodGet, srv.URL+route.path, nil) //nolint:noctx // test-local server, bounded response
				if reqErr != nil {
					t.Fatalf("request %s: %v", route.path, reqErr)
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatalf("fetch %s: %v", route.path, err)
				}

				live, err = io.ReadAll(resp.Body)
				code = resp.StatusCode
				resp.Body.Close()

				if err == nil && code == http.StatusOK {
					break
				}
				if attempt == 2 {
					t.Fatalf("live %s answered %d after retries", route.path, code)
				}
				time.Sleep(100 * time.Millisecond)
			}

			want := normalizePrerender(string(pre))
			got := normalizePrerender(string(live))

			if got != want {
				t.Errorf(
					"prerender drift on %s: static snapshot differs from the live page beyond the by-design CSS link — re-cut with -prerender",
					route.path,
				)
			}
		})
	}
}
