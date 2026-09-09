# Status Report — Kanban Follow-Through Execution (M1–M9 + Final)

**Date:** 2026-09-09 04:20
**Session scope:** Executed the pareto plan's 1%+4%+20% tiers: M1 finish, M2 harvest, M3 goldens, M4 guard+announcement, M5 fuzz+bench, M6 docs, M7 library-wide coarse-pointer guard, M9 demo polish, final verification. Plus one unplanned CI rescue (Website workflow). This report covers THIS session only.
**Branch:** `master` @ `dbc1928` · tree clean · **origin == HEAD (daemon pushed everything, including this session's work and the pin fix)**.
**Verification:** `nix run .#verify` ALL CHECKS PASSED · per-module loop green (6 modules + visualtest compile) · full visual suite green (22s) · kanban e2e 5/5 green in real Chromium · CI green at HEAD (CI + Website).

---

## a) FULLY DONE

| #  | Item                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Proof                                                                                   |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------- |
| ~~1~~  | ~~**M1 finished — coarse-pointer computed-style browser proof.** Empirically established: `Emulation.setEmulatedMedia` IGNORES the pointer/hover features in this Chromium build; `Emulation.setTouchEmulationEnabled(true)` is what flips `matchMedia('(pointer:coarse)')`; and the wrapper's 150ms transition forces POLLING the computed opacity, not a single read. Test = control leg (fine pointer → opacity "0") + flip leg (coarse → settles at "1")~~ done at `553bb88` | ~~`TestKanbanE2ECoarsePointerButtonsVisible` (daemon commit `040e47a`), full verify green~~ |
| ~~2~~  | ~~**Website CI failure root-caused + fixed (unplanned).** Origin's Website workflow was red: daemon commit `c133896` re-flipped `website/package.json` pins (TS `^7.0.2` etc.) — the THIRD documented occurrence of the #126 regression class. Restored lockfile-matched pins per the proven `eb1f0fe` pattern; proved manifest/lockfile sync with `pnpm install --frozen-lockfile` under real nix node. Website workflow **green at HEAD**~~ done at `553bb88` | ~~daemon `7816267`; gh: `success Website 143b8db` onward~~ |
| ~~3~~  | ~~**M2 — TODO_LIST harvest + ID-collision kill.** Kanban follow-through rows #181–#189; owner gates Q1/Q2/Q3 as blocked #190–#192; verified PRs #12/#13/#14 MERGED via gh and dropped their stale rows; renumbered CV survivors → #179/#180 (killed the pre-existing #157–#161 duplication across sections); annotated #159/#160/#168 with the kanban sweep extensions; next-free-ID header → 193~~ done at `553bb88` | ~~`2c7e149`~~ |
| ~~4~~  | ~~**M3 — kanban pixel goldens** light/dark/RTL via `AssertScreenshot`; later **recaptured at the correct 1056px width** after discovering the first capture ran against stale CSS (see d1)~~ done at `553bb88` | ~~`8a111ab` + `5630818`; full visual suite green~~ |
| ~~5~~  | ~~**M4 — cross-board drop guard + post-swap live-region announcement.** dragover no longer claims foreign-card drops; drop returns early before submitting. Every move announces "Moved X to Y." once the re-rendered board lands, via a morph-aware poll (empty re-rendered live region = the signal) that works under htmx outerHTML replacement AND Datastar in-place outer patches. Unit token tests + 2 new Chromium e2e (`TestKanbanE2ECrossBoardDropIgnored`, `TestKanbanE2EAnnouncesMove`); goldens regenerated; CHANGELOG Fixed (backfilled M1 touch entry) + Changed entries~~ done at `553bb88` | ~~`3855004` + daemon `4cf833a`/`76d2ee7`; 5/5 kanban e2e green~~ |
| ~~6~~  | ~~**M5 — FuzzParseKanbanMove + BenchmarkHotPaths_KanbanBoard.** 10s fuzz run: 92,454 execs, 0 failures, 71 interesting inputs; successful decodes must echo card/column verbatim. Benchmarks: readonly ~69µs/op, wired ~81µs/op (3×5 board)~~ done at `553bb88` | ~~`bc9daa9` + daemon `ff07e09`~~ |
| ~~7~~  | ~~**M6 — docs pack.** DOMAIN_LANGUAGE +6 terms (KanbanBoard, Move Contract, Cross-Board Guard, Post-Swap Announcement, Coarse-Pointer Visibility); KanbanBoard added to container-query-strategy Part 6 as a documented **Reject**; transport-wiring kanban section +3 consumer facts (multi-board safety, SR confirmation, touch behavior)~~ done at `553bb88` | ~~`bc5f98a`~~ |
| ~~8~~  | ~~**M7 — library-wide coarse-pointer audit + guard.** Audited every `group-hover:`/`peer-hover:` usage: kanban covered by hook+CSS; tooltip/hovercard = documented exemptions (ADR-0017 family); everything else is color-accent only. New `utils.TestCoarsePointerCompliance` enforces the convention (tc-* hook in custom.css's `@media (pointer: coarse)` block, or explicit exemption), **mutation-tested** (a probe violation fails it). custom.css convention comment, AGENTS.md guard line, skill guard-table row, README/ROADMAP golden counts 114→117~~ done at `553bb88` | ~~`87248d2` + skill row this turn~~ |
| ~~9~~  | ~~**M9 — demo polish.** Dead `heroWireLine` removed; demo lint 0 issues (the plan's "writestring warnings" no longer exist — verified by explicitly linting `./examples/demo/...`); all-route shots captured incl. kanban (demo served on :8901); demo tests green. Transport comparison judged satisfied by the existing side-by-side boards (call flagged in e6)~~ done at `553bb88` | ~~daemon `143b8db`; `/tmp/tc-shots/*`~~ |
| ~~10~~ | ~~**Final verification.** `nix run .#verify` all checks passed; per-module loop green; full visual suite green; stray `examples/demo/demo` binary trashed + path gitignored (daemon staging hazard); **CI green at HEAD on origin** (the transient red at `143b8db` was the stale-golden Visual job — see d1 — resolved by the recapture)~~ done at `553bb88` | ~~`5630818`, `dbc1928`; gh: `success CI dbc1928`~~ |

## b) PARTIALLY DONE

| Item                                   | State                                                                                                         | Remaining                                                                                                                                         |
| -------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Plan execution**                     | 9 of 26 medium tasks done this session (M1, M2, M3, M4, M5, M6, M7, M9, Final); fine tasks f1–f7, f9 complete | M8, M10–M26 untouched (see c)                                                                                                                     |
| **Docs for the new guard**             | AGENTS.md + custom.css + (just now) skill guard table done                                                    | FEATURES.md has no line about the coarse-pointer convention (candidate for the next docs pass; not a drift-guard failure)                         |
| **CHANGELOG coverage of this session** | M1/M4/M7 entries in `[Unreleased]`                                                                            | M3/M5/M9 were test/demo-only diffs — no entry written because test-only changelog policy is literally undecided (TODO #133). Gated, not forgotten |

## c) NOT STARTED (plan IDs; none touched this session)

M8 visualtest 68-finding lint triage (#186) · M10 demo click-through e2e incl. kanban move (#168) · M11 axe-core harness + first sweep (#175) · M12 375px sweep incl. kanban (#159) · M13 RTL browser sweep incl. kanban (#160) · M14 overlay open-state captures (#158) · **M15 release v1.16.0 — gated by Q3/#192** · M16 CI demo smoke (#173) · M17 route goldens (#163) · M18 coverage headroom (#152) · M19 chromedp synthetics (#147) · M20 FormLayoutInline width contract (#166) · M21 DateRange docs+golden (#176) · M22 ErrorPage family goldens (#177) · M23 prerender-vs-live diff (#167) · M24 upstream-watch dry-run (#128) · M25 demo niceties (#189) · M26 Calendar Wire (#157 deferred). Owner-gated: Q1 touch-drag (#190), Q2 kanban props (#191), Q3 release timing (#192). External: BuildFlow #93/#107/#108/#124/#125/#126, branch protection #123, listings #28/#29.

## d) TOTALLY FUCKED UP! (all caught + fixed in-session)

1. ~~**M3 goldens were captured against STALE compiled CSS.** My new visual test introduced `w-[64rem]` in a `.go` file; `demo.css` scans `.go` files BY DESIGN, but the committed `app.css` predated the file — so the class wasn't compiled, the wrapper collapsed to intrinsic 928px, and the goldens enshrined a wrong-width board. It surfaced only at final-suite time as "0 pixels (100%) differ", AND it temporarily reddened origin's CI Visual job (`143b8db`) before my recapture (`5630818`) fixed it. Root cause: I violated the "recompile CSS after class-bearing source changes" rule — and its blind spot: the rule must cover NEW SOURCE FILES, not just edits. `TestCSSFreshness` would have caught it in CI; I got the signal late and from the wrong direction.~~ done (docs-health pass 2026-09-08)
2. ~~**The M4 announcement mechanism was designed wrong twice.** (a) The plan's own idea (htmx `afterSwap` listener) is unsound — outerHTML swaps invalidate the event's target reference (I avoided this one by reasoning, not luck). (b) My first implementation polled for NODE REPLACEMENT — but Datastar outer patches MORPH in place, so `old.isConnected` stays true and the kb-ds e2e timed out at 30s. The morph fact was ALREADY DOCUMENTED (AGENTS + `docs/datastar-runtime-facts.md`) — I designed from memory instead of reading the facts doc first. Final design (empty re-rendered live region as the signal) is transport-correct and e2e-proven.~~ done (docs-health pass 2026-09-08)
3. ~~**Cross-board e2e asserted inverted DOM semantics.** `dispatchEvent` returns `true` when the event was NOT canceled — my labels said the opposite, so the first run "failed" against a WORKING guard. An MDN check before writing would have prevented the round trip.~~ done (docs-health pass 2026-09-08)
4. ~~**Multiedit damaged adjacent code twice, and its reports lied three times.** (a) The M4 e2e edit clobbered the tail of the coarse-pointer test (undefined `visible`; build caught it). (b) The benchmark edit DELETED the Heatmap benchmark body instead of inserting after it (caught by re-reading — the previous session's d1 lesson, repeated). (c) Three separate "Applied 2 of 3 (1 edit failed)" reports were FALSE — the edits had applied; I burned cycles verifying both the failure claim AND the success claim by grep every time.~~ done (docs-health pass 2026-09-08)
5. ~~**Daemon raced ~13 of this session's commits** into generic `chore: auto-commit` messages despite committing at every green gate. Content verified-green every time (I re-checked each daemon commit's diff); only message quality suffered. Known #93/#126 — no in-repo fix exists.~~ done (docs-health pass 2026-09-08)

## e) WHAT WE SHOULD IMPROVE!

1. ~~**Read the runtime-facts doc BEFORE designing transport JS.** `docs/datastar-runtime-facts.md` is exactly the "verify against the counterparty artifact" principle the skill preaches; the morph fact cost me one 30s-timeout debugging cycle + a redesign. Rule: any JS that must survive a swap → check the facts doc + the pinned-bundle contract FIRST.~~ done (docs-health pass 2026-09-08)
2. ~~**Golden-capture runbook: recompile CSS first.** `nix run .#css` must precede ANY `AssertScreenshot -update` whenever the test (or any source) introduces classes — including entirely new files. The freshness guard only fires in CI; local captures are on the honor system.~~ done (docs-health pass 2026-09-08)
3. ~~**Post-multiedit ritual: grep the change AND the neighbors.** This session's tool reports were unreliable in both directions. Cheap ritual: one grep for the intended token + one symbol listing around the edit site before moving on.~~ done (docs-health pass 2026-09-08)
4. ~~**Check DOM/MDN contract semantics before writing event assertions** (dispatchEvent return, event.bubbles, focus behavior). Two of five d-items were "wrote the test wrong against working code".~~ done (docs-health pass 2026-09-08)
5. ~~**Ship the skill's guard-table row in the SAME commit as the guard.** I updated AGENTS.md but forgot the skill table until this report's self-review — the guard's discoverability is half its value. (Fixed now, unprompted gap.)~~ done (docs-health pass 2026-09-08)
6. ~~**Decide #133 (changelog policy for test-only diffs)** so M3/M5-class work has a deterministic changelog home instead of a judgment call every time.~~ **Won't implement — decided N19.**
7. **The daemon regression class needs its root fix** (#126/#93 in `larsartmann/buildflow`): third TS-pin occurrence this session (`c133896`) reddened origin's Website workflow between releases. The pre-commit minified-CSS guard caught nothing here because package.json isn't gated — a `pnpm --frozen-lockfile` gate (or vetted-artifact classifier) would have.

## f) Up to 50 things we should get done next

_(ranked by impact; TODO IDs in parens; the plan file's 92-row fine table stays the exhaustive source)_

1. ~~M11: axe-core a11y scan harness via chromedp + fix first findings (#175) — Impact H, Effort L~~ done — axe N2
2. ~~M8: visualtest 68-finding lint triage — bucket, fix mechanical, exemption decision (#186) — H, M~~ done — visualtest lint zero N4
3. ~~M15: **release v1.16.0** — KanbanBoard + touch fix + guard + everything since M1; gated by Q3/#192 (#188) — H, S~~ done — v1.16.0 cut
4. ~~M10: demo click-through e2e incl. kanban move click (#168) — H, M~~ done — demo flows N3
5. ~~M16: CI demo smoke — build→serve→shots→assert + zero-500 log assert (#173) — H, L~~ done — demo smoke N5
6. ~~M12: 375px mobile sweep incl. kanban (#159) — M, M~~ done — mobile N6
7. ~~M13: RTL browser sweep incl. kanban (#160) — M, M~~ done — rtl N7
8. ~~M14: overlay open-state captures Modal/Drawer/Tooltip/Combobox/Carousel (#158) — M, M~~ done — overlay N8
9. ~~M17: page-level demo route goldens, 7 routes (#163) — M, M~~ done — route goldens N11
10. ~~M19: chromedp synthetics — datastar-fetch→SSEErrorHandling DOM; patch→aria-busy (#147) — M, M~~ done — synthetics N12
11. ~~M20: FormLayoutInline width contract docs+fix+guard (#166) — M, S~~ done — refuted N13
12. ~~M18: coverage headroom over the 70% floor (#152) — M, M~~ done — coverage N9
13. ~~M21: DateRange block-vs-inline docs + adjacent-ranges golden (#176) — L, S~~ done — DateRange N14
14. ~~M22: ErrorPage family matrix goldens, 5 families (#177) — L, S~~ done — errorpage N15
15. ~~M23: prerender-vs-live HTML diff, 7 routes (#167) — L, S~~ done — prerender diff N16
16. ~~M24: upstream-watch workflow_dispatch dry-run green check (#128) — L, S~~ done — upstream watch N17
17. M25: demo niceties — file-backed kanban state, Dashboard-recipe kanban section (#189) — L, M
18. M26: Calendar Wire adoption per the D3 rule (#157 deferred) — M, L
19. ~~#179: `layout.Minimal` head-content support (SEO parity with Base)~~ done — Minimal SEO N10
20. ~~#180: document the SSE-fragment innerHTML-no-scripts fact in datastar/SSE docs~~ done — facts N18
21. Human-eyeball the kanban board PNGs (3, agent-captured — #80 family) and the recaptured widths
22. ~~FEATURES.md: add the coarse-pointer convention line (b-item closure)~~ done — features coarse pointer N18
23. Decide #133: changelog-guard policy for test-only PRs (+ shakedown)
24. #126 root fix in BuildFlow: package.json/lockfile gate (third occurrence this session)
25. Kanban: within-column keyboard reorder button pair — gated by Q2/#191
26. #162: human-eyeball progressbar half_light.png (45% fill suspicion) — still open from earlier reports
27. ~~#133-adjacent: warm CHANGELOG for M3/M5 once policy decided~~ done — unreleased warm
28. ~~Re-verify origin CI stays green through the next daemon push cycle (pin fix is on origin; watch for occurrence #4)~~ done (docs-health pass 2026-09-08)

_(28 concrete items; not padding to 50 — the plan file's fine table carries the rest.)_

## g) Questions I can NOT figure out myself

1. **Q2 (TODO #191) — KanbanBoard feature scope:** which of these belong IN the component vs consumer composition — within-column keyboard reorder, per-column "add card" wiring, WIP limits (red count badge), card/column tone accents? Each is buildable to library standard; only you can rank the scope. (This gates the next kanban-size task.)
2. ~~**Q3 (TODO #188/#192) — Release timing:** cut **v1.16.0 now** — KanbanBoard + touch fix + cross-board guard + coarse-pointer guard are verified end-to-end, `[Unreleased]` is warm, CI green at HEAD — or accumulate the quality tier first (M8 lint triage, M10–M14 e2e/sweep coverage)? The original "accumulate M3–M6 first" branch is moot: M3–M6 shipped this session.~~ **Won't implement — answered v1.16.0.**
3. **Q1 (TODO #190) — Touch-drag story:** is button-based touch UX the accepted end state for KanbanBoard on phones (buttons always visible, per-pointer CSS), or do you want REAL touch drag (long-press / pointer-events DnD) — remembering the dependency budget is deliberately closed (templ + tailwind-merge-go + go-error-family) and a polyfill-free implementation is meaningful bespoke JS?

---

_Point-in-time snapshot. Written by the 2026-09-09 follow-through execution session (M1–M9 + Final). Format: `.md` per explicit user instruction (skill default is HTML — override noted). WAITING FOR INSTRUCTIONS._
