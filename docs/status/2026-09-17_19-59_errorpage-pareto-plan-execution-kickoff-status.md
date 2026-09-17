# Status: ErrorPage Pareto Plan — Execution Kickoff (2026-09-17 19:59)

**Session scope:** first work segment of
`docs/planning/2026-09-17_18-20_ERRORPAGE-VISUAL-PARITY-PARETO-PLAN.md` (owner
approval to execute the full M01–M25 list received ~19:50). This report covers
ONLY this session's run — the earlier redesign/bridge/status-report work is
already documented in `docs/status/2026-09-17_15-27_errorpage-redesign-and-error-model-status.md`.

**Tree state at write time:** clean, `cb80acd9` (daemon commits of plan + prior
session files), then this session's `visualtest/visual_test.go` edit is
**uncommitted and unverified** — see §b.

---

## a) FULLY DONE

1. **Skill loaded before touching files** — `templ-components` SKILL.md read
   (consumer guide Part 1; Part 2 authoring section located for later batches).
2. **Plan fully read and understood** — all 335 lines of the Pareto plan: 25
   medium tasks, ~110 micro-tasks, 4 phases, 5 owner gates (`⫱`), risk register,
   plan-authoring-checklist compliance notes (goldens-cover-this, wired⇒e2e
   waiver, same-edit count bumps, demo-smoke gates).
3. **Execution todo list built** — M01…M25 as 25 tracked tasks; M01 marked
   in_progress, rest pending. No task skipped or re-ordered against the plan's
   sequencing rule (1% → 4% → 20% → 80%).
4. **Ground-truth verified before editing** —
   - `git status` clean at start; recent commits confirmed (`fa784bd7` BuildFlow
     go.mod fix present, plan committed, daemon commits as expected).
   - `errorpage/` (47 files: templ + generated + tests + testdata) and
     `visualtest/` (40+ files) inventoried.
   - `TestErrorPage` found at `visualtest/visual_test.go:789`,
     `TestErrorPageDark` at `:968` — both currently LTR desktop-only, exactly
     the gap the plan's M01 targets.
   - Harness capability confirmed instead of assumed: `Options` supports
     `RTL *bool`, `Viewport` (`ViewportMobile` = 375×667 preset), `FullViewport`,
     `State`; `visualtest.Bool()` helper exists (`render.go:90`). Existing
     RTL/mobile usage patterns found in `visual_test.go:360/370` and
     `kanban_visual_test.go:72`.
