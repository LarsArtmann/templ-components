# Status Report — 2026-09-09 16:19 — Post-Release (v1.16.0), Axe Sweep, Demo E2E In Flight

**Scope:** execution of the N1–N21 quality-tier plan (`docs/planning/2026-09-09_04-24_quality-tier-and-release-plan.md`),
session start at `e2125fc`, currently at `7bbdef1` (+23 unpushed commits, see c-1).
**Format:** `.md` per your explicit instruction (skill default is HTML — fourth override, flagged per skill policy).

---

## a) FULLY DONE

1. ~~**N1 — RELEASE v1.16.0: CUT, TAGGED, PUSHED, PROXY-VERIFIED.**~~ done at `0ac2228`
   ~~- Pre-cut blocker found and fixed: CI's CSS Freshness job red at tip (daemon had un-minified `examples/demo/static/app.css`, +5058 lines) and Website job red (daemon flipped website pins to `astro ^7.3.2 / html-validate ^11.15.0 / typescript ^7.0.2` vs lockfile `7.3.1 / 11.12.0 / 6.0.3`). Recompiled CSS, restored pins, frozen-lockfile install verified, pushed as `6f91b63`, CI + Website green again.~~
   ~~- Full pre-release matrix green: `nix run .#verify` (generate+build+test+lint, 0 issues), per-module loop (7 sub-modules), full visual suite (117 goldens), 4 fast guards, tag-set check, CSS byte-stability, TS pin.~~
   ~~- Cut via `nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh 1.16.0 "display.KanbanBoard …"` → release commit `74b07e1`, 7 SSH-signed tags (root + 6 sub-modules), 8b tree assertions verified by hand (0 replaces in every tagged go.mod, version triple-agreement).~~
   ~~- Daemon-window audit: the mid-cut daemon commit `b0a98cd` contained only the script's own in-flight version bumps — verified per the release checklist before tagging.~~
   ~~- Pushed `master --follow-tags`; GitHub Release page created (first attempt had empty notes — refilled from the CHANGELOG section and edited in place).~~
   ~~- Proxy propagation verified: `go list -m …@v1.16.0` resolves for root, utils, htmx; `go get` in a scratch module succeeds.~~
   ~~- **Post-propagation tidy sweep executed and committed** (`367a5ff`): 5 sub-module go.sums refreshed — exactly the step whose omission left master red for 9 days after v1.11.0.~~
   ~~- pkg.go.dev still showed 1.15.1 at check time (async indexer; versioned URL 404 at first nudge) — left to the standing 24-hour watch; the module proxy itself is authoritative and green.~~

2. ~~**N2 — AXE-CORE A11Y HARNESS + FIRST SWEEP + ALL CRITICAL FIXES: COMPLETE, VERIFY-GREEN.**~~ done at `0ac2228`
   ~~- Vendored axe-core 4.11.1 (MPL-2.0, header intact) into `visualtest/testdata/axe.min.js`; embedded (no CDN, hermetic).~~
   ~~- `RunAxe` harness: inject → `axe.run` → poll-bridge for the promise → JSON decode; chromedp `Run` (never `Do`) per the invalid-context lesson.~~
   ~~- **Positive-control test** (`TestAxeHarnessDetectsViolations`): a page with an unlabeled image MUST fail, proving the guard can fail.~~
   ~~- Baseline ledger `visualtest/testdata/axe_baseline.json` (rule-level acceptance, `-1` budgets, justification comments) — node-count budgets proved flaky on live pages with polled regions.~~
   ~~- Sweep over 9 route passes (7 routes; index+forms light AND dark): **zero unaccepted critical/serious violations**.~~
   ~~- Real library defects found by the sweep and fixed (commit `fafce41` + daemon snapshots `6733356`..`7bbdef1`):~~
     ~~- **Labeled form controls now always associate with their labels** — `fieldID(id,name)` derivation in Input, Textarea, Select, DatePicker, FileInput, Checkbox, Radio, RadioGroup options, Rating stars, Slider, TagsInput, Combobox (previously `for=""`/no association = axe critical `label`).~~
     ~~- ProgressBar never renders unnamed (`aria-label` fallback chain AriaLabel→Label→"Progress").~~
     ~~- BarChart + Heatmap `role="img"` wrappers get fallback names ("Bar chart"/"Heatmap").~~
     ~~- Tabs: `aria-controls` only when a tab owns a panel (was: references to non-existent ids).~~
     ~~- Calendar: invalid bare `role="grid"` → `role="group"`.~~
     ~~- Carousel: scroll-snap track keyboard-focusable (`tabindex` + "Slides" label).~~
     ~~- EChart: `role="img"` alongside `aria-label` (bare-div aria-label was prohibited).~~
     ~~- StatCard trend: `text-green-600`→`-700`, `text-amber-600`→`-700` in light mode (3.2:1 → WCAG-pass).~~
     ~~- Demo content: labels on raw controls, focusable code blocks, underlined in-text links.~~
     ~~- `cmd/tc` embedded `_sources` re-synced for all drifted templates.~~
   ~~- Full `nix run .#verify` green after the batch; CHANGELOG `[Unreleased]` warm (Fixed/Changed/Added sections).~~

