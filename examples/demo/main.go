// Demo application showcasing all templ-components.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils/wire"
)

// writeDatastarPatch writes one Datastar patch-elements SSE event in the
// Datastar v1 wire format:
//
//	event: datastar-patch-elements
//	data: selector #target
//	data: mode inner
//	data: elements <p>...</p>
//
// Format details (all verified against the pinned v1.0.2 runtime bundle):
//
//   - Every line needs a trailing REAL newline, and the event ends with a
//     blank line — without it the browser never dispatches the event.
//   - The client parses each dataline by splitting on the FIRST space, so
//     multi-line HTML must repeat the "elements " key on every line; values
//     for the same key are joined with newlines.
//   - Interior blank lines are PRESERVED: an empty line is emitted as an
//     empty-value dataline ("data: elements " with trailing space), which the
//     runtime joins back to "\n" — <pre> content survives the round-trip.
//     A truly empty byte line would terminate the event early, so lines are
//     never dropped.
//   - selector/mode are single-line protocol values; CR/LF is replaced with
//     spaces so a crafted selector cannot inject datalines.
func writeDatastarPatch(w io.Writer, selector, mode, html string) error {
	selector = strings.ReplaceAll(selector, "\r", " ")
	selector = strings.ReplaceAll(selector, "\n", " ")
	mode = strings.ReplaceAll(mode, "\r", " ")
	mode = strings.ReplaceAll(mode, "\n", " ")

	var b strings.Builder
	b.WriteString("event: datastar-patch-elements\n")
	if selector != "" {
		b.WriteString("data: selector " + selector + "\n")
	}
	if mode != "" {
		b.WriteString("data: mode " + mode + "\n")
	}
	if payload := strings.TrimSpace(html); payload != "" {
		for line := range strings.SplitSeq(payload, "\n") {
			line = strings.TrimSuffix(line, "\r")
			b.WriteString("data: elements " + line + "\n")
		}
	}
	b.WriteString("\n")

	_, err := io.WriteString(w, b.String())

	return err
}

