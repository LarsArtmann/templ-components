# Status Report — Docs Health Audit & Cross-File Count Drift Fixes

**Date:** 2026-07-13 23:14 CEST
**Session scope:** Read all 42 `docs/**/2026-07-0*` files, then execute the docs-health skill (AUDIT mode) — verify every core doc against code, fix all drift.
**Branch:** master (uncommitted — 6 `.md` files changed, +28/-27 lines)
**Version:** 0.17.0

---

## Context

The user asked to (1) read ALL `**/2026-07-0*` files, (2) run the docs-health skill. The 42 files
(5 feedback, 6 planning, 19 status, 4 HTML reviews, 8 remaining) span sessions from 2026-07-04
through 2026-07-09 — versions v0.4.x through v0.11.0. The repo is now at v0.17.0. Every one of
those historical reports has an `AUTO-UPDATED 2026-07-10` overlay noting they're historical.

The docs-health skill calls for: inventory docs → verify each against code → classify findings →
fix drift → check cross-file consistency → report health score.

---

## a) FULLY DONE

### 1. Read and analyzed all 42 `2026-07-0*` files

Read every file completely (not just skimmed). Extracted the common themes:

- **Consumer feedback** (5 files): DiscordSync (2 sessions), browser-history, overview, SwettySwipper. Consistent themes: component discoverability gap, forms package underadopted, CDN dependency hidden, dashboard components (StatCard, Grid, Table) highly valued.
- **Status reports** (19 files): Sessions 6-13 covering feedback-driven improvements, bug hunts (47+ bugs fixed across 2 rounds), dark mode compliance audit (30+ fixes), dedup sprint (6 sub-templates extracted), coverage boost (152 tests), v0.7.0-v0.10.0 releases, CSS integration cleanup (tc-css deleted).
- **Planning docs** (6 files): Pareto breakdowns, hardening plans (v1 + v2), dark mode compliance plan (52 tasks), comprehensive execution plan (89 tasks), SUPERB UI library upgrades.
- **HTML reviews** (4 files): Code quality scan, full code review, data model review, naming review — all Bauhaus Dark dashboard format.

### 2. Verified all core docs against code

| Doc                     | Exists | Verified Against                                                                          | Drift Found?                 |
| ----------------------- | ------ | ----------------------------------------------------------------------------------------- | ---------------------------- |
| README.md               | Yes    | Component counts, icon count, enum count, test count, htmx count                          | **YES — 6 medium issues**    |
| AGENTS.md               | Yes    | htmx component count, icon count, lint paths, package structure                           | **YES — 2 medium issues**    |
| FEATURES.md             | Yes    | Component counts, icon count, enum count, generated files count, htmx/icons detail tables | **YES — 5 medium issues**    |
| TODO_LIST.md            | Yes    | All 52 items spot-checked against code                                                    | **No — best-maintained doc** |
| ROADMAP.md              | Yes    | Component/icon counts, v2.0 items (dialog already shipped)                                | **YES — 2 medium issues**    |
| CHANGELOG.md            | Yes    | `[Unreleased]` has body, heading matches `utils.Version`                                  | **No**                       |
| CONTRIBUTING.md         | Yes    | Lint commands, build commands                                                             | **YES — 1 critical ghost**   |
| skill/SKILL.md          | Yes    | Component count, per-package counts, icon count                                           | **YES — 3 medium issues**    |
| docs/DOMAIN_LANGUAGE.md | Yes    | Domain terms                                                                              | **No**                       |

### 3. Fixed all drift (6 files, +28/-27 lines)

#### Critical fix: CONTRIBUTING.md broken lint commands

| Line | Before                                                                         | After                     | Severity     |
| ---- | ------------------------------------------------------------------------------ | ------------------------- | ------------ |
| 39   | `golangci-lint run ./display/... ./errorpage/... ... ./svg/... ./internal/...` | `golangci-lint run ./...` | **CRITICAL** |
| 45   | Same broken path in full verification command                                  | `golangci-lint run ./...` | **CRITICAL** |

