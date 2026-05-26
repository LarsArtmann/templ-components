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

func TestToggleRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		BaseProps: utils.BaseProps{ID: "notifications"},
		Name:      "notifications",
		Label:     "Enable notifications",
		Checked:   true,
	}))
	utils.AssertContains(t, output, `type="checkbox"`)
	utils.AssertContains(t, output, `name="notifications"`)
	utils.AssertContains(t, output, `id="notifications"`)
	utils.AssertContains(t, output, "Enable notifications")
	utils.AssertContains(t, output, `checked`)
	utils.AssertContains(t, output, `sr-only`)
}

func TestToggleSizesRender(t *testing.T) {
	t.Parallel()
	t.Run("small toggle", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toggle(ToggleProps{
			Name: "sm",
			Size: ToggleSizeSM,
		}))
		utils.AssertContains(t, output, `w-9`)
	})
	t.Run("medium toggle", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toggle(ToggleProps{
			Name: "md",
			Size: ToggleSizeMD,
		}))
		utils.AssertContains(t, output, `w-11`)
	})
	t.Run("large toggle", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toggle(ToggleProps{
			Name: "lg",
			Size: ToggleSizeLG,
		}))
		utils.AssertContains(t, output, `w-14`)
	})
}

func TestFileInputRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		BaseProps: utils.BaseProps{ID: "avatar"},
		Name:      "avatar",
		Label:     "Upload avatar",
		Accept:    "image/*",
	}))
	utils.AssertContains(t, output, `type="file"`)
	utils.AssertContains(t, output, `name="avatar"`)
	utils.AssertContains(t, output, `id="avatar"`)
	utils.AssertContains(t, output, "Upload avatar")
	utils.AssertContains(t, output, `accept="image/*"`)
	utils.AssertContains(t, output, `file:bg-blue-600`)
}

func TestFileInputMultipleRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		Name:     "docs",
		Label:    "Upload documents",
		Multiple: true,
	}))
	utils.AssertContains(t, output, `multiple`)
}

func TestInputGroupRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InputGroup(InputGroupProps{
		BaseProps: utils.BaseProps{ID: "search-group"},
		LeftAddon: Input(InputProps{BaseProps: utils.BaseProps{Class: "pl-10"}}),
	}))
	utils.AssertContains(t, output, `id="search-group"`)
	utils.AssertContains(t, output, `pointer-events-none`)
}

func TestInputGroupPaddingClass(t *testing.T) {
	t.Parallel()
	if got := InputGroupPaddingClass(true, false); got != "pl-10" {
		t.Errorf("InputGroupPaddingClass(true,false) = %q, want pl-10", got)
	}
	if got := InputGroupPaddingClass(false, true); got != "pr-10" {
		t.Errorf("InputGroupPaddingClass(false,true) = %q, want pr-10", got)
	}
	if got := InputGroupPaddingClass(true, true); got != "pl-10 pr-10" {
		t.Errorf("InputGroupPaddingClass(true,true) = %q, want pl-10 pr-10", got)
	}
	if got := InputGroupPaddingClass(false, false); got != "" {
		t.Errorf("InputGroupPaddingClass(false,false) = %q, want empty", got)
	}
}
