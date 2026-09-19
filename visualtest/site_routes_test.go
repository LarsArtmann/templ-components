package visualtest

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"image/png"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// The SITE route tier (#T7/#T8 of the 2026-09-19 sales-page plan): the built
// website (website/dist) is its own largest library consumer, so it gets the
// same regression net as the demo — theme-pinned route goldens, the axe
// sweep, and the touch-target/zoom-reflow audits. Unlike the demo, site
// sections scroll-reveal via [data-animate] (gated on html.js), so every
// capture and audit SCROLLS through the page first — without that, below-fold
// sections sit at opacity 0 and goldens pin blank sections while axe skips
// them entirely.
//
// The dist must be built before running: `nix develop -c bash website/build.sh`
// (the `.#visual` flake app does this automatically; CI's Visual Regression
// job runs through that app and FAILS on skips, so a missing dist is a hard
// failure, never a skip).

// siteDistBase serves ../website/dist once per process with Firebase
// cleanUrls semantics (/foo resolves to foo.html). Empty string = the dist is
// missing (requireSiteDist turns that into a loud failure).
var siteDistBase = sync.OnceValue(func() string {
	distRoot := filepath.Join("..", "website", "dist")

	if _, err := os.Stat(filepath.Join(distRoot, "sales.html")); err != nil {
		return ""
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return ""
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clean, _, _ := strings.Cut(r.URL.Path, "?")
		if !strings.HasSuffix(clean, ".html") && !strings.Contains(clean, ".") {
			cleanPath := filepath.Join(distRoot, clean+".html")
			if _, err := os.Stat(cleanPath); err == nil { //nolint:gosec // test-controlled dist root
				http.ServeFile(w, r, cleanPath) //nolint:gosec // test-controlled dist root

				return
			}
		}

		http.FileServer(http.Dir(distRoot)).ServeHTTP(w, r)
	})

	server := &http.Server{Handler: mux} //nolint:gosec,exhaustruct_v5 // loopback-only test server

	go func() { _ = server.Serve(listener) }()

	return "http://" + listener.Addr().String()
})

// requireSiteDist returns the dist server base URL, failing the test with the
// build command when the dist has not been built.
func requireSiteDist(t *testing.T) string {
	t.Helper()

	base := siteDistBase()
	if base == "" {
		t.Fatalf("website dist not built — run `nix develop -c bash website/build.sh` first (the .#visual flake app does this automatically)")
	}

	return base
}

// siteScrollRevealJS scrolls through the full document in steps (then returns
// to the top) so IntersectionObserver-driven [data-animate] reveals fire for
// every section before a capture or audit. Same pattern as the siteshots
// tool; kept local because that tool is a separate main package.
const siteScrollRevealJS = `(() => new Promise((resolve) => {
	let y = 0;
	const step = innerHeight * 0.8;
	const max = document.body.scrollHeight;
	const tick = () => {
		y += step;
		if (y >= max) { scrollTo(0, 0); setTimeout(resolve, 400); return; }
		scrollTo(0, y);
		setTimeout(tick, 120);
	};
	tick();
}))()`

// siteGoldenRoutes are the pinned site pages. sales is the shipped product
// being protected; index is the sibling every nav/CTA change also touches.
var siteGoldenRoutes = []struct {
	name string
	path string
}{
	{name: "site-sales", path: "/sales"},
	{name: "site-index", path: "/"},
}

// TestSiteRouteGoldens pins light+dark × desktop+mobile full-page captures of
// the site's own routes under routes/site-*. Golden names are namespaced
// "site-" so the demo and site suites never collide in testdata/routes/.
func TestSiteRouteGoldens(t *testing.T) {
	t.Parallel()

	base := requireSiteDist(t)

	viewports := []struct {
		name     string
		viewport Viewport
	}{
		{name: "desktop", viewport: Viewport{Width: 1440, Height: 900}},
		{name: "mobile", viewport: Viewport{Width: 390, Height: 844}},
	}

	for _, route := range siteGoldenRoutes {
		for _, dark := range []bool{false, true} {
			for _, vp := range viewports {
				name := fmt.Sprintf("%s-%s-%s", route.name, themeName(dark), vp.name)

				t.Run(name, func(t *testing.T) {
					t.Parallel()

					assertSiteRouteScreenshot(t, name, base+route.path, dark, vp.viewport)
				})
			}
		}
	}
}

func themeName(dark bool) string {
	if dark {
		return "dark"
	}

	return "light"
}

