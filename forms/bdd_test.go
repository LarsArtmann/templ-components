// Package forms provides behavior-driven tests for form components.
// These tests verify the end-user experience: filling forms, seeing errors, submitting data.
package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- Input Behavior ---

func TestInputUserCanEnterData(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled text input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  "username",
			Type:  InputText,
			Label: "Username",
		}))
		utils.AssertContains(t, output, `name="username"`)
		utils.AssertContains(t, output, "Username")
		utils.AssertContains(t, output, `type="text"`)
	})

	t.Run("user sees required field indicator", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     "email",
			Type:     InputEmail,
			Label:    "Email",
			Required: true,
		}))
		utils.AssertContains(t, output, `required`)
	})

	t.Run("user sees pre-filled input value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  "city",
			Type:  InputText,
			Label: "City",
			Value: "Berlin",
		}))
		utils.AssertContains(t, output, `value="Berlin"`)
	})

	t.Run("user sees help text below input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     "password",
			Type:     InputPassword,
			Label:    "Password",
			HelpText: "Must be at least 8 characters.",
		}))
		utils.AssertContains(t, output, "Must be at least 8 characters.")
	})

	t.Run("user sees field error with accessible attributes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  "email",
			Type:  InputEmail,
			Label: "Email",
			Error: "Email is required",
		}))
		utils.AssertContains(t, output, "Email is required")
		utils.AssertContains(t, output, `aria-invalid="true"`)
	})

	t.Run("user sees disabled input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     "readonly-field",
			Type:     InputText,
			Label:    "Read Only",
			Disabled: true,
		}))
		utils.AssertContains(t, output, `disabled`)
	})
}

// --- Select Behavior ---

func TestSelectUserCanChooseOption(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled select with options", func(t *testing.T) {
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
		utils.AssertContains(t, output, "Germany")
		utils.AssertContains(t, output, "Austria")
	})

	t.Run("user sees pre-selected option", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			Name:  "color",
			Label: "Color",
			Options: []SelectOption{
				{Value: "red", Label: "Red"},
				{Value: "blue", Label: "Blue", Selected: true},
			},
		}))
		utils.AssertContains(t, output, `selected`)
	})

	t.Run("user sees disabled option", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			Name:  "plan",
			Label: "Plan",
			Options: []SelectOption{
				{Value: "free", Label: "Free"},
				{Value: "pro", Label: "Pro (Coming Soon)", Disabled: true},
			},
		}))
		utils.AssertContains(t, output, `disabled`)
	})
}

// --- Textarea Behavior ---

func TestTextareaUserCanEnterMultiLineText(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled textarea", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:  "bio",
			Label: "Biography",
			Rows:  4,
		}))
		utils.AssertContains(t, output, `name="bio"`)
		utils.AssertContains(t, output, "Biography")
		utils.AssertContains(t, output, `rows="4"`)
	})

	t.Run("user sees textarea with pre-filled value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:  "notes",
			Label: "Notes",
			Value: "Existing note content",
		}))
		utils.AssertContains(t, output, "Existing note content")
	})
}

// --- Checkbox Behavior ---

func TestCheckboxUserCanToggle(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled checkbox", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			Name:  "terms",
			Label: "I agree to the terms",
		}))
		utils.AssertContains(t, output, `name="terms"`)
		utils.AssertContains(t, output, "I agree to the terms")
		utils.AssertContains(t, output, `type="checkbox"`)
	})

	t.Run("user sees pre-checked checkbox", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			Name:    "newsletter",
			Label:   "Subscribe to newsletter",
			Checked: true,
		}))
		utils.AssertContains(t, output, `checked`)
	})
}

// --- Label and FieldError Behavior ---

func TestLabelUserSeesFieldLabels(t *testing.T) {
	t.Parallel()

	t.Run("user sees label linked to input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Label("email", "Email address", false))
		utils.AssertContains(t, output, `for="email"`)
		utils.AssertContains(t, output, "Email address")
	})

	t.Run("user sees required indicator on label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Label("name", "Full Name", true))
		utils.AssertContains(t, output, "Full Name")
	})
}

func TestFieldErrorUserSeesValidationFeedback(t *testing.T) {
	t.Parallel()

	t.Run("user sees error message linked to field", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FieldError("email", "Email is required"))
		utils.AssertContains(t, output, "Email is required")
	})

	t.Run("user sees standalone error without field link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FieldError("", "Something went wrong"))
		utils.AssertContains(t, output, "Something went wrong")
	})
}
