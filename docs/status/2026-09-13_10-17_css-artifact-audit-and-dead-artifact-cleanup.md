# Status Report — 2026-09-13 10:17 CEST — CSS Artifact Audit & Dead-Artifact Cleanup

**Session scope:** one user question — *"Why do we have so many .css files!?!"* — answered with a full
audit of every CSS file in the repo, followed by cleanup, a new drift guard, and documentation.
Per instruction, this report covers **only this session's work and what was noticed during it**.
No unrelated research was done; observations about other in-flight work are recorded as-seen.

**State at report time:** working tree clean, all 7 Go modules green, master at `f2bd8906`.
The BuildFlow daemon committed all session work within seconds of each edit (commits
`53e87716`, `43f1522f`, `f2bd8906` — with heuristic messages, as usual).

---

## 0. The answer to the question that started it

The repo had **32 CSS files**: ~10 legitimate hand-written **sources**, 3 legitimate compiled
**artifacts** with real consumers — and **10 dead compiled artifacts** (~208 KB) that no in-repo
code consumed. They existed because:

1. **`e6a0e9f9` (2026-07-18)** committed 4 compiled `.out.css` files justified as "consumers receive
   pre-compiled CSS" — but no docs describe that consumer path and nothing in-repo read them.
2. **v1.8.3 (`dc77b3b6`)** multiplied them for a **`tc new` preset scaffolder that never shipped**
   (preset `.out.css` ×4, starter-kit copies ×3, root theme artifact). The CLI's `tc init` reads
   exactly two starter files (`app.css`, `custom.css`); it never read the other three embeds.
3. `examples/demo/demo.out.css` was a byte-identical duplicate of the demo's real artifact
   `static/app.css` — flagged as a split-brain risk on 2026-08-17, never decided until now.
4. `website/src/styles/global.out.css` was dead: the site compiles `global.css` itself via
   `@tailwindcss/vite`.

Because release.sh **and** the BuildFlow tailwind-build daemon recompile every artifact they find,
these dead files were regenerated on a loop — which is why 8 of them showed up "modified" at
session start (already daemon-committed at `1d1050e2` before I began).

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| a1 | **Full CSS inventory + consumer trace** — every `.css` on disk classified source-vs-artifact; every artifact's consumer searched repo-wide (only `scripts/release.sh` references any; guard tests read only the 3 legitimate files) | session transcript; AGENTS.md "CSS artifact inventory" |
| a2 | **Deleted 10 dead compiled artifacts (~208 KB)**: root `templ-components-theme.out.css`, `templates/presets/{default,emerald,glass,minimal}.out.css`, `cmd/tc/_sources/starter/{styles.css,templ-components-theme.css,templ-components-theme.out.css}` (3 embeds never read by the CLI, ~12.6 KB binary bloat removed), `examples/demo/demo.out.css`, `website/src/styles/global.out.css` | commit `43f1522f` |
| a3 | **`scripts/release.sh` compile block 5 → 3 targets** — recompiles exactly the three legitimate artifacts; comment names the guard + removal date | `scripts/release.sh:282-294`, `bash -n` clean |
| a4 | **New drift guard `utils.TestCompiledCSSInventory`** — asserts the committed `.out.css` set == exactly `templates/templ-components-theme.out.css`, and that `templates/styles.css` + `static/app.css` exist; **negative case proven** (planted zombie `.out.css` → hard FAIL with remediation message; cleaned via `trash`) | `utils/compiled_css_inventory_test.go:26` |
| a5 | **AGENTS.md documentation** — "CSS artifact inventory" bullet at the head of Demo Infrastructure: source/artifact taxonomy, the two same-named theme files disambiguated, deletion decision, daemon-resurrection warning | commit `f2bd8906` |
| a6 | **CHANGELOG `[Unreleased]` warmed** with a `### Removed` entry (release-script contract kept) | commit `f2bd8906` |
| a7 | **docs/theming.md drift fix** — "Three starter presets" → "Four" (trivial fix-on-sight) | commit `f2bd8906` |
| a8 | **Full verification matrix green**: per-module loop `utils icons errorpage charts/echarts datastar htmx` all OK; root `templ generate` (pinned v0.3.1020, **zero diff**) + `go build ./...` + `go test ./...` exit 0; `go vet` + `golangci-lint` on utils clean; `cmd/tc` tests + build green | session transcript |
| a9 | **Concurrent-session work respected, not reverted** — `.githooks/` promotion + `scripts/setup-hooks.sh` + CONTRIBUTING edits appeared mid-session from another actor; read, judged coherent, left alone; my staging stayed scoped | commits `53e87716`/`43f1522f` (mixed by daemon) |

## b) PARTIALLY DONE

