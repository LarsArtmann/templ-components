package errorpage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// ---------------------------------------------------------------------------
// writeFallbackError: currently 0% coverage
// ---------------------------------------------------------------------------

func TestWriteFallbackError(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	writeFallbackError(w, http.StatusNotFound)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Errorf("expected text/plain, got %s", ct)
	}
	if body := w.Body.String(); body == "" {
		t.Error("expected non-empty body")
	}
}

// ---------------------------------------------------------------------------
// writeJSONError: expand coverage for context field
// ---------------------------------------------------------------------------

func TestWriteJSONErrorWithContext(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	writeJSONError(w, http.StatusBadRequest, ErrorPageProps{
		Family:  FamilyRejection,
		Code:    "validation_failed",
		Message: "Email is invalid",
		Title:   "Validation Error",
		Why:     "The email format is incorrect",
		Fix:     "Provide a valid email address",
		Context: []ContextPair{
			{Key: "field", Value: "email"},
			{Key: "value", Value: "not-an-email"},
		},
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	body := w.Body.String()
	utils.AssertContainsAll(t, body, "validation_failed", "Email is invalid", "field", "email")
}

// ---------------------------------------------------------------------------
// ErrorPage: all families, full props
// ---------------------------------------------------------------------------

func TestErrorPageAllFamilies(t *testing.T) {
	t.Parallel()
	for _, family := range []Family{
		FamilyRejection, FamilyConflict, FamilyTransient, FamilyCorruption, FamilyInfrastructure,
	} {
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  family,
			Title:   "Error occurred",
			Message: "Something went wrong",
		}))
		utils.AssertContains(t, output, "Something went wrong")
	}
}

func TestErrorPageFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ErrorPage(ErrorPageProps{
		Family:        FamilyRejection,
		Code:          "auth_failed",
		Title:         "Access Denied",
		Message:       "You do not have permission",
		Why:           "Your session has expired",
		Fix:           "Please log in again",
		WayOut:        "Go to login",
		WayOutHref:    "/login",
		Timestamp:     "2024-01-01T00:00:00Z",
		ShowTimestamp: true,
		Context: []ContextPair{
			{Key: "user_id", Value: "12345"},
		},
		CauseChain: []CauseItem{
			{Message: "Session expired"},
		},
		BaseProps: utils.BaseProps{
			ID:        "err-page",
			AriaLabel: "Error page",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="err-page"`,
		"Access Denied",
		"You do not have permission",
		"Please log in again",
		"Go to login",
		"/login",
		"user_id",
		"12345",
		"Session expired",
	)
}

func TestErrorPageWithNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ErrorPage(ErrorPageProps{
		Family: FamilyTransient,
		Title:  "Server Error",
		BaseProps: utils.BaseProps{
			ID: "err-with-nonce",
		},
	}))
	utils.AssertContains(t, output, "Server Error")
}

// ---------------------------------------------------------------------------
// ErrorDetail: all families, full props
// ---------------------------------------------------------------------------

func TestErrorDetailFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ErrorDetail(ErrorDetailProps{
		Family:  FamilyConflict,
		Code:    "duplicate_entry",
		Title:   "Conflict",
		Message: "Item already exists",
		Fix:     "Use a different name",
		Context: []ContextPair{{Key: "id", Value: "42"}},
		BaseProps: utils.BaseProps{
			ID:        "detail-1",
			AriaLabel: "Error detail",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="detail-1"`,
		"Conflict",
		"Item already exists",
		"Use a different name",
	)
}

func TestErrorDetailAllFamilies(t *testing.T) {
	t.Parallel()
	for _, family := range []Family{
		FamilyRejection, FamilyConflict, FamilyTransient, FamilyCorruption, FamilyInfrastructure,
	} {
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  family,
			Title:   "Detail",
			Message: "Info",
		}))
		utils.AssertContains(t, output, "Info")
	}
}

// ---------------------------------------------------------------------------
// ErrorAlert: all families, dismissible, full props
// ---------------------------------------------------------------------------

func TestErrorAlertFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ErrorAlert(ErrorAlertProps{
		Family:      FamilyTransient,
		Title:       "Warning",
		Message:     "Service is degraded",
		Fix:         "Try again in a few minutes",
		Dismissible: true,
		BaseProps: utils.BaseProps{
			ID:        "alert-1",
			AriaLabel: "Service alert",
			Nonce:     "n-123",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="alert-1"`,
		"Service is degraded",
		"Try again in a few minutes",
	)
}

func TestErrorAlertAllFamilies(t *testing.T) {
	t.Parallel()
	for _, family := range []Family{
		FamilyRejection, FamilyConflict, FamilyTransient, FamilyCorruption, FamilyInfrastructure,
	} {
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  family,
			Message: "Alert message",
		}))
		utils.AssertContains(t, output, "Alert message")
	}
}

// ---------------------------------------------------------------------------
// NotFound404: full props, search form, custom links
// ---------------------------------------------------------------------------

func TestNotFound404FullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NotFound404(NotFound404Props{
		Numeral:           "404",
		Title:             "Page Not Found",
		Message:           "The page you are looking for does not exist.",
		SearchAction:      "/search",
		SearchPlaceholder: "Search...",
		SearchInputName:   "q",
		Links:             DefaultNotFoundLinks(),
		LinksTitle:        "Popular Pages",
		GoHomeHref:        "/home",
		GoHomeText:        "Back Home",
		ShowGoBack:        true,
		BaseProps: utils.BaseProps{
			ID:        "nf-404",
			AriaLabel: "404 error page",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="nf-404"`,
		"Page Not Found",
		"Search...",
		"Popular Pages",
		"Back Home",
		"/home",
	)
}

func TestNotFound404NoSearchNoLinks(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NotFound404(NotFound404Props{
		Title:      "Not Found",
		Message:    "Nothing here",
		ShowGoBack: false,
	}))
	utils.AssertContainsAll(t, output, "Not Found", "Nothing here")
}

func TestNotFound404CustomLinks(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NotFound404(NotFound404Props{
		Links: []NotFoundLink{
			{Text: "Help Center", Href: "/help"},
			{Text: "Contact Us", Href: "/contact"},
		},
	}))
	utils.AssertContainsAll(t, output, "Help Center", "Contact Us")
}

// ---------------------------------------------------------------------------
// ErrorHandler: HTMLShell mode, JSON mode
// ---------------------------------------------------------------------------

func TestErrorHandlerHTMLShellMode(t *testing.T) {
	t.Parallel()
	handler := ErrorHandler(nil, ErrorHandlerConfig{
		HTMLShell: true,
		Nonce:     "test-nonce",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
	handler.ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
	utils.AssertContainsAll(t, w.Body.String(), "<!DOCTYPE html>", "<html")
}

func TestErrorHandlerJSONMode(t *testing.T) {
	t.Parallel()
	handler := ErrorHandler(nil, ErrorHandlerConfig{
		JSON: true,
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
	handler.ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

// ---------------------------------------------------------------------------
// WriteErrorPage with explicit status code
// ---------------------------------------------------------------------------

func TestWriteErrorPageExplicitStatus(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
	WriteErrorPage(w, r, http.StatusTeapot, NotFound(), "")
	if w.Code != http.StatusTeapot {
		t.Errorf("expected 418, got %d", w.Code)
	}
	utils.AssertContains(t, w.Body.String(), "not found")
}

// ---------------------------------------------------------------------------
// FamilyStatusCode: all families
// ---------------------------------------------------------------------------

func TestFamilyStatusCodeAllFamilies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		family Family
		code   int
	}{
		{FamilyRejection, http.StatusBadRequest},
		{FamilyConflict, http.StatusConflict},
		{FamilyTransient, http.StatusServiceUnavailable},
		{FamilyCorruption, http.StatusInternalServerError},
		{FamilyInfrastructure, http.StatusServiceUnavailable},
		{Family("unknown"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		if got := FamilyStatusCode(tt.family); got != tt.code {
			t.Errorf("FamilyStatusCode(%q) = %d, want %d", tt.family, got, tt.code)
		}
	}
}

// ---------------------------------------------------------------------------
// Pre-built constructors
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WriteError convenience
// ---------------------------------------------------------------------------

func TestWriteErrorConvenience(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
	WriteError(w, r, nil, "test-nonce")
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestConstructorsAllRender(t *testing.T) {
	t.Parallel()
	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(NotFound()))
		if len(output) == 0 {
			t.Error("empty output")
		}
	})
	t.Run("Forbidden", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(Forbidden()))
		if len(output) == 0 {
			t.Error("empty output")
		}
	})
	t.Run("BadRequest", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(BadRequest("bad input")))
		if len(output) == 0 {
			t.Error("empty output")
		}
	})
	t.Run("Conflict", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(Conflict("duplicate")))
		if len(output) == 0 {
			t.Error("empty output")
		}
	})
	t.Run("ServiceUnavailable", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ServiceUnavailable()))
		if len(output) == 0 {
			t.Error("empty output")
		}
	})
	t.Run("InternalError", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(InternalError()))
		if len(output) == 0 {
			t.Error("empty output")
		}
	})
}
