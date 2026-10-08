# Research Note: Direction Context Propagation (shadcn-templ's DirectionProvider) — Pattern Comparison

**Created:** 2026-10-08
**Status:** IMPLEMENTED 2026-10-08 (same day) — element-scoped `dir` reads shipped in menu nav, Tabs, and Carousel, pinned by `display.TestRTLDirectionReadsAreSubtreeScoped`. The full provider pattern is deliberately NOT adopted — revisit triggers listed below.

---

## Verification provenance

Per the verify-external-claims gate, every claim about shadcn-templ below was
verified against the primary source on 2026-10-08, not taken from docs prose
alone:

| Claim                                                                                                                                                                                                                                                                                                           | How verified                                                                     |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `utils.WithDirection(ctx, dir)` / `utils.Direction(ctx)` with `directionKey{}` ctx key, `"ltr"` fallback                                                                                                                                                                                                        | Sourcegraph on `github.com/axadrn/shadcn-templ`, `utils/shadcn-templ.go:205-218` |
| `DirectionProvider` renders no element, defaults `ltr`, injects ctx only                                                                                                                                                                                                                                        | `components/direction/direction.templ:21-31`                                     |
| "Set `dir` on `<html>` as well, the layout follows the attribute, not the provider" (their own caveat)                                                                                                                                                                                                          | Comment in `components/direction/direction.templ:19-20`                          |
| Client side: `useDirection(element)` reads the nearest `data-templ-direction` DOM marker, "ltr" otherwise; scripts read it **instead of** `getComputedStyle(...).direction` or `closest('[dir]')`; applies to tabs, toggle group, radio group, menubar, slider, scroll area, navigation menu, popup positioners | `plans/parity-components.md` §Task 16 (Direction entry)                          |
| Build-time RTL inliner (`applyRtlMapping`: physical→logical utilities, `rtl:` variants, `cn-rtl-flip`) with 106 upstream test cases                                                                                                                                                                             | `plans/parity-components.md` §Task 16                                            |

Context: shadcn-templ (formerly templUI) is a shadcn/ui port for templ that
wraps **Base UI** JS primitives. Their components ship real client-side
machinery (positioning engines, composite keyboard navigation), which is why
direction must reach their JavaScript explicitly.

---

## The pattern (three layers)

shadcn-templ's DirectionProvider ports React's context-based direction
propagation to server-rendered templ:

### 1. Server: Go render context

```go
type directionKey struct{}

func WithDirection(ctx context.Context, direction string) context.Context {
    return context.WithValue(ctx, directionKey{}, direction)
}

func Direction(ctx context.Context) string {
    if direction, ok := ctx.Value(directionKey{}).(string); ok && direction != "" {
        return direction
    }
    return "ltr"
}
```

### 2. Provider: renders no element, only a context boundary

```templ
templ DirectionProvider(props ...Props) {
    if p.Direction == "" {
        {{ p.Direction = DirectionLTR }}
    }
    {{ ctx = utils.WithDirection(ctx, string(p.Direction)) }}
    { children... }
}
```

Components below it branch server-side via `utils.Direction(ctx)`.

### 3. Client: Go ctx → DOM marker → nearest-ancestor JS lookup

Provider-consuming components render `data-templ-direction="rtl|ltr"` on
their composite roots; scripts resolve direction with `useDirection(element)`
(walk up to the nearest `data-templ-direction`, else `"ltr"`). This is the
DOM analog of React context — ambient state without prop drilling.

### Why it exists at all

Base UI's JS needs direction for **keyboard navigation** (arrow-key mapping)
and **popup positioning** (anchor side flips). CSS layout mirroring is
explicitly NOT the provider's job — their own source comment says: set
`dir="rtl"` on `<html>` too, "the layout follows the attribute, not the
provider." So even in the context-propagation design, CSS still runs on the
plain HTML attribute; the provider exists _only_ to feed JavaScript.

---

## Our current model (templ-components)

One source of truth: the DOM `dir` attribute.

1. **CSS:** logical properties only (`ms-`/`me-`/`ps-`/`pe-`/`inset-s-*`/
   `inset-e-*`/`text-start`/`border-s-`/`border-e-`). Physical utilities
   (`ml-`/`pl-`/`left-`/`text-left`/`border-l-`, and bare `start-`/`end-`,
   which compile to physical `left`/`right`) are banned at the source and
   enforced by `utils.TestRTLLogicalProperties`. Layout mirrors automatically
   wherever a consumer puts `dir="rtl"`, with zero per-component awareness.
2. **JS:** three event-time reads, subtree-scoped since 2026-10-08 — each
   resolves `(container.closest('[dir]') || document.documentElement)` with
   the component's own root as the container, pinned by
   `display.TestRTLDirectionReadsAreSubtreeScoped`:
   - `display/shared.go:384` — shared menu keyboard nav (Dropdown, ContextMenu)
   - `display/tabs.templ:222` — Tabs
   - `display/carousel.templ:159` — Carousel
