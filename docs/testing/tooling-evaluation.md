# UI Testing Tooling Evaluation — verified facts, spike results, verdicts (2026-10-15)

**Status:** evidence log + verdicts-in-progress for the UI-testing-tooling adoption decision
(`docs/planning/2026-10-15_00-05_UI-TESTING-TOOLING-PARETO-PLAN.md`). Every claim below is
marked **verified** (with source + date) or **blocked/demoted** (with reason). Owner gates
⫱G1–G3 stay with the owner; this doc's job is to make each gate a one-line decision.

**Companion docs:** ADR draft `docs/adr/0048-ui-testing-tooling-adoption.md` · coverage map
`docs/testing/demo-coverage-map.md` · research status
`docs/status/2026-10-14_23-59_ui-testing-tool-research-status.md`.

---

## 1. Verified facts table (M01, updated 2026-10-15)

| # | Claim                                                           | Verdict                                | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - | --------------------------------------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| a | nixpkgs `playwright-driver.browsers` usable for Node Playwright | **VERIFIED — GO (empirically proven)** | Locked nixpkgs ships playwright-driver **1.63.0** with `browsers` = linkFarm of `chromium-1243`, `chromium_headless_shell-1243`, `firefox-1543` (Firefox 155.0), `webkit-2359` (WebKit 26.6), `ffmpeg-1011`. Probe (2026-10-15): `nix build --inputs-from . nixpkgs#playwright-driver.browsers` + npm `playwright-core@1.63.0` + `PLAYWRIGHT_BROWSERS_PATH` + `PLAYWRIGHT_SKIP_VALIDATE_HOST_REQUIREMENTS=true` → **all three engines launched and screenshotted** on this machine. Caveats: npm playwright-core version MUST equal the nixpkgs playwright version (revisions are hardcoded per release — NixOS Wiki); WebKit may need `PLAYWRIGHT_HOST_PLATFORM_OVERRIDE="ubuntu-24.04"` (not needed for this probe). Sources: `pkgs/development/web/playwright/driver.nix`, wiki.nixos.org/wiki/Playwright. |
| b | rustywind maintained + packaged                                 | **VERIFIED — GO**                      | v0.28.0 (2026-08-25); 2026 brought a burst of releases (v0.26 rewrite: template-aware parsing, JSX, Tailwind-prettier parity). Packaged: npm `rustywind` (platform binaries), nixpkgs `pkgs/by-name/ru/rustywind` (0.28.0, nix-update-script), cargo, Homebrew. Source: github.com/avencera/rustywind releases.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| c | stylelint + Tailwind v4 at-rules feasible                       | **VERIFIED — GO**                      | Either `at-rule-no-unknown/ignoreAtRules: [theme, source, utility, variant, custom-variant, plugin, apply, reference, config]` on any stylelint 16.x, or `@dreamsicle.io/stylelint-config-tailwindcss` (stylelint ≥16.17) which registers v4 directives + `theme()` as real syntax (listed on stylelint.io/awesome-stylelint; no official Tailwind config exists). Sources: stylelint docs, tailwindlabs discussion #19252.                                                                                                                                                                                                                                                                                                                                                                                   |
| d | chromedp surface covers CDP vision/media emulation              | **VERIFIED — GO**                      | `github.com/chromedp/cdproto/emulation` (v0.0.0-2026-09…) exports `SetEmulatedVisionDeficiency` with typed values none/blurredVision/reducedContrast/achromatopsia/deuteranopia/protanopia/tritanopia, plus `MediaFeature` for `SetEmulatedMedia` (prefers-reduced-motion, print). Confirmed by direct inspection of the module cache, 2026-10-15.                                                                                                                                                                                                                                                                                                                                                                                                                                                            |

## 2. k6 + xk6-sse probe (M02/μ05, 2026-10-15)

- k6 in locked nixpkgs: **v2.2.0**. Automatic extension resolution is **enabled and works**
  (proven: `k6 x explore` transparently provisioned a custom binary from the registry).
- **BUT `xk6-sse` is NOT in the resolution registry**: import of `k6/x/sse` fails with
  `unknown dependency : k6/x/sse`; `k6 x explore` lists 14 modules, xk6-sse absent.
