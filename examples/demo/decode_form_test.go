package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestDecodeFormEndpoints pins the DecodeForm-migrated demo endpoints: GET
// query decoding, POST body decoding, and the malformed-value paths. Every
// handler in main.go now decodes via wire.DecodeForm — these tests keep that
// migration honest (behavior parity with the old ParseForm chains).
func TestDecodeFormEndpoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		form       url.Values
		wantStatus int
		wantIn     string
	}{
		{
			name:       "load-more decodes cursor query",
			method:     http.MethodGet,
			path:       "/api/items?cursor=1",
			wantStatus: http.StatusOK,
			wantIn:     "/api/items?cursor=2",
		},
		{
			name:       "load-more without cursor defaults to first page",
			method:     http.MethodGet,
			path:       "/api/items",
			wantStatus: http.StatusOK,
			wantIn:     "Item",
		},
		{
			name:       "demo-stats decodes tick and advances",
			method:     http.MethodGet,
			path:       "/api/demo-stats?tick=1",
			wantStatus: http.StatusOK,
			wantIn:     "tick=2",
		},
		{
			name:       "demo-stats with garbage tick restarts at 1",
			method:     http.MethodGet,
			path:       "/api/demo-stats?tick=not-a-number",
			wantStatus: http.StatusOK,
			wantIn:     "tick=1",
		},
		{
			name:       "users decodes status and sort",
			method:     http.MethodGet,
			path:       "/api/users?status=active&sort=name",
			wantStatus: http.StatusOK,
			wantIn:     "Alice",
		},
		{
			name:       "wire wizard advances step 0 with valid email",
			method:     http.MethodPost,
			path:       "/api/wire/wizard",
			form:       url.Values{"step": {"0"}, "email": {"ada@example.com"}},
			wantStatus: http.StatusOK,
			wantIn:     "step",
		},
		{
			name:       "wire wizard restarts on malformed step",
			method:     http.MethodPost,
			path:       "/api/wire/wizard",
			form:       url.Values{"step": {"garbage"}},
			wantStatus: http.StatusOK,
			wantIn:     wireWizardInvalid,
		},
		{
			name:       "wire form decodes name and email",
			method:     http.MethodPost,
			path:       "/api/wire/form",
			form:       url.Values{"name": {"Ada"}, "email": {"ada@example.com"}},
			wantStatus: http.StatusOK,
			wantIn:     "Ada",
		},
		{
			name:       "wire form flags missing name",
			method:     http.MethodPost,
			path:       "/api/wire/form",
			form:       url.Values{"email": {"ada@example.com"}},
			wantStatus: http.StatusOK,
			wantIn:     wireFormNameMissing,
		},
		{
			name:       "wire filter decodes q query",
			method:     http.MethodGet,
			path:       "/api/wire/filter?q=go",
			wantStatus: http.StatusOK,
			wantIn:     "go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(newMux())
			t.Cleanup(server.Close)

			var body io.Reader
			if tt.form != nil {
				body = strings.NewReader(tt.form.Encode())
			}

			req, err := http.NewRequestWithContext(t.Context(), tt.method, server.URL+tt.path, body)
			if err != nil {
				t.Fatal(err)
			}
			if tt.form != nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantIn != "" && !strings.Contains(string(got), tt.wantIn) {
				t.Errorf("response missing %q\nbody:\n%s", tt.wantIn, got)
			}
		})
	}
}
