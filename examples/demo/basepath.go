package main

import (
	"net/http"
	"strings"
)

// demoBasePath is the canonical mount prefix for every demo page, asset, and
// endpoint. The Firebase Hosting site (templcomponents.lars.software) rewrites
// /demo/** to the Cloud Run service, and Cloud Run passes the ORIGINAL path
// through — so the demo must serve its routes under /demo to be reachable at
// https://templcomponents.lars.software/demo/**. withBasePath additionally
// keeps the bare paths (/health, /wire, /api/...) working so the raw run.app
// URL, the visualtest harness, and orchestrator probes keep working unchanged.
const demoBasePath = "/demo"

// demoDocsURL is the documentation site the demo links back to; the same
// site proxies this demo under /demo/**.
const demoDocsURL = "https://templcomponents.lars.software"

// demoGitHubURL is the repository link shown in the shell chrome.
const demoGitHubURL = "https://github.com/larsartmann/templ-components"

// demoURL prefixes a root-relative path with the demo base path. Every URL
// EMITTED by the demo templates (links, CSS, favicons, wire endpoints) goes
// through this so pages work identically on both origins:
//
//	https://templcomponents.lars.software/demo/display
//	https://templcomponents-demo-<hash>.us-central1.run.app/demo/display
func demoURL(path string) string {
	if path == "" {
		return demoBasePath
	}
	if strings.HasPrefix(path, demoBasePath+"/") || path == demoBasePath {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		return demoBasePath + "/" + path
	}

	return demoBasePath + path
}

// withBasePath dual-mounts the mux under /demo: requests to /demo/<path> are
// served by the same handler as /<path>. Bare /demo serves the home page
// directly (no redirect): Firebase Hosting's trailingSlash normalization may
// bounce /demo/ and /demo in either direction, so BOTH must render the home
// page — a redirect here could loop with the edge's own redirect.
func withBasePath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == demoBasePath {
			r.URL.Path = "/"

			next.ServeHTTP(w, r)

			return
		}
		if rest, ok := strings.CutPrefix(r.URL.Path, demoBasePath+"/"); ok {
			r.URL.Path = "/" + rest
		}

		next.ServeHTTP(w, r)
	})
}

// demoCSP is the demo's Content-Security-Policy. It must stay in sync with
// the /demo/** header override in website/firebase.json (if Firebase Hosting
// applies its own CSP to proxied responses, the browser INTERSECTS both —
// equal values make that a no-op). Allowances beyond 'self':
//   - cdn.jsdelivr.net: the ECharts and Datastar SDK scripts (SDKScript
//     components load their runtimes from the jsDelivr CDN)
//   - 'unsafe-eval' in script-src: the pinned Datastar runtime compiles
//     action expressions by evaluating strings (GenerateExpression → eval);
//     without it every data-on:* action dies with EvalError — the documented
//     runtime requirement (docs/datastar-runtime-facts.md). The nonce still
//     gates script LOADING; unsafe-eval only widens runtime-internal
//     expression compilation.
//   - fonts.googleapis.com / fonts.gstatic.com: demoFonts() web fonts
//   - ui-avatars.com: the Avatar/Image demo placeholder images
//   - 'unsafe-inline' in style-src: library components emit style=""
//     attributes (heatmap cells, progress bars); inline scripts stay
//     nonce-gated ('nonce-demo-nonce')
const demoCSP = "default-src 'self'; " +
	"script-src 'self' https://cdn.jsdelivr.net 'unsafe-eval' 'nonce-demo-nonce'; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
	"font-src 'self' https://fonts.gstatic.com; " +
	"img-src 'self' data: https://ui-avatars.com; " +
	"connect-src 'self'; " +
	"object-src 'none'; " +
	"base-uri 'self'; " +
	"form-action 'self'; " +
	"frame-ancestors 'none'"

// withSecurityHeaders stamps the demo's baseline security headers on every
// response, making the demo's CSP guarantee real on BOTH serving origins
// (the raw run.app URL has no Firebase in front of it).
func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", demoCSP)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
