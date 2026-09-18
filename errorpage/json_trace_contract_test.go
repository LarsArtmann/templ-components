package errorpage

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errorfamily "github.com/larsartmann/go-error-family"
)

// tracedCodedError mirrors what go-error-family/bridge promotes from
// samber/oops: a family classification, an error code, and a correlation
// trace ID. One error drives BOTH render paths (HTML page and JSON body) so
// the parity guard can compare them field for field.
type tracedCodedError struct {
	err    error
	family errorfamily.Family
	code   string
	trace  string
}

func (e *tracedCodedError) Error() string { return e.err.Error() }

func (e *tracedCodedError) ErrorFamily() errorfamily.Family { return e.family }

func (e *tracedCodedError) ErrorCode() string { return e.code }

func (e *tracedCodedError) Trace() string { return e.trace }

// TestJSONTraceContract pins the `trace` field's presence/omission behavior:
// a traced error surfaces its correlation ID in the JSON error body; a plain
// error must omit the key entirely (omitempty), so consumers can
// distinguish "no trace" from "empty trace" without heuristics.
func TestJSONTraceContract(t *testing.T) {
	t.Parallel()

	t.Run("traced error carries trace field", func(t *testing.T) {
		t.Parallel()

		handler := ErrorHandler(&tracedCodedError{
			err:    errors.New("boom"),
			family: errorfamily.Conflict,
			code:   "resource.conflict",
			trace:  "trc_abc123",
		}, ErrorHandlerConfig{JSON: true})

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))

		var resp errorResponse
		if err := json.UnmarshalRead(rec.Body, &resp); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if resp.Trace != "trc_abc123" {
			t.Errorf("JSON trace = %q, want %q", resp.Trace, "trc_abc123")
		}
	})

	t.Run("plain error omits trace key", func(t *testing.T) {
		t.Parallel()

		handler := ErrorHandler(errors.New("boom"), ErrorHandlerConfig{JSON: true})

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))

		body := rec.Body.String()
		if strings.Contains(body, `"trace"`) {
			t.Error("JSON body must omit the trace key entirely for untraced errors")
		}
	})
}

// TestChipsJSONParity guards the wire-level contract between the two error
// render paths: for the SAME error, the HTML page's chip row (HTTP status,
// code, trace footer) and the JSON error response (status, code, trace
// field) must agree. A drift here means an API consumer and a browser user
// see different facts about the same failure.
func TestChipsJSONParity(t *testing.T) {
	t.Parallel()

	err := &tracedCodedError{
		err:    errors.New("stale write detected"),
		family: errorfamily.Conflict,
		code:   "resource.conflict",
		trace:  "trc_parity_01",
	}

	htmlRec := httptest.NewRecorder()
	ErrorHandler(err, ErrorHandlerConfig{}).ServeHTTP(
		htmlRec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil),
	)

	jsonRec := httptest.NewRecorder()
	ErrorHandler(err, ErrorHandlerConfig{JSON: true}).ServeHTTP(
		jsonRec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil),
	)

	var resp errorResponse
	if jsonErr := json.UnmarshalRead(jsonRec.Body, &resp); jsonErr != nil {
		t.Fatalf("failed to decode JSON: %v", jsonErr)
	}

	html := htmlRec.Body.String()

	// Status: the HTTP chip (`HTTP 409`) must match the JSON response's
	// written status code.
	if htmlRec.Code != jsonRec.Code {
		t.Errorf("HTML status = %d, JSON status = %d — render paths disagree", htmlRec.Code, jsonRec.Code)
	}

	if !strings.Contains(html, "HTTP 409") {
		t.Errorf("HTML chip row missing `HTTP %d` for conflict family", jsonRec.Code)
	}

	// Code: the HTML `<code>` chip text must equal the JSON `code` field.
	if resp.Code != "resource.conflict" {
		t.Errorf("JSON code = %q, want %q", resp.Code, "resource.conflict")
	}

	if !strings.Contains(html, resp.Code) {
		t.Errorf("HTML code chip missing %q present in JSON", resp.Code)
	}

	// Trace: the HTML footer trace must equal the JSON `trace` field.
	if resp.Trace != "trc_parity_01" {
		t.Errorf("JSON trace = %q, want %q", resp.Trace, "trc_parity_01")
	}

	if !strings.Contains(html, resp.Trace) {
		t.Errorf("HTML trace footer missing %q present in JSON", resp.Trace)
	}
}
