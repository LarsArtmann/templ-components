package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestLoadingDialogContract pins the modal-loading semantics: native dialog
// shell with an accessible name, busy state on the content region, no close
// button (nothing to close while loading), and clamped progress.
func TestLoadingDialogContract(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, LoadingDialog(LoadingDialogProps{
		ID: "d1", Title: "Saving", Message: "Saving...",
		ShowProgress: true, Progress: 140,
	}))

	utils.AssertContains(t, output, `<dialog`)
	utils.AssertContains(t, output, `aria-label="Saving..."`)
	utils.AssertContains(t, output, `aria-busy="true"`)
	utils.AssertContains(t, output, `closedby="any"`)
	utils.AssertContains(t, output, "100%")
	utils.AssertContains(t, output, "width: 100%")
	utils.AssertNotContains(t, output, "140%")
	utils.AssertNotContains(t, output, `aria-modal="true"`)

	// The dialog has no header, so no close button renders — a blocking
	// loading state has nothing to dismiss via a header button. (The shared
	// overlay script still mentions data-tc-close in its delegated listener.)
	utils.AssertNotContains(t, output, `aria-label="Close"`)
	utils.AssertNotContains(t, output, "<button")
}

func TestLoadingDialogDefaults(t *testing.T) {
	t.Parallel()

	props := DefaultLoadingDialogProps()
	if props.Open {
		t.Error("DefaultLoadingDialogProps Open should be false")
	}

	output := utils.Render(t, LoadingDialog(props))
	utils.AssertContains(t, output, "Please wait...")
	utils.AssertNotContains(t, output, `data-tc-open="true"`)
}
