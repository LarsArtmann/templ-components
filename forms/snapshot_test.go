// Package forms provides rendering tests for form components.
package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestInputRender(t *testing.T) {
	t.Parallel()
	t.Run("basic input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  "email",
			Type:  InputEmail,
			Label: "Email address",
		}))
		utils.AssertContains(t, output, `name="email"`)
		utils.AssertContains(t, output, `type="email"`)
		utils.AssertContains(t, output, "Email address")
	})

	t.Run("input with error", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:   "email",
			ID:     "email",
			Label:  "Email address",
			Error:  "Invalid email",
		}))
		utils.AssertContains(t, output, `aria-invalid="true"`)
		utils.AssertContains(t, output, `aria-describedby="email-error"`)
		utils.AssertContains(t, output, "Invalid email")
	})

	t.Run("input without id skips label for", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  "search",
			Label: "Search",
		}))
		utils.AssertNotContains(t, output, `for=""`)
	})
}

func TestCheckboxRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Name:  "terms",
		ID:    "terms",
		Label: "I agree",
	}))
	utils.AssertContains(t, output, `name="terms"`)
	utils.AssertContains(t, output, `id="terms"`)
	utils.AssertContains(t, output, "I agree")
}

func TestSelectRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Select(SelectProps{
		Name:  "country",
		Label: "Country",
		Options: []SelectOption{
			{Value: "de", Label: "Germany"},
			{Value: "at", Label: "Austria"},
		},
	}))
	utils.AssertContains(t, output, `name="country"`)
	utils.AssertContains(t, output, "Country")
	utils.AssertContains(t, output, "Germany")
	utils.AssertContains(t, output, "Austria")
}

func TestTextareaRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Textarea(TextareaProps{
		Name:  "bio",
		Label: "Bio",
		Rows:  4,
	}))
	utils.AssertContains(t, output, `name="bio"`)
	utils.AssertContains(t, output, "Bio")
	utils.AssertContains(t, output, `rows="4"`)
}

func TestLabelRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Label("email", "Email address", true))
	utils.AssertContains(t, output, `for="email"`)
	utils.AssertContains(t, output, "Email address")
	utils.AssertContains(t, output, "*")
}

func TestLabelWithoutForID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Label("", "No for", false))
	utils.AssertNotContains(t, output, `for=""`)
	utils.AssertContains(t, output, "No for")
}

func TestFieldErrorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FieldError("email", "This field is required"))
	utils.AssertContains(t, output, `id="email-error"`)
	utils.AssertContains(t, output, "This field is required")
}

func TestFieldErrorWithoutID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FieldError("", "Generic error"))
	utils.AssertContains(t, output, "Generic error")
	utils.AssertNotContains(t, output, `id="-error"`)
}
