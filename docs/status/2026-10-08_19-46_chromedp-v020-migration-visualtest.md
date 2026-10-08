# Status Report — 2026-10-08 19:46 CEST — chromedp v0.20 migration & visualtest verification

Session scope: break the BuildFlow `go-mod-update [visualtest]` failure loop (6 consecutive
identical failures: chromedp v0.20.1 bump broke the visualtest module compile), migrate the
module to the new chromedp/cdproto APIs, verify at runtime.

## a) Fully done

1. **Root cause fixed: visualtest migrated to chromedp v0.20.1 + cdproto v0.157.9.** ~40 files:
   `Run(ctx, a, b, …)` → `Do(ctx, …)`, typed `Evaluate[T]`/`Poll[T]`, `ActionFunc` → `Func`,
   `Action` → `Action[chromedp.Void]`, `Tasks` → `Steps`/typed slices, `ByQuery`/`ByQueryAll`
   → `chromedp.CSS`/`CSSAll` selectors, `FullScreenshot`/`CaptureScreenshot`/`Screenshot`/
   `Location`/`InnerHTML` restructured to value-returning actions, old cdproto builder+`.Do()`
   commands → generic `cdp.Call`/`chromedp.CallBrowser` (mouse events, touch emulation,
   clipboard permissions, focus emulation), and `tools/shots`' `ListenTarget` + `network.Enable`
   docStatus handshake replaced by `chromedp.RunResponse`.
