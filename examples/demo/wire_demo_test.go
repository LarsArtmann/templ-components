package main

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
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
			wantNotContains: []string{`hx-get="/api/wire/validate`},
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

// TestWireFormEndpointServesBothTransports verifies the dual-transport form
// submission contract: both dialects serialize the form's fields into a
// standard urlencoded POST body, so one ParseForm-driven handler serves both;
// a Datastar caller additionally gets the patch region via response headers
// while an htmx caller does not. Invalid input re-renders the form with a
// ValidationSummary, inline field errors, and the submitted values preserved.
func TestWireFormEndpointServesBothTransports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		datastarRequest bool
		form            url.Values
		wantContains    []string
		wantAbsent      string
	}{
		{
			name:            "datastar caller submits form fields and gets response-header targeting",
			datastarRequest: true,
			form:            url.Values{"name": {"Ada Lovelace"}, "email": {"ada@example.com"}},
			wantContains:    []string{"Subscribed Ada Lovelace (ada@example.com) via datastar.", `data-on:submit="@post(&#39;/api/wire/form&#39;, {contentType: &#39;form&#39;})"`},
		},
		{
			name:            "htmx caller submits the same body without datastar routing headers",
			datastarRequest: false,
			form:            url.Values{"name": {"Grace Hopper"}, "email": {"grace@example.com"}},
			wantContains:    []string{"Subscribed Grace Hopper (grace@example.com) via htmx.", `hx-post="/api/wire/form"`},
		},
		{
			name:            "invalid email re-renders the form with inline errors and preserved values",
			datastarRequest: false,
			form:            url.Values{"name": {"Ada Lovelace"}, "email": {"ada@example"}},
			wantContains: []string{
				"1 error found",
				wireFormEmailBad,
				`value="ada@example"`,
				`aria-invalid="true"`,
			},
			wantAbsent: "Subscribed",
		},
		{
			name:            "missing name re-renders with both errors under datastar",
			datastarRequest: true,
			form:            url.Values{},
			wantContains: []string{
				"2 errors found",
				wireFormNameMissing,
				wireFormEmailBad,
				`data-on:submit="@post(&#39;/api/wire/form&#39;, {contentType: &#39;form&#39;})"`,
			},
			wantAbsent: "Subscribed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				server.URL+"/api/wire/form",
				strings.NewReader(tt.form.Encode()),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			wantSelector := ""
			if tt.datastarRequest {
				req.Header.Set(wire.HeaderDatastarRequest, "true")
				wantSelector = "#wire-form-out"
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}

			if got := resp.Header.Get(wire.HeaderDatastarSelector); got != wantSelector {
				t.Errorf("Datastar-Selector = %q, want %q", got, wantSelector)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(string(body), want) {
					t.Errorf("verdict body missing %q", want)
				}
			}
			if tt.wantAbsent != "" && strings.Contains(string(body), tt.wantAbsent) {
				t.Errorf("error body must not contain %q", tt.wantAbsent)
			}
		})
	}
}

// TestWireDemoFormRendersBothDialects pins the rendered wiring of the
// dual-transport form on the demo page: the htmx dialect carries
// hx-post/hx-trigger/hx-target, the Datastar dialect carries the
// form-serialized submit expression (the pinned v1.0.3 contentType option).
func TestWireDemoFormRendersBothDialects(t *testing.T) {
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
		`hx-post="/api/wire/form"`,
		`hx-trigger="submit"`,
		`hx-target="#wire-form-htmx-region"`,
		`data-on:submit="@post(&#39;/api/wire/form&#39;, {contentType: &#39;form&#39;})"`,
		`id="wire-form-out"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("demo page missing %q", want)
		}
	}
}

// TestWireFilterEndpointServesBothTransports pins the debounced-filter
// endpoint: GET with the query parameter, response-header targeting for
// Datastar callers, and the shared results fragment for both dialects.
func TestWireFilterEndpointServesBothTransports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		query           string
		datastarRequest bool
		wantContains    []string
		wantAbsent      string
	}{
		{
			name:         "empty query lists everything",
			query:        "",
			wantContains: []string{"Alpine.js", "Datastar", "Templ"},
		},
		{
			name:         "query narrows the results case-insensitively",
			query:        "data",
			wantContains: []string{"Datastar"},
			wantAbsent:   "Tailwind",
		},
		{
			name:         "no match renders the empty state",
			query:        "zzzz",
			wantContains: []string{"No matches for"},
			wantAbsent:   "Htmx",
		},
		{
			name:            "datastar caller gets response-header targeting",
			query:           "star",
			datastarRequest: true,
			wantContains:    []string{"Datastar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				server.URL+"/api/wire/filter?q="+url.QueryEscape(tt.query),
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			wantSelector := ""
			if tt.datastarRequest {
				req.Header.Set(wire.HeaderDatastarRequest, "true")
				wantSelector = "#wire-filter-out"
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}

			if got := resp.Header.Get(wire.HeaderDatastarSelector); got != wantSelector {
				t.Errorf("Datastar-Selector = %q, want %q", got, wantSelector)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(string(body), want) {
					t.Errorf("filter body missing %q", want)
				}
			}
			if tt.wantAbsent != "" && strings.Contains(string(body), tt.wantAbsent) {
				t.Errorf("filter body must not contain %q", tt.wantAbsent)
			}
		})
	}
}

// TestWireDemoFilterInputRendersBothDialects pins the demo page's debounced
// filter wiring: the htmx input carries the debounce trigger + target, the
// Datastar input carries the decoded __debounce modifier, and both share
// the results region.
func TestWireDemoFilterInputRendersBothDialects(t *testing.T) {
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
		`hx-get="/api/wire/filter"`,
		`hx-trigger="input changed delay:300ms"`,
		`hx-target="#wire-filter-out"`,
		`data-on:input__debounce.300ms="@get(&#39;/api/wire/filter&#39;, {contentType: &#39;form&#39;})"`,
		`id="wire-filter-out"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("demo page missing %q", want)
		}
	}
}

