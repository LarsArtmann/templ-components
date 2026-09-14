package visualtest

import (
	"context"
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestTouchTargetAudit (M21/F097): at a 375px mobile viewport, every visible
// interactive element on the demo routes must offer at least a 24×24 CSS-px
// target (WCAG 2.2 AA target-size minimum). The 44×44 AAA/Apple bar is
// logged, never gated — crossing it is a deliberate design release.
// Native checkboxes/radios are logged but exempt: their hit area is the
// user agent's, and the library does not resize them.
func TestTouchTargetAudit(t *testing.T) {
	t.Parallel()

	server := StartDemoServer(t)

	for _, route := range axeSweepRoutes {
		if route.dark {
			continue // target geometry is palette-independent
		}

		t.Run(route.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			var raw string

			if err := chromedp.Run(ctx,
				chromedp.EmulateViewport(375, 667),
				chromedp.Navigate(server.BaseURL()+route.path),
				chromedp.WaitReady("body"),
				chromedp.Sleep(settleDelay),
				chromedp.Evaluate(touchTargetProbe, &raw),
			); err != nil {
				t.Fatalf("touch target audit %s: %v", route.path, err)
			}

			var findings []touchTargetFinding

			if err := json.Unmarshal([]byte(raw), &findings); err != nil {
				t.Fatalf("touch target audit %s: decode %v (raw %s)", route.path, err, raw)
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

// touchTargetFinding is one interactive element whose smaller dimension is
// under 24 CSS px.
type touchTargetFinding struct {
	Sel    string `json:"sel"`
	W      int    `json:"w"`
	H      int    `json:"h"`
	Exempt bool   `json:"exempt"`
}

// touchTargetProbe lists visible interactive elements under 24×24 CSS px.
// Spec-backed exemptions (WCAG 2.2 target-size exceptions):
//   - sr-only / clip-hidden controls (e.g. the Base skip link): not a
//     visible target until focused, and focused it is fully sized.
//   - inline text links (display:inline, text-height box): the Inline
//     exception — breadcrumbs, in-sentence links, FieldError links.
const touchTargetProbe = `(() => {
  const out = [];
  const els = document.querySelectorAll('button, a[href], input, select, textarea, summary, [role="button"], [role="tab"]');
  for (const el of els) {
    const r = el.getBoundingClientRect();
    const s = getComputedStyle(el);
    if (s.display === 'none' || s.visibility === 'hidden' || r.width < 1 || r.height < 1) continue;
    if (Math.min(r.width, r.height) >= 24) continue;
    if (el.classList.contains('sr-only')) continue;
    const type = (el.getAttribute('type') || '').toLowerCase();
    const nativeControl = el.tagName === 'INPUT' && (type === 'checkbox' || type === 'radio');
    const inlineTextLink = el.tagName === 'A' && s.display === 'inline';
    if (nativeControl || inlineTextLink) continue;
    let sel = el.tagName.toLowerCase();
    if (el.id) sel += '#' + el.id;
    const al = el.getAttribute('aria-label');
    if (al) sel += '[aria-label="' + al + '"]';
    const txt = (el.textContent || '').trim().slice(0, 24);
    if (txt && !al) sel += ' "' + txt + '"';
    out.push({sel, w: Math.round(r.width), h: Math.round(r.height), exempt: false});
  }
  return JSON.stringify(out);
})()`

// TestZoomReflowAudit (M21/F099): at 200% and 400% effective zoom (640px and
// 320px CSS-px viewports — the WCAG 1.4.10 reflow condition), no demo route
// may scroll horizontally. The worst offender is reported for triage.
func TestZoomReflowAudit(t *testing.T) {
	t.Parallel()

	server := StartDemoServer(t)

	for _, zoom := range []struct {
		name  string
		width int64
	}{
		{name: "zoom200", width: 640},
		{name: "zoom400", width: 320},
	} {
		for _, route := range axeSweepRoutes {
			if route.dark {
				continue
			}

			t.Run(zoom.name+"/"+route.name, func(t *testing.T) {
				t.Parallel()

				ctx, cancel := newTab(t)
				defer cancel()

				ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
				defer cancelTimeout()

				var raw string

				if err := chromedp.Run(ctx,
					chromedp.EmulateViewport(zoom.width, 900),
					chromedp.Navigate(server.BaseURL()+route.path),
					chromedp.WaitReady("body"),
					chromedp.Sleep(settleDelay),
					chromedp.Evaluate(reflowProbe, &raw),
				); err != nil {
					t.Fatalf("reflow audit %s: %v", route.path, err)
				}

				var finding reflowFinding

				if err := json.Unmarshal([]byte(raw), &finding); err != nil {
					t.Fatalf("reflow audit %s: decode %v (raw %s)", route.path, err, raw)
				}

				if finding.Overflow > 1 {
					t.Errorf(
						"%s at %dpx: horizontal overflow of %dpx (worst offender: %s) — WCAG 1.4.10 reflow violation",
						route.name, zoom.width, finding.Overflow, finding.Worst,
					)
				}
			})
		}
	}
}

// reflowFinding reports a page's horizontal overflow and its widest element.
type reflowFinding struct {
	Overflow int    `json:"overflow"`
	Worst    string `json:"worst"`
}

// reflowProbe measures document-level horizontal overflow and names the
// element extending furthest past the viewport.
const reflowProbe = `(() => {
  const de = document.documentElement;
  const over = de.scrollWidth - de.clientWidth;
  if (over <= 1) return JSON.stringify({overflow: 0, worst: ''});
  let worst = null, maxRight = 0;
  for (const el of document.querySelectorAll('body *')) {
    const r = el.getBoundingClientRect();
    if (r.width > 0 && r.right > maxRight) { maxRight = r.right; worst = el; }
  }
  let desc = '?';
  if (worst) {
    desc = worst.tagName.toLowerCase() + (worst.id ? '#' + worst.id : '');
    const cls = (worst.getAttribute('class') || '').split(' ').filter(Boolean).slice(0, 3).join('.');
    if (cls) desc += '.' + cls;
  }
  return JSON.stringify({overflow: Math.round(over), worst: desc});
})()`
