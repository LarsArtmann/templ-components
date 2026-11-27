package visualtest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
)

// The Datastar SignalChip browser E2E proves the client-side filter contract
// against the PINNED runtime. The component's own tests pin the markup it
// emits; only a browser can prove the runtime reads that markup the way the
// component assumes:
//
//   - at load the runtime settles in the server-rendered state, because the
//     selected chip's data-signals declares the signal's initial value;
//   - clicking a chip moves aria-pressed and the pressed look together;
//   - a chip hand-spelled data-on-click (the dash spelling) stays dead: the
//     runtime drops unknown plugin names silently, which is the consumer bug
//     the component exists to prevent.
//
// TestDatastarSignalChipContrast adds the measurable half of the chip's
// accessibility contract: every label and count badge stays at WCAG AA
// contrast (4.5:1) on its own chip background, at rest AND on hover, in light
// and dark mode, for the demo's branded accent and for the library's default
// one. TestDatastarSignalChipDemo smoke-tests the demo's own chip entry on
// /datastar against the live demo server.
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.

// signalChipE2EState renders the observable state of every chip in the group:
// "1"/"0" for aria-pressed, then "P"/"U" for the pressed/unpressed look.
const signalChipE2EState = `(() => {
	const chips = Array.from(document.querySelectorAll('#chips button'));
	return chips.map((b) => (b.getAttribute('aria-pressed') === 'true' ? '1' : '0') + (b.classList.contains('bg-accent-600') ? 'P' : 'U')).join(',');
})()`

// signalChipE2ESignal reads the live signal through a data-text binding.
const signalChipE2ESignal = `document.getElementById('sig').textContent`

// signalChipE2EChip selects the n-th chip of the group.
func signalChipE2EChip(n int) string {
	return "#chips button:nth-of-type(" + strconv.Itoa(n) + ")"
}

func datastarSignalChipE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	accentCSS := libraryAccentCSS(t)

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/library-accent.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write([]byte(accentCSS))
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// ?library-accent=1 layers the library's default accent ramp over the
		// demo's branded one (see libraryAccentCSS).
		page := signalChipE2EPage(r.URL.Query().Has("library-accent"))

		if err := page.Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// signalChipE2EPage is the fixture: a three-chip group with "open" pressed at
