// Package layout provides tests for layout components like Base, Minimal, ThemeScript, and ThemeToggle.
package layout

import "testing"

func TestHtmxIntegrityHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version 2.0.6",
			version: "2.0.6",
			want:    "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
		},
		{name: "unknown version", version: "99.99.99", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := htmxIntegrityHash(tt.version)
			if got != tt.want {
				t.Errorf("htmxIntegrityHash(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestHtmxResponseTargetsIntegrityHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version 2.0.6",
			version: "2.0.6",
			want:    "sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd",
		},
		{name: "unknown version", version: "99.99.99", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := htmxResponseTargetsIntegrityHash(tt.version)
			if got != tt.want {
				t.Errorf(
					"htmxResponseTargetsIntegrityHash(%q) = %q, want %q",
					tt.version,
					got,
					tt.want,
				)
			}
		})
	}
}
