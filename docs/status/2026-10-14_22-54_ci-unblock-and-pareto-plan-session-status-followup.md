# Status: CI-UNBLOCK + PARETO-PLAN SESSION — FOLLOW-UP SNAPSHOT

**Written:** 2026-10-14 22:54 CEST | **Repo:** templ-components |
**My session's tip:** `21642619` — still **origin/master**, still **CI-green**
(run 38629386225: all 5 jobs ✓) | **Local master:** `e3d35692`, **ahead 8**
(new, unpushed, _not mine_ — see Δ) | **Predecessor:**
`2026-10-14_22-32_ci-unblock-and-pareto-plan-session-status.md` (this file
updates it; its evidence tables remain valid — not repeated 1:1).

## Δ What changed since the 22:32 snapshot (22 minutes)

1. The 22:32 status report was **collected by the daemon** and committed
   (`6c0c66a5`, bundled in an auto-commit of 4 files).
2. **The concurrent session is wide awake.** Between 22:32 and 22:54 it
   landed **8 local commits** (unpushed): a `docs(plan)` Pareto plan for
   "NotFound404 hardening (shell, noindex, i18n)" (`3030a84d`), errorpage
   code + golden changes (`72e45ed0`, `e8dfccde`), a lint-fix +
   tc-mirror-resync commit (`85601b1c`), plus auto-commits. Working tree
   clean again.
3. **None of those 8 commits has run CI** (origin has no newer push than
   20:07). When the daemon or the peer pushes, CI will judge errorpage work
   that this session deliberately never touched (⫱3 held all day).
4. My session's own deliverables are unchanged and verified: plan on origin
   (`3fbfa811`), five CI root causes fixed, tip green.

---

## a) FULLY DONE (evidence-cited; details in the 22:32 report, §a A1–A14)

| #   | What                                                                                                                                       | Evidence                                                                        |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------- |
| A1  | Pareto plan (27 M / 98 F, tiered, gated) written, committed, pushed                                                                        | `3fbfa811`; content verified (403 lines, 27+98 rows)                            |
| A2  | **Master CI fully green at my tip** — Lint 1m30s, Build & Test 3m18s, HTML validation 8s, CSS Freshness 23s, Visual Regression 6m13s       | `gh run view 38629386225` → `✓ master CI`; still the newest origin run at 22:54 |
| A3  | ci.yaml drift gate: FileName-form tolerance (Build & Test lane)                                                                            | green in CI on `0465abd4`/`21642619`; script body executed locally, exit 0      |
| A4  | website.yml golangci-lint v2.13.2 → v2.14.0 (Go 1.27 export-data fix)                                                                      | Website job **success** on `a10665b5` (2m53s)                                   |
| A5  | MetaRefresh goldens → full-document `<head>` render (HTML-validation lane)                                                                 | local sweep clean (334 files, 18 ignore classes); CI job green (8s)             |
| A6  | Gate script actionlint/shellcheck-clean (SC2016 ×2 + SC2086 → `norm()` + escaped backticks + directive)                                    | CI's exact pairing (actionlint v1.7.12 + shellcheck) → exit 0; CI Lint green    |
| A7  | MetaRefresh spec-canonical separator `; url=` (component, tests, goldens, inventory)                                                       | `go test ./layout/` green; no other `;url=` renderer in repo                    |
| A8  | Class inventory regenerated same-arc                                                                                                       | `TestClassInventoryFreshness` green `-count=1`; in `0465abd4`                   |
| A9  | aria-controls checker divergence ledgered as documented ignore class #18                                                                   | `scripts/check-html-valid.sh` clean after                                       |
| A10 | M03 ritual: `ci-repro.sh --lint --website` → `VERDICT: PASS (exit 0)` (21:51:55) + tree-identity proof for the delta push                  | ritual output; `git diff 86bcd89c 0465abd4` empty                               |
| A11 | 3 daemon races survived: no force-push, no lost edits, content integrity re-verified after every grab                                      | linear history `3fbfa811 → a10665b5 → 0465abd4 → 21642619`                      |
| A12 | `[Unreleased]` warmed ×3 (drift gate + shellcheck note, lint pin, MetaRefresh)                                                             | CHANGELOG diffs in `a10665b5`/`21642619`                                        |
| A13 | Table-view report-back + owner gates re-asked in chat                                                                                      | prior assistant message                                                         |
| A14 | 22:32 status report written; daemon-collected as `6c0c66a5`                                                                                | `git log -- <report file>`                                                      |
| A15 | ⫱3 respected all session: zero touches on the peer session's files (`utils/darkmode_*`, `cmd/tc/lint.go`, and today-evening `errorpage/*`) | git history — no such paths in my commits                                       |

