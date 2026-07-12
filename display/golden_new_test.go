package display

import (
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenCopyButton(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CopyButton(CopyButtonProps{
		Text:      "npm install foo",
		Label:     "Copy",
		Icon:      true,
		BaseProps: utils.BaseProps{Nonce: "abc123"},
	}))
	golden.Assert(t, "copy_button", output)
}

func TestGoldenRelativeTime(t *testing.T) {
	t.Parallel()
	ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	output := utils.Render(t, RelativeTime(RelativeTimeProps{
		Time: ts,
	}))
	golden.Assert(t, "relative_time", output)
}

func TestGoldenCountBadge(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CountBadge(CountBadgeProps{Count: 5}))
	golden.Assert(t, "count_badge", output)
}

func TestGoldenDefinitionGrid(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
		Cols: GridCols2,
		Items: []DefinitionItem{
			{Term: "CPU", Detail: "42%"},
			{Term: "Memory", Detail: "8.2 GB"},
		},
	}))
	golden.Assert(t, "definition_grid", output)
}

func TestGoldenImage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Image(ImageProps{
		Src:    "/photo.jpg",
		Alt:    "Profile photo",
		Width:  128,
		Height: 128,
	}))
	golden.Assert(t, "image", output)
}

func TestGoldenStatCardHTMX(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StatCard(StatCardProps{
		Value:    "1,204",
		Label:    "Active Users",
		Change:   "+8%",
		Trend:    TrendUp,
		Icon:     icons.Users,
		HxGet:    "/api/stats",
		HxTarget: "#stat-container",
		HxSwap:   htmx.SwapInnerHTML,
	}))
	golden.Assert(t, "stat_card_htmx", output)
}

func TestGoldenStatCardHrefWithHTMX(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, StatCard(StatCardProps{
		Value:  "42",
		Label:  "Pending",
		Href:   "/admin/pending",
		HxGet:  "/api/pending",
		HxSwap: htmx.SwapOuterHTML,
	}))
	golden.Assert(t, "stat_card_href_htmx", output)
}

func TestGoldenCardBodySlot(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Card(CardProps{
		Title: "Installation",
		Body:  templ.Raw("<pre>npm install @larsartmann/templ-components</pre>"),
	}))
	golden.Assert(t, "card_body_slot", output)
}

func TestGoldenCardHeaderSlot(t *testing.T) {
	t.Parallel()
	customHeader := `<div class="flex items-center gap-2">` +
		`<h2 class="text-lg font-bold">Custom Header</h2>` +
		`<span class="badge">New</span></div>`
	output := utils.Render(t, Card(CardProps{
		Header: templ.Raw(customHeader),
		Body:   templ.Raw("<p>Card body content</p>"),
	}))
	golden.Assert(t, "card_header_slot", output)
}

func TestGoldenGridAutoFit(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:        GridColsAutoFit,
		MinColWidth: "200px",
	}))
	golden.Assert(t, "grid_autofit", output)
}

func TestGoldenTableSortableHeaders(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Table(TableProps{
		TypedHeaders: []TableHeader{
			{Label: "Name", Sortable: true, SortDirection: SortAsc, Href: "/sort?col=name&dir=desc"},
			{Label: "Joined", Sortable: true, SortDirection: SortNone, Href: "/sort?col=joined&dir=asc"},
			{Label: "Role"},
		},
		Rows: []TableRow{
			SimpleTableRow("Alice", "2024-01-01", "Admin"),
			SimpleTableRow("Bob", "2024-02-01", "User"),
		},
	}))
	golden.Assert(t, "table_sortable_headers", output)
}

func TestGoldenDataTable(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name", Sortable: true},
			{Label: "Email"},
		},
		Rows: []TableRow{
			SimpleTableRow("Alice", "alice@example.com"),
			SimpleTableRow("Bob", "bob@example.com"),
		},
		Striped: true,
	}))
	golden.Assert(t, "datatable_basic", output)
}

func TestGoldenDataTableSortable(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name", Sortable: true},
			{Label: "Email", Sortable: true, SortKey: "email_address"},
		},
		ActiveSortColumn: "Name",
		ActiveSortDir:    SortAsc,
		SortBaseURL:      "/users",
		Rows: []TableRow{
			SimpleTableRow("Alice", "alice@example.com"),
		},
	}))
	golden.Assert(t, "datatable_sortable", output)
}
