package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// ---------------------------------------------------------------------------
// BaseProps propagation on components where ID/AriaLabel branches are untested
// ---------------------------------------------------------------------------

func TestListNoteBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ListNote(ListNoteProps{
		Shown: 5, Total: 20,
		BaseProps: utils.BaseProps{
			ID:        "list-note-1",
			AriaLabel: "Truncation notice",
		},
	}))
	utils.AssertContainsAll(t, output, `id="list-note-1"`, `aria-label="Truncation notice"`)
}

func TestGridContainerResponsiveWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:                GridCols3,
		ContainerResponsive: true,
		BaseProps: utils.BaseProps{
			ID:        "cgrid",
			AriaLabel: "Responsive grid",
		},
	}))
	utils.AssertContainsAll(t, output, `id="cgrid"`, `aria-label="Responsive grid"`, "@container")
}

func TestGridContainerClassFallback(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:                GridCols("bogus"),
		ContainerResponsive: true,
	}))
	utils.AssertContains(t, output, "grid-cols-1")
}

func TestImageBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Image(ImageProps{
		Src: "photo.jpg",
		BaseProps: utils.BaseProps{
			ID:        "avatar-img",
			AriaLabel: "User avatar",
		},
	}))
	utils.AssertContainsAll(t, output, `id="avatar-img"`, `aria-label="User avatar"`)
}

func TestImageEmptySrc(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Image(ImageProps{Src: ""}))
	utils.AssertNotContains(t, output, "<img")
}

func TestCopyButtonAnchorWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CopyButton(CopyButtonProps{
		Text: "Copy me",
		Href: "/copy",
		Icon: true,
		BaseProps: utils.BaseProps{
			ID:        "copy-link",
			AriaLabel: "Copy link to clipboard",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="copy-link"`,
		`aria-label="Copy link to clipboard"`,
		`href="/copy"`,
		`data-tc-copy`,
	)
}

func TestCopyButtonButtonWithID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CopyButton(CopyButtonProps{
		Text:      "Copy",
		BaseProps: utils.BaseProps{ID: "copy-btn"},
	}))
	utils.AssertContains(t, output, `id="copy-btn"`)
}

func TestCountBadgeWithID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CountBadge(CountBadgeProps{
		Count:     5,
		BaseProps: utils.BaseProps{ID: "notif-count"},
	}))
	utils.AssertContains(t, output, `id="notif-count"`)
}

func TestDefinitionGridWithID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
		Items: []DefinitionItem{
			{Term: "Name", Detail: "Alice"},
		},
		BaseProps: utils.BaseProps{ID: "def-grid"},
	}))
	utils.AssertContains(t, output, `id="def-grid"`)
}

func TestDefinitionGridDetailComponent(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
		Items: []DefinitionItem{
			{
				Term:            "Status",
				DetailComponent: templ.Raw(`<span class="text-green-600">Active</span>`),
			},
		},
	}))
	utils.AssertContains(t, output, "Active")
}

func TestStatCardWithHTMX(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StatCard(StatCardProps{
		Label:    "Revenue",
		Value:    "$42k",
		HxGet:    "/api/stats/revenue",
		HxTarget: "#stat-container",
	}))
	utils.AssertContainsAll(t, output,
		`hx-get="/api/stats/revenue"`,
		`hx-target="#stat-container"`,
	)
}

func TestStatCardHrefWithHTMX(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StatCard(StatCardProps{
		Label: "Revenue",
		Value: "$42k",
		Href:  "/dashboard/revenue",
		HxGet: "/api/stats/revenue",
	}))
	utils.AssertContainsAll(t, output,
		`href="/dashboard/revenue"`,
		`hx-get="/api/stats/revenue"`,
	)
}

// ---------------------------------------------------------------------------
// Default*Props constructor coverage
// ---------------------------------------------------------------------------

func TestDefaultConstructorsRender(t *testing.T) {
	t.Parallel()

	t.Run("DefaultCopyButtonProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(DefaultCopyButtonProps()))
		utils.AssertContains(t, output, "data-tc-copy")
	})

	t.Run("DefaultCountBadgeProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(DefaultCountBadgeProps()))
		_ = output // just ensure no panic
	})

	t.Run("DefaultDefinitionGridProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionGrid(DefaultDefinitionGridProps()))
		utils.AssertContains(t, output, "grid")
	})

	t.Run("DefaultDefinitionListProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionList(DefaultDefinitionListProps()))
		utils.AssertContains(t, output, "<dl")
	})

	t.Run("DefaultPageHeaderProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(DefaultPageHeaderProps()))
		utils.AssertContains(t, output, "<h1")
	})

	t.Run("DefaultImageProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(DefaultImageProps()))
		utils.AssertNotContains(t, output, "<img")
	})
}

// ---------------------------------------------------------------------------
// Additional branch coverage
// ---------------------------------------------------------------------------

func TestSimpleCardBodySlotCoverage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleCard(SimpleCardProps{
		Body: templ.Raw("<p>Simple body slot</p>"),
	}))
	utils.AssertContains(t, output, "Simple body slot")
}

func TestGridAllGapVariants(t *testing.T) {
	t.Parallel()
	for _, gap := range []GridGap{GridGapSM, GridGapMD, GridGapLG, GridGapXL} {
		output := utils.Render(t, Grid(GridProps{Cols: GridCols2, Gap: gap}))
		switch gap {
		case GridGapSM:
			utils.AssertContains(t, output, "gap-2")
		case GridGapMD:
			utils.AssertContains(t, output, "gap-4")
		case GridGapLG:
			utils.AssertContains(t, output, "gap-6")
		case GridGapXL:
			utils.AssertContains(t, output, "gap-8")
		}
	}
}

func TestGridAllColsVariants(t *testing.T) {
	t.Parallel()
	for _, cols := range []GridCols{GridCols1, GridCols2, GridCols3, GridCols4, GridCols5, GridCols6} {
		output := utils.Render(t, Grid(GridProps{Cols: cols}))
		utils.AssertContains(t, output, "grid")
	}
}

func TestImageFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Image(ImageProps{
		Src:         "photo.jpg",
		Alt:         "A scenic photo",
		Width:       800,
		Height:      600,
		Lazy:        false,
		Rounded:     true,
		FallbackSrc: "placeholder.jpg",
		BaseProps: utils.BaseProps{
			ID:    "hero-img",
			Class: "border-2",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="hero-img"`,
		`alt="A scenic photo"`,
		`width="800"`,
		`height="600"`,
		"rounded-full",
		`data-tc-img-fallback`,
		"placeholder.jpg",
		`loading="eager"`,
	)
}

func TestAccordionFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Accordion(AccordionProps{
		Items: []AccordionItem{
			{Title: "Section 1", Content: templ.Raw("Content 1")},
			{Title: "Section 2", Content: templ.Raw("Content 2"), Open: true},
		},
		BaseProps: utils.BaseProps{
			ID:        "accordion-1",
			AriaLabel: "FAQ accordion",
			Class:     "border rounded-lg",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="accordion-1"`,
		`aria-label="FAQ accordion"`,
		"Section 1",
		"Section 2",
		"Content 2",
	)
}

func TestDropdownFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Dropdown(DropdownProps{
		Label: "Actions",
		Items: []DropdownItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Delete", Href: "/delete"},
		},
		BaseProps: utils.BaseProps{
			ID:        "dd-1",
			AriaLabel: "Actions menu",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="dd-1"`,
		"Edit",
		"Delete",
	)
}

func TestButtonAsLinkWithIcon(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Button(ButtonProps{
		Text:    "Settings",
		Href:    "/settings",
		Variant: ButtonPrimary,
		Icon:    icons.Icon(icons.Settings, "w-4 h-4"),
		BaseProps: utils.BaseProps{
			ID:        "settings-link",
			AriaLabel: "Open settings",
		},
	}))
	utils.AssertContainsAll(t, output,
		`href="/settings"`,
		`id="settings-link"`,
		"Settings",
	)
}

func TestDrawerFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Drawer(DrawerProps{
		Title: "Filter Options",
		Side:  "right",
		BaseProps: utils.BaseProps{
			ID:        "drawer-1",
			AriaLabel: "Filter drawer",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="drawer-1"`,
		"Filter Options",
	)
}

func TestModalFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Modal(ModalProps{
		Title: "Confirm Action",
		Size:  ModalSizeLG,
		BaseProps: utils.BaseProps{
			ID:        "modal-1",
			AriaLabel: "Confirmation dialog",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="modal-1"`,
		"Confirm Action",
	)
}

func TestTooltipFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tooltip(TooltipProps{
		Text:      "Helpful tip",
		Position:  TooltipPositionTop,
		BaseProps: utils.BaseProps{ID: "tip-1"},
	}))
	utils.AssertContainsAll(t, output, `id="tip-1"`, "Helpful tip")
}

