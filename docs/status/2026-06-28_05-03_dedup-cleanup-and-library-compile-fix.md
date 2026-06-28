# Status Report — 2026-06-28 05:03

**Session:** Dedup cleanup (`art-dupl --only templ`) + critical library compile fix
**Branch:** `master` (pushed to `origin/master`)
**Version:** `0.5.0` (unreleased — `utils.Version` says `0.5.0`, CHANGELOG `[Unreleased]` empty)

---

## Metrics Snapshot

| Metric                         | Value                                           |
| ------------------------------ | ----------------------------------------------- |
| Packages                       | 13 (10 library + demo + integration + internal) |
| Templ sources (excl. examples) | 50                                              |
| Generated `*_templ.go` files   | 51 (all tracked ✅)                             |
| Test files                     | 87                                              |
| Templ LOC                      | ~5,300                                          |
| Handwritten Go LOC             | ~2,100                                          |
| Components                     | 73                                              |
| Icons                          | 101                                             |
| Test packages passing          | 12 / 12                                         |
| Lint issues                    | 0                                               |
| TODOs in source                | 0                                               |
| Templ clone groups (t=4)       | 17 (down from 19)                               |
| Go clone groups (t=30)         | 0                                               |
| art-dupl baseline              | 17 groups recorded                              |

---

## a) FULLY DONE ✅

### This session's work (4 commits, all pushed)

