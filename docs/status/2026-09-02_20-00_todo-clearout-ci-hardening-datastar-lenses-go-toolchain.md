# Status Report — 2026-09-02 20:00 CEST — TODO-List Clearout, CI Hardening, Datastar Test Lenses, Go 1.26.7

**Session:** single Crush session, 2026-09-02 ~18:00–20:00 CEST
**Base:** master @ `a62c6ca` → HEAD @ `045bbbe` (21 daemon auto-commits, ~75 files changed)
**Verification at close:** `nix run .#verify` all green (build + workspace tests + per-module GOWORK=off tests + lint **0 issues** on all 7 modules) · `nix flake check` green · **89/89 visual goldens pass** · all 4 fast guards + module-layer guard green · TestDocsCountDrift green.

> Format note: user explicitly requested `.md`; the status-report skill's canonical HTML
> output was overridden by that instruction (flagged per skill spec, not propagated).

---

## a) FULLY DONE — implemented AND verified this session

All 19 actionable items from TODO_LIST's "Open" queue were attempted; every item below
is complete with in-session verification (not just "wrote code, assume it works"):

| #          | Item                              | What landed                                                                                                                                                                                                                                                                                                                                      | How verified                                                                                                                                                                                         |
| ---------- | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 112        | visualtest sibling-pin guard      | `check-version-sync.sh` fails if ANY go.mod pins a templ-components sibling ≠ `utils.Version`; release.sh 5b now bumps visualtest pins too; visualtest/go.mod gained local sibling replaces (offline-tidyable, fixes the 9-day-red class structurally); cleanup trap covers visualtest                                                           | Positive test (all guards pass) + negative test (sed-injected v1.10.0 pin → guard exits 1) + all-guard re-run                                                                                        |
| 115        | sibling-pin policy doc            | "visualtest sibling-pin policy" section in `docs/modularization/README.md`                                                                                                                                                                                                                                                                       | doc review                                                                                                                                                                                           |
| 111        | CI fail-fast on go.sum drift      | Premise was **stale** — "Verify no untracked changes" has preceded Test since the initial workflow. Real gap: CI's tidy loop omitted visualtest, so its drift surfaced only in the separate Visual job. Fixed: tidy loop includes visualtest (works offline via the new local replaces)                                                          | git-log archaeology of step order (88f1552 → HEAD)                                                                                                                                                   |
| 113        | visualtest compile flake          | Retry ×3 with backoff after the proxy.golang.org INTERNAL_ERROR abort (run 33399151599)                                                                                                                                                                                                                                                          | YAML parse; logic mirrors ci-repro's (run locally)                                                                                                                                                   |
| 116        | CHANGELOG-diff CI guard           | New `changelog-guard` job: PR touching `*.templ` or library `**/*.go` must diff CHANGELOG.md; uses PR files API (no git history fetch)                                                                                                                                                                                                           | PyYAML parse; regex reviewed against repo layout                                                                                                                                                     |
| 122        | Node-20 action deprecations       | checkout v4→**v7.0.1**, setup-go v5→**v7.0.0**, setup-node v4→**v7.0.0**, upload-artifact v4→**v7.0.1**, download-artifact v4→**v8.0.1** across ci/website/master-red-alert/upstream-watch                                                                                                                                                       | Every version+SHA fetched live from the GitHub API; breaking changes checked (checkout v7 fork-PR blocking is N/A — no pull_request_target here; runner ≥2.327.1 satisfied by ubuntu-latest)         |
| 104        | go-datastar/static upstream watch | New scheduled workflow (weekly + manual): compares `datastar/go.mod` pin vs proxy `@latest` (JSON→jq), opens a single non-duplicating issue with the re-audit checklist                                                                                                                                                                          | Script body executed locally against the real proxy: correctly detects **v0.2.0 pinned vs v0.3.0 published**                                                                                         |
| 97         | SDKScript preconnect lenses       | New `sdk_script_cdn_nonce.golden` (CDN+nonce+preconnect combo) + 3 unit tests: preconnect present (default + custom CDN origin), omitted when `Src` set (no third-party leak), CDN+nonce combination                                                                                                                                             | `go test` + golden regenerated with pinned templ v0.3.1020                                                                                                                                           |
| 100        | LiveRegion URL edge cases         | **Behavior fix:** whitespace-only URL now degrades to a plain container (was emitting `@get('   ')` — same runtime-rejection class as empty); quote-escaping inheritance pinned (`/stream?q='` → `\&#39;` in attribute)                                                                                                                          | New tests pass; templ regenerated via pinned binary; goldens stable                                                                                                                                  |
| 98         | Fuzz writeDatastarPatch           | **Behavior fix:** interior blank lines now PRESERVED via empty-value datalines (verified against pinned bundle's parser bytes: same-key values join with `\n`), so `<pre>` content round-trips; CR/LF in selector/mode replaced with spaces (dataline-injection hardening); `FuzzWriteDatastarPatch` + blank-line regression test                | 2,309,928 fuzz execs, 0 failures; regression test green                                                                                                                                              |
| 105        | Demo polish                       | `testing.Short()` skips the 2s-ticker SSE test (and the 600 ms save test); SSEErrorHandling demo nonce now from `demoBaseProps().Nonce`                                                                                                                                                                                                          | `-short` run + full run                                                                                                                                                                              |
| 94         | SSEErrorHandling BDD + a11y       | 7 behavior-phrased subtests (announcer text on both failure paths, recovery copy, transient retries stay silent, `typeof tcShowToast` guard, error-type toasts, duration wiring) + new `a11y_test.go` (polite announcer, textContent-not-innerHTML, both paths announce, CSP nonce)                                                              | datastar suite green incl. -race                                                                                                                                                                     |
| 101        | aria-busy on LiveRegion           | Auto-started regions render `aria-busy="true"` + `data-tc-live-busy` + singleton-guarded nonced script clearing it on first `datastar-patch-*` or terminal failure; manual/empty-URL regions render nothing; 2 goldens updated (reviewed: only with_url/retry_always changed); unit + BDD + a11y tests                                           | Full datastar suite + race + regenerated goldens reviewed line-by-line                                                                                                                               |
| 96         | SSE heartbeat                     | Demo stream emits `: ping` every 15 s in a second select branch (single goroutine — no writer race); recipe documents `stream.Heartbeat`                                                                                                                                                                                                         | Demo build + endpoint tests                                                                                                                                                                          |
| 99         | Headers-contract audit            | Audit of ALL demo endpoints vs their consumer contracts (LoadMore re-arm cursor chain, ConfirmDelete target, LoadingButton swap=none, PolledRegion tick re-arm, FilterDropdown params) → all `text/html` + new `Cache-Control: no-store` on the 5 fragment endpoints; new `TestHTMXEndpointHeaders` (7 cases) pins headers + core response shape | 7/7 subtests green; full demo suite green                                                                                                                                                            |
| 121        | ci-repro.sh                       | `scripts/ci-repro.sh` mirrors CI's Build & Test step-for-step; `--lint/--css/--visual/--cold` flags; **found+fixed a latent bug while writing it**: CI's own `bc` coverage check silently passes when `bc` is missing — script uses awk                                                                                                          | Ran end-to-end green (all steps, 71.7% coverage); awk comparison unit-checked both directions                                                                                                        |
| 102        | cmd/tc datastar                   | Scaffolder now ships the datastar package (4 post-audit .templ sources) + registry entry + packageDeps (7 sibling files)                                                                                                                                                                                                                         | `tc ls` lists datastar; `tc add live_region --list-deps` resolves 7 deps; cmd/tc tests pass                                                                                                          |
| 118        | release-checklist doc             | `docs/release-checklist.md`: step-by-step with the incident each hardening exists for (v1.8.3/v1.9.0/v1.10.0/v1.11.0), pre-push audit, 24-hour watch, incident index                                                                                                                                                                             | Written from AGENTS.md + release.sh source read in-session                                                                                                                                           |
| 103        | datastar recipe additions         | `NewResponseFromHTTP` one-liner, `stream.Heartbeat` keep-alive section, `datastar-fetch` analytics hooks section (event taxonomy + no-sse-error warning + retry matrix), Further-Reading cross-link to runtime facts                                                                                                                             | **Every documented API verified against module-cache source** (go-datastar response.go:222, go-sse stream.go:184); one invented API (`sse.CloseFromHTTP`) caught and removed before it could mislead |
| 106        | Go toolchain bump                 | **1.26.5 → 1.26.7** everywhere: dedicated `nixpkgs-go` flake input (rolling unstable; preserves the locked nixpkgs' templ v0.3.1020 + golangci pins — same isolation pattern as nixpkgs-chromium), devShell + build/test/coverage/visual apps use `goToolchain`, `go 1.26.7` in all 8 go.mod + go.work                                           | `nix run .#build` + `.#test` green with `go version go1.26.7`; zero templ-regen cosmetic diff (pin held)                                                                                             |
| —          | golangci-lint drift (discovered)  | Full verify exposed: locked nixpkgs' golangci-lint 2.12.2 rejects `.golangci.yml`'s `exhaustruct_v5` (CI pins 2.13.2). `nix run .#lint`/devShell now use 2.13.1 from nixpkgs-go                                                                                                                                                                  | `nix run .#lint` → 0 issues across all 7 modules                                                                                                                                                     |
| —          | treefmt sandbox fix (discovered)  | goimports (gotools) shells out to `go`; the 1.26.5 `go` tried to **download** the 1.26.7 toolchain inside the network-less flake-check sandbox → now uses gotools from nixpkgs-go                                                                                                                                                                | `nix flake check` green                                                                                                                                                                              |
| 114        | pnpm lockfile audit               | `pnpm install --frozen-lockfile` passes in website/ → **no manifest/lockfile splits remain** post-423ea1b                                                                                                                                                                                                                                        | Real frozen install (1 s, "Already up to date")                                                                                                                                                      |
| 119        | local pnpm root cause             | **Root-caused:** `~/.local/bin/node` is a bun 1.3.13 shim lacking `node:sqlite`; pnpm 11 subprocesses resolve `node` via PATH. Workaround verified: prepend real nix nodejs-slim bin → frozen install passes                                                                                                                                     | Reproduced failure, then reproduced success with PATH fix                                                                                                                                            |
| 117        | CodeRabbit                        | **Resolved by investigation:** PR #5 is merged and CodeRabbit DID review (2026-08-31, 1 golines nit, addressed pre-merge) — nothing to re-request                                                                                                                                                                                                | `gh pr view 5`                                                                                                                                                                                       |
| 95/109/110 | visual goldens                    | **16 new PNGs**: statcard yellow/purple (light+dark), outline success/info (light+dark), datastar LiveRegion (filled, light+dark) + Indicator, Eyebrow light/dark, Scrollback light/dark, Eyebrow+PageHeader+Scrollback composition. datastar added as visualtest module dep (+replace)                                                          | Scoped `nix run .#visual -- -run … -update`, then FULL suite: **89/89 OK**                                                                                                                           |

**Docs hygiene:** CHANGELOG `[Unreleased]` warmed with 12 detailed entries (prior entries preserved); TODO_LIST.md rewritten to post-session truth (Open queue now empty, 120/119-notes in Deferred, #80 note extended); AGENTS.md updated (ci-repro entry point, toolchain-input-split gotcha, pnpm shim gotcha, release-checklist pointer). TestDocsCountDrift forced README/ROADMAP golden counts 72→89 — updated, guard green.

---

## b) PARTIALLY DONE — real work shipped, honest caveats

1. **#106 Go bump — toolchain-reproducibility tradeoff.** golangci-lint now comes from the rolling `nixpkgs-go` input (2.13.1) instead of the locked input; CI pins 2.13.2. Config-compatible (verified: lint runs clean) but not bit-identical to CI. The input is documented for fold-back at the next deliberate full-flake update.
2. **Consumer toolchain floor.** `go 1.26.7` in go.mod forces consumer toolchains ≥1.26.7 (auto-download handles older ones). Standard for a security-driven bump, but it IS a floor raise on a published library.
3. **#95 datastar visual coverage — SDKScript deliberately excluded** (renders only `<script>`/`<link>`: no pixels to pin). LiveRegion's visual is the static filled state — actual SSE patching is not pixel-testable by design.
4. **#120 CSS-recompile false negative — root cause NOT identified.** `nix run .#css` is byte-stable today (SHA-256 verified twice); the 2026-08-31 local-vs-CI discrepancy is unreproducible post-`fabd1fb`. Mitigation shipped (`ci-repro.sh --css`), mystery itself moved to Deferred.
5. **#119 pnpm — workaround only.** The bun shim stays until the user removes it; documented in TODO_LIST + AGENTS.md.
6. **New CI workflows are unproven in real CI.** changelog-guard + upstream-watch validated by parse + local logic runs only. Risks: `gh issue list --search` dedup semantics, `--paginate`+`--jq` combination, PR-file list pagination on huge PRs.
7. **New agent-generated PNGs (16) extend blocked #80** — a human must eyeball them (datastar, eyebrow, scrollback are brand-new component surfaces).
8. **CHANGELOG-guard policy choice embedded in regex** — test-only PRs (`_test.go`) currently REQUIRE a changelog diff. Defensible (house rule says every feature/fix), but it's a judgment call the owner hasn't ratified.

---

## c) NOT STARTED — untouched this session (blocked or out of scope)

- **go-datastar/static v0.3.0 bump + wire-format re-audit** — now surfaced by the watch workflow; needs human-audited bundle diff against `docs/datastar-runtime-facts.md`.
- **#80** human eyeball of overlay PNGs (now 16 more agent PNGs on top).
- **#28 / #29** awesome-templ + templ.guide upstream submissions (maintainer approval).
- **#93 / #107 / #108 / #124** BuildFlow daemon family (separate repo; this session produced 21 more hallucinated auto-commit messages — the problem is live).
- **#123** branch protection + required checks (repo-owner GitHub settings).
- **#39** compound overlay components (v2.0/ADR-0023), **#33** `Validate()` methods, **#34** internal/testutil migration — deferred by prior decisions.
- **v1.12.0 release cut** — this session's body of work is unreleased (sits in `[Unreleased]`).

---

## d) TOTALLY FUCKED UP — actual mistakes made this session (all caught, none shipped broken)

1. **Almost deleted shipped CHANGELOG entries.** My first `[Unreleased]` edit REPLACED the prior session's StatTone/CI-recovery entries instead of appending. Caught within one read-back; restored verbatim. Had it slipped through, release notes for landed work would have been silently lost.
2. **Invented an API in consumer-facing docs.** Wrote `sse.CloseFromHTTP` into the datastar recipe — it does not exist in go-sse. Caught by source-checking the doc claims; also had to soften an auto-close claim about `NewStream` that the source does not support. Lesson applied: every recipe API this session was then verified against module-cache source.
3. **First ci-repro.sh version had a silent-failure bug** — the coverage threshold used `bc`, which is absent on this machine, and the check would pass WITHOUT checking. Found only because I actually ran the script. Ironic for a tool whose purpose is trustworthiness; fixed with awk + a comment warning.
4. **Wrong-file edit + templ-syntax-in-Go + wrong-struct-field:** three distinct multi-edit/compile mistakes in one hour (golden-sweep edit aimed at the wrong file; used templ block syntax `{ }` for LiveRegion children inside a Go test; wrote `EyebrowProps.Class` instead of the embedded BaseProps field). All compile-gated and fixed; cost maybe 6 round-trips.
5. **Wrote a nondeterministic visual test despite a known precedent.** TestSpinner already documents that CSS-animated components need raised `MaxMismatch` — I screenshot a rotating spinner at 0.1% anyway; CI-equivalent runs failed 9 px / 0.39%. Fixed with the spinner's own 0.08 pattern. Should have read the precedent first.
6. **The fuzzer earned its keep against its own author:** my round-trip assertion stripped only one trailing `\n`, misreading the event terminator — fuzz crashed my test on seed 2. Fixed (`TrimSuffix "\n\n"`); the final harness then survived 2.3M execs.
7. **~21 daemon auto-commits with hallucinated messages** now sit between `a62c6ca` and HEAD (#93's disease, live). No data lost, but session history is unreadable, and per AGENTS.md the daemon's broad `git add -A` is exactly how stale files get enshrined — nobody has audited these 21 commits diff-by-diff.
8. **Environment friction I failed to pre-empt:** `go install` (for actionlint) is blocked by local security policy; instead of a fallback (nix run) I only PyYAML-parsed the workflows — real GitHub-Actions lint never ran (see b.6).

---

## e) WHAT WE SHOULD IMPROVE — process-level takeaways

1. **Run the thing, don't just write it.** Every bug I caught (bc, fuzz harness, golines, exhaustruct_v5, toolchain-in-sandbox) was caught by EXECUTING; every miss (CloseFromHTTP, wrong-file edit) was caught by re-reading. Both loops are mandatory; neither substitutes for the other.
2. **Precedent-scan before writing tests.** The Spinner MaxMismatch pattern existed; five minutes of `grep` would have saved a failed visual run. For new component tests: grep the target package for similar animated/async cases FIRST.
3. **Verify-external-claims applies to docs I write too.** Recipe/API doc claims need module-cache source checks in the same breath as the writing, not after.
4. **Daemon awareness.** With auto-commits landing mid-session, `git status` is not ground truth — `git log a62c6ca..HEAD` + diff review is. All dirty-tree assumptions should be daemon-proofed (and #93 remains the root fix).
5. **Guard scripts should fail loud on missing dependencies.** bc-lesson generalized: any `check-*` script that shells to an optional binary must hard-fail (or awk/POSIX-fallback) when that binary is absent.
6. **The stale-TODO problem is real.** #111's premise had been false since the initial workflow commit — TODO_LIST items should carry a "verified against HEAD on <date>" stamp, or HARVEST should re-verify premises before re-listing them.

---

## f) NEXT — up to 50 candidate items (brainstorm; most belong in ROADMAP, not TODO_LIST)

**Ship & verify (highest impact)**

1. Human eyeball of all 16 new PNGs + the 5 overlay sets (blocked #80) — 10 minutes, closes #80.
2. Audit the 21 daemon commits diff-by-diff (`git diff a62c6ca..HEAD --stat` per commit) before anything is pushed; enshrined-stale-file check per AGENTS.md.
3. Cut **v1.12.0** via release.sh (release-checklist pre-push audit included) — ships StatTone-era leftovers + this session.
4. Post-push 24-hour watch per checklist: CI/Website green, per-module tidy clean, proxy + pkg.go.dev resolution.
5. **Bump go-datastar/static v0.2.0 → v0.3.0** with the mandated bundle re-audit (event names, CSP, retry matrix vs runtime facts).
6. Real-CI shakedown of the two new workflows (changelog-guard on a test PR; upstream-watch via workflow_dispatch before Monday's cron).
7. Confirm changelog-guard behaves on a docs-only PR (should pass) and a components PR without changelog (should fail) — two throwaway PRs.
8. Add `workflow_dispatch` inputs / dry-run mode to upstream-watch for testing without issue creation.

**Follow-ups on this session's work**
9. golines: decide a max-width policy and add golines locally (or CI-side autofix) instead of hand-fixing violations.
10. golangci-lint: fold 2.13.2 exact-parity into the flake (go-install pin or nixpkgs bump) so local lint is bit-identical to CI.
11. actionlint: add to `nix run .#lint` (it's CI-only today) and run it locally via nixpkgs — closes the untested-YAML gap properly.
12. ci-repro.sh: add a `--actionlint` step + wire the four guards into the default run (currently only under --lint).
13. upstream-watch: also watch `templ` releases + `golangci-lint` releases (same watcher pattern, three jobs).
14. changelog-guard: exempt `_test.go`-only diffs if the owner wants test-only PRs changelog-free (needs decision, see questions).
15. version-sync guard: also assert `go.work`'s `go` directive equals go.mod's (currently unguarded after the 1.26.7 bump).
16. check-module-layers.sh: extend to model visualtest's new datastar dependency explicitly (it passes today; make the DAG include visualtest rather than skip it).
17. Extend `TestHTMXEndpointHeaders` pattern into a shared table the release-checklist references (demo as executable contract doc).
18. Give Indicator a `ReduceMotion`-aware static screenshot mode (or document the MaxMismatch recipe) so future animated components have a canonical pattern.
19. Pre-warm the go1.26.7 toolchain into anything sandboxed that shells to `go` (treefmt lesson — audit OTHER sandboxed checks for toolchain-download traps).
20. Fold nixpkgs-go + nixpkgs into one input at the next deliberate full-flake update (documented TODO in flake.nix).

**Open-source / upstream**
21. File the upstream go-sse feature request if `CloseFromHTTP`-style helper is genuinely wanted (verify-before-filing first).
22. CodeRabbit: enable explicit review request on next PR to confirm rate-limit recovery (verifies the 2026-08-31 hypothesis).
23. Check whether go-datastar v0.3.0 changed the SSE `: comment` heartbeat contract we now rely on in the demo.
24. awesome-templ + templ.guide submissions (#28/#29) — draft + submit.

**Library substance (ROADMAP fuel)**
25. `charts/` visual goldens (LineChart/PieChart/DonutChart/AreaChart have HTML goldens; pixel coverage for dark-mode strokes is thin).
26. Interaction-state (hover/focus) goldens for Dropdown/Popover/ContextMenu top-layer positioning — complements #80's static check.
27. RTL visual goldens for the new components (Eyebrow/Scrollback are direction-agnostic, but PolledRegion/DataTable have RTL variants untested in pixels).
28. `datastar.SSEErrorHandling` — a real chromedp test that dispatches a synthetic `datastar-fetch` event and asserts announcer textContent + toast DOM (the JS paths are only string-pinned today).
29. Same for `aria-busy` clearer script: dispatch `datastar-patch-elements` in chromedp and assert `removeAttribute` ran.
30. LiveRegion: consider `aria-busy` interplay with `RetryAlways` reconnect cycles (busy re-assert on retry?) — design question.
31. PolledRegion: give htmx's counterpart the same `aria-busy` cue for parity.
32. Fuzz `getActionExpr`/`actionExpr` (URL + retry + cancellation combinations) like writeDatastarPatch.
33. cmd/tc: `tc add` for datastar components copies `.templ` but datastar's `.go` siblings import `templ`/`utils` — generate an import-checklist in --list-deps output.
34. cmd/tc tests: add a round-trip test asserting the embedded datastar sources == the real package sources (drift guard for `_sources/`).
35. `_sources/` drift guard generally: sources were copied by hand; a checksum test would catch package-vs-scaffolder drift.
36. Demo: `/api/save` returns "Saved." that nothing displays (hx-swap="none") — either surface a toast or drop the response body.
37. Demo: heartbeat visibility — surface "last ping" text so the keep-alive is demonstrable in the UI.
38. website: exercise `pnpm build` locally with the PATH workaround (only install/frozen-lockfile was validated).
39. Coverage: 71.7% is barely over the 70% gate — the new code paths (busy script, fuzz writer) probably moved it; consider raising the gate or adding missing-coverage tests.
40. Make `TestCSSFreshness` fail locally too (it warns off-CI today) behind a flag, complementing `ci-repro --css`.
41. AGENTS.md is ~40 KB of accreted context — consider splitting "release/release-script" sections to point at the new checklist doc only (started with #118).
42. Introduce `docs/DOMAIN_LANGUAGE.md` entries for the new terms (busy-cue, sibling-pin policy, keep-alive frame).
43. Dependency watch beyond go-datastar: templ v0.3.1036+ release still pending — the moment it publishes, bump go.mod + templ binary in lockstep (AGENTS.md rule).
44. Dependabot/Renovate decision for actions SHAs (bumps are manual today; #122 will recur).
45. Consider `permissions: contents: read` hardening pass over all workflows (upstream-watch has minimal perms; ci/website were not touched).
46. visualtest: `-update` currently rewrites scoped runs fine, but add a guard against accidentally committing `.fail/` artifacts (gitignore + test).
47. Add `scripts/ci-repro.sh --visual` usage note to docs/visual-testing.md.
48. CHANGELOG: consider splitting "Changed" for CI-only vs consumer-visible changes (reviewer ergonomics at release time).
49. Sprint the #34 internal/testutil migration with codemods — it blocks every package-move refactor (deferred item with compounding cost).
50. Schedule the deliberate full `nix flake update` (nixpkgs + chromium + templ re-pin) as its own session with the visual-golden regen procedure from the checklist.

---

## g) Questions I cannot answer myself

1. **go-datastar/static v0.3.0:** do you want me to bump + re-audit the new bundle **now** (it's the only pending functional change and the watch workflow's first real case), or leave it as the watch-issue exercise for a later session?
2. **Daemon history:** the 21 auto-commits between `a62c6ca` and HEAD are unpushed. Should I squash/rewrite them into logical commits before any push (needs your explicit approval since it rewrites history), or leave them as-is per the "daemon is expected" doctrine?
3. **CHANGELOG-guard policy:** should test-only PRs (`_test.go` diffs) keep requiring a CHANGELOG entry (current behavior), or be relaxed to non-test source only?

---

_Awaiting instructions._
