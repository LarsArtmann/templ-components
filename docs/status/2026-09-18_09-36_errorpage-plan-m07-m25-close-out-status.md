# Status — ErrorPage Pareto Plan Execution (M07–M25 close-out)

**2026-09-18 09:36** · Continuation of `2026-09-17_19-59` (kickoff) and
`2026-09-18_07-06` (M01–M06). This session executed the plan's remaining
waves: M07–M17 (partial), M18 (partial), M19–M25.

## a) Fully done (this session)

- **M07** — `ErrorDetail`/`ErrorAlert` visual shield (4 PNGs,
  `errorpage/{detail,alert}_{light,dark}.png`, eyeballed) +
  `TestGoldenHandlerHTMLShell` pinning the exact `<!doctype html>` document
  `ErrorHandler` emits with `HTMLShell: true` (pinned timestamp).
- **M08** — `TestDemoErrorPageGoBack` (chromedp: click `data-tc-go-back` →
  `history.back()` proof; retry-until-needle because the navigation kills
  the poll's JS context), `TestJSONTraceContract` (trace present /
  omitempty), `TestChipsJSONParity` (HTTP chip ↔ status code, code chip ↔
  JSON `code`, trace footer ↔ JSON `trace`).
- **M09** — `FromError` sets `StatusCode = FamilyStatusCode(family)`;
  `TestFromErrorStatusCodePerFamily` + corruption fallback; goldens
  regenerated (HTTP chip now appears in FromError-driven renders).
- **M10** — `ErrorDetailVariant` (`Tinted` default / `Neutral` accent-bar
  shell, ErrorPage parity), `ErrorDetailVariantIsValid` + test, HTML +
  pixel goldens both themes, doc.go table + FEATURES + skill rows.
- **M11** — `SecondaryWayOut`/`SecondaryWayOutHref` ghost action slot
  (family-tinted `ActionButtonGhost` for all 6 families + fallback), link
  and go-back variants, goldens, demo full-model shows "Contact support".
- **M12** — `WayOutAction{Text, Href}` typed bundle (wins entirely, no
  mixing, pinned by `TestResolvedWayOut`) + `ErrorMaxWidth` enum
  (LG/XL default/2XL/4XL, map+fallback), IsValid + test, 4xl golden.
- **M13** — `CopyCode` clipboard button on the code chip (CSP-safe,
  shares display.CopyButton's `data-tc-copy` contract + `tcCopyAttached`
  guard so one listener serves both; documented accepted clone per
  ADR-0009); `errorActionClassScaffold` unifies ErrorPage/NotFound404
  action-button geometry (render byte-identical).
- **M14** — `FuzzParseFamily` (1.25M execs clean), `BenchmarkErrorPage`
  dropped (repo already had `BenchmarkErrorpageRenders` — duplicate
  avoided), coverage matrix + branch-combo renders (errorpage 71.9→72.4%,
  handler code 91%+), whole-repo coverage recomputed into FEATURES
  (root 70.2%, sub-modules 69.6–77.7%).
- **M15** — Website "Error Pages" guide (`guides/error-pages.md` + nav
  registration; link checker caught my bad anchor — the gate works);
  phrasing sweep clean (no "Wix-style" phrases shipped anywhere);
  error-feedback recipe verified current (uses `FromError`/`ErrorHandler`,
  no removed aliases); website goldens refreshed (derived enum count 62).
- **M16** — `.#visual-update` flake app; ci-repro VERDICT + `--quiet-diff`
  verified (parallel session shipped them); visualtest go.mod tidy no-op +
  pin-policy comment; `golden -update` now logs CREATED/CHANGED files
  (visible with `-v`) so same-edit count bumps are obvious.
- **M17.3–M17.5** — check-templ-sync verified repo-wide (website included,
  pre-shipped); golden count single-sourced to the filesystem walk + new
  README `.golden` assert — immediately caught a stale 251 in README.
- **M18 (partial)** — DOMAIN_LANGUAGE gained the Correlation Trace term;
  the enum-table guard named the new enums in FEATURES.
