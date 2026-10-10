# templ-components vs the ecosystem

**External facts verified: 2026-10-08** — star counts fetched from the GitHub
API on that date; product claims read from the linked sources on that date.
Internal claims (component/test/golden counts) are drift-guarded by
`utils.TestDocsCountDrift`. Re-verify the external numbers before citing or
extending this document — they rot fastest.

**Re-verify cadence:** re-check the external facts (stars, component lists,
naming, distribution model) at least QUARTERLY, and immediately after any
known competitor event (a rename/rework like templUI→shadcn-templ, a major
release, or a star-order-of-magnitude move). Record the new verify date in
this header and in the TODO_LIST sweep row; competitor claims shipped stale
twice before this cadence existed (templUI/Alpine.js falsehood in README +
website, found and fixed 2026-10-08).

---

## The contenders

| Project                                                                       | What it is                                                                                                                            | Stars* |
| ----------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| **templ-components** (this repo)                                              | Server-rendered Go component library: templ + Tailwind v4, HTMX-native, Datastar opt-in. Versioned Go module.                         | 4      |
| [shadcn-templ](https://github.com/axadrn/shadcn-templ) (formerly **templUI**) | The shadcn/ui port for templ: same DOM (`data-slot`), same design language, 8 style themes, copy-paste registry + installable blocks. | 1,754  |
| [goshipit](https://github.com/haatos/goshipit) (GoShip.it)                    | DaisyUI-on-Tailwind component library for the GOTH stack, HTMX-driven, with a `gsi` copy CLI and a demo app.                          | 277    |
| [templ examples (a-h)](https://github.com/a-h/templ/tree/main/examples)       | Official educational examples in the templ repo — not a library.                                                                      | —      |

\* GitHub API, 2026-10-08.

Note on naming: templUI v1 (templui.io, 42 components) was renamed and
reworked into shadcn-templ 2.0; templui.io now reads "templUI is now
shadcn-templ". Older third-party comparisons that describe templUI as
"Alpine.js-based" are describing the old v1 lineage — the current project uses
vanilla JS only.

## Snapshot

|                  | templ-components                                                                                                 | shadcn-templ                                                                                                          | goshipit                         |
| ---------------- | ---------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- | -------------------------------- |
| License          | MIT                                                                                                              | MIT                                                                                                                   | MIT                              |
| templ version    | v0.3.1020 (pinned, zero-diff generator)                                                                          | v0.3.1070                                                                                                             | not checked                      |
| Go floor         | 1.27                                                                                                             | 1.26                                                                                                                  | 1.x                              |
| Components       | 128 primitives + 4 recipe screens                                                                                | ~64 + 27 installable blocks                                                                                           | 52                               |
| Icons            | 102 typed names, animated variants                                                                               | icon component (lucide set)                                                                                           | via DaisyUI/inline               |
| CSS              | Tailwind v4 (CSS-first), one look                                                                                | Tailwind v4 + CSS vars, **8 style themes**                                                                            | Tailwind v4 + **DaisyUI 5**      |
| JavaScript       | Per-component CSP-safe inline singletons; native APIs first (`<dialog>`, Popover API, `<details>`, CSS tooltips) | One esbuild-bundled, minified runtime porting Base UI behavior (Floating UI, portals, focus, scroll lock)             | HTMX-driven; no custom framework |
| CSP              | Strict nonce discipline incl. omit-empty rule, render-table guard                                                | Nonce via `templ.GetNonce` on the runtime bundle                                                                      | not documented                   |
| HTMX             | Built-in package; self-hosted runtime embedded by default                                                        | Compatible (htmx fixture-tested)                                                                                      | Native (core design)             |
| Datastar         | Opt-in package + `wire.Action` dual transport                                                                    | none                                                                                                                  | none                             |
| Dark mode        | Built-in, scanner-enforced (`TestDarkModeCompliance`)                                                            | Class strategy + CSS vars                                                                                             | via DaisyUI themes               |
| Distribution     | `go get`, semver, 7-module opt-in DAG                                                                            | copy-paste registry (`shadcn-templ add`), you own the code                                                            | copy CLI (`gsi`)                 |
| Requires Node.js | No                                                                                                               | No                                                                                                                    | CSS build only                   |
| Typed props      | 67 string enums (66 with `IsValid()`, compile-time-checked lookup maps                                           | Base UI-named props, value-only                                                                                       | DaisyUI-shaped props             |
| Testing          | 299 HTML goldens, 200 pixel goldens, axe gate (default-fail), RTL/motion/dark-mode scanners, fuzz, vnu HTML gate | 52 Go test files + Playwright parity harness vs pinned shadcn/ui (DOM/focus/scroll-lock/pixel) + a11y/behavior suites | demo-site level                  |
| Charts           | Native SVG (Line/Area/Pie/Sparkline/Bar/Heatmap) + opt-in ECharts adapter                                        | client-rendered (Recharts-compatible)                                                                                 | —                                |
| Error pages      | `errorpage` package + go-error-family integration                                                                | —                                                                                                                     | —                                |

## Philosophy: who owns the component code?

This is the deepest difference, and it is a genuine tradeoff — neither model is
strictly better.

**Copy-paste registry (shadcn-templ, goshipit).** `shadcn-templ add button`
copies the `.templ`/`.js` source into your repo. You own it: bend any class,
rename any prop, delete what you don't need. There is no version to pin and no
upgrade path — fixes land only when you re-add and re-diff by hand. The
library's own tagline: "Use this to build your own component library."

**Versioned module (templ-components).** `go get` a release; `go get -u` takes
fixes. Semver, CHANGELOG, migration docs, per-module opt-in so an unused
package never enters your dependency graph. The trade: consumers extend via
props, `Attrs`, slots, and Tailwind overrides — you cannot rewrite a
component's internals without forking.

## Where they are ahead

**shadcn-templ**

- **Brand and adoption.** 1,754 stars vs 4. The shadcn/ui design language is
  the most-recognized web UI aesthetic, and their 1:1 DOM fidelity means
  shadcn tutorials, blocks, and themes transfer directly.
- **Design breadth.** 8 complete style themes (luma, lyra, maia, mira, nova,
  rhea, sera, vega). We ship one look (consumer-overridable via `@theme`, but
  one look).
- **Installable blocks.** dashboard01, 5 login pages, 5 signup pages, 16
  sidebar layouts — full screens as one command. Our `recipes` package
  (Dashboard, SettingsLayout, LoginCard, AuthLayout) ships in the module but
  is not a scaffold.
- **Client-rich components.** Command palette (cmdk), Input OTP, Resizable
  panels, Menubar, Navigation Menu, ScrollArea, Toggle Group, chat primitives.
  Several of these are deliberately demand-gated here (TODO_LIST #217); their
  adoption is fresh demand evidence for that gate.
- **Upstream parity harness.** They diff DOM tree, focus, scroll lock, and
  pixels against a pinned shadcn/ui reference build. We have no external
  ground truth of that kind — our goldens pin our own output.

**goshipit**

- **DaisyUI ecosystem.** Themed components and semantic class names out of the
  box; if you already run DaisyUI, it is the native choice.
- **GOTH app breadth.** Ships pricing, timeline, chat, OTP, hero sections —
  plus a runnable demo app and code/data generator around the components.

## Where templ-components is ahead

- **Correctness machinery per component.** Invariant scanners run on every
  source file: dark-mode gaps, missing `motion-reduce:`/`motion-reduce:duration-0`,
  physical RTL properties (logical-properties-only rule), coarse-pointer
  fallbacks, CSP nonce regressions, CSS freshness, docs-count drift — most as
  default-fail CI gates, plus an axe-core sweep over the live demo with a
  documented debt ledger. Neither competitor documents anything comparable.
- **Genuinely zero-JS by default.** Modals/drawers are native `<dialog>`,
  overlays use the native Popover API, accordions are `<details>`, tooltips
  are pure CSS. Less script surface to trust, less to CSP-audit.
- **Dual transport with one spec.** `wire.Action` renders htmx (default) or
  Datastar attributes from one typed struct; `wire.Handler` serves both
  dialects from one endpoint with response-header targeting. Kanban drag-and-drop
  and infinite scroll work under both runtimes. shadcn-templ is htmx-compatible;
  goshipit is htmx-native; nobody else has a Datastar story.
- **Server-facing features.** `errorpage` (family-aware pages + handler
  integration with go-error-family), `FilterInput` debounce, validation
  summary, dirty-guard — infrastructure components, not just widgets.
- **Dependency budget.** 3 runtime deps (templ, tailwind-merge-go,
  go-error-family). shadcn-templ's repo carries esbuild, goldmark, freetype
  (site tooling); goshipit pulls the DaisyUI/Node toolchain for CSS.
- **Release discipline.** One-commit releases, signed tags for root + all
  sub-modules, CHANGELOG enforced warm, generator pinned for zero-diff
  regeneration.

## Which one should you pick?

- You want the shadcn/ui look, themes, and ready-made screens, and you are
  happy to own (and maintain) the component source: **shadcn-templ**.
- You are all-in on DaisyUI aesthetics and HTMX, and Node-for-CSS is fine:
  **goshipit**.
- You want a versioned dependency with strict CSP, accessibility scanners,
  RTL support, HTMX **and** Datastar under one API, error pages, and charts:
  **templ-components**.

They are not mutually exclusive: our components coexist with any of these in
the same app (both are templ + Tailwind v4 at heart).

## Sources (read 2026-10-08)

- shadcn-templ: github.com/axadrn/shadcn-templ (repo tree, `go.mod`,
  `registry.json`, `parity/README.md`, `plans/htmx-616.md`, `AGENTS.md`,
  `llms.txt`), GitHub API star/fork counts.
- templUI → shadcn-templ rename: templui.io banner + shadcn-templ.com docs.
- goshipit: github.com/haatos/goshipit (README, component tree, package.json),
  GitHub API counts.
- This repo: README.md / FEATURES.md numbers are guard-verified
  (`utils.TestDocsCountDrift`), not hand-copied.
