// Package svg_test covers the shared SVG rendering primitives from outside
// the package. External placement is load-bearing: the parent utils package
// now imports utils/svg (DismissButton sources its X path from PathXMark),
// so an in-package test importing utils would form an import cycle.
package svg_test

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/svg"
)

func TestFillIconRender(t *testing.T) {
	t.Parallel()

	t.Run("renders SVG with correct path", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, svg.FillIcon(
			"h-5 w-5",
			"M10 0C4.48 0 0 4.48 0 10s4.48 10 10 10",
			false,
		))
		utils.AssertContains(t, output, `viewBox="0 0 20 20"`)
		utils.AssertContains(t, output, "M10 0C4.48")
		utils.AssertContains(t, output, `aria-hidden="true"`)
		utils.AssertContains(t, output, `class="h-5 w-5"`)
	})

	t.Run("renders with rotation when rotate is true", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, svg.FillIcon("h-4 w-4", "M10 0", true))
		utils.AssertContains(t, output, "rotate-180")
	})

	t.Run("renders without rotation by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, svg.FillIcon("h-4 w-4", "M10 0", false))
		utils.AssertNotContains(t, output, "rotate-180")
	})
}

func TestSpinnerSVGRender(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, svg.SpinnerSVG())
	utils.AssertContains(t, output, `opacity-25`)
	utils.AssertContains(t, output, `opacity-75`)
	utils.AssertContains(t, output, `stroke="currentColor"`)
	utils.AssertContains(t, output, `stroke-width="4"`)
}

func TestPathConstants(t *testing.T) {
	t.Parallel()

	paths := []struct {
		name string
		path string
	}{
		{"PathChevronDown", svg.PathChevronDown},
		{"PathChevronSmall", svg.PathChevronSmall},
		{"PathXMark", svg.PathXMark},
		{"PathArrowUp", svg.PathArrowUp},
		{"PathArrowDown", svg.PathArrowDown},
		{"PathArrowLeft", svg.PathArrowLeft},
		{"PathArrowRight", svg.PathArrowRight},
		{"PathAvatarFill", svg.PathAvatarFill},
	}
	for _, p := range paths {
		t.Run(p.name+" is non-empty", func(t *testing.T) {
			t.Parallel()

			if p.path == "" {
				t.Errorf("%s is empty", p.name)
			}
		})
	}
}

func TestFillIconUsesPathConstants(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, svg.FillIcon("h-5 w-5", svg.PathChevronDown, false))
	utils.AssertContains(t, output, svg.PathChevronDown)
}
