package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ---------------------------------------------------------------------------
// Input: all types, full props, hidden input, BaseProps propagation
// ---------------------------------------------------------------------------

func TestInputAllTypes(t *testing.T) {
	t.Parallel()
	for _, itype := range []InputType{
		InputText, InputEmail, InputPassword, InputNumber,
		InputTel, InputURL, InputSearch,
	} {
		output := utils.Render(t, Input(InputProps{
			Name:  "field",
			Type:  itype,
			Label: "Field",
		}))
		utils.AssertContains(t, output, `type="`+string(itype)+`"`)
	}
}

func TestInputHidden(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		Name:  "csrf",
		Type:  InputHidden,
		Value: "token123",
	}))
	utils.AssertContains(t, output, `type="hidden"`)
	utils.AssertContains(t, output, `value="token123"`)
}

func TestInputFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		Name:        "username",
		Type:        InputText,
		Label:       "Username",
		Value:       "alice",
		Placeholder: "Enter username",
		Required:    true,
		Disabled:    true,
		ReadOnly:    true,
		AutoFocus:   true,
		MaxLength:   50,
		Error:       "Username taken",
		HelpText:    "3-50 characters",
		BaseProps: utils.BaseProps{
			ID:        "input-username",
			AriaLabel: "Username field",
			Class:     "extra-class",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="input-username"`,
		`name="username"`,
		`value="alice"`,
		"Username taken",
		"3-50 characters",
		"required",
		"disabled",
		"readonly",
		"autofocus",
	)
}

func TestInputErrorAndHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		Name:     "email",
		Label:    "Email",
		Error:    "Invalid email",
		HelpText: "e.g. user@example.com",
	}))
	utils.AssertContainsAll(t, output, "Invalid email", "e.g. user@example.com")
}

// ---------------------------------------------------------------------------
// Toggle: all sizes, full props, disabled, checked
// ---------------------------------------------------------------------------

func TestToggleAllSizes(t *testing.T) {
	t.Parallel()
	for _, size := range []ToggleSize{ToggleSizeSM, ToggleSizeMD, ToggleSizeLG} {
		output := utils.Render(t, Toggle(ToggleProps{
			Label: "Switch",
			Size:  size,
		}))
		utils.AssertContains(t, output, "Switch")
	}
}

func TestToggleFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		Name:     "notifications",
		Label:    "Enable notifications",
		Checked:  true,
		Required: true,
		Error:    "Must accept to continue",
		HelpText: "We'll send you updates",
		BaseProps: utils.BaseProps{
			ID:        "toggle-1",
			AriaLabel: "Notifications toggle",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="toggle-1"`,
		"Enable notifications",
		"Must accept to continue",
		"send you updates",
	)
}

func TestToggleDisabledState(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		Label:    "Locked setting",
		Disabled: true,
	}))
	utils.AssertContains(t, output, "Locked setting")
}

func TestToggleInvalidSizeFallback(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		Label: "Fallback",
		Size:  ToggleSize("bogus"),
	}))
	utils.AssertContains(t, output, "Fallback")
}

// ---------------------------------------------------------------------------
// Select: normalizeSelectOptions edge cases, full props
// ---------------------------------------------------------------------------

func TestSelectDisabledAndSelectedContradiction(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Select(SelectProps{
		Name:  "country",
		Label: "Country",
		Options: []SelectOption{
			{Value: "de", Label: "Germany", Selected: true},
			{Value: "at", Label: "Austria", Disabled: true, Selected: true},
		},
	}))
	utils.AssertContainsAll(t, output, "Germany", "Austria")
}

func TestSelectFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Select(SelectProps{
		Name:     "role",
		Label:    "Role",
		Required: true,
		Error:    "Please select a role",
		HelpText: "Choose your access level",
		BaseProps: utils.BaseProps{
			ID:        "select-role",
			AriaLabel: "Role selector",
		},
		Options: []SelectOption{
			{Value: "admin", Label: "Admin"},
			{Value: "user", Label: "User", Disabled: true},
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="select-role"`,
		"Please select a role",
		"Choose your access level",
		"Admin",
		"User",
	)
}

func TestSelectMultipleSelectedOnlyFirstKept(t *testing.T) {
	t.Parallel()
	opts := normalizeSelectOptions([]SelectOption{
		{Value: "a", Label: "A", Selected: true},
		{Value: "b", Label: "B", Selected: true},
		{Value: "c", Label: "C"},
	})
	selected := 0
	for _, o := range opts {
		if o.Selected {
			selected++
		}
	}
	if selected != 1 {
		t.Errorf("expected exactly 1 selected option, got %d", selected)
	}
}

func TestSelectDisabledSelectedCleared(t *testing.T) {
	t.Parallel()
	opts := normalizeSelectOptions([]SelectOption{
		{Value: "x", Label: "X", Disabled: true, Selected: true},
	})
	if opts[0].Selected {
		t.Error("disabled+selected should clear Selected")
	}
}

func TestSelectEmpty(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Select(SelectProps{
		Name:  "empty",
		Label: "Empty Select",
	}))
	utils.AssertContains(t, output, "Empty Select")
}

// ---------------------------------------------------------------------------
// Textarea: full props
// ---------------------------------------------------------------------------

func TestTextareaFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Textarea(TextareaProps{
		Name:        "bio",
		Label:       "Biography",
		Value:       "Software engineer",
		Placeholder: "Tell us about yourself",
		Rows:        6,
		Required:    true,
		Disabled:    true,
		ReadOnly:    true,
		MaxLength:   500,
		Error:       "Too short",
		HelpText:    "Min 50 characters",
		BaseProps: utils.BaseProps{
			ID:        "textarea-bio",
			AriaLabel: "Biography input",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="textarea-bio"`,
		`name="bio"`,
		"Tell us about yourself",
		"Software engineer",
	)
}

func TestTextareaMinimal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Textarea(TextareaProps{
		Name:  "comment",
		Label: "Comment",
	}))
	utils.AssertContains(t, output, "Comment")
}

// ---------------------------------------------------------------------------
// Checkbox: full props
// ---------------------------------------------------------------------------

func TestCheckboxFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Name:     "agree",
		Label:    "I agree to the terms",
		Required: true,
		Error:    "You must agree",
		HelpText: "Read the terms carefully",
		BaseProps: utils.BaseProps{
			ID:        "cb-agree",
			AriaLabel: "Terms agreement",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="cb-agree"`,
		"I agree to the terms",
		"You must agree",
		"Read the terms carefully",
	)
}

func TestCheckboxDisabledState(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Label:    "Pre-checked",
		Disabled: true,
	}))
	utils.AssertContains(t, output, "Pre-checked")
}

// ---------------------------------------------------------------------------
// Radio + RadioGroup
// ---------------------------------------------------------------------------

func TestRadioGroupFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:     "plan",
		Label:    "Select plan",
		Inline:   true,
		Required: true,
		Error:    "Please choose a plan",
		HelpText: "You can upgrade later",
		BaseProps: utils.BaseProps{
			ID:        "rg-plan",
			AriaLabel: "Plan selector",
		},
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro", Disabled: true},
			{Value: "enterprise", Label: "Enterprise"},
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="rg-plan"`,
		"Select plan",
		"Please choose a plan",
		"Free",
		"Pro",
		"Enterprise",
	)
}

func TestRadioGroupVertical(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:  "color",
		Label: "Color",
		Options: []RadioOption{
			{Value: "red", Label: "Red"},
			{Value: "blue", Label: "Blue"},
		},
	}))
	utils.AssertContainsAll(t, output, "Red", "Blue")
}

func TestRadioSingleFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		Name:  "option",
		Value: "val1",
		Label: "Option 1",
		BaseProps: utils.BaseProps{
			ID:        "radio-1",
			AriaLabel: "First option",
		},
	}))
	utils.AssertContainsAll(t, output, `id="radio-1"`, "Option 1")
}

