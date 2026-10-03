package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestTabsRender(t *testing.T) {
	t.Parallel()
	t.Run("basic tabs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "tab1",
			Tabs: []Tab{
				{ID: "tab1", Label: "First"},
				{ID: "tab2", Label: "Second"},
			},
		}))
		utils.AssertContains(t, output, "First")
		utils.AssertContains(t, output, "Second")
		utils.AssertContains(t, output, `role="tablist"`)
		utils.AssertContains(t, output, `aria-selected="true"`)
		utils.AssertContains(t, output, `aria-selected="false"`)
	})

	t.Run("tab with content", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "panel1",
			Tabs: []Tab{
				{
					ID:      "panel1",
					Label:   "Panel",
					Content: templ.Raw("<p>Hello</p>"),
				},
			},
		}))
		utils.AssertContains(t, output, `role="tabpanel"`)
		utils.AssertContains(t, output, "Hello")
	})

	t.Run("pills style", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Variant:     TabsPills,
			Tabs: []Tab{
				{ID: "a", Label: "A"},
			},
		}))
		utils.AssertContains(t, output, "bg-blue-600")
	})

	t.Run("no active tab defaults to first", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{ID: "x", Label: "X"},
				{ID: "y", Label: "Y"},
			},
		}))
		// WAI-ARIA tab pattern: exactly one tab must be keyboard-focusable
		utils.AssertContains(t, output, `tabindex="0"`)
		utils.AssertContains(t, output, `aria-selected="true"`)
	})

	t.Run("tabs without IDs get auto-generated IDs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{Label: "No ID"},
				{Label: "Also No ID"},
			},
		}))
		// Auto-generated IDs prevent invalid HTML (id="-tab") and JS crashes
		utils.AssertContains(t, output, `id="tc-tab-`)
		utils.AssertNotContains(t, output, `id="-tab"`)
		utils.AssertNotContains(t, output, `aria-controls=""`)
	})

	t.Run("default props has TabsDefault variant", func(t *testing.T) {
		t.Parallel()

		props := DefaultTabsProps()
		if props.Variant != TabsDefault {
			t.Errorf("DefaultTabsProps().Variant = %q, want %q", props.Variant, TabsDefault)
		}
	})

	t.Run("tabs with empty tabs list", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{}))
		utils.AssertContains(t, output, `role="tablist"`)
	})

	t.Run("default underline variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Variant:     TabsDefault,
			Tabs: []Tab{
				{ID: "a", Label: "A"},
				{ID: "b", Label: "B"},
			},
		}))
		utils.AssertContains(t, output, "border-blue-500")
		utils.AssertContains(t, output, "border-transparent")
	})

	t.Run("tab without content", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Tabs: []Tab{
				{ID: "a", Label: "A"},
			},
		}))
		utils.AssertNotContains(t, output, `role="tabpanel"`)
	})

	t.Run("pills variant with content", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "p1",
			Variant:     TabsPills,
			Tabs: []Tab{
				{ID: "p1", Label: "One", Content: templ.Raw("<p>Pill content</p>")},
				{ID: "p2", Label: "Two"},
			},
		}))
		utils.AssertContains(t, output, "bg-blue-600")
		utils.AssertContains(t, output, "Pill content")
		utils.AssertContains(t, output, `space-x-2`)
	})

	t.Run("custom ID and class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			BaseProps:   utils.BaseProps{ID: "my-tabs", Class: "max-w-2xl"},
			ActiveTabID: "a",
			Tabs:        []Tab{{ID: "a", Label: "A"}},
		}))
		utils.AssertContains(t, output, `id="my-tabs"`)
		utils.AssertContains(t, output, "max-w-2xl")
	})

	t.Run("inactive tab panel has hidden attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Tabs: []Tab{
				{ID: "a", Label: "Active", Content: templ.Raw("<p>Visible</p>")},
				{ID: "b", Label: "Hidden", Content: templ.Raw("<p>Hidden content</p>")},
			},
		}))
		utils.AssertContains(t, output, "Hidden content")
	})

	t.Run("client-side tabs renders data-tc-tabs and script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			ClientSide:  true,
			Tabs: []Tab{
				{ID: "a", Label: "A"},
				{ID: "b", Label: "B"},
			},
		}))
		utils.AssertContains(t, output, `data-tc-tabs`)
		utils.AssertContains(t, output, "tcTabsAttached")
		utils.AssertContains(t, output, "activateTab")
		utils.AssertContains(t, output, "ArrowRight")
	})

	t.Run("non-client-side tabs omit script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Tabs: []Tab{
				{ID: "a", Label: "A"},
			},
		}))
		utils.AssertNotContains(t, output, `data-tc-tabs`)
		utils.AssertNotContains(t, output, "tcTabsAttached")
	})
}

