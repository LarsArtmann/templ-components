# Status: Datastar ↔ HTMX Parity — Execution Session 1 (T0 + start of T1)

**Date:** 2026-10-02 12:43 CEST
**Plan:** `docs/planning/2026-10-02_11-05_datastar-htmx-parity-master-plan.md`
**Scope of this session:** the whole plan, owner-authorized ("get shit done"); owner
blockers resolved autonomously and recorded.
**Verdict in one line:** Governance done, PolledRegion + LoadingButton shipped
**browser-proven**, wire triggers implemented + tested (ADR/docs pending), go-directive
churn killed at the guard level; 9 auto-commits ahead of origin, full verify not yet
witnessed at tip.

---

## a) Fully done (implemented, tested, verified)

| Item                                                  | Evidence                                                                                                                                                                                                                                                                                                                                                                         |
| ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **L1-32 owner decisions recorded**                    | Master plan §7 rewritten: keep `PatchMode` naming; canonical `go 1.26`; PolledRegion-first order                                                                                                                                                                                                                                                                                 |
| **L1-31 superseding ADR**                             | `docs/adr/0042-datastar-surface-expansion.md` (lifts the ADR-0035 freeze for exactly 4 surfaces; cost/rejected-alternatives honest); ADR-0035 status annotated; ADR-0038 cross-linked                                                                                                                                                                                            |
| **L2-01.1 bundle facts**                              | `on-interval` (default 1000ms; `duration` parses `ms`/`s`/bare-ms; **`leading` must be INSIDE the duration mod set**) + `on-intersect` (`full`/`half`/`threshold.<0-100>`/`exit`/`once`) decoded from pinned bundle v0.6.1 and recorded in `docs/datastar-runtime-facts.md`                                                                                                      |
| **L1-01 `datastar.PolledRegion`**                     | `polled_region.go`/`.templ` + generated; interval `data-on-interval__duration.<n>`; m/h normalized to seconds (both runtimes misparse them as ms); empty URL inert; reuses `LivePoliteness`; timestamp footer; **no Eager prop by design** (leading + outer self-patch = infinite eager loop — the DiscordSync outage class)                                                     |
| **L1-02 PolledRegion lenses**                         | 3 goldens (default/custom/inert), unit+edge tests (interval table, escaping, BaseProps propagation), a11y lens, bdd lens, example, **bundle guard tokens** (`on-interval`, `on-intersect`, `duration`, `leading`)                                                                                                                                                                |
| **L1-03 PolledRegion browser e2e**                    | `visualtest/datastar_polled_region_e2e_test.go` — `TestDatastarPolledRegionBrowser` **GREEN**: ≥2 ticks fetched+patched, exactly one region element (no duplication), request rate interval-bound (loop class fails by two orders)                                                                                                                                               |
| **L1-04 `datastar.LoadingButton`**                    | `loading_button.go`/`.templ` + generated; indicator-signal + `data-show` label swap, ZERO component JS; empty signal degrades inert (bare-`$`-truthy trap); 3 goldens, unit/a11y/bdd lenses, example, **e2e GREEN** (busy visible during 1.2s flight, rest restored)                                                                                                             |
| **L1-05 wire triggers — implementation + tests**      | `utils/wire/triggers.go`: `NormalizedInterval` (single source; `datastar.PolledRegion` now delegates), `Reveal` struct, `Action.Interval` + `Action.Reveal`; htmx merged `hx-trigger` (event+debounce+revealed/intersect+every, comma grammar); Datastar per-source attributes; **interval-only Actions suppress the default click**; 2 guard tests + 10 trigger tests **GREEN** |
| **L1-27+L1-29 go-directive canonicalization + guard** | go.work + all modules normalized to `go 1.26`; `TestGoWorkDirectiveMatchesRootGoMod` now compares major.minor-NORMALIZED (daemon `1.26.0` flips pass by construction); NEW `TestGoDirectivesAlignAcrossWorkspace` walks every go.mod; **proven**: fires on real `1.27.1` skew, tolerates daemon flip, restores green (with `-count=1`)                                           |
| **L1-28 (substance)**                                 | `.buildflow.yml` `go-structure-linter` skip confirmed present (lines 17–29); AGENTS.md doc note lands with the docs pass                                                                                                                                                                                                                                                         |
| **Docs count sync (×3 rounds)**                       | FEATURES/AGENTS/README/ROADMAP/SKILL bumped same-edit each time: now 123 components / 125 generated / 267 goldens; `TestDocsCountDrift` GREEN                                                                                                                                                                                                                                    |

**Verification state:** `datastar` module suite GREEN; `utils` module suite GREEN
(incl. new guards); wire package GREEN; lint 0 issues on datastar + visualtest +
utils/wire; 3 new browser e2e tests GREEN via `nix run .#visual`.

## b) Partially done

1. **L1-05 trigger language (the remaining 20%)**: `wire/triggers_test.go` exists and
   passes, but **ADR-0043 (the trigger-language ADR) is NOT written**, and
   `docs/transport-wiring.md` + FEATURES/skill/website still say "polling/reveal stay
   dialect-local" — now false. `datastar-runtime-facts.md` and the ADR-0042 citation
   already point the right way, but the contract docs need the trigger section.
