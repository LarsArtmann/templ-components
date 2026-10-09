# Status Report — Pareto plan execution: T1–T5 shipped, T6 red on master (daemon), T7–T20 pending

**Date:** 2026-10-08 23:49 CEST
**Session scope:** execute `docs/planning/2026-10-08_21-20_pareto-doc-recovery-plan.md` end-to-end ("the WHOLE list").
**Repo state at report time:** `master...origin/master [ahead 13]`; **`utils.TestDocsCountDrift` is RED on the local master tip** (daemon committed my in-progress T6 guard mid-debug — top priority fix below). A parallel session is active (BuildFlow triage, untracked `docs/status/2026-10-08_23-38_buildflow-triage-*.md` — not mine, not touched).

## Verdict

6 of 20 plan tasks fully executed (T1, T2, T3, T4, T5 + the F109 gate decision), 1 task ~70% done with a red test sitting on master (T6), 13 tasks not started, 1 gated task unblocked (T20 Command palette — the F109 flip). Working quality is good (every shipped task verified green at its own commit), but the daemon's mid-edit commit race turned my unfinished T6 into a red master tip — that is the single most important open item.

## (a) FULLY DONE

| Task | What shipped | Verification |
| --- | --- | --- |
| **T1 — ghost icons closeout (#371)** | The fix itself (iconAliases + `TestAllExportedNameConstsHavePathDataOrAreSpecial`) had already landed in the v1.21.0 window; this session closed it out: TODO #371 struck with resolution, CHANGELOG `[Unreleased]` Fixed entry added. Probed all three constants resolve to real path data (ArrowPath→Refresh, Bars3→Menu, HandThumbUp→ThumbUp); `AllIconNames()`=102 | icons suite + utils guards green; commit `7f37f754` (daemon took the TODO strike as `3558f8e7`) |
| **T2 — site route goldens (#374)** | v1.21.0 had already re-baselined the 8 site goldens; this session verified with a fresh-dist green read (`nix run .#visual -- -run TestSiteRouteGoldens -parallel 4`, ok 6.0s — real captures, not skips) + pixel eyeball of light-mobile bands (matrix wraps in scroll container, proof stats 123/102/63/7). #374 struck | Eyeballed via downscaled crops of the committed PNGs |
| **T3 — llms.txt (#375)** | `website/internal/build/llms.go`: `WriteLLMS`/`LLMSIndex` — llmstxt.org shape (H1, blockquote with canonical counts, Product + Documentation link sections) generated from the SAME parsed docs metadata as the search index (titles/URLs cannot drift). Wired into `run()`, footer link added, `assertLLMS` joined the site-build integrity test, `TestLLMSIndex` fixture test, pages goldens re-baselined. Live artifact verified: 18 docs entries + Product, deterministic. Intentionally NOT in the sitemap (file, not a page — writeSitemaps convention covers rendered top-level pages only) | Full website suite green; `467000fb` is mine |
| **T4 — F109 demand re-check (#217)** | Memo `docs/planning/2026-10-08_f109-demand-recheck.md` with fresh primary-source evidence (shadcn-templ llms.txt re-verified TODAY): **Command palette FLIPPED** (ecosystem demand + two named internal consumers: website docs search, demo navigation) → routed as TODO **#382** with the full testing ladder; DateRangePicker demoted to docs-recipe candidate (competitor composes it from Calendar+Popover+Button; we ship forms.Calendar); MultiSelect/TreeView/FileDrop absent from their registry → HOLD; Toast positions → HOLD. #217 updated with the verdict | Registry check done live via agentic_fetch of shadcn-templ.com/llms.txt |
| **T5 — competitor-claims sweep (#376)** | STANDOUT-IDEAS: dated competitor-landscape banner (templUI v1.9.3 no longer exists → shadcn-templ). javascript-guide: T.A.H. attribution fixed + templUI link → shadcn-templ. ui-library-design-research: dated note. future-architecture.d2: marked SUPERSEDED (pre-dates this library) + SVG re-rendered. comparison.md: quarterly re-verify cadence added to header. Post-sweep grep: zero stale templUI-as-current claims outside archived/dated/status docs | utils guards green; daemon captured as 16851b84/6990e4b1/bf5b2347 |
| **Bonus** | CHANGELOG `[Unreleased]` kept warm across all tasks; go.mod filter helper `isRepoLocalModule` concept proven in countLibraryTests (tests row now computes 1,622/1,647 → "~1,600 + ~1,650" verified honest) | — |

## (b) PARTIALLY DONE

**T6 — guard hand-typed numbers (#377) — ~70%, RED ON MASTER.**

Done:
- README "Packages | 20" → **17** (the 20 was inflated: it counted cmd/examples/internal tooling; live `go list` across the 7 published modules = 17 importable packages) + a definition footnote under the table.
- `utils.TestDocsCountDrift` extended: Packages row (live-computed), "across 7 Go modules" (live-computed), Tests row (`~N test functions + ~M subtests` with a ±49 rounding-band `assertRoundedCount` — currently passes against 1,622/1,647), website UseCases "23 form components" pinned to `packageCounts["forms"]`, and five comparison.md snapshot internals pinned (123 primitives, 102 typed names, 63 IsValid, 272 HTML goldens, 200 pixel goldens).

Broken / missing:
- **`countPublicPackages` returns 23, not 17** — the repo-local-module filter (`isRepoLocalModule`: skip dirs carrying their own go.mod, i.e. visualtest/, website/) was applied to `countLibraryTests` but the matching edit to `countPublicPackages`'s entry loop FAILED to apply during a multiedit batch and I did not re-apply it. Result: `README.md packages says 17; actual is 23`. **The daemon committed this red state to master** (commits `16b49346`/`d2eee934` window).
- 6.3 (statsDirs canonical pin test in website/internal/build) — not written.
- Full utils + website suites not re-run after the last edits (guardrail 1 violation-in-progress).
- Post-re-baseline green re-read of the site goldens (after the footer-link `-update` run) — the update run completed and its PNGs are committed (`d510d83d`), but I never captured the final `ok` verdict line.

## (c) NOT STARTED

T7 (README copy follow-ups), T8 (hero visual ogshot), T9 (comparison onto website), T10 (CHANGELOG Toggle duplicate + duplicate-title lint), T11 (site trust details), T12 (demo prose counts sweep), T13 (AGENTS ownership line + counter-definitions), T14 (htmx-4 probe plan), T15 (installable-blocks memo), T16 (themes gate criteria), T17 (parity-harness sketch), T18 (llms-full.txt — its dependency T3 is now done), T19 (crush-config docs-health PR), **T20 (Command palette build — gate OPENED by T4, TODO #382, full ladder required)**.

## (d) TOTALLY FUCKED UP

1. **Red guard test on master tip.** My unfinished `countPublicPackages` (23 vs 17) was auto-committed by the daemon while I was mid-debug — `utils.TestDocsCountDrift` fails on the current tip. This is the documented "daemon commits without a build gate" hazard (6th/7th classes in AGENTS.md) happening to my own work. Fix is one edit + test + commit.
2. **`\r` mangling in main_test.go (self-inflicted).** I stuffed literal `\r` escape sequences into a multiedit new_string; they landed as real carriage returns and collapsed the new assertLLMS function onto one line. Recovered with `sed 's/\r/\n/g'` + cleanup multiedit; cost ~10 minutes. Same failure class as AGENTS' "never patch Go via python heredocs" — the edit tool is not immune when I hand-roll escapes into strings.
3. **Commit-message loss to daemon races.** The daemon auto-commits within ~60s; 5 of my 6 task commits carry "chore: auto-commit …" heuristic messages instead of my drafted ones (only T3's `467000fb` kept its message). Content landed; history attribution is degraded.
4. **`countPublicPackages` logic shipped wrong.** I added the go.mod-scope rule asymmetrically (one helper got it, its sibling didn't) and verified only the passing half before the interrupt.
5. Minor: one no-op edit round-trip (TestLLMSIndex insertion attempt); an 877-line SVG diff from a d2 binary version mismatch (content validity-checked — note text, dimensions, closing tag — but not pixel-compared against the old render); a throwaway probe test file in icons (deleted immediately).

## (e) WHAT WE SHOULD IMPROVE

- **Commit per file-group immediately.** Batch-across-files editing + a 60s daemon = guaranteed message loss and mid-edit snapshots. New rule for this repo: edit → test → commit within one cycle, per task.
- **Never hand-roll escape sequences inside edit new_strings.** For multi-line content with escapes, `write` the file or use short anchors.
- **Check background job output before stacking work on top.** I started T6 edits while the visual `-update` ran and never read its verdict.
- **Symmetric scope rules.** Every new file-walk helper in this repo must share the repo-local-module (own go.mod → exclude) + nonConsumerDirs rule — ideally one shared helper, not two copies.
- **Guard-in-same-edit, applied to myself.** Yesterday's session typed "Packages | 20" unguarded despite harvesting exactly that lesson (B2 finding: "new doc numbers should enter the guard in the same edit"). Today T6 paid that debt. The guard now computes live so the class is closed.
- **Red-test exposure window.** The daemon snapshots the tree regardless of test state. Working-tree hygiene (never leave a red test uncommitted-and-unfixed across a turn boundary) is the only defense.

## (f) NEXT — ordered list (50)

1. **FIX RED:** apply `isRepoLocalModule` filter to `countPublicPackages`'s entry loop; `go test -run TestDocsCountDrift -count=1` green; commit immediately (beat the daemon).
2. Write 6.3: `TestStatsDirsAreCanonical` pinning `statsDirs` to the 9-package list (order-sensitive) in website/internal/build.
3. Full utils suite + `cd website && go test ./...` (guardrail 1) before the T6 commit.
4. Commit T6 with a real message; strike TODO #377.
5. Scoped green re-read of site route goldens post-re-baseline (`-run TestSiteRouteGoldens`, no `-update`).
6. T7.1 copy-editing pass over README intro→How-It-Compares.
7. T7.2 verify `ThemeScript("")` CSP behavior; soften the "CSP-safe page" hero line if needed.
8. T7.3 nav-link intent (blocked on owner answer Q2).
9. T7.4 GitHub-render eyeball (badges, 4-col table on mobile, anchors).
10. T8.1–8.3 hero visual via ogshot, EVERGREEN text (blocked on owner answer Q1).
11. T9.1 port comparison.md → website/content/docs/comparison.md with frontmatter.
12. T9.2 re-point related-projects.md to the on-site page; T9.3 website suite (link checker).
13. T10.1 reconcile CHANGELOG Toggle duplicate (check whether the parallel session's entry landed first).
14. T10.2 prototype a duplicate-title check for [Unreleased].
15. T11.1 "verified 2026-10-08" caption const under the site matrix.
16. T11.2 sitemap lastmod source for /sales — make it honest (data.go-aware).
17. T11.3 search-index snippet spot-check for related-projects.
18. T12.1/12.2 demo prose counts sweep (.templ/.go hand-typed icon/enum numbers).
19. T13.1 AGENTS.md competitor-facts ownership line (comparison.md + data.go single homes).
20. T13.2 counter-definitions doc section (components/icons/tests definitions + where computed — now partially self-documenting via the guard comments).
21. T14.1 read htmx 4 beta changelog + shadcn-templ #616 notes; T14.2 probe plan into ROADMAP.
22. T15.1 installable-blocks identity memo (module-only vs hybrid; tc add shape); T15.2 route decision to owner.
23. T16.1 style-themes demand-gate criteria → ROADMAP.
24. T17.1 parity-harness analogue sketch (what's pinned, what's diffed) → ROADMAP.
25. T18 llms-full.txt generator (docs corpus; T3's machinery is its foundation).
26. T19.1/19.2 crush-config branch: docs-health VERIFY external-claims + fitness legs; open PR citing the 19:57/20:45 reports.
27. T20 **Command palette** (gate OPEN): spec → templ → props/enum → CSP-safe singleton script → goldens → a11y sweep → e2e → docs counts (124!) → FEATURES/README/skill/demo hero bumps → tc mirror (#382, full ladder).
28. After T20: `TestDocsCountDrift` will fail on 123-vs-124 everywhere it pins — that is the guard working; do the count bumps in the same commit.
29. Final: full guards + `scripts/ci-repro.sh --lint --website` (with the actionlint PATH prepend, NOT nix shell) at the exact tip, then push IMMEDIATELY (13 ahead).
30. Verify v1.21.0 tags actually exist on origin (release.sh does not push; master is 13 ahead — the release commits AND their tags may be local-only; a consumer `go get` would 404).
31. Reconcile the standing BuildFlow preflight FAIL (`workspace/pseudo-version-hygiene` wants zero pseudo-versions; the release script deliberately re-adds real v1.21.0 requires post-cut) — decide which convention owns the post-release state.
32. Resolve the AGENTS.md size-cap disagreement (BuildFlow warns max 220; the repo's go-structure-linter enforces 377) — two tools, two truths, warning noise on every commit.
33. website/go.mod `ignore dist` warning (path doesn't exist at check time) — silence or create.
34. lychee exclude for docs/feedback/archived (standing preflight warn).
35. Rebuild/reinstall BuildFlow binary (ec8d2d3 vs HEAD c773fd9 staleness warn).
36. Triage the 7 standing go-structure-linter findings (assets/ dir suggestion, README install section, testdata placements) once — accept-or-fix.
37. Route the daemon commit-message quality problem upstream (BuildFlow generates from template, not `git diff --stat`).
38. Make the go-licenses wrapper durable (TODO #351/#352 — machine-local `~/.local/bin` wrapper re-breaks on BuildFlow reinstall).
39. Extend assertLLMS: descriptions non-empty for docs pages that have them; parity with search index already structural.
40. Add a self-check pinning countPublicPackages/countLibraryTests against known constants (17 / band) so counter rot is caught independently.
41. Single-source the public-package enumeration (utils guard + website build could share one definition).
42. Document the "llms.txt intentionally not in sitemap" convention in writeSitemaps' comment.
43. Post-deploy smoke: curl /llms.txt in website.yml (new surface = new smoke check).
44. README "Typed enums 64" — no live counter exists for the TOTAL (only the 63 IsValid); add one or drop the 64.
45. Pixel-check the regenerated future-architecture.svg against intent (d2 version re-layout).
46. robots.txt: consider referencing llms.txt.
47. TODO_LIST strikes for #375/#376/#377/#378-#381 as their tasks close (currently none of the new strikes beyond #371/#374 are in).
48. Update the pareto plan doc with per-task status so the plan reflects reality.
49. AGENTS.md Website section: add the llms.txt surface line (one line, mind the cap).
50. After everything: status-report harvest → TODO_LIST/CHANGELOG warmth re-verified before any release cut.

## (g) Questions I cannot answer myself

1. **Push + release completion:** master is 13 commits ahead INCLUDING the v1.21.0 release commits and (possibly) its tags, and release.sh never pushes. A parallel session is active. Do you want me to run the M03 ritual at the tip and push master + tags as part of my final step (completing the release), or is the release owner pushing separately and I should push only after they do?
2. **T8/T7 owner inputs (asked twice, still open):** hero visual — regenerate an EVERGREEN ogshot card and wire it into the README, or keep the README image-free? And the "Why" nav link — external /sales URL or in-page anchor?
3. **Pseudo-version hygiene conflict:** BuildFlow's preflight FAILs every commit since v1.21.0 because requires sit at real versions with local replaces — which is the release script's documented post-release state. Which side changes: BuildFlow's rule (learn the post-release state) or the repo (re-normalize to zero pseudo-versions after every release)?

## Provenance

- Plan: `docs/planning/2026-10-08_21-20_pareto-doc-recovery-plan.md` (T1–T20)
- F109 memo: `docs/planning/2026-10-08_f109-demand-recheck.md`
- Upstream fix consumed: v1.21.0 window (iconAliases + consistency test; site goldens re-baseline)
- Session commits (mine): `7f37f754` (T1 changelog), `467000fb` (T3 llms.txt); daemon-carried: `3558f8e7`, `0e8316fc`, `16851b84`, `6990e4b1`, `bf5b2347`, `d510d83d` (route-golden PNGs + templ regen), `d2eee934`, `16b49346` (T5/T6 files)
- Parallel session: `docs/status/2026-10-08_23-38_buildflow-triage-*.md` (untracked, not touched)
