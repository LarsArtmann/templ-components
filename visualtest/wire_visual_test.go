package visualtest_test

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
	"github.com/larsartmann/templ-components/visualtest"
)

// wireSectionComponent renders the demo's wire section — the same wired Button
// pair and output regions the E2E suite clicks — for pixel-level regression
// coverage of the transport switch UI.
func wireSectionComponent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="grid grid-cols-2 gap-4 w-[40rem]">`); err != nil {
			return err
		}

		htmxButton := display.ButtonProps{
			BaseProps: utils.BaseProps{ID: "btn-wire-htmx"},
			Text:      "Load via htmx",
			Variant:   display.ButtonSecondary,
			Size:      display.ButtonSizeSM,
			Wire: &wire.Action{
				URL:    "/api/wire/fragment",
				Target: "#wire-htmx-out",
			},
		}
		if err := display.Button(htmxButton).Render(ctx, w); err != nil {
			return err
		}

		datastarButton := display.ButtonProps{
			BaseProps: utils.BaseProps{ID: "btn-wire-datastar"},
			Text:      "Load via Datastar",
			Variant:   display.ButtonSecondary,
			Size:      display.ButtonSizeSM,
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/wire/fragment",
			},
		}
		if err := display.Button(datastarButton).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w,
			`<div id="wire-htmx-out" class="col-span-1 text-sm text-gray-500 dark:text-gray-400">htmx target</div>`+
				`<div id="wire-datastar-out" class="col-span-1 text-sm text-gray-500 dark:text-gray-400">datastar target</div>`+
				`</div>`)

		return err
	})
}

// wireFormSectionComponent renders the dual-transport form section in its
// error round-trip state — the re-rendered fragment the validation e2e
// proves functionally. Pins the summary + inline-error + preserved-value
// visual presentation in both modes.
func wireFormSectionComponent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="w-[40rem] space-y-6">`); err != nil {
			return err
		}

		errorState := forms.InputProps{
			Name:  "email",
			Type:  forms.InputEmail,
			Label: "Email address",
			Value: "ada@example",
			Error: "Enter an email address with a domain.",
		}
		summary := forms.ValidationSummary(forms.ValidationSummaryProps{
			Errors: []forms.ValidationError{
				{Field: "email", Message: "Enter an email address with a domain."},
			},
		})

		for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
			if err := summary.Render(ctx, w); err != nil {
				return err
			}

			formProps := forms.FormProps{
				Method: forms.FormPost,
				Wire: &wire.Action{
					Transport: dialect,
					Method:    wire.MethodPost,
					URL:       "/api/wire/form",
				},
			}
			if dialect == wire.TransportHTMX {
				formProps.Wire.Target = "#wire-form-htmx-region"
			}

			children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if err := forms.Input(errorState).Render(ctx, w); err != nil {
					return err
				}

				return display.Button(display.ButtonProps{
					Text:    "Subscribe",
					Type:    display.ButtonHTMLSubmit,
					Size:    display.ButtonSizeSM,
					Variant: display.ButtonPrimary,
				}).Render(ctx, w)
			})

			if err := forms.Form(formProps).Render(templ.WithChildren(ctx, children), w); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})
}

// TestWireSection pins the wired transport-switch section visually in both
// modes; interaction is covered by the browser E2E suite (wire_e2e_test.go).
func TestWireSection(t *testing.T) {
	t.Parallel()

	visualtest.AssertScreenshot(t, "wire/dual_transport_light", wireSectionComponent())
	visualtest.AssertScreenshot(
		t,
		"wire/dual_transport_dark",
		wireSectionComponent(),
		visualtest.Options{Dark: new(true)},
	)
}

// TestWireFormSection pins the dual-transport form's error round-trip state
// (summary + inline error + preserved value + submit button) in both modes;
// the full fill/submit/resubmit behavior is proven by wire_form_e2e_test.go.
func TestWireFormSection(t *testing.T) {
	t.Parallel()

	visualtest.AssertScreenshot(t, "wire/form_roundtrip_light", wireFormSectionComponent())
	visualtest.AssertScreenshot(
		t,
		"wire/form_roundtrip_dark",
		wireFormSectionComponent(),
		visualtest.Options{Dark: new(true)},
	)
}
