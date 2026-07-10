package layout

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestThemeScriptHasNonce verifies CSP nonce is present on the script tag.
func TestThemeScriptHasNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript("test-nonce-123"))
	utils.AssertContains(t, output, `nonce="test-nonce-123"`)
}

// TestThemeScriptHasLocalStorageTryCatch verifies localStorage is wrapped in
// try/catch for Safari private mode (QuotaExceededError).
func TestThemeScriptHasLocalStorageTryCatch(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript(""))
	if !strings.Contains(output, "try") || !strings.Contains(output, "catch") {
		t.Error("expected localStorage wrapped in try/catch")
	}
}

// TestThemeScriptHasColorScheme verifies color-scheme is set.
func TestThemeScriptHasColorScheme(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript(""))
	utils.AssertContains(t, output, "colorScheme")
}

// TestThemeToggleHasRoleSwitch verifies the toggle uses role="switch".
func TestThemeToggleHasRoleSwitch(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("", ""))
	utils.AssertContains(t, output, `role="switch"`)
	utils.AssertContains(t, output, `aria-checked="false"`)
	utils.AssertContains(t, output, "data-theme-toggle")
}

// TestThemeToggleDefaultLabel verifies the default aria-label.
func TestThemeToggleDefaultLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("", ""))
	utils.AssertContains(t, output, "Toggle theme")
}

// TestThemeToggleCustomAriaLabel verifies custom aria-label override.
func TestThemeToggleCustomAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("Switch appearance", ""))
	utils.AssertContains(t, output, "Switch appearance")
}

// TestThemeToggleScriptUsesQuerySelectorAll verifies the toggle syncs all
// instances (not just the first one).
func TestThemeToggleScriptUsesQuerySelectorAll(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("", ""))
	utils.AssertContains(t, output, "querySelectorAll")
}

// TestThemeScriptBeforeHTMX verifies FOUC prevention — the theme script must
// be self-contained (no external dependencies).
func TestThemeScriptSelfContained(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript(""))
	// Must NOT reference external scripts — it runs before page paint
	utils.AssertNotContains(t, output, "src=")
}
