# Status report — SKILL.md audit and progressive-disclosure split

**Date:** 2026-10-27 20:21 CET
**Session scope:** audit, verify, restructure, and repair the templ-components authoring skill (`skill/SKILL.md` + new `skill/references/`)
**Trigger:** operator question — "Is our SKILLS.md superb!!?!?!"
**Report type:** point-in-time snapshot (Markdown per explicit operator path request; the status-report skill's HTML canonical form was overridden by instruction)

---

## Executive summary

The skill was **rich but wrong in places**. Its guarded numbers (136 components, 102 icons, per-package
headings) were correct because `utils.TestDocsCountDrift` enforces them, but the prose around those numbers
had drifted: six live components were missing from the tables, the entire `charts/echarts` module was
undocumented, four build/test commands were wrong in ways that would mislead an agent, one code example
could not compile, and the file was 964 lines (skill-creator's guidance: under 500).

The skill is now a progressive-disclosure bundle:

| File | Lines | Role |
| --- | ---: | --- |
| `skill/SKILL.md` | 499 | Entry point: catalogue index, consumer essentials, authoring rules |
| `skill/references/component-catalogue.md` | 259 | Every signature + one-liner, icons, `utils`/`format`/`wire` |
| `skill/references/consumer-guide.md` | 183 | Quick start, transports, theming, CSS artifact rule, recipes |
| `skill/references/authoring-guide.md` | 113 | Shared helpers, test checklist, decision trees, site registration |
| `skill/references/guard-tests.md` | 84 | Full guard matrix: what each guard prevents, where it lives |

A pre-existing repo breakage was found and repaired: the auto-commit daemon had bumped root `go.mod` to
`templ v0.3.1070` on 2026-10-25 (commit `ef320b2f`), leaving `TestTemplVersionPin` red. The pin is restored
to `v0.3.1020` (with matching `go.sum` hashes) and the full `utils` suite is green.

---

## Scoreboard

| Category | Count | Notes |
| --- | ---: | --- |
| Defects found and fixed in the skill | 14 | Listed in the findings ledger below |
| Components added to the catalogue | 8 | 6 missing rows + both `charts/echarts` components |
| Factual/command corrections | 10 | Test scope, lint loop, hook path, file counts, verify scope, example, helper names |
| New reference files | 4 | Catalogue, consumer guide, authoring guide, guard matrix |
| Skill line count | 964 → 499 | −48%, under the 500-line guidance |
| Pre-existing repo failures fixed | 1 | `templ` pin drift (guard red at session start) |
| Verifications run | 8 distinct lanes | See verification log |
| Verifications NOT run | 5 lanes | `verify`, visual, website, ci-repro, crossbrowser |
| Uncommitted at report time | 4 files + this report | Daemon committed earlier chunks |

---

# a) FULLY DONE

1. **Audited every claim against the code.** Counted exported templ functions per package (`rg` over
   `*.templ`), read the actual APIs (`utils`, `icons`, `utils/wire`, `utils/format`), checked guard-test
   names against the repo, checked every referenced path on disk. Result: 79 referenced paths, 46 referenced
   test names, zero misses.
2. **Completed the component catalogue.** Added the six missing components — `display.StatusBadgeWith`,
   `display.DateRange`, `display.SectionHeading`, `forms.DirtyGuard`, `feedback.CircularProgress`,
   `htmx.ViewTransitions` — and created the missing `charts/echarts` section (`EChart`, `SDKScript`).
   Every per-package list now matches the sources exactly: 52/26/14/11/12/9/6/2.
3. **Fixed the icons reference.** Header claimed "+ 3 functions" while the table listed nine; now 14 exported
   functions including the RTL variants, `AllIconNames`, and `NameIsValid`.
4. **Corrected the lint instructions.** The old root command passed `./htmx/...` and `./datastar/...` (separate
   modules) to `golangci-lint` and the submodule loop skipped both. Replaced with the canonical root package
   set plus all six submodules, matching `flake.nix` and `AGENTS.md`.
5. **Corrected the test-scope claim.** Old text implied `nix run .#test` covered all modules. Verified
   empirically: `go list ./...` from the root expands to 10 root packages and **zero** submodule packages, so
   the per-module `GOWORK=off` loop is the only complete local form. Documented both lanes honestly.
6. **Corrected the `verify` description.** `nix run .#verify` is the complete lane (root test + per-module
   loop + visualtest + per-module lint), not a root-scoped test step.
7. **Corrected the pre-commit hook story.** The repo uses the tracked `.githooks/pre-commit` via
   `git config core.hooksPath .githooks`; `buildflow precommit install` writes to the dead `.git/hooks/`.
   Documented all ten guards the hook actually runs.
8. **Corrected stale numeric/build claims.** 51 generated files → 137; removed the bogus "system templ may be
   v0.3.1036 / unreleased version" guidance in favour of the real pin (`v0.3.1020`, guarded by
   `TestTemplVersionPin`).
9. **Fixed a non-compiling example.** `RightItems: @layout.ThemeToggle(...)` inside a struct literal is invalid
   templ syntax; now `RightItems: layout.ThemeToggle(...)` with an explicit warning about `@`.
10. **Fixed stale helper names** copied from `AGENTS.md`: `lookupFeedbackStyle` / `feedbackIconName` /
    `dismissScript` do not exist. Replaced with `feedbackStyleMap` / `feedbackStyle` / `feedbackIcon` /
    `feedbackRole` and `utils.DismissButton`.
11. **Split the skill into a progressive-disclosure bundle** (499-line entry point + four references), and
    added a routing table plus per-reference pointers so the depth is discoverable.
12. **Fixed the installation instructions.** The old one-line symlink left `references/` dangling. Instructions
    now link `SKILL.md` **and** `references/`; the live fan-out was updated the same way and the
    crush-config guard reports 50 clean links.
13. **Added missing authoring rules** that the skill previously omitted: `IsValid()` + test in the same commit,
    count updates in the same commit, Go-directive lockstep, `cmd/tc/_sources` mirror handling, demo nonce
    constant, CSS recompile after class changes, templ-explicit-import trap, stale-LSP guidance, and the
    auto-commit daemon hazard.
14. **Restored the templ pin.** `go.mod` → `v0.3.1020`, `go.sum` hashes restored from `utils/go.sum`; root
    `go build ./...` green; full `utils` suite green (with the daemon committing the restore in `9c196443`).
15. **Kept `[Unreleased]` warm** with a detailed CHANGELOG entry, then corrected two inaccuracies in that
    entry during self-review (line count, "seven" → eight components).

---

# b) PARTIALLY DONE

1. **Catalogue one-liners are spot-checked, not systematically verified.** I verified `CopyButton.Href`,
   `Tooltip` position values, and `ListNoteCount`; the other ~130 rows were carried over from the old file.
   `Grid`'s row was misleading ("default `true` since v2.0" — only the constructor defaults true, the struct
   literal zero value is viewport) and was corrected during self-review, which proves the class exists.
