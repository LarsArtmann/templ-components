package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// The forms pattern pack e2e suite: browser-level proof for the v1.14.0
// wire surfaces that shipped after the Phase-0 trust pin — FilterInput
// (debounced), FilterDropdown.Wire, the multi-step wizard, multipart upload,
// the GET search form, and DirtyGuard.
//
// The repo's own hardest lesson (string-proven ≠ browser-proven) applies to
// everything added after Phase 0; this suite closes that gap. It follows the
// wire_form_e2e_test.go pattern: library components composed into a dedicated
// page + endpoints under BOTH the self-hosted htmx runtime and the pinned
// Datastar bundle.

const (
	packFilterHTMXOut     = "#pack-filter-htmx-out"
	packFilterDatastarOut = "#pack-filter-datastar-out"

	packDropdownHTMXOut     = "#pack-dropdown-htmx-out"
	packDropdownDatastarOut = "#pack-dropdown-datastar-out"

	packWizardHTMXRegion     = "#pack-wizard-htmx-region"
	packWizardDatastarRegion = "#pack-wizard-datastar-region"

	packUploadHTMXOut     = "#pack-upload-htmx-out"
	packUploadDatastarOut = "#pack-upload-datastar-out"

	packSearchHTMXRegion     = "#pack-search-htmx-region"
	packSearchDatastarRegion = "#pack-search-datastar-region"

	packDirtyRegion = "#pack-dirty-region"

	packWizardEmailBad = "Enter an email address with a domain."
	packWizardNameBad  = "Enter your name."

	packUploadFileName = "e2e-upload.txt"
	packUploadFileSize = 2048
)

// packE2EServer serves the forms pattern pack E2E page, the pinned Datastar
// bundle, the compiled CSS, and the pack endpoints. Datastar callers get
// response-header targeting via wire.Handler — the production server-side
// contract; htmx callers carry their hx-target client-side.
func packE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	props := layout.DefaultPageProps()
	props.Title = "Forms Pattern Pack E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, werr := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return werr
	})

	filterHits := &atomic.Int64{}

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	mux.Handle("GET /api/pack/filter", wire.Handler(wire.PatchTarget{
		Selector: packFilterDatastarOut,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		q := r.URL.Query().Get("q")

		packWriteComponent(w, r, packFilterResults(q, int(filterHits.Add(1))))
	})))

	mux.Handle("GET /api/pack/dropdown", wire.Handler(wire.PatchTarget{
		Selector: packDropdownDatastarOut,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		packWriteComponent(w, r, packDropdownResult(r.URL.Query().Get("framework"), packTransportName(r)))
	})))

	mux.Handle("POST /api/pack/wizard", wire.Handler(wire.PatchTarget{
		Selector: packWizardDatastarRegion,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form body", http.StatusBadRequest)

			return
		}

		dialect := packDialect(r)
		step, _ := strconv.Atoi(r.PostFormValue("step"))
		switch step {
		case 0:
			email := strings.TrimSpace(r.PostFormValue("email"))
			if email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
				packWriteComponent(w, r, packWizardStep(dialect, 0, packWizardEmailBad))

				return
			}
			packWriteComponent(w, r, packWizardStep(dialect, 1, ""))
		case 1:
			if strings.TrimSpace(r.PostFormValue("name")) == "" {
				packWriteComponent(w, r, packWizardStep(dialect, 1, packWizardNameBad))

				return
			}
			packWriteComponent(w, r, packWizardStep(dialect, 2, ""))
		default:
			packWriteComponent(w, r, packWizardStep(dialect, 0, ""))
		}
	})))

	mux.Handle("POST /api/pack/upload", wire.Handler(wire.PatchTarget{
		Selector: packUploadDatastarOut,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		const maxBytes int64 = 8 << 20
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		if err := r.ParseMultipartForm(maxBytes); err != nil {
			http.Error(w, "invalid multipart body", http.StatusBadRequest)

			return
		}

		file, header, err := r.FormFile("attachment")
		if err != nil {
			packWriteComponent(w, r, feedback.InlineError("No attachment received — choose a file and try again."))

			return
		}
		defer file.Close()

		size, err := io.Copy(io.Discard, file)
		if err != nil {
			http.Error(w, "unreadable attachment", http.StatusBadRequest)

			return
		}

		packWriteComponent(w, r, packUploadResult(header.Filename, size, packTransportName(r)))
	})))

	mux.Handle("GET /api/pack/search", wire.Handler(wire.PatchTarget{
		Selector: packSearchDatastarRegion,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		q := r.URL.Query().Get("q")

		packWriteComponent(w, r, packSearchResult(packDialect(r), q))
	})))

	mux.Handle("POST /api/pack/dirty", wire.Handler(wire.PatchTarget{
		Selector: packDirtyRegion,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form body", http.StatusBadRequest)

			return
		}

		packWriteComponent(w, r, packDirtyRoundTrip(strings.TrimSpace(r.PostFormValue("project"))))
	})))

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := packE2EPage(props).Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// packDialect reports which transport carried the request.
func packDialect(r *http.Request) wire.Transport {
	if wire.IsDatastar(r) {
		return wire.TransportDatastar
	}

	return wire.TransportHTMX
}

