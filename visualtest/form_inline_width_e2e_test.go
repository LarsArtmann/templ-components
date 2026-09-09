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
			content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					if err := forms.Input(forms.InputProps{
						BaseProps: utils.BaseProps{ID: "in1"},
						Name:      "q",
						Label:     "Query",
					}).Render(ctx, w); err != nil {
						return err //nolint:wrapcheck // probe
					}

					return forms.Input(forms.InputProps{ //nolint:wrapcheck // probe
						BaseProps: utils.BaseProps{ID: "in2"},
						Name:      "n",
						Label:     "Number",
					}).Render(ctx, w)
				})

				inline := forms.Form(forms.FormProps{
					Action: "/x",
					Layout: forms.FormLayoutInline,
				})

				return inline.Render(templ.WithChildren(ctx, children), w) //nolint:wrapcheck // probe
			})

			props := layout.DefaultPageProps()
			props.Title = "probe"

			return layout.Base(props).Render(templ.WithChildren(ctx, content), w) //nolint:wrapcheck // probe
		})
	}

	server := e2ePageServer(t, page, nil)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	var report string

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`(() => {
			const form = document.querySelector('form.flex');
			if (!form) return 'NO INLINE FORM FOUND';
			const inputs = form.querySelectorAll('input[type="text"], input:not([type])');
			const out = ['form=' + form.offsetWidth + 'px'];
			inputs.forEach(function(input) {
				const wrap = input.closest('div');
				out.push('input=' + input.offsetWidth + 'px wrap=' + (wrap ? wrap.offsetWidth : '?') + 'px');
			});
			return out.join(' | ');
		})()`, &report,
	)); err != nil {
		t.Fatalf("measure: %v", err)
	}

	t.Logf("INLINE WIDTHS: %s", report)
}
