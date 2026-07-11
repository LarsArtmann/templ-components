package display_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/display"
)

func ExampleBadge() {
	props := display.DefaultBadgeProps()
	props.Text = "Beta"
	props.Type = display.BadgeInfo

	var buf bytes.Buffer
	_ = display.Badge(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
	// Output will contain the badge text and Tailwind classes
}

func ExampleCard() {
	props := display.DefaultCardProps()
	props.Title = "Hello World"

	var buf bytes.Buffer
	_ = display.Card(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleStatCard() {
	props := display.DefaultStatCardProps()
	props.Value = "$12,345"
	props.Label = "Total Revenue"
	props.Trend = display.TrendUp
	props.Change = "+12.5%"

	var buf bytes.Buffer
	_ = display.StatCard(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleGrid() {
	props := display.DefaultGridProps()

	var buf bytes.Buffer
	_ = display.Grid(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
	// Output will contain responsive grid classes
}

func ExampleTable() {
	props := display.TableProps{
		Headers: []string{"Name", "Email", "Role"},
		Rows: []display.TableRow{
			display.SimpleTableRow("Alice", "alice@example.com", "Admin"),
			display.SimpleTableRow("Bob", "bob@example.com", "User"),
		},
	}
	var buf bytes.Buffer
	_ = display.Table(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleTable_flushInCard() {
	// When nesting a Table inside a Card(CardPaddingNone), set Flush to true
	// to suppress the table's own border and avoid a double-border defect.
	// Use CellPadding: TableCellPaddingCompact for data-heavy dashboards.
	props := display.TableProps{
		Headers:     []string{"Name", "Status"},
		Rows:        []display.TableRow{display.SimpleTableRow("Alice", "Active")},
		Flush:       true,
		CellPadding: display.TableCellPaddingCompact,
	}
	var buf bytes.Buffer
	_ = display.Table(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}
