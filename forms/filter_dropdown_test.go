package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
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

func TestFilterDropdownWire(t *testing.T) {
	t.Parallel()

	t.Run("wired htmx renders debounce-free change trigger and target", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:    "status",
			Options: testFilterDropdownOptions(),
			Wire: &wire.Action{
				URL:    "/api/filter",
				Target: "#results",
			},
		}))
		utils.AssertContains(t, output, `hx-get="/api/filter"`)
		utils.AssertContains(t, output, `hx-trigger="change"`)
		utils.AssertContains(t, output, `hx-target="#results"`)
	})

	t.Run("wired datastar renders form-encoded change expression", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:    "status",
			Options: testFilterDropdownOptions(),
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/filter",
			},
		}))
		utils.AssertContains(t, output, `data-on:change="@get(&#39;/api/filter&#39;, {contentType: &#39;form&#39;})"`)
	})

	t.Run("wire owns the wiring: legacy HxGet/HxTarget/HxTrigger are ignored", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:      "status",
			Options:   testFilterDropdownOptions(),
			HxGet:     "/legacy",
			HxTarget:  "#legacy",
			HxTrigger: "change delay:99ms",
			Wire: &wire.Action{
				URL:    "/api/filter",
				Target: "#results",
			},
		}))
		utils.AssertContains(t, output, `hx-get="/api/filter"`)
		utils.AssertContains(t, output, `hx-target="#results"`)
		utils.AssertNotContains(t, output, "/legacy")
		utils.AssertNotContains(t, output, "delay:99ms")
	})

	t.Run("htmx extras HxInclude/HxIndicator still render when wired", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:        "status",
			Options:     testFilterDropdownOptions(),
			HxInclude:   "closest form",
			HxIndicator: "#spinner",
			Wire: &wire.Action{
				URL: "/api/filter",
			},
		}))
		utils.AssertContains(t, output, `hx-include="closest form"`)
		utils.AssertContains(t, output, `hx-indicator="#spinner"`)
	})

	t.Run("wired select wraps in a GET form with noscript Apply button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:    "status",
			Options: testFilterDropdownOptions(),
			Action:  "/users",
			Wire: &wire.Action{
				URL: "/api/filter",
			},
		}))
		utils.AssertContains(t, output, `action="/users"`)
		utils.AssertContains(t, output, `method="GET"`)
		utils.AssertContains(t, output, "<noscript>")
		utils.AssertContains(t, output, `type="submit"`)
		utils.AssertContains(t, output, "Apply")
	})

	t.Run("empty URL wires nothing (legacy rendering)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:    "status",
			Options: testFilterDropdownOptions(),
			HxGet:   "/api/filter",
			Wire: &wire.Action{
				URL: "",
			},
		}))
		utils.AssertContains(t, output, `hx-get="/api/filter"`)
		utils.AssertNotContains(t, output, "<noscript>")
		utils.AssertNotContains(t, output, "<form")
	})
}

// testFilterDropdownOptions is the shared option set for FilterDropdown wire
// tests.
func testFilterDropdownOptions() []SelectOption {
	return []SelectOption{
		{Value: "", Label: "All"},
		{Value: "active", Label: "Active"},
	}
}
