package contract

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// TestAttrsPropagateToRoot pins the BaseProps.Attrs contract: consumer-supplied
// attributes are ALWAYS rendered on the component's root element, for every
// flagship component. Components spread `{ props.Attrs... }`; a regression that
// drops the spread silently strips consumer hooks (data-*, aria-*, hx-*), which
// no golden test catches because goldens render with empty Attrs.
//
// The marker value is unique so a stray match elsewhere in the markup is
// impossible.
func TestAttrsPropagateToRoot(t *testing.T) {
	t.Parallel()

	const marker = `data-tc-attrs="z9q"`

	attrs := templ.Attributes{"data-tc-attrs": "z9q"}

	tests := []struct {
		name   string
		render func() string
	}{
		{
			name: "badge",
			render: func() string {
				props := display.DefaultBadgeProps()
				props.Text = "x"
				props.Attrs = attrs

				return utils.Render(t, display.Badge(props))
			},
		},
		{
			name: "card",
			render: func() string {
				props := display.DefaultCardProps()
				props.Title = "x"
				props.Attrs = attrs

				return utils.Render(t, display.Card(props))
			},
		},
		{
			name: "button",
			render: func() string {
				props := display.DefaultButtonProps()
				props.Text = "x"
				props.Attrs = attrs

				return utils.Render(t, display.Button(props))
			},
		},
		{
			name: "statcard",
			render: func() string {
				props := display.DefaultStatCardProps()
				props.Label = "x"
				props.Attrs = attrs

				return utils.Render(t, display.StatCard(props))
			},
		},
		{
			name: "emptystate",
			render: func() string {
				props := display.DefaultEmptyStateProps()
				props.Title = "x"
				props.Attrs = attrs

				return utils.Render(t, display.EmptyState(props))
			},
		},
		{
			name: "modal",
			render: func() string {
				props := display.DefaultModalProps()
				props.Attrs = attrs

				return utils.Render(t, display.Modal(props))
			},
		},
		{
			name: "tabs",
			render: func() string {
				props := display.DefaultTabsProps()
				props.Attrs = attrs

				return utils.Render(t, display.Tabs(props))
			},
		},
		{
			name: "dropdown",
			render: func() string {
				props := display.DefaultDropdownProps()
				props.Attrs = attrs

				return utils.Render(t, display.Dropdown(props))
			},
		},
		{
			name: "input",
			render: func() string {
				props := forms.DefaultInputProps()
				props.Name = "x"
				props.Attrs = attrs

				return utils.Render(t, forms.Input(props))
			},
		},
		{
			name: "select",
			render: func() string {
				props := forms.DefaultSelectProps()
				props.Name = "x"
				props.Attrs = attrs

				return utils.Render(t, forms.Select(props))
			},
		},
		{
			name: "alert",
			render: func() string {
				props := feedback.DefaultAlertProps()
				props.Title = "x"
				props.Attrs = attrs

				return utils.Render(t, feedback.Alert(props))
			},
		},
		{
			name: "progressbar",
			render: func() string {
				props := feedback.DefaultProgressBarProps()
				props.Current = 1
				props.Total = 2
				props.Attrs = attrs

				return utils.Render(t, feedback.ProgressBar(props))
			},
		},
		{
			name: "pagination",
			render: func() string {
				props := navigation.DefaultPaginationProps()
				props.CurrentPage = 1
				props.TotalPages = 3
				props.BaseURL = "/"
				props.Attrs = attrs

				return utils.Render(t, navigation.Pagination(props))
			},
		},
		{
			name: "container",
			render: func() string {
				props := layout.DefaultContainerProps()
				props.Attrs = attrs

				return utils.Render(t, layout.Container(props))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			html := tc.render()
			if !strings.Contains(html, marker) {
				t.Errorf(
					"%s dropped the consumer Attrs %s — the Attrs-propagation contract is broken:\n%s",
					tc.name,
					marker,
					html,
				)
			}
		})
	}
}