3. **Directional icons:** `icons.IconRTL` renders `data-tc-dir-icon`;
   `templates/custom.css` mirrors via `[dir="rtl"] [data-tc-dir-icon] { transform: scaleX(-1); }`.
4. **Verification:** RTL option in `visualtest.AssertScreenshot` (pixel-level
   mirroring captures) + the logical-property scanner.

---

## Comparison

| Dimension                 | shadcn-templ                                                                                        | templ-components                                                                              |
| ------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Layout mirroring          | Physical utilities + build-time `applyRtlMapping` inliner (106 upstream test cases) to convert them | Logical properties banned-in at the source; no build step, no test corpus for the conversion  |
| Direction source of truth | **Two**: ctx marker (`data-templ-direction`) drives JS; `dir` attribute drives CSS                  | **One**: the `dir` attribute drives both                                                      |
| JS direction resolution   | Per-subtree correct (`useDirection` walks to nearest provider marker)                               | Per-subtree correct since 2026-10-08 (`closest('[dir]')` + `<html>` fallback, guard test)     |
| Server-side branching     | Yes — components branch markup on `utils.Direction(ctx)`                                            | No — components are direction-ignorant by construction                                        |
| API surface for consumers | A provider component to wrap content                                                                | None (just the standard `dir` attribute)                                                      |
| Failure mode when omitted | Provider defaults `ltr` — JS disagrees with an RTL layout (their own comment warns about this)      | None left (subtree gap fixed 2026-10-08; was: page-level `dir`, wrong arrows in RTL subtrees) |

---

## Learnings

### Adopted: element-scoped `dir` resolution in JS (TODO_LIST #372)

Shipped 2026-10-08 (TODO_LIST #372). The one real gap our model had: an RTL
widget embedded in an LTR page (or vice versa) got backwards ArrowLeft/Right
in menu nav, Tabs, and Carousel, because all three read `<html>` only. All
three now resolve:

```js
var isRtl = (container.closest('[dir]') || document.documentElement).getAttribute('dir') === 'rtl';
```

with `container` being the component root the handler already resolved (the
menu, the tab, the carousel region). The `dir` attribute IS the equivalent of
their `data-templ-direction` marker — the platform already gives us per-subtree
context propagation for free. We get their correctness without adding a
provider API. Regression guard: `display.TestRTLDirectionReadsAreSubtreeScoped`.

### Rejected: Go-side `utils.WithDirection`/`Direction` (for now)

Their ctx propagation exists because Base UI's JS primitives need the value in
~20 places (positioners, composites, sliders). We have exactly three
event-time reads. A ctx utility would add per-component plumbing and a second
direction source of truth with no consumer benefit today. Their
provider-defaults-to-`ltr` behavior is a documented footgun class: omit the
provider inside an RTL page and JS silently disagrees with the layout.

Their build-time `applyRtlMapping` inliner is likewise compensation for
Base UI's physical utilities. Our source-level ban is the stronger invariant:
the wrong pattern cannot be written, rather than being rewritten at build time.

### Back-pocket pattern: no-element provider + DOM marker

"Provider renders no element, stamps a marker attribute, JS walks to the
nearest marker" is an elegant Go-render-ctx → DOM bridge for ambient
per-subtree state. If we ever need scoped ambient configuration (density,
section-level tone, a subtree scope), this is the template — and it degrades
gracefully (marker absent → default), like our empty-URL wire actions.

---

## Revisit triggers

Reconsider the full provider pattern (Go ctx + marker) when ANY of these hold:

1. A component needs direction-dependent **markup** (not just mirrored
   layout) — e.g. a server-rendered arrow that must point the other way, where
   the `IconRTL` CSS-mirror trick doesn't apply.
2. We ship a JS positioning engine that resolves anchor sides per subtree
   (the reason Base UI needs the ctx) — e.g. if native Popover positioning
   ever proves insufficient and `popoverPositionJS` grows side logic.
3. Consumers report real mixed-direction pages where subtree `dir` attributes
   and the marker approach diverge in behavior (they shouldn't, post-#372 —
   that divergence is the signal the fix was wrong).
4. A mainstream templ peer library converges on ctx-based direction and
   consumers expect the API for interop.

---

## Sources

- shadcn-templ docs, Direction component: https://shadcn-templ.com/docs/components/direction
- shadcn-templ source: `github.com/axadrn/shadcn-templ` — `utils/shadcn-templ.go`,
  `components/direction/direction.templ`, `plans/parity-components.md` (§Task 16)
- Local: `display/shared.go` (menuKeyboardNavJS), `display/tabs.templ`,
  `display/carousel.templ`, `templates/custom.css` (`[data-tc-dir-icon]`),
  `utils` RTL guard tests, `visualtest` RTL capture option
