package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestInputIdentifierTextCorrection (Modern Web Guidance: forms §3) pins
// the autocorrect/autocapitalize/spellcheck trio on identifier fields:
// identifier TYPES always, identifier AUTOCOMPLETE tokens on text inputs,
// and never on plain prose-capable text inputs.
func TestInputIdentifierTextCorrection(t *testing.T) {
	t.Parallel()

	identifierTypes := []InputType{InputEmail, InputURL, InputTel, InputSearch, InputNumber}
	for _, typ := range identifierTypes {
		output := utils.Render(t, Input(InputProps{Type: typ, Name: "f"}))
		utils.AssertContains(t, output, `autocorrect="off"`)
		utils.AssertContains(t, output, `autocapitalize="off"`)
		utils.AssertContains(t, output, `spellcheck="false"`)
	}

	for _, token := range []AutocompleteType{AutocompleteUsername, AutocompleteOneTimeCode} {
		output := utils.Render(t, Input(InputProps{Type: InputText, Name: "f", Autocomplete: token}))
		utils.AssertContains(t, output, `spellcheck="false"`)
	}

	prose := utils.Render(t, Input(InputProps{Type: InputText, Name: "f", Autocomplete: AutocompleteOff}))
	utils.AssertNotContains(t, prose, `autocorrect="off"`)
	utils.AssertNotContains(t, prose, `spellcheck="false"`)
}
