# Status Report — Plan Execution: templ-components ↔ CV Adoption (Session 2)

**Generated:** 2026-09-08 15:35 CEST
**Session scope:** Full-execution mode of the approved Pareto plan (`docs/planning/2026-09-08_07-52_cv-templ-components-adoption-pareto-plan.md`). Two repos: `~/projects/CV` (feature branch `feat/templ-components-adoption`) and `~/projects/templ-components` (master, untouched so far).
**Baseline before work:** CV build + full test suite green on master (`6d156315`).
**Format note:** user-specified `.md` (skill HTML default overridden again).

---

## a) FULLY DONE (verified green before/after each)

| Task | Commit | What landed |
|------|--------|-------------|
| D8 baseline | — | CV `go build ./...` + `go test ./...` green recorded before any change |
| D4 context reads | — | CV `AGENTS.md` read → **found the 2026-09-07 owner surface ruling** (ATS dashboard + chat = low-value, do-not-invest) that reshaped the plan |
| D7 pre-flight | — | No overlap between C1–C6 and TC `TODO_LIST.md`/`FEATURES.md`; TODO #156 (consumer survey) independently validates C1 (another consumer cites missing head-content support) |
| A1 version bump | `48e4f43` | templ-components v1.13.2 → v1.14.0 (all 6 modules), safelist regenerated, output.css rebuilt (Tailwind v4.3.3), build + full tests green. Picks up two CSP-nonce regression fixes + additive wire/forms features |
| A2 error-boundary deletion | `25c40f04` | ~90-line bespoke ATS error overlay deleted (including its CSP-dead `onclick` dismiss button); API-key header stamp + modal listeners + init preserved; single error path = `htmx.GlobalErrorHandling`. Golden regenerated (-86), diff eyeballed |
| A3 skeletons + motion-reduce | `9208e378` | Pipeline `listSkeleton` → `feedback.SkeletonGroup`; 36 motion-reduce violations fixed across 12 core-surface files via exact-token Perl pass; regenerated ATS golden diff carries ONLY motion-reduce tokens (proof of surgical scope) |
| A4 EmptyState adoption | `09cbfc19` | `emptyPanel` → thin delegate to `display.EmptyState` (3 SSE fragments); pipeline `filter-empty` → `EmptyState` with real-button action (JS contracts `#filter-empty`/`#filter-clear` preserved); `interviews-empty` → `EmptyState` (id kept) |
| A6+A7+B4 | `4c725e2d` | Pipeline `statCard`/`statCardLink` → `display.StatCard` (+`ValueID`, aria-live, Href variant — admin-hub pattern); `deadPortalsBody` → `display.Table` (Headers + Body slot, CellPaddingCompact); the stale `RelativeTime` comment replaced with the true mechanism (SSE innerHTML never executes scripts — AutoRefresh could never attach regardless of nonce) |
| B1 A.Team form | `26a25845` | Raw `<form>` → `forms.Form` (CSRF stays JS-populated by design — documented); progress bar → `feedback.ProgressBar` (JS retargeted to `[role="progressbar"]` + aria-valuenow lockstep); submit → `htmx.LoadingButton`; **CSP bug fix en route: dead `onclick=` Reset handler → delegated `data-tc-reset-form` listener** |
| B2 coaching page | `99c933e2` | Inputs/textareas → `forms.Input`/`Textarea`; Run buttons → `display.Button` (`.coaching-submit` hook kept); 503 chip → `display.Badge`; amber notices → `feedback.Alert`. JS contracts untouched |
| D3 PageHeader collision | in tree (`19b19bba` daemon) | `common.PageHeader` → `common.DashboardHero` (definition + 4 call sites) — CV no longer shadows the library's `display.PageHeader` name |
| D2 share.js | in tree (`19b19bba` daemon) | `landing.templ` share.js → nonce'd `tclayout.Script(nonce, src, nil)` |

**Discipline note:** every task ran templ generate → build → package tests; every intentional golden change was diff-eyeballed; the newly-used library classes were spot-checked in the compiled CSS (one false alarm: CSS escapes dots, `h-1\.5`).

## b) PARTIALLY DONE

1. **D2 (CV AGENTS.md adoption table)** — the table rows for the new adoptions (StatCard /pipeline, Table dead-portals, EmptyState pipeline, SkeletonGroup, LoadingButton, forms.Form, coaching components) were being edited when this report was requested; NOT yet updated. The datastar indirect-dep verdict (keep: pulled by root module; `tidy` manages it) is decided but undocumented.
2. **D3 commit hygiene** — rename + share.js changes ride daemon snapshot `19b19bba` with a hallucinated message instead of a reviewed commit (daemon won the race; cleanup scheduled pre-PR).
3. **A3 scope** — intentionally narrowed (ATS/chat surfaces skipped per ruling); the "repo-wide" claim in the plan holds only for core surfaces.
4. **Branch hygiene** — 4 daemon noise commits (`e05bb0dd`, `bbb9479c`, `a1369601`, `36c3acd6`) sit interleaved between reviewed commits; pre-PR rebase surgery (documented detached-worktree procedure) not yet run.
5. **Test-suite state** — final full CV run has one known pre-existing flake (`TestSSEWire_DomainEventTriggersBroadcast`) documented in `9208e378`; 3× green in isolation, 3× green package reruns, green on pristine baseline worktree.

