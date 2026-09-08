# Recipe: Print/PDF-safe pages (A4, no broken sections)

Server-rendered pages that must survive "Save as PDF" and paper: what to
reach for in Tailwind v4 + templ, harvested from a production CV generator
that ships pixel-checked PDF exports.

## CSS primitives

| Need | Tailwind |
|------|----------|
| Never split a block across pages | `break-inside-avoid` (+ a custom `print-break-avoid` alias class if you also target a print stylesheet) |
| Keep a heading with its content | `break-after-avoid` on the heading |
| Start a section on a new page | `break-before-page` |
| Hide screen chrome | `print:hidden` on nav/toolbars/buttons |
| Force light ink | design the print page light-only; do not rely on `print-color-adjust` for brand backgrounds |

## Structure rules that actually prevent bugs

1. **Two layouts, one component tree.** Screen pages want dark mode, nav,
   and interactivity; print pages want A4 width, light-only tokens, zero
   JS. Render BOTH from the same section components (they take text/data
   props), but wrap them in different shells: the library's
   `layout.Base` for screen, a minimal self-owned shell (or
   `layout.Minimal`) with `@media print` geometry for PDF.
2. **`var(--width)` page width.** Drive the content column with a CSS
   variable so screen and PDF renders share one width contract:
   `class="max-w-[var(--width)] w-[var(--width)]"`.
3. **List items are atomic.** `break-inside-avoid` on every `<li>` that
   carries a heading + date row; headers get `break-after-avoid`. This
   kills the "orphan heading at page bottom" class of bugs.
4. **No client JS assumptions.** PDF renderers may run JS inconsistently;
   dates render server-side (`RelativeTime` with `AutoRefresh: false`,
   `<time datetime>` carries the absolute value), icons are inline SVG.
5. **Link fidelity.** External links keep `href` (PDF exporters usually
   linkify anchors); decorate with `print:text-black` when color ink is
   not guaranteed.

## Static export shell

For pre-rendered HTML→PDF pipelines (headless Chromium, wkhtmltopdf),
inline the compiled CSS into a standalone shell so the artifact has no
network dependency:

```templ
templ staticExportShell(lang, personName, cssContent, bodyHTML string) {
    <!DOCTYPE html>
    <html lang={ lang }>
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>CV - { personName }</title>
            @templ.Raw("<style>\n" + cssContent + "\n    </style>")
        </head>
        <body class="max-w-[var(--width)] w-[var(--width)]">
            @templ.Raw(bodyHTML)
        </body>
    </html>
}
```

`templ.Raw` is a deliberate trust boundary: only feed it CSS/HTML your own
pipeline produced.

## Checklist

- [ ] `break-inside-avoid` on atomic blocks; `break-after-avoid` on headings
- [ ] `print:hidden` on all screen-only chrome
- [ ] Light-only palette on print pages (no `.dark` dependency)
- [ ] Zero JS requirements; absolute dates server-rendered
- [ ] Shared section components, two shells (screen vs print)
- [ ] Golden/PDF diff test per renderer version (byte-exact PDFs are
      compiler-version-pinned — see CV's Typst pinning discipline)
