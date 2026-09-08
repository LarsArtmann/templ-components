package visualtest_test

import (
	"context"
	"io"
	"strconv"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
	"github.com/larsartmann/templ-components/visualtest"
)

// The forms pattern pack visual goldens: pixel-level regression coverage for
// the v1.14.0 wire cards' resting/error states in both color modes. The
// interaction (debounce, submit, swap, morphing) is proven by the browser
// e2e suite (wire_forms_pack_e2e_test.go); these goldens pin what the cards
// LOOK like so layout/color regressions surface independently.

// packWrap wraps a card in a fixed-width column so text wrapping and form
// widths render deterministically.
func packWrap(inner templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="w-[36rem] space-y-4 p-2">`); err != nil {
			return err
		}

		if err := inner.Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})
}

// packFilterSection renders the debounced filter input with a populated
// results region — the state after the endpoint's first swap.
func packFilterSection() templ.Component {
	return packWrap(templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := forms.FilterInput(forms.FilterInputProps{
			Name:        "q",
			Label:       "Filter frameworks",
			Placeholder: "Type to filter…",
			HelpText:    "Debounced 300ms — results update as you type",
			DebounceMS:  300,
			Wire: &wire.Action{
				URL:    "/api/wire/filter",
				Target: "#pack-vis-filter-out",
			},
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w,
			`<div id="pack-vis-filter-out" class="mt-3 divide-y divide-gray-200 dark:divide-gray-700 `+
				`text-sm text-gray-900 dark:text-white">`+
				`<div class="py-1.5">Datastar</div><div class="py-1.5">Htmx</div><div class="py-1.5">Templ</div>`+
				`</div>`)

		return err
	}))
}

// packDropdownSection renders the wired filter dropdown with its verdict.
func packDropdownSection() templ.Component {
	return packWrap(templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := forms.FilterDropdown(forms.FilterDropdownProps{
			Name:  "framework",
			Label: "Framework",
			Options: []forms.SelectOption{
				{Value: "alpha", Label: "Alpha"},
				{Value: "beta", Label: "Beta", Selected: true},
				{Value: "gamma", Label: "Gamma"},
			},
			Wire: &wire.Action{
				URL:    "/api/wire/dropdown",
				Target: "#pack-vis-dropdown-out",
			},
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w,
			`<div id="pack-vis-dropdown-out" class="mt-3">`)
		if err != nil {
			return err
		}

		if err := feedback.InlineSuccess("Picked beta via htmx.").Render(ctx, w); err != nil {
			return err
		}

		_, werr := io.WriteString(w, `</div>`)

		return werr
	}))
}

// packWizardSection renders the wizard's profile step in its error state —
// the StepIndicator mid-journey plus the inline field error.
func packWizardSection() templ.Component {
	return packWrap(templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := feedback.StepIndicator(feedback.StepIndicatorProps{
			Steps:       []string{"Account", "Profile", "Done"},
			CurrentStep: 1,
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `<div class="mt-4 space-y-3">`)
		if err != nil {
			return err
		}

		if err := forms.FieldError("pack-vis-wizard-name", "Enter your name.").Render(ctx, w); err != nil {
			return err
		}

		if err := forms.Input(forms.InputProps{
			Name:  "name",
			Label: "Full name",
		}).Render(ctx, w); err != nil {
			return err
		}

		return display.Button(display.ButtonProps{
			Text:    "Next",
			Variant: display.ButtonPrimary,
			Size:    display.ButtonSizeSM,
			Type:    display.ButtonHTMLSubmit,
		}).Render(ctx, w)
	}))
}

// packUploadSection renders the wired multipart upload form with its
// success verdict.
func packUploadSection() templ.Component {
	return packWrap(templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		uploadChildren := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := forms.FileInput(forms.FileInputProps{Name: "attachment", Label: "Attachment"}).
				Render(ctx, w); err != nil {
				return err
			}

			return display.Button(display.ButtonProps{
				Text:    "Upload",
				Variant: display.ButtonPrimary,
				Size:    display.ButtonSizeSM,
				Type:    display.ButtonHTMLSubmit,
			}).Render(ctx, w)
		})

		if err := forms.Form(forms.FormProps{
			Action:  "/api/wire/upload",
			Method:  forms.FormPost,
			Enctype: forms.FormEnctypeMultipart,
			BaseProps: utils.BaseProps{
				Class: "space-y-3",
			},
			Wire: &wire.Action{
				URL:    "/api/wire/upload",
				Target: "#pack-vis-upload-out",
			},
		}).Render(templ.WithChildren(ctx, uploadChildren), w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `<div id="pack-vis-upload-out" class="mt-3">`)
		if err != nil {
			return err
		}

		if err := feedback.InlineSuccess(
			"Uploaded e2e-upload.txt ("+strconv.Itoa(2048)+" bytes) via htmx.",
		).Render(ctx, w); err != nil {
			return err
		}

		_, werr := io.WriteString(w, `</div>`)

		return werr
	}))
}

