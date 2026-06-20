package errorpage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stringFamilyError struct {
	msg    string
	family string
}

func (e *stringFamilyError) Error() string       { return e.msg }
func (e *stringFamilyError) ErrorFamily() string { return e.family }

func TestErrorHandlerCoverage(t *testing.T) {
	t.Parallel()

	t.Run("HTMLShell renders full HTML document", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "shell test"}, ErrorHandlerConfig{HTMLShell: true})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body := rec.Body.String()
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("expected HTML shell with DOCTYPE")
		}
		if !strings.Contains(body, "<html") {
			t.Error("expected HTML shell with <html>")
		}
		if !strings.Contains(body, "</html>") {
			t.Error("expected closing </html>")
		}
	})

	t.Run("HTMLShell with empty title uses status code", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "no title"}, ErrorHandlerConfig{HTMLShell: true})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		body := rec.Body.String()
		if !strings.Contains(body, "<title>") {
			t.Error("expected title element")
		}
	})

	t.Run("JSON mode via ErrorHandler", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "json test"}, ErrorHandlerConfig{JSON: true})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		ct := rec.Header().Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("Content-Type = %q, want json", ct)
		}
		if !strings.Contains(rec.Body.String(), "transient") {
			t.Error("expected family in JSON response")
		}
	})

	t.Run("Override returning nil still renders original props", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "skip"}, ErrorHandlerConfig{
			Override: func(_ error, _ ErrorPageProps) *ErrorPageProps {
				return nil
			},
		})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body := rec.Body.String()
		if !strings.Contains(body, "skip") {
			t.Error("expected original error to still render when Override returns nil")
		}
	})

	t.Run("Nonce propagation via config", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "nonce test"}, ErrorHandlerConfig{
			Nonce: "nonce-abc",
			Override: func(_ error, props ErrorPageProps) *ErrorPageProps {
				props.WayOut = "Return now"
				return &props
			},
		})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body := rec.Body.String()
		if !strings.Contains(body, "nonce-abc") {
			t.Error("expected nonce in output")
		}
	})
}

func TestFamilyFromErrorStringInterface(t *testing.T) {
	t.Parallel()

	t.Run("ErrorFamily() string is parsed", func(t *testing.T) {
		t.Parallel()
		err := &stringFamilyError{msg: "conflict err", family: "conflict"}
		family := familyFromError(err)
		if family != FamilyConflict {
			t.Errorf("family = %q, want %q", family, FamilyConflict)
		}
	})

	t.Run("ErrorFamily() string unknown family falls back", func(t *testing.T) {
		t.Parallel()
		err := &stringFamilyError{msg: "unknown", family: "totally-unknown"}
		family := familyFromError(err)
		if family != FamilyTransient {
			t.Errorf("family = %q, want %q (fallback)", family, FamilyTransient)
		}
	})

	t.Run("plain error falls back to transient", func(t *testing.T) {
		t.Parallel()
		family := familyFromError(&testError{msg: "plain"})
		if family != FamilyTransient {
			t.Errorf("family = %q, want %q", family, FamilyTransient)
		}
	})
}

func TestWriteErrorWrapper(t *testing.T) {
	t.Parallel()

	t.Run("writes error page with nonce", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		WriteError(rec, req, &testError{msg: "wrapped"}, "nonce-xyz")
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "wrapped") {
			t.Error("expected error message in body")
		}
	})
}

func TestWriteErrorPageWrapper(t *testing.T) {
	t.Parallel()

	t.Run("writes pre-configured error page", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		WriteErrorPage(rec, req, http.StatusBadRequest, NotFound(), "")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "page.not_found") {
			t.Error("expected code in body")
		}
	})

	t.Run("writes without nonce when empty", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/x", nil)
		rec := httptest.NewRecorder()
		WriteErrorPage(rec, req, http.StatusBadRequest, BadRequest("bad"), "")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
