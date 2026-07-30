package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestCarouselKeyboardNavigation(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Carousel(CarouselProps{
		BaseProps:      utils.BaseProps{Nonce: "n"},
		ShowArrows:     true,
		ShowIndicators: true,
		Slides: []CarouselSlide{
			{Content: templ.Raw("<div>Slide 1</div>")},
			{Content: templ.Raw("<div>Slide 2</div>")},
			{Content: templ.Raw("<div>Slide 3</div>")},
		},
	}))

	// The carousel region must be focusable so keyboard users can navigate it.
	utils.AssertContains(t, output, `tabindex="0"`)

	// Arrow keys, Home, and End control slide movement.
	utils.AssertContains(t, output, "'ArrowLeft'")
	utils.AssertContains(t, output, "'ArrowRight'")
	utils.AssertContains(t, output, "'Home'")
	utils.AssertContains(t, output, "'End'")

	// RTL-aware arrow-key mapping mirrors the spatial direction of the arrows.
	utils.AssertContains(t, output, "document.documentElement.getAttribute('dir')==='rtl'")
}

func TestCarouselFocusVisible(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Carousel(CarouselProps{
		ShowArrows: true,
		Slides: []CarouselSlide{
			{Content: templ.Raw("<div>Slide 1</div>")},
		},
	}))

	// The focusable carousel region (tabindex="0") must show a visible focus
	// ring so keyboard users can see where focus landed before navigating.
	utils.AssertContains(t, output, "focus-visible:ring-2")
	utils.AssertContains(t, output, "focus:outline-none")
}
