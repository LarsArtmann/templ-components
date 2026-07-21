package forms

import (
	"strings"
	"testing"

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
		utils.AssertContainsAll(t, output,
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
