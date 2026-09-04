package wire

import (
	"strings"
	"testing"
)

// Invariants that must hold for EVERY Action, regardless of field values.
// They pin the contract against silent rot: the same URL appears in both
// dialects, an empty URL stays inert for every field combination, and wire
// rendering never emits a script (the contract is attributes-only by design —
// CSP-safe without a nonce).

func TestEmptyURLIsInertForEveryEnumCombo(t *testing.T) {
	t.Parallel()

	transports := []Transport{TransportUnspecified, TransportHTMX, TransportDatastar, Transport("bogus")}
	methods := []Method{MethodUnspecified, MethodGet, MethodPost, MethodPut, MethodPatch, MethodDelete, Method("bogus")}
	events := []Event{
		EventUnspecified, EventClick, EventSubmit, EventChange, EventInput,
		EventKeyDown, EventKeyUp, EventFocus, EventBlur, Event("hover"),
	}
	targets := []string{"", "#out", "closest div"}

	for _, transport := range transports {
		for _, method := range methods {
			for _, event := range events {
				for _, target := range targets {
					action := Action{
						Transport: transport,
						Method:    method,
						URL:       "",
						Event:     event,
						Target:    target,
					}

					if attrs := action.Attributes(); attrs != nil {
						t.Fatalf("empty URL must wire nothing, got %v for %+v", attrs, action)
					}
				}
			}
		}
	}
}

func TestURLReferencedInBothDialects(t *testing.T) {
	t.Parallel()

	transports := []Transport{TransportUnspecified, TransportHTMX, TransportDatastar}
	methods := []Method{MethodUnspecified, MethodGet, MethodPost, MethodPut, MethodPatch, MethodDelete}
	events := []Event{EventUnspecified, EventClick, EventSubmit, EventChange, EventInput}
	targets := []string{"", "#out"}

	const url = "/api/items?filter=x"

	for _, transport := range transports {
		for _, method := range methods {
			for _, event := range events {
				for _, target := range targets {
					action := Action{
						Transport: transport,
						Method:    method,
						URL:       url,
						Event:     event,
						Target:    target,
					}

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

func TestWireRenderingEmitsNoScript(t *testing.T) {
	t.Parallel()

	actions := []Action{
		{URL: "/api/items"},
		{Method: MethodPost, URL: "/api/items"},
		{Transport: TransportDatastar, URL: "/api/items"},
		{Transport: TransportDatastar, Method: MethodPost, URL: "/api/items?filter=it's", Event: EventInput, Target: "#out"},
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
	datastar := Action{Transport: TransportDatastar, Method: MethodPost, URL: "/api/items", Event: EventClick}

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
