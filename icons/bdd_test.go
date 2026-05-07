// Package icons provides behavior-driven tests for the icon system.
// These tests verify end-user experience: seeing icons, accessibility, theming.
package icons

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- Icon Rendering Behavior ---

func TestIconUserSeesCorrectIcon(t *testing.T) {
	t.Parallel()

	t.Run("user sees a home icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Home, "h-6 w-6"))
		utils.AssertContains(t, output, "<svg")
		utils.AssertContains(t, output, `class="h-6 w-6"`)
	})

	t.Run("user sees icon with custom class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Search, "h-4 w-4 text-blue-500"))
		utils.AssertContains(t, output, "text-blue-500")
	})
}

// --- Icon Accessibility Behavior ---

func TestIconUserGetsAccessibleIcons(t *testing.T) {
	t.Parallel()

	t.Run("decorative icon is hidden from screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Check, "h-5 w-5"))
		utils.AssertContains(t, output, `aria-hidden="true"`)
	})

	t.Run("icon SVG uses currentColor for theming", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Bell, "h-5 w-5"))
		utils.AssertContains(t, output, "currentColor")
	})
}

// --- All Icons Render Behavior ---

func TestAllIconsRenderSuccessfully(t *testing.T) {
	t.Parallel()

	allIcons := []struct {
		name Name
	}{
		{Home},
		{Users},
		{Folder},
		{Document},
		{Search},
		{Settings},
		{Chart},
		{Inbox},
		{Check},
		{X},
		{Plus},
		{Minus},
		{ChevronRight},
		{ChevronLeft},
		{ChevronDown},
		{ChevronUp},
		{ArrowRight},
		{ArrowLeft},
		{Refresh},
		{ExternalLink},
		{Download},
		{Upload},
		{Trash},
		{Edit},
		{Eye},
		{EyeOff},
		{Lock},
		{Unlock},
		{Menu},
		{Bell},
		{Calendar},
		{Clock},
		{Location},
		{Phone},
		{Mail},
		{Globe},
		{Sun},
		{Moon},
		{Spinner},
		{Exclamation},
		{Information},
		{Question},
	}

	for _, tc := range allIcons {
		t.Run("icon "+string(tc.name)+" renders SVG", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Icon(tc.name, "h-5 w-5"))
			if !strings.Contains(output, "<svg") {
				t.Errorf(
					"expected SVG output for icon %s, got: %s",
					tc.name,
					output[:min(len(output), 100)],
				)
			}
			if !strings.Contains(output, "<path") {
				t.Errorf("expected path element for icon %s", tc.name)
			}
		})
	}
}

// --- Unknown Icon Behavior ---

func TestUnknownIconFallsBackGracefully(t *testing.T) {
	t.Parallel()

	t.Run("unknown icon name renders a fallback SVG", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon("nonexistent-icon", "h-5 w-5"))
		utils.AssertContains(t, output, "<svg")
	})
}

// --- Spinner Icon Behavior ---

func TestSpinnerIconIsAnimated(t *testing.T) {
	t.Parallel()

	t.Run("spinner icon has animation class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Spinner, "h-5 w-5"))
		utils.AssertContains(t, output, "animate-spin")
	})
}
