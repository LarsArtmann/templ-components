package contract

import (
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// TestRenderDeterminism pins the library's render-determinism guarantee:
// rendering the same props twice must produce byte-identical HTML. Map
// iteration order, time reads, or rand-derived values leaking into output
// would silently reshuffle class strings or content between renders — the
// golden files only pin ONE render, so this test is the only guard of the
// second.
//
// Scope: components whose props fully determine the output (explicit IDs
// where IDs exist — EnsureID randomness is normalized away for goldens and
// is deliberately out of scope here; rendering those twice differs by
// design). Every family with a class-lookup map, a chart (geometry math),
// or a table (row/column alignment) is represented.
func TestRenderDeterminism(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		render func() string
	}{
		{"badge lookup map", func() string {
			props := display.DefaultBadgeProps()
			props.Text = "beta"
			props.Tone = display.BadgeToneWarning

			return utils.Render(t, display.Badge(props))
		}},
		{"card padding map", func() string {
			props := display.DefaultCardProps()
			props.Title = "Users"
			props.Padding = display.CardPaddingLG

			return utils.Render(t, display.Card(props))
		}},
		{"table typed headers + row href", func() string {
			return utils.Render(t, display.Table(display.TableProps{
				TypedHeaders: []display.TableHeader{
					{Label: "Name", Sortable: true, SortDirection: display.SortAsc},
					{Label: "Email"},
				},
				Rows: []display.TableRow{
					display.SimpleTableRow("Alice", "alice@example.com"),
					display.SimpleTableRow("Bob", "bob@example.com"),
				},
			}))
		}},
		{"line chart smooth path", func() string {
			props := display.DefaultLineChartProps()
			props.Series = []display.LineChartSeries{
				{Label: "a", Values: []float64{1, 4, 2, 8, 3}},
				{Label: "b", Values: []float64{5, 3, 7, 2, 6}},
			}

			return utils.Render(t, display.LineChart(props))
		}},
		{"pie chart arcs", func() string {
			props := display.DefaultPieChartProps()
			props.Slices = []display.PieChartSlice{
				{Label: "a", Value: 12, Color: display.ChartColorEmerald},
				{Label: "b", Value: 7},
			}

			return utils.Render(t, display.PieChart(props))
		}},
		{"stat card", func() string {
			return utils.Render(t, display.StatCard(display.StatCardProps{
				Label: "Revenue", Value: "$42", Trend: display.TrendUp, Delta: "+12%",
			}))
		}},
		{"alert feedback style map", func() string {
			return utils.Render(t, feedback.Alert(feedback.AlertProps{
				Type: feedback.FeedbackWarning,
				Text: "careful",
			}))
		}},
		{"progress bar clamp", func() string {
			return utils.Render(t, feedback.ProgressBar(feedback.ProgressBarProps{
				Current: 45, Total: 100, Label: "uploading",
			}))
		}},
		{"input with error attrs", func() string {
			return utils.Render(t, forms.Input(forms.InputProps{
				Name: "email", Label: "Email", Type: forms.InputEmail,
				Error: "bad", Help: "enter work address",
			}))
		}},
		{"select options", func() string {
			return utils.Render(t, forms.Select(forms.SelectProps{
				Name: "role", Label: "Role",
				Options: []forms.SelectOption{{Value: "admin", Text: "Admin"}, {Value: "user", Text: "User"}},
			}))
		}},
		{"pagination", func() string {
			return utils.Render(t, navigation.Pagination(navigation.PaginationProps{
				CurrentPage: 3, TotalPages: 9, Href: "/page",
			}))
		}},
		{"breadcrumbs", func() string {
			return utils.Render(t, navigation.Breadcrumbs(navigation.BreadcrumbsProps{
				Items: []navigation.Breadcrumb{
					{Text: "Home", Href: "/"},
					{Text: "Settings", Href: "/settings"},
					{Text: "Profile"},
				},
			}))
		}},
		{"theme toggle", func() string {
			return utils.Render(t, layout.ThemeToggle(layout.ThemeToggleProps{Nonce: "n"}))
		}},
		{"relative time", func() string {
			props := display.RelativeTimeProps{Time: mustParseRFC3339("2025-01-15T10:30:00Z")}
			props.Now = mustParseRFC3339("2025-01-15T11:00:00Z")

			return utils.Render(t, display.RelativeTime(props))
		}},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			first, second := tc.render(), tc.render()
			if first != second {
				t.Errorf(
					"render is not deterministic — output differs between two renders of identical props.\nfirst:\n%s\n\nsecond:\n%s",
					first,
					second,
				)
			}
		})
	}
}
