# How to Write a Superb Components UI/UX Library for the Web

> Research synthesis from shadcn/ui, templui, HATEOAS/hypermedia theory, WAI-ARIA
> accessibility practices, and design-token architecture.
>
> **Date:** 2026-07-05 · **Audience:** templ-components maintainers

---

## TL;DR — The 10 Pillars

| #   | Pillar                               | One-sentence thesis                                                                                                                               |
| --- | ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Distribution model**               | Decide explicitly: npm import (templ-components), copy-paste CLI (shadcn), or hybrid (templui) — each has a different ownership/upgrade tradeoff. |
| 2   | **HATEOAS-first**                    | HTML is the source of truth, not a dumb view layer. The server encodes available actions in the hypermedia itself.                                |
| 3   | **Type-safe variant system**         | Variants are typed enums + map lookups (Go) or CVA (TS), never raw strings — make invalid states unrepresentable.                                 |
| 4   | **Design-token theming**             | Three-layer token stack (primitive → semantic → component) via CSS variables; consumers retheme without touching Go/JS.                           |
| 5   | **Compound composition**             | Complex components (Dialog, Card) split into sub-parts that compose freely, mirroring the DOM structure.                                          |
| 6   | **Accessibility by construction**    | Native HTML first, ARIA only when needed, APG keyboard patterns mandatory, motion-reduce on every transition.                                     |
| 7   | **CSP-safe by default**              | Every inline script carries a nonce; no `eval`, no inline handlers, no `javascript:` URLs.                                                        |
| 8   | **Progressive-enhancement JS**       | JavaScript reads state from `data-*` attributes and HTML, not from a parallel client-side state store; idempotent across HTMX re-renders.         |
| 9   | **Multi-layered testing**            | Golden snapshots + a11y assertions + BDD behaviour specs + edge-case coverage — not assertion-only tests.                                         |
| 10  | **Beautiful defaults, zero lock-in** | Components look great out of the box but every visual decision is overridable via CSS variables, class props, or attrs.                           |

---

## 1. Distribution Models Compared

This is the single most consequential architectural decision. Three proven models exist:

### Model A — NPM / Go Module Import (templ-components)

```
go get github.com/larsartmann/templ-components@v0.7.0
import "github.com/larsartmann/templ-components/display"
```

- **Ownership:** Library owns the code; consumer gets updates via version bump.
- **Customization:** Props structs, `Class` field, CSS variables, `Attrs` map. Consumer
  never edits library source.
- **Upgrade story:** `go get -u` — atomic, versioned, semver-guaranteed.
- **Critical requirement:** Generated `*_templ.go` files MUST be committed (the Go module
  proxy serves source as-is; it does not run `templ generate`).
- **Best for:** Teams that want a stable, versioned dependency and are happy to override
  via the provided extension points.

### Model B — Copy-Paste CLI (shadcn/ui)

```
npx shadcn@latest add dialog   # copies source INTO your project
```

- **Ownership:** Consumer owns the code outright after scaffolding.
- **Customization:** Edit the component source directly — no wrappers, no overrides.
- **Upgrade story:** Re-run `add` to pull updates, or stay frozen. Consumer controls when.
- **Registry system:** JSON schema (`registry.json`) + CLI. Anyone can host a custom
  registry. Item types: `registry:ui`, `registry:block`, `registry:theme`, `registry:style`.
- **shadcn's thesis:** "This is not a component library. It is how you build your component
  library." Open code > opaque packages. AI can read, understand, and generate components.
- **Best for:** Teams that want full control and are willing to own maintenance.

### Model C — Hybrid (templui)

templui supports BOTH: import directly OR use their CLI to copy components in.

```
# Option 1: import
import "github.com/templui/templui/components/button"

# Option 2: CLI copy-paste (shadcn-style)
templui add button
```

- **Tradeoff:** Maximum flexibility but doubles the maintenance surface (the import path
  and the copied source can drift).

### Verdict for templ-components

The Go-module model (A) is correct for this library because:

