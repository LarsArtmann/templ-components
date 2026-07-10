package display

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// TestDarkModeBadgeGolden verifies the badge renders correctly inside a
// .dark wrapper — this catches dark-mode class regressions.
func TestDarkModeBadgeGolden(t *testing.T) {
	t.Parallel()
	light := utils.Render(t, Badge(BadgeProps{Text: "Default", Type: BadgePrimary}))
	dark := "<div class=\"dark\">" + light + "</div>"
	golden.Assert(t, "badge_dark", dark)
}

// TestDarkModeCardGolden verifies the card renders correctly in dark mode.
func TestDarkModeCardGolden(t *testing.T) {
	t.Parallel()
	light := utils.Render(t, Card(CardProps{
		BaseProps: utils.BaseProps{},
		Title:     "Test Card",
	}))
	dark := "<div class=\"dark\">" + light + "</div>"
	golden.Assert(t, "card_dark", dark)
}

// TestDarkModeButtonGolden verifies the button renders correctly in dark mode.
func TestDarkModeButtonGolden(t *testing.T) {
	t.Parallel()
	light := utils.Render(t, Button(ButtonProps{
		Text:    "Click me",
		Variant: ButtonPrimary,
	}))
	dark := "<div class=\"dark\">" + light + "</div>"
	golden.Assert(t, "button_dark", dark)
}
