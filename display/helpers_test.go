// Package display provides tests for display components like Badge, Card, Modal, and EmptyState.
package display

import "testing"

func TestMapStatusToBadgeType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		status string
		want   BadgeType
	}{
		{"active", BadgeSuccess},
		{"error", BadgeError},
		{"warning", BadgeWarning},
		{"info", BadgeInfo},
		{"primary", BadgePrimary},
		{"unknown-random", BadgeNeutral},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			t.Parallel()
			got := mapStatusToBadgeType(tt.status)
			if got != tt.want {
				t.Errorf("mapStatusToBadgeType(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestBadgeSizeClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		size BadgeSize
		want string
	}{
		{BadgeSizeSm, "px-2 py-0.5 text-xs"},
		{BadgeSizeMd, "px-2.5 py-0.5 text-xs"},
		{BadgeSizeLg, "px-3 py-1 text-sm"},
	}
	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			t.Parallel()
			got := badgeSizeClass(tt.size)
			if got != tt.want {
				t.Errorf("badgeSizeClass(%q) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestCardPaddingClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		padding string
		want    string
	}{
		{"none", ""},
		{"sm", "px-3 py-3"},
		{"md", "px-4 py-5 sm:p-6"},
		{"lg", "px-6 py-6"},
	}
	for _, tt := range tests {
		t.Run(tt.padding, func(t *testing.T) {
			t.Parallel()
			got := cardPaddingClass(tt.padding)
			if got != tt.want {
				t.Errorf("cardPaddingClass(%q) = %q, want %q", tt.padding, got, tt.want)
			}
		})
	}
}
