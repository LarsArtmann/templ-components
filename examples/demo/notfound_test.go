package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUnknownRouteServesStyledNotFound pins the demo's catch-all: unknown
// paths must answer 404 with the library's own NotFound404 page (demo shell,
// base-path-prefixed escape links, nonce'd scripts) — never Go's plain-text
// "404 page not found" — and must do so through the full middleware stack on
// both mount forms (bare path and the /demo dual-mount the Firebase rewrite
// forwards).
func TestUnknownRouteServesStyledNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
	}{
		{name: "bare mux path", path: "/definitely-not-a-demo-page"},
		{name: "base-path dual-mount", path: "/demo/definitely-not-a-demo-page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			newDemoHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != http.StatusNotFound {
				t.Errorf("GET %s status = %d, want %d", tt.path, rec.Code, http.StatusNotFound)
			}

			body := rec.Body.String()
			for _, want := range []string{
				"Page not found",
				`href="/demo/"`,
				`href="/demo/display"`,
				`href="/demo/error-pages"`,
				`nonce="demo-nonce"`,
			} {
				if !strings.Contains(body, want) {
					t.Errorf("GET %s body does not contain %q", tt.path, want)
				}
			}

			if strings.Contains(body, "404 page not found") {
				t.Errorf("GET %s body contains Go's plain-text 404 fallback", tt.path)
			}
		})
	}
}

// TestUnknownRouteStillServesHome pins the other half of the catch-all: the
// root path keeps rendering the demo home page, so the 404 wiring cannot
// regress into swallowing the landing page.
func TestUnknownRouteStillServesHome(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	newDemoHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("GET / status = %d, want %d", rec.Code, http.StatusOK)
	}

	if body := rec.Body.String(); !strings.Contains(body, "templ-components") {
		t.Errorf("GET / body does not contain the demo home page")
	}
}
