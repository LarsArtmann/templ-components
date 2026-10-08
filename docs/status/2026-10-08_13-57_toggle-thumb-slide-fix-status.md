# Status Report — Toggle Thumb-Slide Fix Session

**Date:** 2026-10-08 13:57 CEST
**Session scope:** Bug report "the apple style button does not move the circle" → diagnose, fix, verify, document. Per operator instruction, this report covers ONLY this session's run and what it noticed — no fresh repo-wide research.

---

## Headline

`forms.Toggle` (the iOS/Apple-style switch) never moved its thumb circle when checked.
Root cause: Tailwind's `peer-checked:` variant compiles to `:is(:where(.peer):checked ~ *)`
— a **sibling** selector. The thumb `<span>` is nested _inside_ the track `<span>`, so its
`peer-checked:translate-x-N` class could never match. The track still turned blue, which
masked the bug, and every golden (HTML + PNG) had captured the broken state — so every
test gate stayed green while the component was visibly broken to users.

**Fix:** translate class moved to the track via Tailwind v4's child variant
`peer-checked:*:translate-x-N` (probe-compiled to `sibling-of-checked-peer > *`).
Zero DOM change. Pixel-proven: `visualtest/testdata/toggle/light_checked.png` now shows
the thumb right of center.

---

## a) FULLY DONE

| #   | Item                                                                                                                                                                                                                                               | Evidence                                                                                                 |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| a1  | Root cause diagnosis with compiled-CSS proof (selector `:is(:where(.peer):checked ~ *)` cannot match a grandchild of the peer)                                                                                                                     | `examples/demo/static/app.css` pre-fix selectors; original `toggle/light_checked.png` showed stuck thumb |
| a2  | Fix implemented in `forms/toggle.templ`: lookup values → `peer-checked:*:translate-x-4/5/6`, class moved from thumb to track span, comments document BOTH failure modes (invisible-token + sibling-matcher)                                        | commit `e2a35054`; `toggle_templ.go` regenerated from repo root with pinned generator                    |
| a3  | Probe-verified fix selector before editing: compiled 3 candidates (`peer-checked:*:`, `group-has-checked:`, `peer-checked:after:`) with the project's own `tailwindcss_4` and picked the zero-DOM-change option                                    | probe in `/tmp/tw-probe`, v4.3.3 output inspected                                                        |
| a4  | Regression guard upgraded: `TestToggleEmitsCompletePeerCheckedClasses` now pins the new literals AND bans the bare `peer-checked:translate-` form via `AssertNotContains`                                                                          | `forms/edge_cases_test.go`                                                                               |
| a5  | HTML goldens re-baselined (only the 3 toggle goldens changed; diff verified class-on-track)                                                                                                                                                        | commit `02783f9a`                                                                                        |
| a6  | Demo CSS recompiled twice (final state: exactly ONE `translate-x-5` rule, the child-variant form; stale old selectors purged after `_sources` mirror sync)                                                                                         | `nix run .#css`; fixed-string + regex verification                                                       |
| a7  | Visual goldens re-baselined and verified: `toggle/light_checked`, `toggle/dark` (thumb move, 6.57% delta), `routes/settings_*` (moved thumbs), plus sub-threshold antialiasing drift baked on 9 route PNGs                                         | full visual suite `ok` in 212s (pixel goldens + axe sweep + dual-transport e2e)                          |
| a8  | Full test matrix green: root module exit 0; all 6 sub-modules OK (utils, icons, errorpage, charts/echarts, datastar, htmx); website module OK; visualtest compiles; `golangci-lint` 0 issues on `./forms/... ./cmd/...`                            | `/tmp/tc-root-test2.log`; per-module runs                                                                |
| a9  | CHANGELOG `[Unreleased]` → `### Fixed` entry added (kept warm per release convention)                                                                                                                                                              | commit `02783f9a`                                                                                        |
| a10 | AGENTS.md lesson extended: complete-literals bullet now covers the sibling-matcher trap, the fix pattern, and the guard test                                                                                                                       | commit `7dd3aead`                                                                                        |
| a11 | `cmd/tc/_sources` mirror synced twice (caught the second drift myself via `TestSourcesMatchPackageFiles` failing after a late comment edit — fixed, root tests re-run green)                                                                       | guard `--fix` output; root exit 0                                                                        |
| a12 | Bonus normalization: 120 generated `*_templ.go` files had `templ.Error` FileName paths encoded from the PARENT directory (`templ-components/...` prefix); my repo-root `templ generate ./...` restored the canonical form (documented drift class) | commit `e2a35054` diff inspection                                                                        |

