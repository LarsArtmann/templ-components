package wire

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

// Headers that let one HTTP handler serve both transports. A datastar fetch
// marks itself with HeaderDatastarRequest; its non-SSE HTML response is
// patched into the region named by HeaderDatastarSelector using the merge
// mode in HeaderDatastarMode (verified against the pinned v1.0.2 runtime —
// see docs/datastar-runtime-facts.md). htmx marks its requests with
// HeaderHXRequest and targets client-side via hx-target.
const (
	HeaderDatastarRequest  = "Datastar-Request"
	HeaderDatastarSelector = "Datastar-Selector"
	HeaderDatastarMode     = "Datastar-Mode"
	HeaderHXRequest        = "Hx-Request"
)

// Transport selects the client-side runtime that executes an Action.
type Transport string

const (
	// TransportUnspecified is the zero value. It resolves to TransportHTMX —
	// the library default (ADR-0030) — so plain Action literals stay ergonomic.
	TransportUnspecified Transport = ""
	// TransportHTMX wires via hx-* attributes (htmx 2.x trigger engine).
	TransportHTMX Transport = "htmx"
	// TransportDatastar wires via data-on:* attributes (@get/@post expressions).
	TransportDatastar Transport = "datastar"
)

// TransportIsValid reports whether t is a defined transport (the zero value
// counts as defined: it means "use the library default").
func TransportIsValid(t Transport) bool {
	switch t {
	case TransportUnspecified, TransportHTMX, TransportDatastar:
		return true
	default:
		return false
	}
}

// Method is the HTTP method of the wired exchange.
type Method string

const (
	// MethodUnspecified is the zero value. It resolves to MethodGet, the
	// safe default in both dialects.
	MethodUnspecified Method = ""
	MethodGet         Method = "get"
	MethodPost        Method = "post"
	MethodPut         Method = "put"
	MethodPatch       Method = "patch"
	MethodDelete      Method = "delete"
)

// MethodIsValid reports whether m is a defined method (the zero value counts
// as defined: it means "use the GET default").
func MethodIsValid(m Method) bool {
	switch m {
	case MethodUnspecified, MethodGet, MethodPost, MethodPut, MethodPatch, MethodDelete:
		return true
	default:
		return false
	}
}

// Event is the DOM event that triggers the wired exchange.
type Event string

const (
	// EventUnspecified is the zero value. Under htmx it is omitted so the
	// trigger engine applies its element defaults (click on buttons, submit
	// on forms, change on inputs). Under Datastar it resolves to click,
	// because data-on: requires an explicit event key.
	EventUnspecified Event = ""
	EventClick       Event = "click"
	EventSubmit      Event = "submit"
	EventChange      Event = "change"
	EventInput       Event = "input"
	EventKeyDown     Event = "keydown"
	EventKeyUp       Event = "keyup"
	EventFocus       Event = "focus"
	EventBlur        Event = "blur"
)

// EventIsValid reports whether e is a defined event (the zero value counts
// as defined: it means "use the dialect default").
func EventIsValid(e Event) bool {
	switch e {
	case EventUnspecified, EventClick, EventSubmit, EventChange, EventInput,
		EventKeyDown, EventKeyUp, EventFocus, EventBlur:
		return true
	default:
		return false
	}
}

// ContentType selects how field values travel to the server. It mirrors the
// Datastar fetch action's contentType option (verified against the pinned
// v1.0.3 bundle — see docs/datastar-runtime-facts.md); htmx ignores it
// because hx-* requests serialize the enclosing form natively.
type ContentType string

const (
	// ContentTypeUnspecified is the zero value. It resolves to ContentTypeJSON
	// (the Datastar runtime default: all signals as a JSON body).
	ContentTypeUnspecified ContentType = ""
	// ContentTypeJSON sends the Datastar signals object as a JSON body — the
	// runtime default.
	ContentTypeJSON ContentType = "json"
	// ContentTypeForm serializes the enclosing <form>'s fields: HTML5
	// constraint validation gates the request, the submitter button's
	// name/value is appended, enctype="multipart/form-data" sends a FormData
	// body (file uploads), anything else sends application/x-www-form-
	// urlencoded, and GET requests carry the fields as query parameters.
	// This is the Datastar twin of htmx's native form serialization — the
	// piece that makes whole-form submission transport-symmetric.
	ContentTypeForm ContentType = "form"
)

// ContentTypeIsValid reports whether c is a defined content type (the zero
// value counts as defined: it means "use the runtime default").
func ContentTypeIsValid(c ContentType) bool {
	switch c {
	case ContentTypeUnspecified, ContentTypeJSON, ContentTypeForm:
		return true
	default:
		return false
	}
}

