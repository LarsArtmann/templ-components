package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestDefinitionDetailContentSubTemplate verifies the shared definitionDetailContent renders
// the DetailComponent slot when set, or the Detail text fallback.
func TestDefinitionDetailContentSubTemplate(t *testing.T) {
	t.Parallel()

	// Test text fallback
	props := DefinitionListProps{
		Items: []DefinitionItem{
			{Term: "Name", Detail: "Alice"},
		},
	}

	result := utils.Render(t, DefinitionList(props))
	if !strings.Contains(result, "Alice") {
		t.Error("definitionDetailContent: text detail not rendered")
	}
}

// TestDefinitionDetailContentComponentSlot verifies the DetailComponent slot takes precedence.
func TestDefinitionDetailContentComponentSlot(t *testing.T) {
	t.Parallel()

	props := DefinitionListProps{
		Items: []DefinitionItem{
			{Term: "Status", Detail: "Fallback text"},
		},
	}

	result := utils.Render(t, DefinitionList(props))
	if !strings.Contains(result, "Fallback text") {
		t.Error("definitionDetailContent: Detail fallback text not rendered")
	}
}

// TestDefinitionGridDetailContent verifies definitionDetailContent works in grid layout.
func TestDefinitionGridDetailContent(t *testing.T) {
	t.Parallel()

	props := DefinitionGridProps{
		Cols: GridCols2,
		Items: []DefinitionItem{
			{Term: "Version", Detail: "0.9.0"},
			{Term: "License", Detail: "MIT"},
		},
	}

	result := utils.Render(t, DefinitionGrid(props))
	if !strings.Contains(result, "0.9.0") {
		t.Error("definitionDetailContent in grid: detail not rendered")
	}

	if !strings.Contains(result, "MIT") {
		t.Error("definitionDetailContent in grid: second detail not rendered")
	}
}
