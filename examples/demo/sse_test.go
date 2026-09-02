package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestDatastarPatchWireFormat pins the Datastar v1 SSE wire format this demo
// emits. Two regressions this guards against (both verified against the
// pinned v1.0.2 runtime bundle):
//
//   - Literal "\n" text instead of real newlines: without a blank line after
//     the datalines the browser never dispatches the event. This exact bug
//     shipped because nothing exercised the handler.
//   - The pre-v1.0 event name "datastar-merge-fragments": the v1 runtime only
//     registers handlers for "datastar-patch-elements"/"datastar-patch-signals"
//     and silently ignores everything else.
func TestDatastarPatchWireFormat(t *testing.T) {
	var b strings.Builder

	err := writeDatastarPatch(&b, "#target", "inner", "<div>\n<p>hello</p>\n</div>")
	if err != nil {
		t.Fatalf("writeDatastarPatch returned error: %v", err)
	}

	want := "event: datastar-patch-elements\n" +
		"data: selector #target\n" +
		"data: mode inner\n" +
		"data: elements <div>\n" +
		"data: elements <p>hello</p>\n" +
		"data: elements </div>\n" +
		"\n"

	if got := b.String(); got != want {
		t.Errorf("writeDatastarPatch output mismatch:\ngot:  %q\nwant: %q", got, want)
	}
}

// TestDatastarPatchPreservesBlankLines pins the <pre> round-trip: interior
// blank lines are emitted as empty-value datalines ("data: elements " with
// trailing space) which the runtime joins back to "\n". Dropping them used to
// silently corrupt preformatted content.
func TestDatastarPatchPreservesBlankLines(t *testing.T) {
	var b strings.Builder

	err := writeDatastarPatch(&b, "#pre", "inner", "<pre>\nline1\n\nline3\n</pre>")
	if err != nil {
		t.Fatalf("writeDatastarPatch returned error: %v", err)
	}

	want := "event: datastar-patch-elements\n" +
		"data: selector #pre\n" +
		"data: mode inner\n" +
		"data: elements <pre>\n" +
		"data: elements line1\n" +
		"data: elements \n" +
		"data: elements line3\n" +
		"data: elements </pre>\n" +
		"\n"

	if got := b.String(); got != want {
		t.Errorf("blank-line round-trip mismatch:\ngot:  %q\nwant: %q", got, want)
	}
}

// FuzzWriteDatastarPatch verifies the writer never emits a corrupt SSE event
// for arbitrary input:
//
//   - the event header is intact and the event ends with exactly one blank
//     line (the dispatch terminator),
//   - no byte line inside the event is empty (an empty line would terminate
//     the event early),
//   - the elements datalines round-trip: joined values equal the payload with
//     outer whitespace trimmed and CR suffixes removed.
func FuzzWriteDatastarPatch(f *testing.F) {
	f.Add("#target", "inner", "<div>\n<p>hello</p>\n</div>")
	f.Add("", "", "")
	f.Add("#pre", "outer", "<pre>\na\n\nb\n</pre>")
	f.Add("#x\ninjected", "inner\nevent: datastar-patch-signals", "<p>ok</p>")
	f.Add("#w", "", "\n\n   \n<p>ws</p>\n\r\n")

	f.Fuzz(func(t *testing.T, selector, mode, html string) {
		var b strings.Builder
		if err := writeDatastarPatch(&b, selector, mode, html); err != nil {
			t.Fatalf("writeDatastarPatch returned error: %v", err)
		}
		out := b.String()

		if !strings.HasPrefix(out, "event: datastar-patch-elements\n") {
			t.Fatalf("event header missing:\n%q", out)
		}
		if !strings.HasSuffix(out, "\n\n") {
			t.Fatalf("event does not end with a blank-line terminator:\n%q", out)
		}

		// No empty byte line inside the event body: an empty line is the SSE
		// event terminator and would truncate the event.
		body := strings.TrimSuffix(out, "\n\n")
		for line := range strings.SplitSeq(body, "\n") {
			if line == "" {
				t.Fatalf("empty byte line inside event (premature terminator):\n%q", out)
			}
		}

		// Round-trip: every payload line appears exactly once, in order, as an
		// elements dataline — including blank lines (empty-value datalines).
		var gotLines []string
		for line := range strings.SplitSeq(body, "\n") {
			if v, ok := strings.CutPrefix(line, "data: elements "); ok {
				gotLines = append(gotLines, v)
			}
		}
		var wantLines []string
		if payload := strings.TrimSpace(html); payload != "" {
			for line := range strings.SplitSeq(payload, "\n") {
				wantLines = append(wantLines, strings.TrimSuffix(line, "\r"))
			}
		}

		if got, want := strings.Join(gotLines, "\n"), strings.Join(wantLines, "\n"); got != want {
			t.Fatalf("elements round-trip mismatch:\ngot:  %q\nwant: %q\nin:\n%q", got, want, out)
		}
	})
}