// packTransportName names the dialect for verdict fragments.
func packTransportName(r *http.Request) string {
	if wire.IsDatastar(r) {
		return "datastar"
	}

	return "htmx"
}

// packWriteComponent renders a fragment or converts the failure to 500.
func packWriteComponent(w http.ResponseWriter, _ *http.Request, c templ.Component) {
	if err := c.Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// packWire builds the dialect's wire action for a pack endpoint. The htmx
// dialect carries its hx-target client-side; Datastar targeting is
// response-driven via wire.Handler.
func packWire(dialect wire.Transport, method wire.Method, url, htmxTarget string) *wire.Action {
	action := &wire.Action{
		Transport: dialect,
		Method:    method,
		URL:       url,
	}
	if dialect == wire.TransportHTMX {
		action.Target = htmxTarget
	}

	return action
}

// packFilterResults is the filter endpoint fragment: the echoed query plus
// the request counter — a single source of truth for both the swap and the
// exactly-once-per-burst debounce proof.
func packFilterResults(q string, hits int) templ.Component {
	return feedback.InlineSuccess(
		"Filter results for “" + q + "” (" + strconv.Itoa(hits) + " request)",
	)
}

// packDropdownResult is the dropdown endpoint fragment.
func packDropdownResult(value, transport string) templ.Component {
	return feedback.InlineSuccess("Picked " + value + " via " + transport + ".")
}

// packWizardStep renders one wizard step: indicator plus that step's wired
// form (step 0 = account email, step 1 = profile name, step 2 = done). The
// server owns advancement; this fragment is the bare step — the page owns
// the region wrappers, so inner-mode swaps never duplicate region ids.
func packWizardStep(dialect wire.Transport, step int, message string) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := feedback.StepIndicator(feedback.StepIndicatorProps{
			Steps:       []string{"Account", "Profile", "Done"},
			CurrentStep: step,
		}).Render(ctx, w); err != nil {
			return err
		}

		if step == 2 {
			return feedback.InlineSuccess("Wizard complete — the server validated each step.").Render(ctx, w)
		}

		formChildren := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := forms.Input(forms.InputProps{
				Name:  "step",
				Type:  forms.InputHidden,
				Value: strconv.Itoa(step),
			}).Render(ctx, w); err != nil {
				return err
			}

			if step == 0 {
				if message != "" {
					if err := forms.FieldError("pack-wizard-email", message).Render(ctx, w); err != nil {
						return err
					}
				}

				if err := forms.Input(forms.InputProps{
					Name:  "email",
					Type:  forms.InputEmail,
					Label: "Email",
				}).Render(ctx, w); err != nil {
					return err
				}
			} else {
				if message != "" {
					if err := forms.FieldError("pack-wizard-name", message).Render(ctx, w); err != nil {
						return err
					}
				}

				if err := forms.Input(forms.InputProps{
					Name:  "name",
					Label: "Full name",
				}).Render(ctx, w); err != nil {
					return err
				}
			}

			return display.Button(display.ButtonProps{
				Text:    "Next",
				Variant: display.ButtonPrimary,
				Size:    display.ButtonSizeSM,
				Type:    display.ButtonHTMLSubmit,
			}).Render(ctx, w)
		})

		return forms.Form(forms.FormProps{
			Action:     "/api/pack/wizard",
			Method:     forms.FormPost,
			BaseProps:  utils.BaseProps{Class: "space-y-3"},
			NoValidate: true,
			Wire:       packWire(dialect, wire.MethodPost, "/api/pack/wizard", packWizardHTMXRegion),
		}).Render(templ.WithChildren(ctx, formChildren), w)
	})

	return body
}

