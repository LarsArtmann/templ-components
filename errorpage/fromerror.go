package errorpage

import (
	"errors"
	"time"

	errorfamily "github.com/larsartmann/go-error-family"
)

// FamilyFromErrorFamily converts a go-error-family Family to an errorpage Family.
// Returns FamilyTransient for unrecognized values.
func FamilyFromErrorFamily(f errorfamily.Family) Family {
	return ParseFamily(f.String())
}

// FromError converts any error into ErrorPageProps.
// Extracts code, family, context, and cause chain from structured errors.
// For go-error-family errors, also extracts Why/Fix defaults.
// Falls back to Infrastructure family for unknown errors.
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

	// Reuse the classified error if familyFromError already resolved it,
	// avoiding a second errors.AsType call.
	if classified, ok := errors.AsType[errorfamily.Classified](err); ok {
		ef := classified.ErrorFamily()
		props.Why = ef.DefaultWhy()
		props.Fix = ef.DefaultFix()
	}

	if titled, ok := err.(interface{ ErrorTitle() string }); ok {
		title := titled.ErrorTitle()
		if title != "" {
			props.Title = title
		}
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
	return FamilyInfrastructure
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
	Family  string            `json:"family"`
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message"`
	Title   string            `json:"title,omitempty"`
	Why     string            `json:"why,omitempty"`
	Fix     string            `json:"fix,omitempty"`
	Context map[string]string `json:"context,omitempty"`
}
