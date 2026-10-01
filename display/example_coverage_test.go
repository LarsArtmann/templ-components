package display_test

import (
	"github.com/larsartmann/templ-components/display"
)

func ExampleAccordion() {
	_ = display.Accordion(display.DefaultAccordionProps())
	// Output:
}

func ExampleAreaChart() {
	_ = display.AreaChart(display.DefaultAreaChartProps())
	// Output:
}

func ExampleAvatar() {
	_ = display.Avatar(display.DefaultAvatarProps())
	// Output:
}

func ExampleStatusBadge() {
	_ = display.StatusBadge("active")
	// Output:
}

func ExampleBarChart() {
	_ = display.BarChart(display.DefaultBarChartProps())
	// Output:
}

func ExampleButton() {
	_ = display.Button(display.DefaultButtonProps())
	// Output:
}

func ExampleSimpleCard() {
	_ = display.SimpleCard(display.DefaultSimpleCardProps())
	// Output:
}

func ExampleCollapsibleSection() {
	_ = display.CollapsibleSection(display.DefaultCollapsibleSectionProps())
	// Output:
}

func ExampleDateRange() {
	_ = display.DateRange(display.DateRangeProps{PresentText: "Present"})
	// Output:
}

func ExampleDefinitionList() {
	_ = display.DefinitionList(display.DefaultDefinitionListProps())
	// Output:
}

func ExampleDrawer() {
	_ = display.Drawer(display.DefaultDrawerProps())
	// Output:
}

func ExampleDropdown() {
	_ = display.Dropdown(display.DefaultDropdownProps())
	// Output:
}

func ExampleEmptyState() {
	_ = display.EmptyState(display.DefaultEmptyStateProps())
	// Output:
}

func ExampleSimpleEmptyState() {
	_ = display.SimpleEmptyState("No results found")
	// Output:
}

func ExampleExternalLink() {
	_ = display.ExternalLink(display.DefaultExternalLinkProps())
	// Output:
}

func ExampleHeatmap() {
	_ = display.Heatmap(display.DefaultHeatmapProps())
	// Output:
}

func ExampleLineChart() {
	_ = display.LineChart(display.DefaultLineChartProps())
	// Output:
}

func ExampleListNote() {
	_ = display.ListNote(display.ListNoteProps{Shown: 10, Total: 100})
	// Output:
}

func ExampleModal() {
	_ = display.Modal(display.DefaultModalProps())
	// Output:
}

func ExamplePageHeader() {
	_ = display.PageHeader(display.DefaultPageHeaderProps())
	// Output:
}

func ExamplePieChart() {
	_ = display.PieChart(display.DefaultPieChartProps())
	// Output:
}

func ExamplePopover() {
	_ = display.Popover(display.DefaultPopoverProps())
	// Output:
}

func ExampleSectionHeading() {
	_ = display.SectionHeading(display.SectionHeadingProps{Title: "Overview"})
	// Output:
}

func ExampleSparkline() {
	_ = display.Sparkline(display.DefaultSparklineProps())
	// Output:
}

func ExampleTabs() {
	_ = display.Tabs(display.DefaultTabsProps())
	// Output:
}

func ExampleTooltip() {
	_ = display.Tooltip(display.DefaultTooltipProps())
	// Output:
}
