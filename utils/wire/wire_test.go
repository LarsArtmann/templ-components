package wire

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestIsValidEnums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		valid bool
		value any
		check func(any) bool
	}{
		{name: "TransportUnspecified", value: TransportUnspecified, check: func(v any) bool { return TransportIsValid(v.(Transport)) }, valid: true},
		{name: "TransportHTMX", value: TransportHTMX, check: func(v any) bool { return TransportIsValid(v.(Transport)) }, valid: true},
		{name: "TransportDatastar", value: TransportDatastar, check: func(v any) bool { return TransportIsValid(v.(Transport)) }, valid: true},
		{name: "TransportUnknown", value: Transport("sse"), check: func(v any) bool { return TransportIsValid(v.(Transport)) }, valid: false},
		{name: "MethodUnspecified", value: MethodUnspecified, check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: true},
		{name: "MethodGet", value: MethodGet, check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: true},
		{name: "MethodPost", value: MethodPost, check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: true},
		{name: "MethodPut", value: MethodPut, check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: true},
		{name: "MethodPatch", value: MethodPatch, check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: true},
		{name: "MethodDelete", value: MethodDelete, check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: true},
		{name: "MethodUnknown", value: Method("fetch"), check: func(v any) bool { return MethodIsValid(v.(Method)) }, valid: false},
		{name: "EventUnspecified", value: EventUnspecified, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventClick", value: EventClick, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventSubmit", value: EventSubmit, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventChange", value: EventChange, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventInput", value: EventInput, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventKeyDown", value: EventKeyDown, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventKeyUp", value: EventKeyUp, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventFocus", value: EventFocus, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventBlur", value: EventBlur, check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: true},
		{name: "EventUnknown", value: Event("hover"), check: func(v any) bool { return EventIsValid(v.(Event)) }, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.check(tt.value); got != tt.valid {
				t.Errorf("IsValid(%v) = %v, want %v", tt.value, got, tt.valid)
			}
		})
	}
}

func TestActionAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		action   Action
		expected templ.Attributes
	}{
		{
			name:     "empty URL wires nothing",
			action:   Action{Transport: TransportHTMX, Method: MethodGet, URL: "", Event: EventUnspecified, Target: ""},
			expected: nil,
		},
		{
			name:     "htmx defaults (zero values resolve to GET)",
			action:   Action{Transport: TransportUnspecified, Method: MethodUnspecified, URL: "/api/items", Event: EventUnspecified, Target: ""},
			expected: templ.Attributes{"hx-get": "/api/items"},
		},
		{
			name:     "htmx full",
			action:   Action{Transport: TransportHTMX, Method: MethodPost, URL: "/api/items", Event: EventSubmit, Target: "#list"},
			expected: templ.Attributes{"hx-post": "/api/items", "hx-trigger": "submit", "hx-target": "#list"},
		},
		{
			name:     "htmx unknown method falls back to get",
			action:   Action{Transport: TransportHTMX, Method: Method("fetch"), URL: "/api/items", Event: EventUnspecified, Target: ""},
			expected: templ.Attributes{"hx-get": "/api/items"},
		},
		{
			name:     "htmx unknown event is omitted (element default applies)",
			action:   Action{Transport: TransportHTMX, Method: MethodDelete, URL: "/api/items/1", Event: Event("hover"), Target: ""},
			expected: templ.Attributes{"hx-delete": "/api/items/1"},
		},
		{
			name:     "datastar defaults (click event injected)",
			action:   Action{Transport: TransportDatastar, Method: MethodUnspecified, URL: "/api/items", Event: EventUnspecified, Target: ""},
			expected: templ.Attributes{"data-on:click": "@get('/api/items')"},
		},
		{
			name:     "datastar full (target is response-driven, never emitted)",
			action:   Action{Transport: TransportDatastar, Method: MethodPost, URL: "/api/items", Event: EventChange, Target: "#list"},
			expected: templ.Attributes{"data-on:change": "@post('/api/items')"},
		},
		{
			name:     "datastar single quotes in URL are escaped",
			action:   Action{Transport: TransportDatastar, Method: MethodGet, URL: "/api/search?q=it's", Event: EventUnspecified, Target: ""},
			expected: templ.Attributes{"data-on:click": `@get('/api/search?q=it\'s')`},
		},
		{
			name:     "datastar unknown event falls back to click",
			action:   Action{Transport: TransportDatastar, Method: MethodDelete, URL: "/api/items/1", Event: Event("hover"), Target: ""},
			expected: templ.Attributes{"data-on:click": "@delete('/api/items/1')"},
		},
		{
			name:     "datastar unknown method falls back to get",
			action:   Action{Transport: TransportDatastar, Method: Method("fetch"), URL: "/api/items", Event: EventUnspecified, Target: ""},
			expected: templ.Attributes{"data-on:click": "@get('/api/items')"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.action.Attributes()

			if tt.expected == nil {
				if got != nil {
					t.Fatalf("Attributes() = %v, want nil", got)
				}

				return
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("Attributes() = %v, want %v", got, tt.expected)
			}

			for key, want := range tt.expected {
				if got[key] != want {
					t.Errorf("Attributes()[%q] = %v, want %v", key, got[key], want)
				}
			}
		})
	}
}

// TestActionAttributesRender proves the rendered attribute spelling: htmx
// dialect keys and the data-on:<event> colon key must survive templ's
// attribute writer verbatim.
func TestActionAttributesRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		action   Action
		contains []string
	}{
		{
			name:     "htmx",
			action:   Action{Method: MethodPost, URL: "/api/save", Target: "#out"},
			contains: []string{`hx-post="/api/save"`, `hx-target="#out"`},
		},
		{
			name:   "datastar colon key (single quotes are HTML-entity encoded in source and decode back in the DOM)",
			action: Action{Transport: TransportDatastar, Method: MethodGet, URL: "/api/items"},
			contains: []string{
				`data-on:click="@get(&#39;/api/items&#39;)"`,
				`data-on:click="@get(`,
				`/api/items&#39;)"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html, err := renderAttributes(tt.action.Attributes())
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			for _, want := range tt.contains {
				if !strings.Contains(html, want) {
					t.Errorf("rendered %q does not contain %q", html, want)
				}
			}
		})
	}
}

// renderAttributes mimics what a templ attribute spread compiles to.
func renderAttributes(attrs templ.Attributes) (string, error) {
	var sb strings.Builder

	if _, err := io.WriteString(&sb, "<button"); err != nil {
		return "", err
	}

	if err := templ.RenderAttributes(context.Background(), &sb, attrs); err != nil {
		return "", err
	}

	if _, err := io.WriteString(&sb, "></button>"); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// FuzzAction verifies attribute rendering never panics and never emits an
// empty-valued attribute on arbitrary input.
func FuzzAction(f *testing.F) {
	f.Add("htmx", "post", "click", "/api/items")
	f.Add("datastar", "", "", "")
	f.Add("", "fetch", "hover", "it's")
	f.Add("unknown", "GET", "load", "#")

	f.Fuzz(func(t *testing.T, transport, method, event, url string) {
		action := Action{
			Transport: Transport(transport),
			Method:    Method(method),
			URL:       url,
			Event:     Event(event),
		}

		attrs := action.Attributes()
		for key, value := range attrs {
			if key == "" {
				t.Fatalf("empty attribute key for %+v", action)
			}

			if value == "" {
				t.Fatalf("attribute %q has empty value for %+v", key, action)
			}
		}
	})
}
