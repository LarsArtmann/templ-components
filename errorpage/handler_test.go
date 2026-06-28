package errorpage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errorfamily "github.com/larsartmann/go-error-family"
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

	t.Run("plain error returns infrastructure with message", func(t *testing.T) {
		t.Parallel()
		props := FromError(&testError{msg: "plain error"})
		if props.Family != FamilyInfrastructure {
			t.Errorf("Family = %q, want %q", props.Family, FamilyInfrastructure)
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

	t.Run("Conflict has conflict family", func(t *testing.T) {
		t.Parallel()
		props := Conflict("Version mismatch")
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
		{"Conflict", Conflict("version clash")},
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

func TestFromErrorWithGoErrorFamily(t *testing.T) {
	t.Parallel()

	t.Run("detects go-error-family rejection", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewRejection("config.invalid", "missing config")
		props := FromError(err)
		if props.Family != FamilyRejection {
			t.Errorf("Family = %q, want %q", props.Family, FamilyRejection)
		}
		if props.Code != "config.invalid" {
			t.Errorf("Code = %q, want %q", props.Code, "config.invalid")
		}
	})

	t.Run("detects go-error-family conflict", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewConflict("resource.conflict", "version clash")
		props := FromError(err)
		if props.Family != FamilyConflict {
			t.Errorf("Family = %q, want %q", props.Family, FamilyConflict)
		}
	})

	t.Run("detects go-error-family transient", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewTransient("db.timeout", "query took too long")
		props := FromError(err)
		if props.Family != FamilyTransient {
			t.Errorf("Family = %q, want %q", props.Family, FamilyTransient)
		}
	})

	t.Run("detects go-error-family corruption", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewCorruption("data.invalid", "parse failed")
		props := FromError(err)
		if props.Family != FamilyCorruption {
			t.Errorf("Family = %q, want %q", props.Family, FamilyCorruption)
		}
	})

	t.Run("detects go-error-family infrastructure", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewInfrastructure("service.down", "no connections")
		props := FromError(err)
		if props.Family != FamilyInfrastructure {
			t.Errorf("Family = %q, want %q", props.Family, FamilyInfrastructure)
		}
	})

	t.Run("extracts Why/Fix from go-error-family defaults", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewCorruption("data.invalid", "parse failed")
		props := FromError(err)
		if props.Why == "" {
			t.Error("expected non-empty Why from go-error-family defaults")
		}
		if props.Fix == "" {
			t.Error("expected non-empty Fix from go-error-family defaults")
		}
	})

	t.Run("extracts context from go-error-family", func(t *testing.T) {
		t.Parallel()
		err := errorfamily.NewTransient("db.timeout", "slow").
			WithContext("host", "db.internal").
			WithContext("port", "5432")
		props := FromError(err)
		if len(props.Context) != 2 {
			t.Fatalf("expected 2 context pairs, got %d", len(props.Context))
		}
	})

	t.Run("extracts cause chain from go-error-family wrapped error", func(t *testing.T) {
		t.Parallel()
		inner := errorfamily.NewTransient("db.timeout", "slow")
		outer := errorfamily.Wrap(inner, errorfamily.Infrastructure, "api.down", "api failed")
		props := FromError(outer)
		if len(props.CauseChain) == 0 {
			t.Fatal("expected non-empty cause chain")
		}
		if props.CauseChain[0].Code != "db.timeout" {
			t.Errorf("CauseChain[0].Code = %q, want %q", props.CauseChain[0].Code, "db.timeout")
		}
	})
}

