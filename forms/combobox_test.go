package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestComboboxRender(t *testing.T) {
	t.Parallel()

	t.Run("basic combobox with options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps:   utils.BaseProps{ID: "country"},
			Name:        "country",
			Label:       "Country",
			Placeholder: "Search...",
			Options: []ComboboxOption{
				{Value: "de", Label: "Germany"},
				{Value: "at", Label: "Austria"},
				{Value: "ch", Label: "Switzerland"},
			},
		}))
		utils.AssertContains(t, output, `role="combobox"`)
		utils.AssertContains(t, output, `aria-autocomplete="list"`)
		utils.AssertContains(t, output, "Germany")
		utils.AssertContains(t, output, "Austria")
		utils.AssertContains(t, output, `type="hidden"`)
		utils.AssertContains(t, output, `name="country"`)
		utils.AssertContains(t, output, `data-combobox-input="country"`)
		utils.AssertContains(t, output, `id="country-listbox"`)
		utils.AssertContains(t, output, `role="listbox"`)
		utils.AssertContains(t, output, `role="option"`)
	})

	t.Run("with preselected value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "city"},
			Name:      "city",
			Value:     "Berlin",
			Options: []ComboboxOption{
				{Value: "berlin", Label: "Berlin"},
				{Value: "munich", Label: "Munich"},
			},
		}))
		utils.AssertContains(t, output, `value="Berlin"`)
	})

	t.Run("with error and help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "tag"},
			Name:      "tag",
			Label:     "Tag",
			Required:  true,
			Error:     "Please select a tag",
			HelpText:  "Choose from the list",
			Options:   []ComboboxOption{{Value: "go", Label: "Go"}},
		}))
		utils.AssertContains(t, output, `required`)
		utils.AssertContains(t, output, `aria-required="true"`)
		utils.AssertContains(t, output, "Please select a tag")
		utils.AssertContains(t, output, "Choose from the list")
	})

	t.Run("with nonce for CSP", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "cb", Nonce: "nonce-cb"},
			Name:      "search",
			Options:   []ComboboxOption{{Value: "a", Label: "A"}},
		}))
		utils.AssertContains(t, output, `nonce="nonce-cb"`)
		utils.AssertContains(t, output, "tcComboboxAttached")
	})

	t.Run("with BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{
				ID:        "cb-custom",
				Class:     "w-96",
				AriaLabel: "Search items",
				Attrs:     templ.Attributes{"data-test": "cb"},
			},
			Name:    "items",
			Options: []ComboboxOption{{Value: "x", Label: "X"}},
		}))
		utils.AssertContains(t, output, `w-96`)
		utils.AssertContains(t, output, `aria-label="Search items"`)
		utils.AssertContains(t, output, `data-test="cb"`)
	})

	t.Run("auto-generates ID when not provided", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			Name:    "no-id",
			Options: []ComboboxOption{{Value: "y", Label: "Y"}},
		}))
		utils.AssertContains(t, output, `id="tc-combobox-`)
		utils.AssertContains(t, output, `data-combobox-input="tc-combobox-`)
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(ComboboxProps{
			BaseProps: utils.BaseProps{ID: "disabled-cb"},
			Name:      "disabled_cb",
			Disabled:  true,
			Options:   []ComboboxOption{{Value: "z", Label: "Z"}},
		}))
		utils.AssertContains(t, output, `disabled`)
		// The hidden submission input must also be disabled so its value is
		// excluded from form submission (HTML spec: disabled controls are not submitted).
		utils.AssertContains(t, output,
			`type="hidden" name="disabled_cb" value="" data-combobox-value="disabled-cb" disabled`)
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Combobox(DefaultComboboxProps()))
		utils.AssertContains(t, output, `role="combobox"`)
	})
}
