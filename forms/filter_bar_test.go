package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestFilterBarRender(t *testing.T) {
	t.Parallel()

	t.Run("default htmx wiring derived from action", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterBar(DefaultFilterBarProps())) // Action empty → inert
		utils.AssertNotContains(t, output, "hx-get")
	})

	t.Run("action wires the bar", func(t *testing.T) {
		t.Parallel()

		props := DefaultFilterBarProps()
		props.Action = "/users"
		props.Target = "#results"
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, `hx-get="/users"`)
		utils.AssertContains(t, output, `hx-target="#results"`)
		utils.AssertContains(t, output, `hx-trigger="change from:find select, change from:find input[type=checkbox]"`)
		utils.AssertContains(t, output, `hx-push-url="true"`)
		utils.AssertContains(t, output, `action="/users"`)
		utils.AssertContains(t, output, `method="get"`)
	})

	t.Run("composite trigger encodes the checkbox footgun", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{Action: "/list"}
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, "change from:find select")
		utils.AssertContains(t, output, "change from:find input[type=checkbox]")
	})

	t.Run("no push url opt-out", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{Action: "/users", NoPushURL: true}
		output := utils.Render(t, FilterBar(props))
		utils.AssertNotContains(t, output, "hx-push-url")
	})

	t.Run("reset link renders when reset text and action set", func(t *testing.T) {
		t.Parallel()

		props := DefaultFilterBarProps()
		props.Action = "/users"
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, `href="/users"`)
		utils.AssertContains(t, output, "Reset")
	})

	t.Run("no reset without action", func(t *testing.T) {
		t.Parallel()

		props := DefaultFilterBarProps()
		output := utils.Render(t, FilterBar(props))
		utils.AssertNotContains(t, output, "Reset")
	})

	t.Run("sticky shell classes", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{Action: "/users", Sticky: true}
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, "sticky")
		utils.AssertContains(t, output, "top-0")
		utils.AssertContains(t, output, "z-20")
	})

	t.Run("indicator only when wired htmx", func(t *testing.T) {
		t.Parallel()

		wired := FilterBarProps{Action: "/users"}
		output := utils.Render(t, FilterBar(wired))
		utils.AssertContains(t, output, "htmx-indicator")
		utils.AssertContains(t, output, "hx-indicator")

		inert := FilterBarProps{}
		output = utils.Render(t, FilterBar(inert))
		utils.AssertNotContains(t, output, "htmx-indicator")
	})

	t.Run("noscript apply button always renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterBar(FilterBarProps{}))
		utils.AssertContains(t, output, "<noscript>")
		utils.AssertContains(t, output, ">Apply</button>")
	})

	t.Run("propagates id class aria-label and attrs", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{
			ID:        "bar",
			Class:     "mb-4",
			AriaLabel: "Filters",
			Attrs:     templ.Attributes{"data-foo": "bar"},
			Action:    "/users",
		}
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, `id="bar"`)
		utils.AssertContains(t, output, "mb-4")
		utils.AssertContains(t, output, `aria-label="Filters"`)
		utils.AssertContains(t, output, `data-foo="bar"`)
	})

	t.Run("children render inside the form", func(t *testing.T) {
		t.Parallel()
		output := renderFilterBarWithChildren(t, FilterBarProps{Action: "/users"}, FilterDropdown(FilterDropdownProps{
			Name:    "status",
			Label:   "Status",
			Options: []SelectOption{{Value: "", Label: "All"}},
		}))
		utils.AssertContains(t, output, `name="status"`)
	})
}

func TestFilterBarWiring(t *testing.T) {
	t.Parallel()

	t.Run("wire htmx keeps composite trigger and form encoding", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{
			Wire:   &wire.Action{Method: wire.MethodGet, URL: "/api/users", Target: "#results"},
			Target: "#results",
		}
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, `hx-get="/api/users"`)
		utils.AssertContains(t, output, `hx-target="#results"`)
		utils.AssertContains(t, output, "change from:find select")
		utils.AssertContains(t, output, `hx-push-url="true"`)
	})

	t.Run("wire datastar renders change event with form content type and no htmx attrs", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{
			Wire: &wire.Action{Transport: wire.TransportDatastar, Method: wire.MethodGet, URL: "/api/users"},
		}
		output := utils.Render(t, FilterBar(props))
		utils.AssertContains(t, output, `data-on:change="@get(&#39;/api/users&#39;, {contentType: &#39;form&#39;})"`)
		utils.AssertNotContains(t, output, "hx-get")
		utils.AssertNotContains(t, output, "hx-push-url")
		utils.AssertNotContains(t, output, "htmx-indicator")
	})

	t.Run("wire with empty url stays inert", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{Wire: &wire.Action{}}
		output := utils.Render(t, FilterBar(props))
		utils.AssertNotContains(t, output, "hx-get")
		utils.AssertNotContains(t, output, "data-on:change")
	})

	t.Run("wire no push url opt-out", func(t *testing.T) {
		t.Parallel()

		props := FilterBarProps{
			Wire:      &wire.Action{Method: wire.MethodGet, URL: "/api/users"},
			NoPushURL: true,
		}
		output := utils.Render(t, FilterBar(props))
		utils.AssertNotContains(t, output, "hx-push-url")
	})

	t.Run("consumer wire action is never mutated", func(t *testing.T) {
		t.Parallel()

		action := &wire.Action{Method: wire.MethodGet, URL: "/api/users"}
		props := FilterBarProps{Wire: action}

		_ = utils.Render(t, FilterBar(props))
		if action.Event != wire.EventUnspecified {
			t.Errorf("consumer action Event mutated to %q", action.Event)
		}

		if action.ContentType != wire.ContentTypeUnspecified {
			t.Errorf("consumer action ContentType mutated to %q", action.ContentType)
		}
	})
}

func TestFilterBarIsHTMX(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props FilterBarProps
		want  bool
	}{
		{"no wire defaults to htmx", FilterBarProps{}, true},
		{"wire htmx", FilterBarProps{Wire: &wire.Action{URL: "/x"}}, true},
		{"wire datastar", FilterBarProps{Wire: &wire.Action{Transport: wire.TransportDatastar, URL: "/x"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := filterBarIsHTMX(tt.props); got != tt.want {
				t.Errorf("filterBarIsHTMX() = %v, want %v", got, tt.want)
			}
		})
	}
}