2. **Reference integrity is checked by an ad-hoc Python script, not a committed guard.** Paths, test names,
   package lists, and fence balance were checked this session only. Nothing prevents the next drift.
3. **The skill's own done-standard was not met by this session.** The skill says "`nix run .#verify` is the
   complete lane"; I ran a narrower set (utils suite + root build + targeted guards), not the full lane.
4. **`docs/count-definitions.md`, `FEATURES.md`, `README.md` were only checked, not updated.** No count changed,
   so no edit was needed, but the skill is now a fourth description of the same surface without a cross-link
   from those docs.
5. **Consumer/author duplication remains.** The minimal page appears twice (SKILL.md and consumer-guide.md) with
   identical content. Deliberate (entry-point convenience), but it is still a second place to edit.
6. **The CHANGELOG entry for the pin restore does not exist.** The daemon committed `go.mod`/`go.sum` with a
   generic message (`9c196443`); the fix is real but unrecorded in `[Unreleased]`.
7. **`AGENTS.md` still carries the stale feedback-helper sentence** (`feedbackStyleSet + lookupFeedbackStyle[T]()`
   + `feedbackIconName()` + `dismissScript()`) that I corrected in the skill. Not fixed because it is outside
   the stated scope, but it is now a known split brain.

---

# c) NOT STARTED

1. **Skill evaluation loop.** No `evals/evals.json`, no with-skill/without-skill runs, no trigger-accuracy
   measurement for the new description. The description is a single unvalidated hypothesis.
2. **Machine guard for the skill bundle.** No `TestSkillReferencesResolve`, no extension of
   `TestDocsCountDrift` to `skill/references/`, no guard that named tests still exist.
