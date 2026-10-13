# Consumer Usage Deep Review — 2026-10-12

**Scope:** how the ~30 LarsArtmann projects on this machine actually use
`github.com/larsartmann/templ-components` — setup, wiring, CSP, duplicates,
version pinning — plus the library's own dogfooding (website/, demo/).
Method: parallel read-only audits of every consumer with ≥5 library
references + the long tail, then source-level verification of every finding
before acting. Cross-references prior art: `docs/integration/extraction-analysis.md`,
TODO_LIST #397/#399/#400 (ADR-0044/0045 already cover theming tokens and CSS
class delivery — those known debts are NOT re-reported here).

## Verdict

The feeling is correct, and it is **two-sided**:

1. **Consumers under-use the library they already depend on.** The same five
   duplicates repeat fleet-wide: hand-rolled `<table>` markup, bespoke
   buttons, hand-written 404/500 pages, inline SVG icons, and raw `hx-*`
   attributes. `utils/wire.Action` — the flagship transport-agnostic wiring
   abstraction — has **1 serious production use across the entire fleet**
   (sperrmuell-direct, plus one site in cqrs-htmx adminui).
2. **The library's sharp edges caused real production bugs** — three of them
   found and FIXED during this review (below). The recurring root cause is
   nonce plumbing: components that render inline scripts silently die under
   the nonce-locked CSPs most of these apps run, because a consumer forgets
   to thread `BaseProps.Nonce`/the `nonce` argument at ONE call site.

## Fixed during this review (verified: build + full test suites green)

### 1. cqrs-htmx `dashboardui` — every CopyButton's clipboard script was dead (nonce bug)

- `dashboardui/components.templ` `copyButton` and `definitions.go`
  `copyButtonComponent` hardcoded `Nonce: ""`. Under the module's
  nonce-locked CSP (`csp_test.go` forbids `unsafe-inline`), the library's
  clipboard singleton rendered untagged and was **silently blocked on every
  page that renders a copy button** (events, audit, snapshots, aggregates,
  time-travel — 12+ sites). The existing CSP sweep only scanned `/dashboard/`
  (overview), which renders no CopyButton — a **vacuous pass**.
- Fix: `nonce` threaded through `copyButton`/`copyButtonComponent`/
  `defItemCopy` and the six detail/metadata builders into all 13 call sites
  (row templates gained a `nonce` param; detail templates pass `p.Nonce`).
  Goldens re-pinned (`nonce="test-nonce"` now in
  `event_detail_page.golden`, `snapshot_detail_page.golden`).
  `TestCSP_LibraryScriptsCarryNonce` now also sweeps the time-travel detail
  page and **fails if the copy singleton is absent** (no more vacuous
  passes). Regenerated with the repo's canonical generator (nixpkgs templ
  v0.3.1020 — go.mod's v0.3.1070 require is the runtime lib, not the
  generator; DO NOT regenerate this repo with 1070, it flips whitespace in
  every golden).

### 2. GmbH — dark-mode toggle and mobile menu dead under CSP

- `server/web/chrome.templ` `navBar` passed `""` as the ThemeToggle nonce and
  never set `SimpleNavProps.BaseProps.Nonce`, so the theme-toggle script AND
  the library's mobile-menu singleton rendered without a nonce — both
  blocked by the server's `script-src 'self' 'nonce-…'` middleware. Probe
  before: 2 live scripts / 2 dead. After: **4 nonce-carrying scripts / 0
  dead**. Also fixed the "Darkestellung" typo → "Darstellung". Goldens
  re-pinned.

### 3. crm-exec-stage — stale fork of `feedback.Alert`

- `internal/httpapi/layout.templ` hand-rolled the success flash banner with
  the comment "Alert hardcodes role=alert" — **stale**: the library has
  derived `role` from `Type` for a long time (sibling repo `crm` already
  uses `Alert` for both kinds). Replaced the fork with
  `tcutils.Ternary(…FeedbackError…FeedbackSuccess)` Alert, inverted
  `TestAlertRoleContract` to pin the new contract (success → `role="status"`,
  error → `role="alert"`). 83/83 behavior specs green; suite run compile-
  verified under go 1.27.2 + tidy (tidy itself reverted — see Caveats).

**Caveats / adjacent discoveries (not fixed):**