1. Go's module system makes versioned imports trivial and reliable.
2. `utils.Class()` (tailwind-merge-go) gives consumers real override power without
   source edits — closing the gap that made shadcn invent copy-paste.
3. The multi-module workspace (svg, utils, icons, errorpage as independent modules)
   gives consumers granular adoption (icons-only, for example) that copy-paste can't
   match.
4. Committing `*_templ.go` is the Go-library equivalent of shipping compiled artifacts —
   it's the right contract for a publishable package.

**Actionable:** Document the "override without forking" story prominently. This is the
key counter-argument to "why not just copy-paste like shadcn?" The answer: `utils.Class()`

- CSS `@theme` variables + `Attrs` map + `templ.Component` slots give you 95% of the
  customization power without owning the maintenance burden.

---

## 2. HATEOAS-First: HTML as the Source of Truth

> "A natural hypermedia such as HTML is a practical necessity for building RESTful systems."
> — [htmx.org/essays/hateoas](https://htmx.org/essays/hateoas/)

### The Core Argument

HATEOAS (Hypermedia as the Engine of Application State) means: **the server encodes all
available actions directly in the HTML response.** The client (browser) needs no prior
knowledge of the application — it just renders hypermedia.

| Aspect                | Hypermedia (HTML)                                   | JSON API + SPA                                          |
| --------------------- | --------------------------------------------------- | ------------------------------------------------------- |
| **Available actions** | Encoded in the response (`<a>`, `<form>`, `hx-*`)   | Client must hardcode URLs/methods                       |
| **State transitions** | Server controls by changing the hypermedia          | Client maintains its own state machine                  |
| **Decoupling**        | Server evolves independently                        | Coordinated server+client updates required              |
| **Self-describing**   | `<form method="post" action="...">` says everything | `{"status": "overdrawn"}` says nothing about what to do |

### Implications for a Component Library

1. **Server-rendered HTML is a complete UI strategy.** A component library that renders
   HTML on the server gets HATEOAS benefits for free — every response is self-describing.

2. **Prefer native HTML primitives.** `<details>`/`<summary>` for disclosure, `<form>` for
   submit, `<a>` for navigation, `<dialog>` for modals. These ARE hypermedia controls.

3. **JavaScript enhances, never replaces.** When JS is needed (modal focus trap, combobox
   filter), it reads state FROM HTML attributes (`data-tc-*`, `datetime`, `aria-*`) —
   progressive enhancement, not SPA replacement.

4. **HTMX extends HTML as hypermedia.** `hx-get`, `hx-post`, `hx-swap` attributes turn any
   element into a hypermedia control that triggers server-side state transitions and
   receives new hypermedia in return. This is HATEOAS in practice.

5. **Avoid client-side state stores.** When the client maintains application state
   independently (React Context, Redux, etc.) and fetches JSON to drive rendering, it
   reintroduces all the coupling that HATEOAS eliminates.

### What This Looks Like in Practice

```
User clicks "Delete" button
  → hx-delete="/items/42" fires (hypermedia control)
  → Server processes deletion
  → Server returns updated item list HTML (new hypermedia)
  → HTMX swaps the old list for the new one
  → The new HTML naturally omits the "Delete" action for the removed item
  → No client-side state synchronization needed
```

Contrast with a JSON API approach:

```
User clicks "Delete"
  → Client sends DELETE /api/items/42
  → Server returns {"success": true}
  → Client must KNOW to remove item 42 from its local state array
  → Client must KNOW to re-render the list component
  → Client must KNOW the item no longer has actions available
  → All of this knowledge is out-of-band coupling
```

---

## 3. Type-Safe Variant Systems

### The Go Pattern (templ-components)

```go
// 1. Typed enum — closed variant set
type BadgeType string
const (
    BadgePrimary BadgeType = "primary"
    BadgeSuccess BadgeType = "success"
    BadgeWarning BadgeType = "warning"
    BadgeDanger  BadgeType = "danger"
    BadgeDefault BadgeType = "default"  // always provide a Default
)

// 2. Map lookup with explicit fallback — graceful degradation
var badgeStyleMap = map[BadgeType]string{
    BadgePrimary: "bg-blue-100 text-blue-800 ...",
    BadgeSuccess: "bg-green-100 text-green-800 ...",
    // ...
}

func badgeClass(t BadgeType) string {
    return utils.Lookup(badgeStyleMap, t, badgeStyleMap[BadgeDefault])
}
```

**Key rules from the templ-components playbook:**

- Pure class/style data → map + `utils.Lookup(map, key, fallback)`.
- Structural DOM differences (which markup to emit) → `if`-branch in the template, NOT a map.
- Enum validation → map+fallback (graceful); never panic in render code.
- Size constants use the `ComponentSize[SM|MD|LG]` suffix pattern.

### The TypeScript Pattern (shadcn/ui via CVA)

```typescript
const buttonVariants = cva("inline-flex items-center justify-center ...", {
  variants: {
    variant: {
      default: "bg-primary text-primary-foreground hover:bg-primary/90",
      destructive: "bg-destructive text-white ...",
      outline: "border bg-background shadow-xs hover:bg-accent ...",
    },
    size: {
      default: "h-9 px-4 py-2",
      sm: "h-8 rounded-md gap-1.5 px-3",
      lg: "h-10 rounded-md px-6",
      icon: "size-9",
    },
  },
  defaultVariants: { variant: "default", size: "default" },
});
```

CVA provides:

- **Type-safe variant combinations** — invalid props are compile-time errors.
- **Compound variants** — styles applied when multiple conditions match simultaneously.
- **Pre-computed class strings** — zero runtime concatenation.

### The templui Pattern (Go, switch-based)

```go
func (b Props) variantClasses() string {
    switch b.Variant {
    case VariantDestructive:
        return "bg-destructive text-white ..."
    case VariantOutline:
        return "border bg-background ..."
    default:
        return "bg-primary text-primary-foreground ..."
    }
}
```

**Note:** templui uses `switch` statements where templ-components mandates map+lookup.
The map approach is preferred because:

1. It's data, not code — easier to audit and extend.
2. The `utils.Lookup` fallback pattern makes graceful degradation explicit.
3. Maps can't accidentally fall through or miss a case.

---

## 4. Design-Token Theming Architecture

### The Three-Layer Token Stack

Every mature design system converges on this architecture:

```
Layer 3: Component Tokens     --button-bg → var(--color-accent)
               ↓ references
Layer 2: Semantic Tokens      --color-accent → var(--blue-500)
               ↓ references
Layer 1: Primitive Tokens     --blue-500 → #3b82f6  (the ONLY layer with raw values)
```

| Layer         | Named after            | Contains                 | Example                            |
| ------------- | ---------------------- | ------------------------ | ---------------------------------- |
| **Primitive** | What they ARE          | Raw values (hex, px)     | `--blue-500: #3b82f6`              |
| **Semantic**  | What they DO           | References to primitives | `--color-accent: var(--blue-500)`  |
| **Component** | What element uses them | References to semantic   | `--button-bg: var(--color-accent)` |

**Critical rule:** Each layer only references the layer directly below. Nothing bleeds
across boundaries. This means:

- **Brand color change** → update one primitive (`--blue-500: #6366f1`).
- **Theme switch** → reassign semantic tokens (dark mode remaps `--color-bg` to a different primitive).
- **Component redesign** → update component tokens only.

### shadcn/ui's Token System

shadcn uses **OKLCH** color values (perceptually uniform) in semantic token pairs:

```css
:root {
  --background: oklch(1 0 0);
  --foreground: oklch(0.145 0 0);
  --primary: oklch(0.205 0 0);
  --primary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.97 0 0);
  --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --border: oklch(0.922 0 0);
  --ring: oklch(0.708 0 0);
}
```

These map to Tailwind via `@theme inline`:

```css
@theme inline {
  --color-background: var(--background);
  --color-primary: var(--primary);
  /* ... */
}
```

So `bg-primary` in component code resolves to `var(--primary)` → `oklch(0.205 0 0)`.

### templ-components' Theming Approach

templ-components takes a **different but compatible** approach:

1. **Components emit standard Tailwind classes** (`bg-blue-600`, `text-gray-900`).
2. **Consumers override via `@theme` CSS variables:**
   ```css
   @theme {
     --color-blue-600: #4f46e5; /* indigo instead of blue */
   }
   ```
3. **Semantic aliases** live in `templ-components-theme.css` (`bg-tc-primary`, `text-tc-danger`).
4. **Dark mode** via class strategy: `@custom-variant dark (&:where(.dark, .dark *))`.

**Gap identified:** templ-components currently uses primitive Tailwind color names
(`blue-600`) directly in components, not semantic tokens (`primary`, `background`).
This means a consumer who wants to change the primary color must remap `--color-blue-600`
globally — which affects ALL uses of blue-600, not just "primary" uses.

**Recommendation:** Consider offering an optional semantic-token layer (like shadcn's
`--primary`/`--background`/`--accent`) that components can opt into. The current
`templ-components-theme.css` is a step in this direction. See §9 for details.

---

## 5. Compound Component Composition

### The Pattern

Complex components split into composable sub-parts that mirror the DOM structure:

```
Dialog (shadcn)          Dialog (templui)         Modal (templ-components)
├── DialogTrigger        ├── Trigger              ├── (single Modal component)
├── DialogContent        ├── Content              │   ├── trigger button
│   ├── DialogHeader     │   ├── Header           │   ├── backdrop
│   │   ├── DialogTitle  │   │   ├── Title        │   ├── content panel
│   │   └── Dialog...    │   │   └── Description  │   └── close button
│   └── DialogFooter     │   ├── Close
│   └── DialogClose      │   └── Footer
                         └── Script()
```

### Tradeoffs

| Approach         | shadcn/templui (compound)                                   | templ-components (monolithic)                |
| ---------------- | ----------------------------------------------------------- | -------------------------------------------- |
| **Flexibility**  | Maximum — compose any structure                             | Constrained to supported configs             |
| **Ergonomics**   | More verbose (many sub-components)                          | Simpler (one props struct)                   |
| **Type safety**  | Each sub-part has its own typed props                       | One props struct covers all                  |
| **HTML control** | Consumer owns the DOM tree                                  | Library owns the DOM tree                    |
| **HATEOAS fit**  | Better — consumer decides what hypermedia controls to embed | Library decides; consumer gets what they get |

**templui's compound pattern** is notable: `Dialog`, `Trigger`, `Content`, `Close`,
`Header`, `Footer`, `Title`, `Description` are all separate `templ` components with their
own props. The `Trigger` and `Close` support a `For` field for external triggering
(triggering a dialog from outside its compound tree).

**templ-components' monolithic pattern** is simpler for consumers but less flexible. For
example, `Modal` renders its own trigger button — you can't easily put an arbitrary
element as the trigger.

**Recommendation:** For NEW components where flexibility matters (e.g., a future Popover,
ContextMenu, HoverCard), consider the compound pattern. For simple, opinionated components
(Badge, Spinner, Skeleton), the monolithic pattern is fine.

---

## 6. Accessibility by Construction

### The Five Rules of ARIA ( distilled from W3C/MDN)

1. **Prefer native HTML.** A `<button>` is universally accessible. A `<div role="button">`
   requires you to manually implement keyboard handling, focus, and activation — and
   you'll never perfectly replicate what the browser gives for free.
2. **Don't change native semantics.** Nest semantic elements inside role containers; don't
   override their meaning.
3. **All interactive ARIA controls must be keyboard-accessible.** `role="button"` must
   respond to both Enter AND Space.
4. **Don't hide focusable elements.** Never `aria-hidden="true"` on a focusable element
   without also `tabindex="-1"`.
5. **"No ARIA is better than bad ARIA."** Pages with ARIA average 41% MORE errors than
   those without (WebAIM). Bad ARIA actively harms users.

### Mandatory Keyboard Patterns (from WAI-ARIA APG)

| Component            | Keyboard pattern                                                          |
| -------------------- | ------------------------------------------------------------------------- |
| **Dialog/Modal**     | Escape to close; focus trapped inside; focus restored to trigger on close |
| **Tabs**             | Arrow Left/Right between tabs; Tab enters panel                           |
| **Menu/Dropdown**    | Arrow keys navigate; Enter/Space activates; Escape closes                 |
| **Accordion**        | Enter/Space toggles; arrows move between headers                          |
| **Combobox/Listbox** | Arrow Up/Down moves selection; Home/End for first/last                    |
| **Tooltip**          | Escape dismisses; focus management on trigger                             |

### Motion Safety

Every transition must carry:

```css
motion-reduce:transition-none motion-reduce:duration-0
```

Every animation must carry:

```css
motion-reduce: animate-none;
```

### What "Accessible by Default" Means

An accessible-by-default library **encapsulates accessibility complexity** so consumers
get it for free. A well-built Dialog automatically:

- Sets `role="dialog"` and `aria-modal="true"`
- Moves focus into the dialog on open
- Traps focus while open (Tab cycles inside, never escapes)
- Restores focus to the trigger on close
- Closes on Escape

The consumer doesn't need to think about any of this. This is the core value proposition.

### templ-components' Accessibility Strengths

- `motion-reduce:*` on all transitions/animations (enforced by a11y tests)
- Focus trap + restore in Modal/Drawer
- `aria-label` propagation via BaseProps on all components
- Keyboard nav in Dropdown, Accordion, Tabs, Combobox
- `role="status"`/`role="alert"` on feedback components
- `utils.EnsureID()` for ARIA wiring (crypto/rand collision-safe IDs)

---

## 7. CSP-Safe by Default

Content Security Policy (CSP) is a defense-in-depth layer against XSS. A strict CSP
forbids: inline scripts without nonces, `eval()`, inline event handlers (`onclick`), and
`javascript:` URLs.

### The templ-components Approach

- **Every inline `<script>` carries `nonce={ props.Nonce }`.** No exceptions.
- **`layout.Script(nonce, src, attrs)`** for external scripts — auto-injects the nonce so
  it can never be forgotten.
- **No `eval()`, no inline handlers, no `javascript:` URLs.**
- **JS attachment is idempotent** via global singleton flags
  (`window.tcXxxAttached`) — safe across HTMX re-renders.

### templui's Approach (Gap)

templui ships **separate `.js` files** (e.g., `dialog.js`, `dialog.min.js`) that are
loaded externally. This works with CSP (external scripts are allowed) but:

1. Requires the consumer to serve the JS files from their own static asset path.
2. The `Script()` component uses `templ.NewOnceHandle()` for idempotent injection but
   doesn't show nonce handling in the examined source.
3. More moving parts (separate JS files to manage, minify, version, serve).

**templ-components' inline-script-with-nonce approach is simpler for consumers** because
there are no external assets to manage — the JS travels with the component.

---

## 8. Progressive-Enhancement JavaScript Strategy

### The Decision Ladder (from templ-components playbook)

Before writing ANY JavaScript, walk this ladder:

1. **Can native HTML do it?**
   - `<details>`/`<summary>` for disclosure/accordion
   - `<form>` for submit
   - `:target`/`:checked` for CSS-only toggles
   - `<dialog>` for native modals
   - If yes → **stop**. No JS needed.

2. **If JS is unavoidable**, write a single inline `<script nonce={ props.Nonce }>` that:
   - Guards with a global singleton flag (`window.tcXxxAttached`) so HTMX re-renders
     are idempotent (re-running the script must not double-bind handlers).
   - Uses event delegation on `document` where practical.
   - Handles Escape-to-dismiss, focus save/restore (`data-tc-prev-focus`), click-outside.
   - Reads state FROM HTML attributes (`data-tc-*`), never from a parallel client store.

3. **Escape IDs** interpolated into JS with `strconv.Quote()` to prevent XSS.

### JavaScript as Progressive Enhancement

The key insight from HATEOAS: JavaScript should **enhance** the hypermedia, not **replace**
it. This means:

- The component renders complete, functional HTML without JS (degrades gracefully).
- JS attaches behaviors that read state from the DOM (`data-tc-*`, `aria-*`, `datetime`).
- When HTMX re-renders a region, the singleton guard prevents double-binding.
- The server remains the source of truth for application state.

### Patterns observed

| Library              | JS idempotency                                 | CSP nonce                               | State source                 |
| -------------------- | ---------------------------------------------- | --------------------------------------- | ---------------------------- |
| **templ-components** | Global singleton flag (`window.tcXxxAttached`) | `nonce={ props.Nonce }` on every script | `data-tc-*` HTML attributes  |
| **templui**          | `templ.NewOnceHandle()`                        | External `.js` files (CSP-allowed)      | `data-tui-*` HTML attributes |
| **shadcn/ui**        | React re-render lifecycle (N/A for SSR)        | N/A (React manages)                     | React state/props            |

---

## 9. Identified Gaps & Opportunities for templ-components

Based on comparing all three libraries and the theory, here are concrete opportunities:

### Gap 1: Semantic Token Layer (High Impact)

**Current:** Components use primitive Tailwind names (`bg-blue-600`, `text-gray-900`).
**shadcn:** Components use semantic tokens (`bg-primary`, `text-foreground`, `bg-accent`).

**Problem:** A consumer who wants to change the "primary" color must remap
`--color-blue-600` globally, which affects ALL uses of blue-600, not just "primary" uses.

**Recommendation:** Introduce an optional semantic token layer:

```css
/* templ-components-theme.css (expanded) */
@theme {
  --color-tc-primary: var(--color-blue-600);
  --color-tc-primary-hover: var(--color-blue-700);
  --color-tc-background: var(--color-white);
  --color-tc-foreground: var(--color-gray-900);
  --color-tc-muted: var(--color-gray-100);
  --color-tc-muted-foreground: var(--color-gray-500);
  --color-tc-accent: var(--color-gray-100);
  --color-tc-destructive: var(--color-red-600);
  --color-tc-border: var(--color-gray-200);
}
```

Then components can use `bg-tc-primary` instead of `bg-blue-600`. Consumers retheme by
reassigning `--color-tc-primary` to any color — no collision with raw Tailwind utilities.

### Gap 2: Compound Components for Complex Widgets (Medium Impact)

**Current:** `Modal`, `Dropdown`, `Drawer` are monolithic — library owns the DOM tree.
**shadcn/templui:** Dialog, Popover, etc. are compound — consumer composes the structure.

**Problem:** Consumer can't put an arbitrary element as a Modal trigger, or compose a
Dialog with a non-standard header layout.

**Recommendation:** For new complex widgets (Popover, ContextMenu, HoverCard, Command
Palette), consider the compound pattern: `Trigger`, `Content`, `Close`, `Header`,
`Footer` sub-components with `data-tc-*` wiring.

### Gap 3: Registry / Copy-Paste Distribution (Exploratory)

**Current:** Go module import only.
**shadcn:** CLI + registry JSON for copy-paste distribution.
**templui:** Hybrid (import OR copy-paste).

**Observation:** Some Go developers prefer owning component source (especially for
one-off customization). A `templ-components add button` CLI that scaffolds a component
into the consumer's project could lower adoption friction for that audience.

**Risk:** Doubles the maintenance surface and can create drift. Probably not worth it
given that `utils.Class()` + CSS variables + `Attrs` already provide strong override
power without source ownership.

### Gap 4: Native `<dialog>` Element (Medium Impact)

**Current:** `Modal` uses div-based overlay with JS focus management.
**templui:** Uses native HTML5 `<dialog>` element with `::backdrop` pseudo-element.

**Benefits of native `<dialog>`:**

- Built-in focus trap (browser-managed, more reliable than JS)
- `::backdrop` pseudo-element for styling
- Top layer rendering (no z-index battles)
- `aria-modal` is implicit
- Better browser testing/compatibility

**Recommendation:** Consider migrating Modal/Drawer to native `<dialog>`/`<aside>` where
browser support allows (all modern browsers support `<dialog>` since 2022).

### Gap 5: Variadic Props (Low Impact, DX Improvement)

**Current:** `Component(props ComponentProps)` — props struct is required.
**templui:** `Button(props ...Props)` — props are optional (variadic), defaults to zero-value.

**templui's pattern:**

```go
templ Button(props ...Props) {
    {{ var p Props }}
    if len(props) > 0 {
        {{ p = props[0] }}
    }
    // ... use p
}
```

This lets consumers write `@button.Button()` instead of `@button.Button(button.Props{})`.
However, it has a downside: zero-value structs may have incorrect defaults unless the
constructor handles them, and it's less explicit.

**Recommendation:** Stick with the explicit props struct pattern. The `DefaultComponentProps()`
constructor pattern is clearer and more type-safe. The variadic pattern trades explicitness
for brevity — not worth the confusion cost.

### Gap 6: Missing Component Categories

Comparing component coverage:

| Category                  | templ-components                     | templui                                       | shadcn/ui    |
| ------------------------- | ------------------------------------ | --------------------------------------------- | ------------ |
| **Data display**          | 25                                   | 12                                            | 20+          |
| **Forms**                 | 16                                   | 12                                            | 15+          |
| **Feedback**              | 13                                   | 3                                             | 5            |
| **Navigation**            | 11                                   | 5                                             | 8            |
| **Layout**                | 5                                    | 3                                             | 4            |
| **Overlays**              | 4 (Modal, Drawer, Dropdown, Tooltip) | 6 (+Popover, HoverCard, Sheet)                | 8+           |
| **Charts**                | 0                                    | 1                                             | 1 (recharts) |
| **Date/Calendar**         | 1 (DatePicker input)                 | 3 (Calendar, DatePicker, TimePicker)          | 2            |
| **Data input (advanced)** | Combobox                             | Combobox, TagsInput, OTPInput, Rating, Slider | Same + more  |

**Notable gaps:** Popover, HoverCard, Slider, Rating, TagsInput, Carousel, Calendar
(full calendar grid), Chart components.

---

## 10. The "Superb" Checklist

A superb web UI library, synthesized from all sources:

### Architecture

- [ ] **Type-safe variants** — enums + map lookups, never raw strings for closed sets
- [ ] **Make invalid states unrepresentable** — typed enums, fallback values, zero panics in render
- [ ] **Compound composition** for complex widgets — sub-parts that mirror DOM structure
- [ ] **Single source of truth** — shared constants for SVG paths, class strings, style data
- [ ] **Composition over configuration** — `templ.Component` slots for extensibility

### Theming

- [ ] **Three-layer token stack** — primitive → semantic → component tokens
- [ ] **CSS variable theming** — consumers retheme without touching Go/JS
- [ ] **Dark mode** — class strategy, consistent color palette (no mixed `slate-*`/`gray-*`)
- [ ] **Beautiful defaults** — components look great out of the box
- [ ] **Zero lock-in** — every visual decision overridable

### Hypermedia Philosophy

- [ ] **HATEOAS-first** — HTML is the source of truth, server controls available actions
- [ ] **Native HTML preferred** — `<details>`, `<form>`, `<dialog>` over JS reimplementations
- [ ] **Progressive enhancement** — JS reads state from HTML, enhances rather than replaces
- [ ] **HTMX-native** — components emit `hx-*` attributes naturally

### Accessibility

- [ ] **Native HTML first** — `<button>` not `<div role="button">`
- [ ] **APG keyboard patterns** — every interactive widget follows WAI-ARIA Authoring Practices
- [ ] **Focus management** — trap in modals, restore on close, logical tab order
- [ ] **Motion safety** — `motion-reduce:*` on every transition/animation
- [ ] **ARIA correct** — roles, states, properties; "no ARIA is better than bad ARIA"
- [ ] **Screen reader text** — `.sr-only` where visual text isn't enough

### Security

- [ ] **CSP-safe by construction** — nonce on every inline script, no eval/handlers
- [ ] **XSS prevention** — `strconv.Quote()` for IDs in JS, `templ.SafeURL` for hrefs
- [ ] **No external assets required** — JS travels with the component (inline + nonce)

### Testing

- [ ] **Golden tests** — exact rendered HTML matches snapshots (with normalization)
- [ ] **A11y tests** — ARIA, roles, keyboard, motion-reduce assertions
- [ ] **BDD tests** — behaviour specs (user-visible, not markup)
- [ ] **Edge-case tests** — empty inputs, unknown enums, ID collisions
- [ ] **Example tests** — godoc `ExampleXxx()` compiles and renders
- [ ] **Coverage tests** — private helpers and branches

### Developer Experience

- [ ] **Zero-config defaults** — `DefaultComponentProps()` for meaningful non-zero defaults
- [ ] **Override without forking** — `Class` prop + `Attrs` map + CSS variables
- [ ] **Godoc examples** — every component has a runnable `ExampleXxx()`
- [ ] **Clear error messages** — graceful degradation, never panic in render path
- [ ] **Progressive disclosure** — simple API surface, deep docs for advanced cases

### Distribution

- [ ] **Versioned releases** — semver, one-commit release convention
- [ ] **Generated code committed** — `*_templ.go` files in the repo (library requirement)
- [ ] **Granular adoption** — multi-module workspace for partial adoption (icons-only, etc.)
- [ ] **CHANGELOG warm** — `[Unreleased]` always has entries

---

## Appendix A: Source Comparison Matrix

| Dimension             | templ-components                          | templui                             | shadcn/ui                           |
| --------------------- | ----------------------------------------- | ----------------------------------- | ----------------------------------- |
| **Language**          | Go + templ                                | Go + templ                          | TypeScript + React                  |
| **Rendering**         | Server-side (SSR)                         | Server-side (SSR)                   | Client-side (CSR/SSR)               |
| **CSS framework**     | Tailwind v4                               | Tailwind v3/v4                      | Tailwind v4                         |
| **Distribution**      | Go module                                 | Go module + CLI copy                | npm CLI copy-paste                  |
| **Props pattern**     | Explicit struct (embeds BaseProps)        | Variadic struct (no BaseProps)      | React props + CVA                   |
| **Variant lookup**    | Map + `utils.Lookup`                      | Switch statement                    | CVA (class-variance-authority)      |
| **Class merge**       | `utils.Class` (tailwind-merge-go)         | `utils.TwMerge` (tailwind-merge-go) | `cn()` (tailwind-merge + clsx)      |
| **Theming**           | Tailwind primitives + `@theme`            | Semantic tokens (`bg-primary`)      | Semantic tokens (OKLCH + `@theme`)  |
| **JS strategy**       | Inline + nonce + singleton guard          | External `.js` files + OnceHandle   | React lifecycle                     |
| **CSP**               | Safe by construction (nonce everywhere)   | External scripts (CSP-allowed)      | React (N/A for inline)              |
| **A11y**              | motion-reduce, focus trap, ARIA, keyboard | ARIA, keyboard                      | Built into Radix/Base UI primitives |
| **Component count**   | 82 + 101 icons                            | 43                                  | 50+                                 |
| **Testing**           | Golden + a11y + BDD + edge + example      | Unknown                             | Vitest + browser + E2E              |
| **HATEOAS alignment** | Strong (HTML-first, HTMX-native)          | Strong (HTML-first)                 | Weak (SPA model, JSON-driven)       |

---

## Appendix B: Key URLs

- [shadcn/ui](https://ui.shadcn.com/) — landing page
- [shadcn/ui docs](https://ui.shadcn.com/docs) — design philosophy, theming, CLI
- [shadcn-ui/ui repo](https://github.com/shadcn-ui/ui) — monorepo, registry system
- [templui repo](https://github.com/templui/templui) — Go templ UI library
- [HATEOAS essay](https://htmx.org/essays/hateoas/) — hypermedia as application state engine
- [WAI-ARIA APG](https://www.w3.org/WAI/ARIA/apg/) — authoritative accessible component patterns
- [Design Tokens spec](https://www.designtokens.org/) — DTCG token format specification
- [CVA](https://cva.style/) — class-variance-authority for type-safe variants
