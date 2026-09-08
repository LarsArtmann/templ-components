package forms

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestFormLayout(t *testing.T) {
	t.Parallel()

	t.Run("default layout is Stack (space-y-6)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x"}))
		utils.AssertContains(t, output, "space-y-6")
	})

	t.Run("Layout: Stack emits space-y-6", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x", Layout: FormLayoutStack}))
		utils.AssertContains(t, output, "space-y-6")
	})

	t.Run("Layout: Inline emits flex flex-wrap items-end gap-3", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x", Layout: FormLayoutInline}))
		utils.AssertContainsAll(t, output, "flex", "flex-wrap", "items-end", "gap-3")
	})

	t.Run("Layout: Grid emits aligned grid with minmax(0,1fr) blowout guard", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x", Layout: FormLayoutGrid}))
		utils.AssertContainsAll(
			t, output,
			"grid", "grid-cols-1", "sm:grid-cols-[auto_minmax(0,1fr)]", "items-start", "gap-x-4", "gap-y-3",
		)
	})

	t.Run("unknown Layout falls back to Stack", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x", Layout: FormLayout("bogus")}))
		utils.AssertContains(t, output, "space-y-6")
	})

	t.Run("Layout=Inline explicit overrides default Stack", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Form(FormProps{Action: "/x", Layout: FormLayoutInline}))

		if strings.Contains(output, "space-y-6") {
			t.Errorf("explicit Inline must not emit space-y-6")
		}
	})
}

func TestFormLayoutIsValid(t *testing.T) {
	t.Parallel()

	for _, l := range []FormLayout{FormLayoutStack, FormLayoutInline, FormLayoutGrid} {
		if !FormLayoutIsValid(l) {
			t.Errorf("FormLayoutIsValid(%q) = false; want true", l)
		}
	}

	if FormLayoutIsValid(FormLayout("bogus")) {
		t.Errorf("FormLayoutIsValid(\"bogus\") = true; want false")
	}
}

// TestFormLayoutInlineFieldGrouping pins the #166 width contract: inside an
// Inline form, each field's label + control + error + help must be grouped
// into ONE flex item. Before the fix, they were separate flex children and
// the w-full controls forced one field per row — the Inline layout silently
// degenerated to a stack (demo audit f25).
func TestFormLayoutInlineFieldGrouping(t *testing.T) {
	t.Parallel()

	page := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return Form(FormProps{Action: "/filter", Method: FormGet, Layout: FormLayoutInline}).
			Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				for _, name := range []string{"status", "q"} {
					input := DefaultInputProps()
					input.Name = name
					input.Label = name

					if err := Input(input).Render(ctx, w); err != nil {
						return err
					}
				}

				return nil
			})), w)
	})

	output := utils.Render(t, page)

	groups := strings.Count(output, `<div class="basis-56 flex-1 min-w-40">`)
	if groups != 2 {
		t.Errorf("Inline form must group each field into one flex item; found %d group divs, want 2\noutput:\n%s", groups, output)
	}

	labels := strings.Count(output, "<label")
	inputs := strings.Count(output, "<input")
	if labels != 2 || inputs != 2 {
		t.Errorf("two labeled inputs expected; got %d labels, %d inputs", labels, inputs)
	}

	// Structural check: the label must sit INSIDE the group div (group div
	// opens before each label, no group close between them). The exact
	// ordering assertion: group div count == label count == input count means
	// nothing was scattered between groups.
	stacked := utils.Render(t, Form(FormProps{Action: "/x"}))
	if strings.Contains(stacked, "min-w-40") {
		t.Error("Stack layout must NOT emit the inline grouping div")
	}
}
