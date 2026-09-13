package contract

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/charts/echarts"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/errorpage"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// propsIDContract is the id every case sets; its presence in the rendered
// output is the contract.
const propsIDContract = `id="tc-id-contract"`

// TestPropsIDRendersInOutput is convention check 4: a consumer-set
// BaseProps.ID must appear in the rendered output of every component. The
// forms.Calendar bug (props.ID silently dropped on the root — found by the
// MonthNav e2e 2026-09-13) proved this failure mode is real: the ID anchors
// consumer CSS, ARIA wiring, AND transport targeting (hx-target="#my-id");
// dropping it silently breaks all three.
//
// The contract checked is "the id appears in the output", not "on the first
// element" — sub-templates legitimately host the root elsewhere. Components
// that cannot host an id by design (external-resource helpers, non-Props
// signatures) are simply absent from the table; scope gaps are caught by
// review, not by an exemption list that rots.
func TestPropsIDRendersInOutput(t *testing.T) {
	t.Parallel()

	table := []struct {
		name   string
		render func() string
	}{
		{"datastar.SDKScript", renderPropsID(t, datastar.SDKScript)},
		{"datastar.LiveRegion", renderPropsID(t, datastar.LiveRegion)},
		{"datastar.Indicator", renderPropsID(t, datastar.Indicator)},

		{"display.Badge", renderPropsID(t, display.Badge)},
		{"display.Avatar", renderPropsID(t, display.Avatar)},
		{"display.Eyebrow", renderPropsID(t, display.Eyebrow)},
		{"display.Scrollback", renderPropsID(t, display.Scrollback)},
		{"display.Tooltip", renderPropsID(t, display.Tooltip)},
		{"display.Accordion", renderPropsID(t, display.Accordion)},
		{"display.Button", renderPropsID(t, display.Button)},
		{"display.Card", renderPropsID(t, display.Card)},
		{"display.SimpleCard", renderPropsID(t, display.SimpleCard)},
		{"display.StatCard", renderPropsID(t, display.StatCard)},
		{"display.Dropdown", renderPropsID(t, display.Dropdown)},
		{"display.Tabs", renderPropsID(t, display.Tabs)},
		{"display.Table", renderPropsID(t, display.Table)},
		{"display.DataTable", renderPropsID(t, display.DataTable)},
		{"display.PageHeader", renderPropsID(t, display.PageHeader)},
		{"display.ListNote", renderPropsID(t, display.ListNote)},
		{"display.EmptyState", renderPropsID(t, display.EmptyState)},
		{"display.DefinitionList", renderPropsID(t, display.DefinitionList)},
		{"display.DefinitionGrid", renderPropsID(t, display.DefinitionGrid)},
		{"display.Modal", renderPropsID(t, display.Modal)},
		{"display.Drawer", renderPropsID(t, display.Drawer)},
		{"display.Grid", renderPropsID(t, display.Grid)},
		{"display.CopyButton", renderPropsID(t, display.CopyButton)},
		{"display.RelativeTime", renderPropsID(t, display.RelativeTime)},
		{"display.CountBadge", renderPropsID(t, display.CountBadge)},
		{"display.Image", renderPropsID(t, display.Image)},
		{"display.HoverCard", renderPropsID(t, display.HoverCard)},
		{"display.ContextMenu", renderPropsID(t, display.ContextMenu)},
		{"display.Carousel", renderPropsID(t, display.Carousel)},
		{"display.Sparkline", renderPropsID(t, display.Sparkline)},
		{"display.BarChart", renderPropsID(t, display.BarChart)},
		{"display.LineChart", renderPropsID(t, display.LineChart)},
		{"display.PieChart", renderPropsID(t, display.PieChart)},
		{"display.AreaChart", renderPropsID(t, display.AreaChart)},
		{"display.Heatmap", renderPropsID(t, display.Heatmap)},
		{"display.ExternalLink", renderPropsID(t, display.ExternalLink)},
		{"display.CollapsibleSection", renderPropsID(t, display.CollapsibleSection)},
		{"display.KanbanBoard", renderPropsID(t, display.KanbanBoard)},

		{"feedback.Alert", renderPropsID(t, feedback.Alert)},
		{"feedback.Toast", renderPropsID(t, feedback.Toast)},
		{"feedback.Spinner", renderPropsID(t, feedback.Spinner)},
		{"feedback.LoadingOverlay", renderPropsID(t, feedback.LoadingOverlay)},
		{"feedback.ProgressBar", renderPropsID(t, feedback.ProgressBar)},
		{"feedback.StepIndicator", renderPropsID(t, feedback.StepIndicator)},

		{"forms.Input", renderPropsID(t, forms.Input)},
		{"forms.Checkbox", renderPropsID(t, forms.Checkbox)},
		{"forms.Select", renderPropsID(t, forms.Select)},
		{"forms.Textarea", renderPropsID(t, forms.Textarea)},
		{"forms.Toggle", renderPropsID(t, forms.Toggle)},
		{"forms.Radio", renderPropsID(t, forms.Radio)},
		{"forms.RadioGroup", renderPropsID(t, forms.RadioGroup)},
		{"forms.Combobox", renderPropsID(t, forms.Combobox)},
		{"forms.DatePicker", renderPropsID(t, forms.DatePicker)},
		{"forms.FileInput", renderPropsID(t, forms.FileInput)},
		{"forms.Form", renderPropsID(t, forms.Form)},
		{"forms.InputGroup", renderPropsID(t, forms.InputGroup)},
		{"forms.ValidationSummary", renderPropsID(t, forms.ValidationSummary)},
		{"forms.FilterDropdown", renderPropsID(t, forms.FilterDropdown)},
		{"forms.FilterInput", renderPropsID(t, forms.FilterInput)},
		{"forms.Slider", renderPropsID(t, forms.Slider)},
		{"forms.Rating", renderPropsID(t, forms.Rating)},
		{"forms.TagsInput", renderPropsID(t, forms.TagsInput)},
		{"forms.Calendar", renderPropsID(t, forms.Calendar)},

		{"layout.ThemeToggle", func() string {
			return utils.Render(t, layout.ThemeToggle("toggle", ""))
		}},
		{"layout.AppShell", renderPropsID(t, layout.AppShell)},
		{"layout.Container", renderPropsID(t, layout.Container)},
		{"layout.Split", renderPropsID(t, layout.Split)},
		{"layout.Stack", renderPropsID(t, layout.Stack)},

		{"navigation.Nav", renderPropsID(t, navigation.Nav)},
		{"navigation.SimpleNav", renderPropsID(t, navigation.SimpleNav)},
		{"navigation.NavLink", func() string {
			props := navigation.DefaultNavLinkProps()
			props.BaseProps = utils.BaseProps{ID: "tc-id-contract"}

			return utils.Render(t, navigation.NavLink(props, "/"))
		}},
		{"navigation.Pagination", renderPropsID(t, navigation.Pagination)},
		{"navigation.Breadcrumbs", renderPropsID(t, navigation.Breadcrumbs)},
		{"navigation.SidebarNav", renderPropsID(t, navigation.SidebarNav)},
		{"navigation.LoadMore", renderPropsID(t, navigation.LoadMore)},
		{"navigation.EndOfList", renderPropsID(t, navigation.EndOfList)},
		{"navigation.Footer", renderPropsID(t, navigation.Footer)},

		{"errorpage.ErrorPage", renderPropsID(t, errorpage.ErrorPage)},
		{"errorpage.NotFound404", renderPropsID(t, errorpage.NotFound404)},
		{"errorpage.ErrorDetail", renderPropsID(t, errorpage.ErrorDetail)},
		{"errorpage.ErrorAlert", renderPropsID(t, errorpage.ErrorAlert)},

		{"echarts.EChart", renderPropsID(t, echarts.EChart)},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			html := tc.render()

			if !strings.Contains(html, propsIDContract) {
				t.Errorf(
					"%s: consumer-set ID missing from output — the props.ID contract is broken (first 300 bytes: %.300s)",
					tc.name,
					html,
				)
			}
		})
	}
}

// renderPropsID sets the contract ID on a single-Props constructor's props
// and renders it. SetBaseProps has a pointer receiver (recvcheck), so the
// constraint goes through PT = *P.
func renderPropsID[P any, PT interface {
	*P
	utils.ComponentProps
}](t *testing.T, ctor func(P) templ.Component) func() string {
	t.Helper()

	return func() string {
		var props P

		PT(&props).SetBaseProps(utils.BaseProps{ID: "tc-id-contract"})

		return utils.Render(t, ctor(props))
	}
}
