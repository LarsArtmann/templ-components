package forms

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestTextareaAutoGrowDefaultTrue(t *testing.T) {
	t.Parallel()

	props := DefaultTextareaProps()
	if !props.AutoGrow {
		t.Error("DefaultTextareaProps should have AutoGrow=true")
	}
}

func TestTextareaAutoGrowAddsClass(t *testing.T) {
	t.Parallel()

	props := DefaultTextareaProps()
	props.ID = "bio"
	props.Name = "bio"

	html := utils.Render(t, Textarea(props))
	if !strings.Contains(html, "tc-auto-grow") {
		t.Error("AutoGrow=true should add tc-auto-grow class")
	}
}

func TestTextareaAutoGrowFalseOmitsClass(t *testing.T) {
	t.Parallel()

	props := DefaultTextareaProps()
	props.AutoGrow = false
	props.ID = "bio"
	props.Name = "bio"

	html := utils.Render(t, Textarea(props))
	if strings.Contains(html, "tc-auto-grow") {
		t.Error("AutoGrow=false should not add tc-auto-grow class")
	}
}

func TestTextareaEnterKeyHint(t *testing.T) {
	t.Parallel()

	props := DefaultTextareaProps()
	props.ID = "chat"
	props.Name = "message"
	props.EnterKeyHint = EnterKeyHintSend

	html := utils.Render(t, Textarea(props))
	if !strings.Contains(html, `enterkeyhint="send"`) {
		t.Error("EnterKeyHint=Send should emit enterkeyhint=\"send\"")
	}
}

func TestTextareaEnterKeyHintEmpty(t *testing.T) {
	t.Parallel()

	props := DefaultTextareaProps()
	props.ID = "bio"
	props.Name = "bio"

	html := utils.Render(t, Textarea(props))
	if strings.Contains(html, "enterkeyhint") {
		t.Error("Empty EnterKeyHint should not emit enterkeyhint attribute")
	}
}

func TestInputSearchWrapsInSearchElement(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Input(InputProps{
		BaseProps: utils.BaseProps{ID: "q"},
		Type:      InputSearch,
		Name:      "q",
	}))
	if !strings.Contains(html, "<search") {
		t.Error("InputSearch should wrap input in <search> element")
	}

	if !strings.Contains(html, `type="search"`) {
		t.Error("InputSearch should emit type=\"search\"")
	}
}

func TestInputTextDoesNotWrapInSearchElement(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Input(InputProps{
		BaseProps: utils.BaseProps{ID: "name"},
		Type:      InputText,
		Name:      "name",
	}))
	if strings.Contains(html, "<search") {
		t.Error("InputText should not wrap in <search> element")
	}
}

func TestFormValidate(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Form(FormProps{
		Action:   "/submit",
		Validate: true,
	}))
	if !strings.Contains(html, `hx-validate="true"`) {
		t.Error("Validate=true should emit hx-validate=\"true\"")
	}
}

func TestFormValidateFalseOmits(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Form(FormProps{Action: "/submit"}))
	if strings.Contains(html, "hx-validate") {
		t.Error("Validate=false should not emit hx-validate")
	}
}

func TestInputEnterKeyHintExplicitOverridesAuto(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Input(InputProps{
		BaseProps:    utils.BaseProps{ID: "msg"},
		Type:         InputText,
		Name:         "msg",
		EnterKeyHint: EnterKeyHintSend,
	}))
	if !strings.Contains(html, `enterkeyhint="send"`) {
		t.Error("Explicit EnterKeyHint=Send should emit enterkeyhint=\"send\" even for InputText")
	}
}

func TestInputEnterKeyHintAutoDerived(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Input(InputProps{
		BaseProps: utils.BaseProps{ID: "email"},
		Type:      InputEmail,
		Name:      "email",
	}))
	if !strings.Contains(html, `enterkeyhint="next"`) {
		t.Error("InputEmail should auto-derive enterkeyhint=\"next\"")
	}
}

func TestSelectStylableEmitsSelectedContent(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Select(SelectProps{
		BaseProps: utils.BaseProps{ID: "country"},
		Name:      "country",
		Stylable:  true,
		Options: []SelectOption{
			{Value: "de", Label: "Germany"},
		},
	}))
	if !strings.Contains(html, "tc-select") {
		t.Error("Stylable=true should add tc-select class")
	}

	if !strings.Contains(html, "<selectedcontent>") {
		t.Error("Stylable=true should emit <selectedcontent> element")
	}

	if !strings.Contains(html, "<button>") {
		t.Error("Stylable=true should emit <button> wrapper")
	}
}

func TestSelectNotStylableOmitsSelectedContent(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, Select(SelectProps{
		BaseProps: utils.BaseProps{ID: "country"},
		Name:      "country",
		Options: []SelectOption{
			{Value: "de", Label: "Germany"},
		},
	}))
	if strings.Contains(html, "tc-select") {
		t.Error("Stylable=false should not add tc-select class")
	}

	if strings.Contains(html, "<selectedcontent>") {
		t.Error("Stylable=false should not emit <selectedcontent>")
	}
}
