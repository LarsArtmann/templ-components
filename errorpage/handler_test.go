package errorpage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestFromError(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns transient", func(t *testing.T) {
		t.Parallel()
		props := FromError(nil)
		if props.Family != FamilyTransient {
			t.Errorf("Family = %q, want %q", props.Family, FamilyTransient)
		}
	})

	t.Run("plain error returns transient with message", func(t *testing.T) {
		t.Parallel()
		props := FromError(&testError{msg: "plain error"})
		if props.Family != FamilyTransient {
			t.Errorf("Family = %q, want %q", props.Family, FamilyTransient)
		}
		if props.Message == "" {
			t.Error("expected non-empty message")
		}
	})

	t.Run("extracts code from coded error", func(t *testing.T) {
		t.Parallel()
		inner := &testCodedError{msg: "coded", code: "db.timeout"}
		props := FromError(inner)
		if props.Code != "db.timeout" {
			t.Errorf("Code = %q, want %q", props.Code, "db.timeout")
		}
	})

	t.Run("extracts family from classified error", func(t *testing.T) {
		t.Parallel()
		err := &testClassifiedError{msg: "bad input", family: FamilyRejection}
		props := FromError(err)
		if props.Family != FamilyRejection {
			t.Errorf("Family = %q, want %q", props.Family, FamilyRejection)
		}
	})

	t.Run("extracts context from contextual error", func(t *testing.T) {
		t.Parallel()
		err := &testContextualError{
			msg: "db error",
			ctx: map[string]string{"host": "db.internal"},
		}
		props := FromError(err)
		if len(props.Context) != 1 {
			t.Fatalf("expected 1 context pair, got %d", len(props.Context))
		}
		if props.Context[0].Key != "host" {
			t.Errorf("Context[0].Key = %q, want %q", props.Context[0].Key, "host")
		}
	})

	t.Run("extracts cause chain", func(t *testing.T) {
		t.Parallel()
		inner := &testError{msg: "leaf"}
		outer := &testError{msg: "wrapper", cause: inner}
		props := FromError(outer)
		if len(props.CauseChain) != 1 {
			t.Fatalf("expected 1 cause, got %d", len(props.CauseChain))
		}
		if props.CauseChain[0].Message != "leaf" {
			t.Errorf("CauseChain[0].Message = %q, want %q", props.CauseChain[0].Message, "leaf")
		}
	})

	t.Run("sets timestamp", func(t *testing.T) {
		t.Parallel()
		props := FromError(&testError{msg: "err"})
		if props.Timestamp == "" {
			t.Error("expected non-empty timestamp")
		}
	})
}

func TestPreBuiltConstructors(t *testing.T) {
	t.Parallel()

	t.Run("NotFound has correct family and code", func(t *testing.T) {
		t.Parallel()
		props := NotFound()
		if props.Family != FamilyRejection {
			t.Errorf("Family = %q, want %q", props.Family, FamilyRejection)
		}
		if props.Code != "page.not_found" {
			t.Errorf("Code = %q, want %q", props.Code, "page.not_found")
		}
		if props.WayOutHref != "/" {
			t.Errorf("WayOutHref = %q, want %q", props.WayOutHref, "/")
		}
	})

	t.Run("Forbidden has rejection family", func(t *testing.T) {
		t.Parallel()
		props := Forbidden()
		if props.Family != FamilyRejection {
			t.Errorf("Family = %q", props.Family)
		}
		if props.Code != "access.forbidden" {
			t.Errorf("Code = %q", props.Code)
		}
	})

	t.Run("BadRequest has rejection family and custom message", func(t *testing.T) {
		t.Parallel()
		props := BadRequest("Invalid email")
		if props.Family != FamilyRejection {
			t.Errorf("Family = %q", props.Family)
		}
		if props.Message != "Invalid email" {
			t.Errorf("Message = %q", props.Message)
		}
	})

	t.Run("BadRequest empty message gets default", func(t *testing.T) {
		t.Parallel()
		props := BadRequest("")
		if props.Message == "" {
			t.Error("expected default message for empty input")
		}
	})

	t.Run("ConflictError has conflict family", func(t *testing.T) {
		t.Parallel()
		props := ConflictError("Version mismatch")
		if props.Family != FamilyConflict {
			t.Errorf("Family = %q", props.Family)
		}
		if props.Message != "Version mismatch" {
			t.Errorf("Message = %q", props.Message)
		}
	})

	t.Run("ServiceUnavailable has transient family", func(t *testing.T) {
		t.Parallel()
		props := ServiceUnavailable()
		if props.Family != FamilyTransient {
			t.Errorf("Family = %q", props.Family)
		}
		if props.Code != "service.unavailable" {
			t.Errorf("Code = %q", props.Code)
		}
	})

	t.Run("InternalError has infrastructure family", func(t *testing.T) {
		t.Parallel()
		props := InternalError()
		if props.Family != FamilyInfrastructure {
			t.Errorf("Family = %q", props.Family)
		}
		if props.Code != "internal.error" {
			t.Errorf("Code = %q", props.Code)
		}
	})
}

