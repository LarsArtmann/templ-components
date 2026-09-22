# Status Report — tc Scaffolder Mirror Guard + art-dupl Baseline Session

**Date:** 2026-09-22 22:38
**Session scope:** (1) explain `cmd/tc/_sources` + art-dupl output, (2) make the `_sources` guard a bidirectional self-healing mirror, (3) clean up the `-t 1` clone-scan noise via a committed baseline, (4) extract genuine test-scaffolding duplication.
**Repo tip at writing:** `191090ea` (daemon-committed; AGENTS.md + CHANGELOG.md sweeps pending in daemon cycle). Nothing pushed.

---

## Opening honesty pass: what did I forget, what could have been better?

### What I FORGOT

1. **The truncated art-dupl capture nearly poisoned everything.** My first `-t 1` run captured only **20 of 39 groups** (101 lines, missing the `Found total` line — I noticed it was missing and *rationalized it away* instead of investigating). I told the user "20 groups, your 41 is wrong" with high confidence. Only a later re-run for a different purpose revealed **39 groups** — the user's 41 was essentially right. Had I recorded the baseline from the first capture, ~19 accepted groups would have been silently unbaselined and `art-dupl check` would have false-failed later. Caught by luck (a re-run), not by my own rigor.
2. **No sync guard between the two hand-copied mirror lists.** `MIRRORED_PKGS` (bash array in `scripts/check-tc-sources-sync.sh`) and `mirroredPackages` (Go slice in `cmd/tc/main_test.go`) must stay in sync and are held together by **comments only**. This repo has an exact precedent for machine-guarding hand-copied lists (`scripts/check-lint-modules.sh` exists because the lint module list drifted across 3 files). I had the pattern in front of me and didn't apply it.
3. **No end-to-end `tc add` smoke test for a rescued component.** I proved registry registration (`tc ls` lists all 22) and read the `_types.go` copy code path (cmd/tc/main.go:382-395), but never actually ran `tc add eyebrow` into a temp dir to watch the `.templ` + `_types.go` land.
4. **Pair-completeness of the 9 rescued `*_types.go` files unverified.** `tc add X` copies `X_types.go` only if `X.templ` is registered. I verified pairs for `auth_layout`/`notfound404` mentally but did not programmatically assert all 9 `_types.go` have a registered `.templ` sibling (e.g. `layout/container_types.go` ↔ `container.templ`).
5. **`packageDeps`/`packageImports` maps unaudited for the 22 rescued components.** These hand-maintained maps drive `tc add --list-deps`. Kanban, the SVG charts, FilterInput, DirtyGuard etc. postdate parts of those maps; dep listings for the newly addable components may be wrong or stale. `TestPackageImportsMatchSources` guards the imports map but only for packages it lists.
6. **Stale `tc new` terminology left in place.** The old guard header said `tc new`; the commands are `tc init`/`tc add`. I used the correct names in new text but did not sweep the stale references (AGENTS.md CSS-inventory bullet, possibly docs/ and website content).
7. **`starter/` contains 5 CSS files but `tc init` writes only 2** (`app.css`, `custom.css`). `templ-components-theme.css`, `styles.css`, and the compiled `templ-components-theme.out.css` inside starter have no in-code consumer that I saw. Possible dead scaffolder weight — noticed, not investigated (out of session scope by user instruction).
8. **The guard's "unfixable problem" path (MISSING PACKAGE) was never scenario-tested** — I tested drift/new/orphan/clean, not the config-error branch.
9. **Full pre-push ritual not run.** I ran the tracked hook + touched-module lanes, but not `scripts/ci-repro.sh --lint --website` (Build & Test + Lint + CSS + Visual + Website lanes). Acceptable because nothing is being pushed, but my "final verify" was narrower than the repo's own definition of green.

### What I could have done BETTER

