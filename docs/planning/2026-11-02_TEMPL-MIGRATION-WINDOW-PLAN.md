# TEMPL MIGRATION WINDOW PLAN — the ONE re-baseline event (TODO #335)

**Date:** 2026-11-02 · **Status:** EXECUTED 2026-11-04 (window opened by user directive; templ-only — #368/#369 unbundled, see TODO_LIST) ·
**Owner:** templ pin, `utils.TestTemplVersionPin`, `scripts/check-templ-gate.sh`
**Bundles:** TODO #368 (size-* sweep) + TODO #369 (`@utility` migration) + the
`templ.ResolveAttributeValue` deletion — all DECIDED-deferred into this window so every
output-shifting re-baseline stays ONE event.

## 0. Why a window at all

The dev shell's templ generator MUST equal every go.mod pin (the zero-diff-generate
invariant). templ ≥ v0.3.1100 DELETED `templ.ResolveAttributeValue`, which the
committed v0.3.1020-generated `*_templ.go` files call ~100× per file — so the pin can
never move as a mechanical bump; it moves as one release event that regenerates and
re-records every golden. Until the window lands, consumers pair templ-components ≤
v1.25.0 with templ ≤ v0.3.1070. History and the daemon's repeated poisoned bumps:
`docs/agent-context-history.md#templ-version`.

## 1. The gate

```bash
scripts/check-templ-gate.sh    # prints GATE: OPEN (with recipe) or CLOSED (correct state)
```

- **CLOSED** = nixpkgs templ ≤ pin. The correct resting state. Re-probe after any
  `nix flake update` or before any planned migration.
- **OPEN** = nixpkgs (locked input) now ships templ newer than the pin. Execute this
  plan within days — the window is the ONLY sanctioned path, and `nix flake update`
  without the plan trips `TestTemplVersionPin` (designed tripwire, not a failure).
- Probe says UNKNOWN when run outside the flake — run inside `nix develop`.

## 2. Preconditions (all must hold before step 1)

1. `git status` clean; `result*` symlinks absent (release blockers per AGENTS.md).
2. Fresh CI green on master (M03: no migration rides an unverified base).
3. Load average sane — the event re-records ~330 HTML goldens + ~206 pixel goldens;
   starved runs misgrade (load-flake rules in AGENTS.md).
4. `scripts/coverage-baseline.txt` re-recording is NOT needed (statement counts are
   toolchain-specific but templ does not change Go's coverage counting — record only if
   the Go toolchain also moves in the same event).

## 3. The event, ordered (est. half a day, mostly machine time)

1. **Flake decision (the "flake-level generator pin" mechanics).** Do NOT add a
   `nixpkgs-templ` split input preemptively — the `nixpkgs-go` split was folded back on
   2026-09-09 precisely because idle splits stop insulating anything (flake.nix NOTE).
   The correct mechanics: ride the locked nixpkgs input; `TestTemplVersionPin` is the
   tripwire that makes a silent generator move LOUD. If, AFTER this migration, templ and
   Go versioning in nixpkgs diverge again, re-split per the existing NOTE text.
2. **Pin bump in lockstep.** `github.com/a-h/templ v0.3.1020` → the nixpkgs version in
   `go.work` + ALL 9 go.mods, ONE commit (the guards `TestTemplVersionPin` +
   `TestGoDirectivesAlign*` hold the lockstep). Never `go get templ@latest` — the target
   is `nix eval --inputs-from . nixpkgs#templ.version`.
3. **Regenerate, canonical form only.** `nix develop -c templ generate ./...` from the
   REPO ROOT — then per-module regeneration from INSIDE each module dir (the
   workspace-relative FileName rewrite hazard; `ci-repro.sh`'s per-module loop is the
   reference implementation). Verify: `git diff --stat` touches only `*_templ.go` +
   goldens.