// packUploadResult is the upload endpoint fragment.
func packUploadResult(name string, size int64, transport string) templ.Component {
	return feedback.InlineSuccess(
		"Uploaded " + name + " (" + strconv.FormatInt(size, 10) + " bytes) via " + transport + ".",
	)
}

// packSearchResult is the search endpoint fragment: the echoed query plus a
// fresh wired form (the round-trip pattern — a correction resubmits in place).
func packSearchResult(dialect wire.Transport, q string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if q == "" {
			if err := feedback.InlineError("Empty query.").Render(ctx, w); err != nil {
				return err
			}
		} else if err := feedback.InlineSuccess(
			"GET received q=“" + q + "” via " + packTransportWireName(dialect) + ".",
		).Render(ctx, w); err != nil {
			return err
		}

		return packSearchForm(dialect, q).Render(ctx, w)
	})
}

// packTransportWireName names the dialect for search verdicts.
func packTransportWireName(dialect wire.Transport) string {
	if dialect == wire.TransportDatastar {
		return "datastar"
	}

	return "htmx"
}

// packSearchForm is the wired GET search form for one dialect.
func packSearchForm(dialect wire.Transport, q string) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := forms.Input(forms.InputProps{
			Name:        "q",
			Type:        forms.InputSearch,
			Label:       "Search",
			Value:       q,
			Placeholder: "Try: ada",
		}).Render(ctx, w); err != nil {
			return err
		}

		return display.Button(display.ButtonProps{
			Text:    "Search",
			Variant: display.ButtonSecondary,
			Size:    display.ButtonSizeSM,
			Type:    display.ButtonHTMLSubmit,
		}).Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return forms.Form(forms.FormProps{
			Action: "/api/pack/search",
			Method: forms.FormGet,
			BaseProps: utils.BaseProps{
				Class: "flex flex-wrap items-end gap-3",
			},
			Wire: packWire(dialect, wire.MethodGet, "/api/pack/search", packSearchHTMXRegion),
		}).Render(templ.WithChildren(ctx, children), w)
	})
}

// packDirtyForm is the guarded subscribe-style form for the DirtyGuard
// lifecycle: the round-trip response re-renders it, so the swapped-in form
// must start clean.
func packDirtyForm(value string) templ.Component {
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := forms.Input(forms.InputProps{
			Name:  "project",
			Label: "Project name",
			Value: value,
		}).Render(ctx, w); err != nil {
			return err
		}

		return display.Button(display.ButtonProps{
			Text:    "Save",
			Type:    display.ButtonHTMLSubmit,
			Size:    display.ButtonSizeSM,
			Variant: display.ButtonPrimary,
		}).Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return forms.Form(forms.FormProps{
			Method:     forms.FormPost,
			DirtyGuard: true,
			BaseProps:  utils.BaseProps{Class: "space-y-3"},
			Wire:       packWire(wire.TransportHTMX, wire.MethodPost, "/api/pack/dirty", packDirtyRegion),
		}).Render(templ.WithChildren(ctx, children), w)
	})
}

