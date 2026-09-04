package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
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

func TestGoldenSweepStatCardTones(t *testing.T) {
	t.Parallel()

	base := StatCardProps{Label: "Uptime", Value: "99.9%", Icon: icons.Check, Change: "0.1%", Trend: TrendUp}

	renderTone := func(tone StatTone) string {
		props := base
		props.Tone = tone

		return utils.Render(t, StatCard(props))
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "stat_card_tone_blue", HTML: renderTone(StatToneBlue)},
		{Name: "stat_card_tone_green", HTML: renderTone(StatToneGreen)},
		{Name: "stat_card_tone_yellow", HTML: renderTone(StatToneYellow)},
		{Name: "stat_card_tone_red", HTML: renderTone(StatToneRed)},
		{Name: "stat_card_tone_purple", HTML: renderTone(StatTonePurple)},
		{Name: "stat_card_tone_unknown_falls_back_blue", HTML: renderTone(StatTone("bogus"))},
	})
}

func TestGoldenSweepButtonOutlineVariants(t *testing.T) {
	t.Parallel()

	renderVariant := func(variant ButtonType) string {
		return utils.Render(t, Button(ButtonProps{Text: "Undo", Variant: variant}))
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "button_outline_danger", HTML: renderVariant(ButtonOutlineDanger)},
		{Name: "button_outline_warning", HTML: renderVariant(ButtonOutlineWarning)},
		{Name: "button_outline_success", HTML: renderVariant(ButtonOutlineSuccess)},
		{Name: "button_outline_info", HTML: renderVariant(ButtonOutlineInfo)},
	})
}

func TestGoldenSweepButtonWired(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "button_wired_htmx_full", HTML: utils.Render(t, Button(ButtonProps{
			Text: "Load more",
			Wire: &wire.Action{
				URL:    "/api/items",
				Event:  wire.EventClick,
				Target: "#items",
			},
		}))},
		{Name: "button_wired_htmx_post", HTML: utils.Render(t, Button(ButtonProps{
			Text: "Save",
			Wire: &wire.Action{Method: wire.MethodPost, URL: "/api/items"},
		}))},
		{Name: "button_wired_datastar", HTML: utils.Render(t, Button(ButtonProps{
			Text: "Load via Datastar",
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/items",
			},
		}))},
		{Name: "button_wired_datastar_custom_event", HTML: utils.Render(t, Button(ButtonProps{
			Text: "Search",
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/search",
				Event:     wire.EventInput,
			},
		}))},
		{Name: "button_wired_empty_url_inert", HTML: utils.Render(t, Button(ButtonProps{
			Text: "Inert",
			Wire: &wire.Action{Target: "#nowhere"},
		}))},
		{Name: "button_wired_link_href", HTML: utils.Render(t, Button(ButtonProps{
			Text: "Wired link",
			Href: "/fallback",
			Wire: &wire.Action{URL: "/api/items"},
		}))},
	})
}
