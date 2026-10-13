# ADR-0046: Button Children Slot (rich button content without a variant explosion)

**Date:** 2026-10-13
**Status:** Accepted
**Decides:** how consumers put non-default content (a spinner + label swap, an
icon-text stack, a badge) inside `display.Button`.
**Context issue:** #405 (cqrs-htmx's hand-copied Primary/Secondary+SM class
strings around `htmx.LoadingButton` content — three copies that drifted from
the library's classes and broke silently on every Button restyle).

## Decision

`ButtonProps.Children templ.Component` — when set, it REPLACES the default
content (Icon + Text) inside the element. The library keeps owning the
element: variants, sizes, disabled state, `AriaLabel`, HTMX attrs, and
`Wire` all still render from props. Consumer-owned: the a11y of the inner
markup (the element still carries `AriaLabel` when set).

## Rejected alternative: a `LoadingButton`-shaped variant

A dedicated `ButtonLoadingState{Default, Loading string, Spinner …}` prop
pair would cover the exact cqrs-htmx case, but:

1. It hard-codes ONE rich-content pattern into the Button API; the next
   consumer (icon + badge, icon stack) needs another prop pair — the
   variant-explosion path the library rejects everywhere else.
2. `htmx.LoadingButton` ALREADY owns the label-swap + spinner composition
   (its own component with its own contract); a Button-level twin would
   duplicate its state machine.

The Children slot composes with the existing piece instead of re-implementing
it: `ButtonProps{Children: htmx.LoadingButton("Save", "Saving…", spinner)}`.

## Consequences

- Rich content is consumer-a11y-owned (documented on the field).
- The goldens pin the slot path with a real child component; the `tc` mirror
  and demo showcase entries carry the pattern (display demo, Buttons section).
- Live proof: cqrs-htmx `adminui/members.templ` renders both submit buttons
  through the slot (was the three-copy class-string site).
- Same trade applies to `TableRow.Class` (#406): an escape-hatch field with a
  documented reason beats a prop explosion for one-off row styling.