- **M19** — QF1003 tagged switches (collapsible_section, animated_icon ×2,
  website docs.templ; render-identical, zero golden churn); retired the
  `errBlankNonRejection` alias indirection; new errorpage tests linted to
  **0 issues** (noctx/nolintlint/unconvert/errname/dupword/gocognit fixed,
  `ErrorHandler` complexity fixed via `writeHTMLShell`/`applyRetrySuggestion`
  extraction).
- **M21** — All three items verified in-tree (parallel session did them):
  SidebarNav light/dark goldens + opt-out test, zero stale
  `formsValidationError` refs, recipe uses `forms.ValidationError`.
- **M22** — Fix card now `aria-describedby`-links the context table
  (grouped for AT; no dangling ref when context is absent) +
  `TestFixCardDescribesContext`; footer semantics audited: meta footer is
  in-card, not a page landmark — correct as-is.
- **M23** — Retry suggestion shipped: `ErrorHandler` auto-fills an empty
  way out with a same-path `Retry` link when `IsRetryable()` (JSON +
  Override suppress; 3 test cases). Validate soft-warning and oops.Time()
  documented as deliberate non-changes (permissive Validate is by-design;
  timestamp adaptation belongs bridge-side).
- **M24** — Lean `/errors/playground`: stateless GET form (family, status,
  title, message) → real ErrorPage at the real status code; sanitized +
  clamped server-side; demo section + form + content components; axe
  audited it in the full visual pass.
- **M25** — Harvest: TODO_LIST #246/#267 consumed (removed), #264 narrowed
  to rate-limit posture, #269 (gef PR ⫱) + #270 (release cut ⫱) added;
  plan annotated with a full execution record + deliberate deviations;
  ROADMAP gained the errorpage-remainder block.

## b) Partially done

- **M14 coverage target**: errorpage settled at 72.4% — the remaining gap
  is generated `_templ.go` wrapper statements, not hand-written code
  (ErrorHandler 91.2%, WriteErrorPage 70%, writeJSONError 90.9%). The 75%
  plan figure assumed render-only gaps; further gains need generator
  changes, not tests.
