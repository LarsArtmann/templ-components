# Status Report — ErrorPage Redesign & Error Model Session

**Date:** 2026-09-17 15:27 CEST
**Session scope:** "Our Error Page is ugly. And do we have Error Models?" → research, redesign, harden, verify.
**Repo:** templ-components @ master (daemon-committing throughout; my work landed in heuristic commits `3fecc0fb`, `7dcc0d92`, `af96082a` and successors).

---

## Executive Summary

The errorpage package **had** a complete, well-typed error model — but the page hid it: the default render showed no title, `StatusCode` was validated yet never displayed, and the family-tinted card washed out in both themes. This session rebuilt the ErrorPage visual design around a neutral card + family accents, made the model visible (HTTP status/code/family chips), fixed an empty-badge rendering bug, repaired four doc-drift spots with a new enforcement guard, and re-verified the entire matrix (7 modules + website + visual/axe suite + 250 HTML-valid goldens). Nothing is broken; one demo-structure oddity was noticed and deliberately left for a decision.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | **Error-model inventory answered** — two layers documented: `errorpage` model (Family ×6, Code, 3 Props types with `Validate()`, CauseItem, ContextPair, FamilyStatusCode, ParseFamily, 6 constructors) + go-error-family bridge (`FromError`, `ErrorHandler`, `WriteError`, JSON mode) | Reported in-session; code at `errorpage/{styles,fromerror,constructors,handler}.go` |
| 2 | **ErrorPage visual redesign** — neutral white card (`rounded-2xl`, dark-mode shadow), family accent bar (`familyVisualStyle.Bar` ×6 + default), 48px icon circle, chip row rendering **`HTTP {StatusCode}`** (previously validated but never rendered), `text-2xl` title, lead message, why paragraph, restyled fix card (wrench icon, uppercase label), neutral context table, cause chain, 404-scale action button, tinted timestamp footer | `errorpage/errorpage.templ`, `errorpage/styles.go`, `errorpage/shared.templ` |
| 3 | **Empty-badge bug fixed** — zero/unknown `Family` used to emit an empty tinted `<span>`; badge (and whole chip row when unset) now renders nothing | `familyBadge`/`errorChips` guards in `shared.templ`; `edge_cases_test.go` assertion flipped to `AssertNotContains("bogus")` |
| 4 | **Shared-template refactor** — `familyIcon(style, iconClass)` size param; `fixCard`/`contextTable` take inset surface (`errorInsetNeutral` white-card / `errorInsetCard` tinted-card, `styles.go`); `diagnosticSection` rhythm via `space-y`; deleted now-redundant `codeAndFamilyBadge` | Zero stale references (swept); all callers regenerated |
| 5 | **New HTML golden coverage** — `TestGoldenSweepErrorPage` (`error_page_full`, `error_page_minimal` incl. go-back script branch); corpus 248 → **250** | `errorpage/golden_sweep_test.go`, `errorpage/testdata/` |
| 6 | **Visual goldens regenerated + inspected** — light & dark PNGs re-captured from full-model props (`TestErrorPage`/`TestErrorPageDark` now exercise the complete API); both verified by eye: dramatic improvement, coherent dark mode | `visualtest/visual_test.go`, `visualtest/testdata/errorpage/{light,dark}.png` |
| 7 | **Demo showcases the full model** — ErrorPage section now renders StatusCode/Code/Title/Message/Why/Fix/WayOut/Context/CauseChain/Timestamp | `examples/demo/errorpage_demo.templ` |
| 8 | **Docs drift repaired** — FEATURES enum table gained `Orchestration`; `TrendWarn` and `LiveOff` documented (4 spots incl. StatCard/LiveRegion rows); README errorpage section rewritten (broken `FamilyNotFound` example replaced); golden counts 248→250 in FEATURES/AGENTS/ROADMAP; CHANGELOG `[Unreleased]` entry (Changed + Fixed) | `FEATURES.md`, `README.md`, `AGENTS.md`, `ROADMAP.md`, `CHANGELOG.md` |
| 9 | **New enforcement guard** — `TestFeaturesEnumValuesExhaustive` (utils) keeps enum VALUES in FEATURES honest; all-or-nothing row resolution, `*Default`/`*Unspecified` exempt; proven red→green (caught 2 real drifts first run) | `utils/features_enum_test.go:88` |
| 10 | **Website generated-drift fix** — my pinned-binary regen flipped `website/internal/pages/base_templ.go` import back to source truth (`encoding/json`, v2→v1 leftover from the 2026-09-14 daemon incident); website goldens + link checker green after | `git show 3fecc0fb -- website/`, website tests ok |
| 11 | **Full verification matrix green** — errorpage suite; all 6 sub-modules; root module; website module; errorpage coverage **71.5%** (≥70 gate); HTML validation **250 goldens clean**; full visual suite incl. axe sweep (74s); lint errorpage+root+utils **0 issues**; `nix flake check` (treefmt) green; templ-sync / replace-directives / version-sync / docs-count guards green | Session command log |