// packSearchSection renders the wired GET search form with its verdict.
func packSearchSection() templ.Component {
	return packWrap(templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := feedback.InlineSuccess("GET received q=“ada” via htmx.").Render(ctx, w); err != nil {
			return err
		}

		searchChildren := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := forms.Input(forms.InputProps{
				Name:        "q",
				Type:        forms.InputSearch,
				Label:       "Search",
				Value:       "ada",
				Placeholder: "Try: ada",
			}).Render(ctx, w); err != nil {
				return err
			}

			return display.Button(display.ButtonProps{
				Text:    "Search",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Type:    display.ButtonHTMLSubmit,
			}).Render(ctx, w)
		})

		if err := forms.Form(forms.FormProps{
			Action: "/api/wire/search",
			Method: forms.FormGet,
			BaseProps: utils.BaseProps{
				Class: "flex flex-wrap items-end gap-3",
			},
			Wire: &wire.Action{
				URL:    "/api/wire/search",
				Target: "#pack-vis-search-out",
			},
		}).Render(templ.WithChildren(ctx, searchChildren), w); err != nil {
			return err
		}

		_, werr := io.WriteString(w, `</div>`)

		return werr
	}))
}

// packDirtySection renders the guarded subscribe form with its verdict —
// the DirtyGuard lifecycle's resting state.
func packDirtySection() templ.Component {
	return packWrap(templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := feedback.InlineSuccess("Saved Graf Zeppelin via htmx.").Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `<div class="mt-3">`)
		if err != nil {
			return err
		}

		children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := forms.Input(forms.InputProps{
				Name:  "project",
				Label: "Project name",
				Value: "Graf Zeppelin",
			}).Render(ctx, w); err != nil {
				return err
			}

			return display.Button(display.ButtonProps{
				Text:    "Save",
				Type:    display.ButtonHTMLSubmit,
				Size:    display.ButtonSizeSM,
				Variant: display.ButtonPrimary,
			}).Render(ctx, w)
		})

		if err := forms.Form(forms.FormProps{
			Method:     forms.FormPost,
			DirtyGuard: true,
			BaseProps:  utils.BaseProps{Class: "space-y-3"},
			Wire: &wire.Action{
				URL:    "/api/wire/dirty",
				Target: "#pack-vis-dirty-out",
			},
		}).Render(templ.WithChildren(ctx, children), w); err != nil {
			return err
		}

		_, werr := io.WriteString(w, `</div>`)

		return werr
	}))
}

// packAssertModes pins a card component in light and dark mode.
func packAssertModes(t *testing.T, name string, c templ.Component) {
	t.Helper()

	visualtest.AssertScreenshot(t, "wire/pack_"+name+"_light", c)
	visualtest.AssertScreenshot(t, "wire/pack_"+name+"_dark", c, visualtest.Options{Dark: new(true)})
}

func TestWirePackFilterSection(t *testing.T) {
	t.Parallel()
	packAssertModes(t, "filter", packFilterSection())
}

func TestWirePackDropdownSection(t *testing.T) {
	t.Parallel()
	packAssertModes(t, "dropdown", packDropdownSection())
}

func TestWirePackWizardSection(t *testing.T) {
	t.Parallel()
	packAssertModes(t, "wizard", packWizardSection())
}

func TestWirePackUploadSection(t *testing.T) {
	t.Parallel()
	packAssertModes(t, "upload", packUploadSection())
}

func TestWirePackSearchSection(t *testing.T) {
	t.Parallel()
	packAssertModes(t, "search", packSearchSection())
}

func TestWirePackDirtySection(t *testing.T) {
	t.Parallel()
	packAssertModes(t, "dirty", packDirtySection())
}