3. **Compile test for the skill's templ snippets.** Both examples are eyeballed, not built.
4. **HARVEST of section (f) into `TODO_LIST.md` / `ROADMAP.md`.** Deliberately deferred — the operator asked
   to wait for instructions after the report.
5. **Full CI reproduction** (`scripts/ci-repro.sh --lint --css --website`) and the visual suite.
6. **Token-budget measurement** for the skill bundle (what it costs to load, entry point vs references).
7. **Verification that `templ generate ./...` from the repo root regenerates submodule templates correctly**
   (documented per `AGENTS.md`, not re-proven this session).
8. **Cross-linking from `README.md` / `CONTRIBUTING.md`** to the skill bundle as the deep reference.

---

# d) TOTALLY FUCKED UP!

Nothing catastrophic happened in this session. The honest version of "fucked up" here is the debt I exposed
and the process traps that remain:

1. **A non-compiling example had been shipping in the skill.** `RightItems: @layout.ThemeToggle(...)` cannot
   compile; it survived because nobody built the snippet. This is the clearest proof that the skill had no
   executable verification, and it is the single most embarrassing finding.
2. **Four commands in the skill were wrong in ways that actively mislead an agent:** root-only test treated as
   all-module, a lint command that lints non-existent root packages while skipping two real modules, the dead
   `.git/hooks` path, and `verify` mis-scoped. An agent following them would report "verified" on partial work.
3. **The repo was red at session start and nobody knew.** The daemon's 2026-10-25 templ bump left
   `TestTemplVersionPin` failing. The bump is a documented recurring adversary; the guard is post-hoc only.
4. **My own CHANGELOG entry was wrong twice** (498 vs 499 lines; "seven" vs eight components). Caught in
   self-review, but it is the same hand-typed-count class the library fights everywhere else.
5. **Work is still uncommitted at report time** (daemon may sweep it): `skill/SKILL.md`,
   `skill/references/component-catalogue.md`, `skill/references/guard-tests.md`, `CHANGELOG.md`, plus this
   report. The daemon has already reverted/consumed session work before; uncommitted = at risk.

---

# e) WHAT WE SHOULD IMPROVE!

1. **Treat the skill as a tested artifact.** Every code block in it should compile or be explicitly marked
   as illustrative. The cheapest fix: a `skill/verify` Go test package that writes each snippet to a temp
   module and runs `go build`.
2. **Automate the integrity checks I ran by hand.** Paths, test names, per-package lists, fence balance, and
   reference reachability belong in one script + flake app, wired into pre-commit and CI.
3. **Extend `TestDocsCountDrift` to `skill/references/`.** Right now the entry point is guarded but the
   catalogue (the actual source of signatures and per-package counts) is not.
4. **Stop hand-typing prose counts anywhere** — including CHANGELOG entries. Prefer "under 500 lines",
   "eight components", or generated counts.
5. **Set a token/size budget for the skill bundle** and enforce it, the way `AGENTS.md` has a 220-line gate.
6. **Decide the canonical home for signatures.** Today: skill references, README catalogue, pkg.go.dev, and
   godoc all describe the same surface. Pick one source and make the others link.
7. **Fix the `AGENTS.md` split brain** found during cross-checking (feedback helper names), and audit the rest
   of AGENTS.md for the same class.