## b) PARTIALLY DONE

1. ~~**N3 — Demo click-through e2e: ~60%.** Five flows written in `visualtest/demo_flows_e2e_test.go` (LoadMore→EndOfList, ConfirmDelete with native confirm dialog, LoadingButton busy gate, multipart upload echo, kanban move click-through on both demo transports) on the new `StartDemoServer` fixture (builds the real demo binary, ephemeral port, health-wait, server-log 500 assertions).~~ done — demo flows green 23-21
   ~~First browser run surfaced two genuine chromedp traps, both fixed in the file: (1) LoadMore's `outerHTML` self-replacement invalidates cached node ids → new `demoClickUntil` re-query-retry helper; (2) a synchronous dialog handler inside `ListenTarget` deadlocked the target event loop → binary hit the 10-minute timeout → moved to an out-of-band goroutine with `Do(ctx)`.~~
   ~~**The fixes compile (`go vet` clean) but have NOT been re-run yet.** The file is deliberately uncommitted (verify-green-then-commit). This is the single next action.~~

## c) NOT STARTED

1. **Push:** 23 commits sit on local master (release + go.sum sweep + axe batch + daemon snapshots). House rule: no push without you — say the word.
2. **N4–N21:** visualtest lint triage (the 68-finding pile grew slightly with the new files), CI demo smoke job, 375px sweep, RTL sweep, overlay open-state captures, coverage margin, `layout.Minimal` head-content, route goldens, Datastar synthetics, FormLayoutInline contract, DateRange docs, ErrorPage family goldens, prerender-vs-live diff, upstream-watch dry-run, docs mini-pack, changelog policy (#133), demo niceties, nix input fold.
3. **Final verification sweep + TODO_LIST harvest** of this report's section f.

## d) TOTALLY FUCKED UP (all caught in-session)

1. ~~**`forms/radio.templ` truncation propagated by the daemon.** My whole-file rewrite of radio.templ dropped the `RadioGroup` template tail; `templ generate` ate it into `radio_templ.go`; the daemon auto-committed the broken pair (`4e65d55`), making HEAD itself uncompilable (`undefined: forms.RadioGroup`) mid-session. Recovered from tag `v1.16.0`, re-applied the fix, diff-verified against the tag before rebuilding. Recurring lesson: whole-file `write` on a partially-read file + a daemon that snapshots every ~60s = repo-wide breakage in one step; only the demo build caught it fast.~~ done (docs-health pass 2026-09-08)
2. ~~**Near-ghost-guard: the axe sweep reported PASS while every subtest silently SKIPPED.** Outside the nix dev shell there is no Chromium, `newTab` skips, and the parent test still prints PASS. I trusted an early "green" run until the positive-control test SKIPped when run alone; `-v | grep -c SKIP` confirmed 9/9 skips. All subsequent a11y runs went through `nix run .#visual`. Prevention idea (section e): fail a suite that ran zero non-skipped subtests.~~ done (docs-health pass 2026-09-08)
3. ~~**10-minute binary timeout burned by a chromedp deadlock.** The synchronous dialog `Run` inside `ListenTarget` blocked the target's event loop; one failing test + three never-run tests + 600s wasted. Root cause understood; fix applied but **unverified** (b-1). Same class as the documented "never `action.Do` on a tab context" trap; deserves a line in the e2e lessons doc.~~ done (docs-health pass 2026-09-08)
4. ~~**LSP diagnostics were garbage all session** (`undefined: KanbanColumn`, templ brace errors, a phantom `exhaustruct` in `bar_chart.go` that golangci-lint contradicts with 0 issues). Builds remain the only ground truth.~~ done (docs-health pass 2026-09-08)

## e) WHAT WE SHOULD IMPROVE

1. ~~**Skip-visibility guard:** a check that a suite with N subtests actually ran ≥1 (would have caught d-2 instantly).~~ done (docs-health pass 2026-09-08)
2. ~~**Document the two chromedp traps** (outerHTML node-id invalidation → re-query retry; no synchronous calls inside `ListenTarget`) in `docs/visual-testing.md` next to the existing lessons.~~ done (docs-health pass 2026-09-08)
3. ~~**N5 should reuse `StartDemoServer`** instead of shell-scripting the demo server in CI.~~ done (docs-health pass 2026-09-08)
4. ~~**Baseline ledger discipline:** the accepted color-contrast debt (gray-400 captions, white-on-blue-500 dark buttons at 4.46:1, amber focus ring) is convention-level and needs an owner decision — it should become a TODO_LIST row, not live only in a test-file comment.~~ done (docs-health pass 2026-09-08)
5. ~~**Daemon mitigation:** for delicate sequences (release cut, template rewrites), consider pausing the daemon; it has now twice committed broken mid-edit states.~~ done (docs-health pass 2026-09-08)
6. ~~**Release cadence:** the a11y batch changes consumer-visible markup (label association!) and StatCard shades — it deserves a tagged release soon after N3 lands, not an open-ended `[Unreleased]`.~~ done (docs-health pass 2026-09-08)

## f) Top #25 things to get done next

| # | Task | Why now |
|---|------|---------|
| ~~1~~ | ~~Re-run `TestDemo*` flows; commit N3~~ done — demo flows green | ~~Fixes are in, unverified (b-1)~~ |
| 2 | Push local master (23+ commits); watch CI + Website | Release + a11y batch invisible to origin until then |
| ~~3~~ | ~~N4: visualtest lint triage (wrapcheck/contextcheck/dupl pile)~~ done — visualtest lint zero N4 | ~~Module is the lint-red outlier~~ |
| ~~4~~ | ~~N5: CI demo smoke job reusing `StartDemoServer`~~ done — demo smoke N5 | ~~CI can't see demo breakage today~~ |
| ~~5~~ | ~~N6: 375px mobile sweep incl. kanban~~ done — mobile N6 | ~~Untested form factor #1~~ |
| ~~6~~ | ~~N7: RTL browser sweep~~ done — rtl N7 | ~~Scanner-verified only~~ |
| ~~7~~ | ~~N8: overlay open-state captures (Click + FullViewport)~~ done — overlay N8 | ~~Top-layer path never screenshot-verified~~ |
| ~~8~~ | ~~N9: coverage margin (70% floor + 2pt)~~ done — coverage N9 | ~~CI stability~~ |
| ~~9~~ | ~~N10: `layout.Minimal` head-content (NoIndex/Canonical/hreflang/JSON-LD)~~ done — Minimal SEO N10 | ~~Real consumer demand (#156 survey)~~ |
| ~~10~~ | ~~N11: route goldens for the 7 demo pages~~ done — route goldens N11 | ~~Would have caught past route regressions~~ |
| ~~11~~ | ~~N12: Datastar JS synthetics (SSEErrorHandling DOM, aria-busy clear)~~ done — synthetics N12 | ~~JS paths string-pinned only~~ |
| ~~12~~ | ~~N13: FormLayoutInline width contract~~ done — refuted N13 | ~~#166~~ |
| ~~13~~ | ~~N14: DateRange block-vs-inline docs + adjacent goldens~~ done — DateRange N14 | ~~#176~~ |
| ~~14~~ | ~~N15: ErrorPage family goldens (5 families)~~ done — errorpage N15 | ~~#177~~ |
| ~~15~~ | ~~N16: prerender vs live HTML diff~~ done — prerender diff N16 | ~~#167~~ |
| ~~16~~ | ~~N17: upstream-watch `workflow_dispatch` dry-run~~ done — upstream watch N17 | ~~#128~~ |
| ~~17~~ | ~~N18: docs mini-pack (SSE innerHTML-no-scripts, FEATURES coarse-pointer, counts)~~ done — docs minipack N18 | ~~#180 + freshness~~ |
| ~~18~~ | ~~N19: changelog policy for test-only diffs~~ **Won't implement — decided N19.** | ~~#133 (owner policy)~~ |
| 19 | N20: demo niceties (file-backed kanban, dashboard recipe section) | #189 |
| ~~20~~ | ~~N21: fold `nixpkgs-go`+`nixpkgs` inputs, verify templ zero-diff~~ done — flake fold N21 | ~~#146~~ |
| 21 | TODO_LIST: add palette-contrast owner row + skip-visibility guard idea | From e-4/e-1 |
| 22 | 24h-watch items: pkg.go.dev shows 1.16.0; CI green on pushed tip; CSS byte-stability | Release follow-up |
| ~~23~~ | ~~Update FEATURES.md/README counts if any catalogue surfaces drifted~~ done — TestDocsCountDrift | ~~Drift guards warn only~~ |
| ~~24~~ | ~~Consider `visualtest` lint waivers formalization (file-level nolint matrix vs config)~~ done — visualtest waiver N4 | ~~Feeds N4~~ |
| 25 | Plan next release (v1.17.0) once N3 lands — the a11y markup changes are consumer-facing | See g |

## g) Top #1 question

**The a11y batch changes what consumers' pages render** (label/control association for every id-less labeled field; StatCard trend shades; named progress/img roles). Do you want **v1.17.0 cut right after N3 lands** (my recommendation — it makes the fixes real for consumers and the tree is verify-green), or should the batch accumulate for a larger release? (This is Q3-style release timing, now for the post-1.16.0 batch.)

---

*Point-in-time snapshot at 2026-09-09 16:19. HEAD `7bbdef1`, origin/master `33cb443`, 1 modified file (`visualtest/demo_flows_e2e_test.go`, N3 fixes awaiting re-run).*
