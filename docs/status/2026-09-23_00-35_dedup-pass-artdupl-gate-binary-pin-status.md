# Status Report — Dedup Pass + art-dupl Gate Repair

**Date:** 2026-09-23 00:35 CEST
**Session scope:** Full dedup pass over the three art-dupl reports (t=1/2/3) in templ-components, triggered by "deduplicate!". Includes gate forensics, five extractions, baseline repair, and doc updates.
**Format note:** status-report skill's canonical format is styled HTML; user explicitly requested `.md` — honored as override. A dedicated `brutal-self-review` HTML report was NOT written (user scoped this session to one report + wait); the self-review answers are folded into sections (d)/(e).

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| a1 | **Root-caused the gate failure**: the "105 new clone groups" verdict was a false alarm — baseline hashes are binary-pinned. The 39-entry morning baseline (pre-actionability detector, also counting 575 files incl. `cmd/tc/_sources`) matched 0 hashes under both installed 0.7.0 and fork v0.7.0-74 | Three-binary comparison: installed (105 new), fork build (103 new), paste-1 (39 total); `/tmp/dupl-*.log` series |
| a2 | **Debunked the `--sort` red herring**: proved `--sort total-tokens` does not change detection (848 vs 847 raw groups with/without); my initial 105-vs-848 alarm was comparing a check-verdict against a scan-total — documented as an AGENTS.md gotcha | `/tmp/dupl-plain.log` vs `/tmp/dupl-plain-sorted.log` vs `/tmp/dupl-sort.log` |
| a3 | **Built the user's art-dupl fork** (`~/projects/art-dupl`, v0.7.0-74-ge7456139, go.mod needs 1.27.1 → `GOTOOLCHAIN=auto`), producing `/tmp/art-dupl-fork` used for all gate operations | Build succeeded; `check` ran green against new baseline |
| a4 | **Extraction 1 — `display/sparklinePointCoords`**: duplicated x/y projection loop between `sparklinePoints`/`sparklineAreaPath` now shared (single source of truth); ADR entry 15 ("already extracted" claim) corrected | display/sparkline.go:60; display tests + goldens green; zero sparkline matches in post-pass scans |
| a5 | **Extraction 2 — `display.headingTag`**: three identical 6-case TitleTag switches (Card, EmptyState, CollapsibleSection) collapsed into one sub-template in display/shared.templ; call-site residue accepted in ADR | display/shared.templ; goldens byte-identical; headingTag-switch clones gone from scans |
| a6 | **Extraction 3 — `display.chartSeriesGroup` + `chartSeriesRenderOpts`**: per-series SVG block (fill path, stroke path, dots) duplicated by LineChart + AreaChart now rendered in one place; stroke/dash/dot can no longer drift | display/chart_shared.templ; area/line chart goldens byte-identical |
| a7 | **Extraction 4 — `forms.calendarMonthNavLink`**: mirrored prev/next month-nav anchors (href fallback + spread transport attrs + aria-label + chevron) unified | forms/calendar.templ; calendar goldens + tests green |
| a8 | **Extraction 5 — `visualtest/tools/internal/{browser,distserver}`**: exact `chromePath()` ×3 (ogshot/shots/siteshots) and the Firebase cleanUrls HTTP handler ×2 consolidated; module builds, `go vet` clean, gci import-order fixed | New packages committed; `GOWORK=off go build ./... && go vet ./...` green in visualtest |
| a9 | **Baseline re-recorded with the fork**: 122 actionable groups (548 files, 224 filtered), gate **GREEN twice** — `✅ No new clones detected (baseline: 122 groups)` — verified again after lint fixes and after `nix fmt` | `/tmp/dupl-final-check.log`; re-run at final tip |
| a10 | **ADR-0009 updated**: new evening-pass section (5 extractions + headingTag-residue acceptance), entry 15 status corrected, "deduplication budget is exhausted" doctrine **retracted** (disproved same day), baseline section rewritten with the tool-version pin + re-record ritual | docs/adr/0009-accepted-clones.md (committed in 506dfeaf) |
| a11 | **AGENTS.md art-dupl paragraph rewritten**: 122-count, binary-pin gotcha, "Found total = not-in-baseline count" trap, fork build command | AGENTS.md (committed in c546c4e9) |
| a12 | **All test lanes green**: all 6 sub-modules + visualtest per-module loop OK; root module 10/10 packages ok; display/forms/utils drift guards (dark-mode, motion-reduce, coarse-pointer, inline-icon, templ-sync) pass | `/tmp/root-tests.log`, per-module loop output |
| a13 | **Lint clean**: display `makezero` + 2× `varnamelen` fixed (sparkline append-order + `coord` renames); visualtest 3× `gci` fixed via `golangci-lint fmt`; final re-lint 0 issues | golangci-lint output after fixes |
| a14 | **`_sources` mirror synced**: guard detected drift, `--fix` re-copied 8 files, recheck green; replace-directives guard green; gofmt clean; `nix fmt` applied | scripts/check-tc-sources-sync.sh output |
| a15 | **Daemon-commit integrity verified**: work survived 5 heuristic auto-commits (14251e80, cc345ed6, 506dfeaf, 9c1db38e, c546c4e9); new internal packages tracked; build+gate+tests re-verified at tip | git log --stat inspection |

