package visualtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/utils/wire"
)

// CSP-mode E2E spike (2026-11-07 Datastar plan T03, gate G1). Proves the
// PINNED runtime's CSP mode against the library's real attribute flows under
// a strict Content-Security-Policy with NO 'unsafe-eval': the page carries
// <html data-nonce="…"> plus a nonce'd runtime script, so the runtime reads
// the nonce at init and executes every expression as a nonce'd inline script
// instead of Function() eval (docs/datastar-runtime-facts.md §CSP).
//
// A securitypolicyviolation collector turns any eval attempt — or any script
// the CSP blocks — into a loud failure instead of a silent DOM stall, and the
// flow assertions fail if the runtime silently stops working. GO unblocks the
// demo/website CSP flip (plan T06); NO-GO means the demo keeps 'unsafe-eval'
// and the verdict is documented in runtime-facts §CSP. Evaluation-cost
// numbers ride the plan's T04 probe, not this suite.
//
// Flows covered (the library's real attribute set): a wire @get fragment
// patch via response headers, raw expression evaluation with signal
// assignment + data-show, a wire debounced input, the indicator signal
// lifecycle around a slow fetch, on-intersect reveal, on-interval polling, a
// wire ContentTypeForm POST (HTML5 validation path + form serialization), a
// LiveRegion SSE stream, and the SSEErrorHandling 500 path.
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.

// cspSpikeHits counts endpoint hits so single-fire guarantees (debounce) and
// repeat-fire guarantees (interval, intersect) are asserted server-side, not
// inferred from DOM side effects.
type cspSpikeHits struct {
	get, debounce, interval, intersect, slow, post atomic.Int64
}

func cspSpikeNonce(t *testing.T) string {
	t.Helper()

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("generate nonce: %v", err)
	}

	return hex.EncodeToString(buf)
}

