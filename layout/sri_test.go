package layout

import "testing"

func TestHtmxIntegrityHash(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"2.0.6", "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm"},
		{"99.99.99", ""},
	}
	for _, tt := range tests {
		got := htmxIntegrityHash(tt.version)
		if got != tt.want {
			t.Errorf("htmxIntegrityHash(%q) = %q, want %q", tt.version, got, tt.want)
		}
	}
}

func TestHtmxResponseTargetsIntegrityHash(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"2.0.6", "sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd"},
		{"99.99.99", ""},
	}
	for _, tt := range tests {
		got := htmxResponseTargetsIntegrityHash(tt.version)
		if got != tt.want {
			t.Errorf("htmxResponseTargetsIntegrityHash(%q) = %q, want %q", tt.version, got, tt.want)
		}
	}
}
