# 50 Big-Picture Improvement Ideas — templ-components

**Date:** 2026-10-14 | **Mode:** brainstorm (point-in-time review artifact)

> A fresh ideation pass, deliberately UNRELATED ideas: 50 directions spanning 12
> axes (API, platform, components, docs, testing, distribution, i18n, perf, CI,
> a11y, charts, strategy). Dupe-checked against `docs/100-IMPROVEMENT-IDEAS.md`,
> `docs/STANDOUT-IDEAS.md`, `docs/SUPERB-FOR-PERSONAL-USE.md`, and the ROADMAP
> harvest tables — items that converge with an existing tracked direction say so
> explicitly and frame only the delta. Nothing here is committed work; harvest
> into `TODO_LIST.md`/ROADMAP only after a demand or owner decision. Several
> component/API ideas are demand-gated per the repo's standing policy (#217 gate).

---

## A. Go API & type system

### 1. Context-injected defaults (nonce, locale, theme)

The #1 adoption papercut is threading `BaseProps.Nonce` (and soon locale/labels)
through every props struct. Add `utils.ContextWithNonce(ctx, nonce)` (plus
locale later): components read the context value as a FALLBACK after the
explicit props field. Explicit wins, context fills, zero value still means the
documented default — the struct-literal zero-value rule holds.
**First step:** ADR on precedence (props > context > zero), then nonce-only
pilot on one script-emitting component + a compliance test that context and
props never silently disagree.

### 2. Typed `Href`/`URL` type

Raw-string `href`s are the last stringly-typed boundary: open-redirect risks,
lost `target=_blank` discipline, inconsistent fragment handling. A
`utils.Href` (validated constructor: internal path / anchored fragment /
external-with-rel) consumed by Button, Badge, ExternalLink, NavLink,
CommandItem makes the unsafe shape unrepresentable. Zero value = no link.
**First step:** survey every `Href string` field, design the type with
`ParseHref`, migrate ONE package behind a type alias for source compat.

### 3. URL round-trip state: `wire.EncodeQuery[T]`

`wire.DecodeForm[T]` shipped the request direction; the missing half is
state-out: filters/sort/page/cursor serialized back to a query string for
FilterChips hrefs, pagination links, reset links, and shareable views. A
generic encoder with round-trip property tests (`Decode(Encode(x)) == x`)
closes the loop and kills hand-rolled `url.Values` soup in consumers.
**First step:** godoc-first design against the filter-bar recipe's real state
struct; round-trip fuzz before any component consumes it.

### 4. `forms.AutoForm[T]` (demand-gated)

