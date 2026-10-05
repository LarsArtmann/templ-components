package wire

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// Invariants that must hold for EVERY Action, regardless of field values.
// They pin the contract against silent rot: the same URL appears in both
// dialects, an empty URL stays inert for every field combination, and wire
// rendering never emits a script (the contract is attributes-only by design —
// CSP-safe without a nonce).

func TestEmptyURLIsInertForEveryEnumCombo(t *testing.T) {
	t.Parallel()

	transports := []Transport{
		TransportUnspecified,
		TransportHTMX,
		TransportDatastar,
		Transport("bogus"),
	}
	methods := []Method{
		MethodUnspecified,
		MethodGet,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		Method("bogus"),
	}
	events := []Event{
		EventUnspecified, EventClick, EventSubmit, EventChange, EventInput,
		EventKeyDown, EventKeyUp, EventFocus, EventBlur, Event("hover"),
	}
	targets := []string{"", "#out", "closest div"}
	contentTypes := []ContentType{
		ContentTypeUnspecified,
		ContentTypeJSON,
		ContentTypeForm,
		ContentType("multipart"),
	}

	for _, transport := range transports {
		for _, method := range methods {
			for _, event := range events {
				for _, target := range targets {
					for _, contentType := range contentTypes {
						action := Action{
							Transport:   transport,
							Method:      method,
							URL:         "",
							Event:       event,
							Target:      target,
							ContentType: contentType,
						}

						if attrs := action.Attributes(); attrs != nil {
							t.Fatalf("empty URL must wire nothing, got %v for %+v", attrs, action)
						}
					}
				}
			}
		}
	}
}

func TestURLReferencedInBothDialects(t *testing.T) {
	t.Parallel()

	transports := []Transport{TransportUnspecified, TransportHTMX, TransportDatastar}
	methods := []Method{
		MethodUnspecified,
		MethodGet,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
	}
	events := []Event{EventUnspecified, EventClick, EventSubmit, EventChange, EventInput}
	targets := []string{"", "#out"}
	contentTypes := []ContentType{ContentTypeUnspecified, ContentTypeJSON, ContentTypeForm}

	const url = "/api/items?filter=x"

	for _, transport := range transports {
		for _, method := range methods {
			for _, event := range events {
				for _, target := range targets {
					for _, contentType := range contentTypes {
						action := Action{
							Transport:   transport,
							Method:      method,
							URL:         url,
							Event:       event,
							Target:      target,
							ContentType: contentType,
						}

						assertURLReferenced(t, action, url)
					}
				}
			}
		}
	}
}

