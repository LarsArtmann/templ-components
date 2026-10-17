package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/tone"
)

const (
	activeBadgeText      = "Active"
	dropdownLabelActions = "Actions"
	dropdownHrefEdit     = "/edit"
	avatarAltAlice       = "Alice"
	tableHeaderName      = "Name"
	dropdownLabelMenu    = "Menu"
	dropdownItemEdit     = "Edit"
	accordionIDFAQ1      = "faq1"
	testNonce            = "test-nonce"
	cssClassMt4          = "mt-4"
)

func TestA11yAttributes(t *testing.T) {
	t.Parallel()

	t.Run("modal uses native dialog element", func(t *testing.T) {
		t.Parallel()

		props := DefaultModalProps()
		props.ID = "test-modal"
		props.Title = "Test"
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "<dialog")
	})

	t.Run("dropdown has proper ARIA", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			ID:    "dd",
			Label: dropdownLabelActions,
			Items: []DropdownItem{
				{Text: dropdownItemEdit, Href: dropdownHrefEdit},
			},
		}))
		utils.AssertContains(t, output, `popovertarget=`)
		utils.AssertContains(t, output, `aria-haspopup="true"`)
		utils.AssertContains(t, output, `role="menu"`)
		utils.AssertContains(t, output, `role="menuitem"`)
	})

	t.Run("tabs have proper ARIA", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "a",
			Tabs: []Tab{
				{ID: "a", Label: activeBadgeText},
				{ID: "b", Label: "Inactive"},
			},
		}))
		utils.AssertContains(t, output, `role="tablist"`)
		utils.AssertContains(t, output, `role="tab"`)
		utils.AssertContains(t, output, `aria-selected="true"`)
	})

	t.Run("accordion uses native details element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "a", Title: "Q1", Open: true},
				{ID: "b", Title: "Q2"},
			},
		}))
		utils.AssertContains(t, output, "<details")
		utils.AssertContains(t, output, "<summary")
		utils.AssertContains(t, output, " open")
	})

	t.Run("avatar image has alt text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src: "/avatar.jpg",
			Alt: avatarAltAlice,
		}))
		utils.AssertContains(t, output, `alt="Alice"`)
	})
}

func TestDarkModeClasses(t *testing.T) {
	t.Parallel()

	t.Run("card has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(DefaultCardProps()))
		utils.AssertContains(t, output, "dark:bg-gray-800")
		utils.AssertContains(t, output, "dark:border-gray-700")
	})

	t.Run("badge has dark mode classes", func(t *testing.T) {
		t.Parallel()

		props := DefaultBadgeProps()
		props.Text = activeBadgeText
		output := utils.Render(t, Badge(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("table has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{tableHeaderName},
			Rows:    []TableRow{SimpleTableRow("Alice")},
		}))
		utils.AssertContains(t, output, "dark:divide-gray-700")
		utils.AssertContains(t, output, "dark:text-gray-300")
	})

	t.Run("dropdown has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			ID:    "dd",
			Label: dropdownLabelMenu,
		}))
		utils.AssertContains(t, output, "dark:bg-gray-800")
	})

	t.Run("avatar initials have dark mode background", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "AB",
		}))
		utils.AssertContains(t, output, "bg-blue-600")
	})
}