- crm-exec-stage is **build-broken for ad-hoc local builds** (pre-existing):
  the committed `go.mod` is out of sync with the sibling checkouts it
  replaces (`../cqrs-htmx/usermgmt` now requires go ≥ 1.27.2, pins say
  1.27.1, and `go mod tidy` has not been run after the sibling bumps). The
  nix path (prepared sources strip dev replaces) is the canonical build. It
  needs a coordinated toolchain-floor migration: go.mod directive,
  `.buildflow.yml` `GOTOOLCHAIN`, `flake.nix` `go_1_27` override (copy
  cqrs-htmx's landed pattern), `scripts/regen-ui.sh`,
  `scripts/check-ui-artifacts.sh`, `scripts/verify-pinned.sh`, and the
  tracked `web/tailwind/library-*` inventories (already stale vs their
  generator) — all in ONE commit.
- github-local-sync and timesheets carry templ-components only as
  transitive indirect pins — **not removable** (an earlier suspicion of
  "dead pins" was a misdiagnosis; `go mod tidy` keeps them). Both are also
  hit by the same go.work 1.27.1-vs-sibling-1.27.2 gate.

## Fleet version matrix (vs current v1.21.0)

| Version                                       | Projects                                                                                                                                                                                                                              |
| --------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| master pseudo-version (`v1.21.1-0.…a30aa632`) | CV (load-bearing: uses `HTMXNone`/`SEOMeta` — a release rollback breaks its build)                                                                                                                                                    |
| v1.21.0                                       | bank-sync, browser-history, cqrs-htmx (all UIs), crm, crush-daily, DiscordSync, file-and-image-renamer, go-health-dashboard, go-taskqueue, library-policy (indirect), nsfw-classifier, overview, pbx-lens, sperrmuell-direct, Zlota44 |
| v1.20.1                                       | crm-exec-stage, GmbH, InboxClean, PapDashboard, webphone, go-cqrs-lite/catalog                                                                                                                                                        |
| v1.19.4                                       | ChastityAPI, Cyberdom, dnsblockd, index, KeyHolderAI, Rolls-Royce-mtuGoHelpCenter-golang, github-local-sync (indirect), plugmarket, timesheets (indirect)                                                                             |

14 consumers lag ≥1 minor version. Bumps are mechanical (`go get … && go mod
tidy` + build/test) but must ride the `go-ecosystem-upgrade` skill loop —
several repos are locally build-broken by the go.work/sibling-pin mismatch
above, which blocks verification today.

## Cross-fleet anti-patterns (ranked by impact)

1. **Nonce hygiene — the fleet's #1 correctness gap.** Real dead-script bugs
   fixed in 2 repos today; the same pattern waits in: KeyHolderAI
   (`TCToastContainer("")`, `tc.templ:20`), pbx-lens
   (`ThemeToggle("…","")`, `calls.templ:36` + `ToastContainer("")`,
   `page.templ:15` — harmless today only because pbx-lens has no CSP yet),
   games/SwettySwipperWeb (**fake literal nonces**
   `nonce="toast"|"vs-keys"|…` — worse than none — plus ThemeToggle called
   with an element ID as its nonce, `layout.templ:135`), PapDashboard
   (un-nonce'd htmx script), crush-daily (`report.templ:51` standalone doc).
   Library-side root cause: nonce consumption is invisible at the call site
   (see TODO #407 below).
2. **`wire.Action` adoption ≈ 0** while every project hand-writes `hx-*`:
   DiscordSync ~45 raw sites, CV 0 usages, crm/crm-exec 0 (pure PRG),
   browser-history 1, bank-sync 0. InboxClean and sperrmuell show the
   target state. This is a discoverability + habit problem, not an API gap.
3. **Raw `<table>` is the most-reinvented wheel** (≈10 projects): crm +
   crm-exec (4 files each), DiscordSync (4), ChastityAPI (5), KeyHolderAI,
   Zlota44 (3), index/readmeindexer (3), CV. `display.Table` with
   `TypedHeaders`/`Flush` covers almost all of them.
4. **Hand-rolled 404/500 pages** despite the errorpage module being a
   dependency in most: sperrmuell (full duplicate of NotFound404 + ErrorPage,
   module already in go.sum as indirect!), webphone, go-taskqueue, SEC,
   index (CLI report variants), crush-daily, KeyHolderAI (partial — uses
   errorpage but bypasses layout), crm/crm-exec (plain-text `http.Error`,
   errorpage only indirect).
5. **Bespoke buttons/cards** flood consumer markup: browser-history (~25 raw
   buttons), GmbH (5 styled anchors), KeyHolderAI (`.btn-primary` classes),
   bank-sync (header dropdown/search/sparkline), CV (DashboardCard/Hero/
   Header), overview (project card + pill badges), DiscordSync
   (`primaryButtonClass` strings + attribute-selector CSS hacks in
   `input.css:193` to rebrand them).
6. **Library forks with drift risk:** DiscordSync `dsSidebar` (visual copy of
   SidebarNav + own JS), go-health-dashboard (ThemeToggle SVG copied from
   v1.16.0), file-and-image-renamer (v1.10.0 stagger code + pollIndicator
   copy), nsfw-classifier (vendored `third_party/_css-scan`),
   ChastityAPI (hand-trimmed `custom.css`). ADR-0044/0045 + #400 address the
   theming/CSS half; the component-copy half needs the swaps in #400.
7. **Physical Tailwind props** (`ml-/mr-/pl-/pr-/left-/right-`) in every
   audited consumer (~40 sites fleet-wide) — RTL-unsafe. The library's
   `TestRTLLogicalProperties` covers only library sources; consumers have no
   guard (a `tc doctor` lint could close this, see TODO #407).
8. **Dark-mode drift in consumer markup:** inverted pairs like
   `text-gray-500 dark:text-gray-600` (dnsblockd dashboard ×4,
   sperrmuell admin, crm-exec), missing ThemeToggle (dnsblockd, Zlota44,
   index-webapp), and `prefers-color-scheme`-only themes fighting the
   library's class strategy.
9. **CSP-less apps with inline handlers** (crm + crm-exec only): `onsubmit=
   "return confirm(…)"` ×4 each (dup of `htmx.ConfirmDelete`), `onclick`
   button attrs, bare `<script src>` — the moment a CSP lands, these die
   exactly like today's nonce bugs. ChastityAPI has the same time bomb
   (`templ.JSFuncCall` drawer trigger at `layout.templ:208`) with CSP_ENFORCE
   already planned.

## Model citizens (copy from these)

- **overview** — full nonce threading, errorpage with family mapping,
  SkeletonCardGrid, `layout.Script` everywhere.
- **go-health-dashboard** — exemplary Base+Container+PageHeader+Grid+StatCard
  - Datastar LiveRegion (only debt: the copied ThemeToggle SVG).
- **Cyberdom** — clean textbook mainstream adoption, zero misuse found.
- **Rolls-Royce-mtuGoHelpCenter-golang** — full module surface, per-library-
  pattern CSS workflow, errorpage wired.
- **go-cqrs-lite/catalog (docserver)** — static-doc embedding done right.
- **sperrmuell-direct** — the only serious `wire.Action` user (debounced
  GET quote check).
- **InboxClean** — real dual-transport (htmx + Datastar) with LoadMore/
  EndOfList working in both.

## Library dogfooding (website/ + demo/) — top gaps

The repo's own site hand-rolls what it sells (full list in the session
report; guards already cover dark/RTL/motion/nonce classes):

1. `website/internal/pages/hero.templ:106-118` + `assets/js/copy-code.js` —
   bespoke code window + copy button + duplicate of the clipboard JS;
   `display.CodeBlock` (shipped 2026-10-10 for THIS exact consumer demand)
   is one page away.
2. `website/internal/pages/header.templ:20-48` + `footer.templ:52-69` +
   `assets/js/header.js` — bespoke nav/mobile-toggle/footer re-implementing
   `navigation.Nav`/`MobileMenu`/`Footer` shipped behavior.
3. `examples/demo/errorpage_demo.templ:82-121` — playground form hand-rolls
   every control while the forms package is demoed next door.
4. `website/internal/pages/sections.templ:82-86,145-176,223-230` — feature
   cards/comparison table/CTA buttons vs `display.Card`/`Table`/`Button`
   (the landing page's "built with the library" claim is true on /sales only).
5. `examples/demo/demo.templ:149-171,238-253` — demo-local shadow
   `demoStatCard`/`demoPageCard`/`recipeLinkCard` vs `display.Card`/`StatCard`.

## Known-debt cross-references (do not re-litigate)

- Theming/CSS-scan fragmentation, DiscordSync inert `@source`, palette
  bridges → ADR-0044/0045, TODO #397/#399/#400.
- ErrorHandler-wrapper shunning → TODO #326.
- Unbounded tables → `docs/recipes/bounding-tables.md`.
- Promoted-BaseProps literal trap, zero-value≠default → AGENTS.md
  conventions + `docs/migration/v1-to-v2.md`.

## Actionable outcomes

Fixed today: items 1–3 above (cqrs-htmx nonce bug, GmbH nonce bugs,
crm-exec-stage fork retirement) + TODO_LIST harvest (#402–#407). The
remaining consumer-side work (raw-table/button/pagination swaps, errorpage
adoptions, physical-prop sweeps) is per-repo migration work → tracked in
TODO #400's wave structure; the library-side enablers (Button children slot,
Table row-class override, nonce/doctor guards) → TODO #405–#407.