- The documented fallback (custom binary via `go run go.k6.io/xk6@latest build --with
  github.com/grafana/xk6-sse@latest`) **fails in this environment** — the Go module proxy
  is unreachable from this machine (same restriction that hit plain `go get`).
- **Consequence for M12:** the k6+xk6-sse scenario is written and checked in, but the load
  numbers in §5 were produced by an equivalent **Go-based SSE load client** (same
  steady-state / reconnect / client-drop profile) because no k6 binary here can carry the
  extension. In CI — where the Go proxy IS reachable — building the custom k6 binary is a
  one-line step and the scenario files run as-is.

## 3. Bencher cost/benefit note (M02/μ06, 2026-10-15)

Self-hosting Bencher = 2 Docker containers (Rust API on :6610 + console on :3000) with
embedded SQLite and no external DB; the OSS edition is free including production use
(public-project feature set, no artificial caps). Bencher Cloud free tier: public
projects only, 65,535 metrics/day. templ-components is a public repo → the **cloud free
tier requires zero infrastructure**; the only cost is an owner account + a CI-wired
`bencher run "go test -bench"` step. vs benchstat+committed baselines: benchstat is
stateless, zero-infra, statistically sound, but manual (no trends, no auto-gating);
Bencher adds time-series history, automatic regression thresholds, and baseline freshness
— at the cost of a service dependency and CI minutes. **Recommendation:** benchstat
remains sufficient for 7 packages at today's benchmark cadence; revisit Bencher only if
benchmark regressions start reaching master (kill-criteria in M24 budget table). Sources:
bencher.dev/pricing, bencher.dev/docs (self-hosted), fetched 2026-10-15.

## 4. Harness inventory — extend-not-duplicate map (M03, 2026-10-15)

| Capability                                                            | Where                                                                              | Extend, don't duplicate                                                         |
| --------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| Element screenshot + pixel compare vs golden, `-update`               | `visualtest/harness.go` `AssertScreenshot`                                         | Any new per-component visual check                                              |
| Full-page route capture (shared demo binary)                          | `visualtest/route_golden_test.go` `TestDemoRouteGoldens` + `assertRouteScreenshot` | New route goldens: one table row per capture                                    |
| Route list for audits                                                 | `visualtest/axe_sweep_test.go` `axeSweepRoutes`                                    | New routes: one entry (+dark) in the same commit                                |
| Vendored axe-core run + baseline ledger                               | `visualtest/axe.go` (`RunAxe`, `AxeViolationKey`) + `testdata/axe_baseline.json`   | Second-engine findings go to their OWN ledger (M20), never shoehorned here      |
| Poll helpers (bool-safe/text)                                         | `visualtest/poll.go` (`pollBool/pollTrue/pollText`)                                | Any in-page predicate waiting                                                   |
| Interaction states (hover/focus/click/context), animation-settle wait | `harness.go` `capture`, `waitAnimationsSettled`                                    | State-dependent captures via `Options.State`                                    |
| Keyboard traversal audit                                              | `keyboard_traversal_test.go` `TestKeyboardTraversalFocusVisibility`                | EXTEND the routes table + per-focus assertions, never a second traversal engine |
| Flake retry (dormant, zero call sites)                                | `flake_retry_test.go` `RetryOnce`                                                  | Leave dormant per flake-policy rule 5                                           |
| Standalone e2e fixture server                                         | `harness.go` `e2ePageServer`                                                       | Per-component e2e suites with custom endpoints                                  |
| Demo server for tests                                                 | `visualtest/demo.go` `StartDemoServer`                                             | All route/e2e tests                                                             |

CDP cheap wins (M11) therefore land as: new Go test files using `cdp.Call` with
`emulation.SetEmulatedVisionDeficiency` / `emulation.SetEmulatedMedia` on top of
`newTab` + `StartDemoServer` — no new browser plumbing.

## 5. Spike results (filled as phases complete)

### M06 — Playwright 3-engine cross-engine screenshot spike

