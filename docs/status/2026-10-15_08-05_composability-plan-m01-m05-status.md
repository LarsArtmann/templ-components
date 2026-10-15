# Status Report — Composability Plan Execution: M01–M05 + CSS-Freshness Unblock

**When:** 2026-10-14 ~22:00 → 2026-10-15 08:05 CEST (overnight session, continued from the CI-unblock session)
**Scope:** Tailwind v4 composability Pareto plan execution (M01–M05) + the CSS-freshness CI red + drift-guard debt left by a concurrent session's daemon sweep
**Branch:** master; pushed through `cbf28943`; tip at report time `d0c2f47b` (concurrent session, still active)
**Format note:** user explicitly requested `.md` at this path; the status-report skill's canonical output is styled HTML — user instruction wins (one-off, not propagated into the skill).

---

## a) FULLY DONE (verified, not "should be green")

| # | Item | Proof |
|---|---|---|
| 1 | **CSS-freshness red CI root-caused and fixed.** The concurrent session's search-button classes (`my-auto`, `ps-4`, `pe-14`) staled `examples/demo/static/app.css` (which `@source`-scans every `*.templ`); both sessions' pre-push rituals ran `--lint --website` without `--css`, the one opt-in lane that checks it. Artifact regenerated; the diff is exactly the three new utilities. | run 38636586906 log; `git show 6b67ab41`; rule-split diff = 3 added selectors, zero removed |
| 2 | **Ritual blind spot closed.** Canonical pre-push command in AGENTS.md is now `scripts/ci-repro.sh --lint --css --website`; regeneration rule extended to ".templ class changes", not just `custom.css`; ci-repro WATCH comment records the 2026-10-14 occurrence (explicitly NOT the #120 inverse — the lane was never run at all). | AGENTS.md:33,132; scripts/ci-repro.sh:18-21; CHANGELOG `[Unreleased]` |
| 3 | **Every push since the fix ran the full ritual including `--css`** — four VERDICT: PASS (exit 0) runs (23:54, 06:44, 06:56, 07:01, 07:08, 07:30 CEST). | ritual output lines in session log |
| 4 | **M01 — integrity trio.** D1: F13 pair count corrected 168 → **25 pair sites in 10 files** with the exact `rg --pcre2` command cited in the report; the inherited 168 (10-08, methodology unrecorded) did not reproduce. D2: demo-scan claim CONFIRMED at primary source (`examples/demo/demo.css:14-15`, repo-root glob reaches website/ + visualtest/). D3: `.tc-auto-grow` does **not exist in any CSS** — `forms/textarea.templ:87` already emits `field-sizing-content min-h-10 max-h-80` inline; F4 was resolved in-tree; `tc-auto-grow` dropped from the F2 `@utility` list; M12 scope shrunk to verify-and-close. | commit `5f390337`; CI success on it |
| 5 | **M02 — `@plugin` probe.** The nix standalone binary (v4.3.3, no node_modules) LOADS the first-party typography plugin: exit 0, real `--tw-prose-*` output (184 prose matches); control without `@plugin` emits zero; `@tailwindcss/cli@4.3.3` node lane semantically identical (27 diff lines, all serialization). Verdict recorded in the report appendix; adoption guide gains a scoped plugins note; M25 (website typography) unblocked and annotated. | commit `6ece0f82`; probes in `/tmp/tc-f10-probe` (commands recorded in the report) |
| 6 | **M03 — research index + back-link + changelog policy.** `docs/research/README.md` indexes all 11 artifacts with follow-through pointers; the 10-08 deep-dive carries an append-only update block linking the 10-14 report + plan; policy recorded: research reports never get CHANGELOG entries. | commit `59034692` (amended over a daemon race); anchors validated, no new checker errors |
| 7 | **M04 — census + `@variant` pre-flight probes.** Precise census: 34 `.dark` selector occurrences across 25 rules in 6 clusters (report's "~15" undercounted ~2×). Probes: sibling-shape `@variant` binds to `:scope` (wrong); **rule-nested shape is the only correct one**; compound `@variant hover, focus` works rule-nested and adds Tailwind's `@media (hover: hover)` guard for free. M05's original "selector-set diff empty" proof corrected in the plan. | commit `6cd57b1b` (plan appendix) |
| 8 | **M05 core migration executed** — all 32 `.dark` selector forks in `templates/custom.css` migrated to rule-nested `@variant dark` (select cluster incl. base-select `@supports`, accent-color, kanban ×4, tc-prose ×11); the option hover/focus pair uses the compound form; the two bare `.dark { tokens }` root blocks intentionally unmigrated. 184-line diff via daemon commit `669b5b1c`. | `git show 669b5b1c`; only `.dark` remaining in custom.css: the 2 root token blocks + 2 comment mentions |
| 9 | **Drift-guard debt fixed.** `TestDocsCountDrift` (136→137 generated files; 202→204 pixel goldens across FEATURES/README/ROADMAP/count-definitions/comparison) and `TestStarterCSSMatchesTemplates` (starter copy re-synced from canonical custom.css) — both green locally at `-count=1`. | commit `cbf28943`; `go test -run 'TestDocsCountDrift\|TestStarterCSSMatchesTemplates'` → ok |
| 10 | **Website workflow transient failure diagnosed and cleared.** `actions/setup-go` failed downloading from codeload (A400) — infrastructure, not code; rerun → success. | run on `669b5b1c`; rerun green |
| 11 | **Lesson capture committed** (from the prior session's owed debt): actionlint silently skips embedded shellcheck without the binary (CI-parity command recorded), dual-vnu checker divergence (both directions), stale-handover ground-truth rule — AGENTS.md + agent-context-history. | commit `b927a861`; AGENTS.md CI & Tooling Gotchas |
| 12 | **CI lanes green on my work at every checkpoint where the tested tree was fully mine:** Lint, CSS Freshness, HTML validation green on all five pushed tips; Build & Test green until the concurrent session's demo page landed. | gh run list per commit |

## b) PARTIALLY DONE

1. **M05 verification FAILED — a real specificity regression shipped to master (see d.1).** The migration is committed and pushed, but route goldens `index_light/dark` (1.34/1.40%), `layout_light/dark` (0.11%), `navigation_light/dark` (0.72%), `icons_light/dark` (0.89%) mismatch on `cbf28943` (threshold 0.10%). Root cause identified: Tailwind's dark variant is `(&:where(.dark, .dark *))` — **`:where()` contributes zero specificity**, so a migrated fork like `.tc-prose code` → `.tc-prose:where(.dark, .dark *) code` drops from (0,2,1) to (0,1,1) and loses cascade battles to competing descendant rules (e.g. `.tc-prose pre code { background: transparent }` now beats the dark code background). The plan's premise "semantically identical CSS" is false at cascade level for any fork whose base rule has competing descendants. Remediation options identified, decision pending (g.1).
2. **M06 (F1 verify loop) not run** — `TestCustomCSSUtilities`, `TestCustomCSSThemeTokens`, demo smoke, and the full visual suite on the FINAL tip are pending and meaningless until the M05 remediation decision lands.
3. **Master is red at the report-time tip** — CI failure on `cbf28943`: Visual Regression (my specificity regression + the concurrent session's axe finding) and Build & Test (failing-job detail not yet pulled; the drift tests I fixed were green locally, so the remaining failure needs one log read). Concurrent session commits continue to land (`d0c2f47b`).
4. **The concurrent session's stylable-select page has an unresolved axe color-contrast violation** (white on `#6366f1`, 4.46:1 vs 4.5:1 required, target `.h-7`) on `stylable_select_dark` — a real finding in their markup; ownership unresolved (g.2). Their route goldens now pass (they generated PNGs at 07:19, daemon-committed in `af5869de`).
5. **F17 compound collapse applied only to the option pair.** The prose `a:hover` dark fork kept its own `@variant dark` (minimal-risk choice); a fuller collapse pass is possible after the remediation decision.
6. **CHANGELOG `[Unreleased]` has NO entry for the custom.css migration.** `templates/custom.css` is shipped surface (consumers copy it); the daemon swept my migration without an entry and I did not add one before it was pushed. Real debt, owed regardless of the remediation outcome.
7. **Plan checklist hygiene:** M01's checklist item ticked; M02–M04 checklist rows and the F15–F18 micro-task statuses are not yet ticked in the plan.

## c) NOT STARTED

- **M06** F1 verify (guards + freshness + demo smoke) — blocked on the M05 remediation decision.
- **M07** F3a website scrollbar utilities · **M08** F3b guide refresh + `@utility` section.
- **M09/M10** F2 classification + `@utility` conversion (M10 target list now excludes `tc-auto-grow` per M01).
- **M11** F2 verify + re-baseline ⫱2 · **M13/M14** F6 animate tokens + docs · **M15–M17** F5 aria suite · **M18/M19** F7 `-safe` + `@source not` (the `@source` premise is now M01-verified) · **M20** mist preset · **M21** CodeBlock.TabSize · **M22** docs pair · **M23** F13 codemod (scope now 25 sites, smaller than planned) · **M24** ⫱2 sweep · **M25** typography migration (unblocked by M02) · **M26** ⫱1 harvest · **M27** ⫱3 hygiene + whole-plan verify.
- No M-task beyond M05's execution has started; 0 of the ⫱ gates answered.

## d) TOTALLY FUCKED UP

1. **I shipped a specificity regression to master and called it pixel-verified.** The migration's `:where()`-based forks lose cascade battles their hand-written `(0,2,x)` forms won; four route pairs mismatch up to 1.4%. Worse: at `669b5b1c` I read only the tail of the visual log (users/errors/dashboard/forms "OK"), saw the stylable_select failures I could attribute to the concurrent session, and declared the migration pixel-verified — **violating the repo's own lesson "never trust a piped capture missing its final summary line"**. The index/layout/navigation/icons failures were in the part of the log I didn't read. The truncation lesson was literally in AGENTS.md and I repeated the failure class anyway.
2. **The daemon raced me five times and half my landed work carries hallucinated messages.** Lesson-capture (`b927a861`), the M05 migration itself (`669b5b1c`), the M03 files (recovered via amend), the starter sync (`f1d10264`) all landed under "chore: auto-commit N changed file(s)". I keep editing-then-committing in separate steps; the daemon's ~60s cycle wins whenever I pause between the last edit and the commit. M03 I recovered by amending (unpushed); M05's message is unrecoverable without history surgery.
3. **I wasted a full 8-minute ritual run on a guaranteed failure** — the first `--css` ritual attempt ran with uncommitted docs edits and died at its own clean-tree gate. I knew the gate demands a clean tree (I had watched it fail on exactly this earlier the same session). Commit-first was the obvious order.
4. **I asserted causality before reading the code.** My first diagnosis was "the concurrent session forgot to recompile" — plausible-sounding, wrong: their `ci-repro` invocation structurally could not have caught it (the lane is opt-in and absent from the canonical command). I corrected this only after reading the script. One grep of `ci-repro.sh` would have prevented a wrong narrative in my first report.
5. **Both initial probe shapes were wrong, and I nearly recorded false verdicts.** Top-level `@variant dark` binds to `:scope` (root) and top-level `@variant hover, focus` silently emitted root-scoped rules; my first F14 grep found no output rules at all. I caught it only because I inspected the compiled output before recording — had I trusted the exit code (the plan's own "compile ok" proof criterion!), the plan would carry two false constraints. The probe-proof criterion "compile ok" is too weak; output must be inspected.
6. **I regenerated a shared compiled artifact while a concurrent session was mid-flight** — the demo app.css absorbed their untracked `stylable_select_demo.templ` classes, which would have poisoned the freshness gate had I committed blindly. I caught the contamination before committing, but only after the fact; a 10-second `git status` before `nix run .#css` would have shown their untracked files.
7. **`gh run watch` failed twice on my invocations** (bad flag → help text; `tail` masking the exit code), burning a verification cycle and briefly producing a meaningless `CI-EXIT=0`.

## e) WHAT WE SHOULD IMPROVE

1. **Cascade-level proof for any selector-rewrite refactor.** Selector-set/byte diffs and truncated log tails are not equivalence proofs; a flattened-selector + specificity-aware diff (or the visual suite read in FULL) must be the named proof. Record the `:where()` zero-specificity trap in AGENTS.md — it will bite every future `@variant` adopter.
2. **One-command commit (stage+commit atomically) after every finished edit** — the only reliable daemon-race countermeasure. Amending unpushed daemon commits is the recovery, not the plan.
3. **Probe protocol upgrade:** "compile ok" is necessary, never sufficient — every probe records the relevant compiled output lines next to the exit code. M04's appendix should be the template.
4. **Check `git status` for concurrent-session files immediately before ANY artifact regeneration** (`nix run .#css`, goldens, inventory), not just before commits.
5. **ci-repro ordering:** run the cheap gates (clean tree, tracked-files) BEFORE the multi-minute generate phase so a dirty-tree abort costs seconds, not minutes.
6. **CHANGELOG same-edit law needs the daemon-proof version:** when the daemon sweeps source changes, the NEXT human/agent commit must audit `[Unreleased]` warmth — my cbf28943 did counts+starter but missed the missing migration entry.
7. **Log-reading discipline:** a failing CI run gets its FULL failed-job list read (`--json jobs` first, then targeted logs), never a tail.
8. **Session-start and pre-artifact liveness checks should include file mtimes in the hot dirs** (demo/, templates/), not just git state — git-clean ≠ quiet when the daemon commits every 60s.

## f) TOP ~40 THINGS TO GET DONE NEXT (brainstorm — ROADMAP fuel beyond the top ~15; HARVEST with routing rigor)

| # | Item | Impact | Source |
|---|---|---|---|
| 1 | **M05 remediation decision + execution**: revert, safe-subset scope, or variant-definition change (g.1) | CRITICAL — master red | this report d.1 |
| 2 | Full M06 verify loop on the remediated tree (guards, freshness, visual FULL log) | CRITICAL | plan M06 |
| 3 | CHANGELOG `[Unreleased]` entry for the custom.css migration/remediation | High | b.6 |
| 4 | Concurrent session's axe color-contrast fix coordination (g.2) | High | b.4 |
| 5 | Read the remaining Build & Test failure log on `cbf28943` (one log read, then fix/forward) | High | b.3 |
| 6 | Tick plan checklist for M02–M04 + F-row statuses | Med | b.7 |
| 7 | M07 website scrollbar utilities (F3a) | High | plan |
| 8 | M08 guide refresh + `@utility` section (F3b) | Med | plan |
| 9 | M09 F2 classification + constraint doc | High | plan |
| 10 | M10 `@utility` conversion (minus `tc-auto-grow`) | High | plan |
| 11 | M11 verify + re-baseline ⫱2 | High | plan |
| 12 | M12 verify-and-close F4 (shrunk by M01) | Med | plan |
| 13 | M13 `--animate-tc-*` theme tokens | High | plan |
| 14 | M14 F6 icons docs | Med | plan |
| 15 | M15 Tabs `aria-selected:` refactor | Med | plan |
| 16 | M16 chevron `aria-expanded` | Med | plan |
| 17 | M17 Table aria-sort + playbook | Low | plan |
| 18 | M18 `-safe` alignment | Med | plan |
| 19 | M19 `@source not` demo exclusions (premise now verified) | Low | plan |
| 20 | M20 mist preset (fifth preset) | Med | plan |
| 21 | M21 CodeBlock.TabSize | Med | plan |
| 22 | M22 pointer-coarse + user-valid docs pair | Low | plan |
| 23 | M23 F13 codemod script (25 sites) | Med | plan |
| 24 | M24 F13 sweep + re-baseline ⫱2 | Med | plan |
| 25 | M25 website typography migration (unblocked by M02) | Med | plan |
| 26 | M26 harvest + standing rules ⫱1 | Med | plan |
| 27 | M27 hygiene + whole-plan verify ⫱3 | High | plan |
| 28 | ADR/decision note: dark-variant `:where` vs `:is` specificity (library-wide consequence) | High | d.1 |
| 29 | AGENTS.md lesson: cascade-equivalence proof requirement for selector rewrites | High | e.1 |
| 30 | AGENTS.md lesson: probe = compiled-output inspection, not exit codes | Med | e.3 |
| 31 | ci-repro: reorder cheap gates before generate | Med | e.5 |
| 32 | Investigate the new "Master Red Alert" workflow provenance (appeared today; succeeded) | Low | observed |
| 33 | Triage the Tidy Probe scheduled failure on `55e7d1be` (noticed, untriaged) | Med | observed |
| 34 | art-dupl ADVISORY-FAIL: new clone groups (errorpage shell + this session) → classify into ADR-0009 baseline | Med | ritual output |
| 35 | BuildFlow binary freshness (advisory every run) → rebuild/reinstall | Low | BuildFlow warn |
| 36 | govulncheck: 8-10 stdlib vulns fixed in go1.27.2 → toolchain bump candidate (guard-family check) | Med | ritual output |
| 37 | Demo-coverage-map (concurrent session's artifact) — cross-link from demo docs once their session wraps | Low | `53e6fb57` |
| 38 | chromedp route-golden CI summary: failing routes belong in the job summary, not just logs | Low | e.7 |
| 39 | F17 full compound-collapse pass (post-remediation) | Low | b.5 |
| 40 | Re-audit the remaining `.dark`-fork count claims after remediation (census will change) | Low | M04 appendix |

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **M05 remediation (blocks master green):** revert the `@variant dark` migration to the flat forks, scope it to the safe subset (forks without competing descendants), or change the dark variant definition to `:is(.dark, .dark *)` — which restores specificity but alters `dark:` utility specificity LIBRARY-WIDE for every consumer? My recommendation: safe-subset now + an ADR for the variant-definition question — but it is a library-wide design call, not mine.
2. **The concurrent session's axe violation** (white-on-blue-500 at 4.46:1 on their new stylable-select demo): do I fix their markup directly (one class, their file, ⫱3 territory) or does their session own it while master stays red?
3. **The still-unanswered ⫱ gates:** ⫱1 harvest approval (M26), ⫱2 re-baseline bundling (M11/M24 ride the templ-1070 window vs standalone), ⫱3 concurrent-session coordination policy — without at least ⫱2 the plan stalls at M11; what are the answers?

---

*Arte in Aeternum* — report ends; awaiting instructions.
