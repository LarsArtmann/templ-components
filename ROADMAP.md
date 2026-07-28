# Roadmap — templ-components

This document tracks the long-term direction of the templ-components library.
Dates are indicative, not committed. Semantic versioning applies: anything
listed under v1.0 is a **freeze**, not a redesign.

---

## v1.x — Current (shipped)

The library is feature-complete for production server-rendered Go web apps.

| Pillar            | Status                                                                                                                                                                                                                                                                                  |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Components        | **98** across 9 packages (display, feedback, forms, layout, navigation, htmx, errorpage, icons, utils) + 3 recipe screens                                                                                                                                                               |
| Icons             | **102** named SVG icons (Heroicons v2 outline + Spinner), typed `icons.Name` constants                                                                                                                                                                                                  |
| Typed enums       | 43 closed-set enums, each with `IsValid()` + test coverage; `map[X]string` + `utils.Lookup` everywhere                                                                                                                                                                                  |
| Layout            | Grid-first 2D layout primitives: `AppShell`, `Container`, `Split`, `Stack` + multi-column `Footer` + `Form.Layout` enum. Rule: grid = 2D, flex = 1D (ADR-0016). `minmax(0,1fr)` mandatory on all flexible grid columns.                                                                 |
| RTL / i18n        | All CSS uses logical properties (`ms-`/`me-`/`start-`/`end-`); auto-mirrors under `dir="rtl"`. Enforced by `TestRTLLogicalProperties` scanner.                                                                                                                                          |
| Dark mode         | Full `dark:` variant coverage on all components; enforced by `TestDarkModeCompliance` + `TestDarkModeSemanticColors` (failing CI tests). Class-based strategy via `layout.ThemeScript()` + `layout.ThemeToggle()`. `color-scheme: light/dark` for native form controls.                 |
| HTMX helpers      | `SwapOOB`, loading indicators, error handling, family-aware toasts                                                                                                                                                                                                                      |
| CSP safety        | Every inline script carries `nonce={ props.Nonce }`; integration test guards regressions                                                                                                                                                                                                |
| Accessibility     | `motion-reduce:*` on all transitions/animations, `aria-sort`, focus trap, `aria-live` regions                                                                                                                                                                                           |
| Error pages       | `errorpage` package: 404, full-page, inline detail, family-aware alert, `http.Handler` integration                                                                                                                                                                                      |
| Theming           | Semantic token layer (`templ-components-theme.css`) + 3 presets, opt-in (v0.22.0). Self-host HTMX opt-in via `PageProps.HTMXSrc`.                                                                                                                                                       |
| Container queries | 8 opt-in container-aware components (`Grid.ContainerResponsive`, `Card`/`Nav`/`Split`/`Form`/`Pagination`/`DefinitionGrid`/`SkeletonCardGrid` `.ContainerAware`) adapt to parent width via `@container` instead of the viewport. ADR-0018. Default flip to opt-out is a v2.0 candidate. |
| Testing & QA      | Golden-file HTML snapshots (`internal/golden` in 5 packages), drift-guard tests (component/enum/version counts), **pixel-level visual regression** (`visualtest/` — chromedp + pixelmatch, `nix run .#visual`; 31 goldens across 11 component types incl. RTL), and compliance scanners (dark mode, motion-reduce, RTL logical properties, container-query, generated-file sync). |

**Current version:** see [`utils/version.go`](utils/version.go) and the latest
heading in [`CHANGELOG.md`](CHANGELOG.md).

---

## v1.0 — API Freeze (SHIPPED 2026-07-21)

v1.0.0 shipped with `ErrorPageProps.Validate()`, deprecated alias removal,
and CI docs-health drift guard. See `CHANGELOG.md` for the full entry.

| Workstream                                | Status      | Notes                                                                                            |
| ----------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------ |
| `Validate()` on `ErrorPageProps`          | ✅ DONE     | v1.0.0 — other props use graceful `utils.Lookup` fallback (no `Validate` needed).                |
| Move test helpers to `internal/testutil/` | ⬜ DEFERRED | 70+ test imports affected; large mechanical migration, deferred post-v1.0.                       |
| Self-host htmx opt-in                     | ✅ DONE     | v0.22.0 — `PageProps.HTMXSrc` opt-in; CDN remains default. Auto-suppresses response-targets ext. |
| Remove deprecated aliases                 | ✅ DONE     | v1.0.0 — `ModalSizeFull`, `DrawerFull`, `FamilyFromErrorFamily`, `FormProps.Inline` removed.     |
| Semantic token layer                      | ✅ DONE     | v0.22.0 — `templates/templ-components-theme.css` + 3 presets, opt-in. See ADR-0008.              |

---

## v1.1+ — Shipped platform work

