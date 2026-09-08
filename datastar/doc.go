// Package datastar provides templ components and helpers for integrating
// [Datastar] — a ~12 KiB zero-dependency frontend framework that unifies
// HTMX-style backend reactivity with Alpine.js-style frontend reactivity via
// Server-Sent Events (SSE) and reactive signals.
//
// This package mirrors the [htmx] package: it emits data-* attributes and
// injects the Datastar runtime <script> tag without importing any server-side
// SDK. The pinned version ([DatastarVersion1_0_3]) is derived from
// [github.com/larsartmann/go-datastar/static].Version so the CDN URL and the
// embedded bundle can never drift.
//
// # Server-side SDK
//
// Consumers who want SSE streaming should use
// [github.com/larsartmann/go-datastar] — a protocol library where every patch
// is a first-class value (not a method call on a live connection). Add it to
// your own go.mod:
//
//	go get github.com/larsartmann/go-datastar
//
// # Self-hosting the Datastar runtime
//
// By default, SDKScript loads the runtime from the jsDelivr CDN. For
// self-hosting, use [github.com/larsartmann/go-datastar/static] (zero
// dependencies — its go.mod requires nothing):
//
//	mux.Handle("GET /datastar.js", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
//	    w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
//	    _, _ = w.Write(static.Bytes())
//	}))
//	@datastar.SDKScript(datastar.SDKScriptProps{Src: "/datastar.js"})
//
// # Quick start
//
// Inject the runtime once per page (in your layout, alongside or instead of
// the HTMX script):
//
//	@datastar.SDKScript(datastar.DefaultSDKScriptProps())
//
// Establish an SSE-powered live region:
//
//	@datastar.LiveRegion(datastar.LiveRegionProps{
//	    URL:       "/stream/metrics",
//	    AutoStart: true,
//	}) {
//	    @display.StatCard(display.StatCardProps{
//	        BaseProps: utils.BaseProps{ID: "metrics"},
//	        Label:     "Active Users",
//	        Value:     "—",
//	    })
//	}
//
// Two options are load-bearing on real deployments:
//
//   - Retry: RetryAlways — the runtime default ('auto') never reconnects
//     after a clean stream EOF, which is exactly what a server restart
//     produces. Without it the region goes permanently stale.
//   - Cancellation: CancellationCleanup — when an HTMX swap replaces the
//     region (e.g. a filter change re-renders it with a different URL),
//     cleanup aborts the old stream so it cannot keep re-patching the
//     region with stale content. The swapped-region combination renders
//     @get(url, {retry: 'always', requestCancellation: 'cleanup'}).
//
// The server endpoint streams patches using go-datastar. Target a child
// element by selector — patch modes other than the default outer mode
// require one, and the default outer mode matches incoming root elements
// by their id (id-less fragments are dropped with a console warning):
//
//	func streamHandler(w http.ResponseWriter, r *http.Request) {
//	    stream := sse.NewStream(w, r)
//	    defer func() { _ = stream.Close() }()
//	    resp := datastar.NewResponse(stream)
//	    _ = resp.PatchElementsTempl(metricsCardContent(),
//	        datastar.WithSelector("#metrics"), datastar.WithModeInner())
//	}
//
// # When to choose Datastar over HTMX
//
// See docs/research/datastar-integration-analysis.md for the full analysis.
// Briefly: HTMX is the default (zero JS dependency). Datastar is opt-in for
// real-time streaming, reactive client-side state, or when the consumer's app
// already uses Datastar.
//
// # Re-audit contract (for contributors)
//
// The runtime bundle this package integrates against is an EXTERNAL
// dependency pinned in go.mod (go-datastar/static). Its wire behavior is
// documented in docs/datastar-runtime-facts.md (with a bundle-provenance
// block: pin version, byte size, sha256, extraction commands) and enforced
// by tests:
//
//   - datastar.TestPinnedRuntimeBundleContract pins the bundle sha256 plus
//     every integration token (event names, retry literals) — a pin bump
//     that changes any pinned byte fails there first.
//   - examples/demo/sse_test.go pins the SSE wire format the demo streams.
//
// When bumping the pin, follow docs/external-dependency-bumps.md end to end:
// verify-at-source, diff old-vs-new bundle, bump, re-audit every fact in
// docs/datastar-runtime-facts.md, update pinnedBundleSHA256 in
// bundle_guard_test.go, and warm [Unreleased] in CHANGELOG.md.
//
// [Datastar]: https://data-star.dev/
// [htmx]: https://htmx.org
package datastar
