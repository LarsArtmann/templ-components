package wire

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// TestConstructors pins the five method constructors: each returns a fresh
// spec with the URL set and the constructor's verb — Get deliberately leaves
// Method at its zero value (the read default renders GET in both dialects).
func TestConstructors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		build      func(url string) *Action
		wantMethod Method
	}{
		{name: "Get is the read default (method unspecified)", build: Get, wantMethod: MethodUnspecified},
		{name: "Post", build: Post, wantMethod: MethodPost},
		{name: "Put", build: Put, wantMethod: MethodPut},
		{name: "Patch", build: Patch, wantMethod: MethodPatch},
		{name: "Delete", build: Delete, wantMethod: MethodDelete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := tt.build("/api/thing")
			if action == nil {
				t.Fatal("constructor returned nil")
			}

			if action.URL != "/api/thing" {
				t.Errorf("URL = %q, want /api/thing", action.URL)
			}

			if action.Method != tt.wantMethod {
				t.Errorf("Method = %q, want %q", action.Method, tt.wantMethod)
			}

			if action.Transport != TransportUnspecified {
				t.Errorf("Transport = %q, want the htmx zero default", action.Transport)
			}
		})
	}
}

// TestGetConstructorRendersGET pins the read default end-to-end: a Get spec
// with no explicit method still renders hx-get (the zero value resolves to
// GET, never to a broken verb).
func TestGetConstructorRendersGET(t *testing.T) {
	t.Parallel()

	attrs := Get("/api/fragment").Attributes()
	if attrs["hx-get"] != "/api/fragment" {
		t.Errorf("hx-get = %v, want /api/fragment", attrs["hx-get"])
	}
}

// TestBuildersReturnCopiesNeverMutate pins the builder contract: every With*
// is a value receiver that returns a modified COPY — the consumer's *Action
// keeps its previous field values (the same non-mutation contract
// kanbanWireAttributes/formWireAttributes pin at the component boundary).
func TestBuildersReturnCopiesNeverMutate(t *testing.T) {
	t.Parallel()

	original := Post("/api/save")
	originalSnapshot := *original

	chained := original.
		WithTarget("#out").
		WithSwap(PatchModeOuter).
		WithContentType(ContentTypeForm).
		WithDebounce(300).
		WithThrottle(1000).
		WithEvent(EventSubmit).
		WithMethod(MethodPut).
		WithSelector("#the-form").
		WithInterval("10s").
		WithReveal(Reveal{ThresholdPercent: 50}).
		WithPreventDefault()

	if chained == originalSnapshot {
		t.Fatal("builder chain returned an identical value — the copy must carry the chain's changes")
	}

	// The copy carries every chained change.
	for _, c := range []struct {
		name string
		got  any
		want any
	}{
		{"Target", chained.Target, "#out"},
		{"Swap", chained.Swap, PatchModeOuter},
		{"ContentType", chained.ContentType, ContentTypeForm},
		{"DebounceMS", chained.DebounceMS, 300},
		{"ThrottleMS", chained.ThrottleMS, 1000},
		{"Event", chained.Event, EventSubmit},
		{"Method", chained.Method, MethodPut},
		{"Selector", chained.Selector, "#the-form"},
		{"Interval", chained.Interval, "10s"},
	} {
		if c.got != c.want {
			t.Errorf("chained.%s = %v, want %v", c.name, c.got, c.want)
		}
	}

	if chained.Reveal == nil || chained.Reveal.ThresholdPercent != 50 {
		t.Errorf("chained.Reveal = %+v, want threshold 50", chained.Reveal)
	}

	if !chained.PreventDefault {
		t.Error("chained.PreventDefault = false, want true")
	}

	// The original is untouched by the chain.
	if original.Target != "" || original.Swap != PatchModeUnspecified ||
		original.ContentType != ContentTypeUnspecified || original.DebounceMS != 0 ||
		original.ThrottleMS != 0 || original.Event != EventUnspecified ||
		original.Method != MethodPost || original.Selector != "" ||
		original.Interval != "" || original.Reveal != nil || original.PreventDefault {
		t.Errorf("builder chain mutated the receiver: %+v", *original)
	}
}

// TestWithFormDefaults pins the shared host-default application: an
// UNSPECIFIED event adopts the host's natural event, an UNSPECIFIED content
// type becomes form encoding — and EXPLICIT values always win over the
// defaults (the consumer's override is never clobbered).
func TestWithFormDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		action      Action
		wantEvent   Event
		wantContent ContentType
	}{
		{
			name:        "unspecified event and content adopt the host defaults",
			action:      Action{URL: "/api/save", Method: MethodPost},
			wantEvent:   EventSubmit,
			wantContent: ContentTypeForm,
		},
		{
			name: "explicit event wins over the host default",
			action: Action{
				URL: "/api/save", Method: MethodPost, Event: EventChange,
			},
			wantEvent:   EventChange,
			wantContent: ContentTypeForm,
		},
		{
			name: "explicit json content wins over the form default",
			action: Action{
				URL: "/api/save", Method: MethodPost, ContentType: ContentTypeJSON,
			},
			wantEvent:   EventSubmit,
			wantContent: ContentTypeJSON,
		},
		{
			name: "explicit everything is passed through untouched",
			action: Action{
				URL: "/api/save", Method: MethodPost,
				Event: EventInput, ContentType: ContentTypeJSON,
			},
			wantEvent:   EventInput,
			wantContent: ContentTypeJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			before := tt.action

			got := tt.action.WithFormDefaults(EventSubmit)

			if got.Event != tt.wantEvent {
				t.Errorf("Event = %q, want %q", got.Event, tt.wantEvent)
			}

			if got.ContentType != tt.wantContent {
				t.Errorf("ContentType = %q, want %q", got.ContentType, tt.wantContent)
			}

			if tt.action != before {
				t.Errorf("WithFormDefaults mutated the receiver: %+v", tt.action)
			}
		})
	}
}

