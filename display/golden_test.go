package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenPageHeader(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, PageHeader(PageHeaderProps{
		Title:    "Users",
		Subtitle: "Manage user accounts",
		Action:   templ.Raw(`<a href="/users/new">New user</a>`),
	}))
	golden.Assert(t, "page_header", output)
}

func TestGoldenDefinitionList(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DefinitionList(DefinitionListProps{
		Items: []DefinitionItem{
			{Term: "Email", Detail: "alice@example.com"},
			{Term: "Plan", Detail: "Pro"},
		},
	}))
	golden.Assert(t, "definition_list", output)
}

func TestGoldenListNote(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ListNote(ListNoteProps{Shown: 50, Total: 127}))
	golden.Assert(t, "list_note", output)
}

func TestGoldenStatCardWithIcon(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StatCard(StatCardProps{
		Label:  "Users",
		Value:  "1,204",
		Icon:   icons.Users,
		Change: "12%",
		Trend:  TrendUp,
	}))
	golden.Assert(t, "stat_card_icon", output)
}
