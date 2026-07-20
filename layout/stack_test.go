package layout

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestStack(t *testing.T) {
	t.Parallel()

	t.Run("default props produce flex flex-col + space-y-4", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stack(DefaultStackProps()))
		utils.AssertContainsAll(t, output, "flex", "flex-col", "space-y-4")
	})

	t.Run("gap SM", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stack(StackProps{Gap: StackGapSM}))
		utils.AssertContains(t, output, "space-y-2")
	})

	t.Run("gap LG", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stack(StackProps{Gap: StackGapLG}))
		utils.AssertContains(t, output, "space-y-6")
	})

	t.Run("gap XL", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stack(StackProps{Gap: StackGapXL}))
		utils.AssertContains(t, output, "space-y-8")
	})

	t.Run("gap None omits space-y class", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Stack(StackProps{Gap: StackGapNone}))
		if strings.Contains(output, "space-y-") {
			t.Errorf("StackGapNone should omit space-y-*; output = %q", output)
		}

		utils.AssertContainsAll(t, output, "flex", "flex-col")
	})

	t.Run("unknown gap falls back to default (MD)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stack(StackProps{Gap: StackGap("bogus")}))
		utils.AssertContains(t, output, "space-y-4")
	})

	t.Run("BaseProps propagate (Class, AriaLabel)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stack(StackProps{
			BaseProps: utils.BaseProps{
				Class:     "data-tc-test",
				AriaLabel: "Card stack",
			},
		}))
		utils.AssertContainsAll(t, output, "data-tc-test", `aria-label="Card stack"`)
	})
}

func TestStackGapIsValid(t *testing.T) {
	t.Parallel()

	for _, g := range []StackGap{StackGapNone, StackGapSM, StackGapMD, StackGapLG, StackGapXL} {
		if !StackGapIsValid(g) {
			t.Errorf("StackGapIsValid(%q) = false; want true", g)
		}
	}

	if StackGapIsValid(StackGap("bogus")) {
		t.Errorf("StackGapIsValid(\"bogus\") = true; want false")
	}
}
