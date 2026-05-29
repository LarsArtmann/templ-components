package errorpage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	errorfamily "github.com/larsartmann/go-error-family"
)

// Pre-built HTTP error page code constants.
const (
	CodePageNotFound    = "page.not_found"
	CodeAccessForbidden = "access.forbidden"
	CodeBadRequest      = "request.bad_request"
	CodeConflict        = "resource.conflict"
	CodeUnavailable     = "service.unavailable"
	CodeInternalError   = "internal.error"

	msgGoHome = "Go home"
)

// FamilyFromErrorFamily converts a go-error-family Family to an errorpage Family.
// Returns FamilyTransient for unrecognized values.
func FamilyFromErrorFamily(f errorfamily.Family) Family {
	return ParseFamily(f.String())
}

// FromError converts any error into ErrorPageProps.
// Extracts code, family, context, and cause chain from structured errors.
// For go-error-family errors, also extracts Why/Fix defaults.
// Falls back to Transient family for unknown errors.
func FromError(err error) ErrorPageProps {
	if err == nil {
		return ErrorPageProps{Family: FamilyTransient} //nolint:exhaustruct // minimal nil response
	}

	family := familyFromError(err)

	props := ErrorPageProps{ //nolint:exhaustruct // filled incrementally
		Family:     family,
		Message:    cleanMessage(err),
		CauseChain: ExtractCauseChain(err, 5),
		Timestamp:  errorTimestamp(err),
	}

	if classified, ok := errors.AsType[errorfamily.Classified](err); ok {
		ef := classified.ErrorFamily()
		props.Why = ef.DefaultWhy()
		props.Fix = ef.DefaultFix()
	}

	if coded, ok := err.(interface{ ErrorCode() string }); ok {
		props.Code = coded.ErrorCode()
	}

	if ctx, ok := err.(interface{ ErrorContext() map[string]string }); ok {
		props.Context = ContextMap(ctx.ErrorContext())
	}

	return props
}

// familyFromError extracts the family from any error, trying go-error-family first.
func familyFromError(err error) Family {
	if classified, ok := errors.AsType[errorfamily.Classified](err); ok {
		return FamilyFromErrorFamily(classified.ErrorFamily())
	}
	if c, ok := err.(interface{ ErrorFamily() Family }); ok {
		return c.ErrorFamily()
	}
	if c, ok := err.(interface{ ErrorFamily() string }); ok {
		return ParseFamily(c.ErrorFamily())
	}
	return FamilyTransient
}

// cleanMessage returns the clean message from a go-error-family error
// (without the [family:code] prefix), or falls back to err.Error().
func cleanMessage(err error) string {
	type messenger interface{ Message() string }
	if m, ok := err.(messenger); ok {
		return m.Message()
	}
	return err.Error()
}

// errorTimestamp returns the error's own timestamp if available,
// otherwise generates one from time.Now().
func errorTimestamp(err error) string {
	type timestamper interface{ Timestamp() time.Time }
	if ts, ok := err.(timestamper); ok {
		return ts.Timestamp().UTC().Format(time.RFC3339)
	}
	return time.Now().UTC().Format(time.RFC3339)
}

