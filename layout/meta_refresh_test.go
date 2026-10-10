package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestMetaRefresh(t *testing.T) {
	t.Parallel()

	t.Run("redirect after delay", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MetaRefresh(5, "/login?retry=1"))
		utils.AssertContains(t, output, `<meta http-equiv="refresh" content="5;url=/login?retry=1">`)
	})

	t.Run("reload same page when url empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MetaRefresh(30, ""))
		utils.AssertContains(t, output, `<meta http-equiv="refresh" content="30">`)
	})

	t.Run("zero delay is immediate", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MetaRefresh(0, "/next"))
		utils.AssertContains(t, output, `content="0;url=/next"`)
	})

	t.Run("negative delay clamps to zero", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MetaRefresh(-3, "/next"))
		utils.AssertContains(t, output, `content="0;url=/next"`)
	})
}

func TestGoldenSweepMetaRefresh(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "meta_refresh_redirect", HTML: utils.Render(t, MetaRefresh(5, "/login?retry=1"))},
		{Name: "meta_refresh_reload", HTML: utils.Render(t, MetaRefresh(30, ""))},
	})
}