| Workstream                   | Status     | Version        | Notes                                                                                                                                                                                                                                                                                  |
| ---------------------------- | ---------- | -------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Popover API migration        | ✅ DONE    | v0.20.0        | Dropdown/Popover/ContextMenu on native `popover="auto`. See ADR-0017.                                                                                                                                                                                                                  |
| Container-aware components   | ✅ DONE    | v0.21.0        | `NavProps.ContainerAware`, `CardProps.ContainerAware`. See ADR-0018.                                                                                                                                                                                                                   |
| Recipes package              | ✅ DONE    | v0.21.0        | `recipes.Dashboard/SettingsLayout/LoginCard`. See ADR-0019.                                                                                                                                                                                                                            |
| `tc` CLI scaffolding tool    | ✅ DONE    | v1.1.0         | `tc init/ls/add` with embedded sources. See `docs/cli.md`.                                                                                                                                                                                                                             |
| Headless / unstyled variants | ❌ WONTFIX | v1.1.0         | ADR-0021 evaluated 3 options; existing `Class` override accepted.                                                                                                                                                                                                                      |
| Post-v1.1.0 defect fixes     | ✅ DONE    | v1.2.0         | Popover top-layer positioning (D1), Tooltip aria-describedby (D3), HTMXSrc CDN leak (D4), `tc add` dependency warning (D6).                                                                                                                                                            |
| `navigation.SidebarNav`      | ✅ DONE    | v1.2.0         | Vertical sidebar nav for admin panels. Permanently-dark surface.                                                                                                                                                                                                                       |
| Recipe demo routes           | ✅ DONE    | v1.2.0         | `/recipes/{dashboard,settings,login}` in `examples/demo`.                                                                                                                                                                                                                              |
| Container query expansion    | ✅ DONE    | `[Unreleased]` | Extended `ContainerAware` from 3 → 8 components (added Split, Form, Pagination, DefinitionGrid, SkeletonCardGrid). ADR-0018.                                                                                                                                                           |
| Visual regression framework  | ✅ DONE    | `[Unreleased]` | `visualtest/` module — chromedp + pixelmatch, separate module (no consumer dep pollution), `nix run .#visual`, CI job. **31 goldens** across 11 component types (Button, Alert, Badge, Card, Modal, Drawer, Dropdown, Popover, ContextMenu, Input, Select + RTL). Shared Chromium process (one browser, per-test tabs). Interaction states: `StateHover`, `StateFocus`, `StateClick`, `StateContext`. See `docs/visual-testing.md`. |
| CSS source scanning fix      | ✅ DONE    | `[Unreleased]` | `@source "**/*.go"` added to CSS templates — Tailwind v4 now scans Go files for class lookup maps. Previously, errorpage family classes (amber/orange/purple) were silently missing from compiled CSS. Enforced by `TestTailwindGoSourceScanning`.                                     |
| Generated-file sync guard    | ✅ DONE    | `[Unreleased]` | `TestTemplGeneratedInSync` verifies every `*_templ.go` file's imports match its `.templ` source. Prevents stale generated artifacts (breadcrumbs drift: json v2 vs v1).                                                                                                                |
| Container-query compliance   | ✅ DONE    | `[Unreleased]` | `TestContainerQueryCompliance` scans `.templ` for structural viewport breakpoints without `ContainerAware`. Mirrors dark-mode/motion-reduce compliance scanners.                                                                                                                       |
| Lint-config regression guard | ✅ DONE    | `[Unreleased]` | 3-layer prevention for the recurring `.golangci.yml` linter regression (5 occurrences): `TestGolangciDisabledLinters` (CI), `scripts/check-lint-config.sh` (pre-commit), CI lint-config guard step. Root cause documented in AGENTS.md.                                                |
| GOWORK=off + .envrc          | ✅ DONE    | `[Unreleased]` | `GOWORK=off` in devShell `shellHook` + `.envrc` (direnv) sets `GOEXPERIMENT=jsonv2` repo-wide for all tools (go, gopls, BuildFlow, IDE).                                                                                                                                               |

---

## v2.0+ — Research (no timeline)

| Direction                       | Description                                                                                                                                                                                                                                                                           |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Compound components             | `Trigger` / `Content` / `Close` sub-component pattern for Modal, Drawer, Dropdown. **ADR-0023** written (compound overlay API design).                                                                                                                                                |
| Per-package modules split       | Independently importable packages. ADR-0020 written; deferred until consumer demand.                                                                                                                                                                                                  |
| Default flip                    | Self-host HTMX becomes default (CDN opt-in); semantic tokens become default; `ContainerAware` flips to opt-out for components commonly placed in constrained containers (Grid, Card, Split). All shipped opt-in. **ADR-0022** written (v2 default-flip migration plan). See TODO #35. |
| Demo / showcase site            | A hosted site rendering every component with live props.                                                                                                                                                                                                                              |
| Visual test coverage expansion  | `StateClick`/`StateContext` + 4 overlay goldens (Dropdown/Popover/ContextMenu) shipped — 31 goldens / 11 component types. Remaining: expand to Combobox, Tabs, Table, Accordion, Tooltip, Carousel, Spinner, Skeleton, ProgressBar, Toast, Avatar, ErrorPage, NotFound404. Calibrate `MaxMismatch` empirically. See TODO #79, #82. |
| Container-aware expansion       | Candidates for `.ContainerAware`: `Container`, `Breadcrumbs`, `EmptyState`, `NotFound404`, `Footer`. Plus a shared `containerAwareWrapper` sub-template to remove the 8× hand-written boilerplate. Fluid typography via container query units (`cqi`/`cqw`) also unexplored. |
| Visualtest API improvements     | `Options` struct uses `*bool` for optional booleans (current `bool` zero-value conflates "false" with "unset"). `InteractionState.String()` for error messages. Viewport presets (iPhone, iPad, desktop). `MaxMismatch` calibration tool. See TODO #84.                                |

---

## Explicitly NOT Planned

| Rejected direction            | Why                                                          |
| ----------------------------- | ------------------------------------------------------------ |
| React / Vue / Svelte wrappers | The library is Go + templ + server-rendered HTML by design.  |
| CSS-in-JS                     | Tailwind v4 utility classes are the styling standard.        |
| Node.js dependency            | Zero Node.js runtime requirement is a hard constraint.       |
| Headless / unstyled variants  | ADR-0021: existing `Class` override suffices. Closed v1.1.0. |