- Scaffold: `tools/crossbrowser-spike/` (Node, playwright-core@1.63.0, zero repo test
  surface — guardrail 1). Route list read from `examples/demo/pages.go` (single source of
  truth) via a tiny generator.
- Result: see §5.1 table below (filled after the run).

<!-- §5.1 engine findings table -->

**Run (2026-10-15):** 16 routes × 3 engines = 48 full-page screenshots (light theme,
1280 viewport, localStorage theme pin + `load` + 1200 ms settle). **0 pageerrors, 0
console errors** in the final pass. Artifacts: `tools/crossbrowser-spike/artifacts/
<engine>-light/` (per-route PNG + report.json); the pre-fix run is preserved LOCALLY
at `tools/crossbrowser-spike/artifacts-m06-initial/` (gitignored — binary weight;
the CSP-violation evidence is quoted in the findings table below).

**Findings (all three were real, all fixed in this repo):**

| # | Finding                                                                                                                                                                                                                                            | Engine(s)                                                                          | Class                  | Fix                                                                                                                                                                                                                                     |
| - | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | ThemeToggle script rendered nonce-less (`Nonce: ""`) → CSP-dead under the demo's nonce CSP (silent feature kill, exactly the `nonce` omit-empty trap)                                                                                              | all                                                                                | real — feature dead    | `navigation_demo.templ`, `recipes_demo.templ` (2 sites): pass `demoNonceConst`                                                                                                                                                          |
| 2 | Avatar images from `picsum.photos` blocked by `img-src`; the fastly.picsum.photos redirect target too                                                                                                                                              | all                                                                                | real — content missing | `examples/demo/basepath.go` demoCSP `img-src += https://picsum.photos https://fastly.picsum.photos`; `website/firebase.json` header kept EQUAL                                                                                          |
| 3 | Demo Modal/Drawer buttons called `onclick="tcOpenOverlay('id')"` — broken TWICE: `tcOpenOverlay` is undefined (the library emits `tcOpenModal`/`tcOpenDrawer`), and inline `onclick` is CSP-dead (`script-src-attr`, proven by the Firefox report) | all (Firefox reported the CSP violation; the undefined global killed every engine) | real — feature dead    | Library-level: CSP-safe declarative `data-tc-overlay-trigger="<dialog-id>"` attribute (document-level click delegation, singleton guard) added to the overlay script; demo buttons now `display.Button` + attribute; recipe doc updated |

**Fix verification:** post-fix runs report **0 errors on all 48 shots**; a click test
(`click [data-tc-overlay-trigger]` → `dialog.open === true`, Escape-close, second
trigger) passes in Chromium, Firefox, AND WebKit. Display-package goldens updated for
the extended overlay script; `/display` route golden re-verified.

**Verdict (μ26 → ⫱G1 material):** multi-engine testing is **practical in this
Nix/Go repo** — the spike went from zero to 48 verified captures in well under a day,
and it caught 3 real CSP bugs that Chromium-only chromedp testing structurally missed
(2 were silent feature deaths). Engine notes: PNG sizes scale FF ≈ +35 % (encoder, not
layout); SSE routes never reach `networkidle` (use `load` + settle). **Recommendation:
adopt as a diagnostic lane (⫱G1); findings-class value proven, integration cost low.**

### M07 — base-select Firefox fallback assert

- Target: `/stylable-select` demo route (M05). Result: see §5.2.

<!-- §5.2 fallback verdict -->

**Run (2026-10-15):** `tools/crossbrowser-spike/stylable-assert.mjs` over
`/stylable-select` in Chromium + Firefox (light theme). Script + raw report:
`tools/crossbrowser-spike/artifacts/m07-fallback/report.json` (per-select crops
alongside). Measured on all four demo selects (`stylable-region`, `stylable-channel`,
`stylable-env`, `native-region`):

| Probe                                            | Chromium (base-select)        | Firefox (no support)                                      |
| ------------------------------------------------ | ----------------------------- | --------------------------------------------------------- |
| computed `appearance` (Stylable selects)         | `base-select`                 | `auto` (invalid keyword dropped at parse → initial value) |
| Stylable select box                              | 384×62 (styled picker button) | 384×29 — **identical to the native comparison select**    |
| pre-selected value / options intact              | yes                           | yes                                                       |
| `selectOption()` functional pick                 | yes                           | yes                                                       |
| horizontal overflow (`scrollWidth > innerWidth`) | none                          | none                                                      |