func TestFamilyFromErrorFamily(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input errorfamily.Family
		want  Family
	}{
		{errorfamily.Rejection, FamilyRejection},
		{errorfamily.Conflict, FamilyConflict},
		{errorfamily.Transient, FamilyTransient},
		{errorfamily.Corruption, FamilyCorruption},
		{errorfamily.Infrastructure, FamilyInfrastructure},
	}
	for _, tc := range tests {
		t.Run(tc.input.String(), func(t *testing.T) {
			t.Parallel()
			got := FamilyFromErrorFamily(tc.input)
			if got != tc.want {
				t.Errorf("FamilyFromErrorFamily(%v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseFamily(t *testing.T) {
	t.Parallel()

	t.Run("valid families including case variants", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			input string
			want  Family
		}{
			{"rejection", FamilyRejection},
			{"conflict", FamilyConflict},
			{"transient", FamilyTransient},
			{"corruption", FamilyCorruption},
			{"infrastructure", FamilyInfrastructure},
			{"REJECTION", FamilyRejection},
			{"Conflict", FamilyConflict},
			{"TRANSIENT", FamilyTransient},
			{"Corruption", FamilyCorruption},
			{"INFRASTRUCTURE", FamilyInfrastructure},
		}
		for _, tc := range tests {
			got := ParseFamily(tc.input)
			if got != tc.want {
				t.Errorf("ParseFamily(%q) = %q, want %q", tc.input, got, tc.want)
			}
		}
	})

	t.Run("unknown returns transient", func(t *testing.T) {
		t.Parallel()
		got := ParseFamily("unknown")
		if got != FamilyTransient {
			t.Errorf("ParseFamily(%q) = %q, want %q", "unknown", got, FamilyTransient)
		}
	})

	t.Run("trims whitespace", func(t *testing.T) {
		t.Parallel()
		got := ParseFamily("  rejection  ")
		if got != FamilyRejection {
			t.Errorf("ParseFamily(%q) = %q, want %q", "  rejection  ", got, FamilyRejection)
		}
	})
}

func TestConflictEmptyMessage(t *testing.T) {
	t.Parallel()

	props := Conflict("")
	if props.Message == "" {
		t.Error("expected default message for empty Conflict input")
	}
}

func TestDefaultErrorAlertProps(t *testing.T) {
	t.Parallel()

	props := DefaultErrorAlertProps()
	if props.Family != FamilyTransient {
		t.Errorf("Family = %q, want %q", props.Family, FamilyTransient)
	}
}

func TestErrorHandlerHTMLShell(t *testing.T) {
	t.Parallel()

	t.Run("HTMLShell wraps in valid HTML", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(
			&testError{msg: "test"},
			ErrorHandlerConfig{HTMLShell: true},
		)
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body := rec.Body.String()
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("expected DOCTYPE in HTML shell")
		}
		if !strings.Contains(body, "<html") {
			t.Error("expected <html> tag in HTML shell")
		}
		if !strings.Contains(body, "<title>") {
			t.Error("expected <title> in HTML shell")
		}
		if !strings.Contains(body, "</html>") {
			t.Error("expected closing </html> in HTML shell")
		}
	})
}

func TestErrorHandlerJSON(t *testing.T) {
	t.Parallel()

	t.Run("JSON renders JSON error response", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(
			errorfamily.NewRejection("config.invalid", "missing config"),
			ErrorHandlerConfig{JSON: true},
		)
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		ct := rec.Header().Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}

		var resp errorResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}
		if resp.Family != "rejection" {
			t.Errorf("Family = %q, want %q", resp.Family, "rejection")
		}
		if resp.Code != "config.invalid" {
			t.Errorf("Code = %q, want %q", resp.Code, "config.invalid")
		}
		if resp.Message == "" {
			t.Error("expected non-empty message")
		}
	})

	t.Run("JSON includes why and fix from go-error-family", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(
			errorfamily.NewCorruption("data.invalid", "broken"),
			ErrorHandlerConfig{JSON: true},
		)
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		var resp errorResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}
		if resp.Why == "" {
			t.Error("expected non-empty why from go-error-family defaults")
		}
		if resp.Fix == "" {
			t.Error("expected non-empty fix from go-error-family defaults")
		}
	})
}

func TestErrorHandlerEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Override customizes props", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(
			&testError{msg: "test"},
			ErrorHandlerConfig{
				Override: func(_ error, props ErrorPageProps) *ErrorPageProps {
					props.Title = "Overridden"
					return &props
				},
			},
		)
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body := rec.Body.String()
		if !strings.Contains(body, "Overridden") {
			t.Error("expected overridden title in response")
		}
	})

	t.Run("nil error returns 503 transient", func(t *testing.T) {
		t.Parallel()
		handler := ErrorHandler(nil, ErrorHandlerConfig{})
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
	})

	t.Run("WriteErrorPage with nonce propagates to props", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rec := httptest.NewRecorder()
		props := ErrorPageProps{
			Family:        FamilyRejection,
			Code:          "test",
			Title:         "Test Error",
			ShowTimestamp: false,
		}
		WriteErrorPage(rec, req, http.StatusBadRequest, props, "test-nonce-abc")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
