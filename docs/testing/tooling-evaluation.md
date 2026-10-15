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

### M07 — base-select Firefox fallback assert

- Target: `/stylable-select` demo route (M05). Result: see §5.2.

<!-- §5.2 fallback verdict -->

### M08 — IBM accessibility-checker spike + axe diff

- Result: see §5.3.

<!-- §5.3 achecker diff table -->

### M10 — aria-snapshot goldens spike

- Result: see §5.4.

<!-- §5.4 aria snapshot verdict -->

### M12 — SSE load profile (Go client stand-in, see §2)

- Result: see §5.5.

<!-- §5.5 load report -->

## 6. Verdicts (M14) — to be ratified at ⫱G1–G3

See ADR 0048 for the adopt/reject matrix. Nothing here is integrated into the repo's
test surface until the owner ratifies a gate.

## 7. Re-verification reminder (M26 seed)

All "verified" rows above are point-in-time (2026-10-15). Re-check if adoption slips
past ~6 months (nixpkgs playwright version moves fast; Firefox base-select lands via
Mozilla 1944403/1958443).
