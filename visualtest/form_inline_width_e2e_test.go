package visualtest

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
)

// TestFormLayoutInlineWidthContract pins the inline-layout width contract in
// a real browser: fields inside a FormLayoutInline form must size to their
// content (side-by-side filter-bar behaviour), never stretch to the form's
// full width. The inputs carry `block w-full`, which browsers resolve to the
// intrinsic input width inside shrink-to-fit flex items — a future wrapper
// gaining `w-full` would silently break every inline filter bar, so this
// guard fails loudly if that regression lands.
func TestFormLayoutInlineWidthContract(t *testing.T) {
	page := func() templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if err := forms.Form(forms.FormProps{
					Action: "/x",
					Layout: forms.FormLayoutInline,
				}).Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					if err := forms.Input(forms.InputProps{
						BaseProps: utils.BaseProps{ID: "in1"},
						Name:      "q",
						Label:     "Query",
					}).Render(ctx, w); err != nil {
						return err
					}

					return forms.Input(forms.InputProps{
						BaseProps: utils.BaseProps{ID: "in2"},
						Name:      "n",
						Label:     "Number",
					}).Render(ctx, w)
				})), w); err != nil {
					return err
				}

				return nil
			})

			props := layout.DefaultPageProps()
			props.Title = "Inline width contract"

			return layout.Base(props).Render(templ.WithChildren(ctx, children), w)
		})
	}

	server := e2ePageServer(t, page, nil)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[forms]: navigate: %v", err)
	}

	var widths struct {
		Form  int64   `json:"form"`
		Count int64   `json:"count"`
		Input int64   `json:"input"`
		TopA  float64 `json:"topA"`
		TopB  float64 `json:"topB"`
	}

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`(() => {
			const form = document.querySelector('form.flex');
			if (!form) return null;
			const inputs = form.querySelectorAll('input');
			return {
				form: form.offsetWidth,
				count: inputs.length,
				input: inputs.length ? inputs[0].offsetWidth : 0,
				topA: inputs.length ? inputs[0].getBoundingClientRect().top : -1,
				topB: inputs.length > 1 ? inputs[1].getBoundingClientRect().top : -1,
			};
		})()`, &widths,
	)); err != nil {
		t.Fatalf("visualtest[forms]: measure inline widths: %v", err)
	}

	if widths.Form == 0 || widths.Count != 2 {
		t.Fatalf("visualtest[forms]: inline form not rendered as expected: %+v", widths)
	}

	// Contract 1: neither input stretches to (half of) the form width —
	// fields are content-sized, not full-width.
	if widths.Input*2 >= widths.Form {
		t.Errorf(
			"visualtest[forms]: input width %dpx means fields render (near) full-width in a %dpx form — the inline width contract is broken",
			widths.Input,
			widths.Form,
		)
	}

	// Contract 2: both fields sit on the SAME row (inline, not stacked).
	if widths.TopA != widths.TopB {
		t.Errorf(
			"visualtest[forms]: inputs wrapped onto separate rows (tops %v vs %v) — inline layout no longer inline",
			widths.TopA,
			widths.TopB,
		)
	}
}