// cspModeE2EServer serves the flow fixture. strict=true adds the CSP header
// (nonce-only script-src, no 'unsafe-eval'), the <html data-nonce> activation
// attribute, and nonce'd runtime/toast scripts — CSP mode. strict=false serves
// the same page WITHOUT the header and WITHOUT data-nonce — the runtime's
// default Function()-eval mode, which is the control arm of the T04 probe.
func cspModeE2EServer(t *testing.T, strict bool) (*httptest.Server, *cspSpikeHits) {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	nonce := cspSpikeNonce(t)
	hits := &cspSpikeHits{}

	csp := ""
	if strict {
		// Strict script policy: the page nonce and NOTHING else. No
		// 'unsafe-eval', no 'unsafe-inline', no 'strict-dynamic' — a nonce'd
		// inline script is the only allowed script source, which is exactly the
		// contract the runtime's CSP mode must satisfy. 'report-sample' names
		// the blocked expression in violation samples.
		csp = "default-src 'none'; script-src 'nonce-" + nonce + "' 'report-sample'; style-src 'self'; " +
			"connect-src 'self'; img-src 'self' data:; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
	}

	patchHeaders := func(w http.ResponseWriter, selector, mode string) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Datastar-Selector", selector)
		w.Header().Set("Datastar-Mode", mode)
	}

	getAttrs := wire.Action{Transport: wire.TransportDatastar, URL: "/api/csp-get"}.Attributes()
	postAttrs := wire.Action{
		Transport:   wire.TransportDatastar,
		Method:      wire.MethodPost,
		Event:       wire.EventSubmit,
		URL:         "/api/csp-post",
		ContentType: wire.ContentTypeForm,
	}.Attributes()
	debounceAttrs := wire.Action{
		Transport:  wire.TransportDatastar,
		Event:      wire.EventInput,
		DebounceMS: 300,
		URL:        "/api/csp-debounce",
	}.Attributes()

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, component := range []templ.Component{
			datastar.SDKScript(datastar.SDKScriptProps{
				Nonce: nonceFor(strict, nonce),
				Src:   "/datastar.js",
			}),
			feedback.ToastContainer(nonceFor(strict, nonce)),
			datastar.SSEErrorHandling(datastar.SSEErrorHandlingConfig{Nonce: nonceFor(strict, nonce)}),
		} {
			if err := component.Render(ctx, w); err != nil {
				return err
			}
		}

		liveProps := datastar.DefaultLiveRegionProps()
		liveProps.BaseProps.Nonce = nonceFor(strict, nonce)
		liveProps.ID = "csp-live"
		liveProps.URL = "/api/csp-good-sse"
		liveProps.AutoStart = true

		if err := datastar.LiveRegion(liveProps).Render(ctx, w); err != nil {
			return err
		}

		// Flow 1: wire @get → response-header-targeted patch.
		if _, err := io.WriteString(w,
			`<div id="csp-get-region"><button id="csp-get-btn" type="button"`); err != nil {
			return err
		}

		if err := templruntime.RenderAttributes(ctx, w, w, "button", getAttrs); err != nil {
			return err
		}

		// Flow 2: raw expression evaluation — signal declaration, negation
		// assignment on click, data-show binding.
		// Flow 3: wire debounce modifier on input.
		if _, err := io.WriteString(w,
			`>load fragment</button><div id="csp-get-target">idle</div></div>`+
				`<div id="csp-show" data-signals-cspvisible="false" data-show="$cspvisible">visible text</div>`+
				`<button id="csp-toggle-btn" type="button" data-on:click="$cspvisible = !$cspvisible">toggle</button>`+
				`<input id="csp-debounce-input" `); err != nil {
			return err
		}

		if err := templruntime.RenderAttributes(ctx, w, w, "input", debounceAttrs); err != nil {
			return err
		}

		// Flow 4: indicator signal lifecycle around a slow fetch.
		// Flow 5: on-intersect reveal. Flow 6: on-interval polling.
		if _, err := io.WriteString(w,
			` /><div id="csp-debounce-out">unset</div>`+
				`<div id="csp-slow-region"><button id="csp-slow-btn" type="button" data-indicator:cspbusy data-on:click="@get('/api/csp-slow')">slow load</button>`+
				`<span id="csp-busy-cue" data-show="$cspbusy">loading…</span><div id="csp-slow-done">pending</div></div>`+
				`<div id="csp-intersect" data-on-intersect="@get('/api/csp-intersect')">reveal region</div>`+
				`<div id="csp-interval" data-on-interval__duration.1s="@get('/api/csp-interval')">ticks: </div>`+
				`<input id="csp-probe-input" data-signals-cspprobe="0" data-on:input="$cspprobe = $cspprobe + 1" />`); err != nil {
			return err
		}

		// Flow 7: wire ContentTypeForm POST on a real form (form+submit
		// auto-preventDefault, HTML5 validation gate, field serialization).
		if _, err := io.WriteString(w, `<form id="csp-post-form"`); err != nil {
			return err
		}

		if err := templruntime.RenderAttributes(ctx, w, w, "form", postAttrs); err != nil {
			return err
		}

		// Flows 8+9: SSE good stream (LiveRegion above) and the SSE error
		// path via a failing fetch.
		if _, err := io.WriteString(w,
			`><input name="cspq" value="csp-value"><button id="csp-post-btn" type="submit">send</button></form>`+
				`<div id="csp-post-out">pending</div>`+
				`<div id="csp-live-child">initial</div>`+
				`<button id="csp-bad-btn" type="button" data-on:click="@get('/api/csp-bad')">boom</button>`); err != nil {
			return err
		}

		return nil
	})

	// Hand-rolled shell (not layout.Base): the spike must control EVERY
	// script tag (nonce'd) and place data-nonce on <html> — the layout-level
	// plumbing for that attribute is plan T05, gated BEHIND this spike.
	head := `<meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">` +
		`<title>Datastar CSP-mode spike — templ-components</title><link rel="stylesheet" href="/app.css">` +
		`<script` + nonceAttrFor(strict, nonce) + `>` +
		`window.__cspViolations=[];` +
		`document.addEventListener('securitypolicyviolation',function(e){window.__cspViolations.push(e.effectiveDirective+' sample='+(e.sample||'')+' at '+(e.sourceFile||'')+':'+e.lineNumber+':'+e.columnNumber+' status='+e.statusCode);});` +
		`window.__dsReady=false;` +
		`document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});` +
		`</script>`

	page := func(w io.Writer) error {
		htmlOpen := `<!doctype html><html lang="en"`
		if strict {
			htmlOpen += ` data-nonce="` + nonce + `"`
		}

		if _, err := io.WriteString(w, htmlOpen+`><head>`+head+`</head><body>`); err != nil {
			return err
		}

		if err := body.Render(context.Background(), w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</body></html>`)

		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	mux.HandleFunc("/api/csp-bad", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})

	mux.HandleFunc("/api/csp-get", func(w http.ResponseWriter, _ *http.Request) {
		hits.get.Add(1)
		patchHeaders(w, "#csp-get-target", "inner")
		_, _ = io.WriteString(w, `<span>patched under csp</span>`)
	})

	mux.HandleFunc("/api/csp-debounce", func(w http.ResponseWriter, _ *http.Request) {
		hits.debounce.Add(1)
		patchHeaders(w, "#csp-debounce-out", "inner")
		_, _ = io.WriteString(w, `debounced`)
	})

	mux.HandleFunc("/api/csp-slow", func(w http.ResponseWriter, _ *http.Request) {
		hits.slow.Add(1)

		// Long enough that the in-flight window is easy to observe, short
		// enough to keep the suite fast.
		time.Sleep(1500 * time.Millisecond)
		patchHeaders(w, "#csp-slow-done", "inner")
		_, _ = io.WriteString(w, `slow done`)
	})

	mux.HandleFunc("/api/csp-intersect", func(w http.ResponseWriter, _ *http.Request) {
		hits.intersect.Add(1)
		patchHeaders(w, "#csp-intersect", "inner")
		_, _ = io.WriteString(w, `seen by intersect`)
	})

	mux.HandleFunc("/api/csp-interval", func(w http.ResponseWriter, _ *http.Request) {
		n := hits.interval.Add(1)
		patchHeaders(w, "#csp-interval", "inner")
		_, _ = fmt.Fprintf(w, `tick %d`, n)
	})

	mux.HandleFunc("/api/csp-post", func(w http.ResponseWriter, r *http.Request) {
		hits.post.Add(1)

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		patchHeaders(w, "#csp-post-out", "inner")
		_, _ = io.WriteString(w, `echo:`+html.EscapeString(r.FormValue("cspq")))
	})

	mux.HandleFunc("/api/csp-good-sse", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)

			return
		}

		_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		select {
		case <-r.Context().Done():
			return
		case <-time.After(400 * time.Millisecond):
		}

		writeSSEPatch(w, "#csp-live-child", "inner", "<span>fresh via csp sse</span>")
		flusher.Flush()

		<-r.Context().Done()
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := page(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if csp != "" {
			w.Header().Set("Content-Security-Policy", csp)
		}

		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)

	return srv, hits
}

// TestDatastarCSPModeBrowser is the G1 gate: every flow must work with ZERO
// securitypolicyviolation events under a strict CSP that allows only nonce'd
// scripts.
func TestDatastarCSPModeBrowser(t *testing.T) {
	t.Parallel()

	srv, hits := cspModeE2EServer(t, true)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 120*time.Second)
	defer cancelTimeout()

	// The runtime must initialize at all: CSP mode with a nonempty nonce is
	// the documented activation, and datastar-ready proves the module
	// executed (a nonce mismatch or init throw would leave it false).
	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL),
		pollTrue(`window.__dsReady===true`, chromedp.WithPollingTimeout(15*time.Second)),
	); err != nil {
		t.Fatalf("runtime failed to initialize under strict CSP (datastar-ready never fired): %v", err)
	}

	var violations string

	// flowStep runs one flow in isolation and, on failure, dumps the
	// violation collector, page state, and endpoint counters — a CSP-mode
	// failure must name its cause instead of timing out opaquely.
	flowStep := func(name string, actions ...chromedp.Action[chromedp.Void]) {
		t.Helper()

		t.Logf("csp flow %s: start", name)

		if err := chromedp.Do(ctx, actions...); err != nil {
			diag := ""

			_ = chromedp.Do(ctx, evalExprInto(
				`JSON.stringify({violations: window.__cspViolations, dsReady: window.__dsReady, url: location.href})`,
				&diag,
			))

			t.Fatalf(
				"csp flow %s failed (hits get=%d debounce=%d interval=%d intersect=%d slow=%d post=%d; diag: %s): %v",
				name, hits.get.Load(), hits.debounce.Load(), hits.interval.Load(),
				hits.intersect.Load(), hits.slow.Load(), hits.post.Load(), diag, err,
			)
		}

		t.Logf("csp flow %s: ok", name)
	}

	// Flow 1: wire @get → response-header patch.
	flowStep("get-patch",
		chromedp.Click("#csp-get-btn", chromedp.NodeVisible),
		pollTrue(`document.getElementById('csp-get-target').textContent.includes('patched under csp')`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// Flow 2: expression evaluation + data-show. The declared signal starts
	// false (hidden), the click assignment flips it visible, the second click
	// flips it back — the full assignment cycle must run in CSP mode exactly
	// as in eval mode.
	flowStep("expr-show-toggle",
		pollTrue(`document.getElementById('csp-show').offsetParent===null`, chromedp.WithPollingTimeout(5*time.Second)),
		chromedp.Click("#csp-toggle-btn", chromedp.NodeVisible),
		pollTrue(`document.getElementById('csp-show').offsetParent!==null`, chromedp.WithPollingTimeout(5*time.Second)),
		chromedp.Click("#csp-toggle-btn", chromedp.NodeVisible),
		pollTrue(`document.getElementById('csp-show').offsetParent===null`, chromedp.WithPollingTimeout(5*time.Second)),
	)

	// Flow 3: wire debounce — one keystroke, one request after quiet.
	flowStep("wire-debounce",
		chromedp.SendKeys("#csp-debounce-input", "x"),
		chromedp.Sleep(700*time.Millisecond),
		pollTrue(`document.getElementById('csp-debounce-out').textContent.includes('debounced')`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// Flow 4: indicator visible mid-flight (server sleeps 1.5s), then hidden
	// after completion.
	flowStep("indicator-lifecycle",
		chromedp.Click("#csp-slow-btn", chromedp.NodeVisible),
		pollTrue(`document.getElementById('csp-busy-cue').offsetParent!==null`, chromedp.WithPollingTimeout(3*time.Second)),
		pollTrue(`document.getElementById('csp-slow-done').textContent.includes('slow done') && document.getElementById('csp-busy-cue').offsetParent===null`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// Flow 5: intersect fired on initial visibility.
	flowStep("intersect",
		pollTrue(`document.getElementById('csp-intersect').textContent.includes('seen by intersect')`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// Flow 6: interval ticked, then kept ticking.
	flowStep("interval",
		pollTrue(`document.getElementById('csp-interval').textContent.includes('tick')`, chromedp.WithPollingTimeout(10*time.Second)),
		chromedp.Sleep(1300*time.Millisecond),
	)

	// Flow 7: wire ContentTypeForm POST — form+submit auto-preventDefault
	// (no navigation) and the field value reached the server.
	flowStep("wire-post-form",
		chromedp.Click("#csp-post-btn", chromedp.NodeVisible),
		pollTrue(`document.getElementById('csp-post-out').textContent.includes('echo:csp-value')`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// Flow 8: LiveRegion SSE stream patched the child and cleared the busy
	// cue.
	flowStep("sse-good-stream",
		pollTrue(`document.getElementById('csp-live-child').textContent.includes('fresh via csp sse')`, chromedp.WithPollingTimeout(10*time.Second)),
		pollTrue(`(()=>{const el=document.getElementById('csp-live');return el.getAttribute('aria-busy')===null && !el.hasAttribute('data-tc-live-busy');})()`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// Flow 9: SSE error path — 500 dispatches datastar-fetch error, the
	// component announces and toasts.
	flowStep("sse-error",
		chromedp.Click("#csp-bad-btn", chromedp.NodeVisible),
		pollTrue(`document.getElementById('tc-datastar-announcer').textContent.includes('Stream error')`, chromedp.WithPollingTimeout(10*time.Second)),
		pollTrue(`(()=>{const toasts=document.querySelectorAll('#tc-toast-container > div');for(const t of toasts){if(t.textContent.includes('live stream endpoint returned an error'))return true;}return false;})()`, chromedp.WithPollingTimeout(10*time.Second)),
	)

	// The whole point: ZERO CSP violations anywhere in the session. The
	// collector is captured BEFORE the assert so a failure can report what
	// actually violated.
	if err := chromedp.Do(ctx,
		evalExprInto(`JSON.stringify(window.__cspViolations)`, &violations),
		pollTrue(`window.__cspViolations.length===0`, chromedp.WithPollingTimeout(5*time.Second)),
	); err != nil {
		t.Fatalf("CSP violations under nonce'd runtime: %s (%v)", violations, err)
	}

	if got := hits.debounce.Load(); got != 1 {
		t.Errorf("debounce must fire exactly once per quiet keystroke; got %d", got)
	}

	if got := hits.interval.Load(); got < 2 {
		t.Errorf("on-interval must keep ticking; got %d hits after the observation window", got)
	}

	if got := hits.post.Load(); got != 1 {
		t.Errorf("form POST must fire exactly once; got %d", got)
	}

	if got := hits.intersect.Load(); got != 1 {
		t.Errorf("on-intersect must fire exactly once on initial visibility; got %d", got)
	}

	if violations != "[]" {
		t.Errorf("CSP violations under nonce'd runtime: %s", violations)
	}
}

// nonceFor returns the nonce string when strict, else "" — components apply
// their omit-empty contract and render no nonce attribute in eval mode.
func nonceFor(strict bool, nonce string) string {
	if !strict {
		return ""
	}

	return nonce
}

func nonceAttrFor(strict bool, nonce string) string {
	if !strict {
		return ""
	}

	return ` nonce="` + nonce + `"`
}

// TestDatastarCSPEvalCostProbe (plan T04) measures per-expression evaluation
// cost in BOTH runtime modes: the default Function()-eval arm (strict=false,
// no CSP header) and the nonce'd inline-script arm (strict=true, CSP mode).
// N synchronous input events drive the SAME trivial expression; the harness
// reports µs per evaluation, median of the configured rounds. The numbers are
// recorded in docs/datastar-runtime-facts.md §CSP — adoption guidance, not a
// gate.
func TestDatastarCSPEvalCostProbe(t *testing.T) {
	t.Parallel()

	const (
		rounds        = 3
		evalsPerRound = 300
	)

	probeExpr := fmt.Sprintf(`(() => {
		const el = document.getElementById('csp-probe-input');
		for (let i = 0; i < 10; i++) el.dispatchEvent(new Event('input', {bubbles: true}));
		const t0 = performance.now();
		for (let i = 0; i < %d; i++) el.dispatchEvent(new Event('input', {bubbles: true}));
		const t1 = performance.now();
		return (t1 - t0) * 1000 / %d;
	})()`, evalsPerRound, evalsPerRound)

	measure := func(strict bool) []float64 {
		srv, _ := cspModeE2EServer(t, strict)

		ctx, cancel := newTab(t)
		defer cancel()

		ctx, cancelTimeout := context.WithTimeout(ctx, 120*time.Second)
		defer cancelTimeout()

		runs := make([]float64, 0, rounds)

		for i := 0; i < rounds; i++ {
			var perEval float64

			if err := chromedp.Do(ctx,
				chromedp.Navigate(srv.URL),
				pollTrue(`window.__dsReady===true`, chromedp.WithPollingTimeout(15*time.Second)),
				evalExprInto(probeExpr, &perEval),
			); err != nil {
				t.Fatalf("probe round %d (strict=%v): %v", i, strict, err)
			}

			runs = append(runs, perEval)
		}

		return runs
	}

	median := func(runs []float64) float64 {
		sorted := append([]float64(nil), runs...)
		sort.Float64s(sorted)

		return sorted[len(sorted)/2]
	}

	evalRuns := measure(false)
	cspRuns := measure(true)

	evalMed, cspMed := median(evalRuns), median(cspRuns)

	t.Logf("per-expression cost, %d rounds of %d dispatches: eval mode median %.1f µs/eval %v; CSP mode median %.1f µs/eval %v",
		rounds, evalsPerRound, evalMed, evalRuns, cspMed, cspRuns)

	if evalMed <= 0 || cspMed <= 0 {
		t.Errorf("probe must produce positive timings; eval=%v csp=%v", evalRuns, cspRuns)
	}
}
