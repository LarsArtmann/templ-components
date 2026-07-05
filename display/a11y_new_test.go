package display

import (
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/utils"
)

// --- CopyButton Accessibility ---

func TestCopyButtonA11y(t *testing.T) {
	t.Parallel()

	t.Run("button has type=button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		utils.AssertContains(t, output, `type="button"`)
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text:      "x",
			BaseProps: utils.BaseProps{AriaLabel: "Copy install command"},
		}))
		utils.AssertContains(t, output, `aria-label="Copy install command"`)
	})

	t.Run("has focus-visible ring", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		utils.AssertContains(t, output, "focus-visible:ring-2")
	})

	t.Run("has motion-reduce transition", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		utils.AssertContains(t, output, "motion-reduce:transition-none")
		utils.AssertContains(t, output, "motion-reduce:duration-0")
	})
}

// --- RelativeTime Accessibility ---

func TestRelativeTimeA11y(t *testing.T) {
	t.Parallel()

	t.Run("has machine-readable datetime", func(t *testing.T) {
		t.Parallel()
		ts := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
		output := utils.Render(t, RelativeTime(RelativeTimeProps{Time: ts}))
		utils.AssertContains(t, output, `datetime="2025-06-15T12:00:00Z"`)
	})

	t.Run("has title for hover context", func(t *testing.T) {
		t.Parallel()
		ts := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
		output := utils.Render(t, RelativeTime(RelativeTimeProps{Time: ts}))
		utils.AssertContains(t, output, `title="Jun 15, 2025`)
	})
}

// --- CountBadge Accessibility ---

func TestCountBadgeA11y(t *testing.T) {
	t.Parallel()

	t.Run("badge is aria-hidden (decorative)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 3}))
		utils.AssertContains(t, output, `aria-hidden="true"`)
	})

	t.Run("propagates aria-label to container", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{
			Count:     1,
			BaseProps: utils.BaseProps{AriaLabel: "3 unread notifications"},
		}))
		utils.AssertContains(t, output, `aria-label="3 unread notifications"`)
	})
}

// --- DefinitionGrid Accessibility ---

func TestDefinitionGridA11y(t *testing.T) {
	t.Parallel()

	t.Run("renders semantic dl/dt/dd", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
			Items: []DefinitionItem{{Term: "X", Detail: "Y"}},
		}))
		utils.AssertContains(t, output, "<dl")
		utils.AssertContains(t, output, "<dt")
		utils.AssertContains(t, output, "<dd")
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
			Items:     []DefinitionItem{{Term: "X", Detail: "Y"}},
			BaseProps: utils.BaseProps{AriaLabel: "System metrics"},
		}))
		utils.AssertContains(t, output, `aria-label="System metrics"`)
	})
}

// --- Image Accessibility ---

func TestImageA11y(t *testing.T) {
	t.Parallel()

	t.Run("includes alt text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{
			Src: "/x.jpg",
			Alt: "Profile photo",
		}))
		utils.AssertContains(t, output, `alt="Profile photo"`)
	})
}

// --- StatCard HTMX Accessibility ---

func TestStatCardHxAttributes(t *testing.T) {
	t.Parallel()

	t.Run("adds hx-get when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label: "Users",
			Value: "42",
			HxGet: "/api/stats",
		}))
		utils.AssertContains(t, output, `hx-get="/api/stats"`)
	})

	t.Run("adds hx-target when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label:    "Users",
			Value:    "42",
			HxGet:    "/api/stats",
			HxTarget: "#stats",
		}))
		utils.AssertContains(t, output, `hx-target="#stats"`)
	})

	t.Run("adds hx-swap when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label:  "Users",
			Value:  "42",
			HxGet:  "/api/stats",
			HxSwap: htmx.SwapInnerHTML,
		}))
		utils.AssertContains(t, output, `hx-swap="innerHTML"`)
	})
}

// --- Card.Body Accessibility ---

func TestCardBodySlot(t *testing.T) {
	t.Parallel()

	t.Run("renders Body when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title: "Test",
			Body:  templ.Raw("<p>Body content from slot</p>"),
		}))
		utils.AssertContains(t, output, "Body content from slot")
	})
}
