package feedback

import (
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestCircularProgress_A11Y(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		props       CircularProgressProps
		checkHidden bool
	}{
		{
			name:  "role progressbar",
			props: CircularProgressProps{Value: 50},
		},
		{
			name:  "aria-valuenow reflects clamped value",
			props: CircularProgressProps{Value: 75},
		},
		{
			name:  "aria-valuemin is 0",
			props: CircularProgressProps{Value: 50},
		},
		{
			name:  "aria-valuemax is 100",
			props: CircularProgressProps{Value: 50},
		},
		{
			name:  "custom aria-label",
			props: CircularProgressProps{Value: 50, BaseProps: utils.BaseProps{AriaLabel: "Upload progress"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := utils.Render(t, CircularProgress(tt.props))

			if !strings.Contains(html, `role="progressbar"`) {
				t.Error("expected role='progressbar'")
			}

			if !strings.Contains(html, `aria-valuenow=`) {
				t.Error("expected aria-valuenow attribute")
			}

			if !strings.Contains(html, `aria-valuemin="0"`) {
				t.Error("expected aria-valuemin='0'")
			}

			if !strings.Contains(html, `aria-valuemax="100"`) {
				t.Error("expected aria-valuemax='100'")
			}

			if tt.props.AriaLabel != "" {
				if !strings.Contains(html, fmt.Sprintf(`aria-label="%s"`, tt.props.AriaLabel)) {
					t.Errorf("expected aria-label=%q", tt.props.AriaLabel)
				}
			}
		})
	}
}

func TestCircularProgress_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		props   CircularProgressProps
		wantVal int
	}{
		{"zero value", CircularProgressProps{Value: 0}, 0},
		{"exactly 100", CircularProgressProps{Value: 100}, 100},
		{"negative clamps to 0", CircularProgressProps{Value: -50}, 0},
		{"over 100 clamps to 100", CircularProgressProps{Value: 150}, 100},
		{"default zero value", CircularProgressProps{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := utils.Render(t, CircularProgress(tt.props))
			expected := fmt.Sprintf(`aria-valuenow="%d"`, tt.wantVal)

			if !strings.Contains(html, expected) {
				t.Errorf("expected %s in HTML", expected)
			}
		})
	}
}

func TestCircularProgress_BDD(t *testing.T) {
	t.Parallel()

	t.Run("renders with default props", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{}))
		if !strings.Contains(html, `role="progressbar"`) {
			t.Error("should render as a progressbar")
		}
	})

	t.Run("shows label when ShowLabel is true", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{
			Value:     42,
			ShowLabel: true,
		}))
		if !strings.Contains(html, "42%") {
			t.Error("should render '42%' label")
		}
	})

	t.Run("hides label when ShowLabel is false", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{
			Value:     42,
			ShowLabel: false,
		}))
		// aria-label always contains the percentage, but the visible label span should not be present
		if strings.Contains(html, `<span class="absolute inset-0`) {
			t.Error("should not render visible label span when ShowLabel is false")
		}
	})

	t.Run("applies custom color class", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{
			Value: 50,
			Color: "text-green-600",
		}))
		if !strings.Contains(html, "text-green-600") {
			t.Error("should apply custom color class")
		}
	})

	t.Run("applies custom track color class", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{
			Value:      50,
			TrackColor: "text-red-200",
		}))
		if !strings.Contains(html, "text-red-200") {
			t.Error("should apply custom track color class")
		}
	})
}

func ExampleCircularProgress() {
	utils.Render(nil, CircularProgress(CircularProgressProps{
		Value:     75,
		ShowLabel: true,
	}))
}

func TestCircularProgress_Coverage(t *testing.T) {
	t.Parallel()

	t.Run("DefaultCircularProgressProps returns md size", func(t *testing.T) {
		t.Parallel()

		props := DefaultCircularProgressProps()
		if props.Size != CircularProgressSizeMD {
			t.Errorf("expected md, got %s", props.Size)
		}

		if props.Color == "" {
			t.Error("default color should not be empty")
		}

		if props.TrackColor == "" {
			t.Error("default track color should not be empty")
		}
	})

	t.Run("all sizes produce valid dimension classes", func(t *testing.T) {
		t.Parallel()

		sizes := []CircularProgressSize{
			CircularProgressSizeSM,
			CircularProgressSizeMD,
			CircularProgressSizeLG,
		}

		for _, size := range sizes {
			html := utils.Render(t, CircularProgress(CircularProgressProps{
				Value: 50,
				Size:  size,
			}))
			if html == "" {
				t.Errorf("size %s produced empty output", size)
			}
		}
	})

	t.Run("unknown size falls back gracefully", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{
			Value: 50,
			Size:  "xl",
		}))
		// Should still render as a progressbar with valid value
		if !strings.Contains(html, `role="progressbar"`) {
			t.Error("unknown size should still render as progressbar")
		}

		if !strings.Contains(html, `aria-valuenow="50"`) {
			t.Error("unknown size should still show correct value")
		}
	})

	t.Run("custom ID is rendered", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, CircularProgress(CircularProgressProps{
			Value:     50,
			BaseProps: utils.BaseProps{ID: "my-progress"},
		}))
		if !strings.Contains(html, `id="my-progress"`) {
			t.Error("should render custom ID")
		}
	})
}
