# Status Report: Demo Overhaul — Final Verification Pass

**Date:** 2026-08-17 15:14 CEST
**Session scope:** Completing the unverified tail of the "Is our demo page superb?" audit & overhaul (13 findings, ~90% implemented in prior session, verification never run). This session: verify everything, fix what verification found, close housekeeping.

---

## a) FULLY DONE

1. **Build verified.** `templ generate ./examples/demo/... && go build ./examples/demo/` → BUILD_OK. The 7 compile errors fixed last session were confirmed good.
2. **Daemon edit reviewed (commit 7c2e3c2).** The BuildFlow daemon replaced my `if pages < 1` clamp in `usersTotalPages` with `max((total+pageSize-1)/pageSize, 1)`. Reviewed: behaviorally identical, arguably cleaner. **Kept.**
3. **CSS recompiled** via `nix run .#css`. Surprise: the flake app minifies (old committed bundle was unminified) — file shrank 114601B → 90642B. I verified this was NOT a content regression: occurrence counts of `tc-modal`, `tc-fluid`, `@container`, `animate-spin`, `popover` identical old vs new. `TestCSSFreshness` in CI is now satisfied.
4. **Drift-guard test caught a REAL bug:** `TestHeroCountsMatchFeatures` failed — FEATURES.md claims **106 icons**, `icons.AllIconNames()` returns **102** (ground truth verified: 101 in `iconPathData` + Spinner). Fixed 4 stale references in FEATURES.md (Totals line, module table, Icon row, section heading + "105 path icons" → 101). The drift-guard written last session paid for itself on its first run.
5. **All tests pass** across all modules (`nix run .#test`), incl. the previously-failing demo package.
6. **Lint clean:** 0 issues across all 7 modules (`nix run .#lint`).
7. **Live probe: all 19 routes return 200** — `/`, `/users` (+sort/dir/page query variants), `/forms`, all 4 recipe pages, `/css/app.css`, `/health`, and all 6 API endpoints (`/api/items`, `/api/items/123`, `/api/save`, `/api/demo-stats?tick=0|3`, `/api/users` + filters).
8. **Stale process killed.** PID 2760026 (from a prior session, bound :8091) was silently serving the OLD binary — my first probe run showed 404s for every new route and I nearly misdiagnosed it as missing handlers. `kill` builtin failed ("unsupported builtin" in mvdan/sh); used `/run/current-system/sw/bin/kill -9`. Free now.
9. **`/users` added to `prerender.go`** (was a known gap — server-only page). Prerender run verified: writes 7 pages incl. `users/index.html`.
10. **CHANGELOG `[Unreleased]` entry added** (repo rule violated by daemon commit 97e7bb0): demo overhaul summary + the FEATURES.md 106→102 correction.
11. **`*_templ.go` tracked:** all 14 demo generated files in `git ls-files`; `.gitignore` `!*_templ.go` unignore intact.

## b) PARTIALLY DONE

- **Demo server left running at :8091** (`/tmp/tc-demo`, fresh binary). Intentional for inspection, but it's another stale-process-in-waiting — should be killed before session end or noted.
- **Working tree uncommitted:** `CHANGELOG.md`, `FEATURES.md`, `examples/demo/prerender.go`, `examples/demo/static/app.css` modified + status reports untracked. The BuildFlow daemon will commit them with a generic message unless I commit first with a proper one.

## c) NOT STARTED (carried over, deliberately out of scope this session)

- `examples/demo/demo.out.css` — purpose still unverified. It's byte-identical to the freshly compiled `static/app.css` now (the daemon's 7c2e3c2 refreshed it), but nothing references it that I found. Likely dead artifact; needs a removal decision.
- `website/src/styles/global.out.css` was recompiled by the daemon in lockstep — not re-verified by me.
- Outstanding user questions from prior session (see g).

## d) TOTALLY FUCKED UP (mistakes this session)

1. **Port hijack misdiagnosis risk:** First live probe returned 404s on all new routes. I initially re-ran the probe without checking whether the server had even started — `/tmp/tc-demo.log` said "address already in use" but the OLD process answered. Wasted one probe cycle. Lesson: check the server log before probing when a port might be occupied.
2. **`kill` is not a builtin in this shell** (mvdan/sh) — first kill attempt silently failed, so the "fresh" server never started and the second probe hit the stale process AGAIN. Should have verified the port was free immediately after kill.
3. **CSS minification judgment call made autonomously:** the flake `css` app minifies, the previously committed CSS did not. I verified class coverage identical and accepted it, but this changes the artifact shape consumers see — borderline rule-breaking (state rule, document tradeoff). Documented here instead of asking first.
4. **Edit tool discipline:** attempted to edit CHANGELOG.md without viewing it first in this session (tool correctly rejected me). Minor, but a wasted round trip.

## e) WHAT WE SHOULD IMPROVE (structural)

