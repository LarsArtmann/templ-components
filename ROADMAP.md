# Roadmap — templ-components

This document tracks the long-term direction of the templ-components library.
Dates are indicative, not committed. Semantic versioning applies: anything
listed under v1.0 is a **freeze**, not a redesign.

---

## v0.x — Current (shipped)

The library is feature-complete for production server-rendered Go web apps.

| Pillar        | Status                                                                                                                                                                                                                                                                  |
| ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Components    | **88** across 9 packages (display, feedback, forms, layout, navigation, htmx, errorpage, icons, utils)                                                                                                                                                                  |
| Icons         | **101** named SVG icons (Heroicons v2 outline + Spinner), typed `icons.Name` constants                                                                                                                                                                                  |
| Typed enums   | 30+ closed-set enums, each with `IsValid()` + test coverage; `map[X]string` + `utils.Lookup` everywhere                                                                                                                                                                 |
| RTL / i18n    | All CSS uses logical properties (`ms-`/`me-`/`start-`/`end-`); auto-mirrors under `dir="rtl"`                                                                                                                                                                           |
| Dark mode     | Full `dark:` variant coverage on all components; enforced by `TestDarkModeCompliance` + `TestDarkModeSemanticColors` (failing CI tests). Class-based strategy via `layout.ThemeScript()` + `layout.ThemeToggle()`. `color-scheme: light/dark` for native form controls. |
| HTMX helpers  | `SwapOOB`, loading indicators, error handling, family-aware toasts                                                                                                                                                                                                      |
| CSP safety    | Every inline script carries `nonce={ props.Nonce }`; integration test guards regressions                                                                                                                                                                                |
| Accessibility | `motion-reduce:*` on all transitions/animations, `aria-sort`, focus trap, `aria-live` regions                                                                                                                                                                           |
| Error pages   | `errorpage` package: 404, full-page, inline detail, family-aware alert, `http.Handler` integration                                                                                                                                                                      |

**Current version:** see [`utils/version.go`](utils/version.go) and the latest
heading in [`CHANGELOG.md`](CHANGELOG.md).

---

## v1.0 — API Freeze

The goal of v1.0 is a **stable, frozen public API**. After v1.0 ships, breaking
changes require a v2.0 major bump.

| Workstream                                | Description                                                                                                                                           |
| ----------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Validate()` error on props               | Components whose props can be validated gain a `Validate() error` method. Render-time guards replace silent fallbacks / panics.                       |
| Move test helpers to `internal/testutil/` | `utils.Render`, `utils.AssertContainsAll`, golden-file helpers move to `internal/testutil/`. Reduces the public surface.                              |
| Self-host htmx by default                 | `PageProps` stops pulling htmx from CDN by default. CDN remains opt-in. See [ADR 0007](docs/adr/0007-self-host-htmx-default.md).                      |
| Remove deprecated aliases                 | `AlertType`/`ToastType` for `FeedbackType`, `ModalSizeFull`/`DrawerFull`, etc. are removed.                                                           |
| Semantic token layer                      | Opt-in `--color-tc-primary` indirection over raw Tailwind colors. See [ADR 0008](docs/adr/0008-semantic-tokens.md). Phased: document → opt-in → flip. |

---

## v2.0+ — Research (no timeline)

| Direction                    | Description                                                                        |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| Compound components          | `Trigger` / `Content` / `Close` sub-component pattern for Modal, Drawer, Dropdown. |
| Native `<dialog>`            | Modal/Drawer migrate to native `HTMLDialogElement`, reducing JS surface.           |
| Headless / unstyled variants | Components that emit semantics + ARIA without Tailwind classes.                    |
| CLI scaffolding tool         | `templ-components add <component>` — Go binary, no Node.                           |
| Demo / showcase site         | A hosted site rendering every component with live props.                           |

---

## Explicitly NOT Planned

| Rejected direction            | Why                                                         |
| ----------------------------- | ----------------------------------------------------------- |
| React / Vue / Svelte wrappers | The library is Go + templ + server-rendered HTML by design. |
| CSS-in-JS                     | Tailwind v4 utility classes are the styling standard.       |
| Node.js dependency            | Zero Node.js runtime requirement is a hard constraint.      |