// packDirtyRoundTrip is the dirty endpoint fragment: the verdict plus a
// fresh guarded form.
func packDirtyRoundTrip(value string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := feedback.InlineSuccess("Saved " + value + " via htmx.").Render(ctx, w); err != nil {
			return err
		}

		return packDirtyForm("").Render(ctx, w)
	})
}

// packFilterInput renders the debounced filter input for one dialect. The
// test shortens the shipped 300ms default to 100ms so bursts settle fast
// without changing what is proven.
func packFilterInput(dialect wire.Transport) templ.Component {
	target := utils.Ternary(dialect == wire.TransportDatastar, "", packFilterDatastarOut)

	return forms.FilterInput(forms.FilterInputProps{
		Name:        "q",
		Label:       "Filter",
		Placeholder: "Type to filter…",
		DebounceMS:  100,
		Wire:        packWire(dialect, wire.MethodGet, "/api/pack/filter", target),
	})
}

// packFilterOutRegion is the filter results region for one dialect.
func packFilterOutRegion(dialect wire.Transport) string {
	if dialect == wire.TransportDatastar {
		return packFilterDatastarOut
	}

	return packFilterHTMXOut
}

// packScopeID is the per-dialect section wrapper for a pane kind — the
// interaction scope for tests (the swap regions live inside it but hold
// only results, never the controls).
func packScopeID(kind string, dialect wire.Transport) string {
	return "#pack-" + kind + "-" + string(dialect) + "-scope"
}

