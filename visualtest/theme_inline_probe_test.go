package visualtest

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestThemeInlineProbe (TODO #432) proves the mode-flipping token strategy —
// ONE utility class per role (`bg-tc-solid`) whose light/dark pair lives in a
// CSS variable overridden under `.dark`, wired via `@theme inline` — renders
// computed colors IDENTICAL to today's hand-maintained `dark:` pair classes
// in a real browser, in both modes, including the opacity-modifier
// (color-mix) path. This is the empirical basis for retiring the ~669
// `dark:` pairs in the v2 class swap: computed-style equality is pixel
// equality for these properties.
//
// Fixtures: testdata/theme-inline-probe/{page.html,compiled.css} — the CSS
// was compiled with the repo-pinned tailwindcss 4.3.3; regeneration and the
// full verdict live in docs/planning/2026-10-16_theme-inline-probe.md.
func TestThemeInlineProbe(t *testing.T) {
	css, err := os.ReadFile(filepath.Join("testdata", "theme-inline-probe", "compiled.css"))
	if err != nil {
		t.Fatalf("theme-inline probe: read compiled.css: %v", err)
	}
	pageTmpl, err := os.ReadFile(filepath.Join("testdata", "theme-inline-probe", "page.html"))
	if err != nil {
		t.Fatalf("theme-inline probe: read page.html: %v", err)
	}

	t.Run("compiled css wires utilities to the mode-flipped vars", func(t *testing.T) {
		for _, needle := range []string{
			"background-color: var(--tc-solid)",
			"color: var(--tc-ink)",
			"color-mix(in oklab, var(--tc-solid) 10%, transparent)",
		} {
			if !strings.Contains(string(css), needle) {
				t.Errorf("compiled.css lacks %q — @theme inline wiring changed; recompile and re-verify the probe", needle)
			}
		}
	})

	page := strings.Replace(string(pageTmpl), "/*__PROBE_CSS__*/", string(css), 1)
	dataURL := "data:text/html;charset=utf-8;base64," + base64.StdEncoding.EncodeToString([]byte(page))

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body"),
		pollTrue(`!!document.getElementById('tok-solid')`),
	); err != nil {
		t.Fatalf("theme-inline probe: navigate: %v", err)
	}

	computed := func(selector, prop string) string {
		t.Helper()

		var value string
		if err := chromedp.Do(ctx, evalExprInto(
			`getComputedStyle(document.getElementById('`+selector+`')).`+prop, &value)); err != nil {
			t.Fatalf("theme-inline probe: computed %s of #%s: %v", prop, selector, err)
		}

		return value
	}

	pairs := []struct {
		name string
		lit  string
		tok  string
		prop string
	}{
		{"solid background", "lit-solid", "tok-solid", "backgroundColor"},
		{"ink text", "lit-ink", "tok-ink", "color"},
		{"solid background at 10% opacity", "lit-solid-op", "tok-solid-op", "backgroundColor"},
	}

	light := map[string]string{}
	for _, pair := range pairs {
		light[pair.name] = computed(pair.lit, pair.prop)
		if light[pair.name] == "" {
			t.Fatalf("theme-inline probe: light-mode literal #%s computed empty — page broken", pair.lit)
		}
		if tok := computed(pair.tok, pair.prop); tok != light[pair.name] {
			t.Errorf(
				"theme-inline probe: light mode %s: literal %q vs token %q — token strategy is NOT pixel-identical",
				pair.name, light[pair.name], tok,
			)
		}
	}

	if err := chromedp.Do(ctx, evalExprInto(
		`document.documentElement.classList.add('dark'); 'ok'`, new(string))); err != nil {
		t.Fatalf("theme-inline probe: enter dark mode: %v", err)
	}

	for _, pair := range pairs {
		lit, tok := computed(pair.lit, pair.prop), computed(pair.tok, pair.prop)
		if tok != lit {
			t.Errorf(
				"theme-inline probe: dark mode %s: literal %q vs token %q — token strategy is NOT pixel-identical",
				pair.name, lit, tok,
			)
		}
		if lit == light[pair.name] {
			t.Errorf(
				"theme-inline probe: dark mode %s did not flip (still %q) — the pair classes on the page are not applying",
				pair.name, lit,
			)
		}
	}
}
