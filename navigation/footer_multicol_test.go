package navigation

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestFooterMultiColumn(t *testing.T) {
	t.Parallel()

	t.Run("legacy: empty Columns renders single-row centered copyright", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{BrandText: "Acme"}))
		utils.AssertContains(t, output, "All rights reserved.")
		utils.AssertContains(t, output, "Acme")

		if strings.Contains(output, "grid-cols") {
			t.Errorf("legacy footer should not render grid; output has grid-cols")
		}
	})

	t.Run("multi-column: Columns non-empty triggers grid layout", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BrandText: "Acme",
			Columns: []FooterColumn{
				{Title: "Product", Links: []FooterLink{{"Features", "/features"}}},
				{Title: "Company", Links: []FooterLink{{"About", "/about"}}},
			},
		}))
		utils.AssertContainsAll(t, output,
			"grid", "grid-cols-2", "md:grid-cols-4", "gap-8",
			"Product", "Features", "/features",
			"Company", "About", "/about",
			"Acme",
		)
	})

	t.Run("multi-column: brand block only when BrandText set", func(t *testing.T) {
		t.Parallel()
		outputNoBrand := utils.Render(t, Footer(FooterProps{
			Columns: []FooterColumn{{Title: "X", Links: []FooterLink{{"Y", "/y"}}}},
		}))
		// Brand block emits text-lg (only on brand); column titles use text-sm.
		if strings.Contains(outputNoBrand, "text-lg") {
			t.Errorf("empty BrandText should not render brand block; got %q", outputNoBrand)
		}
		// Sanity: WITH BrandText, brand block IS rendered.
		outputWithBrand := utils.Render(t, Footer(FooterProps{
			BrandText: "Acme",
			Columns:   []FooterColumn{{Title: "X", Links: []FooterLink{{"Y", "/y"}}}},
		}))
		utils.AssertContains(t, outputWithBrand, "text-lg")
	})

	t.Run("multi-column: BottomBar renders below columns", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BrandText: "Acme",
			Columns:   []FooterColumn{{Title: "X", Links: []FooterLink{{"Y", "/y"}}}},
			BottomBar: templ.Raw(`<p data-test="legal">© 2026 Acme</p>`),
		}))
		utils.AssertContainsAll(t, output, `data-test="legal"`, "border-t")
	})

	t.Run("legacy: BottomBar renders when Columns empty (backward-compat extension)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BrandText: "Acme",
			BottomBar: templ.Raw(`<span data-test="bb">v1.0</span>`),
		}))
		utils.AssertContainsAll(t, output, `data-test="bb"`, "All rights reserved.")
	})

	t.Run("BaseProps propagate (ID, Class, AriaLabel, Attrs)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BaseProps: utils.BaseProps{
				ID:        "foot",
				Class:     "data-tc-test",
				AriaLabel: "Site footer",
				Attrs:     templ.Attributes{"data-testid": "footer"},
			},
		}))
		utils.AssertContainsAll(t, output,
			`id="foot"`, "data-tc-test", `aria-label="Site footer"`, `data-testid="footer"`,
		)
	})

	t.Run("default aria-label is 'Footer' when AriaLabel empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{BrandText: "X"}))
		utils.AssertContains(t, output, `aria-label="Footer"`)
	})

	t.Run("dark mode: every neutral has dark: variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BrandText: "Acme",
			Columns:   []FooterColumn{{Title: "X", Links: []FooterLink{{"Y", "/y"}}}},
		}))
		utils.AssertContainsAll(t, output,
			"dark:border-gray-800", "dark:bg-gray-900",
			"dark:text-white", "dark:text-gray-400",
		)
	})
}
