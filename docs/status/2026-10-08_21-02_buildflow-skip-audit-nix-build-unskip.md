# Status Report — BuildFlow skip audit: nix-build unskipped (root cause fixed), templ skips re-confirmed

- **Date:** 2026-10-08 21:02 CEST
- **Session scope:** User asked "still valid???" about three `.buildflow.yml` `skip_steps` entries (`go-mod-update`, `nix-build`, `go-auto-upgrade`). Session = verification + the documented fix path for `nix-build` + stale-doc repairs. Per instruction, no unrelated research; observations are limited to what this session touched or directly surfaced.
- **Commits this session (daemon-carried):** `6e30a1e6` (flake.nix), `3d34979c` (.buildflow.yml + AGENTS.md). All content verified intact post-commit.

## TL;DR

| Skip              | Verdict     | One-line reason                                                                                                |
| ----------------- | ----------- | -------------------------------------------------------------------------------------------------------------- |
| `go-auto-upgrade` | Still valid | templ pinned v0.3.1020 in all 9 go.mod; TODO #335 open; `utils.TestTemplVersionPin` enforces                   |
| `go-mod-update`   | Still valid | Same root cause as above (added today, so trivially current)                                                   |
| `nix-build`       | **Removed** | Root cause reproduced live, then fixed: treefmt `goimports` now runs through an offline go wrapper; step green |

**Bottom line:** `buildflow -s nix-build` went from structurally-failing (documented 2026-10-07) to **✔ 1 success, 0 failed, 2 checks on x86_64-linux, 4.3s**. The two templ-pin skips remain correct and must NOT be removed until TODO #335 lands.

---

## Self-review (asked explicitly: forgotten / better / improve)

**What did I forget?**

1. **The config in front of me.** After the flake fix, I ran `buildflow -s nix-build` and got `pipeline build failed: no executable nodes after compilation` — and treated it as a BuildFlow behavior mystery. I spent 4 tool calls source-diving the BuildFlow repo (errors.go, pipeline_builder.go, provider listings, commit diffs since the binary build) before concluding the boring truth: `nix-build` was STILL in `skip_steps`, and a skipped step cannot compile into a single-step pipeline. The skip list was literally in the user's first message. **Fix order should have been: remove skip → run step.**
2. **Read-before-edit discipline.** The AGENTS.md multiedit failed once ("you must read the file first") because I sourced the target text from `rg` output instead of `view`. Correct guard, wasted round trip.
3. **Off-by-one on `view` offset** (it is 0-based) when reading TODO_LIST row #335 — cost one extra call.

**What is stupid that we do anyway (observed, not introduced today)?**

1. **Skip rationale comments drift.** The 2026-10-07 `nix-build` rationale named "golines/templ formatters" — neither exists in this repo's treefmt config (only nixfmt/gofumpt/goimports). The _root cause_ (go download in sandbox) was right; the _specifics_ were wrong. Rationale comments are read months later as fact; wrong specifics mislead the next debugger (it misled me for one read today).
2. **Advisory FAIL semantics.** BuildFlow's preflight printed `preflight FAIL [workspace/pseudo-version-hygiene]` in warning-red, then the run passed 1/1. A verdict word that lies ("FAIL" on a passing run) trains humans to ignore preflight output.
3. **Daemon rewrite risk on config files.** The 2026-10-07 daemon rewrite of `.buildflow.yml` silently reshaped entries; today I had to re-verify my own edits survived `3d34979c`. This remains a live hazard for every hand-written config.

**What could I have done better?**

- Grep AGENTS.md's tagging-policy section BEFORE settling the pseudo-version question (3 calls to conclude something the repo policy already answered: sibling requires ride release versions by design).
- State the verification envelope explicitly up front: x86_64-linux only; aarch64 check untested; only the `format` check and the `nix-build` step verified — not the full `nix flake check`, not other BuildFlow lanes.

**What could I still improve?**

- The goimports-offline wrapper changes goimports' failure mode (`GOPROXY=off`): "falls back to stdlib-only fixing" is cqrs-htmx's validated claim, not yet probed on THIS tree (item f-2).
- Two places now record the nix-build story (.buildflow.yml log + AGENTS.md gotcha). Agreed today; they can drift. Accepted with eyes open (repo-wide pattern), but it is a split brain by construction.

**Did I lie?** No. Every green claim names its lane: `nix build .#checks.x86_64-linux.format` (496 files, 0 changed) and `buildflow -s nix-build` (1/1, x86_64-linux, 4.3s). Nothing outside that envelope was called verified.

---

## a) FULLY DONE