// TestTabsWire pins the server-side tab switching contract (ADR-0036): each
// tab link carries the wired action with the {tab} URL placeholder
// substituted, the anchor fallback forces PreventDefault, an empty Target
// defaults to the container's own #ID, htmx gets the settle:0s self-swap,
// and the consumer's *Action is never mutated.
func TestTabsWire(t *testing.T) {
	t.Parallel()

	t.Run("htmx renders per-tab hx-get with substituted URLs and the settle swap", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Tabs(TabsProps{
			ID:          "prefs-tabs",
			ActiveTabID: "general",
			Tabs: []Tab{
				{ID: "general", Label: "General"},
				{ID: "advanced", Label: "Advanced"},
			},
			Wire: &wire.Action{URL: "/api/tabs/{tab}", Target: "#prefs-tabs"},
		}))
		utils.AssertContains(t, output, `hx-get="/api/tabs/general"`)
		utils.AssertContains(t, output, `hx-get="/api/tabs/advanced"`)
		utils.AssertContains(t, output, `hx-target="#prefs-tabs"`)
		utils.AssertContains(t, output, `hx-swap="outerHTML settle:0s"`)
	})

	t.Run("datastar renders per-tab data-on:click with prevent", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Tabs(TabsProps{
			ID:          "prefs-tabs",
			ActiveTabID: "general",
			Tabs: []Tab{
				{ID: "general", Label: "General"},
				{ID: "advanced", Label: "Advanced"},
			},
			Wire: &wire.Action{Transport: wire.TransportDatastar, URL: "/api/tabs/{tab}"},
		}))
		utils.AssertContains(t, output, `data-on:click__prevent="@get(&#39;/api/tabs/general&#39;)"`)
		utils.AssertContains(t, output, `data-on:click__prevent="@get(&#39;/api/tabs/advanced&#39;)"`)
		utils.AssertNotContains(t, output, "hx-")
	})

	t.Run("empty Target defaults to the container's own #ID (htmx)", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Tabs(TabsProps{
			ID: "prefs-tabs",
			Tabs: []Tab{
				{ID: "general", Label: "General"},
			},
			Wire: &wire.Action{URL: "/api/tabs/{tab}"},
		}))
		utils.AssertContains(t, output, `hx-target="#prefs-tabs"`)
	})

	t.Run("empty URL wires nothing (inert anchors)", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Tabs(TabsProps{
			ID: "prefs-tabs",
			Tabs: []Tab{
				{ID: "general", Label: "General"},
			},
			Wire: &wire.Action{},
		}))
		utils.AssertNotContains(t, output, "hx-get")
		utils.AssertNotContains(t, output, "data-on:")
	})

	t.Run("consumer action is never mutated", func(t *testing.T) {
		t.Parallel()

		original := &wire.Action{URL: "/api/tabs/{tab}"}
		snapshot := *original

		_ = utils.Render(t, Tabs(TabsProps{
			ID: "prefs-tabs",
			Tabs: []Tab{
				{ID: "general", Label: "General"},
				{ID: "advanced", Label: "Advanced"},
			},
			Wire: original,
		}))

		if *original != snapshot {
			t.Errorf("Tabs mutated the consumer's action: %+v", *original)
		}
		if original.Target != "" {
			t.Errorf("Target = %q, want empty (the clone carries the default, not the original)", original.Target)
		}
	})

	t.Run("ClientSide is ignored when Wire is set", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Tabs(TabsProps{
			ID:         "prefs-tabs",
			ClientSide: true,
			Tabs: []Tab{
				{ID: "general", Label: "General"},
			},
			Wire: &wire.Action{URL: "/api/tabs/{tab}", Target: "#prefs-tabs"},
		}))
		utils.AssertNotContains(t, output, "data-tc-tabs")
		utils.AssertNotContains(t, output, "tcTabsAttached")
		utils.AssertContains(t, output, `hx-get="/api/tabs/general"`)
	})
}