func TestTabsDefaultVariant(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tabs(TabsProps{
		Variant:     TabsDefault,
		ActiveTabID: "tab1",
		Tabs: []Tab{
			{ID: "tab1", Label: "Overview"},
			{ID: "tab2", Label: "Details"},
		},
		BaseProps: utils.BaseProps{
			ID:        "settings-tabs",
			AriaLabel: "Settings",
		},
	}))
	utils.AssertContainsAll(t, output, `id="settings-tabs"`, "Overview", "Details")
}

func TestDefinitionListDetailComponent(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DefinitionList(DefinitionListProps{
		Items: []DefinitionItem{
			{Term: "Status", DetailComponent: templ.Raw(`<span class="badge">Active</span>`)},
			{Term: "Email", Detail: "alice@example.com"},
		},
		BaseProps: utils.BaseProps{
			ID:        "def-list",
			AriaLabel: "User details",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="def-list"`,
		`aria-label="User details"`,
		"Active",
		"alice@example.com",
	)
}

func TestAvatarAllSizesAndShapes(t *testing.T) {
	t.Parallel()
	for _, size := range []AvatarSize{AvatarSizeXS, AvatarSizeSM, AvatarSizeMD, AvatarSizeLG} {
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "JD",
			Size:     size,
			Status:   AvatarStatusOnline,
		}))
		utils.AssertContains(t, output, "JD")
	}
}

func TestBadgeAllTypesAndSizes(t *testing.T) {
	t.Parallel()
	for _, btype := range []BadgeType{BadgePrimary, BadgeSuccess, BadgeError, BadgeWarning, BadgeInfo, BadgeNeutral} {
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Test",
			Type: btype,
			Size: BadgeSizeMD,
		}))
		utils.AssertContains(t, output, "Test")
	}
}
