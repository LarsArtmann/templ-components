// Package icons provides rendering tests for icon components.
package icons

import (
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
				for _, c := range splitClasses(tt.class) {
					utils.AssertContains(t, output, c)
				}
			}
		})
	}
}

func TestIconAttrs(t *testing.T) {
	t.Parallel()

	t.Run("with aria label returns aria-label", func(t *testing.T) {
		t.Parallel()
		attrs := IconAttrs("Menu")
		if attrs["aria-label"] != "Menu" {
			t.Errorf("aria-label = %v, want Menu", attrs["aria-label"])
		}
		if _, ok := attrs["aria-hidden"]; ok {
			t.Error("should not have aria-hidden when aria-label is set")
		}
	})

	t.Run("without aria label returns aria-hidden", func(t *testing.T) {
		t.Parallel()
		attrs := IconAttrs("")
		if attrs["aria-hidden"] != "true" {
			t.Errorf("aria-hidden = %v, want true", attrs["aria-hidden"])
		}
		if _, ok := attrs["aria-label"]; ok {
			t.Error("should not have aria-label when empty")
		}
	})
}

func TestAllIconsRender(t *testing.T) {
	t.Parallel()
	allIcons := []Name{
		Home, Users, Folder, Document, Search, Settings, Chart, Inbox,
		Check, X, Plus, Minus, ChevronRight, ChevronLeft, ChevronDown, ChevronUp,
		ArrowRight, ArrowLeft, Refresh, ExternalLink, Download, Upload,
		Trash, Edit, Eye, EyeOff, Lock, Unlock, Menu, Bell,
		Calendar, Clock, Location, Phone, Mail, Globe, Sun, Moon,
		Spinner, Exclamation, Information, Question,
	}
	for _, name := range allIcons {
		t.Run(string(name), func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Icon(name, "h-5 w-5"))
			utils.AssertContains(t, output, "<svg")
			utils.AssertContains(t, output, "</svg>")
			utils.AssertContains(t, output, "h-5 w-5")
		})
	}
}

func splitClasses(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}
