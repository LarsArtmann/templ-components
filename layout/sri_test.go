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
			version: "2.0.6",
			ext:     "main",
			want:    "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
		},
		{
			name:    "response-targets 2.0.6",
			version: "2.0.6",
			ext:     "response-targets",
			want:    "sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd",
		},
		{name: "unknown version", version: "99.99.99", ext: "main", want: ""},
		{
			name:    "unknown ext",
			version: "2.0.6",
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
