package main

import (
	"io"
	"net/http/httptest"
	"regexp"
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
// The WITH-nonce twin is TestDemoScriptsCarryDemoNonce below; the
// component-level sweeps are the integration nonce tests
// (integration/csp_nonce_test.go: TestAllInlineScriptsHaveNonce pins
// with-nonce rendering, TestNoEmptyNonceAttribute pins the omit-empty rule).
func TestNoEmptyNonceAcrossDemoPages(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(newDemoHandler())
	t.Cleanup(srv.Close)

	for _, route := range demoSweepRoutes() {
		t.Run(route, func(t *testing.T) {
			t.Parallel()

			resp, err := srv.Client().Get(srv.URL + route)
			if err != nil {
				t.Fatalf("GET %s: %v", route, err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read %s: %v", route, err)
			}

			if strings.Contains(string(body), `nonce=""`) {
				t.Errorf("GET %s renders nonce=\"\" — dead script under the demo CSP", route)
			}
		})
	}
}

// demoSweepRoutes is the full route list both nonce sweeps walk: every demo
// page plus the standalone recipe and error routes.
func demoSweepRoutes() []string {
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

	return routes
}

// TestDemoScriptsCarryDemoNonce is the WITH-nonce twin of the omit-empty
// sweep (#360): every `<script>` on every demo page that carries a nonce
// attribute must carry exactly demoNonceConst. A script with the WRONG
// nonce is exactly as dead as one with an empty nonce under the demo CSP —
// the constant exists so there is exactly one right value. Nonce-less
// scripts (omit-empty rule) are allowed here for non-CSP consumers; the
// demo pages themselves should never produce one.
func TestDemoScriptsCarryDemoNonce(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(newDemoHandler())
	t.Cleanup(srv.Close)

	for _, route := range demoSweepRoutes() {
		t.Run(route, func(t *testing.T) {
			t.Parallel()

			resp, err := srv.Client().Get(srv.URL + route)
			if err != nil {
				t.Fatalf("GET %s: %v", route, err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read %s: %v", route, err)
			}

			for _, m := range scriptNonceRe.FindAllStringSubmatch(string(body), -1) {
				if m[1] != demoNonceConst {
					t.Errorf("GET %s renders a script with nonce=%q, want demoNonceConst", route, m[1])
				}
			}
		})
	}
}

// scriptNonceRe captures the nonce value from SCRIPT/STYLE OPENING TAGS
// only. The [^>]*? arm cannot cross the tag's `>`, so nonce= TEXT inside a
// script BODY (the inlined htmx runtime contains its own
// `nonce="${config.inlineStyleNonce}"` source literal) never matches — the
// first naive `nonce="…"` regex false-positived exactly there. A wrong-value
// style nonce is the same dead-resource class, so styles are included.
var scriptNonceRe = regexp.MustCompile(`<(?:script|style)[^>]*?nonce="([^"]*)"`)

func TestScriptNonceReScopesToTagOpens(t *testing.T) {
	t.Parallel()

	// The htmx-runtime body text that once false-positived.
	const body = `<script nonce="demo">const e = ` + "` nonce=\"${Q.config.inlineStyleNonce}\"`" + `;</script>`
	for _, m := range scriptNonceRe.FindAllStringSubmatch(body, -1) {
		if m[1] != "demo" {
			t.Fatalf("regex matched inside script body: %q", m[1])
		}
	}
}