// packE2EPage composes the full pack page: DirtyGuard script, both-dialect
// filter inputs + result regions, wired dropdowns, wizard regions, upload
// forms, search regions, and the guarded form.
func packE2EPage(props layout.PageProps) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := datastar.SDKScript(datastar.SDKScriptProps{
			BaseProps: utils.BaseProps{Nonce: "pack-e2e-nonce"},
			Src:       "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		if err := forms.DirtyGuard(forms.DirtyGuardProps{
			BaseProps: utils.BaseProps{Nonce: "pack-e2e-nonce"},
		}).Render(ctx, w); err != nil {
			return err
		}

		write := func(s string) error {
			_, err := io.WriteString(w, s)

			return err
		}

		if err := write(`<div class="p-4 grid grid-cols-2 gap-8 items-start">`); err != nil {
			return err
		}

		// Filter pane per dialect.
		for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
			outID := packFilterOutRegion(dialect)
			if err := write(`<section id="` + strings.TrimPrefix(packScopeID("filter", dialect), "#") + `"><h2 class="text-sm font-semibold mb-2">filter ` + string(dialect) + `</h2>`); err != nil {
				return err
			}
			if err := packFilterInput(dialect).Render(ctx, w); err != nil {
				return err
			}
			if err := write(`<div id="` + strings.TrimPrefix(outID, "#") + `" aria-live="polite"></div></section>`); err != nil {
				return err
			}
		}

		// Dropdown pane per dialect.
		for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
			outID := packDropdownDatastarOut
			target := ""
			if dialect == wire.TransportHTMX {
				outID = packDropdownHTMXOut
				target = outID
			}
			if err := write(`<section id="` + strings.TrimPrefix(packScopeID("dropdown", dialect), "#") + `"><h2 class="text-sm font-semibold mb-2">dropdown ` + string(dialect) + `</h2>`); err != nil {
				return err
			}
			if err := forms.FilterDropdown(forms.FilterDropdownProps{
				Name:  "framework",
				Label: "Framework",
				Options: []forms.SelectOption{
					{Value: "alpha", Label: "Alpha"},
					{Value: "beta", Label: "Beta"},
					{Value: "gamma", Label: "Gamma"},
				},
				Wire: packWire(dialect, wire.MethodGet, "/api/pack/dropdown", target),
			}).Render(ctx, w); err != nil {
				return err
			}
			if err := write(`<div id="` + strings.TrimPrefix(outID, "#") + `" aria-live="polite"></div></section>`); err != nil {
				return err
			}
		}

		// Wizard pane per dialect.
		for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
			regionID := packWizardHTMXRegion
			if dialect == wire.TransportDatastar {
				regionID = packWizardDatastarRegion
			}
			if err := write(`<section><h2 class="text-sm font-semibold mb-2">wizard ` + string(dialect) + `</h2><div id="` + strings.TrimPrefix(regionID, "#") + `">`); err != nil {
				return err
			}
			if err := packWizardStep(dialect, 0, "").Render(ctx, w); err != nil {
				return err
			}
			if err := write(`</div></section>`); err != nil {
				return err
			}
		}

		// Upload pane per dialect (results regions per dialect).
		for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
			outID := packUploadDatastarOut
			target := ""
			if dialect == wire.TransportHTMX {
				outID = packUploadHTMXOut
				target = outID
			}
			if err := write(`<section id="` + strings.TrimPrefix(packScopeID("upload", dialect), "#") + `"><h2 class="text-sm font-semibold mb-2">upload ` + string(dialect) + `</h2>`); err != nil {
				return err
			}
			uploadChildren := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if err := forms.FileInput(forms.FileInputProps{Name: "attachment", Label: "Attachment"}).Render(ctx, w); err != nil {
					return err
				}

				return display.Button(display.ButtonProps{
					Text:    "Upload",
					Variant: display.ButtonPrimary,
					Size:    display.ButtonSizeSM,
					Type:    display.ButtonHTMLSubmit,
				}).Render(ctx, w)
			})

			if err := forms.Form(forms.FormProps{
				Action:  "/api/pack/upload",
				Method:  forms.FormPost,
				Enctype: forms.FormEnctypeMultipart,
				BaseProps: utils.BaseProps{
					Class: "space-y-3",
				},
				Wire: packWire(dialect, wire.MethodPost, "/api/pack/upload", target),
			}).Render(templ.WithChildren(ctx, uploadChildren), w); err != nil {
				return err
			}
			if err := write(`<div id="` + strings.TrimPrefix(outID, "#") + `" class="mt-3" aria-live="polite"></div></section>`); err != nil {
				return err
			}
		}

		// Search pane per dialect (region wraps result + form: round-trip).
		for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
			regionID := packSearchHTMXRegion
			if dialect == wire.TransportDatastar {
				regionID = packSearchDatastarRegion
			}
			if err := write(`<section><h2 class="text-sm font-semibold mb-2">search ` + string(dialect) + `</h2><div id="` + strings.TrimPrefix(regionID, "#") + `" aria-live="polite">`); err != nil {
				return err
			}
			if err := packSearchForm(dialect, "").Render(ctx, w); err != nil {
				return err
			}
			if err := write(`</div></section>`); err != nil {
				return err
			}
		}

		// DirtyGuard pane (htmx only — the lifecycle is transport-agnostic).
		if err := write(`<section><h2 class="text-sm font-semibold mb-2">dirty guard</h2><div id="pack-dirty-region" aria-live="polite">`); err != nil {
			return err
		}
		if err := packDirtyForm("").Render(ctx, w); err != nil {
			return err
		}
		if err := write(`</div></section></div>`); err != nil {
			return err
		}

		return nil
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// packGate returns the readiness poll expression for a dialect.
func packGate(dialect wire.Transport) string {
	if dialect == wire.TransportDatastar {
		return `window.__dsReady===true`
	}

	return `document.readyState==='complete' && window.htmx!==undefined`
}

// packDialects enumerates both transports for table-driven subtests.
func packDialects() []wire.Transport {
	return []wire.Transport{wire.TransportHTMX, wire.TransportDatastar}
}

// setSelectValue sets a select's value and fires the change event (the
// trigger both dialects' dropdown wiring listens for).
func setSelectValue(ctx context.Context, region, sel, value string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(cctx context.Context) error {
		expr := `var s=document.querySelector('` + formSel(region, sel) + `'); s.value=` + jsString(value) +
			`; s.dispatchEvent(new Event('change',{bubbles:true})); s.value`
		var out string

		return chromedp.Evaluate(expr, &out).Do(cctx)
	})
}

