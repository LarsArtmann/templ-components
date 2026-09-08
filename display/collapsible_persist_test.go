package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// persistRender renders a CollapsibleSection with the given props.
func persistRender(t *testing.T, props CollapsibleSectionProps) string {
	t.Helper()

	return utils.Render(t, CollapsibleSection(props))
}

func TestCollapsiblePersist_OffByDefault(t *testing.T) {
	t.Parallel()

	output := persistRender(t, CollapsibleSectionProps{
		Title:      "Section",
		StorageKey: "key",
		BaseProps:  utils.BaseProps{Nonce: "n1"},
	})

	if strings.Contains(output, "tcCollapsiblePersist") {
		t.Error("default props must not ship the persistence script")
	}
}

func TestCollapsiblePersist_RequiresStorageKey(t *testing.T) {
	t.Parallel()

	output := persistRender(t, CollapsibleSectionProps{
		Title:        "Section",
		PersistState: true,
		BaseProps:    utils.BaseProps{Nonce: "n1"},
	})

	if strings.Contains(output, "tcCollapsiblePersist") {
		t.Error("PersistState without StorageKey must not ship the script")
	}

	if !strings.Contains(output, `data-collapsible`) && strings.Contains(output, "details") {
		t.Log("no data-collapsible attribute without StorageKey — expected")
	}
}

func TestCollapsiblePersist_RequiresNonce(t *testing.T) {
	t.Parallel()

	output := persistRender(t, CollapsibleSectionProps{
		Title:        "Section",
		StorageKey:   "key",
		PersistState: true,
	})

	if strings.Contains(output, "tcCollapsiblePersist") {
		t.Error("PersistState without Nonce must not ship an un-nonced inline script")
	}
}

func TestCollapsiblePersist_ScriptShipsWithNonce(t *testing.T) {
	t.Parallel()

	output := persistRender(t, CollapsibleSectionProps{
		Title:        "Section",
		StorageKey:   "pipeline-portals",
		PersistState: true,
		BaseProps:    utils.BaseProps{Nonce: "persist-nonce"},
	})

	for _, want := range []string{
		`<script nonce="persist-nonce">`,
		"tcCollapsiblePersistAttached",
		`data-collapsible`,
		`addEventListener("toggle"`,
		"localStorage.setItem",
		"localStorage.getItem",
		", true)", // capture-phase listener: toggle does not bubble
	} {
		if !strings.Contains(output, want) {
			t.Errorf("persistence script should contain %q", want)
		}
	}
}

func TestCollapsiblePersist_RestoresBeforeGuard(t *testing.T) {
	t.Parallel()

	// The applyAll() call must run on every render (restores state for
	// HTMX-swapped re-renders), while the toggle listener attaches once
	// behind the singleton guard.
	output := persistRender(t, CollapsibleSectionProps{
		Title:        "Section",
		StorageKey:   "key",
		PersistState: true,
		BaseProps:    utils.BaseProps{Nonce: "n"},
	})

	applyIdx := strings.Index(output, "applyAll();")
	guardIdx := strings.Index(output, "tcCollapsiblePersistAttached")

	if applyIdx == -1 || guardIdx == -1 {
		t.Fatal("script should contain applyAll() and the singleton guard")
	}

	if applyIdx > guardIdx {
		t.Error("applyAll() must run before the singleton guard so re-renders restore state")
	}
}