- **M18**: `.fail/` naming convention, setup-hooks fresh-clone check, and
  MaxMismatch/viewport audit NOT done — recorded in ROADMAP ("Errorpage
  plan remainder").
- **M17.1/17.2**: BuildFlow-side go-structure-linter rule config — external
  repo work (TODO #93/#231 family), not touched.

## c) Not started

- **M20 execution**: readiness fully verified (guards green, `[Unreleased]`
  warm, tags for 1.18.0 present) but the CUT itself is ⫱ owner (TODO #270).
- Nothing else — every plan item is done, verified-done by the parallel
  session, gated, or explicitly carried in TODO_LIST/ROADMAP.

## d) Totally fucked up

- Nothing destructive. Self-inflicted and fixed in-session: two
  `httptest.NewRequest` noctx findings, a wrong `unparam`-shaped refactor
  that needed a second pass, the promoted-`Nonce`-in-literal gotcha (twice
  — the AGENTS-documented trap), an edit-vs-daemon race on CHANGELOG and
  notfound404.templ (re-read + reapplied), a clobbered test header in
  variants_test.go (caught by build, rewritten cleanly), and the
  enum-table guard catching the unnamed new enums in a full-run only
  (now fixed). Python-heredoc used twice for Go files early on — flagged
  myself, switched back to edit tools per AGENTS.

## e) Improvements made to the repo's immune system

- Docs-count drift guard now also pins README's HTML-golden count
  (M17.4 single-sourcing) — caught a real stale number within one run.
- `golden -update` reports CREATED/CHANGED files — same-edit rule is now
  mechanically visible.
- New closed-set enums ship IsValid + tests in the same commit
  (ErrorDetailVariant, ErrorMaxWidth) and are named in FEATURES per the
  exhaustive-values guard.
- aria grouping + retry suggestion have string-level regression tests.
- tc-sources scaffolder kept in lockstep (sync check green after every
  .templ change).

## f) Next things (prioritized)

1. ⫱ Push `feat/bridge-classified-error-message` + file the gef PR (TODO #269).
2. ⫱ Go/no-go on cutting v1.19.0 (TODO #270) — readiness verified.
3. Re-check CI (lint lane) after the daemon lands this session's commits.
4. Vision-model pass over the new errorpage goldens (`scripts/vision-review-goldens.sh`)
   to eyeball-verify detail/alert/secondary/neutral captures.
5. `.fail/` artifact naming convention + CI cleanup step (ROADMAP).
6. setup-hooks fresh-clone assertion in CI/doctor (ROADMAP).
7. MaxMismatch/viewport audit for errorpage captures (ROADMAP).
8. BuildFlow go-structure-linter rule-level config upstream (TODO #231/#93).
9. Demo endpoint rate-limit posture for `/errors/playground` (TODO #264).
10. Consider `ErrorHandlerConfig.HeadContent` hook so standalone error
    routes can use the library HTMLShell with styles (ROADMAP candidate
    from the M03 redesign discovery).
11. Website error-pages guide: add the playground URL to the demo links.
12. Wire `CopyCode` into the demo full-model route (currently only the
    prop is documented + golden-pinned).
13. Keyboard focus-order test for the primary + secondary action pair.
14. `WayOutAction` adoption in the demo link-card grid (typed path showcase).
15. Extend `TestChipsJSONParity` to cover `context` map ↔ context table.
16. Benchmark the JSON path (`writeJSONError`) and record numbers.
17. Fuzz `sanitizeErrorMessage` (length/control-char handling).
18. Add ErrorAlert `Neutral`-equivalent if inline-alert feedback arrives
    from consumers (open question, no current demand).
19. Route-golden matrix: add `errors_playground` light capture (form-only
    state) if the route list grows again.
20. Docs: mention `MaxWidth` in the website error-pages guide.
21. Sweep the remaining `//nolint:exhaustruct_v5` in errorpage for
    validity (lint config may have moved).
22. `TestErrorPageBranchCombos`: add `CopyCode` × `MaxWidth` interaction.
23. Extract `errorPlaygroundForm` Tailwind classes into shared demo consts
    if a second form appears.
24. Consider promoting the copy listener to `utils/` if a third package
    needs it (goBackScript precedent documented in shared.templ).
25. Annotate AGENTS.md with the " promoted-field-in-literal" recurrence
    (hit twice this session) if it keeps biting.
26. Snapshot the gef bridge probe results into the gef PR body when filed.
27. Check whether `check-html-valid.sh` ignore list can shrink after the
    next vnu upgrade (14 classes currently).
28. Add `/errors/playground` to `docs/` demo endpoints inventory if one
    exists in the demo README.
29. Re-run `nix run .#coverage` after any generator change to refresh the
    FEATURES sentence honestly.
30. `TestFeaturesEnumValuesExhaustive`: confirm the new enums' VALUES are
    covered (guard is conservative — verify it actually parsed them).
31. Review `familyStyleDefault` ActionButtonGhost for gray-family contrast
    in both themes (visual pass said fine; belt-and-suspenders).
32. Consider `SecondaryWayOut` in `NotFound404` (currently ErrorPage-only).
33. Website: link the new guide from the errorpage FEATURE rows' docs column.
34. Verify the flake `visual-update` app inside a clean CI run (locally
    smoke-tested via flake check only).
35. Track upstream templ release for the v0.3.1036 import-style drift
    (bump go.mod + devshell together when it lands).

## g) Unanswerable questions (need owner)

1. **Release timing** — cut v1.19.0 now (all gates green) or batch with
   the gef bridge PR merge so the CHANGELOG can reference the upstream fix?
2. **gef PR authorship** — file from `larsartmann` directly, or via a
   feature fork with the jj-fork-pr-workflow flow?
3. **Parallel session status** — is the release-first session closed? Its
   f-list items are all verified done in-tree; I harvested nothing
   duplicated, but if it's still live we should stop touching the same
   demo/docs files.
