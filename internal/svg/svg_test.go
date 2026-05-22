// Package svg provides tests for shared SVG rendering primitives.
package svg

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestFillIconRender(t *testing.T) {
	t.Parallel()

	t.Run("renders SVG with correct path", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FillIcon(
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
		output := utils.Render(t, FillIcon("h-4 w-4", "M10 0", true))
		utils.AssertContains(t, output, "rotate-180")
	})

	t.Run("renders without rotation by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FillIcon("h-4 w-4", "M10 0", false))
		utils.AssertNotContains(t, output, "rotate-180")
	})
}

func TestSpinnerSVGRender(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SpinnerSVG())
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
		{"PathChevronDown", PathChevronDown},
		{"PathChevronSmall", PathChevronSmall},
		{"PathArrowUp", PathArrowUp},
		{"PathArrowDown", PathArrowDown},
		{"PathArrowLeft", PathArrowLeft},
		{"PathArrowRight", PathArrowRight},
		{"PathAvatarFill", PathAvatarFill},
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
	output := utils.Render(t, FillIcon("h-5 w-5", PathChevronDown, false))
	utils.AssertContains(t, output, PathChevronDown)
}
