package tone_test

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

var docUpdate = flag.Bool("update", false, "rewrite the generated tone sections in the docs")

const (
	themeBridgeDoc = "../../docs/recipes/scoped-theme-bridge.md"
	themeLayerCSS  = "../../templates/templ-components-theme.css"

	paletteBegin = "<!-- BEGIN generated: tone palette table (utils/tone/doc_sync_test.go; go test ./utils/tone/ -update) -->"
	paletteEnd   = "<!-- END generated: tone palette table -->"
)

// tonePaletteTable renders the "Which colors does the library use?" table
// from the tone package's own tables — the single source of truth for the
// library palette (report H4: this was the fifth hand-maintained copy).
func tonePaletteTable() string {
	chartSeries := map[tone.Hue]int{}
	for i, hue := range tone.ChartHues() {
		chartSeries[hue] = i + 1
	}

	var b strings.Builder
	b.WriteString("| Hue | Tones rendered in it | Chart palette |\n")
	b.WriteString("| --- | --- | --- |\n")

	for _, hue := range tone.Hues() {
		prefix := string(hue) + "-"

		var tones []string

		for _, t := range tone.Tones() {
			if toneUsesHue(t, prefix) {
				tones = append(tones, string(t))
			}
		}

		tonesCell := "—"
		if len(tones) > 0 {
			tonesCell = strings.Join(tones, ", ")
		}

		chartCell := "—"
		if series, ok := chartSeries[hue]; ok {
			chartCell = fmt.Sprintf("yes (series %d)", series)
		}

		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", hue, tonesCell, chartCell)
	}

	return b.String()
}

// toneUsesHue reports whether any surface of the tone references the hue
// (classes carry the hue as a `hue-shade` token prefix).
func toneUsesHue(t tone.Tone, prefix string) bool {
	surfaces := tone.For(t)

	for _, pair := range []tone.Pair{
		surfaces.Dot, surfaces.Text, surfaces.TextSoft, surfaces.Border,
		surfaces.BGTint, surfaces.BGSolid, surfaces.Icon,
	} {
		if strings.Contains(pair.Light, prefix) || strings.Contains(pair.Dark, prefix) {
			return true
		}
	}

	return false
}

// TestThemeBridgePaletteTableSync keeps the generated palette table in
// docs/recipes/scoped-theme-bridge.md byte-identical to what the tone
// package derives. Any tone/hue table change fails here until the doc is
// regenerated (`go test ./utils/tone/ -run TestThemeBridgePaletteTableSync -update`).
func TestThemeBridgePaletteTableSync(t *testing.T) {
	t.Parallel()

	doc, err := os.ReadFile(themeBridgeDoc)
	if err != nil {
		t.Fatalf("read theme bridge doc: %v", err)
	}

	want := paletteBegin + "\n" + tonePaletteTable() + paletteEnd

	if *docUpdate {
		updated := replaceGeneratedSection(string(doc), paletteBegin, paletteEnd, want)
		if updated != string(doc) {
			if err := os.WriteFile(themeBridgeDoc, []byte(updated), 0o644); err != nil {
				t.Fatalf("rewrite theme bridge doc: %v", err)
			}
		}

		return
	}

	start := strings.Index(string(doc), paletteBegin)

	end := strings.Index(string(doc), paletteEnd)
	if start < 0 || end < start {
		t.Fatalf(
			"theme bridge doc lacks the generated palette markers %q/%q — run go test ./utils/tone/ -update",
			paletteBegin, paletteEnd,
		)
	}

	if have := string(doc)[start : end+len(paletteEnd)]; have != want {
		t.Errorf(
			"theme bridge doc palette table is stale:\n--- have ---\n%s\n--- want ---\n%s\nRun: go test ./utils/tone/ -update",
			have,
			want,
		)
	}
}

// TestThemeTokenLayerInventory pins the tc-* token inventory of
// templates/templ-components-theme.css (ADR-0008): exactly the four semantic
// ramps (primary/danger/success/warning — the accent, danger, success, and
// warning tones) at every shade 50–950, so a ramp that loses a shade or a
// stray token that appears fails loudly instead of silently re-breaking
// consumer overrides.
func TestThemeTokenLayerInventory(t *testing.T) {
	t.Parallel()

	css, err := os.ReadFile(themeLayerCSS)
	if err != nil {
		t.Fatalf("read theme layer css: %v", err)
	}

	tokenRe := regexp.MustCompile(`--color-tc-([a-z]+)-(\d+):`)
	shades := []string{"50", "100", "200", "300", "400", "500", "600", "700", "800", "900", "950"}
	want := map[string]bool{}

	for _, ramp := range []string{"primary", "danger", "success", "warning"} {
		for _, shade := range shades {
			want[fmt.Sprintf("%s-%s", ramp, shade)] = true
		}
	}

	got := map[string]bool{}
	for _, match := range tokenRe.FindAllStringSubmatch(string(css), -1) {
		got[match[1]+"-"+match[2]] = true
	}

	for token := range want {
		if !got[token] {
			t.Errorf("theme layer css is missing token --color-tc-%s (ramp incomplete)", token)
		}
	}

	for token := range got {
		if !want[token] {
			t.Errorf("theme layer css carries unexpected ramp token --color-tc-%s (inventory drift)", token)
		}
	}
}

func replaceGeneratedSection(doc, begin, end, replacement string) string {
	start := strings.Index(doc, begin)

	finish := strings.Index(doc, end)
	if start < 0 || finish < start {
		return doc + "\n" + replacement + "\n"
	}

	return doc[:start] + replacement + doc[finish+len(end):]
}
