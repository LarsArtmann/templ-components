package visualtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/larsartmann/templ-components/display"
)

func TestZZDebugPopoverProbe(t *testing.T) {
	d := display.DefaultDropdownProps()
	d.Label = "Options"
	d.Nonce = "test-nonce"
	d.Items = []display.DropdownItem{{Text: "Edit", Href: "/edit"}}

	page, err := renderHTML(display.Dropdown(d), Options{Viewport: Viewport{Width: 480, Height: 360}})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel2 := context.WithTimeout(ctx, 30*time.Second)
	defer cancel2()

	probe := func(stage string) chromedp.Action {
		return chromedp.ActionFunc(func(c context.Context) error {
			var out string
			err := chromedp.Evaluate(`(()=>{
				const p=document.querySelector('[popover]');
				return JSON.stringify({attached:!!window.tcPopoverPositionAttached,
					open:p?!!p.matches(':popover-open'):null,
					left:p?p.style.left:null,
					top:p?p.style.top:null});
			})()`, &out).Do(c)
			fmt.Printf("PROBE[%s]: %s err=%v\n", stage, out, err)
			return err
		})
	}

	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(480, 360),
		chromedp.Navigate("about:blank"),
	); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, page)
	}))
	defer srv.Close()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL),
		chromedp.WaitVisible("#tc-root", chromedp.ByQuery),
		probe("after-waitroot"),
		chromedp.Sleep(50*time.Millisecond),
		probe("after-50ms"),
	); err != nil {
		t.Fatal(err)
	}

	if err := chromedp.Run(ctx, clickAction("#tc-root", ""), probe("after-click")); err != nil {
		t.Fatal(err)
	}

	if err := chromedp.Run(ctx, chromedp.Sleep(500*time.Millisecond), probe("after-500ms")); err != nil {
		t.Fatal(err)
	}
}
