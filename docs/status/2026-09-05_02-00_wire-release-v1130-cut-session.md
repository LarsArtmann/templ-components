# Status Report — v1.13.0 Release Cut Session

**Date:** 2026-09-05 02:00 CEST
**Session scope:** Release execution only (T9) — the continuation of the 2026-09-04 wire-Pareto execution session whose T9 was staged but not cut.
**Repo state at writing:** `origin/master` = `23a61ae`, CI green (all 4 jobs, run 33931088366), `v1.13.0` + 6 sub-module tags live on origin and propagated on the module proxy (7/7 first-poll).
**Commits this session:** `cfdb94a` (release-script hardening) → `2cae44a` (release 1.13.0) → `df42cad` (re-add replaces) → `23a61ae` (post-propagation go.sum sweep).

---

## a) FULLY DONE

1. **v1.13.0 is released and live.** Release commit `2cae44a` on master; 7 SSH-signed annotated tags (`v1.13.0`, `utils/v1.13.0`, `icons/v1.13.0`, `errorpage/v1.13.0`, `charts/echarts/v1.13.0`, `datastar/v1.13.0`, `htmx/v1.13.0`) pushed; `scripts/check-release-tags.sh` reports all sub-module tags in lockstep; module proxy resolves all 7 at v1.13.0 (verified via `go list -m …@v1.13.0` against proxy.golang.org, 7/7 on the first poll).
2. **The full in-script verify suite passed inside the cut.** templ generate (117 files), 5 Tailwind distribution targets recompiled, build + `-race` tests for root + 6 sub-modules, golangci-lint for all 7 (0 issues), govulncheck v1.7.0 for all 7 (**zero reachable vulnerabilities**). Evidence: release.sh run log, attempt 3 (background shell 14E).
3. **First-ever live exercise of the vulnerability gate completed clean.** The govulncheck gate was added after v1.12.0 and had never run inside a real cut; it now has, and the run produced two concrete script fixes (next item).
4. **Release-script hardening committed (`cfdb94a`), rolled into the 1.13.0 notes.** Two gate bugs fixed: (a) the per-module govulncheck loop ran `GOWORK=off`, which cannot resolve sibling requires already bumped to the unpushed release version — now workspace-mode like the build/test phase above it; (b) the step-10b tidy sweep aborted cuts because `go mod tidy` ignores go.work and sub-modules carry no replace directives at HEAD — now warns-and-continues, with the canonical go.sum refresh left post-propagation. Plus a new `docs/release-checklist.md` checklist item: fresh worktrees must copy in the gitignored `go.work`.
5. **Every claim in the fix is empirically proven, not guessed.** Four experiments ran against a scratch bump of `icons/go.mod`: workspace-mode govulncheck PASSES on unpushed requires; `GOWORK=off` tidy FAILS (unknown revision); workspace-mode tidy ALSO FAILS (tidy ignores go.work); root tidy with replaces stripped FAILS (proxy has no v1.13.0). The warnings/behavior now match reality.
6. **Pre-flight verification caught a handoff error before it cost anything.** The handoff brief claimed "9 tags lockstep"; actual precedent (v1.12.0) is 7 tags (root + 6 sub-modules). Verified against `git tag -l '*1.12.0'` before cutting.
7. **Release tree independently verified.** In the release commit: `utils/version.go` = 1.13.0, `FEATURES.md` = 1.13.0, `## [1.13.0] — 2026-09-04` heading present, all 7 tagged go.mods replace-free with correct sibling-require counts (root 6, icons 1, errorpage 2, echarts 1, datastar 1, htmx 1 — matches the DAG), the two stale CSS artifacts (`demo.out.css`, `global.out.css`) refreshed, `visualtest/go.mod` sibling pins bumped. Confirmed the CHANGELOG has no link-ref section to update (literal headings — nothing missed).
8. **Post-propagation go.sum sweep committed and pushed (`23a61ae`).** All 8 modules tidied `GOWORK=off`; 5 sub-module go.sums moved v1.12.0 → v1.13.0 sibling checksums; re-run is idempotent; CI tip green within ~3 minutes of the push.
9. **Push done under explicit user authorization** (house rule re-confirm honored), and the second open question resolved: ADR-0037 stays Proposed → T18 (WC module phase 2) remains gated off.
10. **Daemon isolation held.** The entire cut ran in `/tmp/tc-master-wt`; zero daemon auto-commits interleaved; `origin/master` never moved underneath the release.

