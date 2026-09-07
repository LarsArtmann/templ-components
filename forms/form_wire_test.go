package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestFormWire verifies the transport-agnostic submission wiring: the form's
// fields serialize in both dialects (htmx natively, Datastar via
// contentType:'form'), the form-level defaults land (submit event, form
// encoding), and the wiring stays inert when unset or URL-less.
func TestFormWire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		props       FormProps
		contains    []string
		notContains []string
	}{
		{
			name: "nil wire renders no wiring attributes",
			props: FormProps{
				Action: "/save",
				Method: FormPost,
			},
			notContains: []string{"hx-get", "hx-post", "data-on:", "hx-trigger"},
		},
		{
			name: "htmx dialect (unspecified transport resolves to htmx)",
			props: FormProps{
				Action: "/save",
				Method: FormPost,
				Wire: &wire.Action{
					Method: wire.MethodPost,
					URL:    "/api/save",
				},
			},
			contains: []string{
				`hx-post="/api/save"`,
				`hx-trigger="submit"`,
				`method="POST"`,
				`action="/save"`,
			},
		},
		{
			name: "htmx dialect renders the target (client-side patching)",
			props: FormProps{
				Wire: &wire.Action{
					Method: wire.MethodPost,
					URL:    "/api/save",
					Target: "#form-out",
				},
			},
			contains: []string{`hx-target="#form-out"`},
		},
		{
			name: "datastar dialect (form fields serialize via contentType form)",
			props: FormProps{
				Method: FormPost,
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					Method:    wire.MethodPost,
					URL:       "/api/save",
				},
			},
			contains: []string{
				`data-on:submit="@post(&#39;/api/save&#39;, {contentType: &#39;form&#39;})"`,
			},
			notContains: []string{"hx-post", "hx-trigger", "hx-target"},
		},
		{
			name: "datastar explicit json content type is respected (signals submit)",
			props: FormProps{
				Wire: &wire.Action{
					Transport:   wire.TransportDatastar,
					Method:      wire.MethodPost,
					URL:         "/api/save",
					ContentType: wire.ContentTypeJSON,
				},
			},
			contains:    []string{`data-on:submit="@post(&#39;/api/save&#39;)"`},
			notContains: []string{"contentType: &#39;form&#39;"},
		},
		{
			name: "datastar target is response-driven, not emitted",
			props: FormProps{
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					Method:    wire.MethodPost,
					URL:       "/api/save",
					Target:    "#form-out",
				},
			},
			notContains: []string{"target"},
		},
		{
			name: "empty wire URL stays inert (plain HTML form fallback)",
			props: FormProps{
				Action: "/save",
				Method: FormPost,
				Wire: &wire.Action{
					Method: wire.MethodPost,
					URL:    "",
				},
			},
			contains: []string{
				`action="/save"`,
				`method="POST"`,
			},
			notContains: []string{"hx-post", "data-on:"},
		},
		{
			name: "explicit event is not overridden by the form default",
			props: FormProps{
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					Method:    wire.MethodGet,
					URL:       "/api/save",
					Event:     wire.EventChange,
				},
			},
			contains:    []string{`data-on:change=`},
			notContains: []string{"data-on:submit"},
		},
		{
			name: "wire composes with Validate (validation parity in both dialects)",
			props: FormProps{
				Validate: true,
				Wire: &wire.Action{
					Method: wire.MethodPost,
					URL:    "/api/save",
				},
			},
			contains: []string{`hx-validate="true"`, `hx-post="/api/save"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, Form(tt.props))

			for _, want := range tt.contains {
				utils.AssertContains(t, output, want)
			}

			for _, banned := range tt.notContains {
				utils.AssertNotContains(t, output, banned)
			}
		})
	}
}

// TestFormWireDoesNotMutateAction pins the copy semantics: applying the form
// defaults must never write back into the consumer's wire.Action — a shared
// action (e.g. one spec used by a Button and a Form) would otherwise drift.
func TestFormWireDoesNotMutateAction(t *testing.T) {
	t.Parallel()

	action := &wire.Action{Method: wire.MethodPost, URL: "/api/save"}
	_ = utils.Render(t, Form(FormProps{Wire: action}))

	if action.Event != wire.EventUnspecified {
		t.Errorf("consumer action Event mutated to %q", action.Event)
	}

	if action.ContentType != wire.ContentTypeUnspecified {
		t.Errorf("consumer action ContentType mutated to %q", action.ContentType)
	}
}

// TestFormWireCSRFTravel pins that the CSRF hidden input renders inside the
// wired form: under form encoding it serializes with the other fields in both
// dialects, so CSRF protection survives the transport switch.
func TestFormWireCSRFTravel(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Form(FormProps{
		CSRFToken: "tok-123",
		Wire: &wire.Action{
			Transport: wire.TransportDatastar,
			Method:    wire.MethodPost,
			URL:       "/api/save",
		},
	}))
	utils.AssertContains(t, output, `name="csrf_token"`)
	utils.AssertContains(t, output, `value="tok-123"`)
	utils.AssertContains(t, output, `data-on:submit=`)
}

// TestFormWireAttributesDefaults verifies the form-level action defaults in
// isolation: submit event and form encoding for unspecified values, nil and
// URL-less actions staying inert.
func TestFormWireAttributesDefaults(t *testing.T) {
	t.Parallel()

	if attrs := formWireAttributes(nil); attrs != nil {
		t.Fatalf("formWireAttributes(nil) = %v, want nil", attrs)
	}

	if attrs := formWireAttributes(&wire.Action{URL: ""}); attrs != nil {
		t.Fatalf("formWireAttributes(empty URL) = %v, want nil", attrs)
	}

	attrs := formWireAttributes(&wire.Action{
		Transport: wire.TransportDatastar,
		URL:       "/api/save",
	})
	if got := attrs["data-on:submit"]; got != "@post('/api/save', {contentType: 'form'})" {
		t.Fatalf("defaulted datastar expression = %v, want form-encoded submit", got)
	}
}