// setValueQuiet sets an input's value WITHOUT firing events — used by the
// Enter-key test so the debounce never races the native submit.
func setValueQuiet(ctx context.Context, region, sel, value string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(cctx context.Context) error {
		expr := `var i=document.querySelector('` + formSel(region, sel) + `'); i.value=` + jsString(value) + `; i.value`
		var out string

		return chromedp.Evaluate(expr, &out).Do(cctx)
	})
}

// regionExistsExpr builds a predicate: the region contains an element
// matching sel.
func regionExistsExpr(region, sel string) string {
	return `!!document.querySelector('` + formSel(region, sel) + `')`
}

// fireInputBurst sets the value and fires n synchronous input events in one
// JS tick — a working debounce collapses the whole burst into one request.
func fireInputBurst(ctx context.Context, region, sel, value string, n int) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(cctx context.Context) error {
		expr := `var i=document.querySelector('` + formSel(region, sel) + `'); i.value=` + jsString(value) +
			`; for (var k=0;k<` + strconv.Itoa(n) + `;k++) { i.dispatchEvent(new Event('input',{bubbles:true})); } 'burst'`
		var out string

		return chromedp.Evaluate(expr, &out).Do(cctx)
	})
}

// beforeUnloadPrevented dispatches a synthetic cancelable beforeunload on
// window and reports whether the page's guard handler called preventDefault.
// A real dispatch through the real registered listener — it exercises the
// capture-phase marking, the WeakSet membership, and the unload check in a
// live browser without navigating away.
func beforeUnloadPrevented(ctx context.Context) (bool, error) {
	var prevented bool

	err := chromedp.Evaluate(
		`var e=new Event('beforeunload',{cancelable:true}); window.dispatchEvent(e); e.defaultPrevented`,
		&prevented,
	).Do(ctx)

	return prevented, err
}

// regionText reads a region's innerText.
func regionText(ctx context.Context, region string) (string, error) {
	var text string

	err := chromedp.Evaluate(`document.querySelector('` + region + `').innerText`, &text).Do(ctx)

	return text, err
}

// TestWireE2EFilterInputDebouncesAndSwaps proves the debounced filter input
// under both runtimes: a burst of input events collapses into exactly one
// request (the debounce), the response swaps into the results region, and a
// later keystroke fires exactly one more request.
func TestWireE2EFilterInputDebouncesAndSwaps(t *testing.T) {
	t.Parallel()

	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()

			srv := packE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			region := packFilterOutRegion(dialect)
			scope := packScopeID("filter", dialect)
			var ok bool

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				fireInputBurst(ctx, scope, `input[name="q"]`, "fltr", 3),
				chromedp.Poll(regionHasText(region, "Filter results for “fltr” (1 request)"), &ok),
				waitSwapSettled(),
				// A later keystroke fires exactly one more request.
				setFieldValue(ctx, scope, `input[name="q"]`, "templ"),
				chromedp.Poll(regionHasText(region, "Filter results for “templ” (2 request)"), &ok),
			); err != nil {
				t.Fatalf("%s filter input E2E: %v", dialect, err)
			}

			// The burst value survived into the echoed query of the FIRST
			// swap — already asserted via the needle above.
		})
	}
}

// TestWireE2EFilterDropdownWireSwaps proves FilterDropdown.Wire under both
// runtimes: changing the select fires the wired request and the region
// receives the verdict naming the picked value and the transport.
func TestWireE2EFilterDropdownWireSwaps(t *testing.T) {
	t.Parallel()

	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()

			srv := packE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			region := packDropdownDatastarOut
			if dialect == wire.TransportHTMX {
				region = packDropdownHTMXOut
			}
			scope := packScopeID("dropdown", dialect)

			var ok bool

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				setSelectValue(ctx, scope, `select[name="framework"]`, "beta"),
				chromedp.Poll(regionHasText(region, "Picked beta via "+string(dialect)+"."), &ok),
			); err != nil {
				t.Fatalf("%s filter dropdown E2E: %v", dialect, err)
			}
		})
	}
}

