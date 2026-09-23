package layout_test

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
)

// TestBaseHTMXNoneProvesHTMXOff pins the HTMX-off option (backlog #282): with
// HTMXSrc=HTMXNone, Base renders NO htmx runtime — no inline embedded source,
// no CDN script, no preconnect/dns-prefetch hints — while the page shell
// itself still renders normally. A page that wires everything through
// Datastar (or uses no wired components at all) then ships zero framework
// bytes.
func TestBaseHTMXNoneProvesHTMXOff(t *testing.T) {
	t.Parallel()

	props := layout.DefaultPageProps()
	props.Title = "htmx-off"
	props.HTMXSrc = layout.HTMXNone

	output := utils.Render(t, layout.Base(props))

	for _, forbidden := range []string{
		"htmx.min.js",
		"/htmx", // embedded inline script marker, CDN paths
		"preconnect",
		"dns-prefetch",
		"response-targets",
	} {
		if strings.Contains(output, forbidden) {
			t.Errorf("HTMXNone page contains %q; expected no htmx runtime artifacts", forbidden)
		}
	}

	for _, required := range []string{"<main", "Skip to main content"} {
		if !strings.Contains(output, required) {
			t.Errorf("HTMXNone page lost required shell markup %q", required)
		}
	}
}

// TestBaseHTMXSelfHostStillEmbeds guards the default against the new branch:
// HTMXSrc=HTMXSelfHost (the DefaultPageProps value) must still inline the
// runtime.
func TestBaseHTMXSelfHostStillEmbeds(t *testing.T) {
	t.Parallel()

	props := layout.DefaultPageProps()
	props.Title = "htmx-on"

	output := utils.Render(t, layout.Base(props))

	if !strings.Contains(output, "htmx") {
		t.Error("self-hosted Base no longer embeds the htmx runtime")
	}
}