// errorResponse is the JSON structure returned when ErrorHandlerConfig.JSON is true.
type errorResponse struct {
	Family  string `json:"family"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Title   string `json:"title,omitempty"`
	Why     string `json:"why,omitempty"`
	Fix     string `json:"fix,omitempty"`
	Context any    `json:"context,omitempty"`
}

// NotFound returns a 404-style error page.
func NotFound() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyRejection,
		Code:          CodePageNotFound,
		Title:         "Page not found",
		Message:       "The page you're looking for doesn't exist or has been moved.",
		Fix:           "Check the URL for typos or navigate back to the homepage.",
		WayOut:        msgGoHome,
		WayOutHref:    "/",
		ShowTimestamp: true,
	}
}

// Forbidden returns a 403-style error page.
func Forbidden() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyRejection,
		Code:          CodeAccessForbidden,
		Title:         "Access denied",
		Message:       "You don't have permission to access this resource.",
		Fix:           "Contact your administrator if you believe this is an error.",
		WayOut:        msgGoHome,
		WayOutHref:    "/",
		ShowTimestamp: true,
	}
}

// BadRequest returns a 400-style error page.
func BadRequest(message string) ErrorPageProps {
	if message == "" {
		message = "The request was invalid or malformed."
	}
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyRejection,
		Code:          CodeBadRequest,
		Title:         "Bad request",
		Message:       message,
		Fix:           "Check your input and try again.",
		WayOut:        "Go back",
		ShowTimestamp: true,
	}
}

// Conflict returns a 409-style error page.
func Conflict(message string) ErrorPageProps {
	if message == "" {
		message = "A conflict was detected with the current state of the resource."
	}
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyConflict,
		Code:          CodeConflict,
		Title:         "Conflict detected",
		Message:       message,
		Fix:           "Refresh your data and try the operation again.",
		WayOut:        "Go back",
		ShowTimestamp: true,
	}
}

// ServiceUnavailable returns a 503-style error page.
func ServiceUnavailable() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyTransient,
		Code:          CodeUnavailable,
		Title:         "Service temporarily unavailable",
		Message:       "We're performing maintenance or experiencing high traffic.",
		Why:           "This is a temporary issue. No data was lost.",
		Fix:           "Wait a moment and refresh the page.",
		WayOut:        "Retry",
		ShowTimestamp: true,
	}
}

// InternalError returns a 500-style error page.
func InternalError() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyInfrastructure,
		Code:          CodeInternalError,
		Title:         "Something went wrong",
		Message:       "An unexpected error occurred. Our team has been notified.",
		Why:           "This is a system issue, not something you caused.",
		Fix:           "Try again in a few minutes. If the problem persists, contact support.",
		WayOut:        msgGoHome,
		WayOutHref:    "/",
		ShowTimestamp: true,
	}
}

// ErrorHandlerConfig controls how ErrorHandler renders errors.
type ErrorHandlerConfig struct {
	// Nonce is used for CSP-compliant inline scripts.
	Nonce string

	// Override allows per-error customization of the ErrorPageProps
	// before rendering. Return nil to skip rendering (e.g., for custom handling).
	Override func(err error, props ErrorPageProps) *ErrorPageProps

	// HTMLShell wraps the error page in a minimal HTML document with
	// DOCTYPE, html, head, title, and body tags. Use when the error page
	// is served as a standalone HTTP response (not embedded in an existing layout).
	HTMLShell bool

	// JSON renders a JSON error response instead of HTML.
	// The response includes family, code, message, title, why, and fix fields.
	// Use for API endpoints or HTMX error handling.
	JSON bool
}

// ErrorHandler returns an http.Handler that renders a go-error-family
// aware error page. Use it in your HTTP error handling:
//
//	http.HandleFunc("/api/...", func(w http.ResponseWriter, r *http.Request) {
//	    if err := doSomething(); err != nil {
//	        errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{Nonce: nonce}).ServeHTTP(w, r)
//	        return
//	    }
//	    w.WriteHeader(http.StatusOK)
//	})
func ErrorHandler(err error, cfg ErrorHandlerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		props := FromError(err)
		props.ShowTimestamp = true
		if props.Timestamp == "" {
			props.Timestamp = time.Now().UTC().Format(time.RFC3339)
		}

		if cfg.Override != nil {
			if overridden := cfg.Override(err, props); overridden != nil {
				props = *overridden
			}
		}

		if cfg.Nonce != "" {
			props.Nonce = cfg.Nonce
		}

		statusCode := FamilyStatusCode(props.Family)

		if cfg.JSON {
			writeJSONError(w, statusCode, props)
			return
		}

		if cfg.HTMLShell {
			title := props.Title
			if title == "" {
				title = fmt.Sprintf("Error %d", statusCode)
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(statusCode)
			if renderErr := renderWithShell(r.Context(), w, title, props); renderErr != nil {
				slog.Error("error page render failed", "error", renderErr, "original_error", err)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(statusCode)
		renderErr := ErrorPage(props).Render(r.Context(), w) //nolint:contextcheck
		if renderErr != nil {
			slog.Error("error page render failed", "error", renderErr, "original_error", err)
		}
	})
}

// WriteError writes an error page to an http.ResponseWriter.
// Convenience wrapper around ErrorHandler for simpler usage.
func WriteError(w http.ResponseWriter, r *http.Request, err error, nonce string) {
	ErrorHandler(err, ErrorHandlerConfig{Nonce: nonce}).ServeHTTP(w, r) //nolint:exhaustruct // minimal config
}

// WriteErrorPage writes a pre-configured error page with the given HTTP status code.
// Use with the pre-built constructors:
//
//	errorpage.WriteErrorPage(w, r, 404, errorpage.NotFound(), "")
func WriteErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, props ErrorPageProps, nonce string) {
	if nonce != "" {
		props.Nonce = nonce
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	renderErr := ErrorPage(props).Render(r.Context(), w) //nolint:contextcheck
	if renderErr != nil {
		slog.Error("error page render failed", "error", renderErr, "status_code", statusCode)
	}
}

// writeJSONError writes a JSON error response.
func writeJSONError(w http.ResponseWriter, statusCode int, props ErrorPageProps) {
	resp := errorResponse{ //nolint:exhaustruct // Context set conditionally below
		Family:  string(props.Family),
		Code:    props.Code,
		Message: props.Message,
		Title:   props.Title,
		Why:     props.Why,
		Fix:     props.Fix,
	}
	if len(props.Context) > 0 {
		ctx := make(map[string]string, len(props.Context))
		for _, p := range props.Context {
			ctx[p.Key] = p.Value
		}
		resp.Context = ctx
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(resp); err != nil {
		slog.Error("JSON error response encode failed", "error", err)
	}
}

// renderWithShell wraps ErrorPage in a minimal HTML document.
func renderWithShell(ctx context.Context, w io.Writer, title string, props ErrorPageProps) error {
	shell := templ.ComponentFunc(func(_ context.Context, bw io.Writer) error {
		_, _ = fmt.Fprint(bw, `<!DOCTYPE html><html lang="en"><head>`)
		_, _ = fmt.Fprint(bw, `<meta charset="UTF-8">`)
		_, _ = fmt.Fprintf(bw, `<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
		_, _ = fmt.Fprintf(bw, `<title>%s</title>`, htmlEscape(title))
		_, _ = fmt.Fprint(bw, `</head><body>`)
		renderErr := ErrorPage(props).Render(ctx, bw) //nolint:contextcheck // intentional passthrough
		if renderErr != nil {
			return fmt.Errorf("render error page in shell: %w", renderErr)
		}
		_, _ = fmt.Fprint(bw, `</body></html>`)
		return nil
	})
	if err := shell.Render(ctx, w); err != nil {
		return fmt.Errorf("render error page shell: %w", err)
	}
	return nil
}

// htmlEscape escapes a string for safe inclusion in HTML.
func htmlEscape(s string) string {
	var buf []byte
	for _, r := range s {
		switch r {
		case '&':
			buf = append(buf, "&amp;"...)
		case '<':
			buf = append(buf, "&lt;"...)
		case '>':
			buf = append(buf, "&gt;"...)
		case '"':
			buf = append(buf, "&quot;"...)
		case '\'':
			buf = append(buf, "&#39;"...)
		default:
			buf = append(buf, string(r)...)
		}
	}
	return string(buf)
}

// Verify interface compliance.
var _ http.Handler = ErrorHandler(nil, ErrorHandlerConfig{}) //nolint:exhaustruct // type check only