## b) PARTIALLY DONE

1. **ErrorDetail polish** — received the badge guard, inset parametrization, and spacing fix, but its compact family-tinted design is unchanged (deliberate scope cut; no accent-bar/neutral variant parity with ErrorPage).
2. **ErrorAlert family completeness** — docs now say 6 color schemes, but the demo still shows only 5 (no orchestration alert) and no orchestration-only visual golden exists.
3. **`ci-repro.sh --lint` single-command proof** — the run executed but its tail was consumed by a git-diff dump (visualtest/go.sum churned mid-run); I verified every CI lane individually instead of obtaining one clean script verdict.
4. **Errorpage demo section presentation** — the demo wraps full-page components (`min-h-screen` `<main>`) inside bordered boxes; functional but visually odd (nested main, per-section viewport-height). Noticed, not redesigned — needs an owner decision.
5. **AGENTS.md guard table** — the new enum-values guard is implemented but the guard inventory table in AGENTS.md/skill was **not** updated with a row for it (discovered during self-review; see (e)).

## c) NOT STARTED (session-adjacent, identified but untouched)

1. Mobile (375px) and RTL visual goldens for the redesigned ErrorPage.
2. Browser-level e2e proof for the go-back button (`history.back()`) — string test only.
3. `ExampleErrorPage` godoc example refreshed to the full model.
4. ErrorDetail/ErrorAlert have **no** pixel-level visual goldens at all (pre-existing gap).
5. `FromError` auto-populating `StatusCode` from `FamilyStatusCode` (handler does this; `FromError` does not).
6. Handler-level HTML goldens reflecting the new markup.

## d) TOTALLY FUCKED UP

Nothing destroyed, no data loss, no red tests left behind. Two hazards navigated, worth recording:

1. **Daemon raced the work three times** — heuristic commits `3fecc0fb` (19 files, mid-flight, including the *flipped* `base_templ.go` state before my regen corrected it), `7dcc0d92` (8 files: docs counts, golines fix, demo CSS, visual go.sum churn, errorpage PNGs), `af96082a` (14 files: CHANGELOG/README/FEATURES/guard test + PNG re-diff). Net tree state is correct — verified after each snapshot — but the history is a jumble of partial states, exactly the AGENTS-documented daemon pattern.
2. **`visualtest/go.sum` churned during test runs** (testify 1.11.1→1.12.1, difflib dropped) and got committed by the daemon inside my snapshot. Benign (suite green before and after) but unreviewed-by-me; flagging for a tidy check.

**Foreign work noticed, left untouched (per safety rules):** uncommitted edits to `AGENTS.md` (SidebarNav dark-mode exemption removal) and `examples/demo/recipes_demo.go` (migration to `forms.ValidationError`) belong to a separate workstream; both verified consistent with committed code (`forms.ValidationError` exists; `tc-sidebar-*` tokens present; workspace builds).

## e) WHAT WE SHOULD IMPROVE — self-review: forgotten, could-do-better, still-improvable

**What I forgot:**
- The AGENTS.md/skill **guard inventory table** row for `TestFeaturesEnumValuesExhaustive` — the repo rule is "new cross-cutting rule ⇒ document it where guards are listed," not just implement it.
- `ExampleErrorPage` godoc example still shows the old minimal shape.
- A **before/after image pair** in this report (the PNGs exist; the report should carry them).
- Early consideration of **demo presentation** — I updated demo props but didn't question the bordered-box-wrapping of full-page components until late.
- Capturing the **ci-repro verdict explicitly** (`tee` + exit code) — losing the verdict line cost a manual re-verification pass.

**What I could have done better:**
- **Commit per task when the daemon is active** — my changes landed as three heuristic blobs; explicit commits would have kept history reviewable.
- Verify **ErrorDetail/ErrorAlert in a browser**, not only as HTML strings — they got zero new pixel proof; the 404-page lesson says string-green ≠ browser-green.
- Add the **mobile visual golden in the same pass** — the viewport was only ever the default; the chip row's `flex-wrap` path is unproven at 375px.
- Timebox the **enum-guard design** — the first naive version produced 111 false positives; the suffix-convention survey should have preceded implementation.

**What could still be improved (component):**
- ErrorPage secondary-action slot; configurable max-width; `FromError` status-code parity; errorpage docs page on the website (none exists).

## f) Up to 50 things we should get done next

*Brainstorm, not commitment — ROADMAP/TODO_LIST fuel. ★ = do soon.*