// TestDatastarStreamEndpoint exercises the real handler over HTTP: headers,
// framing of the first event, and that a client disconnect ends the stream.
// Skipped in -short mode: the stream's first tick fires after 2s.
func TestDatastarStreamEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip("2s ticker: skipped in -short mode")
	}

	server := httptest.NewServer(newMux())
	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/datastar/stream", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got)
	}
	if got := resp.Header.Get("X-Accel-Buffering"); got != "no" {
		t.Errorf("X-Accel-Buffering = %q, want no", got)
	}

	// Read until the first complete event (blank line terminator). The first
	// tick fires after 2s.
	event, err := readUntilBlankLine(resp.Body)
	if err != nil {
		t.Fatalf("reading first SSE event: %v", err)
	}

	for _, want := range []string{
		"event: datastar-patch-elements",
		"data: selector #datastar-live-content",
		"data: mode inner",
		"data: elements <p",
	} {
		if !strings.Contains(event, want) {
			t.Errorf("first SSE event missing %q\nevent:\n%s", want, event)
		}
	}
}

// TestDatastarActionEndpointHeaders: non-SSE HTML responses are patched via
// Datastar-* response headers — without them the runtime drops id-less
// fragments with a console warning.
func TestDatastarActionEndpointHeaders(t *testing.T) {
	server := httptest.NewServer(newMux())
	t.Cleanup(server.Close)

	resp, err := server.Client().Get(server.URL + "/api/datastar/action")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if got := resp.Header.Get("Datastar-Selector"); got != "#datastar-action-result" {
		t.Errorf("Datastar-Selector = %q, want #datastar-action-result", got)
	}
	if got := resp.Header.Get("Datastar-Mode"); got != "inner" {
		t.Errorf("Datastar-Mode = %q, want inner", got)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Datastar action received") {
		t.Error("action response body missing confirmation fragment")
	}
}

// TestHTMXEndpointHeaders is the headers-contract audit for every HTMX demo
// endpoint: the counterparty is HTMX's swap machinery plus any browser cache
// between them, so each fragment response must (a) declare text/html and
// (b) forbid caching — a cached GET fragment would land stale content in a
// live swap (the same contract-vs-counterparty lens the SSE audit applied).
func TestHTMXEndpointHeaders(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    string
		wantIn  []string // substrings the response body must contain
		skipDur bool     // endpoint has an artificial delay
	}{
		{name: "load-more", method: http.MethodGet, path: "/api/items?cursor=1", wantIn: []string{"hx-get", "/api/items?cursor=2", "Item"}},
		{name: "load-more last page settles with end-of-list", method: http.MethodGet, path: "/api/items?cursor=2", wantIn: []string{"You&#39;ve reached the end"}},
		{name: "confirm-delete target fragment", method: http.MethodDelete, path: "/api/items/123", wantIn: []string{"deleted successfully"}},
		{name: "save acknowledges without swapping", method: http.MethodPost, path: "/api/save", wantIn: []string{"Saved."}, skipDur: true},
		{name: "polled region re-arms with next tick", method: http.MethodGet, path: "/api/demo-stats?tick=1", wantIn: []string{"tick=2", "hx-get"}},
		{name: "polled region settles on final tick", method: http.MethodGet, path: "/api/demo-stats?tick=3", wantIn: []string{"Requests served"}},
		{name: "filter dropdown fragment", method: http.MethodGet, path: "/api/users?status=active&sort=name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if testing.Short() && tt.skipDur {
				t.Skip("artificial delay: skipped in -short mode")
			}

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(context.Background(), tt.method, server.URL+tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}
			if got := resp.Header.Get("Content-Type"); got != "text/html; charset=utf-8" {
				t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
			}
			if got := resp.Header.Get("Cache-Control"); got != "no-store" {
				t.Errorf("Cache-Control = %q, want no-store", got)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tt.wantIn {
				if want != "" && !strings.Contains(string(body), want) {
					t.Errorf("response missing %q\nbody:\n%s", want, body)
				}
			}
		})
	}
}

// readUntilBlankLine reads bytes until "\n\n" (the SSE event terminator) and
// returns everything up to and including it.
func readUntilBlankLine(r io.Reader) (string, error) {
	var b strings.Builder
	buf := make([]byte, 512)

	for {
		n, err := r.Read(buf)
		if n > 0 {
			b.Write(buf[:n])
			if strings.Contains(b.String(), "\n\n") {
				return b.String(), nil
			}
		}
		if err != nil {
			return b.String(), err
		}
	}
}
