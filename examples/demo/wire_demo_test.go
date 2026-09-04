package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// TestWireDemoTransportToggle pins the ?transport= selector contract: the
// chosen dialect renders, the other one does not, an unknown value falls back
// to the both-dialect default, and the segmented control offers all three.
func TestWireDemoTransportToggle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		transport    string
		wantHTMX     bool
		wantDatastar bool
	}{
		{name: "htmx param renders only the htmx dialect", transport: "htmx", wantHTMX: true, wantDatastar: false},
		{name: "datastar param renders only the datastar dialect", transport: "datastar", wantHTMX: false, wantDatastar: true},
		{name: "both param renders both dialects", transport: "both", wantHTMX: true, wantDatastar: true},
		{name: "unknown param falls back to both", transport: "webcomponents", wantHTMX: true, wantDatastar: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			resp, err := server.Client().Get(server.URL + "/?transport=" + tt.transport)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			html := string(body)

			hasHTMX := strings.Contains(html, `hx-get="/api/wire/fragment"`)
			hasDatastar := strings.Contains(html, `data-on:click="@get(&#39;/api/wire/fragment&#39;)"`)

			if hasHTMX != tt.wantHTMX {
				t.Errorf("htmx dialect presence = %v, want %v", hasHTMX, tt.wantHTMX)
			}

			if hasDatastar != tt.wantDatastar {
				t.Errorf("datastar dialect presence = %v, want %v", hasDatastar, tt.wantDatastar)
			}

			for _, want := range []string{
				`?transport=both`, `?transport=htmx`, `?transport=datastar`,
			} {
				if !strings.Contains(html, want) {
					t.Errorf("transport selector missing option link %q", want)
				}
			}
		})
	}
}

// TestWireValidateEndpoint verifies the server-validation endpoint contract:
// wire.Handler routes Datastar callers to #wire-validate-out, and the verdict
// fragment reflects the submitted value.
func TestWireValidateEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         string
		datastarCall  bool
		wantInBody    string
		wantNotInBody string
		wantSelector  string
	}{
		{
			name:         "empty value asks for input",
			value:        "",
			wantInBody:   "Type an email address",
			wantSelector: "",
		},
		{
			name:         "email-looking value passes",
			value:        "you@example.com",
			wantInBody:   "Looks like an email address",
			wantSelector: "",
		},
		{
			name:          "value without @ fails",
			value:         "not-an-email",
			wantInBody:    "not an email address",
			wantNotInBody: "Looks like",
			wantSelector:  "",
		},
		{
			name:         "datastar caller gets response-header targeting",
			value:        "nope",
			datastarCall: true,
			wantInBody:   "not an email address",
			wantSelector: "#wire-validate-out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/api/wire/validate?value="+url.QueryEscape(tt.value), nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.datastarCall {
				req.Header.Set(wire.HeaderDatastarRequest, "true")
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if got := resp.Header.Get(wire.HeaderDatastarSelector); got != tt.wantSelector {
				t.Errorf("Datastar-Selector = %q, want %q", got, tt.wantSelector)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(body), tt.wantInBody) {
				t.Errorf("body %q missing %q", body, tt.wantInBody)
			}

			if tt.wantNotInBody != "" && strings.Contains(string(body), tt.wantNotInBody) {
				t.Errorf("body %q must not contain %q", body, tt.wantNotInBody)
			}
		})
	}
}

// TestWireDemoValidateInputDialects verifies the validation input renders the
// typed wire contract under htmx and the bound-signal escape hatch under
// Datastar.
func TestWireDemoValidateInputDialects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		transport       string
		wantContains    []string
		wantNotContains []string
	}{
		{
			name:      "htmx uses the typed wire contract",
			transport: "htmx",
			wantContains: []string{
				`hx-get="/api/wire/validate"`,
				`hx-trigger="change"`,
				`hx-target="#wire-validate-out"`,
				`name="value"`,
			},
			wantNotContains: []string{"data-bind"},
		},
		{
			name:      "datastar uses the Attrs escape hatch with a bound signal",
			transport: "datastar",
			wantContains: []string{
				`data-bind:value`,
				`data-on:change="@get(&#39;/api/wire/validate?value=&#39; + encodeURIComponent($value || &#39;&#39;))"`,
			},
			wantNotContains: []string{"hx-get"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			resp, err := server.Client().Get(server.URL + "/?transport=" + tt.transport)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			html := string(body)
			for _, want := range tt.wantContains {
				if !strings.Contains(html, want) {
					t.Errorf("page missing %q", want)
				}
			}

			for _, banned := range tt.wantNotContains {
				if strings.Contains(html, banned) {
					t.Errorf("page must not contain %q", banned)
				}
			}
		})
	}
}