// TestWireBusyEndpoint pins the busy-state demo endpoint: a deliberately
// slow fragment for both dialects, with dialect-conditional Datastar
// response-header targeting.
func TestWireBusyEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		datastarRequest bool
		wantSelector    string
		wantContains    string
	}{
		{
			name:         "htmx caller gets the done fragment without routing headers",
			wantContains: "Job finished via htmx",
		},
		{
			name:            "datastar caller gets response-header targeting",
			datastarRequest: true,
			wantSelector:    "#wire-busy-datastar-out",
			wantContains:    "Job finished via datastar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, server.URL+"/api/wire/busy", nil)
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

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(body), tt.wantContains) {
				t.Errorf("busy body missing %q", tt.wantContains)
			}
		})
	}
}

// TestWireDemoBusyCardRendersBothDialects pins the busy-state card wiring:
// the htmx button posts to the slow endpoint, the Datastar button carries
// the data-indicator signal, and the Indicator announces via role="status".
func TestWireDemoBusyCardRendersBothDialects(t *testing.T) {
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
		`hx-post="/api/wire/busy"`,
		`hx-target="#wire-busy-htmx-out"`,
		`data-on:click="@post(&#39;/api/wire/busy&#39;)"`,
		`data-indicator:saving`,
		`role="status"`,
		`id="wire-busy-datastar-out"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("demo page missing %q", want)
		}
	}
}

// TestWireUploadEndpoint pins the multipart upload demo: a file travels the
// wired form under both dialects and the endpoint answers with a 200 OK
// result fragment (errors included — the validation rule).
func TestWireUploadEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		datastarRequest bool
		wantSelector    string
		wantContains    string
	}{
		{
			name:         "htmx caller uploads a file",
			wantContains: "Uploaded “demo.txt”",
		},
		{
			name:            "datastar caller uploads a file with response-header targeting",
			datastarRequest: true,
			wantSelector:    "#wire-upload-out",
			wantContains:    "Uploaded “demo.txt”",
		},
		{
			name:         "missing file renders the error fragment as 200 OK",
			wantContains: "No attachment received",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			if tt.name != "missing file renders the error fragment as 200 OK" {
				part, partErr := writer.CreateFormFile("attachment", "demo.txt")
				if partErr != nil {
					t.Fatal(partErr)
				}
				if _, partErr = part.Write([]byte("hello templ-components")); partErr != nil {
					t.Fatal(partErr)
				}
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, server.URL+"/api/wire/upload", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			wantSelector := ""
			if tt.datastarRequest {
				req.Header.Set(wire.HeaderDatastarRequest, "true")
				wantSelector = tt.wantSelector
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}
			if got := resp.Header.Get(wire.HeaderDatastarSelector); got != wantSelector {
				t.Errorf("Datastar-Selector = %q, want %q", got, wantSelector)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(respBody), tt.wantContains) {
				t.Errorf("upload body missing %q", tt.wantContains)
			}
		})
	}
}

// TestWireSearchEndpoint pins the GET search demo: fields travel as query
// parameters on both dialects and the region re-renders with a fresh form.
func TestWireSearchEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		query           string
		datastarRequest bool
		wantSelector    string
		wantContains    string
	}{
		{
			name:         "htmx caller echoes the query",
			query:        "ada",
			wantContains: `GET received q=“ada” via htmx`,
		},
		{
			name:            "datastar caller echoes the query with response-header targeting",
			query:           "grace",
			datastarRequest: true,
			wantSelector:    "#wire-search-datastar-region",
			wantContains:    `GET received q=“grace” via datastar`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				server.URL+"/api/wire/search?q="+url.QueryEscape(tt.query),
				nil,
			)
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

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(respBody), tt.wantContains) {
				t.Errorf("search body missing %q", tt.wantContains)
			}
			freshFormWiring := `data-on:submit="@get(&#39;/api/wire/search&#39;`
			if !tt.datastarRequest {
				freshFormWiring = `hx-get="/api/wire/search"`
			}
			if !strings.Contains(string(respBody), freshFormWiring) {
				t.Errorf("re-rendered region must contain a fresh wired form (%s)", freshFormWiring)
			}
		})
	}
}

// TestWireDemoUploadAndSearchCards pins the page wiring of the upload and
// GET search cards.
func TestWireDemoUploadAndSearchCards(t *testing.T) {
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
		`enctype="multipart/form-data"`,
		`hx-post="/api/wire/upload"`,
		`hx-target="#wire-upload-out"`,
		`data-on:submit="@post(&#39;/api/wire/upload&#39;, {contentType: &#39;form&#39;})"`,
		`hx-get="/api/wire/search"`,
		`hx-target="#wire-search-htmx-region"`,
		`data-on:submit="@get(&#39;/api/wire/search&#39;, {contentType: &#39;form&#39;})"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("demo page missing %q", want)
		}
	}
}
