package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestInputEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props InputProps
		want  []string
	}{
		{"empty name", InputProps{Type: InputText, Label: "Test"}, []string{`type="text"`}},
		{"all input types", InputProps{Name: "t", Type: InputEmail, Label: "E"}, []string{`type="email"`}},
		{"password type", InputProps{Name: "p", Type: InputPassword, Label: "P"}, []string{`type="password"`}},
		{"search type", InputProps{Name: "s", Type: InputSearch, Label: "S"}, []string{`type="search"`}},
		{"url type", InputProps{Name: "u", Type: InputURL, Label: "U"}, []string{`type="url"`}},
		{"tel type", InputProps{Name: "t", Type: InputTel, Label: "T"}, []string{`type="tel"`}},
		{"number type", InputProps{Name: "n", Type: InputNumber, Label: "N"}, []string{`type="number"`}},
		{"date type", InputProps{Name: "d", Type: InputDate, Label: "D"}, []string{`type="date"`}},
		{"datetime-local type", InputProps{Name: "dt", Type: InputDatetime, Label: "DT"}, []string{`type="datetime-local"`}},
		{"time type", InputProps{Name: "t", Type: InputTime, Label: "T"}, []string{`type="time"`}},
		{"hidden type", InputProps{Name: "h", Type: InputHidden, Label: "H"}, []string{`type="hidden"`}},
		{"custom id/class", InputProps{BaseProps: utils.BaseProps{ID: "inp", Class: "mt-2"}, Name: "x", Type: InputText, Label: "X"}, []string{`id="inp"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Input(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestTextareaMoreEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props TextareaProps
		want  []string
	}{
		{"zero rows", TextareaProps{Name: "t", Label: "T", Rows: 0}, []string{`<textarea`}},
		{"max rows", TextareaProps{Name: "t", Label: "T", Rows: 20}, []string{`rows="20"`}},
		{"custom id/class", TextareaProps{BaseProps: utils.BaseProps{ID: "ta", Class: "mt-2"}, Name: "x", Label: "X"}, []string{`id="ta"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Textarea(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestCheckboxMoreEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props CheckboxProps
		want  []string
	}{
		{"unchecked", CheckboxProps{Name: "c", Label: "C"}, []string{`type="checkbox"`}},
		{"custom id/class", CheckboxProps{BaseProps: utils.BaseProps{ID: "cb", Class: "mt-2"}, Name: "x", Label: "X"}, []string{`id="cb"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Checkbox(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestSelectMoreEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props SelectProps
		want  []string
	}{
		{"disabled select", SelectProps{Name: "s", Label: "S", Disabled: true}, []string{`disabled`}},
		{"select with id", SelectProps{BaseProps: utils.BaseProps{ID: "sel"}, Name: "s"}, []string{`id="sel"`}},
		{"select with aria-label", SelectProps{BaseProps: utils.BaseProps{AriaLabel: "Choose"}, Name: "s"}, []string{`aria-label="Choose"`}},
		{"select with custom class", SelectProps{BaseProps: utils.BaseProps{Class: "custom"}, Name: "s"}, []string{"custom"}},
		{"select with attrs", SelectProps{BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-test": "yes"}}, Name: "s"}, []string{`data-test="yes"`}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Select(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

// TestSelectDoesNotMutateCallerOptions verifies that normalizeSelectOptions
// returns a defensive copy and does not corrupt the caller's []SelectOption.
// This was a real bug: reusing the same options slice across renders silently
// cleared Selected flags on the second render.
func TestSelectDoesNotMutateCallerOptions(t *testing.T) {
	t.Parallel()

	opts := []SelectOption{
		{Value: "a", Label: "A", Selected: true, Disabled: true}, // Disabled+Selected → Selected cleared
		{Value: "b", Label: "B", Selected: true},                 // valid second-selected → cleared (single-value)
	}
	props := SelectProps{Name: "s", Options: opts}

	_ = utils.Render(t, Select(props))
	_ = utils.Render(t, Select(props)) // second render must see the same original data

	if !opts[0].Selected || !opts[0].Disabled {
		t.Error("normalizeSelectOptions mutated caller's slice: opts[0] changed")
	}
	if !opts[1].Selected {
		t.Error("normalizeSelectOptions mutated caller's slice: opts[1].Selected was cleared")
	}
}

// TestCheckboxWithoutIDDoesNotEmitEmptyFor verifies that a checkbox without an
// ID does not render <label for=""> (invalid HTML that breaks label association).
func TestCheckboxWithoutIDDoesNotEmitEmptyFor(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Name:  "agree",
		Label: "I agree",
	}))
	utils.AssertNotContains(t, output, `for=""`)
	utils.AssertContains(t, output, "I agree")
}

// TestToggleEmitsCompletePeerCheckedClasses verifies that the toggle's thumb
// translate classes include the complete "peer-checked:" variant prefix so
// Tailwind's content scanner can detect them. Dynamically concatenated variant
// prefixes are invisible to the scanner and produce no CSS — the thumb would
// not slide when checked.
func TestToggleEmitsCompletePeerCheckedClasses(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		size ToggleSize
		want string
	}{
		{"sm", ToggleSizeSM, "peer-checked:translate-x-4"},
		{"md default", ToggleSizeMD, "peer-checked:translate-x-5"},
		{"lg", ToggleSizeLG, "peer-checked:translate-x-6"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Toggle(ToggleProps{
				Name:  "t",
				Label: "Toggle",
				Size:  tt.size,
			}))
			utils.AssertContains(t, output, tt.want)
		})
	}
}

func TestTextareaFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("readonly textarea", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     "t",
			Label:    "T",
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, `readonly`)
	})

	t.Run("textarea with maxlength", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:      "t",
			Label:     "T",
			MaxLength: 500,
		}))
		utils.AssertContains(t, output, `maxlength="500"`)
	})

	t.Run("textarea with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name: "t",
			BaseProps: utils.BaseProps{
				AriaLabel: "Description",
			},
		}))
		utils.AssertContains(t, output, `aria-label="Description"`)
	})

	t.Run("textarea with attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name: "t",
			BaseProps: utils.BaseProps{
				Attrs: templ.Attributes{"data-test": "yes"},
			},
		}))
		utils.AssertContains(t, output, `data-test="yes"`)
	})

	t.Run("textarea with disabled and value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Textarea(TextareaProps{
			Name:     "t",
			Label:    "T",
			Value:    "pre-filled",
			Disabled: true,
		}))
		utils.AssertContains(t, output, "pre-filled")
		utils.AssertContains(t, output, `disabled`)
	})
}

func TestFormRender(t *testing.T) {
	t.Parallel()

	t.Run("form with action and method", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{
			Action: "/submit",
			Method: FormPost,
		}))
		utils.AssertContains(t, output, `action="/submit"`)
		utils.AssertContains(t, output, `method="POST"`)
	})

	t.Run("form with CSRF token", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{
			CSRFToken: "secret-token",
		}))
		utils.AssertContains(t, output, `name="csrf_token"`)
		utils.AssertContains(t, output, `value="secret-token"`)
	})

	t.Run("form without CSRF omits hidden input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{}))
		utils.AssertNotContains(t, output, `csrf_token`)
	})

	t.Run("form with custom ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{
			BaseProps: utils.BaseProps{ID: "my-form", Class: "max-w-md"},
		}))
		utils.AssertContains(t, output, `id="my-form"`)
		utils.AssertContains(t, output, "max-w-md")
	})

	t.Run("default props uses POST", func(t *testing.T) {
		t.Parallel()
		props := DefaultFormProps()
		if props.Method != FormPost {
			t.Errorf("DefaultFormProps().Method = %q, want %q", props.Method, FormPost)
		}
	})
}
