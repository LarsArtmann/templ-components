---
name: templ-components
description: Authoring playbook for the templ-components Go UI library (templ + HTMX + Tailwind v4). Use this skill whenever working in the github.com/larsartmann/templ-components repo, or when adding, editing, reviewing, or debugging any templ component, props struct, typed enum, style lookup, icon path, or HTMX helper here. Also trigger when the user asks to "add a component", "fix a component", extend the library, wire up a new icon, add a typed enum, set up Tailwind v4 scanning for templ, integrate HTMX loading/error handling, or questions how a component should be structured, themed, made accessible, CSP-safe, or dark-mode aware in this codebase. Load this skill BEFORE writing or editing any .templ file in this repo.
metadata:
  tags: templ, templ-components, htmx, tailwind, tailwind-v4, go, ui, components, accessibility, csp, dark-mode, server-rendered
---

# templ-components

A server-rendered, type-safe, accessible, CSP-ready Go component library built on
[templ](https://templ.guide), [HTMX](https://htmx.org), and Tailwind v4.

This skill serves **two audiences**:

- **Consumers** (most readers): you want to know what components exist, how to
  adopt them, and how to wire them into your app. → Read **Part 1** below.
- **Authors** (maintainers): you want to add, edit, or review components.
  → Read **Part 2**, then keep Part 1 honest.

When a question is about _what exists_ or _how to use it_, Part 1 answers. When
it's about _how to make a new component fit the library_, Part 2 answers.

---

# Part 1: Consumer Guide

## Component catalogue

88 components across 9 packages + 101 icons. If you're about to hand-roll
something, check this table first — 4 of the top 6 consumer "missing components"
already existed.

### By use case (start here)

Don't know what to look for? Find your page type:

| You're building...             | Reach for                                                                                                                                                                                                                          |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Dashboard / metrics page**   | `Grid`, `StatCard`, `Card`, `ProgressBar`, `SkeletonCardGrid`, `PageHeader`                                                                                                                                                        |
| **List / table page**          | `Table` (`Flush` for card nesting, `CellPadding` for compact rows, `Table.Body` for custom rows, `Row.Href` for clickable rows), `Badge`, `StatusBadge`, `Avatar`, `Pagination`, `LoadMore`, `EndOfList`, `EmptyState`, `ListNote` |
| **Detail page**                | `Card`, `DefinitionList`, `DefinitionGrid`, `Tabs`, `PageHeader`, `Breadcrumbs`                                                                                                                                                    |
| **Settings / data-entry form** | `Form`, `Input`, `Select` (supports `Groups` for optgroups), `Textarea`, `Toggle`, `Checkbox`, `RadioGroup`, `ValidationSummary`                                                                                                   |
| **Filter bar (horizontal)**    | Thin custom helper (see `docs/recipes/horizontal-filter-bar.md`) — `forms.Form` targets vertical                                                                                                                                   |
| **Feedback / notifications**   | `Toast`, `ToastContainer`, `Alert`, `Spinner`, `ProgressBar`, `GlobalErrorHandling`                                                                                                                                                |
| **Navigation**                 | `Nav`, `SimpleNav`, `SidebarNav`, `Breadcrumbs`, `Pagination`, `MobileMenu`                                                                                                                                                        |
| **Modal / overlay**            | `Modal`, `Drawer`, `Dropdown`, `Tooltip`, `Popover`, `Accordion`                                                                                                                                                                   |
| **Error pages**                | `ErrorPage`, `NotFound404`, `ErrorDetail`, `ErrorAlert`, `ErrorHandler`                                                                                                                                                            |
| **Full page shell**            | `Base`, `Minimal`, `ThemeScript`, `ThemeToggle`, `Script`                                                                                                                                                                          |

### By package (import path reference)

#### `display` — 26 components

| Component          | Signature                                   | One-liner                                                                                                                                        |
| ------------------ | ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Card`             | `Card(props CardProps)`                     | Bordered card with title, subtitle, footer, header action, 4 padding sizes                                                                       |
| `SimpleCard`       | `SimpleCard(props SimpleCardProps)`         | Minimal card — children only, no header/footer                                                                                                   |
| `StatCard`         | `StatCard(props StatCardProps)`             | Dashboard metric card with value, label, change, trend, icon, optional `Href` link                                                               |
| `Grid`             | `Grid(props GridProps)`                     | Responsive grid — typed `GridCols` enum, `GridGap` enum, `ContainerResponsive`                                                                   |
| `Badge`            | `Badge(props BadgeProps)`                   | Compact status label — 7 types, 3 sizes, pill, dot, optional `Href`                                                                              |
| `StatusBadge`      | `StatusBadge(status string)`                | Auto-maps ~20 status strings to badge types                                                                                                      |
| `Button`           | `Button(props ButtonProps)`                 | Button or link button — variants, sizes, icons, HTMX attrs                                                                                       |
| `Avatar`           | `Avatar(props AvatarProps)`                 | Image avatar with fallback initials, sizes, status dot                                                                                           |
| `Modal`            | `Modal(props ModalProps)`                   | Accessible dialog — focus trap, Escape, backdrop, 5 sizes                                                                                        |
| `Drawer`           | `Drawer(props DrawerProps)`                 | Side panel — left/right slide, focus trap, Escape, backdrop                                                                                      |
| `Dropdown`         | `Dropdown(props DropdownProps)`             | Button-triggered menu — links, buttons, keyboard nav                                                                                             |
| `Tooltip`          | `Tooltip(props TooltipProps)`               | Hover tooltip — 4 positions, arrow, touch support                                                                                                |
| `Popover`          | `Popover(props PopoverProps)`               | Button-triggered floating panel — 4 positions, arbitrary content, click-outside dismiss                                                          |
| `Accordion`        | `Accordion(props AccordionProps)`           | Collapsible sections — open/closed state, keyboard nav                                                                                           |
| `Table`            | `Table(props TableProps)`                   | Responsive data table — striping, hover, caption, bordered, `Body` slot, `Row.Href`, `Flush` (card nesting), `CellPadding` (comfortable/compact) |
| `Tabs`             | `Tabs(props TabsProps)`                     | Tabbed interface — underline/pills, optional client-side JS                                                                                      |
| `EmptyState`       | `EmptyState(props EmptyStateProps)`         | Empty-data placeholder — icon, title, description, action                                                                                        |
| `SimpleEmptyState` | `SimpleEmptyState(message string)`          | Minimal empty state — text only                                                                                                                  |
| `PageHeader`       | `PageHeader(props PageHeaderProps)`         | Page title block — title, subtitle, breadcrumb, action slots                                                                                     |
| `DefinitionList`   | `DefinitionList(props DefinitionListProps)` | Two-column `<dl>` key/value list                                                                                                                 |
| `ListNote`         | `ListNote(props ListNoteProps)`             | "Showing N of M" truncation notice                                                                                                               |
| `CopyButton`       | `CopyButton(props CopyButtonProps)`         | Clipboard copy button or link — CSP-safe, "Copied!" feedback, optional `Href` variant                                                            |
| `RelativeTime`     | `RelativeTime(props RelativeTimeProps)`     | `<time datetime>` with relative text ("2 hours ago")                                                                                             |
| `CountBadge`       | `CountBadge(props CountBadgeProps)`         | Icon + notification count overlay — overflow "N+"                                                                                                |
| `DefinitionGrid`   | `DefinitionGrid(props DefinitionGridProps)` | Responsive grid of term-detail cards                                                                                                             |
| `Image`            | `Image(props ImageProps)`                   | Lazy-loaded `<img>` with CSP-safe fallback, optional `Rounded` for circular                                                                      |

#### `forms` — 16 components

| Component           | Signature                                         | One-liner                                                                                  |
| ------------------- | ------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `Input`             | `Input(props InputProps)`                         | Typed input (text, email, password, search, url, tel, number) with label, error, help text |
| `Checkbox`          | `Checkbox(props CheckboxProps)`                   | Checkbox with label, help text, error                                                      |
| `Select`            | `Select(props SelectProps)`                       | Select dropdown with typed options, disabled, selected, optional `Groups` for optgroups    |
| `Textarea`          | `Textarea(props TextareaProps)`                   | Multi-line text input with label, rows, max length                                         |
| `Toggle`            | `Toggle(props ToggleProps)`                       | Toggle switch with label, required, error, help text                                       |
| `Radio`             | `Radio(props RadioProps)`                         | Radio button with label                                                                    |
| `RadioGroup`        | `RadioGroup(props RadioGroupProps)`               | Group of radios with options                                                               |
| `Combobox`          | `Combobox(props ComboboxProps)`                   | Autocomplete input with filterable dropdown                                                |
| `DatePicker`        | `DatePicker(props DatePickerProps)`               | Date input with label                                                                      |
| `FileInput`         | `FileInput(props FileInputProps)`                 | File upload input with label, accept types                                                 |
| `Form`              | `Form(props FormProps)`                           | Form wrapper with action, method, CSRF token                                               |
| `Label`             | `Label(forID, text string, required bool)`        | Form label element                                                                         |
| `FieldError`        | `FieldError(fieldID, message string)`             | Inline field error message                                                                 |
| `FormFieldWrapper`  | `FormFieldWrapper(props FormFieldProps)`          | Wraps inputs with label, error, help text                                                  |
| `InputGroup`        | `InputGroup(props InputGroupProps)`               | Input with prefix/suffix                                                                   |
| `ValidationSummary` | `ValidationSummary(props ValidationSummaryProps)` | Accessible error summary with icon, count, linked fields                                   |

#### `feedback` — 13 components

| Component          | Signature                                   | One-liner                                           |
| ------------------ | ------------------------------------------- | --------------------------------------------------- |
| `Alert`            | `Alert(props AlertProps)`                   | Inline alert — 4 types, dismissible, icon           |
| `InlineError`      | `InlineError(message string)`               | Compact inline error                                |
| `InlineSuccess`    | `InlineSuccess(message string)`             | Compact inline success                              |
| `Toast`            | `Toast(props ToastProps)`                   | Toast notification — 4 types, dismissible, duration |
| `ToastContainer`   | `ToastContainer(nonce string)`              | Fixed container for toasts (include once per page)  |
| `Spinner`          | `Spinner(props SpinnerProps)`               | Loading spinner — 3 sizes, custom color             |
| `LoadingOverlay`   | `LoadingOverlay(props LoadingOverlayProps)` | Full-screen loading overlay with spinner            |
| `InlineLoading`    | `InlineLoading(message string)`             | Inline "Saving..." indicator                        |
| `Skeleton`         | `Skeleton(variant SkeletonVariant)`         | Pulsing placeholder — 7 variants                    |
| `SkeletonGroup`    | `SkeletonGroup(variants []SkeletonVariant)` | Multiple skeletons under single `role="status"`     |
| `SkeletonCardGrid` | `SkeletonCardGrid(count int)`               | N skeleton cards in responsive grid — loading state |
| `ProgressBar`      | `ProgressBar(props ProgressBarProps)`       | Progress indicator — 3 sizes, indeterminate, label  |
| `StepIndicator`    | `StepIndicator(props StepIndicatorProps)`   | Horizontal/vertical step progress                   |

#### `layout` — 5 components

| Component     | Signature                              | One-liner                                                   |
| ------------- | -------------------------------------- | ----------------------------------------------------------- |
| `Base`        | `Base(props PageProps)`                | Full HTML5 shell — head, meta, theme, HTMX, CSS auto-inject |
| `Minimal`     | `Minimal(props MinimalProps)`          | Minimal HTML doc — no dependencies, for static/PDF          |
| `ThemeScript` | `ThemeScript(nonce string)`            | Dark mode script — prevents FOUC, include in `<head>`       |
| `ThemeToggle` | `ThemeToggle(ariaLabel, nonce string)` | Dark/light toggle button with sun/moon icons                |
| `Script`      | `Script(nonce, src string, attrs)`     | CSP-safe `<script src>` — auto-injects nonce                |

#### `navigation` — 12 components

| Component          | Signature                                               | One-liner                                                                    |
| ------------------ | ------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `Nav`              | `Nav(props NavProps)`                                   | Full nav bar — brand, links, right items slot, mobile menu                   |
| `SimpleNav`        | `SimpleNav(props SimpleNavProps)`                       | Simplified nav — text brand, links, `RightItems` slot                        |
| `NavLink`          | `NavLink(props NavLinkProps, currentPath string)`       | Single nav link with active detection                                        |
| `MobileNavLink`    | `MobileNavLink(props NavLinkProps, currentPath string)` | Mobile menu link                                                             |
| `MobileMenu`       | `MobileMenu(links, currentPath, nonce, menuID)`         | Responsive mobile menu                                                       |
| `MobileMenuToggle` | `MobileMenuToggle(show bool, menuID string)`            | Hamburger toggle button                                                      |
| `Pagination`       | `Pagination(props PaginationProps)`                     | Page-number pagination — ellipsis, prev/next, SEO rel                        |
| `Breadcrumbs`      | `Breadcrumbs(props BreadcrumbsProps)`                   | Breadcrumb trail — separator, JSON-LD                                        |
| `SidebarNav`       | `SidebarNav(props SidebarNavProps)`                     | Vertical sidebar — brand, icon+label items, footer slot                      |
| `LoadMore`         | `LoadMore(props LoadMoreProps)`                         | Cursor-based "Load more" button — hx-get, hx-swap, optional `InfiniteScroll` |
| `EndOfList`        | `EndOfList(props EndOfListProps)`                       | "You've reached the end" indicator — companion to LoadMore/Pagination        |
| `Footer`           | `Footer(props FooterProps)`                             | Simple footer with copyright                                                 |

#### `htmx` — 7 components

| Component              | Signature                                                 | One-liner                   |
| ---------------------- | --------------------------------------------------------- | --------------------------- |
| `LoadingIndicator`     | `LoadingIndicator(spinner templ.Component)`               | Fixed HTMX loading overlay  |
| `InlineLoadingOverlay` | `InlineLoadingOverlay(id string, spinner)`                | Per-target loading overlay  |
| `LoadingButton`        | `LoadingButton(defaultText, loadingText string, spinner)` | Button with loading state   |
| `ConfirmDelete`        | `ConfirmDelete(props ConfirmDeleteProps)`                 | HTMX delete confirmation    |
| `SwapOOB`              | `SwapOOB(props SwapOOBProps)`                             | HTMX out-of-band swap       |
| `CSRFToken`            | `CSRFToken(token string)`                                 | Hidden CSRF input           |
| `GlobalErrorHandling`  | `GlobalErrorHandling(cfg ErrorHandlingConfig)`            | HTMX error → toast pipeline |

#### `errorpage` — 4 components + 6 constructors + handler

| Component / Function           | Signature                             | One-liner                                                                   |
| ------------------------------ | ------------------------------------- | --------------------------------------------------------------------------- |
| `ErrorPage`                    | `ErrorPage(props ErrorPageProps)`     | Full-page error display                                                     |
| `NotFound404`                  | `NotFound404(props NotFound404Props)` | Dedicated 404 page — hero numeral, search, links, configurable `LinksTitle` |
| `ErrorDetail`                  | `ErrorDetail(props ErrorDetailProps)` | Inline error card                                                           |
| `ErrorAlert`                   | `ErrorAlert(props ErrorAlertProps)`   | Family-aware alert                                                          |
| `ErrorHandler`                 | `ErrorHandler(err, cfg) http.Handler` | go-error-family aware HTTP handler                                          |
| `WriteError`                   | `WriteError(w, r, err, nonce)`        | One-call error page from any handler                                        |
| `FromError`                    | `FromError(err) ErrorPageProps`       | Extract family/code/why/fix from error                                      |
| `NotFound` ... `InternalError` | 6 constructors                        | Pre-built error page props by HTTP family                                   |

#### `icons` — 101 icons + 2 functions

| Function                                        | One-liner                                          |
| ----------------------------------------------- | -------------------------------------------------- |
| `Icon(name, class)`                             | Inline SVG icon by name (100 path icons + Spinner) |
| `IconWithStrokeWidth(name, class, strokeWidth)` | Same, custom stroke width (default 1.5)            |
| `IconPathData(name) []string`                   | Raw path data for custom SVG wrapper               |
| `IconPathJS(name) string`                       | Path data formatted for JS injection               |

#### `utils` — shared types + helpers

| Function                   | One-liner                                                        |
| -------------------------- | ---------------------------------------------------------------- |
| `BaseProps`                | Embed in every props struct — ID, Class, Attrs, AriaLabel, Nonce |
| `ComponentProps`           | Interface: `GetBaseProps()` / `SetBaseProps()`                   |
| `Class(classes ...string)` | Tailwind class merge (tailwind-merge-go, mutex-protected)        |
| `Lookup(m, key, fallback)` | Map lookup with fallback for style variants                      |
| `Ternary(cond, a, b)`      | Functional ternary for templ expressions                         |
| `EnsureID(prefix, id)`     | Auto-generate DOM-safe IDs (crypto/rand)                         |
| `Version`                  | Library version string (matches CHANGELOG)                       |

## Quick start: adopting the library

```bash
go get github.com/larsartmann/templ-components@latest
```

### Minimal page

```templ
@layout.Base(layout.DefaultPageProps()) {
    @navigation.SimpleNav(navigation.SimpleNavProps{
        BrandText: "MyApp",
        Links: []navigation.NavLinkProps{
            {Href: "/", Text: "Home"},
        },
        RightItems: @layout.ThemeToggle("Toggle theme", ""),
        CurrentPath: "/",
    })
    <main class="mx-auto max-w-7xl px-4 py-8">
        @display.Grid(display.GridProps{Cols: display.GridCols3}) {
            for _, u := range users {
                @display.StatCard(display.StatCardProps{Label: u.Name, Value: u.Count, Href: "/users/" + u.ID})
            }
        }
    </main>
    @feedback.ToastContainer("")
    @htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
}
```

### Suppressing auto-injected `<head>` tags

`DefaultPageProps()` auto-injects two tags. Override them in your `PageProps`:

```go
props := layout.DefaultPageProps()
props.HTMXVersion = ""  // suppress <script src="...htmx.org@...">
props.CSSPath = ""      // suppress <link rel="stylesheet" href="/app.css">
```

To use a different HTMX CDN or self-host: `props.HTMXCDN = "https://unpkg.com"`.

### CSP-safe script tags

Use `layout.Script` instead of raw `<script>` — the nonce can never be forgotten:

```templ
@layout.Script(nonce, "/static/app.js", templ.Attributes{"defer": true})
```

### Theming

Components emit standard Tailwind classes (`bg-blue-600`). Override globally in your CSS:

```css
@theme {
  --color-blue-600: #4f46e5; /* indigo instead of blue */
}
```

Dark mode: `@custom-variant dark (&:where(.dark, .dark *))` + `layout.ThemeScript()` + `layout.ThemeToggle()`. `color-scheme: light` on `:root`, `color-scheme: dark` on `.dark`.

### Dark mode checklist (for new components)

Every new component MUST pass `TestDarkModeCompliance` + `TestDarkModeSemanticColors`. Checklist:

1. **Neutral colors** — every `text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `ring-gray-*` MUST have a `dark:` variant on the same line.
2. **Semantic colors** — every `bg-blue-600`, `bg-red-600`, `text-blue-600`, `text-green-500`, etc. MUST have a `dark:` variant. Convention: `-600` light → `-500` dark for backgrounds, `-400` for text.
3. **Focus rings** — `focus:ring-blue-500` needs `dark:focus:ring-blue-400`. `focus-visible:outline-blue-600` needs `dark:focus-visible:outline-blue-500`.
4. **Ring offset** — `focus:ring-offset-2` needs `dark:focus:ring-offset-gray-900` (default offset is white, invisible in dark).
5. **Shadows** — overlay/card shadows (`shadow-xl`, `shadow-lg`) should have `dark:shadow-black/20` for visibility on dark backgrounds.
6. **Palette** — use `gray-*` exclusively. No `slate-*`, `zinc-*`, `neutral-*`, `stone-*`.

See `docs/adr/0011-dark-mode-convention.md` for the full convention.

### Consumer `_templ.go` guidance

**Consumers should gitignore `*_templ.go`** and generate at build time via `templ generate`.
**The library commits them** because the Go module proxy serves source as-is and does not
run `templ generate`. This is the standard pattern for publishable templ packages.

## Recipes

| Recipe                                                | When to read                                        |
| ----------------------------------------------------- | --------------------------------------------------- |
| `docs/migration/play-cdn-to-tailwind-v4.md`           | Migrating from Tailwind Play CDN to CSS-first build |
| `docs/recipes/server-rendered-htmx-error-feedback.md` | Wiring HTMX error feedback (toast/alert/page)       |
| `docs/recipes/horizontal-filter-bar.md`               | Horizontal HTMX filter bar vs `forms.Form`          |
| `docs/tailwind-v4-adoption-guide.md`                  | Full Tailwind v4 setup with `@source` scanning      |
| `docs/icons-only-adoption.md`                         | Adopting just the `icons` package (CSS-agnostic)    |
| `docs/recipes/theme-bridge.md`                        | Remap library colors to custom semantic palette     |

## How to know if a component already exists

Before hand-rolling HTML, check in this order:

1. **The "By use case" table above** — find your page type first.
2. **The per-package catalogue below it** — full signatures and one-liners.
3. **`README.md` component catalogue** — code examples grouped by package.
4. **`pkg.go.dev/github.com/larsartmann/templ-components`** — full API reference.
5. **`grep -r "templ [A-Z]" --include="*.templ"`** — find every component definition.

### Consumer tip: track adoption in your AGENTS.md

The #1 discoverability gap across consumers is not knowing which library
components are already adopted vs hand-rolled in your own project. Keep a
grep-able table in your consumer project's `AGENTS.md`:

```markdown
## templ-components adoption

| Library component      | Status  | Where                   |
| ---------------------- | ------- | ----------------------- |
| `display.Grid`         | adopted | dashboard.templ         |
| `display.Table`        | custom  | filters.templ (thinner) |
| `forms.Form`           | custom  | filterForm helper       |
| `feedback.ProgressBar` | adopted | backfill.templ          |
```

This lets every AI session and developer quickly audit what's adopted, what's
hand-rolled, and where the gaps are.

---

# Part 2: Authoring Playbook

## Principles (the why behind every rule)

- **Make invalid states unrepresentable.** Variants are typed enums, not strings; lookups are
  maps with explicit fallbacks; the only permitted runtime panic is the single developer
  data-integrity check in `icons` (stray `|` in a path). Everything else degrades gracefully.
  This is a _library_: a panic in a consumer's render path is a bug we shipped.
- **HATEOAS-first ([htmx.org/essays/hateoas](https://htmx.org/essays/hateoas/)).** HTML is the source of truth.
  Prefer native HTML (`<details>`, forms, links) over scripts. When JavaScript enhances the hypermedia
  (modal, drawer, tooltip, relative time auto-refresh), it reads state FROM HTML attributes
  (`data-tc-*`, `datetime`) — progressive enhancement, not SPA replacement. All JS is guarded by a
  global singleton flag so it is idempotent across HTMX re-renders.
- **Accessibility is not optional.** ARIA attributes, roles, keyboard nav, `motion-reduce:*`
  on every transition/animation, and screen-reader text are part of "done". A component that
  renders but is unusable from a keyboard is not finished.
- **CSP-safe by construction.** Every inline `<script>` carries `nonce={ props.Nonce }`. No
  `eval()`, no inline event handlers, no `javascript:` URLs. This is why consumers can ship a
  strict Content-Security-Policy without forking the library.
- **Single source of truth.** SVG path data lives in `internal/svg`; icon names map to paths in
  `icons/icon_paths.go`; feedback styles live in `feedback/styles.go`; Tailwind class strings are
  never duplicated when a shared constant exists (`cardShellClass`, `mutedTextClass`).
- **Composition over configuration.** Props embed `utils.BaseProps`; consumers override via
  `@theme` CSS variables, never by editing Go. Slots are `templ.Component` children.

## Flake commands (the canonical entry points)

Build automation lives in `flake.nix`, not a Makefile. Per repo policy you run these instead
of raw `go`/`golangci-lint` invocations — they wrap the pinned `templ` (v0.3.1020, matching
`go.mod`) so generation is reproducible.

| Command              | What it does                                                              |
| -------------------- | ------------------------------------------------------------------------- |
| `nix run .#build`    | Regenerate `*_templ.go` + `go build ./...`                                |
| `nix run .#test`     | `go test ./... -count=1 -race`                                            |
| `nix run .#lint`     | `golangci-lint` across all non-example packages                           |
| `nix run .#coverage` | Tests with `-coverprofile` + summary line                                 |
| `nix run .#verify`   | **Generate + build + test + lint in one shot — this is the "done" check** |

Run `nix run .#verify` before considering any component work finished. The full equivalent
manual command (only if you have no Nix) is documented in `AGENTS.md` and `CONTRIBUTING.md`.

## Process (run this when touching any `.templ` file)

1. **Enter the pinned dev shell before generating.** Run `nix develop` (or use the `nix run .#*`
   apps, which pull the same pinned `templ`). The system `templ` binary may be v0.3.1036 and will
   produce a cosmetic import-style diff across all 51 generated files. See `AGENTS.md` "templ
   Version Pin" — do not bump `go.mod` to an unreleased version.
2. **Regenerate and build** via `nix run .#build`. templ generates `*_templ.go` from `.templ`;
   Go never sees the source change until you generate.
3. **Commit `*_templ.go`.** This is a _library_, not an app. The Go module proxy serves source
   as-is from the Git tag and does **not** run `templ generate`. Missing generated files break
   every consumer. The `.gitignore` uses `!*_templ.go` to keep them tracked — never undo that.
   **BuildFlow gotcha:** the pre-commit `templ-generate` step re-appends `*_templ.go` to
   `.gitignore` on every run, which hides NEW generated files from `git status`. After adding a
   new component, run `git add -f yourcomponent_templ.go` to force-track it. Already-tracked
   files are unaffected.
4. **Test the full matrix for the package you touched** — `nix run .#test` (or `go test ./<pkg>/...`).
   Each package pairs golden, a11y, BDD, edge-case, example, and snapshot tests (see the
   testing table below). Golden files: run `go test ./<pkg>/... -update` after an intentional
   output change, then eyeball the `.golden` diff before committing.
5. **Lint** via `nix run .#lint` (`examples/` is excluded by `.golangci.yml`).
6. **Keep `[Unreleased]` warm.** Add the changelog entry in the same commit as the change. The
   release script refuses to cut a version with an empty `[Unreleased]`.
7. **Verify no new dependencies.** The allowed set is `templ`, `tailwind-merge-go`, and
   `go-error-family` (errorpage only). Anything else needs an explicit decision.

The single done-check: **`nix run .#verify`**.

## Component anatomy (the canonical shape)

Every interactive/visual component in this library follows this skeleton. When adding one,
copy the closest existing component and adapt — do not invent a new shape.

```
<package>/
  component.templ        # types, props, render template, private helpers
  component_templ.go     # GENERATED — commit it (git add -f if new)
  component_test.go      # unit + behaviour tests
  testdata/*.golden      # golden HTML snapshots
```

Inside `component.templ`:

1. **Typed enums** for every closed variant set.
   ```go
   type BadgeType string
   const (
       BadgePrimary BadgeType = "primary"
       BadgeSuccess BadgeType = "success"
       // ...
   )
   ```
2. **Size constants** use the `ComponentSize[SM|MD|LG]` suffix pattern.
3. **Props struct** embeds `utils.BaseProps` (the only known exception is `layout.PageProps`).
   ```go
   type BadgeProps struct {
       utils.BaseProps
       Text string
       Type BadgeType
       // ...
   }
   ```
   Embedding auto-satisfies the `utils.ComponentProps` interface via promoted
   `GetBaseProps()`/`SetBaseProps()` (pointer receivers, required by `recvcheck`).
4. **Default constructor** `DefaultComponentProps()` returns meaningful non-zero defaults.
5. **Render template** with a godoc example comment.
6. **Root element** propagates `props.ID`, `props.Class` (via `utils.Class(...)`),
   `props.Attrs`, and `props.AriaLabel`. See `display/badge.templ` for the exact pattern.
7. **Private lookup helpers** (`xxxClass()`, `xxxSizeClass()`) backed by package-level maps.
8. **Shared rendering** across 2+ components is extracted to a private `templ` sub-template
   (e.g. `overlayShell`, `dialogHeader`, `navLinkAnchor`, `statCardInner`, the six
   `errorpage/shared.templ` helpers). Lift duplication into a sub-template rather than copying.
9. **Register the new Props type in the contract inventory.** Open
   `internal/contract/component_props_test.go` and add `yourpackage.YourProps{}` to the
   `componentTypes()` slice (in the right package section). The test
   `TestAllComponentPropsSatisfyInterface` then enforces at CI time that your struct both
   embeds `utils.BaseProps` and satisfies `utils.ComponentProps` — catching silent contract
   breakage for consumers using generic wrappers. Forgetting this step is the #1 way a new
   component slips in without the BaseProps embed.

## Per-component testing checklist

Every new component MUST have all of these. "Assertion tests only" is not acceptable —
that's what made `layout.Script` the worst-tested component in the library.

```
[ ] golden_test.go      — exact rendered HTML matches .golden snapshot
[ ] a11y_test.go        — ARIA, roles, keyboard, motion-reduce, screen-reader text
[ ] rtl_check           — verify no physical properties (ml-, mr-, pl-, pr-, left-, right-)
[ ] bdd_test.go         — behaviour spec (user-visible behaviour, not markup)
[ ] edge_cases_test.go  — empty inputs, unknown enum values, ID collisions
[ ] example_test.go     — godoc ExampleXxx() compiles and renders
[ ] snapshot_test.go    — broader composition snapshot
[ ] coverage_*_test.go  — targeted coverage of private helpers and branches
```

Golden tests use `internal/golden.Assert(t, name, got)` with CSS-class normalization; pass
`-update` to regenerate after an intentional visual change, then review the diff.

## Decision trees

### Map lookup vs if-branch for a variant

- **Pure class/style data** (colors, sizes, padding) → map + `utils.Lookup(map, key, fallback)`.
  This is the dominant pattern: `badgeStyleMap`, `cardPaddingLookup`, `spinnerSizeLookup`,
  `gridColsLookup`, etc.
- **Structural DOM differences** (TabsVariant, DropdownPosition, TrendDirection, StatCardHref)
  → `if`-branch inside the template. Maps are for data, not for choosing which markup to emit.
- **Enum validation** → 15 enums use map+fallback (graceful); `InputType`,
  `ButtonHTMLType`, `FormMethod` fall back to HTML-spec defaults; only the icons path-data
  integrity check is allowed to panic.

### `utils.Class()` vs `utils.Lookup()` — when to use which

- **Merging multiple class strings** → `utils.Class("a", b)` (tailwind-merge resolves conflicts).
- **Selecting a style variant from a map** → `utils.Lookup(map, key, fallback)`.
- **Both** → `utils.Lookup()` first to pick the variant, then pass the result to `utils.Class()`.

### When to introduce a typed enum

If a field takes one of a fixed set of visual/behavioural variants, make it a typed enum
(`type X string` + consts) with a `Default` constant. Never accept a raw string for a variant
that has a known closed set — that reopens the "invalid state" you're trying to prevent.

### When to add an icon

Add the path to `iconPathData` in `icons/icon_paths.go` (single source). Multi-path icons use
a `|` separator. `iconPaths()` validates there are no empty segments. `allIconNames()` is
auto-generated from the map — never hand-maintain a separate icon-name list. Export raw path
strings via `icons.IconPathData` for consumers building their own `<svg>` wrapper.

### When interactivity needs JavaScript

Walk this ladder before writing JS:

1. Can native HTML do it? (`<details>`/`<summary>` for disclosure, `<form>` for submit,
   `:target`/`:checked` for toggles). If yes, stop.
2. If JS is unavoidable, write a single inline `<script nonce={ props.Nonce }>` that:
   - Guards with a global singleton flag (`window.tcXxxAttached`) so HTMX re-renders are
     idempotent — re-running the script must not double-bind handlers.
   - Uses event delegation on `document` where practical.
   - Handles Escape-to-dismiss, focus save/restore (`data-tc-prev-focus`), and click-outside.
3. Escape any IDs interpolated into JS with `strconv.Quote()` to prevent XSS (see
   `validateDropdownID`).

See `docs/adr/0005-js-attachment-patterns.md` for the rationale.

### When to generate an ID

Use `utils.EnsureID("prefix", props.ID)` when a component needs a stable DOM id for ARIA
wiring (Modal, Drawer, Dropdown, Accordion, Combobox). It returns the consumer's id unchanged,
else generates `tc-<prefix>-<16 hex>` via `crypto/rand` (collision-safe across HTMX loads).
Never invent IDs with `time.Now()` alone — predictable under concurrency.

### Theming decisions

- Components emit **standard Tailwind utility classes** (`bg-blue-600`, `text-gray-900`).
- Consumers override colors by setting `@theme { --color-blue-600: #...; }` in their CSS —
  **no Go code change required**. Do not add `Variant`/`Color` props for theming; that defeats
  the CSS-variable model.
- Semantic aliases live in `templ-components-theme.css` (`bg-tc-primary`, `text-tc-danger`).
- Dark mode uses the class strategy: `@custom-variant dark (&:where(.dark, .dark *))`,
  toggled by `layout.ThemeScript()` + `layout.ThemeToggle()`.

## Mandatory conventions (these have no exceptions)

- **Dark mode colors:** `gray-*` exclusively. Never mix `slate-*` and `gray-*` in one
  component — the inconsistency shows in dark mode.
- **Motion safety:** every transition gets
  `motion-reduce:transition-none motion-reduce:duration-0`; every animation gets
  `motion-reduce:animate-none`. If you forget this, you fail the a11y tests.
  **Use the shared motion constants** (`transitionFast`, `transitionNormal`,
  `transitionColors`, `transitionTransform`) from `display/shared.go` instead of
  inline timing strings — they guarantee consistent durations and built-in
  `motion-reduce` fallbacks. `transitionColors` for hover/active color changes,
  `transitionNormal` for overlay panels, `transitionTransform` for sliding/rotating.
- **Class merging:** always go through `utils.Class(...)` so tailwind-merge resolves conflicts
  and consumer overrides win. The only exception is `templ.KV` conditionals where the templ
  runtime must comma-join. `utils.Class` is mutex-protected; do not bypass it.
- **CSP nonce:** every inline script takes `nonce={ props.Nonce }`. For external scripts, use
  `layout.Script(nonce, src, attrs)` — it auto-injects the nonce.
- **Card shell:** use the shared `cardShellClass` constant; `SimpleCard` composes through
  `Card` internally. `StatCard` also uses `cardShellClass`.
- **Muted text:** use `mutedTextClass` (`text-sm text-gray-500 dark:text-gray-400`) plus a
  margin, not a bespoke class string.
- **SVG paths:** reference constants in `internal/svg`, never inline a new path literal.
- **ProgressBar clamp:** use `max(0, min(100, v))` (Go 1.21+ builtins), not manual if-branch.
- **RTL logical properties:** use logical Tailwind utilities (`ms-`, `me-`, `ps-`,
  `pe-`, `start-0`, `end-0`, `text-start`, `border-s-`, `border-e-`) instead of
  physical ones (`ml-`, `mr-`, `pl-`, `pr-`, `left-`, `right-`, `text-left`,
  `border-l-`, `border-r-`). Logical properties render identically in LTR and
  mirror automatically in RTL (Arabic, Hebrew, Persian, Urdu). The ONLY exceptions
  are `left-1/2` (tooltip centering) and `left-0.5` (toggle thumb) — physical
  positioning that must not flip.
- **Container queries:** when a component should respond to its parent's width
  rather than the viewport, use Tailwind v4's `@container` + `@sm:`/`@md:`/`@lg:`
  variants. The `Grid` component supports this via the opt-in `ContainerResponsive`
  field.

## Anti-patterns to refuse on review

- `switch` statements for style lookups → replace with a map + `utils.Lookup`.
- Inline `<script>` without `nonce` → use `layout.Script` or add `nonce={ props.Nonce }`.
- `*_templ.go` deleted or untracked after a `.templ` edit → `git add -f` if new.
- Raw `class={ "a " + b }` string concatenation → `utils.Class("a", b)`.
- A new panic in render code → replace with a documented fallback value.
- A new dependency in `go.mod` outside the allowed three.
- `slate-*` dark-mode colors → switch to `gray-*`.
- Hardcoded `aria-label` that ignores `props.AriaLabel` → propagate via `utils.Ternary`.
- Physical Tailwind properties (`ml-`, `mr-`, `pl-`, `pr-`, `left-`, `right-`, `text-left`)
  → switch to logical (`ms-`, `me-`, `ps-`, `pe-`, `start-`, `end-`, `text-start`).
- Inline transition strings (`"transition-colors motion-reduce:..."`) → use the
  shared `transitionFast`/`transitionNormal`/`transitionColors`/`transitionTransform`
  constants from `display/shared.go`.
- Duplicating an icon path or a Tailwind class string that already has a shared constant.
- Shipping a component with only assertion tests → add golden, BDD, a11y, example lenses.

## Where to read deeper (progressive disclosure)

Load these only when the task needs them — do not read proactively.

- **`AGENTS.md`** (repo root) — the full catalogue of conventions, gotchas, and the
  one-commit release convention. This is the single richest source of context; consult it
  whenever a rule's _why_ is unclear or you hit something surprising.
- **`CONTRIBUTING.md`** — human contributor setup and commit-message format.
- **`docs/migration/play-cdn-to-tailwind-v4.md`** — 7-step migration from Play CDN to v4 CSS-first.
- **`docs/recipes/server-rendered-htmx-error-feedback.md`** — 3 HTMX error render modes.
- **`docs/tailwind-v4-adoption-guide.md`** — Tailwind v4 CSS-first setup, `@source` scanning,
  and `@theme` overrides.
- **`docs/adr/`** — decision records: two icon systems, shared SVG helpers, committing
  generated templ files, filled-vs-stroke icons, JS attachment patterns, feedback-type
  unification. Read the relevant ADR before changing the thing it decided.
- **`docs/javascript-guide.md`** — comprehensive JS patterns reference: the decision
  ladder (native HTML → HTMX → singleton-guard → Alpine → Datastar → islands), CSP
  compliance, templ's built-in JS features (`OnceHandle`, `JSFuncCall`, `JSONString`,
  `JSONScript`), TypeScript workflow, View Transitions API, and anti-patterns.
- **`docs/icons-only-adoption.md`** — adopting just the `icons` package (CSS-agnostic).
- **`docs/DOMAIN_LANGUAGE.md`** — ubiquitous-language glossary for terms used in types.
- **`internal/contract/component_props_test.go`** — the compile-time-enforced Props inventory;
  every new component must be registered here (see Component anatomy step 9).
- **`integration/composition_test.go`** — cross-package composition proof; extend it when a
  change affects how packages combine.
- **`examples/demo/`** — a live, runnable example wiring layout + feedback + display together.
  Read it as the canonical "how a consumer assembles a page" reference.
- **`FEATURES.md`** / **`TODO_LIST.md`** — honest feature inventory and short-term work. Check
  these before proposing a new component so you don't duplicate planned or existing work.
- **`scripts/release.sh`** — the one-command release cut; read it before tagging.

## Installing this skill into Crush

This skill lives in the repo so it versions with the code. Crush only auto-loads skills from
`~/.config/crush/skills/` (repo-level `crush.json` `context_paths` cover `AGENTS.md`, not
skills), so install it once with a symlink — it then stays in sync with the branch you're on:

```bash
mkdir -p ~/.config/crush/skills/templ-components
ln -sf "$PWD/skill/SKILL.md" ~/.config/crush/skills/templ-components/SKILL.md
```

Run this from the repo root so `$PWD` resolves correctly. `-f` makes it idempotent — re-run
after cloning or switching machines. Verify it appears in `available_skills` (named
`templ-components`); on a new session it will auto-load on the triggers in the frontmatter.

Prefer a copy? `cp skill/SKILL.md ~/.config/crush/skills/templ-components/SKILL.md` works too,
but you must re-copy after every update — the symlink is the recommended option.