// assertSiteRouteScreenshot is the site variant of assertRouteScreenshot: it
// pins the theme via localStorage + reload, SCROLLS through the page so
// [data-animate] reveals fire, settles, then captures full-page and compares
// against the routes/<name> golden with the shared compare/update flow.
func assertSiteRouteScreenshot(t *testing.T, name, url string, dark bool, viewport Viewport) {
	t.Helper()

	ctx, cancel := newTab(t)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	theme := "light"
	if dark {
		theme = "dark"
	}

	var shot []byte

	tasks := []chromedp.Action{
		chromedp.EmulateViewport(int64(viewport.Width), int64(viewport.Height)),
		chromedp.Navigate(url),
		chromedp.Evaluate(fmt.Sprintf(`try { localStorage.setItem('theme', %q); } catch (e) {}`, theme), nil),
		chromedp.Reload(),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	}

	if dark {
		tasks = append(tasks,
			chromedp.Evaluate(`document.documentElement.classList.add('dark');`, nil),
			chromedp.Sleep(500*time.Millisecond),
		)
	}

	tasks = append(tasks,
		chromedp.Evaluate(siteScrollRevealJS, nil),
		waitAnimationsSettled(),
		chromedp.Sleep(settleDelay),
		// quality 100 = PNG (any other value yields JPEG).
		chromedp.FullScreenshot(&shot, 100),
	)

	if err := chromedp.Run(timeoutCtx, tasks...); err != nil {
		t.Fatalf("site route golden[%s]: capture: %v", name, err)
	}

	compareSiteGolden(t, name, shot)
}

// compareSiteGolden mirrors assertRouteScreenshot's compare/update tail.
func compareSiteGolden(t *testing.T, name string, shot []byte) {
	t.Helper()

	if *update {
		writeGolden(t, "routes/"+name, shot)

		return
	}

	goldenImg, exists := readGolden(t, "routes/"+name)
	if !exists {
		writeGolden(t, "routes/"+name, shot)
		t.Errorf(
			"site route golden[%s]: no golden yet — wrote %s (re-run without -update to verify)",
			name,
			goldenPath("routes/"+name),
		)

		return
	}

	actualImg, err := png.Decode(bytes.NewReader(shot))
	if err != nil {
		t.Fatalf("site route golden[%s]: decode actual: %v", name, err)
	}

	result, diff := comparePixels(
		goldenImg,
		actualImg,
		defaultOptions(Options{}).Threshold,
		defaultOptions(Options{}).MaxMismatch*percentMultiplier,
	)
	if !result.Match {
		writeFailureArtifacts(t, "routes/"+name, shot, diff)
		t.Errorf("site route golden[%s]: visual mismatch — %s (max %.4f%%).\n"+
			"Inspect testdata/.fail/routes.%s.{actual,diff}.png, then run `nix run .#visual -- -update -run TestSiteRouteGoldens` if intended.",
			name, result, defaultOptions(Options{}).MaxMismatch*percentMultiplier, name)

		return
	}

	cleanFailureArtifacts("routes/" + name)
	t.Logf("site route golden[%s]: OK (%.4f%% mismatched)", name, result.MismatchPct)
}

// axeSiteRoutes are the site pages under a11y audit (light + dark: the site's
// own dark palette is where contrast regressions hide). Names are prefixed
// "site_" inside the shared axe baseline ledger.
var axeSiteRoutes = []struct {
	name string
	path string
	dark bool
}{
	{name: "site_sales", path: "/sales", dark: false},
	{name: "site_sales_dark", path: "/sales", dark: true},
	{name: "site_index", path: "/", dark: false},
	{name: "site_index_dark", path: "/", dark: true},
}

// TestAxeSweepSiteRoutes extends the default-fail axe gate to the built
// website. Same ledger, same policy: critical/serious violations must be
// accepted in testdata/axe_baseline.json (site_* keys) or the build fails.
func TestAxeSweepSiteRoutes(t *testing.T) {
	t.Parallel()

	base := requireSiteDist(t)
	baseline := readAxeBaseline(t)

	for _, route := range axeSiteRoutes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			results := siteAxeAuditRoute(t, ctx, base, route.path, route.dark)

			blocking := results.BlockingViolations()
			t.Logf("axe[%s]: %d rule(s) total, %d blocking", route.name, len(results.Violations), len(blocking))
			assertNoUnacceptedViolations(t, route.name, baseline, results)
		})
	}
}

