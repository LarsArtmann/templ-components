# CORPUS-CLOSE-OUT Plan Execution — Session Status (M2/M14/M15/M3 shipped; M4 nearly)

**Date:** 2026-11-04 05:35 | **Owner:** Lars + this agent session | **Status:** ACTIVE — mid-plan, waiting for instructions
**Plan:** [`docs/planning/2026-11-04_03-43_CORPUS-CLOSE-OUT-DEFECT-FIX-PARETO-PLAN.md`](../planning/2026-11-04_03-43_CORPUS-CLOSE-OUT-DEFECT-FIX-PARETO-PLAN.md)
**Scope claim (single-writer):** this session owns docs/status/, docs/planning/, TODO_LIST, ROADMAP, CHANGELOG, README, AGENTS, DOMAIN_LANGUAGE, docs/research/ + the docs-health skill asset (`~/projects/SKILLS/docs-health/`).

---

## TL;DR

Three commits pushed to master with witnessed `VERDICT: PASS (exit 0)` rituals (M03 honored twice):
`9465f35e` (the 213-file corpus sweep + defect fixes), `482e3a93` (CI lint-pin alignment), `ea5d869c`
(the 9 planning files re-processed + re-archived, prose-item ledger, strict-gate debt measured).
M14 shipped in the docs-health skill repo (`check-rows --strict`, `bare-rows`, `prose-items`, 46
fixture assertions + extended self-test). M4 is 90% done but UNCOMMITTED. The plan's remaining
milestones (M5–M13, M17) have NOT been started. Two environment wars were fought and settled:
a poisoned shared build cache and the templ-pin re-bump daemon (both restored/measured).

## New facts discovered this session (durable value)

1. **Strict-gate debt measured:** `check-rows --strict` flags **189/264 archived files carrying 626
   wholly-unstruck tables** — the class the pre-strict gate could not see (the close-out postmortem's
   admitted blind spot, now quantified). Folded into TODO #498 ① as pipeline fuel
   (close-out report §Follow-through).
2. **Prose ledger measured:** 5,419 bare prose lines across 210/264 archived files (54 clean),
   classified R (TODO-cited residue) / S (era-superseded) / R2 (harvest-wave provenance) —
   close-out report §Prose-item ledger. TODO #212's "remaining prose is narrative" claim is now
   measured, not asserted.
3. **TODO #468 is stale:** HoverCard ALREADY ships `popover="hint"` + `interestfor` +
   shared positioning (`display/hover_card.templ:44-48`) — the TODO row predates the feature.
   NOT yet struck (row scope needs a 3-point verify first: role="tooltip" drop, aria-details).
4. **Shared-cache poison class:** `/mnt/buildcache/go-build` (Forgejo runner + all dev shells share
   one GOCACHE/GOMODCACHE) produced build-cache entries that replay a parse-EOF failure for
   regexp2 even when the extracted module read byte-perfect. Proven by cache bifurcation
   (private GOCACHE green / shared red, GOEXPERIMENT-independent, mtime-touch-immune).
   A runner-raced `go clean -cache` half-nuke cleared it. Root cause (shared cache across
   CI + dev) is UNRESOLVED — recurrence expected.
