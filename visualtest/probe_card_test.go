package visualtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
)

func TestProbeCardLayout(t *testing.T) {
	card := display.DefaultCardProps()
	card.Title = "Monthly revenue"

	page, err := renderHTML(display.Card(card), Options{})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := newTab(t)
	defer cancel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, page)
	}))
	defer srv.Close()

	var out string
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 800),
		chromedp.Navigate(srv.URL),
		chromedp.WaitVisible("#tc-root", chromedp.ByQuery),
		chromedp.Evaluate(`JSON.stringify((() => {
			const boxes = [];
			const walk = (el, depth) => {
				const r = el.getBoundingClientRect();
				const cs = getComputedStyle(el);
				boxes.push({d: depth, tag: el.tagName, cls: (el.className||'').toString().slice(0,50),
					w: Math.round(r.width), h: Math.round(r.height), disp: cs.display});
				if (depth < 4) [...el.children].forEach(ch => walk(ch, depth+1));
			};
			walk(document.getElementById('tc-root'), 0);
			return boxes;
		})())`, &out),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("PROBE:", out)
	_ = os.WriteFile("/tmp/probe.json", []byte(out), 0o644)
}
