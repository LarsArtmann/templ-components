# How to Build a Superb Web UI Component Library

> A deep research synthesis for the templ-components project.
>
> Sources: shadcn/ui (philosophy + registry + blocks), Radix UI (headless primitives),
> React Aria (accessibility architecture), templui (Go templ competitor), HTMX HATEOAS
> theory, WAI-ARIA APG, Tailwind CSS v4, design-token architecture (DTCG spec),
> motion design best practices, and form-handling patterns for server-rendered apps.
>
> **Date:** 2026-07-05 · **Audience:** templ-components maintainers

---

## Table of Contents

1. [The First Principles](#1-the-first-principles)
2. [Distribution Models: Who Owns the Code?](#2-distribution-models-who-owns-the-code)
3. [Architecture: The Headless-Styled Divide](#3-architecture-the-headless-styled-divide)
4. [The Component API Contract](#4-the-component-api-contract)
5. [Variant Systems: Type-Safe Styling](#5-variant-systems-type-safe-styling)
6. [Theming: The Three-Layer Token Stack](#6-theming-the-three-layer-token-stack)
7. [HATEOAS: HTML as the Source of Truth](#7-hateoas-html-as-the-source-of-truth)
8. [Form Handling for Server-Rendered Apps](#8-form-handling-for-server-rendered-apps)
9. [Accessibility: Not a Feature, the Foundation](#9-accessibility-not-a-feature-the-foundation)
10. [Motion Design: Professional Animation](#10-motion-design-professional-animation)
11. [Progressive Enhancement & JavaScript](#11-progressive-enhancement--javascript)
12. [Security: CSP-Safe by Construction](#12-security-csp-safe-by-construction)
13. [Internationalization & RTL](#13-internationalization--rtl)
14. [Tailwind CSS v4: The CSS-First Revolution](#14-tailwind-css-v4-the-css-first-revolution)
15. [Testing Strategy](#15-testing-strategy)
16. [Documentation & Blocks](#16-documentation--blocks)
17. [Performance](#17-performance)
18. [Gap Analysis: templ-components vs the Field](#18-gap-analysis-templ-components-vs-the-field)
19. [The Master Checklist](#19-the-master-checklist)

---

## 1. The First Principles

Before any architecture, these are the non-negotiable principles that separate a superb
library from a mediocre one:

### 1.1 Behavior is universal; styling is local

The W3C ARIA Authoring Practices define exactly how a Dialog, Tabs, Combobox, or Select
should behave — keyboard patterns, ARIA attributes, focus management. This behavior is
**identical across every design system**. A button is still a button whether it's blue or
purple. This means:

> **Extract shared behavior; leave rendering to the consumer.**

Radix UI and React Aria both built their entire architecture on this insight. Radix
provides behavior-only primitives (zero styles); React Aria provides behavior-only hooks
(zero rendering). shadcn/ui then layers beautiful styling on top. The lesson: **the hard
part is the behavior, not the CSS.**

### 1.2 Make impossible states unrepresentable

Variants should be typed enums, not raw strings. Props should encode constraints — if a
field only accepts `"primary" | "secondary" | "destructive"`, the type system should
enforce it at compile time, not at runtime. A consumer who types `variant="primry"`
should get a **compile error**, not a silently broken component.

### 1.3 Accessibility is not optional

A component that renders but is unusable from a keyboard is not finished. A component
without `aria-label` propagation is not finished. A component without `motion-reduce:*`
on its animations is not finished. "Accessible by default" means the consumer gets
accessibility **for free** — they never have to think about ARIA roles, focus traps, or
keyboard patterns.

### 1.4 HTML is the source of truth

Not a dumb view layer. Not a serialization format. HTML is a self-describing hypermedia
that encodes available actions, application state, and interaction semantics — all in the
server's response. JavaScript enhances the hypermedia; it does not replace it.

### 1.5 Beautiful defaults, zero lock-in

Components must look great out of the box AND be overridable without forking. This means
CSS variables for theming, class-merge utilities for style overrides, attribute spreading
for arbitrary HTML, and slot components for composition.

---

## 2. Distribution Models: Who Owns the Code?

This is the single most consequential decision. Three proven models exist, each with
fundamentally different ownership economics.

### Model A — Versioned Module Import (templ-components, traditional libraries)

```
go get github.com/larsartmann/templ-components@v0.7.0
import display "github.com/larsartmann/templ-components/display"
```

| Dimension | Assessment |
|-----------|------------|
| **Ownership** | Library owns; consumer gets updates via `go get -u` |
| **Customization** | Props structs + `Class` field + CSS variables + `Attrs` map |
| **Upgrade story** | Atomic, semver-guaranteed, one command |
| **Critical requirement** | Generated `*_templ.go` MUST be committed (Go module proxy serves source as-is) |
| **Best for** | Teams that want stable, versioned dependencies |

**Why this works for templ-components:** Go's module system makes versioned imports
trivial and reliable. `utils.Class()` (tailwind-merge-go) gives consumers real override
power — closing the gap that made shadcn invent copy-paste. The multi-module workspace
(svg, utils, icons, errorpage) gives granular adoption that copy-paste cannot match.

### Model B — Copy-Paste CLI (shadcn/ui)

```bash
npx shadcn@latest add dialog   # copies source INTO your project
```

| Dimension | Assessment |
|-----------|------------|
| **Ownership** | Consumer owns the code outright after scaffolding |
| **Customization** | Edit the source directly — no wrappers, no `!important` |
| **Upgrade story** | Re-run `add` to pull updates, or stay frozen. Consumer controls when. |
| **Registry system** | JSON schema (`registry.json`) + CLI. Types: `registry:ui`, `registry:block`, `registry:theme` |
| **Best for** | Teams that want full control and accept maintenance ownership |

**shadcn's thesis:** *"This is not a component library. It is how you build your
component library."* Open code > opaque packages. AI can read, understand, and generate
components. The registry JSON format is open and self-hostable.

### Model C — Hybrid (templui)

templui supports both: `import "github.com/templui/templui/components/button"` OR
`templui add button` (copies source). Maximum flexibility but doubles the maintenance
surface (import path and copied source can drift).

### Verdict for templ-components

**Model A is correct.** The key is documenting the "override without forking" story
prominently. The answer to "why not copy-paste like shadcn?" is:

```
utils.Class()      → Tailwind class overrides (tailwind-merge resolves conflicts)
@theme CSS vars    → Color/font/spacing theming without Go changes
Attrs map          → Arbitrary HTML attributes on any component
templ.Component    → Slot-based composition for flexible content injection
```

These four mechanisms give ~95% of shadcn's customization power without owning
maintenance. The remaining 5% (structural DOM changes) is available via Go's
template composition — wrap the component and override what you need.

---

## 3. Architecture: The Headless-Styled Divide

The most important architectural pattern in modern component libraries is the **separation
of behavior from rendering**. Three layers, studied from the best implementations:

### Layer 1: Behavior-Only Primitives (Radix UI)

Radix ships **zero styles** — only behavior, accessibility, and structure:

```jsx
// Radix Dialog — you provide ALL styling
import { Dialog } from "radix-ui";

<Dialog.Root>
  <Dialog.Trigger>Open</Dialog.Trigger>
  <Dialog.Portal>
    <Dialog.Overlay className="fixed inset-0 bg-black/50" />
    <Dialog.Content className="fixed centered bg-white p-6 rounded-lg">
      <Dialog.Title>Are you sure?</Dialog.Title>
      <Dialog.Description>This action cannot be undone.</Dialog.Description>
      <Dialog.Close>Cancel</Dialog.Close>
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
```

**What Radix provides (you don't write this):**
- Focus trapping inside `Content` (Tab cycles, never escapes)
- Focus restoration to `Trigger` on close
- `aria-labelledby` / `aria-describedby` wired from `Title` / `Description`
- Escape key dismissal
- `data-state="open"` / `"closed"` attributes for CSS styling
- Portal rendering (escapes z-index and overflow contexts)
- `asChild` prop for element replacement (use `<a>` instead of `<button>` as trigger)

**What you provide:**
- ALL CSS (even functional styles like overlay coverage)
- Layout and spacing
- Animation (via `data-state` CSS selectors)

**The `data-*` styling bridge:** Radix exposes internal state through data attributes:

```css
.AccordionItem[data-state="open"] { border-bottom-width: 2px; }
.PopoverContent[data-side="top"]  { animation-name: slideUp; }
.PopoverContent[data-side="bottom"] { animation-name: slideDown; }
```

**Controlled/uncontrolled pattern:** Every stateful component supports both:

```jsx
// Uncontrolled (simple — no state management)
<Dialog.Root defaultOpen={false}>...</Dialog.Root>

// Controlled (full control — for async, validation, orchestration)
const [open, setOpen] = useState(false);
<Dialog.Root open={open} onOpenChange={setOpen}>...</Dialog.Root>
```

### Layer 2: Behavior Hooks (React Aria)

React Aria goes further — it separates behavior into **hooks** that return props:

```jsx
function ComboBox(props) {
  let state = useComboBoxState(props);
  let {
    buttonProps, inputProps, listBoxProps, labelProps
  } = useComboBox({...props}, state);

  return (
    <>
      <label {...labelProps}>{props.label}</label>
      <div className="combo-wrapper">
        <input {...inputProps} />
        <button {...buttonProps}>▼</button>
        <ul {...listBoxProps}>
          {state.collection.getItems().map(item => (
            <Option key={item.key} item={item} state={state} />
          ))}
        </ul>
      </div>
    </>
  );
}
```

The hooks handle: keyboard navigation, focus management, ARIA attributes, i18n
(locale-aware), filtered item counting announcements, touch screen reader support —
all invisible to the consumer.

### Layer 3: Styled Components (shadcn/ui, templ-components, templui)

This is where styling meets behavior. Two sub-approaches:

**Monolithic (templ-components):** One component, one props struct, library owns the DOM:
```go
@display.Modal(display.ModalProps{Title: "Delete?", Open: true, Size: display.ModalSizeMD})
```

**Compound (shadcn/ui, templui):** Multiple sub-components, consumer owns the DOM:
```go
// templui compound pattern
@dialog.Dialog() {
    @dialog.Trigger() { @button.Button("Open") }
    @dialog.Content() {
        @dialog.Header() {
            @dialog.Title() { "Are you sure?" }
            @dialog.Description() { "This cannot be undone." }
        }
        @dialog.Footer() { @dialog.Close() { @button.Button("Cancel") } }
    }
}
```

### When to use which

| Component type | Recommended pattern | Why |
|---------------|-------------------|-----|
| Simple primitives (Badge, Spinner, Skeleton) | Monolithic | Structure is fixed; flexibility isn't needed |
| Form inputs (Input, Select, Textarea) | Monolithic | HTML structure is standardized; wrapper adds value |
| Complex overlays (Dialog, Popover, Drawer) | **Compound** | Consumer needs to compose trigger/content/header/footer |
| Data display (Table, Tree, DataTable) | **Compound** | Structure varies dramatically per use case |
| Navigation (Tabs, Accordion, Breadcrumbs) | Either | Depends on structural flexibility needs |

**For templ-components:** The current monolithic pattern is correct for simple components
but limits complex overlays. For future components (Popover, ContextMenu, HoverCard,
Command Palette), consider the compound pattern where `Trigger`, `Content`, `Close` are
separate sub-components wired by `data-tc-*` attributes.

---

## 4. The Component API Contract

Every component in a superb library follows a consistent API contract. This is what makes
the library predictable for consumers and maintainable for authors.

### The templ-components Contract (exemplified by Button)

From the actual source (`display/button.templ` + `display/button_go.go`):

```go
// 1. Props struct embeds utils.BaseProps — auto-satisfies ComponentProps interface
type ButtonProps struct {
    utils.BaseProps          // ID, Class, Attrs, AriaLabel, Nonce
    Text     string
    Type     ButtonHTMLType  // typed enum, not string
    Href     string
    Variant  ButtonType      // typed enum
    Size     ButtonSize      // typed enum
    Disabled bool
    Icon     templ.Component // slot for icon
    External bool
}

// 2. Root element propagates ALL BaseProps
templ Button(props ButtonProps) {
    // Class is merged LAST so consumer overrides win (tailwind-merge)
    {{ classes := utils.Class(
        buttonBaseClass,
        buttonVariantClasses(props.Variant),
        buttonSizeClasses(props.Size),
        props.Class,
    ) }}
    <button
        if props.ID != "" { id={ props.ID } }                          // ID
        class={ classes }                                               // Class (merged)
        if props.AriaLabel != "" { aria-label={ props.AriaLabel } }   // AriaLabel
        { props.Attrs... }                                              // Attrs (spread)
    >
        if props.Icon != nil { @props.Icon }
        { props.Text }
    </button>
}
```

**The four propagation rules (no exceptions):**
1. **ID** → `if props.ID != "" { id={ props.ID } }`
2. **Class** → `utils.Class(libraryClasses..., props.Class)` — consumer always wins
3. **Attrs** → `{ props.Attrs... }` — spread for arbitrary HTML attributes
4. **AriaLabel** → `if props.AriaLabel != "" { aria-label={ props.AriaLabel } }`

### What Radix teaches us about this contract

Radix adds two capabilities that templ-components doesn't have yet:

**`asChild` (element polymorphism):** Any part can become any element:
```jsx
<Dialog.Trigger asChild>
  <a href="/settings">Settings</a>  {/* anchor instead of button */}
</Dialog.Trigger>
```
This is powerful for HATEOAS — a dialog trigger can be a hyperlink, preserving
navigation semantics. In Go templ, this pattern is harder (templates emit fixed tags),
but the spirit can be achieved via conditional rendering (as Button already does with
`Href` → `<a>` vs `<button>`).

**`data-*` state exposure:** Components expose internal state for CSS styling:
```html
<div data-state="open" data-side="top">...</div>
```
templ-components already does this (e.g., `data-tc-dialog-open="true"` in templui
equivalents) but should standardize the naming convention across all interactive
components.

### The ComponentProps Interface

```go
type ComponentProps interface {
    GetBaseProps() *BaseProps
    SetBaseProps(*BaseProps)
}
```

Every props struct that embeds `utils.BaseProps` auto-satisfies this via promoted methods
(pointer receivers, required by `recvcheck`). This enables generic wrappers — a consumer
can write utilities that work with ANY component:

```go
func WithNonce[T utils.ComponentProps](props T, nonce string) T {
    props.GetBaseProps().Nonce = nonce
    return props
}
```

**Critical:** Register every new props type in `internal/contract/component_props_test.go`
so the `TestAllComponentPropsSatisfyInterface` test enforces the contract at CI time.

---

## 5. Variant Systems: Type-Safe Styling

### The Three Approaches Compared

#### A. Map + Lookup (templ-components — recommended)

```go
// Typed enum
type BadgeType string
const (
    BadgeDefault BadgeType = "default"
    BadgePrimary BadgeType = "primary"
    BadgeSuccess BadgeType = "success"
)

// Map lookup with explicit fallback
var badgeStyleMap = map[BadgeType]string{
    BadgePrimary: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
    BadgeSuccess: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
}

func badgeClass(t BadgeType) string {
    return utils.Lookup(badgeStyleMap, t, badgeStyleMap[BadgeDefault])
}
```

**Strengths:** Data not code, graceful fallback, no panic possible, easy to audit.
**When to use:** Pure class/style data (colors, sizes, padding, borders).

#### B. Switch Statement (templui — not recommended)

```go
func (b Props) variantClasses() string {
    switch b.Variant {
    case VariantDestructive: return "bg-destructive text-white ..."
    case VariantOutline:     return "border bg-background ..."
    default:                 return "bg-primary text-primary-foreground ..."
    }
}
```

**Weakness:** Code not data, can't be easily extended, missing case = silent default.
**When to use:** Only when the switch also emits different DOM structure (but then use
an if-branch in the template instead).

#### C. CVA / class-variance-authority (shadcn/ui — TypeScript only)

```typescript
const buttonVariants = cva("inline-flex items-center justify-center ...", {
  variants: {
    variant: {
      default:     "bg-primary text-primary-foreground hover:bg-primary/90",
      destructive: "bg-destructive text-white hover:bg-destructive/90",
    },
    size: {
      default: "h-9 px-4 py-2",
      sm:      "h-8 rounded-md gap-1.5 px-3",
      icon:    "size-9",
    },
  },
  compoundVariants: [
    { variant: "primary", size: "sm", class: "shadow-sm" },
  ],
  defaultVariants: { variant: "default", size: "default" },
})
```

**Strengths:** Type-safe (invalid props = compile error), compound variants for
multi-condition styles, pre-computed at module load.
**Go equivalent:** Not directly translatable, but the map+lookup approach provides the
same data-driven separation. Compound variants are the one feature Go can't easily
replicate — they'd need a struct slice checked at lookup time.

### The Decision Tree

```
Is it pure class/style data (colors, sizes, padding)?
├── YES → Map + utils.Lookup(map, key, fallback)
│         (badgeStyleMap, buttonVariantLookup, cardPaddingLookup, etc.)
└── NO → Does it change the DOM structure?
    ├── YES → if-branch in the template
    │         (TabsVariant, DropdownPosition, TrendDirection, StatCardHref)
    └── NO → Enum validation
              (InputType, ButtonHTMLType, FormMethod → map+fallback)
```

### Compound variants (the CVA feature Go can adopt)

shadcn's `compoundVariants` apply styles when multiple conditions match:

```typescript
{ variant: "primary", disabled: false, class: "hover:bg-primary/90" }
```

In Go, this can be achieved by composing the lookup:

```go
func buttonClasses(p ButtonProps) string {
    base := buttonVariantLookup[p.Variant]
    if !p.Disabled && p.Variant == ButtonPrimary {
        base += " hover:bg-blue-500"
    }
    return utils.Class(base, p.Class)
}
```

This is more verbose than CVA but achieves the same result. For complex compound logic,
extract a dedicated function (`buttonModifierClasses`).

---

## 6. Theming: The Three-Layer Token Stack

Every mature design system converges on the same architecture:

```
┌─────────────────────────────────────────────────────┐
│  Layer 3: Component Tokens                           │
│  --button-bg → var(--color-accent)                   │
│  --input-border → var(--color-border)                │
├─────────────────────────────────────────────────────┤
│  Layer 2: Semantic Tokens (the theming boundary)     │
│  --color-accent → var(--blue-500)                    │
│  --color-bg → var(--white)                           │
│  --color-text → var(--gray-900)                      │
│  --color-border → var(--gray-200)                    │
├─────────────────────────────────────────────────────┤
│  Layer 1: Primitive Tokens (raw values)              │
│  --blue-500: #3b82f6                                  │
│  --gray-900: #111827                                  │
│  --radius-lg: 0.5rem                                  │
└─────────────────────────────────────────────────────┘
```

### The Critical Rule

Each layer only references the layer directly below:

| Change type | Which layer to update | Example |
|-------------|----------------------|---------|
| Brand color change | Primitive | `--blue-500: #6366f1` |
| Theme switch (light/dark) | Semantic | `--color-bg: var(--gray-950)` |
| Component restyle | Component | `--button-bg: var(--color-success)` |

### shadcn/ui's Implementation (OKLCH + semantic pairs)

shadcn uses **OKLCH** color values (perceptually uniform) in semantic token pairs:

```css
:root {
  --primary: oklch(0.205 0 0);
  --primary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.97 0 0);
  --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --border: oklch(0.922 0 0);
  --ring: oklch(0.708 0 0);
}
.dark {
  --primary: oklch(0.922 0 0);
  --primary-foreground: oklch(0.205 0 0);
  /* ...all tokens redefined */
}
```

Mapped to Tailwind via `@theme inline`:
```css
@theme inline {
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
}
```

So `bg-primary` in component code resolves to `var(--primary)` → `oklch(...)`.

### templ-components' Current Approach (Gap Identified)

templ-components currently uses **primitive Tailwind names** directly in components:

```go
// Current: components reference raw color scales
var buttonVariantLookup = map[ButtonType]string{
    ButtonPrimary: "bg-blue-600 text-white hover:bg-blue-500 ...",
    ButtonDanger:  "bg-red-600 text-white hover:bg-red-500 ...",
}
```

**Problem:** A consumer who wants to change the "primary" color must remap
`--color-blue-600` globally, which affects ALL uses of blue-600 — not just "primary"
uses. There's no semantic boundary between "primary" and "this shade of blue."

### Recommended Fix: Semantic Token Layer

Introduce an optional semantic token layer in `templ-components-theme.css`:

```css
/* templ-components-theme.css — semantic tokens */
@theme {
  --color-tc-primary: var(--color-blue-600);
  --color-tc-primary-hover: var(--color-blue-700);
  --color-tc-primary-foreground: var(--color-white);

  --color-tc-secondary: var(--color-gray-100);
  --color-tc-secondary-foreground: var(--color-gray-900);

  --color-tc-destructive: var(--color-red-600);
  --color-tc-success: var(--color-green-600);
  --color-tc-warning: var(--color-yellow-600);

  --color-tc-background: var(--color-white);
  --color-tc-foreground: var(--color-gray-900);
  --color-tc-muted: var(--color-gray-100);
  --color-tc-muted-foreground: var(--color-gray-500);
  --color-tc-border: var(--color-gray-200);

  --color-tc-ring: var(--color-blue-600);
}
```

Then components can use `bg-tc-primary` instead of `bg-blue-600`. Consumers retheme by
reassigning `--color-tc-primary` to any color — no collision with raw Tailwind utilities.

**Migration strategy:** Introduce semantic tokens as an opt-in layer first. Components
continue to use primitive names by default. Consumers who want semantic theming import
`templ-components-theme.css` and use the `tc-*` classes. In a future major version,
flip the default to semantic tokens.

---

## 7. HATEOAS: HTML as the Source of Truth

> *"A natural hypermedia such as HTML is a practical necessity for building RESTful
> systems."* — [htmx.org/essays/hateoas](https://htmx.org/essays/hateoas/)

### The Core Insight

HATEOAS means: **the server encodes all available actions directly in the HTML response.**
The browser needs no prior knowledge of the application — it just renders hypermedia.

| Aspect | Hypermedia (HTML) | JSON API + SPA |
|--------|-------------------|----------------|
| **Available actions** | Encoded in response (`<a>`, `<form>`, `hx-*`) | Client hardcodes URLs/methods |
| **State transitions** | Server controls via hypermedia changes | Client maintains its own state machine |
| **Decoupling** | Server evolves independently | Coordinated server+client updates required |
| **Self-describing** | `<form method="post" action="...">` says everything | `{"status": "overdrawn"}` says nothing |

### Practical Example: Delete Action

**Hypermedia approach (HTMX):**
```
User clicks "Delete"
  → hx-delete="/items/42" fires (hypermedia control in HTML)
  → Server processes deletion
  → Server returns updated item list HTML (new hypermedia)
  → HTMX swaps old list for new
  → New HTML naturally omits "Delete" for the removed item
  → NO client-side state synchronization needed
```

**JSON API approach (SPA):**
```
User clicks "Delete"
  → Client sends DELETE /api/items/42
  → Server returns {"success": true}
  → Client must KNOW to remove item 42 from local state
  → Client must KNOW to re-render the list
  → Client must KNOW the item no longer has actions
  → ALL of this is out-of-band coupling
```

### Implications for Component Libraries

1. **Server-rendered HTML is a complete UI strategy.** A component library rendering HTML
   on the server gets HATEOAS benefits for free.

2. **Prefer native HTML primitives.** `<details>`/`<summary>` for disclosure, `<form>`
   for submit, `<a>` for navigation, `<dialog>` for modals. These ARE hypermedia controls.

3. **JavaScript enhances, never replaces.** When JS is needed, it reads state FROM HTML
   attributes (`data-tc-*`, `datetime`, `aria-*`) — progressive enhancement, not SPA
   replacement.

4. **HTMX extends HTML as hypermedia.** `hx-get`, `hx-post`, `hx-swap` turn any element
   into a hypermedia control. This IS HATEOAS in practice.

5. **templui's Accordion is HATEOAS-native:** It uses `<details>`/`<summary>` — a native
   HTML disclosure element that works without ANY JavaScript. The server controls open/closed
   state via the `open` attribute. This is textbook HATEOAS.

---

## 8. Form Handling for Server-Rendered Apps

Forms are the #1 use case for web applications. A superb server-rendered component library
must handle forms better than any SPA framework.

### 8.1 Progressive Enhancement (Works Without JS)

A form must work via standard HTML POST without any JavaScript, then HTMX layers on top:

```html
<!-- This form works without JS. HTMX enhances it when available. -->
<form action="/register" method="POST"
      hx-post="/register"
      hx-target="#form-container"
      hx-swap="outerHTML">
    @forms.Input(forms.InputProps{
        Name: "email",
        Type: forms.InputTypeEmail,
        Label: "Email",
        Required: true,
    })
    @forms.Input(forms.InputProps{
        Name: "password",
        Type: forms.InputTypePassword,
        Label: "Password",
        Required: true,
    })
    @forms.CSRFToken(csrfToken)
    @display.Button(display.ButtonProps{Text: "Register", Type: display.ButtonSubmit})
</form>
```

With `hx-boost`, forms degrade gracefully: they continue to work, they just don't use AJAX.

### 8.2 Accessible Error Display

The WCAG-compliant pattern for linking errors to fields:

```html
<label for="email">Email *</label>
<input type="email" id="email" name="email" required
       aria-invalid="true"
       aria-describedby="email-error" />
<div id="email-error" class="text-sm text-red-600">
    <span aria-hidden="true">⚠</span> Please enter a valid email address
</div>
```

**Key rules:**
1. `aria-invalid="true"` — set ONLY on fields that fail validation
2. `aria-describedby` — links input to its error message by ID; screen reader announces both
3. Reset ARIA attributes on each submission — stale errors confuse screen readers
4. Never rely on color alone — combine border color with icon + text

**templ-components already does this** via `forms.ErrorAttrs(id, errMsg, helpTextID)`:
```go
// Returns templ.Attributes for aria-invalid/aria-describedby
attrs := forms.ErrorAttrs("email", "Invalid email", "email-help")
```

### 8.3 Validation Summary (Accessibility-Critical)

A validation summary lists ALL errors above the form, with each error linking to its field:

```html
<div id="error-summary" tabindex="-1" role="alert" class="rounded-md bg-red-50 p-4">
    <h3>2 errors found</h3>
    <ul>
        <li><a href="#email">Please enter a valid email address</a></li>
        <li><a href="#password">Password must be at least 8 characters</a></li>
    </ul>
</div>
```

**templ-components' implementation** (`forms/validation.templ`):
```go
type ValidationError struct {
    Field   string
    Message string
}
type ValidationSummaryProps struct {
    utils.BaseProps
    Errors []ValidationError
}
```

The error count uses singular/plural:
```go
fmt.Sprintf("%d error%s found", len(props.Errors),
    utils.Ternary(len(props.Errors) == 1, "", "s"))
```

Each error with a `Field` becomes a link to `"#" + SanitizeID(err.Field)` — clicking jumps
focus directly to the problem field.

### 8.4 CSRF Protection

**Hidden input (preferred — survives page reloads, works without JS):**
```html
<form hx-post="/api/settings">
    <input type="hidden" name="csrf_token" value="{ token }" />
</form>
```

**`hx-headers` (global — but stale after `hx-boost` swaps `<html>`/`<body>`):**
```html
<body hx-headers='{"X-CSRF-TOKEN": "TOKEN_HERE"}'>
```

**templ-components provides** `htmx.CSRFToken(token string)` which renders the hidden input.

### 8.5 HTMX-Specific Form Patterns

**Per-field real-time validation:**
```html
<input type="email" name="email"
    hx-post="/validate?field=email"
    hx-trigger="keyup changed delay:500ms"
    hx-target="#email-error" />
```

**Race condition prevention with `hx-sync`:**
```html
<form hx-post="/register">
    <input name="email"
        hx-post="/validate"
        hx-trigger="change"
        hx-sync="closest form:abort" />  <!-- aborts validation if form submits -->
    <button type="submit">Register</button>
</form>
```

**HTTP 422 for validation errors** (HTMX ignores 4xx by default — configure it):
```html
<meta name="htmx-config" content='{"responseHandling":[{"code":"422","swap":true}]}' />
```

**Eliminating Post/Redirect/Get:** HTMX returns the HTML fragment directly — no 302
redirect needed after successful POST.

---

## 9. Accessibility: Not a Feature, the Foundation

### 9.1 The Five Rules of ARIA

1. **Prefer native HTML.** A `<button>` is universally accessible. A `<div role="button">`
   requires manual keyboard handling, focus management, and activation — and you'll never
   perfectly replicate the browser default.

2. **Don't change native semantics.** Nest semantic elements inside role containers.

3. **All interactive ARIA controls must be keyboard-accessible.** `role="button"` must
   respond to both Enter AND Space.

4. **Don't hide focusable elements.** Never `aria-hidden="true"` on a focusable element
   without `tabindex="-1"`.

5. **"No ARIA is better than bad ARIA."** Pages with ARIA average **41% MORE errors**
   than those without (WebAIM million-page survey). Bad ARIA actively harms users.

### 9.2 Mandatory Keyboard Patterns (WAI-ARIA APG)

| Component | Keyboard Pattern |
|-----------|-----------------|
| **Dialog/Modal** | Escape closes; Tab trapped inside; focus restored to trigger |
| **Tabs** | Arrow Left/Right between tabs; Tab enters panel; Home/End for first/last |
| **Menu/Dropdown** | Arrows navigate; Enter/Space activates; Escape closes |
| **Accordion** | Enter/Space toggles; arrows move between headers (if not using `<details>`) |
| **Combobox** | Arrow Up/Down moves selection; Home/End; Escape clears |
| **Tooltip** | Opens on focus; Escape/hover-out dismisses |

### 9.3 What "Accessible by Default" Means

An accessible-by-default library **encapsulates accessibility complexity**. A well-built
Modal automatically:
- Sets `role="dialog"` and `aria-modal="true"`
- Moves focus into the dialog on open
- Traps focus while open (Tab cycles inside, never escapes)
- Restores focus to trigger on close
- Closes on Escape

The consumer never thinks about any of this. **This is the core value proposition.**

### 9.4 The Radix/React Aria Standard

Radix UI and React Aria are the gold standards because:

- **Every component maps to a specific WAI-ARIA APG pattern** — the implementation IS
  spec compliance.
- **Tested against real assistive technology** — VoiceOver (macOS/iOS), JAWS (Windows),
  NVDA (Windows), TalkBack (Android). Not just automated checkers.
- **Focus management is programmatically correct** — not just `tabindex` hacks but
  context-aware focus movement (e.g., AlertDialog focuses the Cancel button on open to
  anticipate the safe response).
- **Cross-device normalization** — `data-hovered` differs from `:hover` (no sticky touch
  states); `data-pressed` differs from `:active` (cancellable by dragging away);
  `data-focus-visible` differs from `:focus` (keyboard only, not mouse click).

### 9.5 Motion Safety

Every transition must carry:
```css
motion-reduce:transition-none motion-reduce:duration-0
```
Every animation must carry:
```css
motion-reduce:animate-none
```

**templ-components enforces this** via a11y tests that check for `motion-reduce:*` on
every transition and animation.

---

## 10. Motion Design: Professional Animation

Animation separates a polished library from an amateur one. But bad animation is worse
than no animation.

### 10.1 Duration Guidelines

| Interaction | Duration | Rationale |
|-------------|----------|-----------|
| Button press | 80–100ms | Must feel instant |
| Toggle/checkbox | 80–150ms | Quick, clear feedback |
| Button hover | ~200ms | Responsive but not distracting |
| Tooltip | ~100ms | Informational, shouldn't distract |
| Dropdown menu | ~200ms | Predictable, avoid bounce |
| Modal entrance | 250–300ms | Large movement needs gentle timing |
| Drawer/panel | 250–350ms | Spatial reorientation |
| Page transition | 300–500ms | Large context change |
| Success feedback | 400–700ms | Emphasis, deserves time |

**Rule:** Keep animations under 300ms for perceived performance. Exits should be 60–70%
of entrance duration — an exit that lingers blocks the user.

### 10.2 Easing Curves

**The core principle:** Linear motion feels robotic. Real objects accelerate and decelerate.

| Easing | Use For | Avoid |
|--------|---------|-------|
| **`ease-out`** (fast start, slow finish) | Entrances, user-initiated interactions — the workhorse | — |
| **`ease-in`** (slow start, fast finish) | Exits and dismissals only | Never for entrances — feels sluggish |
| **`ease-in-out`** | Elements already moving on screen | Exits (content lingers) |
| **`linear`** | Spinners, progress bars — anything representing passage of time | Almost everything else |

**The asymmetry rule:** Entrances use ease-out (fast-in, slow-out). Exits use ease-in
(slow-in, fast-out). This mirrors physical object behavior.

**Professional custom curves:**
```css
/* Smooth ease-out — professional default for large movements */
cubic-bezier(0.16, 1, 0.3, 1)

/* Spring-like — for playful micro-interactions (buttons, cards) */
cubic-bezier(0.34, 1.56, 0.64, 1)

/* Fast response — for toggles, tooltips, icon transitions */
cubic-bezier(0.4, 0, 0.2, 1)
```

### 10.3 Performant Properties

**Only animate `transform` and `opacity`.** These are GPU-accelerated (composite-only).

```css
/* BAD: triggers layout recalculation every frame */
.dropdown { transition: height 300ms; height: 0; }
.dropdown.open { height: 200px; }

/* GOOD: compositor-only, GPU-accelerated */
.dropdown { transition: transform 300ms; transform: scaleY(0); transform-origin: top; }
.dropdown.open { transform: scaleY(1); }
```

| ❌ Avoid (layout + paint) | ✅ Use (composite only) |
|---|---|
| `width`, `height` | `transform: scale()` |
| `top`, `left`, `right`, `bottom` | `transform: translate()` |
| `margin`, `padding` | `transform` |
| `border-width` | `transform` |

### 10.4 `prefers-reduced-motion`

Not "no motion" — **no motion that triggers vestibular disorders.** Fades and
cross-dissolves are safe replacements.

```css
/* Replace vestibular triggers (scaling) with safe motion (fade) */
.card { transition: transform 300ms var(--ease-enter); }

@media (prefers-reduced-motion: reduce) {
    .card { transition: opacity 150ms linear; }
}
```

In Tailwind, this is expressed as utility classes on every transition/animation:
```
motion-reduce:transition-none motion-reduce:duration-0   (transitions)
motion-reduce:animate-none                                (animations)
```

### 10.5 CSS `@starting-style` (Zero-JS Enter Animations)

Chrome 117+, Firefox, Safari (August 2024 Baseline). Enables enter animations without
JavaScript — critical for server-rendered components:

```css
dialog[open] {
    opacity: 1;
    transform: translateY(0);
    transition: opacity 0.3s, transform 0.3s, overlay 0.3s allow-discrete, display 0.3s allow-discrete;
}

@starting-style {
    dialog[open] {
        opacity: 0;
        transform: translateY(20px);
    }
}

dialog {
    opacity: 0;
    transform: translateY(20px);
}
```

This enables enter/exit animations on `<dialog>`, `popover`, and any `display: none`
element — pure CSS, no JavaScript required. A significant capability for HATEOAS-first
component libraries.

### 10.6 Motion Tokens

Professional libraries define motion as design tokens, not scattered literals:

```css
:root {
    --ease-enter: cubic-bezier(0, 0, 0.3, 1);
    --ease-exit: cubic-bezier(0.4, 0, 1, 1);
    --ease-standard: cubic-bezier(0.4, 0, 0.2, 1);
    --duration-instant: 100ms;
    --duration-fast: 200ms;
    --duration-standard: 300ms;
    --duration-slow: 500ms;
}
```

---

## 11. Progressive Enhancement & JavaScript

### The Decision Ladder

Before writing ANY JavaScript for a component:

```
1. Can native HTML do it?
   ├── <details>/<summary> → disclosure/accordion (NO JS NEEDED)
   ├── <form>              → submit (NO JS NEEDED)
   ├── <dialog>            → modal (native in modern browsers)
   ├── :target / :checked  → CSS-only toggles
   └── If YES → STOP. No JS needed.

2. If JS is unavoidable:
   ├── Guard with global singleton (idempotent across HTMX re-renders)
   ├── Use event delegation on document
   ├── Handle Escape, focus save/restore, click-outside
   ├── Read ALL state from data-* HTML attributes
   └── Escape IDs in JS with strconv.Quote() to prevent XSS

3. Never:
   ├── Maintain a parallel client-side state store
   ├── Replace HTML with JS-rendered content
   ├── Break without JavaScript (progressive enhancement)
   └── Use inline event handlers (onclick — violates CSP)
```

### The Singleton Guard Pattern (templ-components)

Every JS attachment in templ-components uses a global singleton flag:

```javascript
(function() {
    if (window.tcComboboxAttached) return;  // already bound — idempotent
    window.tcComboboxAttached = true;

    document.addEventListener('click', function(e) {
        // event delegation — handles dynamically added elements
        var trigger = e.target.closest('[data-tc-combobox]');
        if (!trigger) return;
        // ... handle combobox interaction
    });
})();
```

This is **critical for HTMX**: when HTMX re-renders a region, the inline script runs
again. Without the guard, handlers would double-bind. The singleton ensures the script
is idempotent — safe to re-run any number of times.

### templui's `templ.NewOnceHandle()` Alternative

templui uses templ's built-in `OnceHandle` for script deduplication:

```go
var scriptOnce = templ.NewOnceHandle()

templ Script() {
    @scriptOnce.Once() {
        @utils.ComponentScript("dialog")
    }
}
```

This prevents the same script from being rendered twice in a single page — useful when
a component appears multiple times. But it doesn't handle HTMX re-renders (the handle
resets on full page swap). The singleton-flag approach is more robust for HTMX.

### JavaScript as Progressive Enhancement

The key principle: JavaScript reads state from HTML and enhances behavior. It never
creates state independently.

```javascript
// GOOD: reads state from HTML attributes
var modal = document.querySelector('[data-tc-modal]');
var open = modal.getAttribute('data-tc-modal-open') === 'true';

// BAD: maintains parallel state in JS
var modalState = { open: false };  // where does this come from? the server!
```

When HTMX swaps new HTML, the JS reads the new `data-*` attributes — it doesn't need to
be notified separately. The HTML IS the state.

---

## 12. Security: CSP-Safe by Construction

Content Security Policy (CSP) is defense against XSS. A strict CSP forbids: inline scripts
without nonces, `eval()`, inline event handlers, and `javascript:` URLs.

### The templ-components Approach (Gold Standard)

- **Every inline `<script>` carries `nonce={ props.Nonce }`.** No exceptions.
- **`layout.Script(nonce, src, attrs)`** for external scripts — auto-injects nonce.
- **No `eval()`, no inline handlers, no `javascript:` URLs.**
- **Nonce flow:** `BaseProps.Nonce` → `overlayShellProps.nonce` → `scriptComponent` →
  `<script nonce="...">` with `html.EscapeString(nonce)`.

The implementation bypasses templ's script-context sanitization deliberately:

```go
func scriptComponent(nonce, js, errLabel string) templ.Component {
    escapedNonce := html.EscapeString(nonce)
    return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
        _, err := fmt.Fprintf(w, "<script nonce=\"%s\">\n%s</script>\n", escapedNonce, js)
        return err
    })
}
```

### Why This Matters

A consumer with `Content-Security-Policy: script-src 'self' 'nonce-{random}'` can use
every templ-components component without forking. The nonce flows automatically through
`BaseProps`. The consumer never has to remember to add it — it's structurally impossible
to emit an unnonce'd script.

### Comparison with templui

templui ships **separate `.js` files** (`dialog.js`, `dialog.min.js`) loaded as external
scripts. This works with CSP (external scripts are allowed without nonces) but:
1. Consumer must serve the JS files from their own static asset path
2. More moving parts (separate files to manage, minify, version, serve)
3. Breaks if the consumer's static file path differs from the library's expectation

**templ-components' inline-script-with-nonce is simpler for consumers** — the JS travels
with the component, no external assets to manage.

---

## 13. Internationalization & RTL

A truly superb component library is i18n-ready. This is where React Aria excels and most
libraries fall short.

### CSS Logical Properties

Replace physical properties with logical ones that automatically flip in RTL:

| Physical (❌) | Logical (✅) | What it does |
|-------------|------------|-------------|
| `margin-left` | `margin-inline-start` | Start margin (left in LTR, right in RTL) |
| `margin-right` | `margin-inline-end` | End margin |
| `padding-left` | `padding-inline-start` | Start padding |
| `text-align: left` | `text-align: start` | Start-aligned text |
| `left: 0` | `inset-inline-start: 0` | Start position |
| `border-left` | `border-inline-start` | Start border |

### Tailwind Logical Utilities (v4)

Tailwind provides logical utility aliases:

| Physical | Logical | Notes |
|----------|---------|-------|
| `ml-4` | `ms-4` | margin-inline-start |
| `mr-4` | `me-4` | margin-inline-end |
| `pl-4` | `ps-4` | padding-inline-start |
| `pr-4` | `pe-4` | padding-inline-end |
| `left-0` | `start-0` | inset-inline-start |
| `right-0` | `end-0` | inset-inline-end |
| `text-left` | `text-start` | text-align: start |
| `text-right` | `text-end` | text-align: end |

### Icon Mirroring

In RTL, **directional icons must flip** — but symmetric icons must not:

| Mirror in RTL | Don't mirror |
|--------------|-------------|
| Arrows (←→↑↓) | Checkmarks |
| Chevrons (‹›) | X / close |
| Play/forward/rewind | Plus, minus |
| "Back" / "forward" icons | Dots, circles |
| Bread crumb separators | Avatars |

CSS for icon flipping:
```css
[dir="rtl"] .icon-directional {
    transform: scaleX(-1);
}
```

### What Component Libraries Must Do

1. **Use logical Tailwind utilities** (`ms-`, `me-`, `ps-`, `pe-`, `start-`, `end-`)
   instead of physical ones (`ml-`, `mr-`, `pl-`, `pr-`, `left-`, `right-`).
2. **Test with `dir="rtl"`** — add RTL test cases to the test suite.
3. **Mirror directional icons** via `[dir="rtl"]` CSS selectors.
4. **Date/number formatting** — use Go's `time.Format` with locale, or ICU formatting.

**Gap in templ-components:** The library uses physical properties (`ml-`, `mr-`) in some
places. Migrating to logical utilities (`ms-`, `me-`) would make the library RTL-ready
with zero behavioral change in LTR contexts.

---

## 14. Tailwind CSS v4: The CSS-First Revolution

v4 is a ground-up rewrite that fundamentally changes component library authoring.

### What Changed

| Feature | v3 | v4 |
|---------|----|----|
| Configuration | `tailwind.config.js` (JS) | `@theme` in CSS (CSS-first) |
| Content detection | `content: [...]` array | Automatic (respects `.gitignore`) |
| Color space | RGB | OKLCH (perceptually uniform, wider gamut) |
| Performance | Baseline | Full builds 5× faster, incremental 182× faster |
| Dark mode | `dark:` with config | `@custom-variant dark (...)` in CSS |
| Custom utilities | Plugin API (`addUtilities`) | `@utility` directive in CSS |
| Custom variants | Plugin API (`addVariant`) | `@custom-variant` directive in CSS |
| External scanning | Manual config | `@source "../node_modules/my-lib"` |

### What This Means for templ-components

1. **CSS variables are the universal contract.** Design tokens (`--color-*`, `--font-*`)
   are real runtime CSS variables. Consumers read them in CSS, inline styles, or JS.

2. **`@theme` as shareable design system.** Publish a design-token package as pure CSS;
   consumers `@import` it. No build step, no JS bundling.

3. **`@source` for library scanning.** Consumers add one line:
   ```css
   @source "../node_modules/templ-components";
   ```
   Wait — templ-components is a Go library, not npm. But the Tailwind scanning still
   needs to find the class strings in the generated `*_templ.go` files. The
   `docs/tailwind-v4-adoption-guide.md` documents this setup.

4. **`@custom-variant` for domain-specific variants.** The library ships:
   ```css
   @custom-variant dark (&:where(.dark, .dark *));
   ```
   Consumers could create custom variants:
   ```css
   @custom-variant tc-error (&:where([data-tc-state="error"]));
   ```

5. **`@utility` for custom utilities.** Libraries can ship first-class utilities:
   ```css
   @utility tc-skeleton {
       background: linear-gradient(90deg, rgba(0,0,0,0.06) 25%, rgba(0,0,0,0.12) 37%, ...);
       background-size: 400% 100%;
       animation: shimmer 1.4s ease infinite;
   }
   ```

6. **OKLCH color space.** The entire default palette switched from RGB to OKLCH —
   perceptually uniform, wider gamut, better dark-mode color selection.

### Breaking Changes to Watch

- Default border color: `gray-200` → `currentColor`
- Ring default width: 3px → 1px (use `ring-3`)
- Shadow scale renamed: `shadow-sm` → `shadow-xs`, `shadow` → `shadow-sm`
- Important modifier: `!flex` → `flex!`
- Variable shorthand: `bg-[--brand]` → `bg-(--brand)`

---

## 15. Testing Strategy

A superb library tests at multiple levels. "Assertion tests only" is not acceptable.

### The templ-components Testing Matrix

Every component MUST have all of these:

```
[✓] golden_test.go      — exact rendered HTML matches .golden snapshot
[✓] a11y_test.go        — ARIA, roles, keyboard, motion-reduce, screen-reader text
[✓] bdd_test.go         — behaviour spec (user-visible behaviour, not markup)
[✓] edge_cases_test.go  — empty inputs, unknown enum values, ID collisions
[✓] example_test.go     — godoc ExampleXxx() compiles and renders
[✓] snapshot_test.go    — broader composition snapshot
[✓] coverage_*_test.go  — targeted coverage of private helpers and branches
```

### Golden Testing

```go
func TestButtonGolden(t *testing.T) {
    component := display.Button(display.ButtonProps{Text: "Click me", Variant: display.ButtonPrimary})
    var buf bytes.Buffer
    component.Render(context.Background(), &buf)
    internal/golden.Assert(t, "button-primary", buf.String())
}
```

- CSS classes are normalized (sorted, deduplicated) before comparison
- Run `go test ./<pkg>/... -update` after intentional visual changes
- Review the `.golden` diff before committing

### Accessibility Testing

```go
func TestModalA11y(t *testing.T) {
    // Assert role="dialog" and aria-modal="true"
    // Assert focus trap (Tab cycles inside)
    // Assert focus restoration to trigger on close
    // Assert Escape key closes
    // Assert motion-reduce classes on all transitions
    // Assert aria-labelledby points to title element
}
```

### Behaviour-Driven Testing (BDD)

```go
var _ = Describe("Modal", func() {
    Describe("opening the modal", func() {
        It("moves focus into the dialog", func() { ... })
        It("traps Tab key within the dialog", func() { ... })
        It("prevents interaction with background content", func() { ... })
    })
    Describe("closing the modal", func() {
        It("restores focus to the trigger element", func() { ... })
        It("removes the dialog from the DOM", func() { ... })
    })
})
```

### What Radix/React Aria Teach About Testing

- **Test with real assistive technology** — VoiceOver, NVDA, JAWS, TalkBack. Not just
  automated checkers (which produce false positives even on correct implementations).
- **Test cross-browser** — focus management behaves differently in Safari vs Chrome.
- **Test touch devices** — hover states don't exist on mobile; long-press replaces
  right-click. Automated checkers miss these entirely.

---

## 16. Documentation & Blocks

### The Three-Tier Documentation Architecture (from shadcn/ui)

1. **Introduction** — Philosophy and principles (why this library exists)
2. **Components** — Individual reference (API, props, examples, live preview)
3. **Blocks** — Pre-composed layouts (dashboards, sidebars, login pages)

### Component Documentation Pattern

Every component doc page should have:
1. **One-line description** — what it is
2. **Installation** — one command to add it
3. **Usage** — minimal import + code snippet
4. **Examples** — multiple named variants with live preview + copy button
5. **API Reference** — prop tables with type and default value
6. **On This Page** sidebar — auto-generated table of contents

### The Blocks Concept

Blocks are **pre-composed multi-component layouts** — full-page or full-section patterns
that demonstrate how primitives compose:

| Block type | Example | What it shows |
|-----------|---------|--------------|
| Dashboard | stat cards + chart + data table + sidebar | Full app shell |
| Login page | card + form inputs + button + link | Authentication flow |
| Sidebar | collapsible nav + user menu + team switcher | Navigation patterns |
| Settings | tabs + form + save button | Form composition |

**Why blocks matter:** They teach consumers HOW to compose primitives. A consumer who
sees a working dashboard block learns more than reading 10 component API pages.

**For templ-components:** The `examples/demo/` directory is a start. Consider adding
formal "block" examples — a dashboard layout, a login page, a settings page — that
demonstrate component composition in real-world contexts.

---

## 17. Performance

### CSS Delivery (Server-Rendered Advantage)

Server-rendered HTML with Tailwind utility classes has inherent performance advantages:

- **No JavaScript hydration** — HTML is immediately interactive (progressive enhancement)
- **Critical CSS is the only CSS** — Tailwind generates only used utilities (tree-shaking)
- **No render-blocking SPA waterfall** — no `bundle.js` → parse → render → fetch API cycle
- **FOUC-free** — no flash of unstyled content (common in SPA hydration)

### CSS Performance (Tailwind v4)

- Only **used** utilities are generated (content scanning tree-shaking)
- **Logical properties** reduce output size (one property for LTR+RTL)
- `color-mix()` replaces opacity-specific utilities
- `@property` registered custom properties enable browser animation optimization
- `space-y-*` selector change (`> :not(:last-child)`) addresses performance on large pages

### Runtime Performance

- **GPU-accelerated animations** — only `transform` and `opacity` (see §10.3)
- **CSS containment** — `contain: layout style paint` isolates rendering
- **No client-side state stores** — HTML IS the state, no reconciliation overhead
- **HTMX partial swaps** — only the changed DOM fragment is transferred and rendered

### Bundle Size

templ-components has **zero JavaScript framework dependencies**. The only runtime deps
are:
- `tailwind-merge-go` — class merging (tiny Go library)
- `go-error-family` — error classification (errorpage module only)
- `templ` — the template runtime (already in the consumer's binary)

Compare with React-based libraries that add React + ReactDOM + Radix + CVA + clsx +
tailwind-merge — easily 100KB+ of JavaScript before a single component renders.

---

## 18. Gap Analysis: templ-components vs the Field

Based on deep comparison with shadcn/ui, Radix UI, React Aria, and templui:

### Gap 1: Semantic Token Layer (HIGH IMPACT)

**Current:** Components use primitive Tailwind names (`bg-blue-600`).
**Better:** Semantic tokens (`bg-tc-primary`) that consumers retheme without collisions.
**See:** §6 for the full migration plan.

### Gap 2: Compound Components for Complex Overlays (MEDIUM IMPACT)

**Current:** `Modal`, `Dropdown`, `Drawer` are monolithic — library owns the DOM.
**Better:** Compound pattern (`Trigger`, `Content`, `Close`, `Header`, `Footer`) for
maximum composition flexibility.
**See:** §3 for the architecture comparison.

### Gap 3: Native `<dialog>` Element (MEDIUM IMPACT)

**Current:** `Modal` uses div-based overlay with JS focus management.
**Better:** Native HTML5 `<dialog>` with `::backdrop`, top-layer rendering, built-in
focus trap. All modern browsers support it since 2022.
**Bonus:** CSS `@starting-style` enables zero-JS enter/exit animations on `<dialog>`.

templui already uses `<dialog>` — see their `dialog.templ`:
```go
<dialog class="fixed left-[50%] top-[50%] ...">
    {/* content */}
</dialog>
```

### Gap 4: RTL / Internationalization (MEDIUM IMPACT)

**Current:** Physical Tailwind utilities (`ml-`, `mr-`) used in some places.
**Better:** Logical utilities (`ms-`, `me-`) everywhere. Add RTL test cases.
**See:** §13 for the migration guide.

### Gap 5: Missing Components

Comparing component coverage across libraries:

| Missing component | Priority | Complexity | Notes |
|-------------------|----------|------------|-------|
| **Popover** | High | Medium | Compound pattern; positioning needed |
| **HoverCard** | Medium | Medium | Like Popover but hover-triggered |
| **Slider** | Medium | Medium | Range input with ARIA slider pattern |
| **Rating** | Low | Low | Star rating with keyboard support |
| **Carousel** | Low | High | Complex; consider if needed |
| **Calendar** | Medium | High | Full calendar grid; date-fns equivalent |
| **DataTable** | Medium | High | Sorting, filtering, pagination, virtualization |
| **ContextMenu** | Low | Medium | Right-click menu; compound pattern |

### Gap 6: Motion Tokens (LOW IMPACT)

**Current:** Durations and easings are inline strings per component.
**Better:** Centralized motion tokens (`--ease-enter`, `--duration-standard`) for consistency.

### Gap 7: Blocks / Composition Examples (LOW IMPACT)

**Current:** `examples/demo/` shows basic composition.
**Better:** Formal "blocks" — dashboard, login, settings, sidebar layouts — that teach
composition patterns.

---

## 19. The Master Checklist

A superb web UI component library, synthesized from all research:

### Architecture
- [ ] Type-safe variants — typed enums + map lookups, never raw strings
- [ ] Make impossible states unrepresentable — typed enums, fallback values, zero panics
- [ ] Compound composition for complex widgets — sub-parts mirroring DOM structure
- [ ] Single source of truth — shared constants for SVG paths, class strings, style data
- [ ] ComponentProps interface — every props struct satisfies it via BaseProps embedding
- [ ] Register new types in contract test — `internal/contract/component_props_test.go`

### Theming
- [ ] Three-layer token stack — primitive → semantic → component tokens
- [ ] Semantic token layer — `bg-tc-primary` not `bg-blue-600`
- [ ] CSS variable theming — consumers retheme without Go/JS changes
- [ ] Dark mode — class strategy, `gray-*` exclusively (no mixed palettes)
- [ ] Beautiful defaults — components look great out of the box
- [ ] Zero lock-in — every visual decision overridable

### Hypermedia Philosophy
- [ ] HATEOAS-first — HTML is source of truth, server controls available actions
- [ ] Native HTML preferred — `<details>`, `<form>`, `<dialog>` over JS reimplementations
- [ ] Progressive enhancement — JS reads state from HTML, enhances rather than replaces
- [ ] HTMX-native — components emit `hx-*` attributes naturally

### Accessibility
- [ ] Native HTML first — `<button>` not `<div role="button">`
- [ ] APG keyboard patterns — every interactive widget follows WAI-ARIA Authoring Practices
- [ ] Focus management — trap in modals, restore on close, logical tab order
- [ ] Motion safety — `motion-reduce:*` on every transition/animation
- [ ] ARIA correct — roles, states, properties; "no ARIA is better than bad ARIA"
- [ ] Screen reader text — `.sr-only` where visual text isn't enough
- [ ] AriaLabel propagation — every component with BaseProps propagates it

### Forms
- [ ] Progressive enhancement — forms work without JavaScript
- [ ] Server-side validation — source of truth, never trust client-only
- [ ] Accessible errors — `aria-invalid`, `aria-describedby`, error summary
- [ ] CSRF token — hidden input (survives page reloads)
- [ ] HTMX integration — `hx-post`, `hx-target`, `hx-swap`, `hx-sync`
- [ ] State preservation — re-populate values after failed submission

### Motion Design
- [ ] Purposeful animation — only where it adds value, not everywhere
- [ ] Professional easing — custom cubic-bezier, not `linear` or default `ease`
- [ ] Asymmetric enter/exit — different curves and durations
- [ ] Only transform/opacity — GPU-accelerated properties only
- [ ] Duration guidelines — 150–350ms for common interactions
- [ ] Motion tokens — centralized duration/easing definitions
- [ ] `prefers-reduced-motion` — always handled with safe alternatives

### Security
- [ ] CSP-safe by construction — nonce on every inline script
- [ ] No eval/handlers — no `eval()`, no inline event handlers, no `javascript:` URLs
- [ ] XSS prevention — `strconv.Quote()` for IDs in JS, `templ.SafeURL` for hrefs
- [ ] No external assets — JS travels with component (inline + nonce)

### Internationalization
- [ ] Logical CSS properties — `ms-`/`me-` not `ml-`/`mr-`
- [ ] RTL testing — `dir="rtl"` test cases
- [ ] Icon mirroring — directional icons flip in RTL
- [ ] Locale-aware formatting — date/number formatting by locale

### Testing
- [ ] Golden tests — exact rendered HTML matches snapshots
- [ ] A11y tests — ARIA, roles, keyboard, motion-reduce
- [ ] BDD tests — behaviour specs (user-visible, not markup)
- [ ] Edge-case tests — empty inputs, unknown enums, ID collisions
- [ ] Example tests — godoc `ExampleXxx()` compiles and renders
- [ ] Coverage tests — private helpers and branches

### Developer Experience
- [ ] Zero-config defaults — `DefaultComponentProps()` for meaningful non-zero defaults
- [ ] Override without forking — `Class` prop + `Attrs` map + CSS variables + slots
- [ ] Godoc examples — every component has a runnable `ExampleXxx()`
- [ ] Clear error messages — graceful degradation, never panic in render path
- [ ] Progressive disclosure — simple API surface, deep docs for advanced cases

### Distribution
- [ ] Versioned releases — semver, one-commit release convention
- [ ] Generated code committed — `*_templ.go` in the repo
- [ ] Granular adoption — multi-module workspace for partial adoption
- [ ] CHANGELOG warm — `[Unreleased]` always has entries

---

## Appendix A: Side-by-Side Button Comparison

### templ-components

```go
type ButtonType string
const (
    ButtonPrimary   ButtonType = "primary"
    ButtonSecondary ButtonType = "secondary"
    ButtonDanger    ButtonType = "danger"
)

type ButtonProps struct {
    utils.BaseProps
    Text     string
    Variant  ButtonType
    Size     ButtonSize
    Href     string
    Disabled bool
    Icon     templ.Component
}

@display.Button(display.ButtonProps{
    Text: "Save",
    Variant: display.ButtonPrimary,
    Size: display.ButtonSizeMD,
})
```

### templui

```go
type Variant string
const (
    VariantDefault Variant = "default"
    VariantOutline Variant = "outline"
)

type Props struct {
    ID         string
    Class      string
    Attributes templ.Attributes
    Variant    Variant
    Size       Size
    Href       string
    Disabled   bool
}

@button.Button(button.Props{
    Variant: button.VariantDefault,
}) {
    Save
}
```

### shadcn/ui (TypeScript)

```tsx
const buttonVariants = cva("inline-flex ...", {
  variants: {
    variant: { default: "...", destructive: "...", outline: "..." },
    size: { default: "h-9 px-4", sm: "h-8", icon: "size-9" },
  },
})

<Button variant="default" size="default">Save</Button>
```

**Key difference:** templ-components embeds `BaseProps` (giving ID/Class/Attrs/AriaLabel/Nonce
to every component); templui does NOT (each props struct repeats ID/Class/Attributes
manually). This means templ-components has a stronger, more consistent API contract.

---

## Appendix B: Source Comparison Matrix

| Dimension | templ-components | templui | shadcn/ui |
|-----------|-----------------|---------|-----------|
| **Language** | Go + templ | Go + templ | TypeScript + React |
| **Rendering** | Server-side | Server-side | Client-side (CSR/SSR) |
| **CSS** | Tailwind v4 | Tailwind v3/v4 | Tailwind v4 |
| **Distribution** | Go module | Module + CLI copy | npm CLI copy-paste |
| **BaseProps** | Embedded (consistent) | Not embedded (repeated) | React props (per-component) |
| **Variant lookup** | Map + `utils.Lookup` | Switch statement | CVA |
| **Class merge** | `utils.Class` (tailwind-merge-go) | `utils.TwMerge` | `cn()` (clsx + tailwind-merge) |
| **Theming** | Primitive Tailwind names | Semantic tokens (`bg-primary`) | Semantic OKLCH tokens |
| **Dark mode** | `gray-*` exclusively | `zinc-*` / mixed | Semantic pairs |
| **JS strategy** | Inline + nonce + singleton | External `.js` + OnceHandle | React lifecycle |
| **CSP** | Safe by construction | External scripts | React (N/A) |
| **Compound components** | Monolithic (mostly) | Compound (Dialog, Accordion) | Compound (Radix-based) |
| **Accordion** | Custom JS | Native `<details>` (zero JS) | Radix (JS required) |
| **Dialog** | Div overlay + JS | Native `<dialog>` | Radix Dialog (JS) |
| **A11y** | Strong (motion-reduce, focus, ARIA) | ARIA, keyboard | Radix primitives (gold standard) |
| **Forms** | 16 components + ValidationSummary | Form.Item + Message | React Hook Form integration |
| **HATEOAS** | Strong (HTML-first, HTMX-native) | Strong (HTML-first) | Weak (SPA, JSON-driven) |
| **Component count** | 82 + 101 icons | 43 | 50+ |
| **Testing** | 7-layer matrix | Unknown | Vitest + browser + E2E |
| **i18n/RTL** | Physical properties (gap) | Physical properties | `migrate rtl` command |
| **Motion tokens** | Inline (gap) | Inline | CSS variables |

---

## Appendix C: Key URLs

- [shadcn/ui](https://ui.shadcn.com/) — landing page
- [shadcn/ui docs](https://ui.shadcn.com/docs) — design philosophy, theming, CLI
- [shadcn/ui blocks](https://ui.shadcn.com/blocks) — pre-composed layouts
- [shadcn-ui/ui repo](https://github.com/shadcn-ui/ui) — monorepo, registry
- [Radix UI primitives](https://www.radix-ui.com/primitives) — headless components
- [React Aria](https://react-spectrum.adobe.com/react-aria/) — behavior hooks
- [templui repo](https://github.com/templui/templui) — Go templ competitor
- [HATEOAS essay](https://htmx.org/essays/hateoas/) — hypermedia theory
- [WAI-ARIA APG](https://www.w3.org/WAI/ARIA/apg/) — accessible patterns
- [Design Tokens spec](https://www.designtokens.org/) — DTCG format
- [CVA](https://cva.style/) — class-variance-authority
- [Tailwind v4](https://tailwindcss.com/blog/tailwindcss-v4) — CSS-first config