## b) PARTIALLY DONE

| Item                          | Works now                                         | Remains open                                                                                                                                                                                                               | Blocker                                      | Effort |
| ----------------------------- | ------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------- | ------ |
| Toggle fix distribution       | Fixed, tested, committed on local master          | NOT pushed (master is 5 commits ahead of origin), NOT released, so the live demo at templcomponents.lars.software/demo still shows the stuck thumb and consumers get it only after the next tag                            | Push/release ownership undecided (see g1/g2) | S      |
| FileName-drift remediation    | Tree is normalized                                | No guard FAILS when a future regen runs from the wrong CWD — the drift class can silently return                                                                                                                           | Not started (was out of session scope)       | S      |
| Visual-golden trustworthiness | Toggle + settings goldens now show correct pixels | The systemic blind spot remains: goldens are self-consistent, so a baked bug passes forever until a human looks; full-page 0.1% threshold dilutes component-level pixel bugs (~0.03% thumb move PASSED on `/forms` routes) | Tooling decision needed (see f-items 9–14)   | M      |

## c) NOT STARTED

| Planned item                                                                                                                                     | Why not started                                                                             | Priority |
| ------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------- | -------- |
| `scripts/ci-repro.sh --lint --website` at the exact tip + push (M03 ritual)                                                                      | I did not push; but the daemon may push autonomously, so the "should be green" risk is live | Critical |
| Interaction-level regression test (chromedp: click toggle → assert computed `translate` ≠ 0)                                                     | Session scoped to fix+verify; needs a small harness decision (g3)                           | High     |
| Patch release cut (`release.sh`) for the fix                                                                                                     | Release cadence is the owner's call (g1)                                                    | High     |
| HARVEST of section (f) into `TODO_LIST.md` / `ROADMAP.md`                                                                                        | Operator said "just report, then wait" — HARVEST is the sanctioned next step                | High     |
| Dead demo template cleanup: `examples/demo/forms_section.templ` is referenced by NO route (found while explaining why `/display` goldens passed) | Out of scope; cleanup task                                                                  | Medium   |
| Static guard for variant-prefix placement (templ-AST scan: `peer-*` classes must sit on siblings of the peer)                                    | Needs design; ROADMAP fuel                                                                  | Medium   |
| `TODO_LIST` #255 strike (toggle goldens now exist AND are correct)                                                                               | HARVEST action                                                                              | Low      |

## d) TOTALLY FUCKED UP

**Nothing in the final tree state is broken.** All suites green, worktree clean, content
verified against the daemon's fragmented commits. Radical-honesty section, so the near
misses and the one systemic scandal:

