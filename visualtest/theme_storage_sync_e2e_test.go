package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/layout"
)

// themeStorageSyncE2EPage renders Base (ThemeScript included, so first paint
// matches the emulated OS preference) plus a nonce'd ThemeToggle.
func themeStorageSyncE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.ThemeToggle("Toggle theme", "e2e-nonce").Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		props := layout.DefaultPageProps()
		props.Nonce = "e2e-nonce"

		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func themeStorageSyncE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = themeStorageSyncE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// darkPreferenceEmulation pins prefers-color-scheme: dark so the OS-fallback
// branch of the storage sync (theme key removed → matchMedia) is OBSERVABLE:
// with a dark emulated OS, removing the key must flip the tab back to dark.
func darkPreferenceEmulation(t *testing.T, ctx context.Context) {
	t.Helper()

	emulateMedia(t, ctx, "",
		&emulation.MediaFeature{Name: "prefers-color-scheme", Value: "dark"},
	)
}

// themeStorageSyncDispatch builds a trusted-enough StorageEvent with the
// given newValue (” → newValue null, the key-removed case).
func themeStorageSyncDispatch(newValue string) string {
	nv := `null`
	if newValue != "" {
		nv = `'` + newValue + `'`
	}

	return `window.dispatchEvent(new StorageEvent('storage', {key:'theme', newValue:` + nv + `}));`
}

// TestThemeStorageSyncAcrossTabs proves the ThemeToggle `storage`-event sync
// in a real browser. The listener contract is OURS and is asserted hard via
// dispatched StorageEvents (light / dark / key-removed → emulated OS
// fallback). The cross-tab hop is the web platform's contract and headless
// Chrome does not reliably deliver storage events across chromedp targets —
// a real two-tab write is therefore evidence-gated: asserted when delivery
// is observable (probe counter > 0), logged-and-continued otherwise.
// Serial — shared browser allocator and cross-tab state.
func TestThemeStorageSyncAcrossTabs(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := themeStorageSyncE2EServer(t)

	observer, cancelObserver := newFlowTab(t)
	defer cancelObserver()

	darkPreferenceEmulation(t, observer)

	// Initial paint follows the emulated dark OS via ThemeScript; the
	// toggle's aria-checked is repaired at script attach.
	if err := chromedp.Do(observer,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
		pollTrue(`document.documentElement.classList.contains('dark')===true`, pollFast()...),
		pollTrue(`document.documentElement.style.colorScheme==='dark'`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='true'`, pollFast()...),
		pollTrue(`window.localStorage.getItem('theme')===null`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: observer initial dark state: %v", err)
	}

	// Listener contract, direction 1: explicit light flips class,
	// color-scheme, and toggle aria-checked.
	if err := chromedp.Do(observer,
		evalVoid(themeStorageSyncDispatch("light")),
		pollTrue(`document.documentElement.classList.contains('dark')===false`, pollFast()...),
		pollTrue(`document.documentElement.style.colorScheme==='light'`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='false'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: storage event sync to light: %v", err)
	}

	// Direction 2: explicit dark flips back.
	if err := chromedp.Do(observer,
		evalVoid(themeStorageSyncDispatch("dark")),
		pollTrue(`document.documentElement.classList.contains('dark')===true`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='true'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: storage event sync to dark: %v", err)
	}

	// Direction 3: key REMOVED (newValue null) falls back to the emulated OS
	// preference (dark) — proving the null branch consults matchMedia
	// instead of skipping non-dark values.
	if err := chromedp.Do(observer,
		evalVoid(themeStorageSyncDispatch("")),
		pollTrue(`document.documentElement.classList.contains('dark')===true`, pollFast()...),
		pollTrue(`document.documentElement.style.colorScheme==='dark'`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='true'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: storage event OS fallback: %v", err)
	}

	// Cross-tab evidence gate: a REAL write from a second tab must flip the
	// observer — when headless Chrome delivers storage events across
	// targets. When it does not (probe stays 0), the platform hop is
	// documented as unobservable here; the listener itself is already
	// proven above.
	writer, cancelWriter := newFlowTab(t)
	defer cancelWriter()

	darkPreferenceEmulation(t, writer)

	if err := chromedp.Do(
		observer,
		evalVoid(
			`window.__probeStorageEvents = 0; window.addEventListener('storage', function(){ window.__probeStorageEvents++; });`,
		),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: observer probe attach: %v", err)
	}

	if err := chromedp.Do(writer,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
		evalVoid(`window.localStorage.setItem('theme','light');`),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: writer light write: %v", err)
	}

	var probeCount string
	if err := chromedp.Do(observer,
		pollText(`String(window.__probeStorageEvents > 0 ? 'delivered' : 'none')`, &probeCount, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: probe read: %v", err)
	}

	if probeCount == "none" {
		t.Log("headless Chrome did not deliver the cross-target storage event; " +
			"listener contract proven via dispatched StorageEvents above")

		return
	}

	if err := chromedp.Do(observer,
		pollTrue(`document.documentElement.classList.contains('dark')===false`, pollFast()...),
		pollTrue(`document.documentElement.style.colorScheme==='light'`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='false'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: cross-tab observer sync to light: %v", err)
	}
}
