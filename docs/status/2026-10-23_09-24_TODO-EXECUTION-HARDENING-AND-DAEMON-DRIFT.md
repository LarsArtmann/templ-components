# TODO Execution, A11y Hardening & Daemon-Drift Report — 2026-10-23 09:24 CEST

**Repo:** `github.com/larsartmann/templ-components` · **Version:** 1.23.0 ·
**HEAD at report time:** `60c74625` (daemon auto-commits; see §6.9) ·
**Branch:** `master` · **Tree:** clean · **Pushed:** yes — `820e2df1` (see Addendum).

**Scope of this session (single-writer claim):** `TODO_LIST.md`, the new plan doc,
`display/heatmap.{templ,go}` + heatmap goldens + `display/a11y_new_test.go`,
`visualtest/axe*.go`, `visualtest/golden.go` + `*_visual_test.go` update helper,
`website/internal/md/md.go`, `website/cmd/site/main.go`, `utils/golden/golden.go`,
`docs/adr/0009-accepted-clones.md`, `AGENTS.md`, `CHANGELOG.md`,
`templates/templ-components-classes.txt`, `tools/crossbrowser-spike/package.json`,
`cmd/tc/_sources/display/heatmap.templ`, and the templ pins in `go.mod`,
`website/go.mod`, `visualtest/go.mod`.

**Primary artifact:** `docs/planning/2026-10-23_09-03_ALL-TODOS-REMAINING-PARETO-PLAN.md`
— complete open-TODO inventory, Pareto tiers, gates G1–G6, and ≤12-minute micro-task
breakdown for every open row.

---

## 1. Executive scoreboard