## c) NOT STARTED

- **C1 `layout.SEOMeta`** upstream (design decision locked: named struct field on `layout.PageProps`)
- **C2** CollapsibleSection nonce'd persistence script
- **C3** `docs/recipes/sse-fragments.md`
- **C4** generic icon renderer
- **C5** `docs/recipes/htmx-modal.md`
- **C6** `docs/recipes/print-pdf.md`
- **D1** RTL/logical-property sweep (motion-reduce ≠ RTL — still open)
- **D5** TC AGENTS.md "v2.0 vs v1.14.0" drift fix
- **D6** full reads of the 5 skimmed regions (~1,800 lines)
- **D7 harvest** into TODO_LIST/ROADMAP (both repos)
- CV branch push + PR (push not requested this session)

## d) TOTALLY FUCKED UP

1. **Empty A4 commit (caught + fixed):** first `commit-tree HEAD^{tree}` captured the HEAD tree, not the index — produced a content-free commit. Rebuilt with `git write-tree` (index tree) and ref-update. Root cause: unfamiliarity with plumbing semantics; verified by `git show --stat` afterwards.
2. **Committed on a red run:** the A6/A7/B4 commit (`4c725e2d`) landed while the suite printed `FAIL` — my `| tail -1` pipeline masked the failure until after the commit. Post-hoc 3× reruns green (the documented pre-existing SSE flake, proven on a pristine baseline worktree), so the CONTENT is fine — but committing before a green gate is a discipline violation I then corrected by adopting `set -o pipefail`.
3. **Plan shipped tasks that died on contact:** the plan (approved with A5/B3/B5/ATS-skeleton work) targeted the ATS dashboard — the very surface the owner had ruled do-not-invest ONE DAY EARLIER in CV's AGENTS.md. D4 was correctly ordered first, but a perfect plan would have read the consumer's AGENTS.md during PLANNING, not execution. ~4 tasks of planned work were void.
4. **B4 was half-wrong in the plan:** my analysis session claimed the RelativeTime nonce comment was "stale" and AutoRefresh could simply be enabled. Execution revealed the deeper truth: SSE innerHTML delivery never executes scripts, so `AutoRefresh:false` is correct BY MECHANISM. The plan item got downgraded to a comment fix. Lesson: the nonce-support fact was right; the adoptability conclusion was not.
5. **Daemon wars cost ~3 commit cycles** (HEAD lock errors, raw snapshots) before switching to the `commit-tree` + `write-tree` + `update-ref` plumbing pattern. The workaround now works reliably; the noise commits it leaves are tracked for pre-PR cleanup.
6. **No lie detected:** all commits verified content-wise; the one red-gate commit is explicitly disclosed above.

## e) WHAT WE SHOULD IMPROVE

1. **Read the consumer's AGENTS.md during planning, not execution** — the single biggest plan-quality miss (see d3).
2. **Never gate commits behind filtered pipelines** — `set -o pipefail` from the first verification command, not after the first incident.
3. **Prefer `write-tree`/`commit-tree` from the start in daemon-infested repos** — the race cost three cycles; the plumbing pattern should have been the default after the first lock error.
4. **Verify adoption feasibility at the delivery layer** (does the JS even execute there?), not just the API layer (does the prop exist?).
5. **Re-read after every scripted bulk edit** — three edit-tool bounces on stale-mod-time files were pure waste.

## f) NEXT — up to 50 things (ordered: finish CV branch → upstream C-tasks → hygiene)

