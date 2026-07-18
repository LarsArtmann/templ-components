package display

import (
	"testing"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestModalEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("modal without title renders no header", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "no-title-modal"},
			Open:      true,
			Size:      ModalSizeMD,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `id="no-title-modal"`)
		utils.AssertNotContains(t, output, `id="no-title-modal-title"`)
		utils.AssertNotContains(t, output, "aria-label=\"Close\"")
	})

	t.Run("dropdown with empty items list renders button only", func(t *testing.T) {
		t.Parallel()

		props := DropdownProps{
			BaseProps: utils.BaseProps{ID: "empty-dd"},
			Items:     []DropdownItem{},
		}
		output := utils.Render(t, Dropdown(props))
		utils.AssertContains(t, output, `id="empty-dd"`)
		utils.AssertContains(t, output, `id="empty-dd-button"`)
	})

	t.Run("dropdown item with both Href and action renders link", func(t *testing.T) {
		t.Parallel()

		props := DropdownProps{
			BaseProps: utils.BaseProps{ID: "both-dd"},
			Items: []DropdownItem{
				{Text: "Link", Href: "/link"},
			},
		}
		output := utils.Render(t, Dropdown(props))
		utils.AssertContains(t, output, "/link")
	})
}

func TestAccordionEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty item ID auto-generates", func(t *testing.T) {
		t.Parallel()

		props := AccordionProps{
			Items: []AccordionItem{
				{ID: "", Title: "Missing"},
			},
		}
		output := utils.Render(t, Accordion(props))
		utils.AssertContains(t, output, `id="tc-accordion-`)
	})

	t.Run("empty items list renders container only", func(t *testing.T) {
		t.Parallel()

		props := AccordionProps{
			Items: []AccordionItem{},
		}
		output := utils.Render(t, Accordion(props))
		utils.AssertContains(t, output, "divide-y")
	})

	t.Run("single item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			Accordion(AccordionProps{Items: []AccordionItem{{ID: "a", Title: "Q1"}}}),
		)
		utils.AssertContains(t, output, "Q1")
	})

	t.Run("closed item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			Accordion(AccordionProps{Items: []AccordionItem{{ID: "a", Title: "Q1", Open: false}}}),
		)
		utils.AssertContains(t, output, `<details`)
	})
}

func TestStatCardEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props StatCardProps
		want  []string
	}{
		{"no change", StatCardProps{Value: "100", Label: "Users", Trend: TrendNone}, []string{"100", "Users"}},
		{"up trend", StatCardProps{Value: "100", Label: "Users", Change: "+12%", Trend: TrendUp}, []string{"100", "+12%", "text-green-600"}},
		{"down trend", StatCardProps{Value: "100", Label: "Users", Change: "-5%", Trend: TrendDown}, []string{"100", "-5%", "text-red-600"}},
		{"empty", StatCardProps{}, []string{"<dl>"}},
		{"custom id/class", StatCardProps{BaseProps: utils.BaseProps{ID: "stat-1", Class: "mt-4"}, Value: "42", Label: "Count"}, []string{`id="stat-1"`, "mt-4"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, StatCard(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestEmptyStateEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		props   EmptyStateProps
		want    []string
		notWant []string
	}{
		{"no icon", EmptyStateProps{Title: "No data", Icon: ""}, []string{"No data"}, []string{"h-12 w-12"}},
		{"no description", EmptyStateProps{Title: "Empty", Description: ""}, []string{"Empty"}, []string{"text-gray-500"}},
		{"no action", EmptyStateProps{Title: "Done", ActionText: ""}, []string{"Done"}, []string{"bg-blue-600"}},
		{"button without href", EmptyStateProps{Title: "Oops", ActionText: "Retry", ActionHref: ""}, []string{"Retry", "<button"}, []string{"<a"}},
		{"link with href", EmptyStateProps{Title: "Oops", ActionText: "Retry", ActionHref: "/retry"}, []string{"Retry", `<a href="/retry"`}, []string{"<button"}},
		{"custom id/class", EmptyStateProps{BaseProps: utils.BaseProps{ID: "empty-1", Class: "my-8"}, Title: "Nothing", Description: "Here"}, []string{`id="empty-1"`, "my-8", "Nothing", "Here"}, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, EmptyState(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)

			for _, nw := range tt.notWant {
				utils.AssertNotContains(t, output, nw)
			}
		})
	}
}

func TestCardEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		props   CardProps
		want    []string
		notWant []string
	}{
		{"minimal", CardProps{Padding: CardPaddingMD}, []string{"<div"}, []string{"<h3", "border-b"}},
		{"title only", CardProps{Title: "Users"}, []string{"<h3", "Users", "border-b"}, nil},
		{"subtitle without title — no header", CardProps{Subtitle: "Details"}, []string{"<div"}, []string{"Details", "<h3"}},
		{"none padding", CardProps{Padding: CardPaddingNone}, []string{"<div"}, []string{"px-4"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Card(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)

			for _, nw := range tt.notWant {
				utils.AssertNotContains(t, output, nw)
			}
		})
	}
}

func TestSimpleCardEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props SimpleCardProps
		want  string
	}{
		{"none padding", SimpleCardProps{Padding: CardPaddingNone}, "rounded-lg"},
		{"sm padding", SimpleCardProps{Padding: CardPaddingSM}, "px-3"},
		{"lg padding", SimpleCardProps{Padding: CardPaddingLG}, "px-6"},
		{"unknown fallback", SimpleCardProps{Padding: CardPadding("unknown")}, "px-4"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, SimpleCard(tt.props))
			utils.AssertContains(t, output, tt.want)
		})
	}
}

func TestBadgeEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props BadgeProps
		want  []string
	}{
		{"neutral", BadgeProps{Text: "Test", Type: BadgeNeutral}, []string{"bg-gray-100"}},
		{"success", BadgeProps{Text: "Test", Type: BadgeSuccess}, []string{"green"}},
		{"error", BadgeProps{Text: "Test", Type: BadgeError}, []string{"red"}},
		{"warning", BadgeProps{Text: "Test", Type: BadgeWarning}, []string{"yellow"}},
		{"info", BadgeProps{Text: "Test", Type: BadgeInfo}, []string{"blue"}},
		{"with dot", BadgeProps{Text: "Live", Type: BadgeSuccess, Dot: true}, []string{"bg-green"}},
		{"custom id/class", BadgeProps{BaseProps: utils.BaseProps{ID: "badge-1", Class: "ml-2"}, Text: "Test"}, []string{`id="badge-1"`, "ml-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Badge(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestAvatarEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props AvatarProps
		want  []string
	}{
		{"initials", AvatarProps{Initials: "AB"}, []string{"AB", "bg-blue-600"}},
		{"image", AvatarProps{Src: "/a.jpg", Alt: "User"}, []string{`src="/a.jpg"`, `alt="User"`}},
		{"online with image", AvatarProps{Src: "/a.jpg", Alt: "User", Status: AvatarStatusOnline}, []string{"bg-green-400"}},
		{"offline with image", AvatarProps{Src: "/a.jpg", Alt: "User", Status: AvatarStatusOffline}, []string{"bg-gray-400"}},
		{"xs size", AvatarProps{Initials: "AB", Size: AvatarSizeXS}, []string{"h-6"}},
		{"custom id/class", AvatarProps{BaseProps: utils.BaseProps{ID: "av", Class: "mr-2"}, Initials: "AB"}, []string{`id="av"`, "mr-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Avatar(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestTableEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		props   TableProps
		want    []string
		notWant []string
	}{
		{"empty rows", TableProps{Headers: []string{"Name"}, Rows: []TableRow{}}, []string{"<table"}, []string{"<tbody><tr"}},
		{"striped", TableProps{Headers: []string{"A"}, Rows: []TableRow{SimpleTableRow("1"), SimpleTableRow("2")}, Striped: true}, []string{"bg-gray-50"}, nil},
		{"hover", TableProps{Headers: []string{"A"}, Rows: []TableRow{SimpleTableRow("1")}, Hover: true}, []string{"hover:bg-gray-100"}, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Table(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)

			for _, nw := range tt.notWant {
				utils.AssertNotContains(t, output, nw)
			}
		})
	}
}

func TestTabsEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props TabsProps
		want  []string
	}{
		{"pills variant", TabsProps{Variant: TabsPills, ActiveTabID: "a", Tabs: []Tab{{ID: "a", Label: "A"}}}, []string{"space-x-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Tabs(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestDropdownEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props DropdownProps
		want  []string
	}{
		{"right position", DropdownProps{BaseProps: utils.BaseProps{ID: "dd"}, Label: "A", Position: DropdownPositionRight}, []string{`data-dropdown-align="right"`}},
		{"left position", DropdownProps{BaseProps: utils.BaseProps{ID: "dd"}, Label: "A", Position: DropdownPositionLeft}, []string{"hidden"}},
		{"default position", DropdownProps{BaseProps: utils.BaseProps{ID: "dd"}, Label: "A"}, []string{"hidden"}},
		{"item with icon", DropdownProps{BaseProps: utils.BaseProps{ID: "dd"}, Label: "A", Items: []DropdownItem{{Text: "Edit", Href: "/edit", Icon: icons.Edit}}}, []string{"Edit"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Dropdown(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestTooltipEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		props   TooltipProps
		want    []string
		notWant []string
	}{
		{"right", TooltipProps{Text: "Hint", Position: TooltipPositionRight}, []string{"left-full"}, []string{"bottom-full"}},
		{"left", TooltipProps{Text: "Hint", Position: TooltipPositionLeft}, []string{"right-full"}, []string{"bottom-full"}},
		{"bottom", TooltipProps{Text: "Hint", Position: TooltipPositionBottom}, []string{"top-full"}, []string{"bottom-full mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Tooltip(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)

			for _, nw := range tt.notWant {
				utils.AssertNotContains(t, output, nw)
			}
		})
	}
}

func TestAriaLabelPropagation(t *testing.T) {
	t.Parallel()

	label := "Custom aria label"

	tests := []struct {
		name  string
		props any
	}{
		{
			"Card",
			CardProps{
				Title:     "T",
				Padding:   CardPaddingMD,
				BaseProps: utils.BaseProps{AriaLabel: label},
			},
		},
		{
			"SimpleCard",
			SimpleCardProps{
				Padding:   CardPaddingSM,
				BaseProps: utils.BaseProps{AriaLabel: label},
			},
		},
		{
			"StatCard",
			StatCardProps{
				Value:     "42",
				Label:     "L",
				Trend:     TrendNone,
				BaseProps: utils.BaseProps{AriaLabel: label},
			},
		},
		{
			"Accordion",
			AccordionProps{
				Items:     []AccordionItem{{ID: "a", Title: "A"}},
				BaseProps: utils.BaseProps{AriaLabel: label},
			},
		},
		{
			"Table",
			TableProps{
				Headers:   []string{"H"},
				BaseProps: utils.BaseProps{AriaLabel: label},
			},
		},
		{
			"Dropdown",
			DropdownProps{
				Label:     "Menu",
				Items:     []DropdownItem{{Text: "I"}},
				BaseProps: utils.BaseProps{ID: "dd", AriaLabel: label},
			},
		},
		{
			"Tabs",
			TabsProps{
				ActiveTabID: "t1",
				Tabs:        []Tab{{ID: "t1", Label: "T"}},
				BaseProps:   utils.BaseProps{AriaLabel: label},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output string

			switch p := tt.props.(type) {
			case CardProps:
				output = utils.Render(t, Card(p))
			case SimpleCardProps:
				output = utils.Render(t, SimpleCard(p))
			case StatCardProps:
				output = utils.Render(t, StatCard(p))
			case AccordionProps:
				output = utils.Render(t, Accordion(p))
			case TableProps:
				output = utils.Render(t, Table(p))
			case DropdownProps:
				output = utils.Render(t, Dropdown(p))
			case TabsProps:
				output = utils.Render(t, Tabs(p))
			}

			utils.AssertContains(t, output, `aria-label="`+label+`"`)
		})
	}
}
