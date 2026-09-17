package visualtest

import (
	"testing"

	"github.com/larsartmann/templ-components/forms"
)

// TestToggleGoldens pins the Toggle switch at pixel level (TODO #255): the
// peer-checked thumb slide, the on/off track colors, dark-mode variants, and
// the disabled state. String/golden HTML tests cannot see whether the thumb
// actually renders in the right place — the August 2025 toggle regressions
// (invisible class, div-in-label) were both invisible to string tests.
func TestToggleGoldens(t *testing.T) {
	t.Parallel()

	off := forms.DefaultToggleProps()
	off.Name = "notifications"
	off.Label = "Notifications"

	on := off
	on.Checked = true

	disabled := off
	disabled.Disabled = true

	AssertScreenshot(t, "toggle/light", forms.Toggle(off), Options{})
	AssertScreenshot(t, "toggle/light_checked", forms.Toggle(on), Options{})
	AssertScreenshot(t, "toggle/dark", forms.Toggle(on), Options{Dark: Bool(true)})
	AssertScreenshot(t, "toggle/dark_disabled", forms.Toggle(disabled), Options{Dark: Bool(true)})
}
