package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultFilterDropdownProps(t *testing.T) {
	t.Parallel()

	_ = DefaultFilterDropdownProps()
}

func TestFilterDropdownBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:  "status",
		Label: "Status",
		Options: []SelectOption{
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
		},
		HxGet:    "/api/users",
		HxTarget: "#user-list",
	}))
	utils.AssertContains(t, output, "<select")
	utils.AssertContains(t, output, `name="status"`)
	utils.AssertContains(t, output, "Active")
	utils.AssertContains(t, output, "Inactive")
	utils.AssertContains(t, output, `hx-get="/api/users"`)
	utils.AssertContains(t, output, `hx-target="#user-list"`)
	utils.AssertContains(t, output, `hx-trigger="change"`)
}

func TestFilterDropdownPreselectValue(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:  "status",
		Value: "active",
		Options: []SelectOption{
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
		},
		HxGet:    "/api/users",
		HxTarget: "#user-list",
	}))
	// The "active" option should appear in the output with its value
	utils.AssertContains(t, output, `value="active"`)
}

func TestFilterDropdownCustomTrigger(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:      "search",
		HxGet:     "/api/search",
		HxTarget:  "#results",
		HxTrigger: "change delay:500ms",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
		},
	}))
	utils.AssertContains(t, output, `hx-trigger="change delay:500ms"`)
}

func TestFilterDropdownHxSwap(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:     "sort",
		HxGet:    "/api/items",
		HxTarget: "#items",
		HxSwap:   "outerHTML",
		Options: []SelectOption{
			{Value: "name", Label: "Name"},
		},
	}))
	utils.AssertContains(t, output, `hx-swap="outerHTML"`)
}

func TestFilterDropdownHxInclude(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:      "status",
		HxGet:     "/api/users",
		HxTarget:  "#user-list",
		HxInclude: "closest form",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
		},
	}))
	utils.AssertContains(t, output, `hx-include="closest form"`)
}

func TestFilterDropdownHxIndicator(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:        "status",
		HxGet:       "/api/users",
		HxTarget:    "#user-list",
		HxIndicator: "#loading",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
		},
	}))
	utils.AssertContains(t, output, `hx-indicator="#loading"`)
}

func TestFilterDropdownHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:     "filter",
		Label:    "Filter",
		HelpText: "Narrow results by category",
		HxGet:    "/api/items",
		HxTarget: "#items",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
		},
	}))
	utils.AssertContains(t, output, "Narrow results by category")
}

func TestFilterDropdownNoLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:     "status",
		HxGet:    "/api/users",
		HxTarget: "#user-list",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
		},
	}))
	utils.AssertNotContains(t, output, "<label")
}

func TestFilterDropdownPreselectEmpty(t *testing.T) {
	t.Parallel()

	opts := []SelectOption{
		{Value: "a", Label: "A"},
		{Value: "b", Label: "B", Selected: true},
	}
	result := filterDropdownPreselect(opts, "")
	// When value is empty, opts returned unchanged
	if result[1].Selected != true {
		t.Error("expected original Selected=true to be preserved when value is empty")
	}
}

func TestFilterDropdownPreselectMatch(t *testing.T) {
	t.Parallel()

	opts := []SelectOption{
		{Value: "a", Label: "A"},
		{Value: "b", Label: "B", Selected: true},
	}

	result := filterDropdownPreselect(opts, "a")
	if !result[0].Selected {
		t.Error("expected option 'a' to be Selected")
	}

	if result[1].Selected {
		t.Error("expected option 'b' to NOT be Selected")
	}
}

func TestFilterDropdownDarkModeCompliance(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:     "status",
		HxGet:    "/api/users",
		HxTarget: "#user-list",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
		},
	}))
	// Select input should have dark mode classes
	utils.AssertContains(t, output, "dark:bg-gray-800")
	utils.AssertContains(t, output, "dark:text-white")
}
