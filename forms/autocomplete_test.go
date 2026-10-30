package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestInputAutocomplete covers the autocomplete attribute contract: recognized
// tokens reach the DOM, an unset field renders no attribute at all (an empty
// autocomplete="" would tell browsers "no idea", not "unspecified"), and an
// unrecognized token is dropped rather than emitted, because browsers silently
// ignore it and the form would look autofill-ready while never autofilling.
func TestInputAutocomplete(t *testing.T) {
	t.Parallel()

	t.Run("valid sign-in tokens render", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Input(InputProps{
			ID:           "ac-user",
			Name:         "username",
			Type:         InputText,
			Autocomplete: AutocompleteUser,
		}))
		utils.AssertContains(t, output, `autocomplete="username"`)

		output = utils.Render(t, Input(InputProps{
			ID:           "ac-pass",
			Name:         "password",
			Type:         InputPassword,
			Autocomplete: AutocompleteCurrentPassword,
		}))
		utils.AssertContains(t, output, `autocomplete="current-password"`)
	})

	t.Run("payment and address tokens render", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Input(InputProps{
			ID:           "ac-cc",
			Name:         "card",
			Type:         InputText,
			Autocomplete: AutocompleteCCNumber,
		}))
		utils.AssertContains(t, output, `autocomplete="cc-number"`)

		output = utils.Render(t, Input(InputProps{
			ID:           "ac-zip",
			Name:         "postal",
			Type:         InputText,
			Autocomplete: AutocompletePostalCode,
		}))
		utils.AssertContains(t, output, `autocomplete="postal-code"`)
	})

	t.Run("unset omits the attribute", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Input(InputProps{ID: "ac-none", Name: "q", Type: InputText}))
		utils.AssertNotContains(t, output, "autocomplete=")
	})

	t.Run("unknown token is dropped instead of emitted", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Input(InputProps{
			ID:           "ac-bogus",
			Name:         "q",
			Type:         InputText,
			Autocomplete: AutocompleteType("not-a-real-token"),
		}))
		utils.AssertNotContains(t, output, "autocomplete=")
	})
}