**Verdict (μ28): degraded-but-usable — no bug, no upstream action beyond tracking.**
Firefox renders the Stylable selects EXACTLY like the demo's native comparison select
(same 384×29 geometry, native chevron), which is precisely what the demo page promises
non-supporting engines will show. No layout breakage, values and keyboard/programmatic
selection fully functional. μ29 follow-up seed = re-check when Mozilla 1944403/1958443
ships (see §7); no repro note needed. Chromium's styled picker measured 62px tall vs
29px native — the docs' claim that `base-select` opts into taller picker buttons is
confirmed.

### M08 — IBM accessibility-checker spike + axe diff

- Result: see §5.3.

<!-- §5.3 achecker diff table -->

**Run (2026-10-15):** `tools/crossbrowser-spike/achecker-scan.mjs` (IBM
accessibility-checker npm, rule pack 2026.07.13, puppeteer on the nix Chromium)
over all 16 demo routes. Raw reports: `artifacts/m08-achecker/` (initial) and
`artifacts/m08-achecker-final/` (post-fix); the README documents the
`PUPPETEER_EXECUTABLE_PATH` + `aChecker.close()` gotchas.

**Value proof — the second engine caught what axe structurally misses** (axe's
duplicate-id rules only cover active/aria elements; achecker checks all ids,
labels, and landmark names):

