// Package layout provides tests for layout components like Base, Minimal, ThemeScript, and ThemeToggle.
package layout

import (
	"fmt"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestHtmxIntegrityHash(t *testing.T) {
	t.Parallel()
	testIntegrityHash(
		t,
		"htmxIntegrityHash",
		htmxIntegrityHash,
		"sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
	)
}

func TestHtmxResponseTargetsIntegrityHash(t *testing.T) {
	t.Parallel()
	testIntegrityHash(
		t,
		"htmxResponseTargetsIntegrityHash",
		htmxResponseTargetsIntegrityHash,
		"sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd",
	)
}

func testIntegrityHash(
	t *testing.T,
	funcName string,
	hashFunc func(string) string,
	version206Hash string,
) {
	t.Helper()
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{name: "version 2.0.6", version: "2.0.6", want: version206Hash},
		{name: "unknown version", version: "99.99.99", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("%s(%q)", funcName, tt.version),
				hashFunc(tt.version),
				tt.want,
			)
		})
	}
}
