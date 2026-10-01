# Status Report — Demo MPA Continuation: CSP Nonce Regression Found & Fixed, Axe Contrast Pass, Live Verification

**Session date:** 2026-10-01, 12:35 CEST (continuation session; started ~08:00 CEST)
**Scope:** Verify the demo-MPA + Firebase Hosting→Cloud Run work shipped by the previous session (see `2026-10-01_07-08`), fix everything the strict gates surfaced, and complete the verification pipeline toward push.
**Repo:** templ-components @ master, 9+ local commits ahead of origin (auto-commit daemon active; it also pushed the MPA earlier — CI deployed it at 05:43Z before this session's fixes existed).
**Predecessor reports:** `docs/status/2026-10-01_07-08_demo-mpa-firebase-cloudrun-status.md` (+ its 08:20 addendum) — this report supersedes its "remaining" list.

---

## Executive summary

The MPA demo **is live** at `https://templcomponents.lars.software/demo/**` (deployed by CI at 05:43Z from the daemon's push). Live probing verified the canonical surface (routes, titles, CSP header through the rewrite, POST fragments, dual-mount, CSS caching) **and found one platform limitation**: SSE cannot stream through a Firebase Hosting rewrite (buffered forever; raw run.app streams fine) — documented, owner decision tracked as TODO #329.

The heavy gates then surfaced a **pipeline of real defects, all fixed forward**:

1. **11 axe color-contrast failures across 6 demo routes** (kanban empty placeholder 2.41:1, scrollback timestamps, errorpage cause list on tinted cards, `InlineSuccess` + outline Warning/Success button text at ~3.2:1, CountBadge white-on-red-500, Nav footer bottom bar, demo filter-results text) — all fixed in the library/demo with `-700`/`-600` light shades.
2. **A harness bug**: the axe sweep audited mid-entrance frames (staggered animations composite translucent text → bogus contrast numbers). It now runs `waitAnimationsSettled()` like the screenshot path.
3. **A touch-target failure** (demo shell footer Documentation/GitHub links 20px tall) — hit-box pattern applied.
4. **A real carousel contrast failure** (white 18px-bold on emerald-600 = 3.65:1) → emerald-700.
5. **`TestDemoKanbanHTTPContracts` still harvested CSRF from `/`** after the MPA split → `pagePath` fixture field, demo targets `/kanban`.
6. **THE BIG ONE — a live CSP regression shipped by the MPA rewrite:** every demo page served **component scripts with `nonce=""`** (CopyButton on every page, Kanban move JS, Tooltip/Popover/ContextMenu/menu-nav JS, Image fallback, DirtyGuard, GlobalErrorHandling, dismissible Alert/Toast scripts, SimpleNav/Nav mobile-menu JS). The demo CSP header requires `'nonce-demo-nonce'`, so **all of that client-side JS was silently dead on the live site since 05:43Z**. The old single-page demo passed `demoNonceConst` into every section builder; the MPA rewrite dropped it at 13 call sites. Root cause found via a kanban e2e failure (`window.tcKanbanAttached === undefined`), confirmed against the LIVE HTML (`<script nonce="">`), fixed at every site, and verified by a temporary inventory probe: **0 empty-nonce scripts across all 16 demo pages** (was 50).

Stat line: **14 defect classes fixed · 1 platform limitation documented (owner decision) · 1 new AGENTS.md gotcha · verification pipeline in progress.**

---

## a) FULLY DONE

1. **Skill load + state recovery** — resumed from the 07:08 report; corrected the stale todo statuses; loaded the templ-components skill.
2. **Working-tree/daemon triage** — confirmed the daemon's `styles.css` "clobber" is the mirror guard's intended state (byte-identical to `templates/styles.css`; `TestCompiledCSSInventory` + tc-sources guard pass); identified that CI had already deployed the MPA at 05:43Z.
3. **Job 08E triage** — the previous session's unverified background suite: goldens regenerated fine; axe failed 6 routes with 11 contrast targets; all enumerated and mapped to source.
4. **Library contrast fixes (with golden + test cascade)** — `display/kanban.templ` (empty placeholder → gray-600/gray-400), `display/scrollback.templ` (timestamps → gray-500/gray-400; TAG tones warning/success → amber-700/green-700 — the sweep caught those on the second pass), `display/count_badge.templ` (`bg-red-600 dark:bg-red-500`), `display/button_go.go` (outline Warning/Success text → amber-700/green-700 light), `errorpage/shared.templ` (causeList label/rows/code spans → gray-600/gray-400 on tinted cards), `feedback/alert.templ` (InlineSuccess → green-700), `navigation/nav.templ` (Footer BottomBar wrappers → gray-500/gray-400). HTML goldens regenerated (`-update` AFTER the package — see gotcha), string-assertion tests updated (`feedback/snapshot_test.go`, `display/scrollback_test.go`, `display/bdd_new_test.go` sentinel).
5. **Demo contrast + touch-target fixes** — forms filter-results text, footer Documentation/GitHub links (`p-1.5 -m-1.5` hit-box, 24×24, layout-neutral), carousel emerald slide → `bg-emerald-700`.
6. **Axe harness fix** — `visualtest/axe_sweep_test.go` runs `waitAnimationsSettled()` after the theme pin; mid-entrance audit frames eliminated (root cause of flapping bogus contrast findings).
7. **Kanban contract test MPA retarget** — `kanbanContractFixture.pagePath` (demo `/kanban`, parity twin `/`); the three verification fetches in the shared probe driver now use `fx.pagePath` instead of hardcoded `/`.
8. **THE CSP NONCE REGRESSION — root-caused and fixed (13 call sites)** — `shell.templ` sidebar CopyButton, `demo.templ` hero + `demoCodeSnippet` CopyButtons, `display_demo.templ` Dropdown/Tooltip ×2/CopyButton/Image ×2, `kanban_demo.go` board props (covers both boards), `errorpage_demo.templ` ErrorAlert, `navigation_demo.templ` Nav + SimpleNav, `wire_demo.templ` DirtyGuard, `htmx_demo.templ` GlobalErrorHandling (helper func wrapping `DefaultErrorHandlingConfig()` + Nonce), `feedback_demo.templ` Toasts ×4. Promoted-field gotcha respected (nested `BaseProps`).
9. **Nonce inventory probe** — temporary test loops all 16 demo pages and asserts every `<script nonce="">` is gone: **0 remaining** (was 50 occurrences / 10 distinct scripts). This temp test is a candidate to KEEP as a permanent demo CSP guard (see e/f).
10. **Live verification of the deployed MPA (read-only probes)** — `/demo/`, `/demo/display`, bare `/demo` → 200 + shell + correct titles; CSP header carries `nonce-demo-nonce` through the rewrite; `/demo/css/app.css` 200 + immutable cache; POST `/demo/api/wire/form` → 200 + wire fragment; raw run.app dual-mount 200; ECharts/Datastar pages serve SDK references.
11. **SSE limitation verified at source** — raw run.app streams `datastar-patch-elements` at ~2.2s; the same endpoint via the Hosting rewrite never delivers bytes (25s+). Research confirmed Firebase Hosting buffers ALL rewrite responses (Firebase team statement + production reports); `X-Accel-Buffering: no` (already emitted by the demo) is ignored. Documented in AGENTS.md + status report addendum + TODO #329 with three bypass options.
12. **Docs** — AGENTS.md: new go1.26 `-update` flag-order gotcha + SSE limitation note in Demo Infrastructure; CHANGELOG `[Unreleased]`: contrast-fix entry (incl. hit-box, carousel, sweep settle) + smoke-probe addition; TODO_LIST: #329 (SSE bypass ⫱) + #330 (latent tone-class contrast sweep) + next-ID bump; CI `website.yml` post-deploy smoke now probes `/demo`, `/demo/`, `/demo/display` (status + shell marker + CSP nonce survival).
13. **Hygiene gates green mid-flight** — `nix run .#lint` 0 issues across all 7 modules (+actionlint), `nix fmt` no-op, `nix flake check` pass, workspace build green, demo package tests green.

## b) PARTIALLY DONE

1. **Visual suite green run** — after each fix round the suite progressed (axe now passes; display goldens regenerated). The final FULL clean witness run has not completed yet — it is the next gate, together with the kanban e2e pair that motivated the nonce hunt (fix applied; e2e rerun in progress as this report is written).
2. **Route golden regeneration** — display-page goldens need one more `-update` pass after the nonce/`-700` changes (the shell hit-box is layout-neutral; emerald-700 and script-tag changes alter bytes).
3. **Per-module test loop** — root/demo/display/feedback/navigation/errorpage verified green individually; the complete loop (utils, icons, charts/echarts, datastar, htmx, visualtest compile, website) plus `ci-repro.sh --lint --website` still pending.
4. **Push + redeploy + live re-verify** — everything is LOCAL; origin is 9+ commits behind the local tip. The live site currently runs the MPA **with the broken nonces + failing axe targets** — the fix batch must deploy to heal it.
5. **Temporary diagnostics cleanup** — `visualtest/tmp_kanban_diag_test.go` + `visualtest/tmp_nonce_inventory_test.go` are in the tree (they compile; the daemon may have committed them). Decision pending: delete, or promote the inventory test to a permanent `TestDemoCSPNonceIntegrity` guard (strongly recommended — it is the guard that would have caught this regression before deploy).

## c) NOT STARTED

1. TODO_LIST/ROADMAP full harvest of the 07-08 report's §f (docs-health pass) — items #329/#330 landed; the remaining ~48 items from that report are still unharvested.
2. The 3 open questions from the 07-08 report (SSE bypass strategy, run.app canonical redirect, site-header Demo link) — user input pending; SSE bypass now also gated on TODO #329.
3. Website `--update-csp` check + `build.sh` + html-validate local lane (covered by ci-repro `--website`, pending).
4. Release cut (v1.19.5) for the library-facing contrast fixes — the changed library components (kanban/scrollback/count_badge/errorpage/feedback/navigation) are consumer-visible; CHANGELOG is warm; release awaits the pipeline + owner.
5. Pre-existing deferred items (stale OG card, bun shim, PR wall-clock budget, etc.) — untouched.

## d) TOTALLY FUCKED UP (all repaired; honest ledger)

1. **I declared the live MPA "verified" while its client-side JS was dead site-wide.** My first live-probe script checked status codes, titles, headers, and a server-side POST — every check green — but never asserted a single piece of CLIENT behavior. The nonce regression (shipped at 05:43Z) was invisible to all of it. The kanban e2e was the only signal, and I initially spent effort treating it as a possible load-flake. **Lesson: "live verification" must include at least one JS-execution assertion per page class (e.g. `window.<singleton>Attached === true`), not just HTTP shapes.** The new inventory guard exists precisely because of this.
2. **Chrome's nonce-masking quirk cost a diagnostic detour.** `script.getAttribute('nonce')` returns null in the page context, so my first in-page probe reported ALL scripts nonce-less — the opposite extreme of the truth (locally they had nonces; the LIVE HTML had the `nonce=""` subset). I nearly "fixed" a non-existent local problem; cross-checking the LIVE HTML source set it straight. **Lesson: verify DOM readings against raw HTML before acting on them.**
3. **`go test -update ./pkg/...` silently targeted the wrong package** — go1.26.7's own `-update` flag consumes the following argument, so the golden update did nothing while printing plausible output (`FAIL . [setup failed]`, easy to misread as "flag not supported"). Wasted several cycles; the correct form is `go test ./pkg/... -update`. Documented in AGENTS.md; the 07-08-era documented form was wrong.
4. **A pipe swallowed a build failure and I printed a false "BUILD-OK".** `go build ./... 2>&1 | head -5 && echo BUILD-OK` reports head's exit status, not the build's — the struct-literal compile errors scrolled past and my own echo lied. Exactly the repo's documented "never trust a piped capture without its final line" trap, self-inflicted. Fixed with `if go build ...; then` + a stderr file.
5. **Struct-literal promoted-field mistakes, twice.** I wrote `Nonce:` directly in literals for `ErrorAlertProps`, `ToastProps`, `NavProps`, `SimpleNavProps`, `DirtyGuardProps` — all embed `utils.BaseProps`, and embedded fields are not promotable into literals (the repo's own documented gotcha). The compiler caught every one; all fixed to nested `BaseProps: utils.BaseProps{Nonce: …}`.
6. **Two syntax errors in my own throwaway diagnostic** (an `if`-init without a condition, twice) — caught by go vet/goformat; repaired. Sloppy for a file that exists only to find other people's bugs.
7. **Axe-scope misjudgment, corrected by iteration.** I first treated the sweep findings as "demo content + one library bug"; the second strict pass proved the scrollback TAG tones and the carousel slide were ALSO failing (real library-visible shades). The iteration was cheap but the first pass should have grepped ALL rendered tone classes up front (TODO #330 now tracks the demo-unrendered tail).
8. **Background-job output truncation bit me twice** — job 08E's 30k-char tail hid which routes/tests failed; a filtered re-run of the axe sweep through plain `go test` silently SKIPPED (no nix browser env — the documented sub-second-"ok"-is-a-skip trap). Both recovered by re-running through `nix run .#visual` with tight grep patterns.

## e) WHAT WE SHOULD IMPROVE

1. **Promote the nonce inventory to a permanent guard** (`TestDemoCSPNonceIntegrity`): every demo page must render zero `<script nonce="">` — plus, ideally, assert one singleton flag (`window.tcKanbanAttached`, `window.tcMobileMenuAttached`, …) is set after load per page class. This is the single highest-leverage guard this session produced: it would have blocked the live regression at PR time.
2. **Demo-level CSP integration test parity**: the library has `integration/csp_nonce_test.go` (asserts `nonce=` on every script) — but the DEMO pages assemble components from many packages; nothing asserted nonce NON-EMPTINESS end-to-end. A demo-package twin of that test closes the class.
3. **Live-probe scripts should exercise behavior, not shapes**: add a JS-execution probe (evaluate a singleton flag) and an SSE canary with a short timeout to any future live-verification pass; record both-limitation results in the report.
4. **The axe baseline ledger still carries `-1` waivers for index/index_dark/forms_dark** — now that pages are split and strictly audited, re-triage those entries: the index waiver should shrink to a documented budget or vanish (the policy prefers budgets; entries that no longer match findings must be deleted).
5. **The `-update` flag-order trap is now documented, but a guard would be better**: `utils/golden` could detect "update requested but zero golden files rewritten while tests passed" and warn — cheap canary for the next Go flag semantic change.
6. **Consider a `templ.GetNonce(ctx)` fallback in `scriptComponent`** (library-side): components rendered under a nonce-carrying context would inherit it without explicit props. Needs a decision (it changes library behavior for empty-nonce callers) — candidate ADR/TODO, not a silent change.
7. **The 07-08 report's §f list (50 items) predates this session's discoveries** — harvest it through docs-health and merge with this report's §f to avoid two competing backlogs.
8. **e2e for shell pages should assert interactivity on EVERY page class** (mobile menu open/close at minimum) — the MPA made the shell the single point of failure for 14 pages; only kanban e2e happened to cover it.

## f) TOP 50 THINGS TO GET DONE NEXT

_Near-term pipeline first (this session's thread), then hardening, then the inherited backlog. IDs: continue TODO_LIST numbering (next free 331)._

**Finish this session's thread (1–8)**
1. Rerun the kanban e2e pair + contract test on the nonce-fixed tree (in progress as this is written) — expect green.
2. Full visual suite `-update` run (regenerate display/route goldens for `-700` + nonce byte changes) — expect axe + all e2e green.
3. Full visual suite CLEAN witness run (no `-update`) — the ritual gate.
4. Eyeball the ~30 route goldens (shell quality check owed from the 07-08 report) + the changed component goldens (scrollback tones, countbadge, buttons).
5. Per-module test loop (utils, icons, errorpage, charts/echarts, datastar, htmx, visualtest compile, website) + website `build.sh`/html-validate lane.
6. `scripts/ci-repro.sh --lint --website` at the exact tip; push immediately on PASS; watch the Website workflow deploy.
7. Live re-verification with the BEHAVIORAL probe set (status codes + CSP + POST + **JS singleton flags** + SSE canary) — confirm the nonce fix heals the live site.
8. Decide + execute the temporary-diagnostics cleanup: promote `tmp_nonce_inventory_test.go` to a permanent guard (item e1), delete `tmp_kanban_diag_test.go`.

**Regression hardening (9–16)**
9. Permanent demo CSP guard (see e1) — page sweep × script inventory × singleton flags, wired into the visual suite.
10. Add the JS-behavior assertion to `website.yml`'s post-deploy smoke (currently HTTP-only even after my probe addition).
11. Add an SSE canary (5s budget) to the post-deploy smoke with an EXPECTED-FAILURE annotation until #329 lands, so the limitation is visible in CI rather than forgotten.
12. Re-triage the axe baseline `-1` entries (e4): split index/index_dark/forms_dark into real budgets or delete matched-no-longer entries; document in `docs/testing/a11y-gate-policy.md`.
13. TODO #330 execution: sweep demo-unrendered sub-4.5:1 TEXT tone classes (ScrollbackToneSuccess-style patterns) with golden cascade.
14. Guard for `-update` flag-order silent no-op (e5) in `utils/golden`.
15. Prerender↔live-render nonce parity: extend `TestPrerenderMatchesLiveServer` to diff script-tag nonce attributes, so static export can never ship nonce-less scripts again.
16. Decide the `templ.GetNonce(ctx)` fallback in `scriptComponent` (e6): ADR or explicit NO with rationale.

**Canonical URL & hosting (17–22)**
17. Owner decision + implementation: SSE bypass (TODO #329 — direct Cloud Run domain / App Hosting / raw-run.app EventSource + CORS + CSP `connect-src`).
18. Owner decision: 301 run.app → canonical `/demo/` or keep dual-serve (07-08 question 2).
19. Owner decision: persistent "Demo" link in the site header (07-08 question 3).
20. Demo sitemap/robots story under the site domain (canonical tags or `Disallow: /demo`).
21. Guard test pinning `demoCSP` (Go) byte-equal to the firebase.json `/demo/**` header (07-08 item c5).
22. Consider `X-Robots-Tag`/cache headers audit for the `/demo/**` rewrite path.

**Inherited from 07-08 §f (23–32, highest-value subset)**
23. Harvest the 07-08 §f items into TODO_LIST via docs-health (single backlog).
24. Release v1.19.5 with the contrast fixes once the pipeline is green (`nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh`).
25. Fix the stale `website/public/og/home.png` OG card via ogshot (pre-existing).
26. Errorpage demo presentation decision (TODO #236 ⫱ — bordered boxes vs standalone routes).
27. `TestPrerenderMatchesLiveServer` determinism (TODO #250).
28. Demo e2e through the Firebase proxy in CI (post-deploy, live URL) — extends item 10.
29. Trim the axe `-1` ledger (same as 12 — kept here for the count only if 12 moves; merge on harvest).
30. Wire the `/demo/**` probe + JS probe into `ci-repro.sh --website` so the smoke lane is reproducible locally.
31. Cold-start check on the demo (min-instances=0): measure first-hit latency after deploy; consider min-instances=1 vs cost (owner).
32. Re-run `scripts/vision-review-goldens.sh` over the new route goldens (human-eyeball assist, TODO #80 family).

**Testing & quality backlog (33–42)**
33. BDD-style e2e for the mobile menu on shell pages (e8) — one flow, both breakpoints.
34. Extend `TestTouchTargetAudit` coverage notes: it caught the footer links only because the audit ran on new pages — sweep ALL pages' footers/headers deliberately.
35. Coverage push toward the 80% awesome-templ gate (TODO #28 blocker) — measure per-package delta after this batch.
36. Keep the axe sweep's new `waitAnimationsSettled()` under observation for runtime cost (21 routes × settle budget) — trim if the suite grows past CI budgets.
37. Consider replacing the two throwaway diag patterns with a reusable `visualtest/probe` helper (script inventory + singleton flags + SSE canary) — the pieces of items 9/10/11 in one tool.
38. Document the Chrome nonce-masking quirk in `docs/javascript-guide.md` (it will burn the next person debugging CSP).
39. Add `go test ./pkg/... -update` correction to any remaining docs that show the old flag order (grep docs/ + recipes).
40. Golden-file count drift guard update if the promoted inventory guard adds testdata.
41. Review whether other Go 1.26 flag-semantics changes affect repo scripts (`-count`, `-run` combos in scripts/*.sh).
42. Feed the "piped-capture false BUILD-OK" lesson into `scripts/ci-repro.sh`-style patterns (it already uses `set -euo` — verify no `| head` gates remain).

**Docs & comms (43–46)**
43. Update `docs/status/2026-10-01_07-08` addendum's "Remaining" list to point at THIS report (avoid two live backlogs).
44. AGENTS.md: record the nonce-passing convention for demo content builders ("every nonce-consuming component gets `demoNonceConst` — enforced by the inventory guard").
45. CHANGELOG: fold the nonce-regression fix into the MPA entry as an explicit "Fixed" line (it is consumer-affecting for anyone who vendored the MPA demo — actually demo-only, but the demo is a sales surface).
46. Note in `docs/testing/a11y-gate-policy.md`: the sweep now runs animation-settled (audit invariant documented).

**Owner-decision queue (47–50)**
47. TODO #329 decision (SSE) — see 17.
48. Release timing: ship v1.19.5 immediately after green pipeline, or batch with the next feature set (owner).
49. The `awesome-templ` PR stays blocked on coverage ≥80% (TODO #28) — decide whether to invest in coverage this quarter.
50. Confirm the three 07-08 questions (SSE/redirect/header link) so the canonical-URL story can close.

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **SSE bypass strategy (TODO #329):** Firebase Hosting will never stream rewrites. Which fix do you want — (a) map a Cloud Run-direct domain (e.g. `demo.lars.software`, needs DNS/domain access) and point the Datastar SSE route there, (b) migrate the demo to Firebase App Hosting (streaming-capable, new deploy pipeline), or (c) keep Hosting for pages but point the EventSource at the raw `*.run.app` URL via a deploy-time env var (+ CORS + CSP `connect-src` addition)? (c) is fully in-repo; (a)/(b) need infra access.
2. **Canonical URL policy:** should the raw `templcomponents-demo-….run.app` origin 301-redirect to `https://templcomponents.lars.software/demo/…` (clean canonical, but the visualtest harness + probes then need the `/demo` prefix — they already send it), or stay dual-serve? And do you want a persistent "Demo" link in the website header nav (currently: hero button, 404 page, README bar)?
3. **Release + guard packaging:** cut **v1.19.5** immediately after the pipeline goes green (the library contrast fixes + CountBadge/kanban/scrollback shades are consumer-visible), or batch with the next feature? And should the nonce-inventory guard land as a permanent visualtest gate in the SAME batch (my recommendation: yes — it is the guard this session's outage justifies)?

---

*Report ends. Verification pipeline (items 1–8) resumes on instruction; the kanban e2e rerun was in flight at writing time.*

---

## Addendum 2 — 13:55 CEST: pipeline COMPLETE, CI fully green, fixes live

All of §b closed; the pipeline ran to the end:

- **Kanban e2e** (the nonce regression's signal): move-buttons both transports + RTL + HTTP contracts + parity — all green after the nonce fix (+ the two remaining hardcoded `/` fetches in the shared contract probe moved to `fx.pagePath`).
- **Full clean visual witness run**: green twice (incl. the new `TestDemoCSPNonceIntegrity` + `TestDemoCSPJSExecutes` guards, which were promoted from the temporary diagnostics; the kanban diag was deleted).
- **ci-repro `--lint --website`**: VERDICT PASS at the pushed tip (actionlint requires the nix dev shell — run it as `nix develop -c scripts/ci-repro.sh --lint --website`).
- **Pushed** (e3b9e6b2→7f282d71 range): nonce fixes, contrast fixes, harness fixes, guards, goldens, go-directive alignment (root `go 1.26.0` — the daemon had flip-flopped it against go.work/visualtest), CSS artifact, vnu ignore class.
- **CI (ci.yaml) FULLY GREEN** on run 36857280299: Build & Test ✓ Lint ✓ Visual Regression ✓ CSS Freshness ✓ HTML validation ✓ — the first fully-green CI run of the day (it had been red on every push since 05:12Z: CSS freshness → HTML validation (two vnu rules) → tc-mirror sync → route goldens, each fixed at the root).
- **Website workflow** deployed the demo + site with the fresh CSS (run 36853607372 ✓).
- **Live re-verification of the fixed deploy**: 0 empty-nonce scripts across all probed canonical pages (was 50), fresh CSS served (emerald/amber/green-700 present), POST fragments 200, site `/` 200, all standalone recipe screens 200 through the proxy, `/demo/errors/404|500` return their intentional status codes, SSE still buffered (documented limitation, TODO #329 — unchanged, owner decision pending).

New defects found and fixed during the pipeline (beyond §d): the daemon's go.mod flip-flop broke the directive-equality guard and workspace toolchain resolution (aligned all to `1.26.0`); `ci-repro` must run inside `nix develop` (actionlint); the shared kanban probe still had two hardcoded `/` verification fetches; the route-golden index captures had baked a stale-CSS render (stacked filter row) and needed one legitimate regeneration; ViewTransitions style/script order violated vnu's first-child rule (fixed forward) and its `::view-transition-*` scoping rule cannot be satisfied by any markup (16th documented ignore class, AGENTS updated).

Open follow-ups unchanged from §f; top of the queue: permanent-guard promotion is DONE (item 9 ✓), items 10–16 (post-deploy JS probe, SSE canary, axe ledger re-triage, #330 tone sweep, `-update` guard, prerender nonce parity, GetNonce ADR) and the owner decisions (#329 SSE, canonical redirect, header Demo link, release timing).
