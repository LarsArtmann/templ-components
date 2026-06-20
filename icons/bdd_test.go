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

	for _, name := range allIconNames() {
		t.Run("icon "+string(name)+" renders SVG", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Icon(name, "h-5 w-5"))
			if !strings.Contains(output, "<svg") {
				t.Errorf(
					"expected SVG output for icon %s, got: %s",
					name,
					output[:min(len(output), 100)],
				)
			}
			if !strings.Contains(output, "<path") && name != Spinner {
				t.Errorf("expected path element for icon %s", name)
			}
		})
	}
}

// --- Unknown Icon Behavior ---

func TestUnknownIconFallback(t *testing.T) {
	t.Parallel()

	t.Run("unknown icon name falls back to Question icon without panic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon("nonexistent-icon", "h-5 w-5"))
		questionOutput := utils.Render(t, Icon(Question, "h-5 w-5"))
		if output != questionOutput {
			t.Error("unknown icon should render the Question fallback icon")
		}
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
