package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDarkModeClasses(t *testing.T) {
	t.Parallel()

	t.Run("input has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultInputProps()
		props.Name = "test"
		output := utils.Render(t, Input(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("select has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultSelectProps()
		props.Name = "test"
		output := utils.Render(t, Select(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("textarea has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultTextareaProps()
		props.Name = "test"
		output := utils.Render(t, Textarea(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("checkbox has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultCheckboxProps()
		props.Name = "test"
		output := utils.Render(t, Checkbox(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("toggle has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultToggleProps()
		props.Name = "test"
		output := utils.Render(t, Toggle(props))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("fileinput has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultFileInputProps()
		props.Name = "test"
		output := utils.Render(t, FileInput(props))
		utils.AssertContains(t, output, "dark:")
	})
}
