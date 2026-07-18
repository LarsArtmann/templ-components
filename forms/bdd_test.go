// Package forms provides behavior-driven tests for form components.
// These tests verify the end-user experience: filling forms, seeing errors, submitting data.
package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	inputNameEmail      = "email"
	inputLabelEmail     = "Email"
	inputNameUsername   = "username"
	inputLabelUsername  = "Username"
	selectFieldCountry  = "country"
	selectLabelCountry  = "Country"
	selectOptionGermany = "Germany"
	selectOptionRed     = "Red"
	selectOptionBlue    = "Blue"
	textareaFieldBio    = "bio"
	textareaLabelBio    = "Bio"
	checkboxFieldTerms  = "terms"
	checkboxLabelTerms  = "Terms"
	checkboxErrorAccept = "You must accept"
)

// --- Input Behavior ---

func TestInputUserCanEnterData(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled text input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(DefaultInputProps()))
		utils.AssertContains(t, output, `type="text"`)
	})

	t.Run("user sees labeled text input with custom props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:  inputNameUsername,
			Type:  InputText,
			Label: inputLabelUsername,
		}))
		utils.AssertContains(t, output, `name="username"`)
		utils.AssertContains(t, output, inputLabelUsername)
		utils.AssertContains(t, output, `type="text"`)
	})

	t.Run("user sees required field indicator", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     inputNameEmail,
			Type:     InputEmail,
			Label:    inputLabelEmail,
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
			Name:  inputNameEmail,
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

	t.Run("user sees input with placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:        "search",
			Type:        InputSearch,
			Label:       "Search",
			Placeholder: "Search...",
		}))
		utils.AssertContains(t, output, `placeholder="Search..."`)
	})

	t.Run("user sees input with autofocus", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:      "first",
			Type:      InputText,
			Label:     "First Name",
			AutoFocus: true,
		}))
		utils.AssertContains(t, output, `autofocus`)
	})

	t.Run("user sees input with both error and help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Input(InputProps{
			Name:     inputNameEmail,
			Type:     InputEmail,
			Label:    inputLabelEmail,
			Error:    "Invalid email",
			HelpText: "Use your work email.",
		}))
		utils.AssertContains(t, output, "Invalid email")
		utils.AssertContains(t, output, "Use your work email.")
	})
}

// --- Select Behavior ---

func TestSelectUserCanChooseOption(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled select with options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(DefaultSelectProps()))
		utils.AssertContains(t, output, "<select")
	})

	t.Run("user sees labeled select with custom options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			Name:  selectFieldCountry,
			Label: selectLabelCountry,
			Options: []SelectOption{
				{Value: "de", Label: selectOptionGermany},
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
				{Value: "red", Label: selectOptionRed},
				{Value: "blue", Label: selectOptionBlue, Selected: true},
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

	t.Run("user sees select with error", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			Name:  "role",
			Label: "Role",
			Error: "Please select a role",
			Options: []SelectOption{
				{Value: "admin", Label: "Admin"},
				{Value: "user", Label: "User"},
			},
		}))
		utils.AssertContains(t, output, "Please select a role")
		utils.AssertContains(t, output, `aria-invalid="true"`)
	})

	t.Run("user sees required select", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			Name:     "country",
			Label:    "Country",
			Required: true,
			Options: []SelectOption{
				{Value: "de", Label: selectOptionGermany},
			},
		}))
		utils.AssertContains(t, output, `required`)
	})

	t.Run("user sees select with help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Select(SelectProps{
			Name:     "timezone",
			Label:    "Timezone",
			HelpText: "Select your local timezone.",
			Options: []SelectOption{
				{Value: "utc", Label: "UTC"},
			},
		}))
		utils.AssertContains(t, output, "Select your local timezone.")
	})
}

// --- Textarea Behavior ---

func TestTextareaUserCanEnterMultiLineText(t *testing.T) {
	t.Parallel()

	t.Run("user sees default textarea", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(DefaultTextareaProps()))
		utils.AssertContains(t, output, "<textarea")
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

	t.Run("user sees textarea with error and required", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     textareaFieldBio,
			Label:    textareaLabelBio,
			Error:    "Biography is required",
			Required: true,
			Rows:     8,
		}))
		utils.AssertContains(t, output, "Biography is required")
		utils.AssertContains(t, output, `aria-invalid="true"`)
		utils.AssertContains(t, output, `required`)
		utils.AssertContains(t, output, `rows="8"`)
	})

	t.Run("user sees textarea with help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     "feedback",
			Label:    "Feedback",
			HelpText: "Tell us what you think.",
		}))
		utils.AssertContains(t, output, "Tell us what you think.")
	})

	t.Run("user sees disabled textarea", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     "locked",
			Label:    "Locked",
			Disabled: true,
		}))
		utils.AssertContains(t, output, `disabled`)
	})
}

// --- Checkbox Behavior ---

