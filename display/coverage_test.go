package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// --- Badge coverage (was 38.2%) ---

func TestBadgeHrefRendersAsAnchor(t *testing.T) {
	t.Parallel()
	t.Run("href renders <a> tag", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Link",
			Href: "/page",
			Type: BadgePrimary,
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, `href="/page"`)
		utils.AssertContains(t, output, "Link")
		utils.AssertNotContains(t, output, `<span`)
	})

	t.Run("href with ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{ID: "link-badge"},
			Text:      "Go",
			Href:      "/go",
		}))
		utils.AssertContains(t, output, `id="link-badge"`)
		utils.AssertContains(t, output, `href="/go"`)
	})

	t.Run("href with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{AriaLabel: "Navigate"},
			Text:      "Nav",
			Href:      "/nav",
		}))
		utils.AssertContains(t, output, `aria-label="Navigate"`)
	})

	t.Run("href with dot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Live",
			Href: "/live",
			Dot:  true,
			Type: BadgeSuccess,
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, "rounded-full")
		utils.AssertContains(t, output, "bg-green-500")
	})

	t.Run("href with pill", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Pill",
			Href: "/pill",
			Pill: true,
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("href with class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{Class: "my-class"},
			Text:      "Styled",
			Href:      "/styled",
		}))
		utils.AssertContains(t, output, "my-class")
	})

	t.Run("href with custom attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-testid": "badge-link"}},
			Text:      "Attr",
			Href:      "/attr",
		}))
		utils.AssertContains(t, output, `data-testid="badge-link"`)
	})

	t.Run("all types with href", func(t *testing.T) {
		t.Parallel()
		for _, bt := range []BadgeType{BadgePrimary, BadgeSuccess, BadgeWarning, BadgeError, BadgeInfo, BadgeNeutral} {
			output := utils.Render(t, Badge(BadgeProps{
				Text: string(bt),
				Type: bt,
				Href: "/" + string(bt),
			}))
			utils.AssertContains(t, output, `<a`)
			utils.AssertContains(t, output, string(bt))
		}
	})

	t.Run("sizes with href", func(t *testing.T) {
		t.Parallel()
		for _, size := range []BadgeSize{BadgeSizeSM, BadgeSizeMD, BadgeSizeLG} {
			output := utils.Render(t, Badge(BadgeProps{
				Text: "Sized",
				Href: "/sized",
				Size: size,
			}))
			utils.AssertContains(t, output, `<a`)
		}
	})

	t.Run("unknown type with href falls back", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Unknown",
			Type: BadgeType("custom"),
			Href: "/custom",
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, "bg-gray-100")
	})
}

// --- Button coverage (was 53.7%) ---

func TestButtonIcon(t *testing.T) {
	t.Parallel()
	t.Run("button with icon component", func(t *testing.T) {
		t.Parallel()
		icon := icons.Icon(icons.Plus, "h-4 w-4")
		output := utils.Render(t, Button(ButtonProps{
			Text: "Add",
			Icon: icon,
		}))
		utils.AssertContains(t, output, "Add")
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("link button with icon", func(t *testing.T) {
		t.Parallel()
		icon := icons.Icon(icons.ArrowRight, "h-4 w-4")
		output := utils.Render(t, Button(ButtonProps{
			Text: "Next",
			Href: "/next",
			Icon: icon,
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, "<svg")
		utils.AssertContains(t, output, "Next")
	})

	t.Run("button with ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			BaseProps: utils.BaseProps{
				ID:    "my-btn",
				Class: "extra-class",
			},
			Text: "Styled",
		}))
		utils.AssertContains(t, output, `id="my-btn"`)
		utils.AssertContains(t, output, "extra-class")
	})

	t.Run("button with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			BaseProps: utils.BaseProps{AriaLabel: "Submit form"},
			Text:      "Submit",
		}))
		utils.AssertContains(t, output, `aria-label="Submit form"`)
	})

	t.Run("button with custom attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-testid": "btn"}},
			Text:      "Test",
		}))
		utils.AssertContains(t, output, `data-testid="btn"`)
	})

	t.Run("link button with ID and aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			BaseProps: utils.BaseProps{
				ID:        "link-btn",
				AriaLabel: "Navigate",
			},
			Text: "Go",
			Href: "/go",
		}))
		utils.AssertContains(t, output, `id="link-btn"`)
		utils.AssertContains(t, output, `aria-label="Navigate"`)
	})

	t.Run("default variant button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Default"}))
		utils.AssertContains(t, output, "bg-blue-600")
	})

	t.Run("unknown variant falls back to primary", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			Text:    "Fallback",
			Variant: ButtonType("unknown"),
		}))
		utils.AssertContains(t, output, "bg-blue-600")
	})
}