// TestThrottleHTMXRendering pins the throttle window under htmx: the
// throttle:<n>ms trigger modifier, nothing at zero, and separate modifiers
// when debounce pipelines alongside (htmx documents them as exclusive — the
// rendering simply emits both; prefer one).
func TestThrottleHTMXRendering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		action    Action
		wantParts []string
		wantNone  []string
	}{
		{
			name:      "throttle renders as the trigger modifier",
			action:    Get("/api/search").WithEvent(EventInput).WithThrottle(500),
			wantParts: []string{"throttle:500ms"},
		},
		{
			name:     "zero throttle renders nothing",
			action:   Get("/api/search").WithEvent(EventClick),
			wantNone: []string{"throttle"},
		},
		{
			name:      "debounce and throttle render as separate modifiers",
			action:    Get("/api/search").WithEvent(EventInput).WithDebounce(200).WithThrottle(1000),
			wantParts: []string{"delay:200ms", "throttle:1000ms", "changed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			attrs := tt.action.Attributes()

			trigger, _ := attrs["hx-trigger"].(string)

			for _, want := range tt.wantParts {
				if !strings.Contains(trigger, want) {
					t.Errorf("hx-trigger = %q, want it to contain %q", trigger, want)
				}
			}

			for _, banned := range tt.wantNone {
				if strings.Contains(trigger, banned) {
					t.Errorf("hx-trigger = %q, want no %q token", trigger, banned)
				}
			}
		})
	}
}

// TestThrottleDatastarRendering pins the throttle window under Datastar: the
// __throttle.<n>ms attribute-name modifier, nothing at zero, and the
// delay-then-throttle pipeline in ONE attribute key.
func TestThrottleDatastarRendering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		action      Action
		wantNeedles []string
		wantBanned  string
	}{
		{
			name:        "throttle renders as the attribute-name modifier",
			action:      Get("/api/search").WithTransport(TransportDatastar).WithEvent(EventInput).WithThrottle(500),
			wantNeedles: []string{"__throttle.500ms"},
		},
		{
			name:       "zero throttle renders nothing",
			action:     Get("/api/search").WithTransport(TransportDatastar).WithEvent(EventClick),
			wantBanned: "__throttle",
		},
		{
			name: "debounce and throttle pipeline in one key",
			action: Get(
				"/api/search",
			).WithTransport(TransportDatastar).
				WithEvent(EventInput).
				WithDebounce(200).
				WithThrottle(1000),
			wantNeedles: []string{"__debounce.200ms", "__throttle.1000ms"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			attrs := tt.action.Attributes()

			keys := make([]string, 0, len(attrs))
			for key := range attrs {
				keys = append(keys, key)
			}

			joined := strings.Join(keys, "\n")

			for _, want := range tt.wantNeedles {
				if !strings.Contains(joined, want) {
					t.Errorf("no datastar attribute carries %s: %v", want, attrs)
				}
			}

			if tt.wantBanned != "" && strings.Contains(joined, tt.wantBanned) {
				t.Errorf("an attribute carries the banned modifier %s: %v", tt.wantBanned, attrs)
			}
		})
	}
}

// TestWithTransportComposesWithConstructors pins the Transport chain: the
// constructors return the htmx zero default and WithTransport flips the
// dialect without touching anything else.
func TestWithTransportComposesWithConstructors(t *testing.T) {
	t.Parallel()

	action := Post("/api/save").WithTransport(TransportDatastar)

	if action.Transport != TransportDatastar {
		t.Errorf("Transport = %q, want datastar", action.Transport)
	}

	if action.Method != MethodPost || action.URL != "/api/save" {
		t.Errorf("WithTransport altered other fields: %+v", action)
	}

	attrs := action.Attributes()

	var datastarKey bool

	for key := range attrs {
		if strings.HasPrefix(key, "data-on:") {
			datastarKey = true
		}

		if strings.HasPrefix(key, "hx-") {
			t.Errorf("datastar action rendered htmx attribute %q", key)
		}
	}

	if !datastarKey {
		t.Error("datastar action rendered no data-on attribute")
	}
}

// TestTriggerAttributesAreStrings pins the attribute value type contract:
// templ.Attributes values rendered from an Action are strings (a typed value
// would escape wrongly in templ's attribute context).
func TestTriggerAttributesAreStrings(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		action Action
	}{
		{name: "htmx interval", action: Get("/u").WithInterval("5s")},
		{name: "htmx reveal", action: Get("/u").WithReveal(Reveal{})},
		{name: "datastar interval", action: Get("/u").WithTransport(TransportDatastar).WithInterval("5s")},
		{name: "datastar reveal", action: Get("/u").WithTransport(TransportDatastar).WithReveal(Reveal{Exit: true})},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for key, value := range tc.action.Attributes() {
				if _, ok := value.(string); !ok {
					t.Errorf("attribute %q has non-string value %T", key, value)
				}
			}
		})
	}
}

// Compile-time guards: the Attributes methods exist on both the value and
// pointer forms the builders traffic in.
var (
	_ func() templ.Attributes        = (Action{}).Attributes
	_ func(*Action) templ.Attributes = (*Action).Attributes
)
