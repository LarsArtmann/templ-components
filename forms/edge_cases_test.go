package forms

import (
	"testing"

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
			for _, w := range tt.want {
				utils.AssertContains(t, output, w)
			}
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
			for _, w := range tt.want {
				utils.AssertContains(t, output, w)
			}
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
			for _, w := range tt.want {
				utils.AssertContains(t, output, w)
			}
		})
	}
}
