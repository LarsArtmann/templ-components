// Package display provides tests for display components like Badge, Card, Modal, and EmptyState.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

//nolint:exhaustruct
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
			props:       testBadgeProps("Active", BadgeSuccess),
			wantContain: []string{"Active", "bg-green-100", "text-green-800"},
		},
		{
			name: "error badge with dot",
			//nolint:exhaustruct
			props: BadgeProps{
				BaseProps: utils.BaseProps{},
				Text:      "Failed",
				Type:      BadgeError,
				Size:      BadgeSizeMD,
				Pill:      false,
				Dot:       true,
			},
			wantContain: []string{"Failed", "bg-red-100", "rounded-full", "bg-red-500"},
		},
		{
			name: "pill badge",
			//nolint:exhaustruct
			props: BadgeProps{
				BaseProps: utils.BaseProps{},
				Text:      "Beta",
				Type:      BadgeInfo,
				Size:      BadgeSizeMD,
				Pill:      true,
				Dot:       false,
			},
			wantContain: []string{"Beta", "rounded-full", "bg-indigo-100"},
		},
		{
			name: "with custom id and class",
			props: BadgeProps{
				BaseProps: utils.BaseProps{
					ID:        "status-badge",
					Class:     "ml-2",
					Attrs:     nil,
					AriaLabel: "",
				},
				Text: "Running",
				Type: BadgePrimary,
				Size: BadgeSizeMD,
				Pill: false,
				Dot:  false,
			},
			wantContain: []string{`id="status-badge"`, "ml-2", "bg-blue-100"},
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
		{name: "active maps to success", status: "active", want: "bg-green-100"},
		{name: "error maps to error", status: "error", want: "bg-red-100"},
		{name: "warning maps to warning", status: "warning", want: "bg-yellow-100"},
		{name: "unknown maps to neutral", status: "unknown", want: "bg-gray-100"},
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
			Text: "Tag",
			Pill: true,
		}))
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("badge with dot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Active",
			Dot:  true,
			Type: BadgeSuccess,
		}))
		utils.AssertContains(t, output, "Active")
		utils.AssertContains(t, output, "rounded-full")
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
			BaseProps: utils.BaseProps{ID: "my-badge", Class: "mt-2"},
			Text:      "Custom",
		}))
		utils.AssertContains(t, output, `id="my-badge"`)
		utils.AssertContains(t, output, "mt-2")
	})

	t.Run("badge with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			BaseProps: utils.BaseProps{AriaLabel: "Status indicator"},
			Text:      "OK",
		}))
		utils.AssertContains(t, output, `aria-label="Status indicator"`)
	})

	t.Run("badge with pill and dot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(BadgeProps{
			Text: "Beta",
			Type: BadgePrimary,
			Pill: true,
			Dot:  true,
		}))
		utils.AssertContains(t, output, "rounded-full")
		utils.AssertContains(t, output, "bg-blue-500")
		utils.AssertContains(t, output, "Beta")
	})
}