1. **FEATURES.md counts drift silently.** The icon count was wrong (106 vs 102) for unknown duration. The new demo test only guards the demo hero against FEATURES.md — it does NOT guard FEATURES.md against the actual icon package. A test in `icons/` asserting `len(AllIconNames())` matches FEATURES.md (or better: generating that line) would close the loop at the source.
2. **Same for component counts:** `componentCount` const in demo.templ is still hand-maintained against FEATURES.md's "116 templ components". One end-to-end guard, but two hand-maintained numbers upstream.
3. **Stale-process hygiene:** prior session left :8080 occupied "never killed" (it was free by this session — someone killed it), this session inherited :8091. A convention (unique ports per session, or a check script) would prevent masked-binary bugs like d.1.
4. **BuildFlow daemon commits unverified work with generic messages** (known T13 issue). This session's verification changes will likely be daemon-committed too. Fix lives in `larsartmann/buildflow`.
5. **The flake `css` app and the committed CSS disagree on minification.** Either the app always minified and someone once committed unminified output, or the app changed. Decide one shape; a checksum/assert in CI (`TestCSSFreshness`?) would catch future divergence before it becomes archaeology.
6. **`demo.out.css` ambiguity** — either delete it or make it the single source that `static/` copies from. Two files, byte-identical, no documented owner = split brain waiting to happen.

## f) NEXT: up to 50 things

**Verification & release hygiene**

1. Commit the working-tree changes with a proper message (before daemon does).
2. Kill or document the :8091 demo server.
3. Decide fate of `examples/demo/demo.out.css` (delete or single-source).
4. Verify `website/src/styles/global.out.css` recompile (daemon's 7c2e3c2) against website sources.
5. Add `[Unreleased]` entry note to `docs/status/2026-08-17_14-26...md` that its "verification pending" state is now resolved.
6. Consider cutting v1.8.5 / v1.9.0: demo overhaul + CircularProgress/SectionHeading/DateRange + StatCard ValueID are sitting in `[Unreleased]`.

**Drift-guard hardening**
7. `icons` package test asserting `len(AllIconNames())` against FEATURES.md (close the loop at source).
8. Generate FEATURES.md Totals line from code instead of hand-editing.
9. Guard `componentCount` (demo) against a real count of components, not just FEATURES.md parity.
10. Extend `TestCSSFreshness` to also assert minification shape consistency.
11. Guard `.gitignore` against BuildFlow re-appending `*_templ.go` (already known; a pre-commit grep exists — maybe move into treefmt).

**Demo depth (from audit, not yet done)**
12. ECharts demo currently hand-authored option strings — consider a `//go:build ignore` example showing the real go-echarts `RenderSnippet()` flow.
13. `/api/save` returns 70B — verify it returns something LoadingButton/HTMX can actually swap (swap="none" means it's ignored; fine, but confirm UX).
14. Confirm `/api/items` EndOfList terminator actually renders (probe returned 662B at cursor=2; check cursor beyond end).
15. Dark-mode visual pass of new pages (`/users`, `/recipes/auth`) via `nix run .#visual` if pages are in scope.
16. RTL pass (`dir=rtl` probe) of new pages.
17. Add `/users` and `/recipes/auth` links to demo nav/TOC on the index page (verify they're reachable by clicking, not just URL).
18. `demo.out.css` vs `static/app.css`: if kept, wire the flake app to emit both or neither.

**Known carried-over questions/infra**
19. Canonical docs URL for hero link (question below).
20. Move demo mock handlers out of `main.go` into `handlers.go`?
21. Revert decision on daemon's `usersTotalPages` edit — I decided KEEP; confirm.
22. BuildFlow generic-commit-message fix (upstream repo).
23. BuildFlow pre-commit re-appending `*_templ.go` to `.gitignore` (upstream fix).
24. BuildFlow 60s budget skipping `go test ./...` (upstream fix).
25. mvdan/sh lacks `kill` builtin — add a note to global memory so future sessions use `/run/current-system/sw/bin/kill` or `pkill` directly.

_(26–50 intentionally unfilled — nothing else concrete surfaced this session; padding would be theater.)_

## g) QUESTIONS (cannot figure out myself)

1. **CSS minification:** the flake `css` app now emits minified `static/app.css` (90KB vs the old unminified 114KB, identical class coverage). Was the old unminified commit intentional (e.g., for debugging/inspection), or is minified the desired steady state? Deciding wrong means churn on every recompile.
2. **`demo.out.css`:** delete it (dead artifact), or is it load-bearing somewhere I haven't found (e.g., BuildFlow DAG target / website import)? You know BuildFlow's contract; I can only see this repo.
3. **Hero docs link:** what is the canonical URL for the demo hero's documentation link — the Firebase-hosted website (`templ-components.lars.software` or similar), `pkg.go.dev/github.com/larsartmann/templ-components`, or the GitHub README? I won't guess a URL.
