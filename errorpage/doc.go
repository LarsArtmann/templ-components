// Package errorpage provides components for presenting structured errors on the web.
//
// Designed as a companion to the go-error-family library, this package renders
// error families (Rejection, Conflict, Transient, Corruption, Infrastructure)
// with family-appropriate visual styling — distinct colors, icons, and tone.
//
// The package integrates with go-error-family for type-safe error extraction.
// FromError() detects errorfamily.Classified errors and extracts family, code,
// context, cause chain, and default Why/Fix messages automatically. It also
// prefers a user-safe message (Public(), the oops convention promoted through
// go-error-family/bridge) over Message()/Error(), and derives a short
// family-appropriate title (FamilyDefaultTitle) when the error carries no
// ErrorTitle — so FromError-driven pages always render a heading.
//
// For errors from other sources, use the string-based bridge:
//
//	family := errorpage.ParseFamily(myError.ErrorFamily())
//
// Components:
//   - ErrorPage:   Full-page error view for HTTP error responses (4xx/5xx)
//   - NotFound404: Dedicated 404 page with gradient numeral, search, quick-links
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
//	Family          | Color  | Icon                | Tone
//	Rejection       | Amber  | ExclamationTriangle | Instructional
//	Conflict        | Orange | ExclamationCircle   | Explanatory
//	Transient       | Blue   | Refresh             | Reassuring
//	Corruption      | Red    | ExclamationTriangle | Urgent
//	Infrastructure  | Gray   | Globe               | Apologetic
//	Orchestration   | Purple | ExclamationTriangle | Factual
//
// Props reference (shared vocabulary across all three props types):
//
//	Field        | ErrorPage | ErrorDetail | ErrorAlert | Meaning
//	Family       | ✓         | ✓           | ✓          | Visual treatment + derived defaults; must be a valid Family constant
//	StatusCode   | ✓         |             |            | HTTP status rendered in the chip row ("HTTP 503"); must be in [400, 599]
//	Code         | ✓         | ✓           |            | Machine-readable error code (open enum, e.g. "service.unavailable")
//	Title        | ✓         | ✓           | ✓          | Page/card heading; FromError fills a family default when absent
//	Message      | ✓         | ✓           | ✓          | What happened, user-safe (never log details here)
//	Why          | ✓         |             |            | Reassurance/context ("this is temporary, no data lost")
//	Fix          | ✓         | ✓           | ✓          | The suggested fix (renders in its own panel)
//	WayOut       | ✓         |             |            | Primary action label ("Retry"); requires WayOutHref
//	WayOutHref   | ✓         |             |            | Primary action destination
//	Context      | ✓         | ✓           |            | Key/value diagnostics (region, request_id, ...)
//	CauseChain   | ✓         | ✓           |            | Ordered causes, deepest first (message + optional code)
//	Timestamp    | ✓         | ✓           |            | RFC3339; FromError uses the error's own when available
//	Trace        | ✓         | ✓           |            | Correlation ID (oops.Trace(), promoted via the bridge); footer/JSON field
//	ShowTimestamp| ✓         |             |            | Render the timestamp footer at all
//
// All three embed utils.BaseProps (Class, Attrs, ID, AriaLabel, Nonce) — set
// them via the embedded struct, not by promoted field name.
//
// Code is an OPEN enum: unlike closed-set enums in this library it ships no
// IsValid, by design. Codes are consumer-namespaced identifiers
// ("service.unavailable", "db.pool") — the set is unbounded, every value is
// renderable, and there is no lookup table to miss. Validation would only
// reject legitimate consumer codes.
package errorpage
