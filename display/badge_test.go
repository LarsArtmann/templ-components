// Package display provides tests for display components like Badge, Card, Modal, and EmptyState.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	badgeTextActive     = activeBadgeText // reuse existing
	badgeTextFailed     = "Failed"
	badgeTextBeta       = "Beta"
	badgeTextRunning    = "Running"
	badgeTextTag        = "Tag"
	badgeTextCustom     = "Custom"
	badgeTextOK         = "OK"
	cssClassRoundedFull = "rounded-full"
	cssClassMl2         = "ml-2"
	cssClassMt2         = "mt-2"
	statusActive        = "active"
	statusError         = "error"
	statusWarning       = "warning"
	statusUnknown       = "unknown"
)

func testBadgeProps(text string, badgeType BadgeType) BadgeProps {
	return BadgeProps{
		BaseProps: utils.BaseProps{},
		Text:      text,
		Type:      badgeType,
		Size:      BadgeSizeMD,
		Pill:      false,
		Dot:       false,
	}
}

func TestBadgeRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		props       BadgeProps
		wantContain []string
	}{
		{
			name:        "basic success badge",
			props:       testBadgeProps(badgeTextActive, BadgeSuccess),
			wantContain: []string{badgeTextActive, "bg-green-100", "text-green-800"},
		},
		{
			name: "error badge with dot",

			props: BadgeProps{
				BaseProps: utils.BaseProps{},
				Text:      badgeTextFailed,
				Type:      BadgeError,
				Size:      BadgeSizeMD,
				Pill:      false,
				Dot:       true,
			},
			wantContain: []string{badgeTextFailed, "bg-red-100", cssClassRoundedFull, "bg-red-500"},
		},
		{
			name: "pill badge",

			props: BadgeProps{
				BaseProps: utils.BaseProps{},
				Text:      badgeTextBeta,
				Type:      BadgeInfo,
				Size:      BadgeSizeMD,
				Pill:      true,
				Dot:       false,
			},
			wantContain: []string{badgeTextBeta, cssClassRoundedFull, "bg-blue-50"},
		},
		{
			name: "with custom id and class",
			props: BadgeProps{
				BaseProps: utils.BaseProps{
					ID:        "status-badge",
					Class:     cssClassMl2,
					Attrs:     nil,
					AriaLabel: "",
				},
				Text: badgeTextRunning,
				Type: BadgePrimary,
				Size: BadgeSizeMD,
				Pill: false,
				Dot:  false,
			},
			wantContain: []string{`id="status-badge"`, cssClassMl2, "bg-blue-100"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, Badge(tt.props))
			for _, want := range tt.wantContain {
				utils.AssertContains(t, output, want)
			}
		})
	}
}

func TestStatusBadgeRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
		want   string
	}{
		{name: "active maps to success", status: statusActive, want: "bg-green-100"},
		{name: "error maps to error", status: statusError, want: "bg-red-100"},
		{name: "warning maps to warning", status: statusWarning, want: "bg-yellow-100"},
		{name: "unknown maps to neutral", status: statusUnknown, want: "bg-gray-100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, StatusBadge(tt.status))
			utils.AssertContains(t, output, tt.want)
		})
	}
}

func TestBadgeFeatures(t *testing.T) {
	t.Parallel()

	t.Run("pill badge renders rounded-full", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: badgeTextTag,
			Pill: true,
		}))
		utils.AssertContains(t, output, cssClassRoundedFull)
	})

	t.Run("badge with dot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: badgeTextActive,
			Dot:  true,
			Type: BadgeSuccess,
		}))
		utils.AssertContains(t, output, badgeTextActive)
		utils.AssertContains(t, output, cssClassRoundedFull)
	})

	t.Run("all badge types render", func(t *testing.T) {
		t.Parallel()

		for _, bt := range []BadgeType{BadgePrimary, BadgeSuccess, BadgeWarning, BadgeError, BadgeInfo, BadgeNeutral} {
			output := utils.Render(t, Badge(BadgeProps{
				Text: string(bt),
				Type: bt,
			}))
			utils.AssertContains(t, output, string(bt))
		}
	})

	t.Run("badge with custom class and id", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{ID: "my-badge", Class: cssClassMt2},
			Text:      badgeTextCustom,
		}))
		utils.AssertContains(t, output, `id="my-badge"`)
		utils.AssertContains(t, output, cssClassMt2)
	})

	t.Run("badge with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{AriaLabel: "Status indicator"},
			Text:      badgeTextOK,
		}))
		utils.AssertContains(t, output, `aria-label="Status indicator"`)
	})

	t.Run("badge with pill and dot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: badgeTextBeta,
			Type: BadgePrimary,
			Pill: true,
			Dot:  true,
		}))
		utils.AssertContains(t, output, cssClassRoundedFull)
		utils.AssertContains(t, output, "bg-blue-500")
		utils.AssertContains(t, output, badgeTextBeta)
	})
}
