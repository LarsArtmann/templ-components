package forms

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestSelectGroups(t *testing.T) {
	t.Parallel()

	t.Run("renders optgroups when Groups is set", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Select(SelectProps{
			Name:      "channel",
			BaseProps: utils.BaseProps{ID: "ch"},
			Groups: []SelectGroup{
				{Label: "Text Channels", Options: []SelectOption{
					{Value: "general", Label: "#general"},
					{Value: "random", Label: "#random"},
				}},
				{Label: "Voice Channels", Options: []SelectOption{
					{Value: "voice-1", Label: "Lounge"},
				}},
			},
		}))
		utils.AssertContains(t, html, `<optgroup label="Text Channels">`)
		utils.AssertContains(t, html, `<optgroup label="Voice Channels">`)
		utils.AssertContains(t, html, `value="general"`)
		utils.AssertContains(t, html, `value="voice-1"`)
	})

	t.Run("flat options when Groups is empty", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Select(SelectProps{
			Name:      "country",
			BaseProps: utils.BaseProps{ID: "ctry"},
			Options:   []SelectOption{{Value: "de", Label: "Germany"}},
		}))
		if strings.Contains(html, "<optgroup") {
			t.Error("flat options should not render optgroup")
		}
		utils.AssertContains(t, html, `value="de"`)
	})

	t.Run("Groups selected option is normalized", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Select(SelectProps{
			Name:      "x",
			BaseProps: utils.BaseProps{ID: "x"},
			Groups: []SelectGroup{
				{Label: "G1", Options: []SelectOption{
					{Value: "a", Label: "A", Disabled: true, Selected: true},
					{Value: "b", Label: "B", Selected: true},
				}},
			},
		}))
		// Disabled+Selected clears Selected, only the first valid Selected stays
		if strings.Contains(html, `value="a"`) && strings.Contains(html, `value="a" selected`) {
			t.Error("disabled option should not be selected")
		}
		utils.AssertContains(t, html, `value="b" selected`)
	})
}
