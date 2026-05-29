package errorpage

import (
	"net/http"
	"time"
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

// FromError converts any error into ErrorPageProps.
// Extracts code, family, context, and cause chain from structured errors.
// Falls back to Transient family for unknown errors.
func FromError(err error) ErrorPageProps {
	if err == nil {
		return ErrorPageProps{Family: FamilyTransient} //nolint:exhaustruct // minimal nil response
	}

	family := FamilyTransient
	if c, ok := err.(interface{ ErrorFamily() Family }); ok {
		family = c.ErrorFamily()
	}

	props := ErrorPageProps{ //nolint:exhaustruct // filled incrementally
		Family:     family,
		Message:    err.Error(),
		CauseChain: ExtractCauseChain(err, 5),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	if coded, ok := err.(interface{ ErrorCode() string }); ok {
		props.Code = coded.ErrorCode()
	}

	if ctx, ok := err.(interface{ ErrorContext() map[string]string }); ok {
		props.Context = ContextMap(ctx.ErrorContext())
	}

	return props
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

// ConflictError returns a 409-style error page.
func ConflictError(message string) ErrorPageProps {
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(statusCode)
		_ = ErrorPage(props).Render(r.Context(), w) //nolint:contextcheck // templ.Render requires context.Context
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
	_ = ErrorPage(props).Render(r.Context(), w) //nolint:contextcheck // templ.Render requires context.Context
}

// Verify interface compliance.
var _ http.Handler = ErrorHandler(nil, ErrorHandlerConfig{}) //nolint:exhaustruct // type check only
