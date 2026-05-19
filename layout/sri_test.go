package layout

import (
	"fmt"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestHtmxSRI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		version string
		ext     string
		want    string
	}{
		{
			name:    "main 2.0.6",
			version: htmxVersion206,
			ext:     "main",
			want:    "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
		},
		{
			name:    "response-targets 2.0.6",
			version: htmxVersion206,
			ext:     "response-targets",
			want:    sriResponseTargets206,
		},
		{name: "unknown version", version: "99.99.99", ext: "main", want: ""},
		{
			name:    "unknown ext",
			version: htmxVersion206,
			ext:     "unknown",
			want:    "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("htmxSRI(%q, %q)", tt.version, tt.ext),
				htmxSRI(tt.version, tt.ext),
				tt.want,
			)
		})
	}
}