**Root cause:** The modularization branch (`modularize/strategic-split`) promoted `internal/svg`
to `svg/` but was never merged. CONTRIBUTING.md still referenced `./svg/...` — a directory that
doesn't exist. Anyone following the lint instructions would get `lstat ./svg/: no such file or
directory`.

#### Medium fixes: Cross-file count drift

The repo advanced from v0.11.0 (the latest the historical reports cover) to v0.17.0, adding
Popover, ContextMenu, Carousel, HoverCard, DataTable, FilterDropdown, Slider, Rating, TagsInput,
Calendar, IconRTL, ViewTransitions. No doc was updated to reflect these additions.

| Metric                       | Before (stale)        | After (verified) | Files Fixed                              |
| ---------------------------- | --------------------- | ---------------- | ---------------------------------------- |
| Components                   | 84-93 (varied by doc) | **97**           | README, FEATURES, ROADMAP, SKILL         |
| Icons                        | 101                   | **102**          | README, FEATURES, ROADMAP, AGENTS, SKILL |
| htmx components              | 7                     | **8**            | README, FEATURES, AGENTS, SKILL          |
| layout components            | 5 (SKILL only)        | **6**            | SKILL                                    |
| Typed enums                  | 26-36 (varied)        | **37**           | README, FEATURES                         |
| Generated `*_templ.go` files | 74 (FEATURES)         | **75**           | FEATURES                                 |
| Test count                   | ~850 (README)         | **~990**         | README                                   |
| Packages                     | 11 (README)           | **9**            | README                                   |
| icons functions              | 2 (SKILL)             | **3**            | SKILL                                    |

#### Missing component rows added

- **FEATURES.md htmx table**: Added `ViewTransitions` row (was in code since v0.11.0, missing from feature inventory)
- **FEATURES.md icons table**: Added `IconRTL` row (RTL mirror icon, shipped in v0.13.0)

### 4. Verified cross-file consistency

After fixes, all 8 docs agree on every count:

```
Components: README=97  FEATURES=97  ROADMAP=97  SKILL=97
Icons:      README=102 FEATURES=102 ROADMAP=102 AGENTS=102 SKILL=102
htmx:       README=8   AGENTS=8     FEATURES=8   SKILL=8
```

### 5. Drift-guard tests pass

- `TestVersionMatchesChangelog` — PASS (0.17.0 matches)
- `TestVersionMatchesFeatures` — PASS (0.17.0 matches)
- `TestDarkModeCompliance` — PASS (0 violations)
- `TestDarkModeSemanticColors` — PASS (0 violations)
- `TestMotionReduceCompliance` — PASS (0 violations)
- `TestSkillComponentCount` — PASS (informational, 94 actual vs 97 documented — within tolerance)

### 6. Build verification

- `go build ./...` — PASS (with `GOEXPERIMENT=jsonv2`)
- `go test ./...` — ALL 14 packages green (with `GOEXPERIMENT=jsonv2`)
- Without `GOEXPERIMENT=jsonv2`: 4 packages fail at setup (`errorpage`, `integration`, `internal/contract`, `examples/demo`) — this is pre-existing and documented in AGENTS.md, NOT caused by this session's changes.

---

## b) PARTIALLY DONE

### 1. FEATURES.md typed enum count discrepancy

FEATURES.md says "37 typed enums (37 with `IsValid()`)". Actual code has 34 `IsValid` methods in
non-generated `.go` files, 37 total including generated files. The "37 typed enums" count is
correct (12 `type X string` definitions + additional enum types defined inline in various files),
but the parenthetical "(37 with IsValid())" may be slightly off — the `IsValid` count depends on
whether you count generated-file methods.

**What I did:** Updated from 36→37 based on actual `IsValid` count across all `.go` files.
**What's uncertain:** Whether the "typed enums" count itself is 37 or higher — there are enums
defined as struct fields, not just `type X string`. A rigorous count would need to grep every
closed-set string type across the codebase.

### 2. SKILL.md detailed component catalogue not fully audited

I verified the header counts (97 components, 102 icons, 8 htmx, 6 layout) and the icons section
detail, but did not line-by-line verify every component in the SKILL.md catalogue tables against
the codebase. Some components may be missing from the catalogue even though the total count is
correct (e.g., new components shipped in v0.13.0-v0.17.0 may not have catalogue rows).

### 3. Historical report overlays not verified for accuracy

All 42 `2026-07-0*` files have an `AUTO-UPDATED 2026-07-10` overlay claiming specific fixes were
completed. I read these overlays but did not independently verify every claim (e.g., "All 7 P0
bugs fixed", "33 regression tests added", "Motion constants centralized"). The TODO_LIST.md
confirms these as DONE with evidence, so they're likely accurate — but I didn't re-verify from
source.

---

## c) NOT STARTED

### 1. Did not run the full `nix run .#verify` pipeline

