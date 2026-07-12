package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// Coverage: exercise low-coverage branches in forms package.
// Input (67.1%), helpText (68.1%), checkboxLabel (68.3%), Combobox (71.1%).

func TestInputCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("input with help text and no error", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			BaseProps: utils.BaseProps{ID: "name"},
			Name:      "name",
			Label:     "Name",
			HelpText:  "Enter your full name",
		}))
		utils.AssertContains(t, output, "Enter your full name")
		utils.AssertContains(t, output, `aria-describedby="name-help"`)
	})

	t.Run("input with placeholder only", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:        "email",
			Type:        InputEmail,
			Placeholder: "you@example.com",
		}))
		utils.AssertContains(t, output, `placeholder="you@example.com"`)
	})

	t.Run("input disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     "locked",
			Disabled: true,
		}))
		utils.AssertContains(t, output, "disabled")
	})

	t.Run("input readonly", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     "ro",
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, "readonly")
	})

	t.Run("input autofocus", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:      "focus",
			AutoFocus: true,
		}))
		utils.AssertContains(t, output, "autofocus")
	})

	t.Run("input maxlength", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:      "limited",
			MaxLength: 50,
		}))
		utils.AssertContains(t, output, `maxlength="50"`)
	})

	t.Run("input with value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  "prefilled",
			Value: "existing value",
		}))
		utils.AssertContains(t, output, `value="existing value"`)
	})

	t.Run("input required", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			BaseProps: utils.BaseProps{ID: "req"},
			Name:      "req",
			Label:     "Required",
			Required:  true,
		}))
		utils.AssertContains(t, output, "required")
	})

	t.Run("input with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			BaseProps: utils.BaseProps{AriaLabel: "Search field"},
			Name:      "q",
			Type:      InputSearch,
		}))
		utils.AssertContains(t, output, `aria-label="Search field"`)
	})

	t.Run("input with custom class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			BaseProps: utils.BaseProps{Class: "my-input"},
			Name:      "x",
		}))
		utils.AssertContains(t, output, "my-input")
	})

	t.Run("input with both error and help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			BaseProps: utils.BaseProps{ID: "dual"},
			Name:      "dual",
			Label:     "Dual",
			Error:     "Wrong",
			HelpText:  "Try again",
		}))
		utils.AssertContains(t, output, "Wrong")
		utils.AssertContains(t, output, "Try again")
	})
}

func TestCheckboxCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("checkbox with help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			BaseProps: utils.BaseProps{ID: "terms"},
			Name:      "terms",
			Label:     "I agree",
			HelpText:  "You must agree to continue",
		}))
		utils.AssertContains(t, output, "You must agree to continue")
	})

	t.Run("checkbox checked", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			BaseProps: utils.BaseProps{ID: "sub"},
			Name:      "sub",
			Label:     "Subscribe",
			Checked:   true,
		}))
		utils.AssertContains(t, output, "checked")
	})

	t.Run("checkbox disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			BaseProps: utils.BaseProps{ID: "locked"},
			Name:      "locked",
			Label:     "Locked",
			Disabled:  true,
		}))
		utils.AssertContains(t, output, "disabled")
	})

	t.Run("checkbox without ID renders without label-for", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			Name:  "noid",
			Label: "No ID",
		}))
		utils.AssertNotContains(t, output, `for=""`)
	})
}

func TestHelpTextRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, helpText("myfield", "This is help"))
	utils.AssertContains(t, output, "This is help")
	utils.AssertContains(t, output, `id="myfield-help"`)
}

func TestComboboxCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("combobox with placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps:   utils.BaseProps{ID: "cb", Nonce: "test-nonce"},
			Name:        "country",
			Label:       "Country",
			Placeholder: "Select a country...",
			Options: []ComboboxOption{
				{Value: "de", Label: "Germany"},
				{Value: "at", Label: "Austria"},
			},
		}))
		utils.AssertContains(t, output, "Select a country...")
	})

	t.Run("combobox disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "cb-disabled", Nonce: "n"},
			Name:      "locked",
			Label:     "Locked",
			Disabled:  true,
			Options: []ComboboxOption{
				{Value: "x", Label: "X"},
			},
		}))
		utils.AssertContains(t, output, "disabled")
	})

	t.Run("combobox with selected value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "cb-selected", Nonce: "n"},
			Name:      "color",
			Value:     "blue",
			Options: []ComboboxOption{
				{Value: "red", Label: "Red"},
				{Value: "blue", Label: "Blue"},
			},
		}))
		utils.AssertContains(t, output, "Blue")
	})

	t.Run("combobox with help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "cb-help", Nonce: "n"},
			Name:      "city",
			Label:     "City",
			HelpText:  "Start typing to search",
			Options: []ComboboxOption{
				{Value: "nyc", Label: "New York"},
			},
		}))
		utils.AssertContains(t, output, "Start typing to search")
	})
}
