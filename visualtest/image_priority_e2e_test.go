package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

// imagePriorityE2EPage renders one lazy (default) and one eager image, the
// two priority personalities the component promises.
func imagePriorityE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		lazy := display.ImageProps{
			Nonce: "e2e-nonce", ID: "e2e-img-lazy",
			Src: "/static/fixture.svg", Alt: "Lazy fixture", Width: 320, Height: 180,
		}
		if err := display.Image(lazy).Render(ctx, w); err != nil {
			return err
		}

		eager := display.ImageProps{
			Nonce: "e2e-nonce", ID: "e2e-img-eager",
			Src: "/static/fixture.svg", Alt: "Eager fixture", Width: 320, Height: 180,
			Lazy: false,
		}

		return display.Image(eager).Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func imagePriorityE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = imagePriorityE2EPage().Render(r.Context(), w)
	})
	mux.HandleFunc("/static/fixture.svg", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = io.WriteString(
			w,
			`<svg xmlns="http://www.w3.org/2000/svg" width="320" height="180"><rect width="100%" height="100%" fill="#3b82f6"/></svg>`,
		)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestImageE2EFetchPriorityLadder proves the loading/priority pairing in the
// DOM of a real page: the lazy default carries loading="lazy" and NO
// fetchpriority attribute (lazy + low would be the documented anti-pattern),
// while the above-the-fold configuration pairs loading="eager" with
// fetchpriority="high". Serial — shared browser allocator.
func TestImageE2EFetchPriorityLadder(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := imagePriorityE2EServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
		pollTrue(
			`(function(){
				var lazy=document.querySelector('img#e2e-img-lazy');
				var eager=document.querySelector('img#e2e-img-eager');
				if(!lazy||!eager)return false;
				if(lazy.getAttribute('loading')!=='lazy')return false;
				if(lazy.hasAttribute('fetchpriority'))return false;
				if(eager.getAttribute('loading')!=='eager')return false;
				if(eager.getAttribute('fetchpriority')!=='high')return false;
				return true;})()`,
			pollFast()...,
		),
	); err != nil {
		t.Fatalf("visualtest[image-priority]: loading/fetchpriority pairing: %v", err)
	}
}

// TestViewTransitionRulePlacement checks the CSS-side of the cross-document
// view-transition opt-in with a string-level contract: the @view-transition
// rule is a DEMO-site concern (it enables MPA cross-document transitions for
// the demo's own navigation) and must not leak into the library's copied CSS
// (templates/custom.css), which consumers ship verbatim.
func TestViewTransitionRulePlacement(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed: cannot locate the visualtest package dir")
	}

	repoRoot := filepath.Join(filepath.Dir(thisFile), "..")

	demoCSS, err := os.ReadFile(filepath.Join(repoRoot, "examples", "demo", "demo.css"))
	if err != nil {
		t.Fatalf("read demo.css: %v", err)
	}

	if !strings.Contains(string(demoCSS), "@view-transition") {
		t.Fatal("examples/demo/demo.css lost the @view-transition opt-in rule")
	}

	customCSS, err := os.ReadFile(filepath.Join(repoRoot, "templates", "custom.css"))
	if err != nil {
		t.Fatalf("read custom.css: %v", err)
	}

	if strings.Contains(string(customCSS), "@view-transition") {
		t.Fatal("templates/custom.css must NOT enable cross-document view transitions — it ships to every consumer")
	}
}
