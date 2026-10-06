# Status Report — 2026-10-06 14:4x CEST — Daemon Templ War ×3, tc-init Root-Cause Rework, Ritual Complete (Unpushed)

Branch `fix/nonce-omit-empty` @ `3230b7f7` — **7 commits ahead of origin, NOT pushed.**
Tree clean. **All three ritual lanes witnessed green at `da8c75a9`**: `nix run .#lint`
EXIT=0, `scripts/ci-repro.sh --lint --website` VERDICT: PASS (exit 0), `nix run .#visual`
EXIT=0. The push command was the literal next action when this report was demanded.

**Note:** the prior summary claimed "9 commits ahead" — reality at session start was 44
(the daemon had pushed and raced continuously). Measure, never trust a carried count.

## a) What was done (this continuation session, `3c250654..HEAD`, 23 commits incl. daemon snapshots)

1. **Templ-bump war, round 3+4 — caught, reverted, durably defended.**
   - Session-start inspection found daemon commit `1e130c1a` had re-applied the templ
     v0.3.1070 go.mod bump across 8 modules — the bump deliberately rolled back
     2026-10-05 (generator-output change; TODO_LIST #335). Verified the dev-shell
     generator is still v0.3.1020 (nixpkgs has NOT caught up; flake.lock nudge irrelevant).
   - Reverted via `go mod edit` + `go mod tidy` per module (never hand-edited manifests).
   - **New guard `utils.TestTemplVersionPin`**: pins every module's go.mod templ version
     to `cmd/tc/doctor.go`'s `templGeneratedWith` constant — fires on partial bumps in
     either direction. Mutation-tested (fails at 1070, passes restored).
   - The daemon re-bumped AGAIN mid-session (`befc13ba`), and a parallel BuildFlow run
     re-bumped root/website/visualtest seconds after my second revert. Third revert
     landed atomically (`da8c75a9`; edit+tidy+commit in one shell).
   - **Durable defense: `go-auto-upgrade` added to `.buildflow.yml` `skip_steps`** —
     same class as the previously skipped `go-structure-linter`/`go-version-auto-configure`
     (the tool's own repair re-creates a state the repo explicitly rejected).
     Documented in `docs/agent-context-history.md#templ-version`.
2. **tc-init scaffold root-cause rework (supersedes the prior session's starter-CSS
   deletion rationale).** The hook's tailwind-build lane failed my templ-guard commit:
   `Can't resolve './templ-components-theme.css'` — the starter `app.css` imports a file
   `tc init` never copied. **The prior session misdiagnosed the starter theme copy as a
   zombie; the byte-equality canon (starter mirrors `templates/`) was right all along,
   and `tc init` has scaffolded a broken first compile since the v2.0
   "semantic aliases included by default" flip.** Coherent fix, all in one commit
   (`b7a27880`): `tc init` scaffolds all three files; starter mirror restored
   (`TestStarterCSSMatchesTemplates` green); the CSS inventory guard refuses only true
   compiled artifacts (`styles.css`, `*.out.css`) and refuses them by GIT-TRACKING so
   the provider's untracked compile residue can never fail a commit; starter compile
   output gitignored; copy instructions (templates/app.css header, adoption guide) now
   say "all three". **Proven end-to-end**: fresh `tc init` in /tmp + Tailwind v4.3.3
   compile → 17.2 KB CSS. AGENTS.md CSS-inventory bullet carries a misdiagnosis note.
3. **AGENTS.md diet landed: 448 → 377 lines** (the preflight gate; warning gone).
   Long-form narratives moved to **`docs/agent-context-history.md`** (anchors:
   `#tagging-policy`, `#tag-tree`, `#templ-version`, `#go-directives`); every operative
   rule stayed. Same pass fixed three stale-canon items: the go-directive rule now
   states the current **bare `go 1.27` lockstep** (the 2026-10-03 "canonical go 1.26.0"
   text was superseded by the 10-05 sweep; verified go.work + all 9 modules at bare
   `1.27`), the Architecture header (`Go: 1.27 floor`), and the ScriptComponent
   convention (#350 unification completed — layout runtime injection is the only
   exception). Formatter ping-pong (prettier reflow 377→378) cost two re-trims.
4. **#349 docs sweep DONE** (`9f0ecf00` + daemon `97bce89c`): README build-flag callout
   ("until Go 1.27" → older-toolchains-only), README requirements (1.26+ + flag →
   1.27+ no flags), `docs/cli.md`, website `installation.md`, ADR-0013
   realized-dated retirement note. Verified the flake dev shell is go1.27.1 (the
   version-support table's "1.27.1" claim is correct). TODO_LIST #349 struck; next-ID
   pointer fixed to 353.
5. **#337 session tripwire DONE** (`0c0d4d39` + `83f86492` + `3230b7f7`):
   `scripts/session-tripwire.sh <start-sha>` — dirty-tree gate, diff stat,
   deleted-then-reappeared (re-tracking class), added-then-deleted (resurrection
   class), and daemon/real-commit overlap flags (torn-snapshot signature). TODO_LIST
   #337 struck. First real run flagged exactly the templ-war daemon overlaps
   (adjudicated: files diff-reviewed; HEAD state proven by the pin guard + full lanes).
   Self-found bug: `grep -v` no-match exits 1 under pipefail — fixed same session.
6. **Carried from the prior session (context, already committed):** #351 license-gate
   root cause (GOROOT prefix mismatch; machine-local wrapper), ScriptComponent
   unification (#350), PolledRegion role=region pin, starter zombie sweep (#290/#346),
   both status reports annotated, SKILL.md canon, LICENSE files on 6 sub-modules.

## b) Verification (witnessed, named lanes)

- `nix run .#lint` — **EXIT=0**, "0 issues" ×7 modules + actionlint (run twice; the
  first capture truncated before its final line, so it was re-run to witness the ending).
- `nix develop -c scripts/ci-repro.sh --lint --website` at `da8c75a9` —
  **"VERDICT: PASS (exit 0) — 2026-10-06 14:30:10 CEST"**. First run at the earlier tip
  FAILED (TestStarterCSSMatchesTemplates) — which is how the tc-init root cause
  surfaced; the failure was the system working.
- `nix run .#visual` at `da8c75a9` — **EXIT=0** (full pixel suite incl. axe sweep).
- `utils.TestTemplVersionPin` + `TestStarterCSSMatchesTemplates` +
  `TestCompiledCSSInventory` + `TestDocsCountDrift` + `cmd/tc` suite — green.
- End-to-end scaffold proof (temp dir `tc init` + `nixpkgs#tailwindcss_4` compile).
- Hook discipline: every manual commit this session went through the full pre-commit
  hook with zero `--no-verify`; two daemon hook-bypassed snapshots were amend-repaired
  (`b7a27880`) or superseded. `scripts/session-tripwire.sh 3c250654` verdict:
  no re-tracked deletions, no resurrected zombies; 3 daemon/real overlap flags on
  go.mod/go.sum — all adjudicated (the templ war itself).

## c) State of the branch

- `fix/nonce-omit-empty` @ `3230b7f7`, **ahead 7**, tree clean, remote untracked (a
  fetch showed origin at `1e130c1a`; the 7 local commits are the war reverts, the
  tc-init rework, tripwire, docs sweep, and AGENTS diet).
- CHANGELOG `[Unreleased]` warm (nonce omit-empty, ScriptComponent, TagsInput,
  playground, trigger-modifier fix, **tc-init scaffold fix**, **pre-1.27 docs
  correction**, Go-1.27 sweep fallout, CSP sweep coverage, site-content docs, new
  guards).
- TODO strikes this session: #349, #337. Open near-term: #335 (templ 0.3.1070
  migration plan), #352 (BuildFlow upstream GOROOT fix), #338 (version-support
  checklist extension), #339 (golden diff readability), #348 (post-merge daemon
  re-verify).

## d) Daemon incidents this session (all handled, all documented)

1. Templ bump re-applied twice more (`1e130c1a`, `befc13ba`) + one parallel-run
   re-bump between my edit and commit → three reverts total; skip-step defense landed.
2. Daemon committed my templ-guard test hook-bypassed (`a81edbed`) — surfaced the
   tailwind-build failure that exposed the tc-init bug (accidental win).
3. Daemon snapshotted the tripwire script mid-write (`0df1e9ad`, 82 lines) before my
   real commit completed it (`0c0d4d39`) — tripwire's overlap check now flags this
   class by design.
4. Daemon committed the 9-file tc-init rework hook-bypassed (`8825a187` → amended to
   `b7a27880` with the real message).
5. Formatter (prettier) reflowed AGENTS.md past the 377 gate twice → re-trim commits.
6. `go run ./cmd/tc init` smoke-test mishap: `go run` executes with the caller's CWD,
   so the subshell `(cd repo && go run ...)` scaffolded INTO the repo root — caught by
   `ls` before the daemon could commit; trashed; re-run properly with a built binary.
   (Lesson: never `go run` a CWD-sensitive tool through a subshell; build first.)

## e) Decisions made (this session)

- **Byte-equality canon restored over the "resolvability" divergence** — the first
  starter fix (opt-in comment + `@tailwindcss-output` directive) traded one invariant
  for another and tripped `TestStarterCSSMatchesTemplates` in full ci-repro. The
  rework keeps BOTH: scaffold all three files, mirror templates/ byte-identically,
  and make the provider's compile residue untracked-by-construction.
- **Inventory refusal is git-tracked-based, not disk-based** — matches the existing
  `*.out.css` guard philosophy ("untracked litter is reported informationally") and
  prevents a commit-doomsday loop against the provider's legitimate rewrites.
- **`go-auto-upgrade` skipped repo-side** pending TODO_LIST #335 / upstream respect
  for pinned rejects.
- **Go-directive canon rewritten** to the current truth (bare `go 1.27` lockstep,
  spelling-stability rule, one-commit bump choreography) — the historical saga moved
  to `docs/agent-context-history.md#go-directives`.

## f) Next tasks (ranked, 1 = do first)

| #     | Task                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | Why it ranks here                                                                              |
| ----- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| 1     | **PUSH the 7 commits now** (`git push origin fix/nonce-omit-empty`) — the ritual completed at `da8c75a9`; every minute unpushed is un-backuped work the daemon can race                                                                                                                                                                                                                                                                                                                                                                                                                       | The literal next action interrupted by this report; zero open verification                     |
| 2     | Open the PR (#336): github-voice skill first; two-section body (nonce omit-empty feature; guards/sweep/playground/docs/unification); `Fixes #N` links; recommend v1.21.0 cut after merge                                                                                                                                                                                                                                                                                                                                                                                                      | Long-running branch; CI + Website workflows must run on the PR                                 |
| 3     | Post-merge daemon re-verify (#348 checklist): `nix run .#css` byte-stability, website typescript pin, CI status                                                                                                                                                                                                                                                                                                                                                                                                                                                                               | Standing checklist after daemon-heavy sessions like this one                                   |
| 4     | Cut v1.21.0 (pending owner timing answer, §g-2): pre-verify lint + touched packages, then `nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh v1.21.0 "<summary>"` inside the dev shell                                                                                                                                                                                                                                                                                                                                                                                       | [Unreleased] is warm; release.sh refuses empty — it is not empty                               |
| 5     | Watch for a 4th templ-bump round: after push, `rg -n 'a-h/templ v' --glob go.mod` should show 1020 everywhere; if not, the skip-step did not reach the parallel session's BuildFlow — consider `nix run .#reinstall` of BuildFlow or pausing the daemon                                                                                                                                                                                                                                                                                                                                       | Three rounds today; the defense is brand-new                                                   |
| 6     | Add `docs/agent-context-history.md` + `scripts/session-tripwire.sh` to the docs-count guard's awareness only if counts apply (they do not) — instead: lychee-linkcheck the new history file's anchors (`#templ-version` etc.)                                                                                                                                                                                                                                                                                                                                                                 | New doc with anchor links; lychee may not know the anchors                                     |
| 7     | #338: extend `docs/version-support.md` "What a floor bump looks like" with the missed steps (goldens, pins, docs floors, AGENTS, consumer-note probe)                                                                                                                                                                                                                                                                                                                                                                                                                                         | The 1.27 sweep shipped the floor but missed goldens+docs on first pass (f15)                   |
| 8     | #339: golden diff readability (windowed/word diff in `utils/golden`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | LCS diffs on templ output near-unreadable during triage                                        |
| 9     | #352 + go-auto-upgrade upstream: file BuildFlow issues (license GOROOT export; skip/respect pinned rejects in go-auto-upgrade) — user-owned upstream, needs §g-3 answer                                                                                                                                                                                                                                                                                                                                                                                                                       | Fleet-wide value; repo-side workarounds are in place but fragile                               |
| 10    | `TestTemplVersionPin` should also pin go.work's effective templ requirement (today it reads the 9 go.mods; a flake-level templ drift is only caught by regen)                                                                                                                                                                                                                                                                                                                                                                                                                                 | Cheap hardening while the war is fresh                                                         |
| 11    | The provider compiles starter/app.css on every hook run (~200 ms × every commit) — consider a BuildFlow exclude-path feature request (compile-check the scaffold ONCE in CI instead)                                                                                                                                                                                                                                                                                                                                                                                                          | Pure waste today; untracked litter by design                                                   |
| 12    | Tripwire: add `--since <date>` convenience (derive start-SHA from `git log --until`) so the habit needs no SHA recall                                                                                                                                                                                                                                                                                                                                                                                                                                                                         | Ergonomics; script is brand-new                                                                |
| 13    | Tripwire: unit-smoke in a scratch repo (its pipefail bug shipped to "done" untested — do not repeat)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | Same class as the §d-6 go-run mishap: tool not proven before trusted                           |
| 14    | AGENTS.md: the CSS-inventory bullet now names starter/{app,custom,templ-components-theme}.css — reconcile the "Compiled artifacts (exactly 3)" count phrasing (starter/styles.css residue is untracked-by-design, not an artifact)                                                                                                                                                                                                                                                                                                                                                            | Text-level canon hygiene; the bullet grew a correction sentence                                |
| 15    | Consider `docs/agent-context-history.md` pruning policy (it grows forever; date-delimited sections suggest an annual archive cut)                                                                                                                                                                                                                                                                                                                                                                                                                                                             | New file; set the convention before it needs one                                               |
| 16–50 | (Carried from `docs/status/2026-10-06_13-12_...md` §f, unchanged in priority by this session: #335 templ migration plan details, #338/#339 above, art-dupl binary-version logging discipline, golines/manual-wrap reconciliation, visualtest `-parallel 4` load-sensitivity documentation, `~48 raw chromedp.Poll` migration (#240), env-guard for `rg -rn` mangling, demo `/demo/api/datastar/stream` SSE-via-Hosting bypass, ogshot regeneration of stale `website/public/og/home.png`, BuildFlow reinstall wrapper fragility, TODO_LIST hygiene pass, FEATURES.md sync for the new guards) | The prior queue remains valid; this session added no consumer-facing debt beyond what is fixed |

## g) Open questions for the owner (3, carried + sharpened)

1. **Push/PR now?** The ritual is complete and witnessed at `da8c75a9` (7 commits now).
   Standing instruction says push immediately — I was about to run
   `git push origin fix/nonce-omit-empty` when this report was demanded. Unless you
   say otherwise, the next action is push, then the #336 PR (github-voice, v1.21.0
   recommendation).
2. **Release timing for v1.21.0** (carried, unanswered): cut immediately after
   merge, or batch with the templ #335 migration / next feature set? The
   `[Unreleased]` section is warm and release-clean; the longer it stays warm the
   more the daemon can regress it.
3. **BuildFlow upstream (#352 + go-auto-upgrade pin-respect)** (carried, unanswered):
   file both upstream in `larsartmann/BuildFlow`, or keep the workarounds repo-side
   for now? Note the machine-local go-licenses wrapper dies on any BuildFlow
   reinstall — the upstream GOROOT fix is the only durable form of #351.
