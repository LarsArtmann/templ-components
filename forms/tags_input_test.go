package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultTagsInputProps(t *testing.T) {
	t.Parallel()

	_ = DefaultTagsInputProps()
}

func TestTagsInputBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, TagsInput(TagsInputProps{
		BaseProps:   utils.BaseProps{ID: "skills"},
		Name:        "skills",
		Label:       "Skills",
		Values:      []string{"Go", "HTMX"},
		Placeholder: "Add a skill...",
	}))
	utils.AssertContains(t, output, "Go")
	utils.AssertContains(t, output, "HTMX")
	utils.AssertContains(t, output, `name="skills"`)
	utils.AssertContains(t, output, `type="hidden"`)
	utils.AssertContains(t, output, `data-tc-tags`)
}

func TestTagsInputGolden(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, TagsInput(TagsInputProps{
		BaseProps:   utils.BaseProps{ID: "tags"},
		Name:        "tags",
		Label:       "Tags",
		Values:      []string{"Go", "Templ"},
		Placeholder: "Add tag...",
	}))
	golden.Assert(t, "tags_input_basic", output)
}

func TestTagsInputA11y(t *testing.T) {
	t.Parallel()

	t.Run("remove buttons have aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "t"},
			Name:      "tags",
			Values:    []string{"Go", "HTMX"},
		}))
		utils.AssertContains(t, output, `aria-label="Remove Go"`)
		utils.AssertContains(t, output, `aria-label="Remove HTMX"`)
	})

	t.Run("error sets aria-invalid", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "err"},
			Name:      "tags",
			Error:     "At least one tag required",
		}))
		utils.AssertContains(t, output, `aria-invalid="true"`)
	})

	t.Run("disabled adds disabled attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "dis"},
			Name:      "tags",
			Disabled:  true,
		}))
		utils.AssertContains(t, output, "disabled")
	})

	t.Run("required shows asterisk", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "req"},
			Name:      "tags",
			Label:     "Tags",
			Required:  true,
		}))
		utils.AssertContains(t, output, "*")
	})

	t.Run("dark mode classes present", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "dm"},
			Name:      "tags",
			Values:    []string{"X"},
		}))
		utils.AssertContains(t, output, "dark:bg-gray-800")
		utils.AssertContains(t, output, "dark:bg-blue-900/50")
	})
}

func TestTagsInputEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty values renders just input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps:   utils.BaseProps{ID: "empty"},
			Name:        "tags",
			Placeholder: "Type...",
		}))
		utils.AssertNotContains(t, output, "data-tc-tag=")
		utils.AssertContains(t, output, "Type...")
	})

	t.Run("max tags renders data attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "max"},
			Name:      "tags",
			MaxTags:   5,
		}))
		utils.AssertContains(t, output, `data-max-tags="5"`)
	})

	t.Run("allow duplicate renders data attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps:      utils.BaseProps{ID: "dup"},
			Name:           "tags",
			AllowDuplicate: true,
		}))
		utils.AssertContains(t, output, `data-allow-duplicate`)
	})

	t.Run("no label omits label element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "nolabel"},
			Name:      "tags",
		}))
		utils.AssertNotContains(t, output, "<label")
	})

	t.Run("help text renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "help"},
			Name:      "tags",
			HelpText:  "Press Enter to add",
		}))
		utils.AssertContains(t, output, "Press Enter to add")
	})

	t.Run("nonce renders script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "nonce", Nonce: "abc123"},
			Name:      "tags",
		}))
		utils.AssertContains(t, output, `nonce="abc123"`)
		utils.AssertContains(t, output, "tcTagsInputAttached")
	})

	t.Run("no nonce omits script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "noscript"},
			Name:      "tags",
		}))
		utils.AssertNotContains(t, output, "tcTagsInputAttached")
	})

	t.Run("custom class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps: utils.BaseProps{ID: "cls", Class: "my-tags"},
			Name:      "tags",
		}))
		utils.AssertContains(t, output, "my-tags")
	})
}

func TestTagsInputSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("full-featured tags input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, TagsInput(TagsInputProps{
			BaseProps:   utils.BaseProps{ID: "full", Nonce: "n"},
			Name:        "keywords",
			Label:       "Keywords",
			Values:      []string{"Go", "HTMX", "Templ"},
			Placeholder: "Add keyword...",
			MaxTags:     10,
			Required:    true,
			HelpText:    "Add up to 10 keywords",
		}))
		utils.AssertContainsAll(t, output,
			"Go", "HTMX", "Templ",
			`name="keywords"`,
			`data-max-tags="10"`,
			"Add up to 10 keywords",
			"required",
			"tcTagsInputAttached",
		)
	})
}

func TestTagsInputContainerClass(t *testing.T) {
	t.Parallel()

	if got := tagsInputContainerClass(false); got == "" {
		t.Error("expected non-empty class string")
	}

	if got := tagsInputContainerClass(true); got == "" {
		t.Error("expected non-empty error class string")
	}
}

func ExampleTagsInput() {
	_ = TagsInput(TagsInputProps{
		Name:        "skills",
		Label:       "Skills",
		Values:      []string{"Go", "HTMX"},
		Placeholder: "Add a skill...",
	})
	// Output:
}
