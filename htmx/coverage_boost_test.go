package htmx

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestConfirmDeleteCoverage(t *testing.T) {
	t.Parallel()
	t.Run("full props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ConfirmDelete(ConfirmDeleteProps{
			BaseProps: utils.BaseProps{
				ID:        "del-btn",
				Class:     "text-sm",
				AriaLabel: "Delete item",
				Attrs:     templ.Attributes{"data-test": "yes"},
			},
			Delete:  "/api/items/42",
			Target:  "#row-42",
			Confirm: "Really delete?",
		}))
		utils.AssertContainsAll(t, output,
			`id="del-btn"`,
			`hx-delete="/api/items/42"`,
			`hx-target="#row-42"`,
			`hx-swap="outerHTML"`,
			`hx-confirm="Really delete?"`,
			`aria-label="Delete item"`,
			"data-test",
			"Delete",
		)
	})
	t.Run("minimal", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ConfirmDelete(ConfirmDeleteProps{
			Delete:  "/delete",
			Target:  "#item",
			Confirm: "Sure?",
		}))
		utils.AssertContains(t, output, `hx-delete="/delete"`)
		utils.AssertNotContains(t, output, `id=`)
	})
}

func TestSwapOOBCoverage(t *testing.T) {
	t.Parallel()
	for _, style := range []SwapStyle{
		SwapInnerHTML, SwapOuterHTML, SwapBeforeBegin, SwapAfterBegin,
		SwapBeforeEnd, SwapAfterEnd, SwapDelete, SwapNone,
	} {
		t.Run("style_"+string(style), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, SwapOOB(SwapOOBProps{
				Selector:  "#target",
				SwapStyle: style,
			}))
			utils.AssertContains(t, output, string(style)+":#target")
		})
	}
	t.Run("invalid style falls back", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SwapOOB(SwapOOBProps{
			Selector:  "#target",
			SwapStyle: SwapStyle("bogus"),
		}))
		utils.AssertContains(t, output, "outerHTML:#target")
	})
	t.Run("with id and attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SwapOOB(SwapOOBProps{
			BaseProps: utils.BaseProps{
				ID:    "oob-1",
				Class: "custom",
				Attrs: templ.Attributes{"data-x": "1"},
			},
			Selector:  "#toast",
			SwapStyle: SwapBeforeEnd,
		}))
		utils.AssertContains(t, output, `id="oob-1"`)
		utils.AssertContains(t, output, "custom")
		utils.AssertContains(t, output, `data-x="1"`)
	})
}

func TestCSRFTokenCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with token", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CSRFToken("abc123secret"))
		utils.AssertContains(t, output, `type="hidden"`)
		utils.AssertContains(t, output, `name="csrf_token"`)
		utils.AssertContains(t, output, `value="abc123secret"`)
	})
}
