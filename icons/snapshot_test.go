// Package icons provides rendering tests for icon components.
package icons

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestIconRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		icon  Name
		class string
	}{
		{"home with class", Home, "h-5 w-5"},
		{"check default", Check, ""},
		{"spinner with color", Spinner, "text-blue-500 animate-spin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Icon(tt.icon, tt.class))
			utils.AssertContains(t, output, "<svg")
			utils.AssertContains(t, output, "</svg>")
			if tt.class != "" {
				for c := range strings.FieldsSeq(tt.class) {
					utils.AssertContains(t, output, c)
				}
			}
		})
	}
}

func TestAllIconsRender(t *testing.T) {
	t.Parallel()
	for _, name := range allIconNames() {
		t.Run(string(name), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Icon(name, "h-5 w-5"))
			utils.AssertContains(t, output, "<svg")
			utils.AssertContains(t, output, "</svg>")
			utils.AssertContains(t, output, "h-5 w-5")
		})
	}
}

func TestIconWithStrokeWidth(t *testing.T) {
	t.Parallel()

	t.Run("custom stroke-width renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, IconWithStrokeWidth(Home, "h-5 w-5", 2.5))
		utils.AssertContains(t, output, `<svg`)
		utils.AssertContains(t, output, `stroke-width="2.5"`)
	})

	t.Run("default Icon uses 1.5 stroke-width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Home, "h-5 w-5"))
		utils.AssertContains(t, output, `stroke-width="1.5"`)
	})

	t.Run("spinner ignores custom stroke-width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, IconWithStrokeWidth(Spinner, "h-5 w-5", 3.0))
		utils.AssertContains(t, output, "animate-spin")
		utils.AssertNotContains(t, output, `stroke-width="3"`)
	})
}

func TestIconPathJS(t *testing.T) {
	t.Parallel()
	t.Run("known icon produces path elements", func(t *testing.T) {
		t.Parallel()
		result := IconPathJS(Home)
		if result == "" {
			t.Error("IconPathJS(Home) returned empty string")
		}
	})
	t.Run("multi-path icon produces multiple paths", func(t *testing.T) {
		t.Parallel()
		result := IconPathJS(ChevronDown)
		if result == "" {
			t.Error("IconPathJS(ChevronDown) returned empty string")
		}
	})
}