// Action describes one client-initiated hypermedia exchange in transport
// dialects' common subset: a method on a URL, triggered by a DOM event,
// patching a target region.
type Action struct {
	// Transport selects the attribute dialect. Zero value renders as htmx.
	Transport Transport
	// Method is the HTTP verb. Zero value renders as GET in BOTH dialects
	// (hx-get / @get): an Action literal that omits Method silently turns a
	// write into a read. Set Method explicitly on every mutating action
	// (usually MethodPost); rely on the default only for reads (filters,
	// searches, fragment loads).
	Method Method
	// URL is the exchange endpoint. An empty URL wires nothing (Attributes
	// returns nil) — a component rendered without a backend endpoint must
	// stay inert, not emit a broken binding.
	URL string
	// Event is the triggering DOM event. Zero value uses the dialect default.
	Event Event
	// Target is the region to patch. Under htmx it renders as hx-target.
	// Under Datastar it is intentionally NOT rendered: Datastar v1.0.2 fetch
	// actions accept no target option — targeting is response-driven (see
	// HeaderDatastarSelector) or id-matched (a fragment root whose id equals
	// the target id patches it in the default outer mode). Handlers honor it
	// by echoing the selector back on HeaderDatastarSelector.
	Target string
	// ContentType selects how field values travel to the server. Under
	// Datastar, ContentTypeForm renders {contentType: 'form'} so the action
	// serializes the enclosing form's fields (validation gate, submitter
	// name/value, enctype-aware body); the zero value and ContentTypeJSON use
	// the runtime default (signals as JSON). htmx ignores this field — hx-*
	// requests serialize the enclosing form natively either way.
	ContentType ContentType
	// Selector is the Datastar-dialect patch target: rendered as the fetch
	// option {selector: '<sel>'}, which patches the response into the
	// element(s) matching the selector client-side — overriding the
	// response-header targeting (Datastar-Selector) when both are present
	// (option wins, verified against the pinned v1.0.3 bundle). Under
	// contentType form it also selects which form serializes, falling back
	// to the action element's closest form. It renders for Datastar only —
	// the htmx twin is Target. Empty renders nothing (response-driven).
	Selector string
	// Swap is the region-merge style for the exchange — how the response is
	// inserted relative to the target. It reuses the server-side PatchMode
	// vocabulary (the SAME enum as PatchTarget.Mode), so one mode word
	// describes both the client request and the server targeting. Under htmx
	// it renders hx-swap (innerHTML/outerHTML/afterbegin/…); under Datastar
	// the fetch option {mode: '…'}, which overrides the Datastar-Mode
	// response header (verified against the pinned bundle). The zero value
	// renders nothing in either dialect: htmx's default swap and the
	// wire.Handler default mode are both "inner", so an unspecified Swap
	// behaves identically under both runtimes.
	//
	// Six modes map 1:1 (inner, outer, before, append, prepend, after, and
	// remove ↔ htmx delete); PatchModeReplace (Datastar-only) degrades to the
	// closest htmx style (outerHTML). Unknown values render nothing.
	Swap PatchMode
	// DebounceMS delays the wired exchange until the event has stopped
	// firing for this many milliseconds — the auto-submit filter-input
	// pattern. htmx renders it as the delay:<n>ms trigger modifier (plus
	// `changed` on value events — input, change, keyup — so an unchanged
	// value never re-requests); Datastar renders it as the __debounce.<n>ms
	// event modifier (spelling decoded from the pinned v1.0.3 bundle — see
	// docs/datastar-runtime-facts.md). Zero (the default) emits no debounce.
	// Under htmx it requires an explicit Event (the zero event renders no
	// hx-trigger at all, so the delay would be silently dropped).
	DebounceMS int
	// PreventDefault suppresses the triggering event's native default action
	// in the Datastar dialect (the __prevent event modifier). The runtime
	// auto-preventDefaults ONLY form+submit — decoded from the pinned v1.0.3
	// bundle — so a wired <a href> would otherwise BOTH patch and navigate,
	// and a wired submit button would double-submit. Set it on actions wired
	// onto anchors whose href is the no-JS fallback (NavLink, Calendar
	// MonthNav): it is the semantic twin of htmx's automatic interception of
	// hx-* clicks, which is why htmx ignores this field. Do NOT set it on
	// checkbox/radio/change wirings: preventDefault there would cancel the
	// native toggle/commit that htmx keeps.
	PreventDefault bool
}

// Attributes renders the action as templ attributes in the transport's
// dialect. It returns nil when the URL is empty (no wiring) so the result
// can always be spread into templ attribute position.
func (a Action) Attributes() templ.Attributes {
	if a.URL == "" {
		return nil
	}

	if a.transport() == TransportDatastar {
		return a.datastarAttributes()
	}

	return a.htmxAttributes()
}

// transport resolves the zero value to the library default.
func (a Action) transport() Transport {
	if a.Transport == TransportDatastar {
		return TransportDatastar
	}

	return TransportHTMX
}