2. **L1-28**: skip confirmed, but the AGENTS.md "daemon flip is now harmless by
   construction" note is not yet written.
3. **Commit hygiene (L1-30)**: 9 daemon auto-commits (`b20e746f`…`95186c1f`) hold this
   session's work in ~9 snapshots; folding into semantic commits not started.
4. **Todo list**: updated through L1-05; a full re-sync to the executed state is pending.

## c) Not started (from the plan)

- **L1-06** Swap/`{mode}` browser e2e extension (closing the ADR-0038 waiver)
- **L1-07–L1-10** wire ergonomics: `Get/Post/...` constructors, builders + dedupe,
  `ThrottleMS`, `doc.go` godoc
- **L1-11** naming decision already recorded (plan §7) — only the doc/addendum remains
- **L1-12** ViewTransition; **L1-13** Tabs `Wire`; **L1-14** SimpleNav `Wire`;
  **L1-15** LoadMore `Swap` migration
- **L1-16** Calendar `settle:0s` doc; **L1-17** Datastar confirm evaluation
- **L1-18–L1-20** demo Swap/remove cards + smoke route
- **L1-21–L1-22** recipes/website/DOMAIN_LANGUAGE/skill updates (incl. the now-stale
  "polling stays dialect-local" rows)
- **L1-23** `nix run .#verify` + `ci-repro --lint --website` witness at tip — **NOT run
  this session** (only per-module suites); **L1-24** full visual run pending
- **L1-25** real fuzz sessions; **L1-26** one-object demo-level guard

## d) Fucked up (own mistakes this session)

1. **Near data-loss**: wrote my PolledRegion e2e to `polled_region_e2e_test.go` — an
   EXISTING htmx e2e file. The write tool refused the overwrite (mtime guard); renamed
   to `datastar_polled_region_e2e_test.go`. Lesson: `ls` before writing new files.
2. **Violated the "never patch Go via shell" rule 3×**: perl one-liners on test files
   produced an illegal rune literal and a mangled assertion; a bash `cat >>` appended
   test helpers. All repaired via the edit tool, but each cost a build/test cycle for
   zero benefit. The rule exists because it keeps being true.
3. **Botched edit on skill/SKILL.md**: replaced a sentence's beginning with a literal
   "121→122 placeholder" string (wrong old_string target). Caught and fixed within one
   edit, but it shipped into a daemon commit.
4. **Both adapted e2e servers dropped the `HeadContent` `__dsReady` ready-flag script**
   (copy-adaptation miss). Cost: two failed 30s browser runs before a throwaway debug
   probe revealed the polling itself worked and only my flag was missing.
5. **Test-design bug caught late**: zero `Reveal` initially rendered `intersect` while
   the doc claimed "revealed" — violated the zero-value-equals-documented-default rule.
   Fixed by inverting to `EveryEntry` (zero = fire once), but only after the first test
   run failed. Should have applied the repo's own mandatory convention at design time.
6. **Identity no-op edit + out-of-order edits** on `triggers_test.go` tangled the import
   fix across three attempts before green.

## e) What could be improved (process)

- **Adaptation checklist for e2e servers**: when copying a server helper, diff the copy
  against the original BEFORE running (HeadContent flags, nonces, routes).
- **Write the ADR in the same pass as the contract change** — L1-05's code is done but
  its ADR/docs are not; that split is how drift starts.
- **Apply "zero value = documented default" at struct-design time**, not at first test
  failure.
- **Debug-probe pattern worked**: a 30-line throwaway chromedp probe (console/errors/
  DOM dump) root-caused in one run what two blind 30s retries missed — reach for it on
  the FIRST e2e failure.
- **golangci-lint run --fix immediately after writing each file** instead of batching —
  the fixes were trivial (gci/wsl/golines) but arrived late.

## f) Next things (ranked; the whole remaining plan)

1. Write ADR-0043 (typed trigger language) + annotate ADR-0036/0042; update
   `docs/transport-wiring.md` trigger section (replace the "stays dialect-local" note)
2. Update FEATURES.md wire scope note + skill Part 1 wire table row for Interval/Reveal
3. Add `datastar.PolledRegion`/`LoadingButton` + wire triggers to the demo (L1-18/19
   cards; `pages.go` meta entry; demo CSS recompile; `/demo/datastar` card copy)
4. L1-20 demo smoke route/golden pinning the Swap surface; L1-26 one-object guard at
   demo level
5. L1-06 extend `visualtest/wire_e2e_test.go` with `Swap`/`{mode}` cases (both runtimes,
   incl. option-overrides-header)
6. L1-07 `wire.Get/Post/Put/Patch/Delete` constructors + table tests
7. L1-08 builders (`WithEvent/WithContentType/WithDebounce`) + delete duplicated
   default-apply helpers in forms/display/navigation (+ their test mirrors)
8. L1-09 `ThrottleMS` (`throttle:<n>ms` ↔ `__throttle.<n>ms`) + tests
9. L1-10 `utils/wire/doc.go` godoc (one-object rule, Selector/Swap, URL-template
   convention)
