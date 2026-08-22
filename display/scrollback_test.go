package display

import (
	"os"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func scrollbackFixtureLines() []ScrollbackLine {
	return []ScrollbackLine{
		{Timestamp: "12:47:03.184", Tag: "query", Text: "ads.example.com A", Tone: ScrollbackToneInfo},
		{Timestamp: "12:47:03.184", Tag: "match", Text: "blocklist: StevenBlack/hosts", Tone: ScrollbackToneDanger},
		{Timestamp: "12:47:03.185", Tag: "action", Text: "NXDOMAIN", Tone: ScrollbackToneWarning},
		{Timestamp: "12:47:03.185", Tag: "ttl", Text: "86400s", Tone: ScrollbackToneSuccess},
	}
}

func TestScrollbackGoldenSweep(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "scrollback_stagger", HTML: utils.Render(t, Scrollback(ScrollbackProps{
			Stagger: true,
			Lines:   scrollbackFixtureLines(),
		}))},
		{Name: "scrollback_instant", HTML: utils.Render(t, Scrollback(ScrollbackProps{
			Lines: scrollbackFixtureLines(),
		}))},
		{Name: "scrollback_labeled", HTML: utils.Render(t, Scrollback(ScrollbackProps{
			Lines: scrollbackFixtureLines(),
			BaseProps: utils.BaseProps{
				AriaLabel: "DNS resolution trace for ads.example.com",
			},
		}))},
		{Name: "scrollback_empty", HTML: utils.Render(t, Scrollback(DefaultScrollbackProps()))},
	})
}

func TestScrollbackBehavior(t *testing.T) {
	t.Parallel()

	t.Run("user sees every log line with timestamp and tag", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{Lines: scrollbackFixtureLines()}))
		utils.AssertContainsAll(t, output,
			"12:47:03.184", "query", "ads.example.com A",
			"match", "blocklist: StevenBlack/hosts",
			"action", "NXDOMAIN", "ttl", "86400s",
		)
	})

	t.Run("lines render in a monospace terminal voice", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{Lines: scrollbackFixtureLines()}))
		utils.AssertContains(t, output, "font-mono")
	})

	t.Run("stagger opt-in emits the animated line class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{
			Stagger: true,
			Lines:   scrollbackFixtureLines(),
		}))
		utils.AssertContains(t, output, "tc-log-line")
		utils.AssertNotContains(t, output, "tc-log-line-still")
	})

	t.Run("zero value renders instantly via the still modifier", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{Lines: scrollbackFixtureLines()}))
		utils.AssertContains(t, output, "tc-log-line-still")
	})

	t.Run("tone colors the tag column per line", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{Lines: scrollbackFixtureLines()}))
		utils.AssertContainsAll(t, output,
			"text-blue-600 dark:text-blue-400",
			"text-red-600 dark:text-red-400",
			"text-amber-600 dark:text-amber-400",
			"text-green-600 dark:text-green-400",
		)
	})
}

func TestScrollbackA11y(t *testing.T) {
	t.Parallel()

	t.Run("decorative by default: aria-hidden hides the trace from screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{Lines: scrollbackFixtureLines()}))
		utils.AssertContains(t, output, `aria-hidden="true"`)
		utils.AssertNotContains(t, output, "aria-label=")
	})

	t.Run("aria-label exposes real log content to screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{
			Lines: scrollbackFixtureLines(),
			BaseProps: utils.BaseProps{
				AriaLabel: "DNS resolution trace",
			},
		}))
		utils.AssertContains(t, output, `aria-label="DNS resolution trace"`)
		utils.AssertNotContains(t, output, `aria-hidden="true"`)
	})

	t.Run("reduced-motion users see all lines immediately (CSS guard)", func(t *testing.T) {
		t.Parallel()
		css, err := os.ReadFile("../templates/custom.css")
		if err != nil {
			t.Fatalf("read custom.css: %v", err)
		}
		styles := string(css)
		for _, want := range []string{
			".tc-log-line",
			".tc-log-line.tc-log-line-still",
			"prefers-reduced-motion: reduce",
			"@keyframes tc-log-in",
		} {
			if !strings.Contains(styles, want) {
				t.Errorf("custom.css missing %q — the scrollback animation depends on it", want)
			}
		}
	})
}

func TestScrollbackEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty lines render nothing", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(DefaultScrollbackProps()))
		if output != "" {
			t.Errorf("expected empty output, got %q", output)
		}
	})

	t.Run("unknown tone falls back to neutral", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{
			Lines: []ScrollbackLine{{Tag: "x", Text: "y", Tone: ScrollbackTone("bogus")}},
		}))
		utils.AssertContains(t, output, "text-gray-500 dark:text-gray-400")
	})

	t.Run("lines without timestamp or tag omit the empty columns", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{
			Lines: []ScrollbackLine{{Text: "bare line"}},
		}))
		if strings.Contains(output, "min-w-12") {
			t.Error("expected no tag column for a line without a tag")
		}
	})

	t.Run("more than 8 lines still render (stagger cap is CSS-only)", func(t *testing.T) {
		t.Parallel()
		lines := make([]ScrollbackLine, 12)
		for i := range lines {
			lines[i] = ScrollbackLine{Tag: "l", Text: strings.Repeat("x", i+1)}
		}
		output := utils.Render(t, Scrollback(ScrollbackProps{Stagger: true, Lines: lines}))
		if got := strings.Count(output, "tc-log-line"); got != 12 {
			t.Errorf("expected 12 line divs, got %d", got)
		}
	})

	t.Run("no physical RTL properties are emitted", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Scrollback(ScrollbackProps{Lines: scrollbackFixtureLines()}))
		for _, physical := range []string{" ml-", " mr-", " pl-", " pr-", "text-left", "text-right"} {
			if strings.Contains(output, physical) {
				t.Errorf("physical property %q emitted — use logical (ms-/me-/text-start)", physical)
			}
		}
	})
}

func TestScrollbackToneLookup(t *testing.T) {
	t.Parallel()

	for tone, classes := range scrollbackToneLookup {
		if !strings.Contains(classes, "dark:") {
			t.Errorf("tone %s lacks a dark: variant: %q", tone, classes)
		}
	}

	if _, ok := scrollbackToneLookup[ScrollbackTone("bogus")]; ok {
		t.Error("unknown tone must not resolve")
	}
}