5. **The templ-pin re-bump daemon is alive and armed:** the FULL-mode BuildFlow daemon
   (`buildflow --fix --build-mode=full --go-upgrade-blame all`, another session's pts/15) re-applied
   templ v0.3.1100 to ALL 18 go.mod/go.sum files during my first ci-repro run (04:47), killing it;
   restored to v0.3.1020 twice (`TestTemplVersionPin` green each time). Run 2 survived because the
   daemon's 15-min budget expired.
6. **AGENTS additions landed (uncommitted):** templ `//`-comment whitespace-leak gotcha (sourced
   from the archived 2026-06-28 06-43 report — comments inside markup leak newlines into output;
   doc comments above `templ` declarations are stripped) and the jsonv2 jsontext/guard/CI-lane
   mentions (errorpage handler + docs_count_test use `encoding/json/jsontext`;
   `TestTemplGeneratedInSync` pins generated imports; ci.yaml exports jsonv2 on all Go lanes).

---

## a) FULLY DONE

| Work                                                    | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| ------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **M2 — ship the sweep** (T2.1–T2.6)                     | TODO_LIST full read; scope-checked staged set (229 entries); 2 commits; M03 ritual ×2 with PASS verdicts; pushed `b8a4bb1a..482e3a93`; origin tip == local                                                                                                                                                                                                                                                                                                                                            |
| **M14 — gate improvements** (T14.1–T14.4)               | `annotate.go`: `--strict` UNTOUCHED flag, `bare-rows` moved-file gate, shared `scanTables` refactor (gates cannot drift), `prose-items` scanner; self-test extended (strict/bare/prose fixtures all green); fixture script 31→46 assertions; SKILL.md gate docs + ARCHIVE gate row updated; gofmt+dprint clean; daemon-committed in SKILLS repo (2 commits, content verified intact)                                                                                                                  |
| **M15 — re-process the 9 planning files** (T15.1–T15.5) | 54 untouched tables inventoried; 169 rows verdicted against source (rg probes per claim); 15 narrative/metrics/decision tables ledger-marked; all 9 files pass `check-rows --strict` + `bare-rows`; re-archived via `git mv`; ADR-0033's 2 citations repointed; TODO #498 ②/④ struck with evidence                                                                                                                                                                                                    |
| **M3 — prose-item ledger** (T3.1–T3.4)                  | Scanner shipped; corpus scanned; 217-line per-file ledger appended to the close-out report; straggler spot-checks verified (Pagination `uint` shipped, SimpleNav BaseProps shipped, testutil = open TODO #34); TODO #498 ③ struck                                                                                                                                                                                                                                                                     |
| **M4 — docs tails, 6 of 7 tasks**                       | T4.1 README icons section (IconRTL + AnimatedIconWithAnimation — REAL signatures verified against `icons/animated_icon_templ.go`); T4.2 "Modern Web Platform, By Default" section; T4.3 "Last verified" row (points at the TestDocsCountDrift recompute command); T4.4 DOMAIN_LANGUAGE nonce omit-empty contract (wire rows already existed); T4.5 AGENTS templ-comment gotcha; T4.6 AGENTS jsonv2 jsontext/guards (211/220 line budget held); T4.7 starfederation refs swept → go-datastar fork pins |
| **Infra fixes along the way**                           | regexp2 module cache re-extracted (zip intact, extracted dir was 0-byte); poisoned shared GOCACHE cleared (half-nuke); templ pins restored ×2; dprint-clean 119 staged md files (verified whitespace-only via `git diff -w`); lint-pin skew resolved in a DEDICATED commit per policy                                                                                                                                                                                                                 |
| **Two green M03 rituals**                               | `VERDICT: PASS (exit 0)` at 05:16:29 (commits 1+2) and 05:30:09 (commit 3), both followed by immediate push + tip re-check                                                                                                                                                                                                                                                                                                                                                                            |

## b) PARTIALLY DONE

| Work                  | State                                                                                                                                                                                          | Rest                        |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| **M4 / T4.7b**        | FEATURES FULLY_FUNCTIONAL spot audit NOT done (the starfederation half of T4.7 is done)                                                                                                        | One focused audit pass      |
| **TODO #468 strike**  | HoverCard popover=hint verified shipped, but the row's other two claims (drop `role="tooltip"`, `aria-details`) not yet verified — row left open deliberately                                  | 3-point verify, then strike |
| **M13 / TODO #498 ①** | The pipeline exists (M14 gates) and the fuel is measured (626 tables / 189 files), but neither the 2026-10/11 non-archived corpus nor the archived-debt burn-down has processed a single batch | The L-tier sweep sessions   |
| **M4 wave commit**    | README/AGENTS/DOMAIN_LANGUAGE/research edits are UNCOMMITTED working-tree changes (made after `ea5d869c` pushed)                                                                               | Commit + ritual + push      |

## c) NOT STARTED

M5 (#493 test/consistency quick wins), M6 (#490 dark-scan ring-offset/shadow rules), M7 (#489
keyboard/a11y bundle), M8 (#495 visual tooling), M9 (#497 demo tails), M10 (#488 release/infra
hardening), M11 (#492 table tests), M12 (#494 NotFound404 polish), M13 sweep batches, M17 (ROADMAP
residue go/no-go). M16 (owner-gated set) untouched by design.

## d) TOTALLY FUCKED UP (honest incident list)

1. **First commit attempt failed 3 BuildFlow checks** (dprint-format, test-compile ×2) — I committed
   without a pre-flight `dprint check` over the staged md set. The 119-file formatting detour was
   self-inflicted; one pre-commit dry-run would have saved ~20 minutes.
2. **The corrupted-cache forensics took the scenic route:** re-extract → still broken → touch
   experiment (wrong theory, wasted a cycle) → GOEXPERIMENT bisect → only THEN the private-GOCACHE
   bifurcation, which was the decisive probe and should have been FIRST for a
   "file reads fine but go sees EOF" symptom.
3. **`go clean -cache` raced the Forgejo runner** → half-nuked cache. I checked runner activity once
   and proceeded; a wait-for-idle loop was the correct move (it worked out — the half-nuke evicted
   the poison — but by luck, not method).
4. **Measurement bug:** `grep -q complete` matches "incomplete" → reported 0/264 strict-fails; the
   earlier `case *complete*` glob had the same substring flaw and mis-reported for a different
   reason (stdout/stderr interleaving). Exit-code checking is the only correct gate read; I
   re-measured properly (189/264) but burned two cycles.
5. **Launched ci-repro run 1 while the FULL-mode BuildFlow daemon was still active** — it re-bumped
   all templ pins mid-run and killed the lane. I knew the daemon existed (AGENTS documents it) but
   did not check for live `buildflow --fix` processes before starting a 15-minute ritual. The
   pin-watchdog only existed from run 2 onward.
6. **Two `rg -rn` strikes** in one session — the exact trap AGENTS.md warns about by name, blocked
   by the hook both times. Third strike host-wide; this is my personal recurrence to kill.
7. **Brief false alarms from sloppy log reading:** misread a truncated log tail as an early ci-repro
   failure; misread an append-mode exit file (`>>`) as "the run aborted". Both cost a re-check
   cycle each; the piped-capture discipline ("never trust a capture missing its final line")
   applies to MY OWN tails too.
8. **Protocol drift:** I edited md files (M4) while ritual 2 was running. Inert to every lane, but
   it violates the clean-tree-during-ritual rule I stated at the start of M2.
9. **AGENTS.md edit had to be re-done** after `gofmt -w` invalidated my read (stale-read guard fired
   correctly; my workflow caused the churn).

## e) WHAT WE SHOULD IMPROVE

1. **Pre-flight before every commit:** `dprint check` + the fast guards over the staged set, not
   after BuildFlow yells. (Cheap: one command.)
2. **Daemon check before rituals:** `ps aux | grep 'buildflow --fix'` before launching ci-repro;
   wait out its budget or accept the race knowingly.
3. **Bifurcation-first debugging:** for "tool A sees X, tool B sees Y" file symptoms, isolate the
   environment (private cache dir) FIRST — it separates data corruption from env poisoning in one
   probe.
4. **Gate reads by exit code only** — never by substring of output.
5. **Cache hygiene:** the shared /mnt/buildcache GOCACHE/GOMODCACHE across CI runner + dev shells is
   a standing corruption risk (see question 3).
6. **Session-start `ps` audit** for sibling BuildFlow/daemon processes, same as the existing
   `git fetch` ritual — cheap and prevents the 04:47 class.
7. **The 626-table debt needs a batching plan** (est. 54 tables ≈ one focused session at M15's care
   level) — decide depth per class before burning sessions on it (see question 2).

## f) NEXT 50 (ordered by the plan's Pareto tiers; ⫱ = needs your input)

1. Commit + ritual + push the M4 wave (uncommitted: README, AGENTS, DOMAIN_LANGUAGE, research doc).
2. T4.7b: FEATURES FULLY_FUNCTIONAL spot audit.
3. Verify + strike TODO #468 (HoverCard: role="tooltip" drop + aria-details).
4. M5.1 `feedback.Alert`: replace `templ.Raw` concat with `utils.ScriptComponent` (+ omit-empty nonce).
5. M5.2 Regenerate Alert goldens; feedback tests green.
6. M5.3 `utils.AssertContainsClass` helper + unit test.
7. M5.4 `FuzzCombobox` + seed corpus.
8. M5.5 Tooltip custom-trigger-element test.
9. M5.6 Accordion >1000px golden + CountBadge `Max=0` case.
10. M5.7 RelativeTime `Intl.RelativeTimeFormat` locale assertion.
11. M5.8 Combinatorial flag tests (top 5 two-bool components).
12. M6.1 dark-scan rules: `ring-offset-*` + `shadow-*` pairs (allowlist Toggle-thumb `bg-white shadow-sm`).
13. M6.2 Fix whatever the new rules flag in library/demo/site.
14. M6.3 Golden re-baseline for touched components; visual suite light+dark.
15. M9.1 Demo: `GridColsAutoFit` section on the layout page.
16. M9.2 Demo: `Card.Header` slot showcase (dashboard recipe header).
17. M9.3 `forms_demo.templ` → `forms_standalone.templ` rename decision (or record NOT-DO).
18. M9.4 Docs code-block filename/lang chrome decision + website md renderer note.
19. M7.1 MobileMenu focus-trap vs Tab-escape decision + implementation.
20. M7.2 Listbox type-ahead (first-letter) for Combobox options.
21. M7.3 Combobox PageUp/PageDown + tests.
22. M7.4 Tabs keyboard parity audit (server-rendered `clientSide` gate).
23. M7.5 Carousel `aria-live` slide announcements.
24. M7.6 Dropdown arrow-key e2e (chromedp poll helpers).
25. M8.1 visualtest `#tc-root` padding configurable (default unchanged).
26. M8.2 50× same-route double-capture variance lane (off by default).
27. M8.3 Overlay/blink-comparison HTML generator for golden-change PR review.
28. M10.1 release.sh `--dry-run`.
29. M10.2 release.sh step-0 preflight (clean tree, go.work, signing key, fetch check).
30. M10.3 Push guidance: single `--follow-tags` vs `--atomic` decision + doc.
31. M10.4 CI shellcheck job (actionlint parity: nix shell nixpkgs#actionlint nixpkgs#shellcheck).
32. M10.5 Release checklist: verify-tag step + tooling-hardening summary rule.
33. M10.6 Worktree lifecycle script (/tmp worktree + go.work copy + preflight).
34. M10.7 `scripts/test-all-modules.sh` (sanitized log paths).
35. M10.8 Compiled-CSS provenance stamp in the 3 tracked artifacts' build targets.
36. M11.1 Table Flush+Bordered interaction golden.
37. M11.2 Pagination `current>total` case in the existing table.
38. M11.3 F109 demand-gated feature set verdicts (record, no build).
39. M12.1 NotFound404 `fullscreenCenter` const + 404+ThemeToggle composition test.
40. M12.2 `HomeHref` alias + `NotFoundPageProps` wrapper decision.
41. M12.3 Globe/Ghost icon option + `notFound404Search` reuse decision.
42. M13.1 2026-10/11 non-archived corpus inventory (~90 files, UNTOUCHED-aware).
43. M13.2–13.4 Verify/annotate/archive batches + manifest append.
44. M13+ Archived-debt burn-down: 626 untouched tables in 189 files (batched; ⫱ depth per question 2).
45. M17.1 ROADMAP website-growth residue go/no-go.
46. M17.2 ROADMAP library-API/testing/docs-ADR ideas triage + TODO promotions.
47. Close the loop on the cache-poison incident in a durable artifact (BuildFlow issue or a
    docs/agent-context-history entry + TODO row — currently it lives only in this report).
48. Watch the FIRST CI run carrying the v2.15.0 lint pin (committed today; CI's `go install` must
    resolve it — local guard green but CI itself hasn't run yet).
49. TODO #498 ⑤: verify the 5 kept-open status files' residue rows are all represented in TODO_LIST.
50. Rebuild the stale BuildFlow binary (doctor warns: built at 74a9212, repo HEAD 783ae81) or record
    the deliberate non-fix.

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **The templ-pin daemon (pts/15 FULL-mode BuildFlow `--go-upgrade-blame all`):** it re-bumped all
   18 module files to templ v0.3.1100 mid-ritual and I restored twice. It is another session's
   process and the bump class is TODO #335-guarded. Kill/coordinate it during push windows, leave it
   alone, or should I file the "respect TestTemplVersionPin-class pins" ask upstream in BuildFlow?
2. **The 626-table archived debt (189 files):** burn it down to zero at M15's per-table
   verification depth (~10+ focused sessions), ledger-mark the narrative-looking tables wholesale
   (fast, shallower judgment), or leave it recorded as debt with no scheduled burn? This sets the
   M13+ lane's size.
3. **Shared-cache split:** the Forgejo runner and every dev shell share one /mnt/buildcache
   GOCACHE/GOMODCACHE, and today's poison proved the coupling can break ALL Go builds on the box.
   Fix now in flake.nix (per-repo/per-user cache dirs — touches the dev env for everything here) or
   record it as a TODO row and keep healing incidents as they surface?

---

_Rituals honored: M03 (two witnessed PASS verdicts, pushed immediately, tips re-checked), single-writer
scope claimed and held, counts/versions machine-guarded (no hand-typed counts landed), annotate.go dry-run
first for every corpus write. Session ends here per instruction — waiting._