// ---------------------------------------------------------------------------
// Combobox: full props
// ---------------------------------------------------------------------------

func TestComboboxFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Combobox(ComboboxProps{
		Name:        "country",
		Label:       "Country",
		Placeholder: "Select a country",
		Value:       "de",
		Required:    true,
		Disabled:    true,
		Error:       "Country required",
		HelpText:    "Pick your home country",
		BaseProps: utils.BaseProps{
			ID:        "cb-country",
			AriaLabel: "Country combobox",
			Nonce:     "nonce123",
		},
		Options: []ComboboxOption{
			{Value: "de", Label: "Germany"},
			{Value: "at", Label: "Austria"},
		},
	}))
	utils.AssertContainsAll(t, output,
		"Country combobox",
		"Germany",
		"Austria",
		"Country required",
		"Pick your home country",
	)
}

func TestComboboxDisplayLabelFallback(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Combobox(ComboboxProps{
		Name:  "item",
		Value: "nonexistent",
		Options: []ComboboxOption{
			{Value: "a", Label: "Alpha"},
		},
	}))
	// When Value doesn't match any option, display the raw value
	utils.AssertContains(t, output, "nonexistent")
}

// ---------------------------------------------------------------------------
// Form: inline, GET method, CSRF
// ---------------------------------------------------------------------------

func TestFormInline(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action: "/search",
		Method: FormGet,
		Inline: true,
	}))
	utils.AssertContainsAll(t, output,
		`action="/search"`,
		`method="GET"`,
	)
}

func TestFormWithCSRF(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action:    "/submit",
		Method:    FormPost,
		CSRFToken: "csrf-abc-123",
	}))
	utils.AssertContainsAll(t, output,
		`action="/submit"`,
		`method="POST"`,
		`value="csrf-abc-123"`,
	)
}

func TestFormWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action: "/save",
		BaseProps: utils.BaseProps{
			ID:        "save-form",
			AriaLabel: "Save form",
			Class:     "mb-4",
		},
	}))
	utils.AssertContainsAll(t, output, `id="save-form"`, `action="/save"`)
}

// ---------------------------------------------------------------------------
// DatePicker, FileInput, InputGroup full props
// ---------------------------------------------------------------------------

func TestDatePickerFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DatePicker(DatePickerProps{
		Name:     "birthdate",
		Label:    "Birth Date",
		Required: true,
		Error:    "Date required",
		HelpText: "YYYY-MM-DD",
		BaseProps: utils.BaseProps{
			ID:        "dp-1",
			AriaLabel: "Birth date picker",
		},
	}))
	utils.AssertContainsAll(t, output, `id="dp-1"`, "Birth Date", "Date required")
}

func TestFileInputFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		Name:     "avatar",
		Label:    "Upload Avatar",
		Required: true,
		HelpText: "PNG or JPG, max 2MB",
		BaseProps: utils.BaseProps{
			ID:        "fi-1",
			AriaLabel: "Avatar upload",
		},
	}))
	utils.AssertContainsAll(t, output, `id="fi-1"`, "Upload Avatar")
}

func TestInputGroupFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InputGroup(InputGroupProps{
		LeftAddon:  templ.Raw("$"),
		RightAddon: templ.Raw("USD"),
		BaseProps: utils.BaseProps{
			ID:        "ig-1",
			AriaLabel: "Amount input",
		},
	}))
	utils.AssertContainsAll(t, output, `id="ig-1"`, "$", "USD")
}

// ---------------------------------------------------------------------------
// ValidationSummary
// ---------------------------------------------------------------------------

func TestValidationSummaryFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		Errors: []ValidationError{
			{Field: "email", Message: "Invalid email"},
			{Field: "name", Message: "Name required"},
		},
	}))
	utils.AssertContainsAll(t, output,
		"Invalid email",
		"Name required",
	)
}

func TestValidationSummaryDefaultTitle(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		Errors: []ValidationError{
			{Field: "x", Message: "Error 1"},
		},
	}))
	utils.AssertContains(t, output, "Error 1")
}
