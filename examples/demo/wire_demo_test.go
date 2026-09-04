package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils/wire"
)

// TestWireFragmentEndpointServesBothTransports verifies the transport-
// branching contract of the shared wire demo endpoint: a Datastar caller
// (marked by the Datastar-Request header) gets the patch region via response
// headers, while an htmx caller relies on client-side hx-target and must NOT
// receive Datastar routing headers.
func TestWireFragmentEndpointServesBothTransports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		datastarRequest bool
		wantSelector    string
		wantMode        string
	}{
		{
			name:            "datastar caller gets response-header targeting",
			datastarRequest: true,
			wantSelector:    "#wire-datastar-out",
			wantMode:        "inner",
		},
		{
			name:            "htmx caller gets no datastar routing headers",
			datastarRequest: false,
			wantSelector:    "",
			wantMode:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/api/wire/fragment", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.datastarRequest {
				req.Header.Set(wire.HeaderDatastarRequest, "true")
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}

			if got := resp.Header.Get(wire.HeaderDatastarSelector); got != tt.wantSelector {
				t.Errorf("Datastar-Selector = %q, want %q", got, tt.wantSelector)
			}

			if got := resp.Header.Get(wire.HeaderDatastarMode); got != tt.wantMode {
				t.Errorf("Datastar-Mode = %q, want %q", got, tt.wantMode)
			}

			if got := resp.Header.Get("Content-Type"); got != "text/html; charset=utf-8" {
				t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(body), "wire contract") {
				t.Error("fragment body missing confirmation text")
			}
		})
	}
}

// TestWireDemoSectionRendersBothDialects verifies the demo page renders the
// same Action shape as htmx attributes on one button and the Datastar
// expression attribute on the other.
func TestWireDemoSectionRendersBothDialects(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newMux())
	t.Cleanup(server.Close)

	resp, err := server.Client().Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)
	for _, want := range []string{
		`hx-get="/api/wire/fragment"`,
		`hx-target="#wire-htmx-out"`,
		`data-on:click="@get(&#39;/api/wire/fragment&#39;)"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("demo page missing %q", want)
		}
	}
}
