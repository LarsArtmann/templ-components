# Demo Coverage Map — components vs demo routes (2026-10-15)

Authoritative answer to "which library components does the demo actually exercise?"
Built from `examples/demo/pages.go` (route source of truth) cross-referenced against the
136 exported templ components (the same counting method as `utils.TestDocsCountDrift`:
`^templ [A-Z]` over display, feedback, forms, navigation, errorpage, layout, htmx,
datastar, charts/echarts).

**Why this map exists:** progressive-enhancement and cross-engine claims are only as
strong as the demo surface they are proven on. The visual suite, axe sweep, route
goldens, and any future cross-engine captures all inherit THIS list — a component absent
from the demo is invisible to every browser-level gate. Companion to
`docs/testing/a11y-gate-policy.md`.

## Route inventory (17 routes)

`/` (home) · `/layout` · `/display` · `/feedback` · `/forms` · `/stylable-select` ·
`/navigation` · `/icons` · `/htmx` · `/datastar` · `/wire` · `/kanban` · `/echarts` ·
`/recipes` · `/users` · `/error-pages`

Plus standalone shells outside the AppShell: recipe detail routes and error routes
(own `layout.Base` shells, per AGENTS.md).

## Result: 136 = 124 demoed + 6 composition-covered + 6 gaps

### Composition-covered (rendered via a parent component — waivable)

| Component | Rendered by | Demo evidence |
| --- | --- | --- |
| `forms.Label` | every field wrapper (Input, Combobox, …) | `/forms` renders dozens of labelled fields |
| `forms.FormFieldWrapper` | Input/Select/Textarea/… internals | same as above |
| `forms.Radio` | `RadioGroup` internals | RadioGroup demoed on `/forms` |
| `layout.Stylesheet` | `Base`/`DefaultPageProps` head injection | every demo page |
| `layout.ThemeScript` | `Base`/`DefaultPageProps` head injection | every demo page |
| `navigation.MobileNavLink` | `MobileMenu` internals | mobile menu on every AppShell route (e2e-covered by `demo_mobile_e2e_test.go`) |

### True gaps (never rendered anywhere in `examples/demo/`)

| Component | Keep/waive | Reason + suggested remedy |
| --- | --- | --- |
| `datastar.LoadingButton` | **KEEP** | The Datastar twin of a heavily-demoed htmx component is invisible; the twin-parity claim is unproven in a browser. Add to `/datastar` with an `Indicator`-driven action. |
| `datastar.PolledRegion` | **KEEP** | Same twin-parity gap (htmx.PolledRegion IS demoed). Add to `/datastar` with a short interval against a stub endpoint. |
| `display.StatusBadgeWith` | **KEEP** | Public API with ZERO usage anywhere in the repo outside its own tests — an API-surface orphan. Either demo it on `/display` (preferred, one-liner) or deprecate it. |
| `layout.MetaRefresh` | waive | Head-level redirect utility with no visual surface; golden + unit tests cover it. Browser-level coverage adds nothing observable. |
| `layout.Minimal` | waive | Infrastructure shell variant (no chrome); exercised by design in unit tests. Not a UI component. |
| `layout.Script` | waive | External-script helper; demo pages use the self-hosted inline HTMX embed and `SDKScript` instead. Golden-covered; rendering an external CDN script in the demo would add a network dependency to every capture. |

### Property-level gap (component demoed, a variant never exercised)

| Component | Missing state | Remedy |
| --- | --- | --- |
| `forms.Select` | `Stylable: true` (`appearance: base-select` path) | **Fixed 2026-10-15**: new `/stylable-select` demo route renders the Stylable path (flat, groups, error) next to its native counterpart with fallback documentation. This is the riskiest cross-engine property in the library — `base-select` is unsupported in Firefox (Mozilla bugs 1944403/1958443) — and was invisible to every browser-level gate before this fix. |

## Maintenance

This map is a point-in-time audit (2026-10-15). When adding a component: the
goldens-cover-this / wired⇒e2e checklist (`docs/plan-authoring-checklist.md`) already
forces demo presence for wired components; use this file's gap table as the model for the
keep/waive verdict. Re-audit triggers: any new root-package component, or any new
cross-engine-sensitive CSS property adopted in `templates/custom.css`.