Net code effect: −240 lines pre-lint-fixes (846+/1086−), five clone families eliminated from scans (verified by grep of post-pass t=1 AND t=2 logs: sparkline, xCoord, chromePath, CHROMEDP_CHROME_PATH, TitleTag switches, calendar prev-nav pattern — all 0 matches).

## b) PARTIALLY DONE

| # | Item | What works | What remains | Effort |
|---|------|-----------|--------------|--------|
| b1 | **Advisory dup-gate lane (TODO #286)** | Gate proven green, ritual documented in ADR/AGENTS | NOT wired into `scripts/ci-repro.sh`/CI; binary version not logged in any lane; determinism double-run not executed | M |
| b2 | **TODO #293 (upstream fingerprint-only baseline)** | Symptom treated: baseline re-recorded under one pinned binary | Upstream `--fingerprint-only`, config auto-discovery, help epilog — untouched (separate repo) | L |
| b3 | **TODO #291 (ADR refresh)** | Counts, tool-version pin, ritual, retracted budget claim — done in ADR | `art-dupl check` step NOT added to `docs/plan-authoring-checklist.md`; TODO_LIST rows #291/#293 now stale (re-record happened) | S |
| b4 | **t=2/t=3 report disposition** | t=1 set fully dispositioned; t=2/t=3 high-confidence groups read + dispositioned (extracted or category-accepted) | No per-group verdict table for the 37+6 "shown" groups committed to the ADR appendix | M |
| b5 | **Tools refactor verification** | Build + vet + import-order green | **Never runtime-smoked** (`go run ./tools/...` against real inputs); wired⇒e2e doctrine satisfied only at compile level | S |
| b6 | **Visual/e2e lanes** | HTML goldens byte-identical (proves rendered markup unchanged, incl. calendar month-nav used by `visualtest/calendar_nav_e2e_test.go` flows) | `nix run .#visual` (pixels + axe + e2e) NOT witnessed this session | M |
| b7 | **CHANGELOG warmth** | — | No `[Unreleased]` entry for the dedup pass; private-internal refactors may be exempt but that was never consciously decided | S |

## c) NOT STARTED

(Session-adjacent backlog observed, untouched — priority per existing TODO_LIST rows.)

1. **#283** Go test pinning `MIRRORED_PKGS` (bash) == `mirroredPackages` (Go) + `*_types.go`↔`.templ` pair test + `tc add` smoke into `t.TempDir()` — still wanted (1% tier).
2. **#284** Guard self-test script: 6 scenarios × idempotent ×2 in detached temp worktree — still wanted.
3. **#285** `packageDeps`/`packageImports` audit for the 22 rescued components + `--list-deps` honesty for kanban/charts — still wanted.
4. **#287** Doctrine verify: `nix run .#visual` kanban e2e + 7-module build + website tests (from the mirror session; still unexecuted).
5. **#288** Docs truth sweep (`tc new` → `tc init`/`tc add`; derived-counts claims; CHANGELOG reframe of the 22-file rescue).
6. **#289** Guard UX: `TC_SKIP_SYNC=1` opt-out + measured runtime + `--fix` pointer.
7. **#290** `starter/` dead-CSS trace → delete-or-wire.
8. **#292** art-dupl upstream: piped-output truncation root-cause (the bug that made this session's morning classification half-blind).
9. **#294** Small polish (`tc ls` derived footer, ci-repro guard-failure output review).
10. **#295** Promote dup-gate to BLOCKING CI — gated on #286 + 2 consecutive green advisory runs.

## d) TOTALLY FUCKED UP

Radical honesty — includes inherited breakage found, and my own failures:

| # | What's broken | Severity | Root cause | Mitigation |
|---|--------------|----------|------------|------------|
| d1 | **The "canonical gate" was a ghost system** (brutal-self-review Q7): committed baseline + documented ritual, but `check` had never once been run green by anyone — its first-ever execution this session failed with 105 false "new" groups. A gate that has never been executed cannot protect anything. | High (false confidence, not data loss) | Baseline recorded from one binary generation; hashes matched nothing under current binaries; no version stamp in the baseline file | FIXED this session: gate green under pinned fork; ADR + AGENTS document the pin. Remaining gap: nothing PREVENTS the next binary bump from silently re-breaking it (see f-2/f-3) |
| d2 | **Morning session's classification was built on a truncated capture** (inherited, TODO #292): 20 of 39 groups, no summary line — and the reports pasted to me were the re-run. My session repeated a variant of this sin (see d4) | High (process) | art-dupl upstream truncation | Untouched (separate repo, blocked) |
| d3 | **My bash exit-code masking, twice**: `cmd \| tail; echo $?` reported tail's exit 0 for a FAILING check ("EXIT: 0" with an error block in the log), and again for the first lint run. Exactly the failure class AGENTS.md already documents | Medium (self-inflicted, could have stopped work on a false green) | Pipeline exit semantics + rushing | Corrected protocol: capture to file + `echo EXIT=$?` directly on the command; used for all subsequent verdicts |
| d4 | **My --sort misdiagnosis**: compared `check`'s "Found total 105" (not-in-baseline count) against a plain scan's "Detected 848" (scan total) and briefly concluded the flag changes detection. ~4 tool rounds wasted before reading raw logs side-by-side | Medium (time, wrong-theory risk) | Like-for-like comparison discipline violated | a2 documents the real semantics; AGENTS.md gotcha added |
| d5 | **Self-inflicted edit bug**: ogshot deletion edit used placeholder text `// run captures` as new_string, leaving a stray comment + truncated file tail; caught only because I re-read the diff | Low (caught, fixed same session) | Careless multiedit new_string choice on deletions | Rule reinforced: deletions get empty/anchored replacements + immediate diff read |
| d6 | **Refactored three executables without running them once** (build+vet only) — violates the spirit of wired⇒e2e; the clones were exact copies so risk is genuinely low, but "verified" overstated the lane | Low-Medium | Runtime smoke not in my checklist for main-package refactors | f-3 |
| d7 | **Daemon fragmentation**: the session's work now lives in 5 heuristic auto-commits with hallucinated messages (14251e80 → c546c4e9); per M03, pushing requires witnessing green at the exact tip, and the history needs squashing/rewording first — that rewrite races the daemon | Medium (history quality, push ritual friction) | Known daemon behavior (T13) | f-2; temp-worktree squash pattern already documented in AGENTS.md |

## e) WHAT WE SHOULD IMPROVE

1. **Exit-code discipline as a hard rule**: never read `$?` after a pipe; every verdict command = `cmd > file 2>&1; echo EXIT=$?` + read the file's final line. This bit me twice and AGENTS.md once before; it should be reflex.
2. **Like-for-like forensics**: before theorizing about tool drift, print the two outputs' *summary lines* side by side. The --sort detour was avoidable in one round.
3. **Gates must self-identify**: a baseline file that cannot name its binary will always eventually lie. Until #293 lands, every lane that runs `check` must echo `art-dupl --version` beside the verdict, and the ADR must name the recording binary (now does).
4. **Ghost-system audit for gates**: "committed + documented" ≠ "enforcing". The dup-gate sat green-in-name-only for hours. Any new gate (there are several guards in this repo) should be executed once, end-to-end, on the day it is introduced — with the green output in the session record.
5. **Main-package refactors need a smoke lane**: `go run ./tools/X --help` (or one real capture) belongs in the definition of done for anything with a `main()`.
6. **CHANGELOG warmth decision**: make the exemption rule explicit (are private-internal refactors exempt?) instead of leaving it to per-session improvisation.
7. **Split-brain guard for counts** (brutal-self-review Q10): the "122 groups" claim now lives in ADR-0009 + AGENTS.md, hand-synced. A `TestVersionMatches*`-style drift guard (baseline entry count ↔ docs claims) would kill the class, same pattern as the version guards.
8. **Enum discipline gap found while working**: `TitleTag` (closed set h1–h6, raw `string`, no IsValid) sits in CardProps/EmptyStateProps/CollapsibleSectionProps despite the repo's "every closed-set enum ships IsValid" convention — consolidation candidate, now that all three switches share `headingTag`.

## f) NEXT TASKS (brainstorm — HARVEST input for TODO_LIST/ROADMAP)

Impact: 🔴 Critical / 🟠 High / 🟡 Medium / ⚪ Low. Effort: S <30min / M 30min–2h / L >2h.

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | Wire `art-dupl check` advisory lane into `scripts/ci-repro.sh` + CI, logging `art-dupl --version` beside the verdict (#286) | 🟠 | M | Quality |
| 2 | Squash/reword the 5 daemon commits (14251e80→c546c4e9) into one properly-messaged commit via detached temp worktree, then witness green at tip per M03 before any push | 🟠 | S | Cleanup |
| 3 | Runtime smoke the three refactored tools: one real `ogshot`/`shots`/`siteshots` run against real inputs | 🟠 | S | Quality |
| 4 | Witness `nix run .#visual` (pixels + axe + e2e incl. calendar month-nav) over the refactored components | 🟠 | M | Quality |
| 5 | Pin the art-dupl fork as a flake input (`nix run .#dupl`) so the gate is hermetic — fork needs go 1.27.1 while the repo pins 1.26.7, a flake input sidesteps that | 🟠 | M | Tooling |
| 6 | Baseline determinism proof: record twice, assert byte-stability (modulo `recordedAt`) in the advisory lane (#286 step) | 🟡 | S | Quality |
| 7 | Drift guard: baseline entry count ↔ ADR/AGENTS "accepted groups" claims (TestVersionMatches-style) | 🟡 | S | Quality |
| 8 | Decide + write the CHANGELOG exemption rule; add `[Unreleased]` entry for the dedup pass if not exempt | 🟡 | S | Documentation |
| 9 | Update stale TODO rows: #293 (re-record DONE; keep fingerprint-only ask), #291 (ADR refresh mostly DONE; checklist step remains) | 🟡 | S | Documentation |
| 10 | Add `art-dupl check` step to `docs/plan-authoring-checklist.md` (#291 remainder) | ⚪ | S | Documentation |
| 11 | Evaluate `--min-lines 2` (or t=2 → 35 shown) for the advisory lane to cut 2-line token noise; record chosen invocation in ADR | 🟡 | S | Quality |
| 12 | Upstream (art-dupl repo): ship stable/fingerprint-only hashes so baselines survive binary bumps (#293) | 🔴 | L | Tooling |
| 13 | Upstream (art-dupl repo): fix piped-output truncation (#292) — root cause of the morning session's half-blind classification | 🔴 | L | Bug |
| 14 | Upstream: make `check` print actionable-vs-raw counts distinctly (retire the "Found total" ambiguity) | 🟡 | M | Tooling |
| 15 | Evaluate `//art-dupl:accept` directives (fork feature) vs hash baseline for accepted residues like the headingTag call pair | ⚪ | M | Quality |
| 16 | #283: Go test pinning `MIRRORED_PKGS`==`mirroredPackages` + pair-completeness + `tc add` smoke into `t.TempDir()` | 🟠 | M | Quality |
| 17 | #284: guard self-test script (6 scenarios, idempotent ×2, detached temp worktree) | 🟡 | M | Quality |
| 18 | #285: `packageDeps`/`packageImports` audit for the 22 rescued components | 🟡 | M | Quality |
| 19 | #287: doctrine verify — kanban e2e + 7-module build + website tests from the mirror session | 🟠 | M | Quality |
| 20 | #288: docs truth sweep (`tc new`→`tc init`/`tc add`, derived-count claims) | 🟡 | M | Documentation |
| 21 | #289: guard UX (`TC_SKIP_SYNC=1` + measured runtime + `--fix` pointer) | ⚪ | M | Tooling |
| 22 | #290: `starter/` dead-CSS trace → delete-or-wire | ⚪ | S | Cleanup |
| 23 | #294: small polish (`tc ls` derived footer, ci-repro guard-failure output review) | ⚪ | S | Cleanup |
| 24 | #295: promote dup-gate to BLOCKING after 2 green advisory runs (#286 gate) | 🟡 | S | Quality |
| 25 | Consolidate `TitleTag` into a typed enum + `IsValid` per repo enum convention (all 3 switches already unified in `headingTag`) | 🟡 | M | Quality |
| 26 | Remove now-redundant `isValidHeadingTag` normalization in collapsible_section if headingTag's default covers it (verify tests first) | ⚪ | S | Cleanup |
| 27 | Name calendar aria-label literals ("Previous month"/"Next month") as constants per goconst discipline | ⚪ | S | Cleanup |
| 28 | ADR-0009: add explicit entries for datastar↔htmx `normalize` script twins and list_note↔end_of_list semantic twins (currently only category-accepted) | ⚪ | S | Documentation |
| 29 | ADR-0009 appendix: per-group verdict table for the t=2 (37 shown) + t=3 (6 shown) sets | ⚪ | M | Documentation |
| 30 | Cross-reference comment (or shared fixture) for `calendarNavQuery` duplicated between demo server and e2e contract test | ⚪ | S | Quality |
| 31 | Evaluate reusing `tools/internal/distserver` from `visualtest/demo.go`'s listener setup | ⚪ | S | Cleanup |
| 32 | Website module: extract the `w-11 h-11 rounded-lg bg-accent-dim` icon-box trio (sales/sections) into a site-local partial — needs scope decision (see g-q3) | ⚪ | S | Cleanup |
| 33 | Website module: unify clamp() hero/sales h1 typography tokens | ⚪ | S | Cleanup |
| 34 | Extend `utils.TestDocsCountDrift` to cover the ADR-0009 group-count claim (supersedes f-7 if preferred) | ⚪ | S | Quality |
| 35 | Add `visualtest/tools/*/README` or doc comments pointing at `tools/internal` so future tools reuse instead of copy | ⚪ | S | Documentation |
| 36 | Consider `tooling` convention: every capture tool gets a `--selftest` flag (allocator + listener come up) as its permanent smoke | ⚪ | M | Tooling |
| 37 | Re-run art-dupl at `-t 2`/`-t 3` post-extraction and prune ADR entries whose groups vanished (ledger hygiene) | ⚪ | S | Cleanup |
| 38 | Sweep `.golangci.yml` varnamelen config: two `coord` renames were needed — consider naming guidance instead of lint-churn | ⚪ | S | Tooling |
| 39 | CI: add a lane that runs the tc-sources guard self-test (#284's script) so guard regressions surface outside pre-commit | 🟡 | S | Quality |
| 40 | Document the five new extraction helpers in the skill/FEATURES "reuse these" lists (ADR Consequences updated; FEATURES/skill table check) | ⚪ | S | Documentation |
| 41 | Investigate why paste-era binary discovered 575 files vs 548 now — confirm the `_sources`-exclude divergence is the whole delta (one command, closes the forensics loop) | ⚪ | S | Quality |
| 42 | teach `check` (upstream) a `--baseline-version-must-match` guard: refuse to compare across binaries | 🟡 | M | Tooling |
| 43 | Add the evening-pass extractions to the demo smoke coverage if any render path is demo-visible (headingTag affects Card/EmptyState/CollapsibleSection demos) | ⚪ | S | Quality |
| 44 | Consider exporting `display.HeadingTag` (public helper) if consumers need heading-level rendering parity — API decision, needs owner input | ⚪ | S | Feature |
| 45 | Backfill CHANGELOG "### Refactored" (or equivalent) for v1.19.x line if the exemption rule (f-8) lands non-exempt | ⚪ | S | Documentation |
| 46 | Run the full `scripts/ci-repro.sh --lint --website` once over the squashed tip (pre-push ritual rehearsal) | 🟠 | M | Quality |
| 47 | Prune `/tmp` forensics logs into `docs/reviews/` evidence appendix if this report needs to stay self-contained (or accept tmp loss) | ⚪ | S | Documentation |
| 48 | Revisit ADR-0009 "Historical full-scan counts" table after #292/#293 land and the detector stabilizes (mark as detector-generation-specific) | ⚪ | S | Documentation |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Canonical art-dupl binary policy**: should I pin your fork (`~/projects/art-dupl`, needs go 1.27.1) as a flake input providing `nix run .#dupl`, so the advisory CI lane is hermetic and version-stamped — or is the system-installed 0.7.0 the intended CI tool until the fork releases? I tried answering from the repo (flake inputs, ci-repro) — nothing pins art-dupl today; this is tool-supply policy only you can set.
2. **History policy for the 5 daemon commits** (14251e80→c546c4e9, all unpushed): squash into one properly-worded commit via the temp-worktree pattern, or leave the daemon trail? Rewriting touches the daemon race window documented in AGENTS.md, so I won't decide it unilaterally.
3. **Scope/exemption policy**: (a) are private-internal refactors CHANGELOG-exempt, and (b) do you want website-module content clones (icon-box trio, clamp typography) actually extracted in a follow-up pass, or blanket-accepted like demo content? Both calls change what I do next; the repo docs don't answer them.

---

*Point-in-time snapshot — goes stale by design. Section (f) is HARVEST input: pull items into `TODO_LIST.md`/`ROADMAP.md` via docs-health HARVEST when resuming. Report per status-report skill; `.md` override per explicit user instruction; no manual commit (harness rule) — daemon will pick this file up.*