## b) PARTIALLY DONE

| #  | Works                                                                                              | Remains                                                                                                                                                                    | Blocker                                  | Effort      |
| -- | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------- | ----------- |
| B1 | Green master at `21642619` on origin                                                               | 8 **unpushed local commits** (peer session's errorpage hardening) have **never run CI**; if red on push, master red-again — outside my scope to pre-verify (⫱3)            | peer's push timing                       | —           |
| B2 | Interim red runs on master during my convergence (`a10665b5` inventory gap, `0465abd4` actionlint) | historical only; both superseded by green `21642619`                                                                                                                       | none                                     | —           |
| B3 | Lessons documented in CHANGELOG + commit messages                                                  | **AGENTS.md / agent-context-history still not updated** (actionlint/shellcheck parity gap; gate FileName tolerance; daemon-race protocol; peer-liveness humility — see d7) | none, pure discipline debt               | S           |
| B4 | Self-review folded into status reports per single-file template                                    | standalone `docs/reviews/*.html` series episode not written                                                                                                                | deliberate (template demanded one `.md`) | S if wanted |
| B5 | Owner gates ⫱1–⫱3 asked 4× (chat ×3 + plan + report)                                               | unanswered → M11, M24, M26, M27 edge-stopped; harvest not run                                                                                                              | Lars's call (g)                          | —           |
| B6 | Dual-checker bridge (18 ignore classes) covers current corpus                                      | no automated CI-parity harness for the local nixpkgs checker                                                                                                               | tooling gap (f#43)                       | M           |

## c) NOT STARTED

| #  | What                                                                                                          | Why                                               | Wanted?                     |
| -- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------- | --------------------------- |
| C1 | **All 27 M-tasks** of the committed Pareto plan (0 executed; Tier-1 `@variant dark` M04→M05→M06 queued first) | CI fire consumed the session; plan gated on ⫱1–⫱3 | Yes — it is the deliverable |
| C2 | M01 integrity trio (D1–D3 in the 10-14 deep-dive HTML)                                                        | plan-sequenced; not pulled forward mid-fire       | Yes (45 min, credibility)   |
| C3 | Harvest to TODO_LIST/ROADMAP (plan seeds + report f)                                                          | ⫱1 + docs-health HARVEST routing                  | After ⫱1                    |
| C4 | M27 hygiene (BuildFlow rebuild, interrogate, govulncheck stdlib triage note GO-2026-66xx wait-for-patch)      | plan-sequenced; ⫱3 check first                    | Yes                         |
| C5 | AGENTS.md lesson capture (B3)                                                                                 | end-of-session rush; NEXT session's first edit    | Yes                         |
| C6 | Release-notes line for MetaRefresh `; url=` delta                                                             | no release cut; `[Unreleased]` covers it          | At next release             |

## d) TOTALLY FUCKED UP (radical honesty; 1–6 from the 22:32 report, still true — 7 is new)

| #  | Fucked up                                                                                                                                                                                                                                                                                                         | Root cause                                                                                                                                                         | Mitigation                                                                                                                    |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| D1 | Executed the handover's stale next-steps (commit→push→ritual) without `gh run list` first — the push had already happened; master was 10+ red runs deep                                                                                                                                                           | trusted summary over ground truth                                                                                                                                  | session-start ritual: fetch + status + `gh run list` BEFORE planning                                                          |
| D2 | Pushed workflow YAML without CI-parity linting; my first gate script was shellcheck-dirty (SC2016 ×2, SC2086); local actionlint silently skips shellcheck without the binary — the ritual's "same as CI" step is weaker than CI                                                                                   | parity theater                                                                                                                                                     | actionlint v1.7.12 + shellcheck together before any YAML push (f#17)                                                          |
| D3 | Daemon out-raced the session 3× (pushed pre-amend `a10665b5` mid-ritual — a commit my own gate had just proven red — and committed over my working edits)                                                                                                                                                         | small-batch editing across ~40 min vs a ~60s daemon budget                                                                                                         | batch-edits → single commit; expect `[behind N]`; delta-commit on remote tip; never amend pushed objects                      |
| D4 | Master red 18:37→22:13 (~3.5 h, ≥10 runs); my arc extended the streak twice while converging                                                                                                                                                                                                                      | breakage pre-dated the session (FileName drift 18:50, MetaRefresh goldens never green, stale website pin) but nobody watched CI — including me, for the first hour | red-cascade alarm rule (e1)                                                                                                   |
| D5 | Two sloppy first drafts of bash-in-YAML in a row (backticked sed; em-dash inside a shellcheck directive → SC1125)                                                                                                                                                                                                 | eyeballed instead of linting in CI-parity mode at write time                                                                                                       | lint the moment it is written (e2)                                                                                            |
| D6 | Cross-package guard blindness: `.templ` source edit without anticipating the class-inventory freshness guard; cost a full 6-min ritual cycle                                                                                                                                                                      | didn't enumerate the source-edit guard fan-out (templ-sync, inventory, `_sources` mirror, goldens)                                                                 | source-edit checklist (e4)                                                                                                    |
| D7 | **I misjudged peer-session liveness.** At ~21:00 I inferred the concurrent session was "probably ended (2 h stale)"; by 22:32 it had landed a 404-hardening plan + errorpage code, and by 22:54 eight commits. Only my own ⫱3 rule (never touch their files) prevented that misjudgment from becoming a collision | treated commit-staleness as session-liveness signal                                                                                                                | peer sessions get liveness-checked by _asking_ (⫱3/g3), never inferred from timestamps; keep hands off their paths regardless |

**Not fucked up (checked):** no force-push; no reverts of others' work; every
daemon grab content-verified; no pins moved; AGENTS.md size gate intact;
`[Unreleased]` warm; no unrelated research — all commands served this
session's own failure surfaces or state checks.

## e) WHAT WE SHOULD IMPROVE

1. **Session-start ground-truth ritual**: `git fetch && git status -sb && gh run list --limit 5` before trusting any handover (kills D1).
2. **CI-parity gate for YAML/bash**: actionlint + shellcheck binary, CI's version, at write time (kills D2/D5); fix ci-repro's Actionlint step to fail loudly when shellcheck is absent (f#17).
3. **Daemon-race protocol**: batch → one commit → expect rejection → fetch → delta-commit; never amend pushed objects; escalate g3 (pause daemon pushes in ritual windows).
4. **Source-edit guard checklist**: `.templ` edit ⇒ generate (root + per-module renorm) + package goldens + class inventory + `_sources` mirror + per-module test loop BEFORE committing.
5. **Dual-checker ledger**: every new golden gets both checker snapshots' verdicts (nixpkgs html5validator vs CI's fresh vnu.jar) or a documented ignore with reason.
6. **BuildFlow SQLITE_BUSY hygiene**: chronic under concurrent sessions (checkpoint losses, budget overruns); F91 rebuild + serialize heavy nix work; else file upstream.
7. **Memory discipline under the ≤220 gate**: one-line AGENTS.md + narrative in agent-context-history, same session as the lesson (B3/C5 debt).
8. **Peer-session humility** (new, from D7): staleness is not liveness evidence; ⫱3 coordination stays a standing question, and hands-off-paths is the only safe default.
9. **Report freshness hygiene** (new): point-in-time claims ("tip green") rot in minutes in this repo — status reports must carry the exact commit + run ID they certify (both reports now do) and a Δ-delta header when superseded.

## f) TOP 50 NEXT TASKS (impact-ranked; feeds docs-health HARVEST)

Impact: Critical/High/Medium/Low · Effort: S <30min / M 30min–2h / L >2h

| #  | Task                                                                                                                                                            | Impact   | Effort | Category      |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 1  | Answer ⫱1–⫱3 (section g) — unblocks M11, M24, M26, M27                                                                                                          | Critical | S      | Decision      |
| 2  | M01: verify D1–D3 claims in the 10-14 report; fix report HTML (F01–F05)                                                                                         | High     | M      | Documentation |
| 3  | M04: `.dark`-fork census + `@variant` probes (F12–F14)                                                                                                          | High     | M      | Feature       |
| 4  | M05: migrate `.dark` forks in custom.css to `@variant dark` (F15–F18; selector-set diff empty)                                                                  | High     | M      | Feature       |
| 5  | M06: verify F1 (guards `-count=1`, `nix run .#css`, freshness, demo smoke) (F19–F21)                                                                            | High     | M      | Quality       |
| 6  | M09: classify `tc-*` @utility-eligible vs component-shaped + constraint doc (F28–F29)                                                                           | High     | M      | Documentation |
| 7  | M10: convert eligible `tc-*` set to `@utility` (F30–F34)                                                                                                        | High     | L      | Feature       |
| 8  | M11: sync guards + re-baseline per ⫱2 (F35–F38)                                                                                                                 | High     | L      | Quality       |
| 9  | M07: website scrollbar native utilities (F22–F24)                                                                                                               | High     | M      | Feature       |
| 10 | Record session lessons in AGENTS.md (1-liners) + agent-context-history (narratives): actionlint parity, gate tolerance, daemon protocol, peer-liveness humility | High     | S      | Documentation |
| 11 | ci-repro parity fix: Actionlint step fails/warns loudly when shellcheck absent                                                                                  | High     | S      | Quality       |
| 12 | M08: adoption-guide refresh + consumer `@utility`/`@variant` section (F25–F27)                                                                                  | Medium   | M      | Documentation |
| 13 | M12: resolve `tc-auto-grow` vs `field-sizing-content` (F39–F41)                                                                                                 | Medium   | M      | Feature       |
| 14 | M13: `--animate-tc-*` theme tokens for 11 icon animations (F42–F45)                                                                                             | High     | L      | Feature       |
| 15 | M14: icons recipe + skill note, wrapper-vs-token (F46–F47)                                                                                                      | Medium   | S      | Documentation |
| 16 | M02: standalone `@plugin` probe, 3 lanes (F06–F08)                                                                                                              | High     | M      | Quality       |
| 17 | M25: typography plugin or documented keep (needs M02; ⚑ not pre-empting FIFTY-IDEAS #13)                                                                        | Medium   | M      | Feature       |
| 18 | M03: research README index + 10-08 forward-pointer + changelog policy (F09–F11)                                                                                 | Medium   | M      | Documentation |
| 19 | M15: Tabs `aria-selected:` refactor, delete JS class writes (F48–F51)                                                                                           | Medium   | L      | Feature       |
| 20 | M16: aria-expanded chevron cleanup, 3 trigger families (F52–F54)                                                                                                | Medium   | L      | Feature       |
| 21 | M17: `aria-[sort=ascending]:` probe + playbook note (F55–F57)                                                                                                   | Low      | M      | Quality       |
| 22 | M18: `-safe` alignment on Pagination/Nav (F58–F61)                                                                                                              | Medium   | M      | Feature       |
| 23 | M19: `@source not` demo exclusions (needs M01/F03; byte-stable proof) (F62–F63)                                                                                 | Low      | S      | Quality       |
| 24 | M20: mist preset + docs + demo + goldens + count bump (F64–F66)                                                                                                 | Medium   | L      | Feature       |
| 25 | M21: `CodeBlock.TabSize` prop + clamp + full ladder (F67–F71)                                                                                                   | Low      | L      | Feature       |
| 26 | M22: pointer-coarse escape hatch + user-valid/invalid docs (F72–F73)                                                                                            | Low      | S      | Documentation |
| 27 | M23: `w-N h-N → size-N` codemod + dry-run (F74–F76)                                                                                                             | Medium   | L      | Cleanup       |
| 28 | M24: tree-wide sweep folded into ⫱2 window (F77–F80)                                                                                                            | Medium   | L      | Cleanup       |
| 29 | M26: harvest plan + report (f) into TODO_LIST/ROADMAP + 3 standing rules (F84–F89; ⫱1)                                                                          | High     | M      | Cleanup       |
| 30 | M27: BuildFlow rebuild, interrogate decision, govulncheck triage note, full ladder (F90–F98; ⫱3)                                                                | High     | L      | Quality       |
| 31 | Split-brain: `ci-repro.sh normalize_form` vs ci.yaml `norm()` — extract shared helper or cross-reference comments                                               | Medium   | S      | Cleanup       |
| 32 | Document the two vnu snapshots' divergence (which checker emits which checks) in docs                                                                           | Medium   | S      | Documentation |
| 33 | Watch next daemon commits for templ v0.3.1070 re-bump (TODO #335; `TestTemplVersionPin` guards)                                                                 | High     | S      | Quality       |
| 34 | Keep website.yml's golangci-lint pin in lockstep with ci.yaml on future bumps                                                                                   | Medium   | S      | Quality       |
| 35 | Demo smoke over layout pages at next demo touch (MetaRefresh output changed)                                                                                    | Medium   | S      | Quality       |
| 36 | Release-notes line for MetaRefresh `; url=` at next release cut                                                                                                 | Low      | S      | Documentation |
| 37 | FIFTY-IDEAS ⚑ check when website typography starts (M25 vs #13 `display.Prose`)                                                                                 | Low      | S      | Process       |
| 38 | gitignore-upserter watch: no `*_templ.go` re-append survives after next pre-commit                                                                              | Medium   | S      | Quality       |
| 39 | Decide daemon push behavior in ritual windows (g3)                                                                                                              | High     | S      | Decision      |
| 40 | `[Unreleased]` warm + drift-count bumps same-commit for every M-task                                                                                            | High     | S      | Process       |
| 41 | `docs-health` HARVEST of report (f) + plan Deferred seeds after ⫱1                                                                                              | Medium   | S      | Cleanup       |
| 42 | Re-run + prune the html5validator ignore list at next vnu/nixpkgs update                                                                                        | Low      | S      | Quality       |
| 43 | Consider a CI lane for the nixpkgs checker (local-only classes currently invisible to CI)                                                                       | Low      | M      | Quality       |
| 44 | `gh run list --limit 5` into the session-start ritual + dispatch-template line                                                                                  | High     | S      | Process       |
| 45 | Annotate the 20-39 status doc (append-only): D1–D3 open, pointers to plan + both reports                                                                        | Low      | S      | Documentation |
| 46 | Review BuildFlow's 9 unavailable tools from preflight (F92: devShell add or N/A)                                                                                | Low      | S      | Cleanup       |
| 47 | Confirm website lane stays green after next dependabot run (templ-1070 PRs to `/website` fail by design — nobody "fix" them)                                    | Medium   | S      | Process       |
| 48 | Retro note: interim red runs `a10665b5`/`0465abd4` self-healed within one push cycle — verify that stays true                                                   | Low      | S      | Process       |
| 49 | After M10: revisit functional `@utility tc-fluid-*` (`--value(--text-*)`, plan Deferred seed)                                                                   | Low      | M      | Feature       |
| 50 | Next-session retro: re-read (d) here; verify rules e1–e9 are actually followed                                                                                  | Medium   | S      | Process       |

## g) THREE QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **⫱1 — Harvest timing:** fold the plan's M/F items + this report's (f) into
   `TODO_LIST.md`/`ROADMAP.md` now, or only after your plan review? (Tried:
   plan encodes it as a gate; repo says nothing about your cadence.)
2. **⫱2 — Re-baseline bundling:** do M11 + M24 share one standalone
   golden/visual re-baseline window, or ride the queued templ-1070 window?
   (Tried: the 1070 landing date is not derivable from the repo.)
3. **Peer session + daemon:** the concurrent session just landed 8 unpushed
   commits (errorpage 404 hardening) while I was reporting. Do you want me to
   (a) stay fully clear of errorpage + let their work push/verify on its own,
   (b) coordinate actively (and how, with which session?), and — related —
   can the daemon's unasked pushes be paused during ritual windows? I cannot
   see other sessions' plans, and I cannot change the daemon from this repo.

---

**WAITING FOR INSTRUCTIONS.** No new code was written this window (report
only); no unrelated research — all checks served my own artifacts (report
file, pushed commits' CI) plus the git state I was asked to notice.
