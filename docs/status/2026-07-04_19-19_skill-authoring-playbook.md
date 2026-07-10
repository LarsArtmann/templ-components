<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Status Update — 2026-07-04 19:19

**Session scope:** Create a proper `SKILL.md` for the templ-components repo, self-critique, execute improvements, push.

**Branch:** `master` — pushed to `origin/master` at `4bd359b`.
**Session commits:** 5 (all on `master`, all pushed).

---

## a) FULLY DONE

1. **`skill/SKILL.md` created and pushed** (282 lines, under the 500-line skill-creator ideal).
   - Valid YAML frontmatter: `name: templ-components`, rich trigger-loaded `description`, `metadata.tags`.
   - Sections: Principles, Flake commands, Process, Component anatomy (9 steps), Decision trees, Mandatory conventions, Anti-pattern rubric, Testing matrix, Deeper reading, Install.
   - Distinct from `AGENTS.md` (context) and `CONTRIBUTING.md` (contributor setup) — this is a _procedural playbook_ ("how to make a component fit").
   - Follows the skill-creator format: progressive disclosure, explains the _why_, points to existing docs rather than duplicating.

2. **Skill installed and verified discoverable.**
   - `~/.config/crush/skills/templ-components/SKILL.md` → symlink to repo `skill/SKILL.md`.
   - Symlink resolves; frontmatter parses; appears in `available_skills` as `templ-components`.
   - Idempotent install command (`ln -sf`) documented in the skill itself.

3. **flake.nix apps are the canonical build/test entry points in the skill.**
   - New "Flake commands" table (`nix run .#{build,test,lint,coverage,verify}`).
   - `nix run .#verify` positioned as the single "done" check.
   - Corrects the original draft which told readers to run raw `go test`/`golangci-lint` — that would bypass the pinned `templ` v0.3.1020 and violate repo policy.

4. **Contract-test step added to "add a component" flow.**
   - Anatomy step 9: register the new Props type in `internal/contract/component_props_test.go`.
   - Cross-linked the contract test, `examples/demo/`, `integration/composition_test.go`, `FEATURES.md`, `TODO_LIST.md` in deeper reading.

5. **Pre-existing `templ fmt` nit fixed as a chore.**
   - BuildFlow's pre-commit `templ-fmt` step normalized `PageProps.HTMXCDN` field alignment in `layout/base.templ` (pre-existing, unrelated to the skill). Committed separately as a chore so the tree stays clean.

---

## b) PARTIALLY DONE

1. **Skill triggering accuracy is unverified.**
   - The `description` is deliberately "pushy" per skill-creator guidance, but I did **not** run the skill-creator's description-optimization loop (`scripts/run_loop.py` with trigger eval queries). That requires the `claude` CLI and a 20-query should-trigger / should-not-trigger eval set.
   - **Status:** description is a best-effort draft; real triggering accuracy is unmeasured.

2. **Self-critique depth.**
   - I caught the flake.nix miss, the contract-test gap, and the install weakness. I did **not** exhaustively diff every factual claim in the skill against the codebase (e.g. the "14 enums use map+fallback" count, the "51 generated files" count) — these were inherited from `AGENTS.md` and may have drifted.

---

## c) NOT STARTED

1. **No evals written for the skill** (`evals/evals.json` with test prompts + assertions).
2. **No baseline comparison run** (skill-creator's with-skill vs without-skill subagent runs).
3. **No integration with the repo's `flake.nix`** as a devShell `skills` entry (if such a mechanism even exists — I didn't explore it).
4. **No `CONTEXT.md` / `FEATURES.md` / `TODO_LIST.md` freshness check** — I referenced these docs but didn't verify they're current.
5. **No mention of the skill in `README.md` or `AGENTS.md`** — a one-line pointer would help humans discover it.

---

## d) TOTALLY FUCKED UP

Nothing this session. All 5 commits are clean, BuildFlow passes 28/28 on every commit, and `git push` succeeded.

One near-miss worth naming: the first commit's pre-commit hook surfaced a `templ fmt` diff in `layout/base.templ` that I didn't author. I committed it as a separate chore rather than reverting or ignoring it — the right call, but worth flagging that BuildFlow _will_ mutate files on commit and you have to decide what to do with the fallout.

---

## e) WHAT WE SHOULD IMPROVE