// --- Dropdown coverage (was 69.1%) ---

func TestDropdownButtonItems(t *testing.T) {
	t.Parallel()
	t.Run("button kind item renders <button>", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-btn"},
			Label:     "Actions",
			Items: []DropdownItem{
				{Text: "Delete", Kind: DropdownItemButton},
			},
		}))
		utils.AssertContains(t, output, "Delete")
		utils.AssertContains(t, output, `type="button"`)
	})

	t.Run("disabled button kind item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-dis"},
			Label:     "Actions",
			Items: []DropdownItem{
				{Text: "Archive", Kind: DropdownItemButton, Disabled: true},
			},
		}))
		utils.AssertContains(t, output, "Archive")
		utils.AssertContains(t, output, `disabled`)
		utils.AssertContains(t, output, `aria-disabled="true"`)
	})

	t.Run("disabled link item renders span", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-dislink"},
			Label:     "Menu",
			Items: []DropdownItem{
				{Text: "Disabled Link", Href: "/link", Disabled: true},
			},
		}))
		utils.AssertContains(t, output, "Disabled Link")
		utils.AssertContains(t, output, `aria-disabled="true"`)
		utils.AssertContains(t, output, `tabindex="-1"`)
	})

	t.Run("disabled link item with icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-disicon"},
			Label:     "Menu",
			Items: []DropdownItem{
				{Text: "Settings", Href: "/settings", Icon: icons.Settings, Disabled: true},
			},
		}))
		utils.AssertContains(t, output, "Settings")
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("IsLink with explicit kind", func(t *testing.T) {
		t.Parallel()
		linkItem := DropdownItem{Kind: DropdownItemLink, Text: "Link"}
		btnItem := DropdownItem{Kind: DropdownItemButton, Text: "Btn"}
		if !linkItem.IsLink() {
			t.Error("expected DropdownItemLink to be link")
		}
		if btnItem.IsLink() {
			t.Error("expected DropdownItemButton to not be link")
		}
	})

	t.Run("dropdown with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-nonce", Nonce: "test123"},
			Label:     "Nonce",
		}))
		utils.AssertContains(t, output, `nonce="test123"`)
	})

	t.Run("dropdown with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-aria", AriaLabel: "Actions menu"},
			Label:     "Actions",
		}))
		utils.AssertContains(t, output, `aria-label="Actions menu"`)
	})

	t.Run("dropdown with custom attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-attr", Attrs: templ.Attributes{"data-testid": "dd"}},
			Label:     "Attr",
		}))
		utils.AssertContains(t, output, `data-testid="dd"`)
	})

	t.Run("button item with icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-btnicon"},
			Label:     "Menu",
			Items: []DropdownItem{
				{Text: "Edit", Kind: DropdownItemButton, Icon: icons.Edit},
			},
		}))
		utils.AssertContains(t, output, "Edit")
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("link item with icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-linkicon"},
			Label:     "Menu",
			Items: []DropdownItem{
				{Text: "Profile", Href: "/profile", Icon: icons.Users},
			},
		}))
		utils.AssertContains(t, output, "Profile")
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("item with custom attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd-itemattr"},
			Label:     "Menu",
			Items: []DropdownItem{
				{Text: "Item", Href: "/item", Attrs: templ.Attributes{"data-action": "click"}},
			},
		}))
		utils.AssertContains(t, output, `data-action="click"`)
	})
}

// --- Modal coverage (was 66.7%) ---

func TestModalSizes(t *testing.T) {
	t.Parallel()
	for _, size := range []ModalSize{ModalSizeSM, ModalSizeMD, ModalSizeLG, ModalSizeXL, ModalSizeFull} {
		t.Run("size_"+string(size), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Modal(ModalProps{
				BaseProps: utils.BaseProps{ID: "modal-" + string(size)},
				Title:     "Test",
				Size:      size,
				Open:      true,
			}))
			utils.AssertContains(t, output, `id="modal-`+string(size)+`"`)
		})
	}
}

func TestModalClosed(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Modal(ModalProps{
		BaseProps: utils.BaseProps{ID: "closed-modal"},
		Title:     "Hidden",
		Open:      false,
	}))
	utils.AssertContains(t, output, "opacity-0")
	utils.AssertContains(t, output, "pointer-events-none")
}

func TestModalWithClass(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Modal(ModalProps{
		BaseProps: utils.BaseProps{ID: "styled-modal", Class: "custom-modal"},
		Title:     "Styled",
		Open:      true,
	}))
	utils.AssertContains(t, output, "custom-modal")
}

func TestModalWithAttrs(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Modal(ModalProps{
		BaseProps: utils.BaseProps{ID: "attr-modal", Attrs: templ.Attributes{"data-testid": "modal"}},
		Open:      true,
	}))
	utils.AssertContains(t, output, `data-testid="modal"`)
}

func TestModalWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Modal(ModalProps{
		BaseProps: utils.BaseProps{ID: "aria-modal", AriaLabel: "Custom label"},
		Open:      true,
	}))
	utils.AssertContains(t, output, `aria-label="Custom label"`)
}

func TestModalWithNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Modal(ModalProps{
		BaseProps: utils.BaseProps{ID: "nonce-modal", Nonce: "nonce-xyz"},
		Title:     "Nonce",
		Open:      true,
	}))
	utils.AssertContains(t, output, `nonce="nonce-xyz"`)
}

// --- Tooltip coverage (was 66.2%) ---

func TestTooltipTopPosition(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tooltip(TooltipProps{
		Text:     "Top tip",
		Position: TooltipPositionTop,
	}))
	utils.AssertContains(t, output, "bottom-full")
}

func TestTooltipWithID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tooltip(TooltipProps{
		Text:      "ID tip",
		Position:  TooltipPositionTop,
		BaseProps: utils.BaseProps{ID: "my-tip"},
	}))
	utils.AssertContains(t, output, `id="my-tip"`)
	utils.AssertContains(t, output, `aria-describedby="my-tip-tooltip"`)
	utils.AssertContains(t, output, `id="my-tip-tooltip"`)
}

func TestTooltipWithClass(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tooltip(TooltipProps{
		Text:      "Class tip",
		Position:  TooltipPositionTop,
		BaseProps: utils.BaseProps{Class: "extra-wrap"},
	}))
	utils.AssertContains(t, output, "extra-wrap")
}

func TestTooltipWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tooltip(TooltipProps{
		Text:      "Aria tip",
		Position:  TooltipPositionTop,
		BaseProps: utils.BaseProps{AriaLabel: "Custom aria"},
	}))
	utils.AssertContains(t, output, `aria-label="Custom aria"`)
}

func TestTooltipUnknownPosition(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Tooltip(TooltipProps{
		Text:     "Fallback",
		Position: TooltipPosition("unknown"),
	}))
	utils.AssertContains(t, output, "bottom-full")
}

// --- Avatar coverage (was 67.0%) ---

func TestAvatarFallbackSVG(t *testing.T) {
	t.Parallel()
	t.Run("no src and no initials renders SVG fallback", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{}))
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("no src with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			BaseProps: utils.BaseProps{AriaLabel: "User avatar"},
		}))
		utils.AssertContains(t, output, `aria-label="User avatar"`)
	})

	t.Run("no src with ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			BaseProps: utils.BaseProps{ID: "fallback-av"},
		}))
		utils.AssertContains(t, output, `id="fallback-av"`)
	})

	t.Run("square shape", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "AB",
			Shape:    AvatarShapeSquare,
		}))
		utils.AssertContains(t, output, "rounded-lg")
	})

	t.Run("image with class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:       "/img.jpg",
			BaseProps: utils.BaseProps{Class: "border-2"},
		}))
		utils.AssertContains(t, output, "border-2")
	})

	t.Run("all sizes", func(t *testing.T) {
		t.Parallel()
		for _, size := range []AvatarSize{AvatarSizeXS, AvatarSizeSM, AvatarSizeMD, AvatarSizeLG, AvatarSizeXL} {
			output := utils.Render(t, Avatar(AvatarProps{
				Initials: "AB",
				Size:     size,
			}))
			utils.AssertContains(t, output, "AB")
		}
	})

	t.Run("image with attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:       "/img.jpg",
			BaseProps: utils.BaseProps{Attrs: templ.Attributes{"loading": "lazy"}},
		}))
		utils.AssertContains(t, output, `loading="lazy"`)
	})
}

// --- EmptyState coverage (was 65.3%) ---

func TestEmptyStateWithIconAndDescription(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, EmptyState(EmptyStateProps{
		Title:       "No data",
		Description: "Add your first item",
		Icon:        icons.Inbox,
	}))
	utils.AssertContains(t, output, "No data")
	utils.AssertContains(t, output, "Add your first item")
	utils.AssertContains(t, output, "<svg")
}

func TestEmptyStateWithActionAttrs(t *testing.T) {
	t.Parallel()
	t.Run("link action with attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:       "Empty",
			Description: "Nothing here",
			Icon:        icons.Folder,
			ActionText:  "Create",
			ActionHref:  "/create",
			ActionAttrs: templ.Attributes{"data-testid": "create-btn"},
		}))
		utils.AssertContains(t, output, `data-testid="create-btn"`)
		utils.AssertContains(t, output, `href="/create"`)
	})

	t.Run("button action with attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:       "Missing",
			Icon:        icons.ExclamationTriangle,
			ActionText:  "Retry",
			ActionAttrs: templ.Attributes{"data-testid": "retry-btn"},
		}))
		utils.AssertContains(t, output, `data-testid="retry-btn"`)
		utils.AssertContains(t, output, `<button`)
	})
}

