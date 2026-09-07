package forms

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDirtyGuardScript(t *testing.T) {
	t.Parallel()

	t.Run("renders a nonce-carrying script with the singleton guard", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DirtyGuard(DirtyGuardProps{BaseProps: utils.BaseProps{Nonce: "test-nonce"}}))
		utils.AssertContains(t, output, `nonce="test-nonce"`)
		utils.AssertContains(t, output, "tcDirtyGuardAttached")
		utils.AssertContains(t, output, "beforeunload")
		utils.AssertContains(t, output, "data-tc-dirty-guard")
	})

	t.Run("script is idempotent across swaps (single attach flag)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DirtyGuard(DirtyGuardProps{}))
		utils.AssertContains(t, output, "if (!window.tcDirtyGuardAttached)")
	})

	t.Run("submissions clear the dirty flag (wired forms included)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DirtyGuard(DirtyGuardProps{}))
		utils.AssertContains(t, output, "addEventListener('submit'")
	})

	t.Run("tracks forms in a WeakSet (swapped-in forms start clean)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DirtyGuard(DirtyGuardProps{}))
		utils.AssertContains(t, output, "new WeakSet()")
	})

	t.Run("no HTML beyond the script element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DirtyGuard(DirtyGuardProps{}))
		if strings.Contains(strings.ReplaceAll(output, "<script", ""), "</script>") &&
			strings.Contains(output, "<div") {
			t.Error("DirtyGuard must render only a script element")
		}
	})
}