func TestPreBuiltConstructorsRender(t *testing.T) {
	t.Parallel()

	constructors := []struct {
		name  string
		props ErrorPageProps
	}{
		{"NotFound", NotFound()},
		{"Forbidden", Forbidden()},
		{"BadRequest", BadRequest("bad input")},
		{"ConflictError", ConflictError("version clash")},
		{"ServiceUnavailable", ServiceUnavailable()},
		{"InternalError", InternalError()},
	}

	for _, tc := range constructors {
		t.Run(tc.name+" renders without panic", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, ErrorPage(tc.props))
			utils.AssertContains(t, output, "min-h-screen")
			utils.AssertContains(t, output, tc.props.Code)
		})
	}
}

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("writes correct HTTP status code", func(t *testing.T) {
		t.Parallel()
		err := &testClassifiedError{msg: "bad input", family: FamilyRejection}
		handler := ErrorHandler(err, ErrorHandlerConfig{})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("sets content-type header", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "err"}, ErrorHandlerConfig{})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		ct := rec.Header().Get("Content-Type")
		if !strings.Contains(ct, "text/html") {
			t.Errorf("Content-Type = %q, want text/html", ct)
		}
	})

	t.Run("renders HTML body", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(&testError{msg: "something broke"}, ErrorHandlerConfig{})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body := rec.Body.String()
		if !strings.Contains(body, "something broke") {
			t.Error("response body should contain error message")
		}
	})

	t.Run("infrastructure error returns 503", func(t *testing.T) {
		t.Parallel()
		err := &testClassifiedError{msg: "down", family: FamilyInfrastructure}
		handler := ErrorHandler(err, ErrorHandlerConfig{})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
	})

	t.Run("override can customize props", func(t *testing.T) {
		t.Parallel()
		called := false
		handler := ErrorHandler(&testError{msg: "err"}, ErrorHandlerConfig{
			Override: func(err error, props ErrorPageProps) *ErrorPageProps {
				called = true
				props.Title = "Custom Title"
				return &props
			},
		})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if !called {
			t.Error("expected Override to be called")
		}
		body := rec.Body.String()
		if !strings.Contains(body, "Custom Title") {
			t.Error("response body should contain custom title")
		}
	})
}

func TestWriteErrorPage(t *testing.T) {
	t.Parallel()

	t.Run("writes correct status and body", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()

		WriteErrorPage(rec, req, http.StatusNotFound, NotFound(), "")

		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "Page not found") {
			t.Error("body should contain title")
		}
	})
}

func TestWriteError(t *testing.T) {
	t.Parallel()

	t.Run("convenience wrapper works", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/test", nil)
		rec := httptest.NewRecorder()

		WriteError(rec, req, &testClassifiedError{msg: "bad", family: FamilyConflict}, "nonce-123")

		if rec.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusConflict)
		}
	})
}

type testClassifiedError struct {
	msg    string
	family Family
}

func (e *testClassifiedError) Error() string       { return e.msg }
func (e *testClassifiedError) ErrorFamily() Family { return e.family }

type testContextualError struct {
	msg string
	ctx map[string]string
}

func (e *testContextualError) Error() string                   { return e.msg }
func (e *testContextualError) ErrorContext() map[string]string { return e.ctx }
