package errorpage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestErrorAlertCoverage(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		family Family
	}{
		{"rejection", FamilyRejection},
		{"conflict", FamilyConflict},
		{"transient", FamilyTransient},
		{"corruption", FamilyCorruption},
		{"infrastructure", FamilyInfrastructure},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, ErrorAlert(ErrorAlertProps{
				Family:  tt.family,
				Title:   "Error",
				Message: "Something went wrong",
				Fix:     "Try again",
			}))
			utils.AssertContains(t, output, "Error")
			utils.AssertContains(t, output, "Something went wrong")
		})
	}

	t.Run("dismissible renders data-dismiss", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:      FamilyConflict,
			Title:       "Conflict",
			Dismissible: true,
			BaseProps:   utils.BaseProps{Nonce: "n1"},
		}))
		utils.AssertContains(t, output, "data-dismiss")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(DefaultErrorAlertProps()))
		utils.AssertContains(t, output, "role=\"alert\"")
	})
}

func TestErrorDetailCoverage(t *testing.T) {
	t.Parallel()

	t.Run("full detail with context and causes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyCorruption,
			Code:    "ERR_404",
			Title:   "Not Found",
			Message: "Resource missing",
			Fix:     "Check the URL",
			Context: []ContextPair{
				{Key: "request_id", Value: "abc123"},
				{Key: "path", Value: "/api/users"},
			},
			CauseChain: []CauseItem{
				{Message: "database query failed"},
				{Message: "connection refused"},
			},
			Timestamp: "2024-01-01T00:00:00Z",
		}))
		utils.AssertContains(t, output, "Not Found")
		utils.AssertContains(t, output, "ERR_404")
		utils.AssertContains(t, output, "request_id")
		utils.AssertContains(t, output, "database query failed")
	})

	t.Run("with timestamp", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:    FamilyInfrastructure,
			Title:     "Down",
			Timestamp: "2024-06-20",
		}))
		utils.AssertContains(t, output, "2024-06-20")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(DefaultErrorDetailProps()))
		utils.AssertContains(t, output, "role=\"alert\"")
	})
}

func TestWriteJSONErrorCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with context", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		writeJSONError(w, http.StatusConflict, ErrorPageProps{
			Family:  FamilyConflict,
			Code:    "CONFLICT",
			Title:   "Conflict",
			Message: "Resource already exists",
			Why:     "Duplicate key",
			Fix:     "Use a unique key",
			Context: []ContextPair{
				{Key: "resource", Value: "user@example.com"},
			},
		})
		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
		}
		ct := w.Header().Get("Content-Type")
		if ct != "application/json; charset=utf-8" {
			t.Errorf("content-type = %q", ct)
		}
		body := w.Body.String()
		utils.AssertContains(t, body, "conflict")
		utils.AssertContains(t, body, "user@example.com")
	})

	t.Run("without context", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		writeJSONError(w, http.StatusBadRequest, ErrorPageProps{
			Family:  FamilyRejection,
			Code:    "BAD",
			Title:   "Bad Request",
			Message: "Invalid input",
		})
		body := w.Body.String()
		utils.AssertContains(t, body, "rejection")
	})
}

func TestErrorPageCoverage(t *testing.T) {
	t.Parallel()

	t.Run("full error page", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyTransient,
			Code:    "503",
			Title:   "Service Unavailable",
			Message: "We'll be right back",
			Why:     "Database maintenance",
			Fix:     "Wait a few minutes",
			WayOut:  "Refresh in 5 min",
			Context: []ContextPair{
				{Key: "region", Value: "us-east-1"},
			},
		}))
		utils.AssertContains(t, output, "Service Unavailable")
		utils.AssertContains(t, output, "us-east-1")
	})
}
