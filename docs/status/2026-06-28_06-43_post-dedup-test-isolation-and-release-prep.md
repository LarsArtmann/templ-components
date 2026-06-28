# Status Report — 2026-06-28 06:43

**Session:** Post-dedup cleanup — golden test isolation fix, accepted-clone documentation, CHANGELOG backfill, clone audit
**Branch:** `master` (up to date with `origin/master`; this session's work uncommitted — 13 files changed)
**Version:** `0.5.0` tagged (`e72338c`); **27 commits unreleased** — CHANGELOG `[Unreleased]` now documents them

---

## Metrics Snapshot

| Metric                         | Value                                                        |
| ------------------------------ | ------------------------------------------------------------ |
| Packages                       | 13 (10 library + demo + integration + internal)              |
| Templ sources (excl. examples) | 50                                                           |
| Generated `*_templ.go` files   | 51 (all tracked ✅)                                          |
| Test files                     | 87                                                           |
| Templ LOC                      | ~5,336                                                       |
| Handwritten Go LOC             | ~2,119                                                       |
| Components                     | 75                                                           |
| Icons                          | 101 (100 path constants + Spinner)                           |
| Test packages passing          | 12 / 12 (`-race -count=1`)                                   |
| Lint issues                    | 0                                                            |
| TODOs in source                | 0                                                            |
| Templ clone groups (t=4)       | 17 (baseline — 0 new)                                        |
| Go clone groups (t=15)         | 0                                                            |
| Commits since v0.5.0           | 27 (unreleased)                                              |
| Breaking changes unreleased    | 3 (FormFieldWrapper, ConfirmDelete, SwapOOB → Props structs) |

---

## a) FULLY DONE ✅

### This session's work (13 files changed, uncommitted)

1. **`fix:` golden test isolation — eliminate the `testdata/` race**
   - `internal/golden.Assert` now delegates to a directory-parameterized `assertInDir(t, dir, name, got)`. The production `Assert` still passes `"testdata"` (zero change for the 12 golden callers in display/feedback/navigation); the package's own unit tests pass `t.TempDir()`.
   - The old tests created/destroyed a shared `testdata/` dir with `os.RemoveAll` under `t.Parallel` — a real race that caused flaky failures under `-count=1` full-suite runs. Now each test gets an isolated temp dir that Go cleans up automatically.
   - Verified: `-race -count=5` passes cleanly; all golden callers (`./display/... ./feedback/... ./navigation/...`) pass.

2. **`docs:` accepted-clone rationale — resolves the "top question" from the last report**
   - Three accepted clone groups now carry explicit rationale in **doc comments** (not in-markup comments — see section d):
     - **`feedback.Alert` ↔ `errorpage.ErrorAlert`** dismiss-button (13 lines). Decision: **accept.** `errorpage` must not import `feedback`; the shared _logic_ (`utils.DismissScript()` JS) is already shared; only the _markup_ is duplicated; no shared markup package exists; the project's own "3+ times" threshold isn't met (2×); a package for 13 lines is premature generalization.
     - **`display.Modal` ↔ `display.Drawer`** panel body (~15 lines of scaffold). Decision: **accept.** The transition classes differ fundamentally (scale/opacity vs translate-x); sizing classes differ; `templ.KV` conditionals are inline; `overlayShell` already owns the shared a11y shell + JS.
     - **`errorpage.ErrorDetail` ↔ `errorpage.ErrorPage`** header scaffold (icon + badge + title + message). Decision: **accept.** Shared logic already extracted (`familyIcon`, `codeAndFamilyBadge`, `diagnosticSection`); the remaining layout differs by visual density (compact `h3` vs full-page `h1`).

3. **`docs:` CHANGELOG `[Unreleased]` backfill — documents all 27 commits since v0.5.0**
   - Added: tooltip touch support + auto-ID, typed `HTMXVersion`, hex color validation, size constants, Toggle fields, `ConfirmDelete`/`SwapOOB` BaseProps, `ErrorHandlerConfig.Lang`.
   - Changed (3 breaking): `FormFieldWrapper` → `FormFieldProps`, `ConfirmDelete` → `ConfirmDeleteProps`, `SwapOOB` → `SwapOOBProps`; handler split; buffer-before-response rendering; Drawer Tailwind classes; `PageProps.HTMXVersion` type.
   - Fixed: missing `*_templ.go` compile break, Button `aria-disabled`, Spinner `role="img"`, Avatar status dot, `errors.Join` in cause chain.
   - Internal: dedup extractions, golden test isolation fix, clone rationale docs.
   - Drift-guard test (`TestVersionMatchesChangelog`) still passes — `Version = "0.5.0"` matches the latest heading `[0.5.0]`.

4. **`chore:` clone audit — art-dupl on templ (t=8) + Go (t=15)**
   - Templ at t=8: **4 groups** — all either documented as accepted (dismiss button, modal/drawer) or not actionable (errordetail/errorpage now documented; demo examples not shipped).
   - Go at t=15: **0 groups** — clean.
   - Baseline check: **0 new clones** vs the 17-group baseline.

### Pre-existing (carried in)

- v0.5.0 tagged (`e72338c`, `release: v0.5.0 — type safety fixes + version single source of truth`).
- `.buildflow.yml` committed (BuildFlow CI config).
- All 75 components propagate `BaseProps` (Class/Attrs/ID/AriaLabel).
- 26 typed string enums with map+fallback validation (zero runtime panics in component code).
- `overlayShell` centralizes Modal/Drawer accessibility shell + JS (focus trap, Escape, backdrop click).
- go-error-family integration in `errorpage` (6 error families, `FromError`, `ErrorHandler`).
- HTMX integration with family-aware error handling.
- 101 icons, CSP-safe scripts, motion-reduce on all transitions/animations.
- Golden file testing infra (`internal/golden`).
- CI pipeline: lint + build + test + coverage + the "all `*_templ.go` tracked" check.
- Pre-commit hook via BuildFlow.

---

## b) PARTIALLY DONE 🟡

1. **CHANGELOG `[Unreleased]` is written but unreleased** — the section documents 27 commits (3 breaking) but no v0.6.0 tag exists. Ready to tag whenever the decision is made (see section g).
2. **`EmptyStateProps` action API** — `ActionText`/`ActionHref`/`ActionAttrs` are three flat string fields; `PageHeaderProps` uses `templ.Component` slots. The merged `emptyStateAction` helper papered over this, but the props shape is still inconsistent. v0.6 breaking-change candidate.
3. **Golden test adoption** — golden files exist for display (4), feedback (7), navigation (1). The remaining 60+ assertion-based tests work fine; converting them is low-value busywork.
4. **Accepted-clone documentation** — **DONE this session** (previously "accepted, undocumented" per the last report's section g).
5. **Tailwind preset/theme config** — pattern documented in `docs/tailwind-v4-adoption-guide.md`; standalone preset file deferred until multiple consumers exist.
6. **`Validate() error` on props structs** — deferred to v1.0; current philosophy is silent fallback. Needs a design decision (replace vs supplement the fallback pattern).

---

## c) NOT STARTED ⬜

1. **Tag v0.6.0** — 27 commits unreleased, including a CRITICAL compile fix (missing `*_templ.go` files that broke `go get` for consumers on v0.5.0) and 3 breaking API changes. CHANGELOG is ready.
2. **Create `ROADMAP.md`** — referenced in AGENTS.md's "Project Documentation Files" table but does not exist.
3. **Submit to awesome-templ** — entry text ready; needs manual PR.
4. **Submit to templ.guide** — needs manual submission to `a-h/templ`.
5. **Document the templ `//` comment whitespace gotcha** — discovered this session (see section d#1); should be added to AGENTS.md so the next person doesn't repeat it.
6. **`EmptyStateProps.Action templ.Component` slot** — v0.6 breaking candidate.
7. **`Validate() error` design spike** — v1.0.
8. **Remove deprecated aliases `AlertType`/`ToastType`** — v1.0 breaking.
9. **Update `TODO_LIST.md`** — last updated 2026-06-27; stale relative to this session's work.

---

## d) TOTALLY FUCKED UP 💥 (things I broke and fixed)

1. **`//` comments inside templ markup leak whitespace into rendered HTML.** I initially placed the rationale comments _inside_ the templ component body (between `if props.Dismissible {` and the `<div class="ml-auto pl-3">`). This is syntactically valid templ — `//` comments are allowed in Go blocks — but templ renders the surrounding whitespace/newlines into the output. The `TestGoldenAlertDismissible` golden test caught it immediately: the rendered HTML had extra whitespace before the dismiss button. **Fixed** by moving all rationale to **doc comments above the `templ` function declarations** (doc comments are stripped from output — this is why every existing `// Foo renders...` comment doesn't leak). This is a non-obvious gotcha that should be documented in AGENTS.md.
2. **The in-markup comment version was caught by a golden test — which is exactly why golden tests exist.** No user-facing damage; the failure was local and immediate. But it reinforces: when you have golden coverage, trust it. When you don't (like `erroralert.templ` and `errordetail.templ` which have no golden tests), a whitespace leak would have shipped silently. The doc-comment approach is correct for all 5 files regardless of test coverage.

---

## e) WHAT WE SHOULD IMPROVE 🎯

### Architecture / Type Model

1. **Unify action/render slots.** `EmptyStateProps` uses flat strings (`ActionText`/`ActionHref`/`ActionAttrs`); `PageHeaderProps` uses `templ.Component` slots. Slots are more composable. v0.6 breaking-change candidate.
2. **`ROADMAP.md` is missing.** AGENTS.md's "Project Documentation Files" table lists it as the home for "long-term direction and raw ideas not yet refined into actionable tasks." It doesn't exist. Either create it or remove the reference.
3. **The accepted clones are now documented — but the "should we extract to `internal/dismiss`?" question will recur.** The rationale comments explain _why_ we accepted, but there's no ADR encoding the decision. A short ADR (`docs/adr-002-accepted-clones.md`) would prevent the next reviewer from re-litigating.

### Process

4. **Templ `//` comment whitespace gotcha** — discovered this session. Comments inside markup render whitespace. Should be documented in AGENTS.md alongside the existing templ conventions.
5. **BuildFlow re-appends `*_templ.go` to `.gitignore` on every commit** (line 32). Harmless (all 51 generated files are tracked; gitignore can't untrack), but confusing for new contributors. Upstream fix needed in `larsartmann/buildflow`.
6. **27 commits unreleased.** The compile fix (missing `*_templ.go`) is critical for consumers who `go get`'d v0.5.0. Every day without v0.6.0 is a day consumers get a broken build.

### Library hygiene

7. **`TODO_LIST.md` is stale** (last updated 2026-06-27). Doesn't reflect this session's golden fix, clone docs, or CHANGELOG backfill.
8. **awesome-templ / templ.guide submissions** still pending — increases project discoverability.

---

## f) Top 25 things to do next (sorted by impact × 1/work)

| #   | Task                                                                                    | Impact | Work    |
| --- | --------------------------------------------------------------------------------------- | ------ | ------- |
| 1   | **Tag v0.6.0** (27 commits unreleased incl. CRITICAL compile fix; CHANGELOG ready)      | Crit   | Trivial |
| 2   | Document templ `//` comment whitespace gotcha in AGENTS.md                              | High   | Trivial |
| 3   | Create `ROADMAP.md` (referenced in AGENTS.md but missing)                               | Med    | Trivial |
| 4   | Update `TODO_LIST.md` (stale — last updated 2026-06-27)                                 | Med    | Trivial |
| 5   | Fix BuildFlow re-adding `*_templ.go` to `.gitignore` (upstream fix)                     | High   | Low     |
| 6   | Submit awesome-templ PR (entry text ready)                                              | Med    | Low     |
| 7   | Submit templ.guide listing (manual)                                                     | Med    | Low     |
| 8   | Write ADR for accepted-clone decisions (`docs/adr-002-accepted-clones.md`)              | Med    | Low     |
| 9   | Update `FEATURES.md` with post-v0.5.0 features (tooltip touch, typed HTMXVersion, etc.) | Med    | Low     |
| 10  | Add `EmptyStateProps.Action templ.Component` slot (breaking — v0.6)                     | High   | Med     |
| 11  | Write ADR for "slot vs flat strings" decision                                           | Med    | Low     |
| 12  | Add `Validate() error` design spike for v1.0                                            | High   | High    |
| 13  | Remove deprecated aliases `AlertType`/`ToastType` (v1.0 breaking)                       | Low    | Trivial |
| 14  | Move test helpers to `internal/testutil/` (v1.0 breaking)                               | Med    | Med     |
| 15  | Consider typed `IconName`/`ComponentName` branded types for stronger safety             | Med    | Med     |
| 16  | Add file-size enforcement (BuildFlow `max_file_size: 350` already set)                  | Low    | Trivial |
| 17  | Audit `examples/demo` for staleness against current API                                 | Low    | Low     |
| 18  | Add CONTRIBUTING note about the BuildFlow gitignore gotcha                              | Low    | Trivial |
| 19  | Consider `internal/svg` → public `svg` (consumers ask for raw paths)                    | Low    | Med     |
| 20  | Add cross-package composition tests (Card+Badge+Table realistic layout)                 | Low    | Low     |
| 21  | Write ADR for the "silent fallback over panic" validation philosophy                    | Med    | Low     |
| 22  | Evaluate `go-error-family` v0.6+ for new error families                                 | Low    | Low     |
| 23  | Run `art-dupl` on Go sources at t=10 (currently clean at t=15)                          | Low    | Trivial |
| 24  | Add integration test for `go get` from clean project (CI already does this)             | Low    | Low     |
| 25  | Plan v0.6.0 scope formally (action slots, composition tests, svg publicity)             | Med    | Med     |

---

## g) Top question I cannot figure out myself 🤔

**#1: Should we tag v0.6.0 immediately, or batch more features first?**

Here's the tension:

- **Argument for tagging NOW:** v0.5.0 has a **critical compile bug** — 4 generated `*_templ.go` files (DefinitionList, ListNote, SidebarNav, PageHeader) were missing from the tag. Anyone who `go get github.com/larsartmann/templ-components@v0.5.0` gets `undefined` errors. The fix (commit `3ce1a31`) is sitting unreleased. Every consumer on v0.5.0 is broken. This alone justifies an immediate patch release.

- **Argument for WAITING:** the 27 unreleased commits include **3 breaking API changes** (`FormFieldWrapper`, `ConfirmDelete`, `SwapOOB` → Props structs). Tagging v0.6.0 means consumers must update their code to get the compile fix. That's friction — they can't just bump the version. In semver pre-1.0, breaking changes in a minor bump are acceptable, but bundling a critical fix with breaking changes forces consumers into an all-or-nothing upgrade.

- **The third option:** tag a **v0.5.1 patch** containing ONLY the compile fix (cherry-pick `3ce1a31`), then tag **v0.6.0** with the breaking changes later. But the compile fix commit (`3ce1a31`) also touches `.gitignore` and is intertwined with the dedup work — cherry-picking cleanly onto v0.5.0 may not be trivial.

What's your call? Tag v0.6.0 now (breaking + fix together), or try to cut a v0.5.1 patch first?
