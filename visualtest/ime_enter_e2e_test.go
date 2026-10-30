package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
)

// imeEnterE2EPage renders the three Enter-commits components — Combobox,
// TagsInput, CommandPalette — on one page so the IME guard can be probed on
// each in a single browser session.
func imeEnterE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		combo := forms.ComboboxProps{
			Nonce: "e2e-nonce", ID: "e2e-ime-combo", Name: "combo",
			Label: "Framework",
			Options: []forms.ComboboxOption{
				{Value: "gs", Label: "Go Server"},
				{Value: "ts", Label: "TS Stack"},
			},
		}
		if err := forms.Combobox(combo).Render(ctx, w); err != nil {
			return err
		}

		tags := forms.TagsInputProps{
			Nonce: "e2e-nonce", ID: "e2e-ime-tags", Name: "tags",
			Label: "Tags", Placeholder: "Add tag",
		}
		if err := forms.TagsInput(tags).Render(ctx, w); err != nil {
			return err
		}

		palette := display.CommandPaletteProps{
			Nonce:        "e2e-nonce",
			ID:           "e2e-ime-palette",
			TriggerLabel: "Commands",
			Groups: []display.CommandGroup{
				{Label: "Actions", Items: []display.CommandItem{{Label: "Deploy"}}},
			},
		}

		return display.CommandPalette(palette).Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func imeEnterE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = imeEnterE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// dispatchKeyJS builds a keydown dispatch expression. The 229 variant patches
// keyCode onto the event because KeyboardEventInit cannot set it — the legacy
// signal some IME engines report instead of isComposing.
func dispatchKeyJS(target, key string, composing, legacy229 bool) string {
	init := "{key:'" + key + "',bubbles:true,cancelable:true}"
	if composing {
		init = "{key:'" + key + "',bubbles:true,cancelable:true,isComposing:true}"
	}

	expr := "(function(){" +
		"var e=new KeyboardEvent('keydown'," + init + ");"

	if legacy229 {
		expr += "Object.defineProperty(e,'keyCode',{get:function(){return 229;}});"
	}

	expr += "document.querySelector('" + target + "').dispatchEvent(e);})()"

	return expr
}

// TestIMEE2ECompositionSafeEnter proves in a real browser that a composition
// Enter — reported as isComposing, or as the legacy keyCode 229 — commits the
// IME string instead of triggering the component action, while a plain Enter
// still triggers it. Covers the three Enter-commits components: Combobox
// (pick option), TagsInput (add tag), CommandPalette (activate result).
// Serial — shared browser allocator.
func TestIMEE2ECompositionSafeEnter(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := imeEnterE2EServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[ime]: navigate: %v", err)
	}

	t.Run("combobox", func(t *testing.T) {
		input := `[data-combobox-input="e2e-ime-combo"]`

		// ArrowDown opens the list and highlights the first option.
		if err := chromedp.Do(ctx,
			evalVoid(`document.querySelector('`+input+`').focus();`),
			evalVoid(dispatchKeyJS(input, "ArrowDown", false, false)),
			pollTrue(
				`document.querySelector('[data-combobox-list="e2e-ime-combo"] [data-selected="true"]')!==null`,
				pollFast()...,
			),
		); err != nil {
			t.Fatalf("visualtest[ime] combobox: open + highlight: %v", err)
		}

		// Composition Enter (isComposing, then legacy 229) must not pick.
		if err := chromedp.Do(ctx,
			evalVoid(dispatchKeyJS(input, "Enter", true, false)),
			evalVoid(dispatchKeyJS(input, "Enter", false, true)),
			pollTrue(
				`document.querySelector('[data-combobox-value="e2e-ime-combo"]').value===''`,
				pollFast()...,
			),
			pollTrue(
				`!document.querySelector('[data-combobox-list="e2e-ime-combo"]').classList.contains('hidden')`,
				pollFast()...,
			),
		); err != nil {
			t.Fatalf("visualtest[ime] combobox: composition Enter picked an option: %v", err)
		}

		// Plain Enter picks the highlighted option and closes the list.
		if err := chromedp.Do(ctx,
			evalVoid(dispatchKeyJS(input, "Enter", false, false)),
			pollTrue(
				`document.querySelector('[data-combobox-value="e2e-ime-combo"]').value==='gs'`,
				pollFast()...,
			),
			pollTrue(
				`document.querySelector('[data-combobox-list="e2e-ime-combo"]').classList.contains('hidden')`,
				pollFast()...,
			),
		); err != nil {
			t.Fatalf("visualtest[ime] combobox: plain Enter did not pick: %v", err)
		}
	})

	t.Run("tags-input", func(t *testing.T) {
		input := `[data-tc-tags-input]`

		// Type a value, then send both composition variants — no tag may be
		// added, and the draft text must survive (Enter belongs to the IME).
		if err := chromedp.Do(ctx,
			evalVoid(`document.querySelector('`+input+`').focus();`),
			chromedp.SetValue(input, "go"),
			evalVoid(dispatchKeyJS(input, "Enter", true, false)),
			evalVoid(dispatchKeyJS(input, "Enter", false, true)),
			pollTrue(`document.querySelectorAll('[data-tc-tag]').length===0`, pollFast()...),
		); err != nil {
			t.Fatalf("visualtest[ime] tags-input: composition Enter added a tag: %v", err)
		}

		var draft string

		if err := chromedp.Do(ctx, evalExprInto(
			`document.querySelector('`+input+`').value`, &draft,
		)); err != nil {
			t.Fatalf("visualtest[ime] tags-input: read draft: %v", err)
		}

		if draft != "go" {
			t.Fatalf("visualtest[ime] tags-input: draft %q consumed by composition Enter, want %q", draft, "go")
		}

		if err := chromedp.Do(ctx,
			evalVoid(dispatchKeyJS(input, "Enter", false, false)),
			pollTrue(`document.querySelectorAll('[data-tc-tag]').length===1`, pollFast()...),
		); err != nil {
			t.Fatalf("visualtest[ime] tags-input: plain Enter did not add tag: %v", err)
		}
	})

	t.Run("command-palette", func(t *testing.T) {
		t.Logf("palette: waiting for trigger visible")

		if err := chromedp.Do(ctx,
			chromedp.WaitVisible(chromedp.CSS("[data-tc-palette-open]")),
		); err != nil {
			t.Fatalf("visualtest[ime] command-palette: trigger never visible: %v", err)
		}

		t.Logf("palette: trigger visible, clicking")

		if err := chromedp.Do(ctx,
			chromedp.Click(chromedp.CSS("[data-tc-palette-open]")),
		); err != nil {
			t.Fatalf("visualtest[ime] command-palette: click: %v", err)
		}

		t.Logf("palette: clicked, polling open")

		if err := chromedp.Do(ctx,
			pollTrue(`document.getElementById('e2e-ime-palette').open===true`, pollFast()...),
		); err != nil {
			t.Fatalf("visualtest[ime] command-palette: did not open: %v", err)
		}

		t.Logf("palette: open, dispatching composition Enters")

		if err := chromedp.Do(ctx,
			evalVoid(dispatchKeyJS(`[data-tc-palette-input]`, "Enter", false, false)),
			pollTrue(`document.getElementById('e2e-ime-palette').open===false`, pollFast()...),
		); err != nil {
			t.Fatalf("visualtest[ime] command-palette: plain Enter did not activate: %v", err)
		}
	})
}
