package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestEmptyStateFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with all props including action link and BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			BaseProps: utils.BaseProps{
				ID:        "empty-1",
				Class:     "mt-8",
				AriaLabel: "No data",
				Attrs:     templ.Attributes{"data-ctx": "dashboard"},
			},
			Title:       "No results",
			Description: "Try adjusting your filters",
			Icon:        icons.Search,
			ActionText:  "Clear filters",
			ActionHref:  "/reset",
			ActionAttrs: templ.Attributes{"data-track": "reset"},
		}))
		utils.AssertContains(t, output, `id="empty-1"`)
		utils.AssertContains(t, output, `aria-label="No data"`)
		utils.AssertContains(t, output, `href="/reset"`)
		utils.AssertContains(t, output, "Clear filters")
		utils.AssertContains(t, output, `data-track="reset"`)
	})

	t.Run("action button without href", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:      "Empty",
			ActionText: "Click me",
		}))
		utils.AssertContains(t, output, `<button`)
		utils.AssertContains(t, output, "Click me")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(DefaultEmptyStateProps()))
		utils.AssertContains(t, output, `role="status"`)
	})
}

func TestTableFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with all options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			BaseProps: utils.BaseProps{
				ID:        "data-table",
				Class:     "w-full",
				AriaLabel: "User data",
			},
			Caption:  "User Statistics",
			Headers:  []string{"Name", "Email", "Role"},
			Striped:  true,
			Hover:    true,
			Bordered: true,
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "Alice"}, {Text: "alice@example.com"}, {Text: "Admin"}}},
				{Cells: []TableCell{{Text: "Bob"}, {Text: "bob@example.com"}, {Text: "User"}}},
			},
		}))
		utils.AssertContains(t, output, "User Statistics")
		utils.AssertContains(t, output, "Alice")
		utils.AssertContains(t, output, `id="data-table"`)
	})
}

func TestTabsFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("pills variant with client-side and BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			BaseProps: utils.BaseProps{
				ID:        "my-tabs",
				Class:     "border-b",
				AriaLabel: "Settings tabs",
			},
			Variant:     TabsPills,
			ClientSide:  true,
			ActiveTabID: "tab2",
			Tabs: []Tab{
				{ID: "tab1", Label: "General"},
				{ID: "tab2", Label: "Security"},
			},
		}))
		utils.AssertContains(t, output, `data-tc-tabs`)
		utils.AssertContains(t, output, "Security")
	})
}

func TestBadgeFullCoverage(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		btype BadgeType
	}{
		{"primary", BadgePrimary},
		{"success", BadgeSuccess},
		{"error", BadgeError},
		{"warning", BadgeWarning},
		{"info", BadgeInfo},
		{"neutral", BadgeNeutral},
	} {
		t.Run(tt.name+" badge", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Badge(BadgeProps{
				BaseProps: utils.BaseProps{ID: "b-" + tt.name},
				Text:      tt.name,
				Type:      tt.btype,
				Size:      BadgeSizeLG,
				Pill:      true,
				Dot:       true,
			}))
			utils.AssertContains(t, output, tt.name)
		})
	}

	t.Run("badge with href", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "New",
			Href: "/new",
		}))
		utils.AssertContains(t, output, `href="/new"`)
	})
}

func TestButtonFullCoverage(t *testing.T) {
	t.Parallel()

	for _, variant := range []ButtonType{ButtonPrimary, ButtonSecondary, ButtonDanger, ButtonGhost, ButtonLink} {
		t.Run(string(variant)+" variant", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Button(ButtonProps{
				BaseProps: utils.BaseProps{
					ID:        "btn-" + string(variant),
					Class:     "extra",
					AriaLabel: string(variant) + " button",
				},
				Text:    string(variant),
				Variant: variant,
				Size:    ButtonSizeLG,
			}))
			utils.AssertContains(t, output, string(variant))
		})
	}

	t.Run("disabled button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			Text:     "Disabled",
			Disabled: true,
		}))
		utils.AssertContains(t, output, `disabled`)
	})
}

func TestCardFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with all options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			BaseProps: utils.BaseProps{
				ID:        "card-1",
				Class:     "shadow-lg",
				AriaLabel: "User card",
			},
			Title:    "User Profile",
			Subtitle: "Admin account",
			Padding:  CardPaddingLG,
		}))
		utils.AssertContains(t, output, "User Profile")
		utils.AssertContains(t, output, "Admin account")
		utils.AssertContains(t, output, `id="card-1"`)
	})
}

func TestStatCardFullCoverage(t *testing.T) {
	t.Parallel()

	for _, trend := range []TrendDirection{TrendUp, TrendDown, TrendNone} {
		t.Run("trend_"+string(trend), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, StatCard(StatCardProps{
				BaseProps: utils.BaseProps{ID: "stat-" + string(trend)},
				Value:     "$1,234",
				Label:     "Revenue",
				Change:    "+12%",
				Trend:     trend,
			}))
			utils.AssertContains(t, output, "$1,234")
			utils.AssertContains(t, output, "Revenue")
		})
	}
}

func TestAvatarFullCoverage(t *testing.T) {
	t.Parallel()

	for _, size := range []AvatarSize{AvatarSizeXS, AvatarSizeSM, AvatarSizeMD, AvatarSizeLG, AvatarSizeXL} {
		t.Run("size_"+string(size), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Avatar(AvatarProps{
				Initials: "AB",
				Size:     size,
			}))
			utils.AssertContains(t, output, "AB")
		})
	}

	t.Run("circle with online status", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:    "/img.jpg",
			Alt:    "User",
			Shape:  AvatarShapeCircle,
			Status: AvatarStatusOnline,
		}))
		utils.AssertContains(t, output, "/img.jpg")
		utils.AssertContains(t, output, "bg-green-400")
	})

	t.Run("square with offline status", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "CD",
			Shape:    AvatarShapeSquare,
			Status:   AvatarStatusOffline,
		}))
		utils.AssertContains(t, output, "CD")
	})
}

func TestTooltipFullCoverage(t *testing.T) {
	t.Parallel()

	for _, pos := range []TooltipPosition{TooltipPositionTop, TooltipPositionBottom, TooltipPositionLeft, TooltipPositionRight} {
		t.Run("position_"+string(pos), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Tooltip(TooltipProps{
				BaseProps: utils.BaseProps{ID: "tip-" + string(pos)},
				Text:      "Helpful tip",
				Position:  pos,
			}))
			utils.AssertContains(t, output, "Helpful tip")
		})
	}
}