1. **Run the skill-creator's description optimizer** to measure and improve triggering accuracy — the current description is a guess, not data.
2. **Add a one-line pointer to `skill/SKILL.md` in `AGENTS.md`** so the next session discovers it via the auto-loaded context. Right now the skill only loads via the symlink, which is a machine-local install.
3. **Add a one-line pointer in `README.md`** (Contributing or a new "AI assistant" section) so human contributors know the skill exists.
4. **Audit the factual counts in `AGENTS.md`** that the skill inherits — "51 generated files", "14 enums use map+fallback", "73 components", "26 typed enums", "101 icons". These drift with every PR; the skill repeats them.
5. **Consider a drift-guard test** for the counts in README/AGENTS (like the existing `utils.TestVersionMatchesChangelog` pattern) so the skill doesn't inherit stale numbers.
6. **The skill currently has no `scripts/` or `references/` subdirectories** — the skill-creator anatomy supports them. A `references/component-template.templ` skeleton file would let the skill say "copy this" instead of describing the anatomy in prose.
7. **The skill's testing matrix table could be turned into a checklist script** that asserts a touched package has the expected test files — turning the rubric into automation.

---

## f) Up to 25 things to get done next

Ordered roughly by impact / effort ratio (highest first):

1. Add a one-line pointer to `skill/SKILL.md` in `AGENTS.md` (under a "Skills" section or similar).
2. Add a pointer in `README.md` Contributing section.
3. Run the skill-creator description optimizer (`run_loop.py`) with a 20-query trigger eval set.
4. Write `skill/evals/evals.json` with 2-3 realistic test prompts + assertions.
5. Run a baseline comparison (with-skill vs without-skill subagents) for the test prompts.
6. Create `skill/references/component-template.templ` — a copy-paste skeleton for new components.
7. Create `skill/references/new-component-checklist.md` — the anatomy steps as a literal checklist.
8. Audit "51 generated files" count in AGENTS.md against actual `*_templ.go` files.
9. Audit "73 components" count in README against actual exported components.
10. Audit "26 typed enums" / "101 icons" counts.
11. Add a drift-guard test for README/AGENTS counts (extend the `utils.TestVersionMatchesChangelog` pattern).
12. Verify `[Unreleased]` in CHANGELOG.md has an entry for the skill addition (it may not — this was a docs change).
13. Check whether the skill should live at repo root (`SKILL.md`) instead of `skill/SKILL.md` — matches how some tools discover skills.
14. Add the skill to `flake.nix` devShell if there's a skills mechanism (explore `flake-parts` options).
15. Write a `skill/references/testing-strategy.md` deep-dive on the golden/a11y/BDD/snapshot test lenses.
16. Write a `skill/references/theming-deep-dive.md` expanding the CSS-variable override model with worked examples.
17. Add a `skill/scripts/new-component.sh` that scaffolds a new component from the template + registers it in the contract test.
18. Cross-link the skill from `CONTRIBUTING.md` "Code Conventions" section.
19. Verify the skill's claim that "only `layout.PageProps` is the exception to BaseProps embedding" — scan for other exceptions.
20. Add a section to the skill on the `errorpage` package's go-error-family integration (currently only a pointer).
21. Add a section on the HTMX package's JS attachment patterns and retry/error-handling config.
22. Run `nix run .#coverage` and capture the current coverage baseline for future comparison.
23. Schedule a periodic docs-freshness-check run against the skill's claims.
24. Consider a second skill: `templ-components-consumer` aimed at _using_ the library (vs authoring in it) — different audience, different triggers.
25. Ask Lars whether the skill should auto-trigger on _every_ session in this repo (via `disable-model-invocation` / `user-invocable` metadata tuning).

---

## g) Top #1 question I can NOT figure out myself

**Should the skill auto-trigger on every session in this repo, or only when explicitly invoked?**

The skill-creator docs mention `disable-model-invocation` and `user-invocable` frontmatter fields, but I can't determine from the Crush tooling alone:

- Whether setting `user-invocable: true` (like the `go-cqrs-lite` skill uses) is the right call here — i.e., the skill only fires when the user says "use the templ-components skill", vs. auto-firing whenever a `.templ` file is touched.
- Whether auto-firing on every `.templ` edit would be annoying (too aggressive) or helpful (enforces the conventions consistently).
- What the actual discovery/loading semantics are for a repo-local skill installed via symlink — does Crush re-read it on every session, or cache it?

This is a product/taste decision about the intended workflow, not something I can derive from the code. I'd want Lars's preference: **"always load when working in this repo"** vs **"only when I ask for it"**.