// assertURLReferenced fails when no rendered attribute value carries the
// action's URL.
func assertURLReferenced(t *testing.T, action Action, url string) {
	t.Helper()

	attrs := action.Attributes()
	if attrs == nil {
		t.Fatalf("URL %q vanished for %+v", url, action)
	}

	found := false

	for _, attrValue := range attrs {
		if value, ok := attrValue.(string); ok && strings.Contains(value, url) {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("no attribute references URL %q for %+v: %v", url, action, attrs)
	}
}

func TestSelectorNeverRenderedForHTMX(t *testing.T) {
	t.Parallel()

	for _, event := range []Event{EventUnspecified, EventClick, EventSubmit, EventChange, EventInput, EventKeyDown, EventKeyUp, EventFocus, EventBlur} {
		action := Action{
			Transport: TransportHTMX,
			Method:    MethodGet,
			URL:       "/api/items",
			Event:     event,
			Selector:  "#out",
		}

		for key, value := range action.Attributes() {
			if vs, isStr := value.(string); isStr && strings.Contains(vs, "#out") {
				t.Fatalf("htmx dialect must never render the datastar selector option, got %q=%q", key, value)
			}

			if strings.Contains(strings.ToLower(key), "selector") {
				t.Fatalf("htmx dialect must never render a selector attribute, got %q", key)
			}
		}
	}
}

func TestTargetNeverRenderedForDatastar(t *testing.T) {
	t.Parallel()

	for _, event := range []Event{EventUnspecified, EventClick, EventSubmit, EventChange, EventInput, EventKeyDown, EventKeyUp, EventFocus, EventBlur} {
		action := Action{
			Transport: TransportDatastar,
			Method:    MethodGet,
			URL:       "/api/items",
			Event:     event,
			Target:    "#out",
		}

		for key := range action.Attributes() {
			if strings.Contains(strings.ToLower(key), "target") {
				t.Fatalf("datastar dialect must never render a target attribute, got %q", key)
			}
		}
	}
}

// TestContentTypeHTMXInert pins that the content-type option never leaks
// into the htmx dialect: hx-* requests serialize forms natively, so the
// option has no hx-* spelling. Any htmx attribute value mentioning
// "contentType" would be a dialect leak.
func TestContentTypeHTMXInert(t *testing.T) {
	t.Parallel()

	for _, contentType := range []ContentType{ContentTypeUnspecified, ContentTypeJSON, ContentTypeForm, ContentType("bogus")} {
		for _, method := range []Method{MethodUnspecified, MethodGet, MethodPost, MethodPut, MethodPatch, MethodDelete} {
			action := Action{
				Transport:   TransportHTMX,
				Method:      method,
				URL:         "/api/items",
				ContentType: contentType,
			}

			for key, value := range action.Attributes() {
				if strings.Contains(key, "contentType") || strings.Contains(fmt.Sprint(value), "contentType") {
					t.Fatalf(
						"htmx dialect must not render the content-type option, got %q=%v for %+v",
						key, value, action,
					)
				}
			}
		}
	}
}

// TestSwapDialectIsolation pins the transport-asymmetric Swap semantics
// (corrected 2026-10-02 against the pinned bundle): the htmx dialect renders
// hx-swap client-side and never a datastar option; the datastar dialect
// renders NO mode option — the patch mode is response-header driven
// (wire.Handler/PatchTarget) — and never an hx-* attribute.
func TestSwapDialectIsolation(t *testing.T) {
	t.Parallel()

	swaps := []PatchMode{
		PatchModeInner, PatchModeOuter, PatchModeBefore, PatchModePrepend,
		PatchModeAppend, PatchModeAfter, PatchModeRemove, PatchModeReplace,
	}

	t.Run("htmx renders hx-swap only", func(t *testing.T) {
		t.Parallel()

		for _, swap := range swaps {
			attrs := Action{
				Transport: TransportHTMX, Method: MethodGet, URL: "/api/items", Swap: swap,
			}.Attributes()

			if _, ok := attrs["hx-swap"]; !ok {
				t.Fatalf("htmx dialect must render hx-swap for %q, got %v", swap, attrs)
			}

			for key, value := range attrs {
				if strings.Contains(strings.ToLower(key), "data-on") ||
					strings.Contains(fmt.Sprint(value), "mode:") {
					t.Fatalf("htmx dialect must not render the datastar mode option, got %q=%v", key, value)
				}
			}
		}
	})

	t.Run("datastar renders no swap at all (response-header driven)", func(t *testing.T) {
		t.Parallel()

		for _, swap := range swaps {
			attrs := Action{
				Transport: TransportDatastar, Method: MethodGet, URL: "/api/items", Swap: swap,
			}.Attributes()

			if datastarAttrsContain(attrs, "mode:") {
				t.Fatalf(
					"datastar dialect must not render a mode option (the runtime ignores it) for %q, got %v",
					swap,
					attrs,
				)
			}

			for key, value := range attrs {
				if strings.HasPrefix(key, "hx-") {
					t.Fatalf("datastar dialect must not render an hx-* attribute, got %q=%v", key, value)
				}
			}
		}
	})
}

// datastarAttrsContain reports whether any attribute value contains needle.
func datastarAttrsContain(attrs templ.Attributes, needle string) bool {
	for _, value := range attrs {
		if strings.Contains(fmt.Sprint(value), needle) {
			return true
		}
	}

	return false
}

// TestUnspecifiedSwapRendersNothing pins zero-value parity: an unspecified
// Swap renders no swap attribute in either dialect (htmx defaults to
// innerHTML, wire.Handler defaults to inner — both "inner").
func TestUnspecifiedSwapRendersNothing(t *testing.T) {
	t.Parallel()

	for _, transport := range []Transport{TransportHTMX, TransportDatastar} {
		action := Action{Transport: transport, Method: MethodGet, URL: "/api/items"}

		for key := range action.Attributes() {
			if key == "hx-swap" {
				t.Fatalf("unspecified Swap must not render hx-swap for %q", transport)
			}
		}

		if transport == TransportDatastar {
			if ds := action.Attributes()["data-on:click"]; ds != "@get('/api/items')" {
				t.Fatalf("unspecified Swap must render a plain expression, got %v", ds)
			}
		}
	}
}

func TestWireRenderingEmitsNoScript(t *testing.T) {
	t.Parallel()

	actions := []Action{
		{URL: "/api/items"},
		{Method: MethodPost, URL: "/api/items"},
		{Transport: TransportDatastar, URL: "/api/items"},
		{
			Transport: TransportDatastar,
			Method:    MethodPost,
			URL:       "/api/items?filter=it's",
			Event:     EventInput,
			Target:    "#out",
		},
	}

	for _, action := range actions {
		rendered := renderAttributes(t, action.Attributes())

		if strings.Contains(strings.ToLower(rendered), "<script") {
			t.Fatalf("wire rendering must never emit a script, got %q", rendered)
		}
	}
}

func BenchmarkActionAttributes(b *testing.B) {
	htmx := Action{Method: MethodPost, URL: "/api/items", Event: EventClick, Target: "#out"}
	datastar := Action{
		Transport: TransportDatastar,
		Method:    MethodPost,
		URL:       "/api/items",
		Event:     EventClick,
	}

	b.Run("htmx full", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = htmx.Attributes()
		}
	})

	b.Run("datastar full", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = datastar.Attributes()
		}
	})
}

