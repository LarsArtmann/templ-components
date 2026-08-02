// Package datastar provides templ components and helpers for integrating
// [Datastar] — a ~12 KiB zero-dependency frontend framework that unifies
// HTMX-style backend reactivity with Alpine.js-style frontend reactivity via
// Server-Sent Events (SSE) and reactive signals.
//
// This package mirrors the [htmx] package: it emits data-* attributes and
// injects the Datastar runtime <script> tag without importing the
// datastar-go SDK. Consumers who want SSE streaming add
// github.com/starfederation/datastar-go to their own go.mod.
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
//	    @display.StatCard(display.StatCardProps{Label: "Active Users", Value: "—"})
//	}
//
// The server endpoint streams patches using the datastar-go SDK:
//
//	func streamHandler(w http.ResponseWriter, r *http.Request) {
//	    sse := datastar.NewSSE(w, r)
//	    for {
//	        sse.PatchElementTempl(metricsCard(currentMetrics()))
//	        time.Sleep(time.Second)
//	    }
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
