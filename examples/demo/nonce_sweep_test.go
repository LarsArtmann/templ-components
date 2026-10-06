package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNoEmptyNonceAcrossDemoPages sweeps every rendered demo page and
// asserts the body never contains `nonce=""`. An empty nonce is the
// dead-script class: the browser rejects the script (the attribute value
// can never match the CSP header) while the markup looks fine — the 2026-10-01
// outage shipped 50 dead component scripts this way when the MPA rewrite
// dropped the explicit nonce pass. The demo CSP makes this an outage, not
// a cosmetic issue, so the sweep is a real test, not AGENTS prose (f22).
//
// Scripts WITHOUT any nonce attribute are fine here: the nonce-omit-empty
// rule renders them nonce-less and the demo's `withSecurityHeaders` CSP
// blocks them — but the demo-wide rule is that demo content passes
// demoNonceConst explicitly, so a nonce-less script on a demo page would
// be the next bug to fix; this sweep pins only the never-ship-empty class.
func TestNoEmptyNonceAcrossDemoPages(t *testing.T) {
	t.Parallel()

	routes := []string{"/"}
	for _, page := range demoPages() {
		routes = append(routes, page.Path)
	}
	routes = append(
		routes,
		"/recipes/dashboard",
		"/recipes/settings",
		"/recipes/login",
		"/recipes/auth",
		"/errors/400",
		"/errors/403",
		"/errors/404",
		"/errors/409",
		"/errors/500",
		"/errors/503",
		"/errors/full",
		"/errors/404-page",
		"/errors/playground",
	)

	srv := httptest.NewServer(newDemoHandler())
	t.Cleanup(srv.Close)

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			t.Parallel()

			resp, err := srv.Client().Get(srv.URL + route)
			if err != nil {
				t.Fatalf("GET %s: %v", route, err)
			}
			defer resp.Body.Close()

			buf := new(strings.Builder)
			if _, err := readAll(resp.Body, buf); err != nil {
				t.Fatalf("read %s: %v", route, err)
			}

			if body := buf.String(); strings.Contains(body, `nonce=""`) {
				t.Errorf("GET %s renders nonce=\"\" — dead script under the demo CSP", route)
			}
		})
	}
}
