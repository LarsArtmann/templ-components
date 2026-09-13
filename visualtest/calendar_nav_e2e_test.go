package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// calendarNavQuery is the month-nav endpoint payload — decoded via
// wire.DecodeForm (GET query), the same decoder consumers use.
type calendarNavQuery struct {
	Year  int `form:"year"`
	Month int `form:"month"`
}

// TestWireE2ECalendarMonthNav proves the Calendar MonthNav wiring end-to-end
// in a real browser, both transports: clicking the next-month arrow fetches
// the substituted {year}/{month} URL and patches the calendar region with
// the next month's grid — no full page load (#157, D3 rule: string-proven
// ≠ browser-proven).
func TestWireE2ECalendarMonthNav(t *testing.T) {
	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			srv := calendarNavServer(t, dialect)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			var ok bool

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
			); err != nil {
				t.Fatalf("%s calendar nav setup: %v", dialect, err)
			}

			if err := chromedp.Run(ctx,
				chromedp.Click(`a[aria-label="Next month"]`, chromedp.ByQuery),
				chromedp.Poll(`document.querySelector('#cal-nav-region h3')?.textContent.includes('August')`, &ok),
			); err != nil || !ok {
				t.Fatalf("%s next-month click did not patch the calendar to August: %v (ok=%v)", dialect, err, ok)
			}

			if err := chromedp.Run(ctx,
				chromedp.Click(`a[aria-label="Previous month"]`, chromedp.ByQuery),
				chromedp.Poll(`document.querySelector('#cal-nav-region h3')?.textContent.includes('July')`, &ok),
			); err != nil || !ok {
				t.Fatalf("%s prev-month click did not patch the calendar back to July: %v (ok=%v)", dialect, err, ok)
			}
		})
	}
}

// calendarNavServer serves a minimal page with one wired Calendar plus the
// month fragment endpoint. One endpoint serves both dialects via
// wire.Handler (Datastar targeting comes from response headers).
func calendarNavServer(t *testing.T, dialect wire.Transport) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	mux.Handle("GET /api/calendar", wire.Handler(wire.PatchTarget{
		Selector: "#cal-nav-region",
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		query, err := wire.DecodeForm[calendarNavQuery](r)
		if err != nil {
			http.Error(w, "invalid query", http.StatusBadRequest)

			return
		}

		if query.Year == 0 || query.Month == 0 {
			query.Year, query.Month = 2026, 7
		}

		if err := wiredCalendar(dialect, query.Year, query.Month).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := calendarNavPage(r.Context(), dialect).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// wiredCalendar renders the Calendar with the MonthNav action in the given
// dialect (placeholders substituted per direction by the component).
func wiredCalendar(dialect wire.Transport, year, month int) templ.Component {
	props := forms.DefaultCalendarProps()
	props.Year = year
	props.Month = time.Month(month)
	props.BaseProps = utils.BaseProps{ID: "cal-nav"}
	props.MonthNav = &wire.Action{ //nolint:exhaustruct // URL is the wiring surface
		Transport: dialect,
		URL:       "/api/calendar?year={year}&month={month}",
	}

	return forms.Calendar(props)
}

// calendarNavPage composes the shell via layout.Base (self-hosted htmx by
// default) + the Datastar bundle when the dialect needs it, with the wired
// calendar inside its patch region.
func calendarNavPage(ctx context.Context, dialect wire.Transport) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if dialect == wire.TransportDatastar {
			if err := datastar.SDKScript(datastar.SDKScriptProps{
				BaseProps: utils.BaseProps{Nonce: "cal-e2e"},
				Src:       "/datastar.js",
			}).Render(ctx, w); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, `<div id="cal-nav-region">`); err != nil {
			return err
		}

		if err := wiredCalendar(dialect, 2026, 7).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})

	pageProps := layout.DefaultPageProps()
	pageProps.Title = "Calendar nav e2e"

	_ = ctx // shell uses its own render context via Base

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(pageProps).Render(templ.WithChildren(ctx, body), w)
	})
}