func TestCheckboxUserCanToggle(t *testing.T) {
	t.Parallel()

	t.Run("user sees labeled checkbox", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(DefaultCheckboxProps()))
		utils.AssertContains(t, output, `type="checkbox"`)
	})

	t.Run("user sees labeled checkbox with custom props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Checkbox(CheckboxProps{
			Name:  checkboxFieldTerms,
			Label: checkboxLabelTerms,
		}))
		utils.AssertContains(t, output, `name="terms"`)
		utils.AssertContains(t, output, checkboxLabelTerms)
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
		output := utils.Render(t, Label(inputNameUsername, inputLabelUsername, false))
		utils.AssertContains(t, output, `for="`+inputNameUsername+`"`)
		utils.AssertContains(t, output, inputLabelUsername)
	})

	t.Run("user sees required indicator on label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Label("name", "Full Name", true))
		utils.AssertContains(t, output, "Full Name")
	})
}

func TestFieldErrorUserSeesValidationFeedback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fieldID string
		message string
	}{
		{"linked to field", inputNameEmail, "Email is required"},
		{"standalone without field link", "", "Something went wrong"},
	}
	for _, tt := range tests {
		t.Run("user sees "+tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, FieldError(tt.fieldID, tt.message))
			utils.AssertContains(t, output, tt.message)
		})
	}
}

func TestSelectEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		props    SelectProps
		contains []string
	}{
		{
			"disabled select",
			SelectProps{Name: "role", Label: "Role", Disabled: true, Options: []SelectOption{{Value: "admin", Label: "Admin"}}},
			[]string{"disabled"},
		},
		{
			"select with disabled option",
			SelectProps{Name: "color", Label: "Color", Options: []SelectOption{{Value: "red", Label: selectOptionRed}, {Value: "blue", Label: selectOptionBlue, Disabled: true}}},
			[]string{selectOptionRed, selectOptionBlue},
		},
		{
			"select with pre-selected option",
			SelectProps{Name: "size", Label: "Size", Options: []SelectOption{{Value: "sm", Label: "Small"}, {Value: "md", Label: "Medium", Selected: true}}},
			[]string{"selected"},
		},
		{
			"select with error and help text",
			SelectProps{Name: "country", Label: "Country", Error: "Required", HelpText: "Select your country", Options: []SelectOption{{Value: "us", Label: "US"}}},
			[]string{"Required", "Select your country", ariaInvalid},
		},
		{
			"select with no options",
			SelectProps{Name: "empty", Label: "Empty", Options: []SelectOption{}},
			[]string{`<select`},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, Select(tt.props))
			for _, want := range tt.contains {
				utils.AssertContains(t, output, want)
			}
		})
	}
}

func TestCheckboxEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		props    CheckboxProps
		contains []string
	}{
		{
			"checked checkbox",
			CheckboxProps{Name: "agree", Label: "I agree", Checked: true},
			[]string{"checked"},
		},
		{
			"checkbox with help text",
			CheckboxProps{Name: "newsletter", Label: "Newsletter", HelpText: "We send weekly updates"},
			[]string{"We send weekly updates"},
		},
		{
			"checkbox with error",
			CheckboxProps{Name: checkboxFieldTerms, Label: checkboxLabelTerms, Error: checkboxErrorAccept},
			[]string{checkboxErrorAccept, ariaInvalid},
		},
		{
			"disabled checkbox",
			CheckboxProps{Name: "locked", Label: "Locked", Disabled: true},
			[]string{"disabled"},
		},
		{
			"checkbox with error and help text",
			CheckboxProps{Name: checkboxFieldTerms, Label: checkboxLabelTerms, Error: checkboxErrorAccept, HelpText: "Check to agree"},
			[]string{checkboxErrorAccept, "Check to agree", ariaInvalid},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, Checkbox(tt.props))
			for _, want := range tt.contains {
				utils.AssertContains(t, output, want)
			}
		})
	}
}

func TestLabelEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name        string
		id          string
		text        string
		contains    string
		notContains string
	}{
		{"label with help text only", "username", "Username", "Username", "text-red"},
		{"label without for ID", "", "Standalone", "Standalone", `for=""`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Label(tt.id, tt.text, false))
			utils.AssertContains(t, output, tt.contains)
			utils.AssertNotContains(t, output, tt.notContains)
		})
	}
}

func TestTextareaEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("textarea with custom rows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:  "bio",
			Label: "Bio",
			Rows:  8,
		}))
		utils.AssertContains(t, output, `rows="8"`)
	})

	t.Run("textarea readonly", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     "terms",
			Label:    "Terms",
			ReadOnly: true,
			Value:    "Read-only content",
		}))
		utils.AssertContains(t, output, "readonly")
	})

	t.Run("textarea with placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:        "note",
			Label:       "Note",
			Placeholder: "Type here...",
		}))
		utils.AssertContains(t, output, "Type here...")
	})

	t.Run("textarea with error and help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     "bio",
			Label:    "Bio",
			Error:    "Too short",
			HelpText: "50 chars minimum",
		}))
		utils.AssertContains(t, output, "Too short")
		utils.AssertContains(t, output, "50 chars minimum")
		utils.AssertContains(t, output, ariaInvalid)
	})
}
