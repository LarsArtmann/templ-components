// Package forms provides rendering tests for form components.
package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	snapshotInputNameEmail      = "email"
	snapshotCheckboxFieldTerms  = "terms"
	snapshotSelectFieldCountry  = "country"
	snapshotSelectLabelCountry  = "Country"
	snapshotSelectOptionGermany = "Germany"
	snapshotTextareaFieldBio    = "bio"
	snapshotTextareaLabelBio    = "Bio"
)

func TestInputRender(t *testing.T) {
	t.Parallel()
	t.Run("basic input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  snapshotInputNameEmail,
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
			BaseProps: utils.BaseProps{ID: snapshotInputNameEmail},
			Name:      snapshotInputNameEmail,
			Label:     "Email address",
			Error:     "Invalid email",
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
		BaseProps: utils.BaseProps{ID: snapshotCheckboxFieldTerms},
		Name:      snapshotCheckboxFieldTerms,
		Label:     "I agree",
	}))
	utils.AssertContains(t, output, `name="terms"`)
	utils.AssertContains(t, output, `id="terms"`)
	utils.AssertContains(t, output, "I agree")
}

func TestSelectRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Select(SelectProps{
		Name:  snapshotSelectFieldCountry,
		Label: snapshotSelectLabelCountry,
		Options: []SelectOption{
			{Value: "de", Label: snapshotSelectOptionGermany},
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
		Name:  snapshotTextareaFieldBio,
		Label: snapshotTextareaLabelBio,
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

func TestRadioRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		BaseProps: utils.BaseProps{ID: "plan-pro"},
		Name:      "plan",
		Value:     "pro",
		Label:     "Pro Plan",
	}))
	utils.AssertContains(t, output, `type="radio"`)
	utils.AssertContains(t, output, `name="plan"`)
	utils.AssertContains(t, output, `value="pro"`)
	utils.AssertContains(t, output, `id="plan-pro"`)
	utils.AssertContains(t, output, "Pro Plan")
}

func TestRadioGroupRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		BaseProps: utils.BaseProps{ID: "plan"},
		Name:      "plan",
		Label:     "Select a plan",
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro"},
		},
	}))
	utils.AssertContains(t, output, `<fieldset`)
	utils.AssertContains(t, output, "Select a plan")
	utils.AssertContains(t, output, `id="plan-free"`)
	utils.AssertContains(t, output, `id="plan-pro"`)
	utils.AssertContains(t, output, "Free")
	utils.AssertContains(t, output, "Pro")
}

func TestRadioGroupInlineRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		BaseProps: utils.BaseProps{ID: "plan"},
		Name:      "plan",
		Label:     "Select a plan",
		Inline:    true,
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
		},
	}))
	utils.AssertContains(t, output, `flex space-x-6`)
}

func TestRadioGroupWithErrorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		BaseProps: utils.BaseProps{ID: "plan"},
		Name:      "plan",
		Label:     "Select a plan",
		Error:     "Please select a plan",
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
		},
	}))
	utils.AssertContains(t, output, `aria-invalid="true"`)
	utils.AssertContains(t, output, "Please select a plan")
}
