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
		version string
		want    string
	}{
		{
			name:    "pinned 2.0.10 returns main SRI hash",
			version: defaultHTMXVersion,
			want:    "sha384-H5SrcfygHmAuTDZphMHqBJLc3FhssKjG7w/CeCpFReSfwBWDTKpkzPP8c+cLsK+V",
		},
		{
			name:    "unknown version falls back to default SRI",
			version: "99.99.99",
			want:    "sha384-H5SrcfygHmAuTDZphMHqBJLc3FhssKjG7w/CeCpFReSfwBWDTKpkzPP8c+cLsK+V",
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
		"htmxScriptURL(defaultHTMXVersion)",
		htmxScriptURL(defaultHTMXVersion),
		"https://unpkg.com/htmx.org@2.0.10",
	)
}

func TestResponseTargetsCDNURL(t *testing.T) {
	t.Parallel()
	utils.AssertEqual(
		t,
		"responseTargetsCDNURL",
		responseTargetsCDNURL,
		"https://unpkg.com/htmx-ext-response-targets@2.0.4/dist/response-targets.min.js",
	)
}