| achecker rule (class)                                                                                                                   | Initial (routes/elements) | After demo fixes | Classification                                                                                                                                                                      |
| --------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `element_id_unique`                                                                                                                     | 2 routes / 17             | **0 — FIXED**    | **new-real** (demo reused `name=`-derived ids: `status`, `name`, `q`, `email`, … across example forms on one page; every demo route now sweeps clean)                               |
| `form_label_unique`                                                                                                                     | 2 / 16                    | **0 — FIXED**    | **new-real** (consequence of the same duplicate ids — controls carried multiple labels)                                                                                             |
| `input_label_before`                                                                                                                    | 2 / 8                     | **0 — FIXED**    | **new-real** (label `for=` attached to the first duplicate control)                                                                                                                 |
| `aria_region_label_unique`                                                                                                              | 1 / 10                    | **0 — FIXED**    | **new-real** (both kanban demo boards defaulted to aria-label "Kanban board"; now transport-suffixed)                                                                               |
| `aria_landmark_name_unique`                                                                                                             | 4 / 18                    | 2 / 14           | known (remaining = 12 `<form>` landmarks on the wire showcase — page nature, not a defect; StepIndicator/kanban duplicates fixed via AriaLabel)                                     |
| `aria_attribute_valid`                                                                                                                  | 10 / 22                   | 2 / 2            | **fixed** (`role="group"` + label on the scrollable code `pre`) + noise (2 remaining: achecker rejects `aria-label` on `<footer aria-label="Footer">` — valid per ARIA, axe agrees) |
| `aria_navigation_label_unique`                                                                                                          | 3 / 6                     | 2 / 4            | known (shell nav landmarks; engine stricter than axe `landmark-unique`)                                                                                                             |
| `aria_descendant_valid`                                                                                                                 | 9 / 28                    | 9 / 28           | **new-real → TODO #421** (CopyButton `role="status"` inside `<button>` — invalid nesting; component change, not spike-safe)                                                         |
| `element_tabbable_role_valid`                                                                                                           | 7 / 18                    | 1 / 3            | **fixed** (code `pre`) + **new-real → TODO #423** (Table scrollable wrapper)                                                                                                        |
| `aria_complementary_labelled`                                                                                                           | 1 / 1                     | 1 / 1            | **new-real → TODO #422** (`layout.Split` unlabeled `<aside>`)                                                                                                                       |
| `label_content_exists`                                                                                                                  | 1 / 5                     | 1 / 5            | noise (Rating star labels carry `sr-only` text the engine ignores)                                                                                                                  |
| `label_ref_valid`                                                                                                                       | 1 / 8                     | 1 / 8            | noise (SegmentedControl `for=` ids verified matching in served HTML; sr-only radio heuristic)                                                                                       |
| `text_contrast_sufficient`                                                                                                              | 3 / 14                    | 3 / 14           | noise (solid-bg buttons: engine resolves the parent bg, not the button's own; axe `color-contrast` passes the same nodes)                                                           |
| `input_label_visible` (48), `label_name_visible` (16), `text_block_heading` (20), `text_sensory_misuse` (17), `style_color_misuse` (16) | —                         | —                | noise (intentional sr-only labels, icon-only buttons, stat numbers, demo copy)                                                                                                      |

**Net effect of the spike-driven fixes:** 330 → 253 total results (−23%),
duplicate ids eliminated on ALL 16 demo routes, violation count down on every
affected route (/kanban 15 → 2, /forms 41 → 24, /wire 46 → 34). Demo fixes:
transport-suffixed ids in the per-dialect wire demos, unique ids for every
reused field, distinct board/step-indicator labels, `role="group"` on code
snippets. Component-level findings are ledgered as TODO #421–#423, never
shoehorned into `axe_baseline.json` (M20 owns the separate-ledger decision,
⫱G2).

### M10 — aria-snapshot goldens spike

- Result: see §5.4.

<!-- §5.4 aria snapshot verdict -->

**Run (2026-10-15):** `tools/crossbrowser-spike/aria-snapshots.mjs` — Playwright
`locator.ariaSnapshot()` (available in playwright-core WITHOUT the
`@playwright/test` runner) over four representative component states; YAML in
`artifacts/m10-aria/`:

| Snapshot                   | Result                                                                                                                                                               |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Table (`aria-sort`)        | Roles/names/selected captured; sort arrow rendered inside the header NAME ("Name ↑")                                                                                 |
| Modal open (`role=dialog`) | Full tree: dialog → heading/button/paragraph — 7-line snapshot                                                                                                       |
| Dropdown trigger           | `- button "Actions"` — **no `aria-expanded` BY DESIGN** (native Popover `popovertarget`; the HTML-validation gate bans the combo, ATs derive state from the popover) |
| Tabs (`role=tablist`)      | tablist → tab "Overview" [selected] + siblings — 4-line snapshot                                                                                                     |

**μ40 regression demo:** removing every `th[aria-sort="ascending"]` from the
/display tables leaves the ariaSnapshot YAML **byte-identical** — Playwright's
ariaSnapshot does NOT expose `aria-sort` values. **Consequence for M21/⫱G1:**
aria goldens guard structure (roles, names, selected, expanded-via-popover),
NOT value-bearing attributes — `aria-sort`, `aria-valuenow` etc. stay the job
of DOM-attribute asserts in the chromedp suite (the existing
`TestTableAriaSort`-style substring asserts already own that axis). Two
capture notes: a CLOSED `<dialog>` snapshots empty (display:none — open it
first); /display hosts TWO sortable tables (mutation asserts must be
count-based, not first-match).

### M11 — CDP cheap wins in the chromedp suite (integrated)

**Landed (2026-10-15):** `visualtest/visual_emulation_test.go` — three tests on the
EXISTING chromedp plumbing (`newFlowTab` + `StartDemoServer` + raw `cdp.Call`):
`TestVisionDeficiencyProbe` (achromatopsia + deuteranopia over /display, /feedback,
/users — asserts emulation applies and captures review artifacts),
`TestReducedMotionEmulationKillsSpinnerAnimation` (prefers-reduced-motion via
`SetEmulatedMedia` — the spinner animation collapses to `none`, proving the
motion-reduce convention end-to-end), `TestPrintMediaProbe` (/users + /recipes under
print media). All green under the nix browser env; they skip cleanly without
`CHROMEDP_CHROME_PATH`.

**Human-review findings** (`tools/crossbrowser-spike/artifacts/m11-cdp/`):
vision-deficiency passes show NO color-only information on the probed routes —
statuses/pagination carry text, contrast survives desaturation (no contrast bombs,
μ42). Print probe (μ45): the demo has ZERO print styling — sidebar, header, theme
toggle, and action buttons all render on paper. Follow-up registered as TODO #424
(a minimal `@media print` layer: hide chrome, drop shadows), not a spike fix.

### M16 / M22 — skip-by-verdict (gates decided upstream)

- **M16 Bencher (μ61–μ62):** not run — M02's verdict (§3: benchstat suffices at
  current cadence) closes it; the revisit trigger (regressions reaching master) is
  recorded for the M24 budget table.
