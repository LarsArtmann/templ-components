package visualtest

import "testing"

// TEMPORARY (M06): targeted refresh of the /display route goldens only.
// TestDemoRouteGoldens is a single unfiltered table, so a blanket -update would
// also rebaseline unrelated routes drifting from another lane's in-flight work.
// Delete after use; do not commit.
func TestDisplayGoldenRefresh(t *testing.T) {
	base := demoRouteBase()
	assertRouteScreenshot(t, "display_light", base+"/display", false, false, ViewportDesktop, true)
	assertRouteScreenshot(t, "display_dark", base+"/display", true, false, ViewportDesktop, true)
}
