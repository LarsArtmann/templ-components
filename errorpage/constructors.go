package errorpage

// Code is a typed error code string for categorizing error pages.
type Code string

// Pre-built HTTP error page code constants.
const (
	CodePageNotFound    Code = "page.not_found"
	CodeAccessForbidden Code = "access.forbidden"
	CodeBadRequest      Code = "request.bad_request"
	CodeConflict        Code = "resource.conflict"
	CodeUnavailable     Code = "service.unavailable"
	CodeInternalError   Code = "internal.error"

	msgGoHome             = "Go home"
	msgGoBack             = notFound404GoBackText
	titleInternalError    = "Something went wrong"
	msgBadRequestDefault  = "The request was invalid or malformed."
	msgConflictDefault    = "A conflict was detected with the current state of the resource."
	msgServiceUnavailable = "We're performing maintenance or experiencing high traffic."
	msgInternalUnexpected = "An unexpected error occurred. Our team has been notified."
)

// NotFound returns a 404-style error page.
func NotFound() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyRejection,
		Code:          CodePageNotFound,
		Title:         notFound404DefaultTitle,
		Message:       notFound404DefaultMessage,
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
		message = msgBadRequestDefault
	}
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyRejection,
		Code:          CodeBadRequest,
		Title:         "Bad request",
		Message:       message,
		Fix:           "Check your input and try again.",
		WayOut:        msgGoBack,
		ShowTimestamp: true,
	}
}

// Conflict returns a 409-style error page.
func Conflict(message string) ErrorPageProps {
	if message == "" {
		message = msgConflictDefault
	}
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyConflict,
		Code:          CodeConflict,
		Title:         "Conflict detected",
		Message:       message,
		Fix:           "Refresh your data and try the operation again.",
		WayOut:        msgGoBack,
		ShowTimestamp: true,
	}
}

// ServiceUnavailable returns a 503-style error page.
func ServiceUnavailable() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyTransient,
		Code:          CodeUnavailable,
		Title:         "Service temporarily unavailable",
		Message:       msgServiceUnavailable,
		Why:           "This is a temporary issue. No data was lost.",
		Fix:           "Wait a moment and refresh the page.",
		WayOut:        "Retry",
		WayOutHref:    "/",
		ShowTimestamp: true,
	}
}

// InternalError returns a 500-style error page.
func InternalError() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // pre-built with intentional defaults
		Family:        FamilyInfrastructure,
		Code:          CodeInternalError,
		Title:         titleInternalError,
		Message:       msgInternalUnexpected,
		Why:           "This is a system issue, not something you caused.",
		Fix:           "Try again in a few minutes. If the problem persists, contact support.",
		WayOut:        msgGoHome,
		WayOutHref:    "/",
		ShowTimestamp: true,
	}
}
