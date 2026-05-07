// Package display provides tests for display components like Badge, Card, Modal, and EmptyState.
package display

import (
	"fmt"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

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
			utils.AssertEqual(
				t,
				fmt.Sprintf("mapStatusToBadgeType(%q)", tt.status),
				mapStatusToBadgeType(tt.status),
				tt.want,
			)
		})
	}
}

func TestBadgeSizeClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		size BadgeSize
		want string
	}{
		{BadgeSizeSM, "px-2 py-0.5 text-xs"},
		{BadgeSizeMD, "px-2.5 py-0.5 text-xs"},
		{BadgeSizeLG, "px-3 py-1 text-sm"},
	}
	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("badgeSizeClass(%q)", tt.size),
				badgeSizeClass(tt.size),
				tt.want,
			)
		})
	}
}

func TestCardPaddingClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		padding CardPadding
		want    string
	}{
		{CardPaddingNone, ""},
		{CardPaddingSM, "px-3 py-3"},
		{CardPaddingMD, "px-4 py-5 sm:p-6"},
		{CardPaddingLG, "px-6 py-6"},
	}
	for _, tt := range tests {
		t.Run(string(tt.padding), func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("cardPaddingClass(%q)", tt.padding),
				cardPaddingClass(tt.padding),
				tt.want,
			)
		})
	}
}