// htmxAttributes renders hx-get/hx-post/..., optional hx-trigger, hx-target.
func (a Action) htmxAttributes() templ.Attributes {
	method := htmxMethod(a.Method)

	attrs := templ.Attributes{
		"hx-" + method: a.URL,
	}

	if a.Event != EventUnspecified && EventIsValid(a.Event) {
		trigger := string(a.Event)
		if a.DebounceMS > 0 {
			if a.Event == EventInput || a.Event == EventChange || a.Event == EventKeyUp {
				trigger += " changed"
			}

			trigger += fmt.Sprintf(" delay:%dms", a.DebounceMS)
		}

		attrs["hx-trigger"] = trigger
	}

	if a.Target != "" {
		attrs["hx-target"] = a.Target
	}

	if style := htmxSwapStyle(a.Swap); style != "" {
		attrs["hx-swap"] = style
	}

	return attrs
}

// datastarAttributes renders data-on:<event>="@<method>('<url>', <opts>)".
func (a Action) datastarAttributes() templ.Attributes {
	event := string(a.Event)
	if !EventIsValid(a.Event) || a.Event == EventUnspecified {
		event = string(EventClick)
	}

	key := "data-on:" + event
	if a.PreventDefault {
		key += "__prevent"
	}

	if a.DebounceMS > 0 {
		key += fmt.Sprintf("__debounce.%dms", a.DebounceMS)
	}

	return templ.Attributes{
		key: datastarActionExpr(a.method(), a.URL, a.ContentType, a.Selector, a.Swap),
	}
}

// method resolves the zero value (and, defensively, unknown values) to GET.
func (a Action) method() Method {
	if MethodIsValid(a.Method) && a.Method != MethodUnspecified {
		return a.Method
	}

	return MethodGet
}

// htmxMethod maps to the hx-<method> suffix, falling back to get.
func htmxMethod(m Method) string {
	if MethodIsValid(m) && m != MethodUnspecified {
		return string(m)
	}

	return string(MethodGet)
}

// datastarActionExpr builds a @<method>('<url>') expression, appending a single
// options object ({selector, mode, contentType}) when the action selects any of
// them — the only non-default options the common subset expresses (JSON is the
// runtime default and is omitted; unknown values degrade to the default).
//
// All options travel in ONE object literal, never several: the pinned bundle's
// action dispatcher invokes apply(ctx, url, opts) with spread arguments, so a
// SECOND object argument is silently ignored (the runtime reads exactly one
// options object). Emitting separate objects — the earlier shape — dropped form
// encoding whenever a Selector was also set.
//
// Single quotes are escaped so a URL, selector, or mode cannot inject into the
// expression (mirrors the datastar package's actionExpr).
func datastarActionExpr(
	method Method,
	url string,
	contentType ContentType,
	selector string,
	swap PatchMode,
) string {
	var opts []string

	if selector != "" {
		opts = append(opts, fmt.Sprintf("selector: '%s'", escapeSingleQuotes(selector)))
	}

	if mode := datastarSwapMode(swap); mode != "" {
		opts = append(opts, fmt.Sprintf("mode: '%s'", mode))
	}

	if contentType == ContentTypeForm {
		opts = append(opts, "contentType: 'form'")
	}

	base := fmt.Sprintf("@%s('%s')", method, escapeSingleQuotes(url))
	if len(opts) == 0 {
		return base
	}

	return fmt.Sprintf("@%s('%s', {%s})", method, escapeSingleQuotes(url), strings.Join(opts, ", "))
}

// escapeSingleQuotes makes a value safe inside a single-quoted JS string
// literal embedded in a data-on expression.
func escapeSingleQuotes(value string) string {
	return strings.ReplaceAll(value, `'`, `\'`)
}

// datastarSwapMode resolves Swap to a Datastar fetch mode, or "" when unset or
// unknown (render nothing, use the runtime/response-header default).
func datastarSwapMode(swap PatchMode) string {
	if !PatchModeIsValid(swap) || swap == PatchModeUnspecified {
		return ""
	}

	return string(swap)
}

// htmxSwapStyles maps the shared PatchMode vocabulary onto htmx hx-swap styles.
// PatchModeReplace is Datastar-only (a morphing replaceWith) and degrades to
// htmx's closest style, outerHTML.
//
//nolint:gochecknoglobals // immutable lookup table
var htmxSwapStyles = map[PatchMode]string{
	PatchModeInner:   "innerHTML",
	PatchModeOuter:   "outerHTML",
	PatchModeBefore:  "beforebegin",
	PatchModePrepend: "afterbegin",
	PatchModeAppend:  "beforeend",
	PatchModeAfter:   "afterend",
	PatchModeRemove:  "delete",
	PatchModeReplace: "outerHTML",
}

// htmxSwapStyle resolves Swap to an htmx hx-swap style, or "" when unset or
// unknown (render no hx-swap, use htmx's innerHTML default).
func htmxSwapStyle(swap PatchMode) string {
	if !PatchModeIsValid(swap) || swap == PatchModeUnspecified {
		return ""
	}

	return htmxSwapStyles[swap]
}
