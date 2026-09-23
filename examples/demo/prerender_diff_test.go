package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// prerenderAutoID matches the EnsureID-generated tokens (tc-<prefix>-<16
// hex>) so the diff tolerates fresh random IDs between the two renders —
// structure must match, randomness must not.
var prerenderAutoID = regexp.MustCompile(`tc-[a-z-]+-[0-9a-f]{16}`)

// prerenderCSRFToken matches the kanban move form's hidden session-CSRF
// input; the live server mints a per-visitor token while the prerender
// mints a build-scoped one, so the value is normalized away.
var prerenderCSRFToken = regexp.MustCompile(`(name="csrf_token" value=")[^"]*(")`)

// prerenderBuildTime matches the Debug Information collapsible's
// server-rendered build timestamp (second precision, UTC); the two renders
// happen at different wall-clock seconds, so the value is normalized away.
var prerenderBuildTime = regexp.MustCompile(`Build: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)

// prerenderPolledUpdated matches the PolledRegion demo's server-time stamp
// (second precision, local); same render-timing tolerance as the build
// timestamp.
var prerenderPolledUpdated = regexp.MustCompile(`Updated \d{2}:\d{2}:\d{2} UTC`)

// prerenderTimeDatetime / prerenderTimeTitle match the wall-clock attribute
// values inside <time> elements (RelativeTime's datetime is RFC3339 with
// second precision, its title minute precision); both are derived from the
// render moment, so they are normalized.
var (
	prerenderTimeDatetime = regexp.MustCompile(`(<time[^>]*\sdatetime=")[^"]*(")`)
	prerenderTimeTitle    = regexp.MustCompile(`(<time[^>]*\stitle=")[^"]*(")`)
)

// normalizePrerender strips the BY-DESIGN differences between a prerendered
// page and the live server's response: the live stylesheet link (prerender
// pages embed fonts only, CSSPath=""), fresh random EnsureID tokens, the
// per-render CSRF tokens, and wall-clock stamps (text and <time>-element
// attributes) that the two renders rarely capture in the same second.
// Everything else must match byte-for-byte — if it does not, the static
// snapshot lies about what the live demo serves.
func normalizePrerender(html string) string {
	html = strings.ReplaceAll(html, `<link rel="stylesheet" href="/css/app.css">`, "")
	html = prerenderCSRFToken.ReplaceAllString(html, `${1}CSRF-NORMALIZED${2}`)
	html = prerenderBuildTime.ReplaceAllString(html, "Build: TIME-NORMALIZED")
	html = prerenderPolledUpdated.ReplaceAllString(html, "Updated TIME-NORMALIZED")
	html = prerenderTimeDatetime.ReplaceAllString(html, `${1}TIME-NORMALIZED${2}`)
	html = prerenderTimeTitle.ReplaceAllString(html, `${1}TIME-NORMALIZED${2}`)

	return prerenderAutoID.ReplaceAllString(html, "tc-NORMALIZED")
}

func TestPrerenderMatchesLiveServer(t *testing.T) {
	dir := t.TempDir()

	if err := prerender(dir); err != nil {
		t.Fatalf("prerender: %v", err)
	}

	srv := httptest.NewServer(newDemoHandler())
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
