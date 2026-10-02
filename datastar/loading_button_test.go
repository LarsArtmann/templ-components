package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestLoadingButtonSignalExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		signal    string
		wantBusy  string
		wantRest  string
	}{
		{"normal signal", "saving", "$saving", "!$saving"},
		{"whitespace signal degrades", "  ", "false", "true"},
		{"empty signal degrades", "", "false", "true"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := loadingButtonBusyExpr(tc.signal); got != tc.wantBusy {
				t.Errorf("loadingButtonBusyExpr(%q) = %q, want %q", tc.signal, got, tc.wantBusy)
			}

			if got := loadingButtonRestExpr(tc.signal); got != tc.wantRest {
				t.Errorf("loadingButtonRestExpr(%q) = %q, want %q", tc.signal, got, tc.wantRest)
			}
		})
	}
}

func TestLoadingButtonMarkup(t *testing.T) {
	t.Parallel()

	props := LoadingButtonProps{
		BaseProps:   utils.BaseProps{ID: "save-label"},
		Signal:      "saving",
		DefaultText: "Save",
		LoadingText: "Saving…",
	}

	output := utils.Render(t, LoadingButton(props))

	utils.AssertContains(t, output, `id="save-label"`)
	utils.AssertContains(t, output, `data-show="$saving"`)
	utils.AssertContains(t, output, `data-show="!$saving"`)
	utils.AssertContains(t, output, "Save")
	utils.AssertContains(t, output, "Saving…")
	utils.AssertContains(t, output, "tc-btn-loading")
}

func TestLoadingButtonEmptySignalInert(t *testing.T) {
	t.Parallel()

	props := LoadingButtonProps{DefaultText: "Save"}
	output := utils.Render(t, LoadingButton(props))

	// Empty signal: rest state always true, busy state always false — the
	// bare "$" object is truthy and would pin the spinner visible.
	utils.AssertContains(t, output, `data-show="false"`)
	utils.AssertContains(t, output, `data-show="true"`)
}

func TestLoadingButtonDefaults(t *testing.T) {
	t.Parallel()

	props := DefaultLoadingButtonProps()
	if props.Signal != "" || props.DefaultText != "" {
		t.Errorf("defaults must be zero-valued; got %+v", props)
	}
}

func TestLoadingButtonSpinnerVariants(t *testing.T) {
	t.Parallel()

	t.Run("default spinner renders", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LoadingButton(LoadingButtonProps{Signal: "s", DefaultText: "Go"}))
		utils.AssertContains(t, output, "animate-spin")
		utils.AssertContains(t, output, "motion-reduce:animate-none")
	})
}
