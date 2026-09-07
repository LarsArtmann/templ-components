package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestFilterInputWiring(t *testing.T) {
	t.Parallel()

	t.Run("htmx dialect renders debounce trigger and target", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:       "q",
			DebounceMS: 300,
			Wire: &wire.Action{
				URL:    "/api/search",
				Target: "#results",
			},
		}))
		utils.AssertContains(t, output, `hx-get="/api/search"`)
		utils.AssertContains(t, output, `hx-trigger="input changed delay:300ms"`)
		utils.AssertContains(t, output, `hx-target="#results"`)
		utils.AssertContains(t, output, `name="q"`)
	})

	t.Run("datastar dialect renders debounce modifier and form encoding", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:       "q",
			DebounceMS: 300,
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/search",
			},
		}))
		utils.AssertContains(t, output, `data-on:input__debounce.300ms=`)
		utils.AssertContains(t, output, `@get('/api/search', {contentType: 'form'})`)
	})

	t.Run("zero debounce emits no delay modifier", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:       "q",
			DebounceMS: 0,
			Wire: &wire.Action{
				URL: "/api/search",
			},
		}))
		utils.AssertContains(t, output, `hx-trigger="input changed"`)
		utils.AssertNotContains(t, output, "delay:")
	})

	t.Run("empty URL wires nothing (inert search form)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(DefaultFilterInputProps()))
		utils.AssertNotContains(t, output, "hx-get")
		utils.AssertNotContains(t, output, "data-on:")
		utils.AssertContains(t, output, `method="GET"`)
	})

	t.Run("method defaults applied (input event, form encoding)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:       "q",
			DebounceMS: 300,
			Wire:       &wire.Action{URL: "/api/search"},
		}))
		utils.AssertContains(t, output, `hx-trigger="input changed delay:300ms"`)
	})
}

func TestFilterInputA11y(t *testing.T) {
	t.Parallel()

	t.Run("search landmark wraps the form", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(DefaultFilterInputProps()))
		utils.AssertContains(t, output, "<search>")
		utils.AssertContains(t, output, `type="search"`)
	})

	t.Run("label is associated via for", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:  "q",
			Label: "Search users",
			ID:    "user-search",
		}))
		utils.AssertContains(t, output, `for="user-search"`)
		utils.AssertContains(t, output, `id="user-search"`)
	})

	t.Run("auto-generated id keeps label associated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:  "q",
			Label: "Search users",
		}))
		utils.AssertContains(t, output, `for="tc-filter-input-`)
	})

	t.Run("aria-label reaches the input when no visible label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:      "q",
			AriaLabel: "Search users",
		}))
		utils.AssertContains(t, output, `aria-label="Search users"`)
	})

	t.Run("help text is linked via aria-describedby", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:     "q",
			ID:       "user-search",
			HelpText: "Matches name or email",
		}))
		utils.AssertContains(t, output, `aria-describedby="user-search-help"`)
		utils.AssertContains(t, output, `id="user-search-help"`)
		utils.AssertContains(t, output, "Matches name or email")
	})
}

func TestFilterInputEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("value is preserved across re-renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:  "q",
			Value: "ada",
		}))
		utils.AssertContains(t, output, `value="ada"`)
	})

	t.Run("no-JS fallback action renders on the form", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:   "q",
			Action: "/users",
		}))
		utils.AssertContains(t, output, `action="/users"`)
		utils.AssertContains(t, output, `method="GET"`)
	})

	t.Run("no action attribute when unset", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(DefaultFilterInputProps()))
		utils.AssertNotContains(t, output, `action="`)
	})

	t.Run("placeholder renders when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:        "q",
			Placeholder: "Type to filter…",
		}))
		utils.AssertContains(t, output, `placeholder="Type to filter…"`)
	})

	t.Run("unknown wire method falls back to get", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterInput(FilterInputProps{
			Name:       "q",
			DebounceMS: 300,
			Wire: &wire.Action{
				Method: wire.Method("fetch"),
				URL:    "/api/search",
			},
		}))
		utils.AssertContains(t, output, `hx-get="/api/search"`)
		utils.AssertNotContains(t, output, "hx-fetch")
	})
}
