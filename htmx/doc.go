// Package htmx provides HTMX integration helpers including error handling patterns
// and response utilities for building dynamic web applications.
//
// # The tc-btn-loading indicator hook
//
// LoadingButton wraps its output in a `<span class="tc-btn-loading">`. That
// class is a deliberately public, stable hook for [hx-indicator] scoping: by
// default HTMX toggles `.htmx-indicator` visibility relative to the REQUESTING
// element, but when the request fires from a different element (e.g. a form
// submit driving a toolbar button), point hx-indicator at the wrapper:
//
//	<button hx-post="/save" hx-indicator="closest .tc-btn-loading">
//	   @htmx.LoadingButton("Save", "Saving…", spinner)
//	</button>
//
// Do not remove or rename the class — consumer CSS and hx-indicator selectors
// depend on it (guarded by htmx snapshot tests + the demo's
// visualtest/loading_button_e2e_test.go).
//
// [hx-indicator]: https://htmx.org/attributes/hx-indicator/
package htmx
