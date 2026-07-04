---
name: templ-components
description: Authoring playbook for the templ-components Go UI library (templ + HTMX + Tailwind v4). Use this skill whenever working in the github.com/larsartmann/templ-components repo, or when adding, editing, reviewing, or debugging any templ component, props struct, typed enum, style lookup, icon path, or HTMX helper here. Also trigger when the user asks to "add a component", "fix a component", extend the library, wire up a new icon, add a typed enum, set up Tailwind v4 scanning for templ, integrate HTMX loading/error handling, or questions how a component should be structured, themed, made accessible, CSP-safe, or dark-mode aware in this codebase. Load this skill BEFORE writing or editing any .templ file in this repo.
metadata:
  tags: templ, templ-components, htmx, tailwind, tailwind-v4, go, ui, components, accessibility, csp, dark-mode, server-rendered
---

# templ-components

Authoring playbook for a server-rendered, type-safe, accessible, CSP-ready Go component
library built on [templ](https://templ.guide), [HTMX](https://htmx.org), and Tailwind v4.

This skill tells you **how to build and modify components correctly** in this repo. It is the
procedural counterpart to `AGENTS.md` (which holds broad project context). When a question is
about _what exists_ or _why something is the way it is_, read `AGENTS.md` and the linked docs.
When it is about _how to make a new or changed component fit the library_, follow this skill.

## Principles (the why behind every rule)

- **Make invalid states unrepresentable.** Variants are typed enums, not strings; lookups are
  maps with explicit fallbacks; the only permitted runtime panic is the single developer
  data-integrity check in `icons` (stray `|` in a path). Everything else degrades gracefully.
  This is a _library_: a panic in a consumer's render path is a bug we shipped.
- **Server-rendered first, zero JS by default.** Prefer native HTML (`<details>`, forms, links)
  over scripts. When interactivity is unavoidable (modal, drawer, dropdown, accordion, combobox,
  tabs, theme toggle), ship minimal vanilla JS guarded by a global singleton flag so it is
  idempotent across HTMX re-renders.
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
   Watch for the BuildFlow pre-commit gotcha documented in `AGENTS.md` that re-appends
   `*_templ.go` to `.gitignore`; for already-tracked files it is harmless, but a NEW
   component's generated file will be invisible to `git status` until `git add -f`.
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
  component_templ.go     # GENERATED — commit it
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
   (e.g. `overlayShell`, `dialogHeader`, `navLinkAnchor`, the six `errorpage/shared.templ`
   helpers). Lift duplication into a sub-template rather than copying.
9. **Register the new Props type in the contract inventory.** Open
   `internal/contract/component_props_test.go` and add `yourpackage.YourProps{}` to the
   `componentTypes()` slice (in the right package section). The test
   `TestAllComponentPropsSatisfyInterface` then enforces at CI time that your struct both
   embeds `utils.BaseProps` and satisfies `utils.ComponentProps` — catching silent contract
   breakage for consumers using generic wrappers. Forgetting this step is the #1 way a new
   component slips in without the BaseProps embed.

## Decision trees

### Map lookup vs if-branch for a variant

- **Pure class/style data** (colors, sizes, padding) → map + `utils.Lookup(map, key, fallback)`.
  This is the dominant pattern: `badgeStyleMap`, `cardPaddingLookup`, `spinnerSizeLookup`, etc.
- **Structural DOM differences** (TabsVariant, DropdownPosition, TrendDirection) → `if`-branch
  inside the template. Maps are for data, not for choosing which markup to emit.
- **Enum validation** → 14 enums use map+fallback (graceful); `InputType`,
  `ButtonHTMLType`, `FormMethod` fall back to HTML-spec defaults; only the icons path-data
  integrity check is allowed to panic.

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
- **Class merging:** always go through `utils.Class(...)` so tailwind-merge resolves conflicts
  and consumer overrides win. The only exception is `templ.KV` conditionals where the templ
  runtime must comma-join. `utils.Class` is mutex-protected; do not bypass it.
- **CSP nonce:** every inline script takes `nonce={ props.Nonce }`.
- **Card shell:** use the shared `cardShellClass` constant; `SimpleCard` composes through
  `Card` internally.
- **Muted text:** use `mutedTextClass` (`text-sm text-gray-500 dark:text-gray-400`) plus a
  margin, not a bespoke class string.
- **SVG paths:** reference constants in `internal/svg`, never inline a new path literal.

## Anti-patterns to refuse on review

- `switch` statements for style lookups → replace with a map + `utils.Lookup`.
- Inline `<script>` without `nonce`.
- `*_templ.go` deleted or untracked after a `.templ` edit.
- Raw `class={ "a " + b }` string concatenation → `utils.Class("a", b)`.
- A new panic in render code → replace with a documented fallback value.
- A new dependency in `go.mod` outside the allowed three.
- `slate-*` dark-mode colors → switch to `gray-*`.
- Hardcoded `aria-label` that ignores `props.AriaLabel` → propagate via `utils.Ternary`.
- Duplicating an icon path or a Tailwind class string that already has a shared constant.

## Testing expectations per component

Each component package carries several complementary test lenses. When you add or change a
component, extend or add the relevant ones:

| Test file            | What it guards                                                        |
| -------------------- | --------------------------------------------------------------------- |
| `golden_test.go`     | Exact rendered HTML matches `.golden` snapshot (regression baseline)  |
| `a11y_test.go`       | ARIA, roles, keyboard, motion-reduce, screen-reader text are present  |
| `bdd_test.go`        | Behaviour spec via Ginkgo/GOmega (user-visible behaviour, not markup) |
| `edge_cases_test.go` | Empty inputs, unknown enum values, ID collisions, concurrency         |
| `example_test.go`    | The godoc examples compile and render (docs stay honest)              |
| `snapshot_test.go`   | Broader composition snapshots                                         |
| `coverage_*_test.go` | Targeted coverage of private helpers and branches                     |

Golden tests use `internal/golden.Assert(t, name, got)` with CSS-class normalization; pass
`-update` to regenerate after an intentional visual change, then review the diff.

## Where to read deeper (progressive disclosure)

Load these only when the task needs them — do not read proactively.

- **`AGENTS.md`** (repo root) — the full catalogue of conventions, gotchas, and the
  one-commit release convention. This is the single richest source of context; consult it
  whenever a rule's _why_ is unclear or you hit something surprising.
- **`CONTRIBUTING.md`** — human contributor setup and commit-message format.
- **`README.md`** — the consumer-facing component catalogue and Tailwind/theming setup.
- **`docs/tailwind-v4-adoption-guide.md`** — Tailwind v4 CSS-first setup, `@source` scanning,
  and `@theme` overrides.
- **`docs/adr/`** — decision records: two icon systems, shared SVG helpers, committing
  generated templ files, filled-vs-stroke icons, JS attachment patterns, feedback-type
  unification. Read the relevant ADR before changing the thing it decided.
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

This skill lives in the repo so it versions with the code. To make Crush auto-trigger it from
any session, symlink (or copy) it into the user skills directory:

```bash
mkdir -p ~/.config/crush/skills/templ-components
ln -s "$PWD/skill/SKILL.md" ~/.config/crush/skills/templ-components/SKILL.md
```

Then it appears in `available_skills` and loads on the triggers in the frontmatter.