Struct-tag-driven form rendering: field type → Input/Select/Textarea, `validate`
tags → ErrorAttrs wiring, one constructor for CRUD admin forms. High leverage,
high surface — per the demand gate, build only on a named consumer need, and
scope v1 to a read-only REFLECT-free codegen variant if tags prove too magic.
**First step:** pre-write the demand-gated ADR (the #81 pattern) so the design
questions are answered the day demand arrives.

## B. Server & runtime platform

### 5. Early Hints (HTTP 103) support in `layout.Base`

`Base` controls every `<head>` asset; it is uniquely positioned to emit
`Link: </app.css>; rel=preload` headers before the body renders, cutting
TTFB-to-paint on slow origin paths. An `EarlyHints()` hook on PageProps (or a
`Base.WriteEarlyHints(w)`) is pure progressive enhancement — browsers without
support ignore it.
**First step:** spike in the demo binary (Cloud Run must be confirmed to
forward 103s); measure, then decide API shape.

### 6. Server-resolved theme (cookie/query-driven, zero-script dark mode)

`ThemeScript` prevents FOUC but costs an inline script on every page. A
server-resolved alternative: read a theme cookie (or `?theme=`) in middleware,
emit the `dark` class server-side, flip the cookie from ThemeToggle's form
endpoint. Fully CSP-minimal, works with JS disabled. Distinct from the ROADMAP
"configurable dark-mode selector" item (that re-targets the class/attribute;
this removes the client round-trip entirely).
**First step:** document the pattern against the existing ThemeToggle, then a
`layout.ThemeMiddleware` reference implementation + recipe.

### 7. Render observability hooks

Server-rendered UIs have no visibility into WHICH component is slow. A
context-carried observer (`utils.RenderObserver(ctx, func(Span))`) letting
middleware time component renders (templ wraps components; a hook can wrap
section-level components) turns "page is slow" into "StatCard grid renders in
80ms". Opt-in, zero overhead when absent.
**First step:** measure the wrapping overhead on the demo's heaviest route; ADR
only if it stays in nanoseconds.

### 8. Fragment memoization for expensive subtrees

Dashboards re-render identical expensive subtrees (charts, big tables) per
request. A `memo.Cached(children, key, ttl)` wrapper (in-process LRU keyed by a
props hash) trades staleness for latency where data is known-slow. Explicitly
opt-in per call site; document the invalidation sharp edges honestly.
**First step:** spike on the demo index route (the known-heavy page); publish
numbers before deciding it's a library feature vs a docs pattern.

### 9. Flash messages (post/redirect/get) as a first-class pattern

Every server-rendered app reinvents "save succeeded" across a redirect: cookie
or session flash → next render as Toast/Alert. Ship `feedback.Flash` (read,
render, clear helpers; http.Cookie plumbing; both Toast and Alert shapes) plus
a recipe. CSP-safe, no JS, works with hx-redirect and plain PRG alike.
**First step:** recipe doc first (zero library surface needed), promote to
helpers if the recipe's boilerplate is real.

### 10. Conditional-GET helper for poll endpoints

`PolledRegion` consumers hand-roll handlers that always render; polling every
5s re-sends identical HTML. A `utils.ETagOf(rendered)` + `http.ServeContent`-style
304 helper (hash the rendered fragment, compare `If-None-Match`) makes poll
endpoints nearly free when data is unchanged.
**First step:** document the pattern with a worked PolledRegion example;
promote to a helper only if the hash/compare boilerplate is genuinely error-prone.

### 11. Typed CSP builder (`utils/csp`)

Nonces are threadthrough-the-components today, but the POLICY header itself is
hand-rolled string soup in every consumer (the demo pins its CSP in two places
that must be kept equal — the repo's own footgun). A typed policy struct
(directives, sources, nonce injection points) with `String()` rendering and a
test that component nonce requirements ⊆ emitted policy closes the loop:
components declare, the builder proves.
**First step:** model the demo's existing CSP as the golden test case; the
firebase.json `keep-equal` invariant becomes a generated artifact.

### 12. Trusted Types compatibility audit

CSP next stop: `require-trusted-types-for 'script'`. Our inline singleton
scripts use `innerHTML`-free DOM APIs today (verify), but consumers composing
htmx + our scripts need a documented stance. An audit pass + docs page ("which
components survive Trusted Types, which sinks to add") is cheap now and a moat
later.
**First step:** flip the demo CSP to report-only with the TT directive and read
the violation report in CI's capture run.

## C. Component gaps

### 13. `display.Prose` typography wrapper

Markdown/CMS-rendered HTML has no home: every consumer hand-writes the
`[&_h2]:… [&_p]:…` soup or ships @tailwindcss/typography. A `.tc-prose` class
family (headings, lists, blockquotes, code, tables, images — dark/RTL/motion
compliant, size variants) + a thin `Prose` wrapper component owns "rendered
document content" the way Card owns cards. The base typography CSS in
`templates/custom.css` is the seed.
**First step:** extract the existing typography block into `.tc-prose`, golden
a markdown-rendered sample, expand coverage from there.

### 14. `forms.FileUploadProgress`

`FileInput` stops at the input; upload PROGRESS (the actual transfer) is where
server-rendered apps get stuck. XHR/`fetch` progress events → width-bound
ProgressBar, wired as a drop-in `Wire`-compatible component (htmx
`hx-post` encoding or Datastar), graceful no-JS fallback to the plain form.
**First step:** spike the event plumbing against a stub endpoint in visualtest
before committing to the API.

### 15. `forms.OtpInput`

Verification-code inputs (2FA, email codes) are a top-10 request class
everywhere: N single-char boxes, auto-advance, paste-splitting, backspace
retreat, auto-submit on complete. Radio/peer patterns + one CSP-safe singleton
script cover it.
**First step:** demand-gate check, then a pre-written ADR so it's build-ready.

### 16. `forms.PasswordInput` (strength + reveal)

Password fields deserve the full treatment: show/hide toggle (hit-box pattern,
aria-pressed), strength meter (server-computed zxcvbn-class score passed as
0-4, never client-rolled crypto), breach-check hook (`hx-post` to a checker),
confirmation-match error state.
**First step:** ADR on where strength scoring lives (server-side, per
security-first policy); component renders the RESULT.

### 17. Global `KeyboardShortcuts` registry + `<kbd>` display

CommandPalette owns ⌘K; the general case (single-letter nav, `?` cheatsheet
dialog, per-page hint registration) is unowned. A CSP-safe singleton registry
(data-attributes declare bindings server-side, script binds them) + a
`ShortcutsDialog` listing what's active + a `Kbd` display atom.
**First step:** extract CommandPalette's hotkey wiring into the shared script
shape, then grow the registry API from it.

### 18. `display.Presence` (live viewers / who's online)

SSE infrastructure exists (LiveRegion); presence is the missing consumer
pattern: avatar stack of currently-connected users, count pill, join/leave
announcements (polite live region). Component + reference SSE handler + recipe.
**First step:** recipe-level first (it's 90% server wiring); componentize the
avatar-stack render once the pattern is proven in the demo.

### 19. Multi-step wizard recipe (Stepper forms)

`StepIndicator` shows progress; the FLOW (per-step server validation, back
without data loss, resume-from-step, final review page) is a recipe gap. One
canonical recipe: state in a signed cookie or server session, step routes,
StepIndicator + Form composition, error re-render per step.
**First step:** write the recipe against the demo (checkout-style flow), then
decide if any helper deserves promotion.

### 20. `Table` sticky header (+ optional sticky first column)

Long tables under a sticky AppShell header need `position: sticky` thead with
a correct top offset (the offset is the hard part: AppShell header height →
CSS var). Zero JS, big usability win on the exact dashboard screens this
library targets.
**First step:** CSS-var contract with AppShell (`--tc-shell-header-h`), visual
golden at a tall-table demo route, RTL + scroll-container edge cases first.

## D. Docs & website

### 21. Docs pages embed the LIVE goldens

Docs can lie; goldens cannot. The docs generator already reads repo files
(DocRef); extend it to inject the actual golden HTML for each component's
example into the docs page — every docs snippet is then the SAME artifact the
test suite asserts on. Drift between docs and behavior becomes impossible by
construction.
**First step:** one pilot component (Card): docs page renders its golden
verbatim, styled by `.tc-prose` (#13) in an iframe or scoped wrapper.

### 22. Versioned demo deployments

"Was that behavior different in v1.19?" is unanswerable today. Deploy the demo
per release line (path- or subdomain-per-version from the existing Cloud Run
pipeline, images tagged with the version). Enables visual diffing across
versions, and gives every blog post/issue a pinned repro URL.
**First step:** tag the existing demo image per release in the deploy workflow
(cheap), version-picker UI later.

### 23. Component state-matrix documentation page

Goldens capture states; humans can't browse goldens. A generated "states"
section per interactive component: default/hover/focus/disabled/open/loading
rendered side by side (from the same StateHover/StateFocus visualtest
machinery), so reviewers and consumers see the full behavioral surface in one
place.
**First step:** generate the matrix for ONE overlay (Dropdown) from existing
capture options; evaluate readability before scaling.

### 24. "Componentshots" animated capture pipeline

ogshot proved the brand-card pipeline; the same chromedp harness can capture
SHORT interaction loops (dropdown opens, tab switches, toast flies in) as
GIF/WebP/MP4 per component — README, docs, and social posts all feed from one
sanctioned generator, always current with the release they document.
**First step:** one manual capture (Modal open) with the existing tools;
pipeline-ize only if the output is good without hand-tuning.

### 25. Migration guides FROM competitor libraries

Adoption friction is "I have 40 templUI/shadcn-templ components in my app —
now what?" Mapping tables (their component → ours → prop translation notes →
gaps) convert the comparison doc's honesty into an on-ramp. Same DocRef
single-sourcing pattern as /comparison.
**First step:** the shadcn-templ table first (biggest overlap), grounded in
their llms.txt inventory with the re-verify cadence from TODO #376.

### 26. Public component-authoring guide

The repo's SKILL.md encodes the entire component playbook (props shape, enum
discipline, testing ladder, guards) — for AI agents. A human-facing,
sanitized version ("How to contribute a component") in CONTRIBUTING or the
site turns that private moat into external contribution capacity.
**First step:** distill SKILL.md Part 2 + plan-authoring-checklist into a
docs page; link it from any first-timer-friendly issue.

## E. Testing & QA

### 27. Screen-reader verification protocol

axe catches violations; it cannot say "the drawer announces itself correctly in
NVDA". A documented manual SR pass (NVDA/Windows, VoiceOver/macOS, Orca/Linux)
per OVERLAY component — script, expected announcements, results ledgered next
to `axe_baseline.json` — upgrades the a11y claim from "no violations" to
"verified with screen readers".
**First step:** protocol doc + one honest pass over Modal and Drawer; ledger
the findings whatever they are.

### 28. Consumer-facing page-golden testing guide

`utils/golden` is battle-tested in-repo but undocumented for consumers testing
THEIR pages composed of our components. A guide + runnable example (golden
your page, update flow, CI wiring) makes "test your UI" a ten-minute adopt, and
anchors issue reports in reproducible HTML.
**First step:** extract an example from an existing package's
`golden_sweep_test.go` into a docs page + a starter template fragment.

### 29. Cross-platform compatibility lanes

The repo is Linux/Nix-only; consumers build on macOS and Windows, and pure-Go
portability bugs (path separators, symlink handling in the mirror guards,
golden line endings) would never surface in-repo. A Windows + macOS
compile-and-unit-test CI lane (no visual/browser jobs) is cheap insurance.
**First step:** Windows runner, `go build ./...` + the per-module unit loop,
visual jobs explicitly excluded.

### 30. `utils.Class` CPU/alloc budget + tailwind-merge fast path

Every component render pays `utils.Class` merging; the 7 benchmark packages
exist but nothing GATES regression. A benchstat-based CI check (fail >X% on
the Class-heavy benchmarks) plus a measured fast path (cache the merge of
known-constant class pairs — the compile-time-known strings dominate) keeps
render cost a managed property, not an accident.
**First step:** profile Class on the demo's heaviest route; the profile
decides whether the fast path is worth its complexity.

### 31. Release visual-impact report (class-diff diffing)

Semver for CSS output: between two versions, diff every golden's class
inventory and CLASSIFY changes (additive variant, reordering noise via the
existing tailwind-merge normalization, removal = potentially breaking). The
release script then prints "this release changes rendering for: Table rows
(removed `hover:bg-gray-50`)" — turning "should be visually compatible" into
an audited statement.
**First step:** one-off script diffing v1.20.1 vs master goldens; validate the
classification against known-intentional changes before wiring into release.sh.

### 32. Attribute-injection fuzz lane

The fuzz family covers inputs and wire actions; the remaining attack surface
is PROP VALUES flowing into attributes (`Attrs` maps, `Title`, `Href`, class
strings carrying quote/angle payloads). templ escapes by contract — the fuzz
proves OUR merging doesn't create smuggling shapes (attribute breakout,
duplicate-attribute override) the escaper wasn't shown.
**First step:** `FuzzHTMLDataAttrs` + `FuzzClass` with adversarial seeds
(quotes, `=`, `<`, newlines); assert rendering stays inside the attribute.

## F. Distribution & ecosystem

### 33. MCP server for AI-assisted adoption

llms.txt serves static docs to AI; an MCP server (`tc mcp` or a small
standalone binary) serves INTERACTIVE lookup: "which component does a date
range picker", "show ComboboxProps", "how do I wire a delete button" — from
the same registry that feeds goldens and `tc ls`. Meets the wave where
developers build with agents.
**First step:** the registry source of truth (#13 in 100-ideas) is the
prerequisite — sequence this behind it, then the MCP surface is a thin query
layer.

### 34. `tc lint` — export the compliance scanners to consumers

The repo's crown jewels include TestDarkModeCompliance, TestRTLLogicalProperties,
TestMotionReduceCompliance, TestInlineIconPathCompliance — all scanning .templ
source. Every CONSUMER writing their own pages against this library wants the
same guards. `tc lint ./...` runs the scanners over a consumer's templ files
(the scanners already do filesystem+regex work; the work is packaging +
documented allowlists).
**First step:** extract one scanner (RTL) behind a `tc lint` subcommand flag
set; dogfood it on the website module's hand-written pages.

### 35. Release provenance attestations

SSH-signed tags exist; SLSA-style build provenance (`actions/attest-build-provenance`
on release artifacts, Go module provenance story documented) is the next supply-chain
tier — increasingly a procurement checkbox for serious adopters. Paired with a
"verify a release" doc (tag → provenance → checksums).
**First step:** enable attestations on the release workflow for the demo
image; document module-proxy verification (GONOSUMDB/GOSUMDB story) for the
modules themselves.

### 36. Security posture dry run

SECURITY.md exists; it has never been exercised. A dry run (private-vuln-report
path, advisory draft, supported-versions table published, govulncheck on a
schedule with the results badge) converts a policy file into a practiced
capability — and the supported-versions table (which releases get fixes) is an
unmade policy decision hiding in plain sight.
**First step:** write the supported-versions table (likely: latest minor only),
schedule govulncheck weekly, draft one fake advisory end-to-end.

### 37. Conference + article distribution engine

The library's differentiators (dual-transport wire, CSP discipline, golden
culture, "100% inert until wired" honesty) are TALK material, and the demo is a
live artifact for any stage. One flagship talk ("Type-safe, CSP-safe
server-rendered UI in Go") + a monthly engineering-article cadence (the SSE
audit and optimistic-UI writeups already prove the genre works) compounds
distribution in a way listing-sites don't.
**First step:** abstract + CFP draft for one Go conference; the article backlog
already exists in ROADMAP harvest tables.

### 38. Framework-adapter starter gallery

templ needs only `http.Handler`, but consumers arrive from chi/echo/gin habits.
A gallery of minimal mounts (one page, auth-protected form, flash message) per
router — each ~50 lines, tested, linked from Getting Started — removes the
"will it work with MY stack" doubt. Stdlib stays the canonical example.
**First step:** chi first (largest Go-router mindshare), as a docs page with
the demo's own mux as the stdlib reference.

## G. i18n depth

### 39. Locale → `dir` auto-resolution + locale-aware formatters

`Base` takes `locale` for `lang`, but `dir="rtl"` is left to the consumer, and
`RelativeTime` is English-only. `utils.Locale` (BCP47 parse → script direction,
plural rules, number/date formatting via `golang.org/x/text` — the one new
dependency this list would ever justify, behind an opt-in subpackage) makes
`dir` unmissable and formatting correct in Arabic/German/Japanese alike.
Complements the ROADMAP's label-props direction (which handles STRINGS; this
handles STRUCTURE).
**First step:** spike direction inference + RelativeTime in de/ar/ja against
ICU ground truths; the ADR lives or dies on x/text's weight budget.

### 40. Message catalog convention for validation + error pages

The label-props direction (ROADMAP g2) solves component labels piecemeal; the
systematic shape is a message-catalog convention: one typed struct of every
user-facing string the library emits (error pages, validation summaries,
list notes), English zero value = current render, consumer swaps the struct.
Plurals handled by function fields, not format-string hacks.
**First step:** inventory every hardcoded user-facing string (scanner-assisted
grep of `templ` files for text nodes), shape the struct from the inventory,
then migrate errorpage first (smallest, most-labeled package).

## H. Performance

### 41. CI benchmark regression gate

Benchmarks exist in 7 packages; nothing fails when they regress. A scheduled
(benchstat old…new) job — comment on PR when render-heavy benchmarks move
beyond threshold, block only on sustained master regressions — keeps the
"server-rendered should be FAST" claim load-bearing as the library grows.
**First step:** benchstat the last 3 releases' benchmark history to calibrate
real variance before picking thresholds.
**Converges:** 100-ideas #99 + TODO #214 already track the benchstat PR
COMMENT (deferred 2026-09-13, needs baseline storage). The delta here is
(1) actually landing #214 and (2) extending it to gate sustained MASTER
regressions, not just comment on PRs.

### 42. `Image` modernization: fetchpriority + LQIP tooling

`Image` has lazy-load, SrcSet, and zero-JS blur-up placeholders; missing are
`fetchpriority` (hero images: `high`), explicit preload hooks for LCP images,
and an ogshot-style LQIP GENERATOR (build-time tiny blur placeholder from the
source asset) so the blur-up path is one command, not a research project.
**First step:** `FetchPriority` prop (trivial, golden-tested) + docs on LCP
pairing with Early Hints (#5); LQIP generator as a `visualtest/tools` follow-up.
**Adjacent:** 100-ideas #86 (preconnect/fetchpriority audit for the htmx asset
path) — same attribute, different surface (asset links vs the Image component).

## I. CI & demo

### 43. Visual-diff PR bot

Visual goldens change on every meaningful PR, but the evidence lives in CI
logs. A workflow step that posts/updates a PR comment with side-by-side
captions for changed goldens (light/dark pairs, mismatch percentages, links to
the `.fail/` artifacts) turns visual review from trust into looking.
**First step:** run on this repo's next component PR manually (actions
artifact + comment step), automate only the proven shape.

