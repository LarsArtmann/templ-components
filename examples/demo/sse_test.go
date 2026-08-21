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

// TestDatastarStreamEndpoint exercises the real handler over HTTP: headers,
// framing of the first event, and that a client disconnect ends the stream.
func TestDatastarStreamEndpoint(t *testing.T) {
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