func main() {
	prerenderDir := flag.String(
		"prerender",
		"",
		"Pre-render demo pages to static HTML in the given directory instead of starting a server",
	)
	flag.Parse()

	if *prerenderDir != "" {
		if err := prerender(*prerenderDir); err != nil {
			fmt.Fprintf(os.Stderr, "Pre-render error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	server := newServer(newMux())

	fmt.Printf("Demo running at http://localhost:%s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}

// demoTransport selects which wire dialect the wire demo section renders.
type demoTransport string

const (
	demoTransportBoth     demoTransport = "both"
	demoTransportHTMX     demoTransport = "htmx"
	demoTransportDatastar demoTransport = "datastar"
)

// parseDemoTransport resolves the ?transport= query value; anything unknown
// falls back to the both-dialect default (repo convention: graceful fallback).
func parseDemoTransport(s string) demoTransport {
	switch demoTransport(s) {
	case demoTransportHTMX, demoTransportDatastar:
		return demoTransport(s)
	default:
		return demoTransportBoth
	}
}

// wireTransport maps the demo selection to the wire package's Transport.
func (t demoTransport) wireTransport() wire.Transport {
	if t == demoTransportDatastar {
		return wire.TransportDatastar
	}

	return wire.TransportHTMX
}

// wireTarget is the htmx target for the single-transport view.
func (t demoTransport) wireTarget() string {
	if t == demoTransportHTMX {
		return "#wire-htmx-out"
	}

	return ""
}

// outID is the region id the single-transport view renders.
func (t demoTransport) outID() string {
	if t == demoTransportHTMX {
		return "wire-htmx-out"
	}

	return "wire-datastar-out"
}

// humanName is the display name used in button labels.
func (t demoTransport) humanName() string {
	if t == demoTransportDatastar {
		return "Datastar"
	}

	return "htmx"
}

// newMux builds the demo routes. Exposed for endpoint tests — the SSE wire
// format broke silently for months because nothing exercised the handlers.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint for Cloud Run / container orchestration
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Serve pre-compiled CSS (embedded in binary, no CDN dependency)
	mux.HandleFunc("/css/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		_, _ = w.Write(embeddedCSS)
	})

	// SVG favicon (matches indigo theme)
	mux.HandleFunc("/favicon.svg", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		fmt.Fprint(
			w,
			`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" rx="20" fill="#4f46e5"/><text x="50" y="70" font-size="56" font-weight="bold" text-anchor="middle" fill="white" font-family="sans-serif">tc</text></svg>`,
		)
	})

	// Mock HTMX endpoints for interactive demo components. Every endpoint the
	// demo templates reference must exist — a 404 in a demo fires an error
	// toast via GlobalErrorHandling.
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		noStore(w)
		componentOr500(w, r, loadMoreResponse(r.URL.Query().Get("cursor")))
	})

	mux.HandleFunc("/api/items/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		noStore(w)
		fmt.Fprint(
			w,
			`<div class="rounded-lg border border-green-200 dark:border-green-800 bg-green-50 dark:bg-green-900/20 p-4">
			<p class="text-sm text-green-800 dark:text-green-200">Item deleted successfully (mock endpoint).</p>
		</div>`,
		)
	})

	mux.HandleFunc("/api/save", func(w http.ResponseWriter, _ *http.Request) {
		// Small delay so the LoadingButton "Saving..." state is visible.
		time.Sleep(600 * time.Millisecond)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		noStore(w)
		fmt.Fprint(w, `<span class="text-sm text-green-700 dark:text-green-300">Saved.</span>`)
	})

	// PolledRegion demo: each poll returns a fresh region (hx-swap=outerHTML
	// replaces the whole region, so the response must re-arm the poll).
	// After the third tick the region is returned without polling — the demo
	// settles instead of hammering the server forever.
	mux.HandleFunc("/api/demo-stats", func(w http.ResponseWriter, r *http.Request) {
		tick := 1
		if v, err := strconv.Atoi(r.URL.Query().Get("tick")); err == nil {
			tick = v + 1
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		noStore(w)
		componentOr500(w, r, polledStatsRegion(tick))
	})

	// FilterDropdown demo: returns a filtered user fragment.
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		sort := r.URL.Query().Get("sort")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		noStore(w)
		componentOr500(w, r, filteredUsersFragment(status, sort))
	})

	// Mock Datastar SSE endpoint — streams periodic updates in Datastar's
	// v1 wire format (datastar-patch-elements events with keyed datalines).
	// The LiveRegion component connects to this via
	// data-init="@get('/api/datastar/stream')" when the Datastar runtime is
	// loaded. See writeDatastarPatch for the format details.
	mux.HandleFunc("/api/datastar/stream", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)

			return
		}

		// The demo server sets WriteTimeout globally (see below), which would
		// cut the stream off mid-flight. SSE connections live as long as the
		// client stays connected, so clear the deadline for this connection.
		_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		// Disable proxy buffering (e.g. nginx) so events arrive as they are
		// written instead of being batched.
		w.Header().Set("X-Accel-Buffering", "no")

		ctx := r.Context()
		patches := time.NewTicker(2 * time.Second)
		defer patches.Stop()
		// SSE comment frames keep reverse proxies and firewalls from killing
		// the "idle" connection; browsers ignore them entirely. The 15s
		// interval mirrors stream.Heartbeat from go-sse (see
		// docs/recipes/datastar-integration.md).
		keepalive := time.NewTicker(15 * time.Second)
		defer keepalive.Stop()

		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-patches.C:
				i++
				err := writeDatastarPatch(
					w,
					"#datastar-live-content",
					"inner",
					fmt.Sprintf(
						`<p class="text-sm text-gray-600 dark:text-gray-400">SSE update #%d — streamed at %s</p>`,
						i,
						time.Now().Format("15:04:05"),
					),
				)
				if err != nil {
					return // client went away
				}
				flusher.Flush()
			case <-keepalive.C:
				if _, err := io.WriteString(w, ": ping\n\n"); err != nil {
					return // client went away
				}
				flusher.Flush()
			}
		}
	})

	// Mock Datastar action endpoint. For non-SSE responses the Datastar v1
	// runtime reads the patch target from response headers: the body is the
	// elements payload, Datastar-Selector picks the target, Datastar-Mode the
	// merge mode. Without the headers the runtime drops id-less fragments
	// with a "PatchElementsNoTargetsFound" console warning.
	mux.HandleFunc("/api/datastar/action", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Datastar-Selector", "#datastar-action-result")
		w.Header().Set("Datastar-Mode", "inner")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(
			w,
			`<div class="rounded-lg border border-green-200 dark:border-green-800 bg-green-50 dark:bg-green-900/20 p-4">
			<p class="text-sm text-green-800 dark:text-green-200">Datastar action received (mock endpoint).</p>
		</div>`,
		)
	})

	// Wire demo: one endpoint serving both transports via the library's own
	// middleware. htmx targets client-side via hx-target; the Datastar runtime
	// is told where to patch via response headers (Datastar v1.0.2 fetch
	// actions have no target option — see docs/datastar-runtime-facts.md).
	mux.Handle("/api/wire/fragment", wire.Handler(wire.PatchTarget{
		Selector: "#wire-datastar-out",
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		noStore(w)
		componentOr500(w, r, wireFragment(time.Now().Format("15:04:05")))
	})))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/forms":
			renderPage(w, r, "Forms Demo - templ-components", "Complete form showcase with validation", formsDemoPage)
		case "/users":
			renderUsersPage(w, r)
		case "/recipes/dashboard":
			renderPage(w, r, "Dashboard Recipe - templ-components", "Dashboard recipe demo", recipesDashboardPage)
		case "/recipes/settings":
			renderPage(w, r, "Settings Recipe - templ-components", "Settings recipe demo", recipesSettingsPage)
		case "/recipes/login":
			renderPage(w, r, "Login Recipe - templ-components", "Login card recipe demo", recipesLoginPage)
		case "/recipes/auth":
			renderPage(w, r, "Auth Layout Recipe - templ-components", "Auth layout recipe demo", recipesAuthPage)
		case "/":
			transport := parseDemoTransport(r.URL.Query().Get("transport"))
			renderPage(w, r, "templ-components Demo", "Showcase of all templ-components", func(props layout.PageProps) templ.Component {
				return demoPage(props, transport)
			})
		default:
			http.NotFound(w, r)
		}
	})

	return mux
}

// newServer wraps the mux in the demo HTTP server. The global WriteTimeout
// is safe for the SSE endpoint because that handler clears the deadline
// per-connection via http.NewResponseController.
//
//nolint:exhaustruct_v5 // Demo code - HTTP server for demonstration only
func newServer(handler http.Handler) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
		// WriteTimeout would kill the SSE stream (/api/datastar/stream);
		// that handler clears the deadline per-connection via
		// http.NewResponseController. Keep the global timeout for everything
		// else so stuck clients cannot pin connections forever.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func renderPage(
	w http.ResponseWriter,
	r *http.Request,
	title, description string,
	page func(layout.PageProps) templ.Component,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	props := layout.DefaultPageProps()
	props.Title = title
	props.Description = description
	props.Nonce = "demo-nonce"
	props.CSSPath = "/css/app.css"
	props.HeadContent = demoFonts("demo-nonce")
	if err := page(props).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// noStore marks a response as uncacheable: HTMX GET fragments must never be
// served from a browser/proxy cache or a stale fragment lands in a live swap.
func noStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
}

// componentOr500 renders a templ fragment handler response.
func componentOr500(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
