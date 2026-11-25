# Status — Consumer-Review Synthesis: Docs, Marketing, AI Skill

**Report time:** 2026-11-25 20:15 CET · **Repo:** templ-components · **HEAD at write time:** `e1fed513`

**Scope — what this report VERIFIED:** all 22 `~/projects/*/docs/review/libs/templ-components.md`
files read in full this session (2026-11-25; the corpus itself is dated 2026-11-04 → 11-17);
the HTML deliverable validated (W3C html5validator: zero content errors, 3 CSS errors proven
template-inherited by diffing the pristine kit template, which produces the same 3; tag balance
100/100 divs, 8/8 sections; template splice byte-faithful: head diff = title line only, tail
identical); the stale-date defect caught and fixed (rename + 5 surgical edits; `grep` confirms
zero full-form `2026-11-17` refs remain); owner rulings applied (cmdguard dropped; GOAL ruling
annotated into the report). **What it did NOT verify:** rendered browser appearance of the
report (no screenshot); corpus quotes re-derived against library source (deliberately — the
report's scope was the corpus's own claims, copied verbatim); the exact count of reviews citing
SKILL.md (published as "10+", not grep-counted); cmdguard corpus files 10–18 (never read; 1–9
read then dropped on owner instruction). Both deliverables are UNTRACKED at write time — no
commit (the harness forbids commits without an explicit request; the auto-commit daemon will
pick them up).

**Owner ruling received this session (2026-11-25, ~20:15 CET):** the "Tailwind-free styling
path" candidate is **RULED OUT**. The goal is the opposite: make downstream consumers use and
leverage TailwindCSS **to the MAX**. Section (f) is reordered accordingly (Lane 0).

## Straight answers first

- **What did I forget?** To run `date`. The HTML report was born dated 2026-11-17 — eight days
  stale, with an appendix claiming the corpus read happened on a day that hadn't occurred yet.
  I inferred "today" from the corpus's newest review date instead of running a one-command
  check. Caught when the status-report skill forced `date`; fixed same session.
- **What could I have done better?** Surface the subject/evidence contradiction immediately.
  The prompt said "cmdguard's docs" but pointed at templ-components reviews (zero cmdguard
  mentions in the corpus — I even grepped). Instead of flagging that as the blocking ambiguity
  it was, I absorbed a second corpus (9 cmdguard files, ~100KB of context) that the owner then
  discarded with "Ignore cmdguard". Also: exact counts over hedges, exit codes not masked by
  `| tail`, single-writer scope claimed at session start.
- **What could I still improve?** Same-session loop closure (harvest into TODO_LIST.md, not
  entomb in a timestamped file), render-level verification for HTML deliverables, and running
  the report skill's own first step (`date`) before ANY dated artifact, not just status reports.

## a) FULLY DONE

1. **Corpus read:** all 22 templ-components consumer reviews, end to end (~2,446 lines / 376KB
   across 18 consumer repos; corpus dated 2026-11-04 → 11-17).
2. **Synthesis:** findings D1–D8 (docs), M1–M3 (marketing), S1–S5 (AI skill), every claim
   cited to the review that raised it; strengths section (what the corpus says is already
   excellent) included so the brief isn't a complaint list.
3. **HTML brief written:** `docs/reviews/2026-11-25_docs-marketing-skill-audit.html` (1,671
   lines; editorial Bauhaus kit variant per the kit's adoption-feedback assignment, deviation
   from the dark sibling noted in the report per the kit's series rule; CSP-safe: zero inline
   handlers, zero style attributes; scrollspy kept as progressive enhancement).
4. **Verification pass:** placeholder-free; tags balanced; template splice byte-faithful
   (head = title line only; tail identical); W3C-valid modulo the kit's 3 documented
   `color-mix()` staleness errors (proven template-inherited, not introduced).
5. **Owner corrections applied cleanly:** (a) cmdguard dropped, report refocused on
   templ-components before any cmdguard content entered it; (b) GOAL ruling annotated into the
   report's routed-elsewhere callout (Tailwind-free path struck through, RULED OUT, direction
   inverted); (c) date defect fixed (file renamed 11-17 → 11-25; top bar, eyebrow, appendix,
   footer corrected; grep-verified).
