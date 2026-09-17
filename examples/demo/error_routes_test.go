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