## b) PARTIALLY DONE

1. **pkg.go.dev indexing.** `https://pkg.go.dev/github.com/larsartmann/templ-components@v1.13.0` returned 404 at report time. The module proxy itself serves v1.13.0 (that is what consumers use, so `go get` works **now**); pkg.go.dev indexes asynchronously, typically within hours. Remaining: re-check within 24h, including the 6 sub-module paths. Effort: S.
2. **The worktree/go.work fix is documented, not automated.** The runbook checklist item exists, but release.sh does not _check_ for go.work at step 0 — a future cut from a fresh worktree still fails, just with a better-documented recovery. Effort to finish: S (add a 5-line guard).
3. **CHANGELOG `[Unreleased]` is empty** (as every release leaves it). House rule says it must be warm _at all times_; the two post-release commits were chores (replaces, go.sum) and precedent says chores don't warm it — but the next feature/fix commit MUST, or the rule is broken. Blocker: none; it warms with the next real commit. Effort: S.
4. **24-hour post-release daemon watch.** Started (origin verified stable at writing), not finished. User-side or a later session. Effort: S spread over 24h.
5. **Handoff question #2 (ownership of the 3 concurrent agent sessions).** Never answered; I worked around it with worktree isolation. Still open, and it now blocks the PR #8 decision (see g).

## c) NOT STARTED