func TestDefaultPropsConstructors(t *testing.T) {
	t.Parallel()

	t.Run("DefaultCardProps returns valid props", func(t *testing.T) {
		t.Parallel()

		props := DefaultCardProps()
		if props.Padding != CardPaddingMD {
			t.Errorf("Padding = %q, want %q", props.Padding, CardPaddingMD)
		}

		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, "<div")
	})

	t.Run("DefaultBadgeProps returns valid props", func(t *testing.T) {
		t.Parallel()

		props := DefaultBadgeProps()
		if props.Type != BadgeNeutral {
			t.Errorf("Type = %q, want %q", props.Type, BadgeNeutral)
		}

		output := utils.Render(t, Badge(props))
		utils.AssertContains(t, output, "<span")
	})

	t.Run("DefaultModalProps returns valid props", func(t *testing.T) {
		t.Parallel()

		props := DefaultModalProps()
		if props.Size != ModalSizeMD {
			t.Errorf("Size = %q, want %q", props.Size, ModalSizeMD)
		}
	})

	// The M26 coverage batch: these constructors were public API with 0%
	// coverage — the consumer-audit lesson made them the documented-defaults
	// promise, so each gets its default pinned.
	t.Run("DefaultStatusDotProps defaults to neutral", func(t *testing.T) {
		t.Parallel()

		if got := DefaultStatusDotProps(); got.Tone != StatusToneNeutral {
			t.Errorf("Tone = %q, want %q", got.Tone, StatusToneNeutral)
		}
	})

	t.Run("DefaultProseProps defaults to MD", func(t *testing.T) {
		t.Parallel()

		if got := DefaultProseProps(); got.Size != ProseSizeMD {
			t.Errorf("Size = %q, want %q", got.Size, ProseSizeMD)
		}
	})

	t.Run("DefaultDataStateProps defaults to Content", func(t *testing.T) {
		t.Parallel()

		if got := DefaultDataStateProps(); got.State != DataStateContent {
			t.Errorf("State = %q, want %q", got.State, DataStateContent)
		}
	})

	t.Run("remaining constructors construct", func(t *testing.T) {
		t.Parallel()
		// Contract: the documented defaults exist; deeper behavior is the
		// component's own test surface.
		_ = DefaultCodeBlockProps()
		_ = DefaultLightboxProps()
		_ = DefaultDonutChartProps()
		_ = DefaultSegmentBarProps()
	})
}

// TestToneAdaptersExhaustive covers the tone Phase 0 adapters (M26 batch):
// every mapped enum value routes to its canonical tone and unknown values
// degrade to neutral — the map+fallback convention applied to the adapter
// layer.
func TestToneAdaptersExhaustive(t *testing.T) {
	t.Parallel()

	statusCases := map[StatusTone]tone.Tone{
		StatusToneNeutral: tone.ToneNeutral,
		StatusToneSuccess: tone.ToneSuccess,
		StatusToneWarning: tone.ToneWarning,
		StatusToneDanger:  tone.ToneDanger,
		StatusToneInfo:    tone.ToneInfo,
	}
	for in, want := range statusCases {
		if got := in.Tone(); got != want {
			t.Errorf("StatusTone(%q).Tone() = %q, want %q", in, got, want)
		}
	}

	scrollbackCases := map[ScrollbackTone]tone.Tone{
		ScrollbackToneNeutral: tone.ToneNeutral,
		ScrollbackToneInfo:    tone.ToneInfo,
		ScrollbackToneSuccess: tone.ToneSuccess,
		ScrollbackToneWarning: tone.ToneWarning,
		ScrollbackToneDanger:  tone.ToneDanger,
	}
	for in, want := range scrollbackCases {
		if got := in.Tone(); got != want {
			t.Errorf("ScrollbackTone(%q).Tone() = %q, want %q", in, got, want)
		}
	}

	if got := StatusTone("bogus").Tone(); got != tone.ToneNeutral {
		t.Errorf("unknown StatusTone degraded to %q, want neutral", got)
	}
}

func TestDropdownXSSSafety(t *testing.T) {
	t.Parallel()

	t.Run("ID with special chars is safely interpolated", func(t *testing.T) {
		t.Parallel()

		maliciousID := `<script>alert('xss')</script>`
		output := utils.Render(t, Dropdown(DropdownProps{
			ID:    maliciousID,
			Label: dropdownLabelMenu,
		}))
		utils.AssertNotContains(t, output, `<script>alert('xss')</script>`)
		utils.AssertContains(t, output, `&lt;script&gt;`)
	})
}

func TestGridA11y(t *testing.T) {
	t.Parallel()

	t.Run("Grid propagates AriaLabel", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Grid(GridProps{
			AriaLabel: "User cards",
			Cols:      GridCols3,
		}))
		utils.AssertContains(t, output, `aria-label="User cards"`)
	})

	t.Run("Grid propagates ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Grid(GridProps{
			ID: "user-grid",
		}))
		utils.AssertContains(t, output, `id="user-grid"`)
	})
}

func TestStatCardHrefA11y(t *testing.T) {
	t.Parallel()

	t.Run("Href anchor propagates AriaLabel", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label:     "Active",
			Value:     "42",
			AriaLabel: "Filter by active users",
			Href:      "/active",
		}))
		utils.AssertContains(t, output, `aria-label="Filter by active users"`)
		utils.AssertContains(t, output, `<a `)
	})

	t.Run("Href anchor has focus-visible ring for keyboard nav", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label: "Total",
			Value: "100",
			Href:  "/total",
		}))
		utils.AssertContains(t, output, "focus-visible:ring-2")
	})
}
