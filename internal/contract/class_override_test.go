package contract

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// TestClassOverrideWins pins the theming contract for flagship components:
// BaseProps.Class must WIN conflicts with the component's default classes
// (tailwind-merge semantics via utils.Class). This is the documented consumer
// theming path — "pass props.Class and it overrides" — and a regression here
// silently breaks every consumer theme.
//
// Strategy: give each component a Class that CONFLICTS with one of its
// default utilities (same tailwind group, different value) plus a unique
// marker class; assert the marker is present and the default is gone.
func TestClassOverrideWins(t *testing.T) {
	t.Parallel()

	const marker = "marker-z9q"

	tests := []struct {
		name        string
		render      func() string
		wantPresent []string // must appear in class="..."
		wantAbsent  []string // must NOT appear (overridden by Class)
	}{
		{
			name: "badge size override",
			render: func() string {
				props := display.DefaultBadgeProps()
				props.Text = "x"
				props.Class = "text-5xl " + marker

				return utils.Render(t, display.Badge(props))
			},
			wantPresent: []string{marker, "text-5xl"},
			wantAbsent:  []string{"text-xs"},
		},
		{
			name: "card padding override",
			render: func() string {
				props := display.DefaultCardProps()
				props.Title = "x"
				props.Class = "p-10 " + marker

				return utils.Render(t, display.Card(props))
			},
			wantPresent: []string{marker, "p-10"},
			wantAbsent:  []string{"p-4"},
		},
		{
			name: "button variant override",
			render: func() string {
				props := display.DefaultButtonProps()
				props.Text = "x"
				props.Class = "bg-purple-700 " + marker

				return utils.Render(t, display.Button(props))
			},
			wantPresent: []string{marker, "bg-purple-700"},
			wantAbsent:  []string{"bg-blue-600"},
		},
		{
			name: "input text-size override",
			render: func() string {
				return utils.Render(t, forms.Input(forms.InputProps{ //nolint:exhaustruct // minimal
					Name:      "x",
					BaseProps: utils.BaseProps{Class: "text-xl " + marker},
				}))
			},
			wantPresent: []string{marker, "text-xl"},
			wantAbsent:  []string{"text-sm"},
		},
		{
			name: "alert border override",
			render: func() string {
				return utils.Render(t, feedback.Alert(feedback.AlertProps{ //nolint:exhaustruct // minimal
					Title:     "x",
					BaseProps: utils.BaseProps{Class: "border-4 " + marker},
				}))
			},
			wantPresent: []string{marker, "border-4"},
			wantAbsent:  []string{"border"},
		},
		{
			name: "pagination gap override",
			render: func() string {
				return utils.Render(t, navigation.Pagination(navigation.PaginationProps{ //nolint:exhaustruct // minimal
					CurrentPage: 1,
					TotalPages:  3,
					BaseURL:     "/",
					BaseProps:   utils.BaseProps{Class: "gap-8 " + marker},
				}))
			},
			wantPresent: []string{marker, "gap-8"},
			wantAbsent:  []string{"gap-1"},
		},
		{
			name: "table wrapper width override",
			render: func() string {
				return utils.Render(t, display.Table(display.TableProps{ //nolint:exhaustruct // minimal
					Headers:   []string{"a"},
					BaseProps: utils.BaseProps{Class: "w-24 " + marker},
				}))
			},
			wantPresent: []string{marker, "w-24"},
			wantAbsent:  []string{"w-full"},
		},
		{
			name: "progressbar height override",
			render: func() string {
				return utils.Render(t, feedback.ProgressBar(feedback.ProgressBarProps{ //nolint:exhaustruct // minimal
					Current:   1,
					Total:     2,
					BaseProps: utils.BaseProps{Class: "h-1 " + marker},
				}))
			},
			wantPresent: []string{marker, "h-1"},
			wantAbsent:  []string{"h-2"},
		},
		{
			name: "container padding override",
			render: func() string {
				return utils.Render(t, layout.Container(layout.ContainerProps{ //nolint:exhaustruct // minimal
					BaseProps: utils.BaseProps{Class: "max-w-md " + marker},
				}))
			},
			wantPresent: []string{marker, "max-w-md"},
			wantAbsent:  []string{"max-w-7xl"},
		},
		{
			name: "empty state icon color override",
			render: func() string {
				return utils.Render(t, display.EmptyState(display.EmptyStateProps{ //nolint:exhaustruct // minimal
					Title:     "x",
					BaseProps: utils.BaseProps{Class: "py-20 " + marker},
				}))
			},
			wantPresent: []string{marker, "py-20"},
			wantAbsent:  []string{"py-12"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			html := tc.render()
			tokens := classTokens(html)

			for _, want := range tc.wantPresent {
				if !tokens[want] {
					t.Errorf(
						"consumer class %q missing from output — the Class-override-wins contract is broken:\n%s",
						want,
						html,
					)
				}
			}

			for _, absent := range tc.wantAbsent {
				if tokens[absent] {
					t.Errorf(
						"default class %q survived a consumer override — tailwind-merge did not resolve the conflict:\n%s",
						absent,
						html,
					)
				}
			}
		})
	}
}

// classTokens extracts every whitespace-separated token from every
// class="..." attribute in an HTML string (token-exact, so "text-sm" does
// not match "sm:text-sm" or "min-w-full" vs "w-full").
func classTokens(html string) map[string]bool {
	tokens := map[string]bool{}

	for attr := range strings.SplitSeq(html, `class="`) {
		// Everything before the closing quote of this attribute.
		value, _, found := strings.Cut(attr, `"`)
		if !found {
			continue
		}

		for token := range strings.SplitSeq(value, " ") {
			if token = strings.TrimSpace(token); token != "" {
				tokens[token] = true
			}
		}
	}

	return tokens
}