| Metric | Value |
| --- | --- |
| Open TODO rows inventoried | 65 (all classes: EXEC/G1–G6/OTHER) |
| Rows fully closed this session | 8 (#29, #189, #294, #357, #437, #435, plus #273/#291/#297/#298/#299/#300/#306/#309/#317 struck as already-complete) |
| Code fixes shipped | 2 (#437, #387) |
| Verification/hygiene strikes | 5 (#29, #189, #294, #357, already-complete rows) |
| Unplanned daemon-drift repairs | 5 (templ pins, playwright pin, tone doc, tc mirror, class inventory) |
| Modules tested green | root + utils + icons + errorpage + charts/echarts + datastar + htmx + website; visualtest compiles |
| Pixel proof | `TestDemoRouteGoldens` green in **92.8 s** (all demo routes, light/dark/mobile/RTL) |
| Lint | 0 issues on display, utils/golden, visualtest, website |
| Full push ritual (`ci-repro --lint --css --visual --website`) | **PASS** (09:41:44 CEST — see Addendum) |
| Full visual suite (`nix run .#visual`, all tests) | **NOT RUN** (route goldens only) |
| Pushed / CI-verified | Pushed `820e2df1`; CI run `39919529944` in progress at report close |
| Root-module coverage | **70.4 %** (70 % gate passes; #28 needs 80 %) |

**Verdict:** the session delivered real correctness and hygiene value, and repaired a
repo that was silently red before I touched it. It did **not** reach "everything works
and is pushed": the full pre-push ritual, the complete visual suite, the achecker
confirmation, and every owner gate remain open. Treat the tree as **green locally,
unproven in CI**.

---

## 2. What is FULLY done

### 2.1 #437 — Heatmap `aria_descendant_valid` (real a11y fix)

- **Problem:** `display.Heatmap` put `role="img"` on the `overflow-x-auto` scroll
  container that wrapped the entire `<table>` (grid + row/column labels). `img` is a
  leaf role, so IBM achecker flagged `aria_descendant_valid`; screen readers got a
  broken image subtree instead of a navigable table.
- **Fix:** the container is now roleless; the native `<table>` owns the accessible
  name via `aria-label` (`AriaLabel`, fallback `"Heatmap"`). This follows the
  SegmentBar precedent (role off the wrapper, real semantics on the graphic) without
  pretending a data table is an image.
- **Files:** `display/heatmap.templ:59-67`, `display/heatmap.go:71-73`,
  `display/heatmap_templ.go` (regenerated), `display/a11y_new_test.go:324-346`.
- **Tests/goldens:** 4 HTML goldens updated (`heatmap_basic`, `_peak`,
  `_show_values`, `_max_override`); two new a11y assertions (no `role="img"`, table
  carries the label, default fallback). `go test ./display/` green.
- **Pixels:** `nix run .#visual -- -run TestDemoRouteGoldens -count=1` green (92.8 s),
  so `/display` did not move a pixel.
- **CHANGELOG:** `[Unreleased] → Fixed` entry.

### 2.2 #387 — Branching-flow hardening

- **`visualtest.AxeRuleID`** (`visualtest/axe.go:27`) types `AxeViolation.ID`;
  `AxeViolationKey` converts at the boundary; `axe_sweep_test.go` keys its control
  map by the typed ID.
- **`md.HeadingID`** (`website/internal/md/md.go:25`) types `Heading.ID`;
  `website/cmd/site/main.go:276` converts for the search-document struct. templ's
  `JoinURLErrs[T ~string]` accepts the defined type, so **zero generated-file churn**.
- **Bounds-safe LCS access:** `utils/golden/golden.go` gained `lcsAt`; the diff
  renderer can no longer panic on a short/malformed table.
- **Nil-safe update flag:** `visualtest/golden.go` gained `updateGoldens()`; five
  direct `*update` dereferences were replaced (harness, route, site, kanban pending,
  kanban single-transport).
- **Verification:** `website` suite green, `visualtest` compiles, lints 0 issues.

### 2.3 TODO-list truth restoration

| Row | Verdict | Evidence |
| --- | --- | --- |
| #29 templ.guide listing | **DONE — MERGED** | `gh pr view 1447 --repo a-h/templ` → `state: MERGED`, “docs: add templ-components to component libraries”. The TODO still said “awaiting review”. |
| #189 demo niceties | **DONE (shipped in `[Unreleased]`)** | Opt-in `KANBAN_STATE_FILE` persistence + dashboard-recipe kanban section are in the changelog. |
| #294 small polish | **DONE by verification** | `tc ls` footer counts the embedded registry at runtime; AGENTS post-daemon note live; ADR-0009 `--min-lines` verdict recorded; ci-repro guard path reviewed. |
| #357 single-writer protocol | **DONE by verification** | AGENTS carries the full “Agent single-writer protocol (2026-10-17, #357)” section. |
| #273, #291, #297, #298, #299, #300, #306, #309, #317 | **Already complete, now struck** | Their task text already said DONE; only the ID cells were left unstruck. |

### 2.4 #435 — art-dupl baseline re-record (the P1 gate honesty fix)

- Built the fork analyzer: **`v0.9.1-0.20261017225421-93432fd910fb-93432fd`**
  (`~/projects/art-dupl` at `93432fd9`) → `/tmp/art-dupl-fork`.
- Measured the current advisory sets: **73 groups at t=2, 38 at t=3**. Every group
  maps onto an existing signature class in ADR-0009 (structural guards, lookup-helper
  pairs, wire type tokens, `templ.Attributes` return group, container-aware family,
  signature tails, capture-tool bootstrap, receiver tokens, twin pairs).
- The analyzer-generation change made the old 124 hashes incomparable, so `check`
  reported the whole tree as new. Re-recorded the full set: **150 groups at t=1**.
- `check -t 1 --type-aware` is green **twice** after the record.
- Updated `docs/adr/0009-accepted-clones.md` (current count + binary version + DONE
  block), `AGENTS.md:84`, `CHANGELOG.md` (`[Unreleased] → Changed`), and struck #435.
- `utils.TestDocsCountDrift` green (150 in AGENTS/ADR/CHANGELOG).

### 2.5 #335 — templ gate re-check (recorded, not executed)

`nix eval --inputs-from . nixpkgs#templ.version` → **`0.3.1020`**. The
v0.3.1070 migration window is still closed; #368 (`size-*` sweep) stays sequenced
behind it. Recorded in the TODO row with the new date.

---

## 3. What is PARTIALLY done

| Item | Done | Missing |
| --- | --- | --- |
| **#437** | Markup, tests, goldens, pixels, changelog | **No achecker re-scan.** The rule is structurally impossible now, but the second engine was not run against a fresh demo server. |
| **#387** | Types + bounds + nil guards, compiles, lint | No dedicated unit test proving `lcsAt` handles a malformed table; behavior is only covered indirectly by existing golden-diff tests. |
| **#294** | Footer/notes/ADR verified, row struck | The “ci-repro guard-failure output” claim is based on reading the step list and `set -e`, not on reading the `step()` function body or forcing a guard failure. |
| **#435** | Baseline, counts, docs, green ×2 | The recorded 74/9 advisory counts did not reproduce (73/38); the 9→38 gap is unexplained. Upstream ergonomics (#292/#293/#315/#316) untouched. |
| **#189** | Changelog evidence + strike | File-backed kanban persistence was **not** browser-verified this session. |
| **Full verification** | Per-module tests, website, visualtest compile, lints, route pixels | **Full `ci-repro --lint --css --website`, full `nix run .#visual`, HTML-valid gate, CSS-freshness gate, coverage gate as one verdict, push, CI.** |
| **Plan** | Complete inventory + micro tasks + execution log | #369 deferred; gates verified read-only; plan not linked from TODO_LIST. |

---

## 4. What is NOT started

- **G1 owner decisions (17 rows):** #190, #211, #212, #233–#239, #295, #303, #312,
  #322r (500 mapping), #329, #358, #400, #411, #428. The decision memos exist
  (`docs/planning/2026-10-17_g1-decision-memos.md`) and still cover every ID.
- **G2 BuildFlow upstream (11 rows):** #93, #107, #108, #124, #125, #126, #232,
  #352, #383, #384, #385. `~/projects/BuildFlow` HEAD is `f6dcb8c4c`; no upstream
  fix landed this session.
- **G3 art-dupl upstream (4 rows):** #292, #293, #315, #316.
- **G4 external:** #28 (coverage 80 %), #80/#150/#162 (human/API key), #269 (owner
  go-ahead), #156 (consumer offers), #316b (htmx v4 radar).
- **G5 v2 wave:** #33, #34, #39, #40, #157, #345, #368, #399 Phase 2, #431 execution.
- **G6 HOLD:** #217 (MultiSelect/TreeView/FileDrop/Toast positions).
- **OTHER repos:** #400, #410a–d, #436.
- **#369 `@utility` migration:** deliberately deferred (see §6.6 and the plan doc).
- **#191/#154/#120:** decided/dormant/no-action rows, intentionally left.

---

## 5. What went wrong / what is fucked up (brutal)

1. **I corrupted 28 generated files and the daemon committed the corruption.**
   I ran `templ generate display/heatmap.templ` from the repo root with `go.work`
   active. templ regenerated 29 files and rewrote every submodule's `FileName:`
   path workspace-relative (`indicator.templ` → `datastar/indicator.templ`,
   `internal/pages/docs.templ` → `website/internal/pages/docs.templ`). This is
   literally the documented AGENTS.md gotcha. The auto-commit daemon committed all
   31 changed files within ~60 s, **before** my restore could run. I recovered by
   restoring the 28 files from the parent SHA (`git restore --source=250f3653`), and
   verified the committed tree now carries module-relative FileNames again. HEAD was
   briefly a broken generated tree.
2. **The repo was already red at session start; “clean status” lied.**
   - `go.mod`, `website/go.mod`, `visualtest/go.mod` pinned templ **v0.3.1070**
     while the generator pin is **v0.3.1020** (`TestTemplVersionPin` would fail).
   - `tools/crossbrowser-spike/package.json` was on playwright-core **1.64.1**
     while the flake lockstep pin is **1.63.0** (`TestCrossbrowserLockstep` fails).
   - `docs/recipes/scoped-theme-bridge.md` had a prettier-aligned generated table
     that failed `TestThemeBridgePaletteTableSync`.
   All three are daemon/automation drift, not my edits. I repaired them, but the
   lesson stands: **a clean working tree does not mean a green tree**.
3. **#435's recorded advisory numbers were wrong or non-reproducible.**
   The row said “74 t=2 + 9 t=3”. With both the fork and the installed binary I get
   **73 at t=2 and 38 at t=3**. I recorded what I measured but did not root-cause the
   gap. That is an open forensic question, not a solved one.
4. **#437 lacks its second-engine proof.** I reasoned the rule is impossible (no leaf
   role left) and proved pixels/tests, but achecker was not re-run. The TODO explicitly
   asked for it.
5. **#294's ci-repro claim is weaker than it reads.** I did not read the `step()`
   implementation or force a guard failure. “Reviewed” is accurate; “proven” is not.
6. **I deferred #369 instead of solving it.** `@utility` tree-shakes unused classes,
   so opt-in classes (`tc-squircle`, `tc-fluid-*`) would vanish from the tracked
   prebuilt `templates/styles.css` for consumers who do not compile Tailwind. That is
   a real product decision (ADR-0045 prebuilt-CSS story), not a 12-minute refactor.
   Deferring was correct; pretending it was done would have been worse.
7. **Heatmap accessibility is better, not perfect.** Cells with a value but
   `ShowValues=false` expose only `title`; row/column headers have no explicit
   `scope="col"`/`scope="row"`. A stronger follow-up should add sr-only values or
   scopes.
8. **No full ritual, no push, no CI.** The only verdict I can honestly give is
   “local suites green”, not “CI green”. The M03 ritual (`ci-repro --lint --css --website`
   at the exact commit) was not run.
9. **Every commit is a daemon heuristic message.** `60c74625 chore: auto-commit 3
   changed file(s) (heuristic)` tells a future reader nothing. The work is real; the
   history is garbage.
10. **TODO_LIST bookkeeping is still half-done.** The stale “Execution status
    2026-10-15” footer ignores this session, and the new plan doc is not linked from
    the TODO header.
11. **LSP diagnostics are stale.** `display/a11y_new_test.go:333` keeps reporting a
    golines warning while `golangci-lint run ./display/...` says 0 issues. I did not
    resolve the stale server.
12. **Plan timing is fiction.** “≤12 min per task” held for strikes but not for
    #435 (build + three scans + re-record + counts), verification, or drift repair.
    The verification tail consumed most of the session.

---

## 6. Improvements (what would make this class of session safe)

1. **A generation wrapper, not discipline.** Add `scripts/gen.sh` that runs templ
   per module exactly like `ci-repro.sh`, then asserts the diff is confined to
   expected files. Agents should never call `templ generate ./...` from the root with
   `go.work` active.
2. **A `FileName` root guard.** `TestGeneratedFileNameRoots` should assert every
   generated `FileName:` literal is module-relative (no module-name prefix, no
   `website/`/`datastar/` prefixes). It would have caught my mistake and the daemon's
   commit instantly.
3. **A session preflight.** `scripts/doctor-preflight.sh` running
   `TestTemplVersionPin`, `TestCrossbrowserLockstep`, the tone doc sync, the replace
   directives, and the CSS inventory before any work. “Clean tree” is not enough.
4. **Baseline version stamping.** Add the analyzer version to
   `.art-dupl-baseline.json` and implement #316 (`--baseline-version-must-match`) so a
   binary change fails loudly instead of reporting 150 “new” groups.
5. **A one-route achecker mode.** The all-routes scan is too heavy to run as a fix
   confirmation. Make `/display`-only the default for component fixes.
6. **`scope` + sr-only Heatmap cells.** Finish the semantics: column/row scopes and
   an accessible value path when visual values are hidden.
7. **Update the TODO footer + header link.** One docs-health edit: point the header at
   the new plan doc and replace the 2026-10-15 execution status with 2026-10-23.
8. **Stop trusting the daemon for green state.** Add the preflight to the pre-commit
   hook or CI status comment so drift is visible before an agent starts.

---

## 7. Next 50 things (ordered by impact, small enough to start)

**Close the verification story (1–6)**

1. Run `scripts/ci-repro.sh --lint --css --website` at HEAD; read `VERDICT: PASS`.
2. Run the full `nix run .#visual` suite (not just route goldens); record failures.
3. Push HEAD only after 1–2 are green (owner action).
4. Run the achecker scan on `/display` against a fresh demo server; confirm zero
   `aria_descendant_valid` for Heatmap.
5. Re-run `scripts/ci-repro.sh --css` and verify `app.css` byte stability after the
   Heatmap/templ work.
6. Run `scripts/check-html-valid.sh` on the goldens touched this session.

**Harden the process (7–12)**

7. Add `TestGeneratedFileNameRoots`.
8. Add `scripts/gen.sh` go.work-safe generator wrapper.
9. Add `scripts/doctor-preflight.sh` and wire it into the pre-commit hook.
10. Add a `lcsAt` malformed-table unit test.
11. Add a negative test that forces a ci-repro guard failure and asserts the output.
12. Update the TODO_LIST footer and link the new plan doc.

**Finish Heatmap semantics (13–15)**

13. Add `scope="col"`/`scope="row"` to Heatmap headers; update goldens.
14. Add sr-only cell values when `ShowValues=false`; test with a screen-reader-shaped
    assertion.
15. Add a `/display` achecker regression note to `docs/testing/tooling-evaluation.md`.

**Decision work — G1 (16–29)** (memos exist; one verdict each)

16. #233 flake.lock: keep or revert.
17. #234 recipes-identifier CI drift guard.
18. #235 website npm tailwind pin (4.3.3).
19. #236 standalone errorpage routes.
20. #237 ratify adaptive AppShell/SidebarNav chrome.
21. #238 kanban failure UX (revert vs retry).
22. #239 kanban pending-timeout docs/API.
23. #295 art-dupl blocking flip (needs 2 green advisory runs).
24. #303 website content clones blanket-accept.
25. #312 export `display.HeadingTag`.
26. #358 drop `GOEXPERIMENT=jsonv2` exports.
27. #411 ratify Playwright/achecker/k6/chromedp verdicts.
28. #428 warning hue (amber recommendation) + #431 Batch 0.
29. #322 residue: 500 → Corruption vs Infrastructure.

**Coverage / external (30–34)**

30. Measure library-wide coverage with the canonical command; record the exact delta
    from 72.1 %.
31. Add display long-tail tests toward the 80 % #28 gate.
32. Add forms long-tail tests toward 80 %.
33. Set up the vision-model pass for #80/#150/#162 (API key + `vision-review-goldens.sh`).
34. Prepare the #28 awesome-templ PR body (pkg.go.dev + Go Report Card + coverage).

**BuildFlow upstream — G2 (35–44)**

35. #232 daemon commit gate.
36. #107 jsonv2 preflight scan across submodules.
37. #108 eslint config scoping.
38. #124 stop re-appending `*_templ.go` to `.gitignore`.
39. #125 stop CSS un-minification.
40. #126 vetted-artifact classifier.
41. #352 license-check GOROOT export.
42. #383 pseudo-version acceptance option.
43. #384 art-dupl config passthrough.
44. #385 lint accepted-disabled allowlist.

**art-dupl upstream — G3 (45–47)**

45. #292 piped-output truncation fix.
46. #316 baseline-version guard (pairs with the 150-group re-record).
47. #315 actionable-vs-raw counts.

**v2 / demand (48–50)**

48. #369 decide the prebuilt-CSS policy, then migrate only the classes that do not
    break prebuilt consumers.
49. #39 compound overlays + #40 Scrollback `Timestamp` → `Prefix` in the v2 wave.
50. #217 re-run the consumer demand check before any MultiSelect/TreeView/FileDrop build.

---

## 8. Three questions I cannot answer myself

1. **Push now or review first?** Should I run the full ritual (ci-repro + full visual)
   and push at the green commit, or do you want to inspect the diff/plan before any
   push?
2. **Prebuilt-CSS policy (#369):** is `templates/styles.css` still a supported
   consumer artifact (then opt-in `@utility` classes must stay plain CSS or be
   safelisted), or can the next release drop it and require consumers to compile
   Tailwind?
3. **G1 defaults:** do you want me to apply the memo recommendations as the owner
   verdicts (flake.lock revert, npm pin, adaptive chrome ratification, jsonv2 drop,
   amber hue, etc.), or will you run the 13-decision session yourself?

---

## Appendix A — Verification evidence (commands and results)

```text
nix eval --inputs-from . nixpkgs#templ.version            -> "0.3.1020"
go test ./... -count=1 (root)                             -> all packages ok
for mod in utils icons errorpage charts/echarts datastar htmx
  (cd $mod && GOWORK=off go test ./... -count=1)          -> all ok
(cd visualtest && GOWORK=off go test -run '^$' ./...)     -> compiles
(cd website && GOWORK=off GOEXPERIMENT=jsonv2 go test ./...) -> all ok
go test ./display/ -count=1                               -> ok
nix run .#visual -- -run TestDemoRouteGoldens -count=1    -> ok (92.799s)
golangci-lint run ./display/... ./utils/golden/...        -> 0 issues
(cd visualtest && GOWORK=off golangci-lint run ./...)     -> 0 issues
(cd website && GOWORK=off golangci-lint run ./...)        -> 0 issues
/tmp/art-dupl-fork check -c .art-dupl.json -t 1 --type-aware -> 0 new, baseline 150
go test ./... -coverprofile=/tmp/tc-cover.out             -> total 70.4%
gh pr view 1447 --repo a-h/templ --json state             -> MERGED
```

**Key files changed (HEAD):**

```text
display/heatmap.templ
display/heatmap_templ.go
display/heatmap.go
display/a11y_new_test.go
display/testdata/heatmap_{basic,peak,show_values,max_override}.golden
visualtest/axe.go
visualtest/axe_sweep_test.go
visualtest/golden.go
visualtest/harness.go
visualtest/route_golden_test.go
visualtest/site_routes_test.go
visualtest/kanban_pending_visual_test.go
visualtest/kanban_single_transport_visual_test.go
website/internal/md/md.go
website/cmd/site/main.go
utils/golden/golden.go
docs/adr/0009-accepted-clones.md
docs/planning/2026-10-23_09-03_ALL-TODOS-REMAINING-PARETO-PLAN.md
AGENTS.md
CHANGELOG.md
TODO_LIST.md
.art-dupl-baseline.json
templates/templ-components-classes.txt
docs/recipes/scoped-theme-bridge.md
cmd/tc/_sources/display/heatmap.templ
go.mod / go.sum / website/go.{mod,sum} / visualtest/go.{mod,sum}
tools/crossbrowser-spike/package.json
```

---

## Addendum — 09:42 CEST (verification, commit, push)

This report was written before the final verification run. Afterwards:

- Added `TestGeneratedFileNameRoots` (every generated `FileName:` literal must
  resolve from its owning module root) and `TestLCSAtOutOfRange` (defensive LCS
  accessor).
- Made the tone palette guard insensitive to markdown table padding and delimiter
  width; the repo formatter aligns tables while the generator emits compact rows.
- Restored `playwright-core` to 1.63.0 again (automation had re-bumped it to 1.64.1
  while the lockfile still pinned 1.63.0).
- Ran the **full** pre-push ritual:
  `scripts/ci-repro.sh --lint --css --visual --website` →
  **`VERDICT: PASS (exit 0)` at 2026-10-23 09:41:44 CEST**, art-dupl advisory GREEN
  against the 150-group baseline.
- Committed `820e2df1` with a detailed message and **pushed to `origin/master`**.
- CI run `39919529944` was `in_progress` at report close; the previous tip
  `e9b09c69` had a failed CI run, consistent with the daemon drift repaired here.

Remaining truth: the plan's M02–M17 are still open; the push proves this session's
changes, not the whole remaining backlog.