### 44. Public per-route quality badges

The axe sweep, HTML validation, and visual suite run on every demo route —
the results are invisible. A small build step emits a JSON quality report per
route (violations: 0, validated: yes, goldens: fresh) rendered as a
"quality ledger" page on the site. Honest marketing: the claims become
checkable.
**First step:** emit the JSON from the existing sweep job; the page is a
trivial table over it.

### 45. `hx-boost` decision + support

hx-boost (progressive enhancement turning normal links into hx requests with
body swapping) is htmx's highest-leverage page-level feature and the library
has no stance: no Base flag, no docs, no tested pattern. Decide: support
(`Base.HxBoost` + scroll/focus/heading-swap recipes) or reject with a written
reason (ViewTransitions overlap). Either answer is fine; none is not.
**First step:** spike in the demo (boost the docs section), catalog the
sharp edges (scripts, focus, scroll), then the ADR writes itself.

## J. Accessibility

### 46. Cognitive accessibility pass (plain language)

WCAG beyond the checklist: reading level, plain-language labels, consistent
term usage across ALL user-facing strings (error pages, list notes, empty
states). One pass cataloging current strings against plain-language guidance
(+ the message catalog #40 as the fixing mechanism) addresses the a11y
dimension no scanner can see — and may be the differentiator in the a11y
story: "readable by default".
**First step:** dump all user-facing strings (#40's inventory doubles here),
flag jargon/idiom/length outliers, fix the worst ten.

## K. Charts

### 47. Chart data-table fallback

Charts are `role="img"` + aria-label today: the DATA is invisible to screen
readers and copy-paste. Every chart already has the data in Go — emitting a
visually-hidden (or `<details>`-collapsed, or sr-only-until-focus) `<table>`
with the series costs ~20 lines per chart and upgrades chart a11y from
"described" to "readable". Progressive: `DataTableFallback *bool`.
**First step:** LineChart pilot with a golden asserting the table shape;
A11y review decides sr-only vs details presentation.

## L. Strategy & meta

### 48. Public roadmap page (ROADMAP.md → site docs)

ROADMAP.md is the most honest document in the repo and it's invisible. The
DocRef single-sourcing pattern (comparison.md precedent) renders it as a site
page: signals momentum, sets expectations, and lets contributors find
demand-gated work to pick up. Harvest tables become the "ideas we're thinking
about" section, openly.
**First step:** DocRef entry + frontmatter synthesis (the comparison.md
mechanics already exist); decide what (if anything) stays repo-private first.

### 49. Consumer fleet dashboard

The fleet-closure work proved the loop: consumer repos (crm, sperrmuell-direct,
PapDashboard, …) each pin a version, use a component subset, and surface
adoption blockers. A small local dashboard (scan consumer go.mods + imports →
version matrix + component-usage heatmap) turns "what do consumers need" from
memory into data — feeding the demand gates (#4, #217) with evidence instead
of recollection.
**First step:** a throwaway script scanning the known consumer checkouts;
if the output changes one roadmap decision, promote it to a maintained tool.

### 50. Email-safe render mode (research spike)

Transactional email is the sibling surface server-rendered web apps can't
escape, and the library's value system (type-safe, CSP-safe, graceful
degradation) maps onto email's constraints (table layout, inline styles, no
JS) better than any Go option today. NOT a component family yet: a research
spike — what breaks (Tailwind utilities in email clients), what a
`render.EmailMode` would even mean, whether a 5-component email subset
(Button, Card-shell, Alert, Layout-table) is honest or a trap. Written up as
an ADR, likely ending in "not planned" — but the honest answer, not the
default one.
**First step:** time-boxed spike rendering Alert + Button with inlined
styles through a litmus-style checker; the ADR is the deliverable, not the code.

---

## Reading order (if only five)

1. **#1 Context-injected defaults** — removes the single largest per-props papercut.
2. **#13 Prose** — closes the most-obvious catalogue gap for real apps.
3. **#34 tc lint** — turns internal quality infrastructure into adoption surface.
4. **#21 Docs embed live goldens** — makes the docs structurally incapable of lying.
5. **#33 MCP server** — meets adoption where it now happens (agent-assisted).

## provenance

- Generated 2026-10-14, big-picture ideation session.
- Dupe-checked against: `docs/100-IMPROVEMENT-IDEAS.md` (2026-09-13 corpus),
  `docs/STANDOUT-IDEAS.md`, `docs/SUPERB-FOR-PERSONAL-USE.md`,
  `ROADMAP.md` harvest tables, `TODO_LIST.md` #1-#400.
- Explicit convergence flags: #6 vs ROADMAP "configurable dark-mode selector";
  #39/#40 vs ROADMAP "i18n label props"; #33 vs 100-ideas #13 (registry);
  #31 vs ROADMAP parity-sketch; #32 extends the existing fuzz family;
  #30 complements 100-ideas #54/#55 (byte budgets vs CPU gate);
  #41 converges with 100-ideas #99 / TODO #214 (benchstat comment — delta is
  master-regression gating); #42 adjacent to 100-ideas #86 (fetchpriority audit).