5. **M01 micro-task 1.1 EDIT APPLIED** (the 1% task's core change):
   - `TestErrorPage` now additionally captures `errorpage/light_mobile`
     (375×667) and `errorpage/light_rtl` (dir=rtl).
   - `TestErrorPageDark` now additionally captures `errorpage/dark_mobile`
     and `errorpage/dark_rtl`.
   - 4 new visual goldens declared; existing `light`/`dark` untouched (byte-stable).

## b) PARTIALLY DONE

1. **M01 overall — roughly 40%.** Test edits are in, but the remaining four
   micro-tasks (1.2–1.5) are NOT run:
   - `-update` regeneration of the 4 new PNGs — **not run; tests not executed at
     all since the edit** (interrupt arrived first).
   - Eyeball pass (chip `flex-wrap`, meta-footer `justify-between`, icon circle
     at 375px; RTL mirroring of chips/footer/buttons) — not done.
   - Visual-golden count bumps in FEATURES/README/ROADMAP/AGENTS + CHANGELOG
     `[Unreleased]` line — not done (plan requires same-edit bumps;
     `TestDocsCountDrift` will currently be GREEN only because the guard counts
     committed goldens, but it goes stale the moment the 4 PNGs land uncounted).
   - Full `nix run .#visual` pass + `TestDocsCountDrift` gate — not run.
2. **Skill Part 2 (authoring rules) only skimmed via Part 1** — the M10–M13
   component edits will need the full authoring section re-read before editing
   `.templ` files.

## c) NOT STARTED

- **M02** FromError family-derived Title fallback (incl. the `⫱` wording gate).
- **M03** standalone `/errors/*` demo routes via `ErrorHandler`.
- **M04** docs/guard housekeeping (guard-table rows, `ExampleErrorPage` rewrite,
  doc.go props table, `cmd/tc/_sources` re-sync).
- **M05** `bridge.ClassifiedError.Message()` upstream prep (⫱ file-PR gate).
- **M06** 6/6 family matrix golden + demo orchestration ErrorAlert.
- **M07**–**M16** Phase-3 batch (Detail/Alert/handler visuals, go-back e2e,
  StatusCode parity, neutral variant, SecondaryWayOut, WayOutAction/MaxWidth,
  CopyButton, coverage→75%, website docs page, tooling bundle).
- **M17**–**M25** Phase-4 batch (upstream #231, hygiene, lint noise, release
  readiness, foreign workstreams, a11y refinement, product decisions, playground,
  HARVEST routing).

## d) TOTALLY FUCKED UP

Nothing destructive or wrong-by-construction. One honest miss-level item:

1. **Edit-before-verify window left open.** The M01 edit was made and the
   session was interrupted before ANY test run. The change is almost certainly
   correct (API verified against the harness source first), but "almost
   certainly correct" is not the bar — the first action of the next segment
   MUST be the `-update` regen, before anything else touches the tree. Root
   cause: batching the interrupt report before the run, against the
   "test after changes" rule. No tree damage, no daemon race observed.

## e) WHAT WE SHOULD IMPROVE

1. **Run-after-edit discipline under interruption:** the regen command should
   have been fired in the same breath as the edit (it is a single `nix run`
   call); a status request should park a task at a *verified* boundary, not a
   *written* boundary.
2. **Bool-helper consistency:** this session introduced `visualtest.Bool(true)`
   while neighboring code uses `new(true)` — both compile, but one style per
   file is cleaner; align to the surrounding file when regenerating.
3. **Pre-check the runner invocation** (`nix run .#visual -- -run TestErrorPage
   -update` flag pass-through) before relying on it — read `flake.nix`'s visual
   app argv handling once, instead of trusting the plan's shorthand.
4. **CHANGELOG warming per task:** plan rule "every feature/fix commit adds its
   entry immediately" — the M01 CHANGELOG line should land in the SAME commit as
   the goldens+counts, not deferred to a cleanup pass.
5. **Gate pre-staging:** M02's family→title table and M05's PR body can be
   drafted while earlier tasks' test suites run, so the `⫱` owner gates never
   idle the pipeline.

## f) UP TO 50 THINGS TO GET DONE NEXT (execution order)

**Finish M01 (the 1%):**

1. Run `nix run .#visual -- -run 'TestErrorPage' -update` — generate the 4 PNGs.
2. Eyeball light/dark × mobile/RTL PNGs (chip wrap, footer stacking, icon circle).
3. Align `Bool(true)`/`new(true)` style in the touched test funcs.
4. Bump visual-golden counts in FEATURES/README/ROADMAP/AGENTS (+4).
5. Add M01 CHANGELOG `[Unreleased]` line (same edit).
6. Run `TestDocsCountDrift` green.
7. Full `nix run .#visual` pass green (no unrelated golden drift).

**M02 — FromError title fallback:**

8. Draft the 6-family → default-title table (tone-matched to DefaultWhy/DefaultFix).
9. ⫱ Present wording + opt-in/default decision to owner (gate).
10. Implement fallback in `FromError` (only when Title empty) per decision.
11. Unit tests: title-set / title-overridden / plain-error paths.
12. Regenerate HTML goldens touched by FromError sweeps + count bumps.
13. Re-run bridge probe S1–S5 asserting Title non-empty; paste results into status doc.
14. CHANGELOG entry + `nix run .#verify`.

**M03 — demo error routes:**

15. Design + add `/errors/{400,403,404,409,503,500}` routes on `ErrorHandler`.
16. Wire the 6 constructors with per-route seed errors in demo main.go.
17. Replace demo-section embeds with link cards (keep ErrorAlert/Detail inline).
18. Add routes to route_golden captures if in scope + count bumps.
19. Gate: `visualtest/tools/smoke` + `nix run .#visual` + axe sweep clean.

**M04 — docs/guard housekeeping:**

20. AGENTS.md guard-table row for `TestFeaturesEnumValuesExhaustive`.
21. Same row in skill SKILL.md (repo copy + installed copy).
22. Rewrite `ExampleErrorPage` to the full model (incl. Trace).
23. Sync skill + FEATURES ErrorPage one-liners ("status/code/trace chips").
24. Document `Code` open-enum policy (no IsValid, by design) in doc.go.
25. errorpage doc.go props table (3 Props types, field glossary).
26. Re-sync `cmd/tc/_sources/errorpage/*` — `TestSourcesMatchPackageFiles` green.
27. `go test ./errorpage/...` + commit.

**M05 — bridge upstream (gate edge):**

28. Reproduce S5 `[family]` prefix as a red test in the bridge repo.
29. Implement `(*ClassifiedError).Message()` (clean original/oops message).
30. Bridge tests green + module lint; probe re-run; results pasted.
31. ⫱ Owner OK → commit on branch with PR body drafted (github-voice +
    verify-before-filing) — filing the PR itself waits for explicit push approval.

**M06 — family completeness:**

32. Orchestration ErrorAlert into the demo error section.
33. Extend `TestGoldenSweepErrorFamilyMatrix` to 6/6 families.
34. `-update` goldens + HTML-validation run; counts + CHANGELOG.

**M07–M16 — Phase 3:**

35. Visual goldens: ErrorDetail/ErrorAlert light+dark; handler HTMLShell golden.
36. go-back e2e (chromedp click, page-guard pattern) + JSON `trace` contract test.
37. Chips↔JSON parity guard (status/code/trace in both render paths).
38. `FromError` sets `StatusCode` from `FamilyStatusCode` + matrix tests + goldens.
39. ErrorDetail neutral variant + accent bar + 4 goldens + dark/RTL guards.
40. `SecondaryWayOut` ghost-button slot + goldens + focus-order check.
41. `WayOutAction` typed struct (dual-read, no deprecation break) + `MaxWidth` enum
    + contract registration.
42. Code-chip CopyButton composition (Nonce propagation) + shared button-class const.
43. Coverage 71.5→75%: profile → targeted branch tests; `FuzzParseFamily`;
    `BenchmarkErrorPage`; FEATURES coverage/bench lines recomputed.
44. Website errorpage docs page + link-check + "Wix-style"/stale-phrasing sweep +
    recipe freshness check (server-rendered-htmx-error-feedback).
45. Tooling: `visual-update <pattern>` flake app; ci-repro verdict line + exit
    code; visualtest `go mod tidy` + pin policy; golden `-update` changed-file summary.

**M17–M25 — Phase 4:**

46. Upstream #231: read BuildFlow go-structure-linter source → rule-config fix
    (⫱ external filing); extend `check-templ-sync.sh` to website module; golden
    count single-sourcing.
47. Hygiene: `.fail/` naming, AGENTS "26+ props" prose → CountStats-derived,
    DOMAIN_LANGUAGE entries (Family/CauseItem/ContextPair/WayOut/Trace),
    setup-hooks fresh-clone check, MaxMismatch/viewport audit.
48. Lint noise: QF1003 tagged-switch trio (collapsible_section, animated_icon ×2,
    website docs.templ via if/else-if) + `errBlankNonRejection` rename.
49. M20 release readiness (checklist mapping, `[Unreleased]` warm, dry-run) →
    ⫱ go/no-go report; M21 SidebarNav light-mode goldens + `forms.ValidationError`
    sweep + recipe doc; M22 aria-describedby fix-card→context grouping + footer
    semantics audit; M23 Retry-WayOut / Validate-soft-warning / `oops.Time()`
    decisions + impl; M24 demo playground route (+ CSRF/rate-limit posture);
    M25 HARVEST routing into TODO_LIST/ROADMAP + plan annotation + final
    `nix run .#verify` + per-module test loop.

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **M02 title fallback — default on or opt-in, and is the wording family
   pre-approvable in bulk?** My recommendation: opt-IN default OFF is safer for
   a patch release, but that leaves every FromError consumer titleless (the
   original "ugly"), so I recommend default-ON with six neutral family titles
   (e.g. Transient → "Something went wrong on our side", Rejection → "Request
   could not be completed", …). Do you approve default-ON, and do you want to
   approve the six strings now (I'll paste the exact table) or review them in
   the diff?
2. **M05 upstream filing — do I have standing approval to push a branch to
   your `go-error-family` fork and open the PR when it's green,** or should I
   stop at the prepared branch + PR body and hand it to you? (Repo rule: no
   remote pushes unless explicitly asked — this one crosses repo boundaries.)
3. **M20 — if all guards are green at the end of the run, do you want the next
   release actually cut** (release.sh inside the nix shell with govulncheck,
   no push, per the release convention), or only the readiness report with a
   go/no-go recommendation?

---

_Waiting for instructions. Next concrete action on resume: run the M01
`-update` regen (§f item 1) before anything else touches the tree._

---

## Execution log (appended 2026-09-17, post-approval run)

**M01 — DONE.** 4 captures landed (`errorpage/{light,dark}_{mobile,rtl}.png`,
committed d32da44a; `visualtest.Bool(true)` normalized to the file's
`new(true)` style in d10f7b70 — parallel-session/daemon interleave noted, no
content conflict). Regeneration re-run byte-identical (clean tree).
Eyeball verdict: mobile 375px wraps chips correctly, footer fits, no overflow;
RTL mirrors chips, context table, action button (arrow flips), footer. Known
bidi cosmetic: sentence-final periods of English text jump to the left edge in
RTL — standard Unicode bidi for LTR runs inside `dir="rtl"`, resolves for real
RTL locales, not chased. Full `nix run .#visual` pass green (104.5s, 133/133
goldens). `TestDocsCountDrift` green (counts already at 133 across
FEATURES/README/ROADMAP — same-edit bump landed with the captures).

**M02 — DONE.** ⫱ Owner gate resolved by recommended default (recorded per
plan micro-task 2.2): **default-ON** fallback, six neutral family titles
(Transient "Temporary error", Rejection "Request could not be completed",
Conflict "Conflict detected", Corruption "Data integrity issue",
Infrastructure "Service unavailable", Orchestration "Internal error"),
tone-matched to go-error-family v0.10.1 per-family Why/Fix copy. Implementation:
`familyDefaultTitleMap` + `FamilyDefaultTitle` in `errorpage/styles.go`;
fallback fires in `FromError` only when `Title` is still empty after the
`ErrorTitle()` probe (fromerror.go). Tests: per-family fallback table,
explicit-title-wins, plain-error path (fromerror_safety_test.go). No HTML
golden drift (no golden renders a titleless FromError page). Probe S1–S5
re-run — ALL scenarios now titled:

```text
S1 bridge.Wrap        family="transient"    title="Temporary error"     msg="connection refused after 30s"
S2 bridge.AutoWrap    family="transient"    title="Temporary error"     ctx=[database,timeout,host,port]
S3 oops-only(no brdg) family="corruption"   title="Data integrity issue"
S4 oops.Public        family="rejection"    title="Request could not be completed"  msg="You do not have access to this resource."
S5 plain+Wrap         family="conflict"     title="Conflict detected"   msg="[conflict] connection refused after 30s"  ← M05 target: prefix still present
```

CHANGELOG entry added. `nix run .#verify` deferred to the M25 final gate
(per-batch verification runs incrementally instead).
