package icons

import (
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

func TestIconPathsFallbackOnUnknown(t *testing.T) {
	t.Parallel()

	paths := iconPaths(Name("nonexistent"))

	questionPaths := iconPaths(Question)
	if len(paths) != len(questionPaths) {
		t.Errorf("unknown icon should fall back to Question, got %d paths, want %d", len(paths), len(questionPaths))
	}
}

func TestSpinnerIsSpecialAndFallsBackToQuestion(t *testing.T) {
	t.Parallel()

	if !specialIcons[Spinner] {
		t.Fatal("Spinner must be in specialIcons so it doesn't silently return wrong path data")
	}

	spinnerPaths := IconPathData(Spinner)

	questionPaths := IconPathData(Question)
	if len(spinnerPaths) != len(questionPaths) {
		t.Errorf(
			"IconPathData(Spinner) should fall back to Question, got %d paths, want %d",
			len(spinnerPaths),
			len(questionPaths),
		)
	}

	for i, p := range spinnerPaths {
		if p != questionPaths[i] {
			t.Errorf("IconPathData(Spinner)[%d] = %q, want Question fallback %q", i, p, questionPaths[i])
		}
	}
}

func TestAllExportedNameConstsHavePathDataOrAreSpecial(t *testing.T) {
	t.Parallel()
	// Every exported Name const must either be in iconPathData, be an alias
	// of a name in iconPathData, or be a special icon (Spinner). A const added
	// without a map entry would silently render as Question — this test catches
	// that regression.
	allNames := allIconNames()
	for _, name := range allNames {
		if specialIcons[name] {
			continue
		}

		if _, ok := iconPathData[name]; ok {
			continue
		}

		if alias, ok := iconAliases[name]; ok {
			if _, ok := iconPathData[alias]; ok {
				continue
			}
		}

		t.Errorf(
			"Name const %q has no iconPathData entry and is not special — would silently fall back to Question",
			name,
		)
	}
}