- **Trust no truncated output.** The `Found total N` line is the ground truth of an art-dupl run; its absence means the capture is incomplete. I should have re-run immediately instead of building a classification on 51% of the data. This is the session's only genuine process failure, and it was self-caught late, by accident.
- **Written the list-sync guard in the same commit as the lists** (item 2 above) — one small test, zero drift class.
- **Made the extractions browser-proven.** The kanban e2e helper refactor is compile-checked (`go vet`) but the 5 affected tests were not run in Chromium (`nix run .#visual`). The wire_demo tests are self-proving (HTTP tests ran green); the kanban ones are not.
- **Measured the guard** instead of inheriting the old "<100ms" claim — the new script does more work (two directions, all files); the header should state a measured number.
- **Recorded per-threshold counts for the ADR** — I documented only the t=1 baseline count; the ADR's historical per-threshold table (-t 5/7/8) is now stale-ish and labeled historical without current numbers.

### What I could STILL improve (forward-looking)

- Turn the baseline from *documentation* into an *enforced gate* (CI/pre-commit wiring — currently nothing fails if someone introduces new duplication).
- Regression-test the guard itself (a script that spins the 4 scenarios in a temp worktree instead of my one-off manual run).
- Give the guard an explicit, loud opt-out (e.g. `TC_SKIP_SYNC=1`) for the rare "I really meant to edit `_sources` directly" case — currently the guard would clobber such edits silently-by-design.
- Reduce 1-line noise classes (`defer cancel()`) from future baselines via `--min-lines`, or accept and document them.

---

## a) FULLY DONE

| # | Work | Evidence |
|---|------|----------|
| 1 | **Bidirectional mirror guard** — `scripts/check-tc-sources-sync.sh` rewritten: direction 1 (embedded→library: drift + orphans), direction 2 (library→embedded: `.templ` + `*_types.go`), full `--fix` (re-copy / add / `git rm`), fix-mode exits 0, check-mode exits 1 | Worktree scenarios: clean→0, drift→caught+refixed, new→caught+added, orphan→caught+removed; all exit codes verified |
| 2 | **Pre-commit Guard 8 auto-fixes and stages** (`.githooks/pre-commit`) — runs `--fix`, greps `^synced `, `git add`s `cmd/tc/_sources`, prints loud summary; fails only if unfixable | Simulated hook block in worktree: drift fixed + staged (`STAGED:` list printed); full tracked hook run at tip: exit 0 |
| 3 | **Go guard direction 2** — `TestSourcesShipEveryMirrorableFile` in cmd/tc/main_test.go (recursive os.DirFS walk, `.templ`+`*_types.go`, fs.Stat against embed FS) | `go test -count=1 ./cmd/tc/...` green; would have failed before the 22-file sync |
| 4 | **22 never-embedded files rescued** — kanban, line/pie/area charts, eyebrow, scrollback, section_heading, collapsible_section, date_range, circular_progress, filter_input, dirty_guard, auth_layout + 9 `*_types.go` — `tc add` rejected all of them before ("unknown component") | `bash scripts/check-tc-sources-sync.sh --fix` → `synced 22 file(s) (added 22…)`; recheck exit 0; `tc ls` lists them |
| 5 | **Committed hash baseline** `.art-dupl-baseline.json` (39 groups, t=1, type-aware) + ADR-0009 rewritten (baseline section, canonical invocation, consequences) | `art-dupl check -c .art-dupl.json -t 1 --type-aware` → "No new clones detected (baseline: 39 groups)", exit 0 |
| 6 | **Classification of all 39 groups** — component-idiom/templ-DSL, demo content, CLI-tool boilerplate (`visualtest/tools/*`), website page content; zero `_sources` leaks (exclusion config verified correct) | Full inventory captured and classed; documented in ADR-0009 |
| 7 | **Two genuine 5-site duplications extracted instead of baselined** — `newKanbanReadyTab` (visualtest/kanban_e2e_test.go, 5×19-line bootstrap) and `fetchDemoHTML` (examples/demo/wire_demo_test.go, 5×17-line fetch scaffold) | demo tests green (3.4s), `go vet ./...` in visualtest clean, both groups gone from scan |
| 8 | **Docs**: ADR-0009 updated; AGENTS.md dedup bullet rewritten (canonical gate) + new self-healing-mirror gotcha bullet; CHANGELOG `[Unreleased]` warmed with both entries | In tree; daemon committed ADR/baseline in `191090ea` |
| 9 | **Verification suite**: tracked pre-commit hook end-to-end (exit 0), cmd/tc + demo tests with `-count=1`, visualtest vet, golangci-lint on cmd/... and visualtest (0 issues), `nix fmt` (0 changed), `art-dupl check` (0 new) | All green at tip `191090ea` |

