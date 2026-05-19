package display

import (
	"bytes"
	"context"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	activeBadgeText      = "Active"
	dropdownLabelActions = "Actions"
	dropdownHrefEdit     = "/edit"
	avatarAltAlice       = "Alice"
	tableHeaderName      = "Name"
	dropdownLabelMenu    = "Menu"
	dropdownItemEdit     = "Edit"
	accordionIDFAQ1      = "faq1"
	testNonce            = "test-nonce"
	cssClassMt4          = "mt-4"
)

func TestA11yAttributes(t *testing.T) {
	t.Parallel()

	t.Run("modal has role=dialog", func(t *testing.T) {
		t.Parallel()
		props := DefaultModalProps()
		props.ID = "test-modal"
		props.Title = "Test"
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `role="dialog"`)
		utils.AssertContains(t, output, `aria-modal="true"`)
	})

	t.Run("dropdown has proper ARIA", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd"},
			Label:     dropdownLabelActions,
			Items: []DropdownItem{
				{Text: dropdownItemEdit, Href: dropdownHrefEdit},
			},
		}))
		utils.AssertContains(t, output, `aria-expanded="false"`)
		utils.AssertContains(t, output, `aria-haspopup="true"`)
		utils.AssertContains(t, output, `role="menu"`)
		utils.AssertContains(t, output, `role="menuitem"`)
	})

	t.Run("tabs have proper ARIA", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Tabs: []Tab{
				{ID: "a", Label: activeBadgeText},
				{ID: "b", Label: "Inactive"},
			},
		}))
		utils.AssertContains(t, output, `role="tablist"`)
		utils.AssertContains(t, output, `role="tab"`)
		utils.AssertContains(t, output, `aria-selected="true"`)
	})

	t.Run("tooltip has role=tooltip", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "Help text",
			Position: TooltipPositionTop,
		}))
		utils.AssertContains(t, output, `role="tooltip"`)
	})

	t.Run("accordion has aria-expanded", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "a", Title: "Q1", Open: true},
				{ID: "b", Title: "Q2"},
			},
		}))
		utils.AssertContains(t, output, `aria-expanded="true"`)
		utils.AssertContains(t, output, `aria-expanded="false"`)
	})

	t.Run("avatar image has alt text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src: "/avatar.jpg",
			Alt: avatarAltAlice,
		}))
		utils.AssertContains(t, output, `alt="Alice"`)
	})
}

func TestDarkModeClasses(t *testing.T) {
	t.Parallel()

	t.Run("card has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(DefaultCardProps()))
		utils.AssertContains(t, output, "dark:bg-slate-800")
		utils.AssertContains(t, output, "dark:border-slate-700")
	})

	t.Run("badge has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultBadgeProps()
		props.Text = activeBadgeText
		output := utils.Render(t, Badge(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("table has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{tableHeaderName},
			Rows:    []TableRow{SimpleTableRow("Alice")},
		}))
		utils.AssertContains(t, output, "dark:divide-slate-700")
		utils.AssertContains(t, output, "dark:text-gray-300")
	})

	t.Run("dropdown has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd"},
			Label:     dropdownLabelMenu,
		}))
		utils.AssertContains(t, output, "dark:bg-slate-800")
	})

	t.Run("avatar initials have dark mode background", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "AB",
		}))
		utils.AssertContains(t, output, "bg-blue-600")
	})
}

func TestDefaultPropsConstructors(t *testing.T) {
	t.Parallel()

	t.Run("DefaultCardProps returns valid props", func(t *testing.T) {
		t.Parallel()
		props := DefaultCardProps()
		if props.Padding != CardPaddingMD {
			t.Errorf("Padding = %q, want %q", props.Padding, CardPaddingMD)
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, "<div")
	})

	t.Run("DefaultBadgeProps returns valid props", func(t *testing.T) {
		t.Parallel()
		props := DefaultBadgeProps()
		if props.Type != BadgeNeutral {
			t.Errorf("Type = %q, want %q", props.Type, BadgeNeutral)
		}
		output := utils.Render(t, Badge(props))
		utils.AssertContains(t, output, "<span")
	})

	t.Run("DefaultModalProps returns valid props", func(t *testing.T) {
		t.Parallel()
		props := DefaultModalProps()
		if props.Size != ModalSizeMD {
			t.Errorf("Size = %q, want %q", props.Size, ModalSizeMD)
		}
	})
}

func TestDropdownXSSSafety(t *testing.T) {
	t.Parallel()

	t.Run("ID with special chars is safely interpolated", func(t *testing.T) {
		t.Parallel()
		maliciousID := `<script>alert('xss')</script>`
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: maliciousID},
			Label:     dropdownLabelMenu,
		}))
		utils.AssertNotContains(t, output, `<script>alert('xss')</script>`)
		utils.AssertContains(t, output, `&lt;script&gt;`)
	})
}

func BenchmarkHotPaths(b *testing.B) {
	// Note: benchmarks use render directly since utils.Render requires *testing.T
	b.Run("Class merge", func(b *testing.B) {
		for b.Loop() {
			utils.Class("px-4 py-2", "px-6", "bg-red-500", "bg-blue-500")
		}
	})

	b.Run("Badge render", func(b *testing.B) {
		props := DefaultBadgeProps()
		props.Text = activeBadgeText
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Badge(props).Render(context.Background(), &buf)
		}
	})
}
