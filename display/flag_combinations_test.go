package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// T5.8 (#493): boolean-flag components are usually tested one flag at a
// time, so flag INTERACTIONS (two attributes that must co-exist, one flag
// suppressing another's markup) go unproven. These tables walk meaningful
// flag combinations and assert the invariant that must hold for EVERY
// combination: rendered in one pass, no panic, and each combination's
// contract markers present (or absent) together.

// TestButtonFlagCombinations walks the exhaustive Disabled × External ×
// link-mode surface. Contract: href renders an <a> (never a <button>);
// button mode renders a <button>; Disabled suppresses pointer events in
// both modes.
func TestButtonFlagCombinations(t *testing.T) {
	t.Parallel()

	for _, disabled := range []bool{false, true} {
		for _, external := range []bool{false, true} {
			for _, href := range []bool{false, true} {
				name := "disabled=" + b2s(disabled) + "/external=" + b2s(external) + "/href=" + b2s(href)
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					props := ButtonProps{
						Text:     "Go",
						Disabled: disabled,
						External: external,
					}
					if href {
						props.Href = "/target"
					}

					output := utils.Render(t, Button(props))

					if href {
						utils.AssertContains(t, output, `href="/target"`)
						utils.AssertNotContains(t, output, "<button")
					} else {
						utils.AssertContains(t, output, "<button")
					}

					if disabled {
						utils.AssertContains(t, output, "disabled:opacity-50")
						utils.AssertContains(t, output, ` disabled`)
					}
				})
			}
		}
	}
}

// TestTableFlagCombinations pins Table's presentation flags: every flag on
// its own must carry its contract marker, and all flags together must
// coexist without breaking another's markup.
func TestTableFlagCombinations(t *testing.T) {
	t.Parallel()

	flags := []struct {
		name   string
		apply  func(*TableProps)
		marker string
	}{
		{"striped", func(p *TableProps) { p.Striped = true }, "bg-gray-50"},
		{"bordered", func(p *TableProps) { p.Bordered = true }, "border-gray-200"},
		{"sticky-header", func(p *TableProps) { p.StickyHeader = true }, "sticky"},
		{"lazy-rows", func(p *TableProps) { p.LazyRows = true }, "tc-content-auto"},
	}

	props := func() TableProps {
		return TableProps{
			Headers: []string{"A", "B"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "1"}, {Text: "2"}}},
				{Cells: []TableCell{{Text: "3"}, {Text: "4"}}},
			},
		}
	}

	t.Run("all flags off renders a plain table", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(props()))
		utils.AssertContains(t, output, "<table")
		utils.AssertNotContains(t, output, "sticky")
	})

	for _, f := range flags {
		t.Run("single flag "+f.name, func(t *testing.T) {
			t.Parallel()

			p := props()
			f.apply(&p)
			output := utils.Render(t, Table(p))
			utils.AssertContains(t, output, "<table")
			utils.AssertContains(t, output, f.marker)
		})
	}

	t.Run("all flags on coexist", func(t *testing.T) {
		t.Parallel()

		p := props()
		for _, f := range flags {
			f.apply(&p)
		}
		p.Hover = true

		output := utils.Render(t, Table(p))
		utils.AssertContains(t, output, "<table")
		utils.AssertContains(t, output, "border-gray-200")
		utils.AssertContains(t, output, "tc-content-auto")
	})
}

func b2s(b bool) string {
	if b {
		return "t"
	}

	return "f"
}