## b) PARTIALLY DONE

| # | Work | Done | Missing |
|---|------|------|---------|
| 1 | Final verification | Hook + touched-module tests/lint/vet/fmt at tip | Full `scripts/ci-repro.sh --lint --website` (all CI lanes) not run — nothing is being pushed, so ritual not triggered |
| 2 | Browser-proof of extraction | Compile + vet for kanban e2e; real HTTP tests for wire_demo | The 5 kanban e2e tests not executed in Chromium (`nix run .#visual`) |
| 3 | `tc add` completeness for rescued components | Files embedded, registry lists them, `_types.go` copy path read and confirmed | No runtime smoke (`tc add` into temp dir); no pair-completeness assertion for all 9 `_types.go`; `packageDeps`/`packageImports` accuracy for the 22 unaudited |
| 4 | Dup-gate institutionalization | Baseline recorded + canonical invocation documented in ADR + AGENTS.md | Gate is not enforced anywhere (no CI lane, no hook step) — today it's documentation, not a gate |
| 5 | Terminology hygiene | New text uses `tc init`/`tc add` correctly | Old `tc new` references not swept (AGENTS.md line ~280, possibly docs/, website content) |
| 6 | Guard performance claim | Old "<100ms" header inherited | New guard does strictly more work; no measurement taken |

## c) NOT STARTED (identified this session, deliberately out of scope)

1. `art-dupl check` wired into `ci-repro.sh` and/or a CI lane (provisioning question — see Questions).
2. art-dupl packaged in `flake.nix` (deterministic local+CI availability).
3. Sync guard between `MIRRORED_PKGS` and `mirroredPackages` (comment-only today).
4. Guard self-test script (4 scenarios automated in a temp worktree).
5. `starter/` dead-CSS investigation (3 of 5 starter files appear unconsumed by `tc init`).
6. Stale `tc new` → `tc init`/`tc add` doc sweep.
7. `packageDeps`/`packageImports` audit for the 22 rescued components.
8. Website content check: do docs/api-reference pages claim a component list/count that `tc add` coverage just changed?
9. Opt-out escape hatch (`TC_SKIP_SYNC=1`) for the guard.
10. art-dupl upstream (own tool): investigate why a piped/tee'd run truncated its report (my 20-vs-39 incident) and whether `Found total` can be made flush-reliable; consider `--fingerprint-only` baselines (no `recordedAt` churn).
11. Project AGENTS.md gotcha: "never trust a piped CLI capture whose final summary line is missing" (session lesson, not yet written down anywhere).
12. Per-threshold count refresh in ADR-0009 consequences (current -t 5/7/8 numbers vs historical).

## d) TOTALLY FUCKED UP

**Nothing is broken or destroyed.** No data loss, no pushed bad code, no reverted others' work, all guard changes verified against four hand-built failure scenarios plus the real hook. The one candidate for this section, stated plainly:

- **The 20-vs-39 misclassification.** I built my first complete world-model on a truncated tool output, dismissed the user's contradicting number (41), and only stumbled onto the truth via an unrelated re-run. If the session had ended one command earlier, the committed baseline would have been wrong (20 groups) and silently so. **Damage: zero (caught in-session). Lesson: unresolved.**

## e) WHAT WE SHOULD IMPROVE (process + code, distilled)

