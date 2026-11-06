package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestGridAutoFitWithMinColWidth(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:        GridColsAutoFit,
		MinColWidth: "190px",
	}))
	utils.AssertContainsAll(t, output, "auto-fit", "minmax(var(--tc-grid-min),1fr)")
	utils.AssertContains(t, output, "--tc-grid-min:190px")
}

func TestGridAutoFitDefaultMinWidth(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols: GridColsAutoFit,
	}))
	utils.AssertContainsAll(t, output, "auto-fit", "minmax(var(--tc-grid-min),1fr)")
	utils.AssertContains(t, output, "--tc-grid-min:240px")
}

func TestGridAutoFitWithGap(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:        GridColsAutoFit,
		MinColWidth: "12rem",
		Gap:         GridGapLG,
	}))
	utils.AssertContainsAll(t, output, "auto-fit", "minmax(var(--tc-grid-min),1fr)", "gap-6")
	utils.AssertContains(t, output, "--tc-grid-min:12rem")
}

func TestGridAutoFitDoesNotEmitFixedBreakpoints(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:        GridColsAutoFit,
		MinColWidth: "190px",
	}))
	utils.AssertNotContains(t, output, "grid-cols-3")
	utils.AssertNotContains(t, output, "sm:grid-cols-2")
}

func TestGridColsAutoFitIsValid(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(t, "GridColsAutoFit valid", GridColsIsValid(GridColsAutoFit), true)
}

func TestGridAutoFitTakesPrecedenceOverContainerAware(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Grid(GridProps{
		Cols:           GridColsAutoFit,
		MinColWidth:    "190px",
		ContainerAware: true,
	}))
	utils.AssertContainsAll(t, output, "auto-fit", "minmax(var(--tc-grid-min),1fr)")
	utils.AssertContains(t, output, "--tc-grid-min:190px")
	utils.AssertNotContains(t, output, "@sm:grid-cols")
	utils.AssertNotContains(t, output, "@lg:grid-cols")
}

func TestCardHeaderSlotReplacesDefaultHeader(t *testing.T) {
	t.Parallel()

	customHeader := templ.Raw("<div data-test='custom-header'><h2>My Header</h2></div>")
	output := utils.Render(t, Card(CardProps{
		Header: customHeader,
	}))
	utils.AssertContains(t, output, "data-test='custom-header'")
	utils.AssertContains(t, output, "My Header")
}

func TestCardHeaderSlotSkipsDefaultTitle(t *testing.T) {
	t.Parallel()

	customHeader := templ.Raw("<div data-test='custom-header'>Custom</div>")
	output := utils.Render(t, Card(CardProps{
		Title:  "This Should Not Appear",
		Header: customHeader,
	}))
	utils.AssertContains(t, output, "data-test='custom-header'")
	utils.AssertNotContains(t, output, "This Should Not Appear")
}

func TestCardHeaderNilRendersDefaultTitle(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Card(CardProps{
		Title: "Default Title",
	}))
	utils.AssertContains(t, output, "Default Title")
	utils.AssertContains(t, output, "<h3")
}

func TestCardPaddingNoneSkipsWrappingDiv(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Card(CardProps{
		Title:   "Test",
		Padding: CardPaddingNone,
		Body:    templ.Raw("<table data-test='inner-table'><tr><td>data</td></tr></table>"),
	}))
	utils.AssertContains(t, output, "data-test='inner-table'")
	utils.AssertNotContains(t, output, "px-4 py-5")
}

func TestCardPaddingMDStillWrapsInDiv(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Card(CardProps{
		Padding: CardPaddingMD,
		Body:    templ.Raw("<p>content</p>"),
	}))
	utils.AssertContains(t, output, "px-4 py-5")
}

// The min column width must travel as a CSS custom property, not as part of a
// runtime-composed arbitrary class: Tailwind's scanner only sees literal class
// strings, so `[grid-template-columns:...minmax(190px,1fr)]` built at render
// time never produced any CSS and every auto-fit grid collapsed to one column.
func TestGridAutoFitClassIsScannerVisible(t *testing.T) {
	t.Parallel()

	t.Run("class literal references the custom property", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Grid(GridProps{Cols: GridColsAutoFit, MinColWidth: "190px"}))
		utils.AssertContains(t, output, "[grid-template-columns:repeat(auto-fit,minmax(var(--tc-grid-min),1fr))]")
		utils.AssertNotContains(t, output, "minmax(190px")
	})

	t.Run("no style attribute for fixed column counts", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Grid(GridProps{Cols: GridCols3}))
		utils.AssertNotContains(t, output, "--tc-grid-min")
	})
}
