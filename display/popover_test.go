package display

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestPopoverRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		props   PopoverProps
		wantAll []string
	}{
		{
			name:  "bottom popover with trigger text",
			props: PopoverProps{TriggerText: "Details"},
			wantAll: []string{
				"Details",
				`role="dialog"`,
				`aria-expanded="false"`,
				`aria-haspopup="dialog"`,
				"top-full",
			},
		},
		{
			name:    "top position",
			props:   PopoverProps{TriggerText: "Info", Position: PopoverPositionTop},
			wantAll: []string{"bottom-full"},
		},
		{
			name:    "left position",
			props:   PopoverProps{TriggerText: "Left", Position: PopoverPositionLeft},
			wantAll: []string{"right-full"},
		},
		{
			name:    "right position",
			props:   PopoverProps{TriggerText: "Right", Position: PopoverPositionRight},
			wantAll: []string{"left-full"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, Popover(tt.props))
			for _, want := range tt.wantAll {
				utils.AssertContains(t, output, want)
			}
		})
	}
}

func TestPopoverWithContent(t *testing.T) {
	t.Parallel()
	t.Run("renders child content in panel", func(t *testing.T) {
		t.Parallel()

		child := templ.Raw(`<p data-test="content">Popover body</p>`)

		var buf bytes.Buffer

		_ = Popover(PopoverProps{TriggerText: "Open"}).Render(
			templ.WithChildren(context.Background(), child), &buf,
		)
		output := strings.TrimSpace(buf.String())
		utils.AssertContains(t, output, "Popover body")
		utils.AssertContains(t, output, `data-test="content"`)
	})
}

func TestDefaultPopoverProps(t *testing.T) {
	t.Parallel()

	props := DefaultPopoverProps()
	if props.Position != PopoverPositionBottom {
		t.Errorf("DefaultPopoverProps().Position = %q, want %q", props.Position, PopoverPositionBottom)
	}
}

func TestPopoverA11y(t *testing.T) {
	t.Parallel()

	t.Run("trigger has aria-haspopup dialog", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{TriggerText: "Open"}))
		utils.AssertContains(t, output, `aria-haspopup="dialog"`)
	})

	t.Run("trigger has aria-expanded false by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{TriggerText: "Open"}))
		utils.AssertContains(t, output, `aria-expanded="false"`)
	})

	t.Run("trigger aria-controls links to content panel", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			BaseProps:   utils.BaseProps{ID: "my-pop"},
			TriggerText: "Open",
		}))
		utils.AssertContains(t, output, `aria-controls="my-pop-content"`)
	})

	t.Run("content panel has role dialog", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{TriggerText: "Open"}))
		utils.AssertContains(t, output, `role="dialog"`)
	})

	t.Run("content panel has aria-labelledby trigger", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			BaseProps:   utils.BaseProps{ID: "p1"},
			TriggerText: "Open",
		}))
		utils.AssertContains(t, output, `aria-labelledby="p1-trigger"`)
	})

	t.Run("auto-generated ID links trigger and content", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{TriggerText: "Open"}))
		utils.AssertContains(t, output, `id="`)
		utils.AssertContains(t, output, "-trigger")
		utils.AssertContains(t, output, "-content")
	})
}

func TestPopoverEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("unknown position falls back to bottom", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			TriggerText: "X",
			Position:    PopoverPosition("bogus"),
		}))
		utils.AssertContains(t, output, "top-full")
	})

	t.Run("custom class on root element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			BaseProps:   utils.BaseProps{Class: "custom-pop"},
			TriggerText: "X",
		}))
		utils.AssertContains(t, output, "custom-pop")
	})

	t.Run("aria-label propagated to root", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			BaseProps:   utils.BaseProps{AriaLabel: "More info popover"},
			TriggerText: "X",
		}))
		utils.AssertContains(t, output, `aria-label="More info popover"`)
	})

	t.Run("nonce has no script tag (popover API is JS-free)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			BaseProps:   utils.BaseProps{Nonce: "test-nonce-123"},
			TriggerText: "X",
		}))
		utils.AssertNotContains(t, output, "<script")
		utils.AssertNotContains(t, output, "test-nonce-123")
	})

	t.Run("trigger has popovertarget attribute (native invoker)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{
			BaseProps:   utils.BaseProps{ID: "native-pop"},
			TriggerText: "X",
		}))
		utils.AssertContains(t, output, `popovertarget="native-pop-content"`)
		utils.AssertContains(t, output, `popover="auto"`)
	})
}

func TestPopoverGolden(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Popover(PopoverProps{
		BaseProps:   utils.BaseProps{ID: "golden-popover"},
		TriggerText: "Details",
		Position:    PopoverPositionBottom,
	}))
	golden.Assert(t, "popover_bottom", output)
}

func TestPopoverDarkModeCompliance(t *testing.T) {
	t.Parallel()
	t.Run("panel has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{TriggerText: "X"}))
		utils.AssertContains(t, output, "dark:bg-gray-800")
		utils.AssertContains(t, output, "dark:text-white")
		utils.AssertContains(t, output, "dark:shadow-black/20")
	})
	t.Run("trigger button has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Popover(PopoverProps{TriggerText: "X"}))
		utils.AssertContains(t, output, "dark:bg-gray-800")
		utils.AssertContains(t, output, "dark:ring-gray-600")
	})
}
