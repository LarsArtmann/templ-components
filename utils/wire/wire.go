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
	// Method is the HTTP verb. Zero value renders as GET.
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

	return attrs
}

// datastarAttributes renders data-on:<event>="@<method>('<url>', <opts>)".
func (a Action) datastarAttributes() templ.Attributes {
	event := string(a.Event)
	if !EventIsValid(a.Event) || a.Event == EventUnspecified {
		event = string(EventClick)
	}

	key := "data-on:" + event
	if a.DebounceMS > 0 {
		key += fmt.Sprintf("__debounce.%dms", a.DebounceMS)
	}

	return templ.Attributes{
		key: datastarActionExpr(a.method(), a.URL, a.ContentType, a.Selector),
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

// datastarActionExpr builds a @<method>('<url>') expression, appending
// {contentType: 'form'} / {selector: '…'} when the action selects them —
// the only non-default options the common subset expresses (JSON is the
// runtime default and is omitted; unknown values degrade to the default).
// Single quotes are escaped so a URL or selector cannot inject into the
// expression (mirrors the datastar package's actionExpr).
func datastarActionExpr(method Method, url string, contentType ContentType, selector string) string {
	escaped := strings.ReplaceAll(url, `'`, `\'`)

	var opts []string
	if selector != "" {
		opts = append(opts, fmt.Sprintf("{selector: '%s'}", strings.ReplaceAll(selector, `'`, `\'`)))
	}
	if contentType == ContentTypeForm {
		opts = append(opts, "{contentType: 'form'}")
	}

	if len(opts) == 0 {
		return fmt.Sprintf("@%s('%s')", method, escaped)
	}

	return fmt.Sprintf("@%s('%s', %s)", method, escaped, strings.Join(opts, ", "))
}
