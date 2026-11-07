package datastar

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/a-h/templ"
)

// Typed builders for the Datastar runtime's reactive plugins — the CLIENT
// half of the attribute surface (the request half is utils/wire, ADR-0036).
// Every builder returns templ.Attributes: spread them into any element
// (`<span ${ datastar.Text("$user.name") }>`) and templ renders the value
// escaped. Empty names/expressions degrade to nil — a builder without its
// inputs renders nothing instead of a broken binding (the graceful-
// degradation rule every library component follows).
//
// Scope note (JS decision ladder, docs/javascript-guide.md): these builders
// exist for consumers' OWN glue — dashboard widgets, exploratory UIs,
// reactive islands. The library's components stay native-first/htmx-first;
// no component is re-ported to signals.

// signalSegmentGrammar is the attribute-suffix-safe signal grammar: lowercase
// segments (the HTML parser lowercases attribute NAMES, so an uppercase
// segment would silently bind a different signal than the expression reads —
// see runtime-facts §CSP caveats), digits and underscores allowed, dotted
// paths for nested signals. A leading underscore marks a LOCAL signal (never
// sent to the backend: the runtime's default filterSignals exclude
// /(^|\.)_/). Double underscores are banned — they are the MODIFIER delimiter
// (data-on:click__debounce…), so a signal carrying them can never round-trip
// through an attribute suffix.
var signalSegmentGrammar = regexp.MustCompile(`^[a-z0-9_]+$`)

// SignalName is a validated Datastar signal name for attribute-suffix
// bindings (data-bind:<name>, data-computed:<name>). Not a closed enum — an
// open validated type: there is no fixed value set, so there is no
// SignalNameDefault constant; IsValid is the contract, its test lives in
// attributes_test.go (same commit rule).
type SignalName string

// IsValid reports whether the name is safe to embed in an attribute suffix:
// one or more lowercase word segments separated by dots, no "__" anywhere
// (the modifier delimiter), no empty segments.
func (s SignalName) IsValid() bool {
	name := string(s)
	if name == "" || strings.Contains(name, "__") {
		return false
	}

	for segment := range strings.SplitSeq(name, ".") {
		if !signalSegmentGrammar.MatchString(segment) {
			return false
		}
	}

	return true
}

// eventGrammar constrains handler events to real DOM/custom-event names:
// letters, digits, dashes and dots (e.g. click, keydown, datastar-fetch).
// Anything else cannot be an event and would corrupt the attribute NAME —
// rejected (renders nothing) rather than escaped, unlike expression VALUES
// which escape safely.
var eventGrammar = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]*$`)

// normalizeExpr makes a consumer expression safe AS AN ATTRIBUTE VALUE: CR/LF
// collapse to spaces (one statement line). Quotes and angle brackets pass
// through VERBATELY — they are expression CODE (string literals like
// @get('/x'), comparisons) that the runtime parses, and the HTML attribute
// layer is escaped by templ at render time. The builders never interpolate
// user input INTO an expression; they only carry one.
func normalizeExpr(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")

	return s
}

// attrsOrNil collapses an empty attribute map to nil so `…${ nil }…` renders
// nothing.
func attrsOrNil(attrs templ.Attributes) templ.Attributes {
	if len(attrs) == 0 {
		return nil
	}

	return attrs
}

// Show binds an element's visibility to an expression (data-show): the
// element renders when the expression is truthy.
//
//	<span ${ datastar.Show("$user.admin") }>secrets</span>
func Show(expr string) templ.Attributes {
	if expr == "" {
		return nil
	}

	return templ.Attributes{"data-show": normalizeExpr(expr)}
}

// Text binds an element's text content to an expression (data-text): the
// node re-renders whenever any referenced signal changes.
//
//	<dd ${ datastar.Text("$metrics.users") }>0</dd>
func Text(expr string) templ.Attributes {
	if expr == "" {
		return nil
	}

	return templ.Attributes{"data-text": normalizeExpr(expr)}
}

// Bind two-ways an input element to a signal (data-bind:<name>): element
// edits write the signal, signal patches write the element.
//
//	<input ${ datastar.Bind("email") }>
func Bind(name SignalName) templ.Attributes {
	if !name.IsValid() {
		return nil
	}

	return templ.Attributes{"data-bind:" + string(name): ""}
}

// Computed declares a derived signal (data-computed:<name>): the expression
// re-evaluates when its signal dependencies change.
//
//	<div ${ datastar.Computed("total", "$price * $qty") }>
func Computed(name SignalName, expr string) templ.Attributes {
	if !name.IsValid() || expr == "" {
		return nil
	}

	return templ.Attributes{"data-computed:" + string(name): normalizeExpr(expr)}
}

// Signals declares initial signal values (data-signals): the JSON object
// deep-merges into the signal tree at init. Keys are VALUE-borne (JSON), so
// camelCase keys are safe here — the attribute-name-lowercasing gotcha does
// not apply. A payload that cannot marshal to JSON degrades to nil.
func Signals(values map[string]any) templ.Attributes {
	if len(values) == 0 {
		return nil
	}

	buf, err := json.Marshal(values)
	if err != nil {
		return nil
	}

	return templ.Attributes{"data-signals": string(buf)}
}

// On wires an event to an expression (data-on:<event>) with the runtime's
// modifier groups spelled the pinned-bundle way (data-on:click__debounce.300ms,
// data-on:submit__prevent — groups only parse after __, never as a dotted
// suffix). Chain modifiers with value receivers; the zero-chained form is the
// plain handler. An empty event or expression renders nothing.
//
//	<button ${ datastar.On("click", "$count++").Debounce(300) }>go</button>
type OnExpr struct {
	event      string
	expr       string
	debounceMS int
	throttleMS int
	prevent    bool
}

// On starts an event-handler builder.
func On(event, expr string) OnExpr {
	return OnExpr{event: event, expr: expr}
}

// Debounce adds the __debounce.<n>ms modifier (trailing-edge; the runtime's
// leading/notrailing sub-modifiers are not modeled — use raw attrs for those).
func (o OnExpr) Debounce(ms int) OnExpr {
	o.debounceMS = ms

	return o
}

// Throttle adds the __throttle.<n>ms modifier (leading-edge default).
func (o OnExpr) Throttle(ms int) OnExpr {
	o.throttleMS = ms

	return o
}

// Prevent adds the __prevent modifier (suppresses the event's native
// default — required for wired anchors whose href is the no-JS fallback).
func (o OnExpr) Prevent() OnExpr {
	o.prevent = true

	return o
}

// Attrs renders the builder. Modifiers require the explicit event (an orphan
// modifier is dropped by design — see wire.Action's documented rule). An
// event that is not a plausible DOM/custom-event name renders nothing.
func (o OnExpr) Attrs() templ.Attributes {
	if o.event == "" || o.expr == "" || !eventGrammar.MatchString(o.event) {
		return nil
	}

	key := "data-on:" + o.event

	if o.prevent {
		key += "__prevent"
	}

	if o.debounceMS > 0 {
		key += fmt.Sprintf("__debounce.%dms", o.debounceMS)
	}

	if o.throttleMS > 0 {
		key += fmt.Sprintf("__throttle.%dms", o.throttleMS)
	}

	return attrsOrNil(templ.Attributes{key: normalizeExpr(o.expr)})
}