| # | Item | Gap |
|---|------|-----|
| b1 | **Guard-test documentation sync** — new guard is in AGENTS.md + release.sh, but **not** in the repo-wide guard table in `skill/SKILL.md` (the "Repo-wide guard & compliance tests" table). Forgotten until writing this report. | one table row + one line in the guard-scripts prose |
| b2 | **Dead-reference sweep for the unshipped scaffolder** — `docs/cross-project-analysis.md:144` still references "tc new tests +34". I found it during report writing; not yet cleaned. | grep + prune that one reference (and audit for siblings) |
| b3 | **Daemon-resurrection watch** — guard exists, but the first BuildFlow tailwind-build cycle *after* the deletions has not been observed yet. Claim "daemon can no longer litter" is unproven until then. | one daemon cycle + `git status` check |
| b4 | **Keep-decision on `templates/styles.css` + `templates/templ-components-theme.out.css`** — kept because release.sh (post-v1.10.0 hardening) treats them as distribution targets; but I found **no actual consumer** for either. Decision rests on convention, not proven use. See question 1. | consumer evidence or deletion |
| b5 | **CHANGELOG entry ordering** — I placed `### Removed` before `### Fixed`; Keep-a-Changelog canonical order is Added/Changed/Deprecated/Removed/Fixed/Security. Cosmetic, easy reorder. | 3-line move |

## c) NOT STARTED