1. **The systemic one (pre-existing, now fixed but instructive):** the visual golden
   `toggle/light_checked.png` had the bug BAKED IN. The pixel test did exactly what it
   was built for (TODO #255, "the August 2025 regressions were invisible to string
   tests") and still passed — because it compared a broken render against a broken
   golden. A user-visible bug sailed through HTML goldens, string tests, the a11y sweep,
   CI, and the visual suite simultaneously. Severity: was user-facing; now mitigated by
   the fix + re-baseline + HTML guard. Residual risk: the same self-consistency trap
   exists for EVERY component golden (see f9–f14).
2. **Near miss (mine):** I ran `nix develop -c templ generate ./forms/...` WHILE the full
   visual suite was running. AGENTS.md explicitly warns concurrent nix processes fight
   over the eval-cache SQLite DB. It worked. It was a gamble I should not have taken.
3. **Near miss (mine):** an over-escaped rg pattern made me print "old selectors gone"
   based on ZERO matches that were actually a broken pattern, not a verified absence.
   Caught on re-verification with fixed-string search + full rule enumeration.

## e) WHAT WE SHOULD IMPROVE

1. **Golden self-consistency is a structural blind spot.** For stateful visual components
   (toggle, slider, rating, accordion chevron), add INTERACTION-level proofs (click →
   computed style) — string guards and re-baselined PNGs both survive a baked-golden bug.
   Impact: this exact bug class already escaped twice (Aug 2025, today).
2. **Full-page route threshold dilution.** 0.1% of a 1280×N page is ~2–4k px; a moved
   toggle thumb is ~640 px and PASSED on `/forms`. Either per-region diffing, component
   goldens with tighter thresholds (they exist — 0.1% on a 175×56 crop caught it), or a
   rule: interactive-state pixels are pinned at component level, never route level.
3. **`-update` granularity is per-test-function.** Re-baselining `TestDemoRouteGoldens`
   re-baked 9 unrelated route PNGs with antialiasing drift. A `golden:-update=name`
   filter (or per-golden update mode) would keep re-baselines surgical.
4. **No guard for the FileName CWD drift.** `check-templ-sync.sh` checks imports, not
   FileName paths. One regex (`FileName: \`templ-components/`) in the existing guard
   makes the wrong-CWD regen class fail in <2s.
5. **I hit the `rg -rn` gotcha AGAIN** (it is documented in AGENTS.md, I hit it anyway,
   third+ occurrence repo-wide). Candidate fix: a shell alias/wrapper or a pre-commit
   linter for `-rn` in command history is unrealistic — but the AGENTS.md bullet could
   move higher / into the skills. Same session also produced a false-negative rg via
   over-escaping; rule: verify absence claims with a SECOND method.
6. **Piping long-running suites through `head` loses data.** The visual inventory run's
   `| head -40` cut off TestToggleGoldens + every suite after routes → extra rerun.
   Rule: tee to a file, filter afterwards.
7. **Daemon commit fragmentation mid-session (4 commits) plus mid-flight `git status`**
   made the settings-PNG investigation loop through 3 theories before evidence (git log
   dates + viewing the PNG) settled it. Lesson already half-documented; the decisive
   move is ALWAYS: evidence on disk (dates, content), never theory about timing.
8. **CI-repro ritual gap:** work landed on master (via daemon) without
   `scripts/ci-repro.sh --lint --website` ever producing a VERDICT at this tip. Local
   per-lane tests all pass, but the ritual exists precisely because "should be green"
   is not a verdict.

## f) TOP 50 NEXT TASKS (brainstorm — HARVEST fodder, not commitments)

Impact: Critical/High/Medium/Low · Effort: S <30min, M 30min–2h, L >2h

| #  | Task                                                                                                                                                                                          | Impact   | Effort | Category      |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 1  | Run `scripts/ci-repro.sh --lint --website` at current tip; only push on `VERDICT: PASS`                                                                                                       | Critical | S      | Quality       |
| 2  | Build chromedp interaction test: click toggle, assert computed `translate`/`translate-x` ≠ 0 in both themes                                                                                   | High     | M      | Quality       |
| 3  | Cut patch release shipping the Toggle fix (after g1 decision); run release inside nix shell with govulncheck on PATH                                                                          | High     | S      | Release       |
| 4  | Redeploy demo (Cloud Run) so templcomponents.lars.software/demo stops showing the stuck thumb                                                                                                 | High     | S      | Bug           |
| 5  | HARVEST this report: strike TODO #255, fold f-items into TODO_LIST/ROADMAP                                                                                                                    | High     | S      | Documentation |
| 6  | Add FileName-prefix check to `scripts/check-templ-sync.sh` (fail on `FileName: \`templ-components/`)                                                                                          | High     | S      | Quality       |
| 7  | Audit ALL `peer-*` class usages repo-wide (.templ AND .go): confirm every target is a sibling of its `.peer` (rating + hover_card verified this session; sweep the rest)                      | High     | S      | Quality       |
| 8  | Add the "variant-prefix classes: complete literal AND sibling-placed" item to `docs/plan-authoring-checklist.md`                                                                              | High     | S      | Documentation |
| 9  | Route-golden improvement: per-region or per-component pixel pinning for interactive states (kills threshold dilution)                                                                         | Medium   | L      | Quality       |
| 10 | `-update` scoping: allow updating a single golden file, not the whole test function                                                                                                           | Medium   | M      | Quality       |
| 11 | Investigate the settings-route 2576-px transient (scroll-reveal freeze?): pin `[data-animate]` completion before route captures                                                               | Medium   | M      | Quality       |
| 12 | Visual toggle goldens for SM + LG sizes (only MD pinned today; translate distances differ per size)                                                                                           | Medium   | S      | Quality       |
| 13 | Toggle RTL visual golden (`start-0.5` flips; thumb must move toward inline-start end)                                                                                                         | Medium   | S      | Quality       |
| 14 | Toggle focus-visible visual golden (keyboard state unpinned)                                                                                                                                  | Medium   | S      | Quality       |
| 15 | Document the `peer-checked:*:` child-variant pattern in `docs/` (JavaScript-guide-adjacent or a small CSS-patterns doc)                                                                       | Medium   | S      | Documentation |
| 16 | Mirror the sibling-matcher lesson into the templ-components SKILL.md (in-repo source, fan-out delivers)                                                                                       | Medium   | S      | Documentation |
| 17 | Delete or wire `examples/demo/forms_section.templ` (dead template, contains the only checked toggle on a non-route)                                                                           | Medium   | S      | Cleanup       |
| 18 | Audit `examples/demo` for other unwired templ funcs                                                                                                                                           | Medium   | S      | Cleanup       |
| 19 | Static analysis: templ-AST scanner failing when a `peer-*` class lands on an element nested inside a sibling-of-peer (ROADMAP fuel)                                                           | Medium   | L      | Quality       |
| 20 | Antialiasing drift policy for route PNGs: perceptual-hash diffing or scheduled sanctioned re-baseline lane                                                                                    | Medium   | M      | Quality       |
| 21 | visualtest freshness canary: assert the served `/app.css` hash equals the compiled artifact at harness start                                                                                  | Medium   | S      | Quality       |
| 22 | Placement guard for Rating's label-sibling invariant (same class as Toggle; today only a comment protects it)                                                                                 | Medium   | S      | Quality       |
| 23 | Sweep `.go` files for dynamically composed variant prefixes (`"peer-checked:" + x` style) — the complete-literal rule is currently enforced only by tests/convention                          | Medium   | S      | Quality       |
| 24 | Daemon hygiene: 4 mid-session commits split one logical change across e2a35054/02783f9a/7dd3aead/6d237b8d — evaluate squash-before-push or BuildFlow upstream change                          | Medium   | M      | Cleanup       |
| 25 | Record the "git status mid-flight lies during visual -update runs" gotcha in AGENTS.md daemon section                                                                                         | Low      | S      | Documentation |
| 26 | Find the commit that introduced the FileName CWD drift (`git log -S 'FileName: \`templ-components/`) — bounds how long the invariant was broken                                               | Low      | S      | Cleanup       |
| 27 | Pre-release ritual: re-verify `[Unreleased]` warmth + `git cat-file -p HEAD \| grep -c parent` (daemon 6th/7th class checks) before any tag                                                   | High     | S      | Release       |
| 28 | After release: verify pkg.go.dev serves the fixed source and the sub-module tag set (check-release-tags.sh)                                                                                   | Medium   | S      | Release       |
| 29 | Consider an issue in this repo for the Toggle bug (traceability: CHANGELOG entry ↔ issue link)                                                                                                | Low      | S      | Documentation |
| 30 | Evaluate `group-has-checked:` as documented alternative for consumers pinning older Tailwind (we standardize on v4, but the doc note is cheap)                                                | Low      | S      | Documentation |
| 31 | Golden sweep: verify every OTHER stateful-variant component's checked/active golden actually shows the ACTIVE state (audit for more baked-in bugs: Slider fill, Rating stars, Accordion open) | High     | M      | Quality       |
| 32 | Add computed-style assertions to the existing wire/forms e2e harness for transition components (reuse newFlowTab patterns)                                                                    | Medium   | M      | Quality       |
| 33 | AGENTS.md: move the `rg -rn` bullet into a "tool misuse" cluster at the top of CI gotchas (I hit it AGAIN this session)                                                                       | Low      | S      | Documentation |
| 34 | Session tooling rule → skill: "verify absence with a second method" (over-escaped rg false-negative) into references/lessons.md via crush-config commit                                       | Low      | S      | Documentation |
| 35 | Visual suite runtime: 212s full run; consider `-run` lane split in CI (pixel vs e2e) for faster feedback                                                                                      | Low      | M      | Quality       |
| 36 | Confirm `TestCompiledCSSInventory` covers the `_sources` copies' CSS implications (mirror sync changed scan inputs this session)                                                              | Low      | S      | Quality       |
| 37 | CHANGELOG entry could link the pixel proof (before/after PNG paths) for future archaeology                                                                                                    | Low      | S      | Documentation |
| 38 | Pre-existing backlog, adjacent: TODO #240 — migrate ~48 raw chromedp `Poll` sites to the poll helpers (this session's e2e reads relied on them)                                               | Medium   | L      | Quality       |
| 39 | Pre-existing: TODO #292 — art-dupl piped-capture truncation upstream bug                                                                                                                      | Medium   | M      | Bug           |
| 40 | Pre-existing: TODO #335 — evaluate templ v0.3.1070 (generator output change; needs re-baseline plan like this session's, but repo-wide)                                                       | High     | L      | Feature       |
| 41 | Pre-existing: TODO #351/#352 — go-licenses GOROOT wrapper durability (BuildFlow reinstall re-breaks it)                                                                                       | Medium   | S      | Bug           |
| 42 | This session's probe pattern ("compile candidate selectors before editing") is reusable → document in the templ-components skill as the standard move for variant work                        | Medium   | S      | Documentation |
| 43 | Demo forms page: toggle demo is UNCHECKED-only; add a checked + disabled demo pair so the live page shows all four states                                                                     | Low      | S      | Feature       |
| 44 | Verify no OTHER package CSS artifacts (`templates/styles.css`, `templ-components-theme.out.css`) contain the dead sibling-form selectors (release.sh refreshes them; check timing)            | Low      | S      | Cleanup       |
| 45 | Add `AssertContainsAll`-based class-co-location helper to utils (assert two classes share one attribute) — makes placement guards one-liners                                                  | Medium   | S      | Quality       |
| 46 | Post-mortem note for docs/agent-context-history.md: the settings-PNG daemon race (evidence-first debugging saved it; record the signature)                                                    | Low      | S      | Documentation |
| 47 | Check whether the visual harness should hard-fail when `head`-style truncation is detected in CI logs (ties to AGENTS.md "final summary line missing" rule)                                   | Low      | S      | Quality       |
| 48 | Idea: golden files embed the component's `utils.Version` as a comment — makes stale-golden detection trivial (ROADMAP fuel)                                                                   | Low      | S      | Quality       |
| 49 | Idea: demo page badge showing the library version the live demo runs (would have made the "live demo still broken" state visible)                                                             | Low      | M      | Feature       |
| 50 | Review whether other Tailwind v4 child/descendant variants (`*`, `**`) should replace any remaining nested-DOM class patterns repo-wide                                                       | Low      | M      | Cleanup       |

## g) TOP 3 QUESTIONS I CANNOT ANSWER MYSELF

1. **Release cadence:** Cut a patch release for the Toggle fix NOW, or batch it with the
   other warm `[Unreleased]` entries? (Tried: CHANGELOG/ROADMAP give no cadence policy;
   the decision is yours. Determines whether I run the release ritual + item 4's redeploy
   immediately.)
2. **Push ownership:** Master is 5 commits ahead of origin (fix included), and AGENTS.md
   documents the daemon pushing master autonomously. Do you want me to run the ci-repro
   ritual at the tip and push explicitly, or is pushing exclusively daemon/your territory?
   (Tried: no doc states who owns push; I must not push unasked.)
3. **Interaction-test investment:** Should I build the permanent chromedp interaction
   harness (click → computed style, per f2 — the only test shape that would have caught
   this bug despite baked goldens), or do you consider re-baselined PNGs + the HTML guard
   sufficient? (Cost: ~half a day + a new test pattern; benefit: kills the golden
   self-consistency blind spot for stateful components.)

---

_Point-in-time snapshot. Section (f) is HARVEST input for `TODO_LIST.md`/`ROADMAP.md`
(docs-health → HARVEST) — do not let it die in this timestamped file._
