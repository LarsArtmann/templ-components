# ADR 0048: UI Testing Tooling Adoption — additive Playwright, vetted lanes

## Date

2026-10-15

## Status

**Proposed — pending owner gates ⫱G1–G3** (posed in
`docs/testing/tooling-evaluation.md` §6). Nothing in this ADR changes the repo's
test surface until the corresponding gate is ratified; the spike scaffolds live
outside the Go test surface by design (guardrail 1 of the plan).

## Context

The repo's browser-level testing is mature but Chromium-only (chromedp +
vendored axe-core + pixel goldens in `visualtest/`). Three claims are therefore
Chromium-proven only: cross-engine rendering, a second a11y opinion, and
load/reconnect behaviour. The 2026-10-14 research session produced a verified
shortlist; the 2026-10-15 pareto plan executed the spikes. Evidence lives in
`docs/testing/tooling-evaluation.md`; per-task artifacts are linked from there.

Key facts driving this decision (all verified 2026-10-15, sources in the
evaluation doc):

- nixpkgs `playwright-driver.browsers` (1.63.0: chromium-1243, firefox-1543,
  webkit-2359) works with npm `playwright-core@1.63.0` via
  `PLAYWRIGHT_BROWSERS_PATH` — empirically proven on this machine, all three
  engines.
- The M06 spike proved the value claim, not just the packaging claim: over the
  existing demo routes it surfaced **3 real CSP bugs that Chromium-only
  chromedp testing structurally missed** (two silent feature deaths: nonce-less
  ThemeToggle scripts, picsum `img-src` gap; one double-broken Modal/Drawer
  trigger that also hid an undefined JS wrapper). All three are fixed — the
  Modal/Drawer fix landed as a CSP-safe `data-tc-overlay-trigger` library
  feature (eval doc §5.1).
- `forms.Select` `Stylable` (`appearance: base-select`) is unsupported in
  Firefox (Mozilla 1944403/1958443) and now has a dedicated demo route
  (`/stylable-select`, 2026-10-15) with axe + route-golden coverage. The M07
  cross-engine assert measured the fallback **degraded-but-usable**: Firefox
  renders the Stylable selects identical to the native comparison select
  (384×29 vs Chromium's styled 384×62), zero layout breakage, selection fully
  functional (eval doc §5.2).
- xk6-sse is not in k6's automatic-extension registry; the custom-binary build
  is the viable path where the Go proxy is reachable (CI, not this machine).
- IBM accessibility-checker is maintained and brings an ACT-mapped ruleset
  independent of axe-core.
- Lost Pixel is dead (archived 2026-04-22) and stays rejected.

## Decision (recommended)

**1. chromedp stays the primary pipeline; Playwright is additive,
cross-engine-only** (⫱G3 recommendation: ADDITIVE). Concretely:

- `visualtest/` (chromedp) keeps: pixel goldens, axe sweep, e2e flows, CDP
  emulations, route goldens. No migration.
- Playwright lands (⫱G1) as `tools/crossbrowser-spike` promoted to a
  `nix run .#crossbrowser` app: 3-engine full-page screenshots of the demo
  routes + console/pageerror collection. It is a **diagnostic lane** — it
  produces artifacts and findings, it does not gate CI and it never becomes a
  second long-term pixel-diff suite (never two pixel suites).
- aria-snapshot goldens (⫱G1 for the class) would be the ONLY Playwright
  assertion type allowed to gate, IF M21 promotes them — they occupy a
  different axis (accessibility tree, not pixels) and don't duplicate
  chromedp's job.

**2. IBM accessibility-checker (⫱G2): TRIAL**, not adopt-by-default. Run it
sporadically as a second opinion over demo routes; its findings go to a
separate ledger with its own semantics (Violation / Needs review /
Recommendation), never into `axe_baseline.json`. Adopt permanently only if it
finds real issues axe missed across two consecutive runs.

**3. k6 + xk6-sse: REJECT for local, KEEP scenario files for CI.** Load
testing SSE is a nightly concern (⫱G4), never pre-push. The scenario files are
committed; the custom k6 binary builds in CI where the Go proxy is reachable.

**4. Bencher: REJECT for now** (benchstat + committed baselines suffice at
current cadence; revisit via the M24 budget table trigger).
**rustywind/stylelint: TRIAL on demand** (both verified healthy; adopt only if
findings justify a formatter/linter lane — M17 artifacts are the input).

## Alternatives considered

- **Full chromedp → Playwright migration**: single toolchain, native
  multi-engine, trace viewer; but throws away a mature ~30-file suite, and the
  pixel-golden + axe + CDP machinery would need porting with zero behavioral
  gain today. Rejected for now; re-decidable at any ⫱G3 without sunk cost —
  the Playwright scaffolding already exists.
- **Two maintained pixel suites** (chromedp goldens + Playwright screenshots
  both gating): double the flake surface and double the baselines for the same
  guarantee. Rejected on principle (plan guardrail 1).
- **Do nothing**: cross-engine and second-opinion claims stay unproven while
  the library sells graceful degradation. Rejected.

## Consequences

- One new flake app (`.#crossbrowser`), one Node tool directory with a pinned
  playwright-core version matching the nixpkgs playwright-driver — the version
  pair moves in lockstep (same-edit rule, like the templ pin).
- The crossbrowser lane runs on-demand/nightly, uploads artifacts on failure,
  carries an owner + kill-criteria row in the M24 budget table before its
  first CI run.
- Findings from any new lane flow into TODO_LIST as numbered items — a finding
  without a registered item is a leak, not a result.
- This ADR is the boundary record: if ⫱G3 ever flips to REPLACE, the migration
  plan gets its own pareto plan and this ADR is superseded, not edited.
