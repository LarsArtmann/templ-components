package visualtest_test

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
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
