package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestButtonWire verifies the transport-agnostic wiring field renders the
// right attribute dialect (htmx default, Datastar opt-in) and stays inert
// when unset.
func TestButtonWire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		props       ButtonProps
		contains    []string
		notContains []string
	}{
		{
			name: "nil wire renders no wiring attributes",
			props: ButtonProps{
				BaseProps: utils.BaseProps{ID: "btn-plain"},
				Text:      "Plain",
			},
			notContains: []string{"hx-get", "hx-post", "data-on:", "hx-target"},
		},
		{
			name: "htmx dialect (unspecified transport resolves to htmx)",
			props: ButtonProps{
				BaseProps: utils.BaseProps{ID: "btn-htmx"},
				Text:      "Load",
				Wire: &wire.Action{
					Method: wire.MethodGet,
					URL:    "/api/wire/fragment",
					Target: "#wire-htmx-out",
				},
			},
			contains: []string{`hx-get="/api/wire/fragment"`, `hx-target="#wire-htmx-out"`},
		},
		{
			name: "datastar dialect (target is response-driven, not emitted)",
			props: ButtonProps{
				BaseProps: utils.BaseProps{ID: "btn-datastar"},
				Text:      "Load",
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					URL:       "/api/wire/fragment",
					Target:    "#wire-datastar-out",
				},
			},
			contains:    []string{`data-on:click=`},
			notContains: []string{"hx-get", "hx-target"},
		},
		{
			name: "wire composes with Attrs (Attrs wins on conflict)",
			props: ButtonProps{
				BaseProps: utils.BaseProps{
					ID: "btn-both",
					Attrs: map[string]any{
						"hx-target": "#consumer-override",
					},
				},
				Text: "Both",
				Wire: &wire.Action{URL: "/api/wire/fragment"},
			},
			contains: []string{`hx-target="#consumer-override"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := utils.Render(t, Button(tt.props))

			for _, want := range tt.contains {
				utils.AssertContains(t, html, want)
			}

			for _, absent := range tt.notContains {
				utils.AssertNotContains(t, html, absent)
			}
		})
	}
}
