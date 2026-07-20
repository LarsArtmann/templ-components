package layout

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestContainer(t *testing.T) {
	t.Parallel()

	t.Run("default props produce max-w-7xl + padding + centering", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(DefaultContainerProps()))
		// tailwind-merge reorders classes; assert each token individually.
		utils.AssertContainsAll(t, output, "mx-auto", "w-full", "max-w-7xl", "px-4", "sm:px-6", "lg:px-8")
	})

	t.Run("prose width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidthProse}))
		utils.AssertContains(t, output, "max-w-prose")
	})

	t.Run("sm width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidthSM}))
		utils.AssertContains(t, output, "max-w-3xl")
	})

	t.Run("md width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidthMD}))
		utils.AssertContains(t, output, "max-w-5xl")
	})

	t.Run("xl width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidthXL}))
		utils.AssertContains(t, output, "max-w-[90rem]")
	})

	t.Run("full width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidthFull}))
		utils.AssertContains(t, output, "max-w-full")
	})

	t.Run("unknown width falls back to default (lg)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidth("bogus")}))
		utils.AssertContains(t, output, "max-w-7xl")
	})

	t.Run("pad=false omits pad class", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Container(ContainerProps{Width: ContainerWidthDefault, Pad: false}))
		if strings.Contains(output, "px-4") || strings.Contains(output, "sm:px-6") ||
			strings.Contains(output, "lg:px-8") {
			t.Errorf("expected pad class omitted; output = %q", output)
		}

		utils.AssertContains(t, output, "max-w-7xl")
	})

	t.Run("BaseProps propagate (ID, Class, AriaLabel, Attrs)", func(t *testing.T) {
		t.Parallel()

		props := ContainerProps{
			BaseProps: utils.BaseProps{
				ID:        "main-wrap",
				Class:     "data-tc-test",
				AriaLabel: "Main content wrapper",
				Attrs:     templ.Attributes{"data-testid": "container"},
			},
		}
		output := utils.Render(t, Container(props))
		utils.AssertContainsAll(t, output,
			`id="main-wrap"`,
			"data-tc-test",
			`aria-label="Main content wrapper"`,
			`data-testid="container"`,
		)
	})

	t.Run("consumer Class overrides default (tailwind-merge)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Container(ContainerProps{
			BaseProps: utils.BaseProps{Class: "max-w-2xl"},
		}))
		utils.AssertContains(t, output, "max-w-2xl")

		if strings.Contains(output, "max-w-7xl") {
			t.Errorf("consumer max-w-2xl should override default max-w-7xl via tailwind-merge")
		}
	})
}

func TestContainerWidthIsValid(t *testing.T) {
	t.Parallel()

	valid := []ContainerWidth{
		ContainerWidthSM, ContainerWidthMD, ContainerWidthLG,
		ContainerWidthXL, ContainerWidthFull, ContainerWidthProse,
	}
	for _, w := range valid {
		if !ContainerWidthIsValid(w) {
			t.Errorf("ContainerWidthIsValid(%q) = false; want true", w)
		}
	}

	if ContainerWidthIsValid(ContainerWidth("nope")) {
		t.Errorf("ContainerWidthIsValid(\"nope\") = true; want false")
	}
}