10. L1-13 Tabs `Wire` field + both-dialect goldens; L2-13.3 scope note
11. L1-14 SimpleNav `Wire` (TODO #155) + goldens + docs row
12. L1-15 LoadMore `Swap: PatchModeOuter` migration (keep `hx-target="this"`, reveal
    trigger via the new typed trigger!) + goldens
13. L1-16 Calendar `settle:0s` exception doc (code comment + transport-wiring note)
14. L1-17 Datastar confirm evaluation (bundle decode first; implement only if it
    degrades gracefully; ADR otherwise)
15. L1-11 naming addendum (already decided: keep `PatchMode`; one paragraph in
    ADR-0038/0042 refs)
16. L1-12 `ViewTransition` (`useViewTransition` ↔ `transition:true`) + tests
17. L1-21 recipes: transport-migration row for triggers; website api-reference +
    wire-guide counts/sections
18. L1-22 DOMAIN_LANGUAGE.md Patch-Mode/trigger rows + skill Part 2 interop note
19. L1-28 finish: AGENTS.md daemon-flip note (normalized guard makes flips harmless)
20. L1-25 real fuzz sessions (`-fuzz=FuzzAction -fuzztime=30s`, `-fuzz=FuzzDecodeForm
    -fuzztime=30s`) + corpus notes
21. L1-23 `nix run .#verify` at tip, then `scripts/ci-repro.sh --lint --website`;
    record witness output in a status doc (RITUAL M03: witness before any push)
22. L1-24 `nix run .#visual` full pass; triage load-flakes vs real mismatches; fix
    goldens only if genuinely stale
23. L1-30 commit hygiene: fold the 9 daemon snapshots into semantic commits (Read ADR
    hygiene rules first; `--no-verify` windows per L2-30.2)
24. HTML-valid gate: `nix shell nixpkgs#html5validator -c scripts/check-html-valid.sh`
    over the new goldens (PolledRegion/LoadingButton/trigger pages)
25. Update `docs/datastar-runtime-facts.md` provenance table (v0.6.1 pin line — the
    table still says v0.5.0; bundle is byte-identical so hash unchanged)
26. `TestCompiledCSSInventory` / CSS freshness re-check after demo CSS recompile
27. Replace-directive guard re-check (`scripts/check-replace-directives.sh`) after any
    go.mod edits (visualtest/go.mod was edited today)
28. go.sum tidiness: `go mod tidy` in visualtest after the go-directive edit
29. Extend `integration/csp_nonce_test.go` inventory if any new component renders
    scripts (PolledRegion/LoadingButton render none — assert that stays true)
30. `TestFeaturesEnumTableExhaustive`/`ValuesExhaustive` re-run after docs pass
31. Demo `wire_demo.templ` e2e for the new demo cards (L2-18.2/19.2 test rows)
32. Website docs prose counts (`api-reference.md`, `guides/invariants.md`) sync after
    all adds settle
33. Final drift-guard sweep: `TestDocsCountDrift`, `TestVersionMatches*` (version
    still unreleased — CHANGELOG `[Unreleased]` must gain the trigger/PolledRegion/
    LoadingButton entries in the same pass)
34. CHANGELOG `[Unreleased]` entries: Added (PolledRegion, LoadingButton, wire
    Interval/Reveal triggers), decision notes (ADR-0042, ADR-0043)
35. Todo list full re-sync + plan checkboxes
36. Consider `htmx.PolledRegion` trigger consolidation onto `wire.Action` (optional,
    separate decision — NOT silently)
37. Optional: `charts/echarts` untouched — sanity `go build` in the final verify only
38. Optional: benchmark suite run for the two new components (repo convention)

## g) Questions I cannot answer myself

1. **Push policy for the 9 daemon commits**: fold into semantic commits first (L1-30),
   then witness `ci-repro` at the folded tip and push — or push the daemon snapshots as
   they are and clean up in history later? (RITUAL M03 says witness-then-push; L1-30
   says fold first. Both before push is my recommendation — confirm.)
2. **`wire.Action.Reveal` API shape**: current = `Reveal *Reveal` (pointer, nil = none;
   zero Reveal = "revealed"/fire-once via inverted `EveryEntry`). Confirm before ADR-0043
   freezes the public surface — after the ADR, renaming is churn.
3. **`Interval`/`Reveal` vs `htmx.PolledRegion`**: leave `htmx.PolledRegion` as-is
   (component owns its trigger) while `wire.Action` gains the generic triggers (current
   direction), or migrate `htmx.PolledRegion` onto the wire trigger too in this batch?
   (Migration is cleaner long-term but touches a shipped component for no consumer
   demand — my recommendation is leave it; asked because it shapes ADR-0043's scope.)

---

**Current tree state:** master ahead 9 (daemon snapshots); dirty: `module_surface_test.go`
(new workspace guard), `triggers_test.go` (formatting), `visualtest/go.mod` (canonical
`go 1.26`). Per-module suites green; full `.#verify` witness pending. NO push (not
requested; ritual not yet satisfied).
