package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestDatePickerRender(t *testing.T) {
	t.Parallel()

	t.Run("basic date picker with label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DatePicker(DatePickerProps{
			BaseProps: utils.BaseProps{ID: "birthday"},
			Name:      "birthday",
			Label:     "Date of Birth",
			Value:     "1990-01-15",
		}))
		utils.AssertContains(t, output, `type="date"`)
		utils.AssertContains(t, output, `name="birthday"`)
		utils.AssertContains(t, output, `value="1990-01-15"`)
		utils.AssertContains(t, output, "Date of Birth")
	})

	t.Run("with min and max constraints", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DatePicker(DatePickerProps{
			BaseProps: utils.BaseProps{ID: "event-date"},
			Name:      "event_date",
			Min:       "2024-01-01",
			Max:       "2024-12-31",
		}))
		utils.AssertContains(t, output, `min="2024-01-01"`)
		utils.AssertContains(t, output, `max="2024-12-31"`)
	})

	t.Run("required with error and help text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DatePicker(DatePickerProps{
			BaseProps: utils.BaseProps{ID: "deadline"},
			Name:      "deadline",
			Label:     "Deadline",
			Required:  true,
			Error:     "Date is required",
			HelpText:  "Pick a future date",
		}))
		utils.AssertContains(t, output, `required`)
		utils.AssertContains(t, output, `aria-required="true"`)
		utils.AssertContains(t, output, "Date is required")
		utils.AssertContains(t, output, "Pick a future date")
		utils.AssertContains(t, output, `aria-invalid="true"`)
	})

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DatePicker(DatePickerProps{
			BaseProps: utils.BaseProps{ID: "locked"},
			Name:      "locked",
			Disabled:  true,
		}))
		utils.AssertContains(t, output, `disabled`)
	})

	t.Run("with BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DatePicker(DatePickerProps{
			BaseProps: utils.BaseProps{
				ID:        "dp-1",
				Class:     "w-full",
				AriaLabel: "Pick a date",
				Attrs:     templ.Attributes{"data-test": "date"},
			},
			Name:  "date",
			Label: "Date",
		}))
		utils.AssertContains(t, output, `id="dp-1"`)
		utils.AssertContains(t, output, "w-full")
		utils.AssertContains(t, output, `aria-label="Pick a date"`)
		utils.AssertContains(t, output, `data-test="date"`)
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DatePicker(DefaultDatePickerProps()))
		utils.AssertContains(t, output, `type="date"`)
	})
}