- **M22 k6 nightly (μ74):** not run — μ49's verdict (§5.5: no signal at 50 clients)
  closes it; kill-criteria drafted there.

### M17 — stylelint + rustywind trial (μ63–μ64)

- **stylelint 17.16 + `@dreamsicle.io/stylelint-config-tailwindcss`** over
  `templates/custom.css` + `examples/demo/demo.css`: **ZERO findings** (exit 0). The
  Tailwind-v4 at-rules (`@theme`, `@custom-variant`, `@source`, `@utility`) parse as
  real syntax with zero config. Verdict: the CSS is healthy; adoption value would be
  regression-guard only — **defer until the CSS surface grows** (M24 budget row).
- **rustywind 0.28 dry-run over the repo: 183 files would change** — including
  `*.golden` files (whose class order is set by the library's OWN class sorter) and
  Go test files. Verdict: **REJECT** — it fights the golden sorting convention and
  would churn goldens for zero behavioral gain. → rejected-tools log (M23).

### M18 — trace viewer as debugging aid (μ65–μ66)

- Recorded a kanban-interaction trace via `context.tracing` in playwright-core
  (`artifacts/m18-trace/kanban-move.trace.zip`, gitignored) — no `@playwright/test`
  runner needed. Verdict: **KEEP as a debugging skill, not a lane**: it becomes
  useful exactly when the crossbrowser lane (⫱G1) reports an engine-specific
  failure; chromedp's equivalent (flake-retry harness) stays dormant per policy.

### M19 — engine hardening (μ67–μ69)

- **Fontconfig:** the `.#crossbrowser` flake app pins the SAME pure fontconfig
  (Inter + DejaVu only) as the visual app — WebKit/Gecko/Chromium renders are
  host-font-proof (μ67 done in the app wiring).
- **Determinism (μ69):** consecutive runs are **pixel-identical (0/13,949,440)**;
  raw PNG bytes differ (encoder chunks) — compare pixels, never bytes.
- **Viewport policy (μ68):** engine captures run at 1280 width — the same
  below-`xl:` caveat as route goldens applies (any `xl:`/`2xl:`-gated layout change
  needs a throwaway custom-viewport probe before concluding "no pixel change");
  documented here + in the spike README.

### M20 — achecker ledger: gate-edge schema (⫱G2)

Wiring waits for ⫱G2; the schema is decided now so ratification is one line:
`visualtest/testdata/achecker_baseline.json` keyed `route` → `ruleId` → `{count,
level}` with the achecker semantics (`violation` / `needsreview` / `recommendation` /
`potentialviolation`), NEVER merged into `axe_baseline.json`; noise classes
documented in §5.3 stay OUT of the ledger (they are engine quirks, not debt). The
finding flow: new class → verify in served HTML → fix or ledger with justification →
delete the entry when it stops matching (mirrors the axe ledger rules).

### M21 — aria-golden class promotion path (⫱G1)

Promotion is ready but GATED: aria snapshots would become a 5th golden class
(`testdata/aria/`, `.yaml`, `-update` workflow identical to HTML goldens,
`TestDocsCountDrift` bump in the same edit) with the §5.4 boundary — they guard
STRUCTURE only (roles/names/selected), never value-bearing attributes. Until ⫱G1:
`tools/crossbrowser-spike/aria-snapshots.mjs` remains the manual capture path.

### M12 — SSE load profile (Go client stand-in, see §2)

- Result: see §5.5.

<!-- §5.5 load report -->

**Run (2026-10-15):** Go-client load profiles (`visualtest/tools/sseload`, the §2
stand-in for k6+xk6-sse) against the demo binary's `/api/datastar/stream` DIRECTLY
(2 s patch cadence, 15 s keep-alives):

| Profile   | Setup                                   | Result                                                                                                                                |
| --------- | --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| steady    | 10 clients × 20 s                       | 91 events; arrival gaps **mean 2000.05 ms, p95 2008 ms, max 2008 ms — 0 continuity violations**; keep-alives observed on every stream |
| reconnect | 10 immediate connect→first-event cycles | time-to-first-event **2002–2012 ms** every attempt (next 2 s tick — no handler leaks, no drift)                                       |
| drop      | 50 synchronized client drops            | server reaped all handlers; fresh post-drop client patched in **2005 ms**, service healthy                                            |

**Verdict (μ49): no nightly SSE lane NOW.** The single-binary Go SSE handler shows
sub-10 ms jitter at 50 clients with zero continuity breaks — a k6/nightly lane would
add toolchain cost for no signal. **Kill-criteria draft (→ M24 budget table):** adopt a
nightly SSE lane if (a) consumers report live-region disconnects, (b) the stream
handler grows real fan-out logic, or (c) #329's canonical-URL buffering is fixed and
prod-parity load becomes testable. The `sseload` tool stays in the unpublished
visualtest module as the on-demand re-run path.

## 6. Verdicts (M14) — to be ratified at ⫱G1–G3

Per-tool verdicts with evidence (μ55); the adopt/reject matrix + pipeline boundary
live in ADR 0048 (μ56/μ58: recommendation = **ADDITIVE** — chromedp stays primary,
Playwright is cross-engine-diagnostic-only, never a second pixel suite):

| Tool | Verdict | Evidence | Gate |
| ---- | ------- | -------- | ---- |
| Playwright (nixpkgs driver + npm playwright-core) | **ADOPT as diagnostic lane** (`nix run .#crossbrowser`) | §5.1: 48/48 shots, 3 real CSP bugs fixed; README: pixel-deterministic (0.0000% repeat delta); lane + dispatch-only CI draft ready | ⫱G1 |
| aria snapshots (`locator.ariaSnapshot()`) | **TRIAL verdict reached: structure-only goldens** | §5.4: guards roles/names/selected; does NOT expose `aria-sort` values — attribute asserts stay in chromedp | ⫱G1 (with M21) |
| IBM accessibility-checker | **TRIAL** (second opinion, own ledger semantics) | §5.3: caught duplicate ids/labels across 16 routes that axe structurally misses; −23% findings after demo fixes; 3 component seeds | ⫱G2 |
| k6 + xk6-sse | **REJECT local / KEEP scenario path for CI** | §2: extension not in registry; Go proxy blocked locally | — |
| SSE nightly lane | **REJECT now** (kill-criteria drafted) | §5.5: 2000±8 ms jitter at 50 clients, 0 violations, clean reconnect/drop | ⫱G4 |
| Go SSE load client (`visualtest/tools/sseload`) | **KEEP** (on-demand re-run path) | §5.5 | — |
| Bencher | **REJECT for now** (benchstat suffices; revisit trigger recorded) | §3 | — |
| stylelint | **DEFER** (CSS healthy — 0 findings; regression-guard value only) | M17: stylelint 17.16 + TW4 config, zero findings on custom.css + demo.css | — |
| rustywind | **REJECT** (fights golden class-sorting; 183-file churn) | M17: dry-run rewrites 183 files incl. `*.golden`; → rejected-tools log | — |
| Playwright trace viewer | **KEEP as debugging skill, not a lane** | M18: trace via playwright-core `context.tracing` (no test runner); useful when the crossbrowser lane fails | — |
| Lost Pixel | **REJECTED — dead** (archived 2026-04-22) | plan context; rejected-tools log (M23) | — |

Nothing here is integrated into the repo's gating test surface beyond the M11
chromedp tests (which extend the EXISTING pipeline, not a new one) until the
owner ratifies a gate.

## 7. Re-verification reminder (M26 seed)

All "verified" rows above are point-in-time (2026-10-15). Re-check if adoption slips
past ~6 months (nixpkgs playwright version moves fast; Firefox base-select lands via
Mozilla 1944403/1958443).