1. **Gate rule:** a CLI report is only "read" when its final summary line is present in the capture; otherwise re-run. (Personal tooling rule worth encoding in project AGENTS.md gotchas.)
2. **Never ship twin lists guarded by comments** — this repo already learned that lesson with lint module sets; apply it to the mirror package lists immediately.
3. **"Extraction or baseline" decisions deserve their proof class attached:** test-scaffolding extractions should name their verification level (compile-only vs browser-run) in the ADR/CHANGELOG entry.
4. **Enforced > documented:** the baseline only pays off when `art-dupl check` runs automatically somewhere that blocks.
5. **Guards that auto-fix need a loud, documented escape hatch** — silent self-healing plus a daemon that auto-commits is a combination that can bury a user's intentional action without a trace.
6. **Rescuing surface area changes product behavior** (`tc add` accepts 22 more components) — that warrants the same doc-count drift attention as adding a component to the library itself.

## f) TOP THINGS TO GET DONE NEXT (brainstorm, impact-ordered; ROADMAP fuel — harvest into TODO_LIST.md)

**Scaffolder / guard hardening**
1. Sync-guard `MIRRORED_PKGS` (bash) vs `mirroredPackages` (Go) — parse the bash array in a Go test, fail on mismatch.
2. `tc add` smoke test in cmd/tc: scaffold a rescued component into `t.TempDir()`, assert `.templ` + `_types.go` land.
3. Pair-completeness test: every embedded `*_types.go` must have a registered `.templ` sibling (and vice versa).
4. Audit `packageDeps`/`packageImports` for the 22 newly addable components (kanban, charts, filter_input, dirty_guard, auth_layout…).
5. Guard self-test: `scripts/test-tc-sources-guard.sh` — spins the 4 scenarios in a temp detached worktree.
6. Test the guard's MISSING-PACKAGE (unfixable) branch.
7. Add `TC_SKIP_SYNC=1` opt-out to the guard + hook, with a loud stderr banner.
8. Rename-component scenario test (delete+add in one commit → orphan+unembedded simultaneously).
9. Measure guard runtime on the full tree; replace the inherited "<100ms" header claim.
10. Sweep all docs for stale `tc new` references → `tc init`/`tc add`.
11. Investigate the 3 apparently unconsumed `starter/` CSS files (incl. the compiled `.out.css` living inside starter); delete or wire, per the CSS-inventory rule.

**Duplication gate**
12. Provision art-dupl for CI (flake input of the art-dupl repo or setup step) — prerequisite for everything below.
13. Wire `art-dupl check -c .art-dupl.json -t 1 --type-aware` into `scripts/ci-repro.sh` + a CI lane.
14. Decide canonical threshold: keep t=1 (accepts 1-line noise groups in baseline) vs `--min-lines 2/3` or t=3 as the gate with t=1 informational.
15. Decide class policy: exclude `website/` content and/or `visualtest/tools/*` boilerplate via `.art-dupl.json` (vendored-class argument) vs keep them baselined for churn visibility.
16. Evaluate `//art-dupl:accept` directives for point-suppression vs repo-wide baseline for a few truly-local clones.
17. Re-record per-threshold counts (-t 5/7/8) for ADR-0009 consequences; mark historical table as such.
18. Add "run `art-dupl check` before claiming dedup done" to `docs/plan-authoring-checklist.md`.
19. Baseline-regeneration ritual documented (update ADR table first, then re-baseline) — link from release checklist.
20. Verify baseline determinism: two consecutive `check` runs produce identical results.
21. File upstream issue on art-dupl (own repo): truncated/piped output lost ~19 of 39 groups in one capture; investigate flush/summary-line reliability.

**Verification debt**
22. Run `scripts/ci-repro.sh --lint --website` at tip before the next push (ritual; my verify was narrower).
23. Browser-proof the extracted kanban e2e helper (`nix run .#visual`, targeted kanban e2e).
24. Explicit `go build ./...` (root + all modules) at tip — embed-set changed, nothing depends on it breaking, but prove it.
25. Run website module tests for completeness (`cd website && GOWORK=off go test ./...`).