// first paint, two chips carrying a count badge, a hand-spelled dash-spelling
// control, and a data-text read-out of the live signal.
func signalChipE2EPage(libraryAccent bool) templ.Component {
	openCount, doneCount := 12, 7

	props := layout.DefaultPageProps()
	props.Title = "Datastar SignalChip E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		); err != nil {
			return err
		}

		if !libraryAccent {
			return nil
		}

		// Rendered after the demo stylesheet link. Unlayered, so it wins over
		// the @layer theme block the demo's accent override lives in.
		_, err := io.WriteString(w, `<link rel="stylesheet" href="/library-accent.css">`)

		return err
	})

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := datastar.SDKScript(datastar.SDKScriptProps{Src: "/datastar.js"}).Render(ctx, w); err != nil {
			return err
		}

		if err := datastar.SignalChipGroup(datastar.SignalChipGroupProps{
			BaseProps: utils.BaseProps{ID: "chips", AriaLabel: "Filter by state"},
			Signal:    "stateFilter",
			Selected:  "open",
			Items: []datastar.SignalChipItem{
				{Label: "All"},
				{Label: "Open", Value: "open", Count: &openCount},
				{Label: "Done", Value: "done", Count: &doneCount},
			},
		}).Render(ctx, w); err != nil {
			return err
		}

		// Negative control: the exact markup a consumer shipped by mistake. It
		// assigns a value NO chip carries, so any change would be unambiguous.
		_, err := io.WriteString(w,
			`<button id="dash-spelled" type="button" data-on-click="$stateFilter='dash'">dash spelled</button>`+
				`<span id="sig" data-text="$stateFilter"></span>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// TestDatastarSignalChipBrowser proves the pressed state, the look, and the
// signal move together on click, and that the dash spelling stays inert.
func TestDatastarSignalChipBrowser(t *testing.T) {
	t.Parallel()

	srv := datastarSignalChipE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollTrue(`window.__dsReady===true && document.querySelectorAll('#chips button').length===3`),
		// The runtime settled in the server-rendered state: the pressed chip
		// kept both its aria-pressed state and its pressed look.
		pollTrue(signalChipE2EState+`==='0U,1P,0U'`),
		pollTrue(signalChipE2ESignal+`==='open'`),
		// The dash spelling names an unknown plugin: the click is delivered
		// and nothing happens, with no console error to notice.
		chromedp.Click("#dash-spelled", chromedp.NodeVisible),
		chromedp.Sleep(settleDelay),
		pollTrue(signalChipE2EState+`==='0U,1P,0U'`),
		pollTrue(signalChipE2ESignal+`==='open'`),
		// The colon spelling moves the state, the look and the signal.
		chromedp.Click(signalChipE2EChip(3), chromedp.NodeVisible),
		pollTrue(signalChipE2EState+`==='0U,0U,1P'`),
		pollTrue(signalChipE2ESignal+`==='done'`),
		// The empty-value chip clears the filter.
		chromedp.Click(signalChipE2EChip(1), chromedp.NodeVisible),
		pollTrue(signalChipE2EState+`==='1P,0U,0U'`),
		pollTrue(signalChipE2ESignal+`===''`),
	); err != nil {
		t.Fatalf("SignalChip pressed-state proof: %v", err)
	}
}

// signalChipDemoVisibleTasks lists the data-status of every task the demo's
// chip group currently leaves visible (data-show hides the rest).
const signalChipDemoVisibleTasks = `Array.from(document.querySelectorAll('#task-filter-list > li'))
	.filter((li) => li.getClientRects().length > 0)
	.map((li) => li.dataset.status).join(',')`

// TestDatastarSignalChipDemo smoke-tests the demo's own chip section on
// /datastar: the route loads the Datastar runtime (in-content SDKScript), the
// runtime settles in the server-rendered state, and a click on a chip moves
// the pressed state AND filters the list bound to the same signal. A demo
// entry whose data-on:* markup is inert because the SDK never loaded is the
// 2026-10-03 outage class (pages.go: NeedsDatastar).
func TestDatastarSignalChipDemo(t *testing.T) {
	t.Parallel()

	server := StartDemoServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(server.BaseURL()+"/datastar"),
		// The runtime settled in the server-rendered state: "open" pressed,
		// only the open tasks listed.
		pollTrue(`(`+signalChipDemoVisibleTasks+`)==='open,open'`),
		pollTrue(`document.querySelector('#task-filter-chips button:nth-of-type(2)').getAttribute('aria-pressed')==='true'`),
		chromedp.Click("#task-filter-chips button:nth-of-type(3)", chromedp.NodeVisible),
		pollTrue(`(`+signalChipDemoVisibleTasks+`)==='done'`),
		pollTrue(`document.querySelector('#task-filter-chips button:nth-of-type(3)').getAttribute('aria-pressed')==='true'`),
		pollTrue(`document.querySelector('#task-filter-chips button:nth-of-type(2)').getAttribute('aria-pressed')==='false'`),
		// The empty-value chip clears the filter.
		chromedp.Click("#task-filter-chips button:nth-of-type(1)", chromedp.NodeVisible),
		pollTrue(`(`+signalChipDemoVisibleTasks+`)==='open,open,done'`),
		pollTrue(`document.querySelector('#task-filter-chips button:nth-of-type(1)').getAttribute('aria-pressed')==='true'`),
	); err != nil {
		t.Fatalf("demo /datastar SignalChip section: %v", err)
	}

	server.FailIfServerErrors(t)
}

// wcagAATextContrast is WCAG 2.2's AA threshold for text below large size.
// The chip label and its count badge are 14px.
const wcagAATextContrast = 4.5

// contrastRecord is one browser-measured WCAG contrast ratio. Name spells
// out what was measured; see contrastRecordName.
type contrastRecord struct {
	Name    string  `json:"name"`
	Ratio   float64 `json:"ratio"`
	FG      string  `json:"fg"`
	BG      string  `json:"bg"`
	Pressed string  `json:"pressed"`
}

//nolint:err113 // test-fixture sentinel, never wrapped
var errContrastNoRecords = errors.New("contrast measurement returned no records")

// contrastRecordName identifies one measurement:
// "<palette>/<rest|hover>/<light|dark>/chip<n>/<label|count>".
func contrastRecordName(palette, state, theme string, chip int, part string) string {
	return palette + "/" + state + "/" + theme + "/chip" + strconv.Itoa(chip) + "/" + part
}

// TestDatastarSignalChipContrast measures, in the browser, the contrast of
// every chip label and count badge against its own chip background. The
// string tests can pin class names; only the browser can pin what those
// classes are worth once the theme, the hover state, and the compositor are
// involved (the count badge once shipped at 70% opacity: 3.6:1 on the
// pressed chip, invisible to every string assertion).
//
// Two palettes are measured: the demo's branded accent (the stylesheet every
// visual test loads) and the library's default one (libraryAccentCSS). The
// suite's shared browser reports (hover: none), so the test drives a private
// hover-capable Chromium (newHoverCapableTab).
func TestDatastarSignalChipContrast(t *testing.T) {
	t.Parallel()

	srv := datastarSignalChipE2EServer(t)

	ctx, cancel := newHoverCapableTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)
	defer cancelTimeout()

	// selectedBG remembers each palette's pressed-chip background, to prove
	// the two palettes really rendered different accents.
	selectedBG := map[string]string{}

	for _, palette := range []struct{ name, query string }{
		{name: "demo-accent", query: ""},
		{name: "library-accent", query: "?library-accent=1"},
	} {
		var (
			hoverCapable bool
			records      []contrastRecord
		)

		steps := []chromedp.Action[chromedp.Void]{
			chromedp.Navigate(srv.URL + "/" + palette.query),
			pollTrue(`window.__dsReady===true && document.querySelectorAll('#chips button').length===3`),
			evalExprInto(`matchMedia('(hover: hover)').matches`, &hoverCapable),
		}

		for _, theme := range []string{"light", "dark"} {
			steps = append(steps,
				evalVoid(`document.documentElement.classList.toggle('dark', `+strconv.FormatBool(theme == "dark")+`)`),
				waitAnimationsSettled(),
				collectContrast(palette.name, "rest", theme, 0, &records),
				// Hover the pressed chip (label + count), then the first chip
				// (the unpressed look).
				hoverAction(signalChipE2EChip(2)), waitAnimationsSettled(),
				collectContrast(palette.name, "hover", theme, 2, &records),
				hoverAction(signalChipE2EChip(1)), waitAnimationsSettled(),
				collectContrast(palette.name, "hover", theme, 1, &records),
				resetCaptureNeutralState(), waitAnimationsSettled(),
			)
		}

		if err := chromedp.Do(ctx, steps...); err != nil {
			t.Fatalf("%s: contrast measurement: %v", palette.name, err)
		}

		if !hoverCapable {
			t.Fatalf("%s: the private browser does not match (hover: hover); hover styles cannot be measured", palette.name)
		}

		byName := make(map[string]contrastRecord, len(records))
		for _, rec := range records {
			byName[rec.Name] = rec
		}

		for _, theme := range []string{"light", "dark"} {
			rest, ok := byName[contrastRecordName(palette.name, "rest", theme, 2, "label")]
			if !ok {
				t.Fatalf("%s/%s: no rest record for the pressed chip", palette.name, theme)
			}

			if rest.Pressed != "true" {
				t.Fatalf("%s/%s: chip 2 is not the pressed chip (aria-pressed=%q)", palette.name, theme, rest.Pressed)
			}

			selectedBG[palette.name] = rest.BG
		}

		assertSignalChipContrast(t, palette.name, records, byName)
	}

	if selectedBG["demo-accent"] == "" || selectedBG["demo-accent"] == selectedBG["library-accent"] {
		t.Fatalf(
			"the library-accent override did not change the pressed chip's background (demo %q, library %q)",
			selectedBG["demo-accent"], selectedBG["library-accent"],
		)
	}
}

// collectContrast appends the page's contrast records for one pointer state
// to dest. only > 0 restricts the measurement to that chip.
func collectContrast(
	palette, state, theme string, only int, dest *[]contrastRecord,
) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tgt *chromedp.Target) error {
		recs, err := chromedp.Evaluate[[]contrastRecord](signalChipContrastJS(palette, state, theme, only))(ctx, tgt)
		if err != nil {
			return err
		}

		if len(recs) == 0 {
			return errContrastNoRecords
		}

		*dest = append(*dest, recs...)

		return nil
	})
}

// assertSignalChipContrast fails on any ratio under WCAG AA, and first proves
// the measurement is not vacuous: every hovered chip must have changed its
// background between rest and hover (otherwise hover: styles never applied
// and the hover records would merely repeat the rest state).
func assertSignalChipContrast(t *testing.T, palette string, records []contrastRecord, byName map[string]contrastRecord) {
	t.Helper()

	for _, theme := range []string{"light", "dark"} {
		for _, chip := range []int{1, 2} {
			rest := byName[contrastRecordName(palette, "rest", theme, chip, "label")]
			hover := byName[contrastRecordName(palette, "hover", theme, chip, "label")]

			if rest.BG == "" || hover.BG == "" {
				t.Fatalf("%s/%s/chip%d: missing rest or hover record", palette, theme, chip)
			}

			if rest.BG == hover.BG {
				t.Errorf(
					"%s/%s/chip%d: hovering left the chip background at %s — hover styles did not apply, "+
						"so the hover records are not a measurement",
					palette, theme, chip, rest.BG,
				)
			}
		}
	}

	for _, rec := range records {
		t.Logf("%s: %.2f:1 (text %s on %s, aria-pressed=%s)", rec.Name, rec.Ratio, rec.FG, rec.BG, rec.Pressed)

		if rec.Ratio < wcagAATextContrast {
			t.Errorf(
				"%s: contrast %.2f:1 (text %s on %s) is below WCAG AA (%.1f:1)",
				rec.Name, rec.Ratio, rec.FG, rec.BG, wcagAATextContrast,
			)
		}
	}
}

// signalChipContrastJS measures, inside the page, the WCAG 2.x contrast ratio
// of each chip's label (and count badge when it carries one) against the
// chip's own background, in the CURRENT theme and pointer state:
//
//   - colors are canvas-resolved: getComputedStyle hands back non-sRGB
//     notations verbatim (the default theme's tokens are oklch()), so every
//     color is painted into a 1x1 canvas and read back as sRGB bytes — the
//     resolution the compositor performs;
//   - an element's CSS opacity is applied the way the compositor does (alpha
//     over the backdrop), so a faded badge measures as its on-screen color;
//   - the chip background is composited over the page background first.
func signalChipContrastJS(palette, state, theme string, only int) string {
	return strings.NewReplacer(
		"__PALETTE__", palette,
		"__STATE__", state,
		"__THEME__", theme,
		"__ONLY__", strconv.Itoa(only),
	).Replace(`(() => {
	const only = __ONLY__;
	const prefix = '__PALETTE__/__STATE__/__THEME__/';
	const canvas = document.createElement('canvas');
	canvas.width = canvas.height = 1;
	const ctx = canvas.getContext('2d', { willReadFrequently: true });
	const hex = (d) => '#' + [d[0], d[1], d[2]].map((v) => v.toString(16).padStart(2, '0')).join('');
	// Paint one CSS color over a backdrop; returns the composited sRGB bytes.
	const paint = (color, alpha, backdrop) => {
		ctx.clearRect(0, 0, 1, 1);
		ctx.globalAlpha = 1;
		if (backdrop) { ctx.fillStyle = backdrop.css; ctx.fillRect(0, 0, 1, 1); }
		ctx.globalAlpha = alpha;
		ctx.fillStyle = color;
		ctx.fillRect(0, 0, 1, 1);
		const d = ctx.getImageData(0, 0, 1, 1).data;
		return { css: hex(d), rgb: [d[0], d[1], d[2]] };
	};
	const lum = (rgb) => {
		const f = (v) => { const c = v / 255; return c <= 0.04045 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4); };
		return 0.2126 * f(rgb[0]) + 0.7152 * f(rgb[1]) + 0.0722 * f(rgb[2]);
	};
	const ratio = (a, b) => {
		const x = lum(a), y = lum(b);
		return (Math.max(x, y) + 0.05) / (Math.min(x, y) + 0.05);
	};
	// The page backdrop: html, then body, each over white (a transparent
	// background leaves the layer beneath unchanged).
	let page = paint('#ffffff', 1, null);
	page = paint(getComputedStyle(document.documentElement).backgroundColor, 1, page);
	page = paint(getComputedStyle(document.body).backgroundColor, 1, page);
	const out = [];
	document.querySelectorAll('#chips button').forEach((btn, idx) => {
		if (only > 0 && idx + 1 !== only) return;
		const cs = getComputedStyle(btn);
		const bg = paint(cs.backgroundColor, 1, page);
		const pressed = btn.getAttribute('aria-pressed');
		const name = prefix + 'chip' + (idx + 1) + '/';
		const label = paint(cs.color, parseFloat(cs.opacity), bg);
		out.push({ name: name + 'label', ratio: ratio(label.rgb, bg.rgb), fg: label.css, bg: bg.css, pressed: pressed });
		const span = btn.querySelector('span');
		if (!span) return;
		const scs = getComputedStyle(span);
		const count = paint(scs.color, parseFloat(scs.opacity), bg);
		out.push({ name: name + 'count', ratio: ratio(count.rgb, bg.rgb), fg: count.css, bg: bg.css, pressed: pressed });
	});
	return out;
})()`)
}

// newHoverCapableTab returns a tab in a PRIVATE Chromium that reports a fine,
// hovering pointer. The suite's shared browser reports (hover: none) — CDP
// cannot emulate the feature — and Tailwind gates every hover: utility in
// @media (hover:hover), so on the shared browser a real pointer hover never
// applies hover: styles. The extra process is torn down with the test.
func newHoverCapableTab(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	if !browserConfigured() {
		t.Skipf("visual tests skipped: %v (set CHROMEDP_CHROME_PATH)", errNoBrowser)
	}

	//nolint:gosec // path is from CHROMEDP_CHROME_PATH env var, not user-controlled
	chromePath := os.Getenv("CHROMEDP_CHROME_PATH")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("font-render-hinting", "none"),
		chromedp.Flag("disable-background-timer-throttling", true),
		// Blink hover/pointer types: 2 = can hover, 4 = fine pointer (the
		// flags Playwright launches with for the same reason).
		chromedp.Flag("blink-settings", "primaryHoverType=2,availableHoverTypes=2,primaryPointerType=4,availablePointerTypes=4"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	t.Cleanup(allocCancel)

	ctx, cancel := chromedp.NewContext(allocCtx)

	bounded, cancelBounded := context.WithTimeout(ctx, tabLifetime)

	return bounded, func() {
		cancelBounded()
		cancel()
	}
}

// libraryAccentTokenRe matches one accent token definition in
// templates/custom.css (ADR-0044/0049), e.g.
// `--color-accent-600: oklch(54.6% 0.245 262.881);`.
//
//nolint:gochecknoglobals // compiled once for the contrast test above
var libraryAccentTokenRe = regexp.MustCompile(`--color-accent-\d+:\s*[^;]+;`)

// libraryAccentCSS resets the accent family to the library's OWN defaults,
// read from templates/custom.css (the token source of truth), as one
// unlayered :root block. The demo stylesheet brands accent as indigo, so
// without this override a default-theme contrast failure would be invisible
// in a test that loads the demo CSS.
func libraryAccentCSS(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed: cannot locate the visualtest package dir")
	}

	src, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "templates", "custom.css"))
	if err != nil {
		t.Fatalf("read templates/custom.css: %v", err)
	}

	var tokens []string

	for _, line := range strings.Split(string(src), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "//") {
			continue // the doc comment shows an example override
		}

		if token := libraryAccentTokenRe.FindString(line); token != "" {
			tokens = append(tokens, token)
		}
	}

	joined := strings.Join(tokens, "")

	for _, shade := range []string{"--color-accent-600:", "--color-accent-700:"} {
		if !strings.Contains(joined, shade) {
			t.Fatalf("templates/custom.css defines no %s token — the pressed chip's shades are gone", shade)
		}
	}

	return ":root{" + joined + "}"
}
