package datastar

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestIndicatorDefault(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Indicator(IndicatorProps{
		Signal: "saving",
	}))

	utils.AssertContains(t, output, `data-show="$saving"`)
	utils.AssertContains(t, output, `role="status"`)
	utils.AssertContains(t, output, `animate-spin`)
	utils.AssertContains(t, output, `motion-reduce:animate-none`)
}

func TestIndicatorWithNopSpinner(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Indicator(IndicatorProps{
		Signal:  "loading",
		Spinner: templ.NopComponent,
	}))

	utils.AssertContains(t, output, `data-show="$loading"`)
	utils.AssertNotContains(t, output, "animate-spin")
}

func TestIndicatorDarkModeCompliance(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Indicator(IndicatorProps{
		Signal: "saving",
	}))

	utils.AssertContains(t, output, "dark:border-blue-400")
}

func TestIndicatorDefaultProps(t *testing.T) {
	t.Parallel()

	props := DefaultIndicatorProps()

	if props.Signal != "" {
		t.Errorf("expected empty Signal, got %s", props.Signal)
	}
}
