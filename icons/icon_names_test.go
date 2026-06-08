package icons

import (
	"fmt"
	"strings"
	"testing"
)

func TestIconNames(t *testing.T) {
	t.Parallel()

	for _, name := range allIconNames() {
		t.Run(string(name), func(t *testing.T) {
			t.Parallel()
			if name == "" {
				t.Errorf("icon name should not be empty")
			}
		})
	}
}

func TestIconCount(t *testing.T) {
	t.Parallel()
	// allIconNames includes Spinner which is handled specially (not in iconPathData)
	if len(allIconNames()) != len(iconPathData)+1 {
		t.Errorf(
			"allIconNames has %d entries, iconPathData has %d (+1 Spinner) — they must stay in sync",
			len(allIconNames()),
			len(iconPathData),
		)
	}
}

func TestAllIconNamesCoversIconPathData(t *testing.T) {
	t.Parallel()
	nameSet := make(map[Name]bool, len(allIconNames()))
	for _, n := range allIconNames() {
		nameSet[n] = true
	}
	for name := range iconPathData {
		if !nameSet[name] {
			t.Errorf("iconPathData has %q but allIconNames does not", name)
		}
	}
}

func TestIconPathsNoEmptySegments(t *testing.T) {
	t.Parallel()
	for name, data := range iconPathData {
		parts := strings.SplitSeq(data, "|")
		for p := range parts {
			if p == "" {
				t.Errorf("icon %q has empty path segment in iconPathData", name)
			}
		}
	}
}

func TestIconPathDataNoPipeInSVGPaths(t *testing.T) {
	t.Parallel()
	for name, data := range iconPathData {
		if strings.HasPrefix(data, "|") || strings.HasSuffix(data, "|") {
			t.Errorf("icon %q has leading/trailing | separator: %q", name, data)
		}
		if strings.Contains(data, "||") {
			t.Errorf("icon %q has double | separator: %q", name, data)
		}
	}
}

func TestIconPathJSProducesValidHTML(t *testing.T) {
	t.Parallel()
	for name := range iconPathData {
		result := IconPathJS(name)
		if !strings.HasPrefix(result, `<path`) {
			t.Errorf("IconPathJS(%q) should start with <path, got: %s", name, result[:50])
		}
		if !strings.HasSuffix(result, `"/>`) {
			t.Errorf("IconPathJS(%q) should end with \"/>, got: %s", name, result[len(result)-20:])
		}
	}
}

func TestIconPathsPanicsOnUnknown(t *testing.T) {
	t.Parallel()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for unknown icon name")
		}
		msg := fmt.Sprintf("%v", r)
		if !strings.Contains(msg, "unknown icon name") {
			t.Errorf("panic message should mention unknown icon name, got: %s", msg)
		}
	}()
	iconPaths(Name("nonexistent"))
}
