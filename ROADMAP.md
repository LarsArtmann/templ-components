# Roadmap — templ-components

This document tracks the long-term direction of the templ-components library.
Dates are indicative, not committed. Semantic versioning applies: anything
listed under v1.0 is a **freeze**, not a redesign.

---

## v0.x — Current (shipped)

The library is feature-complete for production server-rendered Go web apps.

| Pillar        | Status                                                                                                                                                                                                                                                                  |
| ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Components    | **98** across 9 packages (display, feedback, forms, layout, navigation, htmx, errorpage, icons, utils)                                                                                                                                                                  |
| Icons         | **102** named SVG icons (Heroicons v2 outline + Spinner), typed `icons.Name` constants                                                                                                                                                                                  |
| Typed enums   | 43 closed-set enums, each with `IsValid()` + test coverage; `map[X]string` + `utils.Lookup` everywhere                                                                                                                                                                  |
| Layout        | Grid-first 2D layout primitives: `AppShell`, `Container`, `Split`, `Stack` + multi-column `Footer` + `Form.Layout` enum. Rule: grid = 2D, flex = 1D (ADR-0016). `minmax(0,1fr)` mandatory on all flexible grid columns.                                                 |
| RTL / i18n    | All CSS uses logical properties (`ms-`/`me-`/`start-`/`end-`); auto-mirrors under `dir="rtl"`. Enforced by `TestRTLLogicalProperties` scanner.                                                                                                                          |
| Dark mode     | Full `dark:` variant coverage on all components; enforced by `TestDarkModeCompliance` + `TestDarkModeSemanticColors` (failing CI tests). Class-based strategy via `layout.ThemeScript()` + `layout.ThemeToggle()`. `color-scheme: light/dark` for native form controls. |
| HTMX helpers  | `SwapOOB`, loading indicators, error handling, family-aware toasts                                                                                                                                                                                                      |
| CSP safety    | Every inline script carries `nonce={ props.Nonce }`; integration test guards regressions                                                                                                                                                                                |
| Accessibility | `motion-reduce:*` on all transitions/animations, `aria-sort`, focus trap, `aria-live` regions                                                                                                                                                                           |
| Error pages   | `errorpage` package: 404, full-page, inline detail, family-aware alert, `http.Handler` integration                                                                                                                                                                      |

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

| Workstream                   | Status     | Version | Notes                                                                 |
| ---------------------------- | ---------- | ------- | --------------------------------------------------------------------- |
| Popover API migration        | ✅ DONE    | v0.20.0 | Dropdown/Popover/ContextMenu on native `popover="auto`. See ADR-0017. |
| Container-aware components   | ✅ DONE    | v0.21.0 | `NavProps.ContainerAware`, `CardProps.ContainerAware`. See ADR-0018.  |
| Recipes package              | ✅ DONE    | v0.21.0 | `recipes.Dashboard/SettingsLayout/LoginCard`. See ADR-0019.           |
| `tc` CLI scaffolding tool    | ✅ DONE    | v1.1.0  | `tc init/ls/add` with embedded sources. See `docs/cli.md`.            |
| Headless / unstyled variants | ❌ WONTFIX | v1.1.0  | ADR-0021 evaluated 3 options; existing `Class` override accepted.     |

---

## v2.0+ — Research (no timeline)

| Direction                 | Description                                                                          |
| ------------------------- | ------------------------------------------------------------------------------------ |
| Compound components       | `Trigger` / `Content` / `Close` sub-component pattern for Modal, Drawer, Dropdown.   |
| Per-package modules split | Independently importable packages. ADR-0020 written; deferred until consumer demand. |
| Demo / showcase site      | A hosted site rendering every component with live props.                             |

---

## Explicitly NOT Planned

| Rejected direction            | Why                                                          |
| ----------------------------- | ------------------------------------------------------------ |
| React / Vue / Svelte wrappers | The library is Go + templ + server-rendered HTML by design.  |
| CSS-in-JS                     | Tailwind v4 utility classes are the styling standard.        |
| Node.js dependency            | Zero Node.js runtime requirement is a hard constraint.       |
| Headless / unstyled variants  | ADR-0021: existing `Class` override suffices. Closed v1.1.0. |