1. Finish D2: update CV AGENTS.md adoption table (StatCard /pipeline, Table dead-portals, EmptyState pipeline, SkeletonGroup, LoadingButton, forms.Form, coaching set)
2. Land D2/D3 content as one reviewed commit (daemon snapshot `19b19bba` → amend/squash via plumbing)
3. D1: CV RTL sweep — inventory `ml-`/`pl-`/`left-`/`text-left` on core surfaces, convert to logical properties
4. Pre-PR branch surgery: drop the 4 daemon noise commits (detached temp worktree, `git rebase -i` drop, verify `git diff origin/master --stat` shows only semantic content)
5. Run CV lint (golangci-lint) on all changed packages
6. Run CV full test suite final gate (with `pipefail`, flake excluded/document)
7. Push CV branch + open PR (split: one adoption PR vs per-task PRs — see Q1)
8. Verify the Interviews StatCard link still receives the operator API-key `?key=` stamp (JS `querySelector("a[href='/calendar/interviews.ics']")` must match StatCard's `<a href>` — expected yes, unverified)
9. Run CV e2e (playwright) for A.Team form + pipeline + landing after the markup swaps
10. Re-run ATS visual golden regen if e2e/visual harness flags drift (expected none beyond A2/A3)
11. Manual/axe pass: pipeline page after EmptyState swaps (role=status nesting, heading levels)
12. Landing page smoke: share.js still loads via nonce'd Script; theme toggle intact
13. Confirm print/PDF CV is byte-identical (all changes were screen-side; assert no `@media print` surface touched)
14. CV AGENTS.md: document the B1 CSRF-contract decision (meta-tag JS fill vs FormProps.CSRFToken)
15. CV AGENTS.md: document the SSE-innerHTML-no-scripts mechanism note (source of truth for future adopters)
16. File CV issue: `TestSSEWire_DomainEventTriggersBroadcast` flake (SSE broadcast timing under parallel load)
17. Evaluate remaining raw buttons on A.Team (Validate) and ATS (Refresh/NewAnalysis — ruling-blocked) for `display.Button` swaps when surfaces reopen
18. Evaluate `stat-updated` SSE scalar → server-driven RelativeTime fragment (bigger refactor; needs SSE contract change)
19. Add pipeline equivalence goldens (currently only ATS has one — the swaps had no golden safety net)
20. C1: implement `SEOMeta` struct (NoIndex/Canonical/Alternates/JSONLD) on TC `layout.PageProps` + `Alternate` type
21. C1: render in `layout.Base` head + golden + unit + a11y tests
22. C1: CHANGELOG `[Unreleased]` + FEATURES.md + skill catalogue update
23. C1: `nix run .#verify` green in TC
24. C1 scope question: extend head-content support to `layout.Minimal` (TODO #156 demand from `nsfw-classifier`) — same PR or follow-up (see Q3)
25. C2: CollapsibleSection optional nonce'd persistence singleton (ThemeScript pattern) + tests + golden
26. C3: write `docs/recipes/sse-fragments.md` from CV's pipeline pattern (named events → fragment swap, JSON scalars, reconnect banner, the innerHTML-no-scripts gotcha)
27. C3: cross-link from htmx/datastar docs + transport-wiring doc
28. C4: `icons.Render(viewBox, paths, class, fill)` generic renderer + tests + consumer-set recipe (deletes CV's 3 hand-rolled SVG scaffolds later)
29. C5: `docs/recipes/htmx-modal.md` (dialog shell + HTMX swap + open/close helpers)
30. C6: `docs/recipes/print-pdf.md` harvested from CV's print stack
31. D5: fix TC AGENTS.md "v2.0" vs `Version=1.14.0` drift (annotate, docs-health ANNOTATE discipline)
32. D6.1: read `admin_page.templ` 760–1669 fully
33. D6.2: read `approvals_fragment.templ` 200–297 + `applications_fragment.templ` 200–313 fully
34. D6.3: read `dashboard_fragments.templ` 200–285 + `psychological_impact.templ` 80–247 fully
35. D7: HARVEST this report's section f into CV + TC `TODO_LIST.md`/`ROADMAP.md`
36. TC: adopt `htmx.PolledRegion` internally in the demo where `HtmxCard`-style patterns exist (dogfood the recommendation given to CV)
37. TC: consider documenting the A4 pattern (EmptyState + ActionAttrs JS hooks) as a recipe example
38. TC: RelativeTime docs — state the SSE/innerHTML limitation explicitly (library-side truth surfaced by CV)
39. CV: decide ateam Skills free-text vs `forms.TagsInput` (evaluation deferred from B1)
40. CV: pipeline `Refresh` dead-portals button → `display.Button` (small follow-up to A7)
41. CV: revisit ATS adoption items (A5 HtmxCard→PolledRegion, B3 buttons, B5 Modal) ONLY if the surface ruling is lifted — re-entry trigger documented in CV AGENTS.md
42. CV: measure output.css size delta post-adoption (safelist superset risk is bounded; confirm)
43. TC: record the "daemon vs feature-branch commits" workaround (commit-tree/write-tree/update-ref) in TC AGENTS.md daemon section
44. TC: BuildFlow upstream issues (#93 family) — unchanged, still blocked
45. CV: after PR merge, watch CI for the flake (#16) and safelist/CSS parity checks
46. TC: re-run `nix run .#css` byte-stability after C-task class additions
47. CV: keep `data/last-eval-pass.json` runtime churn out of PRs (daemon attractor file — consider .gitignore proposal)
48. Proposal: CV PR description carries the ruling-compliance note (which planned items were dropped and why)
49. TC: consumer case-study doc ("What CV taught us") distilling the 6 gaps
50. Final status report v3 after C-tasks + CV PR

## g) QUESTIONS (asked via native tool, blocking)

See the question form — (1) CV PR strategy + daemon handling, (2) upstream push policy for TC code changes, (3) SSE flake ownership.

---

*Point-in-time snapshot. Section (f) is HARVEST input for `TODO_LIST.md`/`ROADMAP.md`.*
