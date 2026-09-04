package wire

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestTransportIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value Transport
		valid bool
	}{
		{name: "unspecified is the defined zero value", value: TransportUnspecified, valid: true},
		{name: "htmx", value: TransportHTMX, valid: true},
		{name: "datastar", value: TransportDatastar, valid: true},
		{name: "unknown", value: Transport("sse"), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := TransportIsValid(tt.value); got != tt.valid {
				t.Errorf(
					"TransportIsValid(%q) = %v, want %v", tt.value, got, tt.valid,
				)
			}
		})
	}
}

func TestMethodIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value Method
		valid bool
	}{
		{name: "unspecified is the defined zero value", value: MethodUnspecified, valid: true},
		{name: "get", value: MethodGet, valid: true},
		{name: "post", value: MethodPost, valid: true},
		{name: "put", value: MethodPut, valid: true},
		{name: "patch", value: MethodPatch, valid: true},
		{name: "delete", value: MethodDelete, valid: true},
		{name: "unknown", value: Method("fetch"), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := MethodIsValid(tt.value); got != tt.valid {
				t.Errorf("MethodIsValid(%q) = %v, want %v", tt.value, got, tt.valid)
			}
		})
	}
}

func TestEventIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value Event
		valid bool
	}{
		{name: "unspecified is the defined zero value", value: EventUnspecified, valid: true},
		{name: "click", value: EventClick, valid: true},
		{name: "submit", value: EventSubmit, valid: true},
		{name: "change", value: EventChange, valid: true},
		{name: "input", value: EventInput, valid: true},
		{name: "keydown", value: EventKeyDown, valid: true},
		{name: "keyup", value: EventKeyUp, valid: true},
		{name: "focus", value: EventFocus, valid: true},
		{name: "blur", value: EventBlur, valid: true},
		{name: "unknown", value: Event("hover"), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := EventIsValid(tt.value); got != tt.valid {
				t.Errorf("EventIsValid(%q) = %v, want %v", tt.value, got, tt.valid)
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
			name: "empty URL wires nothing",
			action: Action{
				Transport: TransportHTMX,
				Method:    MethodGet,
				URL:       "",
				Event:     EventUnspecified,
				Target:    "",
			},
			expected: nil,
		},
		{
			name: "htmx defaults (zero values resolve to GET)",
			action: Action{
				Transport: TransportUnspecified,
				Method:    MethodUnspecified,
				URL:       "/api/items",
				Event:     EventUnspecified,
				Target:    "",
			},
			expected: templ.Attributes{"hx-get": "/api/items"},
		},
		{
			name: "htmx full",
			action: Action{
				Transport: TransportHTMX,
				Method:    MethodPost,
				URL:       "/api/items",
				Event:     EventSubmit,
				Target:    "#list",
			},
			expected: templ.Attributes{
				"hx-post":    "/api/items",
				"hx-trigger": "submit",
				"hx-target":  "#list",
			},
		},
		{
			name: "htmx unknown method falls back to get",
			action: Action{
				Transport: TransportHTMX,
				Method:    Method("fetch"),
				URL:       "/api/items",
				Event:     EventUnspecified,
				Target:    "",
			},
			expected: templ.Attributes{"hx-get": "/api/items"},
		},
		{
			name: "htmx unknown event is omitted (element default applies)",
			action: Action{
				Transport: TransportHTMX,
				Method:    MethodDelete,
				URL:       "/api/items/1",
				Event:     Event("hover"),
				Target:    "",
			},
			expected: templ.Attributes{"hx-delete": "/api/items/1"},
		},
		{
			name: "datastar defaults (click event injected)",
			action: Action{
				Transport: TransportDatastar,
				Method:    MethodUnspecified,
				URL:       "/api/items",
				Event:     EventUnspecified,
				Target:    "",
			},
			expected: templ.Attributes{"data-on:click": "@get('/api/items')"},
		},
		{
			name: "datastar full (target is response-driven, never emitted)",
			action: Action{
				Transport: TransportDatastar,
				Method:    MethodPost,
				URL:       "/api/items",
				Event:     EventChange,
				Target:    "#list",
			},
			expected: templ.Attributes{"data-on:change": "@post('/api/items')"},
		},
		{
			name: "datastar single quotes in URL are escaped",
			action: Action{
				Transport: TransportDatastar,
				Method:    MethodGet,
				URL:       "/api/search?q=it's",
				Event:     EventUnspecified,
				Target:    "",
			},
			expected: templ.Attributes{"data-on:click": `@get('/api/search?q=it\'s')`},
		},
		{
			name: "datastar unknown event falls back to click",
			action: Action{
				Transport: TransportDatastar,
				Method:    MethodDelete,
				URL:       "/api/items/1",
				Event:     Event("hover"),
				Target:    "",
			},
			expected: templ.Attributes{"data-on:click": "@delete('/api/items/1')"},
		},
		{
			name: "datastar unknown method falls back to get",
			action: Action{
				Transport: TransportDatastar,
				Method:    Method("fetch"),
				URL:       "/api/items",
				Event:     EventUnspecified,
				Target:    "",
			},
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

			html := renderAttributes(t, tt.action.Attributes())

			for _, want := range tt.contains {
				if !strings.Contains(html, want) {
					t.Errorf("rendered %q does not contain %q", html, want)
				}
			}
		})
	}
}

// renderAttributes mimics what a templ attribute spread compiles to.
func renderAttributes(t *testing.T, attrs templ.Attributes) string {
	t.Helper()

	var sb strings.Builder

	if _, err := io.WriteString(&sb, "<button"); err != nil {
		t.Fatalf("write open tag: %v", err)
	}

	if err := templ.RenderAttributes(context.Background(), &sb, attrs); err != nil {
		t.Fatalf("render attributes: %v", err)
	}

	if _, err := io.WriteString(&sb, "></button>"); err != nil {
		t.Fatalf("write close tag: %v", err)
	}

	return sb.String()
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
