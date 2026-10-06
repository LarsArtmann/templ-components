package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestErrorRoutesServeStatusAndBody pins the standalone error demo routes:
// each route must answer with its own HTTP status code and a body rendered by
// the errorpage handler machinery (not the demo layout). The 404-page route
// additionally proves the dedicated NotFound404 component serves standalone.
func TestErrorRoutesServeStatusAndBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path         string
		wantStatus   int
		wantContains string
	}{
		{"/errors/400", http.StatusBadRequest, "Bad request"},
		{"/errors/403", http.StatusForbidden, "Access denied"},
		{"/errors/404", http.StatusNotFound, "Page not found"},
		{"/errors/409", http.StatusConflict, "Conflict detected"},
		{"/errors/500", http.StatusInternalServerError, "Something went wrong"},
		{"/errors/503", http.StatusServiceUnavailable, "Service temporarily unavailable"},
		{"/errors/full", http.StatusServiceUnavailable, "Trace: trc_9f3a1c2d"},
		{"/errors/404-page", http.StatusNotFound, "Page not found"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("GET %s status = %d, want %d", tt.path, rec.Code, tt.wantStatus)
			}

			if body := rec.Body.String(); !strings.Contains(body, tt.wantContains) {
				t.Errorf("GET %s body does not contain %q", tt.path, tt.wantContains)
			}
		})
	}
}

// TestErrorRoutesPlaygroundPinsQueryWiring pins the playground's query-param
// contract: sanitized status passthrough, the ErrorMaxWidth clamp (unknown →
// XL), the always-on way-out action back to the components page, and the
// CopyCode wiring (button renders when a code is supplied; its code!=\"\"
// gating is pinned by the errorpage package's own golden/matrix tests).
func TestErrorRoutesPlaygroundPinsQueryWiring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		wantStatus   int
		wantContains []string
	}{
		{
			name:       "full query renders width, way-out, and copy button",
			query:      "?family=transient&status=503&code=MAINT_WINDOW&width=2xl&title=Down",
			wantStatus: http.StatusServiceUnavailable,
			wantContains: []string{
				"max-w-2xl",
				`data-tc-copy="MAINT_WINDOW"`,
				"Back to the error page components",
				"/demo/error-pages",
			},
		},
		{
			name:       "empty query defaults to 200 and XL width",
			query:      "",
			wantStatus: http.StatusOK,
			wantContains: []string{
				"max-w-xl",
				"Back to the error page components",
			},
		},
		{
			name:       "unknown width clamps to XL",
			query:      "?width=banana",
			wantStatus: http.StatusOK,
			wantContains: []string{
				"max-w-xl",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/errors/playground"+tt.query, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("GET /errors/playground%s status = %d, want %d", tt.query, rec.Code, tt.wantStatus)
			}

			body := rec.Body.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(body, want) {
					t.Errorf("GET /errors/playground%s body does not contain %q", tt.query, want)
				}
			}
		})
	}
}
