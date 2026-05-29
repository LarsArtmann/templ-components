// Package errorpage provides components for presenting structured errors on the web.
//
// Designed as a companion to the go-error-family library, this package renders
// error families (Rejection, Conflict, Transient, Corruption, Infrastructure)
// with family-appropriate visual styling — distinct colors, icons, and tone.
//
// The package has zero dependency on go-error-family. Consumers bridge the gap
// with trivial string constants:
//
//	// Convert go-error-family → errorpage
//	family := errorpage.Family(myError.ErrorFamily().String())
//
// Components:
//   - ErrorPage:  Full-page error view for HTTP error responses (4xx/5xx)
//   - ErrorDetail: Inline card with context table, cause chain, and fix
//   - ErrorAlert:  Alert banner with family-aware styling
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