8. **Make the templ-pin bump impossible to land.** The daemon re-applies it; a pre-commit BuildFlow exclusion
   or a flake-level pin is the durable fix (TODO #335 migration plan).
9. **Commit explicitly at the end of agent sessions that produce reviewable artifacts**, rather than relying
   on the daemon's heuristic messages.
10. **Verify the full lane for docs changes at least once.** "Docs-only" still means the CI run must be green;
    this session verified narrower than the skill it was editing claims.

---

# f) Up to 50 things we should get done next

Ordered by impact/effort within groups. P1 = do next, P2 = this week, P3 = when convenient.

## Skill verification (highest value)

| # | Pri | Item |
| ---: | --- | --- |
| 1 | P1 | Add a compile test that builds every templ/Go snippet in `skill/` (temp module + `go build`); start with the two examples in SKILL.md and consumer-guide. |
| 2 | P1 | Extend `TestDocsCountDrift` to read `skill/references/component-catalogue.md` per-package counts and icon count. |
| 3 | P1 | Add `TestSkillReferencesResolve`: every `references/...` link in SKILL.md opens, and every file in `references/` is referenced. |
| 4 | P1 | Add `TestSkillNamedTestsExist`: every backticked `TestXxx` in the skill files resolves to a real test function. |
| 5 | P1 | Add `TestSkillPathsExist`: every backticked repo path in the skill files exists. |
| 6 | P2 | Commit the ad-hoc Python integrity script as `scripts/check-skill-integrity.py` + flake app. |
| 7 | P2 | Wire that script into `.githooks/pre-commit` and `.github/workflows/ci.yaml`. |
| 8 | P2 | Run the full `nix run .#verify` at the final HEAD and record the verdict in a follow-up status note. |
| 9 | P2 | Run `scripts/ci-repro.sh --lint --css --website` once at this HEAD. |
| 10 | P2 | Run the website test suite (`cd website && GOWORK=off GOEXPERIMENT=jsonv2 go test ./...`) to confirm nothing reads the old skill layout. |

## Skill content depth

| # | Pri | Item |
| ---: | --- | --- |
| 11 | P1 | Systematically verify the ~130 catalogue one-liners against godoc/source; record which were checked. |
| 12 | P1 | Add a `wire.Action` cheat sheet (transport × target × swap × event) — the highest-frequency consumer confusion. |
| 13 | P2 | Add an "HTMX vs Datastar: which do I pick?" decision page. |
| 14 | P2 | Document the demo CSP equality rule (`firebase.json` ⇄ `withSecurityHeaders`) in the authoring guide. |
| 15 | P2 | Add the HTML-validation gate workflow (vnu ignore list, local vs CI divergence) to the authoring guide. |
| 16 | P2 | Summarise the one-commit release convention inside the skill instead of only linking `docs/release-checklist.md`. |
| 17 | P2 | Add `nix run .#visual-update` and `nix run .#ogshot` to the flake-command table (both exist, neither is listed). |
| 18 | P2 | Verify how `nix run .#build` regenerates submodule templates, then document the exact behaviour. |
| 19 | P2 | Add a "was this golden update legitimate?" checklist (theme pinned, browser env, load flakes). |
| 20 | P3 | Decide whether the minimal page lives once (SKILL) or twice (SKILL + consumer guide). |
| 21 | P3 | Add a short "common consumer mistakes from real outages" section (stale compiled CSS, `nonce=""`, missing package in scan list). |
| 22 | P3 | Add the `utils/tone` color language to the catalogue index (currently one line). |
| 23 | P3 | Add a per-package "what does NOT exist here" note where consumers repeatedly expect it (e.g. no charts in `display` beyond Tier 1). |
| 24 | P3 | Add a one-page adoption checklist (CSS compile, nonce, HTMX source, theme script order). |
| 25 | P3 | Document the `ContainerAware` zero-value trap in the catalogue rows themselves (Grid fixed; audit Split/Card/Nav/Form/Pagination/DefinitionGrid/SkeletonCardGrid). |

## Skill process and evaluation

| # | Pri | Item |
| ---: | --- | --- |
| 26 | P1 | Create `skill/evals/evals.json` with 3 realistic prompts (add a component, review a component, find/adopt a component). |
| 27 | P2 | Run with-skill vs without-skill comparisons and grade them (skill-creator loop). |
| 28 | P2 | Measure token cost of SKILL.md alone vs SKILL + references; set a budget. |
| 29 | P2 | Test trigger accuracy with ~20 queries, including should-not-trigger near-misses (Tailwind-only, HTMX-only, other Go UI libs). |
| 30 | P3 | Add a "skill maintenance" note to `CONTRIBUTING.md` pointing at the bundle and its guards. |
| 31 | P3 | Add a lightweight `skill/CHANGELOG`-style note? (Probably no — repo CHANGELOG suffices.) |

## Repo hygiene discovered this session

| # | Pri | Item |
| ---: | --- | --- |
| 32 | P1 | Fix `AGENTS.md`'s stale feedback-helper sentence (`lookupFeedbackStyle` etc. do not exist). |
| 33 | P3 | Document that the sold golden count excludes `website/` and `visualtest/` (raw `git ls-files '*.golden'` = 335; guard scope = 330) so the raw number does not look like drift. |
| 34 | P3 | Clarify `AGENTS.md`'s "7-module workspace" vs "9 modules" wording (7 published + 2 repo-local). |
| 35 | P2 | Add a CHANGELOG entry for the templ-pin restore (the daemon's commit message hides it). |
| 36 | P2 | Investigate why BuildFlow `go-mod-update` re-applies templ `v0.3.1070`; add a permanent block (TODO #335 migration). |
| 37 | P2 | Add a session-start check that fails loudly when any `go.mod` contradicts the templ pin. |
| 38 | P2 | Decide whether `docs/status/` needs an index or is intentionally append-only. |
| 39 | P3 | Consider removing `go.mod` naive `require`s drift: root requires submodules at v1.24.0 with replaces; verify release lockstep after this session. |
| 40 | P3 | Check whether any consumer docs (website content) duplicate the component tables and now drift from the skill catalogue. |

## Bookkeeping and follow-through

| # | Pri | Item |
| ---: | --- | --- |
| 41 | P1 | Commit the remaining session work (or confirm the daemon did) and verify the diff survived. |
| 42 | P1 | Harvest section (f) into `TODO_LIST.md` / `ROADMAP.md` via docs-health HARVEST after operator approval. |
| 43 | P2 | Re-read the skill through the fan-out symlink in a NEW session to confirm it appears in `available_skills` with working references. |
| 44 | P2 | Verify the fan-out install instructions on a machine without the existing directory (fresh-clone test). |
| 45 | P2 | Add the skill bundle to the repo's doc inventory (`docs/count-definitions.md` or `CONTRIBUTING.md`) so future sessions know it exists. |
| 46 | P3 | Add `skill/` to any docs-health VERIFY scope so claim drift is caught like other docs. |
| 47 | P3 | Record how to regenerate the fan-out links in the release checklist (skills version with the repo). |
| 48 | P3 | Check art-dupl relevance: docs are not scanned, but confirm the new references do not trip any textual guard. |
| 49 | P3 | Re-run `TestDocsCountDrift` after the next component add to confirm the new split did not break the guard's assumptions. |
| 50 | P3 | Keep this report's findings ledger as the baseline for the next skill audit; annotate (never rewrite) when resolved. |

---

# g) Questions I cannot answer myself

1. **Canonical home for component signatures.** Should `skill/references/component-catalogue.md` become the
   single source (with `README.md`'s catalogue demoted to a link and `TestDocsCountDrift` extended to the
   reference file), or do you want the README table kept as a second, independently guarded surface?
   I cannot tell which audience path you value more.
2. **Verification bar for docs-only changes.** Is `TestDocsCountDrift` + the touched module's tests enough for
   a skill/docs change, or must every such change also run `nix run .#verify` / `scripts/ci-repro.sh`?
   The skill I just wrote claims the full lane; I did not meet that bar and cannot decide whether the claim or
   the practice should change.
3. **Permanent templ-pin enforcement.** Do you want the daemon's `v0.3.1070` re-bump blocked now at the flake
   level (starting the TODO #335 migration), or is "restore on sight + guard" the accepted policy until v2?
   This decides whether item 36/37 above is a one-hour fix or a planned migration.

---

# Verification log

| Command | Result |
| --- | --- |
| `nix develop -c bash -c 'go list ./... \| wc -l'` | `10` packages; `/utils$` and `/htmx$` counts `0` — proves root-only pattern |
| `wc -l skill/SKILL.md skill/references/*.md` | 499 + 259 + 183 + 113 + 84 |
| `cd utils && GOWORK=off GOEXPERIMENT=jsonv2 go test -run TestDocsCountDrift -count=1 .` | `ok` (repeated 4×) |
| `cd utils && GOWORK=off GOEXPERIMENT=jsonv2 go test ./... -count=1` | First run **FAIL** `TestTemplVersionPin` (pre-existing pin drift); after restore: all `ok` |
| `GOWORK=off GOEXPERIMENT=jsonv2 go build ./...` (root) | no output = success, after pin restore |
| Python integrity scan of all five skill files | 79 paths exist, 46 test names exist, 0 missing, fences balanced |
| Per-package set comparison vs `.templ` sources | exact match 52/26/14/11/12/9/6/2 |
| `bash ~/projects/crush-config/scripts/check-skill-fanout.sh` | `OK: 50 skill symlink(s) ... resolve cleanly` |
| Not run | `nix run .#verify`, `nix run .#visual`, `scripts/ci-repro.sh`, website tests, `nix run .#crossbrowser` |

---

# Findings ledger (what was wrong, what changed)

| # | Defect | Status |
| ---: | --- | --- |
| 1 | `display` table: 49 rows vs 52 claimed (missing StatusBadgeWith, DateRange, SectionHeading) | Fixed |
| 2 | `forms` table: 25 vs 26 (DirtyGuard) | Fixed |
| 3 | `feedback` table: 13 vs 14 (CircularProgress) | Fixed |
| 4 | `htmx` table: 8 vs 9 (ViewTransitions) | Fixed |
| 5 | `charts/echarts` module absent (EChart, SDKScript) | Fixed |
| 6 | Icons table: header "+ 3 functions", body 9 (now 14) | Fixed |
| 7 | Root lint command included `./htmx/... ./datastar/...`; submodule loop omitted both | Fixed |
| 8 | `nix run .#test` implied all modules | Fixed + proof recorded |
| 9 | `nix run .#verify` described as root-scoped | Fixed |
| 10 | Hook path `.git/hooks/pre-commit`; only 4 guards named | Fixed (`.githooks` + 10 guards) |
| 11 | "51 generated files"; "templ v0.3.1036 / unreleased" guidance | Fixed (137; v0.3.1020) |
| 12 | `RightItems: @layout.ThemeToggle(...)` does not compile | Fixed |
| 13 | Stale feedback helper names (did not exist) | Fixed |
| 14 | Catalogue `Grid` row implied struct-literal default true | Fixed during self-review |
| 15 | My CHANGELOG entry: 498 vs 499 lines; "seven" vs eight components | Fixed during self-review |

---

# Files changed this session

| File | Change |
| --- | --- |
| `skill/SKILL.md` | Rewritten: 964 → 499 lines; corrected commands/facts; new rules; references routing |
| `skill/references/component-catalogue.md` | New: complete catalogues (was 224 lines inside SKILL; now 259 standalone, 8 components added) |
| `skill/references/consumer-guide.md` | New: adoption walkthrough extracted from SKILL |
| `skill/references/authoring-guide.md` | New: helpers, test checklist, decision trees, site registration |
| `skill/references/guard-tests.md` | New: guard matrix (was a 25-row table inside SKILL, now curated + expanded) |
| `CHANGELOG.md` | `[Unreleased]` entry for the skill bundle |
| `go.mod` / `go.sum` | `templ v0.3.1070` → `v0.3.1020` restore (daemon had bumped; guard was red) |
| `~/.config/crush/skills/templ-components/references` | New symlink → repo `skill/references` (fan-out install) |
| `docs/status/2026-10-27_20-21_SKILL-MD-AUDIT-AND-PROGRESSIVE-DISCLOSURE-SPLIT.md` | This report |

---

# Risks and open hazards

1. **Daemon writes during reporting.** It already committed four times this session (`5209a272`, `a00620b7`,
   `1ca16e43`, `9c196443`) and may sweep the remaining edits. Verify current `git log`/`git status` before
   trusting this report's file list.
2. **Daemon re-bumps templ.** The pin restore is durable only until the next `go-mod-update` pass.
3. **Guarded vs unguarded counts.** SKILL.md's numbers are enforced; the new catalogue's per-package counts are
   not. A component add can now drift in the reference file silently.
4. **Skill depth now depends on references being read.** If an agent stops at SKILL.md it loses signatures and
   the guard matrix — acceptable by design, but the entry point must keep enough routing to make that likely.
5. **Examples remain unbuilt.** Until item 1 in section (f) lands, another non-compiling snippet can ship.
6. **`AGENTS.md` split brain** on feedback helpers is live until item 32 lands.

---

# Appendix: verified counts at report time

- Published modules: 7 (root, utils, icons, errorpage, charts/echarts, htmx, datastar) + 2 repo-local
  (visualtest, website).
- Exported templ components: **136** across 9 counted packages — display 52, forms 26, feedback 14, layout 11,
  navigation 12, htmx 9, datastar 6, errorpage 4, charts/echarts 2.
- Icons: **102** (101 path + Spinner); the `icons` package exposes 8 of its own templ components, which are
  **not** part of the 136.
- Generated files (`*_templ.go`, website excluded): **137**.
- Golden HTML files: **330** under the guard's scope (`.git`, `website/`, `visualtest/` excluded) — matches
  the `AGENTS.md` claim; the raw tracked total is 335 because the website and visualtest goldens are excluded
  from the sold count.
- Skill bundle: **1138** lines total (SKILL 499 + references 639).
