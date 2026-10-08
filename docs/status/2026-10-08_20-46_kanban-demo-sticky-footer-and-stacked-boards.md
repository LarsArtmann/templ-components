# Status Report — 2026-10-08 20:46 CEST

## Kanban demo UX fixes: sticky footer at every viewport + one full-width row per transport board

**Session scope:** `https://templcomponents.lars.software/demo/kanban` complaints — (1) footer floating mid-screen, (2) HTMX and Datastar boards crammed side by side.

**Tree state at report time:** `master`, **3 unpushed commits** (daemon-authored snapshots of this session's work; `CHANGELOG.md` edit still unstaged-staged `M`). Nothing pushed by me; CI on origin/master does NOT yet contain this work.

---

## a) FULLY DONE

| Work | Evidence |
| --- | --- |
| **Root cause found for footer bug (two layers).** The demo nested its footer inside AppShell's *Content* slot, so the `mt-auto` pin never applied; independently, AppShell's `Footer` slot only stuck at the grid breakpoint (below `lg:` the shell was a plain block, content column never stretched). | Code read of `examples/demo/shell.templ`, `layout/appshell.templ`, live HTML fetch |
| **AppShell library fix:** shell wrapper is now `flex flex-col min-h-dvh` (all four breakpoint entries + the no-sidebar class) and the content column gained `flex-1` — footer pins at mobile too; desktop grid rendering byte-unchanged (flex properties are inert on grid items). | `layout/appshell.templ`, `layout/appshell_types.go` godoc |
| **Demo fix:** footer moved out of `demoShellContent` into `AppShellProps.Footer`. | `examples/demo/shell.templ` |
| **Kanban boards stacked:** `xl:grid-cols-2` removed → one full-width row per transport, all four columns visible at desktop, clearly labeled HTMX/Datastar rows. | `examples/demo/kanban_demo.templ` |
| **Regression pins:** two new AppShell test cases (mobile flex-column classes + footer-stick chain `flex flex-1 min-w-0 flex-col` + `<footer class="mt-auto">`). | `layout/appshell_test.go` |
| **Full verify:** `nix run .#verify` = generate (zero diff) + build + race tests + lint **0 issues on all 7 modules**. Layout/utils/demo suites green. | verify output |
| **Full visual suite green** (goldens, axe a11y sweeps, all e2e flows) — 321s, `-parallel 4`. | `ok visualtest 321.539s` |
| **5 route goldens legitimately re-baselined** (`forms_{light,dark,rtl}`, `users_{light,dark}`) — short pages where the footer moved to the bottom. | `git status visualtest/testdata/routes/` |
| **Browser-proven at the user's real viewport** (the step goldens could not do): temp probe test captured `/kanban` + `/users` at 1920×1080 and 375×667 — footer at the very bottom on a short page, boards on separate labeled rows, mobile fold correct. Probe deleted after review. | viewed PNGs during session |
| **Hygiene:** templ regenerated with the pinned generator from repo root; demo CSS recompiled (`nix run .#css`); `cmd/tc/_sources` mirror synced in both files; stale `forms/toggle_templ.go` generated-comment drift healed by the regeneration; `TestTemplGeneratedInSync` implicitly green. | git log, mirror script output |
| **CHANGELOG `[Unreleased]` warmed** with both fixes (library footer fix + kanban layout fix). **AGENTS.md gotcha recorded** (route goldens render just below `xl:`; demo footer must ride the Footer slot). `TestDocsCountDrift` green after. | CHANGELOG.md, AGENTS.md |
| **Daemon-race ritual:** re-verified all session edits survived the daemon's mid-session commits (5 commits observed); tree now clean except the report + staged CHANGELOG. | grep counts on every touched file |

## b) PARTIALLY DONE

1. **Visual proof below `lg:` for a genuinely SHORT mobile page.** My 375px probe was `/users`, whose content nearly fills the fold — the footer looked right, but it wasn't the clean mt-auto proof (a near-empty AppShell page at phone width would be). The library fix is logically airtight and goldens are green; the airtight *visual* proof for that exact case is missing.
2. **"Better way to demonstrate HTMX vs Datastar."** Stacking satisfies the literal request and reads far better, but the comparison is still only textual labels; the 4th column ("Done") is still partially cut at 1920 (board content 1200px vs ~1088px container) — I silently judged internal scroll acceptable (idiomatic kanban) without surfacing the call.
3. **Push/deploy status.** Work is committed locally (daemon) but NOT pushed; the live site still runs the old layout. The pre-push ritual (`scripts/ci-repro.sh --lint --website`) has NOT been run — `nix run .#verify` covered lint/build/test but not the website CI lane.
4. **xl:-blindness of the golden suite: documented, not fixed.** AGENTS.md now records that `ViewportDesktop` (1280) renders ~1265px CSS-wide (scrollbar), so `xl:`/`2xl:`-gated layouts are invisible to every route golden. The structural fix (a wide-viewport golden tier or scrollbar-stable capture width) is not built.
5. **TODO_LIST.md not updated.** Follow-up items from this session live in this report only; the repo convention is actionable items belong in TODO_LIST.md.

## c) NOT STARTED

- CI status check for the eventual push; deployment of the demo image (Cloud Run) so the live page actually shows the fix.
- Any permanent tooling/guard from this session's lessons (see f).
- Release planning: these fixes ride the next version's `[Unreleased]`; no version cut touched (correctly).

## d) TOTALLY FUCKED UP

Nothing shipped broken — but one **near-miss worth naming honestly**: I initially trusted the golden re-baseline as the visual verdict. The kanban goldens came back **byte-identical** after an intentional layout change, and had I stopped there I would have declared the stacking change "verified" without ever having SEEN it — the golden tier is structurally blind to `xl:` layouts. Only the custom-viewport probe caught reality. Also minor: my first AppShell test assertion used a non-adjacent substring (`"flex flex-col min-h-dvh"`) that tailwind-merge's class reordering broke — caught by the test itself in seconds.

## e) WHAT WE SHOULD IMPROVE

1. **Trust the cheapest oracle last, not first.** A byte-identical re-baseline of an intentionally-changed layout is a red flag, not a pass. Make "intentional change + unchanged golden = investigate" a reflex.
2. **The golden suite needs a viewport it can't lie about.** One wide tier (or scrollbar-stabilized capture) closes the entire `xl:` blind spot permanently instead of per-incident probes.
3. **Structural guard for the footer-slot rule.** The 2026-10-01 nonce outage class and this footer bug share a shape: demo content nesting broke a slot contract no test watches. A cheap served-HTML assertion beats another incident post-mortem.
4. **Silent judgment calls should be surfaced.** The "4th column scrolls at 1920" decision was made silently; design tradeoffs with visible consequences belong in the PR/report, not in my head.

## f) NEXT — up to 50 things to get done (ordered, ~impact first)

1. Run `scripts/ci-repro.sh --lint --website` at the exact tip, then push the 3 daemon commits (push ritual).
2. Verify CI green on master after push (ci.yaml + website.yml).
3. Redeploy/verify the demo (Cloud Run) and eyeball `/demo/kanban` live — footer bottom, boards stacked (the user's original page).
4. Add a wide-viewport (1920×1080) route-golden tier to `visualtest` (or make captures scrollbar-stable) — kills the `xl:` blind spot class.
5. Add a served-HTML structural test: demo pages render the footer as AppShell's `<footer class="mt-auto">` child of the content column, never inside the content container.
6. Clean mobile proof: render a near-empty AppShell page at 375×667 and pin footer-at-bottom in a test.
7. Decide the 4-column overflow at wide screens: accept internal scroll vs narrower demo columns vs wider container — make it an explicit decision, not a silent one.
8. Extend AppShell breakpoint tests: assert the flex-column classes also for `Breakpoint: MD` and `XL` variants (currently only default).
9. Audit other demo dual-transport sections for the same side-by-side cramming (calendar "both dialects side by side" is documented in old status docs).
10. Write `docs/recipes/sticky-footer.md` — the consumer recipe for AppShell footer pinning (body → main → shell → mt-auto chain).
11. Update `docs/visual-testing.md` with the scrollbar-sub-viewport gotcha (AGENTS.md has it; the public testing doc should too).
12. Consider a transport badge/dialect chip (`hx-post` vs `@post`) in each kanban board header for at-a-glance comparison.
13. Drag-and-drop smoke e2e against the actual demo page (current kanban e2e drives its own harness page, not `/demo/kanban`).
14. Dark-mode wide-viewport captures for key pages (goldens are 1280 light/dark + mobile; no wide-dark human review happened).
15. RTL golden at wide viewport (stacked full-width boards under `dir="rtl"`).
16. Regenerate the stale `website/public/og/home.png` ogshot-style (pre-existing, flagged in AGENTS.md).
17. Delete the stale untracked `examples/demo/demo` binary from disk (hygiene; rebuildable).
18. Check whether website docs (`content/docs/`) describe AppShell's Footer slot; update if the behavior wording is stale.
19. Check FEATURES.md AppShell wording for footer behavior; amend if it implies content-slot placement.
20. Triage items 1–19 + below into `TODO_LIST.md` (repo convention) and strike them as they land.
21. Guard-test the demo footer rule at the HTML-validation layer too (vnu run over goldens already green; add an explicit `<footer>`-placement rule if cheap).
22. Add `TestNoFooterInsideDemoContent`-style sweep across ALL demo pages, not just kanban.
23. Evaluate `scrollbar-gutter: stable` in the demo shell so 1280-wide sessions don't flip across the `xl:` breakpoint live (user-facing consistency, not just captures).
24. Re-check axe sweep viewport policy — confirm the a11y gate runs at the width where layouts actually diverge; consider a wide-viewport axe pass.
25. Consider extracting the probe pattern into a sanctioned `visualtest/tools/probe` (per skill rule: new capture tools must use `internal/browser`).
26. Release planning: fold both fixes into the next cut; remember `release.sh` steps and the `[Unreleased]` warmth rule (already warm).
27. After next release, verify updated godoc renders on pkg.go.dev (Footer slot wording).
28. AppShell godoc: add the sticky-footer example snippet to the component doc (currently only in types + CHANGELOG).
29. Sanity-check embedded-shell consumers (`min-h-0` override path from AppShellProps docs) against the new flex-1 content column — covered by suite, but one explicit test would pin it.
30. Review whether `min-h-dvh` triplication (body class + `<main>` flex-1 + shell) deserves a documented ownership note in `layout/base.templ`.
31. Kanban demo: show the pending/failed states link (ADR-0041) near the boards — the 800ms delay is documented in prose only.
32. Consider `KanbanBoard` container-query audit (ADR-0018 gate) — stacked full-width boards are width-hungry; rejected candidates list says don't expand casually, but kanban was never evaluated.
33. File/track the "golden-blind-to-xl" lesson as a possible upstream chromedp capture improvement (overlay scrollbars in EmulateViewport).
34. Re-run `nix run .#shots` after deploy for a fresh full-route human review set.
35. Verify `tc add layout` scaffolder still round-trips after the mirror sync (smoke the CLI once).
36. Sweep for other daemon-era drift: `templ generate` healed `toggle_templ.go` this session; run `check-templ-sync.sh` fresh at next touch to catch siblings.
37. Keep an eye on the 3 daemon snapshot commits' messages — they are generic heuristic snapshots; consider one squash/amend before push if history tidiness matters (note: daemon may push first).
38. Add the "intentional change + byte-identical golden = investigate" rule to the visual-testing docs checklist.
39. Consider bumping `TestAppShellBreakpointRendering` to pin the full class string per breakpoint (tighter than substrings).
40. Evaluate whether `demoShellFooter` should also carry a "Built with" line linking the kanban/wire pages (tiny nav value; optional).
41. Run the fuzz/bench suites once before the next release (unchanged this session, cheap ritual).
42. Confirm `TestDemoKanbanHTTPContracts` parity table still mirrors after any future kanban_demo.go edits (unchanged this session — no action now, standing rule).
43. Review `visualtest/testdata/.fail/` cleanup state (ensure no stale failure artifacts from this session's runs got staged).
44. Consider caching the demo binary build in the visual harness keyed by source hash (suite builds it 2×+ per run).
45. Check CI's visual-regression job passed with the 5 re-baselined goldens (CI runs them on push; local suite green is strong but CI is the gate).
46. Look at `/users` page's large empty region above the pinned footer at 1920 (aesthetic acceptability of flex-grown whitespace) — same class as item 7.
47. Consider documenting the probe-then-delete workflow in `docs/visual-testing.md` so future sessions don't commit scratch tests (one almost got daemon-committed this session).
48. Sweep other `xl:grid-cols-2` uses in the repo for the same cramming pattern (`grep -rn "xl:grid-cols-2"`).
49. AppShell: evaluate emitting `aria-label` on the `<footer>` when AriaLabel set (minor a11y polish, consistent with BaseProps propagation rule).
50. Post-deploy: fetch the live `/demo/kanban` HTML and assert `flex flex-col min-h-dvh` + stacked grid made it through the Firebase/Cloud Run path (guards against stale image deploys).

## g) Questions I cannot answer myself

1. **Deployment:** the fix is committed locally but unpushed, and the live page still shows the bug. Do you want me to run the pre-push ritual and push now (and is the demo redeployed automatically by the website workflow on push, or do you cut it manually with releases)?
2. **Golden tier:** should I build the permanent 1920×1080 golden lane (closes the `xl:` blind spot, costs longer CI + more flake surface), or keep 1280 goldens + manual probes as the policy?
3. **Kanban columns at wide screens:** is the slightly cut 4th column with internal scroll acceptable as idiomatic kanban, or do you want all four columns fully visible at 1920 (narrower columns / three-column demo board / wider container)?
