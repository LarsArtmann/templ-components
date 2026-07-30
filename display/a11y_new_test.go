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

// --- Sparkline Accessibility ---

func TestSparklineA11y(t *testing.T) {
	t.Parallel()

	t.Run("hidden from screen readers by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Sparkline(SparklineProps{
			Values: []float64{1, 3, 2, 5},
		}))
		utils.AssertContains(t, output, `aria-hidden="true"`)
	})

	t.Run("aria-label when provided for meaningful trends", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Sparkline(SparklineProps{
			Values:    []float64{1, 3, 2, 5},
			BaseProps: utils.BaseProps{AriaLabel: "Daily active users trend"},
		}))
		utils.AssertContains(t, output, `aria-label="Daily active users trend"`)
		utils.AssertNotContains(t, output, `aria-hidden`)
	})
}

// --- BarChart Accessibility ---

func TestBarChartA11y(t *testing.T) {
	t.Parallel()

	t.Run("has role=img for screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{
				{Label: "general", Value: 100},
			},
		}))
		utils.AssertContains(t, output, `role="img"`)
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{
				{Label: "general", Value: 100},
			},
			BaseProps: utils.BaseProps{AriaLabel: "Messages by channel"},
		}))
		utils.AssertContains(t, output, `aria-label="Messages by channel"`)
	})
}

// --- ExternalLink Accessibility ---

func TestExternalLinkA11y(t *testing.T) {
	t.Parallel()

	t.Run("has target=_blank", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ExternalLink(ExternalLinkProps{
			Href: "https://example.com",
			Text: "Link",
		}))
		utils.AssertContains(t, output, `target="_blank"`)
	})

	t.Run("has rel=noopener noreferrer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ExternalLink(ExternalLinkProps{
			Href: "https://example.com",
			Text: "Link",
		}))
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})

	t.Run("arrow icon hidden from screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ExternalLink(ExternalLinkProps{
			Href:     "https://example.com",
			Text:     "Link",
			ShowIcon: true,
		}))
		utils.AssertContains(t, output, `aria-hidden="true"`)
	})
}
