package wire

import (
	"fmt"
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

// TestSwapDialectIsolation pins that the shared Swap/PatchMode vocabulary
// renders in the RIGHT dialect and never leaks across: the htmx dialect
// renders hx-swap and no datastar mode option; the datastar dialect renders
// the mode fetch option and no hx-* attribute.
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

	t.Run("datastar renders the mode option only", func(t *testing.T) {
		t.Parallel()

		for _, swap := range swaps {
			attrs := Action{
				Transport: TransportDatastar, Method: MethodGet, URL: "/api/items", Swap: swap,
			}.Attributes()

			if !datastarAttrsContain(attrs, "mode: '"+string(swap)+"'") {
				t.Fatalf("datastar dialect must render the mode option for %q, got %v", swap, attrs)
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