| # | Work                                                                                                                                                                                                                                                                                      | Evidence                                                                                                                                                                                                                                                                                                                                                                                               |
| - | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1 | Templ-pin state verified → both templ skips re-confirmed valid                                                                                                                                                                                                                            | 9/9 go.mod @ v0.3.1020; TODO_LIST #335 row open (unmigrated); `utils/templ_pin_test.go:22` present and pinning vs `cmd/tc/doctor.go`                                                                                                                                                                                                                                                                   |
| 2 | `nix-build` skip root cause reproduced live                                                                                                                                                                                                                                               | `nix build .#checks.x86_64-linux.format` failed with exactly the documented error: 6× `go: downloading go1.27.0` → `proxy.golang.org` DNS refused in sandbox                                                                                                                                                                                                                                           |
| 3 | Root cause fixed in flake.nix (the 2026-10-07 note's own fix path)                                                                                                                                                                                                                        | `goimports-offline` writeShellApplication wrapper: `goToolchain` first on PATH + `GOTOOLCHAIN=local` + `GOPROXY=off` (cqrs-htmx pattern); `lib` added to perSystem args; stale "(1.26.7)" comment replaced — flake.nix:411                                                                                                                                                                             |
| 4 | Fix verified at both levels                                                                                                                                                                                                                                                               | Check derivation green: `formatted 496 files (0 changed)`; BuildFlow step green: `✔ nix-build 4.3s`, "1 success, 0 failed, 2 check(s) on x86_64-linux"                                                                                                                                                                                                                                                 |
| 5 | Skip removed from `.buildflow.yml` + rationale converted to `UNSKIPPED 2026-10-08` note (file's documented-history style)                                                                                                                                                                 | commit `3d34979c`; no `nix-build` skip entry remains                                                                                                                                                                                                                                                                                                                                                   |
| 6 | Stale AGENTS.md claim #1 fixed: "eslint-fix / nix-build are in skip_steps" → both UNSKIPPED with the real blocker + the new GOTCHA (future go-shelling formatters must use the offline wrapper)                                                                                           | AGENTS.md CI & Tooling Gotchas                                                                                                                                                                                                                                                                                                                                                                         |
| 7 | Stale AGENTS.md claim #2 fixed: "go-structure-linter and go-version-auto-configure are skipped" → corrected (go-structure-linter re-enabled with the 2026-10-05 lockstep bump; gvac unskipped 2026-10-07 via `respect_patch_floor: true`; go-auto-upgrade + go-mod-update remain skipped) | AGENTS.md Go-directives bullet                                                                                                                                                                                                                                                                                                                                                                         |
| 8 | BuildFlow `workspace/pseudo-version-hygiene` preflight FAIL investigated and documented as a known false positive for this repo                                                                                                                                                           | Root go.mod requires were v1.20.0 pre-drift; daemon `go mod tidy` (commit `6de0ae7b`, 2026-10-05 21:15) refreshed to v1.20.1 after the release — this is the repo's release-version pin convention (tagging policy 2026-09-25, issue #27), not a regression; gate is advisory (today's nix-build passed 1/1 with it present). Documented in AGENTS.md with "do NOT fix go.mod to zero pseudo-versions" |
| 9 | Daemon-commit integrity check                                                                                                                                                                                                                                                             | `6e30a1e6` (flake.nix, 20:55:56) + `3d34979c` (2 files, 21:00:52) both contain this session's edits verbatim (post-commit `rg` spot-checks green)                                                                                                                                                                                                                                                      |

## b) PARTIALLY DONE

| # | Work                         | What is missing                                                                                                                                                                                                                                          |
| - | ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | nix-build fix verification   | x86_64-linux only. `checks.aarch64-linux.format` never built (would need an aarch64 builder/emulation); full `nix flake check` (both systems, both checks) not run; other BuildFlow nix lanes (nix-flake-check etc.) not re-run after the wrapper change |
| 2 | `GOPROXY=off` fallback claim | The "goimports falls back to stdlib-only fixing" behavior is adopted from cqrs-htmx's validation, not probed on this tree (no deliberately-broken-import test was run)                                                                                   |
| 3 | Pseudo-version drift history | Conclusion based on 2 of 6 drifted requires (utils, icons) + the policy section; the other 4 modules' pre-drift states were not individually diffed                                                                                                      |
| 4 | Session todo list            | 4/4 completed — but items f-1..f-6 (this report's tail) are the unfinished verification fringe by design                                                                                                                                                 |
| 5 | This report                  | Written in `.md` per explicit user instruction (skill default is HTML) — see closing note                                                                                                                                                                |

## c) NOT STARTED (noticed; deliberately out of scope)

| # | Item                                                                                                                       | Why not started                                                                                                   |
| - | -------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| 1 | TODO #335 (templ v0.3.1070 migration: flake generator pin, 9 go.mod bumps, golden + pixel re-baseline, source-trim review) | Own plan; gates re-enabling the two remaining skips                                                               |
| 2 | Re-enabling `go-auto-upgrade` / `go-mod-update`                                                                            | Blocked on #335 (or upstream BuildFlow learning pinned rejects)                                                   |
| 3 | BuildFlow binary rebuild/reinstall                                                                                         | Doctor advises (`binary built at ec8d2d3, HEAD e4ad71a`); advisory only, and BuildFlow is its own project/session |
| 4 | Upstream BuildFlow: pseudo-version-hygiene gate vs release-version pin policy                                              | Upstream decision, not a repo edit                                                                                |
| 5 | lychee excludes for `docs/feedback/archived` (preflight warn, exact fix printed)                                           | Preflight warn only; touches docs tooling untouched this session                                                  |
| 6 | AGENTS.md size budget (378 lines vs BuildFlow's 220; excess 158)                                                           | Recurring warn; a docs-restructuring pass, not this session                                                       |
| 7 | HARVEST of section (f) into TODO_LIST/ROADMAP                                                                              | Awaits user instruction (skill: report first, then wait)                                                          |

## d) TOTALLY FUCKED UP

**Nothing shipped broken.** Every change was verified green before stopping, and all three commits landed intact. Honest ledger of own-goals and hazards:

1. **The `no executable nodes` detour (biggest own-goal).** 4 wasted tool calls source-diving BuildFlow for a behavior whose cause was the skip list in the user's own first message. Lesson recorded in section e.
2. **Rationale-comment rot is systemic, not incidental.** The 2026-10-07 nix-build note carried wrong formatter names; the flake's treefmt comment still claimed go.mod toolchain "1.26.7" two days after the 1.27 lockstep bump. Nobody's re-verification pass had caught either — including the pass that wrote the note.
3. **Parallel-session exposure observed live.** During this session, `charts/echarts/echarts_templ.go` went modified and the 20-42/20-45/20-46 status docs were edited by another session — while the daemon committed my files between my edit and my verification. The AGENTS.md "daemon reverted session edits mid-session" class is a standing hazard; today's spot-checks passed, but the window exists on every edit.
4. **Small process stumbles:** read-before-edit guard trip (1 call), 0-based view off-by-one (1 call), `PIPESTATUS` unsupported in the shell interpreter (cosmetic, redone cleanly). No data at risk at any point; nothing useful was deleted; no ghost system was created (the wrapper is wired into `checks.format` and verified end-to-end; both docs point at real config).

## e) WHAT WE SHOULD IMPROVE

1. **"Check the config in hand before source-diving."** The first reflex on a tool-behavior surprise must be the project's own config (skip lists, excludes, env), not the tool's source.
2. **Re-verification notes must re-verify their own specifics.** When writing "still fails because X", name the exact formatter/config lines as of THAT day, or the note becomes a trap. Pair rule: every unskip updates its rationale in the same commit (done today — keep it).
3. **Verify the envelope, then say the envelope.** x86_64-only, single-check verification is fine — as long as the report says so (this one does).
4. **Kill the advisory-FAIL confusion upstream** (BuildFlow preflight verdict wording) so operators stop training themselves to ignore red.
5. **Record the new coupling:** the go-directive lockstep bump checklist must now include "flake's `go_1_27` rev ≥ new floor, or `goimports-offline` (GOTOOLCHAIN=local) fails fast in the check sandbox" — the wrapper made the formatter chain floor-sensitive on purpose.
6. **Config-file daemon rewrites:** consider a tiny guard that every `skip_steps` entry carries an adjacent dated rationale comment — cheap to check, catches silent daemon reshaping.

## f) Things to get done next (ranked; ≤ 50 — brainstorm, not commitment; HARVEST input)

**Verify this session's fix envelope**

1. Run full `nix flake check` once (both systems, both checks) to prove no other nix lane broke from the wrapper. (~5 min)
2. Probe the `GOPROXY=off` fallback on THIS tree: add a deliberately-unresolvable import in a scratch branch, run the format check, observe fail-fast/fallback, revert. Pins the adopted claim locally.
3. Resolve aarch64: either verify `checks.aarch64-linux.format` once or trim `aarch64-linux` from flake `systems` with rationale (needs g-question 1).
4. Rebuild + reinstall the BuildFlow binary (`cd ~/projects/BuildFlow && nix build . && nix run .#reinstall`) to kill the binary-freshness WARN (needs g-question 2).
5. Re-run `buildflow -s nix-build` on the fresh binary to re-pin the green verdict.
6. Watch one daemon cycle for resurrection of the `nix-build` skip line or reverts of the new AGENTS.md content (resurrection is precedented: the CSS-zombie class).
7. Re-run `buildflow doctor` after items 4–6 and confirm 3 warn / 1 fail → 0 warn (or documented residue).

**Skip re-enable gate (TODO #335 window)**
8. Execute TODO #335: pin the templ generator at flake level (`github:a-h/templ` input), bump all 9 go.mod pins to v0.3.1070, regenerate from repo root in the dev shell.
9. Re-baseline HTML goldens for the icon+label spacing change; review whether sources should trim instead of accepting the space.
10. Re-baseline pixel visual goldens (full `nix run .#visual -- -update` per the visual ritual).
11. Source-trim review of icon+label call sites; update the AGENTS.md templ-pin section in the same event.
12. Re-enable `go-auto-upgrade` in `.buildflow.yml` (remove entry + UNSKIPPED note, today's pattern) after #335.
13. Re-enable `go-mod-update` likewise.
14. Optional de-risk: prototype the flake generator pin on a branch BEFORE the migration window so the window stays one event.

**BuildFlow upstream (separate project/session)**
15. Teach `go-auto-upgrade`/`go-mod-update` to respect repo-level pinned rejects (version floors / reject lists) — the permanent fix making both skips unnecessary.
16. Teach the `pseudo-version-hygiene` gate per-repo pin policy (or a config opt-out) — it false-fails release-version-pinning repos like this one.
17. Fix preflight verdict semantics: advisory findings must not render as FAIL on a passing run.
18. Improve the `-s <skipped-step>` error to name the skip and the config line instead of `no executable nodes after compilation` (this session's detour becomes zero calls).
19. Land the upstream go-licenses GOROOT fix (TODO #352) so the `~/.local/bin/go-licenses` wrapper requirement disappears; then simplify the `.buildflow.yml` header note.

**Repo hygiene surfaced by this session**
20. AGENTS.md 378 lines vs the 220 budget: move gotcha detail into `docs/` topic files, keep AGENTS.md as the index (kills the recurring `docs/agents-md-size` warn).
21. Add the lychee exclude for `docs/feedback/archived` (create `lychee.toml` / `.lycheeignore` as needed — neither exists today per BuildFlow's debug line); then run the link scan once to confirm archived-quote links are intentionally excluded, not broken.
22. Audit the remaining `.buildflow.yml` rationale comments for the same rot class (license-check wrapper note, pnpm-audit note, gvac note — each dated, each still true?).
23. Add the GOTOOLCHAIN=local coupling to the go-directive lockstep bump checklist (AGENTS.md Go-directives bullet + the bump ritual): flake `go_1_27` rev must be ≥ the new floor or the offline formatter fails fast by design.
24. Add a `packages.x86_64-linux.default` to the flake (e.g. the demo binary) so bare `nix build` works and BuildFlow's missing-attribute info line disappears.
25. Consider a `skip_steps` rationale guard (every entry needs an adjacent dated comment) — pre-commit-grade, ~10 lines.
26. Consider adding the templ formatter to treefmt (cqrs-htmx parity) — ONLY through the same offline-wrapper pattern if it shells out to go.
27. Annotate the same-day 20-42 status doc (`chromedp-v020-verification-buildflow-pin-war-resolution`): its nix-build references predate today's unskip — docs-health ANNOTATE mode, one pass, don't rewrite.
28. HARVEST this f-list into TODO_LIST.md / ROADMAP.md via docs-health (explicit follow-up once instructions arrive).
29. After the next release cycle, confirm requires stay at the new release version and the pseudo-version-hygiene note in AGENTS.md survived daemon rewrites.
30. Fold the "check config-in-hand first" lesson into the project's AGENTS.md tooling-gotchas only if it generalizes beyond this session (it is a session-process lesson; candidate for crush-config `references/lessons.md` instead — cross-project).

_(Stopped at 30: items 31–50 would be filler — the honest backlog from this session's observations is exhausted. Per the skill: a larger N is a brainstorm, and unharvested items belong in ROADMAP, not TODO_LIST.)_

## g) Questions I cannot figure out myself

1. **aarch64-linux:** keep the system in the flake and verify its check once (needs an aarch64 builder/emulation on this box), or trim `systems` to `x86_64-linux` with a rationale? I cannot infer the infra intent (is there an aarch64 consumer at all?).
2. **BuildFlow binary refresh:** OK to rebuild + reinstall BuildFlow from `~/projects/BuildFlow` HEAD now (shared machine-local tool used by other live sessions — the doctor flags built-at `ec8d2d3` vs HEAD `e4ad71a`), or is that owned by a dedicated BuildFlow session?
3. **TODO #335 timing:** should I prep the templ v0.3.1070 migration now (flake generator pin prototype on a branch, so the eventual window is one event), or hold everything until nixpkgs ships 0.3.1070 and the two skips just sit?

---

_Report per explicit instruction: `.md` format at the requested path (skill default is a styled HTML dashboard; the user's format demand wins). No commit made — harness forbids unrequested commits; the auto-commit daemon picks this file up. WAITING FOR INSTRUCTIONS._
