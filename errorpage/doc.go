// Package errorpage provides components for presenting structured errors on the web.
//
// Designed as a companion to the go-error-family library, this package renders
// error families (Rejection, Conflict, Transient, Corruption, Infrastructure)
// with family-appropriate visual styling — distinct colors, icons, and tone.
//
// The package integrates with go-error-family for type-safe error extraction.
// FromError() detects errorfamily.Classified errors and extracts family, code,
// context, cause chain, and default Why/Fix messages automatically.
//
// For errors from other sources, use the string-based bridge:
//
//	family := errorpage.ParseFamily(myError.ErrorFamily())
//
// Components:
//   - ErrorPage:  Full-page error view for HTTP error responses (4xx/5xx)
//   - ErrorDetail: Inline card with context table, cause chain, and fix
//   - ErrorAlert:  Alert banner with family-aware styling
//
// HTTP Handlers:
//   - ErrorHandler(err, cfg): returns http.Handler with correct status code
//   - WriteError(w, r, err, nonce): convenience wrapper
//   - WriteErrorPage(w, r, status, props, nonce): pre-configured page
//   - HTMLShell option: wraps in valid HTML document for standalone responses
//   - JSON option: renders JSON for API/HTMX endpoints
//
// Pre-built constructors: NotFound(), Forbidden(), BadRequest(msg),
// Conflict(msg), ServiceUnavailable(), InternalError()
//
// Each family maps to a distinct visual treatment:
//
//	Family          | Color   | Icon                | Tone
//	Rejection       | Amber   | ExclamationTriangle | Instructional
//	Conflict        | Orange  | ExclamationCircle   | Explanatory
//	Transient       | Blue    | Refresh             | Reassuring
//	Corruption      | Red     | ExclamationTriangle | Urgent
//	Infrastructure  | Slate   | Globe               | Apologetic
package errorpage