// --- Tabs coverage ---

func TestTabsClientSide(t *testing.T) {
	t.Parallel()
	t.Run("client-side tabs with JS", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "tab1",
			Tabs:        []Tab{{ID: "tab1", Label: "First"}, {ID: "tab2", Label: "Second"}},
			ClientSide:  true,
		}))
		utils.AssertContains(t, output, `data-tc-tabs`)
		utils.AssertContains(t, output, `nonce=`)
	})

	t.Run("inactive tab renders without active classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "tab2",
			Tabs:        []Tab{{ID: "tab1", Label: "First"}, {ID: "tab2", Label: "Second"}},
		}))
		utils.AssertContains(t, output, "First")
		utils.AssertContains(t, output, "Second")
	})
}

// --- Table coverage ---

func TestTableBordered(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Table(TableProps{
		Headers:  []string{"A", "B"},
		Rows:     []TableRow{SimpleTableRow("1", "2")},
		Bordered: true,
	}))
	utils.AssertContains(t, output, "border")
}

func TestTableCaption(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Table(TableProps{
		Headers: []string{"Name"},
		Rows:    []TableRow{SimpleTableRow("Alice")},
		Caption: "User list",
	}))
	utils.AssertContains(t, output, "<caption")
	utils.AssertContains(t, output, "User list")
}

// --- Drawer ---

func TestDrawerRender(t *testing.T) {
	t.Parallel()
	t.Run("right drawer with title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "test-drawer"},
			Title:     "Settings",
			Open:      true,
			Side:      DrawerRight,
		}))
		utils.AssertContains(t, output, `id="test-drawer"`)
		utils.AssertContains(t, output, "Settings")
		utils.AssertContains(t, output, `role="dialog"`)
		utils.AssertContains(t, output, `aria-modal="true"`)
		utils.AssertContains(t, output, `translate-x-0`)
	})

	t.Run("left drawer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "left-drawer"},
			Open:      true,
			Side:      DrawerLeft,
		}))
		utils.AssertContains(t, output, `style="inset-y: 0; left: 0;"`)
	})

	t.Run("closed drawer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "closed-drawer"},
			Open:      false,
			Side:      DrawerRight,
		}))
		utils.AssertContains(t, output, "opacity-0")
		utils.AssertContains(t, output, "translate-x-full")
	})

	t.Run("empty ID panics", func(t *testing.T) {
		t.Parallel()
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic for Drawer with empty ID")
			}
		}()
		utils.Render(t, Drawer(DrawerProps{}))
	})

	t.Run("all sizes", func(t *testing.T) {
		t.Parallel()
		for _, size := range []DrawerSize{DrawerSizeSM, DrawerSizeMD, DrawerSizeLG, DrawerSizeXL, DrawerFull} {
			output := utils.Render(t, Drawer(DrawerProps{
				BaseProps: utils.BaseProps{ID: "size-" + string(size)},
				Open:      true,
				Size:      size,
			}))
			utils.AssertContains(t, output, `id="size-`+string(size)+`"`)
		}
	})

	t.Run("with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "aria-drawer", AriaLabel: "Navigation panel"},
			Open:      true,
		}))
		utils.AssertContains(t, output, `aria-label="Navigation panel"`)
	})

	t.Run("with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "nonce-drawer", Nonce: "test-nonce"},
			Open:      true,
		}))
		utils.AssertContains(t, output, `nonce="test-nonce"`)
	})

	t.Run("without title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "no-title-drawer"},
			Open:      true,
		}))
		utils.AssertNotContains(t, output, `id="no-title-drawer-title"`)
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		props := DefaultDrawerProps()
		if props.Side != DrawerRight {
			t.Error("expected right side default")
		}
		if props.Size != DrawerSizeMD {
			t.Error("expected MD size default")
		}
	})
}

func TestAvatarImageBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Avatar(AvatarProps{
		Src:       "/photo.jpg",
		Alt:       "Profile photo",
		BaseProps: utils.BaseProps{ID: "user-avatar", Class: "ring-2", AriaLabel: "Alice avatar"},
	}))
	utils.AssertContains(t, output, `id="user-avatar"`)
	utils.AssertContains(t, output, "ring-2")
	utils.AssertContains(t, output, `aria-label="Alice avatar"`)
}

func TestAvatarImageWithoutID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Avatar(AvatarProps{
		Src: "/photo.jpg",
		Alt: "Profile photo",
	}))
	utils.AssertContains(t, output, `src="/photo.jpg"`)
	utils.AssertNotContains(t, output, `id=""`)
}