AGENTS.md says `nix run .#verify` is the canonical done-check. I used `go build ./...` and
`go test ./...` with `GOEXPERIMENT=jsonv2` directly instead. The result is equivalent, but I
didn't follow the prescribed process.

### 2. Did not add a CHANGELOG `[Unreleased]` entry for the doc fixes

The AGENTS.md convention says "every feature/fix commit that lands on master must add its
changelog entry to the [Unreleased] section immediately." Documentation-only fixes are a gray
area, but the convention suggests I should have added an entry. I did not.

### 3. Did not verify all per-component entries in FEATURES.md

FEATURES.md has detailed per-component tables (Status, Description, Key Features). I updated
counts and added missing rows for ViewTransitions and IconRTL, but did not verify every row's
status claim against the code. Some components marked FULLY_FUNCTIONAL might have edge-case bugs
or missing features that would warrant PARTIALLY_FUNCTIONAL.

### 4. Did not verify all ADR cross-references

AGENTS.md, README.md, CONTRIBUTING.md, and ROADMAP.md all reference ADRs (0007-0015). I did not
verify that every referenced ADR file actually exists at the referenced path. Some ADRs may have
been renamed or moved.

### 5. Did not verify all recipe/doc cross-links

README.md has a "Further reading" table linking to docs (javascript-guide, motion-design,
container-queries recipe, etc.). I did not verify each linked file exists.

### 6. Did not commit the changes

6 files are modified but uncommitted. The user has not said "commit."

---

## d) TOTALLY FUCKED UP

### 1. Read the HTML review files unnecessarily

The 4 HTML review files (`docs/reviews/2026-07-05_18-*.html`) are large self-contained dashboard
reports (~500+ lines each). The user asked me to "READ ALL `**/2026-07-0*` files" so I did — but
these are point-in-time visual reports with embedded CSS. Reading them consumed significant
context window for zero actionable content. They're historical artifacts, not working docs.
I should have noted they exist, confirmed they're HTML reports, and skipped the full read.

**Impact:** Wasted context window capacity that could have been used for deeper verification.

### 2. First `multiedit` on CONTRIBUTING.md failed (1 of 2 edits)

The first attempt to fix both lint commands in CONTRIBUTING.md via `multiedit` only applied 1 of
2 edits. The second edit's `old_string` didn't match because the first edit had already changed
the text. Had to re-read and apply the second fix separately.

**Root cause:** The two lint commands in CONTRIBUTING.md are on different lines (39 and 45) with
different surrounding context, but I tried to match them as if they were adjacent. Should have
used separate `edit` calls from the start.

### 3. Initial component count was wrong (84 → 94 → 97)

I first set README "Components" to 94 based on a quick grep, then had to correct to 97 after
verifying FEATURES.md and ROADMAP.md had a different methodology (including icons' 3 templ
components). Three edits to the same field in 10 minutes — embarrassing.

**Root cause:** I counted `templ [A-Z]` in `.templ` files from the 8 main packages (94) but
forgot the 3 icons components (Icon, IconWithStrokeWidth, IconRTL), which brought the total to 97.
FEATURES.md and ROADMAP.md already used the higher count. I should have reconciled BEFORE editing.

### 4. Did not notice the build was broken without GOEXPERIMENT=jsonv2