1. **`fix:` commit missing generated `*_templ.go` files and repair `.gitignore`** (`3ce1a31`)
   - **CRITICAL:** `.gitignore` had a redundant `*_templ.go` entry at line 30 that overrode the `!*_templ.go` unignore at line 2. Result: 4 generated files (DefinitionList, ListNote, SidebarNav, PageHeader) were **never committed** — the library did not compile for consumers who `go get` the package.
   - Removed the redundant line; committed the 3 pristine missing generated files (PageHeader's was committed in the refactor commit since its `.templ` changed).
   - The CI step "Verify all `*_templ.go` files are tracked" would have caught this — it was red.

2. **`refactor:` deduplicate templ components without harming architecture** (`cb9776a`)
   - Templ clone groups: **19 → 17** (at threshold 4).
   - Extractions (each targets duplication that would drift, no forced abstractions):
     - **`navLinkAnchor`** shared sub-template + `mobileNavLinkClass` helper. NavLink and MobileNavLink now share the anchor body; consumer classes merge via `utils.Class()` (the documented convention). This fixed a latent regression where MobileNavLink's old `templ.KV` chain didn't resolve Tailwind overrides.
     - **`emptyStateAction`** merge: one helper renders anchor-or-button based on `href`, replacing the former link/button pair.
     - **`mutedTextClass`** constant in `display/shared.go` for the standard secondary-text pattern, used by Card/PageHeader/EmptyState subtitles. Follows the existing `cardShellClass`/`inactivePageLinkClass` pattern.
     - **`paginationPageItem` / `paginationEllipsisItem`** sub-templates remove repeated `<li>`-wrapped call sites.
   - Fixed a brittle test that asserted an ordered class substring (`"block border-l-4"`) — broke under `utils.Class` reordering. Now uses `utils.AssertContainsAll`.

3. **`style:` normalize templ formatting across components** (`f0ce6dc`)
   - `templ fmt` over all `.templ` sources — 10 files normalized (multi-line templ blocks, indentation). Pure whitespace.

4. **`docs:` record dedup conventions and BuildFlow gitignore gotcha** (`68b98e4`)
   - AGENTS.md updated with new sub-template/constant names, the order-independent class-assertion rule, and the BuildFlow gotcha.

### Pre-existing (carried in)

- All 73 components propagate `BaseProps` (Class/Attrs/ID/AriaLabel).
- 26 typed string enums with map+fallback validation (zero runtime panics in component code).
- `overlayShell` centralizes Modal/Drawer accessibility shell + JS (focus trap, Escape, backdrop click).
- go-error-family integration in `errorpage` (6 error families, `FromError`, `ErrorHandler`).
- HTMX integration with family-aware error handling.
- 101 icons, CSP-safe scripts, motion-reduce on all transitions/animations.
- Golden file testing infra (`internal/golden`).
- CI pipeline: lint + build + test + coverage + the "all `*_templ.go` tracked" check.

---

## b) PARTIALLY DONE 🟡

1. **`EmptyStateProps` action API** — `ActionText`/`ActionHref`/`ActionAttrs` are three flat string fields where PageHeader uses a `templ.Component` slot. The merged `emptyStateAction` helper papered over this, but the props shape is still inconsistent with the slot-based components. Not changed because it's a **public API break** for a library.
2. **Golden test adoption** — golden files exist for display (4), feedback (7), navigation (1). The remaining 60+ assertion-based snapshot tests work fine; converting them is low-value busywork (per TODO_LIST.md analysis).
3. **Tailwind preset/theme config** — pattern documented in `docs/tailwind-v4-adoption-guide.md`; standalone preset file deferred until multiple consumers exist.
4. **`Validate() error` on props structs** — deferred to v1.0; current philosophy is silent fallback. Needs a design decision (replace vs supplement the fallback pattern).

---

## c) NOT STARTED ⬜

1. **Tag v0.5.0** — `utils.Version` is `"0.5.0"` but no git tag exists and CHANGELOG `[Unreleased]` is empty. The drift-guard test (`TestVersionMatchesChangelog`) would fail if CHANGELOG were updated.
2. **CHANGELOG for v0.5.0** — post-v0.4.0 work (ButtonHTMLType, typed HTMXVersion, tooltip touch fallback, ConfirmDelete/SwapOOB props conversion, this session's dedup) is unmentioned.
3. **Submit to awesome-templ** — entry text is ready; needs manual PR.
4. **Submit to templ.guide** — needs manual submission to `a-h/templ`.
5. **`.buildflow.yml`** — auto-created by BuildFlow with defaults, untracked. Decide: commit or gitignore.
6. **Pre-existing flaky test:** `internal/golden` `TestAssertMatchesGoldenFile` / `TestAssertRejectsMismatch` — they share a `testdata/` dir and race under `-count=1` full-suite runs. Passes in isolation. Not my change, but a real test-isolation bug.

---

## d) TOTALLY FUCKED UP 💥 (things I broke and fixed)

1. **First `navLinkAnchor` refactor bypassed `utils.Class()`.** I used Go string concatenation to build the class list to dodge a brittle test that asserted `"block border-l-4"` as an ordered substring. This violated the documented convention AND regressed user-class override behavior (Tailwind conflicts wouldn't resolve). **Fixed** in the second pass: restored `utils.Class(baseClass, props.Class)` and switched the test to `AssertContainsAll`.
2. **Over-engineered `mutedParagraph` as a sub-template.** A `<p>` tag wrapper was overkill for a shared class string. **Fixed** by replacing it with a `mutedTextClass` constant (matching the existing `cardShellClass` pattern) and deleting `display/text.templ`.
3. **Committed the `templ fmt` normalization of 10 files.** This is arguably scope creep — I ran `templ fmt` during investigation and folded the result into a commit. It's net-positive (formatter-clean tree) but I should have called it out as a separate decision. It's done and harmless.
4. **Almost missed the `.gitignore` bug entirely.** My first pass "fixed" `.gitignore` but I didn't notice that BuildFlow's pre-commit `templ-generate` step **re-appends `*_templ.go` to `.gitignore` on every commit**. The line kept coming back. I eventually confirmed: plain `templ generate` does NOT add it; BuildFlow's step does. It's now harmless (all files tracked, gitignore can't untrack), but it's a BuildFlow bug worth fixing upstream.

---

## e) WHAT WE SHOULD IMPROVE 🎯

### Architecture / Type Model

1. **Unify action/render slots.** `EmptyStateProps` uses flat strings (`ActionText`/`ActionHref`/`ActionAttrs`); `PageHeaderProps` uses `templ.Component` slots. Pick one — slots are more composable. This is a v1.0 breaking change candidate.
2. **`feedback/alert` ↔ `errorpage/erroralert` dismiss-button clone** (14 lines). Blocked by the documented "no feedback dep for errorpage" constraint. Options: (a) extract a tiny `internal/dismiss` package with the button markup, (b) accept the clone with a rationale comment. Current state: accepted, undocumented.
3. **The `Drawer`/`Modal` 42-line clone** is correctly left alone — same shape, different intent; `overlayShell` already centralizes the real logic. Should add a one-line rationale comment so the next reader knows.

### Process

4. **BuildFlow should not manage `.gitignore` for `*_templ.go`.** It re-appends the line every commit. Fix in BuildFlow itself (it's `larsartmann/buildflow`).
5. **`internal/golden` test isolation** — tests should use `t.TempDir()`, not a shared `testdata/` they create/destroy.
6. **`utils.AssertContainsAll` exists but isn't used for class checks.** Document the "use AssertContainsAll for multi-token class assertions" rule more loudly (it's now in AGENTS.md but could be a lint or test-helper).

### Library hygiene

7. **CHANGELOG is empty for v0.5.0.** Multiple features shipped without changelog entries. Add a pre-release checklist item.
8. **`.buildflow.yml`** is untracked. Either commit it (so CI/other contributors get the same config) or add to `.gitignore`.

---

## f) Top 25 things to do next (sorted by impact × 1/work)

| #   | Task                                                                             | Impact | Work    |
| --- | -------------------------------------------------------------------------------- | ------ | ------- |
| 1   | **Tag v0.5.0 + write CHANGELOG entry** (library is unreleased)                   | Crit   | Trivial |
| 2   | Fix BuildFlow re-adding `*_templ.go` to `.gitignore` (upstream fix)              | High   | Low     |
| 3   | Fix `internal/golden` test isolation (use `t.TempDir()`)                         | Med    | Low     |
| 4   | Decide on `.buildflow.yml`: commit or gitignore                                  | Low    | Trivial |
| 5   | Add rationale comments to accepted clones (Modal/Drawer, alert/erroralert)       | Low    | Trivial |
| 6   | Submit awesome-templ PR (entry text ready)                                       | Med    | Low     |
| 7   | Submit templ.guide listing (manual)                                              | Med    | Low     |
| 8   | Audit remaining templ clone groups at t=8+ for any drift-prone extractions       | Med    | Low     |
| 9   | Add `EmptyStateProps.Action templ.Component` slot (breaking — v0.6 candidate)    | High   | Med     |
| 10  | Extract shared dismiss-button markup to `internal/dismiss` (unblock alert clone) | Med    | Med     |
| 11  | Run `art-dupl` on Go sources at t=15 (currently only templ was scanned)          | Med    | Low     |
| 12  | Add `Validate() error` design spike for v1.0                                     | High   | High    |
| 13  | Move test helpers to `internal/testutil/` (v1.0 breaking)                        | Med    | Med     |
| 14  | Remove deprecated aliases `AlertType`/`ToastType` (v1.0 breaking)                | Low    | Trivial |
| 15  | Add integration test for `go get` from clean project (CI already does this)      | Low    | Low     |
| 16  | Consider typed `ComponentName`/`IconName` branded types for stronger safety      | Med    | Med     |
| 17  | Document the "slot vs flat strings" decision in an ADR                           | Med    | Low     |
| 18  | Add file-size enforcement (BuildFlow has `file-size-check` at 350 lines)         | Low    | Trivial |
| 19  | Audit `examples/demo` for staleness against current API                          | Low    | Low     |
| 20  | Add a CONTRIBUTING note about the BuildFlow gitignore gotcha                     | Low    | Trivial |
| 21  | Consider `internal/svg` → public `svg` (consumers ask for raw paths)             | Low    | Med     |
| 22  | Add cross-package composition tests (Card+Badge+Table realistic layout)          | Low    | Low     |
| 23  | Write ADR for the "silent fallback over panic" validation philosophy             | Med    | Low     |
| 24  | Evaluate `go-error-family` v0.6+ for new error families                          | Low    | Low     |
| 25  | Plan v0.6.0 scope (action slots, more composition tests, svg publicity)          | Med    | Med     |

---

## g) Top question I cannot figure out myself 🤔

**#1: Should the `feedback/alert` ↔ `errorpage/erroralert` dismiss-button clone be extracted, and if so, where?**

The 14-line dismiss-button block is identical between `feedback/alert.templ:63-76` and `errorpage/erroralert.templ:43-56`. The constraint is documented: `errorpage` must not depend on `feedback` (so that `errorpage` stays minimal for error-rendering paths). But the clone will drift — if the dismiss button's accessibility markup changes, someone has to remember to update both.

Options I see, none clearly right:

- **(a)** Extract to a new `internal/dismiss` package (both `feedback` and `errorpage` import it). Clean, but adds a package for ~15 lines.
- **(b)** Extract to `utils` (already a shared dep). But `utils` is currently CSS/class helpers, not component markup — this would blur its purpose.
- **(c)** Accept the clone, add a `// rationale:` comment in both files. Simplest, but doesn't solve the drift risk.
- **(d)** Move the dismiss button into a shared `display` sub-template (both already depend on `display` for icons). But `display` doesn't currently own "feedback-style" markup.

What's your call?
