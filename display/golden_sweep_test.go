package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// Golden sweep for display components that previously lacked golden tests.
// Uses golden.AssertSnapshots for table-driven coverage with zero boilerplate.
// ID normalization (golden.normalizeIDs) makes EnsureID-using components
// (Accordion, Tabs, Dropdown, Tooltip, Carousel, ContextMenu) safe for
// golden testing without explicit ID props.

func TestGoldenSweepAccordion(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "accordion_default", HTML: utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{
					ID:      "section1",
					Title:   "Getting Started",
					Content: templ.Raw("<p>Welcome to the docs.</p>"),
					Open:    true,
				},
				{ID: "section2", Title: "Installation", Content: templ.Raw("<p>Run <code>go get</code>.</p>")},
				{ID: "section3", Title: "Configuration", Content: templ.Raw("<p>Edit <code>config.yaml</code>.</p>")},
			},
		}))},
		{Name: "accordion_empty", HTML: utils.Render(t, Accordion(DefaultAccordionProps()))},
	})
}

func TestGoldenSweepTabs(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "tabs_default", HTML: utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{ID: "overview", Label: "Overview", Content: templ.Raw("<p>Overview content</p>")},
				{ID: "activity", Label: "Activity", Content: templ.Raw("<p>Activity feed</p>")},
				{ID: "settings", Label: "Settings", Content: templ.Raw("<p>Settings panel</p>")},
			},
			ActiveTabID: "overview",
		}))},
		{Name: "tabs_pills", HTML: utils.Render(t, Tabs(TabsProps{
			Variant: TabsPills,
			Tabs: []Tab{
				{ID: "tab1", Label: "First", Content: templ.Raw("<p>First</p>")},
				{ID: "tab2", Label: "Second", Content: templ.Raw("<p>Second</p>")},
			},
			ActiveTabID: "tab1",
		}))},
	})
}

func TestGoldenSweepDropdown(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "dropdown_basic", HTML: utils.Render(t, Dropdown(DropdownProps{
			Label: "Actions",
			Items: []DropdownItem{
				{Text: "Edit", Href: "/edit", Icon: icons.Edit},
				{Text: "Duplicate", Href: "/duplicate", Icon: icons.DocumentDuplicate},
				{Text: "Delete", Href: "/delete", Icon: icons.Trash},
			},
		}))},
	})
}

func TestGoldenSweepTooltip(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "tooltip_top", HTML: utils.Render(t, Tooltip(TooltipProps{
			Text: "Edit this item",
		}))},
		{Name: "tooltip_bottom", HTML: utils.Render(t, Tooltip(TooltipProps{
			Text:     "More details here",
			Position: TooltipPositionBottom,
		}))},
	})
}

func TestGoldenSweepCarousel(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "carousel_basic", HTML: utils.Render(t, Carousel(CarouselProps{
			ShowIndicators: true,
			Slides: []CarouselSlide{
				{Content: templ.Raw(`<div class="p-8 text-center"><h3>Slide 1</h3></div>`)},
				{Content: templ.Raw(`<div class="p-8 text-center"><h3>Slide 2</h3></div>`)},
				{Content: templ.Raw(`<div class="p-8 text-center"><h3>Slide 3</h3></div>`)},
			},
		}))},
	})
}

func TestGoldenSweepContextMenu(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "contextmenu_basic", HTML: utils.Render(t, ContextMenu(ContextMenuProps{
			Items: []ContextMenuItem{
				{Text: "Copy", Href: "#copy"},
				{Text: "Share", Href: "#share"},
				{Text: "Delete", Href: "#delete"},
			},
		}))},
	})
}

func TestGoldenSweepAvatar(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "avatar_image", HTML: utils.Render(t, Avatar(AvatarProps{
			Src:      "/img/jane.jpg",
			Alt:      "Jane Doe",
			Initials: "JD",
			Size:     AvatarSizeMD,
			Status:   AvatarStatusOnline,
		}))},
		{Name: "avatar_initials", HTML: utils.Render(t, Avatar(AvatarProps{
			Alt:      "Bob Smith",
			Initials: "BS",
			Size:     AvatarSizeLG,
		}))},
	})
}

func TestGoldenSweepEmptyState(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "empty_state_default", HTML: utils.Render(t, EmptyState(EmptyStateProps{
			Title:       "No results found",
			Description: "Try adjusting your search or filters.",
			ActionText:  "Clear filters",
			ActionHref:  "/reset",
		}))},
		{Name: "empty_state_minimal", HTML: utils.Render(t, EmptyState(DefaultEmptyStateProps()))},
	})
}