When I first ran `go build ./...` and `go test ./...`, 4 packages failed. I initially thought my
changes broke something (they didn't — I only edited `.md` files). It took me a moment to realize
this is the pre-existing `encoding/json/v2` issue documented in AGENTS.md. I should have known
this from reading AGENTS.md at the start of the session.

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Reconcile counts BEFORE editing, not after.** I edited README to 94, then had to fix it to 97. Always grep the actual code AND check what other docs say before writing a number. The
   "single source of truth" approach: compute once, apply everywhere.

2. **Skip HTML report files in future "read all" commands.** The 4 HTML files in `docs/reviews/`
   are self-contained dashboards with no actionable content for docs-health. Future sessions
   should read the `.md` files and note the HTML files exist without full reads.

3. **Run `nix run .#verify` as the canonical check.** I used manual `go build` + `go test` with
   explicit `GOEXPERIMENT=jsonv2`. The Nix verify app handles this automatically. Using it would
   have avoided the moment of panic when 4 packages failed.

4. **Add a count-drift CI check.** The drift-guard tests (`TestVersionMatchesChangelog`,
   `TestVersionMatchesFeatures`) catch version drift. We need equivalent tests for component
   counts, icon counts, and enum counts — so this kind of drift can't accumulate across 6
   versions again. A `TestREADMECountsMatchCode` that greps README "By the Numbers" against
   actual `templ` function count, `icon_names.go` entry count, etc.

5. **CHANGELOG [Unreleased] entry for doc fixes.** The convention is clear: "every change that
   lands on master must add its changelog entry." I skipped this. Doc fixes are changes.

### Documentation system

6. **The historical report overlay is excellent but could be automated.** The `AUTO-UPDATED
2026-07-10` block at the top of every `2026-07-0*` file is a manual annotation. A script that
   reads `TODO_LIST.md` statuses and auto-generates these overlays would keep them accurate as
   TODO items are resolved.

7. **FEATURES.md enum count methodology is ambiguous.** "37 typed enums" — does this mean 37
   `type X string` definitions? 37 enums with `IsValid()`? 37 closed-set string types including
   inline definitions? The count varies depending on methodology. Document the counting method.

8. **SKILL.md component catalogue doesn't list every component.** The header says 97 but the
   per-package tables may not have a row for every one (especially v0.13.0-v0.17.0 additions like
   Popover, ContextMenu, Carousel, HoverCard, Calendar, etc.). The catalogue should be exhaustive.

---

## f) Up to 50 Things We Should Get Done Next

### Immediate (before commit)

| #   | Task                                                   | Effort |
| --- | ------------------------------------------------------ | ------ |
| 1   | Add CHANGELOG `[Unreleased]` entry for doc count fixes | 3m     |
| 2   | Review the 6-file diff for accuracy                    | 5m     |
| 3   | Commit the doc fixes                                   | 2m     |

### Count-drift prevention

| #   | Task                                                                                  | Impact | Effort |
| --- | ------------------------------------------------------------------------------------- | ------ | ------ |
| 4   | Add `TestREADMECountsMatchCode` — verify README "By the Numbers" against actual code  | HIGH   | 15m    |
| 5   | Add `TestFeaturesCountsMatchCode` — verify FEATURES.md totals against actual code     | HIGH   | 15m    |
| 6   | Add `TestRoadmapCountsMatchCode` — verify ROADMAP.md component/icon counts            | MED    | 10m    |
| 7   | Add `TestAgentsCountsMatchCode` — verify AGENTS.md per-package counts                 | MED    | 10m    |
| 8   | Add `TestSKILLCountsMatchCode` — verify SKILL.md per-package counts                   | MED    | 10m    |
| 9   | Document the enum counting methodology in FEATURES.md (what counts as a "typed enum") | LOW    | 5m     |

### SKILL.md catalogue completeness

| #   | Task                                                                                                                                                  | Impact | Effort |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 10  | Audit SKILL.md display section — verify all 30 components have catalogue rows                                                                         | MED    | 15m    |
| 11  | Audit SKILL.md forms section — verify all 21 components have catalogue rows                                                                           | MED    | 10m    |
| 12  | Audit SKILL.md feedback section — verify all 13 components have catalogue rows                                                                        | LOW    | 10m    |
| 13  | Audit SKILL.md navigation section — verify all 12 components have catalogue rows                                                                      | LOW    | 10m    |
| 14  | Audit SKILL.md htmx section — verify all 8 components have catalogue rows                                                                             | LOW    | 5m     |
| 15  | Add new components (Popover, ContextMenu, Carousel, HoverCard, Calendar, DataTable, FilterDropdown, Slider, Rating, TagsInput) to SKILL.md if missing | MED    | 30m    |

### FEATURES.md completeness

| #   | Task                                                                                                                                | Impact | Effort |
| --- | ----------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 16  | Audit FEATURES.md for all v0.13.0-v0.17.0 components (Popover, ContextMenu, etc.)                                                   | MED    | 20m    |
| 17  | Verify every FULLY_FUNCTIONAL claim in FEATURES.md against actual test results                                                      | LOW    | 30m    |
| 18  | Add `templates/app.css` details to FEATURES.md (modern CSS features: `@starting-style`, `field-sizing`, `content-visibility`, etc.) | LOW    | 15m    |

### README.md improvements

| #   | Task                                                                                     | Impact | Effort |
| --- | ---------------------------------------------------------------------------------------- | ------ | ------ |
| 19  | Add ViewTransitions to README htmx code example                                          | LOW    | 3m     |
| 20  | Add IconRTL to README icons section                                                      | LOW    | 3m     |
| 21  | Add modern web standards section to README (dialog, stylable select, field-sizing, etc.) | MED    | 15m    |
| 22  | Update README "Dependencies" count if go-error-family version changed                    | LOW    | 2m     |

### ADR and cross-link verification

| #   | Task                                                                                    | Impact | Effort |
| --- | --------------------------------------------------------------------------------------- | ------ | ------ |
| 23  | Verify all ADR references in docs resolve (ADR 0007-0015 all exist)                     | MED    | 10m    |
| 24  | Verify all recipe doc links in README resolve                                           | LOW    | 5m     |
| 25  | Verify all "Further reading" links in README resolve                                    | LOW    | 5m     |
| 26  | Add ADR entries for v0.11.0-v0.17.0 features (if any architectural decisions were made) | LOW    | 15m    |

### Historical report cleanup

| #   | Task                                                                                      | Impact | Effort |
| --- | ----------------------------------------------------------------------------------------- | ------ | ------ |
| 27  | Consider archiving `2026-07-0*` files to `docs/archive/2026-07/` — they're all historical | LOW    | 10m    |
| 28  | Or add a `docs/status/README.md` index explaining the historical nature of these files    | LOW    | 5m     |
| 29  | Verify the `AUTO-UPDATED 2026-07-10` overlays are still accurate (spot-check 3 claims)    | LOW    | 10m    |

### Domain language

| #   | Task                                                                                                                       | Impact | Effort |
| --- | -------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 30  | Update `docs/DOMAIN_LANGUAGE.md` with terms from v0.11.0-v0.17.0 (View Transitions, Stylable Select, ContainerQuery, etc.) | LOW    | 15m    |
| 31  | Verify DOMAIN_LANGUAGE.md terms are used consistently across all docs                                                      | LOW    | 10m    |

### ROADMAP.md improvements

| #   | Task                                                                                       | Impact | Effort |
| --- | ------------------------------------------------------------------------------------------ | ------ | ------ |
| 32  | Add v0.11.0-v0.17.0 milestones to ROADMAP.md "Current" section                             | LOW    | 10m    |
| 33  | Update ROADMAP.md v2.0 section — "Native `<dialog>`" already shipped, remove from research | LOW    | 2m     |
| 34  | Add "Modern web standards adoption" to ROADMAP.md current status                           | LOW    | 5m     |

### Testing improvements

| #   | Task                                                                    | Impact | Effort |
| --- | ----------------------------------------------------------------------- | ------ | ------ |
| 35  | Add a test that asserts every package directory has a FEATURES.md entry | LOW    | 10m    |
| 36  | Add a test that asserts every component appears in FEATURES.md          | MED    | 20m    |
| 37  | Add a test that asserts every component appears in SKILL.md catalogue   | MED    | 20m    |

### CONTRIBUTING.md improvements

| #   | Task                                                                   | Impact | Effort |
| --- | ---------------------------------------------------------------------- | ------ | ------ |
| 38  | Add `GOEXPERIMENT=jsonv2` to CONTRIBUTING.md build instructions        | MED    | 3m     |
| 39  | Add `nix develop` mention to CONTRIBUTING.md for templ version pinning | LOW    | 5m     |
| 40  | Cross-link CONTRIBUTING.md and AGENTS.md more prominently              | LOW    | 3m     |

### AGENTS.md improvements

| #   | Task                                                                                                       | Impact | Effort |
| --- | ---------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 41  | Add ViewTransitions to AGENTS.md conventions section                                                       | LOW    | 5m     |
| 42  | Add Stylable Select to AGENTS.md conventions section                                                       | LOW    | 5m     |
| 43  | Add modern web standards section to AGENTS.md (dialog migration summary, field-sizing, content-visibility) | MED    | 15m    |
| 44  | Verify "75 generated files" count is current after any new component additions                             | LOW    | 2m     |

### Process improvements

| #   | Task                                                                                                          | Impact | Effort |
| --- | ------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 45  | Create a `scripts/docs-health.sh` that automates this audit (count drift, ghost refs, cross-file consistency) | HIGH   | 45m    |
| 46  | Add docs-health check to CI (run `scripts/docs-health.sh`, fail on drift)                                     | HIGH   | 15m    |
| 47  | Add pre-commit hook that warns when component/icon counts change but docs aren't updated                      | MED    | 30m    |

### Polish

| #   | Task                                                                               | Impact | Effort |
| --- | ---------------------------------------------------------------------------------- | ------ | ------ |
| 48  | Run `go test ./... -count=1` with `GOEXPERIMENT=jsonv2` to get uncached test count | LOW    | 5m     |
| 49  | Update README "By the Numbers" test count to exact (not ~990) after full test run  | LOW    | 2m     |
| 50  | Consider adding a "Last verified" date to README "By the Numbers" section          | LOW    | 2m     |

---

## g) Top 2 Questions I Cannot Answer Myself

### 1. Should the `2026-07-0*` historical files be archived, deleted, or left as-is?

There are 42 files in `docs/status/`, `docs/planning/`, `docs/feedback/`, and `docs/reviews/`
from July 4-9, 2026. They all have the `AUTO-UPDATED 2026-07-10` overlay marking them as
historical. They reference versions v0.4.x-v0.11.0; the repo is now v0.17.0.

**Arguments for archiving** (move to `docs/archive/2026-07/`):

- Reduces noise in `docs/` directory listing
- Makes it clear which docs are current vs historical
- Prevents future sessions from acting on stale "open" items

**Arguments for leaving as-is:**

- The overlay already marks them as historical
- Moving them breaks any external links
- They're useful context for understanding why decisions were made
- The git history preserves them regardless

**Arguments for deleting:**

- The TODO_LIST.md has extracted all actionable items with evidence citations
- The CHANGELOG captures what shipped
- 42 files is a lot of context noise for future sessions

I lean toward **archiving** (move to `docs/archive/2026-07/`) but this is a taste/filing decision.

### 2. Should count-drift tests be failing (block CI) or informational (log only)?

The existing `TestSkillComponentCount` is informational — it logs actual vs documented but
doesn't fail. I propose adding `TestREADMECountsMatchCode`, `TestFeaturesCountsMatchCode`, etc.
that verify component/icon/enum counts match between docs and code.

**Failing tests (block CI):** Forces docs to stay in sync. Any PR adding a component MUST update
all docs in the same commit. Pro: zero drift. Con: adds friction, especially for quick component
additions.

**Informational tests (log only):** Surfaces drift during CI but doesn't block. Pro: low friction.
Con: drift accumulates silently (exactly what happened here — 6 versions of drift).

**Hybrid:** Failing on component count (structural, changes rarely), informational on test count
(grows organically, hard to keep exact).

I lean toward **failing on structural counts** (components, icons, enums, packages) and
**informational on volatile counts** (tests, coverage %). But this is a maintenance philosophy
decision.
