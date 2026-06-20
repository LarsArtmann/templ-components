package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestInputFullCoverage(t *testing.T) {
	t.Parallel()
	for _, it := range []InputType{InputText, InputEmail, InputPassword, InputNumber, InputTel, InputURL, InputSearch} {
		t.Run(string(it)+" with error and help", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Input(InputProps{
				BaseProps:   utils.BaseProps{ID: "in-" + string(it), Class: "w-full", AriaLabel: string(it)},
				Type:        it,
				Name:        "field_" + string(it),
				Label:       string(it) + " field",
				Placeholder: "Enter " + string(it),
				Required:    true,
				Error:       "Invalid value",
				HelpText:    "Must be valid " + string(it),
				MaxLength:   100,
			}))
			utils.AssertContains(t, output, string(it)+" field")
		})
	}
	t.Run("disabled and readonly", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     "locked",
			Label:    "Locked",
			Disabled: true,
			Value:    "fixed",
		}))
		utils.AssertContains(t, output, `disabled`)
	})
}

func TestCheckboxFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("checked with error", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			BaseProps: utils.BaseProps{ID: "cb-1", AriaLabel: "Accept terms"},
			Name:      "agree",
			Label:     "I agree",
			Checked:   true,
			Required:  true,
			Error:     "You must agree",
			HelpText:  "Required to continue",
		}))
		utils.AssertContains(t, output, "I agree")
		utils.AssertContains(t, output, `checked`)
	})
}

func TestRadioGroupFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("inline with options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, RadioGroup(RadioGroupProps{
			BaseProps: utils.BaseProps{ID: "rg-1", AriaLabel: "Choose color"},
			Name:      "color",
			Label:     "Color",
			Inline:    true,
			Required:  true,
			Error:     "Select a color",
			Options: []RadioOption{
				{Value: "red", Label: "Red"},
				{Value: "blue", Label: "Blue", Disabled: true},
			},
		}))
		utils.AssertContains(t, output, "Red")
	})
}

func TestSelectFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with options and error", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			BaseProps: utils.BaseProps{ID: "sel-1", AriaLabel: "Choose role"},
			Name:      "role",
			Label:     "Role",
			Required:  true,
			Error:     "Select required",
			HelpText:  "Pick one",
			Options: []SelectOption{
				{Value: "admin", Label: "Admin"},
				{Value: "user", Label: "User", Disabled: true},
				{Value: "guest", Label: "Guest", Selected: true},
			},
		}))
		utils.AssertContains(t, output, "Admin")
		utils.AssertContains(t, output, "selected")
	})
}

func TestTextareaComprehensive(t *testing.T) {
	t.Parallel()
	t.Run("with error and maxlength", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			BaseProps:   utils.BaseProps{ID: "ta-1", AriaLabel: "Bio"},
			Name:        "bio",
			Label:       "Bio",
			Placeholder: "Tell us about yourself",
			Rows:        5,
			Required:    true,
			MaxLength:   500,
			Error:       "Too short",
			HelpText:    "Min 10 chars",
		}))
		utils.AssertContains(t, output, "Bio")
		utils.AssertContains(t, output, `rows="5"`)
	})
}

func TestToggleFullCoverage(t *testing.T) {
	t.Parallel()
	for _, size := range []ToggleSize{ToggleSizeSM, ToggleSizeMD, ToggleSizeLG} {
		t.Run("size_"+string(size)+" checked", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Toggle(ToggleProps{
				BaseProps: utils.BaseProps{ID: "tg-" + string(size), AriaLabel: "Notifications"},
				Name:      "notifications",
				Label:     "Enable notifications",
				Checked:   true,
				Size:      size,
			}))
			utils.AssertContains(t, output, "Enable notifications")
		})
	}
}

func TestFileInputFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with multiple and accept", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FileInput(FileInputProps{
			BaseProps: utils.BaseProps{ID: "fi-1", AriaLabel: "Upload"},
			Name:      "documents",
			Label:     "Documents",
			Required:  true,
			Accept:    ".pdf,.docx",
			Multiple:  true,
			Error:     "File too large",
			HelpText:  "Max 10MB",
		}))
		utils.AssertContains(t, output, "Documents")
		utils.AssertContains(t, output, `accept=".pdf,.docx"`)
		utils.AssertContains(t, output, `multiple`)
	})
}

func TestValidationSummaryFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with multiple errors", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
			BaseProps: utils.BaseProps{ID: "vs-1", AriaLabel: "Errors"},
			Errors: []ValidationError{
				{Field: "email", Message: "Email is required"},
				{Field: "password", Message: "Password too short"},
			},
		}))
		utils.AssertContains(t, output, "Email is required")
		utils.AssertContains(t, output, `role="alert"`)
	})
}

func TestFormFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("POST with CSRF", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{ //nolint:gosec // G101 false positive on test CSRF token
			BaseProps: utils.BaseProps{
				ID:        "form-1",
				AriaLabel: "Contact form",
				Attrs:     templ.Attributes{"data-track": "submit"},
			},
			Action:    "/submit",
			Method:    FormPost,
			CSRFToken: "test-csrf-token",
		}))
		utils.AssertContains(t, output, `action="/submit"`)
		utils.AssertContains(t, output, `method="POST"`)
		utils.AssertContains(t, output, "test-csrf-token")
		utils.AssertContains(t, output, `data-track="submit"`)
	})
}
