# File Uploads Over Wire (Both Transports)

File uploads work over both HTMX and Datastar with one wired `forms.Form`.
The pieces:

1. `Enctype: forms.FormEnctypeMultipart` on the form (a typed enum added for
   this pattern).
2. A `forms.FileInput` field.
3. `Wire` on the form — htmx serializes the FormData natively; the pinned
   Datastar v1.0.3 bundle sends a **FormData body exactly when enctype is
   multipart** (bundle-verified, see `docs/datastar-runtime-facts.md`).

```templ
@forms.Form(forms.FormProps{
	Action:  "/api/attachments",
	Method:  forms.FormPost,
	Enctype: forms.FormEnctypeMultipart,
	Wire: &wire.Action{
		Method: wire.MethodPost,
		URL:    "/api/attachments",
		Target: "#upload-result", // htmx only; Datastar is response-driven
	},
}) {
	@forms.FileInput(forms.FileInputProps{
		Name:  "attachment",
		Label: "Attachment",
	})
	@display.Button(display.ButtonProps{
		Text:    "Upload",
		Variant: display.ButtonPrimary,
		Type:    display.ButtonHTMLSubmit,
	})
}
<div id="upload-result" aria-live="polite"></div>
```

## The endpoint (one handler, both dialects)

```go
mux.HandleFunc("POST /api/attachments", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	r.Body = http.MaxBytesReader(w, r.Body, 8<<20) // cap request size
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		// 200 OK + error fragment — the validation rule (see
		// docs/recipes/server-side-validation.md). 4xx does not swap
		// under either runtime's defaults.
		renderFragment(w, uploadErrorFragment("File too large (max 8 MiB)."))

		return
	}

	file, header, err := r.FormFile("attachment")
	if err != nil {
		renderFragment(w, uploadErrorFragment("No attachment received."))

		return
	}
	defer file.Close()

	size, _ := file.Seek(0, io.SeekEnd)
	_, _ = file.Seek(0, io.SeekStart)

	renderFragment(w, uploadSuccessFragment(header.Filename, size))
})
```

Wrap it in `wire.Handler(wire.PatchTarget{Selector: "#upload-result", Mode:
wire.PatchModeInner}, ...)` (or set the `Datastar-Selector` /
`Datastar-Mode` headers yourself) so the Datastar response patches the
result region; htmx targets it client-side via `hx-target`.

## What travels with the file

- **The CSRF hidden input** — `contentType: 'form'` serializes the whole
  form, hidden inputs included, under both runtimes. Nothing extra to wire.
- **The submitter button's name/value** — appended by both runtimes when
  the submission comes from a button (verified in the pinned bundle).
- **Other form fields** — free; FormData serializes siblings.

## Limits and gotchas

- **Size limits**: enforce with `http.MaxBytesReader` before
  `ParseMultipartForm`; render the overflow as a 200 OK error fragment, not
  a 413 (a 413 will not swap under either runtime's default handling).
- **Accept filtering**: `FileInputProps.Accept` (e.g. `image/*`) is a
  client-side hint only — re-validate the content type server-side.
- **No-JS**: the form degrades to a native multipart POST; keep the
  endpoint a plain form handler and the pattern works without either
  runtime.
- **Enter key**: submitting via Enter inside a text field performs the
  no-JS full-page POST under Datastar if no submit expression intercepts
  the `submit` event — prefer a real submit button for uploads.

Live demo: the wire page's "File upload" card
(`examples/demo/wire_demo.templ`, `/api/wire/upload`).