4. **Source-trim review (the v0.3.1070 spacing change).** The generator emits a
   separating space after self-closing `@icons` calls: either trim the affected sources
   or accept the new golden bytes — DECIDE PER SITE, never blanket-accept
   (`rg -n 'icons\.Icon' --glob '*.templ'` for the affected set).
5. **Regenerate the CSS artifact** if markup changed: `nix run .#css` (demo
   `@source` scan picks up every `.templ`; stale app.css = red CI, 2026-10-14).
6. **Re-record goldens, same event:** HTML goldens per package
   (`go test ./pkg/... -update` — flag AFTER the package list) AND pixel goldens
   (`nix run .#visual -- -update`; NEVER bare `go test -update` in visualtest). Load
   flakes get patience (`-parallel 4`), never `-update`.
7. **Bundled cleanups (each is its own commit INSIDE the same window, ordered after the
   generator event so their diffs stay attributable):**
   - **#368 size-***: 168 `w-N h-N` pairs → `size-N` across packages + goldens
     re-baseline (visually identical; burns the golden budget once because it rides
     this window's `-update`).
   - **#369 `@utility`**: migrate utility-shaped `.tc-*` classes (`tc-squircle`,
     `tc-content-auto/-compact`, `tc-log*`, `tc-fluid-*`); component-shaped CSS stays
     plain; delete the `loadCustomCSSClasses` boilerplate. Expect byte-level compiled-CSS
     churn (utilities layer vs components layer) — covered by this window's visual
     re-baseline.
8. **ResolveAttributeValue verification:** `rg -l ResolveAttributeValue --glob
   '*_templ.go' | wc -l` → `0`. Then `GOWORK=off go build ./...` per module (standalone
   builds prove the generated code no longer needs the deleted symbol).
9. **Guards green, then release:** full verify (`nix run .#verify`), the M03
   ci-repro ritual, ONE release cut (`scripts/release.sh`) — the pairing table below
   dies with this release, so the release notes MUST state the new minimum templ.

## 4. Rollback, per stage

- After step 2 (pin bump, no regen yet): revert the single commit — zero output changed.
- After step 3/4 (regenerated, goldens not yet re-recorded): revert the commit(s);
  committed `*_templ.go` still reference `ResolveAttributeValue` only if the NEW
  generator wrote them — reverting the pin AND the files together restores v0.3.1020
  output byte-for-byte (generated files are committed artifacts, see AGENTS.md).
- After step 6 (goldens re-recorded): rollback is now a RELEASE decision, not a
  `git revert` — the visual baseline matches the new generator only. Re-cut by reverting
  the whole branch and re-running the event; never mix old goldens with new generated
  files.
- The master-tip `result*`/daemon rules apply throughout: re-check `git status -sb`
  before every push (the daemon races commits).

## 5. Verification checklist (exit criteria)

- [ ] `scripts/check-templ-gate.sh` → `GATE: CLOSED` (pin == nixpkgs again)
- [ ] `rg -l ResolveAttributeValue --glob '*_templ.go' | wc -l` → `0`
- [ ] `utils.TestTemplVersionPin`, `TestGeneratedFileNameRoots`, `TestTemplGeneratedInSync` green (`-count=1`)
- [ ] HTML golden sweeps green in every package; pixel suite green (`nix run .#visual`)
- [ ] axe sweep + touch-target + zoom-reflow audits green (visual env rides them)
- [ ] `#368`/`#369` closed; `loadCustomCSSClasses` deleted if #369 landed
- [ ] CHANGELOG `[Unreleased]` states the new generator + consumer pairing change
- [ ] Release cut tags every published module (check-release-tags.sh green)

## 6. Provenance

- Gate probe: `scripts/check-templ-gate.sh` (plan C18)
- Incident history: `docs/agent-context-history.md#templ-version`
- AGENTS.md templ section: "templ Version Pin" + "Tailwind lanes move in LOCKSTEP"
- Prior decisions: TODO #335 (window), #368/#369 rows (DECIDED-DEFERRED, C18 2026-10-30),
  ADR-0018-adjacent CSS policy for #369
