package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// FuzzInputType verifies inputType() never panics on arbitrary input.
func FuzzInputType(f *testing.F) {
	f.Add("text")
	f.Add("email")
	f.Add("")
	f.Add("invalid-type")
	f.Fuzz(func(t *testing.T, input string) {
		result := inputType(InputType(input))
		if result == "" {
			t.Errorf("inputType(%q) returned empty string", input)
		}
	})
}

// FuzzCombobox verifies Combobox renders without panicking on arbitrary
// string input and keeps the submission-critical invariants: the hidden
// value input stays present and the combobox role stays on the visible
// input, even when fuzzed names/values carry markup or quotes.
func FuzzCombobox(f *testing.F) {
	f.Add("country", "Country", "Search...", "de", "", false, false)
	f.Add("", "", "", "", "required", true, true)
	f.Add(`"><script>alert(1)</script>`, `L"abel`, "<b>x</b>", "val", "E", false, true)
	f.Add("name-with-'quote", "O'Brien", "it's", "x'y", "err'r", true, false)

	f.Fuzz(func(t *testing.T, name, label, placeholder, value, errMsg string, required, disabled bool) {
		output := utils.Render(t, Combobox(ComboboxProps{
			Name:        name,
			Label:       label,
			Placeholder: placeholder,
			Value:       value,
			Error:       errMsg,
			Required:    required,
			Disabled:    disabled,
			Options: []ComboboxOption{
				{Value: value, Label: label},
				{Value: "fixed", Label: "Fixed"},
			},
		}))

		utils.AssertContains(t, output, `type="hidden"`)
		utils.AssertContains(t, output, `role="combobox"`)
	})
}