// TestWireE2EWizardStepsAdvances proves the server-owned step machine under
// both runtimes: an invalid email re-renders step 0 with the inline error,
// a valid one advances to the profile step, and a valid name completes the
// wizard — every response swapped into the same region.
func TestWireE2EWizardStepsAdvances(t *testing.T) {
	t.Parallel()

	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()

			srv := packE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			region := packWizardHTMXRegion
			if dialect == wire.TransportDatastar {
				region = packWizardDatastarRegion
			}

			var ok bool

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				// Step 0: empty email → inline error, no advance.
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(region, packWizardEmailBad), &ok),
				chromedp.Poll(regionExistsExpr(region, `input[name="email"]`), &ok),
				waitSwapSettled(),
				// Step 0: valid email → advance to profile.
				setFieldValue(ctx, region, `input[name="email"]`, "ada@example.com"),
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionExistsExpr(region, `input[name="name"]`), &ok),
				waitSwapSettled(),
				// Step 1: empty name → inline error.
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(region, packWizardNameBad), &ok),
				waitSwapSettled(),
				// Step 1: valid name → wizard complete.
				setFieldValue(ctx, region, `input[name="name"]`, "Ada Lovelace"),
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(region, "Wizard complete"), &ok),
			); err != nil {
				text, textErr := regionText(ctx, region)
				if textErr != nil {
					text = "<region read failed: " + textErr.Error() + ">"
				}

				t.Fatalf("%s wizard E2E: %v\nregion text: %s", dialect, err, text)
			}
		})
	}
}

// TestWireE2EUploadFileRoundTrip proves the multipart upload under both
// runtimes: submitting without a file shows the inline error, then a real
// file selected via the browser's file input uploads and the shared results
// region names the file, its size, and the transport.
func TestWireE2EUploadFileRoundTrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, packUploadFileName)
	content := strings.Repeat("a", packUploadFileSize)
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write upload fixture: %v", err)
	}

	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()

			srv := packE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			out := packUploadDatastarOut
			if dialect == wire.TransportHTMX {
				out = packUploadHTMXOut
			}
			scope := packScopeID("upload", dialect)

			var ok bool

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				// No file yet → the endpoint's inline error fragment.
				chromedp.Click(formSel(scope, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(out, "No attachment received"), &ok),
				waitSwapSettled(),
				// Pick a real file and upload it.
				chromedp.SetUploadFiles(formSel(scope, `input[type="file"]`), []string{filePath}),
				chromedp.Click(formSel(scope, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(
					out,
					"Uploaded "+packUploadFileName+" ("+strconv.Itoa(packUploadFileSize)+" bytes) via "+string(dialect)+".",
				), &ok),
			); err != nil {
				t.Fatalf("%s upload E2E: %v", dialect, err)
			}
		})
	}
}

// TestWireE2EGETSearchRoundTrip proves the wired GET form under both
// runtimes: the field travels as a query parameter, the region shows the
// verdict naming the transport, the submitted value survives into the fresh
// form, and a correction resubmits in place.
func TestWireE2EGETSearchRoundTrip(t *testing.T) {
	t.Parallel()

	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()

			srv := packE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			region := packSearchHTMXRegion
			if dialect == wire.TransportDatastar {
				region = packSearchDatastarRegion
			}

			var (
				ok        bool
				preserved string
			)

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				setFieldValue(ctx, region, `input[name="q"]`, "ada"),
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(region, "GET received q=“ada” via "+string(dialect)+"."), &ok),
				waitSwapSettled(),
				// The submitted value survived the round-trip re-render.
				chromedp.Evaluate(formValueExpr(region, `input[name="q"]`), &preserved),
				// A correction resubmits in place.
				setFieldValue(ctx, region, `input[name="q"]`, "grace"),
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(region, "GET received q=“grace” via "+string(dialect)+"."), &ok),
			); err != nil {
				t.Fatalf("%s search E2E: %v", dialect, err)
			}

			if preserved != "ada" {
				t.Fatalf("%s: submitted query not preserved across re-render; got %q", dialect, preserved)
			}
		})
	}
}

