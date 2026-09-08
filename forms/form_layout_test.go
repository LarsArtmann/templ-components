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

	// Count by a single stable token (utils.Class/tailwind-merge reorders
	// class tokens — never assert ordered substrings, AGENTS.md convention).
	groups := strings.Count(output, "basis-56")
	if groups != 2 {
		t.Errorf("Inline form must group each field into one flex item; found %d group divs, want 2\noutput:\n%s", groups, output)
	}

	utils.AssertContainsAll(t, output, "min-w-40", "flex-1", "basis-56")

	// Structural: each group div must open immediately before its label, so
	// label + control live inside ONE flex item (not scattered siblings).
	if n := strings.Count(output, `basis-56"><label`); n != 2 {
		t.Errorf("group div must wrap the label; found %d label-adjacent groups, want 2", n)
	}

	labels := strings.Count(output, "<label")
	inputs := strings.Count(output, "<input")
	if labels != 2 || inputs != 2 {
		t.Errorf("two labeled inputs expected; got %d labels, %d inputs", labels, inputs)
	}

	stacked := utils.Render(t, Form(FormProps{Action: "/x"}))
	if strings.Contains(stacked, "basis-56") {
		t.Error("Stack layout must NOT emit the inline grouping div")
	}
}
