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

// TestThemeStorageSyncAcrossTabs proves the `storage`-event tab sync in a
// real browser with two tabs sharing one origin's localStorage: a write in
// tab B flips tab A's theme class, color-scheme, and toggle aria-checked;
// removing the key falls back to the (emulated) OS preference. The writer
// tab itself receives no storage event — tab A is the observer throughout.
// Serial — shared browser allocator and cross-tab state.
func TestThemeStorageSyncAcrossTabs(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := themeStorageSyncE2EServer(t)

	observer, cancelObserver := newFlowTab(t)
	defer cancelObserver()
	darkPreferenceEmulation(t, observer)

	writer, cancelWriter := newFlowTab(t)
	defer cancelWriter()
	darkPreferenceEmulation(t, writer)

	// Tab A (observer): initial paint follows the emulated dark OS via
	// ThemeScript; the toggle's aria-checked is repaired at script attach.
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

	// Tab B (writer): same page, then write theme=light. The storage event
	// fires ONLY in tab A — the observed flip is the sync listener at work.
	if err := chromedp.Do(writer,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
		evalVoid(`window.localStorage.setItem('theme','light');`),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: writer light write: %v", err)
	}
	if err := chromedp.Do(observer,
		pollTrue(`document.documentElement.classList.contains('dark')===false`, pollFast()...),
		pollTrue(`document.documentElement.style.colorScheme==='light'`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='false'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: observer sync to light: %v", err)
	}

	// Explicit dark write flips the observer back on.
	if err := chromedp.Do(writer,
		evalVoid(`window.localStorage.setItem('theme','dark');`),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: writer dark write: %v", err)
	}
	if err := chromedp.Do(observer,
		pollTrue(`document.documentElement.classList.contains('dark')===true`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='true'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: observer sync to dark: %v", err)
	}

	// Removing the key (newValue null) falls back to the emulated OS
	// preference (dark) — proving the null branch actually consults
	// matchMedia instead of skipping non-dark values.
	if err := chromedp.Do(writer,
		evalVoid(`window.localStorage.removeItem('theme');`),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: writer key removal: %v", err)
	}
	if err := chromedp.Do(observer,
		pollTrue(`window.localStorage.getItem('theme')===null`, pollFast()...),
		pollTrue(`document.documentElement.classList.contains('dark')===true`, pollFast()...),
		pollTrue(`document.documentElement.style.colorScheme==='dark'`, pollFast()...),
		pollTrue(`document.querySelector('[data-theme-toggle]').getAttribute('aria-checked')==='true'`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[theme-sync]: observer OS-preference fallback: %v", err)
	}
}
