package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultSliderProps(t *testing.T) {
	t.Parallel()
	p := DefaultSliderProps()
	if p.Min != 0 || p.Max != 100 || p.Step != 1 {
		t.Error("expected default Min=0, Max=100, Step=1")
	}
}

func TestSliderBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:  "volume",
		Label: "Volume",
		Min:   0,
		Max:   100,
		Value: 50,
		Step:  5,
	}))
	utils.AssertContains(t, output, `type="range"`)
	utils.AssertContains(t, output, `name="volume"`)
	utils.AssertContains(t, output, `min="0"`)
	utils.AssertContains(t, output, `max="100"`)
	utils.AssertContains(t, output, `value="50"`)
	utils.AssertContains(t, output, `step="5"`)
}

func TestSliderShowValue(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:      "volume",
		Label:     "Volume",
		Value:     75,
		ShowValue: true,
	}))
	utils.AssertContains(t, output, "75")
}

func TestSliderDisabled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:     "vol",
		Value:    50,
		Disabled: true,
	}))
	utils.AssertContains(t, output, "disabled")
}

func TestSliderHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:     "vol",
		Value:    50,
		HelpText: "Adjust the volume",
	}))
	utils.AssertContains(t, output, "Adjust the volume")
}

func TestSliderError(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:  "vol",
		Value: 50,
		Error: "Volume too high",
		BaseProps: utils.BaseProps{
			ID: "vol-slider",
		},
	}))
	utils.AssertContains(t, output, "Volume too high")
	utils.AssertContains(t, output, `aria-invalid="true"`)
}

func TestSliderDarkMode(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:  "vol",
		Value: 50,
	}))
	utils.AssertContains(t, output, "dark:bg-gray-700")
	utils.AssertContains(t, output, "dark:accent-blue-400")
}

func TestSliderDefaultStepWhenZero(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:  "vol",
		Value: 50,
		Step:  0,
	}))
	utils.AssertContains(t, output, `step="1"`)
}