**Errorpage components**
1. ★ Update AGENTS.md + skill guard tables with `TestFeaturesEnumValuesExhaustive`.
2. ★ Refresh `ExampleErrorPage` to the full model.
3. ★ Errorpage demo: render full-page components via standalone routes (`ErrorHandler`) instead of bordered boxes.
4. Add orchestration ErrorAlert to demo + family matrix golden (6/6).
5. Mobile viewport visual golden for ErrorPage (375px).
6. RTL visual golden for ErrorPage (chip row, footer).
7. Visual goldens for ErrorDetail and ErrorAlert (none exist).
8. Browser e2e for the WayOut go-back button.
9. `FromError` sets `StatusCode` via `FamilyStatusCode` (parity with handler).
10. Handler HTML goldens (`WriteError`/`HTMLShell`) reflecting new markup.
11. Secondary (ghost) action slot on ErrorPage (e.g., "Contact support").
12. Configurable card width (`max-w-*`) for context-heavy pages.
13. ErrorDetail neutral/accent-bar variant for parity with ErrorPage.
14. Copy-button composition for `Code` chip (debugging UX).
15. Unify ErrorPage/NotFound404 button-class construction into a shared constant.
16. Consider typed `WayOutAction` struct (text+href pair) replacing two loose strings.
17. Fuzz test for `ParseFamily`.
18. Benchmark for `ErrorPage` render if missing.
19. Coverage push: errorpage 71.5% → 75%+ (chips guards, fmt branch).
20. Document `Code` open-enum policy (no IsValid by design?).
21. Website docs page for errorpage package (none found).
22. Check `docs/recipes/server-rendered-htmx-error-feedback.md` for stale visuals references.
23. Grep for leftover "Wix-style"/old-design phrasing anywhere.
24. JSON `errorResponse` ↔ chips data parity guard.

**Guards/tooling (session-observed)**
25. ★ `go mod tidy` + review `visualtest/go.sum` churn (testify bump) — decide pin policy.
26. ci-repro.sh: explicit verdict line + exit code; optionally `--quiet-diff`.
27. Flake app `visual-update <pattern>` for targeted golden regeneration.
28. Investigate daemon/templ-watcher race that resurrects `base_templ.go` import flips (BuildFlow upstream).
29. Extend check-templ-sync.sh to assert website module generated files explicitly.
30. Consider single-sourcing the golden count (3 docs carry 250) behind the guard only.
31. Add before/after screenshot pairs to visual-test failure artifacts naming convention.

**Docs**
32. AGENTS.md "26+ props structs" phrasing vs 121-component reality — refresh stale counts prose.
33. docs/DOMAIN_LANGUAGE.md: add Family/CauseItem/ContextPair/WayOut terms if absent.
34. Skill SKILL.md ErrorPage one-liner: mention status chip + redesign.
35. ROADMAP: route (f) items through docs-health HARVEST with routing rigor.

**Repo hygiene noticed this session (not researched further, per scope)**
36. templ-LSP QF1003 infos in `collapsible_section.templ`/`animated_icon.templ` — recurring noise, either fix or suppress.
37. `docs_count_test` claims vs FEATURES "72.3% coverage" line — recompute at next release.
38. Confirm errorpage demo section stays axe-clean after any demo restructuring (#3).
39. Verify pre-commit hook installed on fresh clones (`scripts/setup-hooks.sh`) post-daemon churn.
40. Read `docs/release-checklist.md` before next cut; CHANGELOG `[Unreleased]` is warm.

**Follow-ups on foreign workstream (only with owner's go-ahead)**
41. SidebarNav theming change (`AGENTS.md` edit in tree) — needs its own tests/goldens pass.
42. Demo migration to `forms.ValidationError` — check remaining demo files for the old `formsValidationError`.
43. Server-side-validation recipe doc update to match `forms.ValidationError`.

**Nice-to-have**
44. ErrorPage `MaxMismatch`/viewport options audit for new overlays.
45. Consider `aria-describedby` wiring from chips to fix card (screen-reader grouping).
46. `ErrorPageProps.Validate()` could warn (not error) when StatusCode set without Title.
47. Golden-file diff tooling: LCS line numbers already exist — expose in `-update` output summary.
48. Demo: live error-page playground (pick family/status, render page).
49. Consider printing rendered-props table into errorpage package doc.go.
50. Retire `errBlankNonRejection` alias indirection if still accurate (naming hygiene check).

## g) Questions I cannot figure out myself

1. **Ownership of the foreign working-tree edits** — the SidebarNav theming (`AGENTS.md` + `templates/custom.css`) and `forms.ValidationError` demo migration are uncommitted and not mine. Are they an active parallel session's work I should keep avoiding, and who lands them?
2. **Demo architecture decision** — should full-page components (ErrorPage/NotFound404, `min-h-screen` `<main>`) be demoed inside bordered boxes (today) or via standalone routes through `ErrorHandler`? This changes the demo's structure and the axe/route-golden surface, so I want the intent before restructuring.
3. **"Retry" semantics** — `ServiceUnavailable()` renders WayOut "Retry" with `WayOutHref: "/"`. Is sending users home the intended product behavior, or should retry re-request the current URL (a component-level pattern worth adding)?

---

*Point-in-time snapshot; goes stale. Section (f) is HARVEST input for `TODO_LIST.md`/`ROADMAP.md` via docs-health when instructed.*