// BenchmarkFormWireExpression pins the cost of the form-submission
// expression path — the option-building hot path behind forms.Form,
// FilterInput, and FilterDropdown.Wire (defaults applied per render).
func BenchmarkFormWireExpression(b *testing.B) {
	unwired := Action{URL: ""}

	htmxForm := Action{Method: MethodPost, URL: "/api/submit", Target: "#form-out"}
	datastarForm := Action{
		Transport:   TransportDatastar,
		Method:      MethodPost,
		URL:         "/api/submit",
		Selector:    "#form-out",
		ContentType: ContentTypeForm,
		DebounceMS:  300,
	}

	b.Run("inert empty URL", func(b *testing.B) {
		b.ReportAllocs()

		var sink templ.Attributes

		for b.Loop() {
			sink = formWireAttributesLike(&unwired)
		}

		if sink == nil {
			b.Fatal("inert path must return nil")
		}
	})

	b.Run("htmx form defaults", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = formWireAttributesLike(&htmxForm)
		}
	})

	b.Run("datastar form expression (selector+contentType+debounce)", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = datastarForm.Attributes()
		}
	})
}

// formWireAttributesLike mirrors forms.formWireAttributes (the copy +
// defaults + render path) without importing the forms package (which would
// cycle the module graph in tests).
func formWireAttributesLike(w *Action) templ.Attributes {
	if w == nil {
		return nil
	}

	action := *w

	if action.Event == EventUnspecified {
		action.Event = EventSubmit
	}

	if action.ContentType == ContentTypeUnspecified {
		action.ContentType = ContentTypeForm
	}

	return action.Attributes()
}

// TestTriggerModifiersRequireAnEventTrigger pins the debounce/throttle
// orphan-modifier invariant: a delay or throttle window with NO explicit,
// valid event trigger renders nothing in EITHER dialect. htmx cannot hang a
// modifier on its implicit default trigger (hx-trigger needs the event
// name), and Datastar must match — the earlier behavior attached the
// orphan modifier to the defaulted click, so the same Action silently
// debounce'd under Datastar but not under htmx. The action itself still
// fires (htmx element default / data-on:click); only the orphan modifier
// vanishes.
func TestTriggerModifiersRequireAnEventTrigger(t *testing.T) {
	t.Parallel()

	orphanModifiers := []struct {
		name         string
		mutate       func(*Action)
		bannedNeedle string
	}{
		{"debounce without event", func(a *Action) { a.DebounceMS = 150 }, "__debounce"},
		{"throttle without event", func(a *Action) { a.ThrottleMS = 150 }, "__throttle"},
		{"debounce with interval, no event", func(a *Action) {
			a.DebounceMS = 150
			a.Interval = "5s"
		}, "__debounce"},
		{"invalid event counts as no trigger", func(a *Action) {
			a.Event = Event("hover")
			a.DebounceMS = 150
		}, "__debounce"},
	}

	for _, tt := range orphanModifiers {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, transport := range []Transport{TransportHTMX, TransportDatastar} {
				action := Action{Transport: transport, Method: MethodGet, URL: "/api/items"}
				tt.mutate(&action)

				rendered := renderAttributes(t, action.Attributes())

				if transport == TransportDatastar && strings.Contains(rendered, tt.bannedNeedle) {
					t.Fatalf("orphan %s modifier must be dropped under datastar, got %s", tt.bannedNeedle, rendered)
				}

				if transport == TransportHTMX {
					if strings.Contains(rendered, "delay:") || strings.Contains(rendered, "throttle:") {
						t.Fatalf("orphan modifier must be dropped under htmx, got %s", rendered)
					}
				}
			}
		})
	}
}