**Docs / product surface**
26. Check website content + FEATURES/README for `tc add` coverage claims now changed by the 22 rescues (docs-count drift class).
27. CHANGELOG entries: consider splitting the 22-file rescue into its own titled entry for release-notes visibility (it's the user-facing headline).
28. Document the mirror's auto-fix behavior in the consumer-facing docs (if `tc` has a docs page on the website).
29. AGENTS.md gotcha: "piped CLI capture without its summary line = truncated; re-run" (session lesson).
30. Consider listing per-package `tc add` examples in the scaffolder docs using freshly rescued components (kanban is a good showcase).

**Daemon-resilience / hygiene**
31. Post-daemon-commit sanity ritual: after daemon sweeps, re-run `art-dupl check` + `bash scripts/check-tc-sources-sync.sh` once (daemon has a history of regressing same-day fixes).
32. Confirm `.art-dupl-baseline.json` survives future `*.json` gitignore patterns (currently tracked; fine — just don't add a broad `*-baseline.json` ignore later).
33. Consider `git rm` fallback in guard: for untracked strays it uses `rm -f` — acceptable, but a one-line comment in AGENTS.md would prevent a future "why does a script rm?" audit flag.
34. Guard output: emit `synced N file(s)` also when invoked with `--fix` AND already-in-sync-but-untracked-files-exist (edge: daemon staged something odd) — review edge interplay with `git add` in hook.

**Upstream / tooling**
35. art-dupl feature: `--fingerprint-only` baseline (drop `recordedAt` → diff-stable re-baselines).
36. art-dupl feature: config auto-discovery or a well-known config name to kill the "scans without `-c` re-report" foot-gun documented in AGENTS.md.
37. art-dupl: document that `baseline`/`check` flag sets must match exactly (hash stability) in `--help` epilog.

**Backlog candidates noticed in passing (not researched — one-liners only)**
38. `starter/templ-components-theme.out.css` — tracked exception to the `*.out.css` gitignore; re-justify or delete.
39. `scripts/pre-commit.sh` (full pre-push) calls the guard in check mode — verify its failure output is as actionable as the hook's.
40. `TestSourcesMatchPackageFiles` direction-1 error message still says "re-copy it (and re-run this test)" — point it at the script for consistency.
41. `datastar/*docs*` and `src/datastar/...` vestigial exclusion patterns were removed from the bash guard — confirm no other script/test copied those dead patterns.
42. `tc add --list-deps` for chart components: the SVG charts share `chart_geometry.go` — dep listing should say so (audit belongs with #4).
43. Consider surfacing "N components addable via `tc add`" as a derived count in `tc ls` footer (kept drift-guarded by tests, not prose).
44. Worktree-based guard tests (#5) could double as the fixture for #6/#8 — one script, three scenarios.
45. `CHANGELOG` [Unreleased] currently has no `### Fixed` — if the 22-file rescue is framed as a bugfix for consumers, move/alias it there.

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **How should `art-dupl check` be provisioned and enforced?** It's currently nowhere in CI. Is the art-dupl repo public/packagable as a flake input (deterministic devShell + CI), or should CI build/install it another way — and do you want it as a blocking lane, or advisory in `ci-repro.sh` only?
2. **What is the intended canonical threshold/class policy?** Keep t=1 (and baseline the 1-line `defer cancel()` noise), or gate on `--min-lines`/t=3 with t=1 informational? And should `website/` content + `visualtest/tools/*` boilerplate be excluded via config (like `_sources`) instead of baselined?
3. **Product intent for the 22 rescued components:** do you actually want all of them `tc add`-able, given the documented caveat that big ones (kanban, SVG charts) reference sibling files and won't compile standalone — or should some be deliberately unshippable/hidden from the registry (which would mean a `tc add` allowlist rather than my complete-mirror semantics)?

---

*Point-in-time snapshot. Section (f) is HARVEST fuel for `TODO_LIST.md`/`ROADMAP.md` (docs-health), not an entombed promise list.*