// siteAxeAuditRoute is axeAuditRoute plus the scroll-through: the site's
// [data-animate] sections sit at opacity 0 until intersected, and axe skips
// invisible elements — an unscrolled audit would silently exempt every
// below-fold section from contrast checks.
func siteAxeAuditRoute(t *testing.T, ctx context.Context, baseURL, path string, dark bool) AxeResults {
	t.Helper()

	theme := "light"
	if dark {
		theme = "dark"
	}

	actions := []chromedp.Action{
		chromedp.Navigate(baseURL + path),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(fmt.Sprintf(
			`localStorage.setItem('theme', %q); document.documentElement.classList.toggle('dark', %t); document.documentElement.style.colorScheme = %q; true`,
			theme,
			dark,
			theme,
		), nil),
		chromedp.Evaluate(siteScrollRevealJS, nil),
		// The scroll-through FIRES the [data-animate] entrance animations;
		// axe samples computed color+opacity mid-animation and reports bogus
		// blended-contrast findings (measured: elements at ~38% opacity).
		// Settle every finite animation before auditing.
		waitAnimationsSettled(),
		chromedp.Sleep(settleDelay),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for range 2 {
				var got bool
				if err := chromedp.Evaluate(
					`document.documentElement.classList.contains('dark')`, &got,
				).Do(ctx); err != nil {
					return err
				}

				if got == dark {
					return nil
				}

				if err := chromedp.Evaluate(fmt.Sprintf(
					`localStorage.setItem('theme', %q); document.documentElement.classList.toggle('dark', %t); document.documentElement.style.colorScheme = %q; true`,
					theme,
					dark,
					theme,
				), nil).Do(ctx); err != nil {
					return err
				}

				if err := chromedp.Sleep(settleDelay).Do(ctx); err != nil {
					return err
				}
			}

			return fmt.Errorf("visualtest[axe]: theme pin did not stick on %s%s (wanted dark=%t)", baseURL, path, dark)
		}),
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		t.Fatalf("visualtest[axe]: load %s%s: %v", baseURL, path, err)
	}

	results, err := RunAxe(ctx)
	if err != nil {
		t.Fatalf("visualtest[axe]: audit %s%s: %v", baseURL, path, err)
	}

	return results
}

// TestSiteTouchTargetAudit runs the 24px target audit against the site's
// light routes at a 375px viewport. The site header (search box, toggle,
// text links) is its own UI surface, not covered by the demo audit.
func TestSiteTouchTargetAudit(t *testing.T) {
	t.Parallel()

	base := requireSiteDist(t)

	for _, route := range siteGoldenRoutes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			var raw string

			if err := chromedp.Run(ctx,
				chromedp.EmulateViewport(375, 667),
				chromedp.Navigate(base+route.path),
				chromedp.WaitReady("body"),
				chromedp.Evaluate(siteScrollRevealJS, nil),
				chromedp.Sleep(settleDelay),
				chromedp.Evaluate(touchTargetProbe, &raw),
			); err != nil {
				t.Fatalf("site touch target audit %s: %v", route.path, err)
			}

			var findings []touchTargetFinding

			if err := json.Unmarshal([]byte(raw), &findings); err != nil {
				t.Fatalf("site touch target audit %s: decode %v (raw %s)", route.path, err, raw)
			}

			for _, f := range findings {
				if f.Exempt {
					t.Logf("%s: sub-24 native control (exempt, logged): %s %dx%d", route.name, f.Sel, f.W, f.H)

					continue
				}

				t.Errorf(
					"%s: %s renders a %dx%d target — below the WCAG 2.2 AA 24px minimum. "+
						"Enlarge the control (padding/size), not the audit.",
					route.name, f.Sel, f.W, f.H,
				)
			}
		})
	}
}

// TestSiteZoomReflowAudit runs the WCAG 1.4.10 reflow audit against the
// site's light routes at 640px/320px.
func TestSiteZoomReflowAudit(t *testing.T) {
	t.Parallel()

	base := requireSiteDist(t)

	for _, zoom := range []struct {
		name  string
		width int64
	}{
		{name: "zoom200", width: 640},
		{name: "zoom400", width: 320},
	} {
		for _, route := range siteGoldenRoutes {
			t.Run(zoom.name+"/"+route.name, func(t *testing.T) {
				t.Parallel()

				ctx, cancel := newTab(t)
				defer cancel()

				ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
				defer cancelTimeout()

				var raw string

				if err := chromedp.Run(ctx,
					chromedp.EmulateViewport(zoom.width, 900),
					chromedp.Navigate(base+route.path),
					chromedp.WaitReady("body"),
					chromedp.Sleep(settleDelay),
					chromedp.Evaluate(reflowProbe, &raw),
				); err != nil {
					t.Fatalf("site reflow audit %s: %v", route.path, err)
				}

				var finding reflowFinding

				if err := json.Unmarshal([]byte(raw), &finding); err != nil {
					t.Fatalf("site reflow audit %s: decode %v (raw %s)", route.path, err, raw)
				}

				if finding.Overflow > 1 {
					t.Errorf(
						"%s at %dpx: horizontal overflow of %dpx (worst offender: %s) — WCAG 1.4.10 reflow violation",
						zoom.name+"/"+route.name, zoom.width, finding.Overflow, finding.Worst,
					)
				}
			})
		}
	}
}
