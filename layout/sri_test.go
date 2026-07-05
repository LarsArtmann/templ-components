package layout

import (
	"fmt"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestHtmxMainSRI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		version HTMXVersion
		want    string
	}{
		{
			name:    "pinned 2.0.10 returns main SRI hash",
			version: defaultHTMXVersion,
			want:    htmxMainSRIDefault,
		},
		{
			name:    "unknown version returns empty (no wrong hash)",
			version: "99.99.99",
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("htmxMainSRI(%q)", tt.version),
				htmxMainSRI(tt.version),
				tt.want,
			)
		})
	}
}

func TestResponseTargetsSRI(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"sriResponseTargets",
		sriResponseTargets,
		"sha384-T41oglUPvXLGBVyRdZsVRxNWnOOqCynaPubjUVjxhsjFTKrFJGEMm3/0KGmNQ+Pg",
	)
}

func TestHtmxScriptURL(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"htmxScriptURL(defaultHTMXVersion, \"\")",
		htmxScriptURL(defaultHTMXVersion, ""),
		"https://cdn.jsdelivr.net/npm/htmx.org@2.0.10",
	)
}

func TestResponseTargetsURL(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"responseTargetsURL(\"\")",
		responseTargetsURL(""),
		"https://cdn.jsdelivr.net/npm/htmx-ext-response-targets@2.0.4/dist/response-targets.min.js",
	)
}

func TestHtmxScriptURLOverride(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"htmxScriptURL with unpkg override",
		htmxScriptURL(defaultHTMXVersion, "https://unpkg.com"),
		"https://unpkg.com/htmx.org@2.0.10",
	)
}

func TestHtmxScriptURLTrailingSlash(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"htmxScriptURL with trailing slash trimmed",
		htmxScriptURL(defaultHTMXVersion, "https://unpkg.com/"),
		"https://unpkg.com/htmx.org@2.0.10",
	)
}

func TestResponseTargetsURLOverride(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"responseTargetsURL with unpkg override",
		responseTargetsURL("https://unpkg.com"),
		"https://unpkg.com/htmx-ext-response-targets@2.0.4/dist/response-targets.min.js",
	)
}
