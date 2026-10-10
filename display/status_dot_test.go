package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestStatusDotRender(t *testing.T) {
	t.Parallel()

	t.Run("static dot is decorative without label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusDot(StatusDotProps{Tone: StatusToneSuccess}))
		utils.AssertContains(t, output, "bg-green-500")
		utils.AssertContains(t, output, `aria-hidden="true"`)
		utils.AssertNotContains(t, output, "animate-ping")
	})

	t.Run("pulse renders ring with motion-reduce fallback", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusDot(StatusDotProps{Tone: StatusToneSuccess, Pulse: true, Label: "Live"}))
		utils.AssertContains(t, output, "animate-ping")
		utils.AssertContains(t, output, "motion-reduce:animate-none")
		utils.AssertContains(t, output, `role="status"`)
		utils.AssertContains(t, output, "sr-only")
		utils.AssertContains(t, output, "Live")
	})

	t.Run("unknown tone degrades to neutral", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusDot(StatusDotProps{Tone: StatusTone("bogus")}))
		utils.AssertContains(t, output, "bg-gray-400")
		utils.AssertNotContains(t, output, "bg-green-500")
	})

	t.Run("propagates base props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusDot(StatusDotProps{
			BaseProps: utils.BaseProps{
				ID:        "dot",
				Class:     "mt-1",
				AriaLabel: "Conn",
				Attrs:     templ.Attributes{"data-x": "1"},
			},
			Tone: StatusToneDanger,
		}))
		utils.AssertContains(t, output, `id="dot"`)
		utils.AssertContains(t, output, "mt-1")
		utils.AssertContains(t, output, `aria-label="Conn"`)
		utils.AssertContains(t, output, `data-x="1"`)
		utils.AssertContains(t, output, "bg-red-500")
	})
}

func TestLivePillRender(t *testing.T) {
	t.Parallel()

	t.Run("pill renders dot and text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LivePill(DefaultLivePillProps()))
		utils.AssertContains(t, output, "Live")
		utils.AssertContains(t, output, "bg-green-500")
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("warning tone", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LivePill(LivePillProps{Tone: StatusToneWarning, Text: "Reconnecting…"}))
		utils.AssertContains(t, output, "Reconnecting…")
		utils.AssertContains(t, output, "bg-amber-500")
	})

	t.Run("unknown tone degrades to neutral", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LivePill(LivePillProps{Tone: StatusTone("bogus"), Text: "Paused"}))
		utils.AssertContains(t, output, "bg-gray-400")
	})

	t.Run("propagates base props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LivePill(LivePillProps{
			BaseProps: utils.BaseProps{ID: "conn", Class: "mb-2"},
			Tone:      StatusToneSuccess,
			Text:      "Live",
		}))
		utils.AssertContains(t, output, `id="conn"`)
		utils.AssertContains(t, output, "mb-2")
	})
}

func TestStatusToneBg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tone StatusTone
		want string
	}{
		{"neutral", StatusToneNeutral, "bg-gray-400 dark:bg-gray-500"},
		{"success", StatusToneSuccess, "bg-green-500 dark:bg-green-400"},
		{"warning", StatusToneWarning, "bg-amber-500 dark:bg-amber-400"},
		{"danger", StatusToneDanger, "bg-red-500 dark:bg-red-400"},
		{"info", StatusToneInfo, "bg-blue-500 dark:bg-blue-400"},
		{"bogus falls back neutral", StatusTone("bogus"), "bg-gray-400 dark:bg-gray-500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := statusToneBg(tt.tone); got != tt.want {
				t.Errorf("statusToneBg(%q) = %q, want %q", tt.tone, got, tt.want)
			}
		})
	}
}

func TestStatusToneIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tone StatusTone
		want bool
	}{
		{"neutral", StatusToneNeutral, true},
		{"success", StatusToneSuccess, true},
		{"warning", StatusToneWarning, true},
		{"danger", StatusToneDanger, true},
		{"info", StatusToneInfo, true},
		{"bogus", StatusTone("bogus"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := StatusToneIsValid(tt.tone); got != tt.want {
				t.Errorf("StatusToneIsValid(%q) = %v, want %v", tt.tone, got, tt.want)
			}
		})
	}
}
