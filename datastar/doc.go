// Package datastar provides templ components and helpers for integrating
// [Datastar] — a ~12 KiB zero-dependency frontend framework that unifies
// HTMX-style backend reactivity with Alpine.js-style frontend reactivity via
// Server-Sent Events (SSE) and reactive signals.
//
// This package mirrors the [htmx] package: it emits data-* attributes and
// injects the Datastar runtime <script> tag without importing any server-side
// SDK. The pinned version ([DatastarVersion1_0_2]) is derived from
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
// [Datastar]: https://data-star.dev/
// [htmx]: https://htmx.org
package datastar
