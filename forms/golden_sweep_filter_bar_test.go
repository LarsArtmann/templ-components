package forms

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
)

// renderFilterBarWithChildren renders a FilterBar with child components.
func renderFilterBarWithChildren(t *testing.T, props FilterBarProps, children ...templ.Component) string {
	t.Helper()

	var buf bytes.Buffer

	bar := FilterBar(props)

	ctx := templ.WithChildren(context.Background(), templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, child := range children {
			if err := child.Render(ctx, w); err != nil {
				return fmt.Errorf("render child: %w", err)
			}
		}

		return nil
	}))

	if err := bar.Render(ctx, &buf); err != nil {
		t.Fatalf("failed to render filter bar with children: %v", err)
	}

	return strings.TrimSpace(buf.String())
}

// Golden sweep for FilterBar (TODO_LIST #390): the DiscordSync filterForm
// shape productized, both transports.

func TestGoldenSweepFilterBar(t *testing.T) {
	t.Parallel()

	statusDropdown := FilterDropdown(FilterDropdownProps{
		Name:  "status",
		Label: "Status",
		Options: []SelectOption{
			{Value: "", Label: "All"},
			{Value: "active", Label: "Active"},
		},
	})

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "filter_bar_default_htmx", HTML: renderFilterBarWithChildren(t, FilterBarProps{
			Action: "/users",
			Target: "#results",
		}, statusDropdown)},
		{Name: "filter_bar_sticky", HTML: renderFilterBarWithChildren(t, FilterBarProps{
			Action: "/users",
			Sticky: true,
		}, statusDropdown)},
		{Name: "filter_bar_reset", HTML: renderFilterBarWithChildren(t, FilterBarProps{
			Action:    "/users",
			ResetText: "Clear filters",
		}, statusDropdown)},
		{Name: "filter_bar_no_push_url", HTML: utils.Render(t, FilterBar(FilterBarProps{
			Action:    "/users",
			NoPushURL: true,
		}))},
		{Name: "filter_bar_inert", HTML: utils.Render(t, FilterBar(FilterBarProps{}))},
		{Name: "filter_bar_wire_htmx", HTML: utils.Render(t, FilterBar(FilterBarProps{
			Wire: &wire.Action{Method: wire.MethodGet, URL: "/api/users", Target: "#results"},
		}))},
		{Name: "filter_bar_wire_datastar", HTML: utils.Render(t, FilterBar(FilterBarProps{
			Wire: &wire.Action{Transport: wire.TransportDatastar, Method: wire.MethodGet, URL: "/api/users"},
		}))},
	})
}