1. **T18 — Web Components module phase 2.** Gated off: D1 = "not ratified", and you confirmed this session that ADR-0037 stays Proposed. Correctly not started; do not start without ratification.
2. **govulncheck durable provisioning.** v1.7.0 exists only as a hand-built `/tmp/tc-bin` binary — no flake input, no CI pin. The shipped 1.13.0 notes literally say "pinned v1.7.0"; that is currently aspirational. Next cut re-hits the mismatch (and `/tmp` dies on reboot). Priority: high.
3. **CSS Freshness CI scope extension.** The CI job gates exactly one artifact (`examples/demo/static/app.css`). The other four distribution targets are unchecked — `demo.out.css` and `global.out.css` had drifted stale on master and nobody noticed until this cut recompiled them. Extension not started. Priority: high.
4. **1.13.1 with PR #8** (LiveRegion busy-script nonce, other session's branch, currently red CI on that branch). Not started, not mine to start without the ownership answer.
5. **The harvested backlog** (TODO_LIST #150–154 wire block: human-eyeball wire goldens, go-datastar/static v0.5.0 bump + re-audit, interval/intersect trigger ADR, transport-symmetric Wire candidate survey, prerender sync; plus the older datastar-block items). Harvested by the prior session, execution not started.

## d) TOTALLY FUCKED UP

1. **Two aborted release attempts before the third succeeded — both self-inflicted, both preventable at pre-flight.** Attempt 1: I ran the cut from a worktree I did not verify had `go.work` — despite having _just read_ the script's header comment "GOWORK is NOT set to off — we use go.work (workspace mode)". The gitignored file exists only in the main checkout; the icons build died on a missing go.sum entry. Attempt 2: I fixed go.work but ran the govulncheck gate that had never been live-tested, without first simulating it — same failure class, same module, one full verify cycle later. Cost: ~2 wasted verify cycles (templ generate + CSS + full race-test/lint/govulncheck matrix each) plus two rollback cycles. Mitigation that existed: the rollback EXIT trap worked perfectly both times — zero tree damage, zero shipped damage. Severity: process, not product.
2. **The release notes ship a claim that is not true yet: "govulncheck (pinned v1.7.0)".** Nothing pins it. The nix flake provides 1.6.0 (too old for the Go 1.26 toolchain pairing); I hand-built v1.7.0 into `/tmp/tc-bin` via `go install`. That binary evaporates on reboot, and the next release engineer (or next session) hits the same wall with no durable answer. I shipped the prior session's wording without reconciling it against the environment I actually used. Doc-reality splitbrain, published in a tag (immutable). Severity: medium (misleads the next releaser); fix: item f#1/f#2.
3. **I recreated a documented race I knew about.** The v1.10.0 incident notes and the script's own final output say push master + tags carefully; the script even suggests `git push origin master --follow-tags`. I pushed master first and tags second as two separate commands because explicit tag listing _felt_ safer. Result: CI started on `df42cad` before the tags existed and went red on the release-adjacent commit — the exact red-race class prior releases documented. Fixed forward by the sweep commit (tip green), but the release window carries a red run it did not need. Severity: cosmetic-to-process; root cause: procedural; mitigation exists (f#5).
4. **What is NOT fucked up (verified):** the tags themselves are clean (replace-free go.mods, correct requires, signed), consumers are unaffected by the stale sub-module go.sums inside the published modules (a consumer's own go.sum governs; `go get` records what it needs — same shipped state as v1.12.0), and no history was rewritten. Product damage: none found.

## e) WHAT WE SHOULD IMPROVE

1. **A machine-checked pre-flight for cuts.** Today's checklist is prose; three of this session's four problems (missing go.work, unprovisioned govulncheck, version mismatch) were environment-readiness facts a 20-line script can assert: go.work present, `govulncheck -version` matches the required pin, signing key agent-loaded, `git status` clean, daemon idle N seconds, tags-absence on origin. Fix: `scripts/preflight-release.sh` + wire into release.sh step 0.
2. **Simulate destructive gates before the real run.** The 4-experiment pattern (scratch-bump a go.mod, run the gate, observe, restore) cost minutes and turned guesswork into proof — do it as an explicit step whenever a gate has never survived a live cut, not after the first abort.
3. **Single-transaction pushes.** `git push --atomic origin master --follow-tags` (or refs listed in one command) makes the master-before-tags race structurally impossible. `--atomic` also guarantees all-or-nothing across refs. One line in the checklist + muscle memory.
4. **Never-live-tested gates are broken gates.** Both of this session's script bugs were in code paths added _after_ the last successful cut (the vuln gate). A `--dry-run` mode (bump → verify → strip → assert on a throwaway branch, then roll back) would exercise them per-release without a real cut.
5. **Doc claims about tooling must be environment-verified before shipping.** "pinned v1.7.0" shipped unverified. Rule: a release-note claim about tooling gets checked against `command -v` / `--version` during the cut, or reworded to the requirement ("requires govulncheck v1.7.0; provision via …").
6. **Verify handoff claims against git.** The brief's "9 tags" was wrong (7). One `git tag -l` query settled it. The rule generalizes: handoffs are pointers, not ground truth.
7. **CSS Freshness should gate what ships.** One of five compiled artifacts is checked in CI; two of the unchecked ones shipped stale this release. Extend the job to diff all five outputs.
8. **Worktree lifecycle should be scripted.** Create worktree + copy go.work + verify + (optionally) preflight in one command — removes the entire "fresh worktree" failure class and the /tmp mortality footgun.
9. **TODO_LIST has a numbering splitbrain** — two separate #150–154 sequences (datastar-era block and wire-session block). Harvest numbering needs a re-baseline before the next HARVEST adds a third collision.

## f) 50 things to get done next

_Brainstorm ranked by impact — HARVEST fuel for TODO_LIST/ROADMAP, not a commitment list. Impact: Critical/High/Medium/Low. Effort: S <30min, M 30min–2h, L >2h._

| #  | Task                                                                                                                                     | Impact   | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 1  | Pin govulncheck v1.7.0 in flake.nix as a dedicated input (nixpkgs-go pattern) so the release gate is reproducible                        | Critical | S      | Quality       |
| 2  | Reconcile the shipped "pinned v1.7.0" claim with reality — amend docs/notes wording until #1 lands                                       | Critical | S      | Documentation |
| 3  | Merge PR #8 (LiveRegion nonce fix) once its CI is green                                                                                  | High     | S      | Bug           |
| 4  | Cut v1.13.1 after PR #8 (first post-release patch; warms the release cadence)                                                            | High     | S      | Release       |
| 5  | Single-transaction release push: `git push --atomic origin master --follow-tags`; document in release-checklist Push section             | High     | S      | Quality       |
| 6  | Extend CSS Freshness CI to diff all 5 distribution CSS artifacts, not just static/app.css                                                | High     | S      | Quality       |
| 7  | Add step-0 environment gate to release.sh: go.work present, govulncheck version, signing agent, daemon idle                              | High     | S      | Quality       |
| 8  | Build `release.sh --dry-run` (bump→verify→strip→assert on throwaway branch, roll back) so gates get exercised without a real cut         | High     | L      | Quality       |
| 9  | Daemon pause mechanism for release windows (the runbook says "wait or pause it" but no pause exists)                                     | High     | M      | Quality       |
| 10 | scripts/release-worktree.sh: create /tmp worktree + copy go.work + preflight in one command                                              | Medium   | S      | Quality       |
| 11 | Warm CHANGELOG `[Unreleased]` with the next feature/fix commit (it is empty post-release)                                                | High     | S      | Documentation |
| 12 | HARVEST this report's section f into TODO_LIST.md / ROADMAP.md (docs-health)                                                             | Medium   | S      | Documentation |
| 13 | Fix TODO_LIST duplicate numbering: two #150–154 sequences (datastar-era lines 54–58, wire-session lines 81–85)                           | Medium   | S      | Cleanup       |
| 14 | Bump go-datastar/static v0.4.0 → v0.5.0 with the full bundle re-audit protocol (TODO_LIST wire #151)                                     | High     | M      | Feature       |
| 15 | ADR: typed interval/intersect triggers in wire.Event (TODO_LIST wire #152)                                                               | Medium   | L      | Feature       |
| 16 | Survey next transport-symmetric Wire candidates (SimpleNav links, Form action) per the D3 rule (TODO_LIST wire #153)                     | Medium   | M      | Feature       |
| 17 | Sweep the codebase for components still hand-writing htmx/datastar attributes that should use wire (D3 symmetric rule)                   | Medium   | M      | Feature       |
| 18 | Route more demo endpoints through wire.Handler (round 2 of the adoption)                                                                 | Medium   | M      | Feature       |
| 19 | shellcheck scripts/ suite + add to CI Lint job (release.sh-class bugs found live twice this session)                                     | Medium   | S      | Quality       |
| 20 | Fixture-test the step-6 awk CHANGELOG transform (only the 8b tree assertions are fixture-tested today)                                   | Medium   | M      | Quality       |
| 21 | Verify ci-repro.sh --vuln uses the same govulncheck the release gate requires (PATH-pinning mismatch class)                              | Medium   | S      | Quality       |
| 22 | Confirm the Website workflow deployed transport-wiring.mdx to the live site post-release                                                 | Medium   | S      | Documentation |
| 23 | Decide tag-CI behavior: `gh run list --branch v1.13.0` is empty — should tags trigger CI? If yes, add `on.push.tags`                     | Medium   | S      | Quality       |
| 24 | 24h daemon watch: confirm no daemon push disturbs the release window on origin                                                           | Medium   | S      | Quality       |
| 25 | Human-eyeball the wire visual goldens (agent-captured caveat, TODO_LIST wire #150)                                                       | Low      | S      | Quality       |
| 26 | Coverage gate margin: 71.7% vs 70% floor — add missing-coverage tests or re-pin (datastar block #152)                                    | Medium   | M      | Quality       |
| 27 | TestCSSFreshness fail-capable local flag + guard against committing `.fail/` artifacts (datastar block #150)                             | Medium   | S      | Quality       |
| 28 | Fuzz `getActionExpr`/`actionExpr` URL+retry+cancellation combos (datastar block #154)                                                    | Medium   | M      | Quality       |
| 29 | DOMAIN_LANGUAGE: add busy-cue, sibling-pin policy, keep-alive frame terms (datastar block #151)                                          | Low      | S      | Documentation |
| 30 | Standing layout test: tags keep sub-module go.mods replace-free while master root keeps exactly 7 replaces (catches layout drift)        | Medium   | S      | Quality       |
| 31 | wire.Handler HTTP-level benchmark (only `Attributes()` is benchmarked today)                                                             | Low      | S      | Quality       |
| 32 | Cross-link the transport-wiring guide from htmx and datastar package docs                                                                | Low      | S      | Documentation |
| 33 | Make visualtest's untagged status explicit in check-release-tags.sh output                                                               | Low      | S      | Documentation |
| 34 | Add post-push `git verify-tag` step to the release checklist                                                                             | Low      | S      | Quality       |
| 35 | Document /tmp worktree mortality (reboot) + the recreate procedure (or fold into #10's script)                                           | Low      | S      | Documentation |
| 36 | Document whether go.work.sum should be copied into release worktrees alongside go.work (mine worked without it — is that guaranteed?)    | Low      | S      | Documentation |
| 37 | Resolve handoff question #2 (concurrent-session ownership) — blocks PR #8/1.13.1 handling                                                | Medium   | S      | Process       |
| 38 | Compiled-CSS provenance: record "last compiled at release X" so one-release-stale artifacts are visible (demo.out.css drifted unnoticed) | Low      | S      | Cleanup       |
| 39 | Annotate the 2026-09-05_01-11 status report: T9 now done (docs-health ANNOTATE, non-destructive)                                         | Low      | S      | Documentation |
| 40 | Include tooling hardening in the release summary line when the notes carry it (this cut's summary omitted it)                            | Low      | S      | Documentation |
| 41 | Fold the 3-release streak of cut lessons into one consolidated runbook / the go-release skill                                            | Low      | M      | Documentation |
| 42 | Drift-guard the counts: assert 91 visual goldens / 91 HTML goldens like TestSkillComponentCount does for components                      | Low      | S      | Quality       |
| 43 | Upstream BuildFlow fix: hallucinated daemon commit messages (documented 5+ sessions, root cause known)                                   | Medium   | L      | Cleanup       |
| 44 | Make the rule explicit: post-release chore commits (re-add, go.sum) do not warm [Unreleased] — write it in AGENTS.md                     | Low      | S      | Documentation |
| 45 | Apply the check-lint-config.sh guard pattern to release-script invariants (e.g. "no GOWORK=off in the govulncheck loop")                 | Low      | S      | Quality       |
| 46 | Re-verify prerender.go stays in sync the next time the wire demo grows (TODO_LIST wire #154)                                             | Low      | S      | Quality       |
| 47 | Add `ci-repro.sh --release-dry` running the dry-run mode (#8) locally before any push                                                    | Medium   | M      | Quality       |
| 48 | Sweep ROADMAP wire ideas for anything promoted by wire.Handler's adoption this release                                                   | Low      | S      | Documentation |
| 49 | Track pkg.go.dev indexing of v1.13.0 + all 6 sub-modules within 24h                                                                      | Low      | S      | Documentation |
| 50 | Retire the /tmp/tc-bin govulncheck after #1 lands; update the runbook install instructions                                               | Low      | S      | Cleanup       |

## g) Three questions I cannot answer myself

1. **The three concurrent agent sessions (`fix/alert-nonce-fallback`, `fix/live-region-nonce-guard`/-`pr`, `fix/statcard-dl-structure`/-`clean`) — are they yours and intentional, and should they be paused during release windows?** I tried: `gh pr list` (only PR #8 open), branch inspection, daemon-log inference. I cannot know intent or whether their sessions are still alive; this blocks PR #8 / 1.13.1 handling (items f#3, f#4, f#37).
2. **Release cadence: cut 1.13.1 immediately after PR #8 merges, or batch it with more work?** This is a product-ownership call (tag budget vs. fix latency); the answer decides whether I queue the cut now or harvest backlog while it batches.
3. **May I add a dedicated govulncheck flake input (the `nixpkgs-go`-style input split) to flake.nix?** AGENTS.md reserves flake-input changes for deliberate, sanctioned updates. The alternative is keeping the ad-hoc `go install` documented in the runbook — your policy call (items f#1, f#2, f#50).

---

_Format note: skill spec is a styled HTML dashboard; user explicitly requested `.md` for this report — override honored per skill rules. Section f is HARVEST fuel: it should route into TODO_LIST.md/ROADMAP.md via docs-health HARVEST, not die in this timestamped file._
