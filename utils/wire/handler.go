package wire

import "net/http"

// PatchMode is the merge mode the Datastar runtime applies when patching a
// non-SSE HTML response into the region named by HeaderDatastarSelector. The
// value set is verified against the pinned v1.0.2 runtime bundle — see
// docs/datastar-runtime-facts.md.
type PatchMode string

const (
	// PatchModeUnspecified is the zero value. Handler resolves it to
	// PatchModeInner — the least destructive mode that still replaces the
	// region's content.
	PatchModeUnspecified PatchMode = ""
	PatchModeInner       PatchMode = "inner"
	PatchModeOuter       PatchMode = "outer"
	PatchModePrepend     PatchMode = "prepend"
	PatchModeAppend      PatchMode = "append"
	PatchModeBefore      PatchMode = "before"
	PatchModeAfter       PatchMode = "after"
	PatchModeReplace     PatchMode = "replace"
)

// PatchModeIsValid reports whether m is a merge mode of the pinned Datastar
// runtime (the zero value counts as defined: it means "use the inner
// default").
func PatchModeIsValid(m PatchMode) bool {
	switch m {
	case PatchModeUnspecified, PatchModeInner, PatchModeOuter, PatchModePrepend,
		PatchModeAppend, PatchModeBefore, PatchModeAfter, PatchModeReplace:
		return true
	default:
		return false
	}
}

// mode resolves the zero value to the Handler default.
func (m PatchMode) mode() PatchMode {
	if PatchModeIsValid(m) && m != PatchModeUnspecified {
		return m
	}

	return PatchModeInner
}

// PatchTarget describes where a Datastar caller patches a fragment response.
// htmx callers do not need it: hx-target already named the region client-side.
type PatchTarget struct {
	// Selector is the CSS selector of the patch region, e.g. "#out".
	Selector string
	// Mode is the merge mode. Zero value renders as inner.
	Mode PatchMode
}

// IsDatastar reports whether the request was issued by the Datastar runtime
// (any fetch action marks itself with HeaderDatastarRequest).
func IsDatastar(r *http.Request) bool {
	return r.Header.Get(HeaderDatastarRequest) != ""
}

// IsHTMX reports whether the request was issued by htmx (every htmx AJAX
// request carries HeaderHXRequest).
func IsHTMX(r *http.Request) bool {
	return r.Header.Get(HeaderHXRequest) != ""
}

// Handler wraps next so one endpoint serves every caller of a wired fragment:
//
//   - Datastar callers receive HeaderDatastarSelector/HeaderDatastarMode
//     response headers targeting target — Datastar fetch actions have no
//     client-side target option, so the response owns the targeting.
//   - htmx and plain callers pass through untouched: htmx targets client-side
//     via hx-target, and a plain navigation just renders the fragment.
//
// An empty target.Selector degrades gracefully: no routing headers are set,
// so the Datastar runtime falls back to its default id-matched patching (a
// fragment root whose id equals the caller's element id).
//
//	mux.Handle("/api/fragment", wire.Handler(wire.PatchTarget{
//		Selector: "#out",
//	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "text/html; charset=utf-8")
//		componentOr500(w, r, fragment())
//	})))
func Handler(target PatchTarget, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if target.Selector != "" && IsDatastar(r) {
			w.Header().Set(HeaderDatastarSelector, target.Selector)
			w.Header().Set(HeaderDatastarMode, string(target.Mode.mode()))
		}

		next.ServeHTTP(w, r)
	})
}
