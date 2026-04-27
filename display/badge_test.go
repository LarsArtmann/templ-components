package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestBadgeRender(t *testing.T) {
	tests := []struct {
		name string
		props BadgeProps
		wantContain []string
	}{
		{
			name: "basic success badge",
			props: BadgeProps{Text: "Active", Type: BadgeSuccess},
			wantContain: []string{"Active", "bg-green-100", "text-green-800"},
		},
		{
			name: "error badge with dot",
			props: BadgeProps{Text: "Failed", Type: BadgeError, Dot: true},
			wantContain: []string{"Failed", "bg-red-100", "rounded-full", "bg-red-500"},
		},
		{
			name: "pill badge",
			props: BadgeProps{Text: "Beta", Type: BadgeInfo, Pill: true},
			wantContain: []string{"Beta", "rounded-full", "bg-indigo-100"},
		},
		{
			name: "with custom id and class",
			props: BadgeProps{
				BaseProps: utils.BaseProps{ID: "status-badge", Class: "ml-2"},
				Text: "Running",
				Type: BadgePrimary,
			},
			wantContain: []string{`id="status-badge"`, "ml-2", "bg-blue-100"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := utils.Render(t, Badge(tt.props))
			for _, want := range tt.wantContain {
				utils.AssertContains(t, output, want)
			}
		})
	}
}

func TestStatusBadgeRender(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   string
	}{
		{"active maps to success", "active", "bg-green-100"},
		{"error maps to error", "error", "bg-red-100"},
		{"warning maps to warning", "warning", "bg-yellow-100"},
		{"unknown maps to neutral", "unknown", "bg-gray-100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := utils.Render(t, StatusBadge(tt.status))
			utils.AssertContains(t, output, tt.want)
		})
	}
}
