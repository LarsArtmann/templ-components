package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestButtonRender(t *testing.T) {
	t.Parallel()

	t.Run("primary button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(DefaultButtonProps()))
		utils.AssertContains(t, output, `<button`)
		utils.AssertContains(t, output, `type="button"`)
		utils.AssertContains(t, output, `bg-blue-600`)
	})

	t.Run("button with text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			BaseProps: utils.BaseProps{ID: "save-btn"},
			Text:      "Save",
			Type:      ButtonHTMLSubmit,
		}))
		utils.AssertContains(t, output, `id="save-btn"`)
		utils.AssertContains(t, output, "Save")
		utils.AssertContains(t, output, `type="submit"`)
	})

	t.Run("link button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			Text:    "Docs",
			Href:    "/docs",
			Variant: ButtonSecondary,
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, `href="/docs"`)
		utils.AssertContains(t, output, "Docs")
	})

	t.Run("disabled button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			Text:     "Delete",
			Variant:  ButtonDanger,
			Disabled: true,
		}))
		utils.AssertContains(t, output, `disabled`)
		utils.AssertContains(t, output, `bg-red-600`)
	})

	t.Run("external link button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{
			Text:     "GitHub",
			Href:     "https://github.com",
			External: true,
		}))
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})
}

func TestButtonHTMLType(t *testing.T) {
	t.Parallel()

	t.Run("submit type passes through", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Save", Type: ButtonHTMLSubmit}))
		utils.AssertContains(t, output, `type="submit"`)
	})

	t.Run("reset type passes through", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Clear", Type: ButtonHTMLReset}))
		utils.AssertContains(t, output, `type="reset"`)
	})

	t.Run("empty falls back to button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Go"}))
		utils.AssertContains(t, output, `type="button"`)
	})

	t.Run("invalid type falls back to button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Go", Type: ButtonHTMLType("destroy-everything")}))
		utils.AssertContains(t, output, `type="button"`)
		utils.AssertNotContains(t, output, `destroy-everything`)
	})
}

func TestButtonSizes(t *testing.T) {
	t.Parallel()
	t.Run("small button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "SM", Size: ButtonSizeSM}))
		utils.AssertContains(t, output, `px-2.5`)
	})
	t.Run("large button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "LG", Size: ButtonSizeLG}))
		utils.AssertContains(t, output, `px-4`)
	})
}

func TestButtonVariants(t *testing.T) {
	t.Parallel()
	t.Run("ghost button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Ghost", Variant: ButtonGhost}))
		utils.AssertContains(t, output, `bg-transparent`)
	})
	t.Run("link button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Button(ButtonProps{Text: "Link", Variant: ButtonLink}))
		utils.AssertContains(t, output, `text-blue-600`)
	})
}