6. **Scope discipline:** exactly one new file under `docs/reviews/`; the sibling session's
   `cmd/tc/*` working-tree changes untouched; nothing staged or committed.

## b) PARTIALLY DONE

1. **cmdguard corpus read:** 9 of 18 files before the owner correction; all cmdguard content
   dropped from the deliverable. The reads themselves are now unused context (see d2).
2. **Single-writer protocol:** scope claimed at write time ("one new file under
   docs/reviews/"), not at session start; `TODO_LIST.md` not re-read before the edit
   (defensible — the edit was a new untracked file with zero overlap — but the rule says
   re-read before every edit).
3. **Report actions harvested:** the P1–P3 lanes are ranked in the report but NOT yet written
   into `TODO_LIST.md` (docs-health HARVEST is the standing next step; see g1).
4. **Claim precision:** "10+ reviews cite SKILL.md by line number" is hedged, not
   grep-verified (f42).
5. **Validator exit-code discipline:** the html5validator rc was masked by a `| tail` pipeline
   ("exit: 0" was tail's status, not the validator's); mitigated by proving the same 3 errors
   exist in the pristine template, but the check itself was sloppy.

## c) NOT STARTED

1. **The GOAL program** (owner directive 2026-11-25): making downstream leverage Tailwind v4
   to the max — Lane 0 in section (f), nothing begun.
2. **Every P1/P2/P3 doc fix** in the report (countgen over README counts, module map, pairing
   matrix, first-screen cost block, replacement map, upstream-requests ledger, CHANGELOG
   impact lines, guide harvest, presets/demo surfacing, guard recipe).
3. **Skill additions:** adopt/reject quick profile, "what's in my pin?" protocol, gotcha
   harvest into the consumer guide.
4. **Product filings:** StatCard value slot, DataTable tfoot, day-unit `CompactDuration`,
   popover `position-area` fix, Badge prefix slot, errorpage stylesheet slot.
5. **Harvest** of this session's outputs into `TODO_LIST.md` / `ROADMAP.md`.
6. **Release-cadence decision** (v1.28.0 cut vs rc tags; five consumers gate on "the next
   tag", one pinned a pseudo-version).

## d) TOTALLY FUCKED UP

1. **The deliverable was born with a false date.** File name, top bar, eyebrow, appendix, and
   footer all said 2026-11-17; today is 2026-11-25. The appendix literally claimed the corpus
   was "read in full on 2026-11-17" — a day that hadn't happened when the read occurred
   (2026-11-25). Root cause: I never ran `date` during the report work and silently equated
   the corpus's newest review date (profileagent, 11-17) with today. Caught only because the
   status-report skill mandates `date` first. Fixed: renamed to
   `2026-11-25_docs-marketing-skill-audit.html`, 5 surgical internal edits, grep-verified.
   Lesson: one command would have prevented it; the rule exists because this exact failure
   mode recurs.
2. **The subject/evidence contradiction was absorbed instead of surfaced.** "cmdguard's docs"
   vs a templ-components review corpus — I resolved it autonomously by reading a second corpus
   (9 cmdguard reviews) instead of stopping at the cheapest disambiguation (state the conflict,
   ask one question, or present a scoped both-lane plan). Owner corrected mid-stream; damage
   bounded (~100KB context, 9 file reads; no cmdguard content reached the deliverable), but
   the wasted work was avoidable.
3. **Minor, same class:** the masked validator exit code (b5), and hedges where a one-line
   grep would have produced an exact number (b4).

## e) WHAT WE SHOULD IMPROVE

1. Run `date` + `git rev-parse HEAD` at the START of every dated deliverable; pin facts,
   never infer them from data recency.
2. When named subject ≠ named evidence, surface the contradiction in the first reply; choose
   the cheapest disambiguation before absorbing more evidence.
3. Exact counts over hedges — `grep -c` before publishing any "N+" claim.
4. Claim the single-writer scope in the session's first message, not at write time.
5. Close the loop same-session: report section (f) exists to be harvested into
   `TODO_LIST.md` (docs-health HARVEST), not entombed.
6. Never let a downstream pipe eat a validator's exit code; capture rc explicitly.
7. Add a render-level check (headless screenshot) for HTML deliverables when practical;
   syntax-valid ≠ visually correct.
8. Re-check `git status -sb` before finishing — the auto-commit daemon races session output,
   and its messages are adversarial-by-accident (known repo gotcha).

## f) Up to 50 things to get done next

**Lane 0 — GOAL: Tailwind-to-the-max (owner directive 2026-11-25; "Tailwind-free path" ruled out)**

1. Friction audit of the consumer Tailwind pipeline end to end (`go get` → styled page);
   write the friction list as the program backlog.
2. `tc css init|build|check` — scaffold + build + drift-check the consumer's `app.css`
   (`@source`, `@theme`, `@custom-variant` dark) in one command.
3. Canonical CSS-drift guard recipe (docs page + CI snippet) — D6; five consumers rebuilt
   this wheel independently (CV safelist gate, dnsblockd regen-ui, quotaboard corpus-fresh,
   cqrs-htmx check-css-bundles, cvagent gotcha 24).
4. Surface the four theme presets (default/minimal/glass/emerald) in README + tc scaffolding —
   M1; the two unhappiest adoptions were aesthetic mismatches that presets answer.
5. `@theme` token-completeness audit: every library color family overridable via semantic
   tokens (kills the ChastityAPI gray-vs-token conflict class that blocked DataTable and
   ValidationSummary adoption).
6. "Remap, don't replace" bridge recipe with copy-paste `@theme` blocks — `docs/theming.md`
   exists; surface and modernize it.
7. `@source` scanning contract doc for all three consumer layouts (vendored / replace /
   workspace) — CV `_tc-src`, dnsblockd `regen-ui`, and invoices vendor paths each reinvented it.
8. Tailwind version lockstep doc (dev shell / demo Dockerfile / website npm lanes, currently
   4.3.3, F-20) written for consumers, not just maintainers.
9. Dark-mode integration recipe: class-strategy `@custom-variant` per consumer setup
   (plain / BuildFlow / nix).
10. Generalize CV's `gen-tc-safelist` pipeline into a fleet-shared, documented recipe.
11. Decision doc: precompiled `styles.css` distribution artifact for zero-build consumers
    (still 100% Tailwind output, no pipeline) — ADR if adopted; see g2.
12. Demo: a "Tailwind-max" showcase page (theme override, container queries, fluid type, dark
    mode) — market the ceiling of what leveraging Tailwind here looks like.
13. Metric: time-to-first-styled-page for a fresh consumer; measure before/after Lane 0.

**Lane 1 — P1 doc quick wins (one afternoon; report findings D1–D4, S2–S3)**

14. countgen writes every count-bearing README cell, including the "By the Numbers" /
    comparison table (D1; six consumers reconciled counts by hand).
15. Define the icon count as names; reconcile 102 vs 105 vs 111 (D1).
16. State the module count once: 9 go.mods = 7 published + 2 repo-local (D1/D2).
17. Module-map table above the fold in README (D2 — the hallucination anchor; one reviewer
    invented the module list wholesale).
18. templ pairing matrix in README Requirements + "bump both pins in one commit" (D3; caused
    live breakage in index/invoices/dnsblockd).
19. "What you take on" first-screen block — worded to the GOAL: Tailwind v4 is the point, and
    here is the paved pipeline for it (D4).
20. Link `docs/icons-only-adoption.md` where pipeline-less consumers look (D4/PapDashboard;
    orthogonal to the GOAL, still true).
21. SKILL.md adopt/reject quick profile (S3) — worded so rejections cite the paved Tailwind
    path, not a Tailwind escape.
22. SKILL.md "what's in my pin?" protocol (S2; four reviews reinvented it).
23. Annotate the comparison-table count basis inline (D1/M3).

**Lane 2 — P2 (about a day; D5–D8, M1, S4)**

24. Hand-rolled → library replacement map, seeded from the seven duplicate-reporting reviews
    (D5: ThemeScript, MetaRefresh, Toast, EmptyState, Pagination, DataTable, LoadingButton,
    errorpage, RelativeTime, LivePill).
25. `docs/upstream-requests.md` ledger (gap, requester, status, issue) cross-linked to the
    CHANGELOG (D7).
26. "Consumer impact:" line convention for behavior-changing CHANGELOG entries (D8; v1.27.0's
    nonce-less validity singleton bit pbx-lens in the browser).
27. Link `docs/version-support.md` from the README counts block (D8).
28. Harvest the five recurring gotchas into `references/consumer-guide.md` (S4: explicit
    BaseProps.ID determinism, same-commit CSS rebuild, nonce for new singletons, amber hue
    semantics, paired templ bump).
29. Refresh README verification numbers (tests/goldens/guards), countgen-fed (M3 — the
    paragraph that flips SKIP→ADOPT).
30. Link the demo site + presets from the README theming section (M1).

**Lane 3 — P3 / process**

31. Release cadence: cut v1.28.0 from the ~91-commit unreleased wave OR adopt `-rc` preview
    tags; document the train policy (M2) — see g3.
32. Decide the `docs/reviews` series variant once (editorial vs dark) and record it (kit
    assignment says editorial for adoption feedback; newest sibling is dark).

**Lane 4 — product filings (GitHub issues, not prose)**

33. StatCard value slot (CV's proposal).
34. DataTable `tfoot` (journal).
35. day-unit `CompactDuration` (quotaboard).
36. Popover `position-area` fix for Chromium 156 (cqrs-htmx; their vendored corrected layer
    is upstreamable).
37. Badge prefix/icon slot (webphone).
38. errorpage shell stylesheet slot (browser-history draft staged).
39. Next-release workaround-retirement notes for #48/#49/#51/#60/#63 consumers (pbx-lens
    queue).

**Lane 5 — session/report hygiene**

40. Harvest Lanes 0–4 into `TODO_LIST.md` (+ `ROADMAP.md` for P3/product items) — see g1.
41. grep-verify the "10+ reviews cite SKILL.md" claim; tighten the report line.
42. Render-level check (headless screenshot) of the HTML brief.
43. Verify the daemon's auto-commit of the two new files carries sane messages; fix if
    hallucinated (known daemon behavior).
44. Before any push: `scripts/ci-repro.sh --lint --css --website` at the exact tip (ritual;
    docs-only changes should sail, but green must name its lanes).

(45 items total; "up to 50" satisfied without padding.)

## g) Three questions I cannot figure out myself

1. **Harvest now or after triage?** Shall I run docs-health HARVEST and write Lanes 0–4 into
   `TODO_LIST.md` (ROADMAP for the P3/product items) right away — and if so, does Lane 0
   (the GOAL) take the top priority slot — or do you want to triage section (f) first?
2. **GOAL boundary:** does "leverage Tailwind to the max" include shipping a **precompiled
   `styles.css` artifact** so pipeline-less consumers get full library styling with zero build
   (still pure Tailwind output — but it changes the "consumer runs the pipeline" contract,
   D4/D6, and the compiled-CSS inventory guards), or must the consumer-run Tailwind pipeline
   stay mandatory-by-design?
3. **Release cadence:** cut v1.28.0 from the current ~91-commit unreleased wave soon (five
   consumers gate work on "the next tag"; Code-To-CV-Agent pinned a pseudo-version), or adopt
   `-rc` preview tags and keep the minor cadence as is?

## Deliverables of this session

| File                                                               | Status                                                           |
| ------------------------------------------------------------------ | ---------------------------------------------------------------- |
| `docs/reviews/2026-11-25_docs-marketing-skill-audit.html`          | New, untracked, validated, date-corrected, GOAL ruling annotated |
| `docs/status/2026-11-25_20-15_consumer-review-synthesis-status.md` | This report                                                      |

Per the harness contract ("NEVER COMMIT unless the user says commit") neither file is
committed; the auto-commit daemon will pick them up — check its messages (f43).

**WAITING FOR INSTRUCTIONS.**
