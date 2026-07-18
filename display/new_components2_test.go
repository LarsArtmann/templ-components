package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// --- HoverCard ---

func TestDefaultHoverCardProps(t *testing.T) {
	t.Parallel()

	p := DefaultHoverCardProps()
	if p.Position != HoverCardPositionBottom {
		t.Error("expected default position Bottom")
	}
}

func TestHoverCardPositionIsValid(t *testing.T) {
	t.Parallel()

	for _, pos := range []HoverCardPosition{HoverCardPositionTop, HoverCardPositionBottom, HoverCardPositionStart, HoverCardPositionEnd} {
		if !HoverCardPositionIsValid(pos) {
			t.Errorf("expected %q to be valid", pos)
		}
	}

	if HoverCardPositionIsValid(HoverCardPosition("invalid")) {
		t.Error("expected invalid to be invalid")
	}
}

func TestHoverCardBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, HoverCard(HoverCardProps{
		BaseProps: utils.BaseProps{ID: "hc"},
		Content:   templ.Raw("<p>Card content</p>"),
	}))
	utils.AssertContains(t, output, "Card content")
	utils.AssertContains(t, output, `role="tooltip"`)
	utils.AssertContains(t, output, `aria-describedby="hc"`)
}

func TestHoverCardGolden(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, HoverCard(HoverCardProps{
		BaseProps: utils.BaseProps{ID: "hcg"},
		Position:  HoverCardPositionTop,
		Content:   templ.Raw("<p>Tooltip text</p>"),
	}))
	golden.Assert(t, "hover_card_basic", output)
}

func TestHoverCardA11y(t *testing.T) {
	t.Parallel()

	t.Run("trigger is focusable", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, HoverCard(DefaultHoverCardProps()))
		utils.AssertContains(t, output, `tabindex="0"`)
	})

	t.Run("card has role tooltip", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, HoverCard(HoverCardProps{Content: templ.Raw("x")}))
		utils.AssertContains(t, output, `role="tooltip"`)
	})

	t.Run("dark mode classes present", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, HoverCard(HoverCardProps{Content: templ.Raw("x")}))
		utils.AssertContains(t, output, "dark:bg-gray-800")
		utils.AssertContains(t, output, "dark:border-gray-700")
	})
}

func TestHoverCardEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("invalid position falls back to bottom", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, HoverCard(HoverCardProps{
			BaseProps: utils.BaseProps{ID: "inv"},
			Position:  HoverCardPosition("bad"),
			Content:   templ.Raw("x"),
		}))
		utils.AssertContains(t, output, "top-full")
	})

	t.Run("nil content renders empty card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, HoverCard(HoverCardProps{}))
		utils.AssertContains(t, output, `role="tooltip"`)
	})

	t.Run("custom class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, HoverCard(HoverCardProps{
			BaseProps: utils.BaseProps{Class: "my-hc"},
		}))
		utils.AssertContains(t, output, "my-hc")
	})
}

// --- ContextMenu ---

func TestDefaultContextMenuProps(t *testing.T) {
	t.Parallel()

	_ = DefaultContextMenuProps()
}

func TestContextMenuBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ContextMenu(ContextMenuProps{
		BaseProps: utils.BaseProps{ID: "ctx", Nonce: "n"},
		Items: []ContextMenuItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Delete", Href: "/delete"},
		},
	}))
	utils.AssertContains(t, output, "Edit")
	utils.AssertContains(t, output, "Delete")
	utils.AssertContains(t, output, `role="menu"`)
	utils.AssertContains(t, output, `data-tc-ctxmenu-trigger`)
}

func TestContextMenuA11y(t *testing.T) {
	t.Parallel()

	t.Run("items have role menuitem", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ContextMenu(ContextMenuProps{
			BaseProps: utils.BaseProps{ID: "a11y", Nonce: "n"},
			Items:     []ContextMenuItem{{Text: "Edit", Href: "/edit"}},
		}))
		utils.AssertContains(t, output, `role="menuitem"`)
	})

	t.Run("disabled items have pointer-events-none", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ContextMenu(ContextMenuProps{
			BaseProps: utils.BaseProps{ID: "dis", Nonce: "n"},
			Items:     []ContextMenuItem{{Text: "Edit", Disabled: true}},
		}))
		utils.AssertContains(t, output, "pointer-events-none")
	})
}

func TestContextMenuEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("no nonce omits script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ContextMenu(ContextMenuProps{
			BaseProps: utils.BaseProps{ID: "nononce"},
			Items:     []ContextMenuItem{{Text: "X", Href: "/x"}},
		}))
		utils.AssertNotContains(t, output, "tcCtxMenuAttached")
	})

	t.Run("item without href renders span", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ContextMenu(ContextMenuProps{
			BaseProps: utils.BaseProps{ID: "span", Nonce: "n"},
			Items:     []ContextMenuItem{{Text: "Label only"}},
		}))
		utils.AssertContains(t, output, "Label only")
	})
}

// --- Carousel ---

func TestDefaultCarouselProps(t *testing.T) {
	t.Parallel()

	p := DefaultCarouselProps()
	if !p.ShowArrows {
		t.Error("expected ShowArrows=true by default")
	}
}

func TestCarouselBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Carousel(CarouselProps{
		Slides: []CarouselSlide{
			{Content: templ.Raw("<div>Slide 1</div>")},
			{Content: templ.Raw("<div>Slide 2</div>")},
		},
		ShowArrows:     true,
		ShowIndicators: true,
	}))
	utils.AssertContains(t, output, "Slide 1")
	utils.AssertContains(t, output, "Slide 2")
	utils.AssertContains(t, output, `role="region"`)
	utils.AssertContains(t, output, `aria-roledescription="carousel"`)
}

func TestCarouselA11y(t *testing.T) {
	t.Parallel()

	t.Run("slides have aria-roledescription", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Carousel(CarouselProps{
			Slides: []CarouselSlide{{Content: templ.Raw("x")}},
		}))
		utils.AssertContains(t, output, `aria-roledescription="slide"`)
	})

	t.Run("arrows have aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Carousel(CarouselProps{
			Slides: []CarouselSlide{
				{Content: templ.Raw("a")},
				{Content: templ.Raw("b")},
			},
			ShowArrows: true,
		}))
		utils.AssertContains(t, output, `aria-label="Previous slide"`)
		utils.AssertContains(t, output, `aria-label="Next slide"`)
	})

	t.Run("single slide hides arrows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Carousel(CarouselProps{
			Slides:     []CarouselSlide{{Content: templ.Raw("only")}},
			ShowArrows: true,
		}))
		utils.AssertNotContains(t, output, `aria-label="Previous slide"`)
	})
}

func TestCarouselEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty slides renders empty carousel", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Carousel(CarouselProps{}))
		utils.AssertContains(t, output, `data-tc-carousel`)
	})

	t.Run("nonce renders script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Carousel(CarouselProps{
			BaseProps: utils.BaseProps{Nonce: "n"},
			Slides: []CarouselSlide{
				{Content: templ.Raw("a")},
				{Content: templ.Raw("b")},
			},
		}))
		utils.AssertContains(t, output, "tcCarouselAttached")
	})

	t.Run("no nonce omits script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Carousel(CarouselProps{
			Slides: []CarouselSlide{
				{Content: templ.Raw("a")},
				{Content: templ.Raw("b")},
			},
		}))
		utils.AssertNotContains(t, output, "tcCarouselAttached")
	})
}

func ExampleHoverCard() {
	_ = HoverCard(HoverCardProps{
		Position: HoverCardPositionTop,
		Content:  templ.Raw("<p>Info</p>"),
	})
	// Output:
}

func ExampleContextMenu() {
	_ = ContextMenu(ContextMenuProps{
		Items: []ContextMenuItem{
			{Text: "Edit", Href: "/edit"},
		},
	})
	// Output:
}

func ExampleCarousel() {
	_ = Carousel(CarouselProps{
		Slides: []CarouselSlide{
			{Content: templ.Raw("<div>Slide 1</div>")},
		},
		ShowArrows: true,
	})
	// Output:
}
