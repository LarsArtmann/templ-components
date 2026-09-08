# Recipe: HTMX-loaded modal (native `<dialog>`)

The right way to "open a dialog whose content HTMX fetched" with
`display.Modal` — keeping the native focus trap, Escape handling, backdrop
dismiss, and top-layer rendering instead of regressing to a
`role="dialog"` `<div>` (which has none of those, and whose hand-rolled
dismiss `onclick=` is dead under a strict CSP).

## The pattern

1. Render the `display.Modal` shell server-side with `Open: false` — the
   `<dialog>` exists but stays closed.
2. Give the swap target an id **inside** the dialog body.
3. Trigger buttons `hx-get` the fragment into that target, then open.
4. Close with the library's `tcCloseModal(id)` helper (wired via data
   attributes and the page script — never inline `onclick`).

```templ
// Shell — render once, hidden.
@display.Modal(display.ModalProps{
    BaseProps: utils.BaseProps{
        ID:    "analysis-modal",
        Nonce: nonce,
        Attrs: templ.Attributes{"data-tc-modal": "analysis-modal"},
    },
    Title: "Analyze CV",
}) {
    <div id="analysis-modal-body"><!-- fragment lands here --></div>
}

// Trigger — fetch content, swap into the body, then open via the
// htmx:afterSwap listener or an hx-on-free delegated script:
//   document.body.addEventListener("htmx:afterSwap", function (evt) {
//     if (evt.detail.target.id === "analysis-modal-body") tcOpenModal("analysis-modal");
//   });
@display.Button(display.ButtonProps{
    BaseProps: utils.BaseProps{Attrs: templ.Attributes{
        "hx-get":     "/analyses/new-form",
        "hx-target":  "#analysis-modal-body",
        "hx-swap":    "innerHTML",
    }},
    Text: "New Analysis",
})

// Fragment — the server renders ONLY the inner content; close buttons
// carry data attributes dispatched by delegation:
//   <button type="button" data-tc-close-modal="analysis-modal">Cancel</button>
```

```js
// Delegated close (CSP-safe; survives HTMX swaps):
document.addEventListener("click", function (e) {
  var btn = e.target.closest("[data-tc-close-modal]");
  if (btn) tcCloseModal(btn.getAttribute("data-tc-close-modal"));
});
```

## Why not a div with `role="dialog"`

- No focus trap → Tab escapes to the page behind (WCAG 2.4.3 failure).
- No Escape-to-close, no `inert` background, no top layer (z-index wars).
- `::backdrop` and `@starting-style` open/close animations don't apply.
- Inline `onclick=` dismiss handlers are blocked by `script-src 'self'`.

## Checklist

- [ ] `display.Modal` shell rendered server-side, `Open: false`
- [ ] Fragment swaps into an inner `id`ed div, never over the `<dialog>`
- [ ] Open via `tcOpenModal(id)` after `htmx:afterSwap`
- [ ] Close via `tcCloseModal(id)` + delegated `data-tc-close-modal`
- [ ] Server fragment endpoints return 200-OK error fragments (htmx 2 does not swap 4xx/5xx by default — see `docs/recipes/server-rendered-htmx-error-feedback.md`)