// TestWireE2EDirtyGuardLifecycle proves the unsaved-changes guard in a live
// browser: the script attaches, a clean form does not trigger beforeunload,
// typing marks the form dirty, the wired submit clears the flag, and the
// swapped-in form starts clean — then can be dirtied again.
func TestWireE2EDirtyGuardLifecycle(t *testing.T) {
	t.Parallel()

	srv := packE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var (
		ok       bool
		attached bool
	)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(wire.TransportHTMX), &ok),
	); err != nil {
		t.Fatalf("dirty guard setup: %v", err)
	}

	// The singleton attached (the page really executed the guard script).
	if err := chromedp.Evaluate(`window.tcDirtyGuardAttached===true`, &attached).Do(ctx); err != nil {
		t.Fatalf("read guard singleton: %v", err)
	}
	if !attached {
		t.Fatal("tcDirtyGuardAttached is not set — the guard script did not run")
	}

	assertPrevented := func(want bool, stage string) {
		t.Helper()

		got, err := beforeUnloadPrevented(ctx)
		if err != nil {
			t.Fatalf("%s: beforeunload probe: %v", stage, err)
		}
		if got != want {
			t.Fatalf("%s: beforeunload preventDefault = %v, want %v", stage, got, want)
		}
	}

	// Clean form: no unload prompt.
	assertPrevented(false, "clean form")

	// Typing marks the form dirty (capture-phase input listener).
	if err := setFieldValue(ctx, packDirtyRegion, `input[name="project"]`, "Graf Zeppelin").Do(ctx); err != nil {
		t.Fatalf("dirty the form: %v", err)
	}
	assertPrevented(true, "dirty form")

	// The wired submit dispatches submit, which clears the flag; the
	// response swaps in a fresh guarded form.
	if err := chromedp.Run(ctx,
		chromedp.Click(formSel(packDirtyRegion, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionHasText(packDirtyRegion, "Saved Graf Zeppelin via htmx."), &ok),
		waitSwapSettled(),
	); err != nil {
		t.Fatalf("submit dirty form: %v", err)
	}
	assertPrevented(false, "after submit")

	// The swapped-in form starts clean — and tracks dirt again.
	assertPrevented(false, "swapped-in form")
	if err := setFieldValue(ctx, packDirtyRegion, `input[name="project"]`, "Hindenburg").Do(ctx); err != nil {
		t.Fatalf("dirty the swapped-in form: %v", err)
	}
	assertPrevented(true, "swapped-in form dirty")
}

// TestWireE2EFilterInputEnterKeySubmitsNatively proves the documented
// Enter-key degradation: the wired input does not intercept Enter — the
// native GET form submits full-page with the field as a query parameter and
// the reloaded page stays interactive.
func TestWireE2EFilterInputEnterKeySubmitsNatively(t *testing.T) {
	t.Parallel()

	for _, dialect := range packDialects() {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()

			srv := packE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			var ok bool

			scope := packScopeID("filter", dialect)

			// Fire the native submit, then poll in a fresh Run: the
			// full-page navigation invalidates the current execution
			// context, so the location poll must run after it commits.
			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				setValueQuiet(ctx, scope, `input[name="q"]`, "enter-test"),
				chromedp.SendKeys(formSel(scope, `input[name="q"]`), kb.Enter),
			); err != nil {
				t.Fatalf("%s Enter-key submit: %v", dialect, err)
			}

			if err := chromedp.Run(ctx,
				chromedp.Poll(`window.location.search.indexOf('q=enter-test')>-1 && document.readyState==='complete'`, &ok),
			); err != nil {
				var current string
				_ = chromedp.Location(&current).Do(ctx)
				t.Fatalf("%s Enter-key E2E: %v (location=%s)", dialect, err, current)
			}
		})
	}
}