// actionFieldDialects is the drift-guard contract for every `wire.Action`
// field: which dialect(s) the field must render in, given its probe context.
// TestEveryActionFieldHasADialectContract reflects over Action's exported
// fields, so a NEWLY ADDED field fails here until it gets a row — a dialect
// wired in one spelling only can no longer ship silently (the htmx v4
// event-renames class, TODO #316b, is exactly the drift this pins).
//
// Fields whose probe needs context (a modifier needs an event to attach to;
// the selector needs ContentTypeForm to be read at all) carry contextFields
// merged into the probed action AND the baseline it is diffed against.
var actionFieldDialects = map[string]struct {
	htmx          bool
	datastar      bool
	contextFields []string
}{
	"Method":         {htmx: true, datastar: true},
	"Event":          {htmx: true, datastar: true},
	"Target":         {htmx: true, datastar: false},
	"ContentType":    {htmx: false, datastar: true},
	"Selector":       {htmx: false, datastar: true, contextFields: []string{"ContentType"}},
	"Swap":           {htmx: true, datastar: false},
	"DebounceMS":     {htmx: true, datastar: true, contextFields: []string{"Event"}},
	"ThrottleMS":     {htmx: true, datastar: true, contextFields: []string{"Event"}},
	"PreventDefault": {htmx: false, datastar: true},
	"Interval":       {htmx: true, datastar: true},
	"Reveal":         {htmx: true, datastar: true},
}

// skippedActionFields are not probeable one-field-at-a-time: Transport IS
// the dialect selector (the test renders both dialects itself) and URL is
// the base every probe stands on (its parity is pinned by
// TestURLReferencedInBothDialects).
var skippedActionFields = map[string]string{
	"Transport": "dialect selector — rendered by the test itself",
	"URL":       "the base field; parity pinned by TestURLReferencedInBothDialects",
}

// TestEveryActionFieldHasADialectContract reflects over every exported
// Action field, probes it one field at a time in both dialects, and fails
// when (a) a field has no row in actionFieldDialects, (b) a mapped field
// stops rendering where the map says it must, or (c) a field starts
// rendering where the map says it must not (a dialect leak).
func TestEveryActionFieldHasADialectContract(t *testing.T) {
	t.Parallel()

	base := Action{Method: MethodGet, URL: "/api/probe"}
	baseHTMX := canonicalAttrs(t, base.Attributes())
	baseDatastar := canonicalAttrs(t, base.TransportDatastarAction().Attributes())

	fields := reflect.VisibleFields(reflect.TypeOf(base))

	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		name := field.Name

		if why, skip := skippedActionFields[name]; skip {
			t.Logf("skipped %s: %s", name, why)

			continue
		}

		contract, mapped := actionFieldDialects[name]
		if !mapped {
			t.Errorf("Action field %q has no dialect contract: add it to actionFieldDialects "+
				"and render it in both dialects (or document the asymmetry there) — "+
				"a field wired into one dialect only ships a silent consumer bug", name)

			continue
		}

		probe := base
		for _, ctx := range contract.contextFields {
			applyProbeValue(&probe, ctx)
		}
		applyProbeValue(&probe, name)

		gotHTMX := canonicalAttrs(t, probe.Attributes()) != baseHTMX
		gotDatastar := canonicalAttrs(t, probe.TransportDatastarAction().Attributes()) != baseDatastar

		if gotHTMX != contract.htmx || gotDatastar != contract.datastar {
			t.Errorf("Action field %q dialect drift: contract says htmx=%v datastar=%v, probe rendered htmx=%v datastar=%v",
				name, contract.htmx, contract.datastar, gotHTMX, gotDatastar)
		}
	}
}

// TransportDatastarAction returns a copy of the action pinned to the
// Datastar dialect — the probe's second rendering.
func (a Action) TransportDatastarAction() Action {
	a.Transport = TransportDatastar

	return a
}

// probeValues are representative values per field name; every field in
// actionFieldDialects must have one.
var probeValues = map[string]func(*Action){
	"Method":         func(a *Action) { a.Method = MethodPost },
	"Event":          func(a *Action) { a.Event = EventSubmit },
	"Target":         func(a *Action) { a.Target = "#probe-target" },
	"ContentType":    func(a *Action) { a.ContentType = ContentTypeForm },
	"Selector":       func(a *Action) { a.Selector = "#probe-form" },
	"Swap":           func(a *Action) { a.Swap = PatchModeOuter },
	"DebounceMS":     func(a *Action) { a.DebounceMS = 150 },
	"ThrottleMS":     func(a *Action) { a.ThrottleMS = 150 },
	"PreventDefault": func(a *Action) { a.PreventDefault = true },
	"Interval":       func(a *Action) { a.Interval = "5s" },
	"Reveal":         func(a *Action) { a.Reveal = &Reveal{ThresholdPercent: 50} },
}

func applyProbeValue(action *Action, field string) {
	mutate, ok := probeValues[field]
	if !ok {
		panic("no probe value for Action field " + field + " — add one to probeValues")
	}

	mutate(action)
}

// canonicalAttrs renders attributes as a deterministic sorted string so
// two renderings compare by content.
func canonicalAttrs(t *testing.T, attrs templ.Attributes) string {
	t.Helper()

	if attrs == nil {
		return "<nil>"
	}

	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var b strings.Builder
	for _, key := range keys {
		fmt.Fprintf(&b, "%s=%v;", key, attrs[key])
	}

	return b.String()
}
