package display

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestDataStateRender(t *testing.T) {
	t.Parallel()

	t.Run("content renders children", func(t *testing.T) {
		t.Parallel()
		output := renderDataStateWithChildren(t, DataStateProps{}, plainChild("<p>table goes here</p>"))
		utils.AssertContains(t, output, "table goes here")
		utils.AssertNotContains(t, output, "inbox")
	})

	t.Run("empty renders placeholder with inbox icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataState(DataStateProps{
			State:       DataStateEmpty,
			Title:       "No tracked device activity",
			Description: "Nothing landed in the last 30 days.",
		}))
		utils.AssertContains(t, output, "No tracked device activity")
		utils.AssertContains(t, output, "Nothing landed in the last 30 days.")
	})

	t.Run("unavailable renders placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataState(DataStateProps{
			State:       DataStateUnavailable,
			Title:       "Unavailable — retrying",
			Description: "The query failed; this region keeps retrying every 30s.",
		}))
		utils.AssertContains(t, output, "Unavailable — retrying")
	})

	t.Run("disabled renders placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataState(DataStateProps{
			State: DataStateDisabled,
			Title: "Tracking is disabled",
		}))
		utils.AssertContains(t, output, "Tracking is disabled")
	})

	t.Run("unknown state falls back to content", func(t *testing.T) {
		t.Parallel()
		output := renderDataStateWithChildren(
			t,
			DataStateProps{State: DataStateState("bogus")},
			plainChild("<p>still content</p>"),
		)
		utils.AssertContains(t, output, "still content")
	})

	t.Run("explicit icon overrides state default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataState(DataStateProps{
			State: DataStateEmpty,
			Title: "Custom",
			Icon:  icons.Question,
		}))
		utils.AssertContains(t, output, "Custom")
	})

	t.Run("propagates base props on both rungs", func(t *testing.T) {
		t.Parallel()
		content := renderDataStateWithChildren(t, DataStateProps{
			ID: "region", Class: "mb-4", AriaLabel: "Devices",
		}, plainChild("<p>x</p>"))
		utils.AssertContains(t, content, `id="region"`)
		utils.AssertContains(t, content, "mb-4")
		utils.AssertContains(t, content, `aria-label="Devices"`)

		empty := utils.Render(t, DataState(DataStateProps{
			ID: "region2", Class: "mb-6",
			State: DataStateEmpty,
			Title: "Empty",
		}))
		utils.AssertContains(t, empty, `id="region2"`)
		utils.AssertContains(t, empty, "mb-6")
	})
}

func TestDataStateIcon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props DataStateProps
		want  icons.Name
	}{
		{"empty default", DataStateProps{State: DataStateEmpty}, icons.Inbox},
		{"unavailable default", DataStateProps{State: DataStateUnavailable}, icons.ExclamationTriangle},
		{"disabled default", DataStateProps{State: DataStateDisabled}, icons.NoSymbol},
		{"explicit wins", DataStateProps{State: DataStateEmpty, Icon: icons.Question}, icons.Question},
		{"content rung has inbox fallback", DataStateProps{}, icons.Inbox},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := dataStateIcon(tt.props); got != tt.want {
				t.Errorf("dataStateIcon(%+v) = %q, want %q", tt.props.State, got, tt.want)
			}
		})
	}
}

func TestDataStateStateIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    DataStateState
		want bool
	}{
		{"content", DataStateContent, true},
		{"empty", DataStateEmpty, true},
		{"unavailable", DataStateUnavailable, true},
		{"disabled", DataStateDisabled, true},
		{"bogus", DataStateState("bogus"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DataStateStateIsValid(tt.s); got != tt.want {
				t.Errorf("DataStateStateIsValid(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

// plainChild renders raw HTML — a test stand-in for real children.
func plainChild(html string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, html)

		return err
	})
}

// renderDataStateWithChildren renders a DataState content rung with children.
func renderDataStateWithChildren(t *testing.T, props DataStateProps, children ...templ.Component) string {
	t.Helper()

	var buf bytes.Buffer

	ctx := templ.WithChildren(context.Background(), children[0])
	if err := DataState(props).Render(ctx, &buf); err != nil {
		t.Fatalf("render data state: %v", err)
	}

	return strings.TrimSpace(buf.String())
}