2. **New helper file `visualtest/eval.go`** (`evalInto`/`evalExprInto`/`evalVoid`) — the adapters
   that let evaluation steps stay inside single `Do` chains; **`poll.go` rewritten onto typed
   `Poll[T]` with unchanged public helper signatures** (the pollBool/pollTrue/pollText contract
   from the TODO #193 lesson survives).
3. **templ v0.3.1020 pin restored across all 9 modules** — the daemon re-applied the forbidden
   v0.3.1070 bump for the third time today (commit ee2f10bd side-effect + earlier 02f1863f);
   my `go get` for chromedp also dragged templ up. Leaf-first downgrade + tidy landed all 9 at
   v0.3.1020; `scripts/check-replace-directives.sh` OK; `utils.TestTemplVersionPin` OK.
4. **All static gates green on the migrated module**: `go build ./...` OK, `go vet ./...` OK
   (compiles all test files), `golangci-lint run ./...` = **0 issues** (after exhaustruct nolints
   with justification, a funlen split `captureModeShots`, ST1023, unconvert, nonamedreturns
   fixes), gofmt clean, `go mod tidy` clean.
5. **Run-#1 failure forensics (the three golden mismatches) fully root-caused** — none were the
   migration (details under "Totally fucked up" + evidence): image-diff inspection, a 4-way
   Chromium paint probe (old/new CSS × old/new class — all combinations paint as expected),
   and an A/B run in a detached temp worktree at `ec37a3fe` (old tree passes).

## b) Partially done

1. **Runtime verification of the migrated module.** The suite RUNS end-to-end on chromedp
   v0.20 (218 s + 233 s full runs, real headless Chromium; the e2e pack and the axe a11y
   sweep passed in run #1). Not yet green overall:
   - Run #1 component failures (toggle ×4, carousel ×2, notfound404 ×2) were a **mid-run
     daemon race**: my test binary built 19:23 with the parallel session's renamed templates
     (committed 19:26) but the pre-recompile demo CSS (recompiled 19:28). Run #2 confirmed:
     **all three component-level failures now pass.**
   - Run #2 still fails on **4 demo route goldens: kanban (light/dark), recipes (light/dark),
     users (light/dark)** — demo pages the parallel session is still actively churning
     (`appshell.templ` modified mid-run, `docs/agent-context-history.md` being edited right
     now). Their `.actual.png` artifacts sit in `visualtest/testdata/.fail/routes/`.
2. **BuildFlow loop break.** The code-level root cause is fixed, but the failing step
   (`buildflow -s "go-mod-update [visualtest]"`) has not been re-run to flip the 6-run loop
   to PASS. Deliberately sequenced after a green visual suite.
3. **Memory/lesson recording.** All analysis done (chromedp v0.20 API map, the daemon's
   split-commit race pattern, the gofmt -r silent no-op trap), nothing written yet to
   AGENTS.md / CHANGELOG `[Unreleased]` / TODO_LIST.

## c) Not started

- CHANGELOG `[Unreleased]` entry for the migration (repo rule: warm at all times).
- AGENTS.md updates: chromedp v0.20 migration notes in the visualtest sections.
- TODO_LIST entries: axe-sweep modernization to `EvalAwaitPromise` (deletes the
  `window.__tcAxeJSON` polling handshake); note that TODO #240's raw-Poll migration is now
  cheaper with typed `Poll[T]`.
- `scripts/ci-repro.sh --lint --website` ritual (required before any push).
- Docs sweep for v0.16-era chromedp API mentions (`docs/visual-testing.md` etc.).
- Pin-policy decision: visualtest `go.mod` says "track latest upstream" but BuildFlow does not
  respect pins — encode the resulting contract (migrate-on-break, or teach BuildFlow rejects).
- LSP restart (gopls served stale compiler diagnostics the entire session).

## d) Totally fucked up (own failures this session)

1. **Run #1 raced the daemon's split commits.** I started a 4-minute full visual suite while
   the tree was actively churning (daemon commits at 19:11, 19:14, 19:19, 19:21, 19:23,
   19:26, 19:28, 19:30). The test binary straddled the template/CSS split. Cost: one burned
   suite run + a long forensics detour. The tree's churn was visible in `git log` before I
   started — I didn't look.
2. **First `gofmt -r` pass silently did nothing.** Two stacked causes: (a) my loop built
   patterns with spaces instead of commas (`a1 a2`), (b) gofmt -r does not match
   qualified-name patterns with wildcards (`chromedp.Run(c, a1, a2)`) at all — it exited 0
   both times. I only caught it when the compiler surfaced the "unmigrated" sites. Silent
   no-op tooling failures must be spot-verified on first application.
3. **Destroyed run-#2's own failure detail.** I piped the re-run through `| tail -15`, so
   when it failed I could not say WHAT failed without re-running — the exact
   "never trust a capture whose final summary is missing" lesson from AGENTS (TODO #292),
   violated hours after re-reading it. (Recovered the fact via `testdata/.fail/` artifacts.)
4. **Kept reading stale gopls diagnostics all session** despite AGENTS explicitly warning
   that the LSP misleads after edits ("restart it; build is ground truth"). Cost: repeated
   noise in every tool result.

## e) What we should improve (systemic)

- **Long verification runs need a quiescent tree.** Either pause the daemon (BuildFlow has no
  pause?) or add a pre-flight check to the visual workflow: refuse to start if HEAD is
  younger than N minutes / working tree dirty; fail loud if HEAD moves mid-run.
- **The daemon's template commit and CSS recompile land as separate commits** (19:26 vs
  19:28). Every consumer building in that window gets the invisible-numeral class of breakage.
  BuildFlow's tailwind-build should be atomic with templ changes, or the demo CSS freshness
  guard should run at commit time, not just in CI.
- **Parallel sessions need a coordination signal** (a lock file, or a convention: announce
  "tree churning" in a known place). Two sessions + daemon on one worktree produced three
  overlapping races today.
- gofmt -r: add a repo note that qualified-name rewrite patterns silently never match; use
  helper-function rewrite targets or plain codemods instead.

## f) Next tasks (prioritized, ≤50)

**Migration closure**

1. Diagnose run-#2's route-golden failures (kanban/recipes/users) once the parallel session's
   demo churn settles; pipe full output to a file this time.
2. If those routes are genuine parallel-session redesigns: re-baseline via
   `nix run .#visual -- -update -run TestRoutes…` after eyeballing the `.actual.png`s.
3. Re-run the full visual suite to green on a quiescent tree.
4. Re-run `buildflow -s "go-mod-update [visualtest]"` → confirm the 6-run loop is broken.
5. Full `buildflow --fix --build-mode=full` → exit 0.
6. Warm `CHANGELOG.md` `[Unreleased]`: chromedp v0.20.1/cdproto v0.157.9 migration of the
   visual-test infrastructure (consumer-irrelevant, dev-dep only).
7. Update AGENTS.md: chromedp v0.20 API map (Run/Do/Steps/Evaluate[T]/Events/RunResponse/
   cdp.Call/CallBrowser/Void), eval.go adapter rationale, gofmt -r gotcha.
8. TODO_LIST: axe sweep → `EvalAwaitPromise` modernization; link typed-Poll to #240.
9. Restart the LSP.
10. Spot-verify daemon commits didn't mangle migrated files (AGENTS post-commit check).
11. Decide + document the chromedp update policy (track-latest vs pin) in visualtest/go.mod
    comment and AGENTS.md; optionally teach/ask BuildFlow for dep excludes.
12. Sweep `docs/` for v0.16-era chromedp snippets (visual-testing.md, agent guides).
13. Run `scripts/ci-repro.sh --lint --website` at the exact pre-push commit (ritual).
14. Push/PR only on explicit go-ahead.

**Daemons/guardrails**
15. Investigate which BuildFlow step keeps re-applying the templ v0.3.1070 bump (3rd time
today); gate it at the source, not just via TestTemplVersionPin post-hoc.
16. Make demo-CSS recompile atomic with template changes (or commit-time freshness guard).
17. Visual-suite pre-flight: refuse dirty/churning trees; detect mid-run HEAD movement.
18. Always tee long test runs to a file; enforce the "final summary line present" rule.
19. Re-verify pseudo-version hygiene after today's go.mod surgery (BuildFlow preflight
flagged 6 drifted requires; check they are still canonical).

**Parallel session's changes (respect, don't revert)**
20. Verify the toggle thumb fix's visual goldens land in ITS session (or offer my re-baseline).
21. Same for the carousel extra slide (check HTML goldens updated in same edit — count rule).
22. Confirm `inset-s-0` is an intended Tailwind v4 utility (compiled fine; verify intent —
v4 idiom is `start-0`/`end-*`).
23. Confirm notfound404 passes now (probe says the consistent tree paints; run #2 suggests yes).
24. Review the renamed classes for the demo brand overrides (`--color-blue-500` indigo) —
axe color-contrast ledger entries still valid?

**Repo hygiene surfaced by BuildFlow preflight today**
25. AGENTS.md 377 lines > 220 budget — split content into docs/ per the buildflow hint.
26. lychee.toml: add `exclude_path = ['docs/feedback/archived']` (+ prune the 30 dead
file:// links it flags).
27. shellcheck: fix the 13 findings in scripts (SC2046 word splitting first).
28. Add dprint/prettier/hadolint/lychee/shellcheck to the devShell (removes the nix-run
WITHOUT project deps warnings).
29. interrogate: install or add to BuildFlow config as not-applicable.
30. Rebuild/reinstall BuildFlow (binary at ec8d2d3, repo HEAD moved past it).
31. vulnix: 56 CVEs in the nix closure (binutils, apr-util, avahi, ada…) — evaluate a
nixpkgs bump; likely upstream-side, but record the decision.
32. golangci-lint-auto-configure nags the 3 deliberately disabled linters every run —
document in .buildflow.yml or file upstream.
33. `.fail/` artifacts: add to .gitignore check (they are transient; make sure never staged).
34. Re-check `TestCompiledCSSInventory`/`TestCSSFreshness` after today's daemon recompiles.
35. Run root-module test suite (go test ./...) — untouched by my work but the demo CSS and
templates churned; cheap insurance.
36. `scripts/check-templ-sync.sh` after the daemon's templ regen bursts.

**Migration quality follow-ups (non-blocking)**
37. Migrate remaining ~48 raw `chromedp.Poll` call sites (TODO #240) onto typed `Poll[T]` —
now mechanical via evalExprInto-style helpers.
38. Consider `RunResponse` where e2e tests hand-roll navigation-response checks.
39. Re-check the kanban e2e timing-flake classes under the new action execution model.
40. ogshot quality constant: confirm 95 (JPEG-in-.png gotcha) is still intended.
41. Evaluate deleting `evalVoid` in favor of direct `chromedp.Evaluate[chromedp.Void]` once
call sites stop churning (one indirection fewer).
42. Update `docs/agent-context-history.md` pointer for the 2026-10-08 daemon race incident
(its entry is being written by the other session — coordinate).
43. After green: run `nix run .#verify` once end-to-end.
44. Revisit `.buildflow.yml` skip list: `go-auto-upgrade` skip note mentions "when the step
learns to respect pinned rejects" — same applies to go-mod-update now (loop evidence).
45. Confirm the A/B temp worktree is cleaned up (`git worktree remove /tmp/tc-ab-404`).
46. Clean `/tmp/tc-404-probe` artifacts (not in repo, but tidy).
47. Check `TestDemoKanbanHTTPContracts` & friends still green post CSS churn.
48. Consider a `visualtest/README` or doc.go note about the v0.20 API conventions used.
49. Re-check docs-count drift guards if the carousel slide changed demo counts.
50. Only after ALL gates green + user approval: push (never before).

## g) Questions I cannot resolve myself (max 3)

1. **Ownership of the parallel changes:** the toggle thumb fix, carousel slide, notfound404
   class rename, and the in-flight recipes/users/kanban demo edits belong to another session
   (or you). Should I re-baseline their pixel goldens as part of THIS verification, or does
   that session own golden updates in its own commits? I can coordinate either way; I chose
   not to touch their goldens without an answer.
2. **Run #2's route failures:** with the tree still churning as I type
   (`docs/agent-context-history.md` modified right now), do you want me to keep retrying the
   visual suite until the other session goes quiet, or is there a way to pause the
   daemon/other session so I can get one clean run?
3. **chromedp update policy:** visualtest/go.mod says "track latest upstream", but BuildFlow
   treats any breakage as a 6-run failure loop and does not respect pins. Is
   "migrate-on-break within the same day" the standing policy you want encoded, or should I
   add a BuildFlow exclusion for visualtest deps until you say otherwise?

## Self-review (asked explicitly)

- **What did I forget?** To look at `git log` freshness before launching the long suite; to
  warm the CHANGELOG the moment the migration compiled (repo rule); to restart the LSP; to
  preserve full run output; and to announce/coordinate that I was about to occupy the tree
  for ~4 minutes of browser tests.
- **What could I have done better?** Freeze a commit (detached worktree) for run #1 the way
  I eventually did for the A/B test — one worktree run from the start would have saved the
  entire forensics detour. Also: verify the first application of any generated rewrite rule
  (the gofmt -r silence), and treat "suite green" as blocked-not-passing while the tree is
  churning, communicating that earlier instead of pushing through two runs.
- **What could I still improve?** The migration itself is compile+lint proven and
  e2e/axe-proven, but not full-suite-green-proven — that is the one remaining acceptance
  gate, and it is blocked on tree quiescence more than on code. Memory writes (AGENTS/
  CHANGELOG/TODO_LIST) are drafted in my head but not on disk — that is the next action
  after your go-ahead, alongside task list items 1–5.