| # | Item | Why it's on the list |
|---|------|----------------------|
| c1 | **BuildFlow upstream skip-list for `*.out.css`** (TODO #125 family) — the *root* fix; in-repo guard is the seatbelt, not the fix | lives in `larsartmann/buildflow`, not this repo |
| c2 | **HARVEST of this report's section (f) into `TODO_LIST.md` / `ROADMAP.md`** (docs-health HARVEST mode) | skill contract: section (f) must not stay entombed in a timestamped file |
| c3 | **`nix run .#visual` run post-cleanup** — deletions shouldn't touch rendering (demo embed untouched), but one pixel-suite pass is cheap confidence | not run this session |
| c4 | **`scripts/ci-repro.sh --lint`** — I ran the per-module loop + targeted lint, not the full step-for-step CI clone | pre-push gate per AGENTS.md |
| c5 | **Minification guard scope** — `scripts/check-css-minified.sh` guards only `static/app.css`; the two `templates/` targets have no minification guard (daemon un-minified artifacts before, per 2026-09-08 report) | small extension |
| c6 | **Single-source the 3-target list** — the list now lives in release.sh (bash) *and* the guard test (Go); they can drift independently | extract to one shared source |
| c7 | **Document (or kill) the "pre-compiled CSS for consumers" story** — currently an implicit convention with zero docs | README/adoption-guide decision |
| c8 | **Compiled-artifact naming normalization** — kept artifacts now use three different conventions: `*.out.css`, `styles.css`, `app.css` | rename or document |
| c9 | **`tc init` theme gap** — README:344 tells consumers to copy `templ-components-theme.css`, but `tc init` doesn't scaffold it (pre-existing gap; the starter used to embed a copy nothing read — deleted) | product decision + small CLI change |
| c10 | **Annotate the 2026-09-08 status report** — it claims "daemon-regression-check.sh … done", but no such file exists in `scripts/`; `check-css-minified.sh` covers only part of the claimed trio. Ghost/renamed item; needs docs-health ANNOTATE after verification | verify-then-annotate |

## d) TOTALLY FUCKED UP

**Nothing from this session's changes.** Honest accounting of the nearest misses and non-mine breakage:

| # | Item | Severity | Detail |
|---|------|----------|--------|
| d1 | **Concurrent session broke the utils test compile mid-session** (not my change, self-healed) | was-🔴 now-green | `utils/docs_count_test.go` briefly referenced `countIconNames` before its definition existed + had invalid `\(` escapes in an interpreted string → whole utils package uncompilable for several minutes. The other actor fixed it; the daemon then merged their second test file (`skill_count_test.go`, −72 lines) into `docs_count_test.go` (+61) at `f2bd8906`. My guard test was unaffected (targeted run passed pre-incident; full module green after). Lesson re-confirmed: cross-check any "broken" tool output against a fresh run before acting — I re-ran `go vet` instead of "fixing" their file, which would have clobbered in-flight work. |
| d2 | **My process slips** (minor, self-caught) | 🟡 | (1) malformed `rg --type-add` invocation wasted a round trip; (2) `rg -rn` flag misuse (`-r` = replace!) produced garbage output before I spotted and redid it; (3) one edit-before-read tool rejection on release.sh. Total cost: ~3 wasted tool calls, zero incorrect state. |
| d3 | **Pre-existing fuck-up this session exposed** | 🟠 | Two files coincidentally named `templ-components-theme.css` with completely different purposes (root: `@theme`/dark-mode example linked from README:344; templates/: ADR-0008 semantic tokens). Also the v1.8.3 commit added embedded starter files the CLI never read — committed infrastructure for an unshipped feature with a drift guard that explicitly skips the `starter/` dir, so nothing would ever have caught their staleness or absence. |

## e) WHAT WE SHOULD IMPROVE

1. **Kill zombie artifacts at commit time, not audit time.** This session needed git archaeology
   (two contradicting commits on the same day in July: one gitignored these artifacts, one committed
   them). The new guard now enforces the set — but the *rule* "a committed artifact needs a named
   consumer" should be in CONTRIBUTING, not just AGENTS.md.
2. **Daemon contract.** Three failure classes recurred this session alone: heuristic commit messages
   that hide semantic work (the 10-file deletion landed as "auto-commit 15 changed file(s)"),
   mid-edit snapshot commits (the docs_count breakage window), and artifact resurrection risk.
   All three are BuildFlow-side; the in-repo guards are seatbelts. TODO #125 and the message-quality
   fix belong upstream in `larsartmann/buildflow`.
3. **"Done" claims in status reports must point at an artifact.** The 2026-09-08 report's
   "daemon-regression-check.sh done" doesn't resolve to a file. Verify-external-claims applies to
   our own past reports too (c10).
4. **Same-name-different-file traps.** The two `templ-components-theme.css` files cost real audit
   time. Rename one (e.g. root → `theme-example.css`) or at minimum add cross-referencing header
   comments. Cheap now, expensive later.
5. **Release script ↔ guard duplication** (c6): one list, two consumers, or drift is guaranteed.

## f) UP TO 50 THINGS WE SHOULD GET DONE NEXT

*Brainstorm per skill contract — most items below #10 are ROADMAP fuel, not commitments. Grouped;
✱ = direct follow-up from this session.*

**Decisions needed (questions in section g feed these)**
1. ✱ Decide fate of `templates/styles.css` + `templates/templ-components-theme.out.css` (Q1).
2. ✱ BuildFlow upstream: implement `*.out.css` skip-list / test-gated auto-commits (Q2, TODO #125).
3. ✱ Confirm + finish the tracked-hooks promotion from the concurrent session (Q3).
4. ✱ Presets: bless `templates/presets/*.css` as a documented consumer feature (theming.md) — or
   remove if the scaffolder story was their only purpose (I lean keep; docs document them).
5. Document or delete the "pre-compiled CSS distribution" consumer path (c7).
6. Decide `tc init --theme` (use the preset sources) vs. status-quo (c9).
7. Decide compiled-artifact naming convention (c8).

**Direct cleanup follow-ups (small, this week)**
8. ✱ Sync `skill/SKILL.md` guard table with `TestCompiledCSSInventory` (b1).
9. ✱ Prune "tc new tests" reference from `docs/cross-project-analysis.md` + audit siblings (b2).
10. ✱ Reorder CHANGELOG `### Removed` into canonical position (b5).
11. Extend minification guard to the two `templates/` targets or fold into the Go guard (c5).
12. Single-source the compiled-target list shared by release.sh + guard test (c6).
13. Annotate the 2026-09-08 report's ghost "daemon-regression-check.sh done" item (c10).
14. HARVEST this report's section (f) into TODO_LIST/ROADMAP (c2).
15. One `nix run .#visual` pass post-cleanup (c3).
16. One `scripts/ci-repro.sh --lint` before next push (c4).
17. Observe one daemon cycle for `.out.css` resurrection; if it happens, delete + escalate to BuildFlow (b3).
18. Post-cleanup `nix run .#css` byte-stability check per the post-daemon-commit trio.
19. Add cross-referencing header comments (or rename) to disambiguate the two `templ-components-theme.css` files (e4).
20. Add "committed artifact ⇒ named consumer" rule to CONTRIBUTING.md (e1).

**Guard/test hardening**
21. Guard test: also assert the three targets are minified (move check-css-minified.sh's logic into Go, one place).
22. Guard test: fail-loud audit — confirm `TestCompiledCSSInventory` cannot `t.Skipf` on any path (it currently can't; keep it that way).
23. Consider a release.sh unit test (bash-bats or golden) pinning the compile-command list (supports #12).
24. Re-check `TestCSSFreshness` CI-scope extension (older report flagged only 1 of 5 targets gated; now 3 targets exist — re-evaluate scope post-cleanup).
25. Drift-guard the AGENTS.md inventory bullet: guard test asserts the three artifact paths it names still exist (docs vs reality single-sourcing).

**Daemon/BuildFlow (upstream)**
26. Commit messages from `git diff --stat`, not templates (message-quality fix, upstream).
27. No mid-edit snapshot commits: debounce + green-test gate before auto-commit (upstream).
28. Stop re-appending `*_templ.go` to `.gitignore` on every pre-commit run (known, still open).
29. Pin-flip suppression for manifest/lockfile pairs (from 2026-09-08 report; still relevant).
30. After any upstream BuildFlow fix, re-verify the in-repo guards are still needed/accurate (remove seatbelts when the car gets airbags).

**Docs hygiene noticed in passing**
31. Reconcile component counts: AGENTS.md says display=43, SKILL.md catalogue says display=42 ("118 components" headline vs per-package sums) — same drift family the old `skill_count_test.go` watched; verify the daemon's merge into `docs_count_test.go` still covers it.
32. Audit remaining planning docs for `tc new` promises (found 1 hit in cross-project-analysis; ROADMAP/planning dirs excluded from my sweep — one deliberate pass).
33. Status-report archive policy: 144 files in `docs/status/`, growing ~10/day on daemon-heavy days; define prune/archive rule.
34. FEATURES.md/ROADMAP: reflect the artifact-inventory decision (harvest will do part of this).
35. README CSS section: decide whether to mention the pre-compiled artifacts at all (ties to #5).

**CLI/consumer surface**
36. `tc init`: optionally scaffold `templ-components-theme.css` (align with README:344 flow).
37. `tc presets` or `tc init --preset <name>`: the original unshipped idea, built on preset *sources* — roadmap decision (#4, #6).
38. `tc` binary size: verify embed shrank as expected (~12.6 KB) and note in release notes if consumer-visible.
39. Consider `tc add --css` hint: when adding components with `.tc-*` classes, remind about `custom.css` (idea-grade).

**Verification/process**
40. Run the full `nix run .#verify` once before next release cut (session used per-module loop; the flake app is the canonical done-check).
41. govulncheck via the nix-shell wrapper at next release (standing convention, unchanged).
42. Keep the two-registries lesson: when another actor's changes appear mid-session, read + judge + leave (this session's d1 handling worked; codify in AGENTS.md collaboration note if not already implied).
43. Add the negative-case proof pattern to the testing docs: this session's guard was proven by planting a zombie file — worth a line in the guard-test conventions.
44. Re-verify `.golangci.yml` disabled-linters guard still green after the concurrent session's commits (TestGolangciDisabledLinters ran green in the module loop — spot-check only).
45. Check whether `check-master-green.sh` / `verify-local.sh` reference any deleted artifact paths (my sweep covered release.sh only; these two scripts touch CSS flows).

**Farther out (roadmap fuel)**
46. Website: `global.out.css` deletion means the site's only CSS path is Vite — document that in website/README if one exists.
47. Docker CSS stage: confirm it never depended on deleted artifacts (it compiles fresh from demo.css — verify once).
48. Consider `.gitignore` hardening: ignore `*.out.css` globally *except* the tracked one, as belt-and-suspenders against future stray artifacts (gitignore allowlist juggling — only if #2 upstream stalls).
49. Demos: `demo.out.css` deletion closes the 2026-08-17 "delete or single-source" item — annotate that old report's open item as resolved.
50. Revisit in 30 days: if the inventory guard never fired and the daemon never resurrected anything, consider demoting it from hard FAIL to informational (guards that never earn their keep are noise).

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **`templates/styles.css` + `templates/templ-components-theme.out.css`: real consumers or next deletions?**
   I kept them because release.sh has treated them as distribution targets since the v1.10.0
   hardening — but I found no doc, test, or code that consumes either. Do you know of an actual
   consumer flow (e.g., people fetching pre-compiled CSS from GitHub raw or the Go module cache),
   or should they join the other 10? If kept, #5/#7 (document the path) becomes mandatory.

2. **BuildFlow: fix upstream now, or is the in-repo guard the accepted endgame?**
   The root fix for the entire churn class (resurrected artifacts, heuristic messages, mid-edit
   snapshots) lives in `larsartmann/buildflow` (TODO #125 family). Do you want me to prepare the
   upstream change/skip-list, or is `TestCompiledCSSInventory` + daemon-watch the intended
   permanent control plane? This decides items #2, #26–30.

3. **Is the concurrent session's tracked-hooks promotion complete and blessed?**
   Mid-session, `.githooks/pre-commit` (tracked) + `scripts/setup-hooks.sh` + CONTRIBUTING edits
   appeared from another actor and the daemon committed them mixed with my deletions in `43f1522f`.
   It looks coherent (core.hooksPath activation for fresh clones), but I can't know whether that
   session considers it finished — or whether CI/fresh-clone flows still assume `.git/hooks/`.
   Should I verify/complete it, or is its author still mid-flight (in which case I stay out)?

---

*Point-in-time snapshot — 2026-09-13 10:17 CEST — master `f2bd8906`, tree clean, all modules green.
Format note: user explicitly requested `.md`; the status-report skill's canonical HTML dashboard was
skipped for this report per that instruction. No manual commit made (harness forbids commits without
explicit request); the auto-commit daemon will pick this file up.*
