package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestCommandPaletteSizeIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size CommandPaletteSize
		want bool
	}{
		{"small", CommandPaletteSizeSM, true},
		{"medium", CommandPaletteSizeMD, true},
		{"large", CommandPaletteSizeLG, true},
		{"empty", "", false},
		{"unknown", "xl", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := CommandPaletteSizeIsValid(tt.size); got != tt.want {
				t.Errorf("CommandPaletteSizeIsValid(%q) = %v, want %v", tt.size, got, tt.want)
			}
		})
	}
}

func TestGoldenSweepCommandPalette(t *testing.T) {
	t.Parallel()

	linkGroups := []CommandGroup{
		{
			Label: "Pages",
			Items: []CommandItem{
				{Label: "Dashboard", Href: "/dashboard", Hint: "Go to dashboard"},
				{Label: "Settings", Href: "/settings"},
			},
		},
		{
			Label: "Actions",
			Items: []CommandItem{
				{Label: "Sign out", Href: "/logout"},
			},
		},
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{
			Name: "command_palette_default",
			HTML: utils.Render(t, CommandPalette(DefaultCommandPaletteProps())),
		},
		{
			Name: "command_palette_trigger_hotkey",
			HTML: utils.Render(t, CommandPalette(CommandPaletteProps{
				ID:           "cmdk",
				TriggerLabel: "Commands",
				Hotkey:       true,
				Groups:       linkGroups,
			})),
		},
		{
			Name: "command_palette_link_items",
			HTML: utils.Render(t, CommandPalette(CommandPaletteProps{
				ID:     "links",
				Groups: linkGroups,
			})),
		},
		{
			Name: "command_palette_no_groups",
			HTML: utils.Render(t, CommandPalette(CommandPaletteProps{
				ID: "empty",
			})),
		},
		{
			Name: "command_palette_large",
			HTML: utils.Render(t, CommandPalette(CommandPaletteProps{
				ID:     "wide",
				Size:   CommandPaletteSizeLG,
				Groups: linkGroups[:1],
			})),
		},
	})
}

func TestGoldenSweepCommandPaletteWire(t *testing.T) {
	t.Parallel()

	wired := wire.Post("/api/items/42/delete").WithTarget("#item-42")
	action := &wired

	golden.AssertSnapshots(t, []golden.Snapshot{
		{
			Name: "command_palette_wire_item",
			HTML: utils.Render(t, CommandPalette(CommandPaletteProps{
				ID: "wired",
				Groups: []CommandGroup{
					{
						Label: "Danger zone",
						Items: []CommandItem{
							{Label: "Delete item", Hint: "Cannot be undone", Wire: action},
						},
					},
				},
			})),
		},
		{
			Name: "command_palette_inert_item",
			HTML: utils.Render(t, CommandPalette(CommandPaletteProps{
				ID: "inert",
				Groups: []CommandGroup{
					{Label: "Read only", Items: []CommandItem{{Label: "Just a label"}}},
				},
			})),
		},
	})
}

func TestCommandPaletteStructure(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, CommandPalette(CommandPaletteProps{
		ID: "palette", Nonce: "test-nonce",
		TriggerLabel: "Commands",
		Hotkey:       true,
		Groups: []CommandGroup{
			{
				Label: "Pages",
				Items: []CommandItem{
					{Label: "Dashboard", Href: "/dashboard"},
					{Label: "Reports", Href: "/reports"},
				},
			},
			{
				Label: "Actions",
				Items: []CommandItem{{Label: "Sign out", Href: "/logout"}},
			},
		},
	}))

	// Item ids must be unique ACROSS groups (running offset, not per-group).
	utils.AssertContains(t, html, `id="palette-item-0"`)
	utils.AssertContains(t, html, `id="palette-item-1"`)
	utils.AssertContains(t, html, `id="palette-item-2"`)
	utils.AssertNotContains(t, html, `id="palette-item-1-0"`)

	// ARIA combobox wiring between input and listbox.
	utils.AssertContains(t, html, `role="combobox"`)
	utils.AssertContains(t, html, `aria-controls="palette-list"`)
	utils.AssertContains(t, html, `aria-autocomplete="list"`)
	utils.AssertContains(t, html, `role="listbox"`)
	utils.AssertContains(t, html, `role="option"`)

	// Hotkey affordances: dialog flag + visible kbd hint.
	utils.AssertContains(t, html, `data-tc-hotkey="true"`)
	utils.AssertContains(t, html, "<kbd")

	// Open trigger + singleton script carrying the nonce.
	utils.AssertContains(t, html, `data-tc-palette-open`)
	utils.AssertContains(t, html, `nonce="test-nonce"`)
	utils.AssertContains(t, html, "tcPaletteAttached")
}

func TestCommandPaletteWireAttributesRender(t *testing.T) {
	t.Parallel()

	wired := wire.Post("/api/run").WithTarget("#output")
	action := &wired
	html := utils.Render(t, CommandPalette(CommandPaletteProps{
		ID: "wired",
		Groups: []CommandGroup{
			{Items: []CommandItem{{Label: "Run job", Wire: action}}},
		},
	}))

	// The wire item must render the htmx dialect attributes on the button.
	utils.AssertContains(t, html, `hx-post="/api/run"`)
	utils.AssertContains(t, html, `hx-target="#output"`)
}

func TestCommandPaletteEmptyMessage(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, CommandPalette(CommandPaletteProps{
		ID:           "msg",
		EmptyMessage: "Nothing matches.",
	}))

	utils.AssertContains(t, html, "Nothing matches.")
	utils.AssertContains(t, html, `data-tc-palette-empty`)
}
